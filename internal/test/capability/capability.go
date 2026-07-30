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

// Package capability is the capability layer over testing.TB. It owns the two
// abstractions every other layer builds on:
//
//   - Runner[X], the self-referential constraint expressing "a testing.TB
//     that can spawn subtests of its own type" (testing.TB deliberately
//     omits Run because T.Run and B.Run take callbacks of different types);
//   - As[C], the errors.As-style capability probe that resolves the
//     benchmark-only (b.Loop, b.ReportMetric, timer control, ...) and
//     test-only (t.Parallel) surfaces through any chain of wrappers.
//
// The named helpers (IsBenchmark, Loop, Measured, ReportMetric, ...) are
// thin conveniences over As: each performs exactly one capability check
// against the narrow interface it needs — never against the concrete
// *testing.B / *testing.T types, so wrappers embedding them keep working —
// and degrades to a documented fallback when the capability is absent.
// Their constraint is plain testing.TB (looser than Runner[X]) so both
// Runner-generic orchestration code and interface-typed leaf helpers can
// call them. New capability checks do not require touching this package:
// callers probe for their own narrow interface via As directly.
package capability

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
)

// Runner constrains the concrete testing entry types the test framework can
// drive. testing.TB deliberately omits Run (T.Run and B.Run take callbacks
// of their own concrete type, so no single method signature fits the
// interface), which is why the constraint is self-referential: X must both
// behave like testing.TB and spawn subtests of its own type. *testing.T and
// *testing.B satisfy it; *testing.F does not (it has Fuzz, not Run).
type Runner[X testing.TB] interface {
	testing.TB
	Run(name string, f func(X)) bool
}

// maxUnwrapDepth bounds the wrapper chain As is willing to follow. Wrappers
// holding closures are not comparable, so a same-value cycle check is not an
// option; the depth bound keeps a misbehaving self-returning Unwrap from
// spinning.
const maxUnwrapDepth = 8

// As reports whether t — or any testing.TB it wraps — satisfies the
// capability interface C, returning the first value in the chain that does.
// Wrappers participate by exposing Unwrap() testing.TB (the convention Node
// implements, mirroring errors.Unwrap). It is the single extension point of
// the capability layer: code needing a capability this package has no named
// helper for probes for its own narrow interface, e.g.
//
//	if b, ok := capability.As[interface{ Elapsed() time.Duration }](t); ok { ... }
func As[C any](tb testing.TB) (C, bool) {
	tb.Helper()
	for range maxUnwrapDepth {
		if c, ok := any(tb).(C); ok {
			return c, true
		}
		u, ok := tb.(interface{ Unwrap() testing.TB })
		if !ok {
			break
		}
		inner := u.Unwrap()
		if inner == nil {
			break
		}
		tb = inner
	}
	var zero C
	return zero, false
}

// IsBenchmark reports whether t is driven by the benchmark harness,
// detected through the Loop capability rather than the concrete type.
func IsBenchmark[X testing.TB](t X) bool {
	_, ok := As[interface{ Loop() bool }](t)
	return ok
}

// Loop executes body once per measured iteration when t exposes the
// benchmark Loop capability (b.Loop, which also confines the benchmark
// timer to the loop) and exactly once on any other testing.TB value. It is
// the unified "measured region" iteration primitive.
func Loop[X testing.TB](t X, body func()) {
	t.Helper()
	if l, ok := As[interface{ Loop() bool }](t); ok {
		for l.Loop() {
			body()
		}
		return
	}
	body()
}

// Measured runs fn as t's measured unit: each Loop iteration executes fn
// once with its own fresh timeout window when timeout > 0 (a single shared
// window would start expiring before the first iteration and starve later
// ones), and per-iteration errors are joined so an early failure is not
// masked by later successes. On non-benchmark Runners fn runs exactly once
// under the same per-run window semantics.
func Measured[X testing.TB](
	ctx context.Context, t X, timeout time.Duration, fn func(context.Context) error,
) (err error) {
	t.Helper()
	Loop(t, func() {
		ierr := func() error {
			ictx := ctx
			if timeout > 0 {
				var cancel context.CancelFunc
				ictx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}
			return fn(ictx)
		}()
		if ierr != nil {
			err = errors.Join(err, ierr)
		}
	})
	return err
}

// ReportMetric exposes value on t's benchmark result line (benchstat
// compatible); it is a no-op when t cannot report metrics.
func ReportMetric[X testing.TB](t X, value float64, unit string) {
	if r, ok := As[interface{ ReportMetric(float64, string) }](t); ok {
		r.ReportMetric(value, unit)
	}
}

// ReportAllocs enables allocation reporting when t supports it.
func ReportAllocs[X testing.TB](t X) {
	if r, ok := As[interface{ ReportAllocs() }](t); ok {
		r.ReportAllocs()
	}
}

// SetBytes records the number of bytes processed per iteration when t
// supports it.
func SetBytes[X testing.TB](t X, n int64) {
	if r, ok := As[interface{ SetBytes(int64) }](t); ok {
		r.SetBytes(n)
	}
}

// ResetTimer zeroes the benchmark timer when t supports it (no-op
// otherwise), so generic setup code can keep itself out of the measured
// window without knowing whether it runs under a test or a benchmark.
func ResetTimer[X testing.TB](t X) {
	if r, ok := As[interface{ ResetTimer() }](t); ok {
		r.ResetTimer()
	}
}

// StartTimer resumes the benchmark timer when t supports it (no-op
// otherwise); pair it with StopTimer around unmeasured teardown work.
func StartTimer[X testing.TB](t X) {
	if r, ok := As[interface{ StartTimer() }](t); ok {
		r.StartTimer()
	}
}

// StopTimer pauses the benchmark timer when t supports it (no-op
// otherwise), keeping generic teardown work out of the measured window.
func StopTimer[X testing.TB](t X) {
	if r, ok := As[interface{ StopTimer() }](t); ok {
		r.StopTimer()
	}
}

// Parallel signals that the test may run in parallel with (and only with)
// other parallel tests; benchmarks have no such phase, so it is a no-op
// there (b.RunParallel is a different, intra-benchmark concept).
func Parallel[X testing.TB](t X) {
	if p, ok := As[interface{ Parallel() }](t); ok {
		p.Parallel()
	}
}
