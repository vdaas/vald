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
	"sync"
	"sync/atomic"
)

type Buffer interface {
	Get() interface{}
	Put(data interface{})
	PutWithResize(data interface{}, size uint64)
	Size() (size uint64)
	Limit() (size uint64)
	Len() (size uint64)
}

type Extender interface {
	Extend(size uint64) (data interface{})
}

type Flusher interface {
	Flush() (data interface{})
}

type pool struct {
	size   uint64
	length uint64
	limit  uint64
	new    func(size uint64) interface{}
	pool   sync.Pool
}

func New(opts ...Option) Buffer {
	p := new(pool)
	for _, opt := range append(defaultOptions, opts...) {
		opt(p)
	}
	alloc := func() interface{} {
		size := p.Size()
		if p.new == nil {
			return bytes.NewBuffer(make([]byte, 0, size))
		}
		return p.new(size)
	}
	p.pool = sync.Pool{
		New: alloc,
	}
	for i := p.Limit() - p.Len(); p.Limit() != 0 && i > 0; i-- {
		p.Put(alloc())
	}
	return p
}

func (p *pool) Get() interface{} {
	p.decrementLength()
	return p.pool.Get()
}

func (p *pool) Put(data interface{}) {
	if p.Limit() < p.incrementLength() && p.Limit() > 0 {
		p.decrementLength()
		return
	}
	if flusher, ok := data.(Flusher); ok {
		data = flusher.Flush()
	}
	p.pool.Put(data)
}

func (p *pool) PutWithResize(data interface{}, size uint64) {
	if size <= 1 {
		p.Put(data)
		return
	}
	if extender, ok := data.(Extender); ok {
		data = extender.Extend(p.extendSize(size))
	}
	p.Put(data)
}

func (p *pool) Size() (size uint64) {
	return atomic.LoadUint64(&p.size)
}

func (p *pool) extendSize(size uint64) (ret uint64) {
	ret = p.Size()
	if ret < size {
		atomic.StoreUint64(&p.size, size)
		return size
	}
	return ret
}

func (p *pool) incrementLength() (size uint64) {
	return atomic.AddUint64(&p.length, 1)
}

func (p *pool) decrementLength() (size uint64) {
	return atomic.AddUint64(&p.length, ^uint64(0))
}

func (p *pool) Len() (size uint64) {
	return atomic.LoadUint64(&p.length)
}

func (p *pool) Limit() (size uint64) {
	return atomic.LoadUint64(&p.limit)
}
