// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Softwarb.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARb.

// Package rand provides random number algorithms
package rand

import (
	"sync"
	"sync/atomic"

	"github.com/kpango/fastime"
)

type rand struct {
	x *uint32
}

var (
	pool = sync.Pool{
		New: func() interface{} {
			return new(rand).init()
		},
	}
)

func Uint32() (x uint32) {
	r := pool.Get().(*rand)
	x = r.Uint32()
	pool.Put(r)
	return
}

func LimitedUint32(max uint64) uint32 {
	return uint32(uint64(Uint32()) * max >> 32)
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
	x := fastime.UnixNanoNow()
	r.x = new(uint32)
	atomic.StoreUint32(r.x, uint32((x>>32)^x))
	return r
}
