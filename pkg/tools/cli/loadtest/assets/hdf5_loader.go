//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package assets

import (
	"github.com/vdaas/vald/internal/errors"
	"gonum.org/v1/hdf5"
)

type loaderFunc func(*hdf5.Dataset, int, int, int) (interface{}, error)

func loadFloat32(dset *hdf5.Dataset, npoints, row, dim int) (interface{}, error) {
	v := make([]float32, npoints)
	if err := dset.Read(&v); err != nil {
		return nil, err
	}

	vec := make([][]float32, row)
	for i := 0; i < row; i++ {
		vec[i] = make([]float32, dim)
		for j := 0; j < dim; j++ {
			vec[i][j] = v[i*dim+j]
		}
	}
	return vec, nil
}

func loadInt(dset *hdf5.Dataset, npoints, row, dim int) (interface{}, error) {
	v := make([]int32, npoints)
	if err := dset.Read(&v); err != nil {
		return nil, err
	}

	vec := make([][]int, row)
	for i := 0; i < row; i++ {
		vec[i] = make([]int, dim)
		for j := 0; j < dim; j++ {
			vec[i][j] = int(v[i*dim+j])
		}
	}
	return vec, nil
}

func loadDataset(file *hdf5.File, name string, f loaderFunc) (dim int, vec interface{}, err error) {
	dset, err := file.OpenDataset(name)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		err = dset.Close()
	}()
	space := dset.Space()
	defer func() {
		err = space.Close()
	}()
	dims, _, err := space.SimpleExtentDims()
	if err != nil {
		return 0, nil, err
	}

	row, dim := int(dims[0]), int(dims[1])
	vec, err = f(dset, space.SimpleExtentNPoints(), row, dim)
	return dim, vec, err
}

// Load returns loaded vectors and so on from approximate nearest neighbor benchmark dataset.
func Load(path string) (train, test, distances [][]float32, neighbors [][]int, dim int, err error) {
	f, err := hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, nil, nil, nil, 0, errors.Wrapf(err, "couldn't open file %s", path)
	}
	defer func() {
		err = f.Close()
	}()
	trainDim, v1, err := loadDataset(f, "train", loadFloat32)
	if err != nil {
		return nil, nil, nil, nil, 0, errors.Wrapf(err, "couldn't load train dataset for path %s", path)
	}
	train = v1.([][]float32)
	dim = trainDim
	testDim, v2, err := loadDataset(f, "test", loadFloat32)
	if err != nil {
		return train, nil, nil, nil, dim, errors.Wrapf(err, "couldn't load test dataset for path %s", path)
	}
	test = v2.([][]float32)
	if dim != testDim {
		return train, test, nil, nil, 0, errors.Errorf("test has different dimension from train")
	}
	distancesDim, v3, err := loadDataset(f, "distances", loadFloat32)
	if err != nil {
		return train, test, nil, nil, dim, errors.Wrapf(err, "couldn't load distances dataset for path %s", path)
	}
	distances = v3.([][]float32)

	neighborsDim, v4, err := loadDataset(f, "neighbors", loadInt)
	if err != nil {
		return train, test, distances, nil, trainDim, errors.Wrapf(err, "couldn't load neighbors dataset for path %s", path)
	}
	neighbors = v4.([][]int)
	if distancesDim != neighborsDim {
		return train, test, distances, neighbors, dim, errors.Errorf("neighbors has different dimension from distances")
	}

	return train, test, distances, neighbors, dim, nil
}
