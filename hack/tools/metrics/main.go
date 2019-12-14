package main

import (
	"os"

	"github.com/vdaas/vald/internal/log"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
)

var (
	title string
	xLabel string
	yLabel string
	output string
)

func main() {
	log.Init(log.DefaultGlg())

	p, err := plot.New()
	if err != nil {
		log.Error(err)
	}
	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.X.Max = 1.0
	p.X.Min = 0.0
	p.Y.Label.Text = yLabel
	p.Y.Scale = plot.LogScale{}

	points := make(plotter.XYs, m.Len())
	for i := 0; i < m.Len(); i++ {
		points[i].X = m.Recall[i]
		points[i].Y = m.Qps[i]
	}
	if err := plotutil.AddLinePoints(p, points); err != nil {
		log.Error(err)
	}

	canvas := vgsvg.New(1280, 960)
	p.Draw(draw.New(canvas))
	out, err := os.OpenFile(output, os.O_CREATE | os.O_WRONLY | os.O_TRUNC, os.ModeTemporary)
	if err != nil {
		log.Error(err)
	}
	defer func() {
		err := out.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	_, err = canvas.WriteTo(out)
	if err != nil {
		log.Error(err)
	}
}
