//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	num     = 10
	radius  = -1.0
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
	log.Init()

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
