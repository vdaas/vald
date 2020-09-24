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
	"sync/atomic"

	"github.com/vdaas/vald/apis/grpc/agent/core"
	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/config"
)

func searchRequestProvider(dataset assets.Dataset) (func() interface{}, int, error) {
	size := dataset.QuerySize()
	idx := int32(-1)
	return func() (ret interface{}) {
		if i := int(atomic.AddInt32(&idx, 1)); i < size {
			v, err := dataset.Query(i)
			if err != nil {
				return nil
			}
			ret = &payload.Search_Request{
				Vector: v.([]float32),
			}
		}
		return ret
	}, size, nil
}

func (l *loader) newSearch() (loadFunc, error) {
	switch l.service {
	case config.Agent:
		return func(ctx context.Context, conn *grpc.ClientConn, i interface{}, copts ...grpc.CallOption) (interface{}, error) {
			return core.NewAgentClient(conn).Search(ctx, i.(*payload.Search_Request), copts...)
		}, nil
	case config.Gateway:
		return func(ctx context.Context, conn *grpc.ClientConn, i interface{}, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).Search(ctx, i.(*payload.Search_Request), copts...)
		}, nil
	default:
		return nil, errors.Errorf("undefined service: %s", l.service.String())
	}
}

func (l *loader) newStreamSearch() (loadFunc, error) {
	switch l.service {
	case config.Agent:
		return func(ctx context.Context, conn *grpc.ClientConn, i interface{}, copts ...grpc.CallOption) (interface{}, error) {
			return core.NewAgentClient(conn).StreamSearch(ctx, copts...)
		}, nil
	case config.Gateway:
		return func(ctx context.Context, conn *grpc.ClientConn, i interface{}, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).StreamSearch(ctx, copts...)
		}, nil
	default:
		return nil, errors.Errorf("undefined service: %s", l.service.String())
	}
}
