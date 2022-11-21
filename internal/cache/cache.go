//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

type cache struct {
	cacher         cacher.Type
	expireDur      time.Duration
	expireCheckDur time.Duration
	expiredHook    func(context.Context, string)
}

// New returns the Cache instance or error.
func New(opts ...Option) (cc cacher.Cache, err error) {
	c := new(cache)
	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}
	switch c.cacher {
	case cacher.GACHE:
		return gache.New(
			gache.WithExpireDuration(c.expireDur),
			gache.WithExpireCheckDuration(c.expireCheckDur),
			gache.WithExpiredHook(c.expiredHook),
		), nil
	}
	return nil, errors.ErrInvalidCacherType
}
