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

package memstore

import (
	"context"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/pkg/agent/internal/kvs"
	"github.com/vdaas/vald/pkg/agent/internal/vqueue"
)

func Exists(kv kvs.BidiMap, vq vqueue.Queue, uuid string) (oid uint32, ok bool) {
	var its, dts int64
	_, its, dts, ok = vq.GetVectorWithTimestamp(uuid)
	if !ok {
		oid, _, ok = kv.Get(uuid)
		if !ok {
			return 0, false
		}
		if dts > its {
			log.Debugf(
				"Exists\tuuid: %s's data found in kvsdb and not found in insert vqueue, but delete vqueue data exists. the object will be delete soon\terror: %v",
				uuid,
				errors.ErrObjectIDNotFound(uuid),
			)
			return 0, false
		}
	}
	return oid, ok
}

func GetObject(
	kv kvs.BidiMap, vq vqueue.Queue, uuid string, getVectorFn func(oid uint32) ([]float32, error),
) (vec []float32, timestamp int64, err error) {
	vec, its, dts, exists := vq.GetVectorWithTimestamp(uuid)
	if exists {
		return vec, its, nil
	}

	oid, kts, ok := kv.Get(uuid)
	if !ok {
		log.Debugf("GetObject\tuuid: %s's data not found in kvsdb and insert vqueue", uuid)
		return nil, 0, errors.ErrObjectIDNotFound(uuid)
	}

	if ok && dts > its {
		log.Debugf("GetObject\tuuid: %s's data found in kvsdb and not found in insert vqueue, but delete vqueue data exists. the object will be delete soon", uuid)
		return nil, 0, errors.ErrObjectIDNotFound(uuid)
	}

	if getVectorFn == nil {
		return nil, kts, nil
	}

	vec, err = getVectorFn(oid)
	if err != nil {
		log.Debugf("GetObject\tuuid: %s oid: %d's vector not found in ngt index", uuid, oid)
		return nil, 0, errors.ErrObjectNotFound(err, uuid)
	}

	return vec, kts, nil
}

// ListObjectFunc applies the input function on each index stored in the kvs and vqueue.
// Use this function for performing something on each object with caring about the memory usage.
// If the vector exists in the vqueue, this vector is not indexed so the oid(object ID) is processed as 0.
func ListObjectFunc(
	ctx context.Context,
	kv kvs.BidiMap,
	vq vqueue.Queue,
	f func(uuid string, oid uint32, ts int64) bool,
) {
	dup := make(map[string]bool, max(vq.DVQLen(), 3)/3)
	vq.Range(ctx, func(uuid string, vec []float32, ts int64) (ok bool) {
		oid, kts, ok := kv.Get(uuid)
		if ok { // exists same data ikv
			if ts > kts { // exist ikv but vq is newer
				dup[uuid] = true
				return f(uuid, oid, ts)
			}
			// exist in kv and kvs data is newer thavqueue skip and process it at kvs.Range
			return true
		}
		// not exist in kv
		return f(uuid, 0, ts)
	})
	kv.Range(ctx, func(uuid string, oid uint32, ts int64) (ok bool) {
		if dup[uuid] {
			return true
		}
		// if delete vqueue data exists and timestamp of dvq is newer which means data will be delete soon, then skip process
		dts, ok := vq.DVExists(uuid)
		if ok && dts != 0 {
			return true
		}
		return f(uuid, oid, ts)
	})
}

func UUIDs(ctx context.Context, kv kvs.BidiMap, vq vqueue.Queue) (uuids []string) {
	// The capacity is only a hint. kv.Len() can briefly lag behind pending
	// deletes (kvs is decremented before the delete entry is popped), so clamp
	// instead of letting the uint64 subtraction wrap and make panic.
	n := kv.Len() + uint64(vq.IVQLen())  //nolint:gosec // IVQLen is a queue length converted from a uint64 counter, never negative
	if d := uint64(vq.DVQLen()); d < n { //nolint:gosec // DVQLen is a queue length converted from a uint64 counter, never negative
		n -= d
	} else {
		n = 0
	}
	uuids = make([]string, 0, n)
	var mu sync.Mutex
	ListObjectFunc(ctx, kv, vq, func(uuid string, oid uint32, _ int64) bool {
		mu.Lock()
		uuids = append(uuids, uuid)
		mu.Unlock()
		return true
	})
	return uuids
}

