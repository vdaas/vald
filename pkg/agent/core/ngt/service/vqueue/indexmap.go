// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package vqueue

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type indexMap struct {
	mu     sync.Mutex
	read   atomic.Value
	dirty  map[string]*entryIndexMap
	misses int
}

type readOnlyIndexMap struct {
	m       map[string]*entryIndexMap
	amended bool
}

// skipcq: GSC-G103
var expungedIndexMap = unsafe.Pointer(new(index))

type entryIndexMap struct {
	// skipcq: GSC-G103
	p unsafe.Pointer
}

func newEntryIndexMap(i index) *entryIndexMap {
	// skipcq: GSC-G103
	return &entryIndexMap{p: unsafe.Pointer(&i)}
}

func (m *indexMap) Load(key string) (value index, ok bool) {
	read, _ := m.read.Load().(readOnlyIndexMap)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()

		read, _ = m.read.Load().(readOnlyIndexMap)
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

func (e *entryIndexMap) load() (value index, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedIndexMap {
		return value, false
	}
	return *(*index)(p), true
}

func (m *indexMap) Store(key string, value index) {
	read, _ := m.read.Load().(readOnlyIndexMap)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyIndexMap)
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
			m.read.Store(readOnlyIndexMap{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryIndexMap(value)
	}
	m.mu.Unlock()
}

func (e *entryIndexMap) tryStore(i *index) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedIndexMap {
			return false
		}
		// skipcq: GSC-G103
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryIndexMap) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedIndexMap, nil)
}

func (e *entryIndexMap) storeLocked(i *index) {
	// skipcq: GSC-G103
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *indexMap) LoadOrStore(key string, value index) (actual index, loaded bool) {
	read, _ := m.read.Load().(readOnlyIndexMap)
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyIndexMap)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		actual, loaded, _ = e.tryLoadOrStore(value)
	} else if e, ok := m.dirty[key]; ok {
		actual, loaded, _ = e.tryLoadOrStore(value)
		m.missLocked()
	} else {
		if !read.amended {

			m.dirtyLocked()
			m.read.Store(readOnlyIndexMap{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryIndexMap(value)
		actual, loaded = value, false
	}
	m.mu.Unlock()

	return actual, loaded
}

func (e *entryIndexMap) tryLoadOrStore(i index) (actual index, loaded, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == expungedIndexMap {
		return actual, false, false
	}
	if p != nil {
		return *(*index)(p), true, true
	}

	ic := i
	for {

		// skipcq: GSC-G103
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
			return i, false, true
		}
		p = atomic.LoadPointer(&e.p)
		if p == expungedIndexMap {
			return actual, false, false
		}
		if p != nil {
			return *(*index)(p), true, true
		}
	}
}

func (m *indexMap) LoadAndDelete(key string) (value index, loaded bool) {
	read, _ := m.read.Load().(readOnlyIndexMap)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyIndexMap)
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			delete(m.dirty, key)

			m.missLocked()
		}
		m.mu.Unlock()
	}
	if ok {
		return e.delete()
	}
	return value, false
}

func (m *indexMap) Delete(key string) {
	m.LoadAndDelete(key)
}

func (e *entryIndexMap) delete() (value index, ok bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedIndexMap {
			return value, false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return *(*index)(p), true
		}
	}
}

func (m *indexMap) Range(f func(key string, value index) bool) {
	read, _ := m.read.Load().(readOnlyIndexMap)
	if read.amended {

		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyIndexMap)
		if read.amended {
			read = readOnlyIndexMap{m: m.dirty}
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

func (m *indexMap) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyIndexMap{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *indexMap) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyIndexMap)
	m.dirty = make(map[string]*entryIndexMap, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entryIndexMap) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedIndexMap) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedIndexMap
}
