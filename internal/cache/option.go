//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Option represents the functional option for cache.
type Option[V any] func(*cache[V])

// var defaultOptions = []Option{
// 	WithType(cacher.GACHE.String()),
// 	WithExpireDuration("30m"),
// 	WithExpireCheckDuration("5m"),
// }

func defaultOptions[V any]() []Option[V] {
	return []Option[V]{
		WithType[V](cacher.GACHE.String()),
		WithExpireDuration[V]("30m"),
		WithExpireCheckDuration[V]("5m"),
	}
}

// WithExpiredHook returns Option after set expiredHook when f is not nil.
func WithExpiredHook[V any](f func(context.Context, string)) Option[V] {
	return func(c *cache[V]) {
		if f != nil {
			c.expiredHook = f
		}
	}
}

// WithType returns Option after set cacher when len(mo string) is not nil.
func WithType[V any](mo string) Option[V] {
	return func(c *cache[V]) {
		if len(mo) == 0 {
			return
		}

		c.cacher = cacher.ToType(mo)
	}
}

// WithExpireDuration returns Option after set expireDur when dur is cprrect param.
func WithExpireDuration[V any](dur string) Option[V] {
	return func(c *cache[V]) {
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

// WithExpireCheckDuration returns Option after set expireCheckDur when dur is cprrect param.
func WithExpireCheckDuration[V any](dur string) Option[V] {
	return func(c *cache[V]) {
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
