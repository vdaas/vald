package ngt

import (
	context "context"
	"flag"
	"strings"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/core/benchmark"
	"github.com/vdaas/vald/hack/benchmark/core/benchmark/strategy"
	"github.com/vdaas/vald/internal/log"
)

var (
	targets []string
)

func init() {
	testing.Init()
	log.Init()

	var dataset string
	flag.StringVar(&dataset, "dataset", "", "available dataset(choice with comma)")
	flag.Parse()

	targets = strings.Split(strings.TrimSpace(dataset), ",")
}

func BenchmarkNGT(b *testing.B) {
	for _, target := range targets {
		benchmark.New(
			b,
			benchmark.WithName(target),
			benchmark.WithNGT(nil),
			benchmark.WithStrategy(
				strategy.NewInsert(),
			),
		).Run(context.Background(), b)
	}
}
