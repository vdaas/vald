// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
package tikv

import (
	"context"
	"time"

	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
	"github.com/vdaas/vald/internal/client/v1/client/meta"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

const (
	apiName      = "vald/internal/client/meta/v1/client/meta/tikv"
	waitInterval = time.Second * 2
)

type Client interface {
	meta.MetadataClient
	Close() error
}

type client struct {
	addrs []string
	c     *rawkv.Client
}

func New(opts ...Option) (Client, error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	var err error
	c.c, err = rawkv.NewClient(
		context.Background(),
		c.addrs,
		config.DefaultConfig().Security,
	)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *client) Close() error {
	return c.c.Close()
}

func (c *client) Get(ctx context.Context, key []byte) ([]byte, error) {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/Get"), apiName+"/Get")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.c.Get(ctx, key)
}

func (c *client) Put(ctx context.Context, key, val []byte) error {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/Put"), apiName+"/Put")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.c.Put(ctx, key, val)
}

func (c *client) Delete(ctx context.Context, key []byte) error {
	ctx, span := trace.StartSpan(grpc.WrapGRPCMethod(ctx, "internal/client/meta/Delete"), apiName+"/Delete")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return c.c.Delete(ctx, key)
}
