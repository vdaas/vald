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

package grpc

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type checkList struct {
	mu     sync.Mutex
	read   atomic.Value
	dirty  map[string]*entryCheckList
	misses int
}

type readOnlyCheckList struct {
	m       map[string]*entryCheckList
	amended bool
}

var expungedCheckList = unsafe.Pointer(new(struct{}))

type entryCheckList struct {
	p unsafe.Pointer
}

func (m *checkList) Exists(key string) bool {
	read, _ := m.read.Load().(readOnlyCheckList)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyCheckList)
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			m.missLocked()
		}
		m.mu.Unlock()
	}
	if !ok {
		return false
	}
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedCheckList {
		return false
	}
	return true
}

func (m *checkList) Check(key string) {
	value := struct{}{}
	read, _ := m.read.Load().(readOnlyCheckList)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyCheckList)
	if e, ok := read.m[key]; ok {
		if e.unexpungeLocked() {
			m.dirty[key] = e
		}
		atomic.StorePointer(&e.p, unsafe.Pointer(&value))
	} else if e, ok := m.dirty[key]; ok {
		atomic.StorePointer(&e.p, unsafe.Pointer(&value))
	} else {
		if !read.amended {
			m.dirtyLocked()
			m.read.Store(readOnlyCheckList{m: read.m, amended: true})
		}
		m.dirty[key] = &entryCheckList{p: unsafe.Pointer(&value)}
	}
	m.mu.Unlock()
}

func (e *entryCheckList) tryStore(i *struct{}) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedCheckList {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryCheckList) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedCheckList, nil)
}

func (m *checkList) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyCheckList{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *checkList) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyCheckList)
	m.dirty = make(map[string]*entryCheckList, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entryCheckList) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedCheckList) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedCheckList
}
