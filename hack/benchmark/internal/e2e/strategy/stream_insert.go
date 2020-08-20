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

// Package strategy provides strategy for e2e testing functions
package strategy

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/internal/client"
)

type streamInsert struct{}

func NewStreamInsert(opts ...StreamInsertOption) e2e.Strategy {
	s := new(streamInsert)
	for _, opt := range append(defaultStreamInsertOptions, opts...) {
		opt(s)
	}
	return s
}

func (sisrt *streamInsert) dataProvider(total *uint32, b *testing.B, dataset assets.Dataset) func() *client.ObjectVector {
	var cnt uint32

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()

	return func() *client.ObjectVector {
		n := int(atomic.AddUint32(&cnt, 1)) - 1
		if n >= b.N {
			return nil
		}

		total := int(atomic.AddUint32(total, 1)) - 1
		v, err := dataset.Train(total % dataset.TrainSize())
		if err != nil {
			return nil
		}
		return &client.ObjectVector{
			Id:     fmt.Sprint(n),
			Vector: v.([]float32),
		}
	}
}

func (sisrt *streamInsert) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	var total uint32
	b.Run("StreamInsert", func(bb *testing.B) {
		c.StreamInsert(ctx, sisrt.dataProvider(&total, bb, dataset), func(err error) {
			if err != nil {
				bb.Error(err)
			}
		})
	})
}
