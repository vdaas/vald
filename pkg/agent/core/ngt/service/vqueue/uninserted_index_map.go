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
package vqueue

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type uiim struct {
	mu sync.Mutex

	read atomic.Value

	dirty map[string]*entryUiim

	misses int
}

type readOnlyUiim struct {
	m       map[string]*entryUiim
	amended bool
}

var expungedUiim = unsafe.Pointer(new(index))

type entryUiim struct {
	p unsafe.Pointer
}

func newEntryUiim(i index) *entryUiim {
	return &entryUiim{p: unsafe.Pointer(&i)}
}

func (m *uiim) Load(key string) (value index, ok bool) {
	read, _ := m.read.Load().(readOnlyUiim)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()

		read, _ = m.read.Load().(readOnlyUiim)
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

func (e *entryUiim) load() (value index, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedUiim {
		return value, false
	}
	return *(*index)(p), true
}

func (m *uiim) Store(key string, value index) {
	read, _ := m.read.Load().(readOnlyUiim)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyUiim)
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
			m.read.Store(readOnlyUiim{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryUiim(value)
	}
	m.mu.Unlock()
}

func (e *entryUiim) tryStore(i *index) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedUiim {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryUiim) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedUiim, nil)
}

func (e *entryUiim) storeLocked(i *index) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *uiim) LoadOrStore(key string, value index) (actual index, loaded bool) {
	read, _ := m.read.Load().(readOnlyUiim)
	if e, ok := read.m[key]; ok {
		actual, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, loaded
		}
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyUiim)
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
			m.read.Store(readOnlyUiim{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryUiim(value)
		actual, loaded = value, false
	}
	m.mu.Unlock()

	return actual, loaded
}

func (e *entryUiim) tryLoadOrStore(i index) (actual index, loaded, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == expungedUiim {
		return actual, false, false
	}
	if p != nil {
		return *(*index)(p), true, true
	}

	ic := i
	for {
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
			return i, false, true
		}
		p = atomic.LoadPointer(&e.p)
		if p == expungedUiim {
			return actual, false, false
		}
		if p != nil {
			return *(*index)(p), true, true
		}
	}
}

func (m *uiim) LoadAndDelete(key string) (value index, loaded bool) {
	read, _ := m.read.Load().(readOnlyUiim)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyUiim)
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

func (m *uiim) Delete(key string) {
	m.LoadAndDelete(key)
}

func (e *entryUiim) delete() (value index, ok bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedUiim {
			return value, false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return *(*index)(p), true
		}
	}
}

func (m *uiim) Range(f func(key string, value index) bool) {
	read, _ := m.read.Load().(readOnlyUiim)
	if read.amended {

		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyUiim)
		if read.amended {
			read = readOnlyUiim{m: m.dirty}
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

func (m *uiim) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyUiim{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *uiim) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyUiim)
	m.dirty = make(map[string]*entryUiim, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entryUiim) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedUiim) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedUiim
}
