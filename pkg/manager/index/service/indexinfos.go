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

package service

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
)

type indexInfos struct {
	mu     sync.Mutex
	read   atomic.Value
	dirty  map[string]*entryIndexInfos
	misses int
}

type readOnlyIndexInfos struct {
	m       map[string]*entryIndexInfos
	amended bool
}

// skipcq: GSC-G103
var expungedIndexInfos = unsafe.Pointer(new(*payload.Info_Index_Count))

type entryIndexInfos struct {
	p unsafe.Pointer
}

func newEntryIndexInfos(i *payload.Info_Index_Count) *entryIndexInfos {
	// skipcq: GSC-G103
	return &entryIndexInfos{p: unsafe.Pointer(&i)}
}

func (m *indexInfos) Load(key string) (value *payload.Info_Index_Count, ok bool) {
	read, _ := m.read.Load().(readOnlyIndexInfos)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyIndexInfos)
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

func (e *entryIndexInfos) load() (value *payload.Info_Index_Count, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedIndexInfos {
		return value, false
	}
	return *(**payload.Info_Index_Count)(p), true
}

func (m *indexInfos) Store(key string, value *payload.Info_Index_Count) {
	read, _ := m.read.Load().(readOnlyIndexInfos)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyIndexInfos)
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
			m.read.Store(readOnlyIndexInfos{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryIndexInfos(value)
	}
	m.mu.Unlock()
}

func (e *entryIndexInfos) tryStore(i **payload.Info_Index_Count) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedIndexInfos {
			return false
		}
		// skipcq: GSC-G103
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryIndexInfos) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedIndexInfos, nil)
}

func (e *entryIndexInfos) storeLocked(i **payload.Info_Index_Count) {
	// skipcq: GSC-G103
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *indexInfos) Delete(key string) {
	read, _ := m.read.Load().(readOnlyIndexInfos)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyIndexInfos)
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

func (e *entryIndexInfos) delete() (hadValue bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedIndexInfos {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
		}
	}
}

func (m *indexInfos) Range(f func(key string, value *payload.Info_Index_Count) bool) {
	read, _ := m.read.Load().(readOnlyIndexInfos)
	if read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyIndexInfos)
		if read.amended {
			read = readOnlyIndexInfos{m: m.dirty}
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

func (m *indexInfos) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyIndexInfos{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *indexInfos) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyIndexInfos)
	m.dirty = make(map[string]*entryIndexInfos, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entryIndexInfos) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedIndexInfos) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedIndexInfos
}