// UpdateTimestamp updates memstore(kvs, vqueue) data's timestamp.
func UpdateTimestamp(
	kv kvs.BidiMap, vq vqueue.Queue, uuid string, ts int64, force bool,
) (err error) {
	if len(uuid) == 0 {
		return errors.ErrUUIDNotFound(0) // invalid uuid, we can't check any object without uuid
	}
	if !force && ts <= 0 {
		return errors.ErrZeroTimestamp
	}
	vec, its, dts, vqok := vq.GetVectorWithTimestamp(uuid) // read insert/delete vqueue data
	oid, kts, kvok := kv.Get(uuid)                         // read kvs data
	if !vqok && !kvok {
		return errors.ErrObjectNotFound(nil, uuid) // no object in memstore then return NotFound
	}
	if !force && (ts <= kts || ts <= its) {
		return errors.ErrNewerTimestampObjectAlreadyExists(uuid, ts) // no old object found in this memstore
	}
	switch {
	case vqok && !kvok && dts != 0 && dts < ts && (force || its < ts):
		// if only found from vqueue and timestamp is newer than delete-vqueue-timestamp(dts)
		// update insert-vqueue first
		err = vq.PushInsert(uuid, vec, ts)
		if err != nil {
			return err
		}
		pdts, ok := vq.PopDelete(uuid) // there is no kvs data and ts is newer than dts which means we don't need to delete processing for uuid this time
		if ok && pdts != dts {
			// if time difference detected the data might be changed by another thread so we need to rollback
			return vq.PushDelete(uuid, pdts)
		}
		return nil // successfully update the vqueue
	case vqok && kvok && its <= kts && (force || kts < ts):
		// The queued insert is not newer than the indexed object, so CreateIndex skips it
		// (and the older delete); the indexed object is the object. Update it and drop
		// the stale entries instead of re-queuing the insert at ts, which would let it
		// replace the newer indexed object.
		if !kv.SetTimestamp(uuid, oid, ts) {
			return UpdateTimestamp(kv, vq, uuid, ts, force)
		}
		return dropStaleQueued(vq, uuid, its, dts)
	case vqok && kvok && dts < ts && (force || (kts < ts && its < ts)):
		// if vqueue data exists and new timestamp never delete and force-update or timestamp is newer than insert queue timestamp
		// update insert-vqueue first
		err = vq.PushInsert(uuid, vec, ts)
		if err != nil {
			return err
		}
		// The kvs timestamp is left as is: it belongs to the indexed object, which the queued
		// delete replaces. Raising it to ts would make CreateIndex skip that delete (it is older
		// than the indexed timestamp) and then the insert at ts (not newer than it).
		if dts == 0 { // if kvs data exists but not found delete-vqueue data it would be better to add delete vqueue for update
			return vq.PushDelete(uuid, ts-1)
		}
		return nil // successfully update the vqueue and kvs
	case !vqok && its == 0 && kvok && (force || kts < ts):
		// if insert-vqueue not found and kvs data found just update kvs data
		if !kv.SetTimestamp(uuid, oid, ts) {
			// CreateIndex deleted or replaced the object after kv.Get; a plain Set would
			// re-create the mapping to a removed oid, so decide again on the current state.
			return UpdateTimestamp(kv, vq, uuid, ts, force)
		}
		if dts != 0 && (force || dts < ts) {
			// if delete-vqueue found and ts is newer than delete timestamp and kvs timestamp, should update kvs and remove delete-vqueue
			pdts, ok := vq.PopDelete(uuid)
			if ok && pdts != dts {
				// if time difference detected the data might be changed by another thread so we need to rollback
				return vq.PushDelete(uuid, pdts) // successfully update the kvs but failed to dequeue delete-vqueue and rollbacked them
			}
			return nil // successfully update the kvs and delete-vqueue
		}
		return nil // successfully update the kvs
	case !vqok && its != 0 && kvok && kts <= dts && dts > its:
		// its != 0 means a queued vector exists (getVector reports its=0 for a nil vector), so
		// vqok=false means the queued delete is not older than it. A newer delete cancels the
		// insert and removes the indexed object at CreateIndex: the object is pending deletion
		// (GetObject and Exists report it as not found). Raising the kvs timestamp past dts
		// would make CreateIndex skip that delete and resurrect the indexed object.
		return errors.ErrObjectNotFound(nil, uuid)
	case !vqok && its != 0 && kvok && kts <= dts && dts < ts && (force || (kts < ts && its < ts)):
		// dts == its: CreateIndex applies the delete and then the queued insert, so the queued
		// vector is the object. Re-queue it at ts and leave the kvs timestamp as is, as above.
		return vq.PushInsert(uuid, vec, ts)
	case !vqok && its != 0 && kvok && kts > dts && (force || kts < ts):
		// The indexed object is newer than both queued entries, so CreateIndex skips them
		// (newest wins) and the indexed object stays. Raise its timestamp and drop the stale
		// insert; re-queuing that insert at ts would let it replace the newer indexed object.
		if !kv.SetTimestamp(uuid, oid, ts) {
			return UpdateTimestamp(kv, vq, uuid, ts, force)
		}
		return dropStaleQueued(vq, uuid, its, dts)
	}
	return errors.ErrNothingToBeDoneForUpdate(uuid)
}

// dropStaleQueued removes the queued insert (its) and delete (dts, 0 for none) of
// uuid that CreateIndex would skip because the indexed object is newer. An entry
// changed by another thread since it was read is put back.
func dropStaleQueued(vq vqueue.Queue, uuid string, its, dts int64) error {
	if pvec, pits, ok := vq.PopInsert(uuid); ok && pvec != nil && pits != its {
		return vq.PushInsert(uuid, pvec, pits)
	}
	if dts != 0 {
		if pdts, ok := vq.PopDelete(uuid); ok && pdts != dts {
			return vq.PushDelete(uuid, pdts)
		}
	}
	return nil
}
