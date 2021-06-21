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

// Package vqueue manages the vector cache layer for reducing FFI overhead for fast Agent processing.
package vqueue

import (
	"context"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
)

// Queue
type Queue interface {
	Start(ctx context.Context) (<-chan error, error)
	PushInsert(uuid string, vector []float32, date int64) error
	PushDelete(uuid string, date int64) error
	GetVector(uuid string) ([]float32, bool)
	RangePopInsert(ctx context.Context, f func(uuid string, vector []float32) bool)
	RangePopDelete(ctx context.Context, f func(uuid string) bool)
	IVExists(uuid string) bool
	DVExists(uuid string) bool
	IVQLen() int
	IVCLen() int
	DVQLen() int
	DVCLen() int
}

type vqueue struct {
	ich              chan index // ich is insert channel
	uii              []index    // uii is un inserted index
	imu              sync.Mutex // insert mutex
	uiim             uiim       // uiim is un inserted index map (this value is used for GetVector operation to return queued vector cache data)
	dch              chan key   // dch is delete channel
	udk              []key      // udk is un deleted key
	dmu              sync.Mutex // delete mutex
	udim             udim       // udim is un deleted index map (this value is used for Exists operation to return cache data existence)
	eg               errgroup.Group
	finalizingInsert atomic.Value
	finalizingDelete atomic.Value
	closed           atomic.Value

	// buffer config
	ichSize  int
	dchSize  int
	iBufSize int
	dBufSize int
}

type index struct {
	uuid   string
	vector []float32
	date   int64
}

type key struct {
	uuid   string
	vector []float32
	date   int64
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
	vq.ich = make(chan index, vq.ichSize)
	vq.uii = make([]index, 0, vq.iBufSize)
	vq.dch = make(chan key, vq.dchSize)
	vq.udk = make([]key, 0, vq.dBufSize)
	vq.finalizingInsert.Store(false)
	vq.finalizingDelete.Store(false)
	vq.closed.Store(true)
	return vq, nil
}

func (v *vqueue) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 1)
	v.eg.Go(safety.RecoverFunc(func() (err error) {
		v.closed.Store(false)
		defer v.closed.Store(true)
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				v.finalizingInsert.Store(true)
				close(v.ich)
				for i := range v.ich {
					v.addInsert(i)
				}
				v.finalizingInsert.Store(false)
				err := ctx.Err()
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return nil
				}
				return err
			case i := <-v.ich:
				v.addInsert(i)
			}
		}
	}))
	v.eg.Go(safety.RecoverFunc(func() (err error) {
		v.closed.Store(false)
		defer v.closed.Store(true)
		for {
			select {
			case <-ctx.Done():
				v.finalizingDelete.Store(true)
				close(v.dch)
				for d := range v.dch {
					v.addDelete(d)
				}
				v.finalizingDelete.Store(false)
				err := ctx.Err()
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return nil
				}
				return err
			case d := <-v.dch:
				v.addDelete(d)
			}
		}
	}))
	return ech, nil
}

func (v *vqueue) PushInsert(uuid string, vector []float32, date int64) error {
	// we have to check this instance's channel bypass daemon is finalizing or not, if in finalizing process we should not send new index to channel
	if v.finalizingInsert.Load().(bool) || v.closed.Load().(bool) {
		return errors.ErrVQueueFinalizing
	}
	if date == 0 {
		date = time.Now().UnixNano()
	}
	idx := index{
		uuid:   uuid,
		vector: vector,
		date:   date,
	}
	v.uiim.Store(uuid, idx)
	v.ich <- idx
	return nil
}

func (v *vqueue) PushDelete(uuid string, date int64) error {
	// we have to check this instance's channel bypass daemon is finalizing or not, if in finalizing process we should not send new index to channel
	if v.finalizingDelete.Load().(bool) || v.closed.Load().(bool) {
		return errors.ErrVQueueFinalizing
	}
	if date == 0 {
		date = time.Now().UnixNano()
	}
	v.udim.Store(uuid, date)
	v.dch <- key{
		uuid: uuid,
		date: date,
	}
	return nil
}

func (v *vqueue) RangePopInsert(ctx context.Context, f func(uuid string, vector []float32) bool) {
	// if finalizing, wait for all insert channel queue processed
	for v.finalizingInsert.Load().(bool) {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Millisecond * 100):
		}
	}
	v.flushAndRangeInsert(f)
}

func (v *vqueue) RangePopDelete(ctx context.Context, f func(uuid string) bool) {
	// if finalizing, wait for all insert & delete channel queue processed
	for v.finalizingDelete.Load().(bool) || v.finalizingInsert.Load().(bool) {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Millisecond * 100):
		}
	}
	v.flushAndRangeDelete(f)
}

func (v *vqueue) GetVector(uuid string) ([]float32, bool) {
	vec, ok := v.uiim.Load(uuid)
	if !ok {
		// data not in the insert queue then return not exists(false)
		return nil, false
	}
	di, ok := v.udim.Load(uuid)
	if !ok {
		// data not in the delete queue but exists in insert queue then return exists(true)
		return vec.vector, true
	}
	// data exists both queue, compare data timestamp if insert queue timestamp is newer than delete one, this function returns exists(true)
	if di <= vec.date {
		return vec.vector, true
	}
	return nil, false
}

