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

// Package cache provides implementation of cache
package cache

import (
	"context"

	"github.com/vdaas/vald/internal/cache/cacher"
	"github.com/vdaas/vald/internal/timeutil"
)

type Option func(*cache)

var (
	defaultOpts = []Option{
		WithType(cacher.GACHE.String()),
		WithExpireDuration("30m"),
		WithExpireCheckDuration("5m"),
	}
)

func WithExpiredHook(f func(context.Context, string)) Option {
	return func(c *cache) {
		if f != nil {
			c.expiredHook = f
		}
	}
}

func WithType(mo string) Option {
	return func(c *cache) {
		if len(mo) == 0 {
			return
		}

		c.cacher = cacher.ToType(mo)
	}
}

func WithExpireDuration(dur string) Option {
	return func(c *cache) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return
		}
		c.expireDur = d
	}
}

func WithExpireCheckDuration(dur string) Option {
	return func(c *cache) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			return
		}
		c.expireCheckDur = d
	}
}
