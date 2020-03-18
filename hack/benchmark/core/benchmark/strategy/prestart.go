package strategy

import (
	"context"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/core/ngt"
)

type PreStart func(context.Context, *testing.B, ngt.NGT, assets.Dataset) []uint

type preStart struct{}

func (p *preStart) Func(
	ctx context.Context, b *testing.B, ngt ngt.NGT, dataset assets.Dataset,
) (ids []uint) {
	ids = make([]uint, 0, len(dataset.Train()))
	isrtIds, errs := ngt.BulkInsert(dataset.Train())
	for i, err := range errs {
		if err != nil {
			b.Error(err)
		} else {
			ids = append(ids, isrtIds[i])
		}
	}
	return
}
