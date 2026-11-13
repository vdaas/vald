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
package logging

import (
	"context"
	"os"
	"testing"

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
)

// The default logger level is DEBUG, which renders every success access log. Production runs
// at INFO, where the success entry is not rendered: run the benchmarks with
// BENCH_LOG_LEVEL=info to measure that case (BENCH_LOG_LEVEL=debug for the rendered case).
func init() {
	if lv := os.Getenv("BENCH_LOG_LEVEL"); lv != "" {
		log.Init(log.WithLevel(lv))
	}
}

var (
	benchAccessLogInfo    = &grpc.UnaryServerInfo{FullMethod: "/vald.v1.Search/Search"}
	benchAccessLogHandler = func(context.Context, any) (any, error) { return nil, nil }
)

// BenchmarkAccessLogInterceptor measures the per-RPC overhead the access log
// interceptor adds around a handler that does nothing (successful call).
func BenchmarkAccessLogInterceptor(b *testing.B) {
	in := AccessLogInterceptor()
	ctx := context.Background()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := in(ctx, nil, benchAccessLogInfo, benchAccessLogHandler); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAccessLogInterceptor_Parallel(b *testing.B) {
	in := AccessLogInterceptor()
	ctx := context.Background()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := in(ctx, nil, benchAccessLogInfo, benchAccessLogHandler); err != nil {
				b.Fatal(err)
			}
		}
	})
}
