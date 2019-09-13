//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
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
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/vdaas/vald/hack/core/ngt/assets"
	"github.com/yahoojapan/gongt"
)

func BenchmarkGoNGTFashionMNIST(b *testing.B) {
	benchmarkGoNGTAll(b, "Fashion-MNIST", "fashion-mnist-784-euclidean")
}

func BenchmarkGoNGTGlove25(b *testing.B) {
	benchmarkGoNGTAll(b, "Glove-25", "glove-25-angular")
}

func BenchmarkGoNGTGlove50(b *testing.B) {
	benchmarkGoNGTAll(b, "Glove-50", "glove-50-angular")
}

func BenchmarkGoNGTGlove100(b *testing.B) {
	benchmarkGoNGTAll(b, "Glove-100", "glove-100-angular")
}

//func BenchmarkGoNGTGlove200(b *testing.B) {
//	benchmarkGoNGTAll(b, "Glove-200", "glove-200-angular")
//}

func BenchmarkGoNGTMIST(b *testing.B) {
	benchmarkGoNGTAll(b, "MNIST", "mnist-784-euclidean")
}

func BenchmarkGoNGTNYTimes(b *testing.B) {
	benchmarkGoNGTAll(b, "NYTimes", "nytimes-256-angular")
}

func BenchmarkGoNGTSIFT(b *testing.B) {
	benchmarkGoNGTAll(b, "SIFT", "sift-128-euclidean")
}

func benchmarkGoNGTAll(b *testing.B, indexName, datasetName string) {
	d, err := assets.LoadDataset(fmt.Sprintf("%s/%s.hdf5", indexBasePath, datasetName))
	if err != nil {
		b.Error(err)
	}
	defer d.Close()
	benchmarkGoNGTInsert(b, d)
	benchmarkGoNGTSearch(b, fmt.Sprintf("%s/%s", indexBasePath, indexName), d)
}

func benchmarkGoNGTInsert(b *testing.B, d *assets.Dataset) {
	dataset, err := d.LoadTrain()
	if err != nil {
		b.Error(err)
	}
	b.Run("Insert", func(sb *testing.B) {
		tmpdir, err := ioutil.TempDir("", "tmpdir")
		if err != nil {
			sb.Error(err)
		}
		defer os.RemoveAll(tmpdir)

		n := gongt.New(tmpdir).SetObjectType(gongt.Float).SetDimension(len(dataset[0])).Open()
		defer n.Close()

		sb.ReportAllocs()
		sb.ResetTimer()
		sb.StartTimer()
		for i := 0; i < sb.N; i++ {
			_, err = n.Insert(dataset[i%len(dataset)])
			if err != nil {
				sb.Error(err)
			}
		}
		sb.StopTimer()
	})

	b.Run("InsertParallel", func(sb *testing.B) {
		tmpdir, err := ioutil.TempDir("", "tmpdir")
		if err != nil {
			sb.Error(err)
		}
		defer os.RemoveAll(tmpdir)

		n := gongt.New(tmpdir).SetObjectType(gongt.Float).SetDimension(len(dataset[0])).Open()
		defer n.Close()

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
}

func benchmarkGoNGTSearch(b *testing.B, indexPath string, d *assets.Dataset) {
	dataset, err := d.LoadTest()
	if err != nil {
		b.Error(err)
	}

	n := gongt.New(indexPath).Open()
	if err != nil {
		b.Error(err)
	}
	defer n.Close()

	size := 10
	b.Run("Search", func(sb *testing.B) {
		sb.ReportAllocs()
		sb.ResetTimer()
		sb.StartTimer()
		for i := 0; i < sb.N; i++ {
			_, err = n.Search(dataset[i%len(dataset)], size, 0.1)
			if err != nil {
				sb.Error(err)
			}
		}
		sb.StopTimer()
	})

	b.Run("SearchParallel", func(sb *testing.B) {
		sb.ReportAllocs()
		sb.ResetTimer()
		sb.StartTimer()
		sb.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				_, err = n.Search(dataset[i%len(dataset)], size, 0.1)
				if err != nil {
					sb.Error(err)
				}
				i++
			}
		})
		sb.StopTimer()
	})
}
