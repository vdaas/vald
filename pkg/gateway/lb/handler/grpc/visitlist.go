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

type visitList struct {
	mu     sync.Mutex
	read   atomic.Value
	dirty  map[string]*entryVisitList
	misses int
}

type readOnlyVisitList struct {
	m       map[string]*entryVisitList
	amended bool
}

var expungedVisitList = unsafe.Pointer(new(bool))

type entryVisitList struct {
	p unsafe.Pointer
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
func (m *visitList) Load(key string) (value bool, ok bool) {
	read, _ := m.read.Load().(readOnlyVisitList)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyVisitList)
		e, ok = read.m[key]
		if !ok && read.amended {
			e, ok = m.dirty[key]
			m.misses++
			if m.misses >= len(m.dirty) {
				m.read.Store(readOnlyVisitList{m: m.dirty})
				m.dirty = nil
				m.misses = 0
			}
		}
		m.mu.Unlock()
	}
	if !ok {
		return value, false
	}
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedVisitList {
		return value, false
	}
	return *(*bool)(p), true
}

// Store sets the value for a key.
func (m *visitList) Store(key string, value bool) {
	read, _ := m.read.Load().(readOnlyVisitList)
	if e, ok := read.m[key]; ok {
		for {
			p := atomic.LoadPointer(&e.p)
			if p == expungedVisitList {
				break
			}
			if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(&value)) {
				return
			}
		}
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyVisitList)
	if e, ok := read.m[key]; ok {
		if atomic.CompareAndSwapPointer(&e.p, expungedVisitList, nil) {
			m.dirty[key] = e
		}
		atomic.StorePointer(&e.p, unsafe.Pointer(&value))
	} else if e, ok := m.dirty[key]; ok {
		atomic.StorePointer(&e.p, unsafe.Pointer(&value))
	} else {
		if !read.amended {
			if m.dirty == nil {

				read, _ := m.read.Load().(readOnlyVisitList)
				m.dirty = make(map[string]*entryVisitList, len(read.m))
				for k, e := range read.m {
					skip := false
					p := atomic.LoadPointer(&e.p)
					for p == nil {
						if atomic.CompareAndSwapPointer(&e.p, nil, expungedVisitList) {
							skip = true
							break
						}
						p = atomic.LoadPointer(&e.p)
					}
					if !skip && p != expungedVisitList {
						m.dirty[k] = e
					}
				}
			}
			m.read.Store(readOnlyVisitList{m: read.m, amended: true})
		}
		m.dirty[key] = &entryVisitList{p: unsafe.Pointer(&value)}
	}
	m.mu.Unlock()
}
