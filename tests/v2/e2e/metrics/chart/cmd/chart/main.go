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

// Command chart renders a Recall-QPS Pareto SVG chart from one or more
// tests/v2/e2e/metrics.GlobalSnapshot JSON files (as produced by
// metrics.SnapshotPresenter.AsJSON / json.Marshal), grouping them into
// labeled series via repeatable -series flags:
//
//	chart \
//	  -series "glove-100-angular-K10=run1.json,run2.json" \
//	  -series "sift-128-euclidean-K10=run3.json" \
//	  -title "Recall-QPS Pareto" -x Recall -y "Achieved QPS" \
//	  -width 1280 -height 960 \
//	  -output pareto.svg
//
// It is the CLI counterpart to the tests/v2/e2e/metrics/chart library
// (LoadParetoSeries/RenderParetoSVG), modeled after hack/tools/metrics/main.go
// but designed for multiple labeled series (rather than a single gob-encoded
// metrics file) since a Pareto comparison is inherently multi-series.
package main

import (
	"flag"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/v2/e2e/metrics/chart"
)

// outputFileMode is the permission used for the rendered SVG file. 0o600
// (owner read/write only) satisfies gosec's G306 and matches the permission
// already used for GlobalSnapshot JSON fixtures in chart_test.go.
const outputFileMode = 0o600

// seriesSpec pairs a chart series label with the GlobalSnapshot JSON file
// paths that belong to it.
type seriesSpec struct {
	label string
	paths []string
}

// seriesSpecs implements flag.Value so -series can be repeated on the
// command line, one occurrence per labeled group of input files, e.g.:
//
//	-series "glove-100-angular-K10=run1.json,run2.json" \
//	-series "sift-128-euclidean-K10=run3.json"
type seriesSpecs []seriesSpec

// String implements flag.Value.
func (s *seriesSpecs) String() string {
	if s == nil {
		return ""
	}
	parts := make([]string, 0, len(*s))
	for _, spec := range *s {
		parts = append(parts, spec.label+"="+strings.Join(spec.paths, ","))
	}
	return strings.Join(parts, " ")
}

// Set implements flag.Value. It parses a single "label=path1,path2,..."
// occurrence and appends it, so -series may be passed multiple times to
// build up multiple labeled series.
func (s *seriesSpecs) Set(value string) error {
	label, rawPaths, ok := strings.Cut(value, "=")
	label = strings.TrimSpace(label)
	if !ok || label == "" {
		return errors.Errorf("invalid -series value %q, want format label=path1,path2,...", value)
	}
	var paths []string
	for _, p := range strings.Split(rawPaths, ",") {
		if p = strings.TrimSpace(p); p != "" {
			paths = append(paths, p)
		}
	}
	if len(paths) == 0 {
		return errors.Errorf("invalid -series value %q: no non-empty paths given for label %q", value, label)
	}
	*s = append(*s, seriesSpec{label: label, paths: paths})
	return nil
}

// nolint:gochecknoglobals // flag bindings in main.go are exempted, see .golangci.json
var (
	specs  seriesSpecs
	title  = flag.String("title", chart.DefaultConfig().Title, "chart title")
	xLabel = flag.String("x", chart.DefaultConfig().XLabel, "X axis label")
	yLabel = flag.String("y", chart.DefaultConfig().YLabel, "Y axis label")
	width  = flag.Int("width", chart.DefaultConfig().Width, "chart image width in pixels")
	height = flag.Int("height", chart.DefaultConfig().Height, "chart image height in pixels")
	output = flag.String("output", "pareto.svg", "output SVG file path")
)

func main() {
	log.Init()
	flag.Var(&specs, "series",
		"labeled group of GlobalSnapshot JSON files to render as one series, "+
			"format: label=path1,path2,... (repeatable, one flag occurrence per series)")
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
		return
	}
}

// run performs the actual load -> render -> write pipeline and returns an
// error instead of exiting directly, so main can be the single point that
// terminates the process (Vald Law: no panic/log.Fatal outside the
// top-level entry point boundary).
func run() error {
	if len(specs) == 0 {
		return errors.New("at least one -series label=path1,path2,... flag is required")
	}

	series := make([]chart.Series, 0, len(specs))
	for _, spec := range specs {
		s, err := chart.LoadParetoSeries(spec.label, spec.paths...)
		if err != nil {
			return errors.Wrapf(err, "failed to load series %q", spec.label)
		}
		series = append(series, s)
	}

	res, err := chart.RenderParetoSVG(series, chart.Config{
		Title:  *title,
		XLabel: *xLabel,
		YLabel: *yLabel,
		Width:  *width,
		Height: *height,
	})
	if err != nil {
		return errors.Wrap(err, "failed to render Pareto SVG")
	}

	if err := os.WriteFile(*output, res.SVG, outputFileMode); err != nil {
		return errors.Wrapf(err, "failed to write output SVG to %s", *output)
	}

	log.Infof("rendered %d series to %s (points per series: %v)", res.SeriesCount, *output, res.PointCount)
	return nil
}
