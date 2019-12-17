//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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

	"google.golang.org/grpc"
)

type gRPCConns struct {
	mu     sync.Mutex
	read   atomic.Value // readOnly
	dirty  map[string]*entryGRPCConns
	misses int
}

type readOnlyGRPCConns struct {
	m       map[string]*entryGRPCConns
	amended bool
}

var expungedGRPCConns = unsafe.Pointer(new(*grpc.ClientConn))

type entryGRPCConns struct {
	p unsafe.Pointer // *interface{}
}

func newEntryGRPCConns(i *grpc.ClientConn) *entryGRPCConns {
	return &entryGRPCConns{p: unsafe.Pointer(&i)}
}

func (m *gRPCConns) Load(key string) (value *grpc.ClientConn, ok bool) {
	read, _ := m.read.Load().(readOnlyGRPCConns)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyGRPCConns)
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

func (e *entryGRPCConns) load() (value *grpc.ClientConn, ok bool) {
	p := atomic.LoadPointer(&e.p)
	if p == nil || p == expungedGRPCConns {
		return value, false
	}
	return *(**grpc.ClientConn)(p), true
}

func (m *gRPCConns) Store(key string, value *grpc.ClientConn) {
	read, _ := m.read.Load().(readOnlyGRPCConns)
	if e, ok := read.m[key]; ok && e.tryStore(&value) {
		return
	}

	m.mu.Lock()
	read, _ = m.read.Load().(readOnlyGRPCConns)
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
			m.read.Store(readOnlyGRPCConns{m: read.m, amended: true})
		}
		m.dirty[key] = newEntryGRPCConns(value)
	}
	m.mu.Unlock()
}

func (e *entryGRPCConns) tryStore(i **grpc.ClientConn) bool {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == expungedGRPCConns {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
			return true
		}
	}
}

func (e *entryGRPCConns) unexpungeLocked() (wasExpunged bool) {
	return atomic.CompareAndSwapPointer(&e.p, expungedGRPCConns, nil)
}

func (e *entryGRPCConns) storeLocked(i **grpc.ClientConn) {
	atomic.StorePointer(&e.p, unsafe.Pointer(i))
}

func (m *gRPCConns) Delete(key string) {
	read, _ := m.read.Load().(readOnlyGRPCConns)
	e, ok := read.m[key]
	if !ok && read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyGRPCConns)
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

func (e *entryGRPCConns) delete() (hadValue bool) {
	for {
		p := atomic.LoadPointer(&e.p)
		if p == nil || p == expungedGRPCConns {
			return false
		}
		if atomic.CompareAndSwapPointer(&e.p, p, nil) {
			return true
		}
	}
}

func (m *gRPCConns) Range(f func(key string, value *grpc.ClientConn) bool) {
	read, _ := m.read.Load().(readOnlyGRPCConns)
	if read.amended {
		m.mu.Lock()
		read, _ = m.read.Load().(readOnlyGRPCConns)
		if read.amended {
			read = readOnlyGRPCConns{m: m.dirty}
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

func (m *gRPCConns) missLocked() {
	m.misses++
	if m.misses < len(m.dirty) {
		return
	}
	m.read.Store(readOnlyGRPCConns{m: m.dirty})
	m.dirty = nil
	m.misses = 0
}

func (m *gRPCConns) dirtyLocked() {
	if m.dirty != nil {
		return
	}

	read, _ := m.read.Load().(readOnlyGRPCConns)
	m.dirty = make(map[string]*entryGRPCConns, len(read.m))
	for k, e := range read.m {
		if !e.tryExpungeLocked() {
			m.dirty[k] = e
		}
	}
}

func (e *entryGRPCConns) tryExpungeLocked() (isExpunged bool) {
	p := atomic.LoadPointer(&e.p)
	for p == nil {
		if atomic.CompareAndSwapPointer(&e.p, nil, expungedGRPCConns) {
			return true
		}
		p = atomic.LoadPointer(&e.p)
	}
	return p == expungedGRPCConns
}
