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
package kvs

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type uo struct {
	mu     sync.Mutex
	read   atomic.Value
	dirty  map[string]*entryUo
	misses int
}

type readOnlyUo struct {
	m       map[string]*entryUo
	amended bool
}

// skipcq: GSC-G103
var expungedUo = unsafe.Pointer(new(ValueStructUo))

type entryUo struct {
	p unsafe.Pointer
}

type ValueStructUo struct {
	value     uint32
	timestamp int64
}

func newEntryUo(i ValueStructUo) *entryUo {
	// skipcq: GSC-G103
	return &entryUo{p: unsafe.Pointer(&i)}
}

func (m *uo) Load(key string) (value uint32, timestamp int64, ok bool) {
	read, _ := m.read.Load().(readOnlyUo)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyUo)
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return value, timestamp, false
	}
	return e.load()
}

func (e *entryUo) load() (value uint32, timestamp int64, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedUo {
		return value, timestamp, false
	}
	return (*ValueStructUo)(p).value, (*ValueStructUo)(p).timestamp, true
}

func (m *uo) Store(key string, value ValueStructUo) {
	read, _ := m.read.Load().(readOnlyUo)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}
	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyUo)
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
			m.read.Store(readOnlyUo{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryUo(value)
	}
	m.mu.Unlock()
}

func (e *entryUo) tryStore(i *ValueStructUo) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedUo {
			return false
		}
		// skipcq: GSC-G103
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryUo) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedUo, nil)
}

func (e *entryUo) storeLocked(i *ValueStructUo) {
	// skipcq: GSC-G103
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *uo) LoadOrStore(key string, value ValueStructUo) (actual uint32, at int64, loaded bool) {
	read, _ := m.read.Load().(readOnlyUo)
	if e, ok := read.m[key]; ok {
		actual, at, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, at, loaded
		}
	}
	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyUo)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		actual, at, loaded, _ = e.tryLoadOrStore(value)
	} else if e, ok := m.dirty[key]; ok {
		actual, at, loaded, _ = e.tryLoadOrStore(value)
		m.missLocked()
	} else {
		if !read.amended {
			m.dirtyLocked()
			m.read.Store(readOnlyUo{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryUo(value)
		actual, at, loaded = value.value, value.timestamp, false
	}
	m.mu.Unlock()
	return actual, at, loaded
}

func (e *entryUo) tryLoadOrStore(i ValueStructUo) (actual uint32, at int64, loaded, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == expungedUo {
		return actual, at, false, false
	}
	if p != nil {
		return (*ValueStructUo)(p).value, (*ValueStructUo)(p).timestamp, true, true
	}
	ic := i
	for {
		// skipcq: GSC-G103
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
			return i.value, i.timestamp, false, true
		}
		p = atomic.LoadPointer(&e.p)
		if p == expungedUo {
			return actual, at, false, false
		}
		if p != nil {
			return (*ValueStructUo)(p).value, (*ValueStructUo)(p).timestamp, true, true
		}
	}
}

func (m *uo) LoadAndDelete(key string) (value uint32, timestamp int64, loaded bool) {
	read, _ := m.read.Load().(readOnlyUo)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyUo)
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
	return value, timestamp, false
}

func (m *uo) Delete(key string) {
	m.LoadAndDelete(key)
}

func (e *entryUo) delete() (value uint32, timestamp int64, ok bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedUo {
			return value, timestamp, false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return (*ValueStructUo)(p).value, (*ValueStructUo)(p).timestamp, true
		}
	}
}

func (m *uo) Range(f func(key string, value ValueStructUo) bool) {
	read, _ := m.read.Load().(readOnlyUo)
	if read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyUo)
		if read.amended {
			read = readOnlyUo{m: m.dirty}
			m.read.Store(read)
			m.dirty = nil
			m.misses = 0
		}
		m.mu.Unlock()
	}

	for k, e := range read.m {
		v, t, ok := e.load()
		if !ok {
			continue
		}
		if !f(k, ValueStructUo{value: v, timestamp: t}) {
			break
		}
	}
}

func (m *uo) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyUo{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *uo) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyUo)
	m.dirty = make(map[string]*entryUo, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entryUo) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedUo) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedUo
}
