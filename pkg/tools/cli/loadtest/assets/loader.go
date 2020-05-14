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
package assets

import (
	"strconv"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
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

func Load(path string) (train, test, distances [][]float32, neighbors [][]int, dim int, err error) {
	f, err := hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		log.Error(path)
		return nil, nil, nil, nil, 0, err
	}
	defer func() {
		err = f.Close()
	}()
	trainDim, v1, err := loadDataset(f, "train", loadFloat32)
	if err != nil {
		log.Error(path)
		return nil, nil, nil, nil, 0, err
	}
	train = v1.([][]float32)
	dim = trainDim
	testDim, v2, err := loadDataset(f, "test", loadFloat32)
	if err != nil {
		log.Error(path)
		return train, nil, nil, nil, dim, err
	}
	test = v2.([][]float32)
	if dim != testDim {
		return train, test, nil, nil, 0, errors.Errorf("test has different dimension from train")
	}
	distancesDim, v3, err := loadDataset(f, "distances", loadFloat32)
	if err != nil {
		log.Error(path)
		return train, test, nil, nil, dim, err
	}
	distances = v3.([][]float32)

	neighborsDim, v4, err := loadDataset(f, "neighbors", loadInt)
	if err != nil {
		log.Error(path)
		return train, test, distances, nil, trainDim, err
	}
	neighbors = v4.([][]int)
	if distancesDim != neighborsDim {
		return train, test, distances, neighbors, dim, errors.Errorf("neighbors has different dimension from distances")
	}

	return train, test, distances, neighbors, dim, nil
}

func CreateRandomIDs(n int) (ids []string) {
	ids = make([]string, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, fuid.String())
	}
	return ids
}

func CreateRandomIDsWithLength(n, l int) (ids []string) {
	ids = make([]string, 0, n)
	for i := 0; i < n; i++ {
		id := fuid.String()
		for len(id) < l {
			id = id + fuid.String()
		}
		ids = append(ids, id[:l])
	}
	return ids
}

func CreateSequentialIDs(n int) []string {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, strconv.Itoa(i))
	}
	return ids
}

func LoadDataWithRandomIDs(path string) (ids []string, train, test, distances [][]float32, neighbors [][]int, dim int, err error) {
	train, test, distances, neighbors, dim, err = Load(path)
	if err != nil {
		return nil, train, test, distances, neighbors, dim, err
	}
	return CreateRandomIDs(len(train)), train, test, distances, neighbors, dim, nil
}

func LoadDataWithSequentialIDs(path string) (ids []string, train, test, distances [][]float32, neighbors [][]int, dim int, err error) {
	train, test, distances, neighbors, dim, err = Load(path)
	if err != nil {
		return nil, train, test, distances, neighbors, dim, err
	}
	return CreateSequentialIDs(len(train)), train, test, distances, neighbors, dim, nil
}
