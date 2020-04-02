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

// Package redis provides implementation of Go API for redis interface
package redis

import (
	"context"
	"crypto/tls"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*redisClient) error

var (
	defaultOpts = []Option{
		WithInitialPingDuration("30ms"),
		WithInitialPingTimeLimit("5m"),
	}
)

func WithDialer(der func(ctx context.Context, addr, port string) (net.Conn, error)) Option {
	return func(r *redisClient) error {
		if der != nil {
			r.dialer = der
		}
		return nil
	}
}

func WithAddrs(addrs ...string) Option {
	return func(r *redisClient) error {
		if len(addrs) == 0 {
			return nil
		}
		if r.addrs == nil {
			r.addrs = addrs
		} else {
			r.addrs = append(r.addrs, addrs...)
		}
		return nil
	}
}

func WithDB(db int) Option {
	return func(r *redisClient) error {
		r.db = db
		return nil
	}
}

func WithClusterSlots(f func() ([]redis.ClusterSlot, error)) Option {
	return func(r *redisClient) error {
		if f != nil {
			r.clusterSlots = f
		}
		return nil
	}
}

func WithDialTimeout(dur string) Option {
	return func(r *redisClient) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute * 30
		}
		r.dialTimeout = d
		return nil
	}
}

func WithIdleCheckFrequency(dur string) Option {
	return func(r *redisClient) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute
		}
		r.idleCheckFrequency = d
		return nil
	}
}

func WithIdleTimeout(dur string) Option {
	return func(r *redisClient) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute
		}
		r.idleTimeout = d
		return nil
	}
}

func WithKeyPrefix(prefix string) Option {
	return func(r *redisClient) error {
		if prefix != "" {
			r.keyPref = prefix
		}
		return nil
	}
}

func WithMaximumConnectionAge(dur string) Option {
	return func(r *redisClient) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return nil
		}
		r.maxConnAge = d
		return nil
	}
}

func WithRedirectLimit(maxRedirects int) Option {
	return func(r *redisClient) error {
		r.maxRedirects = maxRedirects
		return nil
	}
}

func WithRetryLimit(maxRetries int) Option {
	return func(r *redisClient) error {
		r.maxRetries = maxRetries
		return nil
	}
}

func WithMaximumRetryBackoff(dur string) Option {
	return func(r *redisClient) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute * 2
		}
		r.maxRetryBackoff = d
		return nil
	}
}

func WithMinimumIdleConnection(minIdleConns int) Option {
	return func(r *redisClient) error {
		r.minIdleConns = minIdleConns
		return nil
	}
}

func WithMinimumRetryBackoff(dur string) Option {
	return func(r *redisClient) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Millisecond * 5
		}
		r.minRetryBackoff = d
		return nil
	}
}

func WithOnConnectFunction(f func(*redis.Conn) error) Option {
	return func(r *redisClient) error {
		if f != nil {
			r.onConnect = f
		}
		return nil
	}
}

func WithOnNewNodeFunction(f func(*redis.Client)) Option {
	return func(r *redisClient) error {
		if f != nil {
			r.onNewNode = f
		}
		return nil
	}
}

func WithPassword(password string) Option {
	return func(r *redisClient) error {
		if password != "" {
			r.password = password
		}
		return nil
	}
}

func WithPoolSize(poolSize int) Option {
	return func(r *redisClient) error {
		r.poolSize = poolSize
		return nil
	}
}

func WithPoolTimeout(dur string) Option {
	return func(r *redisClient) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute * 5
		}
		r.poolTimeout = d
		return nil
	}
}

func WithReadOnlyFlag(readOnly bool) Option {
	return func(r *redisClient) error {
		r.readOnly = readOnly
		return nil
	}
}

func WithReadTimeout(dur string) Option {
	return func(r *redisClient) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute
		}
		r.readTimeout = d
		return nil
	}
}

func WithRouteByLatencyFlag(routeByLatency bool) Option {
	return func(r *redisClient) error {
		r.routeByLatency = routeByLatency
		return nil
	}
}

func WithRouteRandomlyFlag(routeRandomly bool) Option {
	return func(r *redisClient) error {
		r.routeRandomly = routeRandomly
		return nil
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(r *redisClient) error {
		if cfg != nil {
			r.tlsConfig = cfg
		}
		return nil
	}
}

func WithWriteTimeout(dur string) Option {
	return func(r *redisClient) error {
		if dur == "" {
			return nil
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			d = time.Minute
		}
		r.writeTimeout = d
		return nil
	}
}

func WithInitialPingTimeLimit(lim string) Option {
	return func(r *redisClient) error {
		if lim == "" {
			return nil
		}
		pd, err := timeutil.Parse(lim)
		if err != nil {
			pd = time.Second * 30
		}
		r.initialPingTimeLimit = pd
		return nil
	}
}

func WithInitialPingDuration(dur string) Option {
	return func(r *redisClient) error {
		if dur == "" {
			return nil
		}
		pd, err := timeutil.Parse(dur)
		if err != nil {
			pd = time.Millisecond * 50
		}
		r.initialPingDuration = pd
		return nil
	}
}
