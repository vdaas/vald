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
	"math"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
)

const (
	testExecIdxInfo        = "IndexInfo"
	testTypeIdxInfo        = "index_info"
	testScenarioExact30    = "exact_30_boundary"
	testScenarioShiftedIdx = "shifted_idx"
	testSubLatency         = "latency"
)

func TestCheckNoise(t *testing.T) {
	t.Parallel()

	cfg := DefaultNoiseConfig()
	base := makeSyntheticReport(t, "stream_crud", 1.0, 1.0)
	within := makeSyntheticReport(t, "stream_crud", 1.1, 1.1)
	exact30Base := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      testScenarioExact30,
		Summary:       makeSyntheticSnapshot(t, 100, 10*time.Second, 2e6, 5e6),
	}
	baseDurNS := float64(10 * time.Second)
	exact30Truncated := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      testScenarioExact30,
		// 10s / 1.3 truncates to 7692307692ns (~0.300000000052 QPS delta).
		Summary: makeSyntheticSnapshot(t, 100, time.Duration(baseDurNS/1.3), 2e6, 5e6),
	}
	over30 := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      testScenarioExact30,
		Summary:       makeSyntheticSnapshot(t, 100, time.Duration(baseDurNS/1.301), 2e6, 5e6),
	}
	noisyQPS := makeSyntheticReport(t, "stream_crud", 2.0, 1.0)
	noisyLat := makeSyntheticReport(t, "stream_crud", 1.0, 2.5)
	tiny := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      "tiny",
		Summary:       makeSyntheticSnapshot(t, 1, 10*time.Millisecond, 1e6, 1e6),
	}
	dupZeroIdx := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      "dup_zero_idx",
		Summary:       makeSyntheticSnapshot(t, 200, 10*time.Second, 2e6, 5e6),
		Strategies: []StrategyReport{
			{
				Name: "S",
				Operations: []OperationReport{
					{
						Name: "O",
						Executions: []ExecutionReport{
							{Index: 0, Name: testExecIdxInfo, Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 2e6, 5e6)},
							{Index: 0, Name: testExecIdxInfo, Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 10e6, 25e6)},
						},
					},
				},
			},
		},
	}
	dupZeroIdxSecondNoisy := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      "dup_zero_idx",
		Summary:       makeSyntheticSnapshot(t, 200, 10*time.Second, 2e6, 5e6),
		Strategies: []StrategyReport{
			{
				Name: "S",
				Operations: []OperationReport{
					{
						Name: "O",
						Executions: []ExecutionReport{
							{Index: 0, Name: testExecIdxInfo, Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 2e6, 5e6)},
							{Index: 0, Name: testExecIdxInfo, Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 20e6, 50e6)},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		curr      *Report
		base      *Report
		name      string
		wantSub   string
		cfg       NoiseConfig
		wantRetry bool
		wantErr   bool
	}{
		{name: "invalid threshold errors", curr: within, base: base, cfg: NoiseConfig{Threshold: 0}, wantErr: true},
		{name: "NaN min duration errors", curr: within, base: base, cfg: NoiseConfig{Threshold: 0.3, MinDurationSec: math.NaN()}, wantErr: true},
		{name: "nil current errors", curr: nil, base: base, cfg: cfg, wantErr: true},
		{name: "nil baseline skips cleanly", curr: within, base: nil, cfg: cfg, wantRetry: false, wantSub: "nil"},
		{name: "below MinTotal skips cleanly", curr: tiny, base: tiny, cfg: cfg, wantRetry: false, wantSub: "guards"},
		{name: "within threshold passes", curr: within, base: base, cfg: cfg, wantRetry: false},
		{name: "exact 30pct target truncated to integer ns does not spuriously retry", curr: exact30Truncated, base: exact30Base, cfg: cfg, wantRetry: false},
		{name: "30.1pct QPS delta exceeds threshold and triggers retry", curr: over30, base: exact30Base, cfg: cfg, wantRetry: true, wantSub: "delta"},
		{name: "high QPS delta triggers retry", curr: noisyQPS, base: base, cfg: cfg, wantRetry: true, wantSub: "delta"},
		{name: "high latency delta triggers retry", curr: noisyLat, base: base, cfg: cfg, wantRetry: true, wantSub: testSubLatency},
		{name: "duplicate Index self-comparison falls back to occurrence", curr: dupZeroIdx, base: dupZeroIdx, cfg: cfg, wantRetry: false},
		{name: "duplicate Index mismatch on second occurrence triggers retry", curr: dupZeroIdxSecondNoisy, base: dupZeroIdx, cfg: cfg, wantRetry: true, wantSub: testSubLatency},
		{
			name: "shifted execution indices match distinct baseline entries 1-to-1",
			curr: &Report{
				SchemaVersion: ReportSchemaVersion,
				Scenario:      testScenarioShiftedIdx,
				Summary:       makeSyntheticSnapshot(t, 200, 10*time.Second, 6e6, 15e6),
				Strategies: []StrategyReport{
					{
						Name: "S",
						Operations: []OperationReport{
							{
								Name: "O",
								Executions: []ExecutionReport{
									{Index: 0, Name: "A", Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 2e6, 5e6)},
									{Index: 1, Name: "A", Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 10e6, 25e6)},
								},
							},
						},
					},
				},
			},
			base: &Report{
				SchemaVersion: ReportSchemaVersion,
				Scenario:      testScenarioShiftedIdx,
				Summary:       makeSyntheticSnapshot(t, 200, 10*time.Second, 6e6, 15e6),
				Strategies: []StrategyReport{
					{
						Name: "S",
						Operations: []OperationReport{
							{
								Name: "O",
								Executions: []ExecutionReport{
									{Index: 1, Name: "A", Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 2e6, 5e6)},
									{Index: 2, Name: "A", Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 10e6, 25e6)},
								},
							},
						},
					},
				},
			},
			cfg:       cfg,
			wantRetry: false,
		},
		{
			name: "shifted execution indices detect genuine regression on second entry",
			curr: &Report{
				SchemaVersion: ReportSchemaVersion,
				Scenario:      testScenarioShiftedIdx,
				Summary:       makeSyntheticSnapshot(t, 200, 10*time.Second, 6e6, 15e6),
				Strategies: []StrategyReport{
					{
						Name: "S",
						Operations: []OperationReport{
							{
								Name: "O",
								Executions: []ExecutionReport{
									{Index: 0, Name: "A", Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 2e6, 5e6)},
									{Index: 1, Name: "A", Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 20e6, 50e6)},
								},
							},
						},
					},
				},
			},
			base: &Report{
				SchemaVersion: ReportSchemaVersion,
				Scenario:      testScenarioShiftedIdx,
				Summary:       makeSyntheticSnapshot(t, 200, 10*time.Second, 6e6, 15e6),
				Strategies: []StrategyReport{
					{
						Name: "S",
						Operations: []OperationReport{
							{
								Name: "O",
								Executions: []ExecutionReport{
									{Index: 1, Name: "A", Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 2e6, 5e6)},
									{Index: 2, Name: "A", Type: testTypeIdxInfo, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 10e6, 25e6)},
								},
							},
						},
					},
				},
			},
			cfg:       cfg,
			wantRetry: true,
			wantSub:   testSubLatency,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			retry, reason, err := CheckNoise(tc.curr, tc.base, tc.cfg)
			if (err != nil) != tc.wantErr {
				t.Fatalf("CheckNoise err = %v, wantErr=%v", err, tc.wantErr)
			}
			if retry != tc.wantRetry {
				t.Fatalf("CheckNoise retry = %v (reason=%q), want %v", retry, reason, tc.wantRetry)
			}
			if tc.wantSub != "" && !strings.Contains(reason, tc.wantSub) {
				t.Fatalf("CheckNoise reason = %q, want substring %q", reason, tc.wantSub)
			}
		})
	}
}

func TestCheckNoiseFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	currReport := makeSyntheticReport(t, "s1", 1.0, 1.0)
	currPath := filepath.Join(dir, "current.json")
	if err := WriteReport(t.Context(), currPath, currReport); err != nil {
		t.Fatalf("WriteReport(curr) failed: %v", err)
	}
	emptyBase := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(emptyBase, nil, 0o644); err != nil {
		t.Fatalf("WriteFile(empty) failed: %v", err)
	}
	noisyBase := filepath.Join(dir, "noisy_base.json")
	if err := WriteReport(t.Context(), noisyBase, makeSyntheticReport(t, "s1", 3.0, 1.0)); err != nil {
		t.Fatalf("WriteReport(noisy) failed: %v", err)
	}

	if _, _, err := CheckNoiseBaselineFile(nil, "", DefaultNoiseConfig()); err == nil {
		t.Fatalf("expected error when current report is nil")
	}
	if retry, reason, err := CheckNoiseBaselineFile(currReport, noisyBase, DefaultNoiseConfig()); err != nil || !retry || reason == "" {
		t.Fatalf("CheckNoiseBaselineFile noisy = (%v, %q, %v), want (true, non-empty, nil)", retry, reason, err)
	}

	if retry, _, err := CheckNoiseFile(currPath, "", DefaultNoiseConfig()); err != nil || retry {
		t.Fatalf("empty baseline path = (%v, %v), want (false, nil)", retry, err)
	}
	if retry, _, err := CheckNoiseFile(currPath, filepath.Join(dir, "missing.json"), DefaultNoiseConfig()); err != nil || retry {
		t.Fatalf("missing baseline = (%v, %v), want (false, nil)", retry, err)
	}
	if retry, _, err := CheckNoiseFile(currPath, emptyBase, DefaultNoiseConfig()); err != nil || retry {
		t.Fatalf("empty baseline = (%v, %v), want (false, nil)", retry, err)
	}
	if _, _, err := CheckNoiseFile(currPath, dir, DefaultNoiseConfig()); err == nil {
		t.Fatalf("expected error when baseline path is a directory")
	}
	if retry, reason, err := CheckNoiseFile(currPath, noisyBase, DefaultNoiseConfig()); err != nil || !retry || reason == "" {
		t.Fatalf("noisy baseline = (%v, %q, %v), want (true, non-empty, nil)", retry, reason, err)
	}

	malformedBase := filepath.Join(dir, "malformed.json")
	if err := os.WriteFile(malformedBase, []byte("{not-json"), 0o644); err != nil {
		t.Fatalf("WriteFile(malformed) failed: %v", err)
	}
	if _, _, err := CheckNoiseFile(filepath.Join(dir, "missing_curr.json"), "", DefaultNoiseConfig()); err == nil {
		t.Fatalf("expected error when current report path is missing")
	}
	if _, _, err := CheckNoiseFile(malformedBase, "", DefaultNoiseConfig()); err == nil {
		t.Fatalf("expected error when current report JSON is malformed")
	}
	if _, _, err := CheckNoiseFile(currPath, malformedBase, DefaultNoiseConfig()); err == nil {
		t.Fatalf("expected error when baseline JSON is malformed")
	}
}

// NOT IMPLEMENTED BELOW
