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

package capability

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
)

// TestAs exercises the capability probe directly: a direct hit on the
// concrete entry, a hit through a wrapper chain (Node), a miss returning
// the zero value, and a caller-defined capability interface this package
// has no named helper for.
func TestAs(t *testing.T) {
	if _, ok := As[interface{ Helper() }](t); !ok {
		t.Error("As must resolve a capability the entry itself satisfies")
	}
	if _, ok := As[interface{ Loop() bool }](t); ok {
		t.Error("As must miss capabilities *testing.T does not have")
	}

	// Through a wrapper: Deadline is *testing.T-only (not part of
	// testing.TB), so Node cannot satisfy this probe by method promotion —
	// As must follow Unwrap to the concrete entry to resolve it.
	n := NewNode(t)
	if _, ok := As[interface{ Deadline() (time.Time, bool) }](n); !ok {
		t.Error("As must follow Unwrap through Node to the concrete entry")
	}

	// Caller-defined capability without a named helper in this package —
	// the open/closed extension path.
	if run, ok := As[interface {
		Run(string, func(*testing.T)) bool
	}](n); !ok {
		t.Error("As must expose caller-defined capability interfaces")
	} else if !run.Run("sub", func(*testing.T) {}) {
		t.Error("capability resolved via As must be callable")
	}
}

// TestCapabilities_T exercises every capability helper with X = *testing.T:
// Loop runs the body exactly once, Measured applies its timeout window to
// the single run and returns fn's error as-is via errors.Join, and all
// benchmark-only controls are safe no-ops.
func TestCapabilities_T(t *testing.T) {
	if IsBenchmark(t) {
		t.Error("IsBenchmark(*testing.T) must be false")
	}

	var runs int
	Loop(t, func() { runs++ })
	if runs != 1 {
		t.Errorf("Loop on *testing.T must run the body exactly once, ran %d times", runs)
	}

	if err := Measured(t.Context(), t, 100*time.Millisecond, func(ctx context.Context) error {
		if _, ok := ctx.Deadline(); !ok {
			return errors.New("measured context must carry the per-run timeout deadline")
		}
		return nil
	}); err != nil {
		t.Errorf("Measured returned unexpected error: %v", err)
	}

	wantErr := errors.New("measured failure")
	if err := Measured(t.Context(), t, 0, func(context.Context) error {
		return wantErr
	}); !errors.Is(err, wantErr) {
		t.Errorf("Measured must surface fn's error, got: %v", err)
	}

	// Benchmark-only controls must be no-ops on *testing.T.
	ReportMetric(t, 1.0, "noop")
	ReportAllocs(t)
	SetBytes(t, 1)
	ResetTimer(t)
	StartTimer(t)
	StopTimer(t)

	t.Run("parallel capability", func(tt *testing.T) {
		Parallel(tt) // *testing.T supports it; single subtest, so harmless.
	})

	// The constraint is plain testing.TB, so an interface-typed value must
	// also instantiate the helpers (X = testing.TB) with the capability
	// detection working off the dynamic type — the shape interface-typed
	// leaf helpers such as tests/v2/e2e/crud's logRecallAndQPS rely on.
	var itb testing.TB = t
	if IsBenchmark(itb) {
		t.Error("IsBenchmark must inspect the dynamic type behind testing.TB")
	}
	runs = 0
	Loop(itb, func() { runs++ })
	if runs != 1 {
		t.Errorf("Loop with X = testing.TB must run the body exactly once, ran %d times", runs)
	}
	ReportMetric(itb, 1.0, "noop")
}

// TestCapabilities_B exercises the helpers with X = *testing.B through
// testing.Benchmark. b.Loop consumes the benchmark's whole iteration
// budget, so Loop and Measured (which iterates via Loop) each get their own
// benchmark invocation — calling Loop twice in one invocation would make
// the second loop exit immediately, which is exactly the hazard Measured's
// single-Loop design avoids.
func TestCapabilities_B(t *testing.T) {
	var loops int
	res := testing.Benchmark(func(b *testing.B) {
		b.Helper()
		if !IsBenchmark(b) {
			b.Error("IsBenchmark(*testing.B) must be true")
		}
		ReportAllocs(b)
		ResetTimer(b)
		loops = 0
		Loop(b, func() { loops++ })
		StopTimer(b)
		SetBytes(b, 1)
		ReportMetric(b, float64(loops), "loops")
		StartTimer(b)
		Parallel(b) // no-op: *testing.B has no Parallel phase.
	})
	if loops < 1 {
		t.Errorf("Loop on *testing.B must run the body at least once, result: %s", res.String())
	}

	var iterations, misses int
	res = testing.Benchmark(func(b *testing.B) {
		b.Helper()
		iterations, misses = 0, 0
		if err := Measured(b.Context(), b, time.Second, func(ctx context.Context) error {
			iterations++
			if _, ok := ctx.Deadline(); !ok {
				misses++
			}
			return nil
		}); err != nil {
			b.Errorf("Measured returned unexpected error: %v", err)
		}
	})
	if iterations < 1 {
		t.Errorf("Measured on *testing.B must run fn at least once, result: %s", res.String())
	}
	if misses != 0 {
		t.Errorf("Measured must give every iteration a deadline-carrying context, %d/%d missed", misses, iterations)
	}
}

