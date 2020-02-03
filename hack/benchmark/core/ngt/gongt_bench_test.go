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

// Package ngt_test provides implementation of Go API for https://github.com/yahoojapan/NGT
package ngt

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/log"
	"github.com/yahoojapan/gongt"
)

const (
	size    = 10
	epsilon = 0.1
)

var (
	targets    []string
	datasetVar string
	once       sync.Once
)

func init() {
	log.Init()

	flag.StringVar(&datasetVar, "dataset", "", "available dataset(choice with comma)")
}

func parseArgs(tb testing.TB) {
	tb.Helper()

	once.Do(func() {
		flag.Parse()
		targets = strings.Split(strings.TrimSpace(datasetVar), ",")
	})
}

func BenchmarkGoNGTSequential(b *testing.B) {
	parseArgs(b)

	for _, target := range targets {
		b.Run(target, func(bb *testing.B) {
			tmpdir, err := ioutil.TempDir("", "tmpdir")
			if err != nil {
				bb.Error(err)
			}
			defer os.RemoveAll(tmpdir)

			d := assets.Data(target)(b)
			n := gongt.SetIndexPath(tmpdir).SetObjectType(gongt.Float).SetDimension(d.Dimension()).Open()
			defer n.Close()

			bb.Run("Insert", func(sb *testing.B) {
				dataset := d.TrainAsFloat64()

				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				for i := 0; i < sb.N; i++ {
					_, err := n.Insert(dataset[i%len(dataset)])
					if err != nil {
						sb.Error(err)
					}
				}
				sb.StopTimer()
			})

			bb.Run("CreateIndex", func(sb *testing.B) {
				err := n.CreateIndex(10000)
				if err != nil {
					sb.Error(err)
				}
			})

			bb.Run("Search", func(sb *testing.B) {
				dataset := d.QueryAsFloat64()

				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				for i := 0; i < sb.N; i++ {
					_, err := n.Search(dataset[i%len(dataset)], size, epsilon)
					if err != nil {
						sb.Error(err)
					}
				}
				sb.StopTimer()
			})
		})
	}
}

func BenchmarkGoNGTParallel(b *testing.B) {
	parseArgs(b)

	for _, target := range targets {
		b.Run(target, func(bb *testing.B) {
			tmpdir, err := ioutil.TempDir("", "tmpdir")
			if err != nil {
				bb.Error(err)
			}
			defer os.RemoveAll(tmpdir)

			d := assets.Data(target)(b)
			n := gongt.SetIndexPath(tmpdir).SetObjectType(gongt.Float).SetDimension(d.Dimension()).Open()
			defer n.Close()

			bb.Run("Insert", func(sb *testing.B) {
				dataset := d.TrainAsFloat64()
				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				sb.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						_, err = n.Insert(dataset[i%len(dataset)])
						if err != nil {
							sb.Error(err)
						}
						i++
					}
				})
				sb.StopTimer()
			})

			bb.Run("CreateIndex", func(sb *testing.B) {
				err := n.CreateIndex(10000)
				if err != nil {
					sb.Error(err)
				}
			})

			bb.Run("Search", func(sb *testing.B) {
				dataset := d.QueryAsFloat64()
				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				sb.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						_, err := n.Search(dataset[i%len(dataset)], size, epsilon)
						if err != nil {
							sb.Error(err)
						}
						i++
					}
				})
				sb.StopTimer()
			})
		})
	}
}
