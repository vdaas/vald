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
	"math"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Update used to fail as soon as one agent failed to remove the old copy, before
// inserting the new one anywhere.
func Test_server_Update_partialAgentFailure(t *testing.T) {
	t.Parallel()
	internal := status.Error(codes.Internal, "agent internal")
	vec := []float32{1, 2, 3}
	tests := []struct {
		ctx         func(t *testing.T) context.Context
		agents      func() []*pwAgent
		name        string
		wantName    string
		wantInserts int64
		wantCode    codes.Code // codes.OK: success expected
		connMissing bool
		// noInsertCall makes DoMulti skip every agent without calling f, as it
		// does for agents whose circuit breaker is open.
		noInsertCall bool
	}{
		{
			name: "remove fails on one agent, insert succeeds",
			agents: func() []*pwAgent {
				return []*pwAgent{newPWAgent("a", internal, nil, nil), newPWAgent("b", nil, nil, nil)}
			},
			wantName:    "a,b",
			wantInserts: 2,
		},
		{
			name: "remove fails on every agent, insert succeeds",
			agents: func() []*pwAgent {
				return []*pwAgent{newPWAgent("a", internal, nil, nil), newPWAgent("b", internal, nil, nil)}
			},
			wantName:    "a,b",
			wantInserts: 2,
		},
		{
			name: "remove fails on one agent, insert succeeds on one",
			agents: func() []*pwAgent {
				return []*pwAgent{newPWAgent("a", internal, internal, nil), newPWAgent("b", nil, nil, nil)}
			},
			wantName:    "b",
			wantInserts: 2,
		},
		{
			name: "remove fails on one agent, insert fails everywhere",
			agents: func() []*pwAgent {
				return []*pwAgent{newPWAgent("a", internal, internal, nil), newPWAgent("b", nil, internal, nil)}
			},
			wantCode:    codes.Internal,
			wantInserts: 2,
		},
		{
			name: "remove failure with a canceled request is fatal",
			ctx: func(t *testing.T) context.Context {
				t.Helper()
				ctx, cancel := context.WithCancel(t.Context())
				cancel()
				return ctx
			},
			agents: func() []*pwAgent {
				return []*pwAgent{newPWAgent("a", internal, nil, nil), newPWAgent("b", nil, nil, nil)}
			},
			wantCode: codes.Internal,
		},
		{
			name: "remove succeeds, insert reaches no agent",
			agents: func() []*pwAgent {
				return []*pwAgent{newPWAgent("a", nil, nil, nil), newPWAgent("b", nil, nil, nil)}
			},
			noInsertCall: true,
			wantCode:     codes.Unavailable,
		},
		{
			name:        "missing connection is fatal",
			connMissing: true,
			agents: func() []*pwAgent {
				return []*pwAgent{newPWAgent("a", nil, nil, nil)}
			},
			wantCode: codes.Internal,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			agents := test.agents()
			gw := pwGateway(agents...)
			if test.connMissing {
				gw.broadCast = func(context.Context, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					return errors.ErrGRPCClientConnNotFound("*")
				}
			}
			if test.noInsertCall {
				gw.doMulti = func(context.Context, int, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					return nil
				}
			}
			ctx := t.Context()
			if test.ctx != nil {
				ctx = test.ctx(t)
			}
			s := &server{gateway: gw, replica: len(agents)}
			loc, err := s.Update(ctx, &payload.Update_Request{
				Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec},
				Config: &payload.Update_Config{SkipStrictExistCheck: true, Timestamp: 10},
			})
			var inserts int64
			for _, a := range agents {
				inserts += a.inserts.Load()
			}
			if inserts != test.wantInserts {
				t.Errorf("agents got %d inserts, want %d", inserts, test.wantInserts)
			}
			if test.wantCode != codes.OK {
				if err == nil {
					t.Fatalf("Update returned %v, want an error", loc)
				}
				if got := pwCode(t, err); got != test.wantCode {
					t.Errorf("status code = %v, want %v: %v", got, test.wantCode, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Update returned %v, want success", err)
			}
			if loc.GetName() != test.wantName {
				t.Errorf("Name = %q, want %q", loc.GetName(), test.wantName)
			}
		})
	}
}

