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

package redis

import (
	"context"
	"crypto/tls"
	"reflect"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
)

var (
	// Nil is a type alias of redis.Nil.
	Nil = redis.Nil
)

// Redis is an interface to manipulate Redis server.
type Redis interface {
	TxPipeline() redis.Pipeliner
	Ping() *StatusCmd
	Close() error
	Lister
	Getter
	Setter
	Deleter
}

type (
	// Conn is a type alias of redis.Conn.
	Conn = redis.Conn
	// IntCmd is a type alias of redis.IntCmd.
	IntCmd = redis.IntCmd
	// StringCmd is a type alias of redis.StringCmd.
	StringCmd = redis.StringCmd
	// StatusCmd is a type alias of redis.StatusCmd.
	StatusCmd = redis.StatusCmd
)

type redisClient struct {
	addrs                []string
	clusterSlots         func() ([]redis.ClusterSlot, error)
	db                   int
	dialTimeout          time.Duration
	dialer               func(ctx context.Context, network, addr string) (net.Conn, error)
	idleCheckFrequency   time.Duration
	idleTimeout          time.Duration
	initialPingDuration  time.Duration
	initialPingTimeLimit time.Duration
	keyPref              string
	maxConnAge           time.Duration
	maxRedirects         int
	maxRetries           int
	maxRetryBackoff      time.Duration
	minIdleConns         int
	minRetryBackoff      time.Duration
	onConnect            func(*redis.Conn) error
	onNewNode            func(*redis.Client)
	password             string
	poolSize             int
	poolTimeout          time.Duration
	readOnly             bool
	readTimeout          time.Duration
	routeByLatency       bool
	routeRandomly        bool
	tlsConfig            *tls.Config
	writeTimeout         time.Duration
	client               Redis
}

// New returns Redis implementation if no error occurs.
func New(ctx context.Context, opts ...Option) (rc Redis, err error) {
	r := new(redisClient)
	for _, opt := range append(defaultOpts, opts...) {
		if err = opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	r, err = r.newRedisClient(ctx)
	if err != nil {
		return nil, err
	}

	return r.ping(ctx)
}

func (rc *redisClient) newRedisClient(ctx context.Context) (*redisClient, error) {
	switch len(rc.addrs) {
	case 0:
		return nil, errors.ErrRedisAddrsNotFound
	case 1:
		if len(rc.addrs[0]) == 0 {
			return nil, errors.ErrRedisAddrsNotFound
		}
		rc.client = redis.NewClient(&redis.Options{
			Addr:               rc.addrs[0],
			Password:           rc.password,
			Dialer:             rc.dialer,
			OnConnect:          rc.onConnect,
			DB:                 rc.db,
			MaxRetries:         rc.maxRetries,
			MinRetryBackoff:    rc.minRetryBackoff,
			MaxRetryBackoff:    rc.maxRetryBackoff,
			DialTimeout:        rc.dialTimeout,
			ReadTimeout:        rc.readTimeout,
			WriteTimeout:       rc.writeTimeout,
			PoolSize:           rc.poolSize,
			MinIdleConns:       rc.minIdleConns,
			MaxConnAge:         rc.maxConnAge,
			PoolTimeout:        rc.poolTimeout,
			IdleTimeout:        rc.idleTimeout,
			IdleCheckFrequency: rc.idleCheckFrequency,
			TLSConfig:          rc.tlsConfig,
		})
	default:
		rc.client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:              rc.addrs,
			Dialer:             rc.dialer,
			MaxRedirects:       rc.maxRedirects,
			ReadOnly:           rc.readOnly,
			RouteByLatency:     rc.routeByLatency,
			RouteRandomly:      rc.routeRandomly,
			ClusterSlots:       rc.clusterSlots,
			OnNewNode:          rc.onNewNode,
			OnConnect:          rc.onConnect,
			Password:           rc.password,
			MaxRetries:         rc.maxRetries,
			MinRetryBackoff:    rc.minRetryBackoff,
			MaxRetryBackoff:    rc.maxRetryBackoff,
			DialTimeout:        rc.dialTimeout,
			ReadTimeout:        rc.readTimeout,
			WriteTimeout:       rc.writeTimeout,
			PoolSize:           rc.poolSize,
			MinIdleConns:       rc.minIdleConns,
			MaxConnAge:         rc.maxConnAge,
			PoolTimeout:        rc.poolTimeout,
			IdleTimeout:        rc.idleTimeout,
			IdleCheckFrequency: rc.idleCheckFrequency,
			TLSConfig:          rc.tlsConfig,
		}).WithContext(ctx)
	}

	return rc, nil
}

func (rc *redisClient) ping(ctx context.Context) (r Redis, err error) {
	pctx, cancel := context.WithTimeout(ctx, rc.initialPingTimeLimit)
	defer cancel()
	tick := time.NewTicker(rc.initialPingDuration)
	for {
		select {
		case <-pctx.Done():
			return nil, errors.Wrap(errors.Wrap(err, errors.ErrRedisConnectionPingFailed.Error()), pctx.Err().Error())
		case <-tick.C:
			err = rc.client.Ping().Err()
			if err == nil {
				return rc.client, nil
			}
			log.Error(err)
		}
	}
}
