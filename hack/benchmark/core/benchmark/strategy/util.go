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

// Package strategy provides benchmark strategy
package strategy

import (
	"context"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/core"
	"github.com/vdaas/vald/internal/errors"
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

func insertAndCreateIndex32(ctx context.Context, c core.Core32, dataset assets.Dataset) (ids []uint, err error) {
	n := 1000

	train := dataset.Train()
	ids = make([]uint, 0, len(train)*n)

	for i := 0; i < n; i++ {
		inserted, errs := c.BulkInsert(train)
		err = wrapErrors(errs)
		if err != nil {
			return nil, err
		}
		ids = append(ids, inserted...)
	}

	err = c.CreateIndex(10)
	if err != nil {
		return nil, err
	}
	return
}

func insertAndCreateIndex64(ctx context.Context, c core.Core64, dataset assets.Dataset) (ids []uint, err error) {
	n := 1000

	train := dataset.TrainAsFloat64()
	ids = make([]uint, 0, len(train)*n)

	for i := 0; i < n; i++ {
		inserted, errs := c.BulkInsert(train)
		err = wrapErrors(errs)
		if err != nil {
			return nil, err
		}
		ids = append(ids, inserted...)
	}

	err = c.CreateIndex(10)
	if err != nil {
		return nil, err
	}
	return
}
