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

// Package service manages the main logic of server.
package service

import (
	"context"

	"github.com/kpango/gache"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/kvs/redis"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/tcp"
	"github.com/vdaas/vald/internal/tls"
)

type Redis interface {
	Connect(context.Context) error
	Disconnect() error
	Get(string) (string, error)
	GetMultiple(...string) ([]string, error)
	GetInverse(string) (string, error)
	GetInverseMultiple(...string) ([]string, error)
	Set(string, string) error
	SetMultiple(map[string]string) error
	Delete(string) (string, error)
	DeleteMultiple(...string) ([]string, error)
	DeleteInverse(string) (string, error)
	DeleteInverseMultiple(...string) ([]string, error)
}

type client struct {
	db      redis.Redis
	opts    []redis.Option
	topts   []tcp.DialerOption
	kvIndex int
	vkIndex int
}

func New(cfg *config.Redis) (Redis, error) {
	c := &client{
		kvIndex: cfg.KVIndex,
		vkIndex: cfg.VKIndex,
	}

	if c.kvIndex == c.vkIndex {
		return nil, errors.ErrRedisInvalidKVVKIndex(c.kvIndex, c.vkIndex)
	}

	c.topts = make([]tcp.DialerOption, 0, 8)
	if cfg.TCP != nil {
		if cfg.TCP.DNS.CacheEnabled {
			c.topts = append(c.topts,
				tcp.WithCache(gache.New()),
				tcp.WithEnableDNSCache(),
				tcp.WithDNSCacheExpiration(cfg.TCP.DNS.CacheExpiration),
				tcp.WithDNSRefreshDuration(cfg.TCP.DNS.RefreshDuration),
			)
		}
		if cfg.TCP.Dialer.DualStackEnabled {
			c.topts = append(c.topts, tcp.WithEnableDialerDualStack())
		} else {
			c.topts = append(c.topts, tcp.WithDisableDialerDualStack())
		}
		if cfg.TCP.TLS != nil && cfg.TCP.TLS.Enabled {
			tcfg, err := tls.New(
				tls.WithCert(cfg.TCP.TLS.Cert),
				tls.WithKey(cfg.TCP.TLS.Key),
				tls.WithCa(cfg.TCP.TLS.CA),
			)
			if err != nil {
				return nil, err
			}
			c.topts = append(c.topts, tcp.WithTLS(tcfg))
		}
		c.topts = append(c.topts,
			tcp.WithDialerKeepAlive(cfg.TCP.Dialer.KeepAlive),
			tcp.WithDialerTimeout(cfg.TCP.Dialer.Timeout),
		)
	}

	c.opts = make([]redis.Option, 0, 25)
	c.opts = append(c.opts,
		redis.WithAddrs(cfg.Addrs...),
		redis.WithDialTimeout(cfg.DialTimeout),
		redis.WithIdleCheckFrequency(cfg.IdleCheckFrequency),
		redis.WithIdleTimeout(cfg.IdleTimeout),
		redis.WithKeyPrefix(cfg.KeyPref),
		redis.WithMaximumConnectionAge(cfg.MaxConnAge),
		redis.WithRetryLimit(cfg.MaxRetries),
		redis.WithMaximumRetryBackoff(cfg.MaxRetryBackoff),
		redis.WithMinimumIdleConnection(cfg.MinIdleConns),
		redis.WithMinimumRetryBackoff(cfg.MinRetryBackoff),
		redis.WithOnConnectFunction(func(conn *redis.Conn) error {
			return nil
		}),
		// redis.WithOnNewNodeFunction(f func(*redis.Client)) ,
		redis.WithPassword(cfg.Password),
		redis.WithPoolSize(cfg.PoolSize),
		redis.WithPoolTimeout(cfg.PoolTimeout),
		// redis.WithReadOnlyFlag(readOnly bool) ,
		redis.WithReadTimeout(cfg.ReadTimeout),
		redis.WithRouteByLatencyFlag(cfg.RouteByLatency),
		redis.WithRouteRandomlyFlag(cfg.RouteRandomly),
		redis.WithWriteTimeout(cfg.WriteTimeout),
	)

	if cfg.TLS != nil && cfg.TLS.Enabled {
		tcfg, err := tls.New(
			tls.WithCert(cfg.TLS.Cert),
			tls.WithKey(cfg.TLS.Key),
			tls.WithCa(cfg.TLS.CA),
		)
		if err != nil {
			return nil, err
		}
		c.opts = append(c.opts, redis.WithTLSConfig(tcfg))
	}

	if len(cfg.Addrs) > 1 {
		c.opts = append(c.opts,
			redis.WithRedirectLimit(cfg.MaxRedirects),
		)

	} else {
		c.opts = append(c.opts,
			redis.WithDB(cfg.DB),
		)
	}

	return c, nil
}

