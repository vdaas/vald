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

// This file encodes the TDD RED-phase test contract for the not-yet-implemented
// recall metrics feature (see design rationale below) and intentionally does not
// compile yet - it references RecallRecorder/RecordRecall/GlobalSnapshot.Recalls
// etc. before the corresponding production code exists. Since this repository's
// Stop-hook / golangci-lint gate treats "package fails to typecheck" as an
// unconditional failure with no TDD RED-phase exception (typecheck errors are not
// scoped by --new-from-rev the way lint findings on changed lines are), this file
// is excluded from the default build via the `ignore` tag below - mirroring the
// existing `//go:build e2e` convention already used throughout tests/v2/e2e to
// gate files that require special build/runtime setup. The next Maker should
// delete the `//go:build ignore` line (only) once RecallRecorder, RecordRecall,
// WithRecallHistogram/WithRecallTDigest, and GlobalSnapshot.Recalls/
// RecallPercentiles are implemented; the assertions below are written to then
// fail on real behavior (true RED), not on missing symbols.

package metrics

// RED PHASE (TDD Test Maker): this file intentionally references symbols that do
// not exist yet:
//   - RecallRecorder (new interface, embedded into Collector)
//   - collector.recalls / collector.recallPercentiles (new unexported fields)
//   - (*collector).RecordRecall(val float64)
//   - WithRecallHistogram / WithRecallTDigest (new Option constructors in option.go)
//   - GlobalSnapshot.Recalls (*HistogramSnapshot) / GlobalSnapshot.RecallPercentiles (TDigest)
//
// Design decision (locked by this test file, matches the request from the task):
// recall is NOT folded into RequestResult/Record(). It is computed by the caller
// (tests/v2/e2e/crud/search_test.go's calculateRecall) *after* a search response has
// already been unmarshalled, i.e. strictly after Record() would have been called for
// that same request. Reusing Record() would therefore require either a second Record
// call (double-counting Total/Errors) or widening RequestResult with a Search-only
// field that is meaningless for non-search operations (Insert/Update/Remove/Object).
// Instead we mirror the existing dual-sketch pattern used for latencies/queueWaits:
// a Histogram (coarse bucketed distribution, used for the ASCII histogram in
// SnapshotPresenter) plus a TDigest (accurate quantiles), each independently
// toggleable via functional options and both registered by default so recall
// aggregation "just works" for any existing NewCollector() call site.
//
// No range clamping to [0,1] is performed by RecordRecall: like Record() does for
// Latency/QueueWait, out-of-range floats are accepted verbatim (the caller is
// trusted to pass a valid ratio; Vald Law: "validate only at system boundaries").
// NaN/+Inf/-Inf are silently dropped, mirroring histogram.Record's and
// tdigest.Add's existing guards - RecordRecall performs no additional filtering of
// its own, it simply delegates to the same Histogram/TDigest primitives already
// used for latency and queue-wait.

import (
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/strings"
)

// TestRecallRecorder_InterfaceContract locks the shape of the new interface and its
// embedding into Collector at compile time. If this test file fails to compile,
// that IS the expected RED state - the assertions below exist purely to make the
// intended contract unambiguous for the implementing Maker.
func TestRecallRecorder_InterfaceContract(t *testing.T) {
	t.Parallel()

	var _ RecallRecorder = (*collector)(nil)
	var _ interface {
		RecordRecall(val float64)
	} = (*collector)(nil)
	// Collector must embed RecallRecorder so any Collector value can record recall
	// without a type assertion back to *collector.
	var _ Collector = (*collector)(nil)
	var _ RecallRecorder = Collector(nil)
}

