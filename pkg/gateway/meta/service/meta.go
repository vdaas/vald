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

// Package service provides meta service
package service

import (
	"context"
	"reflect"

	gmeta "github.com/vdaas/vald/apis/grpc/meta"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

type Meta interface {
	Start(ctx context.Context) (<-chan error, error)
	Exists(context.Context, string) (bool, error)
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
	addr                string
	client              grpc.Client
	cache               cache.Cache
	enableCache         bool
	expireCheckDuration string
	expireDuration      string
}

const (
	uuidCacheKeyPref = "uuid-"
	metaCacheKeyPref = "meta-"
)

func NewMeta(opts ...Option) (mi Meta, err error) {
	m := new(meta)
	for _, opt := range append(defaultOpts, opts...) {
		if err = opt(m); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	if m.enableCache {
		if m.cache == nil {
			m.cache, err = cache.New(
				cache.WithExpireDuration(m.expireDuration),
				cache.WithExpireCheckDuration(m.expireCheckDuration),
			)
			if err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

func (m *meta) Start(ctx context.Context) (<-chan error, error) {
	if m.enableCache && m.cache != nil {
		m.cache.Start(ctx)
	}
	return m.client.StartConnectionMonitor(ctx)
}

func (m *meta) Exists(ctx context.Context, meta string) (bool, error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-vald/service/Meta.Exists")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if m.enableCache {
		_, ok := m.cache.Get(uuidCacheKeyPref + meta)
		if ok {
			return true, nil
		}
	}
	key, err := m.client.Do(ctx, m.addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		key, err := gmeta.NewMetaClient(conn).GetMetaInverse(ctx, &payload.Meta_Val{
			Val: meta,
		}, copts...)
		if err != nil {
			if status.Code(err) == status.NotFound {
				return "", nil
			}
			return nil, err
		}
		return key.GetKey(), nil
	})
	if err != nil {
		return false, err
	}

	k := key.(string)
	if k == "" {
		return false, nil
	}

	if m.enableCache {
		m.cache.Set(uuidCacheKeyPref+meta, k)
		m.cache.Set(metaCacheKeyPref+k, meta)
	}
	return true, nil
}

func (m *meta) GetMeta(ctx context.Context, uuid string) (v string, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-vald/service/Meta.GetMeta")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if m.enableCache {
		data, ok := m.cache.Get(metaCacheKeyPref + uuid)
		if ok {
			return data.(string), nil
		}
	}
	val, err := m.client.Do(ctx, m.addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
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
	v = val.(string)

	if m.enableCache {
		m.cache.Set(metaCacheKeyPref+uuid, v)
		m.cache.Set(uuidCacheKeyPref+v, uuid)
	}
	return v, nil
}

func (m *meta) GetMetas(ctx context.Context, uuids ...string) ([]string, error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-vald/service/Meta.GetMetas")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if m.enableCache {
		metas, ok := func() (metas []string, ok bool) {
			for _, uuid := range uuids {
				data, ok := m.cache.Get(metaCacheKeyPref + uuid)
				if !ok {
					return nil, false
				}
				metas = append(metas, data.(string))
			}
			return metas, true
		}()
		if ok {
			return metas, nil
		}
	}
	vals, err := m.client.Do(ctx, m.addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
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
			if m.enableCache {
				for i, v := range vs {
					uuid := uuids[i]
					m.cache.Set(metaCacheKeyPref+uuid, v)
					m.cache.Set(uuidCacheKeyPref+v, uuid)
				}
			}
			return vs, err
		}
	}
	return nil, err
}

func (m *meta) GetUUID(ctx context.Context, meta string) (k string, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-vald/service/Meta.GetUUID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if m.enableCache {
		data, ok := m.cache.Get(uuidCacheKeyPref + meta)
		if ok {
			return data.(string), nil
		}
	}
	key, err := m.client.Do(ctx, m.addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
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

	k = key.(string)
	if m.enableCache {
		m.cache.Set(uuidCacheKeyPref+meta, k)
		m.cache.Set(metaCacheKeyPref+k, meta)
	}
	return k, nil
}

func (m *meta) GetUUIDs(ctx context.Context, metas ...string) ([]string, error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-vald/service/Meta.GetUUIDs")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if m.enableCache {
		uuids, ok := func() (uuids []string, ok bool) {
			for _, meta := range metas {
				data, ok := m.cache.Get(uuidCacheKeyPref + meta)
				if !ok {
					return nil, false
				}
				uuids = append(uuids, data.(string))
			}
			return uuids, true
		}()
		if ok {
			return uuids, nil
		}
	}
	keys, err := m.client.Do(ctx, m.addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
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
			if m.enableCache {
				for i, k := range ks {
					meta := metas[i]
					m.cache.Set(uuidCacheKeyPref+meta, k)
					m.cache.Set(metaCacheKeyPref+k, meta)
				}
			}
			return ks, err
		}
	}
	return nil, err
}

