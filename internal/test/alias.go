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

package test

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/test/capability"
	"github.com/vdaas/vald/internal/test/table"
	"github.com/vdaas/vald/internal/test/testdata"
)

// This file is the compatibility facade of the layered test framework: every
// historical internal/test name is preserved as a generic type alias or a
// thin forwarding function, so the packages below stay free to have narrow,
// single-responsibility APIs while existing call sites keep compiling
// unchanged. New code may import the layer packages directly.

// Layer re-exports (generic type aliases, Go 1.24+).
type (
	// Runner constrains the concrete testing entry types the framework can
	// drive; see capability.Runner.
	Runner[X testing.TB] = capability.Runner[X]
	// Node is the type-erased face of Runner[X]; see capability.Node.
	Node = capability.Node

	// CaseFor is one row of a table-driven test or benchmark; see
	// table.CaseFor.
	CaseFor[X Runner[X], T, A any] = table.CaseFor[X, T, A]
	// Result carries either a success value or an error; see table.Result.
	Result[T any] = table.Result[T]
	// BeforeFuncFor prepares a case's Args before do runs.
	BeforeFuncFor[X Runner[X], A any] = table.BeforeFuncFor[X, A]
	// AfterFuncFor observes the case outcome for cleanup or extra assertions.
	AfterFuncFor[X Runner[X], T, A any] = table.AfterFuncFor[X, T, A]
	// CheckFuncFor compares the wanted and actual Result of a case.
	CheckFuncFor[X Runner[X], T any] = table.CheckFuncFor[X, T]
	// DoFor is the function under test, executed once per case.
	DoFor[X Runner[X], T, A any] = table.DoFor[X, T, A]
)

// The historical *testing.T-based names are kept as generic type aliases so
// existing call sites compile unchanged, alongside the *testing.B
// instantiations for table-driven benchmarks.
type (
	// Case is CaseFor instantiated for plain tests (X = *testing.T).
	Case[T, A any] = table.CaseFor[*testing.T, T, A]
	// BeforeFunc is BeforeFuncFor instantiated for plain tests.
	BeforeFunc[A any] = table.BeforeFuncFor[*testing.T, A]
	// AfterFunc is AfterFuncFor instantiated for plain tests.
	AfterFunc[T, A any] = table.AfterFuncFor[*testing.T, T, A]
	// CheckFunc is CheckFuncFor instantiated for plain tests.
	CheckFunc[T any] = table.CheckFuncFor[*testing.T, T]
	// Do is DoFor instantiated for plain tests.
	Do[T, A any] = table.DoFor[*testing.T, T, A]

	// BenchmarkCase is CaseFor instantiated for benchmarks (X = *testing.B).
	BenchmarkCase[T, A any] = table.CaseFor[*testing.B, T, A]
	// BenchmarkBeforeFunc is BeforeFuncFor instantiated for benchmarks.
	BenchmarkBeforeFunc[A any] = table.BeforeFuncFor[*testing.B, A]
	// BenchmarkAfterFunc is AfterFuncFor instantiated for benchmarks.
	BenchmarkAfterFunc[T, A any] = table.AfterFuncFor[*testing.B, T, A]
	// BenchmarkCheckFunc is CheckFuncFor instantiated for benchmarks.
	BenchmarkCheckFunc[T any] = table.CheckFuncFor[*testing.B, T]
	// BenchmarkDo is DoFor instantiated for benchmarks.
	BenchmarkDo[T, A any] = table.DoFor[*testing.B, T, A]
)

// ValidIndex is the fixture path of a valid NGT agent index; see
// testdata.ValidIndex.
const ValidIndex = testdata.ValidIndex

// GetTestdataPath returns the test data file path under
// `internal/test/data`; see testdata.GetTestdataPath.
func GetTestdataPath(filename string) string {
	return testdata.GetTestdataPath(filename)
}

// Run executes each case as a subtest (or sub-benchmark) of t; see
// table.Run.
func Run[X Runner[X], T, A any](
	ctx context.Context, t X, do DoFor[X, T, A], tests ...CaseFor[X, T, A],
) error {
	t.Helper()
	return table.Run(ctx, t, do, tests...)
}

// DefaultCheck is the CheckFunc used when a case does not provide one; see
// table.DefaultCheck.
func DefaultCheck[X Runner[X], T any](tt X, want, got Result[T]) error {
	tt.Helper()
	return table.DefaultCheck(tt, want, got)
}

// NewNode wraps t into a Node; see capability.NewNode.
func NewNode[X Runner[X]](t X) Node {
	return capability.NewNode(t)
}

// As reports whether tb — or any testing.TB it wraps — satisfies the
// capability interface C; see capability.As.
func As[C any](tb testing.TB) (C, bool) {
	tb.Helper()
	return capability.As[C](tb)
}

// IsBenchmark reports whether t is driven by the benchmark harness; see
// capability.IsBenchmark.
func IsBenchmark[X testing.TB](t X) bool {
	return capability.IsBenchmark(t)
}

// Loop executes body once per measured iteration on benchmarks and exactly
// once otherwise; see capability.Loop.
func Loop[X testing.TB](t X, body func()) {
	t.Helper()
	capability.Loop(t, body)
}

// Measured runs fn as t's measured unit with a fresh per-iteration timeout
// window; see capability.Measured.
func Measured[X testing.TB](
	ctx context.Context, t X, timeout time.Duration, fn func(context.Context) error,
) error {
	t.Helper()
	return capability.Measured(ctx, t, timeout, fn)
}

// ReportMetric exposes value on t's benchmark result line when supported;
// see capability.ReportMetric.
func ReportMetric[X testing.TB](t X, value float64, unit string) {
	capability.ReportMetric(t, value, unit)
}

// ReportAllocs enables allocation reporting when t supports it; see
// capability.ReportAllocs.
func ReportAllocs[X testing.TB](t X) {
	capability.ReportAllocs(t)
}

// SetBytes records the number of bytes processed per iteration when t
// supports it; see capability.SetBytes.
func SetBytes[X testing.TB](t X, n int64) {
	capability.SetBytes(t, n)
}

// ResetTimer zeroes the benchmark timer when t supports it; see
// capability.ResetTimer.
func ResetTimer[X testing.TB](t X) {
	capability.ResetTimer(t)
}

// StartTimer resumes the benchmark timer when t supports it; see
// capability.StartTimer.
func StartTimer[X testing.TB](t X) {
	capability.StartTimer(t)
}

// StopTimer pauses the benchmark timer when t supports it; see
// capability.StopTimer.
func StopTimer[X testing.TB](t X) {
	capability.StopTimer(t)
}

// Parallel signals that the test may run in parallel when t supports it;
// see capability.Parallel.
func Parallel[X testing.TB](t X) {
	capability.Parallel(t)
}
