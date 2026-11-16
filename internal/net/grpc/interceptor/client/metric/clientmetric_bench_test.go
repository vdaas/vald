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

	"google.golang.org/grpc"
)

// BenchmarkClientInterceptors_Unary measures the per-RPC cost of the client metric
// interceptor around an invoker that returns immediately.
func BenchmarkClientInterceptors_Unary(b *testing.B) {
	in, _, err := ClientInterceptors()
	if err != nil {
		b.Fatal(err)
	}
	invoker := func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error { return nil }
	ctx := context.Background()
	b.ReportAllocs()
	for b.Loop() {
		if err := in(ctx, "/vald.v1.Search/Search", nil, nil, nil, invoker); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkClientInterceptors_Unary_Parallel(b *testing.B) {
	in, _, err := ClientInterceptors()
	if err != nil {
		b.Fatal(err)
	}
	invoker := func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error { return nil }
	ctx := context.Background()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := in(ctx, "/vald.v1.Search/Search", nil, nil, nil, invoker); err != nil {
				b.Fatal(err)
			}
		}
	})
}
