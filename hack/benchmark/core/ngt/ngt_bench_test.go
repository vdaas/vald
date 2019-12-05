//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package ngt_test provides implementation of Go API for https://github.com/yahoojapan/NGT
package ngt_test

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
	"github.com/vdaas/vald/internal/log"
)

const (
	size    = 10
	epsilon = 0.1
	radius  = -1
)

var (
	targets    []string
	datasetVar string
	once       sync.Once
)

func init() {
	log.Init(log.DefaultGlg())

	flag.StringVar(&datasetVar, "dataset", "", "available dataset(choice with comma)")
}

func parseArgs(tb testing.TB) {
	once.Do(func() {
		flag.Parse()
		targets = strings.Split(strings.TrimSpace(datasetVar), ",")
	})
}

func BenchmarkNGTSequential(b *testing.B) {
	parseArgs(b)

	for _, target := range targets {
		b.Run(target, func(bb *testing.B) {
			d := assets.Data(target)(b)
			tmpdir, err := ioutil.TempDir("", "tmpdir")
			if err != nil {
				bb.Error(err)
			}
			defer os.RemoveAll(tmpdir)

			n, err := ngt.New(
				ngt.WithIndexPath(tmpdir),
				ngt.WithObjectType(ngt.Float),
				ngt.WithDimension(d.Dimension()),
			)
			if err != nil {
				bb.Error(err)
			}
			// defer n.Close()

			bb.Run("Insert", func(sb *testing.B) {
				train := d.Train()
				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				for i := 0; i < sb.N; i++ {
					_, err = n.Insert(train[i%len(train)])
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
				query := d.Query()
				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				for i := 0; i < sb.N; i++ {
					_, err := n.Search(query[i%len(query)], size, epsilon, radius)
					if err != nil {
						sb.Error(err)
					}
				}
				sb.StopTimer()
			})
		})
	}
}

func BenchmarkNGTSequentialBulk(b *testing.B) {
	parseArgs(b)

	for _, target := range targets {
		b.Run(target, func(bb *testing.B) {
			d := assets.Data(target)(b)
			tmpdir, err := ioutil.TempDir("", "tmpdir")
			if err != nil {
				bb.Error(err)
			}
			defer os.RemoveAll(tmpdir)

			n, err := ngt.New(
				ngt.WithIndexPath(tmpdir),
				ngt.WithObjectType(ngt.Float),
				ngt.WithDimension(d.Dimension()),
			)
			if err != nil {
				bb.Error(err)
			}
			// defer n.Close()

			bb.Run("BulkInsert", func(sb *testing.B) {
				train := d.Train()
				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				for i := 0; i < sb.N; i++ {
					_, err = n.BulkInsert(train)
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
				query := d.Query()
				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				for i := 0; i < sb.N; i++ {
					_, err := n.Search(query[i%len(query)], size, epsilon, radius)
					if err != nil {
						sb.Error(err)
					}
				}
				sb.StopTimer()
			})
		})
	}
}

func BenchmarkNGTParallel(b *testing.B) {
	parseArgs(b)

	for _, target := range targets {
		d := assets.Data(target)(b)
		b.Run(target, func(bb *testing.B) {
			tmpdir, err := ioutil.TempDir("", "tmpdir")
			if err != nil {
				bb.Error(err)
			}
			defer os.RemoveAll(tmpdir)

			n, err := ngt.New(
				ngt.WithIndexPath(tmpdir),
				ngt.WithObjectType(ngt.Float),
				ngt.WithDimension(d.Dimension()),
			)
			if err != nil {
				bb.Error(err)
			}
			// defer n.Close()

			bb.Run("InsertParallel", func(sb *testing.B) {
				train := d.Train()

				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				sb.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						_, err = n.Insert(train[i%len(train)])
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

			bb.Run("SearchParallel", func(sb *testing.B) {
				query := d.Query()
				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				sb.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						_, err = n.Search(query[i%len(query)], size, epsilon, radius)
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

func BenchmarkNGTParallelBulk(b *testing.B) {
	parseArgs(b)

	for _, target := range targets {
		d := assets.Data(target)(b)
		b.Run(target, func(bb *testing.B) {
			tmpdir, err := ioutil.TempDir("", "tmpdir")
			if err != nil {
				bb.Error(err)
			}
			defer os.RemoveAll(tmpdir)

			n, err := ngt.New(
				ngt.WithIndexPath(tmpdir),
				ngt.WithObjectType(ngt.Float),
				ngt.WithDimension(d.Dimension()),
			)
			if err != nil {
				bb.Error(err)
			}
			// defer n.Close()

			bb.Run("BulkInsertParallel", func(sb *testing.B) {
				train := d.Train()

				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				sb.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						_, err = n.BulkInsert(train)
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

			bb.Run("SearchParallel", func(sb *testing.B) {
				query := d.Query()
				sb.ReportAllocs()
				sb.ResetTimer()
				sb.StartTimer()
				sb.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						_, err = n.Search(query[i%len(query)], size, epsilon, radius)
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
