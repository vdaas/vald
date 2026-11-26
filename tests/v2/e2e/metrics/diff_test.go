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
	"math"
	"path/filepath"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
)

const (
	testModeUnary        = "unary"
	testModeStream       = "stream"
	testTypeSearch       = "search"
	testExecEscaped      = "Exec`A|1"
	testExecName         = "Exec"
	testTypeEscaped      = "op`t"
	testModeEscaped      = "m|1"
	testScenarioBaseExec = "baseline_execs"

	envBenchDiffOutput   = "E2E_BENCH_DIFF_OUTPUT"
	envPRMetricsDir      = "E2E_PR_METRICS_DIR"
	envLatestMetricsDir  = "E2E_LATEST_METRICS_DIR"
	envNightlyMetricsDir = "E2E_NIGHTLY_METRICS_DIR"
)

func runBenchDiff(ctx context.Context, outPath string) error {
	prDir := strings.TrimSpace(os.Getenv(envPRMetricsDir))
	if prDir == "" {
		return errors.Errorf("%s is required when %s is set", envPRMetricsDir, envBenchDiffOutput)
	}
	md, err := GenerateBenchmarkDiff(
		prDir,
		strings.TrimSpace(os.Getenv(envLatestMetricsDir)),
		strings.TrimSpace(os.Getenv(envNightlyMetricsDir)),
		DefaultMaxCommentBytes,
	)
	if err != nil {
		return err
	}
	if _, err := file.OverWriteFile(ctx, outPath, strings.NewReader(md), defaultReportFilePerm); err != nil {
		return errors.Wrapf(err, "failed to write benchmark diff to %s", outPath)
	}
	return nil
}

// TestGenerateBenchmarkDiff is the pure-Go entry point for `make e2e/v2/diff`.
// When E2E_BENCH_DIFF_OUTPUT is unset (normal unit test runs), it is a no-op.
func TestGenerateBenchmarkDiff(t *testing.T) {
	outPath := strings.TrimSpace(os.Getenv(envBenchDiffOutput))
	if outPath == "" {
		return
	}
	if err := runBenchDiff(t.Context(), outPath); err != nil {
		t.Fatalf("failed to generate benchmark diff: %v", err)
	}
}

func TestRunBenchDiff(t *testing.T) {
	t.Setenv(envPRMetricsDir, "")
	if err := runBenchDiff(t.Context(), filepath.Join(t.TempDir(), "diff.md")); err == nil {
		t.Fatalf("expected error when %s is empty", envPRMetricsDir)
	}

	prDir := t.TempDir()
	latestDir := t.TempDir()
	diffOut := filepath.Join(t.TempDir(), "bench-diff.md")
	if err := WriteReport(t.Context(), filepath.Join(prDir, "stream_crud.json"), makeSyntheticReport(t, "stream_crud", 1.25, 0.8)); err != nil {
		t.Fatalf("WriteReport(pr) failed: %v", err)
	}
	if err := WriteReport(t.Context(), filepath.Join(latestDir, "stream_crud.json"), makeSyntheticReport(t, "stream_crud", 1.0, 1.0)); err != nil {
		t.Fatalf("WriteReport(latest) failed: %v", err)
	}

	t.Setenv(envPRMetricsDir, prDir)
	t.Setenv(envLatestMetricsDir, latestDir)
	t.Setenv(envNightlyMetricsDir, "")
	if err := runBenchDiff(t.Context(), diffOut); err != nil {
		t.Fatalf("runBenchDiff failed: %v", err)
	}
	raw, err := os.ReadFile(diffOut)
	if err != nil {
		t.Fatalf("ReadFile(diffOut) failed: %v", err)
	}
	md := string(raw)
	if !strings.HasPrefix(md, BenchmarkDiffMarker) || !strings.Contains(md, "`stream_crud`") {
		t.Fatalf("unexpected diff markdown:\n%s", md)
	}

	blocker := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile(blocker) failed: %v", err)
	}
	if err := runBenchDiff(t.Context(), filepath.Join(blocker, "diff.md")); err == nil {
		t.Fatalf("expected error on unwritable diff output path")
	}
}

