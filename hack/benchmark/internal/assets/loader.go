//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
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
	"fmt"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/internal/log"
	"gonum.org/v1/hdf5"
)

func loadDataset(f *hdf5.File, name string) (dim int, vec [][]float64, err error) {
	dset, err := f.OpenDataset(name)
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

	var row int
	row, dim = int(dims[0]), int(dims[1])
	v := make([]float32, space.SimpleExtentNPoints())
	if err := dset.Read(&v); err != nil {
		return dim, nil, err
	}

	vec = make([][]float64, row)
	for i := 0; i < row; i++ {
		vec[i] = make([]float64, dim)
		for j := 0; j < dim; j++ {
			vec[i][j] = float64(v[i*dim+j])
		}
	}
	return dim, vec, nil
}

func Load(path string) (train [][]float64, test [][]float64, dim int, err error) {
	f, err := hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		log.Error(path)
		return nil, nil, 0, err
	}
	defer func() {
		err = f.Close()
	}()
	var trainDim int
	trainDim, train, err = loadDataset(f, "train")
	if err != nil {
		log.Error(path)
		return nil, nil, 0, err
	}
	var testDim int
	testDim, test, err = loadDataset(f, "test")
	if err != nil {
		log.Error(path)
		return train, nil, trainDim, err
	}
	if trainDim != testDim {
		return train, test, 0, fmt.Errorf("test has different dimension from train")
	}
	return train, test, trainDim, nil
}

func CreateIDs(n int) []string {
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, fuid.String())
	}
	return ids
}

func LoadDataAndIDs(path string) (ids []string, train [][]float64, test [][]float64, dim int, err error) {
	train, test, dim, err = Load(path)
	if err != nil {
		return nil, train, test, dim, err
	}
	return CreateIDs(len(train)), train, test, dim, nil
}
