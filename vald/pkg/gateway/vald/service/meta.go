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
	"sync/atomic"
	"time"

	gmeta "github.com/vdaas/vald/apis/grpc/meta"
	"github.com/vdaas/vald/apis/grpc/payload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Meta interface {
	Start(ctx context.Context) <-chan error
	GetMeta(context.Context, string) string
	GetMetas(context.Context, ...string) []string
	SetMeta(context.Context, string) string
	SetMetas(context.Context, ...string) []string
	DelMeta(context.Context, string) error
	DelMetas(context.Context, ...string) error
}

type meta struct {
	hcDur time.Duration
	host  string
	port  int
	mc    atomic.Value
	gopts []grpc.DialOption
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
	go func() {
		tick := time.NewTicker(m.hcDur)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
			case <-tick.C:
				if conn == nil ||
					conn.GetState() == connectivity.Shutdown ||
					conn.GetState() == connectivity.TransientFailure {
					conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", m.host, m.port), m.gopts...)
					if err != nil {
						ech <- err
					}
					m.mc.Store(gmeta.NewMetaClient(conn))
				}
			}
		}
	}()
	return ech
}

func (m *meta) GetMeta(ctx context.Context, key string) string {
	m.mc.Load().(gmeta.MetaClient).GetMeta(ctx, &payload.Object_ID{
		Id: key,
	})
	return ""
}

func (m *meta) GetMetas(ctx context.Context, keys ...string) []string {
	return nil
}

func (m *meta) SetMeta(ctx context.Context, key string) string {
	return ""
}

func (m *meta) SetMetas(ctx context.Context, keys ...string) []string {
	return nil
}

func (m *meta) DelMeta(ctx context.Context, key string) error {
	return nil
}

func (m *meta) DelMetas(ctx context.Context, keys ...string) error {
	return nil
}