func (m *meta) SetUUIDandMeta(ctx context.Context, uuid, meta string) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-vald/service/Meta.SetUUIDandMeta")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, err = m.client.Do(ctx, m.addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		_, err := gmeta.NewMetaClient(conn).SetMeta(ctx, &payload.Meta_KeyVal{
			Key: uuid,
			Val: meta,
		}, copts...)

		return nil, err
	})
	if err != nil {
		return err
	}

	if m.enableCache {
		m.cache.Set(uuidCacheKeyPref+meta, uuid)
		m.cache.Set(metaCacheKeyPref+uuid, meta)
	}
	return nil
}

func (m *meta) SetUUIDandMetas(ctx context.Context, kvs map[string]string) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-vald/service/Meta.SetUUIDandMetas")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	data := make([]*payload.Meta_KeyVal, 0, len(kvs))
	for uuid, meta := range kvs {
		data = append(data, &payload.Meta_KeyVal{
			Key: uuid,
			Val: meta,
		})
	}
	_, err = m.client.Do(ctx, m.addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
		_, err := gmeta.NewMetaClient(conn).SetMetas(ctx, &payload.Meta_KeyVals{
			Kvs: data,
		}, copts...)

		return nil, err
	})
	if err != nil {
		return err
	}

	if m.enableCache {
		for uuid, meta := range kvs {
			m.cache.Set(uuidCacheKeyPref+meta, uuid)
			m.cache.Set(metaCacheKeyPref+uuid, meta)
		}
	}
	return nil
}

func (m *meta) DeleteMeta(ctx context.Context, uuid string) (v string, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway-vald/service/Meta.DeleteMeta")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if m.enableCache {
		meta, ok := m.cache.GetAndDelete(metaCacheKeyPref + uuid)
		if ok {
			m.cache.Delete(uuidCacheKeyPref + meta.(string))
		}
	}
	val, err := m.client.Do(ctx, m.addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
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
	ctx, span := trace.StartSpan(ctx, "vald/gateway-vald/service/Meta.DeleteMetas")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if m.enableCache {
		for _, uuid := range uuids {
			meta, ok := m.cache.GetAndDelete(metaCacheKeyPref + uuid)
			if ok {
				m.cache.Delete(uuidCacheKeyPref + meta.(string))
			}
		}
	}
	vals, err := m.client.Do(ctx, m.addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
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
	ctx, span := trace.StartSpan(ctx, "vald/gateway-vald/service/Meta.DeleteUUID")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if m.enableCache {
		uuid, ok := m.cache.GetAndDelete(uuidCacheKeyPref + meta)
		if ok {
			m.cache.Delete(metaCacheKeyPref + uuid.(string))
		}
	}
	key, err := m.client.Do(ctx, m.addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
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
	ctx, span := trace.StartSpan(ctx, "vald/gateway-vald/service/Meta.DeleteUUIDs")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if m.enableCache {
		for _, meta := range metas {
			uuid, ok := m.cache.GetAndDelete(uuidCacheKeyPref + meta)
			if ok {
				m.cache.Delete(metaCacheKeyPref + uuid.(string))
			}
		}
	}
	keys, err := m.client.Do(ctx, m.addr, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (interface{}, error) {
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
