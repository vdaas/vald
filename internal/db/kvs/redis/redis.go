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

package redis

import (
	"context"
	"reflect"
	"time"

	redis "github.com/go-redis/redis/v8"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/tls"
)

// Nil is a type alias of redis.Nil.
var Nil = redis.Nil

// Connector is an interface to connect to Redis servers.
type Connector interface {
	Connect(ctx context.Context) (Redis, error)
}

// Redis is an interface to communicate with Redis servers.
type Redis interface {
	TxPipeline() redis.Pipeliner
	Ping(context.Context) *StatusCmd
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
	// Dialer is a net.Dialer.
	dialer               net.Dialer
	// client is a Redis client.
	client               Redis
	// limiter is a Limiter.
	limiter              Limiter
	// tlsConfig is a tls.Config.
	tlsConfig            *tls.Config
	// dialerFunc is a function to dial.
	dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
	// clusterSlots is a function to get cluster slots.
	clusterSlots         func(context.Context) ([]redis.ClusterSlot, error)
	// onConnect is a function to call on connect.
	onConnect            func(context.Context, *redis.Conn) error
	// network is a network.
	network              string
	// username is a username.
	username             string
	// sentinelPassword is a sentinel password.
	sentinelPassword     string
	// keyPref is a key prefix.
	keyPref              string
	// sentinelMasterName is a sentinel master name.
	sentinelMasterName   string
	// password is a password.
	password             string
	// hooks is a list of hooks.
	hooks                []Hook
	// addrs is a list of addresses.
	addrs                []string
	// minIdleConns is a minimum number of idle connections.
	minIdleConns         int
	// initialPingDuration is an initial ping duration.
	initialPingDuration  time.Duration
	// minRetryBackoff is a minimum retry backoff.
	minRetryBackoff      time.Duration
	// maxRetries is a max number of retries.
	maxRetries           int
	// maxRedirects is a max number of redirects.
	maxRedirects         int
	// poolSize is a pool size.
	poolSize             int
	// poolTimeout is a pool timeout.
	poolTimeout          time.Duration
	// db is a database.
	db                   int
	// readTimeout is a read timeout.
	readTimeout          time.Duration
	// dialTimeout is a dial timeout.
	dialTimeout          time.Duration
	// idleCheckFrequency is an idle check frequency.
	idleCheckFrequency   time.Duration
	// maxConnAge is a max connection age.
	maxConnAge           time.Duration
	// initialPingTimeLimit is an initial ping time limit.
	initialPingTimeLimit time.Duration
	// maxRetryBackoff is a max retry backoff.
	maxRetryBackoff      time.Duration
	// idleTimeout is an idle timeout.
	idleTimeout          time.Duration
	// writeTimeout is a write timeout.
	writeTimeout         time.Duration
	// routeRandomly is a flag to route randomly.
	routeRandomly        bool
	// routeByLatency is a flag to route by latency.
	routeByLatency       bool
	// readOnly is a flag to read only.
	readOnly             bool
}

// New returns Connector if no error occurs.
func New(opts ...Option) (c Connector, err error) {
	r := new(redisClient)
	for _, opt := range append(defaultOptions, opts...) {
		if err = opt(r); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return r, nil
}

func (rc *redisClient) setClient(ctx context.Context) (err error) {
	switch len(rc.addrs) {
	case 0:
		return errors.ErrRedisAddrsNotFound
	case 1:
		rc.client, err = rc.newClient(ctx)
		if err != nil {
			return err
		}
	default:
		rc.client, err = rc.newClusterClient(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rc *redisClient) newClient(ctx context.Context) (c *redis.Client, err error) {
	if len(rc.addrs) == 0 || len(rc.addrs[0]) == 0 {
		return nil, errors.ErrRedisAddrsNotFound
	}

	if len(rc.sentinelMasterName) != 0 {
		if rc.routeRandomly || rc.routeByLatency {
			return nil, errors.ErrRedisInvalidOption
		}
		c = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:         rc.sentinelMasterName,
			SentinelAddrs:      rc.addrs,
			SentinelPassword:   rc.sentinelPassword,
			Username:           rc.username,
			Password:           rc.password,
			Dialer:             rc.dialerFunc,
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
		}).WithContext(ctx)
	} else {
		c = redis.NewClient(&redis.Options{
			Addr: rc.addrs[0],
			Network: func() string {
				nt := net.NetworkTypeFromString(rc.network)
				if nt == 0 || nt == net.Unknown || strings.EqualFold(nt.String(), net.Unknown.String()) {
					return net.TCP.String()
				}
				return nt.String()
			}(),
			Username:           rc.username,
			Password:           rc.password,
			Dialer:             rc.dialerFunc,
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
			Limiter:            rc.limiter,
		}).WithContext(ctx)
	}
	for _, hk := range rc.hooks {
		c.AddHook(hk)
	}

	return c, nil
}

func (rc *redisClient) newClusterClient(ctx context.Context) (c *redis.ClusterClient, err error) {
	if len(rc.addrs) == 0 || len(rc.addrs[0]) == 0 {
		return nil, errors.ErrRedisAddrsNotFound
	}

	if len(rc.sentinelMasterName) != 0 {
		c = redis.NewFailoverClusterClient(&redis.FailoverOptions{
			MasterName:         rc.sentinelMasterName,
			SentinelAddrs:      rc.addrs,
			SentinelPassword:   rc.sentinelPassword,
			Dialer:             rc.dialerFunc,
			RouteByLatency:     rc.routeByLatency,
			RouteRandomly:      rc.routeRandomly,
			OnConnect:          rc.onConnect,
			Password:           rc.password,
			Username:           rc.username,
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
	} else {
		c = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:              rc.addrs,
			Dialer:             rc.dialerFunc,
			MaxRedirects:       rc.maxRedirects,
			ReadOnly:           rc.readOnly,
			RouteByLatency:     rc.routeByLatency,
			RouteRandomly:      rc.routeRandomly,
			ClusterSlots:       rc.clusterSlots,
			OnConnect:          rc.onConnect,
			Password:           rc.password,
			Username:           rc.username,
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

	for _, hk := range rc.hooks {
		c.AddHook(hk)
	}

	return c, nil
}

// Connect returns Redis instance that has connection to servers.
func (rc *redisClient) Connect(ctx context.Context) (Redis, error) {
	if rc.dialer != nil {
		rc.dialer.StartDialerCache(ctx)
		rc.dialerFunc = rc.dialer.GetDialer()
	}

	if err := rc.setClient(ctx); err != nil {
		return nil, err
	}

	return rc.ping(ctx)
}

func (rc *redisClient) ping(ctx context.Context) (r Redis, err error) {
	pctx, cancel := context.WithTimeout(ctx, rc.initialPingTimeLimit)
	defer cancel()
	tick := time.NewTicker(rc.initialPingDuration)
	for {
		select {
		case <-pctx.Done():
			err = errors.Join(err, errors.ErrRedisConnectionPingFailed, pctx.Err())
			log.Error(err)
			return nil, err
		case <-tick.C:
			err = rc.client.Ping(ctx).Err()
			if err == nil {
				return rc.client, nil
			}
			log.Warn(err)
		}
	}
}
