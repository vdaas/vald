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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/vdaas/vald/internal/net/grpc/pool"
)

type grpcConns struct {
	mu     sync.Mutex
	read   atomic.Value
	dirty  map[string]*entryGrpcConns
	misses int
}

type readOnlyGrpcConns struct {
	m       map[string]*entryGrpcConns
	amended bool
}

var expungedGrpcConns = unsafe.Pointer(new(pool.Conn))

type entryGrpcConns struct {
	p unsafe.Pointer
}

func newEntryGrpcConns(i pool.Conn) *entryGrpcConns {
	return &entryGrpcConns{p: unsafe.Pointer(&i)}
}

func (m *grpcConns) Load(key string) (value pool.Conn, ok bool) {
	read, _ := m.read.Load().(readOnlyGrpcConns)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyGrpcConns)
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return value, false
	}
	return e.load()
}

func (e *entryGrpcConns) load() (value pool.Conn, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedGrpcConns {
		return value, false
	}
	return *(*pool.Conn)(p), true
}

func (m *grpcConns) Store(key string, value pool.Conn) {
	read, _ := m.read.Load().(readOnlyGrpcConns)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyGrpcConns)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		e.storeLocked(&value)
	} else if e, ok := m.dirty[key]; ok {
		e.storeLocked(&value)
	} else {
		if !read.amended {
			m.dirtyLocked()
			m.read.Store(readOnlyGrpcConns{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryGrpcConns(value)
	}
	m.mu.Unlock()
}

func (e *entryGrpcConns) tryStore(i *pool.Conn) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedGrpcConns {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryGrpcConns) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedGrpcConns, nil)
}

func (e *entryGrpcConns) storeLocked(i *pool.Conn) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *grpcConns) Delete(key string) {
	read, _ := m.read.Load().(readOnlyGrpcConns)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyGrpcConns)
		e, ok = read.m[key]
		if !ok && read.amended {
			delete(m.dirty, key)
		}
		m.mu.Unlock()
	}
	if ok {
		e.delete()
	}
}

func (e *entryGrpcConns) delete() (hadValue bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedGrpcConns {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
		}
	}
}

func (m *grpcConns) Range(f func(key string, value pool.Conn) bool) {
	read, _ := m.read.Load().(readOnlyGrpcConns)
	if read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyGrpcConns)
		if read.amended {
			read = readOnlyGrpcConns{m: m.dirty}
			m.read.Store(read)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	for k, e := range read.m {
		v, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

func (m *grpcConns) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyGrpcConns{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *grpcConns) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyGrpcConns)
	m.dirty = make(map[string]*entryGrpcConns, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entryGrpcConns) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedGrpcConns) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedGrpcConns
}
