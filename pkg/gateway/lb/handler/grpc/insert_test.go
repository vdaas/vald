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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Test_server_MultiInsert_joinedErrors pins the error MultiInsert returns now
// that per-item failures are collected in a slice and joined once after the
// workers finish: distinct failures are all kept, identical messages collapse
// to one, and the status code is still derived from the members.
func Test_server_MultiInsert_joinedErrors(t *testing.T) {
	t.Parallel()
	const n = 5
	tests := []struct {
		msg         func(id string) string
		name        string
		wantMembers int // 0: a single (non-joined) error
	}{
		{
			name:        "distinct messages are all kept",
			msg:         func(id string) string { return "agent failure: " + id },
			wantMembers: n,
		},
		{
			name: "identical messages collapse to one",
			msg:  func(string) string { return "agent failure" },
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c := &regClient{
				insert: func(_ context.Context, in *payload.Insert_Request) (*payload.Object_Location, error) {
					return nil, status.Error(codes.NotFound, test.msg(in.GetVector().GetId()))
				},
			}
			s := &server{gateway: syncGateway(c), replica: 1, multiConcurrency: n}
			reqs := &payload.Insert_MultiRequest{Requests: make([]*payload.Insert_Request, 0, n)}
			for i := range n {
				reqs.Requests = append(reqs.Requests, &payload.Insert_Request{
					Vector: &payload.Object_Vector{Id: "id-" + strconv.Itoa(i), Vector: []float32{1, 2, 3}},
					Config: &payload.Insert_Config{SkipStrictExistCheck: true},
				})
			}
			_, err := s.MultiInsert(t.Context(), reqs)
			if err == nil {
				t.Fatal("MultiInsert returned nil, want the per-item failures")
			}
			if st, ok := status.FromError(err); !ok || st == nil || st.Code() != codes.NotFound {
				t.Errorf("status code = %v (ok=%t), want NotFound: %v", st.Code(), ok, err)
			}
			var u interface{ Unwrap() []error }
			joined := errors.As(err, &u)
			switch {
			case test.wantMembers == 0 && joined:
				t.Errorf("got joined error with %d members, want a single error: %v", len(u.Unwrap()), err)
			case test.wantMembers != 0 && !joined:
				t.Errorf("got %T, want a joined error: %v", err, err)
			case joined && len(u.Unwrap()) != test.wantMembers:
				t.Errorf("members = %d, want %d: %v", len(u.Unwrap()), test.wantMembers, err)
			}
			for i := range n {
				if m := test.msg("id-" + strconv.Itoa(i)); !strings.Contains(err.Error(), m) {
					t.Errorf("error %q does not contain %q", err.Error(), m)
				}
			}
		})
	}
}

// Insert used to fail when one agent failed although the other agents stored all
// the replicas, and counted an agent that timed out as a stored replica.
func Test_server_Insert_partialAgentFailure(t *testing.T) {
	t.Parallel()
	// survivorsName is the location name when agent a stores no replica and b and c do.
	const survivorsName = "b,c"
	internal := status.Error(codes.Internal, "agent internal")
	tests := []struct {
		name     string
		wantName string
		errs     []error // per agent, in DoMulti order
		replica  int
		wantCode codes.Code // codes.OK: success expected
		// noCall makes DoMulti skip every agent without calling f, as it does for
		// agents whose circuit breaker is open: no agent stored the vector.
		noCall bool
	}{
		{
			name:     "one agent fails, replicas satisfied by the others",
			errs:     []error{internal, nil, nil},
			replica:  2,
			wantName: survivorsName,
		},
		{
			name:     "one agent fails, replicas short",
			errs:     []error{internal, nil},
			replica:  2,
			wantName: "b",
		},
		{
			name:     "timed-out agent is not counted as a replica",
			errs:     []error{context.DeadlineExceeded, nil, nil},
			replica:  2,
			wantName: survivorsName,
		},
		{
			name:     "canceled agent is not counted as a replica",
			errs:     []error{context.Canceled, nil, nil},
			replica:  2,
			wantName: survivorsName,
		},
		{
			name:     "every agent fails",
			errs:     []error{internal, internal},
			replica:  2,
			wantCode: codes.Internal,
		},
		{
			name:     "every agent times out",
			errs:     []error{context.DeadlineExceeded},
			replica:  1,
			wantCode: codes.DeadlineExceeded,
		},
		{
			name:     "no agent called",
			errs:     []error{nil, nil},
			replica:  2,
			noCall:   true,
			wantCode: codes.Unavailable,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			agents := make([]*pwAgent, 0, len(test.errs))
			for i, err := range test.errs {
				agents = append(agents, newPWAgent(string(rune('a'+i)), nil, err, nil))
			}
			gw := pwGateway(agents...)
			if test.noCall {
				gw.doMulti = func(context.Context, int, func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					return nil
				}
			}
			s := &server{gateway: gw, replica: test.replica}
			loc, err := s.Insert(t.Context(), &payload.Insert_Request{
				Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: []float32{1, 2, 3}},
				Config: &payload.Insert_Config{SkipStrictExistCheck: true, Timestamp: 1},
			})
			if test.wantCode != codes.OK {
				if err == nil {
					t.Fatalf("Insert returned %v, want an error", loc)
				}
				if got := pwCode(t, err); got != test.wantCode {
					t.Errorf("status code = %v, want %v: %v", got, test.wantCode, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Insert returned %v, want success", err)
			}
			if loc.GetName() != test.wantName {
				t.Errorf("Name = %q, want %q", loc.GetName(), test.wantName)
			}
		})
	}
}

