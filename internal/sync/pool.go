// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package sync

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// noCopy may be added to structs which must not be copied after the first use.
type noCopy struct{}

// poolHeader holds the GC-managed local and victim pool arrays shared with
// runtime STW poolCleanup regardless of T or build mode.
//
// All fields are accessed through sync/atomic even inside poolCleanup so the race
// detector observes atomic synchronization across runtime.GC() and mutators.
type poolHeader struct {
	local      unsafe.Pointer
	victim     unsafe.Pointer
	allNext    atomic.Pointer[poolHeader]
	oldNext    atomic.Pointer[poolHeader]
	localSize  uintptr
	victimSize uintptr
	inAll      atomic.Bool
	noCopy     noCopy
}

var (
	//nolint:gochecknoglobals // protects per-pool cold initialization against redundant allocation.
	allPoolsMu sync.Mutex

	//nolint:gochecknoglobals // lock-free intrusive stack of pools holding primary caches.
	allPoolsHead atomic.Pointer[poolHeader]

	//nolint:gochecknoglobals // intrusive stack of pools holding previous-cycle victim caches.
	oldPoolsHead atomic.Pointer[poolHeader]
)

// Lock is a no-op used by vet's copylocks checker.
func (*noCopy) Lock() {}

// Unlock is a no-op used by vet's copylocks checker.
func (*noCopy) Unlock() {}

//nolint:gochecknoinits // chains runtime poolCleanup so both Pool[T] and stdlib sync.Pool drain on GC.
func init() {
	runtimeRegisterPoolCleanup(poolCleanup)
}

func registerPool(h *poolHeader) {
	if h.inAll.Swap(true) {
		return
	}
	for {
		old := allPoolsHead.Load()
		h.allNext.Store(old)
		if allPoolsHead.CompareAndSwap(old, h) {
			return
		}
	}
}

// poolCleanup is invoked by the Go runtime during STW at GC start.
// It chains stdlib sync.Pool cleanup and rotates Pool[T] primary -> victim -> cleared.
func poolCleanup() {
	// Runtime stores a single global poolcleanup hook, so stdlib sync.Pool must be chained here.
	stdPoolCleanup()

	for h := oldPoolsHead.Swap(nil); h != nil; {
		next := h.oldNext.Swap(nil)
		atomic.StorePointer(&h.victim, nil)
		atomic.StoreUintptr(&h.victimSize, 0)
		h = next
	}

	for h := allPoolsHead.Swap(nil); h != nil; {
		next := h.allNext.Swap(nil)
		h.inAll.Store(false)

		atomic.StorePointer(&h.victim, atomic.LoadPointer(&h.local))
		atomic.StoreUintptr(&h.victimSize, atomic.LoadUintptr(&h.localSize))
		atomic.StorePointer(&h.local, nil)
		atomic.StoreUintptr(&h.localSize, 0)

		h.oldNext.Store(oldPoolsHead.Load())
		oldPoolsHead.Store(h)
		h = next
	}
}

//go:linkname runtimeRegisterPoolCleanup sync.runtime_registerPoolCleanup

func runtimeRegisterPoolCleanup(cleanup func())

//go:linkname stdPoolCleanup sync.poolCleanup

func stdPoolCleanup()
