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

// Package service manages the main logic of server.
package service

import (
	"context"

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
	db              redis.Redis
	opts            []redis.Option
	topts           []tcp.DialerOption
	kvPrefix        string
	vkPrefix        string
	prefixDelimiter string
}

func New(cfg *config.Redis) (Redis, error) {
	c := &client{
		kvPrefix:        cfg.KVPrefix,
		vkPrefix:        cfg.VKPrefix,
		prefixDelimiter: cfg.PrefixDelimiter,
	}

	if c.kvPrefix == "" {
		c.kvPrefix = "kv"
	}
	if c.vkPrefix == "" {
		c.vkPrefix = "vk"
	}
	if c.kvPrefix == c.vkPrefix {
		return nil, errors.ErrRedisInvalidKVVKPrefix(c.kvPrefix, c.vkPrefix)
	}
	if c.prefixDelimiter == "" {
		c.prefixDelimiter = "-"
	}

	c.topts = make([]tcp.DialerOption, 0, 8)
	if cfg.TCP != nil {
		if cfg.TCP.DNS != nil && cfg.TCP.DNS.CacheEnabled {
			c.topts = append(c.topts,
				tcp.WithEnableDNSCache(),
				tcp.WithDNSCacheExpiration(cfg.TCP.DNS.CacheExpiration),
				tcp.WithDNSRefreshDuration(cfg.TCP.DNS.RefreshDuration),
			)
		}
		if cfg.TCP.Dialer != nil && cfg.TCP.Dialer.DualStackEnabled {
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
		redis.WithInitialPingDuration(cfg.InitialPingDuration),
		redis.WithInitialPingTimeLimit(cfg.InitialPingTimeLimit),
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
	der, err := tcp.NewDialer(c.topts...)
	if err != nil {
		return err
	}
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
	return c.get(c.kvPrefix, key)
}

func (c *client) GetMultiple(keys ...string) (vals []string, err error) {
	return c.getMulti(c.kvPrefix, keys...)
}

func (c *client) GetInverse(val string) (string, error) {
	return c.get(c.vkPrefix, val)
}

func (c *client) GetInverseMultiple(vals ...string) ([]string, error) {
	return c.getMulti(c.vkPrefix, vals...)
}

func (c *client) appendPrefix(prefix, key string) string {
	return prefix + c.prefixDelimiter + key
}

func (c *client) get(prefix, key string) (val string, err error) {
	pipe := c.db.TxPipeline()
	res := pipe.Get(c.appendPrefix(prefix, key))
	_, err = pipe.Exec()
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

func (c *client) getMulti(prefix string, keys ...string) (vals []string, err error) {
	pipe := c.db.TxPipeline()
	ress := make(map[string]*redis.StringCmd, len(keys))
	for _, k := range keys {
		ress[k] = pipe.Get(c.appendPrefix(prefix, k))
	}
	if _, err = pipe.Exec(); err != nil {
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
				err = errors.Wrap(err, errors.ErrRedisNotFound(k).Error())
			} else {
				err = errors.Wrap(err, errors.ErrRedisGetOperationFailed(k, err).Error())
			}
			vals = append(vals, "")
			continue
		}
		vals = append(vals, res.Val())
	}
	return vals, err
}

func (c *client) Set(key, val string) (err error) {
	kvKey := c.appendPrefix(c.kvPrefix, key)
	vkKey := c.appendPrefix(c.vkPrefix, val)
	pipe := c.db.TxPipeline()
	kv := pipe.Set(kvKey, val, 0)
	vk := pipe.Set(vkKey, key, 0)
	_, err = pipe.Exec()
	if err != nil {
		return err
	}
	err = kv.Err()
	if err != nil {
		return errors.Wrap(c.db.Del(vkKey).Err(), errors.ErrRedisSetOperationFailed(kvKey, err).Error())
	}
	err = vk.Err()
	if err != nil {
		return errors.Wrap(c.db.Del(kvKey).Err(), errors.ErrRedisSetOperationFailed(vkKey, err).Error())
	}
	return nil
}

func (c *client) SetMultiple(kvs map[string]string) (err error) {
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
		kvress[vkKey] = pipe.Set(kvKey, v, 0)
		vkress[kvKey] = pipe.Set(vkKey, k, 0)
	}
	if _, err = pipe.Exec(); err != nil {
		return err
	}
	for vkKey, res := range kvress {
		if err = res.Err(); err != nil {
			err = errors.Wrap(c.db.Del(vkKey).Err(), errors.ErrRedisSetOperationFailed(vks[vkKey], err).Error())
		}
	}
	for kvKey, res := range vkress {
		if err = res.Err(); err != nil {
			err = errors.Wrap(c.db.Del(kvKey).Err(), errors.ErrRedisSetOperationFailed(kvs[kvKey], err).Error())
		}
	}
	return err
}

func (c *client) Delete(key string) (string, error) {
	return c.delete(c.kvPrefix, c.vkPrefix, key)
}

func (c *client) DeleteMultiple(keys ...string) ([]string, error) {
	return c.deleteMulti(c.kvPrefix, c.vkPrefix, keys...)
}

func (c *client) DeleteInverse(val string) (string, error) {
	return c.delete(c.vkPrefix, c.kvPrefix, val)
}

func (c *client) DeleteInverseMultiple(vals ...string) ([]string, error) {
	return c.deleteMulti(c.vkPrefix, c.kvPrefix, vals...)
}

func (c *client) delete(pfx, pfxInv, key string) (val string, err error) {
	val, err = c.get(pfx, key)
	if err != nil {
		return "", err
	}
	pipe := c.db.TxPipeline()
	k := pipe.Del(c.appendPrefix(pfx, key))
	v := pipe.Del(c.appendPrefix(pfxInv, val))
	if _, err = pipe.Exec(); err != nil {
		return "", err
	}
	if err = k.Err(); err != nil {
		if pfx == c.kvPrefix {
			return "", errors.Wrap(c.Set(key, val), err.Error())
		}
		return "", errors.Wrap(c.Set(val, key), err.Error())
	}
	if err = v.Err(); err != nil {
		if pfx == c.kvPrefix {
			return "", errors.Wrap(c.Set(key, val), err.Error())
		}
		return "", errors.Wrap(c.Set(val, key), err.Error())

	}
	return val, nil
}

func (c *client) deleteMulti(pfx, pfxInv string, keys ...string) (vals []string, err error) {
	vals, err = c.getMulti(pfx, keys...)
	if err != nil {
		return nil, err
	}
	pipe := c.db.TxPipeline()
	ress := make(map[string]*redis.IntCmd, len(keys)*2)
	for _, k := range keys {
		key := c.appendPrefix(pfx, k)
		ress[key] = pipe.Del(key)
	}
	for _, v := range vals {
		key := c.appendPrefix(pfxInv, v)
		ress[key] = pipe.Del(key)
	}
	if _, err = pipe.Exec(); err != nil {
		return nil, err
	}
	var errs error
	for k, res := range ress {
		if err = res.Err(); err != nil {
			errs = errors.Wrap(errs, errors.ErrRedisDeleteOperationFailed(k, err).Error())
			continue
		}
	}
	return vals, errs
}
