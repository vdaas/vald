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
package kvs

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type ou struct {
	mu     sync.Mutex
	read   atomic.Value
	dirty  map[uint32]*entryOu
	misses int
}

type readOnlyOu struct {
	m       map[uint32]*entryOu
	amended bool
}

// skipcq: GSC-G103
var expungedOu = unsafe.Pointer(new(ValueStructOu))

type entryOu struct {
	p unsafe.Pointer
}

type ValueStructOu struct {
	value     string
	timestamp int64
}

func newEntryOu(i ValueStructOu) *entryOu {
	// skipcq: GSC-G103
	return &entryOu{p: unsafe.Pointer(&i)}
}

func (m *ou) Load(key uint32) (value string, timestamp int64, ok bool) {
	read, _ := m.read.Load().(readOnlyOu)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyOu)
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

func (e *entryOu) load() (value string, timestamp int64, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedOu {
		return value, timestamp, false
	}
	return (*ValueStructOu)(p).value, (*ValueStructOu)(p).timestamp, true
}

func (m *ou) Store(key uint32, value ValueStructOu) {
	read, _ := m.read.Load().(readOnlyOu)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}
	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyOu)
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
			m.read.Store(readOnlyOu{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryOu(value)
	}
	m.mu.Unlock()
}

func (e *entryOu) tryStore(i *ValueStructOu) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedOu {
			return false
		}
		// skipcq: GSC-G103
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryOu) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedOu, nil)
}

func (e *entryOu) storeLocked(i *ValueStructOu) {
	// skipcq: GSC-G103
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *ou) LoadOrStore(key uint32, value ValueStructOu) (actual string, at int64, loaded bool) {
	read, _ := m.read.Load().(readOnlyOu)
	if e, ok := read.m[key]; ok {
		actual, at, loaded, ok := e.tryLoadOrStore(value)
		if ok {
			return actual, at, loaded
		}
	}
	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyOu)
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
			m.read.Store(readOnlyOu{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryOu(value)
		actual, at, loaded = value.value, value.timestamp, false
	}
	m.mu.Unlock()
	return actual, at, loaded
}

func (e *entryOu) tryLoadOrStore(i ValueStructOu) (actual string, at int64, loaded, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == expungedOu {
		return actual, at, false, false
	}
	if p != nil {
		return (*ValueStructOu)(p).value, (*ValueStructOu)(p).timestamp, true, true
	}
	ic := i
	for {
		// skipcq: GSC-G103
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
			return i.value, i.timestamp, false, true
		}
		p = atomic.LoadPointer(&e.p)
		if p == expungedOu {
			return actual, at, false, false
		}
		if p != nil {
			return (*ValueStructOu)(p).value, (*ValueStructOu)(p).timestamp, true, true
		}
	}
}

func (m *ou) LoadAndDelete(key uint32) (value string, timestamp int64, loaded bool) {
	read, _ := m.read.Load().(readOnlyOu)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyOu)
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

func (m *ou) Delete(key uint32) {
	m.LoadAndDelete(key)
}

func (e *entryOu) delete() (value string, timestamp int64, ok bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedOu {
			return value, timestamp, false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return (*ValueStructOu)(p).value, (*ValueStructOu)(p).timestamp, true
		}
	}
}

func (m *ou) Range(f func(key uint32, value ValueStructOu) bool) {
	read, _ := m.read.Load().(readOnlyOu)
	if read.amended {

		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyOu)
		if read.amended {
			read = readOnlyOu{m: m.dirty}
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
		if !f(k, ValueStructOu{value: v, timestamp: t}) {
			break
		}
	}
}

func (m *ou) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyOu{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *ou) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyOu)
	m.dirty = make(map[uint32]*entryOu, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entryOu) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedOu) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedOu
}
