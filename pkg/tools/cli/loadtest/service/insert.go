//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package service

import (
	"context"
	"sync/atomic"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
)

func insertRequestProvider(dataset assets.Dataset, batchSize int) (f func() interface{}, size int, err error) {
	switch {
	case batchSize == 1:
		f, size = objectVectorProvider(dataset)
	case batchSize >= 2:
		f, size = objectVectorsProvider(dataset, batchSize)
	default:
		err = errors.New("batch size must be natural number.")
	}
	if err != nil {
		return nil, 0, err
	}
	return f, size, nil
}

func objectVectorProvider(dataset assets.Dataset) (func() interface{}, int) {
	idx := int32(-1)
	size := dataset.TrainSize()
	return func() (ret interface{}) {
		if i := int(atomic.AddInt32(&idx, 1)); i < size {
			v, err := dataset.Train(i)
			if err != nil {
				return nil
			}
			ret = &payload.Insert_Request{
				Vector: &payload.Object_Vector{
					Id:     fuid.String(),
					Vector: v.([]float32),
				},
			}
		}
		return ret
	}, size
}

func objectVectorsProvider(dataset assets.Dataset, n int) (func() interface{}, int) {
	provider, s := objectVectorProvider(dataset)
	size := s / n
	if s%n != 0 {
		size = size + 1
	}
	return func() (ret interface{}) {
		r := make([]*payload.Insert_Request, 0, n)
		for i := 0; i < n; i++ {
			d := provider()
			if d == nil {
				break
			}
			r = append(r, d.(*payload.Insert_Request))
		}
		if len(r) == 0 {
			return nil
		}
		return &payload.Insert_MultiRequest{
			Requests: r,
		}
	}, size
}

func (l *loader) newInsert() (f loadFunc, err error) {
	switch {
	case l.batchSize == 1:
		f = func(ctx context.Context, conn *grpc.ClientConn, i interface{}, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewInsertClient(conn).Insert(ctx, i.(*payload.Insert_Request), copts...)
		}
	case l.batchSize >= 2:
		f = func(ctx context.Context, conn *grpc.ClientConn, i interface{}, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewInsertClient(conn).MultiInsert(ctx, i.(*payload.Insert_MultiRequest), copts...)
		}
	default:
		err = errors.New("batch size must be natural number.")
	}
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (l *loader) newStreamInsert() (f loadFunc, err error) {
	l.batchSize = 1
	return func(ctx context.Context, conn *grpc.ClientConn, i interface{}, copts ...grpc.CallOption) (interface{}, error) {
		return vald.NewValdClient(conn).StreamInsert(ctx, copts...)
	}, nil
}
