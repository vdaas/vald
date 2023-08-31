//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package rand provides random number algorithms
package rand

import (
	"sync/atomic"

	"github.com/kpango/fastime"
	"github.com/vdaas/vald/internal/sync"
)

type rand struct {
	x *uint32
}

var pool = sync.Pool{
	New: func() interface{} {
		return (&rand{
			x: new(uint32),
		}).init()
	},
}

func Uint32() (x uint32) {
	r := pool.Get().(*rand)
	x = r.Uint32()
	pool.Put(r)
	return
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
	return
}

func (r *rand) init() *rand {
	if r.x == nil {
		r.x = new(uint32)
	}
	x := fastime.UnixNanoNow()
	atomic.StoreUint32(r.x, uint32((x>>32)^x))
	return r
}
