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
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Remove used to write the exist-check flag and the timestamp into the caller's request.
func Test_server_Remove_doesNotMutateRequest(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	tests := []struct {
		cfg    *payload.Remove_Config
		name   string
		wantTS int64 // 0: any non-zero timestamp
	}{
		{name: "nil config (exist check runs)"},
		{name: "skip check, no timestamp", cfg: &payload.Remove_Config{SkipStrictExistCheck: true}},
		{name: "exist check requested", cfg: &payload.Remove_Config{}},
		{name: "timestamp already set", cfg: &payload.Remove_Config{SkipStrictExistCheck: true, Timestamp: 7}, wantTS: 7},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c := newOKClient(nil)
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: syncGateway(c), replica: 1, eg: eg}

			id := &payload.Object_ID{Id: uuid}
			req := &payload.Remove_Request{Id: id, Config: test.cfg}
			var before payload.Remove_Config
			if test.cfg != nil {
				before.SkipStrictExistCheck, before.Timestamp = test.cfg.SkipStrictExistCheck, test.cfg.Timestamp
			}
			if _, err := s.Remove(t.Context(), req); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if req.Config != test.cfg {
				t.Error("Remove replaced the caller's req.Config")
			}
			if test.cfg != nil && (test.cfg.SkipStrictExistCheck != before.SkipStrictExistCheck || test.cfg.Timestamp != before.Timestamp) {
				t.Errorf("Remove modified the caller's config: got %v", test.cfg)
			}
			got := c.removeIn
			if got == nil {
				t.Fatal("no request reached the agent")
			}
			if got.GetId() != id {
				t.Error("the forwarded request does not carry the caller's id")
			}
			if !got.GetConfig().GetSkipStrictExistCheck() {
				t.Error("the forwarded request must skip the agent's exist check")
			}
			if ts := got.GetConfig().GetTimestamp(); ts == 0 || (test.wantTS != 0 && ts != test.wantTS) {
				t.Errorf("forwarded timestamp = %d, want %d (0: any non-zero)", ts, test.wantTS)
			}
		})
	}
}

// A NotFound answer from an agent is tolerated by Remove and ParseError never attaches details to
// it, so the (costly) error details must not be built for it: Addrs is queried only while
// building them. A real agent failure must still carry them.
func Test_server_Remove_buildsErrorDetailsOnlyForReturnedErrors(t *testing.T) {
	t.Parallel()
	const agents = 4
	tests := []struct {
		agentErr      error
		name          string
		wantAddrsCall bool
		wantErr       bool
	}{
		{name: "NotFound agents are tolerated without details", agentErr: status.Error(codes.NotFound, "not found")},
		{name: "Internal agent error carries details", agentErr: status.Error(codes.Internal, "boom"), wantAddrsCall: true, wantErr: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ok := newOKClient(nil)
			gw := &regGateway{
				agents: agents,
				broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					errs := make([]error, 0, agents)
					for i := range agents {
						if i == agents-1 {
							// the last agent holds the object
							errs = append(errs, f(ctx, "127.0.0.1:8081", ok))
							continue
						}
						errs = append(errs, f(ctx, "127.0.0.1:8081", &regClient{
							remove: func(context.Context, *payload.Remove_Request) (*payload.Object_Location, error) {
								return nil, test.agentErr
							},
						}))
					}
					return errors.Join(errs...)
				},
			}
			s := &server{gateway: gw, replica: agents}

			_, err := s.Remove(t.Context(), &payload.Remove_Request{
				Id:     &payload.Object_ID{Id: fixtureUUID},
				Config: &payload.Remove_Config{SkipStrictExistCheck: true},
			})
			if (err != nil) != test.wantErr {
				t.Fatalf("got error: %v, want error: %t", err, test.wantErr)
			}
			if got := gw.addrsCalls.Load(); (got != 0) != test.wantAddrsCall {
				t.Errorf("Addrs called %d times, want called: %t", got, test.wantAddrsCall)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_server_Remove(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Remove_Request
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
// 		wantLocs *payload.Object_Location
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Location, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 			gotLocs, err := s.Remove(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamRemove(t *testing.T) {
// 	type args struct {
// 		stream vald.Remove_StreamRemoveServer
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
// 			err := s.StreamRemove(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiRemove(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Remove_MultiRequest
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
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 			gotLocs, err := s.MultiRemove(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_RemoveByTimestamp(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Remove_TimestampRequest
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
// 		wantLocs *payload.Object_Locations
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		fields     fields
// 		want       want
// 		checkFunc  func(want, *payload.Object_Locations, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotLocs *payload.Object_Locations, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotLocs, w.wantLocs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotLocs, w.wantLocs)
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
// 			gotLocs, err := s.RemoveByTimestamp(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
