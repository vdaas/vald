//
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

package rand

import (
	"math/bits"
	"sync/atomic"

	"github.com/kpango/fastime"
	"github.com/vdaas/vald/internal/sync"
)

type rand struct {
	x *uint32
}

var pool = sync.Pool{
	New: func() any {
		return (&rand{
			x: new(uint32),
		}).init()
	},
}

func Uint32() (x uint32) {
	r := pool.Get().(*rand)
	x = r.Uint32()
	pool.Put(r)
	return x
}

func LimitedUint32(max uint64) uint32 {
	return uint32(uint64(Uint32()) * max >> 32)
}

func Float32() float32 {
	return float32(Uint32()) / (1 << 32)
}

func (r *rand) Uint32() (x uint32) {
	if atomic.LoadUint32(r.x) == 0 {
		r.init()
	}
	x = atomic.LoadUint32(r.x)
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	atomic.StoreUint32(r.x, x)
	return x
}

func (r *rand) init() *rand {
	if r.x == nil {
		r.x = new(uint32)
	}
	for {
		seed := uint32((fastime.UnixNanoNow() >> 32) ^ fastime.UnixNanoNow())
		if seed != 0 {
			atomic.StoreUint32(r.x, seed)
			break
		}
	}
	return r
}

type rand64 struct {
	x *uint64
}

var pool64 = sync.Pool{
	New: func() any {
		return (&rand64{
			x: new(uint64),
		}).init()
	},
}

func Uint64() (x uint64) {
	r := pool64.Get().(*rand64)
	x = r.Uint64()
	pool64.Put(r)
	return x
}

func LimitedUint64(max uint64) uint64 {
	hi, _ := bits.Mul64(Uint64(), max)
	return hi
}

func Float64() float64 {
	return float64(Uint64()>>11) / (1 << 53)
}

func (r *rand64) Uint64() (x uint64) {
	if *r.x == 0 {
		r.init()
	}
	x = *r.x
	x ^= x << 13
	x ^= x >> 7
	x ^= x << 17
	*r.x = x
	return x
}

func (r *rand64) init() *rand64 {
	if r.x == nil {
		r.x = new(uint64)
	}
	for {
		seed := uint64(fastime.UnixNanoNow())
		if seed != 0 {
			atomic.StoreUint64(r.x, seed)
			break
		}
	}
	return r
}