// Update with DisableBalancedUpdate and UpdateTimestamp used to fail when one agent
// failed although another applied the change.
func Test_server_broadCastWrite_partialAgentFailure(t *testing.T) {
	t.Parallel()
	internal := status.Error(codes.Internal, "agent internal")
	unavailable := status.Error(codes.Unavailable, "agent unavailable")
	calls := []struct {
		call func(ctx context.Context, s *server) (*payload.Object_Location, error)
		name string
	}{
		{
			name: "Update with DisableBalancedUpdate",
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.Update(ctx, &payload.Update_Request{
					Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: []float32{1, 2, 3}},
					Config: &payload.Update_Config{DisableBalancedUpdate: true, Timestamp: 10},
				})
			},
		},
		{
			name: vald.UpdateTimestampRPCName,
			call: func(ctx context.Context, s *server) (*payload.Object_Location, error) {
				return s.UpdateTimestamp(ctx, &payload.Update_TimestampRequest{Id: fixtureUUID, Timestamp: 10})
			},
		},
	}
	tests := []struct {
		name     string
		wantName string
		errs     []error    // per agent, in BroadCast order
		wantCode codes.Code // codes.OK: success expected
	}{
		{name: "one agent fails", errs: []error{internal, nil}, wantName: "b"},
		{name: "every agent fails", errs: []error{unavailable, internal}, wantCode: codes.Unavailable},
	}
	for _, c := range calls {
		for _, test := range tests {
			t.Run(c.name+"/"+test.name, func(t *testing.T) {
				t.Parallel()
				agents := make([]*pwAgent, 0, len(test.errs))
				for i, err := range test.errs {
					agents = append(agents, newPWAgent(string(rune('a'+i)), nil, nil, err))
				}
				// replica 1: no top-up insert is needed after one agent applied the change.
				s := &server{gateway: pwGateway(agents...), replica: 1}
				loc, err := c.call(t.Context(), s)
				if test.wantCode != codes.OK {
					if err == nil {
						t.Fatalf("returned %v, want an error", loc)
					}
					if got := pwCode(t, err); got != test.wantCode {
						t.Errorf("status code = %v, want %v (first target's error): %v", got, test.wantCode, err)
					}
					return
				}
				if err != nil {
					t.Fatalf("returned %v, want success", err)
				}
				if loc.GetName() != test.wantName {
					t.Errorf("Name = %q, want %q", loc.GetName(), test.wantName)
				}
			})
		}
	}
}

// The balanced Update removes at `now` and re-inserts at now+1. The client's vector
// timestamp used to be forwarded with the re-insert; an agent then took it over the
// config timestamp and dropped the insert behind its own delete.
func Test_server_Update_reinsertNotDroppedBehindDelete(t *testing.T) {
	t.Parallel()
	vec := []float32{1, 2, 3}
	tests := []struct {
		call     func(s *server) (*payload.Object_Location, error)
		name     string
		configTS int64 // 0: the gateway picks time.Now()
	}{
		{
			name:     "Update with config timestamp, older vector timestamp",
			configTS: 100,
			call: func(s *server) (*payload.Object_Location, error) {
				return s.Update(t.Context(), &payload.Update_Request{
					Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec, Timestamp: 5},
					Config: &payload.Update_Config{SkipStrictExistCheck: true, Timestamp: 100},
				})
			},
		},
		{
			name: "Update without config timestamp, vector timestamp set",
			call: func(s *server) (*payload.Object_Location, error) {
				return s.Update(t.Context(), &payload.Update_Request{
					Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec, Timestamp: 5},
					Config: &payload.Update_Config{SkipStrictExistCheck: true},
				})
			},
		},
		{
			name:     "Update with exist check, older vector timestamp",
			configTS: 100,
			call: func(s *server) (*payload.Object_Location, error) {
				return s.Update(t.Context(), &payload.Update_Request{
					Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec, Timestamp: 5},
					Config: &payload.Update_Config{Timestamp: 100},
				})
			},
		},
		{
			name:     "Upsert of an existing object, older vector timestamp",
			configTS: 100,
			call: func(s *server) (*payload.Object_Location, error) {
				return s.Upsert(t.Context(), &payload.Upsert_Request{
					Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec, Timestamp: 5},
					Config: &payload.Upsert_Config{SkipStrictExistCheck: true, Timestamp: 100},
				})
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// The stored vector differs from the requested one, so Update removes and re-inserts.
			a := newTSAgent("agent-0", loopbackIP, []float32{4, 5, 6}, 1)
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: syncGateway(a), replica: 1, eg: eg}

			if _, err := test.call(s); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			stored, ts, delTS, inserts := a.state()
			if inserts != 1 {
				t.Fatalf("agent got %d inserts, want 1", inserts)
			}
			if !stored {
				t.Fatalf("re-insert was dropped behind the delete at %d (vector timestamp forwarded)", delTS)
			}
			if ts != delTS+1 {
				t.Errorf("stored timestamp = %d, want delete timestamp + 1 = %d", ts, delTS+1)
			}
			if test.configTS != 0 && delTS != test.configTS {
				t.Errorf("delete timestamp = %d, want %d", delTS, test.configTS)
			}
		})
	}
}

