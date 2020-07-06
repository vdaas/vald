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
	"sync"

	"github.com/vdaas/vald/apis/grpc/agent/core"
	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/assets"
	"github.com/vdaas/vald/pkg/tools/cli/loadtest/config"
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
	v := dataset.Train()
	ids := dataset.IDs()
	i := 0
	size := len(v)
	m := &sync.Mutex{}
	return func() (ret interface{}) {
		m.Lock()
		defer m.Unlock()
		if i < size {
			ret = &payload.Object_Vector{
				Id:     ids[i],
				Vector: v[i],
			}
			i++
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
		v := make([]*payload.Object_Vector, 0, n)
		for i := 0; i < n; i++ {
			d := provider()
			if d == nil {
				break
			}
			v = append(v, d.(*payload.Object_Vector))
		}
		return &payload.Object_Vectors{
			Vectors: v,
		}
	}, size
}

type inserter interface {
	Insert(context.Context, *payload.Object_Vector, ...grpc.CallOption) (*payload.Empty, error)
	MultiInsert(context.Context, *payload.Object_Vectors, ...grpc.CallOption) (*payload.Empty, error)
}

func agent(conn *grpc.ClientConn) inserter {
	return core.NewAgentClient(conn)
}

func gateway(conn *grpc.ClientConn) inserter {
	return vald.NewValdClient(conn)
}

func insert(c func(*grpc.ClientConn) inserter) loadFunc {
	return func(ctx context.Context, conn *grpc.ClientConn, i interface{}, copts ...grpc.CallOption) (interface{}, error) {
		return c(conn).Insert(ctx, i.(*payload.Object_Vector), copts...)
	}
}

func bulkInsert(c func(*grpc.ClientConn) inserter) loadFunc {
	return func(ctx context.Context, conn *grpc.ClientConn, i interface{}, copts ...grpc.CallOption) (interface{}, error) {
		return c(conn).MultiInsert(ctx, i.(*payload.Object_Vectors), copts...)
	}
}

func (l *loader) newInsert() (f loadFunc, err error) {
	switch {
	case l.batchSize == 1:
		switch l.service {
		case config.Agent:
			f = insert(agent)
		case config.Gateway:
			f = insert(gateway)
		default:
			err = errors.Errorf("undefined service: %s", l.service.String())
		}
	case l.batchSize >= 2:
		switch l.service {
		case config.Agent:
			f = bulkInsert(agent)
		case config.Gateway:
			f = bulkInsert(gateway)
		default:
			err = errors.Errorf("undefined service: %s", l.service.String())
		}
	default:
		err = errors.Errorf("batch size must be natural number.")
	}
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (l *loader) newStreamInsert() (f loadFunc, err error) {
	l.batchSize = 1
	switch l.service {
	case config.Agent:
		f = func(ctx context.Context, conn *grpc.ClientConn, i interface{}, copts ...grpc.CallOption) (interface{}, error) {
			return core.NewAgentClient(conn).StreamInsert(ctx, copts...)
		}
	case config.Gateway:
		f = func(ctx context.Context, conn *grpc.ClientConn, i interface{}, copts ...grpc.CallOption) (interface{}, error) {
			return vald.NewValdClient(conn).StreamInsert(ctx, copts...)
		}
	default:
		err = errors.Errorf("undefined service: %s", l.service.String())
	}
	if err != nil {
		return nil, err
	}
	return f, nil
}