// scriptedInsertClient answers the i-th Insert with attempts[i]: a nil entry is a
// success located on host, anything else is returned as the agent's error.
type scriptedInsertClient struct {
	*regClient
	host     string
	attempts []error
	calls    int
}

func (c *scriptedInsertClient) Insert(
	_ context.Context, in *payload.Insert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	i := c.calls
	c.calls++
	if i < len(c.attempts) && c.attempts[i] != nil {
		return nil, c.attempts[i]
	}
	return &payload.Object_Location{Uuid: in.GetVector().GetId(), Name: c.host, Ips: []string{c.host}}, nil
}

// Insert used to keep an agent's first failure after a retry of the same agent
// (client backoff, or DoMulti's second pass) succeeded, failing a request whose
// replicas were all written.
func Test_server_Insert_retriedTargetError(t *testing.T) {
	t.Parallel()
	unavailable := status.Error(codes.Unavailable, "agent unavailable")
	internal := status.Error(codes.Internal, "agent internal")
	exists := status.Error(codes.AlreadyExists, "already exists")
	const (
		uuid     = "uuid"
		wantName = "a,b"
	)
	tests := []struct {
		attempts map[string][]error
		name     string
		wantName string
		order    []string // targets in the order DoMulti calls f
		wantIPs  int
		wantErr  bool
	}{
		{
			name:     "backoff retry of the same target succeeds",
			attempts: map[string][]error{"a": {unavailable, nil}, "b": {nil}},
			order:    []string{"a", "a", "b"},
			wantName: wantName,
			wantIPs:  2,
		},
		{
			name:     "second pass retries a failed target and it succeeds",
			attempts: map[string][]error{"a": {internal, nil}, "b": {nil}},
			order:    []string{"a", "b", "a"},
			wantName: wantName,
			wantIPs:  2,
		},
		{
			name:     "retry answers AlreadyExists",
			attempts: map[string][]error{"a": {unavailable, exists}, "b": {nil}},
			order:    []string{"a", "a", "b"},
			wantName: wantName,
			wantIPs:  2,
		},
		{
			// another agent stored the vector, so the request succeeds under-replicated.
			name:     "every attempt of a target fails",
			attempts: map[string][]error{"a": {unavailable, internal}, "b": {nil}},
			order:    []string{"a", "a", "b"},
			wantName: "b",
			wantIPs:  1,
		},
		{
			name:     "every attempt of every target fails",
			attempts: map[string][]error{"a": {unavailable, internal}, "b": {internal}},
			order:    []string{"a", "a", "b"},
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			clients := make(map[string]vald.Client, len(test.attempts))
			for target, attempts := range test.attempts {
				clients[target] = &scriptedInsertClient{regClient: new(regClient), host: target, attempts: attempts}
			}
			gw := &regGateway{
				doMulti: func(ctx context.Context, _ int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					// The real DoMulti/OrderedRange only log per-target failures.
					for _, target := range test.order {
						if err := f(ctx, target, clients[target]); err != nil {
							t.Logf("target %s: %v", target, err)
						}
					}
					return nil
				},
			}
			s := &server{gateway: gw, replica: 2}
			loc, err := s.Insert(t.Context(), &payload.Insert_Request{
				Vector: &payload.Object_Vector{Id: uuid, Vector: []float32{1, 2, 3}},
				Config: &payload.Insert_Config{SkipStrictExistCheck: true, Timestamp: 1},
			})
			if test.wantErr {
				if err == nil {
					t.Fatalf("Insert returned %v, want the failed target's error", loc)
				}
				if st, ok := status.FromError(err); !ok || st == nil || st.Code() != codes.Internal {
					t.Errorf("status code = %v (ok=%t), want the last attempt's Internal: %v", st.Code(), ok, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Insert returned %v, want success after the retry", err)
			}
			if loc.GetName() != test.wantName {
				t.Errorf("Name = %q, want %q", loc.GetName(), test.wantName)
			}
			if len(loc.GetIps()) != test.wantIPs {
				t.Errorf("Ips = %v, want %d entries", loc.GetIps(), test.wantIPs)
			}
		})
	}
}

// The replica top-up inserts used to forward a vector carrying a different timestamp
// from the one the other replicas got, and agents prefer the vector's timestamp.
func Test_server_shortageInsertUsesIntendedTimestamp(t *testing.T) {
	t.Parallel()
	vec := []float32{1, 2, 3}
	tests := []struct {
		call       func(s *server) (*payload.Object_Location, error)
		name       string
		wantTS     int64
		matchFirst bool // the top-up replica must hold the updated replica's timestamp
	}{
		{
			name:       "Update with DisableBalancedUpdate",
			wantTS:     101, // the updated replica stores config timestamp + 1
			matchFirst: true,
			call: func(s *server) (*payload.Object_Location, error) {
				return s.Update(t.Context(), &payload.Update_Request{
					Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec, Timestamp: 5},
					Config: &payload.Update_Config{DisableBalancedUpdate: true, Timestamp: 100},
				})
			},
		},
		{
			name:   "Update with DisableBalancedUpdate, no config timestamp",
			wantTS: 0, // the agent falls back to its own clock, as the Update replicas do
			call: func(s *server) (*payload.Object_Location, error) {
				return s.Update(t.Context(), &payload.Update_Request{
					Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: vec, Timestamp: 5},
					Config: &payload.Update_Config{DisableBalancedUpdate: true},
				})
			},
		},
		{
			name:   vald.UpdateTimestampRPCName,
			wantTS: 500,
			call: func(s *server) (*payload.Object_Location, error) {
				return s.UpdateTimestamp(t.Context(), &payload.Update_TimestampRequest{Id: fixtureUUID, Timestamp: 500})
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			// The first agent answers GetObject with the stored (old) timestamp.
			first := newTSAgent("agent-0", loopbackIP, []float32{4, 5, 6}, 1)
			second := newTSAgent("agent-1", "127.0.0.2", nil, 0)
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: twoAgentGateway(first, second), replica: 2, eg: eg}

			loc, err := test.call(s)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			stored, ts, _, inserts := second.state()
			if inserts != 1 || !stored {
				t.Fatalf("top-up insert: inserts=%d stored=%v, want 1 stored insert", inserts, stored)
			}
			if ts != test.wantTS {
				t.Errorf("top-up replica timestamp = %d, want %d", ts, test.wantTS)
			}
			if test.matchFirst {
				if _, fts, _, _ := first.state(); fts != ts {
					t.Errorf("top-up replica timestamp = %d, updated replica timestamp = %d", ts, fts)
				}
			}
			if got := loc.GetIps(); !slices.Equal(got, []string{loopbackIP, "127.0.0.2"}) {
				t.Errorf("location ips = %v", got)
			}
		})
	}
}

