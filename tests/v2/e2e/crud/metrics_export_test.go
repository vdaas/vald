//go:build e2e

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

package crud

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"github.com/vdaas/vald/tests/v2/e2e/metrics"
)

type spyTB struct {
	testing.TB
	errs []string
}

func (*spyTB) Helper() {}

func (s *spyTB) Errorf(format string, args ...any) {
	s.errs = append(s.errs, fmt.Sprintf(format, args...))
}

func buildMetricsFixture(t *testing.T, latScale float64) *config.Data {
	t.Helper()
	newCol := func(counters ...string) metrics.Collector {
		t.Helper()
		col, err := metrics.NewCollector(metrics.WithCustomCounters(counters...))
		if err != nil {
			t.Fatalf("NewCollector failed: %v", err)
		}
		return col
	}

	base := time.Date(2026, 11, 24, 21, 0, 0, 0, time.UTC)
	insertCol := newCol()
	for i := range 3 {
		lat := time.Duration(float64(4*time.Millisecond) * latScale)
		insertCol.Record(t.Context(), uint64(i+1), &metrics.RequestResult{
			StartedAt: base.Add(time.Duration(i) * 100 * time.Millisecond),
			EndedAt:   base.Add(time.Duration(i)*100*time.Millisecond + lat),
			Latency:   lat,
			Status:    codes.OK,
		})
	}
	waitCol := newCol()
	searchCol := newCol(recallSumCounterName, recallSamplesCounterName)
	searchCol.Record(t.Context(), 1, &metrics.RequestResult{
		StartedAt: base.Add(300 * time.Millisecond),
		EndedAt:   base.Add(300*time.Millisecond + time.Duration(float64(10*time.Millisecond)*latScale)),
		Latency:   time.Duration(float64(10*time.Millisecond) * latScale),
		Status:    codes.OK,
	})
	searchCol.IncCounter(recallSumCounterName, 950_000)
	searchCol.IncCounter(recallSamplesCounterName, 1)

	opCol := newCol()
	for _, c := range []metrics.Collector{insertCol, waitCol, searchCol} {
		if err := c.MergeInto(opCol); err != nil {
			t.Fatalf("MergeInto(opCol) failed: %v", err)
		}
	}
	stCol := newCol()
	if err := opCol.MergeInto(stCol); err != nil {
		t.Fatalf("MergeInto(stCol) failed: %v", err)
	}
	rootCol := newCol()
	if err := stCol.MergeInto(rootCol); err != nil {
		t.Fatalf("MergeInto(rootCol) failed: %v", err)
	}

	return &config.Data{
		FilePath:  "tests/v2/e2e/assets/stream_crud.yaml",
		Collector: rootCol,
		Strategies: []*config.Strategy{
			{
				Name:      "Initial",
				Collector: stCol,
				Operations: []*config.Operation{
					{
						Name:      "CRUD",
						Collector: opCol,
						Executions: []*config.Execution{
							{Name: "Insert", Type: config.OpInsert, Mode: config.OperationStream, Collector: insertCol},
							{Name: "Wait", Type: config.OpWait, Mode: config.OperationUnary, Collector: waitCol},
							{Name: "Search", Type: config.OpSearch, Mode: config.OperationStream, Collector: searchCol},
						},
					},
				},
			},
		},
	}
}

func TestMetricsExportAndBenchDiff(t *testing.T) {
	t.Run("exportMetricsReport writes GlobalSnapshot report and handles noise retry markers", func(t *testing.T) {
		dir := t.TempDir()
		outPath := filepath.Join(dir, "current.json")
		basePath := filepath.Join(dir, "baseline.json")

		t.Setenv(envScenarioName, "stream_crud-helm-chart-k3d-null")
		exportMetricsReport(t.Context(), t, buildMetricsFixture(t, 1.0), outPath)
		loaded, err := metrics.LoadReport(outPath)
		if err != nil {
			t.Fatalf("LoadReport failed: %v", err)
		}
		if loaded.Scenario != "stream_crud-helm-chart-k3d-null" || loaded.Summary.Total != 4 {
			t.Fatalf("unexpected exported report: %+v", loaded)
		}
		if gotSum, gotSamples := loaded.Summary.CustomCounters[recallSumCounterName], loaded.Summary.CustomCounters[recallSamplesCounterName]; gotSum != 950_000 || gotSamples != 1 {
			t.Fatalf("expected CustomCounters recall=(950000, 1), got (%d, %d)", gotSum, gotSamples)
		}

		// Write slower baseline so current triggers .retry marker.
		exportMetricsReport(t.Context(), t, buildMetricsFixture(t, 4.0), basePath)
		t.Setenv(envMetricsBaseline, basePath)
		t.Setenv(envNoiseThreshold, "0.20")
		exportMetricsReport(t.Context(), t, buildMetricsFixture(t, 1.0), outPath)
		if _, statErr := os.Stat(outPath + retryFileSuffix); statErr != nil {
			t.Fatalf("expected %s to exist on noisy run: %v", outPath+retryFileSuffix, statErr)
		}

		// Re-run with matching baseline clears stale .retry marker.
		exportMetricsReport(t.Context(), t, buildMetricsFixture(t, 1.0), basePath)
		exportMetricsReport(t.Context(), t, buildMetricsFixture(t, 1.0), outPath)
		if _, statErr := os.Stat(outPath + retryFileSuffix); !errors.Is(statErr, os.ErrNotExist) {
			t.Fatalf("expected stale .retry marker removed, statErr=%v", statErr)
		}
	})

	t.Run("exportMetricsReport rejects invalid E2E_NOISE_THRESHOLD and removes stale .retry marker", func(t *testing.T) {
		for _, raw := range []string{"abc", "-1", "0", "NaN", "+Inf"} {
			t.Run(raw, func(t *testing.T) {
				dir := t.TempDir()
				outPath := filepath.Join(dir, "current.json")
				basePath := filepath.Join(dir, "baseline.json")
				exportMetricsReport(t.Context(), t, buildMetricsFixture(t, 1.0), basePath)

				retryPath := outPath + retryFileSuffix
				if err := os.WriteFile(retryPath, []byte("stale\n"), 0o644); err != nil {
					t.Fatalf("WriteFile(stale retry) failed: %v", err)
				}

				t.Setenv(envMetricsBaseline, basePath)
				t.Setenv(envNoiseThreshold, raw)
				spy := &spyTB{TB: t}
				exportMetricsReport(t.Context(), spy, buildMetricsFixture(t, 1.0), outPath)
				if len(spy.errs) == 0 || !strings.Contains(spy.errs[0], envNoiseThreshold) {
					t.Fatalf("expected %s error for %q, got %v", envNoiseThreshold, raw, spy.errs)
				}
				if _, err := metrics.LoadReport(outPath); err != nil {
					t.Fatalf("expected report still written before threshold validation: %v", err)
				}
				if _, statErr := os.Stat(retryPath); !errors.Is(statErr, os.ErrNotExist) {
					t.Fatalf("expected stale .retry marker removed on invalid threshold, statErr=%v", statErr)
				}
			})
		}
	})

	t.Run("buildMetricsReport and resolveScenarioName handle nil elements and fallbacks", testMetricsReportBuildAndResolve)
	t.Run("exportMetricsReport unwritable output and missing baseline bootstrap", testExportMetricsReportEdgeCases)
	t.Run("hasParamsFlag", testHasParamsFlag)
}

