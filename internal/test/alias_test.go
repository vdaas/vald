//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package test_test

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test"
)

// TestFacade_Run exercises the compatibility facade exactly as historical
// call sites do: test.Run over test.Case rows with the historical hook and
// check aliases, including the Want mismatch path via DefaultCheck.
func TestFacade_Run(t *testing.T) {
	wantErr := errors.New("expected failure")
	var before test.BeforeFunc[int] = func(_ context.Context, t *testing.T, in int) int {
		t.Helper()
		return in + 1
	}
	var check test.CheckFunc[int] = func(t *testing.T, want, got test.Result[int]) error {
		t.Helper()
		return test.DefaultCheck(t, want, got)
	}
	var after test.AfterFunc[int, int] = func(_ context.Context, t *testing.T, _, _ int, err error) error {
		t.Helper()
		return err
	}
	var do test.Do[int, int] = func(t *testing.T, in int) (int, error) {
		t.Helper()
		if in < 0 {
			return 0, wantErr
		}
		return in * 2, nil
	}
	if err := test.Run(t.Context(), t, do, []test.Case[int, int]{
		{Name: "default check", Args: 20, Want: test.Result[int]{Val: 42}, BeforeFunc: before},
		{Name: "expected error", Args: -2, Want: test.Result[int]{Err: wantErr}},
		{Name: "explicit hooks", Args: 2, Want: test.Result[int]{Val: 6}, BeforeFunc: before, CheckFunc: check, AfterFunc: after},
	}...); err != nil {
		t.Errorf("facade Run returned unexpected error: %v", err)
	}
}

// TestFacade_Benchmark drives the *testing.B instantiations and every
// benchmark-capability wrapper through one testing.Benchmark invocation,
// plus Measured through its own (b.Loop's budget is consumed per
// invocation).
func TestFacade_Benchmark(t *testing.T) {
	var loops int
	res := testing.Benchmark(func(b *testing.B) {
		b.Helper()
		if !test.IsBenchmark(b) {
			b.Error("IsBenchmark(*testing.B) must be true through the facade")
		}
		var do test.BenchmarkDo[int, int] = func(b *testing.B, in int) (int, error) {
			b.Helper()
			test.ReportAllocs(b)
			test.ResetTimer(b)
			loops = 0
			test.Loop(b, func() { loops++ })
			test.StopTimer(b)
			test.SetBytes(b, 1)
			test.ReportMetric(b, float64(loops), "loops")
			test.StartTimer(b)
			return in, nil
		}
		var check test.BenchmarkCheckFunc[int] = func(b *testing.B, want, got test.Result[int]) error {
			b.Helper()
			return test.DefaultCheck(b, want, got)
		}
		var before test.BenchmarkBeforeFunc[int] = func(_ context.Context, b *testing.B, in int) int {
			b.Helper()
			return in
		}
		var after test.BenchmarkAfterFunc[int, int] = func(_ context.Context, b *testing.B, _, _ int, err error) error {
			b.Helper()
			return err
		}
		if err := test.Run(b.Context(), b, do, []test.BenchmarkCase[int, int]{
			{Name: "loop", Args: 1, Want: test.Result[int]{Val: 1}, BeforeFunc: before, CheckFunc: check, AfterFunc: after},
		}...); err != nil {
			b.Errorf("facade Run over *testing.B returned unexpected error: %v", err)
		}
	})
	if loops < 1 {
		t.Errorf("facade Loop must iterate under testing.Benchmark, result: %s", res.String())
	}

	testing.Benchmark(func(b *testing.B) {
		b.Helper()
		if err := test.Measured(b.Context(), b, time.Second, func(ctx context.Context) error {
			if _, ok := ctx.Deadline(); !ok {
				return errors.New("facade Measured must hand fn a deadline-carrying context")
			}
			return nil
		}); err != nil {
			b.Errorf("facade Measured returned unexpected error: %v", err)
		}
	})
}

// TestFacade_NodeAndCapabilities covers Node construction and the
// remaining wrappers with X = *testing.T (no-op paths), plus the As
// extension probe and the testdata helpers.
func TestFacade_NodeAndCapabilities(t *testing.T) {
	n := test.NewNode(t)
	if test.IsBenchmark(n) {
		t.Error("IsBenchmark(Node{*testing.T}) must be false through the facade")
	}
	ran := false
	n.Run("child", func(child test.Node) {
		child.Helper()
		ran = true
	})
	if !ran {
		t.Error("facade Node.Run must execute the child")
	}
	// A parallel child is parked until the parent body returns, so it gets
	// its own subtree without a completion assertion.
	t.Run("parallel capability", func(tt *testing.T) {
		test.Parallel(tt)
	})

	// Deadline is *testing.T-only (not part of testing.TB), so resolving it
	// through a Node requires the facade's As to follow Unwrap.
	if _, ok := test.As[interface{ Deadline() (time.Time, bool) }](n); !ok {
		t.Error("facade As must resolve capabilities through Node's Unwrap")
	}
	if _, ok := test.As[interface{ Loop() bool }](t); ok {
		t.Error("facade As must miss benchmark capabilities on *testing.T")
	}

	if got, want := test.GetTestdataPath("tls/ca.pem"), "/internal/test/data/tls/ca.pem"; len(got) <= len(want) {
		t.Errorf("facade GetTestdataPath returned unexpectedly short path: %q", got)
	}
	if test.ValidIndex == "" {
		t.Error("facade ValidIndex must re-export the fixture path")
	}
}
