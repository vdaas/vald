//
// Copyright (c) 2009-2023 The Go Authors. All rights resered.
// Modified by vdaas.org vald team <vald@vdaas.org>
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

//go:build go1.7
// +build go1.7

package semaphore_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/sync/semaphore"
)

// weighted is an interface matching a subset of *Weighted.  It allows
// alternate implementations for testing and benchmarking.
type weighted interface {
	Acquire(context.Context, int64) error
	TryAcquire(int64) bool
	Release(int64)
}

// semChan implements Weighted using a channel for
// comparing against the condition variable-based implementation.
type semChan chan struct{}

func newSemChan(n int64) semChan {
	return semChan(make(chan struct{}, n))
}

func (s semChan) Acquire(_ context.Context, n int64) error {
	for i := int64(0); i < n; i++ {
		s <- struct{}{}
	}
	return nil
}

func (s semChan) TryAcquire(n int64) bool {
	if int64(len(s))+n > int64(cap(s)) {
		return false
	}

	for i := int64(0); i < n; i++ {
		s <- struct{}{}
	}
	return true
}

func (s semChan) Release(n int64) {
	for i := int64(0); i < n; i++ {
		<-s
	}
}

// acquireN calls Acquire(size) on sem N times and then calls Release(size) N times.
func acquireN(b *testing.B, sem weighted, size int64, N int) {
	b.Helper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < N; j++ {
			sem.Acquire(context.Background(), size)
		}
		for j := 0; j < N; j++ {
			sem.Release(size)
		}
	}
}

// tryAcquireN calls TryAcquire(size) on sem N times and then calls Release(size) N times.
func tryAcquireN(b *testing.B, sem weighted, size int64, N int) {
	b.Helper()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < N; j++ {
			if !sem.TryAcquire(size) {
				b.Fatalf("TryAcquire(%v) = false, want true", size)
			}
		}
		for j := 0; j < N; j++ {
			sem.Release(size)
		}
	}
}

func BenchmarkNewSeq(b *testing.B) {
	for _, cap := range []int64{1, 128} {
		b.Run(fmt.Sprintf("Weighted-%d", cap), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = semaphore.NewWeighted(cap)
			}
		})
		b.Run(fmt.Sprintf("semChan-%d", cap), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = newSemChan(cap)
			}
		})
	}
}

func BenchmarkAcquireSeq(b *testing.B) {
	for _, c := range []struct {
		cap, size int64
		N         int
	}{
		{1, 1, 1},
		{2, 1, 1},
		{16, 1, 1},
		{128, 1, 1},
		{2, 2, 1},
		{16, 2, 8},
		{128, 2, 64},
		{2, 1, 2},
		{16, 8, 2},
		{128, 64, 2},
	} {
		for _, w := range []struct {
			name string
			w    weighted
		}{
			{"Weighted", semaphore.NewWeighted(c.cap)},
			{"semChan", newSemChan(c.cap)},
		} {
			b.Run(fmt.Sprintf("%s-acquire-%d-%d-%d", w.name, c.cap, c.size, c.N), func(b *testing.B) {
				acquireN(b, w.w, c.size, c.N)
			})
			b.Run(fmt.Sprintf("%s-tryAcquire-%d-%d-%d", w.name, c.cap, c.size, c.N), func(b *testing.B) {
				tryAcquireN(b, w.w, c.size, c.N)
			})
		}
	}
}
