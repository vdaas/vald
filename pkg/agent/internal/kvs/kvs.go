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

package kvs

import (
	"context"

	"github.com/vdaas/vald/internal/hash"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"golang.org/x/sys/cpu"
)

// BidiMap represents an interface for operating kvs.
type BidiMap interface {
	Get(string) (uint32, int64, bool)
	GetInverse(uint32) (string, int64, bool)
	Set(string, uint32, int64)
	SetTimestamp(string, uint32, int64) bool
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

type ouShard struct {
	m  map[uint32]valueStructOu
	mu sync.RWMutex
	_  cpu.CacheLinePad
}

type uoShard struct {
	m  map[string]ValueStructUo
	mu sync.RWMutex
	_  cpu.CacheLinePad
}

type bidi struct {
	ou          [slen]*ouShard
	uo          [slen]*uoShard
	concurrency int
	l           uint64
}

const (
	// slen is shards length.
	slen = 512
	// slen = 4096
	// mask is slen-1 Hex value.
	mask = 0x1FF
	// mask = 0xFFF.

	maxHashKeyLength = slen / 2
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
	ouBuf := make([]ouShard, slen)
	uoBuf := make([]uoShard, slen)
	for i := range slen {
		b.ou[i] = &ouBuf[i]
		b.uo[i] = &uoBuf[i]
	}

	return b
}

// Get returns the value and boolean from the given key.
// If the value does not exist, it returns nil and false.
func (b *bidi) Get(key string) (oid uint32, timestamp int64, exists bool) {
	s := b.uo[getShardID(key)]
	s.mu.RLock()
	vs, ok := s.m[key]
	s.mu.RUnlock()
	if !ok {
		return 0, 0, false
	}
	return vs.value, vs.timestamp, true
}

// GetInverse returns the key and the boolean from the given val.
// If the key does not exist, it returns nil and false.
func (b *bidi) GetInverse(val uint32) (string, int64, bool) {
	s := b.ou[val&mask]
	s.mu.RLock()
	vs, ok := s.m[val]
	s.mu.RUnlock()
	if !ok {
		return "", 0, false
	}

	return vs.value, vs.timestamp, true
}

// Set sets the key and val to the bidi.
func (b *bidi) Set(key string, val uint32, ts int64) {
	uo := b.uo[getShardID(key)]
	uo.mu.Lock()
	if uo.m == nil {
		uo.m = make(map[string]ValueStructUo)
	}
	vs, loaded := uo.m[key]
	uo.m[key] = ValueStructUo{value: val, timestamp: ts}
	uo.mu.Unlock()

	if loaded {
		if vs.value != val {
			oldOu := b.ou[vs.value&mask]
			oldOu.mu.Lock()
			if cur, ok := oldOu.m[vs.value]; ok && cur.value == key {
				delete(oldOu.m, vs.value)
			}
			oldOu.mu.Unlock()
		}
	} else { // increase the count only if the key did not exist before
		atomic.AddUint64(&b.l, 1)
	}

	ou := b.ou[val&mask]
	ou.mu.Lock()
	if ou.m == nil {
		ou.m = make(map[uint32]valueStructOu)
	}
	ou.m[val] = valueStructOu{value: key, timestamp: ts}
	ou.mu.Unlock()
}

// SetTimestamp sets the timestamp of key to ts only while key is mapped to val, and reports
// whether it did. Unlike Set it never re-creates a pair removed by a concurrent Delete nor
// overwrites a pair repointed by a concurrent Set.
func (b *bidi) SetTimestamp(key string, val uint32, ts int64) bool {
	uo := b.uo[getShardID(key)]
	uo.mu.Lock()
	vs, ok := uo.m[key]
	if !ok || vs.value != val {
		uo.mu.Unlock()
		return false
	}
	uo.m[key] = ValueStructUo{value: val, timestamp: ts}
	uo.mu.Unlock()

	ou := b.ou[val&mask]
	ou.mu.Lock()
	ovs, ok := ou.m[val]
	if ok && ovs.value == key {
		ou.m[val] = valueStructOu{value: key, timestamp: ts}
	}
	ou.mu.Unlock()
	return true
}

// Delete deletes the key and the value from the bidi by the given key and returns val and true.
// If the value for the key does not exist, it returns nil and false.
func (b *bidi) Delete(key string) (val uint32, ok bool) {
	uo := b.uo[getShardID(key)]
	uo.mu.Lock()
	vs, ok := uo.m[key]
	if ok {
		delete(uo.m, key)
	}
	uo.mu.Unlock()

	val = vs.value
	if ok {
		ou := b.ou[val&mask]
		ou.mu.Lock()
		if cur, ok2 := ou.m[val]; ok2 && cur.value == key {
			delete(ou.m, val)
		}
		ou.mu.Unlock()
		atomic.AddUint64(&b.l, ^uint64(0))
	}
	return val, ok
}

// DeleteInverse deletes the key and the value from the bidi by the given val and returns the key and true.
// If the key for the val does not exist, it returns nil and false.
func (b *bidi) DeleteInverse(val uint32) (key string, ok bool) {
	ou := b.ou[val&mask]
	ou.mu.Lock()
	vs, ok := ou.m[val]
	if ok {
		delete(ou.m, val)
	}
	ou.mu.Unlock()

	key = vs.value
	if ok {
		uo := b.uo[getShardID(key)]
		uo.mu.Lock()
		if cur, ok2 := uo.m[key]; ok2 && cur.value == val {
			delete(uo.m, key)
		}
		uo.mu.Unlock()
		atomic.AddUint64(&b.l, ^uint64(0))
	}
	return key, ok
}

// Range retrieves all set keys and values and calls the callback function f.
func (b *bidi) Range(ctx context.Context, f func(string, uint32, int64) bool) {
	type entry struct {
		key string
		val ValueStructUo
	}
	eg, ctx := errgroup.New(ctx)
	if b.concurrency > 0 {
		eg.SetLimit(b.concurrency)
	}
	for i := range b.uo {
		idx := i
		eg.Go(safety.RecoverFunc(func() (err error) {
			shard := b.uo[idx]
			shard.mu.RLock()
			if len(shard.m) == 0 {
				shard.mu.RUnlock()
				return nil
			}
			snapshot := make([]entry, 0, len(shard.m))
			for k, v := range shard.m {
				snapshot = append(snapshot, entry{key: k, val: v})
			}
			shard.mu.RUnlock()

			for _, e := range snapshot {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					if !f(e.key, e.val.value, e.val.timestamp) {
						return nil
					}
				}
			}
			return nil
		}))
	}
	err := eg.Wait()
	if err != nil {
		log.Error(err)
	}
}

// Len returns the length of the cache that is set in the bidi.
func (b *bidi) Len() uint64 {
	if b == nil {
		return 0
	}
	return atomic.LoadUint64(&b.l)
}

func (b *bidi) Close() error {
	return nil
}

func getShardID(key string) (id uint64) {
	if len(key) > maxHashKeyLength {
		return hash.String(key[:maxHashKeyLength]) & mask
	}
	return hash.String(key) & mask
}