func (c *client) Disconnect() error {
	return c.db.Close()
}

func (c *client) Connect(ctx context.Context) error {
	der := tcp.NewDialer(c.topts...)
	der.StartDialerCache(ctx)
	r, err := redis.New(ctx, append(c.opts,
		redis.WithDialer(der.GetDialer()),
	)...)

	if err != nil {
		return err
	}
	c.db = r

	return nil
}

func (c *client) Get(key string) (string, error) {
	return c.get(c.kvIndex, key)
}

func (c *client) GetMultiple(keys ...string) (vals []string, err error) {
	return c.getMulti(c.kvIndex, keys...)
}

func (c *client) GetInverse(val string) (string, error) {
	return c.get(c.vkIndex, val)
}

func (c *client) GetInverseMultiple(vals ...string) ([]string, error) {
	return c.getMulti(c.vkIndex, vals...)
}

func (c *client) get(idx int, key string) (string, error) {
	pipe := c.db.TxPipeline()
	pipe.Select(idx)
	res := pipe.Get(key)
	if _, err := pipe.Exec(); err != nil {
		return "", err
	}
	return res.Result()
}

func (c *client) getMulti(idx int, keys ...string) (vals []string, err error) {
	pipe := c.db.TxPipeline()
	pipe.Select(idx)
	res := pipe.MGet(keys...)
	if _, err = pipe.Exec(); err != nil {
		return nil, err
	}
	if err = res.Err(); err != nil {
		return nil, err
	}
	vals = make([]string, 0, len(keys))
	for _, r := range res.Val() {
		val, ok := r.(string)
		if ok {
			vals = append(vals, val)
		}
	}
	return vals[:len(vals)], nil
}

func (c *client) Set(key, val string) error {
	pipe := c.db.TxPipeline()
	pipe.Select(c.kvIndex)
	kv := pipe.Set(key, val, 0)
	pipe.Select(c.vkIndex)
	vk := pipe.Set(val, key, 0)
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	if err := kv.Err(); err != nil {
		return err
	}
	return vk.Err()
}

func (c *client) SetMultiple(kvs map[string]string) error {
	kvl := make([]interface{}, 0, len(kvs)*2)
	vkl := make([]interface{}, 0, len(kvs)*2)
	for k, v := range kvs {
		if len(k) == 0 || len(v) == 0 {
			continue
		}
		kvl = append(kvl, k, v)
		vkl = append(kvl, v, k)
	}
	pipe := c.db.TxPipeline()
	pipe.Select(c.kvIndex)
	kv := pipe.MSet(kvl...)
	pipe.Select(c.vkIndex)
	vk := pipe.MSet(vkl...)
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	if err := kv.Err(); err != nil {
		return err
	}
	return vk.Err()
}

func (c *client) Delete(key string) (string, error) {
	return c.delete(c.kvIndex, key)
}

func (c *client) DeleteMultiple(keys ...string) ([]string, error) {
	return c.deleteMulti(c.kvIndex, keys...)
}

func (c *client) DeleteInverse(val string) (string, error) {
	return c.delete(c.vkIndex, val)
}

func (c *client) DeleteInverseMultiple(vals ...string) ([]string, error) {
	return c.deleteMulti(c.vkIndex, vals...)
}

func (c *client) delete(idx int, key string) (val string, err error) {
	val, err = c.get(idx, key)
	if err != nil {
		return "", err
	}
	idxInv := c.kvIndex
	if c.kvIndex == idx {
		idxInv = c.vkIndex
	}
	pipe := c.db.TxPipeline()
	pipe.Select(idx)
	k := pipe.Del(key)
	pipe.Select(idxInv)
	v := pipe.Del(val)
	if _, err = pipe.Exec(); err != nil {
		return "", err
	}
	if err = k.Err(); err != nil {
		return "", err
	}
	if err = v.Err(); err != nil {
		return "", err
	}
	return val, nil
}

func (c *client) deleteMulti(idx int, keys ...string) (vals []string, err error) {
	vals, err = c.getMulti(idx, keys...)
	if err != nil {
		return nil, err
	}
	idxInv := c.kvIndex
	if c.kvIndex == idx {
		idxInv = c.vkIndex
	}
	pipe := c.db.TxPipeline()
	pipe.Select(idx)
	k := pipe.Del(keys...)
	pipe.Select(idxInv)
	v := pipe.Del(vals...)
	if _, err = pipe.Exec(); err != nil {
		return nil, err
	}
	if err = k.Err(); err != nil {
		return nil, err
	}
	if err = v.Err(); err != nil {
		return nil, err
	}
	return vals[:len(vals)], nil
}
