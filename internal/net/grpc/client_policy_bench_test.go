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

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newPolicyBenchClient(b *testing.B, addr string) (*gRPCClient, fanOutConnPool) {
	b.Helper()
	log.Init(log.WithLoggerType("nop"))
	bo := backoff.New(
		backoff.WithInitialDuration("5ms"),
		backoff.WithBackOffFactor(1.1),
		backoff.WithJitterLimit("100ms"),
		backoff.WithMaximumDuration("5s"),
		backoff.WithBackOffTimeLimit("5s"),
		backoff.WithRetryCount(100),
		backoff.WithEnableErrorLog(),
	)
	cb, err := circuitbreaker.New(
		circuitbreaker.WithBreakerOpts(
			circuitbreaker.WithClosedErrorRate(0.7),
			circuitbreaker.WithHalfOpenErrorRate(0.5),
			circuitbreaker.WithMinSamples(1000),
			circuitbreaker.WithOpenTimeout("1s"),
			circuitbreaker.WithClosedRefreshTimeout("10s"),
		),
	)
	if err != nil {
		b.Fatal(err)
	}
	g, ok := New("policy-bench", WithBackoff(bo), WithCircuitBreaker(cb)).(*gRPCClient)
	if !ok {
		b.Fatal("New did not return *gRPCClient")
	}
	conn, err := ggrpc.NewClient("passthrough:///policy-bench",
		ggrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		if cerr := conn.Close(); cerr != nil {
			b.Error(cerr)
		}
		if cerr := g.Close(context.Background()); cerr != nil {
			b.Error(cerr)
		}
	})
	p := fanOutConnPool{conn: conn, healthy: true}
	g.storeConn(addr, p)
	return g, p
}

func BenchmarkDoVoid_HealthyPolicy(b *testing.B) {
	const addr = "10.0.0.1:8081"
	g, p := newPolicyBenchClient(b, addr)
	ctx := WithGRPCMethod(context.Background(), "Exists")
	noop := func(context.Context, string, *ClientConn, ...CallOption) error { return nil }

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := g.doVoid(ctx, p, addr, noop); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDoVoid_NotFoundPolicy(b *testing.B) {
	const addr = "10.0.0.1:8081"
	g, p := newPolicyBenchClient(b, addr)
	ctx := WithQuietErrorLog(WithGRPCMethod(context.Background(), "Exists"))
	notFound := status.Error(codes.NotFound, "not found")
	fn := func(context.Context, string, *ClientConn, ...CallOption) error { return notFound }

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := g.doVoid(ctx, p, addr, fn); err == nil {
			b.Fatal("expected error")
		}
	}
}

func BenchmarkHandleExecuteErr_NotFound(b *testing.B) {
	const addr = "10.0.0.1:8081"
	g, p := newPolicyBenchClient(b, addr)
	ctx := WithQuietErrorLog(context.Background())
	notFound := status.Error(codes.NotFound, "not found")

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_, _, err := g.handleExecuteErr(ctx, p, addr, notFound)
		if err == nil {
			b.Fatal("expected error")
		}
	}
}