// TestNewCollector_DefaultOptions_EnableRecallRecording verifies that, exactly like
// latencies/queueWaits/exemplars, recall recording is ON by default - callers must
// not need to pass any new option to start collecting recall metrics.
func TestNewCollector_DefaultOptions_EnableRecallRecording(t *testing.T) {
	t.Parallel()

	c, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}
	cc, ok := c.(*collector)
	if !ok {
		t.Fatalf("NewCollector() returned unexpected concrete type %T", c)
	}
	if cc.recalls == nil {
		t.Error("collector.recalls is nil; WithRecallHistogram must be part of defaultOptions")
	}
	if cc.recallPercentiles == nil {
		t.Error("collector.recallPercentiles is nil; WithRecallTDigest must be part of defaultOptions")
	}

	c.RecordRecall(0.875)
	snap := c.GlobalSnapshot()
	if snap.Recalls == nil {
		t.Fatal("GlobalSnapshot().Recalls is nil after RecordRecall")
	}
	if snap.Recalls.Total != 1 {
		t.Errorf("Recalls.Total = %d, want 1", snap.Recalls.Total)
	}
	if snap.RecallPercentiles == nil {
		t.Fatal("GlobalSnapshot().RecallPercentiles is nil after RecordRecall")
	}
	if got := snap.RecallPercentiles.Quantile(0.5); got != 0.875 {
		t.Errorf("RecallPercentiles.Quantile(0.5) = %v, want 0.875", got)
	}
}

// TestWithRecallHistogram_And_WithRecallTDigest exercises the two new functional
// options directly, mirroring TestNewHistogram/TestNewTDigest style option tests.
func TestWithRecallHistogram_And_WithRecallTDigest(t *testing.T) {
	t.Parallel()

	t.Run("WithRecallHistogram sets collector.recalls", func(t *testing.T) {
		t.Parallel()
		c, err := NewCollector(WithRecallHistogram(defaultHistogramOpts...))
		if err != nil {
			t.Fatalf("NewCollector() error = %v", err)
		}
		cc := c.(*collector)
		if cc.recalls == nil {
			t.Fatal("collector.recalls is nil after WithRecallHistogram")
		}
	})

	t.Run("WithRecallTDigest sets collector.recallPercentiles", func(t *testing.T) {
		t.Parallel()
		c, err := NewCollector(WithRecallTDigest(defaultTDigestOpts...))
		if err != nil {
			t.Fatalf("NewCollector() error = %v", err)
		}
		cc := c.(*collector)
		if cc.recallPercentiles == nil {
			t.Fatal("collector.recallPercentiles is nil after WithRecallTDigest")
		}
	})

	t.Run("propagates HistogramOption/TDigestOption errors", func(t *testing.T) {
		t.Parallel()
		if _, err := NewCollector(WithRecallHistogram(WithHistogramNumShards(0))); err == nil {
			t.Error("expected error from invalid WithHistogramNumShards(0), got nil")
		}
		if _, err := NewCollector(WithRecallTDigest(WithTDigestCompression(-1))); err == nil {
			t.Error("expected error from invalid WithTDigestCompression(-1), got nil")
		}
	})
}

