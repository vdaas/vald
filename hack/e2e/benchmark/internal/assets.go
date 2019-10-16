//
// Copyright (C) 2019 kpango (Yusuke Kato)
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
package internal

import (
	"github.com/vdaas/vald/internal/log"
	"gonum.org/v1/hdf5"
)

func loadDataset(f *hdf5.File, name string) (vec [][]float64, err error) {
	dset, err := f.OpenDataset(name)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	v := make([]float32, space.SimpleExtentNPoints())
	if err := dset.Read(&v); err != nil {
		return nil, err
	}

	row, col := int(dims[0]), int(dims[1])

	vec = make([][]float64, row)
	for i := 0; i < row; i++ {
		vec[i] = make([]float64, col)
		for j := 0; j < col; j++ {
			vec[i][j] = float64(v[i*col+j])
		}
	}
	return vec, nil

}

func Load(path string) (train [][]float64, test [][]float64, err error) {
	f, err := hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		log.Error(path)
		return nil, nil, err
	}
	defer func() {
		err = f.Close()
	}()
	train, err = loadDataset(f, "train")
	if err != nil {
		log.Error(path)
		return nil, nil, err
	}
	test, err = loadDataset(f, "test")
	if err != nil {
		log.Error(path)
		return train, nil, err
	}
	return train, test, nil
}
