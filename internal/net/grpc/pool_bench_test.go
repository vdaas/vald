//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
package grpc

import (
	"context"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/discoverer"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/net"
	"google.golang.org/grpc"
)

const (
	DefaultServerAddr = ":5001"
	DefaultPoolSize   = 100
)

type server struct {
	discoverer.DiscovererServer
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

	pool, err := NewPool(context.Background(), DefaultServerAddr, DefaultPoolSize, grpc.WithInsecure())
	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		conn, shared := pool.Get()
		do(b, conn)
		if !shared {
			pool.Put(conn)
		}
	}
}

func Benchmark_StaticDial(b *testing.B) {
	defer ListenAndServe(b, DefaultServerAddr)()

	conn, err := grpc.DialContext(context.Background(), DefaultServerAddr, grpc.WithInsecure())
	if err != nil {
		b.Error(err)
	}

	conns := new(sync.Map)
	conns.Store(DefaultServerAddr, conn)

	for i := 0; i < b.N; i++ {
		val, _ := conns.Load(DefaultServerAddr)
		do(b, val.(*ClientConn))
	}
}

func BenchmarkParallel_ConnPool(b *testing.B) {
	defer ListenAndServe(b, DefaultServerAddr)()

	pool, err := NewPool(context.Background(), DefaultServerAddr, DefaultPoolSize, grpc.WithInsecure())
	if err != nil {
		b.Error(err)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			conn, shared := pool.Get()
			do(b, conn)
			if !shared {
				pool.Put(conn)
			}
		}
	})
}

func BenchmarkParallel_StaticDial(b *testing.B) {
	defer ListenAndServe(b, DefaultServerAddr)()

	conn, err := grpc.DialContext(context.Background(), DefaultServerAddr, grpc.WithInsecure())
	if err != nil {
		b.Error(err)
	}

	conns := new(sync.Map)
	conns.Store(DefaultServerAddr, conn)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			val, _ := conns.Load(DefaultServerAddr)
			do(b, val.(*ClientConn))
		}
	})
}
