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

	"github.com/tikv/client-go/v2/config"
	"github.com/tikv/client-go/v2/rawkv"
	"github.com/tikv/client-go/v2/txnkv"
	"github.com/vdaas/vald/internal/errors"
)

type (
	Client interface {
		Set(ctx context.Context, key, val []byte) error
		Get(ctx context.Context, key []byte) ([]byte, error)
		Delete(ctx context.Context, key []byte) error
		Close() error
	}

	client struct {
		addrs []string
		txn   bool
		rcli  *rawkv.Client
		tcli  *txnkv.Client
	}
)

func New(ctx context.Context, opts ...Option) (Client, error) {
	var (
		c   = new(client)
		err error
	)

	defer func() {
		if err != nil {
			c.Close()
		}
	}()

	for _, opt := range append(defaultOptions, opts...) {
		if err = opt(c); err != nil {
			return nil, errors.ErrTiKVOptionFailed(err)
		}
	}

	if c.txn {
		c.tcli, err = txnkv.NewClient(c.addrs)
		if err != nil {
			return nil, errors.ErrNewTiKVTxnClientFailed(err)
		}
	}	else {
		c.rcli, err = rawkv.NewClient(ctx, c.addrs, config.DefaultConfig().Security)
		if err != nil {
			return nil, errors.ErrNewTiKVRawClientFailed(err)
		}
	}
	return c, nil
}

func (c *client) Set(ctx context.Context, key, val []byte) error {
	if c.txn {
		txn, err := c.tcli.Begin()
		if err != nil {
			return errors.ErrTiKVBeginOperationFailed(err)
		}
		txn.SetEnableAsyncCommit(true)
		err = txn.Set(key, val)
		if err != nil {
			return errors.ErrTiKVSetOperationFailed(key, val, err)
		}
		err = txn.Commit(ctx)
		if err != nil {
			return errors.ErrTiKVCommitOperationFailed(err)
		}
	} else {
		err := c.rcli.Put(ctx, key, val)
		if err != nil {
			return errors.ErrTiKVSetOperationFailed(key, val, err)
		}
	}

	return nil
}

func (c *client) Get(ctx context.Context, key []byte) (val []byte, err error) {
	if c.txn {
		txn, err := c.tcli.Begin()
		if err != nil {
			return nil, errors.ErrTiKVBeginOperationFailed(err)
		}
		val, err = txn.Get(ctx, key)
		if err != nil {
			return nil, errors.ErrTiKVGetOperationFailed(key, err)
		}
	} else {
		val, err = c.rcli.Get(ctx, key)
		if err != nil {
			return nil, errors.ErrTiKVGetOperationFailed(key, err)
		}
	}

	return val, nil
}

func (c *client) Delete(ctx context.Context, key []byte) error {
	if c.txn {
		txn, err := c.tcli.Begin()
		if err != nil {
			return errors.ErrTiKVBeginOperationFailed(err)
		}
		txn.SetEnableAsyncCommit(true)
		err = txn.Delete(key)
		if err != nil {
			return errors.ErrTiKVDeleteOperationFailed(key, err)
		}
		err = txn.Commit(ctx)
		if err != nil {
			return errors.ErrTiKVCommitOperationFailed(err)
		}
	} else {
		err := c.rcli.Delete(ctx, key)
		if err != nil {
			return errors.ErrTiKVDeleteOperationFailed(key, err)
		}
	}

	return nil
}

func (c *client) Close() error {
	if c.rcli != nil {
		if err := c.rcli.Close(); err != nil {
			return errors.ErrTiKVRawClientCloseOperationFailed(err)
		}
	}

	if c.tcli != nil {
		if err := c.tcli.Close(); err != nil {
			return errors.ErrTiKVTxnClientCloseOperationFailed(err)
		}
	}

	return nil
}
