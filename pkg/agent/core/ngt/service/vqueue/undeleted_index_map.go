// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

type udim struct {
	mu     sync.Mutex
	read   atomic.Value
	dirty  map[string]*entryUdim
	misses int
}

type readOnlyUdim struct {
	m       map[string]*entryUdim
	amended bool
}

var expungedUdim = unsafe.Pointer(new(int64))

type entryUdim struct {
	p unsafe.Pointer
}

func newEntryUdim(i int64) *entryUdim {
	return &entryUdim{p: unsafe.Pointer(&i)}
}

func (m *udim) Load(key string) (value int64, ok bool) {
	read, _ := m.read.Load().(readOnlyUdim)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyUdim)
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

func (e *entryUdim) load() (value int64, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedUdim {
		return value, false
	}
	return *(*int64)(p), true
}

func (m *udim) Store(key string, value int64) {
	read, _ := m.read.Load().(readOnlyUdim)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyUdim)
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
			m.read.Store(readOnlyUdim{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryUdim(value)
	}
	m.mu.Unlock()
}

func (e *entryUdim) tryStore(i *int64) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedUdim {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryUdim) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedUdim, nil)
}

func (e *entryUdim) storeLocked(i *int64) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *udim) LoadOrStore(key string, value int64) (actual int64, loaded bool) {
	read, _ := m.read.Load().(readOnlyUdim)
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyUdim)
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
			m.read.Store(readOnlyUdim{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryUdim(value)
		actual, loaded = value, false
	}
	m.mu.Unlock()

	return actual, loaded
}

func (e *entryUdim) tryLoadOrStore(i int64) (actual int64, loaded, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == expungedUdim {
		return actual, false, false
	}
	if p != nil {
		return *(*int64)(p), true, true
	}

	ic := i
	for {
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
			return i, false, true
		}
		p = atomic.LoadPointer(&e.p)
		if p == expungedUdim {
			return actual, false, false
		}
		if p != nil {
			return *(*int64)(p), true, true
		}
	}
}

func (m *udim) LoadAndDelete(key string) (value int64, loaded bool) {
	read, _ := m.read.Load().(readOnlyUdim)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyUdim)
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

func (m *udim) Delete(key string) {
	m.LoadAndDelete(key)
}

func (e *entryUdim) delete() (value int64, ok bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedUdim {
			return value, false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return *(*int64)(p), true
		}
	}
}

func (m *udim) Range(f func(key string, value int64) bool) {
	read, _ := m.read.Load().(readOnlyUdim)
	if read.amended {

		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyUdim)
		if read.amended {
			read = readOnlyUdim{m: m.dirty}
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

func (m *udim) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyUdim{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *udim) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyUdim)
	m.dirty = make(map[string]*entryUdim, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entryUdim) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedUdim) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedUdim
}
