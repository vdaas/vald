//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package config

import (
	"context"

	"github.com/vdaas/vald/internal/db/kvs/redis"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/tls"
)

// Redis represents the configuration for redis cluster.
type Redis struct {
	// Net represents the network configuration.
	Net *Net `json:"tcp,omitempty" yaml:"net"`
	// TLS represents the TLS configuration.
	TLS *TLS `json:"tls,omitempty" yaml:"tls"`
	// MinRetryBackoff specifies the minimum duration to wait before retrying a failed operation.
	MinRetryBackoff string `json:"min_retry_backoff,omitempty" yaml:"min_retry_backoff"`
	// SentinelPassword specifies the password for Redis Sentinel authentication.
	SentinelPassword string `json:"sentinel_password,omitempty" yaml:"sentinel_password"`
	// IdleTimeout specifies the duration after which an idle connection is closed.
	IdleTimeout string `json:"idle_timeout,omitempty" yaml:"idle_timeout"`
	// InitialPingDuration specifies the timeout for the initial ping check when establishing a connection.
	InitialPingDuration string `json:"initial_ping_duration,omitempty" yaml:"initial_ping_duration"`
	// InitialPingTimeLimit specifies the total time limit for initial ping checks during startup.
	InitialPingTimeLimit string `json:"initial_ping_time_limit,omitempty" yaml:"initial_ping_time_limit"`
	// KVPrefix specifies the key prefix used for Key-Value data (e.g. metadata).
	KVPrefix string `json:"kv_prefix,omitempty" yaml:"kv_prefix"`
	// KeyPref specifies the global key prefix applied to all Redis keys.
	KeyPref string `json:"key_pref,omitempty" yaml:"key_pref"`
	// MaxConnAge specifies the maximum duration a connection can be reused before being closed.
	MaxConnAge string `json:"max_conn_age,omitempty" yaml:"max_conn_age"`
	// WriteTimeout specifies the timeout for write operations.
	WriteTimeout string `json:"write_timeout,omitempty" yaml:"write_timeout"`
	// VKPrefix specifies the key prefix used for Vector-Key data.
	VKPrefix string `json:"vk_prefix,omitempty" yaml:"vk_prefix"`
	// MaxRetryBackoff specifies the maximum duration to wait before retrying a failed operation.
	MaxRetryBackoff string `json:"max_retry_backoff,omitempty" yaml:"max_retry_backoff"`
	// Username specifies the username for Redis ACL authentication.
	Username string `json:"username,omitempty" yaml:"username"`
	// DialTimeout specifies the timeout for establishing a new connection.
	DialTimeout string `json:"dial_timeout,omitempty" yaml:"dial_timeout"`
	// Network specifies the network type (e.g., "tcp", "unix").
	Network string `json:"network,omitempty" yaml:"network"`
	// IdleCheckFrequency specifies the interval at which the pool checks for and closes idle connections.
	IdleCheckFrequency string `json:"idle_check_frequency,omitempty" yaml:"idle_check_frequency"`
	// ReadTimeout specifies the timeout for read operations.
	ReadTimeout string `json:"read_timeout,omitempty" yaml:"read_timeout"`
	// Password specifies the password for Redis authentication.
	Password string `json:"password,omitempty" yaml:"password"`
	// PrefixDelimiter specifies the delimiter used between key parts (e.g., ":").
	PrefixDelimiter string `json:"prefix_delimiter,omitempty" yaml:"prefix_delimiter"`
	// SentinelMasterName specifies the name of the master set in Redis Sentinel configuration.
	SentinelMasterName string `json:"sentinel_main_name,omitempty" yaml:"sentinel_main_name"`
	// PoolTimeout specifies the timeout for waiting for a connection from the pool.
	PoolTimeout string `json:"pool_timeout,omitempty" yaml:"pool_timeout"`
	// Addrs specifies the list of Redis server addresses (host:port).
	Addrs []string `json:"addrs,omitempty" yaml:"addrs"`
	// PoolSize specifies the maximum number of socket connections.
	PoolSize int `json:"pool_size,omitempty" yaml:"pool_size"`
	// DB specifies the Redis database index to select.
	DB int `json:"db,omitempty" yaml:"db"`
	// MinIdleConns specifies the minimum number of idle connections to maintain in the pool.
	MinIdleConns int `json:"min_idle_conns,omitempty" yaml:"min_idle_conns"`
	// MaxRetries specifies the maximum number of retries for failed operations.
	MaxRetries int `json:"max_retries,omitempty" yaml:"max_retries"`
	// MaxRedirects specifies the maximum number of redirects (e.g., MOVED, ASK) to follow in a cluster.
	MaxRedirects int `json:"max_redirects,omitempty" yaml:"max_redirects"`
	// RouteByLatency enables routing read operations to the nearest replica based on latency monitoring.
	RouteByLatency bool `json:"route_by_latency,omitempty" yaml:"route_by_latency"`
	// RouteRandomly enables routing read operations to random replicas for load balancing.
	RouteRandomly bool `json:"route_randomly,omitempty" yaml:"route_randomly"`
	// ReadOnly enables read-only mode, preventing write operations.
	ReadOnly bool `json:"read_only,omitempty" yaml:"read_only"`
}

