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

package redis

import (
	"context"
	"crypto/tls"
	"net"
	"reflect"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/vdaas/vald/internal/errors"
)

type Redis interface {
	TxPipeline() redis.Pipeliner
	Close() error
	Lister
	Getter
	Setter
	Deleter
}

var (
	a = map[string]func() Redis{
		"aa": nil,
	}
)

type Conn = redis.Conn

type redisClient struct {
	addrs              []string
	clusterSlots       func() ([]redis.ClusterSlot, error)
	db                 int
	dialTimeout        time.Duration
	dialer             func(ctx context.Context, network, addr string) (net.Conn, error)
	idleCheckFrequency time.Duration
	idleTimeout        time.Duration
	keyPref            string
	maxConnAge         time.Duration
	maxRedirects       int
	maxRetries         int
	maxRetryBackoff    time.Duration
	minIdleConns       int
	minRetryBackoff    time.Duration
	onConnect          func(*redis.Conn) error
	onNewNode          func(*redis.Client)
	password           string
	poolSize           int
	poolTimeout        time.Duration
	readOnly           bool
	readTimeout        time.Duration
	routeByLatency     bool
	routeRandomly      bool
	tlsConfig          *tls.Config
	writeTimeout       time.Duration
}

func New(ctx context.Context, opts ...Option) (Redis, error) {
	r := new(redisClient)
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	redis.NewClient(nil)
	switch len(r.addrs) {
	case 0:
		return nil, errors.ErrAddrsNotFound
	case 1:
		if len(r.addrs[0]) == 0 {
			return nil, errors.ErrAddrsNotFound
		}
		return redis.NewClient(&redis.Options{
			Addr:               r.addrs[0],
			Password:           r.password,
			Dialer:             r.dialer,
			OnConnect:          r.onConnect,
			DB:                 r.db,
			MaxRetries:         r.maxRetries,
			MinRetryBackoff:    r.minRetryBackoff,
			MaxRetryBackoff:    r.maxRetryBackoff,
			DialTimeout:        r.dialTimeout,
			ReadTimeout:        r.readTimeout,
			WriteTimeout:       r.writeTimeout,
			PoolSize:           r.poolSize,
			MinIdleConns:       r.minIdleConns,
			MaxConnAge:         r.maxConnAge,
			PoolTimeout:        r.poolTimeout,
			IdleTimeout:        r.idleTimeout,
			IdleCheckFrequency: r.idleCheckFrequency,
			TLSConfig:          r.tlsConfig,
		}), nil
	default:
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:              r.addrs,
			Dialer:             r.dialer,
			MaxRedirects:       r.maxRedirects,
			ReadOnly:           r.readOnly,
			RouteByLatency:     r.routeByLatency,
			RouteRandomly:      r.routeRandomly,
			ClusterSlots:       r.clusterSlots,
			OnNewNode:          r.onNewNode,
			OnConnect:          r.onConnect,
			Password:           r.password,
			MaxRetries:         r.maxRetries,
			MinRetryBackoff:    r.minRetryBackoff,
			MaxRetryBackoff:    r.maxRetryBackoff,
			DialTimeout:        r.dialTimeout,
			ReadTimeout:        r.readTimeout,
			WriteTimeout:       r.writeTimeout,
			PoolSize:           r.poolSize,
			MinIdleConns:       r.minIdleConns,
			MaxConnAge:         r.maxConnAge,
			PoolTimeout:        r.poolTimeout,
			IdleTimeout:        r.idleTimeout,
			IdleCheckFrequency: r.idleCheckFrequency,
			TLSConfig:          r.tlsConfig,
		}).WithContext(ctx), nil
	}
	return nil, nil
}
