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

// Package gache provides implementation of cache using gache
package gache

import (
	"context"
	"time"

	"github.com/kpango/gache"
)

// Option represents the functional option for cache.
type Option func(*cache)

// defaultOptions returns []Option with gache.New().
func defaultOptions() []Option {
	return []Option{
		WithGache(gache.New()),
	}
}

// WithGache returns Option after set gache to cache.
func WithGache(g gache.Gache) Option {
	return func(c *cache) {
		c.gache = g
	}
}

// WithExpiredHook returns Option after set expiredHook when f is not nil.
func WithExpiredHook(f func(context.Context, string)) Option {
	return func(c *cache) {
		if f != nil {
			c.expiredHook = f
		}
	}
}

// WithExpireDuration returns Option after set expireDur when dur is not 0.
func WithExpireDuration(dur time.Duration) Option {
	return func(c *cache) {
		if dur != 0 {
			c.expireDur = dur
		}
	}
}

// WithExpireCheckDuration returns Option after set expireCheckDur when dur is not 0.
func WithExpireCheckDuration(dur time.Duration) Option {
	return func(c *cache) {
		if dur != 0 {
			c.expireCheckDur = dur
		}
	}
}
