package main

import (
	"flag"
	"strings"

	"github.com/vdaas/vald/hack/core/ngt/assets"
	"github.com/vdaas/vald/internal/core/ngt"
	"github.com/vdaas/vald/internal/log"
)

var (
	distanceType = flag.String("distance-type", "L2", "choice indexing distance type")
)

func main() {
	log.Init(log.DefaultGlg())

	flag.Parse()
	distanceTypeOption := ngt.L2
	switch strings.ToLower(*distanceType) {
	case "l1":
		distanceTypeOption = ngt.L1
	case "l2":
		distanceTypeOption = ngt.L2
	case "angle":
		distanceTypeOption = ngt.Angle
	case "hamming":
		distanceTypeOption = ngt.Hamming
	case "cos":
		distanceTypeOption = ngt.Cosine
	}
	indexName := flag.Arg(0)

	d, err := assets.LoadDataset(flag.Arg(1))
	if err != nil {
		log.Warn(err)
		return
	}
	defer log.Infof("[%s] Done.", indexName)
	opts := []ngt.Option{
		ngt.WithDistanceType(distanceTypeOption),
	}
	loader := func() ([][]float64, error) {
		vectors, err := d.LoadTrain()
		if err != nil {
			return nil, err
		}
		defer log.Infof("[%s] %d items", d.Path, len(vectors))
		return vectors, nil
	}
	if err := assets.CreateIndex(indexName, loader, opts...); err != nil {
		log.Warn(err)
		return
	}
}
