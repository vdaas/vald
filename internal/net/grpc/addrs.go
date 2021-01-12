//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"sync"
	"sync/atomic"
)

type AtomicAddrs interface {
	GetAll() ([]string, bool)
	Range(func(addr string) bool)
	Add(addr string)
	Delete(addr string)
	Next() (string, bool)
}

type atomicAddrs struct {
	addrs      atomic.Value
	dupCheck   map[string]bool
	mu         sync.RWMutex
	addrSeeker uint64
	l          uint64
}

func newAddr(addrList map[string]struct{}) AtomicAddrs {
	a := new(atomicAddrs)
	a.dupCheck = make(map[string]bool)
	if addrList == nil {
		a.addrs.Store(make([]string, 0, 10))
	} else {
		addrs := make([]string, 0, len(addrList))
		for addr := range addrList {
			a.dupCheck[addr] = true
			addrs = append(addrs, addr)
		}
		a.addrs.Store(addrs)
		atomic.StoreUint64(&a.l, uint64(len(addrs)))
	}
	return a
}

func (a *atomicAddrs) GetAll() ([]string, bool) {
	aas := a.addrs.Load()
	if aas == nil {
		return nil, false
	}
	addrs, ok := aas.([]string)
	if !ok {
		return nil, false
	}
	return addrs, true
}

func (a *atomicAddrs) Range(f func(addr string) bool) {
	addrs, ok := a.GetAll()
	if ok {
		for _, addr := range addrs {
			if !f(addr) {
				return
			}
		}
	}
	a.mu.RLock()
	defer a.mu.RUnlock()
	for addr := range a.dupCheck {
		if !f(addr) {
			return
		}
	}
}

func (a *atomicAddrs) Add(addr string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.dupCheck[addr] {
		addrs, ok := a.GetAll()
		if !ok {
			addrs = make([]string, 0, 100)
			for addr := range a.dupCheck {
				addrs = append(addrs, addr)
			}
		}
		a.dupCheck[addr] = true
		a.addrs.Store(append(addrs, addr))
		atomic.StoreUint64(&a.l, uint64(len(addrs)))
	}
}

func (a *atomicAddrs) Delete(addr string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.dupCheck[addr] {
		delete(a.dupCheck, addr)
		addrs := make([]string, 0, len(a.dupCheck))
		for addr := range a.dupCheck {
			addrs = append(addrs, addr)
		}
		a.addrs.Store(addrs)
		atomic.StoreUint64(&a.l, uint64(len(addrs)))
	}
}

func (a *atomicAddrs) Next() (string, bool) {
	addrs, ok := a.GetAll()
	if !ok {
		return "", false
	}
	for range addrs {
		addr := addrs[atomic.AddUint64(&a.addrSeeker, 1)%atomic.LoadUint64(&a.l)]
		if len(addr) != 0 {
			return addr, true
		}
	}
	return "", false
}
