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

// Package vqueue manages the vector cache layer for reducing FFI overhead for fast Agent processing.
package vqueue

import (
	"context"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/slices"
	"github.com/vdaas/vald/internal/sync"
)

// Queue represents vector queue cache interface
type Queue interface {
	PushInsert(uuid string, vector []float32, date int64) error
	PushDelete(uuid string, date int64) error
	GetVector(uuid string) (vec []float32, timestamp int64, exists bool)
	Range(ctx context.Context, f func(uuid string, vector []float32, ts int64) bool)
	RangePopInsert(ctx context.Context, now int64, f func(uuid string, vector []float32, date int64) bool)
	RangePopDelete(ctx context.Context, now int64, f func(uuid string) bool)
	IVExists(uuid string) bool
	DVExists(uuid string) bool
	IVQLen() int
	DVQLen() int
}

type vqueue struct {
	il, dl sync.Map[string, *index]
	ic, dc uint64
}

type index struct {
	date   int64
	vector []float32
	uuid   string
}

func New(opts ...Option) (Queue, error) {
	vq := new(vqueue)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(vq); err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))

			e := new(errors.ErrCriticalOption)
			if errors.As(err, &e) {
				log.Error(werr)
				return nil, werr
			}
			log.Warn(werr)
		}
	}
	return vq, nil
}

func (v *vqueue) PushInsert(uuid string, vector []float32, date int64) error {
	if date == 0 {
		date = time.Now().UnixNano()
	}
	didx, ok := v.dl.Load(uuid)
	if ok && didx.date > date {
		return nil
	}
	idx := index{
		uuid:   uuid,
		vector: vector,
		date:   date,
	}
	oidx, loaded := v.il.LoadOrStore(uuid, &idx)
	if loaded {
		if date > oidx.date { // if data already exists and existing index is older than new one
			v.il.Store(uuid, &idx)
		}
	} else {
		_ = atomic.AddUint64(&v.ic, 1)
	}
	return nil
}

func (v *vqueue) PushDelete(uuid string, date int64) error {
	if date == 0 {
		date = time.Now().UnixNano()
	}
	idx := index{
		uuid: uuid,
		date: date,
	}
	oidx, loaded := v.dl.LoadOrStore(uuid, &idx)
	if loaded {
		if date > oidx.date { // if data already exists and existing index is older than new one
			v.dl.Store(uuid, &idx)
		}
	} else {
		_ = atomic.AddUint64(&v.dc, 1)
	}
	return nil
}

// GetVector returns the vector stored in the queue.
// If the same UUID exists in the insert queue and the delete queue, the timestamp is compared.
// And the vector is returned if the timestamp in the insert queue is newer than the delete queue.
func (v *vqueue) GetVector(uuid string) (vec []float32, timestamp int64, exists bool) {
	idx, ok := v.il.Load(uuid)
	if !ok {
		// data not in the insert queue then return not exists(false)
		return nil, 0, false
	}
	didx, ok := v.dl.Load(uuid)
	if !ok {
		// data not in the delete queue but exists in insert queue then return exists(true)
		return idx.vector, idx.date, true
	}
	// data exists both queue, compare data timestamp if insert queue timestamp is newer than delete one, this function returns exists(true)
	if didx.date <= idx.date {
		return idx.vector, idx.date, true
	}
	return nil, 0, false
}

// IVExists returns true if there is the UUID in the insert queue.
// If the same UUID exists in the insert queue and the delete queue, the timestamp is compared.
// And the true is returned if the timestamp in the insert queue is newer than the delete queue.
func (v *vqueue) IVExists(uuid string) bool {
	idx, ok := v.il.Load(uuid)
	if !ok {
		// data not in the insert queue then return not exists(false)
		return false
	}
	didx, ok := v.dl.Load(uuid)
	if !ok {
		// data not in the delete queue but exists in insert queue then return exists(true)
		return true
	}
	// data exists both queue, compare data timestamp if insert queue timestamp is newer than delete one, this function returns exists(true)
	// However, if insert and delete are sent by the update instruction, the timestamp will be the same
	return didx.date <= idx.date
}

// DVExists returns true if there is the UUID in the delete queue.
// If the same UUID exists in the insert queue and the delete queue, the timestamp is compared.
// And the true is returned if the timestamp in the delete queue is newer than the insert queue.
func (v *vqueue) DVExists(uuid string) bool {
	didx, ok := v.dl.Load(uuid)
	if !ok {
		return false
	}
	idx, ok := v.il.Load(uuid)
	if !ok {
		// data not in the insert queue then return not exists(false)
		return true
	}

	// data exists both queue, compare data timestamp if insert queue timestamp is newer than delete one, this function returns exists(true)
	return didx.date > idx.date
}

func (v *vqueue) RangePopInsert(ctx context.Context, now int64, f func(uuid string, vector []float32, date int64) bool) {
	uii := make([]index, 0, atomic.LoadUint64(&v.ic))
	defer func() {
		uii = nil
	}()
	v.il.Range(func(uuid string, idx *index) bool {
		if idx.date > now {
			return true
		}
		didx, ok := v.dl.Load(uuid)
		if ok {
			if idx.date < didx.date {
				v.il.Delete(idx.uuid)
				atomic.AddUint64(&v.ic, ^uint64(0))
			}
			return true
		}
		uii = append(uii, *idx)
		select {
		case <-ctx.Done():
			return false
		default:
		}
		return true
	})
	slices.SortFunc(uii, func(left, right index) bool {
		// sort by latest unix time order
		return left.date > right.date
	})
	for _, idx := range uii {
		if !f(idx.uuid, idx.vector, idx.date) {
			return
		}
		v.il.Delete(idx.uuid)
		atomic.AddUint64(&v.ic, ^uint64(0))
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (v *vqueue) RangePopDelete(ctx context.Context, now int64, f func(uuid string) bool) {
	udi := make([]index, 0, atomic.LoadUint64(&v.dc))
	defer func() {
		udi = nil
	}()
	v.dl.Range(func(_ string, idx *index) bool {
		if idx.date > now {
			return true
		}
		udi = append(udi, *idx)
		select {
		case <-ctx.Done():
			return false
		default:
		}
		return true
	})
	slices.SortFunc(udi, func(left, right index) bool {
		// sort by latest unix time order
		return left.date > right.date
	})
	for _, idx := range udi {
		if !f(idx.uuid) {
			return
		}
		v.dl.Delete(idx.uuid)
		atomic.AddUint64(&v.dc, ^uint64(0))
		iidx, ok := v.il.Load(idx.uuid)
		if ok && idx.date > iidx.date {
			v.il.Delete(idx.uuid)
			atomic.AddUint64(&v.ic, ^uint64(0))
		}
		select {
		case <-ctx.Done():
			return
		default:
		}

	}
}

func (v *vqueue) Range(ctx context.Context, f func(uuid string, vector []float32, ts int64) bool) {
	v.il.Range(func(uuid string, idx *index) bool {
		if idx == nil {
			return true
		}
		didx, ok := v.dl.Load(uuid)
		if !ok || (didx != nil && idx.date > didx.date) {
			return f(uuid, idx.vector, idx.date)
		}
		return true
	})
}

// IVQLen returns the number of uninserted indexes stored in the insert queue.
func (v *vqueue) IVQLen() (l int) {
	return int(atomic.LoadUint64(&v.ic))
}

// DVQLen returns the number of undeleted keys stored in the delete queue.
func (v *vqueue) DVQLen() (l int) {
	return int(atomic.LoadUint64(&v.dc))
}
