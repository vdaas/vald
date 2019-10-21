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
	SetMeta(context.Context, string, string) error
	SetMetas(context.Context, map[string]string) error
	DelMeta(context.Context, string) error
	DelMetas(context.Context, ...string) error
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

	mi = m

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
	m.mc.Store(gmeta.NewMetaClient(conn))
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

func (m *meta) GetMeta(ctx context.Context, key string) (val string, err error) {
	ids, err := m.mc.Load().(gmeta.MetaClient).GetMeta(ctx, &payload.Meta_Key{
		Key: &payload.Object_ID{
			Id: key,
		},
	}, m.copts...)
	if err != nil {
		return "", err
	}
	return ids.Val.GetId(), nil
}

func (m *meta) GetMetas(ctx context.Context, keys ...string) (vals []string, err error) {
	ids, err := m.mc.Load().(gmeta.MetaClient).GetMetas(ctx, &payload.Meta_Keys{
		Keys: &payload.Object_IDs{
			Ids: func() []*payload.Object_ID {
				ids := make([]*payload.Object_ID, 0, len(keys))
				for _, key := range keys {
					ids = append(ids, &payload.Object_ID{
						Id: key,
					})
				}
				return ids
			}(),
		},
	}, m.copts...)
	if err != nil {
		return nil, err
	}
	vals = make([]string, 0, len(ids.Vals.Ids))
	for _, id := range ids.Vals.Ids {
		vals = append(vals, id.GetId())
	}
	return vals, nil
}

func (m *meta) SetMeta(ctx context.Context, key, val string) error {
	_, err := m.mc.Load().(gmeta.MetaClient).SetMeta(ctx, &payload.Meta_KeyVal{
		Key: &payload.Object_ID{
			Id: key,
		},
		Val: &payload.Object_ID{
			Id: val,
		},
	}, m.copts...)
	return err
}

func (m *meta) SetMetas(ctx context.Context, kvs map[string]string) error {
	_, err := m.mc.Load().(gmeta.MetaClient).SetMetas(ctx, &payload.Meta_KeyVals{
		Kvs: func() []*payload.Meta_KeyVal {
			data := make([]*payload.Meta_KeyVal, 0, len(kvs))
			for k, v := range kvs {
				data = append(data, &payload.Meta_KeyVal{
					Key: &payload.Object_ID{
						Id: k,
					},
					Val: &payload.Object_ID{
						Id: v,
					},
				})
			}
			return data
		}(),
	}, m.copts...)
	return err
}

func (m *meta) DelMeta(ctx context.Context, key string) error {
	_, err := m.mc.Load().(gmeta.MetaClient).DeleteMeta(ctx, &payload.Meta_Key{
		Key: &payload.Object_ID{
			Id: key,
		},
	}, m.copts...)
	return err
}

func (m *meta) DelMetas(ctx context.Context, keys ...string) error {
	_, err := m.mc.Load().(gmeta.MetaClient).GetMetas(ctx, &payload.Meta_Keys{
		Keys: &payload.Object_IDs{
			Ids: func() []*payload.Object_ID {
				ids := make([]*payload.Object_ID, 0, len(keys))
				for _, key := range keys {
					ids = append(ids, &payload.Object_ID{
						Id: key,
					})
				}
				return ids
			}(),
		},
	}, m.copts...)
	return err
}