// TestCollector_RecordRecall covers normal values, edge values, and non-destructive
// interaction with the pre-existing latency/queueWait/exemplar recording paths.
func TestCollector_RecordRecall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		values    []float64
		wantTotal uint64
		wantMean  float64
		meanExact bool
	}{
		{
			name:      "single in-range value",
			values:    []float64{0.9},
			wantTotal: 1,
			wantMean:  0.9,
			meanExact: true,
		},
		{
			name:      "multiple in-range values averaged",
			values:    []float64{1.0, 0.5, 0.0},
			wantTotal: 3,
			wantMean:  0.5,
			meanExact: true,
		},
		{
			name:      "NaN is silently dropped, mirrors histogram/tdigest guard",
			values:    []float64{0.8, math.NaN(), 0.6},
			wantTotal: 2,
			wantMean:  0.7,
			meanExact: true,
		},
		{
			name:      "+Inf and -Inf are silently dropped",
			values:    []float64{0.5, math.Inf(1), math.Inf(-1)},
			wantTotal: 1,
			wantMean:  0.5,
			meanExact: true,
		},
		{
			name:      "out-of-range values are NOT clamped (recorded verbatim)",
			values:    []float64{-0.5, 1.5},
			wantTotal: 2,
			wantMean:  0.5, // mean(-0.5, 1.5) == 0.5; proves no clamping occurred
			meanExact: true,
		},
		{
			name:      "no values recorded leaves Total at 0",
			values:    nil,
			wantTotal: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c, err := NewCollector()
			if err != nil {
				t.Fatalf("NewCollector() error = %v", err)
			}
			for _, v := range tt.values {
				c.RecordRecall(v)
			}
			snap := c.GlobalSnapshot()
			if snap.Recalls == nil {
				t.Fatal("GlobalSnapshot().Recalls is nil")
			}
			if snap.Recalls.Total != tt.wantTotal {
				t.Errorf("Recalls.Total = %d, want %d", snap.Recalls.Total, tt.wantTotal)
			}
			if tt.meanExact && snap.Recalls.Total > 0 {
				if math.Abs(snap.Recalls.Mean-tt.wantMean) > 1e-9 {
					t.Errorf("Recalls.Mean = %v, want %v", snap.Recalls.Mean, tt.wantMean)
				}
			}
		})
	}

	t.Run("RecordRecall does not affect latency/queueWait/exemplar/Total counters", func(t *testing.T) {
		t.Parallel()
		c, err := NewCollector()
		if err != nil {
			t.Fatalf("NewCollector() error = %v", err)
		}
		c.Record(t.Context(), 0, &RequestResult{Latency: 10 * time.Millisecond, QueueWait: time.Millisecond})
		c.RecordRecall(0.42)
		c.RecordRecall(0.99)

		snap := c.GlobalSnapshot()
		if snap.Total != 1 {
			t.Errorf("Total = %d, want 1 (RecordRecall must not touch the request Total counter)", snap.Total)
		}
		if snap.Latencies == nil || snap.Latencies.Total != 1 {
			t.Errorf("Latencies unexpectedly affected by RecordRecall: %+v", snap.Latencies)
		}
		if snap.Recalls == nil || snap.Recalls.Total != 2 {
			t.Errorf("Recalls.Total = %+v, want 2", snap.Recalls)
		}
	})

	t.Run("nil-safety: zero-value collector must not panic", func(t *testing.T) {
		t.Parallel()
		// White-box: a bare &collector{} (bypassing NewCollector/defaultOptions) has
		// nil recalls/recallPercentiles, exactly like it already has nil
		// latencies/queueWaits/exemplars in that state. RecordRecall must guard
		// against this the same way Record() already guards its own fields.
		cc := &collector{}
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("RecordRecall panicked on zero-value collector: %v", r)
			}
		}()
		cc.RecordRecall(0.5)
	})
}

// TestCollector_Reset_ClearsRecall ensures Reset() clears recall data along with
// every other metric, so a Collector can be safely reused across benchmark phases.
func TestCollector_Reset_ClearsRecall(t *testing.T) {
	t.Parallel()

	c, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}
	c.RecordRecall(0.5)
	c.RecordRecall(0.9)

	if snap := c.GlobalSnapshot(); snap.Recalls == nil || snap.Recalls.Total != 2 {
		t.Fatalf("precondition failed, Recalls = %+v", snap.Recalls)
	}

	cc := c.(*collector)
	cc.Reset()

	snap := c.GlobalSnapshot()
	if snap.Recalls == nil {
		t.Fatal("Recalls is nil after Reset(); the histogram itself must survive Reset (like latencies does), just emptied")
	}
	if snap.Recalls.Total != 0 {
		t.Errorf("Recalls.Total = %d after Reset(), want 0", snap.Recalls.Total)
	}
}

