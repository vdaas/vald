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

// Package redis provides implementation of Go API for redis interface
package redis

import (
	"context"
	"crypto/tls"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/tcp"
	"github.com/vdaas/vald/internal/timeutil"
)

// Option represents the functional option for redisClient.
type Option func(*redisClient) error

var defaultOptions = []Option{
	WithInitialPingDuration("30ms"),
	WithInitialPingTimeLimit("5m"),
	WithNetwork("tcp"),
}

// WithNetwork returns the option to set the network like tcp or unix.
func WithNetwork(network string) Option {
	return func(r *redisClient) error {
		if network != "" {
			r.network = network
		}
		return nil
	}
}

// WithDialer returns the option to set the dialer.
func WithDialer(der tcp.Dialer) Option {
	return func(r *redisClient) error {
		if der != nil {
			r.dialer = der
		}
		return nil
	}
}

// WithDialerFunc returns the option to set the dialer func.
func WithDialerFunc(der func(ctx context.Context, addr, port string) (net.Conn, error)) Option {
	return func(r *redisClient) error {
		if der != nil {
			r.dialerFunc = der
		}
		return nil
	}
}

// WithAddrs returns the option to set the addrs.
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

// WithDB returns the option to set the db.
func WithDB(db int) Option {
	return func(r *redisClient) error {
		r.db = db
		return nil
	}
}

// WithClusterSlots returns the option to set the clusterSlots.
func WithClusterSlots(f func(context.Context) ([]redis.ClusterSlot, error)) Option {
	return func(r *redisClient) error {
		if f != nil {
			r.clusterSlots = f
		}
		return nil
	}
}

// WithDialTimeout returns the option to set the dialTimeout.
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

// WithIdleCheckFrequency returns the option to set the idleCheckFrequency.
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

// WithIdleTimeout returns the option to set the idleTimeout.
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

// WithKeyPrefix returns the option to set the keyPref.
func WithKeyPrefix(prefix string) Option {
	return func(r *redisClient) error {
		if prefix != "" {
			r.keyPref = prefix
		}
		return nil
	}
}

// WithMaximumConnectionAge returns the option to set the maxConnAge.
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

// WithRedirectLimit returns the option to set the maxRedirects.
func WithRedirectLimit(maxRedirects int) Option {
	return func(r *redisClient) error {
		r.maxRedirects = maxRedirects
		return nil
	}
}

// WithRetryLimit returns the option to set the maxRetries.
func WithRetryLimit(maxRetries int) Option {
	return func(r *redisClient) error {
		r.maxRetries = maxRetries
		return nil
	}
}

// WithMaximumRetryBackoff returns the option to set the maxRetryBackoff.
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

// WithMinimumIdleConnection returns the option to set the minIdleConns.
func WithMinimumIdleConnection(minIdleConns int) Option {
	return func(r *redisClient) error {
		r.minIdleConns = minIdleConns
		return nil
	}
}

// WithMinimumRetryBackoff returns the option to set the minRetryBackoff.
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

// WithOnConnectFunction returns the option to set the onConnect.
func WithOnConnectFunction(f func(context.Context, *redis.Conn) error) Option {
	return func(r *redisClient) error {
		if f != nil {
			r.onConnect = f
		}
		return nil
	}
}

// WithUsername returns the option to set the username.
func WithUsername(name string) Option {
	return func(r *redisClient) error {
		if name != "" {
			r.username = name
		}
		return nil
	}
}

// WithSentinelMasterName returns the option to set the password.
func WithSentinelMasterName(name string) Option {
	return func(r *redisClient) error {
		if name != "" {
			r.sentinelMasterName = name
		}
		return nil
	}
}

// WithSentinelPassword returns the option to set the password.
func WithSentinelPassword(password string) Option {
	return func(r *redisClient) error {
		if password != "" {
			r.sentinelPassword = password
		}
		return nil
	}
}

// WithPassword returns the option to set the password.
func WithPassword(password string) Option {
	return func(r *redisClient) error {
		if password != "" {
			r.password = password
		}
		return nil
	}
}

// WithPoolSize returns the option to set the poolSize.
func WithPoolSize(poolSize int) Option {
	return func(r *redisClient) error {
		r.poolSize = poolSize
		return nil
	}
}

// WithPoolTimeout returns the option to set the poolTimeout.
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

// WithReadOnlyFlag returns the option to set the readOnly.
func WithReadOnlyFlag(readOnly bool) Option {
	return func(r *redisClient) error {
		r.readOnly = readOnly
		return nil
	}
}

// WithReadTimeout returns the option to set the readTimeout.
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

// WithRouteByLatencyFlag returns the option to set the routeByLatency.
func WithRouteByLatencyFlag(routeByLatency bool) Option {
	return func(r *redisClient) error {
		r.routeByLatency = routeByLatency
		return nil
	}
}

// WithRouteRandomlyFlag returns the option to set the routeRandomly.
func WithRouteRandomlyFlag(routeRandomly bool) Option {
	return func(r *redisClient) error {
		r.routeRandomly = routeRandomly
		return nil
	}
}

// WithTLSConfig returns the option to set the tlsConfig.
func WithTLSConfig(cfg *tls.Config) Option {
	return func(r *redisClient) error {
		if cfg != nil {
			r.tlsConfig = cfg
		}
		return nil
	}
}

// WithWriteTimeout returns the option to set the writeTimeout.
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

// WithInitialPingTimeLimit returns the option to set the initialPingTimeLimit.
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

// WithInitialPingDuration returns the option to set the initialPingDuration.
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

// WithHooks returns the option to add hooks.
func WithHooks(hooks ...Hook) Option {
	return func(r *redisClient) error {
		if hooks == nil {
			return nil
		}

		if r.hooks != nil {
			r.hooks = append(r.hooks, hooks...)
			return nil
		}

		r.hooks = hooks

		return nil
	}
}

// WithLimiter returns the option to limiter.
func WithLimiter(limiter Limiter) Option {
	return func(r *redisClient) error {
		if limiter == nil {
			return nil
		}
		r.limiter = limiter

		return nil
	}
}
