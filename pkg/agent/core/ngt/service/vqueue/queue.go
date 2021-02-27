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
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/safety"
)

type Queue interface {
	Start(ctx context.Context) (<-chan error, error)
	PushInsert(uuid string, vector []float32, date int64) error
	PushDelete(uuid string, date int64) error
	PopInsert() (uuid string, vector []float32)
	PopDelete() (uuid string)
	GetVector(uuid string) ([]float32, bool)
	RangePopInsert(ctx context.Context, f func(uuid string, vector []float32) bool)
	RangePopDelete(ctx context.Context, f func(uuid string) bool)
	IVQLen() int
	DVQLen() int
}

type vqueue struct {
	ich        chan index
	uii        []index // un inserted index
	imu        sync.Mutex
	uiil       map[string][]float32
	dch        chan key
	udk        []key // un deleted key
	dmu        sync.Mutex
	eg         errgroup.Group
	finalizing atomic.Value
	closed     atomic.Value
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

func New(eg errgroup.Group) Queue {
	vq := &vqueue{
		ich:  make(chan index, 1000),
		uii:  make([]index, 0, 10000),
		uiil: make(map[string][]float32, 100),
		dch:  make(chan key, 1000),
		udk:  make([]key, 0, 10000),
		eg:   eg,
	}
	vq.finalizing.Store(false)
	vq.closed.Store(true)
	return vq
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
				v.finalizing.Store(true)
				close(v.dch)
				close(v.ich)
				for d := range v.dch {
					v.addDelete(d)
				}
				for i := range v.ich {
					v.addInsert(i)
				}
				v.finalizing.Store(false)
				return ctx.Err()
			case i := <-v.ich:
				v.addInsert(i)
			case d := <-v.dch:
				v.addDelete(d)
			}
		}
	}))
	return ech, nil
}

func (v *vqueue) PushInsert(uuid string, vector []float32, date int64) error {
	// we have to check this instance's channel bypass daemon is finalizing or not, if in finalizing process we should not send new index to channel
	if v.finalizing.Load().(bool) || v.closed.Load().(bool) {
		return errors.ErrVQueueFinalizing
	}
	if date == 0 {
		date = time.Now().UnixNano()
	}

	v.ich <- index{
		uuid:   uuid,
		vector: vector,
		date:   date,
	}
	return nil
}

func (v *vqueue) PushDelete(uuid string, date int64) error {
	// we have to check this instance's channel bypass daemon is finalizing or not, if in finalizing process we should not send new index to channel
	if v.finalizing.Load().(bool) || v.closed.Load().(bool) {
		return errors.ErrVQueueFinalizing
	}
	if date == 0 {
		date = time.Now().UnixNano()
	}
	v.dch <- key{
		uuid: uuid,
		date: date,
	}
	return nil
}

func (v *vqueue) PopInsert() (uuid string, vector []float32) {
	i := v.popInsert()
	return i.uuid, i.vector
}

func (v *vqueue) PopDelete() (uuid string) {
	d := v.popDelete()
	return d.uuid
}

func (v *vqueue) RangePopInsert(ctx context.Context, f func(uuid string, vector []float32) bool) {
	if v.finalizing.Load().(bool) {
		for !v.finalizing.Load().(bool) {
			time.Sleep(time.Millisecond * 100)
		}
	}
	for _, idx := range v.flushAndLoadInsert() {
		select {
		case <-ctx.Done():
			return
		default:
			if !f(idx.uuid, idx.vector) {
				return
			}
		}
	}
}

func (v *vqueue) RangePopDelete(ctx context.Context, f func(uuid string) bool) {
	if v.finalizing.Load().(bool) {
		for !v.finalizing.Load().(bool) {
			time.Sleep(time.Millisecond * 100)
		}
	}
	for _, key := range v.flushAndLoadDelete() {
		select {
		case <-ctx.Done():
			return
		default:
			if !f(key.uuid) {
				return
			}
		}
	}
}

func (v *vqueue) GetVector(uuid string) ([]float32, bool) {
	v.imu.Lock()
	vec, ok := v.uiil[uuid]
	v.imu.Unlock()
	return vec, ok
}

func (v *vqueue) addInsert(i index) {
	v.imu.Lock()
	v.uii = append(v.uii, i)
	v.uiil[i.uuid] = i.vector
	v.imu.Unlock()
}

func (v *vqueue) addDelete(d key) {
	v.dmu.Lock()
	v.udk = append(v.udk, d)
	v.dmu.Unlock()
}

func (v *vqueue) popInsert() (i index) {
	v.imu.Lock()
	i = v.uii[0]
	v.uii = v.uii[1:]
	delete(v.uiil, i.uuid)
	v.imu.Unlock()
	return i
}

func (v *vqueue) popDelete() (d key) {
	v.dmu.Lock()
	d = v.udk[0]
	v.udk = v.udk[1:]
	v.dmu.Unlock()
	return d
}

func (v *vqueue) flushAndLoadInsert() (uii []index) {
	v.imu.Lock()
	uii = make([]index, len(v.uii))
	copy(uii, v.uii)
	v.uii = v.uii[:0]
	v.imu.Unlock()
	sort.Slice(uii, func(i, j int) bool {
		return uii[i].date > uii[j].date
	})
	dup := make(map[string]bool, len(uii)/2)
	dl := make([]int, 0, len(uii)/2)
	for i, idx := range uii {
		v.imu.Lock()
		delete(v.uiil, idx.uuid)
		v.imu.Unlock()
		if dup[idx.uuid] {
			dl = append(dl, i)
		} else {
			dup[idx.uuid] = true
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(dl)))
	for _, i := range dl {
		uii = append(uii[:i], uii[i+1:]...)
	}
	sort.Slice(uii, func(i, j int) bool {
		return uii[i].date < uii[j].date
	})
	return uii
}

func (v *vqueue) flushAndLoadDelete() (udk []key) {
	v.dmu.Lock()
	udk = make([]key, len(v.udk))
	copy(udk, v.udk)
	v.udk = v.udk[:0]
	v.dmu.Unlock()
	sort.Slice(udk, func(i, j int) bool {
		return udk[i].date > udk[j].date
	})
	dup := make(map[string]bool, len(udk)/2)
	dl := make([]int, 0, len(udk)/2)
	for i, idx := range udk {
		if dup[idx.uuid] {
			dl = append(dl, i)
		} else {
			dup[idx.uuid] = true
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(dl)))
	for _, i := range dl {
		udk = append(udk[:i], udk[i+1:]...)
	}
	sort.Slice(udk, func(i, j int) bool {
		return udk[i].date < udk[j].date
	})

	udm := make(map[string]int64, len(udk))
	for _, d := range udk {
		udm[d.uuid] = d.date
	}
	dl = dl[:0]
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
		// remove unnecessary insert vector queue data
		v.uii = append(v.uii[:i], v.uii[i+1:]...)
		// remove from existing map
		delete(v.uiil, v.uii[i].uuid)
		v.imu.Unlock()
	}
	return udk
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
