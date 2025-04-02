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
package service

import (
	"context"
	"sync/atomic"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
)

func searchRequestProvider(dataset assets.Dataset) (func() (proto.Message, bool), int, error) {
	size := dataset.QuerySize()
	idx := int32(-1)
	return func() (ret proto.Message, ok bool) {
		if i := int(atomic.AddInt32(&idx, 1)); i < size {
			v, err := dataset.Query(i)
			if err != nil {
				return nil, false
			}
			obj := any(&payload.Search_Request{
				Vector: v.([]float32),
				Config: &payload.Search_Config{
					Num:     10,
					Radius:  -1,
					Epsilon: 0.1,
				},
			})
			ret = &obj
		}
		return ret, true
	}, size, nil
}

func (l *loader) newSearch() (loadFunc, error) {
	return func(ctx context.Context, conn *grpc.ClientConn, i any, copts ...grpc.CallOption) (any, error) {
		return vald.NewSearchClient(conn).Search(ctx, i.(*payload.Search_Request), copts...)
	}, nil
}

func (l *loader) newStreamSearch() (loadFunc, error) {
	return func(ctx context.Context, conn *grpc.ClientConn, i any, copts ...grpc.CallOption) (any, error) {
		return vald.NewSearchClient(conn).StreamSearch(ctx, copts...)
	}, nil
}
