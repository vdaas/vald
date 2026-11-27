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

package metrics

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
)

const (
	testCounterRecallSum     = "recall_sum_ppm"
	testCounterRecallSamples = "recall_samples"
)

func buildTestCollectorSnapshot(t *testing.T, count int, custom map[string]int64) *GlobalSnapshot {
	t.Helper()
	opts := make([]Option, 0, len(custom))
	if len(custom) > 0 {
		names := make([]string, 0, len(custom))
		for k := range custom {
			names = append(names, k)
		}
		opts = append(opts, WithCustomCounters(names...))
	}
	col, err := NewCollector(opts...)
	if err != nil {
		t.Fatalf("NewCollector failed: %v", err)
	}
	base := time.Date(2026, 11, 24, 21, 0, 0, 0, time.UTC)
	for i := range count {
		lat := time.Duration(i+1) * time.Millisecond
		col.Record(t.Context(), uint64(i+1), &RequestResult{
			StartedAt: base.Add(time.Duration(i) * 100 * time.Millisecond),
			EndedAt:   base.Add(time.Duration(i)*100*time.Millisecond + lat),
			Latency:   lat,
			Status:    codes.OK,
		})
	}
	for k, v := range custom {
		col.IncCounter(k, v)
	}
	return SnapshotFromCollector(col)
}

func TestGlobalSnapshotHelpers(t *testing.T) {
	t.Parallel()

	t.Run("nil and zero-request snapshots", func(t *testing.T) {
		t.Parallel()
		var nilSnap *GlobalSnapshot
		if nilSnap.Duration() != 0 || nilSnap.QPS() != 0 || nilSnap.LatencyMinNS() != 0 ||
			nilSnap.LatencyMeanNS() != 0 || nilSnap.LatencyMaxNS() != 0 ||
			nilSnap.LatencyQuantileNS(0.50) != 0 || nilSnap.LatencyP99NS() != 0 {
			t.Fatalf("expected all zero values for nil GlobalSnapshot")
		}

		zeroSnap := buildTestCollectorSnapshot(t, 0, nil)
		if zeroSnap.Total != 0 || zeroSnap.Duration() != 0 || zeroSnap.QPS() != 0 ||
			zeroSnap.LatencyMinNS() != 0 || zeroSnap.LatencyMaxNS() != 0 || zeroSnap.LatencyP99NS() != 0 {
			t.Fatalf("expected zero metrics for Total==0 snapshot, got %+v", zeroSnap)
		}
		if SnapshotFromCollector(nil) == nil {
			t.Fatalf("expected non-nil empty GlobalSnapshot from nil Collector")
		}
	})

	t.Run("Total==1 and Total>1 with CustomCounters", func(t *testing.T) {
		t.Parallel()
		oneSnap := buildTestCollectorSnapshot(t, 1, nil)
		if oneSnap.Total != 1 || oneSnap.Duration() != time.Millisecond || oneSnap.QPS() != 1000 {
			t.Fatalf("unexpected Total==1 stats: total=%d dur=%v qps=%v", oneSnap.Total, oneSnap.Duration(), oneSnap.QPS())
		}
		if oneSnap.LatencyP99NS() != oneSnap.LatencyMeanNS() || oneSnap.LatencyQuantileNS(0.50) != oneSnap.LatencyMeanNS() {
			t.Fatalf("expected Total==1 percentiles to equal mean (%v), got p50=%v p99=%v",
				oneSnap.LatencyMeanNS(), oneSnap.LatencyQuantileNS(0.50), oneSnap.LatencyP99NS())
		}

		multiSnap := buildTestCollectorSnapshot(t, 5, map[string]int64{
			testCounterRecallSum:     1_900_000,
			testCounterRecallSamples: 2,
		})
		if multiSnap.Total != 5 || multiSnap.LatencyMinNS() <= 0 || multiSnap.LatencyMaxNS() < multiSnap.LatencyMinNS() {
			t.Fatalf("unexpected multi-request stats: total=%d min=%v max=%v",
				multiSnap.Total, multiSnap.LatencyMinNS(), multiSnap.LatencyMaxNS())
		}
		if multiSnap.LatencyQuantileNS(0.50) <= 0 || multiSnap.LatencyQuantileNS(0.90) <= 0 ||
			multiSnap.LatencyQuantileNS(0.95) <= 0 || multiSnap.LatencyP99NS() <= 0 {
			t.Fatalf("expected positive percentiles on multi-request snapshot")
		}
		if multiSnap.CustomCounters[testCounterRecallSum] != 1_900_000 || multiSnap.CustomCounters[testCounterRecallSamples] != 2 {
			t.Fatalf("CustomCounters = %+v, want recall_sum_ppm=1900000 recall_samples=2", multiSnap.CustomCounters)
		}

		// Total==1 fallback when LatPercentiles is nil uses LatencyMeanNS.
		singleNoDigest := &GlobalSnapshot{
			Total:     1,
			Latencies: &HistogramSnapshot{Total: 1, Mean: 12345},
		}
		if got := singleNoDigest.LatencyP99NS(); got != 12345 {
			t.Fatalf("LatencyP99NS with nil LatPercentiles and Total==1 = %v, want 12345", got)
		}
	})
}