func makeSyntheticSnapshot(
	t *testing.T, total uint64, dur time.Duration, meanNS, p99NS float64,
) *GlobalSnapshot {
	t.Helper()
	base := time.Date(2026, 11, 24, 21, 0, 0, 0, time.UTC)
	td, err := NewTDigest(WithTDigestCompression(100), WithTDigestNumShards(1))
	if err != nil {
		t.Fatalf("NewTDigest failed: %v", err)
	}
	if total > 0 && p99NS > 0 {
		for range 99 {
			td.Add(meanNS)
		}
		td.Add(p99NS)
	}
	return &GlobalSnapshot{
		Total:       total,
		StartTime:   base,
		LastUpdated: base.Add(dur),
		Latencies: &HistogramSnapshot{
			Total: total,
			Min:   meanNS / 2,
			Mean:  meanNS,
			Max:   p99NS * 1.1,
		},
		LatPercentiles: td,
	}
}

func makeSyntheticReport(
	t *testing.T, scenario string, qpsScale float64, latScale float64,
) *Report {
	t.Helper()
	dur := time.Duration(float64(10*time.Second) / qpsScale)
	sumSnap := makeSyntheticSnapshot(t, 1000, dur, 2e6*latScale, 5e6*latScale)
	insSnap := makeSyntheticSnapshot(t, 800, dur, 1.5e6*latScale, 4e6*latScale)
	schSnap := makeSyntheticSnapshot(t, 200, dur, 3e6*latScale, 8e6*latScale)
	return &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      scenario,
		Summary:       sumSnap,
		Strategies: []StrategyReport{
			{
				Name:    "CRUD",
				Summary: sumSnap,
				Operations: []OperationReport{
					{
						Name:    "Core",
						Summary: sumSnap,
						Executions: []ExecutionReport{
							{Index: 0, Name: "Insert", Type: "insert", Mode: testModeStream, Summary: insSnap},
							{Index: 1, Name: "Search", Type: testTypeSearch, Mode: testModeUnary, Summary: schSnap},
						},
					},
				},
			},
		},
	}
}

