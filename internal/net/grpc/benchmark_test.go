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
package grpc_test

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	agent "github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/servers"
	"github.com/vdaas/vald/internal/servers/server"
	handler "github.com/vdaas/vald/pkg/agent/core/ngt/handler/grpc"
	"github.com/vdaas/vald/pkg/agent/core/ngt/service"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// mockNGT implements service.NGT for benchmarking purposes.
type mockNGT struct {
	searchLatency time.Duration
}

func (m *mockNGT) Start(ctx context.Context) <-chan error { return nil }
func (m *mockNGT) Search(ctx context.Context, vec []float32, size uint32, epsilon, radius float32) (*payload.Search_Response, error) {
	if m.searchLatency > 0 {
		time.Sleep(m.searchLatency)
	}
	return &payload.Search_Response{
		Results: []*payload.Object_Distance{
			{Id: "id-1", Distance: 0.1},
			{Id: "id-2", Distance: 0.2},
			{Id: "id-3", Distance: 0.3},
		},
	}, nil
}
func (m *mockNGT) SearchByID(ctx context.Context, uuid string, size uint32, epsilon, radius float32) ([]float32, *payload.Search_Response, error) {
	return nil, nil, nil
}
func (m *mockNGT) LinearSearch(ctx context.Context, vec []float32, size uint32) (*payload.Search_Response, error) {
	return nil, nil
}
func (m *mockNGT) LinearSearchByID(ctx context.Context, uuid string, size uint32) ([]float32, *payload.Search_Response, error) {
	return nil, nil, nil
}
func (m *mockNGT) Insert(uuid string, vec []float32) (err error) { return nil }
func (m *mockNGT) InsertWithTime(uuid string, vec []float32, t int64) (err error) { return nil }
func (m *mockNGT) InsertMultiple(vecs map[string][]float32) (err error) { return nil }
func (m *mockNGT) InsertMultipleWithTime(vecs map[string][]float32, t int64) (err error) { return nil }
func (m *mockNGT) Update(uuid string, vec []float32) (err error) { return nil }
func (m *mockNGT) UpdateWithTime(uuid string, vec []float32, t int64) (err error) { return nil }
func (m *mockNGT) UpdateMultiple(vecs map[string][]float32) (err error) { return nil }
func (m *mockNGT) UpdateMultipleWithTime(vecs map[string][]float32, t int64) (err error) { return nil }
func (m *mockNGT) UpdateTimestamp(uuid string, ts int64, force bool) (err error) { return nil }
func (m *mockNGT) Delete(uuid string) (err error) { return nil }
func (m *mockNGT) DeleteWithTime(uuid string, t int64) (err error) { return nil }
func (m *mockNGT) DeleteMultiple(uuids ...string) (err error) { return nil }
func (m *mockNGT) DeleteMultipleWithTime(uuids []string, t int64) (err error) { return nil }
func (m *mockNGT) RegenerateIndexes(ctx context.Context) (err error) { return nil }
func (m *mockNGT) GetObject(uuid string) (vec []float32, timestamp int64, err error) {
	return nil, 0, nil
}
func (m *mockNGT) ListObjectFunc(ctx context.Context, f func(uuid string, oid uint32, timestamp int64) bool) {
}
func (m *mockNGT) Exists(uuid string) (uint32, bool) { return 0, false }
func (m *mockNGT) CreateIndex(ctx context.Context, poolSize uint32) (err error) { return nil }
func (m *mockNGT) SaveIndex(ctx context.Context) (err error) { return nil }
func (m *mockNGT) CreateAndSaveIndex(ctx context.Context, poolSize uint32) (err error) { return nil }
func (m *mockNGT) IsIndexing() bool { return false }
func (m *mockNGT) IsFlushing() bool { return false }
func (m *mockNGT) IsSaving() bool { return false }
func (m *mockNGT) Len() uint64 { return 0 }
func (m *mockNGT) NumberOfCreateIndexExecution() uint64 { return 0 }
func (m *mockNGT) NumberOfProactiveGCExecution() uint64 { return 0 }
func (m *mockNGT) UUIDs(context.Context) (uuids []string) { return nil }
func (m *mockNGT) InsertVQueueBufferLen() uint64 { return 0 }
func (m *mockNGT) DeleteVQueueBufferLen() uint64 { return 0 }
func (m *mockNGT) GetDimensionSize() int { return 128 }
func (m *mockNGT) BrokenIndexCount() uint64 { return 0 }
func (m *mockNGT) IndexStatistics() (*payload.Info_Index_Statistics, error) { return nil, nil }
func (m *mockNGT) IsStatisticsEnabled() bool { return false }
func (m *mockNGT) IndexProperty() (*payload.Info_Index_Property, error) { return nil, nil }
func (m *mockNGT) Close(ctx context.Context) error { return nil }

