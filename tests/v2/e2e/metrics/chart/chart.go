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

// Package chart renders Recall-QPS Pareto charts (SVG) for
// tests/v2/e2e/metrics.GlobalSnapshot data.
//
// It is deliberately independent from hack/tools/metrics and
// hack/benchmark/metrics: it is scoped to GlobalSnapshot (recall +
// AchievedQPS) rather than the hack/benchmark-specific gob-encoded
// Metrics/SearchMetrics types, but reuses the same rendering approach
// (gonum/plot scatter, log-scale Y axis, HCL-spaced colors via go-colorful -
// see hack/tools/metrics/main.go for the reference implementation this
// design is modeled after).
package chart

import (
	"bytes"
	"image/color"
	"math"
	"path/filepath"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/tests/v2/e2e/metrics"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
)

const (
	defaultWidth  = 1280
	defaultHeight = 960

	// hueMax/saturation/luminance are the HCL parameters used to spread
	// per-series colors, mirroring hack/tools/metrics/main.go.
	hueMax     = 360.0
	saturation = 0.6
	luminance  = 0.6

	// fullAlpha is the fully-opaque alpha channel value for color.RGBA.
	fullAlpha = 255

	// degenerateYAxisSpanFactor widens a single-value (min == max) Y range
	// by this factor on each side so plot.LogTicks has a non-degenerate,
	// strictly-positive span to derive ticks from (see the comment above
	// its use in RenderParetoSVG).
	degenerateYAxisSpanFactor = 10.0
)

// Point is a single (Recall, AchievedQPS) sample plotted on a Pareto chart.
type Point struct {
	// Recall is the x-axis value, expected in [0,1] but not enforced here.
	Recall float64
	// QPS is the y-axis value, plotted on a log scale - must be > 0 to be
	// plottable.
	QPS float64
}

// Series is a named group of Points, rendered with its own color/legend
// entry.
type Series struct {
	Label  string
	Points []Point
}

// Config configures RenderParetoSVG's output.
type Config struct {
	Title, XLabel, YLabel string
	// Width/Height are in pixels; <= 0 falls back to DefaultConfig's values.
	Width, Height int
}

// DefaultConfig returns the default chart rendering configuration.
func DefaultConfig() Config {
	return Config{
		Title:  "Recall-QPS Pareto",
		XLabel: "Recall",
		YLabel: "Achieved QPS",
		Width:  defaultWidth,
		Height: defaultHeight,
	}
}

// RenderResult is the outcome of a successful RenderParetoSVG call.
type RenderResult struct {
	// PointCount is the per-label count of rendered (finite, QPS>0) points.
	PointCount map[string]int
	SVG        []byte
	// SeriesCount is the number of series that contributed >=1 rendered
	// point.
	SeriesCount int
}

// pointRenderable reports whether p can be plotted on a log-scale Y axis:
// both coordinates must be finite, and QPS must be strictly positive (a
// log-scale axis cannot represent zero or negative values).
func pointRenderable(p Point) bool {
	if math.IsNaN(p.Recall) || math.IsInf(p.Recall, 0) {
		return false
	}
	if math.IsNaN(p.QPS) || math.IsInf(p.QPS, 0) {
		return false
	}
	return p.QPS > 0
}

