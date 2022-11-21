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

// Package gache provides implementation of cache using gache
package gache

import (
	"context"
	"time"

	"github.com/kpango/gache"
	"github.com/vdaas/vald/internal/cache/cacher"
)

type cache struct {
	gache          gache.Gache
	expireDur      time.Duration
	expireCheckDur time.Duration
	expiredHook    func(context.Context, string)
}

// New loads a cache model and returns a new cache struct.
func New(opts ...Option) cacher.Cache {
	c := new(cache)
	for _, opt := range append(defaultOptions(), opts...) {
		opt(c)
	}
	c.gache.SetDefaultExpire(c.expireDur)
	if c.expiredHook != nil {
		c.gache = c.gache.
			SetExpiredHook(c.expiredHook).
			EnableExpiredHook()
	}
	return c
}

// Start calls StartExpired func of c.gache.
func (c *cache) Start(ctx context.Context) {
	c.gache.StartExpired(ctx, c.expireCheckDur)
}

// Get calls StartExpired func of c.gache and returns (interface{}, bool) according to key.
func (c *cache) Get(key string) (interface{}, bool) {
	return c.gache.Get(key)
}

// Set calls Set func of c.gache.
func (c *cache) Set(key string, val interface{}) {
	c.gache.Set(key, val)
}

// Delete calls Delete func of c.gache.
func (c *cache) Delete(key string) {
	c.gache.Delete(key)
}

// GetAndDelete returns (interface{}, bool) and delete value according to key when value of key is set.
// When value of key is not set, returns (nil, false).
func (c *cache) GetAndDelete(key string) (interface{}, bool) {
	v, ok := c.gache.Get(key)
	if !ok {
		return nil, false
	}
	c.gache.Delete(key)
	return v, true
}