// Bind binds the actual data from the Redis receiver fields.
func (r *Redis) Bind() *Redis {
	if r.TLS != nil {
		r.TLS.Bind()
	} else {
		r.TLS = new(TLS)
	}
	if r.Net != nil {
		r.Net.Bind()
	} else {
		r.Net = new(Net)
	}

	r.Addrs = GetActualValues(r.Addrs)
	r.DialTimeout = GetActualValue(r.DialTimeout)
	r.IdleCheckFrequency = GetActualValue(r.IdleCheckFrequency)
	r.IdleTimeout = GetActualValue(r.IdleTimeout)
	r.InitialPingDuration = GetActualValue(r.InitialPingDuration)
	r.InitialPingTimeLimit = GetActualValue(r.InitialPingTimeLimit)
	r.KVPrefix = GetActualValue(r.KVPrefix)
	r.KeyPref = GetActualValue(r.KeyPref)
	r.MaxConnAge = GetActualValue(r.MaxConnAge)
	r.MaxRetryBackoff = GetActualValue(r.MaxRetryBackoff)
	r.MinRetryBackoff = GetActualValue(r.MinRetryBackoff)
	r.Network = GetActualValue(r.Network)
	r.Password = GetActualValue(r.Password)
	r.SentinelMasterName = GetActualValue(r.SentinelMasterName)
	r.PoolTimeout = GetActualValue(r.PoolTimeout)
	r.SentinelPassword = GetActualValue(r.SentinelPassword)
	r.PrefixDelimiter = GetActualValue(r.PrefixDelimiter)
	r.ReadTimeout = GetActualValue(r.ReadTimeout)
	r.Username = GetActualValue(r.Username)
	r.VKPrefix = GetActualValue(r.VKPrefix)
	r.WriteTimeout = GetActualValue(r.WriteTimeout)
	return r
}

// Opts creates the functional option list from the Redis.
// If the error occurs, it will return no functional options and the error.
func (r *Redis) Opts() (opts []redis.Option, err error) {
	nt := net.NetworkTypeFromString(r.Network)
	if nt == 0 || nt == net.Unknown || strings.EqualFold(nt.String(), net.Unknown.String()) {
		nt = net.TCP
	}
	r.Network = nt.String()
	opts = []redis.Option{
		redis.WithAddrs(r.Addrs...),
		redis.WithDialTimeout(r.DialTimeout),
		redis.WithIdleCheckFrequency(r.IdleCheckFrequency),
		redis.WithIdleTimeout(r.IdleTimeout),
		redis.WithKeyPrefix(r.KeyPref),
		redis.WithMaximumConnectionAge(r.MaxConnAge),
		redis.WithRetryLimit(r.MaxRetries),
		redis.WithMaximumRetryBackoff(r.MaxRetryBackoff),
		redis.WithMinimumIdleConnection(r.MinIdleConns),
		redis.WithMinimumRetryBackoff(r.MinRetryBackoff),
		redis.WithOnConnectFunction(func(ctx context.Context, conn *redis.Conn) error {
			log.Infof("redis connection succeed to %s", conn.ClientGetName(ctx).String())
			return nil
		}),
		redis.WithUsername(r.Username),
		redis.WithPassword(r.Password),
		redis.WithPoolSize(r.PoolSize),
		redis.WithPoolTimeout(r.PoolTimeout),
		// In the current implementation, we do not need to use the read only flag for redis usages.
		// This implementation is to only align to the redis interface.
		// We will remove this comment out if we need to use this.
		// redis.WithReadOnlyFlag(readOnly bool) ,
		redis.WithNetwork(r.Network),
		redis.WithReadTimeout(r.ReadTimeout),
		redis.WithRouteByLatencyFlag(r.RouteByLatency),
		redis.WithRouteRandomlyFlag(r.RouteRandomly),
		redis.WithWriteTimeout(r.WriteTimeout),
		redis.WithInitialPingDuration(r.InitialPingDuration),
		redis.WithInitialPingTimeLimit(r.InitialPingTimeLimit),
		redis.WithSentinelPassword(r.SentinelPassword),
		redis.WithSentinelMasterName(r.SentinelMasterName),
	}

	if r.TLS != nil && r.TLS.Enabled {
		tls, err := tls.NewClientConfig(r.TLS.Opts()...)
		if err != nil {
			return nil, err
		}
		opts = append(opts, redis.WithTLSConfig(tls))
	}

	if r.Net != nil {
		netOpts, err := r.Net.Opts()
		if err != nil {
			return nil, err
		}
		dialer, err := net.NewDialer(netOpts...)
		if err != nil {
			return nil, err
		}
		opts = append(opts, redis.WithDialer(dialer))
	}

	if len(r.Addrs) > 1 {
		opts = append(opts,
			redis.WithRedirectLimit(r.MaxRedirects),
		)
	} else {
		opts = append(opts,
			redis.WithDB(r.DB),
		)
	}

	return opts, nil
}
