// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package grpc

import (
	"context"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/net/grpc/types"
	lbmetrics "github.com/vdaas/vald/internal/observability/metrics/gateway/lb"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// linearClient serves LinearSearch with one result so that every query asking
// for more (MinNum) fails with its own "not enough" NotFound error.
type linearClient struct {
	*regClient
}

func (*linearClient) LinearSearch(
	context.Context, *payload.Search_Request, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return probeResponse(), nil
}

// LinearSearchByID backs the LinearSearchByID fallback with the same single result.
func (*linearClient) LinearSearchByID(
	context.Context, *payload.Search_IDRequest, ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return probeResponse(), nil
}

// Test_server_MultiLinearSearch_joinedErrors pins the error MultiLinearSearch and
// MultiLinearSearchByID return now that per-query failures are collected in a
// slice and joined once after the workers finish: N distinct failures give one
// flat joined error with N members and the members' status code.
func Test_server_MultiLinearSearch_joinedErrors(t *testing.T) {
	t.Parallel()
	const n = 8
	// The i-th query asks for MinNum 2+i, so each failure message differs.
	cfgs := make([]*payload.Search_Config, 0, n)
	for minNum := uint32(2); len(cfgs) < n; minNum++ {
		cfgs = append(cfgs, &payload.Search_Config{Num: 10, MinNum: minNum, Timeout: int64(time.Minute)})
	}
	tests := []struct {
		call func(context.Context, *server) error
		name string
	}{
		{
			name: "MultiLinearSearch",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Search_MultiRequest{Requests: make([]*payload.Search_Request, 0, n)}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Search_Request{Vector: []float32{1, 2, 3}, Config: cfgs[i]})
				}
				_, err := s.MultiLinearSearch(ctx, reqs)
				return err
			},
		},
		{
			name: "MultiLinearSearchByID",
			call: func(ctx context.Context, s *server) error {
				reqs := &payload.Search_MultiIDRequest{Requests: make([]*payload.Search_IDRequest, 0, n)}
				for i := range n {
					reqs.Requests = append(reqs.Requests, &payload.Search_IDRequest{Id: "id-" + strconv.Itoa(i), Config: cfgs[i]})
				}
				_, err := s.MultiLinearSearchByID(ctx, reqs)
				return err
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var c vald.Client = &linearClient{regClient: &regClient{
				getObject: func(_ context.Context, in *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
					return &payload.Object_Vector{Id: in.GetId().GetId(), Vector: []float32{1, 2, 3}}, nil
				},
			}}
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: syncGateway(c), replica: 1, multiConcurrency: n, eg: eg}
			err := test.call(t.Context(), s)
			if err == nil {
				t.Fatal("returned nil, want the per-query failures")
			}
			if st, ok := status.FromError(err); !ok || st == nil || st.Code() != codes.NotFound {
				t.Errorf("status code = %v (ok=%t), want NotFound: %v", st.Code(), ok, err)
			}
			var u interface{ Unwrap() []error }
			if !errors.As(err, &u) {
				t.Fatalf("got %T, want a joined error: %v", err, err)
			}
			if got := len(u.Unwrap()); got != n {
				t.Errorf("members = %d, want %d (flat): %v", got, n, err)
			}
			for _, m := range u.Unwrap() {
				var nested interface{ Unwrap() []error }
				if errors.As(m, &nested) {
					t.Errorf("member %T is itself a joined error, want flat: %v", m, m)
				}
			}
			for i := range n {
				if want := "required: " + strconv.Itoa(2+i) + ","; !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not contain %q", err.Error(), want)
				}
			}
		})
	}
}

type funcLinearClient struct {
	*regClient
	fn func(context.Context, *payload.Search_Request) (*payload.Search_Response, error)
}

func (c *funcLinearClient) LinearSearch(
	ctx context.Context, in *payload.Search_Request, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	return c.fn(ctx, in)
}

