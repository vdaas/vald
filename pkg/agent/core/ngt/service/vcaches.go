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

package service

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type vcaches struct {
	length uint64
	mu     sync.Mutex
	read   atomic.Value // readOnly
	dirty  map[string]*entryVCache
	misses int
}

type readOnlyVCache struct {
	m       map[string]*entryVCache
	amended bool
}

var expungedVCache = unsafe.Pointer(new(vcache))

type entryVCache struct {
	p unsafe.Pointer
}

func newEntryVCache(i vcache) *entryVCache {
	return &entryVCache{p: unsafe.Pointer(&i)}
}

func (m *vcaches) Load(key string) (value vcache, ok bool) {
	read, _ := m.read.Load().(readOnlyVCache)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyVCache)
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

func (e *entryVCache) load() (value vcache, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedVCache {
		return value, false
	}
	return *(*vcache)(p), true
}

func (m *vcaches) Store(key string, value vcache) {
	defer atomic.AddUint64(&m.length, 1)
	read, _ := m.read.Load().(readOnlyVCache)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyVCache)
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
			m.read.Store(readOnlyVCache{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryVCache(value)
	}
	m.mu.Unlock()
}

func (e *entryVCache) tryStore(i *vcache) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedVCache {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryVCache) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedVCache, nil)
}

func (e *entryVCache) storeLocked(i *vcache) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *vcaches) Delete(key string) {
	atomic.AddUint64(&m.length, ^uint64(0))
	read, _ := m.read.Load().(readOnlyVCache)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyVCache)
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

func (e *entryVCache) delete() (hadValue bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedVCache {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
		}
	}
}

func (m *vcaches) Range(f func(key string, value vcache) bool) {
	read, _ := m.read.Load().(readOnlyVCache)
	if read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyVCache)
		if read.amended {
			read = readOnlyVCache{m: m.dirty}
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

func (m *vcaches) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyVCache{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *vcaches) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyVCache)
	m.dirty = make(map[string]*entryVCache, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entryVCache) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedVCache) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedVCache
}

func (m *vcaches) Len() uint64 {
	return atomic.LoadUint64(&m.length)
}