func TestWriteLoadReports(t *testing.T) {
	t.Parallel()

	t.Run("WriteReport and LoadReport round-trip including Total==0 and CustomCounters", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		outPath := filepath.Join(dir, "nested", "stream_crud.json")

		zeroSnap := buildTestCollectorSnapshot(t, 0, nil)
		searchSnap := buildTestCollectorSnapshot(t, 4, map[string]int64{
			testCounterRecallSum:     950_000,
			testCounterRecallSamples: 1,
		})
		rep := &Report{
			Scenario: "stream_crud",
			Summary:  searchSnap,
			Strategies: []StrategyReport{
				{
					Name:    "Main",
					Summary: searchSnap,
					Operations: []OperationReport{
						{
							Name:    "Ops",
							Summary: searchSnap,
							Executions: []ExecutionReport{
								{Index: 0, Name: "Wait", Type: "wait", Mode: testModeUnary, Summary: zeroSnap},
								{Index: 1, Name: "Search", Type: testTypeSearch, Mode: "stream", Summary: searchSnap},
							},
						},
					},
				},
			},
		}

		if err := WriteReport(t.Context(), outPath, rep); err != nil {
			t.Fatalf("WriteReport failed: %v", err)
		}

		loaded, err := LoadReport(outPath)
		if err != nil {
			t.Fatalf("LoadReport failed: %v", err)
		}
		if loaded.SchemaVersion != ReportSchemaVersion || loaded.Scenario != "stream_crud" || loaded.Timestamp == "" {
			t.Fatalf("unexpected loaded report header: %+v", loaded)
		}
		if loaded.Summary.Total != 4 || loaded.Summary.CustomCounters[testCounterRecallSum] != 950_000 {
			t.Fatalf("unexpected loaded summary: %+v", loaded.Summary)
		}
		if gotZero := loaded.Strategies[0].Operations[0].Executions[0].Summary; gotZero.Total != 0 || gotZero.LatencyMinNS() != 0 {
			t.Fatalf("expected zero-request execution to round-trip cleanly, got %+v", gotZero)
		}

		// LoadReportsFromDir ignores subdirectories and non-JSON files.
		if mkErr := os.MkdirAll(filepath.Join(dir, "nested", "subdir"), 0o755); mkErr != nil {
			t.Fatalf("MkdirAll failed: %v", mkErr)
		}
		if wErr := os.WriteFile(filepath.Join(dir, "nested", "notes.txt"), []byte("ignore"), 0o644); wErr != nil {
			t.Fatalf("WriteFile failed: %v", wErr)
		}
		all, err := LoadReportsFromDir(filepath.Join(dir, "nested"))
		if err != nil {
			t.Fatalf("LoadReportsFromDir failed: %v", err)
		}
		if len(all) != 1 || all["stream_crud"] == nil {
			t.Fatalf("unexpected LoadReportsFromDir map: %+v", all)
		}
	})
}