func testMetricsReportBuildAndResolve(t *testing.T) {
	nilRep := buildMetricsReport(nil, "custom")
	if nilRep.SchemaVersion != metrics.ReportSchemaVersion || nilRep.Scenario != "custom" || nilRep.Summary == nil || len(nilRep.Strategies) != 0 {
		t.Fatalf("unexpected buildMetricsReport(nil): %+v", nilRep)
	}

	sparseCfg := &config.Data{
		FilePath: "tests/v2/e2e/assets/stream_crud.yaml",
		Strategies: []*config.Strategy{
			nil,
			{
				Name: "S1",
				Operations: []*config.Operation{
					nil,
					{
						Name:       "O1",
						Executions: []*config.Execution{nil, {Name: "Wait", Type: config.OpWait, Mode: config.OperationUnary}},
					},
				},
			},
		},
	}
	sparseRep := buildMetricsReport(sparseCfg, resolveScenarioName("  ", sparseCfg.FilePath))
	if sparseRep.Scenario != "stream_crud" || len(sparseRep.Strategies) != 1 ||
		len(sparseRep.Strategies[0].Operations) != 1 || len(sparseRep.Strategies[0].Operations[0].Executions) != 1 {
		t.Fatalf("unexpected sparse buildMetricsReport result: %+v", sparseRep)
	}

	for _, tc := range []struct {
		name, env, cfgPath, want string
	}{
		{name: "both empty", env: "", cfgPath: "", want: ""},
		{name: "whitespace env and empty path", env: "  \t ", cfgPath: "  ", want: ""},
		{name: "fallback to cfgData stem", env: "", cfgPath: "tests/v2/e2e/assets/unary_crud.yaml", want: "unary_crud"},
		{name: "explicit env wins", env: "override-name", cfgPath: "tests/v2/e2e/assets/unary_crud.yaml", want: "override-name"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := resolveScenarioName(tc.env, tc.cfgPath); got != tc.want {
				t.Fatalf("resolveScenarioName(%q, %q) = %q, want %q", tc.env, tc.cfgPath, got, tc.want)
			}
		})
	}
}

func testExportMetricsReportEdgeCases(t *testing.T) {
	dir := t.TempDir()
	blocker := filepath.Join(dir, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(blocker) failed: %v", err)
	}
	spy := &spyTB{TB: t}
	exportMetricsReport(t.Context(), spy, buildMetricsFixture(t, 1.0), filepath.Join(blocker, "out.json"))
	if len(spy.errs) == 0 {
		t.Fatalf("expected spyTB error on unwritable output path")
	}

	outPath := filepath.Join(dir, "bootstrap.json")
	t.Setenv(envMetricsBaseline, filepath.Join(dir, "nonexistent-baseline.json"))
	t.Setenv(envNoiseThreshold, "")
	spyClean := &spyTB{TB: t}
	exportMetricsReport(t.Context(), spyClean, buildMetricsFixture(t, 1.0), outPath)
	if len(spyClean.errs) != 0 {
		t.Fatalf("unexpected errors on missing baseline bootstrap: %v", spyClean.errs)
	}
	if _, err := os.Stat(outPath + retryFileSuffix); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected no .retry marker on missing baseline bootstrap, got %v", err)
	}
}

func testHasParamsFlag(t *testing.T) {
	const testCfgFile = "cfg.yaml"
	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "short -c", args: []string{"-c", testCfgFile}, want: true},
		{name: "long --config=", args: []string{"--config=" + testCfgFile}, want: true},
		{name: "short -f", args: []string{"-f", testCfgFile}, want: true},
		{name: "short -h", args: []string{"-h"}, want: true},
		{name: "long --version", args: []string{"--version"}, want: true},
		{name: "stops at --", args: []string{"--", "-c", testCfgFile}, want: false},
		{name: "test flags only", args: []string{"-test.run=^$"}, want: false},
		{name: "positional args only", args: []string{"positional"}, want: false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := hasParamsFlag(tc.args); got != tc.want {
				t.Fatalf("hasParamsFlag(%v) = %v, want %v", tc.args, got, tc.want)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
