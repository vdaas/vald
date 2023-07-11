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
	"time"

	"github.com/vdaas/vald/internal/cache/cacher"
	"github.com/vdaas/vald/internal/cache/gache"
	"github.com/vdaas/vald/internal/errors"
)

type cache[V any] struct {
	cacher         cacher.Type
	expireDur      time.Duration
	expireCheckDur time.Duration
	expiredHook    func(context.Context, string)
}

// New returns the Cache instance or error.
func New[V any](opts ...Option[V]) (cc cacher.Cache[V], err error) {
	c := new(cache[V])
	for _, opt := range append(defaultOptions[V](), opts...) {
		opt(c)
	}
	switch c.cacher {
	case cacher.GACHE:
		return gache.New[V](
			gache.WithExpireDuration[V](c.expireDur),
			gache.WithExpireCheckDuration[V](c.expireCheckDur),
			gache.WithExpiredHook[V](c.expiredHook),
		), nil
	default:
		return nil, errors.ErrInvalidCacherType
	}
}
