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

package circuitbreaker

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

func newBenchBreakerManager(b *testing.B) CircuitBreaker {
	b.Helper()
	log.Init(log.WithLoggerType("nop"))
	cb, err := New(
		WithBreakerOpts(
			WithClosedErrorRate(0.7),
			WithHalfOpenErrorRate(0.5),
			WithMinSamples(1000),
			WithOpenTimeout("1s"),
			WithClosedRefreshTimeout("10s"),
		),
	)
	if err != nil {
		b.Fatal(err)
	}
	return cb
}

func BenchmarkBreakerManager_Do_Healthy(b *testing.B) {
	cb := newBenchBreakerManager(b)
	ctx := context.Background()
	const key = "Exists/10.0.0.1:8081"
	fn := func(context.Context) (any, error) { return struct{}{}, nil }

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := cb.Do(ctx, key, fn); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBreakerManager_Do_IgnorableNotFound(b *testing.B) {
	cb := newBenchBreakerManager(b)
	ctx := context.Background()
	const key = "Exists/10.0.0.1:8081"
	notFound := errors.New("rpc error: code = NotFound desc = not found")
	fn := func(context.Context) (any, error) { return nil, errors.NewErrCircuitBreakerIgnorable(notFound) }

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := cb.Do(ctx, key, fn); err == nil {
			b.Fatal("expected error")
		}
	}
}

func BenchmarkBreakerManager_Handle_Healthy(b *testing.B) {
	cb := newBenchBreakerManager(b)
	bm := cb.(*breakerManager)
	const key = "Exists/10.0.0.1:8081"
	h, err := bm.Resolve(key)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := h.Allow(); err != nil {
			b.Fatal(err)
		}
		h.FinishSuccess()
	}
}

func BenchmarkBreakerManager_Handle_NonRetryableNotFound(b *testing.B) {
	cb := newBenchBreakerManager(b)
	bm := cb.(*breakerManager)
	const key = "Exists/10.0.0.1:8081"
	h, err := bm.Resolve(key)
	if err != nil {
		b.Fatal(err)
	}
	notFound := errors.New("rpc error: code = NotFound desc = not found")

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := h.Allow(); err != nil {
			b.Fatal(err)
		}
		if err := h.FinishError(notFound, false); err == nil {
			b.Fatal("expected error")
		}
	}
}

func BenchmarkBreaker_currentState_Closed(b *testing.B) {
	br, err := newBreaker("bench-closed", WithClosedRefreshTimeout("10s"))
	if err != nil {
		b.Fatal(err)
	}
	_ = br.currentState()

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if st := br.currentState(); st != StateClosed {
			b.Fatalf("state = %s, want Closed", st)
		}
	}
}
