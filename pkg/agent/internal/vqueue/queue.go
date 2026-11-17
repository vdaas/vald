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

package vqueue

import (
	"cmp"
	"context"
	"reflect"
	"slices"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
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
	// RangePopDelete calls f for each delete queued at or before now, newest first. f reports
	// whether it applied the delete; see the method for what each answer pops.
	RangePopDelete(ctx context.Context, now int64, f func(uuid string, timestamp int64) (applied bool))
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
	push(&v.il, &v.ic, uuid, &idx)
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
	push(&v.dl, &v.dc, uuid, &idx)
	return nil
}

// push stores idx unless a newer entry for uuid is already queued. The CAS loop
// keeps the newest timestamp under concurrent pushes (the map compares entries with
// reflect.DeepEqual, which is cheap for the same pointer), and a key popped in between
// is re-created through LoadOrStore so cnt stays in step with the map.
// cnt is incremented before the entry is published (and the increment is undone when
// the key already existed), so a concurrent pop of the new entry can never decrement
// cnt below zero; cnt may only over-count briefly.
func push(m *sync.Map[string, *index], cnt *uint64, uuid string, idx *index) {
	for {
		// An already queued uuid never changes cnt, so check it first and leave the shared counter
		// alone: bumping and undoing it for every re-push made all pushers contend on one cache line.
		if oidx, ok := m.Load(uuid); ok {
			if !newer(idx.timestamp, oidx.timestamp) || m.CompareAndSwap(uuid, oidx, idx) {
				return
			}
			continue
		}
		atomic.AddUint64(cnt, 1)
		oidx, loaded := m.LoadOrStore(uuid, idx)
		if !loaded {
			return
		}
		atomic.AddUint64(cnt, ^uint64(0))
		if !newer(idx.timestamp, oidx.timestamp) || m.CompareAndSwap(uuid, oidx, idx) {
			return
		}
	}
}

// popIfSame removes uuid only while it still maps to idx, so an entry pushed after
// idx was read is never dropped by mistake. The map compares with reflect.DeepEqual, so
// an entry equal to idx in uuid, vector and timestamp counts as idx too; such entries
// are interchangeable, so which of them is removed does not matter.
func popIfSame(m *sync.Map[string, *index], cnt *uint64, uuid string, idx *index) {
	if m.CompareAndDelete(uuid, idx) {
		atomic.AddUint64(cnt, ^uint64(0))
	}
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
	uii := make([]*index, 0, atomic.LoadUint64(&v.ic))
	defer func() {
		uii = nil
	}()
	v.il.Range(func(uuid string, idx *index) bool {
		if newer(idx.timestamp, now) {
			return true
		}
		// A newer delete cancels the insert. An older (or equal) delete still queued
		// has not been applied yet, e.g. a concurrent Update pushed it after the
		// RangePopDelete(now) snapshot. Applying the insert first would repoint the uuid and let that delete
		// remove the new object on the next run, so keep the insert queued until the
		// delete has been applied.
		if dts, ok := v.loadDVQ(uuid); ok {
			if newer(dts, idx.timestamp) {
				popIfSame(&v.il, &v.ic, uuid, idx)
			}
			return true
		}
		uii = append(uii, idx)
		select {
		case <-ctx.Done():
			return false
		default:
		}
		return true
	})
	slices.SortFunc(uii, func(left, right *index) int {
		return cmp.Compare(right.timestamp, left.timestamp)
	})
	for _, idx := range uii {
		if !f(idx.uuid, idx.vector, idx.timestamp) {
			return
		}

		popIfSame(&v.il, &v.ic, idx.uuid, idx)
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

// RangePopDelete hands each delete in the snapshot to f with its timestamp. When f applied it
// (the uuid is gone from the index), every delete queued for the uuid is satisfied and popped,
// together with the inserts older than the newest one. When f skipped it (the indexed object
// is newer than the delete), only that entry is popped: a delete pushed while f ran was not
// compared with the index yet and stays queued for the next run.
func (v *vqueue) RangePopDelete(
	ctx context.Context, now int64, f func(uuid string, timestamp int64) (applied bool),
) {
	udi := make([]*index, 0, atomic.LoadUint64(&v.dc))
	defer func() {
		udi = nil
	}()
	v.dl.Range(func(_ string, idx *index) bool {
		if newer(idx.timestamp, now) {
			return true
		}
		udi = append(udi, idx)
		select {
		case <-ctx.Done():
			return false
		default:
		}
		return true
	})
	slices.SortFunc(udi, func(left, right *index) int {
		return cmp.Compare(right.timestamp, left.timestamp)
	})
	for _, didx := range udi {
		if !f(didx.uuid, didx.timestamp) {
			popIfSame(&v.dl, &v.dc, didx.uuid, didx)
		} else {
			// f already removed uuid from the index, which also satisfies any newer
			// delete pushed meanwhile, so pop the delete entry by key. Cancel only
			// inserts older than the newest delete popped; a newer insert must stay.
			dts, ok := v.PopDelete(didx.uuid)
			if !ok || dts < didx.timestamp {
				dts = didx.timestamp
			}
			if iidx, ok := v.il.Load(didx.uuid); ok && iidx != nil && newer(dts, iidx.timestamp) {
				popIfSame(&v.il, &v.ic, didx.uuid, iidx)
			}
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