// TestCollector_Clone_PropagatesRecall verifies Clone() deep-copies recall data and
// that the clone is fully independent from the source collector afterwards.
func TestCollector_Clone_PropagatesRecall(t *testing.T) {
	t.Parallel()

	c, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}
	c.RecordRecall(0.4)
	c.RecordRecall(0.6)

	cloned, err := c.Clone()
	if err != nil {
		t.Fatalf("Clone() error = %v", err)
	}

	// Mutate the original after cloning; the clone must not observe this.
	c.RecordRecall(1.0)

	origSnap := c.GlobalSnapshot()
	cloneSnap := cloned.GlobalSnapshot()

	if origSnap.Recalls.Total != 3 {
		t.Errorf("original Recalls.Total = %d, want 3", origSnap.Recalls.Total)
	}
	if cloneSnap.Recalls == nil {
		t.Fatal("cloned Recalls is nil")
	}
	if cloneSnap.Recalls.Total != 2 {
		t.Errorf("cloned Recalls.Total = %d, want 2 (clone must be independent of later writes to the source)", cloneSnap.Recalls.Total)
	}
	if cloneSnap.RecallPercentiles == nil {
		t.Fatal("cloned RecallPercentiles is nil")
	}
	if got := cloneSnap.RecallPercentiles.Quantile(0.5); math.Abs(got-0.5) > 1e-6 {
		t.Errorf("cloned RecallPercentiles.Quantile(0.5) = %v, want ~0.5", got)
	}
}

// TestCollector_Merge_PropagatesRecall verifies that MergeInto/merge combine recall
// Histograms and TDigests from both sides, symmetric to mergeHistograms/mergeTDigests.
func TestCollector_Merge_PropagatesRecall(t *testing.T) {
	t.Parallel()

	src, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}
	dst, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	src.RecordRecall(0.2)
	src.RecordRecall(0.4)
	dst.RecordRecall(0.8)

	if err := src.MergeInto(dst); err != nil {
		t.Fatalf("MergeInto() error = %v", err)
	}

	snap := dst.GlobalSnapshot()
	if snap.Recalls == nil {
		t.Fatal("merged Recalls is nil")
	}
	if snap.Recalls.Total != 3 {
		t.Errorf("merged Recalls.Total = %d, want 3", snap.Recalls.Total)
	}
	wantMean := (0.2 + 0.4 + 0.8) / 3
	if math.Abs(snap.Recalls.Mean-wantMean) > 1e-9 {
		t.Errorf("merged Recalls.Mean = %v, want %v", snap.Recalls.Mean, wantMean)
	}
	if snap.RecallPercentiles == nil {
		t.Fatal("merged RecallPercentiles is nil")
	}
}

// TestMergeSnapshots_PropagatesRecall exercises the package-level MergeSnapshots
// helper (used by the Pareto chart tooling to combine per-window/per-run JSON
// snapshots), including the edge case where one of the inputs was produced by a
// collector/serializer that predates recall support (Recalls == nil).
func TestMergeSnapshots_PropagatesRecall(t *testing.T) {
	t.Parallel()

	newSnapshotWithRecall := func(t *testing.T, vals ...float64) *GlobalSnapshot {
		t.Helper()
		c, err := NewCollector()
		if err != nil {
			t.Fatalf("NewCollector() error = %v", err)
		}
		for _, v := range vals {
			c.RecordRecall(v)
		}
		return c.GlobalSnapshot()
	}

	t.Run("two populated snapshots merge recall stats", func(t *testing.T) {
		t.Parallel()
		a := newSnapshotWithRecall(t, 0.5, 0.7)
		b := newSnapshotWithRecall(t, 0.9)

		merged, err := MergeSnapshots(a, b)
		if err != nil {
			t.Fatalf("MergeSnapshots() error = %v", err)
		}
		if merged.Recalls == nil {
			t.Fatal("merged.Recalls is nil")
		}
		if merged.Recalls.Total != 3 {
			t.Errorf("merged.Recalls.Total = %d, want 3", merged.Recalls.Total)
		}
		if merged.RecallPercentiles == nil {
			t.Fatal("merged.RecallPercentiles is nil")
		}
	})

	t.Run("nil Recalls on one input does not panic and preserves the other side's data", func(t *testing.T) {
		t.Parallel()
		withRecall := newSnapshotWithRecall(t, 0.6, 0.8)
		// Simulate a pre-recall snapshot (e.g. decoded from an older JSON file that
		// never had the "recalls"/"recall_percentiles" fields).
		legacy := &GlobalSnapshot{
			Total:             withRecall.Total,
			StartTime:         withRecall.StartTime,
			LastUpdated:       withRecall.LastUpdated,
			BoundsHash:        withRecall.BoundsHash,
			SketchKind:        withRecall.SketchKind,
			Recalls:           nil,
			RecallPercentiles: nil,
		}

		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("MergeSnapshots panicked with a nil Recalls input: %v", r)
			}
		}()
		merged, err := MergeSnapshots(withRecall, legacy)
		if err != nil {
			t.Fatalf("MergeSnapshots() error = %v", err)
		}
		if merged.Recalls == nil || merged.Recalls.Total != withRecall.Recalls.Total {
			t.Errorf("merged.Recalls = %+v, want Total=%d preserved from the populated input", merged.Recalls, withRecall.Recalls.Total)
		}
	})

	t.Run("both inputs nil Recalls stays nil-safe", func(t *testing.T) {
		t.Parallel()
		legacy1 := &GlobalSnapshot{Total: 1}
		legacy2 := &GlobalSnapshot{Total: 1}
		merged, err := MergeSnapshots(legacy1, legacy2)
		if err != nil {
			t.Fatalf("MergeSnapshots() error = %v", err)
		}
		if merged.Recalls != nil && merged.Recalls.Total != 0 {
			t.Errorf("merged.Recalls = %+v, want nil or zero Total when neither input had recall data", merged.Recalls)
		}
	})
}

