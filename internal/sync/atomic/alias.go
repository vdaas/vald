// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package atomic

import "sync/atomic"

type (
	Bool           = atomic.Bool
	Int32          = atomic.Int32
	Int64          = atomic.Int64
	Uint32         = atomic.Uint32
	Uint64         = atomic.Uint64
	Uintptr        = atomic.Uintptr
	Value          = atomic.Value
	Pointer[T any] = atomic.Pointer[T]
)

var (
	AddInt32              = atomic.AddInt32
	AddInt64              = atomic.AddInt64
	AddUint32             = atomic.AddUint32
	AddUint64             = atomic.AddUint64
	AddUintptr            = atomic.AddUintptr
	CompareAndSwapInt32   = atomic.CompareAndSwapInt32
	CompareAndSwapInt64   = atomic.CompareAndSwapInt64
	CompareAndSwapUint32  = atomic.CompareAndSwapUint32
	CompareAndSwapUint64  = atomic.CompareAndSwapUint64
	CompareAndSwapUintptr = atomic.CompareAndSwapUintptr
	CompareAndSwapPointer = atomic.CompareAndSwapPointer
	LoadInt32             = atomic.LoadInt32
	LoadInt64             = atomic.LoadInt64
	LoadUint32            = atomic.LoadUint32
	LoadUint64            = atomic.LoadUint64
	LoadUintptr           = atomic.LoadUintptr
	LoadPointer           = atomic.LoadPointer
	StoreInt32            = atomic.StoreInt32
	StoreInt64            = atomic.StoreInt64
	StoreUint32           = atomic.StoreUint32
	StoreUint64           = atomic.StoreUint64
	StoreUintptr          = atomic.StoreUintptr
	StorePointer          = atomic.StorePointer
	SwapInt32             = atomic.SwapInt32
	SwapInt64             = atomic.SwapInt64
	SwapUint32            = atomic.SwapUint32
	SwapUint64            = atomic.SwapUint64
	SwapUintptr           = atomic.SwapUintptr
	SwapPointer           = atomic.SwapPointer
)
