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

package kvs

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/safety"
	"github.com/zeebo/xxh3"
)

// BidiMap represents an interface for operating kvs.
type BidiMap interface {
	Get(string) (uint32, bool)
	GetInverse(uint32) (string, bool)
	Set(string, uint32)
	Delete(string) (uint32, bool)
	DeleteInverse(uint32) (string, bool)
	Range(ctx context.Context, f func(string, uint32) bool)
	Len() uint64
	Close() error
}

type bidi struct {
	concurrency int
	l           uint64
	ou          [slen]*ou
	uo          [slen]*uo
	eg          errgroup.Group
}

const (
	// slen is shards length.
	slen = 512
	// slen = 4096
	// mask is slen-1 Hex value.
	mask = 0x1FF
	// mask = 0xFFF.
)

// New returns the bidi that satisfies the BidiMap interface.
func New(opts ...Option) BidiMap {
	b := &bidi{
		l:           0,
		concurrency: 0,
	}
	for _, opt := range append(defaultOptions, opts...) {
		opt(b)
	}
	for i := range b.ou {
		b.ou[i] = new(ou)
		b.uo[i] = new(uo)
	}

	if b.eg == nil {
		b.eg, _ = errgroup.New(context.Background())
	}

	if b.concurrency > 0 {
		b.eg.Limitation(b.concurrency)
	}

	return b
}

// Get returns the value and boolean from the given key.
// If the value does not exist, it returns nil and false.
func (b *bidi) Get(key string) (uint32, bool) {
	return b.uo[xxh3.HashString(key)&mask].Load(key)
}

// GetInverse returns the key and the boolean from the given val.
// If the key does not exist, it returns nil and false.
func (b *bidi) GetInverse(val uint32) (string, bool) {
	return b.ou[val&mask].Load(val)
}

// Set sets the key and val to the bidi.
func (b *bidi) Set(key string, val uint32) {
	id := xxh3.HashString(key) & mask
	old, loaded := b.uo[id].LoadOrStore(key, val)
	if !loaded { // increase the count only if the key is not exists before
		atomic.AddUint64(&b.l, 1)
	} else {
		b.ou[val&mask].Delete(old) // delete paired map value using old value_key
		b.uo[id].Store(key, val)   // store if loaded for overwrite new value
	}
	b.ou[val&mask].Store(val, key) // store anytime
}

// Delete deletes the key and the value from the bidi by the given key and returns val and true.
// If the value for the key does not exist, it returns nil and false.
func (b *bidi) Delete(key string) (val uint32, ok bool) {
	val, ok = b.uo[xxh3.HashString(key)&mask].LoadAndDelete(key)
	if ok {
		b.ou[val&mask].Delete(val)
		atomic.AddUint64(&b.l, ^uint64(0))
	}
	return val, ok
}

// DeleteInverse deletes the key and the value from the bidi by the given val and returns the key and true.
// If the key for the val does not exist, it returns nil and false.
func (b *bidi) DeleteInverse(val uint32) (key string, ok bool) {
	key, ok = b.ou[val&mask].LoadAndDelete(val)
	if ok {
		b.uo[xxh3.HashString(key)&mask].Delete(key)
		atomic.AddUint64(&b.l, ^uint64(0))
	}
	return key, ok
}

// Range retrieves all set keys and values and calls the callback function f.
func (b *bidi) Range(ctx context.Context, f func(string, uint32) bool) {
	var wg sync.WaitGroup
	for i := range b.uo {
		idx := i
		wg.Add(1)
		b.eg.Go(safety.RecoverFunc(func() (err error) {
			b.uo[idx].Range(func(uuid string, oid uint32) bool {
				f(uuid, oid)
				select {
				case <-ctx.Done():
					return false
				default:
					return true
				}
			})
			wg.Done()
			return nil
		}))
	}
	wg.Wait()
}

// Len returns the length of the cache that is set in the bidi.
func (b *bidi) Len() uint64 {
	if b == nil {
		return 0
	}
	return atomic.LoadUint64(&b.l)
}

func (b *bidi) Close() error {
	if b == nil {
		return nil
	}
	return b.eg.Wait()
}