// An agent answering AlreadyExists counts toward the replicas, but its address used to
// be left out of the returned location (Update records it).
func Test_server_Insert_alreadyExistsRecordsLocation(t *testing.T) {
	t.Parallel()
	const (
		existingIP = "10.0.0.1"
		portless   = "agent-host"
	)
	exists := &regClient{
		insert: func(context.Context, *payload.Insert_Request) (*payload.Object_Location, error) {
			return nil, status.WrapWithAlreadyExists("Insert API: uuid already exists", errors.ErrUUIDAlreadyExists(fixtureUUID))
		},
	}
	fresh := newTSAgent("agent-1", "10.0.0.2", nil, 0)
	tests := []struct {
		agents   map[string]vald.Client
		name     string
		wantName string
		wantIPs  []string
	}{
		{
			name:     "all agents already hold the uuid",
			agents:   map[string]vald.Client{"10.0.0.1:8081": exists},
			wantIPs:  []string{existingIP},
			wantName: existingIP,
		},
		{
			name:     "target without a port",
			agents:   map[string]vald.Client{portless: exists},
			wantIPs:  []string{portless},
			wantName: portless,
		},
		{
			name:     "one existing, one inserted",
			agents:   map[string]vald.Client{"10.0.0.1:8081": exists, "10.0.0.2:8081": fresh},
			wantIPs:  []string{existingIP, "10.0.0.2"},
			wantName: "10.0.0.1,agent-1",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			targets := make([]string, 0, len(test.agents))
			for target := range test.agents {
				targets = append(targets, target)
			}
			slices.Sort(targets)
			gw := &regGateway{
				agents: len(targets),
				doMulti: func(ctx context.Context, _ int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
					var errs error
					for _, target := range targets {
						if err := f(ctx, target, test.agents[target]); err != nil {
							errs = errors.Join(errs, err)
						}
					}
					return errs
				},
			}
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: gw, replica: len(targets), eg: eg}

			loc, err := s.Insert(t.Context(), &payload.Insert_Request{
				Vector: &payload.Object_Vector{Id: fixtureUUID, Vector: []float32{1, 2, 3}},
				Config: &payload.Insert_Config{SkipStrictExistCheck: true, Timestamp: 7},
			})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			ips := slices.Clone(loc.GetIps())
			slices.Sort(ips)
			if !slices.Equal(ips, test.wantIPs) {
				t.Errorf("ips = %v, want %v", ips, test.wantIPs)
			}
			if loc.GetName() != test.wantName {
				t.Errorf("name = %q, want %q", loc.GetName(), test.wantName)
			}
		})
	}
}

