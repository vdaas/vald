//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

package kvs

import (
	"context"
	"sync/atomic"

	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/zeebo/xxh3"
)

// BidiMap represents an interface for operating kvs.
type BidiMap interface {
	Get(string) (uint32, int64, bool)
	GetInverse(uint32) (string, int64, bool)
	Set(string, uint32, int64)
	Delete(string) (uint32, bool)
	DeleteInverse(uint32) (string, bool)
	Range(ctx context.Context, f func(string, uint32, int64) bool)
	Len() uint64
	Close() error
}

type valueStructOu struct {
	value     string
	timestamp int64
}

type ValueStructUo struct {
	value     uint32
	timestamp int64
}

type bidi struct {
	concurrency int
	l           uint64
	ou          [slen]*sync.Map[uint32, valueStructOu]
	uo          [slen]*sync.Map[string, ValueStructUo]
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
		b.ou[i] = new(sync.Map[uint32, valueStructOu])
		b.uo[i] = new(sync.Map[string, ValueStructUo])
	}

	if b.eg == nil {
		b.eg, _ = errgroup.New(context.Background())
	}

	if b.concurrency > 0 {
		b.eg.SetLimit(b.concurrency)
	}

	return b
}

// Get returns the value and boolean from the given key.
// If the value does not exist, it returns nil and false.
func (b *bidi) Get(key string) (oid uint32, timestamp int64, exists bool) {
	vs, ok := b.uo[getShardID(key)].Load(key)
	if !ok {
		return 0, 0, false
	}
	return vs.value, vs.timestamp, true
}

// GetInverse returns the key and the boolean from the given val.
// If the key does not exist, it returns nil and false.
func (b *bidi) GetInverse(val uint32) (string, int64, bool) {
	vs, ok := b.ou[val&mask].Load(val)
	if !ok {
		return "", 0, false
	}

	return vs.value, vs.timestamp, true
}

// Set sets the key and val to the bidi.
func (b *bidi) Set(key string, val uint32, ts int64) {
	id := getShardID(key)
	vs, loaded := b.uo[id].LoadOrStore(key, ValueStructUo{value: val, timestamp: ts})
	old := vs.value
	if !loaded { // increase the count only if the key is not exists before
		atomic.AddUint64(&b.l, 1)
	} else {
		b.ou[val&mask].Delete(old)                                    // delete paired map value using old value_key
		b.uo[id].Store(key, ValueStructUo{value: val, timestamp: ts}) // store if loaded for overwrite new value
	}
	b.ou[val&mask].Store(val, valueStructOu{value: key, timestamp: ts}) // store anytime
}

// Delete deletes the key and the value from the bidi by the given key and returns val and true.
// If the value for the key does not exist, it returns nil and false.
func (b *bidi) Delete(key string) (val uint32, ok bool) {
	vs, ok := b.uo[getShardID(key)].LoadAndDelete(key)
	val = vs.value
	if ok {
		b.ou[val&mask].Delete(val)
		atomic.AddUint64(&b.l, ^uint64(0))
	}
	return val, ok
}

// DeleteInverse deletes the key and the value from the bidi by the given val and returns the key and true.
// If the key for the val does not exist, it returns nil and false.
func (b *bidi) DeleteInverse(val uint32) (key string, ok bool) {
	vs, ok := b.ou[val&mask].LoadAndDelete(val)
	key = vs.value
	if ok {
		b.uo[getShardID(key)].LoadAndDelete(key)
		atomic.AddUint64(&b.l, ^uint64(0))
	}
	return key, ok
}

// Range retrieves all set keys and values and calls the callback function f.
func (b *bidi) Range(ctx context.Context, f func(string, uint32, int64) bool) {
	var wg sync.WaitGroup
	for i := range b.uo {
		idx := i
		wg.Add(1)
		b.eg.Go(safety.RecoverFunc(func() (err error) {
			b.uo[idx].Range(func(uuid string, val ValueStructUo) bool {
				select {
				case <-ctx.Done():
					return false
				default:
					return f(uuid, val.value, val.timestamp)
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

func getShardID(key string) (id uint64) {
	if len(key) > 128 {
		return xxh3.HashString(key[:128]) & mask
	}
	return xxh3.HashString(key) & mask
}
