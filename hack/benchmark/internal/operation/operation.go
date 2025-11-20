// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package operation

import (
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/hack/benchmark/internal/assets"
	"github.com/vdaas/vald/internal/client/v1/client"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"google.golang.org/grpc/codes"
)

type Operation interface {
	Search(b *testing.B, ds assets.Dataset)
	SearchByID(b *testing.B, maxIdNum int)
	StreamSearch(b *testing.B, ds assets.Dataset)
	StreamSearchByID(b *testing.B, maxIdNum int)
	Insert(b *testing.B, ds assets.Dataset) (insertedNum int)
	StreamInsert(b *testing.B, ds assets.Dataset) (insertedNum int)
	Remove(b *testing.B, maxIdNum int)
	StreamRemove(b *testing.B, maxIdNum int)
	CreateIndex(b *testing.B)
}

type operation struct {
	client   client.Client
	indexerC client.Indexer
}

func New(opts ...Option) Operation {
	o := &operation{}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func (o *operation) CreateIndex(b *testing.B) {
	req := &payload.Control_CreateIndexRequest{
		PoolSize: 10000,
	}
	b.ResetTimer()
	b.Run("CreateIndex", func(b *testing.B) {
		ctx := b.Context()
		for i := 0; i < b.N; i++ {
			_, err := o.indexerC.CreateIndex(ctx, req)
			if err != nil {
				st, _ := status.FromError(err)
				if st.Code() != codes.FailedPrecondition {
					b.Error(err)
				}
			}
		}
	})
}
