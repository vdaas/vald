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
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// tsAgent mimics the NGT agent's timestamp handling for a single uuid: Insert takes
// the vector's timestamp and falls back to the config's only when it is 0
// (pkg/agent/core/ngt/handler/grpc/insert.go), and an insert older than a queued
// delete is silently dropped while still reporting success (vqueue.PushInsert).
type tsAgent struct {
	*regClient
	updateTS func(ctx context.Context, in *payload.Update_TimestampRequest) (*payload.Object_Location, error)

	mu       sync.Mutex
	delTS    int64
	stored   bool
	storedTS int64
	inserts  int
}

func newTSAgent(name, ip string, stored []float32, storedTS int64) *tsAgent {
	a := &tsAgent{stored: stored != nil, storedTS: storedTS}
	loc := &payload.Object_Location{Uuid: fixtureUUID, Name: name, Ips: []string{ip}}
	a.regClient = &regClient{
		insert: func(_ context.Context, in *payload.Insert_Request) (*payload.Object_Location, error) {
			ts := in.GetVector().GetTimestamp()
			if ts == 0 {
				ts = in.GetConfig().GetTimestamp()
			}
			a.mu.Lock()
			defer a.mu.Unlock()
			a.inserts++
			if a.delTS > ts {
				return loc, nil // dropped behind the newer delete
			}
			a.stored, a.storedTS = true, ts
			return loc, nil
		},
		remove: func(_ context.Context, in *payload.Remove_Request) (*payload.Object_Location, error) {
			a.mu.Lock()
			defer a.mu.Unlock()
			a.delTS = in.GetConfig().GetTimestamp()
			a.stored = false
			return loc, nil
		},
		exists: func(context.Context, *payload.Object_ID) (*payload.Object_ID, error) {
			return &payload.Object_ID{Id: fixtureUUID}, nil
		},
		getObject: func(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error) {
			return &payload.Object_Vector{Id: fixtureUUID, Vector: stored, Timestamp: storedTS}, nil
		},
		// Like the agent's Update: delete at the config timestamp, insert at +1.
		update: func(_ context.Context, in *payload.Update_Request) (*payload.Object_Location, error) {
			if ts := in.GetConfig().GetTimestamp(); ts > 0 {
				a.mu.Lock()
				defer a.mu.Unlock()
				a.delTS = ts
				a.stored, a.storedTS = true, ts+1
			}
			return loc, nil
		},
	}
	a.updateTS = func(context.Context, *payload.Update_TimestampRequest) (*payload.Object_Location, error) {
		return loc, nil
	}
	return a
}

func (a *tsAgent) UpdateTimestamp(
	ctx context.Context, in *payload.Update_TimestampRequest, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return a.updateTS(ctx, in)
}

func (a *tsAgent) state() (stored bool, ts, delTS int64, inserts int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.stored, a.storedTS, a.delTS, a.inserts
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

// twoAgentGateway broadcasts to `first` only and runs DoMulti against `second`, so the
// handler sees a replica shortage it tops up on `second`.
func twoAgentGateway(first, second vald.Client) *regGateway {
	return &regGateway{
		agents: 2,
		broadCast: func(ctx context.Context, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return f(ctx, "127.0.0.1:8081", first)
		},
		doMulti: func(ctx context.Context, _ int, f func(ctx context.Context, target string, vc vald.Client, copts ...grpc.CallOption) error) error {
			return f(ctx, "127.0.0.2:8081", second)
		},
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
