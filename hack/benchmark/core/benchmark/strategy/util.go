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
