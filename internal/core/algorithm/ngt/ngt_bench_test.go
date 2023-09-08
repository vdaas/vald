//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package ngt provides implementation of Go API for https://github.com/yahoojapan/NGT
package ngt

import (
	"runtime"
	"testing"

	"github.com/vdaas/vald/internal/core/malloc"
	"gonum.org/v1/hdf5"
)

var (
	vectors [][]float32
	n       NGT
	ids     []uint
)

func init() {
	vectors, _, _ = load("sift-128-euclidean.hdf5")
	n, _ = New(
		WithDimension(len(vectors[0])),
		WithDefaultPoolSize(8),
		WithObjectType(Float),
		WithDistanceType(L2),
	)
	runtime.GC()
}

func BenchmarkNGT(b *testing.B) {
	log := func() {
		mem := new(runtime.MemStats)
		runtime.ReadMemStats(mem)
		b.Logf("           heap in use: %v", mem.HeapInuse)
		b.Logf("           total alloc: %v", mem.TotalAlloc)

		m, _ := malloc.GetMallocInfo()
		b.Logf("       total fast size: %v", m.Total[0].Size)
		b.Logf("       total rest size: %v", m.Total[1].Size)
		b.Logf("   system current size: %v", m.System[0].Size)
		b.Logf("       system max size: %v", m.System[1].Size)
		b.Logf("     aspace total size: %v", m.Aspace[0].Size)
		b.Logf("  aspace mprotect size: %v", m.Aspace[1].Size)
	}

	b.Log("start")
	log()
	defer func() {
		b.Log("end")
		log()
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ids = make([]uint, len(vectors))
		for i, vector := range vectors {
			id, err := n.Insert(vector)
			if err != nil {
				b.Fatal(err)
			}
			ids[i] = id
		}
		b.Log("insert")
		log()

		if err := n.CreateIndex(8); err != nil {
			b.Fatal(err)
		}
		b.Log("create index")
		log()

		for _, id := range ids {
			if err := n.Remove(id); err != nil {
				b.Fatal(err)
			}
		}
		b.Log("remove")
		log()
	}
}

// load function loads training and test vector from hdf file. The size of ids is same to the number of training data.
// Each id, which is an element of ids, will be set a random number.
func load(path string) (train, test [][]float32, err error) {
	var f *hdf5.File
	f, err = hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	// readFn function reads vectors of the hierarchy with the given the name.
	readFn := func(name string) ([][]float32, error) {
		// Opens and returns a named Dataset.
		// The returned dataset must be closed by the user when it is no longer needed.
		d, err := f.OpenDataset(name)
		if err != nil {
			return nil, err
		}
		defer d.Close()

		// Space returns an identifier for a copy of the dataspace for a dataset.
		sp := d.Space()
		defer sp.Close()

		// SimpleExtentDims returns dataspace dimension size and maximum size.
		dims, _, _ := sp.SimpleExtentDims()
		row, dim := int(dims[0]), int(dims[1])

		// Gets the stored vector. All are represented as one-dimensional arrays.
		// The type of the slice depends on your dataset.
		// For fashion-mnist-784-euclidean.hdf5, the datatype is float32.
		vec := make([]float32, sp.SimpleExtentNPoints())
		if err := d.Read(&vec); err != nil {
			return nil, err
		}

		// Converts a one-dimensional array to a two-dimensional array.
		// Use the `dim` variable as a separator.
		vecs := make([][]float32, row)
		for i := 0; i < row; i++ {
			vecs[i] = make([]float32, dim)
			for j := 0; j < dim; j++ {
				vecs[i][j] = float32(vec[i*dim+j])
			}
		}

		return vecs, nil
	}

	// Gets vector of `train` hierarchy.
	train, err = readFn("train")
	if err != nil {
		return nil, nil, err
	}

	// Gets vector of `test` hierarchy.
	test, err = readFn("test")
	if err != nil {
		return nil, nil, err
	}

	return
}
