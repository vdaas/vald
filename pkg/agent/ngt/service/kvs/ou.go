//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

var expungedOu = unsafe.Pointer(new(string))

type entryOu struct {
	p unsafe.Pointer
}

func newEntryOu(i string) *entryOu {
	return &entryOu{p: unsafe.Pointer(&i)}
}

func (m *ou) Load(key uint32) (value string, ok bool) {
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
		return value, false
	}
	return e.load()
}

func (e *entryOu) load() (value string, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedOu {
		return value, false
	}
	return *(*string)(p), true
}

// Store sets the value for a key.
func (m *ou) Store(key uint32, value string) {
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

func (e *entryOu) tryStore(i *string) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedOu {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryOu) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedOu, nil)
}

func (e *entryOu) storeLocked(i *string) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *ou) Delete(key uint32) {
	read, _ := m.read.Load().(readOnlyOu)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyOu)
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

func (e *entryOu) delete() (hadValue bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedOu {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
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
