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
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// healthyConnPool is a non-IP, always-healthy fixedConnPool, so executeRPC
// treats a retryable status as worth retrying.
type healthyConnPool struct {
	fixedConnPool
}

func (healthyConnPool) IsIPConn() bool                 { return false }
func (healthyConnPool) IsHealthy(context.Context) bool { return true }

// BenchmarkDo_SearchBroadcastResourceExhausted measures one search-broadcast
// call (the ctx aggregationSearch builds) to an agent that keeps answering
// ResourceExhausted, under the charts/vald/values.yaml default client backoff
// (5ms initial, x1.1, 100ms jitter, 5s cap, 100 retries) and a 3s deadline.
func BenchmarkDo_SearchBroadcastResourceExhausted(b *testing.B) {
	log.Init(log.WithLoggerType("nop"))
	const addr = "10.0.0.1:8081"
	g, ok := New("bench", WithBackoff(backoff.New(
		backoff.WithInitialDuration("5ms"),
		backoff.WithBackOffFactor(1.1),
		backoff.WithJitterLimit("100ms"),
		backoff.WithMaximumDuration("5s"),
		backoff.WithBackOffTimeLimit("5s"),
		backoff.WithRetryCount(100),
	))).(*gRPCClient)
	if !ok {
		b.Fatal("New did not return *gRPCClient")
	}
	conn, err := ggrpc.NewClient("passthrough:///retry-bench",
		ggrpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		if cerr := conn.Close(); cerr != nil {
			b.Error(cerr)
		}
	})
	p := healthyConnPool{fixedConnPool{conn: conn}}
	exhausted := status.Error(codes.ResourceExhausted, "agent overloaded")
	f := func(context.Context, *ClientConn, ...CallOption) (any, error) {
		return nil, exhausted
	}
	b.ResetTimer()
	for range b.N {
		ctx, cancel := context.WithTimeout(b.Context(), 3*time.Second)
		ctx = WithQuietErrorLog(WithoutResourceExhaustedRetry(ctx))
		if _, err := g.do(ctx, p, addr, f); err == nil {
			b.Fatal("do() must fail")
		}
		cancel()
	}
}