func startTestServer(t testing.TB, host string, port uint16, ngt service.NGT) (servers.Listener, string) {
	t.Helper()
	h, err := handler.New(handler.WithNGT(ngt))
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	s, err := server.New(
		server.WithName("agent-grpc"),
		server.WithServerMode(server.GRPC),
		server.WithGRPCRegisterar(func(srv *ggrpc.Server) {
			agent.RegisterAgentServer(srv, h)
			vald.RegisterValdServer(srv, h)
		}),
		server.WithHost(host),
		server.WithPort(port),
	)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	l := servers.New(
		servers.WithServer(s),
	)

	go l.ListenAndServe(context.Background())
	return l, addr
}

func BenchmarkSearchRPC(b *testing.B) {
	host := "127.0.0.1"
	port := uint16(9091)

	mock := &mockNGT{
		searchLatency: 0,
	}

	l, addr := startTestServer(b, host, port, mock)
	defer l.Shutdown(context.Background())

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)

	req := &payload.Search_Request{
		Vector: make([]float32, 128),
		Config: &payload.Search_Config{
			Num:     10,
			Radius:  -1,
			Epsilon: 0.01,
			Timeout: 1000,
		},
	}

	b.Run("Sequential", func(b *testing.B) {
		ctx := context.Background()
		client, err := grpc.New(
			"benchmark-client-seq",
			grpc.WithAddrs(addr),
			grpc.WithDialOptions(
				ggrpc.WithTransportCredentials(insecure.NewCredentials()),
			),
			grpc.WithConnectionPoolSize(1),
		)
		if err != nil {
			b.Fatalf("failed to create client: %v", err)
		}
		defer client.Close(ctx)
		if _, err := client.Connect(ctx, addr); err != nil {
			b.Fatalf("failed to connect: %v", err)
		}

		b.ResetTimer()
		for b.Loop() {
			_, err := client.Search(ctx, req)
			if err != nil {
				b.Errorf("Search failed: %v", err)
			}
		}
	})

	b.Run("Parallel", func(b *testing.B) {
		ctx := context.Background()
		client, err := grpc.New(
			"benchmark-client-par",
			grpc.WithAddrs(addr),
			grpc.WithDialOptions(
				ggrpc.WithTransportCredentials(insecure.NewCredentials()),
			),
			grpc.WithConnectionPoolSize(uint64(runtime.GOMAXPROCS(0))),
		)
		if err != nil {
			b.Fatalf("failed to create client: %v", err)
		}
		defer client.Close(ctx)
		if _, err := client.Connect(ctx, addr); err != nil {
			b.Fatalf("failed to connect: %v", err)
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := client.Search(ctx, req)
				if err != nil {
					b.Errorf("Search failed: %v", err)
				}
			}
		})
	})
}

func BenchmarkInsertRPC(b *testing.B) {
	host := "127.0.0.1"
	port := uint16(9092)

	mock := &mockNGT{}

	l, addr := startTestServer(b, host, port, mock)
	defer l.Shutdown(context.Background())

	time.Sleep(100 * time.Millisecond)

	req := &payload.Insert_Request{
		Vector: &payload.Object_Vector{
			Id:     "test-id",
			Vector: make([]float32, 128),
		},
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: true,
		},
	}

	b.Run("Parallel", func(b *testing.B) {
		ctx := context.Background()
		client, err := grpc.New(
			"benchmark-client-insert",
			grpc.WithAddrs(addr),
			grpc.WithDialOptions(
				ggrpc.WithTransportCredentials(insecure.NewCredentials()),
			),
			grpc.WithConnectionPoolSize(uint64(runtime.GOMAXPROCS(0))),
		)
		if err != nil {
			b.Fatalf("failed to create client: %v", err)
		}
		defer client.Close(ctx)

		if _, err := client.Connect(ctx, addr); err != nil {
			b.Fatalf("failed to connect: %v", err)
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := client.Insert(ctx, req)
				if err != nil {
					b.Errorf("Insert failed: %v", err)
				}
			}
		})
	})
}
