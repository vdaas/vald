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

package service

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/sync"
)

// Test_index_createIndex_queuesOnError pins that an agent whose CreateIndex fails with anything
// but FailedPrecondition is still queued for the next save tick, because the agent may have
// built a new index generation before returning the error.
func Test_index_createIndex_queuesOnError(t *testing.T) {
	t.Parallel()
	const addr = "agent-0"
	tests := []struct {
		createErr  error
		name       string
		wantQueued bool
		wantErr    bool
	}{
		{
			name:       "successful CreateIndex queues the agent",
			wantQueued: true,
		},
		{
			name:       "Internal error (e.g. statistics or export failure after build) queues the agent",
			createErr:  status.Error(codes.Internal, "failed to export index info"),
			wantQueued: true,
			wantErr:    true,
		},
		{
			name:       "Canceled after the delete phase queues the agent",
			createErr:  status.Error(codes.Canceled, "canceled"),
			wantQueued: true,
			wantErr:    true,
		},
		{
			name:       "Unavailable (lost reply) queues the agent",
			createErr:  status.Error(codes.Unavailable, "connection lost"),
			wantQueued: true,
			wantErr:    true,
		},
		{
			name:      "FailedPrecondition (nothing to index) does not queue the agent",
			createErr: status.Error(codes.FailedPrecondition, "uncommitted index is empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := &fakeIndexAgent{createErr: tt.createErr}
			a.uncommitted.Store(100)
			conn := serveFakeIndexAgent(t, addr, a)
			idx := &index{
				client: &fakeIndexDiscoverer{client: &fakeIndexGRPCClient{
					conns: map[string]*grpc.ClientConn{addr: conn},
					order: []string{addr},
				}},
				createIndexConcurrency: 1,
			}

			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			if err := idx.createIndex(ctx, false); (err != nil) != tt.wantErr {
				t.Errorf("createIndex error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := a.creates.Load(); got != 1 {
				t.Errorf("CreateIndex calls = %d, want 1", got)
			}
			if _, got := idx.shouldSaveList.Load(addr); got != tt.wantQueued {
				t.Errorf("queued after createIndex = %v, want %v", got, tt.wantQueued)
			}
		})
	}
}

// Test_index_saveOnTick_forcedSaveRetriedWhenBusy pins that a forced-save tick that finds a
// save already running is rescheduled after saveIndexDuration instead of being dropped until
// the next saveIndexDurationLimit tick, and that a regular tick is still simply dropped.
func Test_index_saveOnTick_forcedSaveRetriedWhenBusy(t *testing.T) {
	t.Parallel()
	const (
		addr     = "agent-0"
		retry    = 10 * time.Millisecond
		longTick = time.Hour
		// waitFire bounds the wait for a rescheduled tick; quietFor is how long a ticker that
		// was not rescheduled is observed (it cannot fire before longTick).
		waitFire = 10 * time.Second
		quietFor = 100 * time.Millisecond
	)
	tests := []struct {
		name      string
		busy      bool
		force     bool
		wantFire  bool
		wantSaves uint32
	}{
		{
			name:     "busy forced tick is rescheduled after saveIndexDuration",
			busy:     true,
			force:    true,
			wantFire: true,
		},
		{
			name:  "busy regular tick is dropped",
			busy:  true,
			force: false,
		},
		{
			name:      "idle forced tick saves and resets to saveIndexDurationLimit",
			force:     true,
			wantSaves: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := new(fakeIndexAgent)
			conn := serveFakeIndexAgent(t, addr, a)
			idx := &index{
				client: &fakeIndexDiscoverer{client: &fakeIndexGRPCClient{
					conns: map[string]*grpc.ClientConn{addr: conn},
					order: []string{addr},
				}},
				saveIndexConcurrency:   1,
				saveIndexDuration:      retry,
				saveIndexDurationLimit: longTick,
			}
			if !tt.force {
				idx.saveIndexDuration = longTick
			}
			tick := time.NewTicker(longTick)
			defer tick.Stop()
			var mu sync.Mutex
			if tt.busy {
				mu.Lock()
				defer mu.Unlock()
			}
			ech := make(chan error, 1)

			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			idx.saveOnTick(ctx, &mu, tick, tt.force, ech)

			if got := a.saves.Load(); got != tt.wantSaves {
				t.Errorf("SaveIndex calls = %d, want %d", got, tt.wantSaves)
			}
			select {
			case err := <-ech:
				t.Errorf("unexpected error reported: %v", err)
			default:
			}
			wait := quietFor
			if tt.wantFire {
				wait = waitFire
			}
			timer := time.NewTimer(wait)
			defer timer.Stop()
			select {
			case <-tick.C:
				if !tt.wantFire {
					t.Error("ticker fired, want it to keep its long period")
				}
			case <-timer.C:
				if tt.wantFire {
					t.Errorf("ticker did not fire within %v, want the forced save rescheduled after %v", waitFire, retry)
				}
			}
		})
	}
}
