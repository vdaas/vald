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
	"unsafe"

	"github.com/kpango/fastime"
	"github.com/vdaas/vald/internal/sync"
)

type number interface {
	uint32 | uint64
}

type rng[T number] struct {
	x *T
}

func (r *rng[T]) init() *rng[T] {
	if r.x == nil {
		r.x = new(T)
	}
	for {
		var seed T
		if unsafe.Sizeof(seed) == 4 {
			s := uint32((fastime.UnixNanoNow() >> 32) ^ fastime.UnixNanoNow())
			seed = T(s)
		} else {
			s := uint64(fastime.UnixNanoNow())
			seed = T(s)
		}
		if seed != 0 {
			*r.x = seed
			break
		}
	}
	return r
}

func (r *rng[T]) Value() (x T) {
	if *r.x == 0 {
		r.init()
	}
	x = *r.x
	if unsafe.Sizeof(x) == 4 {
		x ^= x << 13
		x ^= x >> 17
		x ^= x << 5
	} else {
		x ^= x << 13
		x ^= x >> 7
		x ^= x << 17
	}
	*r.x = x
	return x
}

var pool32 = sync.Pool{
	New: func() any {
		return (&rng[uint32]{
			x: new(uint32),
		}).init()
	},
}

var pool64 = sync.Pool{
	New: func() any {
		return (&rng[uint64]{
			x: new(uint64),
		}).init()
	},
}

func Uint32() uint32 {
	r := pool32.Get().(*rng[uint32])
	x := r.Value()
	pool32.Put(r)
	return x
}

func LimitedUint32(max uint64) uint32 {
	return uint32(uint64(Uint32()) * max >> 32)
}

func Float32() float32 {
	return float32(Uint32()) / (1 << 32)
}

func Uint64() uint64 {
	r := pool64.Get().(*rng[uint64])
	x := r.Value()
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
