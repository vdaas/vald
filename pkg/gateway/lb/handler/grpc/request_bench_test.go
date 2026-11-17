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
	"fmt"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/pkg/gateway/lb/service"
)

// The request benchmarks drive the LB handler end to end (validation, fan-out,
// aggregation, response) against in-memory agents and a gateway that calls them
// one after another, so they measure the handler's own cost per request and not
// the network or the scheduler.

const (
	benchAgents = 3
	benchDim    = 128
	benchTopK   = 10
)

type benchAgent struct {
	vald.Client
	name string
}

func (a *benchAgent) loc(uuid string) *payload.Object_Location {
	return &payload.Object_Location{Uuid: uuid, Name: a.name, Ips: []string{a.name}}
}

func (a *benchAgent) Insert(
	_ context.Context, in *payload.Insert_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return a.loc(in.GetVector().GetId()), nil
}

func (a *benchAgent) Update(
	_ context.Context, in *payload.Update_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return a.loc(in.GetVector().GetId()), nil
}

func (a *benchAgent) Remove(
	_ context.Context, in *payload.Remove_Request, _ ...grpc.CallOption,
) (*payload.Object_Location, error) {
	return a.loc(in.GetId().GetId()), nil
}

func (a *benchAgent) Exists(
	_ context.Context, in *payload.Object_ID, _ ...grpc.CallOption,
) (*payload.Object_ID, error) {
	return &payload.Object_ID{Id: in.GetId()}, nil
}

func (a *benchAgent) Search(
	_ context.Context, _ *payload.Search_Request, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	res := &payload.Search_Response{Results: make([]*payload.Object_Distance, benchTopK)}
	for i := range res.Results {
		res.Results[i] = &payload.Object_Distance{Id: fmt.Sprintf("%s-%d", a.name, i), Distance: float32(i) * 0.01}
	}
	return res, nil
}

type benchGateway struct {
	addrs  []string
	agents []vald.Client
}

func newBenchGateway() *benchGateway {
	g := new(benchGateway)
	for i := range benchAgents {
		name := fmt.Sprintf("10.0.0.%d:8081", i+1)
		g.addrs = append(g.addrs, name)
		g.agents = append(g.agents, &benchAgent{name: name})
	}
	return g
}

func (g *benchGateway) Start(context.Context) (<-chan error, error) { return nil, nil }
func (g *benchGateway) GetAgentCount(context.Context) int           { return len(g.agents) }
func (g *benchGateway) Addrs(context.Context) []string              { return g.addrs }

func (g *benchGateway) DoMulti(
	ctx context.Context,
	num int,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) error {
	done := 0
	for i, a := range g.agents {
		if done >= num {
			break
		}
		if err := f(ctx, g.addrs[i], a); err == nil {
			done++
		}
	}
	return nil
}

func (g *benchGateway) BroadCast(
	ctx context.Context,
	_ service.BroadCastKind,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) error {
	for i, a := range g.agents {
		if err := f(ctx, g.addrs[i], a); err != nil {
			return err
		}
	}
	return nil
}

func newBenchServer(b *testing.B) vald.Server {
	b.Helper()
	eg, _ := errgroup.New(b.Context())
	return New(
		WithName("bench"),
		WithIP("127.0.0.1"),
		WithGateway(newBenchGateway()),
		WithErrGroup(eg),
		WithReplicationCount(benchAgents),
	)
}

func benchVector() []float32 {
	v := make([]float32, benchDim)
	for i := range v {
		v[i] = float32(i) * 0.001
	}
	return v
}

func BenchmarkRequest_Insert(b *testing.B) {
	s := newBenchServer(b)
	req := &payload.Insert_Request{
		Vector: &payload.Object_Vector{Id: "uuid-0001", Vector: benchVector()},
		Config: &payload.Insert_Config{SkipStrictExistCheck: true},
	}
	ctx := b.Context()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := s.Insert(ctx, req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRequest_Update(b *testing.B) {
	s := newBenchServer(b)
	req := &payload.Update_Request{
		Vector: &payload.Object_Vector{Id: "uuid-0001", Vector: benchVector()},
		Config: &payload.Update_Config{SkipStrictExistCheck: true},
	}
	ctx := b.Context()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := s.Update(ctx, req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRequest_Remove(b *testing.B) {
	s := newBenchServer(b)
	req := &payload.Remove_Request{
		Id:     &payload.Object_ID{Id: "uuid-0001"},
		Config: &payload.Remove_Config{SkipStrictExistCheck: true},
	}
	ctx := b.Context()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := s.Remove(ctx, req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRequest_Search(b *testing.B) {
	s := newBenchServer(b)
	req := &payload.Search_Request{
		Vector: benchVector(),
		Config: &payload.Search_Config{Num: benchTopK, Radius: -1, Epsilon: 0.1, Timeout: 3_000_000_000},
	}
	ctx := b.Context()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := s.Search(ctx, req); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRequest_Search_Parallel(b *testing.B) {
	s := newBenchServer(b)
	req := &payload.Search_Request{
		Vector: benchVector(),
		Config: &payload.Search_Config{Num: benchTopK, Radius: -1, Epsilon: 0.1, Timeout: 3_000_000_000},
	}
	ctx := b.Context()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := s.Search(ctx, req); err != nil {
				b.Fatal(err)
			}
		}
	})
}
