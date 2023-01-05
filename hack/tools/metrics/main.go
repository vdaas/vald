// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"image/color"
	"io/fs"
	"os"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/vdaas/vald/hack/benchmark/metrics"
	"github.com/vdaas/vald/internal/log"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
)

var (
	title  = flag.String("title", "metrics", "metrics chart title")
	xLabel = flag.String("x", "x", "x axis label")
	yLabel = flag.String("y", "y", "y axis label")
	input  = flag.String("input", "metrics.gob", "input gob file path")
	output = flag.String("output", "chart.svg", "output chart file path")
	width  = flag.Int("width", 1280, "chart image width")
	height = flag.Int("height", 960, "chart image height")
)

func main() {
	flag.Parse()

	log.Init()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	in, err := os.OpenFile(*input, os.O_RDONLY, fs.ModeDevice)
	if err != nil {
		return err
	}
	defer func() {
		if err := in.Close(); err != nil {
			log.Error(err)
		}
	}()

	var ms []metrics.Metrics
	if err := gob.NewDecoder(in).Decode(&ms); err != nil {
		return err
	}

	p := plot.New()
	p.Title.Text = *title
	p.X.Label.Text = *xLabel
	p.X.Max = 1.0
	p.X.Min = 0.0
	p.Y.Label.Text = *yLabel
	p.Y.Tick.Marker = plot.LogTicks{}
	p.Y.Scale = plot.LogScale{}
	p.Add(plotter.NewGrid())

	min := 0.0
	max := 270.0
	var step float64
	if len(ms) == 1 {
		step = 0
	} else {
		step = 1 / float64(len(ms)-1)
	}
	for i, m := range ms {
		sort.Slice(m.Search, func(i, j int) bool {
			return m.Search[i].Recall < m.Search[j].Recall
		})
		xys := make(plotter.XYs, len(m.Search))
		for i, s := range m.Search {
			xys[i].X = s.Recall
			xys[i].Y = s.Qps
		}
		points, _, err := plotter.NewLinePoints(xys)
		if err != nil {
			log.Error(err)
		}
		r, g, b := colorful.Hcl((max-min)*float64(i)*step+min, 0.6, 0.8).RGB255()
		points.Color = color.RGBA{R: r, G: g, B: b, A: 255}
		p.Add(points)
		p.Legend.Add(fmt.Sprintf("%s-Recall@%d", m.DatasetName, m.K), points)
	}

	canvas := vgsvg.New(vg.Length(*width), vg.Length(*height))
	p.Draw(draw.New(canvas))
	out, err := os.OpenFile(*output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fs.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		err := out.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	_, err = canvas.WriteTo(out)
	if err != nil {
		return err
	}

	return nil
}
