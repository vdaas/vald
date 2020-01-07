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
	"reflect"

	gmeta "github.com/vdaas/vald/apis/grpc/meta"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
)

type Meta interface {
	Start(ctx context.Context) (<-chan error, error)
	GetMeta(context.Context, string) (string, error)
	GetMetas(context.Context, ...string) ([]string, error)
	GetUUID(context.Context, string) (string, error)
	GetUUIDs(context.Context, ...string) ([]string, error)
	SetUUIDandMeta(context.Context, string, string) error
	SetUUIDandMetas(context.Context, map[string]string) error
	DeleteMeta(context.Context, string) (string, error)
	DeleteMetas(context.Context, ...string) ([]string, error)
	DeleteUUID(context.Context, string) (string, error)
	DeleteUUIDs(context.Context, ...string) ([]string, error)
}

type meta struct {
	addr   string
	client grpc.Client
}

func NewMeta(opts ...MetaOption) (mi Meta, err error) {
	m := new(meta)
	for _, opt := range append(defaultMetaOpts, opts...) {
		if err = opt(m); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return m, nil
}

func (m *meta) Start(ctx context.Context) (<-chan error, error) {
	return m.client.StartConnectionMonitor(ctx)
}

func (m *meta) GetMeta(ctx context.Context, uuid string) (v string, err error) {
	val, err := m.client.Do(ctx, m.addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		val, err := gmeta.NewMetaClient(conn).GetMeta(ctx, &payload.Meta_Key{
			Key: uuid,
		}, copts...)
		if err != nil {
			return nil, err
		}
		return val.GetVal(), nil
	})
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

func (m *meta) GetMetas(ctx context.Context, uuids ...string) ([]string, error) {
	vals, err := m.client.Do(ctx, m.addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		vals, err := gmeta.NewMetaClient(conn).GetMetas(ctx, &payload.Meta_Keys{
			Keys: uuids,
		}, copts...)
		if vals != nil {
			return vals.GetVals(), err
		}
		return nil, err
	})
	if vals != nil {
		vs, ok := vals.([]string)
		if ok {
			return vs, err
		}
	}
	return nil, err
}

func (m *meta) GetUUID(ctx context.Context, meta string) (string, error) {
	key, err := m.client.Do(ctx, m.addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		key, err := gmeta.NewMetaClient(conn).GetMetaInverse(ctx, &payload.Meta_Val{
			Val: meta,
		}, copts...)
		if err != nil {
			return nil, err
		}
		return key.GetKey(), nil
	})
	if err != nil {
		return "", err
	}
	return key.(string), nil
}

func (m *meta) GetUUIDs(ctx context.Context, metas ...string) ([]string, error) {
	keys, err := m.client.Do(ctx, m.addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		keys, err := gmeta.NewMetaClient(conn).GetMetasInverse(ctx, &payload.Meta_Vals{
			Vals: metas,
		}, copts...)
		if keys != nil {
			return keys.GetKeys(), err
		}
		return nil, err
	})
	if keys != nil {
		ks, ok := keys.([]string)
		if ok {
			return ks, err
		}
	}
	return nil, err
}

func (m *meta) SetUUIDandMeta(ctx context.Context, uuid, meta string) (err error) {
	_, err = m.client.Do(ctx, m.addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		_, err := gmeta.NewMetaClient(conn).SetMeta(ctx, &payload.Meta_KeyVal{
			Key: uuid,
			Val: meta,
		}, copts...)

		return nil, err
	})
	return
}

func (m *meta) SetUUIDandMetas(ctx context.Context, kvs map[string]string) (err error) {
	data := make([]*payload.Meta_KeyVal, len(kvs))
	for uuid, meta := range kvs {
		data = append(data, &payload.Meta_KeyVal{
			Key: uuid,
			Val: meta,
		})
	}
	_, err = m.client.Do(ctx, m.addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		_, err := gmeta.NewMetaClient(conn).SetMetas(ctx, &payload.Meta_KeyVals{
			Kvs: data,
		}, copts...)

		return nil, err
	})
	return
}

func (m *meta) DeleteMeta(ctx context.Context, uuid string) (v string, err error) {
	val, err := m.client.Do(ctx, m.addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		val, err := gmeta.NewMetaClient(conn).DeleteMeta(ctx, &payload.Meta_Key{
			Key: uuid,
		}, copts...)
		if err != nil {
			return nil, err
		}
		return val.GetVal(), nil
	})
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

func (m *meta) DeleteMetas(ctx context.Context, uuids ...string) ([]string, error) {
	vals, err := m.client.Do(ctx, m.addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		vals, err := gmeta.NewMetaClient(conn).DeleteMetas(ctx, &payload.Meta_Keys{
			Keys: uuids,
		}, copts...)
		if err != nil {
			return nil, err
		}
		return vals.GetVals(), nil
	})
	if err != nil {
		return nil, err
	}
	return vals.([]string), nil
}

func (m *meta) DeleteUUID(ctx context.Context, meta string) (string, error) {
	key, err := m.client.Do(ctx, m.addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		key, err := gmeta.NewMetaClient(conn).DeleteMetaInverse(ctx, &payload.Meta_Val{
			Val: meta,
		}, copts...)
		if err != nil {
			return nil, err
		}
		return key.GetKey(), nil
	})
	if err != nil {
		return "", err
	}
	return key.(string), nil
}

func (m *meta) DeleteUUIDs(ctx context.Context, metas ...string) ([]string, error) {
	keys, err := m.client.Do(ctx, m.addr, func(conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		keys, err := gmeta.NewMetaClient(conn).DeleteMetasInverse(ctx, &payload.Meta_Vals{
			Vals: metas,
		}, copts...)
		if err != nil {
			return nil, err
		}
		return keys.GetKeys(), nil
	})
	if err != nil {
		return nil, err
	}
	return keys.([]string), nil
}
