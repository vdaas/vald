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

	"github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/atomic"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// fakeIndexAgent serves IndexInfo with a live uncommitted count and counts CreateIndex calls.
type fakeIndexAgent struct {
	core.UnimplementedAgentServer
	vald.UnimplementedIndexServer

	saveErr     error // returned by SaveIndex
	createErr   error // returned by CreateIndex
	uncommitted atomic.Uint32
	creates     atomic.Uint32
	saves       atomic.Uint32
}

func (a *fakeIndexAgent) IndexInfo(
	context.Context, *payload.Empty,
) (*payload.Info_Index_Count, error) {
	return &payload.Info_Index_Count{Uncommitted: a.uncommitted.Load()}, nil
}

func (a *fakeIndexAgent) CreateIndex(
	context.Context, *payload.Control_CreateIndexRequest,
) (*payload.Empty, error) {
	a.creates.Add(1)
	if a.createErr != nil {
		return nil, a.createErr
	}
	a.uncommitted.Store(0)
	return new(payload.Empty), nil
}

// fakeIndexGRPCClient ranges over fixed conns sequentially.
type fakeIndexGRPCClient struct {
	grpc.Client
	conns map[string]*grpc.ClientConn
	order []string
}

func (f *fakeIndexGRPCClient) RangeConcurrent(
	ctx context.Context,
	_ int,
	fn func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	return f.OrderedRangeConcurrent(ctx, f.order, 1, fn)
}

func (f *fakeIndexGRPCClient) OrderedRangeConcurrent(
	ctx context.Context,
	order []string,
	_ int,
	fn func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) error {
	for _, addr := range order {
		if conn, ok := f.conns[addr]; ok {
			if err := fn(ctx, addr, conn); err != nil {
				return err
			}
		}
	}
	return nil
}

type fakeIndexDiscoverer struct {
	discoverer.Client
	client *fakeIndexGRPCClient
}

func (f *fakeIndexDiscoverer) GetClient() grpc.Client            { return f.client }
func (f *fakeIndexDiscoverer) GetAddrs(context.Context) []string { return f.client.order }

//nolint:unparam // addr is the passthrough target name and log label of the served agent; every current caller serves a single agent
func serveFakeIndexAgent(t *testing.T, addr string, a *fakeIndexAgent) *grpc.ClientConn {
	t.Helper()
	lis := bufconn.Listen(1 << 20)
	srv := ggrpc.NewServer()
	core.RegisterAgentServer(srv, a)
	vald.RegisterIndexServer(srv, a)
	served := make(chan struct{})
	go func() {
		defer close(served)
		if err := srv.Serve(lis); err != nil {
			t.Logf("fake agent %s stopped: %v", addr, err)
		}
	}()
	conn, err := ggrpc.NewClient("passthrough:///"+addr,
		ggrpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		ggrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to create client for %s: %v", addr, err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Logf("failed to close conn to %s: %v", addr, err)
		}
		srv.Stop()
		select {
		case <-served:
		case <-time.After(10 * time.Second):
			t.Errorf("fake agent %s did not stop", addr)
		}
	})
	return conn
}

// Test_index_createIndex_refreshesInfosBeforeSkip pins that the per-agent skip decision uses
// the agent's current uncommitted count, not the snapshot left by the previous run.
func Test_index_createIndex_refreshesInfosBeforeSkip(t *testing.T) {
	t.Parallel()
	const addr = "agent-0"
	tests := []struct {
		name            string
		stale           uint32 // snapshot left by the previous run
		live            uint32 // what the agent reports now
		lowIndexSkip    bool
		wantCreateCalls uint32
	}{
		{
			name:            "stale zero but live backlog creates the index",
			stale:           0,
			live:            100,
			lowIndexSkip:    true,
			wantCreateCalls: 1,
		},
		{
			name:            "stale low but live backlog over minUncommitted creates the index",
			stale:           1,
			live:            100,
			lowIndexSkip:    true,
			wantCreateCalls: 1,
		},
		{
			name:            "stale backlog but live zero skips the agent",
			stale:           100,
			live:            0,
			lowIndexSkip:    false,
			wantCreateCalls: 0,
		},
		{
			name:            "stale backlog but live low skips the agent on a low-index-skip run",
			stale:           100,
			live:            1,
			lowIndexSkip:    true,
			wantCreateCalls: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			a := new(fakeIndexAgent)
			a.uncommitted.Store(tt.live)
			conn := serveFakeIndexAgent(t, addr, a)
			idx := &index{
				client: &fakeIndexDiscoverer{client: &fakeIndexGRPCClient{
					conns: map[string]*grpc.ClientConn{addr: conn},
					order: []string{addr},
				}},
				createIndexConcurrency: 1,
				minUncommitted:         10,
			}
			idx.indexInfos.Store(addr, &payload.Info_Index_Count{Uncommitted: tt.stale})

			ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
			defer cancel()
			if err := idx.createIndex(ctx, tt.lowIndexSkip); err != nil {
				t.Fatalf("createIndex error = %v", err)
			}
			if got := a.creates.Load(); got != tt.wantCreateCalls {
				t.Errorf("CreateIndex calls = %d, want %d", got, tt.wantCreateCalls)
			}
			// The post-run refresh still records the state after indexing.
			info, ok := idx.indexInfos.Load(addr)
			if !ok {
				t.Fatal("indexInfos entry is missing after createIndex")
			}
			if got, want := info.GetUncommitted(), a.uncommitted.Load(); got != want {
				t.Errorf("indexInfos uncommitted = %d, want %d", got, want)
			}
		})
	}
}
