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

// Package singleflight represents zero time caching.
package singleflight

import (
	"context"
	"sync"
	"sync/atomic"
)

type call[V any] struct {
	wg   sync.WaitGroup
	val  V
	err  error
	dups uint64
}

// Group represents interface for zero time cache.
type Group[V any] interface {
	Do(ctx context.Context, key string, fn func() (V, error)) (v V, shared bool, err error)
}

type group[V any] struct {
	m sync.Map
}

// New returns Group implementation.
func New[V any]() Group[V] {
	return new(group[V])
}

// Do execute the given function and return the result.
// It makes sure only one execution of the function for each given key.
// If duplicate comes, the duplicated call with the same key will wait for the first caller return.
// It returns the result and the error of the given function, and whether the result is shared from the first caller.
func (g *group[V]) Do(_ context.Context, key string, fn func() (V, error)) (v V, shared bool, err error) {
	actual, loaded := g.m.LoadOrStore(key, new(call[V]))
	c := actual.(*call[V])
	if loaded {
		atomic.AddUint64(&c.dups, 1)
		c.wg.Wait()
		v, err = c.val, c.err
		return v, true, err
	}
	c.wg.Add(1)
	c.val, c.err = fn()
	c.wg.Done()

	g.m.LoadAndDelete(key)

	return c.val, atomic.LoadUint64(&c.dups) > 0, c.err
}