func TestBenchmarkDiff(t *testing.T) {
	t.Parallel()

	t.Run("FormatBenchmarkDiff empty and 3-way table + UTF-8 truncation", func(t *testing.T) {
		t.Parallel()
		emptyMD := FormatBenchmarkDiff(nil, nil, nil, 0)
		if !strings.HasPrefix(emptyMD, BenchmarkDiffMarker) || !strings.Contains(emptyMD, "No benchmark reports found") {
			t.Fatalf("unexpected empty diff markdown:\n%s", emptyMD)
		}

		prDir := t.TempDir()
		latestDir := t.TempDir()
		nightlyDir := t.TempDir()
		if err := WriteReport(t.Context(), filepath.Join(prDir, "stream_crud.json"), makeSyntheticReport(t, "stream_crud", 1.25, 0.8)); err != nil {
			t.Fatalf("WriteReport(pr) failed: %v", err)
		}
		if err := WriteReport(t.Context(), filepath.Join(latestDir, "stream_crud.json"), makeSyntheticReport(t, "stream_crud", 1.0, 1.0)); err != nil {
			t.Fatalf("WriteReport(latest) failed: %v", err)
		}
		if err := WriteReport(t.Context(), filepath.Join(nightlyDir, "baseline_only.json"), makeSyntheticReport(t, "baseline_only", 1.0, 1.0)); err != nil {
			t.Fatalf("WriteReport(nightly) failed: %v", err)
		}

		md, err := GenerateBenchmarkDiff(prDir, latestDir, nightlyDir, DefaultMaxCommentBytes)
		if err != nil {
			t.Fatalf("GenerateBenchmarkDiff failed: %v", err)
		}
		for _, want := range []string{
			BenchmarkDiffMarker,
			"### Scenario Summary",
			"`stream_crud`",
			"+25.0%",
			"`baseline_only`",
			"| CRUD / Core / [0] Insert | `insert` | `stream` |",
			naValue,
		} {
			if !strings.Contains(md, want) {
				t.Fatalf("generated markdown missing %q:\n%s", want, md)
			}
		}

		capped := FormatBenchmarkDiff(
			map[string]*Report{"a_café_λ": makeSyntheticReport(t, "a_café_λ", 1.0, 1.0), "b_naïve_μ": makeSyntheticReport(t, "b_naïve_μ", 1.0, 1.0)},
			nil,
			nil,
			250,
		)
		if len(capped) > 250 || !utf8.ValidString(capped) {
			t.Fatalf("expected valid UTF-8 truncated markdown <= 250 bytes, got len=%d valid=%v", len(capped), utf8.ValidString(capped))
		}
		if !strings.HasPrefix(capped, BenchmarkDiffMarker) {
			t.Fatalf("expected truncated markdown to retain marker")
		}
	})

	t.Run("FormatBenchmarkDiff escapes special characters and duplicate indices", func(t *testing.T) {
		t.Parallel()
		badName := "bad|name`with\nnewline<script>"
		prRep := makeSyntheticReport(t, badName, 1.0, 1.0)
		prRep.Strategies[0].Operations[0].Executions = []ExecutionReport{
			{Index: 0, Name: testExecEscaped, Type: testTypeEscaped, Mode: testModeEscaped, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 3e6, 7.5e6)},
			{Index: 0, Name: testExecEscaped, Type: testTypeEscaped, Mode: testModeEscaped, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 4e6, 8e6)},
		}
		latRep := makeSyntheticReport(t, badName, 1.0, 1.0)
		latRep.Strategies[0].Operations[0].Executions = []ExecutionReport{
			{Index: 0, Name: testExecEscaped, Type: testTypeEscaped, Mode: testModeEscaped, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 2e6, 5e6)},
			{Index: 0, Name: testExecEscaped, Type: testTypeEscaped, Mode: testModeEscaped, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 8e6, 16e6)},
		}

		md := FormatBenchmarkDiff(map[string]*Report{badName: prRep}, map[string]*Report{badName: latRep}, nil, DefaultMaxCommentBytes)
		for _, want := range []string{
			"<code>bad\\|name&#96;with newline&lt;script&gt;</code>",
			"<summary><code>bad|name`with newline&lt;script&gt;</code> — Summary &amp; Per-RPC Breakdown</summary>",
			"CRUD / Core / [0] Exec\\`A\\|1 | <code>op&#96;t</code> | `m\\|1` | 100 (0) |",
		} {
			if !strings.Contains(md, want) {
				t.Fatalf("escaped benchmark diff missing %q:\n%s", want, md)
			}
		}
		var dupRows []string
		for _, line := range strings.Split(md, "\n") {
			if strings.HasPrefix(line, "| CRUD / Core / [0] Exec\\`A\\|1 |") {
				dupRows = append(dupRows, line)
			}
		}
		if len(dupRows) != 2 {
			t.Fatalf("expected 2 duplicate-index execution rows, got %d:\n%s", len(dupRows), md)
		}
		if !strings.Contains(dupRows[0], "| +50.0% / +50.0% |") {
			t.Fatalf("first duplicate-index row expected +50.0%% / +50.0%% delta, got:\n%s", dupRows[0])
		}
		if !strings.Contains(dupRows[1], "| -50.0% / -50.0% |") {
			t.Fatalf("second duplicate-index row expected -50.0%% / -50.0%% delta, got:\n%s", dupRows[1])
		}
		if strings.Contains(md, "<script>") {
			t.Fatalf("expected raw <script> tag to be escaped in HTML summary:\n%s", md)
		}
	})

	t.Run("FormatBenchmarkDiff mixed duplicate indices", func(t *testing.T) {
		t.Parallel()
		prDup := makeSyntheticReport(t, "mixed", 1.0, 1.0)
		prDup.Strategies[0].Operations[0].Executions = []ExecutionReport{
			{Index: 0, Name: testExecName, Type: testTypeSearch, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 3e6, 7.5e6)},
			{Index: 0, Name: testExecName, Type: testTypeSearch, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 4e6, 8e6)},
		}
		latDistinct := makeSyntheticReport(t, "mixed", 1.0, 1.0)
		latDistinct.Strategies[0].Operations[0].Executions = []ExecutionReport{
			{Index: 0, Name: testExecName, Type: testTypeSearch, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 2e6, 5e6)},
			{Index: 1, Name: testExecName, Type: testTypeSearch, Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 100, 10*time.Second, 8e6, 16e6)},
		}
		md := FormatBenchmarkDiff(map[string]*Report{"mixed": prDup}, map[string]*Report{"mixed": latDistinct}, nil, DefaultMaxCommentBytes)
		var rows []string
		for _, line := range strings.Split(md, "\n") {
			if strings.HasPrefix(line, "| CRUD / Core / [0] Exec |") {
				rows = append(rows, line)
			}
		}
		if len(rows) != 2 {
			t.Fatalf("expected 2 rows in mixed duplicate-index diff, got %d:\n%s", len(rows), md)
		}
		if !strings.Contains(rows[0], "| +50.0% / +50.0% |") {
			t.Fatalf("first mixed row expected +50.0%% / +50.0%%, got:\n%s", rows[0])
		}
		if !strings.Contains(rows[1], "| -50.0% / -50.0% |") {
			t.Fatalf("second mixed row expected -50.0%% / -50.0%%, got:\n%s", rows[1])
		}
	})
}

