//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"

	xxhash "github.com/cespare/xxhash/v2"
)

type BidiMap interface {
	Get(string) (uint32, bool)
	GetInverse(uint32) (string, bool)
	Set(string, uint32)
	Delete(string) (uint32, bool)
	DeleteInverse(uint32) (string, bool)
	Range(ctx context.Context, f func(string, uint32) bool)
	Len() uint64
}

type bidi struct {
	ou [slen]*ou
	uo [slen]*uo
	l  uint64
}

const (
	// slen is shards length.
	slen = 512
	// slen = 4096
	// mask is slen-1 Hex value.
	mask = 0x1FF
	// mask = 0xFFF.
)

func New() BidiMap {
	b := &bidi{
		l: 0,
	}
	for i := range b.ou {
		b.ou[i] = new(ou)
	}
	for i := range b.uo {
		b.uo[i] = new(uo)
	}
	return b
}

func (b *bidi) Get(key string) (uint32, bool) {
	return b.uo[xxhash.Sum64(stringToBytes(key))&mask].Load(key)
}

func (b *bidi) GetInverse(val uint32) (string, bool) {
	return b.ou[val&mask].Load(val)
}

func (b *bidi) Set(key string, val uint32) {
	b.uo[xxhash.Sum64(stringToBytes(key))&mask].Store(key, val)
	b.ou[val&mask].Store(val, key)
	atomic.AddUint64(&b.l, 1)
}

func (b *bidi) Delete(key string) (val uint32, ok bool) {
	idx := xxhash.Sum64(stringToBytes(key)) & mask
	val, ok = b.uo[idx].Load(key)
	if !ok {
		return 0, false
	}
	b.uo[idx].Delete(key)
	b.ou[val&mask].Delete(val)
	atomic.AddUint64(&b.l, ^uint64(0))
	return val, true
}

func (b *bidi) DeleteInverse(val uint32) (key string, ok bool) {
	idx := val & mask
	key, ok = b.ou[idx].Load(val)
	if !ok {
		return "", false
	}
	b.uo[xxhash.Sum64(stringToBytes(key))&mask].Delete(key)
	b.ou[val&mask].Delete(val)
	atomic.AddUint64(&b.l, ^uint64(0))
	return key, true
}

func (b *bidi) Range(ctx context.Context, f func(string, uint32) bool) {
	wg := new(sync.WaitGroup)
	for i := range b.uo {
		wg.Add(1)
		go func(c context.Context, idx int) {
			b.uo[idx].Range(func(uuid string, oid uint32) bool {
				select {
				case <-c.Done():
					return false
				default:
					f(uuid, oid)
					return true
				}
			})
			wg.Done()
		}(ctx, i)
	}
	wg.Wait()
}

func (b *bidi) Len() uint64 {
	return atomic.LoadUint64(&b.l)
}

func stringToBytes(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}))
}
