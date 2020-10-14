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

// Package service
package service

import (
	"fmt"

	"github.com/vdaas/vald/internal/cache"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(m *meta) error

var (
	defaultOpts = []Option{
		WithMetaCacheEnabled(true),
		WithMetaCacheExpireDuration("30m"),
		WithMetaCacheExpiredCheckDuration("2m"),
	}
)

func WithMetaAddr(addr string) Option {
	return func(m *meta) error {
		m.addr = addr
		return nil
	}
}

func WithMetaHostPort(host string, port int) Option {
	return func(m *meta) error {
		m.addr = fmt.Sprintf("%s:%d", host, port)
		return nil
	}
}

func WithMetaClient(client grpc.Client) Option {
	return func(m *meta) error {
		if client != nil {
			m.client = client
		}
		return nil
	}
}

func WithMetaCacheEnabled(flg bool) Option {
	return func(m *meta) error {
		m.enableCache = flg
		return nil
	}
}

func WithMetaCache(c cache.Cache) Option {
	return func(m *meta) error {
		if c != nil {
			m.cache = c
		}
		return nil
	}
}

func WithMetaCacheExpireDuration(dur string) Option {
	return func(m *meta) error {
		_, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		m.expireDuration = dur
		return nil
	}
}

func WithMetaCacheExpiredCheckDuration(dur string) Option {
	return func(m *meta) error {
		_, err := timeutil.Parse(dur)
		if err != nil {
			return err
		}
		m.expireCheckDuration = dur
		return nil
	}
}
