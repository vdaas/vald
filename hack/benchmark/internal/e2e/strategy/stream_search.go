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
	"sync/atomic"
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/internal/client/v1/client"
	"github.com/vdaas/vald/internal/net/grpc"
)

type streamSearch struct {
	cfg *client.SearchConfig
}

func NewStreamSearch(opts ...StreamSearchOption) e2e.Strategy {
	s := new(streamSearch)
	for _, opt := range append(defaultStreamSearchOptions, opts...) {
		opt(s)
	}
	return s
}

func (s *streamSearch) dataProvider(total *uint32, b *testing.B, dataset assets.Dataset) func() *client.SearchRequest {
	var cnt uint32

	b.StopTimer()
	b.ReportAllocs()
	b.ResetTimer()
	b.StartTimer()

	return func() *client.SearchRequest {
		n := int(atomic.AddUint32(&cnt, 1)) - 1
		if n >= b.N {
			return nil
		}

		total := int(atomic.AddUint32(total, 1)) - 1
		v, err := dataset.Query(total % dataset.QuerySize())
		if err != nil {
			return nil
		}
		return &client.SearchRequest{
			Vector: v.([]float32),
			Config: s.cfg,
		}
	}
}

func (s *streamSearch) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	var total uint32
	b.Run("StreamSearch", func(bb *testing.B) {
		srv, err := c.StreamSearch(ctx)
		if err != nil {
			bb.Error(err)
		}
		grpc.BidirectionalStreamClient(srv, func() interface{} {
			return s.dataProvider(&total, bb, dataset)()
		}, func() interface{} {
			return new(client.SearchRequest)
		}, func(msg interface{}, err error) {
		})
	})
}
