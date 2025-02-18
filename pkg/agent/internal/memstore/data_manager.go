//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
	var its, dts, kts int64
	_, its, dts, ok = vq.GetVectorWithTimestamp(uuid)
	if !ok {
		oid, kts, ok = kv.Get(uuid)
		if !ok {
			return 0, false
		}
		if kts < its {
			kv.Set(uuid, oid, its)
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

	if kts < its {
		kv.Set(uuid, oid, its)
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
	uuids = make([]string, 0, kv.Len()+uint64(vq.IVQLen())-uint64(vq.DVQLen()))
	var mu sync.Mutex
	ListObjectFunc(ctx, kv, vq, func(uuid string, oid uint32, _ int64) bool {
		mu.Lock()
		uuids = append(uuids, uuid)
		mu.Unlock()
		return true
	})
	return uuids
}

// UpdateTimestamp updates memstore(kvs, vqueue) data's timestamp
func UpdateTimestamp(
	kv kvs.BidiMap,
	vq vqueue.Queue,
	uuid string,
	ts int64,
	force bool,
	getVectorFn func(oid uint32) ([]float32, error),
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
		return nil // succesfully update the vqueue
	case vqok && kvok && dts < ts && (force || (kts < ts && its < ts)):
		// if vqueue data exists and new timestamp never delete and force-update or timestamp is newer than insert queue timestamp
		// update insert-vqueue first
		err = vq.PushInsert(uuid, vec, ts)
		if err != nil {
			return err
		}
		// if updated insert-vqueue and data exists ikvdb and it's timestamp is older than query, update kvs data
		kv.Set(uuid, oid, ts)
		if dts == 0 { // if kvs data exists but not found delete-vqueue data it would be better to add delete vqueue for update
			return vq.PushDelete(uuid, ts-1)
		}
		return nil // succesfully update the vqueue and kvs
	case !vqok && its == 0 && kvok && (force || kts < ts):
		// if insert-vqueue not found and kvs data found just update kvs data
		kv.Set(uuid, oid, ts)
		if dts != 0 && (force || dts < ts) {
			// if delete-vqueue found and ts is newer than delete timestamp and kvs timestamp, should update kvs and remove delete-vqueue
			pdts, ok := vq.PopDelete(uuid)
			if ok && pdts != dts {
				// if time difference detected the data might be changed by another thread so we need to rollback
				return vq.PushDelete(uuid, pdts) // succesfully update the kvs but failed to dequeue delete-vqueue and rollbacked them
			}
			return nil // succesfully update the kvs and delete-vqueue
		}
		return nil // succesfully update the kvs
	case !vqok && its != 0 && kvok && (force || kts < ts):
		// if insert-vqueue found there are 2 case of vqok=false are vec==nil or dts > its so check kvok and update it and remove insert-vqueue
		kv.Set(uuid, oid, ts)
		if vec == nil && its > dts && getVectorFn != nil {
			ovec, err := getVectorFn(oid)
			if err == nil && ovec != nil {
				return vq.PushInsert(uuid, ovec, ts)
			}
		}
		pvec, pits, ok := vq.PopInsert(uuid)
		if pvec != nil && ok && pits != its {
			// if time difference detected the data might be changed by another thread so we need to rollback
			return vq.PushInsert(uuid, pvec, pits)
		}
		return nil // succesfully update the kvs
	}
	return errors.ErrNothingToBeDoneForUpdate(uuid)
}