// Insert used to write the exist-check flag and the timestamp into the caller's request,
// which MultiInsert shares with its other readers. It must forward a new request instead.
func Test_server_Insert_doesNotMutateRequest(t *testing.T) {
	t.Parallel()
	const uuid = "uuid"
	filters := &payload.Filter_Config{Targets: []*payload.Filter_Target{{Host: "filter", Port: 1}}}
	tests := []struct {
		cfg        *payload.Insert_Config
		wantFilter *payload.Filter_Config
		name       string
		wantTS     int64 // 0: any non-zero timestamp
	}{
		{name: "nil config (exist check runs)"},
		{name: "skip check, no timestamp", cfg: &payload.Insert_Config{SkipStrictExistCheck: true, Filters: filters}, wantFilter: filters},
		{name: "exist check requested, filters kept", cfg: &payload.Insert_Config{Filters: filters}, wantFilter: filters},
		{name: "timestamp already set", cfg: &payload.Insert_Config{SkipStrictExistCheck: true, Timestamp: 7}, wantTS: 7},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			c := newOKClient([]float32{4, 5, 6})
			// Exists finds nothing, so Insert may proceed.
			c.exists = func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
				return &payload.Object_ID{}, nil
			}
			eg, _ := errgroup.New(t.Context())
			s := &server{gateway: syncGateway(c), replica: 1, eg: eg}

			vector := &payload.Object_Vector{Id: uuid, Vector: []float32{1, 2, 3}}
			req := &payload.Insert_Request{Vector: vector, Config: test.cfg}
			var before payload.Insert_Config
			if test.cfg != nil {
				before.SkipStrictExistCheck, before.Filters, before.Timestamp = test.cfg.SkipStrictExistCheck, test.cfg.Filters, test.cfg.Timestamp
			}
			if _, err := s.Insert(t.Context(), req); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if req.Config != test.cfg {
				t.Error("Insert replaced the caller's req.Config")
			}
			if test.cfg != nil && (test.cfg.SkipStrictExistCheck != before.SkipStrictExistCheck ||
				test.cfg.Filters != before.Filters || test.cfg.Timestamp != before.Timestamp) {
				t.Errorf("Insert modified the caller's config: got %v", test.cfg)
			}
			got := c.insertIn
			if got == nil {
				t.Fatal("no request reached the agent")
			}
			if got.GetVector() != vector {
				t.Error("the forwarded request does not carry the caller's vector")
			}
			if !got.GetConfig().GetSkipStrictExistCheck() {
				t.Error("the forwarded request must skip the agent's exist check")
			}
			if ts := got.GetConfig().GetTimestamp(); ts == 0 || (test.wantTS != 0 && ts != test.wantTS) {
				t.Errorf("forwarded timestamp = %d, want %d (0: any non-zero)", ts, test.wantTS)
			}
			if got.GetConfig().GetFilters() != test.wantFilter {
				t.Errorf("forwarded filters = %v, want %v", got.GetConfig().GetFilters(), test.wantFilter)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_server_Insert(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		req *payload.Insert_Request
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
// 		wantCe *payload.Object_Location
// 		err    error
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
// 	defaultCheckFunc := func(w want, gotCe *payload.Object_Location, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotCe, w.wantCe) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotCe, w.wantCe)
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
// 			gotCe, err := s.Insert(test.args.ctx, test.args.req)
// 			if err := checkFunc(test.want, gotCe, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_StreamInsert(t *testing.T) {
// 	type args struct {
// 		stream vald.Insert_StreamInsertServer
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
// 			err := s.StreamInsert(test.args.stream)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_server_MultiInsert(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		reqs *payload.Insert_MultiRequest
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
// 			gotLocs, err := s.MultiInsert(test.args.ctx, test.args.reqs)
// 			if err := checkFunc(test.want, gotLocs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
