//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package service manages the main logic of server.
package service

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/db/kvs/redis"
	"github.com/vdaas/vald/internal/errors"
)

type Redis interface {
	Connect(context.Context) error
	Disconnect() error
	Get(context.Context, string) (string, error)
	GetMultiple(context.Context, ...string) ([]string, error)
	GetInverse(context.Context, string) (string, error)
	GetInverseMultiple(context.Context, ...string) ([]string, error)
	Set(context.Context, string, string) error
	SetMultiple(context.Context, map[string]string) error
	Delete(context.Context, string) (string, error)
	DeleteMultiple(context.Context, ...string) ([]string, error)
	DeleteInverse(context.Context, string) (string, error)
	DeleteInverseMultiple(context.Context, ...string) ([]string, error)
}

type client struct {
	connector       redis.Connector
	db              redis.Redis
	kvPrefix        string
	vkPrefix        string
	prefixDelimiter string
}

func New(opts ...Option) (Redis, error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

func (c *client) Disconnect() error {
	return c.db.Close()
}

func (c *client) Connect(ctx context.Context) (err error) {
	if c.connector != nil {
		c.db, err = c.connector.Connect(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *client) Get(ctx context.Context, key string) (string, error) {
	return c.get(ctx, c.kvPrefix, key)
}

func (c *client) GetMultiple(ctx context.Context, keys ...string) (vals []string, err error) {
	return c.getMulti(ctx, c.kvPrefix, keys...)
}

func (c *client) GetInverse(ctx context.Context, val string) (string, error) {
	return c.get(ctx, c.vkPrefix, val)
}

func (c *client) GetInverseMultiple(ctx context.Context, vals ...string) ([]string, error) {
	return c.getMulti(ctx, c.vkPrefix, vals...)
}

func (c *client) appendPrefix(prefix, key string) string {
	return prefix + c.prefixDelimiter + key
}

func (c *client) get(ctx context.Context, prefix, key string) (val string, err error) {
	pipe := c.db.TxPipeline()
	res := pipe.Get(ctx, c.appendPrefix(prefix, key))
	_, err = pipe.Exec(ctx)
	if err != nil {
		if err == redis.Nil {
			return "", errors.ErrRedisNotFound(key)
		}
		return "", errors.ErrRedisGetOperationFailed(key, err)
	}
	err = res.Err()
	if err != nil {
		if err == redis.Nil {
			return "", errors.ErrRedisNotFound(key)
		}
		return "", errors.ErrRedisGetOperationFailed(key, err)
	}
	return res.Val(), nil
}

func (c *client) getMulti(ctx context.Context, prefix string, keys ...string) (vals []string, err error) {
	pipe := c.db.TxPipeline()
	ress := make(map[string]*redis.StringCmd, len(keys))
	for _, k := range keys {
		ress[k] = pipe.Get(ctx, c.appendPrefix(prefix, k))
	}
	if _, err = pipe.Exec(ctx); err != nil {
		for _, key := range keys {
			err = errors.Wrap(errors.ErrRedisGetOperationFailed(key, err), err.Error())
		}
		return nil, err
	}
	vals = make([]string, 0, len(ress))
	for _, k := range keys {
		res := ress[k]
		err = res.Err()
		if err != nil {
			if err == redis.Nil {
				err = errors.Wrap(errors.ErrRedisNotFound(k), err.Error())
			} else {
				err = errors.Wrap(errors.ErrRedisGetOperationFailed(k, err), err.Error())
			}
			vals = append(vals, "")
			continue
		}
		vals = append(vals, res.Val())
	}
	return vals, err
}

func (c *client) Set(ctx context.Context, key, val string) (err error) {
	kvKey := c.appendPrefix(c.kvPrefix, key)
	vkKey := c.appendPrefix(c.vkPrefix, val)
	pipe := c.db.TxPipeline()
	kv := pipe.Set(ctx, kvKey, val, 0)
	vk := pipe.Set(ctx, vkKey, key, 0)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}
	err = kv.Err()
	if err != nil {
		err = errors.ErrRedisSetOperationFailed(kvKey, err)
		dberr := c.db.Del(ctx, vkKey).Err()
		if dberr != nil {
			err = errors.Wrap(err, dberr.Error())
		}
		return err
	}
	err = vk.Err()
	if err != nil {
		err = errors.ErrRedisSetOperationFailed(vkKey, err)
		dberr := c.db.Del(ctx, kvKey).Err()
		if dberr != nil {
			err = errors.Wrap(err, dberr.Error())
		}
		return err
	}
	return nil
}

func (c *client) SetMultiple(ctx context.Context, kvs map[string]string) (err error) {
	pipe := c.db.TxPipeline()
	vks := make(map[string]string, len(kvs))
	kvress := make(map[string]*redis.StatusCmd, len(kvs))
	vkress := make(map[string]*redis.StatusCmd, len(kvs))
	for k, v := range kvs {
		if len(k) == 0 || len(v) == 0 {
			continue
		}
		vks[v] = k
		kvKey := c.appendPrefix(c.kvPrefix, k)
		vkKey := c.appendPrefix(c.vkPrefix, v)
		kvress[vkKey] = pipe.Set(ctx, kvKey, v, 0)
		vkress[kvKey] = pipe.Set(ctx, vkKey, k, 0)
	}
	if _, err = pipe.Exec(ctx); err != nil {
		return err
	}
	for vkKey, res := range kvress {
		if err = res.Err(); err != nil {
			err = errors.ErrRedisSetOperationFailed(vks[vkKey], err)
			dberr := c.db.Del(ctx, vkKey).Err()
			if dberr != nil {
				err = errors.Wrap(err, dberr.Error())
			}
		}
	}
	for kvKey, res := range vkress {
		if err = res.Err(); err != nil {
			err = errors.ErrRedisSetOperationFailed(kvs[kvKey], err)
			dberr := c.db.Del(ctx, kvKey).Err()
			if dberr != nil {
				err = errors.Wrap(err, dberr.Error())
			}
		}
	}
	return err
}

func (c *client) Delete(ctx context.Context, key string) (string, error) {
	return c.delete(ctx, c.kvPrefix, c.vkPrefix, key)
}

func (c *client) DeleteMultiple(ctx context.Context, keys ...string) ([]string, error) {
	return c.deleteMulti(ctx, c.kvPrefix, c.vkPrefix, keys...)
}

func (c *client) DeleteInverse(ctx context.Context, val string) (string, error) {
	return c.delete(ctx, c.vkPrefix, c.kvPrefix, val)
}

func (c *client) DeleteInverseMultiple(ctx context.Context, vals ...string) ([]string, error) {
	return c.deleteMulti(ctx, c.vkPrefix, c.kvPrefix, vals...)
}

func (c *client) delete(ctx context.Context, pfx, pfxInv, key string) (val string, err error) {
	val, err = c.get(ctx, pfx, key)
	if err != nil {
		return "", err
	}
	pipe := c.db.TxPipeline()
	k := pipe.Del(ctx, c.appendPrefix(pfx, key))
	v := pipe.Del(ctx, c.appendPrefix(pfxInv, val))
	if _, err = pipe.Exec(ctx); err != nil {
		return "", err
	}
	if err = k.Err(); err != nil {
		if pfx == c.kvPrefix {
			return "", errors.Wrap(c.Set(ctx, key, val), err.Error())
		}
		return "", errors.Wrap(c.Set(ctx, val, key), err.Error())
	}
	if err = v.Err(); err != nil {
		if pfx == c.kvPrefix {
			return "", errors.Wrap(c.Set(ctx, key, val), err.Error())
		}
		return "", errors.Wrap(c.Set(ctx, val, key), err.Error())
	}
	return val, nil
}

func (c *client) deleteMulti(ctx context.Context, pfx, pfxInv string, keys ...string) (vals []string, err error) {
	vals, err = c.getMulti(ctx, pfx, keys...)
	if err != nil {
		return nil, err
	}
	pipe := c.db.TxPipeline()
	ress := make(map[string]*redis.IntCmd, len(keys)*2)
	for _, k := range keys {
		key := c.appendPrefix(pfx, k)
		ress[key] = pipe.Del(ctx, key)
	}
	for _, v := range vals {
		key := c.appendPrefix(pfxInv, v)
		ress[key] = pipe.Del(ctx, key)
	}
	if _, err = pipe.Exec(ctx); err != nil {
		return nil, err
	}
	var errs error
	for k, res := range ress {
		if err = res.Err(); err != nil {
			err = errors.ErrRedisDeleteOperationFailed(k, err)
			if errs == nil {
				errs = err
			} else {
				errs = errors.Wrap(err, errs.Error())
			}
			continue
		}
	}
	return vals, errs
}