func TestGenerateBenchmarkDiffErrors(t *testing.T) {
	t.Parallel()
	prDir := t.TempDir()
	if err := WriteReport(t.Context(), filepath.Join(prDir, "ok.json"), makeSyntheticReport(t, "ok", 1.0, 1.0)); err != nil {
		t.Fatalf("WriteReport failed: %v", err)
	}
	regFile := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(regFile, []byte("x"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	badJSONDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(badJSONDir, "broken.json"), []byte("{bad"), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	tests := []struct {
		name       string
		prDir      string
		latestDir  string
		nightlyDir string
	}{
		{name: "empty prDir", prDir: ""},
		{name: "regular file prDir", prDir: regFile},
		{name: "regular file latestDir", prDir: prDir, latestDir: regFile},
		{name: "regular file nightlyDir", prDir: prDir, nightlyDir: regFile},
		{name: "malformed JSON in latestDir", prDir: prDir, latestDir: badJSONDir},
		{name: "malformed JSON in nightlyDir", prDir: prDir, nightlyDir: badJSONDir},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := GenerateBenchmarkDiff(tc.prDir, tc.latestDir, tc.nightlyDir, DefaultMaxCommentBytes); err == nil {
				t.Fatalf("expected GenerateBenchmarkDiff error for %s", tc.name)
			}
		})
	}
}

