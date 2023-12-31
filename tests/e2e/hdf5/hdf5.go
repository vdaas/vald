//go:build e2e

//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// package hdf5 provides hdf5 utilities for e2e testing
package hdf5

import (
	"gonum.org/v1/hdf5"
)

type Dataset struct {
	Train     [][]float32
	Test      [][]float32
	Neighbors [][]int
}

func HDF5ToDataset(name string) (*Dataset, error) {
	file, err := hdf5.OpenFile(name, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	train, err := ReadDatasetF32(file, "train")
	if err != nil {
		return nil, err
	}

	test, err := ReadDatasetF32(file, "test")
	if err != nil {
		return nil, err
	}

	neighbors32, err := ReadDatasetI32(file, "neighbors")
	if err != nil {
		return nil, err
	}
	neighbors := make([][]int, len(neighbors32))
	for i, ns := range neighbors32 {
		neighbors[i] = make([]int, len(ns))
		for j, n := range ns {
			neighbors[i][j] = int(n)
		}
	}

	return &Dataset{
		Train:     train,
		Test:      test,
		Neighbors: neighbors,
	}, nil
}

func ReadDatasetF32(file *hdf5.File, name string) ([][]float32, error) {
	data, err := file.OpenDataset(name)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	dataspace := data.Space()
	defer dataspace.Close()

	dims, _, err := dataspace.SimpleExtentDims()
	if err != nil {
		return nil, err
	}
	height, width := int(dims[0]), int(dims[1])

	rawFloats := make([]float32, dataspace.SimpleExtentNPoints())
	if err := data.Read(&rawFloats); err != nil {
		return nil, err
	}

	vecs := make([][]float32, height)
	for i := 0; i < height; i++ {
		vecs[i] = rawFloats[i*width : i*width+width]
	}

	return vecs, nil
}

func ReadDatasetI32(file *hdf5.File, name string) ([][]int32, error) {
	data, err := file.OpenDataset(name)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	dataspace := data.Space()
	defer dataspace.Close()

	dims, _, err := dataspace.SimpleExtentDims()
	if err != nil {
		return nil, err
	}
	height, width := int(dims[0]), int(dims[1])

	rawFloats := make([]int32, dataspace.SimpleExtentNPoints())
	if err := data.Read(&rawFloats); err != nil {
		return nil, err
	}

	vecs := make([][]int32, height)
	for i := 0; i < height; i++ {
		vecs[i] = rawFloats[i*width : i*width+width]
	}

	return vecs, nil
}
