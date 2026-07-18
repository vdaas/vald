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
// tests/v2/e2e/metrics/chart package and intentionally does not compile yet. It is
// excluded from the default build via the `ignore` tag below for the same reason
// documented in ../recall_test.go: this repository's Stop-hook / golangci-lint
// gate has no TDD RED-phase exception for typecheck failures on a brand new
// package. The next Maker should delete the `//go:build ignore` line (only) and
// add chart.go implementing the contract documented below.

// Package chart is a NEW package (does not exist yet - RED phase of TDD).
//
// It is deliberately independent from hack/tools/metrics and hack/benchmark/metrics
// (per the task's constraint: "hack/benchmark 配下は一切変更しない"). It reuses the
// same rendering approach (gonum/plot scatter/line-points, log-scale Y axis,
// HCL-spaced colors via go-colorful - see hack/tools/metrics/main.go for the
// reference implementation this design is modeled after) but is scoped to
// tests/v2/e2e/metrics.GlobalSnapshot (recall + AchievedQPS) instead of the
// hack/benchmark-specific gob-encoded Metrics/SearchMetrics types.
//
// Expected API contract (to be implemented by the next Maker in a new
// tests/v2/e2e/metrics/chart/chart.go):
//
//	type Point struct {
//		Recall float64 // x-axis, expected in [0,1] but not enforced here
//		QPS    float64 // y-axis, plotted on a log scale - must be > 0 to be plottable
//	}
//
//	type Series struct {
//		Label  string
//		Points []Point
//	}
//
//	type Config struct {
//		Title, XLabel, YLabel string
//		Width, Height         int // pixels; <= 0 falls back to DefaultConfig()'s values
//	}
//
//	func DefaultConfig() Config
//
//	type RenderResult struct {
//		SVG         []byte
//		SeriesCount int            // number of series that contributed >=1 rendered point
//		PointCount  map[string]int // per-label count of rendered (finite, QPS>0) points
//	}
//
//	func RenderParetoSVG(series []Series, cfg Config) (*RenderResult, error)
//
//	func LoadParetoSeries(label string, jsonPaths ...string) (Series, error)
//
// Design rationale for RenderResult carrying SeriesCount/PointCount explicitly:
// gonum/plot's SVG backend (vgsvg, via the ajstarks/svgo writer) renders scatter
// glyphs as generic stroked <path> elements (see RingGlyph.DrawGlyph in
// gonum.org/v1/plot/vg/draw) that are visually indistinguishable, at the raw SVG
// tag level, from grid lines and axis ticks. Asserting "N points were drawn" by
// counting <path> elements in the output would therefore be both unstable across
// gonum/plot versions AND unable to distinguish data points from decorative
// elements. Returning the counts computed at render time (ground truth, no
// re-parsing needed) is both a more robust test seam AND a more useful library API
// (callers - e.g. a future `make e2e/pareto-chart` helper - can log "rendered N
// points across M series" without parsing their own SVG output). The tests below
// still additionally validate the SVG is well-formed XML and spot-check
// human-readable content (title/axis labels/series legend text/distinct per-series
// colors) as a secondary regression guard.
package chart

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/v2/e2e/metrics"
)

// --- helpers -----------------------------------------------------------------

// mustBeWellFormedSVG decodes data as a full XML token stream (not into a typed
// struct) so that leading <?xml ...?> / <!-- comment --> preludes emitted by
// gonum/plot's vgsvg backend do not need special-casing, and returns the local
// name of the outermost element plus the raw byte content for further substring
// inspection by callers.
func mustBeWellFormedSVG(t *testing.T, data []byte) (rootLocal string) {
	t.Helper()
	if len(data) == 0 {
		t.Fatal("SVG output is empty")
	}
	dec := xml.NewDecoder(bytes.NewReader(data))
	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("SVG output is not well-formed XML: %v\n--- output ---\n%s", err, data)
		}
		if se, ok := tok.(xml.StartElement); ok && rootLocal == "" {
			rootLocal = se.Name.Local
		}
	}
	if rootLocal == "" {
		t.Fatal("SVG output has no root element")
	}
	return rootLocal
}