// Update removes at ts and inserts at ts+1; ts == MaxInt64 used to wrap the
// insert timestamp to MinInt64, so the agent dropped the insert and Update
// still reported success. It must be rejected before anything is sent.
func Test_server_Update_timestamp(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	tests := []struct {
		name       string
		ts         int64
		wantCode   codes.Code
		wantRemove int64
		wantInsert int64
	}{
		{name: "max int64 is rejected", ts: math.MaxInt64, wantCode: codes.InvalidArgument},
		{name: "normal timestamp", ts: 7, wantCode: codes.OK, wantRemove: 7, wantInsert: 8},
		{name: "max int64 - 1 still fits", ts: math.MaxInt64 - 1, wantCode: codes.OK, wantRemove: math.MaxInt64 - 1, wantInsert: math.MaxInt64},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// The stored vector differs from the requested one, so Update proceeds to remove+insert.
			c := newOKClient([]float32{4, 5, 6})
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: syncGateway(c), replica: 1, eg: eg}

			_, err := s.Update(t.Context(), &payload.Update_Request{
				Vector: &payload.Object_Vector{Id: uuid, Vector: []float32{1, 2, 3}},
				Config: &payload.Update_Config{SkipStrictExistCheck: true, Timestamp: test.ts},
			})
			if test.wantCode != codes.OK {
				st, ok := status.FromError(err)
				if err == nil || !ok || st == nil || st.Code() != test.wantCode {
					t.Fatalf("got error %v, want code %v", err, test.wantCode)
				}
				if c.removeIn != nil || c.insertIn != nil {
					t.Errorf("a rejected Update reached the agent: remove=%v insert=%v", c.removeIn, c.insertIn)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := c.removeIn.GetConfig().GetTimestamp(); got != test.wantRemove {
				t.Errorf("remove timestamp = %d, want %d", got, test.wantRemove)
			}
			if got := c.insertIn.GetConfig().GetTimestamp(); got != test.wantInsert {
				t.Errorf("insert timestamp = %d, want %d", got, test.wantInsert)
			}
		})
	}
}

// The replica top-up failure of Update (disable balanced update) used to be discarded.
func Test_server_Update_shortageDoMultiError(t *testing.T) {
	t.Parallel()
	errDoMulti := errors.New("doMulti failed")
	gw := &regGateway{
		agents: 3,
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return f(ctx, "127.0.0.1:8081", &regClient{
				update: func(context.Context, *payload.Update_Request) (*payload.Object_Location, error) {
					return &payload.Object_Location{Name: "agent-0", Ips: []string{loopbackIP}}, nil
				},
			})
		},
		doMulti: func(context.Context, int, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return errDoMulti
		},
	}
	s := &server{gateway: gw, replica: 3}

	loc, err := s.Update(t.Context(), &payload.Update_Request{
		Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: []float32{1, 2, 3}},
		Config: &payload.Update_Config{DisableBalancedUpdate: true},
	})
	if !errors.Is(err, errDoMulti) {
		t.Errorf("got error: %v, want: %v", err, errDoMulti)
	}
	if loc != nil {
		t.Errorf("unexpected location on error: %v", loc)
	}
}

// Update used to set SkipStrictExistCheck on the caller's request after its own exist
// check, a write nothing read back (it builds fresh Remove/Insert requests).
func Test_server_Update_doesNotMutateRequest(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	tests := []struct {
		cfg  *payload.Update_Config
		name string
	}{
		{name: "nil config"},
		{name: "exist check requested", cfg: &payload.Update_Config{Timestamp: 3}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// The stored vector differs from the requested one, so Update proceeds to remove+insert.
			c := newOKClient([]float32{4, 5, 6})
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: syncGateway(c), replica: 1, eg: eg}

			req := &payload.Update_Request{
				Vector: &payload.Object_Vector{Id: uuid, Vector: []float32{1, 2, 3}},
				Config: test.cfg,
			}
			if _, err := s.Update(t.Context(), req); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if req.Config != test.cfg {
				t.Error("Update replaced the caller's req.Config")
			}
			if test.cfg != nil && test.cfg.SkipStrictExistCheck {
				t.Error("Update set SkipStrictExistCheck on the caller's config")
			}
			if c.removeIn == nil || c.insertIn == nil {
				t.Fatalf("remove/insert did not reach the agent: remove=%v insert=%v", c.removeIn, c.insertIn)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_server_Update(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Update_Request
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
// 		wantRes *payload.Object_Location
// 		err     error
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
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Location, err error) error {
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
// 			gotRes, err := s.Update(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamUpdate(t *testing.T) {
// 	type args struct {
// 		stream vald.Update_StreamUpdateServer
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
// 			err := s.StreamUpdate(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiUpdate(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Update_MultiRequest
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
// 			gotLocs, err := s.MultiUpdate(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_UpdateTimestamp(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Update_TimestampRequest
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
// 		wantRes *payload.Object_Location
// 		err     error
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
// 	defaultCheckFunc := func(w want, gotRes *payload.Object_Location, err error) error {
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
// 			gotRes, err := s.UpdateTimestamp(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotRes, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