func Test_server_LinearSearch_subsetReplicaCover(t *testing.T) {
	t.Parallel()
	addrs := []string{"a0", "a1", "a2", "a3", "a4"}
	vec := []float32{1, 2, 3}

	t.Run("healthy 3/5 subset rotates and records 0 unreached metrics", func(t *testing.T) {
		t.Parallel()
		clients := make(map[string]vald.Client, len(addrs))
		for _, addr := range addrs {
			clients[addr] = &funcLinearClient{
				regClient: new(regClient),
				fn: func(context.Context, *payload.Search_Request) (*payload.Search_Response, error) {
					return &payload.Search_Response{
						Results: []*payload.Object_Distance{{Id: "hit-" + addr, Distance: 0.1}},
					}, nil
				},
			}
		}
		gw := newSubsetRegGateway(addrs, clients)
		sm := lbmetrics.New()
		s := &server{gateway: gw, replica: 3, timeout: time.Second, searchMetrics: sm}

		// Asking for Num=10 with MinNum=0 and getting 3 (< 10) from 3/5 healthy agents must NOT fall back.
		res1, err := s.LinearSearch(t.Context(), &payload.Search_Request{
			Vector: vec,
			Config: &payload.Search_Config{Num: 10},
		})
		if err != nil || len(res1.GetResults()) != 3 {
			t.Fatalf("first LinearSearch = (%d, %v), want 3 results", len(res1.GetResults()), err)
		}
		res2, err := s.LinearSearch(t.Context(), &payload.Search_Request{
			Vector: vec,
			Config: &payload.Search_Config{Num: 10},
		})
		if err != nil || len(res2.GetResults()) != 3 {
			t.Fatalf("second LinearSearch = (%d, %v), want 3 results", len(res2.GetResults()), err)
		}

		gw.mu.Lock()
		if len(gw.waves) != 2 || gw.bcCalls != 0 {
			gw.mu.Unlock()
			t.Fatalf("waves = %+v, bcCalls = %d, want 2 waves and 0 full broadcasts", gw.waves, gw.bcCalls)
		}
		if !slices.Equal(gw.waves[0].targets, []string{"a0", "a1", "a2"}) ||
			!slices.Equal(gw.waves[1].targets, []string{"a1", "a2", "a3"}) {
			gw.mu.Unlock()
			t.Fatalf("rotated waves = %+v, want [a0 a1 a2] then [a1 a2 a3]", gw.waves)
		}
		gw.mu.Unlock()

		if metrics := collectSearchMetrics(t, sm); len(metrics) != 0 {
			t.Errorf("expected 0 failure/partial metrics on healthy subset, got %v", metrics)
		}
	})

	t.Run("Ratio != nil uses full broadcast", func(t *testing.T) {
		t.Parallel()
		clients := make(map[string]vald.Client, len(addrs))
		for _, addr := range addrs {
			clients[addr] = &linearClient{regClient: new(regClient)}
		}
		gw := newSubsetRegGateway(addrs, clients)
		s := &server{gateway: gw, replica: 3, timeout: time.Second}
		res, err := s.LinearSearch(t.Context(), &payload.Search_Request{
			Vector: vec,
			Config: &payload.Search_Config{Num: 5, Ratio: types.Float(0.5)},
		})
		if err != nil || len(res.GetResults()) == 0 {
			t.Fatalf("LinearSearch = (%v, %v)", res, err)
		}
		gw.mu.Lock()
		defer gw.mu.Unlock()
		if len(gw.waves) != 0 || gw.bcCalls != 1 {
			t.Fatalf("waves = %d, bcCalls = %d, want 0 waves and 1 full broadcast", len(gw.waves), gw.bcCalls)
		}
	})
}