func TestBenchmarkDiffFormattingHelpers(t *testing.T) {
	t.Parallel()
	prEmptyExec := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      "empty_execs",
		Summary:       makeSyntheticSnapshot(t, 0, 0, 0, 0),
	}
	nightlyOnlyExec := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      "nightly_only",
		Summary:       makeSyntheticSnapshot(t, 10, 500*time.Millisecond, 500, 1.5e9),
		Strategies: []StrategyReport{
			{
				Name:    "S",
				Summary: makeSyntheticSnapshot(t, 10, 500*time.Millisecond, 500, 1.5e9),
				Operations: []OperationReport{
					{
						Name:    "",
						Summary: makeSyntheticSnapshot(t, 10, 500*time.Millisecond, 500, 1.5e9),
						Executions: []ExecutionReport{
							{Index: 0, Name: "", Type: "wait", Mode: testModeUnary, Summary: makeSyntheticSnapshot(t, 10, 500*time.Millisecond, 500, 1.5e9)},
						},
					},
				},
			},
		},
	}
	latestEmptyExec := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      "latest_empty",
		Summary:       makeSyntheticSnapshot(t, 10, 2*time.Second, 2e5, 5e6),
	}

	md := FormatBenchmarkDiff(
		map[string]*Report{"empty_execs": prEmptyExec},
		map[string]*Report{"latest_empty": latestEmptyExec},
		map[string]*Report{"nightly_only": nightlyOnlyExec},
		DefaultMaxCommentBytes,
	)
	for _, want := range []string{
		"_No per-RPC executions recorded._",
		"| S / [0] wait | `wait` | `unary` | N/A | N/A |",
		"500.0ms",
		"500ns",
		"200.00µs",
		"1.65s",
	} {
		if !strings.Contains(md, want) {
			t.Fatalf("execution edge markdown missing %q:\n%s", want, md)
		}
	}

	// Explicit coverage for formatDurationValue, formatLatencyNS, formatDeltaPercent,
	// formatExecutionLabel, and truncateUTF8 / capMarkdownToBytes tight-cap fallback.
	if got := formatExecutionLabel("StrategyA", "", ExecutionReport{Index: 2, Type: testTypeSearch}); got != "StrategyA / [2] search" {
		t.Fatalf("formatExecutionLabel = %q, want %q", got, "StrategyA / [2] search")
	}
	for _, tc := range []struct {
		want string
		in   float64
	}{
		{"0s", 0},
		{"0s", -1},
		{"250.0ms", 0.25},
		{"2.50s", 2.5},
	} {
		if got := formatDurationValue(tc.in); got != tc.want {
			t.Fatalf("formatDurationValue(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
	for _, tc := range []struct {
		want string
		in   float64
	}{
		{"0ns", 0},
		{"0ns", -5},
		{"500ns", 500},
		{"250.00µs", 250_000},
		{"5.00ms", 5_000_000},
		{"1.50s", 1_500_000_000},
	} {
		if got := formatLatencyNS(tc.in); got != tc.want {
			t.Fatalf("formatLatencyNS(%v) = %q, want %q", tc.in, got, tc.want)
		}
	}
	for _, tc := range []struct {
		want       string
		curr, base float64
	}{
		{naValue, 10, 0},
		{naValue, 10, math.NaN()},
		{naValue, 10, math.Inf(1)},
		{naValue, -1, 10},
		{naValue, math.NaN(), 10},
		{"0.0%", 100.01, 100.0},
		{"+20.0%", 120, 100},
	} {
		if got := formatDeltaPercent(tc.curr, tc.base); got != tc.want {
			t.Fatalf("formatDeltaPercent(%v, %v) = %q, want %q", tc.curr, tc.base, got, tc.want)
		}
	}

	// Tight maxBytes below len(header)+len(truncationNotice) exercises truncateUTF8 via boundedMarkdownBuilder.finish.
	utf8Header := "éééé"
	mb := newBoundedMarkdownBuilder(utf8Header, 5)
	mb.append("block")
	tight := mb.finish()
	if len(tight) > 5 || !utf8.ValidString(tight) || tight != "éé" {
		t.Fatalf("boundedMarkdownBuilder.finish tight UTF-8 = %q (len=%d, valid=%v), want %q", tight, len(tight), utf8.ValidString(tight), "éé")
	}
	if got := truncateUTF8("abc", 10); got != "abc" {
		t.Fatalf("truncateUTF8 short string = %q, want %q", got, "abc")
	}
	if got := truncateUTF8("abc", 0); got != "" {
		t.Fatalf("truncateUTF8 zero max = %q, want empty", got)
	}
}

func TestBenchmarkDiffEmptyPRExecutions(t *testing.T) {
	t.Parallel()
	prNoExec := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      testScenarioBaseExec,
		Summary:       makeSyntheticSnapshot(t, 100, 10*time.Second, 2e6, 5e6),
	}
	latestWithExec := makeSyntheticReport(t, testScenarioBaseExec, 1.0, 1.0)
	md := FormatBenchmarkDiff(
		map[string]*Report{testScenarioBaseExec: prNoExec},
		map[string]*Report{testScenarioBaseExec: latestWithExec},
		nil,
		DefaultMaxCommentBytes,
	)
	if strings.Contains(md, "_No per-RPC executions recorded._") {
		t.Fatalf("expected baseline executions to be shown even when PR has none:\n%s", md)
	}
	if !strings.Contains(md, "| CRUD / Core / [0] Insert | `insert` | `stream` | N/A | N/A | N/A | N/A | N/A | N/A | N/A | N/A | N/A | N/A |") {
		t.Fatalf("expected Insert row with N/A PR values and deltas:\n%s", md)
	}
}

func TestBenchmarkDiffUnavailableVsMeasuredZero(t *testing.T) {
	t.Parallel()
	baseline := makeSyntheticReport(t, "s", 1.0, 1.0)

	// Unavailable: Total == 0, no timestamps, no latency samples.
	prUnavail := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      "s",
		Summary:       makeSyntheticSnapshot(t, 0, 0, 0, 0),
		Strategies: []StrategyReport{
			{
				Name: "CRUD",
				Operations: []OperationReport{
					{
						Name: "Core",
						Executions: []ExecutionReport{
							{Index: 0, Name: "Insert", Type: "insert", Mode: testModeStream, Summary: makeSyntheticSnapshot(t, 0, 0, 0, 0)},
						},
					},
				},
			},
		},
	}
	unavailMD := FormatBenchmarkDiff(
		map[string]*Report{"s": prUnavail},
		map[string]*Report{"s": baseline},
		nil,
		DefaultMaxCommentBytes,
	)
	if strings.Contains(unavailMD, "-100.0%") {
		t.Fatalf("unavailable PR measurements must not render as -100.0%%:\n%s", unavailMD)
	}
	for _, want := range []string{
		"| QPS | N/A | 100.00 | N/A |",
		"| Duration | N/A | 10.00s | N/A |",
		"| Latency Min | N/A | 1.00ms | N/A |",
		"| Latency Mean | N/A | 2.00ms | N/A |",
		"| Latency Max | N/A | 5.50ms | N/A |",
		"| Latency P99 | N/A | 3.50ms | N/A |",
	} {
		if !strings.Contains(unavailMD, want) {
			t.Fatalf("unavailable markdown missing %q:\n%s", want, unavailMD)
		}
	}

	// Measured zeroes: Total > 0, StartTime == LastUpdated (0s duration), and 0ns latency samples.
	tdZero, err := NewTDigest(WithTDigestCompression(100), WithTDigestNumShards(1))
	if err != nil {
		t.Fatalf("NewTDigest failed: %v", err)
	}
	tdZero.Add(0)
	baseTime := time.Date(2026, 11, 24, 21, 0, 0, 0, time.UTC)
	zeroMeasuredSnap := &GlobalSnapshot{
		Total:       10,
		StartTime:   baseTime,
		LastUpdated: baseTime,
		Latencies: &HistogramSnapshot{
			Total: 10,
			Min:   0,
			Mean:  0,
			Max:   0,
		},
		LatPercentiles: tdZero,
	}
	prMeasuredZero := &Report{
		SchemaVersion: ReportSchemaVersion,
		Scenario:      "s",
		Summary:       zeroMeasuredSnap,
	}
	zeroMD := FormatBenchmarkDiff(
		map[string]*Report{"s": prMeasuredZero},
		map[string]*Report{"s": baseline},
		nil,
		DefaultMaxCommentBytes,
	)
	for _, want := range []string{
		"| QPS | N/A | 100.00 | N/A |",
		"| Duration | 0s | 10.00s | -100.0% |",
		"| Latency Min | 0ns | 1.00ms | -100.0% |",
		"| Latency Mean | 0ns | 2.00ms | -100.0% |",
		"| Latency Max | 0ns | 5.50ms | -100.0% |",
		"| Latency P99 | 0ns | 3.50ms | -100.0% |",
	} {
		if !strings.Contains(zeroMD, want) {
			t.Fatalf("measured-zero markdown missing %q:\n%s", want, zeroMD)
		}
	}
}

// NOT IMPLEMENTED BELOW
