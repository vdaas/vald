// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package ngt_test provides implementation of Go API for https://github.com/yahoojapan/NGT
package ngt_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/vdaas/vald/hack/core/ngt/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

const (
	indexBasePath = "../../../hack/core/ngt/assets"
)

func BenchmarkNGTFashionMNIST(b *testing.B) {
	benchmarkNGTAll(b, "Fashion-MNIST", "fashion-mnist-784-euclidean")
}

func BenchmarkNGTGlove25(b *testing.B) {
	benchmarkNGTAll(b, "Glove-25", "glove-25-angular")
}

func BenchmarkNGTGlove50(b *testing.B) {
	benchmarkNGTAll(b, "Glove-50", "glove-50-angular")
}

func BenchmarkNGTGlove100(b *testing.B) {
	benchmarkNGTAll(b, "Glove-100", "glove-100-angular")
}

//func BenchmarkNGTGlove200(b *testing.B) {
//	benchmarkNGTAll(b, "Glove-200", "glove-200-angular")
//}

func BenchmarkNGTMIST(b *testing.B) {
	benchmarkNGTAll(b, "MNIST", "mnist-784-euclidean")
}

func BenchmarkNGTNYTimes(b *testing.B) {
	benchmarkNGTAll(b, "NYTimes", "nytimes-256-angular")
}

func BenchmarkNGTSIFT(b *testing.B) {
	benchmarkNGTAll(b, "SIFT", "sift-128-euclidean")
}

func benchmarkNGTAll(b *testing.B, indexName, datasetName string) {
	d, err := assets.LoadDataset(fmt.Sprintf("%s/%s.hdf5", indexBasePath, datasetName))
	if err != nil {
		b.Error(err)
	}
	defer d.Close()
	benchmarkNGTInsert(b, d)
	benchmarkNGTSearch(b, fmt.Sprintf("%s/%s", indexBasePath, indexName), d)
}

func benchmarkNGTInsert(b *testing.B, d *assets.Dataset) {
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

		n, err := ngt.New(
			ngt.WithIndexPath(tmpdir),
			ngt.WithObjectType(ngt.Float),
			ngt.WithDimension(len(dataset[0])),
		)
		if err != nil {
			sb.Error(err)
		}
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

		n, err := ngt.New(
			ngt.WithIndexPath(tmpdir),
			ngt.WithObjectType(ngt.Float),
			ngt.WithDimension(len(dataset[0])),
		)
		if err != nil {
			sb.Error(err)
		}
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

func benchmarkNGTSearch(b *testing.B, indexPath string, d *assets.Dataset) {
	dataset, err := d.LoadTest()
	if err != nil {
		b.Error(err)
	}

	n, err := ngt.Load(
		ngt.WithIndexPath(indexPath),
	)
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
			_, err = n.Search(dataset[i%len(dataset)], size, 0.1, -1.0)
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
				_, err = n.Search(dataset[i%len(dataset)], size, 0.1, -1.0)
				if err != nil {
					sb.Error(err)
				}
				i++
			}
		})
		sb.StopTimer()
	})
}
