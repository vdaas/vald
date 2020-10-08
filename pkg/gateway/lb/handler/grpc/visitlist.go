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

// Package grpc provides grpc server logic
package grpc

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type visitlist struct {
	mu     sync.Mutex
	read   atomic.Value
	dirty  map[string]*entryVisitlist
	misses int
}

type readOnlyVisitlist struct {
	m       map[string]*entryVisitlist
	amended bool
}

var expungedVisitlist = unsafe.Pointer(new(struct{}))

type entryVisitlist struct {
	p unsafe.Pointer
}

func (m *visitlist) Visited(key string) (loaded bool) {
	value := struct{}{}
	read, _ := m.read.Load().(readOnlyVisitlist)
	if e, ok := read.m[key]; ok {
		loaded, ok := e.visited(value)
		if ok {
			return loaded
		}
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyVisitlist)
	if e, ok := read.m[key]; ok {
		if atomic.CompareAndSwapPointer(&e.p, expungedVisitlist, nil) {
			m.dirty[key] = e
		}
		loaded, _ = e.visited(value)
	} else if e, ok := m.dirty[key]; ok {
		loaded, _ = e.visited(value)
		m.misses++
		if m.misses < len(m.dirty) {
			return
		}
		m.read.Store(readOnlyVisitlist{m: m.dirty})
		m.dirty = nil
		m.misses = 0
	} else {
		if !read.amended {
			if m.dirty == nil {
				read, _ := m.read.Load().(readOnlyVisitlist)
				m.dirty = make(map[string]*entryVisitlist, len(read.m))
				for k, e := range read.m {
					if !e.tryExpungeLocked() {
						m.dirty[k] = e
					}
				}
			}
			m.read.Store(readOnlyVisitlist{m: read.m, amended: true})
		}
		m.dirty[key] = &entryVisitlist{p: unsafe.Pointer(&value)}
		loaded = false
	}
	m.mu.Unlock()
	return loaded
}

func (e *entryVisitlist) visited(i struct{}) (visited, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == expungedVisitlist {
		return false, false
	}
	if p != nil {
		return true, true
	}
	ic := i
	for {
		if atomic.CompareAndSwapPointer(&e.p, nil, unsafe.Pointer(&ic)) {
			return false, true
		}
		p = atomic.LoadPointer(&e.p)
		if p == expungedVisitlist {
			return false, false
		}
		if p != nil {
			return true, true
		}
	}
}

func (e *entryVisitlist) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedVisitlist) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedVisitlist
}
