//go:build e2e

// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// package hdf5 provides hdf5 utilities for e2e testing
package hdf5

import (
	"unsafe"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/iter"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test/data/vector/noise"
	"gonum.org/v1/hdf5"
)

type Dataset struct {
	Train     [][]float32
	Test      [][]float32
	Neighbors [][]int
	once      sync.Once
	noiseFunc noise.Func
	maxLen    uint64
}

func (d *Dataset) TrainCycle(num, offset uint64) iter.Cycle[[][]float32, []float32] {
	if num > d.maxLen && d.noiseFunc == nil {
		d.InitNoiseFunc(num)
	}
	return iter.NewCycle(d.Train, num, offset, d.noiseFunc)
}

func (d *Dataset) TestCycle(num, offset uint64) iter.Cycle[[][]float32, []float32] {
	if num > d.maxLen && d.noiseFunc == nil {
		d.InitNoiseFunc(num)
	}
	return iter.NewCycle(d.Test, num, offset, d.noiseFunc)
}

func (d *Dataset) NeighborsCycle(num, offset uint64) iter.Cycle[[][]int, []int] {
	return iter.NewCycle(d.Neighbors, num, offset, nil)
}

func (d *Dataset) InitNoiseFunc(num uint64, opts ...noise.Option) noise.Func {
	if num > d.maxLen && d.noiseFunc == nil {
		d.once.Do(func() {
			data := d.Train
			if len(data) == 0 || len(d.Test) > len(data) {
				data = d.Test
			}
			d.noiseFunc = noise.New(data, num, opts...).Mod()
		})
	}
	return d.noiseFunc
}

// New builds a Dataset directly from already in-memory slices, bypassing
// ToDataset. maxLen is unexported so callers outside this package (e.g.
// synthetic, fixture-less datasets) cannot set it via a struct literal; this
// constructor exists so they can still control it explicitly instead of
// silently getting the zero value, which would make every TrainCycle/
// TestCycle call route through InitNoiseFunc's noise generation path.
func New(train, test [][]float32, neighbors [][]int, maxLen uint64) *Dataset {
	return &Dataset{
		Train:     train,
		Test:      test,
		Neighbors: neighbors,
		maxLen:    maxLen,
	}
}

func ToDataset(name string) (*Dataset, error) {
	file, err := hdf5.OpenFile(name, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	train, err := ReadDataset[float32](file, "train")
	if err != nil {
		return nil, err
	}

	test, err := ReadDataset[float32](file, "test")
	if err != nil {
		return nil, err
	}

	// ann-benchmarks files store neighbors as 4-byte integers; read with a
	// matching element type (see the size guard in ReadDataset) and widen to
	// int afterwards.
	neighbors32, err := ReadDataset[int32](file, "neighbors")
	if err != nil {
		return nil, err
	}
	neighbors := make([][]int, len(neighbors32))
	for i, row := range neighbors32 {
		r := make([]int, len(row))
		for j, v := range row {
			r[j] = int(v)
		}
		neighbors[i] = r
	}

	return &Dataset{
		Train:     train,
		Test:      test,
		Neighbors: neighbors,
		maxLen:    uint64(max(len(train), len(test), len(neighbors))),
	}, nil
}

func ReadDataset[T any](file *hdf5.File, name string) ([][]T, error) {
	data, err := file.OpenDataset(name)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	// gonum/hdf5's Dataset.Read passes the FILE datatype as the memory type,
	// so no element conversion ever happens: reading e.g. a 4-byte integer
	// dataset into []int64 silently packs two file values per Go element and
	// leaves the tail zeroed. Reject size mismatches loudly instead.
	dtype, err := data.Datatype()
	if err != nil {
		return nil, err
	}
	defer dtype.Close()
	var zero T
	if got, want := dtype.Size(), uint(unsafe.Sizeof(zero)); got != want {
		return nil, errors.Errorf(
			"dataset %s element size mismatch: file stores %d-byte elements but requested Go type is %d bytes",
			name, got, want,
		)
	}

	dataspace := data.Space()
	defer dataspace.Close()

	dims, _, err := dataspace.SimpleExtentDims()
	if err != nil {
		return nil, err
	}
	height, width := int(dims[0]), int(dims[1])

	rawFloats := make([]T, dataspace.SimpleExtentNPoints())
	if err := data.Read(&rawFloats); err != nil {
		return nil, err
	}

	vecs := make([][]T, height)
	for i := 0; i < height; i++ {
		vecs[i] = rawFloats[i*width : i*width+width]
	}

	return vecs, nil
}
