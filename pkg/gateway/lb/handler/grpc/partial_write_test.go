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
	"github.com/vdaas/vald/internal/sync/atomic"
)

// pwAgent answers every write with err (nil: success located on its name) and
// counts the inserts it was asked for.
type pwAgent struct {
	*regClient
	removeErr error
	insertErr error
	updateErr error
	name      string
	inserts   atomic.Int64
}

func newPWAgent(name string, removeErr, insertErr, updateErr error) *pwAgent {
	a := &pwAgent{name: name, removeErr: removeErr, insertErr: insertErr, updateErr: updateErr}
	loc := func(uuid string) *payload.Object_Location {
		return &payload.Object_Location{Uuid: uuid, Name: name, Ips: []string{name}}
	}
	a.regClient = &regClient{
		remove: func(_ context.Context, in *payload.Remove_Request) (*payload.Object_Location, error) {
			if a.removeErr != nil {
				return nil, a.removeErr
			}
			return loc(in.GetId().GetId()), nil
		},
		insert: func(_ context.Context, in *payload.Insert_Request) (*payload.Object_Location, error) {
			a.inserts.Add(1)
			if a.insertErr != nil {
				return nil, a.insertErr
			}
			return loc(in.GetVector().GetId()), nil
		},
		update: func(_ context.Context, in *payload.Update_Request) (*payload.Object_Location, error) {
			if a.updateErr != nil {
				return nil, a.updateErr
			}
			return loc(in.GetVector().GetId()), nil
		},
	}
	return a
}

func (a *pwAgent) UpdateTimestamp(
	_ context.Context, in *payload.Update_TimestampRequest, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	if a.updateErr != nil {
		return nil, a.updateErr
	}
	return &payload.Object_Location{Uuid: in.GetId(), Name: a.name, Ips: []string{a.name}}, nil
}

// pwGateway mimics the lb service gateway: BroadCast calls every agent and joins
// their errors in target order, DoMulti calls the agents in order until num of
// them succeeded (a failed agent is not counted) and only reports a missing
// connection, as OrderedRange logs the per-agent failures.
func pwGateway(agents ...*pwAgent) *regGateway {
	return &regGateway{
		agents: len(agents),
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			var errs []error
			for _, a := range agents {
				if err := f(ctx, a.name, a); err != nil {
					errs = append(errs, err)
				}
			}
			return errors.Join(errs...)
		},
		doMulti: func(ctx context.Context, num int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			var cur int
			for _, a := range agents {
				if cur >= num {
					break
				}
				if err := f(ctx, a.name, a); err == nil {
					cur++
				}
			}
			return nil
		},
	}
}

func pwCode(t *testing.T, err error) codes.Code {
	t.Helper()
	st, ok := status.FromError(err)
	if !ok || st == nil {
		t.Fatalf("not a status error: %v", err)
	}
	return st.Code()
}

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
