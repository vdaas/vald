//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	gmeta "github.com/vdaas/vald/apis/grpc/meta"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/safety"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Meta interface {
	Start(ctx context.Context) <-chan error
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
	hcDur time.Duration
	host  string
	port  int
	mc    atomic.Value
	eg    errgroup.Group
	bo    backoff.Backoff
	gopts []grpc.DialOption
	copts []grpc.CallOption
}

func NewMeta(opts ...MetaOption) (mi Meta, err error) {
	m := new(meta)
	for _, opt := range opts {
		err = opt(m)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *meta) Start(ctx context.Context) <-chan error {
	ech := make(chan error)
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", m.host, m.port), m.gopts...)
	if err != nil {
		ech <- err
	} else {
		m.mc.Store(gmeta.NewMetaClient(conn))
	}
	m.eg.Go(safety.RecoverFunc(func() (err error) {
		tick := time.NewTicker(m.hcDur)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				close(ech)
				return ctx.Err()
			case <-tick.C:
				if conn == nil ||
					conn.GetState() == connectivity.Shutdown ||
					conn.GetState() == connectivity.TransientFailure {
					conn, err = grpc.DialContext(ctx, fmt.Sprintf("%s:%d", m.host, m.port), m.gopts...)
					if err != nil {
						ech <- err
					} else {
						m.mc.Store(gmeta.NewMetaClient(conn))
						runtime.Gosched()
					}
				}
			}
		}
		return nil
	}))
	return ech
}

func (m *meta) GetMeta(ctx context.Context, key string) (string, error) {
	val, err := m.mc.Load().(gmeta.MetaClient).GetMeta(ctx, &payload.Meta_Key{
		Key: key,
	}, m.copts...)
	if err != nil {
		return "", err
	}
	return val.GetVal(), nil
}

func (m *meta) GetMetas(ctx context.Context, keys ...string) ([]string, error) {
	vals, err := m.mc.Load().(gmeta.MetaClient).GetMetas(ctx, &payload.Meta_Keys{
		Keys: keys,
	}, m.copts...)
	if err != nil {
		return nil, err
	}
	return vals.GetVals(), nil
}

func (m *meta) GetUUID(ctx context.Context, val string) (string, error) {
	key, err := m.mc.Load().(gmeta.MetaClient).GetMetaInverse(ctx, &payload.Meta_Val{
		Val: val,
	}, m.copts...)
	if err != nil {
		return "", err
	}
	return key.GetKey(), nil
}

func (m *meta) GetUUIDs(ctx context.Context, vals ...string) ([]string, error) {
	keys, err := m.mc.Load().(gmeta.MetaClient).GetMetasInverse(ctx, &payload.Meta_Vals{
		Vals: vals,
	}, m.copts...)
	if err != nil {
		return nil, err
	}
	return keys.GetKeys(), nil
}

func (m *meta) SetUUIDandMeta(ctx context.Context, key, val string) error {
	_, err := m.mc.Load().(gmeta.MetaClient).SetMeta(ctx, &payload.Meta_KeyVal{
		Key: key,
		Val: val,
	}, m.copts...)
	return err
}

func (m *meta) SetUUIDandMetas(ctx context.Context, kvs map[string]string) error {
	data := make([]*payload.Meta_KeyVal, len(kvs))
	for k, v := range kvs {
		data = append(data, &payload.Meta_KeyVal{
			Key: k,
			Val: v,
		})
	}
	_, err := m.mc.Load().(gmeta.MetaClient).SetMetas(ctx, &payload.Meta_KeyVals{
		Kvs: data,
	})
	return err
}

func (m *meta) DeleteMeta(ctx context.Context, key string) (string, error) {
	val, err := m.mc.Load().(gmeta.MetaClient).DeleteMeta(ctx, &payload.Meta_Key{
		Key: key,
	}, m.copts...)
	if err != nil {
		return "", err
	}
	return val.GetVal(), nil
}

func (m *meta) DeleteMetas(ctx context.Context, keys ...string) ([]string, error) {
	vals, err := m.mc.Load().(gmeta.MetaClient).DeleteMetas(ctx, &payload.Meta_Keys{
		Keys: keys,
	}, m.copts...)
	if err != nil {
		return nil, err
	}
	return vals.GetVals(), nil
}

func (m *meta) DeleteUUID(ctx context.Context, val string) (string, error) {
	key, err := m.mc.Load().(gmeta.MetaClient).DeleteMetaInverse(ctx, &payload.Meta_Val{
		Val: val,
	}, m.copts...)
	if err != nil {
		return "", err
	}
	return key.GetKey(), nil
}

func (m *meta) DeleteUUIDs(ctx context.Context, vals ...string) ([]string, error) {
	keys, err := m.mc.Load().(gmeta.MetaClient).DeleteMetasInverse(ctx, &payload.Meta_Vals{
		Vals: vals,
	}, m.copts...)
	if err != nil {
		return nil, err
	}
	return keys.GetKeys(), nil
}
