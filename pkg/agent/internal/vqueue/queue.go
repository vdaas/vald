//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	"cmp"
	"context"
	"reflect"
	"slices"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/sync"
)

// Queue represents vector queue cache interface.
type Queue interface {
	PushInsert(uuid string, vector []float32, timestamp int64) error
	PushDelete(uuid string, timestamp int64) error
	PopInsert(uuid string) (vector []float32, timestamp int64, ok bool)
	PopDelete(uuid string) (timestamp int64, ok bool)
	GetVector(uuid string) (vec []float32, timestamp int64, exists bool)
	Range(ctx context.Context, f func(uuid string, vector []float32, ts int64) bool)
	GetVectorWithTimestamp(uuid string) (vec []float32, its, dts int64, exists bool)
	RangePopInsert(ctx context.Context, now int64, f func(uuid string, vector []float32, timestamp int64) bool)
	RangePopDelete(ctx context.Context, now int64, f func(uuid string) bool)
	IVExists(uuid string) (timestamp int64, ok bool)
	DVExists(uuid string) (timestamp int64, ok bool)
	IVQLen() int
	DVQLen() int
}

type vqueue struct {
	il, dl sync.Map[string, *index]
	ic, dc uint64
}

type index struct {
	uuid      string
	vector    []float32
	timestamp int64
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

func (v *vqueue) PushInsert(uuid string, vector []float32, timestamp int64) error {
	if len(uuid) == 0 || vector == nil {
		return nil
	}
	if timestamp == 0 {
		timestamp = time.Now().UnixNano()
	}
	dts, ok := v.loadDVQ(uuid)
	if ok && newer(dts, timestamp) {
		return nil
	}
	idx := index{
		uuid:      uuid,
		vector:    vector,
		timestamp: timestamp,
	}
	oidx, loaded := v.il.LoadOrStore(uuid, &idx)
	if loaded {
		if newer(timestamp, oidx.timestamp) { // if data already exists and existing index is older than new one
			v.il.Store(uuid, &idx)
		}
	} else {
		_ = atomic.AddUint64(&v.ic, 1)
	}
	return nil
}

func (v *vqueue) PushDelete(uuid string, timestamp int64) error {
	if len(uuid) == 0 {
		return nil
	}
	if timestamp == 0 {
		timestamp = time.Now().UnixNano()
	}
	idx := index{
		uuid:      uuid,
		timestamp: timestamp,
	}
	oidx, loaded := v.dl.LoadOrStore(uuid, &idx)
	if loaded {
		if newer(timestamp, oidx.timestamp) { // if data already exists and existing index is older than new one
			v.dl.Store(uuid, &idx)
		}
	} else {
		_ = atomic.AddUint64(&v.dc, 1)
	}
	return nil
}

func (v *vqueue) PopInsert(uuid string) (vector []float32, timestamp int64, ok bool) {
	var idx *index
	idx, ok = v.il.LoadAndDelete(uuid)
	if !ok || idx == nil || idx.timestamp == 0 {
		return nil, 0, false
	}
	_ = atomic.AddUint64(&v.ic, ^uint64(0))
	return idx.vector, idx.timestamp, ok
}

func (v *vqueue) PopDelete(uuid string) (timestamp int64, ok bool) {
	var idx *index
	idx, ok = v.dl.LoadAndDelete(uuid)
	if !ok || idx == nil || idx.timestamp == 0 {
		return 0, false
	}
	_ = atomic.AddUint64(&v.dc, ^uint64(0))
	return idx.timestamp, ok
}

// GetVector returns the vector stored in the queue.
func (v *vqueue) GetVector(uuid string) (vec []float32, timestamp int64, exists bool) {
	vec, timestamp, _, exists = v.getVector(uuid, false)
	return vec, timestamp, exists
}

// GetVectorWithTimestamp returns the vector and timestamps stored in the queue.
func (v *vqueue) GetVectorWithTimestamp(uuid string) (vec []float32, its, dts int64, exists bool) {
	return v.getVector(uuid, true)
}

// getVector returns the vector and timestamps stored in the queue.
// If the same UUID exists in the insert queue and the delete queue, the timestamp is compared.
// And the vector is returned if the timestamp in the insert queue is newer than the delete queue.
func (v *vqueue) getVector(
	uuid string, enableDeleteTimestamp bool,
) (vec []float32, its, dts int64, ok bool) {
	vec, its, ok = v.loadIVQ(uuid)
	if !ok || vec == nil {
		if !enableDeleteTimestamp {
			// data not in the insert queue then return not exists(false)
			return nil, 0, 0, false
		}
		dts, ok = v.loadDVQ(uuid)
		if !ok || dts == 0 {
			// data not in the delete queue and insert queue then return not exists(false)
			return nil, 0, 0, false
		}
		// data not in theinsert queue and exists in delete queue then return not exists(false) with delete index timestamp
		return nil, 0, dts, false
	}
	dts, ok = v.loadDVQ(uuid)
	if !ok || dts == 0 {
		// data not in the delete queue but exists in insert queue then return exists(true)
		return vec, its, 0, vec != nil // usually vec is non-nil which means true
	}
	// data exists both queue, compare data timestamp if insert queue timestamp is newer than delete one last value will true
	// However, if insert and delete are sent by the update instruction, the timestamp will be the same
	return vec, its, dts, vec != nil && newer(its, dts) // ususaly vec is non-nil
}

// IVExists returns timestamp of iv and true if there is the UUID in the insert queue.
// If the same UUID exists in the insert queue and the delete queue, the timestamp is compared.
// And the true is returned if the timestamp in the insert queue is newer than the delete queue.
func (v *vqueue) IVExists(uuid string) (its int64, ok bool) {
	_, its, _, ok = v.getVector(uuid, false)
	if !ok || its == 0 {
		return 0, false
	}
	return its, true
}

// DVExists returns timestamp of dv and true if there is the UUID in the delete queue.
// If the same UUID exists in the insert queue and the delete queue, the timestamp is compared.
// And the true is returned if the timestamp in the delete queue is newer than the insert queue.
func (v *vqueue) DVExists(uuid string) (dts int64, ok bool) {
	_, _, dts, ok = v.getVector(uuid, true)
	if ok || dts == 0 {
		return 0, false
	}
	return dts, true
}

func (v *vqueue) RangePopInsert(
	ctx context.Context, now int64, f func(uuid string, vector []float32, timestamp int64) bool,
) {
	uii := make([]index, 0, atomic.LoadUint64(&v.ic))
	defer func() {
		uii = nil
	}()
	v.il.Range(func(uuid string, idx *index) bool {
		if newer(idx.timestamp, now) {
			return true
		}
		dts, ok := v.loadDVQ(uuid)
		if ok && newer(dts, idx.timestamp) {
			_, _, _ = v.PopInsert(uuid)
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
	slices.SortFunc(uii, func(left, right index) int {
		return cmp.Compare(right.timestamp, left.timestamp)
	})
	for _, idx := range uii {
		if !f(idx.uuid, idx.vector, idx.timestamp) {
			return
		}

		_, _, _ = v.PopInsert(idx.uuid)
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
		if newer(idx.timestamp, now) {
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
	slices.SortFunc(udi, func(left, right index) int {
		return cmp.Compare(right.timestamp, left.timestamp)
	})
	for _, didx := range udi {
		if !f(didx.uuid) {
			return
		}
		_, _ = v.PopDelete(didx.uuid)
		_, its, ok := v.loadIVQ(didx.uuid)
		if ok && newer(didx.timestamp, its) {
			_, _, _ = v.PopInsert(didx.uuid)
		}
		select {
		case <-ctx.Done():
			return
		default:
		}

	}
}

// Range calls f sequentially for each key and value present in the vqueue.
func (v *vqueue) Range(_ context.Context, f func(uuid string, vector []float32, ts int64) bool) {
	v.il.Range(func(uuid string, idx *index) bool {
		if idx == nil {
			return true
		}
		dts, ok := v.loadDVQ(uuid)
		if !ok || newer(idx.timestamp, dts) {
			return f(uuid, idx.vector, idx.timestamp)
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

func (v *vqueue) loadIVQ(uuid string) (vec []float32, ts int64, ok bool) {
	var idx *index
	idx, ok = v.il.Load(uuid)
	if !ok || idx == nil {
		return nil, 0, false
	}
	return idx.vector, idx.timestamp, true
}

func (v *vqueue) loadDVQ(uuid string) (ts int64, ok bool) {
	var idx *index
	idx, ok = v.dl.Load(uuid)
	if !ok || idx == nil {
		return 0, false
	}
	return idx.timestamp, true
}

func newer(ts1, ts2 int64) bool {
	return ts1 > ts2
}
