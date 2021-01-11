//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package pool provides grpc client connection pool
package pool

import (
	"context"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/discoverer"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/level"
	"github.com/vdaas/vald/internal/net"
	"google.golang.org/grpc"
)

const (
	DefaultServerAddr = "localhost:5001"
	DefaultPoolSize   = 10
)

type server struct {
	discoverer.DiscovererServer
}

func init() {
	testing.Init()
	log.Init(log.WithLevel(level.ERROR.String()))
}

func (s *server) Pods(context.Context, *payload.Discoverer_Request) (*payload.Info_Pods, error) {
	return &payload.Info_Pods{
		Pods: []*payload.Info_Pod{
			{
				Name: "vald is high scalable distributed high-speed approximate nearest neighbor search engine",
			},
		},
	}, nil
}

func (s *server) Nodes(context.Context, *payload.Discoverer_Request) (*payload.Info_Nodes, error) {
	return new(payload.Info_Nodes), nil
}

func ListenAndServe(b *testing.B, addr string) func() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		b.Error(err)
	}

	s := grpc.NewServer()
	discoverer.RegisterDiscovererServer(s, new(server))

	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		wg.Done()
		if err := s.Serve(lis); err != nil {
			b.Error(err)
		}
	}()

	wg.Wait()
	return func() {
		s.Stop()
	}
}

func do(b *testing.B, conn *ClientConn) {
	b.Helper()
	_, err := discoverer.NewDiscovererClient(conn).Nodes(context.Background(), new(payload.Discoverer_Request))
	if err != nil {
		b.Error(err)
	}
}

func Benchmark_ConnPool(b *testing.B) {
	defer ListenAndServe(b, DefaultServerAddr)()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool, err := New(ctx,
		WithAddr(DefaultServerAddr),
		WithSize(DefaultPoolSize),
		WithDialOptions(grpc.WithInsecure()),
	)
	if err != nil {
		b.Error(err)
	}
	pool, err = pool.Connect(ctx)
	if err != nil {
		b.Error(err)
	}

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		conn, ok := pool.Get()
		if ok {
			do(b, conn)
		}
	}
	b.StopTimer()
}

func Benchmark_StaticDial(b *testing.B) {
	defer ListenAndServe(b, DefaultServerAddr)()

	conn, err := grpc.DialContext(context.Background(), DefaultServerAddr, grpc.WithInsecure())
	if err != nil {
		b.Error(err)
	}

	conns := new(sync.Map)
	conns.Store(DefaultServerAddr, conn)

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		val, ok := conns.Load(DefaultServerAddr)
		if ok {
			do(b, val.(*ClientConn))
		}
	}
	b.StopTimer()
}

func BenchmarkParallel_ConnPool(b *testing.B) {
	defer ListenAndServe(b, DefaultServerAddr)()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pool, err := New(ctx,
		WithAddr(DefaultServerAddr),
		WithSize(DefaultPoolSize),
		WithDialOptions(grpc.WithInsecure()),
	)
	if err != nil {
		b.Error(err)
	}
	pool, err = pool.Connect(ctx)
	if err != nil {
		b.Error(err)
	}

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, ok := pool.Get()
			if ok {
				do(b, conn)
			}
		}
	})
	b.StopTimer()
}

func BenchmarkParallel_StaticDial(b *testing.B) {
	defer ListenAndServe(b, DefaultServerAddr)()

	conn, err := grpc.DialContext(context.Background(), DefaultServerAddr, grpc.WithInsecure())
	if err != nil {
		b.Error(err)
	}

	conns := new(sync.Map)
	conns.Store(DefaultServerAddr, conn)

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			val, ok := conns.Load(DefaultServerAddr)
			if ok {
				do(b, val.(*ClientConn))
			}
		}
	})
	b.StopTimer()
}