func (v *vqueue) IVExists(uuid string) bool {
	vec, ok := v.uiim.Load(uuid)
	if !ok {
		// data not in the insert queue then return not exists(false)
		return false
	}
	di, ok := v.udim.Load(uuid)
	if !ok {
		// data not in the delete queue but exists in insert queue then return exists(true)
		return true
	}
	// data exists both queue, compare data timestamp if insert queue timestamp is newer than delete one, this function returns exists(true)
	// However, if insert and delete are sent by the update instruction, the timestamp will be the same
	return di <= vec.date
}

func (v *vqueue) DVExists(uuid string) bool {
	di, ok := v.udim.Load(uuid)
	if !ok {
		return false
	}
	vec, ok := v.uiim.Load(uuid)
	if !ok {
		// data not in the insert queue then return not exists(false)
		return true
	}

	// data exists both queue, compare data timestamp if insert queue timestamp is newer than delete one, this function returns exists(true)
	return di > vec.date
}

func (v *vqueue) addInsert(i index) {
	v.imu.Lock()
	v.uii = append(v.uii, i)
	v.imu.Unlock()
}

func (v *vqueue) addDelete(d key) {
	v.dmu.Lock()
	v.udk = append(v.udk, d)
	v.dmu.Unlock()
}

func (v *vqueue) flushAndRangeInsert(f func(uuid string, vector []float32) bool) {
	v.imu.Lock()
	uii := make([]index, len(v.uii))
	copy(uii, v.uii)
	v.uii = v.uii[:0]
	v.imu.Unlock()
	sort.Slice(uii, func(i, j int) bool {
		// sort by latest unix time order
		return uii[i].date > uii[j].date
	})
	dup := make(map[string]bool, len(uii)/2)
	for i, idx := range uii {
		// if the same uuid is detected in the delete map during insert phase, which means the data is not processed in the delete phase.
		// we need to add it back to insert map to process it in next create index process.
		if _, ok := v.udim.Load(idx.uuid); ok {
			v.imu.Lock()
			v.uii = append(v.uii, idx)
			v.imu.Unlock()
			continue
		}

		// if duplicated data exists current loop's data is old due to the uii's sort order
		if !dup[idx.uuid] {
			dup[idx.uuid] = true
			if !f(idx.uuid, idx.vector) {
				v.imu.Lock()
				v.uii = append(uii[i:], v.uii...)
				v.imu.Unlock()
				return
			}
			v.uiim.Delete(idx.uuid)
		}
	}
}

func (v *vqueue) flushAndRangeDelete(f func(uuid string) bool) {
	v.dmu.Lock()
	udk := make([]key, len(v.udk))
	copy(udk, v.udk)
	v.udk = v.udk[:0]
	v.dmu.Unlock()
	sort.Slice(udk, func(i, j int) bool {
		return udk[i].date > udk[j].date
	})
	dup := make(map[string]bool, len(udk)/2)
	udm := make(map[string]int64, len(udk))
	for i, idx := range udk {
		if !dup[idx.uuid] {
			dup[idx.uuid] = true
			if !f(idx.uuid) {
				v.dmu.Lock()
				v.udk = append(udk[i:], v.udk...)
				v.dmu.Unlock()
				return
			}
			v.udim.Delete(idx.uuid)
			udm[idx.uuid] = idx.date
		}
	}

	dl := make([]int, 0, len(udk)/2)

	// In the CreateIndex operation of the NGT Service, the Delete Queue is processed first, and then the Insert Queue is processed,
	// so the Insert Queue still contains the old Insert Operation older than the Delete Queue,
	// and it is possible that data that was intended to be deleted is registered again.
	// For this reason, the data is deleted from the Insert Queue only when retrieving data from the Delete Queue.
	// we should check insert vqueue if insert vqueue exists and delete operation date is newer than insert operation date then we should remove insert vqueue's data.
	v.imu.Lock()
	for i, idx := range v.uii {
		// check same uuid & operation date
		// if date is equal, it may update operation we shouldn't remove at that time
		date, exists := udm[idx.uuid]
		if exists && date > idx.date {
			dl = append(dl, i)
		}
	}
	v.imu.Unlock()
	sort.Sort(sort.Reverse(sort.IntSlice(dl)))
	for _, i := range dl {
		v.imu.Lock()
		// load removal target uuid
		uuid := v.uii[i].uuid
		// remove unnecessary insert vector queue data
		v.uii = append(v.uii[:i], v.uii[i+1:]...)
		v.imu.Unlock()
		// remove from existing map
		v.uiim.Delete(uuid)
	}
}

func (v *vqueue) IVQLen() (l int) {
	v.imu.Lock()
	l = len(v.uii)
	v.imu.Unlock()
	return l
}

func (v *vqueue) DVQLen() (l int) {
	v.dmu.Lock()
	l = len(v.udk)
	v.dmu.Unlock()
	return l
}

func (v *vqueue) IVCLen() int {
	return len(v.ich)
}

func (v *vqueue) DVCLen() int {
	return len(v.dch)
}
