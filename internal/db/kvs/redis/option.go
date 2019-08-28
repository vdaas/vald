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

// Package redis provides implementation of Go API for redis interface
package redis

import "github.com/vdaas/vald/internal/net/tcp"

type Option func(*redisClient) error

var (
	defaultOpts = []Option{}
)

func WithDialer(der tcp.Dialer) Option {
	return func(r *redisClient) error {
		if der != nil {
			r.dialer = der
		}
		return nil
	}
}

func WithAddr(addr string) Option {
	return func(r *redisClient) error {
		if r.addrs == nil {
			r.addrs = []string{addr}
		} else {
			r.addrs = append(r.addrs, addr)
		}
		return nil
	}
}

func WithHosts(addrs []string) Option {
	return func(r *redisClient) error {
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
