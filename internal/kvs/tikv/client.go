//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package tikv

import (
	"context"

	"github.com/tikv/client-go/config"
	"github.com/tikv/client-go/rawkv"
	"github.com/vdaas/vald/internal/errors"
)

type(
	Client interface {
		Set(key, val []byte) error
		Get(key []byte) ([]byte, error)
		Delete(key []byte) error
		Close() error
	}
	
	client struct {
		addrs []string
		rcli rawkv.Cilnet
	}
)

func New(ctx context.Context, opts ...Option) (Clinet, error) {
	var (
		c = new(client)
		err error
	)

	defer func() {
		if err != nil {
			c.Close()
		}
	}()

	for _, opt := range append(defaultOptions, opts...) {
		if err = opt(c); err != nil {
			return nil, errors.NewTiKVError("TiKV option error")
		}
	}

	c.rcli = rawkv.NewClient(ctx, c.addrs, config.DefaultConfig().Security)

	return c, nil
}

func (c *client) Set(ctx context.Context, key, val []byte) error {
	err := c.rcli.Set(ctx, key, val)
	if err != nil {
		return errors.NewTiKVError("failed to set key-value. key: %s, value: %s", key, val)
	}

	return nil
}

func (c *client) Get(ctx context.Context, key []byte) ([]byte, error) {
	val, err := c.rcli.Get(ctx, key)
	if err != nil {
		return errors.NewTiKVError("failed to get value for key. key: %s", key)
	}

	return val, nil
}

func (c *client) Delete(ctx context.Context, key []byte) error {
	err := c.rcli.Delete(ctx, key)
	if err != nil {
		return errors.NewTiKVError("failed to delete key. key: %s", key)
	}

	return nil
}

func (c *client) Close() error {
	if c.rcli != nil {
		err := c.rcli.Close()
	}
	if err != nil {
		return errors.NewTiKVError("failed to close TiKV raw client")
	}

	return nil
}
