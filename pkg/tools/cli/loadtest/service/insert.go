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
package service

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
)

func NewInsert(opts ...Option) (Loader, error) {
	l, err := newLoader(opts...)
	if err != nil {
		return nil, err
	}
	l.requestsFunc = func(dataset assets.Dataset) ([]interface{}, error) {
		vectors := dataset.Train()
		ids := dataset.IDs()
		requests := make([]interface{}, len(vectors))
		for j, v := range vectors {
			requests[j] = &payload.Object_Vector{
				Id:     ids[j],
				Vector: v,
			}
		}
		return requests, nil
	}
	l.loaderFunc = func(ctx context.Context, c vald.ValdClient, i interface{}, copts ...grpc.CallOption) error {
		_, err := c.Insert(ctx, i.(*payload.Object_Vector), copts...)
		return err
	}
	return l, nil
}
