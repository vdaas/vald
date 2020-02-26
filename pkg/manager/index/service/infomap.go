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

package service

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/vdaas/vald/apis/grpc/payload"
)

type infoMap struct {
	mu     sync.Mutex
	read   atomic.Value
	dirty  map[string]*entryInfoMap
	misses int
}

type readOnlyInfoMap struct {
	m       map[string]*entryInfoMap
	amended bool
}

var expungedInfoMap = unsafe.Pointer(new(*payload.Info_Index))

type entryInfoMap struct {
	p unsafe.Pointer
}

func newEntryInfoMap(i *payload.Info_Index) *entryInfoMap {
	return &entryInfoMap{p: unsafe.Pointer(&i)}
}

func (m *infoMap) Load(key string) (value *payload.Info_Index, ok bool) {
	read, _ := m.read.Load().(readOnlyInfoMap)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyInfoMap)
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

func (e *entryInfoMap) load() (value *payload.Info_Index, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedInfoMap {
		return value, false
	}
	return *(**payload.Info_Index)(p), true
}

func (m *infoMap) Store(key string, value *payload.Info_Index) {
	read, _ := m.read.Load().(readOnlyInfoMap)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyInfoMap)
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
			m.read.Store(readOnlyInfoMap{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryInfoMap(value)
	}
	m.mu.Unlock()
}

func (e *entryInfoMap) tryStore(i **payload.Info_Index) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedInfoMap {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryInfoMap) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedInfoMap, nil)
}

func (e *entryInfoMap) storeLocked(i **payload.Info_Index) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *infoMap) Delete(key string) {
	read, _ := m.read.Load().(readOnlyInfoMap)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyInfoMap)
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

func (e *entryInfoMap) delete() (hadValue bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedInfoMap {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
		}
	}
}

func (m *infoMap) Range(f func(key string, value *payload.Info_Index) bool) {
	read, _ := m.read.Load().(readOnlyInfoMap)
	if read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyInfoMap)
		if read.amended {
			read = readOnlyInfoMap{m: m.dirty}
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

func (m *infoMap) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyInfoMap{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *infoMap) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyInfoMap)
	m.dirty = make(map[string]*entryInfoMap, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entryInfoMap) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedInfoMap) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedInfoMap
}
