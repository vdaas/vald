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
package pool

import (
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Benchmark_GetHealthyConn(b *testing.B) {
	addr := "localhost:5003"
	defer ListenAndServe(b, addr)()

	ctx := b.Context()
	p, err := New(ctx,
		WithAddr(addr),
		WithSize(10),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		b.Fatal(err)
	}
	p, err = p.Connect(ctx)
	if err != nil {
		b.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond) // Wait for connections to be ready

	po, ok := p.(*pool)
	if !ok {
		b.Fatal("not a *pool")
	}

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()

	for b.Loop() {
		po.getHealthyConn(ctx)
	}
}

func BenchmarkParallel_GetHealthyConn(b *testing.B) {
	addr := "localhost:5004"
	defer ListenAndServe(b, addr)()

	ctx := b.Context()
	p, err := New(ctx,
		WithAddr(addr),
		WithSize(10),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		b.Fatal(err)
	}
	p, err = p.Connect(ctx)
	if err != nil {
		b.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond) // Wait for connections to be ready

	po, ok := p.(*pool)
	if !ok {
		b.Fatal("not a *pool")
	}

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			po.getHealthyConn(ctx)
		}
	})
}
