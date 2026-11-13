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
package metric

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/net/grpc"
)

// BenchmarkMetricInterceptors_Unary measures the per-RPC cost of the server metric
// interceptor (latency histogram + completed-RPC counter) around an empty handler.
func BenchmarkMetricInterceptors_Unary(b *testing.B) {
	in, _, err := MetricInterceptors()
	if err != nil {
		b.Fatal(err)
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/vald.v1.Search/Search"}
	handler := func(context.Context, any) (any, error) { return nil, nil }
	ctx := context.Background()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := in(ctx, nil, info, handler); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMetricInterceptors_Unary_Parallel(b *testing.B) {
	in, _, err := MetricInterceptors()
	if err != nil {
		b.Fatal(err)
	}
	info := &grpc.UnaryServerInfo{FullMethod: "/vald.v1.Search/Search"}
	handler := func(context.Context, any) (any, error) { return nil, nil }
	ctx := context.Background()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := in(ctx, nil, info, handler); err != nil {
				b.Fatal(err)
			}
		}
	})
}
