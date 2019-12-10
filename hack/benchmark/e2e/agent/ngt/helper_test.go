package ngt

import (
	"flag"
	"strings"
	"sync"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/log"
)

const (
	num = 10
	radius = -1.0
	epsilon = 0.01
)

var (
	searchConfig = &payload.Search_Config{
		Num:     num,
		Radius:  radius,
		Epsilon: epsilon,
	}
	targets    []string
	addresses  []string
	datasetVar string
	addressVar string
	numVar     uint
	radiusVar  float64
	epsilonVar float64
	once       sync.Once
)

func init() {
	log.Init(log.DefaultGlg())

	flag.StringVar(&datasetVar, "dataset", "", "available dataset(choice with comma)")
	flag.StringVar(&addressVar, "address", "", "vald agent address")
	flag.UintVar(&numVar, "num", num, "search response size")
	flag.Float64Var(&radiusVar, "radius", radius, "search radius size")
	flag.Float64Var(&epsilonVar, "epsilon", epsilon, "search epsilon size")
}

func parseArgs(tb testing.TB) {
	tb.Helper()
	once.Do(func() {
		flag.Parse()
		targets = strings.Split(strings.TrimSpace(datasetVar), ",")
		addresses = strings.Split(strings.TrimSpace(addressVar), ",")
		if len(targets) != len(addresses) {
			tb.Fatal("address and dataset must have same length.")
		}
		searchConfig = &payload.Search_Config{
			Num:     uint32(numVar),
			Radius:  float32(radiusVar),
			Epsilon: float32(epsilonVar),
		}
	})
}
