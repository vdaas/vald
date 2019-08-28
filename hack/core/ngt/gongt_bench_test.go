// MIT License
//
// Copyright (c) 2019 Kosuke Morimoto
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
