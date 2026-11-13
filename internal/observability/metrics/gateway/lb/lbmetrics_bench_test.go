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
package lb

import (
	"testing"

	"github.com/vdaas/vald/internal/net/grpc/codes"
)

// The record methods run once per failed agent call while an agent is unhealthy,
// i.e. on every request that reaches it, from all request goroutines at once.

func BenchmarkSearchMetrics_RecordAgentFailure(b *testing.B) {
	sm := New()
	b.ReportAllocs()
	for b.Loop() {
		sm.RecordAgentFailure("Search", codes.Unavailable)
	}
}

func BenchmarkSearchMetrics_RecordAgentFailure_Parallel(b *testing.B) {
	sm := New()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sm.RecordAgentFailure("Search", codes.Unavailable)
		}
	})
}

func BenchmarkSearchMetrics_RecordWritePartialSuccess_Parallel(b *testing.B) {
	sm := New()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sm.RecordWritePartialSuccess("Insert")
		}
	})
}
