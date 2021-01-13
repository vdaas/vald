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
	"testing"

	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/hack/benchmark/internal/e2e"
	"github.com/vdaas/vald/internal/client/v1/client"
)

type createIndex struct {
	poolSize uint32
	client.Indexer
}

func NewCreateIndex(opts ...CreateIndexOption) e2e.Strategy {
	ci := new(createIndex)
	for _, opt := range append(defaultCreateIndexOptions, opts...) {
		opt(ci)
	}
	return ci
}

func (ci *createIndex) Run(ctx context.Context, b *testing.B, c client.Client, dataset assets.Dataset) {
	b.Run("CreateIndex", func(bb *testing.B) {
		for i := 0; i < bb.N; i++ {
			ci.do(ctx, b)
		}
	})
}

func (ci *createIndex) do(ctx context.Context, b *testing.B) {
	if _, err := ci.Indexer.CreateIndex(ctx, &client.ControlCreateIndexRequest{
		PoolSize: ci.poolSize,
	}); err != nil {
		b.Error(err)
	}
}