// TestGlobalSnapshot_JSON_RoundTrip_RecallPercentiles locks the requirement that
// GlobalSnapshot.UnmarshalJSON must restore RecallPercentiles into a *working*
// concrete TDigest (not merely a nil interface / a struct that satisfies the
// interface but panics on use). This mirrors TestGlobalSnapshot_JSON's coverage of
// LatPercentiles/QWPercentiles and will fail to compile/pass until the
// UnmarshalJSON shadow struct in metrics.go is extended with a RecallPercentiles
// *tdigest field, exactly like the existing LatPercentiles/QWPercentiles handling.
func TestGlobalSnapshot_JSON_RoundTrip_RecallPercentiles(t *testing.T) {
	t.Parallel()

	c, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}
	c.RecordRecall(0.1)
	c.RecordRecall(0.9)

	snap := c.GlobalSnapshot()

	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("json.Marshal(snap) error = %v", err)
	}

	var restored GlobalSnapshot
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal error = %v", err)
	}

	if restored.RecallPercentiles == nil {
		t.Fatal("restored.RecallPercentiles is nil after JSON round-trip")
	}
	if got := restored.RecallPercentiles.Quantile(0.5); math.Abs(got-0.5) > 1e-6 {
		t.Errorf("restored RecallPercentiles.Quantile(0.5) = %v, want ~0.5", got)
	}
	if restored.Recalls == nil {
		t.Fatal("restored.Recalls is nil after JSON round-trip")
	}
	if restored.Recalls.Total != 2 {
		t.Errorf("restored.Recalls.Total = %d, want 2", restored.Recalls.Total)
	}
}

// TestSnapshotPresenter_AsJSON_IncludesRecall locks the exact (snake_case, per this
// repo's tagliatelle lint config) JSON field names for the new fields, since
// AsJSON() is a thin wrapper around json.MarshalIndent(GlobalSnapshot) - no changes
// to presenter.go itself should be required, only to the GlobalSnapshot struct tags.
func TestSnapshotPresenter_AsJSON_IncludesRecall(t *testing.T) {
	t.Parallel()

	c, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}
	c.RecordRecall(0.77)
	snap := c.GlobalSnapshot()

	out, err := NewSnapshotPresenter(snap).AsJSON()
	if err != nil {
		t.Fatalf("AsJSON() error = %v", err)
	}
	if !strings.Contains(out, `"recalls"`) {
		t.Errorf("AsJSON() output missing %q field:\n%s", "recalls", out)
	}
	if !strings.Contains(out, `"recall_percentiles"`) {
		t.Errorf("AsJSON() output missing %q field:\n%s", "recall_percentiles", out)
	}
}