// countDistinctColors returns the number of distinct "#RRGGBB" stroke/fill color
// literals in an SVG document, excluding pure black (#000000, used for axes/text/
// grid by gonum/plot's default theme) and pure white (#FFFFFF, the canvas
// background fill). This is used as a coarse proxy for "each series got its own
// color", without depending on which specific hue-generation algorithm the
// implementation uses (go-colorful HCL spacing, as in hack/tools/metrics, or
// otherwise).
func countDistinctColors(data []byte) int {
	seen := make(map[string]struct{})
	s := string(data)
	for i := 0; i+7 <= len(s); i++ {
		if s[i] != '#' {
			continue
		}
		hex := s[i : i+7]
		isHex := true
		for _, r := range hex[1:] {
			if !isHexDigit(r) {
				isHex = false
				break
			}
		}
		if !isHex {
			continue
		}
		upper := toUpperHex(hex)
		if upper == "#000000" || upper == "#FFFFFF" {
			continue
		}
		seen[upper] = struct{}{}
	}
	return len(seen)
}

// isHexDigit reports whether r is a valid hexadecimal digit character.
func isHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func toUpperHex(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'f' {
			b[i] = c - 'a' + 'A'
		}
	}
	return string(b)
}

func samplePoints() []Point {
	return []Point{
		{Recall: 0.5, QPS: 100},
		{Recall: 0.8, QPS: 50},
		{Recall: 0.95, QPS: 10},
	}
}

// --- RenderParetoSVG -----------------------------------------------------------

func TestRenderParetoSVG_SingleSeries(t *testing.T) {
	t.Parallel()

	cfg := Config{
		Title:  "Recall-QPS Pareto",
		XLabel: "Recall",
		YLabel: "Achieved QPS",
		Width:  800,
		Height: 600,
	}
	series := []Series{
		{Label: "glove-100-angular-K10", Points: samplePoints()},
	}

	res, err := RenderParetoSVG(series, cfg)
	if err != nil {
		t.Fatalf("RenderParetoSVG() error = %v", err)
	}
	if res == nil {
		t.Fatal("RenderParetoSVG() returned nil result with nil error")
	}

	root := mustBeWellFormedSVG(t, res.SVG)
	if root != "svg" {
		t.Errorf("root element = %q, want %q", root, "svg")
	}

	if res.SeriesCount != 1 {
		t.Errorf("SeriesCount = %d, want 1", res.SeriesCount)
	}
	if got := res.PointCount["glove-100-angular-K10"]; got != len(series[0].Points) {
		t.Errorf("PointCount[label] = %d, want %d", got, len(series[0].Points))
	}

	for _, want := range []string{cfg.Title, cfg.XLabel, cfg.YLabel, "glove-100-angular-K10"} {
		if !bytes.Contains(res.SVG, []byte(want)) {
			t.Errorf("SVG output does not contain expected text %q", want)
		}
	}
}

func TestRenderParetoSVG_MultipleSeries_DistinctColors(t *testing.T) {
	t.Parallel()

	series := []Series{
		{Label: "dataset-A-K1", Points: samplePoints()},
		{Label: "dataset-B-K10", Points: []Point{
			{Recall: 0.6, QPS: 200},
			{Recall: 0.9, QPS: 20},
		}},
	}
	cfg := DefaultConfig()

	res, err := RenderParetoSVG(series, cfg)
	if err != nil {
		t.Fatalf("RenderParetoSVG() error = %v", err)
	}

	mustBeWellFormedSVG(t, res.SVG)

	if res.SeriesCount != 2 {
		t.Errorf("SeriesCount = %d, want 2", res.SeriesCount)
	}
	wantTotalPoints := len(series[0].Points) + len(series[1].Points)
	gotTotalPoints := 0
	for _, n := range res.PointCount {
		gotTotalPoints += n
	}
	if gotTotalPoints != wantTotalPoints {
		t.Errorf("sum(PointCount) = %d, want %d", gotTotalPoints, wantTotalPoints)
	}

	for _, s := range series {
		if !bytes.Contains(res.SVG, []byte(s.Label)) {
			t.Errorf("SVG legend is missing label %q", s.Label)
		}
	}

	if n := countDistinctColors(res.SVG); n < 2 {
		t.Errorf("countDistinctColors() = %d, want >= 2 (one per series, excluding axis black/background white)", n)
	}
}

