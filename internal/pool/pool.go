//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package pool provides pool functionality for pooling buffer or others
package pool

import (
	"bytes"
	"context"
	"sync"
	"sync/atomic"
)

type Buffer interface {
	Get(ctx context.Context) interface{}
	Put(ctx context.Context, data interface{})
	PutWithResize(ctx context.Context, data interface{}, size uint64)
	Size(ctx context.Context) (size uint64)
	Limit(ctx context.Context) (size uint64)
	Len(ctx context.Context) (size uint64)
}

type Extender interface {
	Extend(ctx context.Context, size uint64) (edata interface{})
}

type Flusher interface {
	Flush(ctx context.Context) (fdata interface{})
}

type pool struct {
	size   uint64
	length uint64
	limit  uint64
	new    func(size uint64) interface{}
	pool   sync.Pool
}

func New(ctx context.Context, opts ...Option) Buffer {
	p := new(pool)
	for _, opt := range append(defaultOptions, opts...) {
		opt(p)
	}
	alloc := func() interface{} {
		size := p.Size(ctx)
		if p.new == nil {
			return bytes.NewBuffer(make([]byte, 0, size))
		}
		return p.new(size)
	}
	p.pool = sync.Pool{
		New: alloc,
	}
	for i := p.Limit(ctx) - p.Len(ctx); i > 0; i-- {
		p.Put(ctx, alloc())
	}
	return p
}

func (p *pool) Get(ctx context.Context) interface{} {
	p.decrementLength(ctx)
	return p.pool.Get()
}

func (p *pool) Put(ctx context.Context, data interface{}) {
	if p.Limit(ctx) < p.incrementLength(ctx) && p.Limit(ctx) > 0 {
		p.decrementLength(ctx)
		return
	}
	if flusher, ok := data.(Flusher); ok {
		data = flusher.Flush(ctx)
	}
	p.pool.Put(data)
}

func (p *pool) PutWithResize(ctx context.Context, data interface{}, size uint64) {
	if size <= 1 || p.extender == nil {
		p.Put(ctx, data)
		return
	}
	if extender, ok := data.(Extender); ok {
		data = extender.Extend(ctx, p.extendSize(ctx, size))
	}
	p.Put(ctx, data)
}

func (p *pool) Size(ctx context.Context) (size uint64) {
	return atomic.LoadUint64(&p.size)
}

func (p *pool) extendSize(ctx context.Context, size uint64) (ret uint64) {
	ret = p.Size(ctx)
	if ret < size {
		atomic.StoreUint64(&p.size, size)
		return size
	}
	return ret
}

func (p *pool) incrementLength(ctx context.Context) (size uint64) {
	return atomic.AddUint64(&p.length, 1)
}

func (p *pool) decrementLength(ctx context.Context) (size uint64) {
	return atomic.AddUint64(&p.length, ^uint64(0))
}

func (p *pool) Len(ctx context.Context) (size uint64) {
	return atomic.LoadUint64(&p.length)
}

func (p *pool) Limit(ctx context.Context) (size uint64) {
	return atomic.LoadUint64(&p.limit)
}