// RenderParetoSVG renders one or more Series onto a single Recall (x) vs
// Achieved QPS (y, log-scale) scatter chart and returns the resulting SVG
// document along with render-time ground-truth counts (see package doc for
// why counts are returned rather than re-derived from the SVG).
//
// Points that are not finite or have QPS <= 0 are silently filtered out (a
// log-scale Y axis cannot represent them) rather than causing an error;
// RenderParetoSVG only errors if there is nothing at all to render.
func RenderParetoSVG(series []Series, cfg Config) (*RenderResult, error) {
	if len(series) == 0 {
		return nil, errors.New("chart: RenderParetoSVG requires at least one series")
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		def := DefaultConfig()
		if cfg.Width <= 0 {
			cfg.Width = def.Width
		}
		if cfg.Height <= 0 {
			cfg.Height = def.Height
		}
	}

	p := plot.New()
	p.Title.Text = cfg.Title
	p.X.Label.Text = cfg.XLabel
	p.Y.Label.Text = cfg.YLabel
	p.Y.Scale = plot.LogScale{}
	p.Y.Tick.Marker = plot.LogTicks{}
	p.Add(plotter.NewGrid())

	result := &RenderResult{PointCount: make(map[string]int, len(series))}

	n := len(series)
	var hueStep float64
	if n > 1 {
		hueStep = 1.0 / float64(n-1)
	}
	minY, maxY := math.Inf(1), math.Inf(-1)
	for i, s := range series {
		xys := make(plotter.XYs, 0, len(s.Points))
		for _, pt := range s.Points {
			if !pointRenderable(pt) {
				continue
			}
			xys = append(xys, plotter.XY{X: pt.Recall, Y: pt.QPS})
			if pt.QPS < minY {
				minY = pt.QPS
			}
			if pt.QPS > maxY {
				maxY = pt.QPS
			}
		}
		if len(xys) == 0 {
			continue
		}

		sc, err := plotter.NewScatter(xys)
		if err != nil {
			return nil, errors.Wrapf(err, "chart: failed to build scatter plot for series %q", s.Label)
		}
		hue := 0.0
		if n > 1 {
			hue = hueMax * float64(i) * hueStep
		}
		r, g, b := colorful.Hcl(hue, saturation, luminance).Clamped().RGB255()
		sc.Color = color.RGBA{R: r, G: g, B: b, A: fullAlpha}
		p.Add(sc)
		p.Legend.Add(s.Label, sc)

		result.PointCount[s.Label] = len(xys)
		result.SeriesCount++
	}
	if result.SeriesCount == 0 {
		return nil, errors.New("chart: no series had any renderable point (all points were non-finite or QPS <= 0)")
	}

	// plot.Plot.Add derives the Y axis range via
	// `p.Y.Min = math.Min(p.Y.Min, ymin)`, but Axis.Min's zero value is 0,
	// so the computed range's lower bound never rises above 0 even when
	// every rendered QPS value is strictly positive. That is fatal for a
	// log-scale axis (plot.LogTicks.Ticks panics on min <= 0), so the range
	// is set explicitly here from the actual rendered points instead of
	// relying on plot's auto-ranging. A degenerate (single-value) range is
	// widened symmetrically in log-space so LogTicks has a real span to
	// derive ticks from.
	if minY == maxY {
		minY /= degenerateYAxisSpanFactor
		maxY *= degenerateYAxisSpanFactor
	}
	p.Y.Min = minY
	p.Y.Max = maxY

	canvas := vgsvg.New(vg.Length(cfg.Width), vg.Length(cfg.Height))
	p.Draw(draw.New(canvas))
	var buf bytes.Buffer
	if _, err := canvas.WriteTo(&buf); err != nil {
		return nil, errors.Wrap(err, "chart: failed to encode SVG")
	}
	result.SVG = buf.Bytes()
	return result, nil
}

// LoadParetoSeries loads one Point per GlobalSnapshot JSON file (as produced
// by tests/v2/e2e/metrics.SnapshotPresenter.AsJSON / json.Marshal), using the
// snapshot's recorded recall mean as Point.Recall and its
// (*GlobalSnapshot).AchievedQPS() as Point.QPS, and returns them as a single
// Series labeled with label.
func LoadParetoSeries(label string, jsonPaths ...string) (Series, error) {
	if len(jsonPaths) == 0 {
		return Series{}, errors.New("chart: LoadParetoSeries requires at least one jsonPath")
	}

	series := Series{Label: label, Points: make([]Point, 0, len(jsonPaths))}
	for _, path := range jsonPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			return Series{}, errors.Wrapf(err, "chart: failed to read snapshot file %s", path)
		}

		var snap metrics.GlobalSnapshot
		if err := json.Unmarshal(data, &snap); err != nil {
			return Series{}, errors.Wrapf(err, "chart: failed to parse snapshot JSON in %s", filepath.Base(path))
		}

		if snap.Recalls == nil || snap.Recalls.Total == 0 {
			return Series{}, errors.Errorf("chart: snapshot %s has no recorded recall data", filepath.Base(path))
		}

		series.Points = append(series.Points, Point{
			Recall: snap.Recalls.Mean,
			QPS:    snap.AchievedQPS(),
		})
	}
	return series, nil
}