// loopWrapper wraps a testing.TB and provides its own Loop implementation,
// pinning As's shallowest-match-wins semantics: the probe asserts at every
// level of the chain (like errors.As) instead of blindly unwrapping to the
// terminal entry first, so a wrapper's own capability takes precedence over
// whatever it wraps.
type loopWrapper struct {
	testing.TB
	loops int
}

func (w *loopWrapper) Unwrap() testing.TB { return w.TB }

func (w *loopWrapper) Loop() bool {
	w.loops++
	return w.loops <= 1
}

func TestAsShallowestMatchWins(t *testing.T) {
	w := &loopWrapper{TB: t}
	got, ok := As[interface{ Loop() bool }](w)
	if !ok {
		t.Fatal("As must resolve the wrapper's own capability")
	}
	if got != any(w) {
		t.Errorf("As must return the shallowest match (the wrapper itself), got %T", got)
	}
	// The named helpers inherit the same semantics: the wrapper's Loop
	// drives the iteration even though the wrapped *testing.T has none.
	var runs int
	Loop(w, func() { runs++ })
	if runs != 1 {
		t.Errorf("Loop must consume the wrapper's Loop budget exactly once, ran %d times", runs)
	}
	if !IsBenchmark(w) {
		t.Error("IsBenchmark must report true for a wrapper providing Loop itself")
	}
}

// selfWrapper is a pathological testing.TB wrapper whose Unwrap returns
// itself, exercising As's depth bound: the probe must terminate and report
// a miss instead of spinning (wrappers holding closures are not comparable,
// so a same-value cycle check is not an option).
type selfWrapper struct {
	testing.TB
}

func (w *selfWrapper) Unwrap() testing.TB { return w }

func TestAsDepthBound(t *testing.T) {
	w := &selfWrapper{TB: t}
	if _, ok := As[interface{ Loop() bool }](w); ok {
		t.Error("As through a self-returning wrapper must terminate with a miss")
	}
	if IsBenchmark(w) {
		t.Error("IsBenchmark must stay false for a wrapper that never resolves to a benchmark")
	}
	// Loop must still fall back to running the body exactly once.
	var runs int
	Loop(w, func() { runs++ })
	if runs != 1 {
		t.Errorf("Loop through an unresolvable wrapper must run once, ran %d times", runs)
	}
}

// TestNode verifies the type-erasure contract: NewNode captures the
// concrete Runner type exactly once, Run keeps spawning correctly-typed
// children arbitrarily deep, a Node is a valid testing.TB, and every
// capability helper resolves the underlying entry through Unwrap — so
// IsBenchmark/Loop/ReportMetric behave identically whether they receive
// the raw *testing.B or a Node wrapping it.
func TestNode(t *testing.T) {
	root := NewNode(t)
	if IsBenchmark(root) {
		t.Error("IsBenchmark(Node{*testing.T}) must be false via Unwrap")
	}

	var depth2 bool
	root.Run("child", func(child Node) {
		child.Helper() // promoted testing.TB method
		var runs int
		Loop(child, func() { runs++ })
		if runs != 1 {
			child.Errorf("Loop through a Node over *testing.T must run once, ran %d times", runs)
		}
		child.Run("grandchild", func(gc Node) {
			depth2 = true
			var itb testing.TB = gc // Node satisfies testing.TB
			if itb.Name() == "" {
				gc.Error("promoted Name must identify the subtest")
			}
		})
	})
	if !depth2 {
		t.Error("nested Node.Run must execute the grandchild")
	}

	var loops int
	benchNode := false
	res := testing.Benchmark(func(b *testing.B) {
		b.Helper()
		n := NewNode(b)
		benchNode = IsBenchmark(n)
		n.Run("measured", func(child Node) {
			loops = 0
			Loop(child, func() { loops++ })
			ReportMetric(child, float64(loops), "loops")
		})
	})
	if !benchNode {
		t.Error("IsBenchmark(Node{*testing.B}) must be true via Unwrap")
	}
	if loops < 1 {
		t.Errorf("Loop through a Node over *testing.B must iterate, result: %s", res.String())
	}
}