func TestRenderParetoSVG_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("empty series list is an error", func(t *testing.T) {
		t.Parallel()
		if _, err := RenderParetoSVG(nil, DefaultConfig()); err == nil {
			t.Error("expected error for nil series, got nil")
		}
		if _, err := RenderParetoSVG([]Series{}, DefaultConfig()); err == nil {
			t.Error("expected error for empty series slice, got nil")
		}
	})

	t.Run("series with only empty Points is an error", func(t *testing.T) {
		t.Parallel()
		series := []Series{{Label: "empty", Points: nil}}
		if _, err := RenderParetoSVG(series, DefaultConfig()); err == nil {
			t.Error("expected error when no series has any point, got nil")
		}
	})

	t.Run("non-finite or non-positive QPS points are filtered, not fatal", func(t *testing.T) {
		t.Parallel()
		series := []Series{
			{
				Label: "mixed",
				Points: []Point{
					{Recall: 0.5, QPS: 100},         // valid
					{Recall: 0.6, QPS: 0},           // invalid: log-scale can't plot 0
					{Recall: 0.7, QPS: -5},          // invalid: negative
					{Recall: math.NaN(), QPS: 10},   // invalid: NaN recall
					{Recall: 0.8, QPS: math.Inf(1)}, // invalid: +Inf QPS
					{Recall: 0.9, QPS: 42},          // valid
				},
			},
		}
		res, err := RenderParetoSVG(series, DefaultConfig())
		if err != nil {
			t.Fatalf("RenderParetoSVG() error = %v, want the valid points to still render", err)
		}
		if res.SeriesCount != 1 {
			t.Errorf("SeriesCount = %d, want 1", res.SeriesCount)
		}
		if got := res.PointCount["mixed"]; got != 2 {
			t.Errorf("PointCount[mixed] = %d, want 2 (only the two finite, QPS>0 points)", got)
		}
	})

	t.Run("a single data point across a single series does not panic (degenerate color-step case)", func(t *testing.T) {
		t.Parallel()
		series := []Series{{Label: "solo", Points: []Point{{Recall: 0.5, QPS: 1}}}}
		res, err := RenderParetoSVG(series, DefaultConfig())
		if err != nil {
			t.Fatalf("RenderParetoSVG() error = %v", err)
		}
		if res.SeriesCount != 1 || res.PointCount["solo"] != 1 {
			t.Errorf("unexpected result for single point/single series: %+v", res)
		}
	})

	t.Run("non-positive Width/Height falls back to DefaultConfig dimensions instead of erroring", func(t *testing.T) {
		t.Parallel()
		series := []Series{{Label: "s", Points: samplePoints()}}
		res, err := RenderParetoSVG(series, Config{Title: "t", XLabel: "x", YLabel: "y", Width: 0, Height: -1})
		if err != nil {
			t.Fatalf("RenderParetoSVG() error = %v, want graceful fallback to default dimensions", err)
		}
		if len(res.SVG) == 0 {
			t.Error("expected non-empty SVG output with fallback dimensions")
		}
	})
}

// --- LoadParetoSeries ------------------------------------------------------------

