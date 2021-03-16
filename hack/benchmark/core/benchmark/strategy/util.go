//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package strategy provides benchmark strategy
package strategy

import (
	"context"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
)

const (
	bulkInsertCnt = 1000
)

func wrapErrors(errs []error) (wrapped error) {
	for _, err := range errs {
		if err != nil {
			if wrapped == nil {
				wrapped = err
			} else {
				wrapped = errors.Wrap(wrapped, err.Error())
			}
		}
	}
	return
}

func insertAndCreateIndex32(ctx context.Context, c algorithm.Bit32, dataset assets.Dataset) (ids []uint, err error) {
	ids = make([]uint, 0, dataset.TrainSize()*bulkInsertCnt)

	n := 0
	for i := 0; i < bulkInsertCnt; i++ {
		train := make([][]float32, 0, dataset.TrainSize()/bulkInsertCnt)
		for j := 0; j < len(train); j++ {
			v, err := dataset.Train(n)
			if err != nil {
				n = 0
				break
			}
			train = append(train, v.([]float32))
			n++
		}
		inserted, errs := c.BulkInsert(train)
		err = wrapErrors(errs)
		if err != nil {
			return nil, err
		}
		ids = append(ids, inserted...)
	}

	err = c.CreateIndex(uint32((dataset.TrainSize() * bulkInsertCnt) / 100))
	if err != nil {
		return nil, err
	}
	return
}

func insertAndCreateIndex64(ctx context.Context, c algorithm.Bit64, dataset assets.Dataset) (ids []uint, err error) {
	ids = make([]uint, 0, dataset.TrainSize()*bulkInsertCnt)

	n := 0
	for i := 0; i < bulkInsertCnt; i++ {
		train := make([][]float64, 0, dataset.TrainSize()/bulkInsertCnt)
		for j := 0; j < len(train); j++ {
			v, err := dataset.Train(n)
			if err != nil {
				n = 0
				break
			}
			train = append(train, float32To64(v.([]float32)))
			n++
		}
		inserted, errs := c.BulkInsert(train)
		err = wrapErrors(errs)
		if err != nil {
			return nil, err
		}
		ids = append(ids, inserted...)
	}

	err = c.CreateIndex(uint32((dataset.TrainSize() * bulkInsertCnt) / 100))
	if err != nil {
		return nil, err
	}
	return
}

func float32To64(x []float32) (y []float64) {
	y = make([]float64, len(x))
	for i, a := range x {
		y[i] = float64(a)
	}
	return y
}
