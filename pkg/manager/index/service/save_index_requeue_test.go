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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
)

func (a *fakeIndexAgent) SaveIndex(context.Context, *payload.Empty) (*payload.Empty, error) {
	a.saves.Add(1)
	if a.saveErr != nil {
		return nil, a.saveErr
	}
	return new(payload.Empty), nil
}

// Test_index_saveIndex_requeueOnError pins that an agent whose CreateIndex
// result is queued for saving stays queued when SaveIndex fails, so the next
// non-forced save tick retries it instead of waiting for the next CreateIndex
// or the force-save limit.
func Test_index_saveIndex_requeueOnError(t *testing.T) {
	t.Parallel()
	const addr = "agent-0"
	tests := []struct {
		saveErr    error
		name       string
		wantSaves  uint32
		queued     bool
		force      bool
		wantQueued bool
		wantErr    bool
	}{
		{
			name:       "failed save of a queued agent keeps it queued",
			queued:     true,
			saveErr:    status.Error(codes.Internal, "save failed"),
			wantSaves:  1,
			wantQueued: true,
			wantErr:    true,
		},
		{
			name:       "failed forced save of a queued agent keeps it queued",
			queued:     true,
			force:      true,
			saveErr:    status.Error(codes.Unavailable, "agent unavailable"),
			wantSaves:  1,
			wantQueued: true,
			wantErr:    true,
		},
		{
			name:      "successful save dequeues the agent",
			queued:    true,
			wantSaves: 1,
		},
		{
			name:      "skipped save (FailedPrecondition) dequeues the agent",
			queued:    true,
			saveErr:   status.Error(codes.FailedPrecondition, "nothing to save"),
			wantSaves: 1,
		},
		{
			name:      "failed forced save of an agent that was not queued does not queue it",
			force:     true,
			saveErr:   status.Error(codes.Internal, "save failed"),
			wantSaves: 1,
			wantErr:   true,
		},
		{
			name: "not queued and not forced does not call SaveIndex",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := &fakeIndexAgent{saveErr: tt.saveErr}
			conn := serveFakeIndexAgent(t, addr, a)
			idx := &index{
				client: &fakeIndexDiscoverer{client: &fakeIndexGRPCClient{
					conns: map[string]*grpc.ClientConn{addr: conn},
					order: []string{addr},
				}},
				saveIndexConcurrency: 1,
			}
			if tt.queued {
				idx.shouldSaveList.Store(addr, true)
			}

			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			if err := idx.saveIndex(ctx, tt.force); (err != nil) != tt.wantErr {
				t.Errorf("saveIndex error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := a.saves.Load(); got != tt.wantSaves {
				t.Errorf("SaveIndex calls = %d, want %d", got, tt.wantSaves)
			}
			if _, got := idx.shouldSaveList.Load(addr); got != tt.wantQueued {
				t.Errorf("queued after saveIndex = %v, want %v", got, tt.wantQueued)
			}
		})
	}
}
