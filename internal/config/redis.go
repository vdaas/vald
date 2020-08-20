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

// Package config providers configuration type and load configuration logic
package config

import (
	"github.com/vdaas/vald/internal/db/kvs/redis"
	"github.com/vdaas/vald/internal/net/tcp"
	"github.com/vdaas/vald/internal/tls"
)

type Redis struct {
	Addrs                []string `json:"addrs" yaml:"addrs"`
	DB                   int      `json:"db" yaml:"db"`
	DialTimeout          string   `json:"dial_timeout" yaml:"dial_timeout"`
	IdleCheckFrequency   string   `json:"idle_check_frequency" yaml:"idle_check_frequency"`
	IdleTimeout          string   `json:"idle_timeout" yaml:"idle_timeout"`
	InitialPingTimeLimit string   `json:"initial_ping_time_limit" yaml:"initial_ping_time_limit"`
	InitialPingDuration  string   `json:"initial_ping_duration" yaml:"initial_ping_duration"`
	KeyPref              string   `json:"key_pref" yaml:"key_pref"`
	MaxConnAge           string   `json:"max_conn_age" yaml:"max_conn_age"`
	MaxRedirects         int      `json:"max_redirects" yaml:"max_redirects"`
	MaxRetries           int      `json:"max_retries" yaml:"max_retries"`
	MaxRetryBackoff      string   `json:"max_retry_backoff" yaml:"max_retry_backoff"`
	MinIdleConns         int      `json:"min_idle_conns" yaml:"min_idle_conns"`
	MinRetryBackoff      string   `json:"min_retry_backoff" yaml:"min_retry_backoff"`
	Password             string   `json:"password" yaml:"password"`
	PoolSize             int      `json:"pool_size" yaml:"pool_size"`
	PoolTimeout          string   `json:"pool_timeout" yaml:"pool_timeout"`
	ReadOnly             bool     `json:"read_only" yaml:"read_only"`
	ReadTimeout          string   `json:"read_timeout" yaml:"read_timeout"`
	RouteByLatency       bool     `json:"route_by_latency" yaml:"route_by_latency"`
	RouteRandomly        bool     `json:"route_randomly" yaml:"route_randomly"`
	TLS                  *TLS     `json:"tls" yaml:"tls"`
	TCP                  *TCP     `json:"tcp" yaml:"tcp"`
	WriteTimeout         string   `json:"write_timeout" yaml:"write_timeout"`
	KVPrefix             string   `json:"kv_prefix" yaml:"kv_prefix"`
	VKPrefix             string   `json:"vk_prefix" yaml:"vk_prefix"`
	PrefixDelimiter      string   `json:"prefix_delimiter" yaml:"prefix_delimiter"`
}

func (r *Redis) Bind() *Redis {
	if r.TLS != nil {
		r.TLS.Bind()
	} else {
		r.TLS = new(TLS)
	}
	if r.TCP != nil {
		r.TCP.Bind()
	} else {
		r.TCP = new(TCP)
	}

	r.Addrs = GetActualValues(r.Addrs)
	r.DialTimeout = GetActualValue(r.DialTimeout)
	r.DialTimeout = GetActualValue(r.DialTimeout)
	r.IdleCheckFrequency = GetActualValue(r.IdleCheckFrequency)
	r.IdleTimeout = GetActualValue(r.IdleTimeout)
	r.KeyPref = GetActualValue(r.KeyPref)
	r.MaxConnAge = GetActualValue(r.MaxConnAge)
	r.MaxRetryBackoff = GetActualValue(r.MaxRetryBackoff)
	r.MinRetryBackoff = GetActualValue(r.MinRetryBackoff)
	r.Password = GetActualValue(r.Password)
	r.PoolTimeout = GetActualValue(r.PoolTimeout)
	r.ReadTimeout = GetActualValue(r.ReadTimeout)
	r.WriteTimeout = GetActualValue(r.WriteTimeout)
	r.KVPrefix = GetActualValue(r.KVPrefix)
	r.VKPrefix = GetActualValue(r.VKPrefix)
	r.PrefixDelimiter = GetActualValue(r.PrefixDelimiter)
	r.InitialPingTimeLimit = GetActualValue(r.InitialPingTimeLimit)
	r.InitialPingDuration = GetActualValue(r.InitialPingDuration)
	return r
}

func (r *Redis) Opts() (opts []redis.Option, err error) {
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
		redis.WithOnConnectFunction(func(conn *redis.Conn) error {
			return nil
		}),
		// redis.WithOnNewNodeFunction(f func(*redis.Client)) ,
		redis.WithPassword(r.Password),
		redis.WithPoolSize(r.PoolSize),
		redis.WithPoolTimeout(r.PoolTimeout),
		// redis.WithReadOnlyFlag(readOnly bool) ,
		redis.WithReadTimeout(r.ReadTimeout),
		redis.WithRouteByLatencyFlag(r.RouteByLatency),
		redis.WithRouteRandomlyFlag(r.RouteRandomly),
		redis.WithWriteTimeout(r.WriteTimeout),
		redis.WithInitialPingDuration(r.InitialPingDuration),
		redis.WithInitialPingTimeLimit(r.InitialPingTimeLimit),
	}

	if r.TLS != nil && r.TLS.Enabled {
		tls, err := tls.New(r.TLS.Opts()...)
		if err != nil {
			return nil, err
		}
		opts = append(opts, redis.WithTLSConfig(tls))
	}

	if r.TCP != nil {
		dialer, err := tcp.NewDialer(r.TCP.Opts()...)
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
