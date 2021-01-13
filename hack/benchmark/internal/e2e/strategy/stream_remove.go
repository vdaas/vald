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

// Package strategy provides strategy for e2e testing functions
package strategy

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/internal/client/v1/client"
	"github.com/vdaas/vald/internal/net/grpc"
)

type streamRemove struct{}

func NewStreamRemove(opts ...StreamRemoveOption) e2e.Strategy {
	sr := new(streamRemove)
	for _, opt := range append(defaultStreamRemoveOptions, opts...) {
		opt(sr)
	}
	return sr
}

func (sr *streamRemove) dataProvider(total *uint32, b *testing.B, dataset assets.Dataset) func() *client.RemoveRequest {
	var cnt uint32

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()
	defer b.StopTimer()

	return func() *client.RemoveRequest {
		n := int(atomic.AddUint32(&cnt, 1)) - 1
		if n >= b.N {
			return nil
		}

		total := int(atomic.AddUint32(total, 1)) - 1
		return &client.RemoveRequest{
			Id: &client.ObjectID{
				Id: fmt.Sprint(total % dataset.TrainSize()),
			},
		}
	}
}

func (sr *streamRemove) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	var total uint32
	b.Run("StreamRemove", func(bb *testing.B) {
		srv, err := c.StreamRemove(ctx)
		if err != nil {
			bb.Error(err)
		}
		grpc.BidirectionalStreamClient(srv, func() interface{} {
			return sr.dataProvider(&total, bb, dataset)()
		}, func() interface{} {
			return new(client.RemoveRequest)
		}, func(msg interface{}, err error) {
		})
	})
}