func TestWriteLoadReportsValidationAndNormalize(t *testing.T) {
	t.Parallel()

	t.Run("WriteReport and LoadReport validation errors", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		canceledCtx, cancel := context.WithCancel(t.Context())
		cancel()

		if err := WriteReport(t.Context(), "", &Report{}); err == nil {
			t.Fatalf("expected error on empty path")
		}
		if err := WriteReport(t.Context(), filepath.Join(dir, "a.json"), nil); err == nil {
			t.Fatalf("expected error on nil report")
		}
		if err := WriteReport(canceledCtx, filepath.Join(dir, "a.json"), &Report{}); err == nil {
			t.Fatalf("expected error on canceled context")
		}
		if _, err := LoadReport(filepath.Join(dir, "missing.json")); err == nil {
			t.Fatalf("expected error on missing report")
		}
		emptyFile := filepath.Join(dir, "empty.json")
		if err := os.WriteFile(emptyFile, nil, 0o644); err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}
		if _, err := LoadReport(emptyFile); err == nil {
			t.Fatalf("expected error on empty file")
		}
		badSchema := filepath.Join(dir, "bad.json")
		if err := os.WriteFile(badSchema, []byte(`{"schema_version":"v99"}`), 0o644); err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}
		if _, err := LoadReport(badSchema); err == nil || !strings.Contains(err.Error(), "unsupported report schema_version") {
			t.Fatalf("expected schema_version error, got %v", err)
		}
		if _, err := LoadReportsFromDir(filepath.Join(dir, "missing-dir")); err == nil {
			t.Fatalf("expected error on missing directory")
		}
		// Additional validation paths: nil ctx, unwritable path, empty/dir/malformed LoadReport & LoadReportsFromDir.
		//nolint:staticcheck // SA1012: intentionally testing WriteReport nil-context guard.
		if err := WriteReport(nil, filepath.Join(dir, "nilctx.json"), &Report{}); err == nil {
			t.Fatalf("expected error on nil context")
		}
		blocker := filepath.Join(dir, "blocker")
		if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
			t.Fatalf("WriteFile(blocker) failed: %v", err)
		}
		if err := WriteReport(t.Context(), filepath.Join(blocker, "child.json"), &Report{}); err == nil {
			t.Fatalf("expected error on unwritable report path")
		}
		if _, err := LoadReport("   "); err == nil {
			t.Fatalf("expected error on empty LoadReport path")
		}
		if _, err := LoadReport(dir); err == nil {
			t.Fatalf("expected error when LoadReport path is a directory")
		}
		malformed := filepath.Join(dir, "malformed.json")
		if err := os.WriteFile(malformed, []byte("{bad"), 0o644); err != nil {
			t.Fatalf("WriteFile(malformed) failed: %v", err)
		}
		if _, err := LoadReport(malformed); err == nil {
			t.Fatalf("expected error on malformed JSON in LoadReport")
		}
		if _, err := LoadReportsFromDir("   "); err == nil {
			t.Fatalf("expected error on empty LoadReportsFromDir path")
		}
		if _, err := LoadReportsFromDir(blocker); err == nil {
			t.Fatalf("expected error when LoadReportsFromDir path is a regular file")
		}
		badDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(badDir, "bad.json"), []byte("{bad"), 0o644); err != nil {
			t.Fatalf("WriteFile(badDir) failed: %v", err)
		}
		if _, err := LoadReportsFromDir(badDir); err == nil {
			t.Fatalf("expected error when LoadReportsFromDir encounters malformed JSON")
		}
	})

	t.Run("normalizeReport fills nil summaries and slices across hierarchy", func(t *testing.T) {
		t.Parallel()
		nilStrat := normalizeReport(&Report{})
		if nilStrat.Summary == nil || nilStrat.Strategies == nil || len(nilStrat.Strategies) != 0 {
			t.Fatalf("unexpected normalizeReport(&Report{}): %+v", nilStrat)
		}

		nested := normalizeReport(&Report{
			Strategies: []StrategyReport{
				{Name: "S_nil_ops", Operations: nil},
				{
					Name: "S_with_ops",
					Operations: []OperationReport{
						{Name: "O_nil_execs", Executions: nil},
						{
							Name: "O_with_execs",
							Executions: []ExecutionReport{
								{Index: 0, Name: "E_nil_sum", Summary: nil},
							},
						},
					},
				},
			},
		})
		if nested.Summary == nil || nested.Strategies[0].Summary == nil || nested.Strategies[0].Operations == nil {
			t.Fatalf("expected nil Strategy Summary/Operations initialized: %+v", nested.Strategies[0])
		}
		if nested.Strategies[1].Operations[0].Summary == nil || nested.Strategies[1].Operations[0].Executions == nil {
			t.Fatalf("expected nil Operation Summary/Executions initialized: %+v", nested.Strategies[1].Operations[0])
		}
		if nested.Strategies[1].Operations[1].Executions[0].Summary == nil {
			t.Fatalf("expected nil Execution Summary initialized")
		}
	})
}

// NOT IMPLEMENTED BELOW