func Test_server_LinearSearch_subsetFallback(t *testing.T) {
	t.Parallel()
	addrs := []string{"a0", "a1", "a2", "a3", "a4"}
	vec := []float32{1, 2, 3}
	newClients := func(fn func(idx int, addr string) (*payload.Search_Response, error)) map[string]vald.Client {
		clients := make(map[string]vald.Client, len(addrs))
		for i, addr := range addrs {
			idx := i
			target := addr
			clients[addr] = &funcLinearClient{
				regClient: new(regClient),
				fn: func(context.Context, *payload.Search_Request) (*payload.Search_Response, error) {
					return fn(idx, target)
				},
			}
		}
		return clients
	}

	for _, tc := range []struct {
		setup func() (*subsetRegGateway, *payload.Search_Config)
		name  string
	}{
		{
			name: "Wave 1 error triggers Wave 2 fallback",
			setup: func() (*subsetRegGateway, *payload.Search_Config) {
				clients := newClients(func(idx int, addr string) (*payload.Search_Response, error) {
					if idx == 0 {
						return nil, status.Error(codes.Internal, "boom")
					}
					return &payload.Search_Response{
						Results: []*payload.Object_Distance{{Id: "hit-" + addr, Distance: 0.1}},
					}, nil
				})
				return newSubsetRegGateway(addrs, clients), &payload.Search_Config{Num: 1}
			},
		},
		{
			name: "unreached Wave 1 agent triggers Wave 2 fallback",
			setup: func() (*subsetRegGateway, *payload.Search_Config) {
				clients := newClients(func(_ int, addr string) (*payload.Search_Response, error) {
					return &payload.Search_Response{
						Results: []*payload.Object_Distance{{Id: "hit-" + addr, Distance: 0.1}},
					}, nil
				})
				gw := newSubsetRegGateway(addrs, clients)
				gw.skipped = map[string]bool{"a1": true}
				return gw, &payload.Search_Config{Num: 1}
			},
		},
		{
			name: "per-target DeadlineExceeded triggers Wave 2 fallback",
			setup: func() (*subsetRegGateway, *payload.Search_Config) {
				clients := newClients(func(idx int, addr string) (*payload.Search_Response, error) {
					if idx == 0 {
						return nil, status.Error(codes.DeadlineExceeded, "agent timeout")
					}
					return &payload.Search_Response{
						Results: []*payload.Object_Distance{{Id: "hit-" + addr, Distance: 0.1}},
					}, nil
				})
				return newSubsetRegGateway(addrs, clients), &payload.Search_Config{Num: 1}
			},
		},
		{
			name: "0 results in Wave 1 triggers Wave 2 fallback",
			setup: func() (*subsetRegGateway, *payload.Search_Config) {
				clients := newClients(func(idx int, addr string) (*payload.Search_Response, error) {
					if idx < 3 {
						return &payload.Search_Response{}, nil
					}
					return &payload.Search_Response{
						Results: []*payload.Object_Distance{{Id: "hit-" + addr, Distance: 0.1}},
					}, nil
				})
				return newSubsetRegGateway(addrs, clients), &payload.Search_Config{Num: 1}
			},
		},
		{
			name: "fewer unique results than MinNum triggers Wave 2 fallback",
			setup: func() (*subsetRegGateway, *payload.Search_Config) {
				clients := newClients(func(idx int, addr string) (*payload.Search_Response, error) {
					if idx < 3 {
						return &payload.Search_Response{
							Results: []*payload.Object_Distance{{Id: "dup", Distance: 0.1}},
						}, nil
					}
					return &payload.Search_Response{
						Results: []*payload.Object_Distance{{Id: "unique-" + addr, Distance: 0.2}},
					}, nil
				})
				return newSubsetRegGateway(addrs, clients), &payload.Search_Config{Num: 5, MinNum: 2}
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gw, cfg := tc.setup()
			s := &server{gateway: gw, replica: 3, timeout: time.Second}
			res, err := s.LinearSearch(t.Context(), &payload.Search_Request{
				Vector: vec,
				Config: cfg,
			})
			if err != nil {
				t.Fatalf("LinearSearch error = %v", err)
			}
			if len(res.GetResults()) == 0 {
				t.Fatal("expected results after Wave 2 fallback")
			}
			gw.mu.Lock()
			defer gw.mu.Unlock()
			if len(gw.waves) != 2 || gw.waves[0].count != 3 || gw.waves[1].count != 2 {
				t.Fatalf("waves = %+v, want [3, 2]", gw.waves)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_server_LinearSearch(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_Request
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotRes, err := s.LinearSearch(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_LinearSearchByID(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Search_IDRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Response
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Response, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Response, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           req:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotRes, err := s.LinearSearchByID(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamLinearSearch(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamLinearSearchServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			err := s.StreamLinearSearch(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamLinearSearchByID(t *testing.T) {
// 	type args struct {
// 		stream vald.Search_StreamLinearSearchByIDServer
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           stream:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           stream:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			err := s.StreamLinearSearchByID(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiLinearSearch(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotRes, err := s.MultiLinearSearch(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiLinearSearchByID(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Search_MultiIDRequest
// 	}
// 	type fields struct {
// 		UnimplementedValdServer vald.UnimplementedValdServer
// 		eg                      errgroup.Group
// 		gateway                 service.Gateway
// 		name                    string
// 		ip                      string
// 		timeout                 time.Duration
// 		replica                 int
// 		streamConcurrency       int
// 		multiConcurrency        int
// 	}
// 	type want struct {
// 		wantRes *payload.Search_Responses
// 		err     error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Search_Responses, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotRes *payload.Search_Responses, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotRes, w.wantRes) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotRes, w.wantRes)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		       },
// 		       fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           reqs:nil,
// 		           },
// 		           fields: fields {
// 		           UnimplementedValdServer:nil,
// 		           eg:nil,
// 		           gateway:nil,
// 		           name:"",
// 		           ip:"",
// 		           timeout:nil,
// 		           replica:0,
// 		           streamConcurrency:0,
// 		           multiConcurrency:0,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
// 			s := &server{
// 				UnimplementedValdServer: test.fields.UnimplementedValdServer,
// 				eg:                      test.fields.eg,
// 				gateway:                 test.fields.gateway,
// 				name:                    test.fields.name,
// 				ip:                      test.fields.ip,
// 				timeout:                 test.fields.timeout,
// 				replica:                 test.fields.replica,
// 				streamConcurrency:       test.fields.streamConcurrency,
// 				multiConcurrency:        test.fields.multiConcurrency,
// 			}
//
// 			gotRes, err := s.MultiLinearSearchByID(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
