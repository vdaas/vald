package ngt

import (
	"flag"
	"strings"
	"testing"
)

var targets []string

func init() {
	testing.Init()

	var dataset string
	flag.StringVar(&dataset, "dataset", "", "available dataset(choice with comma)")
	flag.Parse()
	targets = strings.Split(strings.TrimSpace(dataset), ",")
}

func BenchmarkNGT_Insert(b *testing.B) {
}