func writeSnapshotJSON(t *testing.T, dir, name string, snap *metrics.GlobalSnapshot) string {
	t.Helper()
	data, err := json.Marshal(snap)
	if err != nil {
		t.Fatalf("json.Marshal(snapshot) error = %v", err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("os.WriteFile(%s) error = %v", path, err)
	}
	return path
}

// buildSnapshot constructs a GlobalSnapshot the way a real Collector would (via
// metrics.NewCollector + RecordRecall + Record), rather than a bare struct
// literal, so this test also exercises the real recall/AchievedQPS wiring
// end-to-end together with the recall_test.go / snapshot_qps_test.go contracts.
func buildSnapshot(
	t *testing.T, recall float64, total int, window time.Duration,
) *metrics.GlobalSnapshot {
	t.Helper()
	c, err := metrics.NewCollector()
	if err != nil {
		t.Fatalf("metrics.NewCollector() error = %v", err)
	}
	start := time.Now().Add(-window)
	for i := 0; i < total; i++ {
		c.Record(t.Context(), 0, &metrics.RequestResult{
			StartedAt: start,
			EndedAt:   start.Add(time.Millisecond),
		})
	}
	c.RecordRecall(recall)
	return c.GlobalSnapshot()
}

func TestLoadParetoSeries(t *testing.T) {
	t.Parallel()

	t.Run("loads one Point per JSON file, labeled as requested", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		p1 := writeSnapshotJSON(t, dir, "run1.json", buildSnapshot(t, 0.9, 100, 10*time.Second))
		p2 := writeSnapshotJSON(t, dir, "run2.json", buildSnapshot(t, 0.7, 50, 5*time.Second))

		series, err := LoadParetoSeries("my-label", p1, p2)
		if err != nil {
			t.Fatalf("LoadParetoSeries() error = %v", err)
		}
		if series.Label != "my-label" {
			t.Errorf("series.Label = %q, want %q", series.Label, "my-label")
		}
		if len(series.Points) != 2 {
			t.Fatalf("len(series.Points) = %d, want 2", len(series.Points))
		}
		for _, p := range series.Points {
			if p.Recall <= 0 || p.Recall > 1 {
				t.Errorf("Point.Recall = %v, want in (0,1]", p.Recall)
			}
			if p.QPS <= 0 {
				t.Errorf("Point.QPS = %v, want > 0", p.QPS)
			}
		}
	})

	t.Run("no paths is an error", func(t *testing.T) {
		t.Parallel()
		if _, err := LoadParetoSeries("label"); err == nil {
			t.Error("expected error for zero jsonPaths, got nil")
		}
	})

	t.Run("nonexistent file returns a wrapped error, not a panic", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("LoadParetoSeries panicked on a missing file: %v", r)
			}
		}()
		if _, err := LoadParetoSeries("label", filepath.Join(t.TempDir(), "does-not-exist.json")); err == nil {
			t.Error("expected error for a nonexistent file, got nil")
		}
	})

	t.Run("malformed JSON returns an error, not a panic", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := filepath.Join(dir, "bad.json")
		if err := os.WriteFile(path, []byte("{not valid json"), 0o600); err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}
		if _, err := LoadParetoSeries("label", path); err == nil {
			t.Error("expected error for malformed JSON, got nil")
		}
	})

	t.Run("snapshot with no recorded recall data is an error naming the offending file", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		c, err := metrics.NewCollector()
		if err != nil {
			t.Fatalf("metrics.NewCollector() error = %v", err)
		}
		c.Record(t.Context(), 0, &metrics.RequestResult{Latency: time.Millisecond})
		path := writeSnapshotJSON(t, dir, "no-recall.json", c.GlobalSnapshot())

		_, err = LoadParetoSeries("label", path)
		if err == nil {
			t.Fatal("expected error for a snapshot with no recall data, got nil")
		}
		if !strings.Contains(err.Error(), filepath.Base(path)) {
			t.Errorf("error %q does not mention the offending file %q", err.Error(), filepath.Base(path))
		}
	})
}

// TestLoadParetoSeries_RenderParetoSVG_Integration wires LoadParetoSeries and
// RenderParetoSVG together across two labeled groups of snapshot files, matching
// the end-to-end usage the task describes: several GlobalSnapshot JSON files (+ a
// label) in, one combined multi-color Pareto SVG out.
func TestLoadParetoSeries_RenderParetoSVG_Integration(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	aPaths := []string{
		writeSnapshotJSON(t, dir, "a1.json", buildSnapshot(t, 0.95, 100, 20*time.Second)),
		writeSnapshotJSON(t, dir, "a2.json", buildSnapshot(t, 0.80, 100, 8*time.Second)),
	}
	bPaths := []string{
		writeSnapshotJSON(t, dir, "b1.json", buildSnapshot(t, 0.60, 100, 4*time.Second)),
	}

	seriesA, err := LoadParetoSeries("ngt-k10", aPaths...)
	if err != nil {
		t.Fatalf("LoadParetoSeries(ngt-k10) error = %v", err)
	}
	seriesB, err := LoadParetoSeries("hnsw-k10", bPaths...)
	if err != nil {
		t.Fatalf("LoadParetoSeries(hnsw-k10) error = %v", err)
	}

	res, err := RenderParetoSVG([]Series{seriesA, seriesB}, DefaultConfig())
	if err != nil {
		t.Fatalf("RenderParetoSVG() error = %v", err)
	}
	mustBeWellFormedSVG(t, res.SVG)

	if res.SeriesCount != 2 {
		t.Errorf("SeriesCount = %d, want 2", res.SeriesCount)
	}
	if res.PointCount["ngt-k10"] != 2 {
		t.Errorf(`PointCount["ngt-k10"] = %d, want 2`, res.PointCount["ngt-k10"])
	}
	if res.PointCount["hnsw-k10"] != 1 {
		t.Errorf(`PointCount["hnsw-k10"] = %d, want 1`, res.PointCount["hnsw-k10"])
	}
	if n := countDistinctColors(res.SVG); n < 2 {
		t.Errorf("countDistinctColors() = %d, want >= 2", n)
	}
}
