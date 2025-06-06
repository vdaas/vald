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

package gache

import (
	"context"
	"time"

	gache "github.com/kpango/gache/v2"
)

// Option represents the functional option for cache.
type Option[V any] func(*cache[V])

// defaultOptions returns []Option with gache.New().
func defaultOptions[V any]() []Option[V] {
	return []Option[V]{
		WithGache(gache.New[V]()),
	}
}

// WithGache returns Option after set gache to cache.
func WithGache[V any](g gache.Gache[V]) Option[V] {
	return func(c *cache[V]) {
		c.gache = g
	}
}

// WithExpiredHook returns Option after set expiredHook when f is not nil.
func WithExpiredHook[V any](f func(context.Context, string, V)) Option[V] {
	return func(c *cache[V]) {
		if f != nil {
			c.expiredHook = f
		}
	}
}

// WithExpireDuration returns Option after set expireDur when dur is not 0.
func WithExpireDuration[V any](dur time.Duration) Option[V] {
	return func(c *cache[V]) {
		if dur != 0 {
			c.expireDur = dur
		}
	}
}

// WithExpireCheckDuration returns Option after set expireCheckDur when dur is not 0.
func WithExpireCheckDuration[V any](dur time.Duration) Option[V] {
	return func(c *cache[V]) {
		if dur != 0 {
			c.expireCheckDur = dur
		}
	}
}
