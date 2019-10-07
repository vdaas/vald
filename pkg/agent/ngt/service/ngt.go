//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package service manages the main logic of server.
package service

import (
	"os"
	"sync/atomic"

	"github.com/vdaas/vald/internal/config"
	core "github.com/vdaas/vald/internal/core/ngt"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/agent/ngt/model"
	"github.com/vdaas/vald/pkg/agent/ngt/service/kvs"
)

type NGT interface {
	Search(vec []float64, size uint32, epsilon, radius float32) ([]model.Distance, error)
	SearchByID(uuid string, size uint32, epsilon, radius float32) ([]model.Distance, error)
	Insert(uuid string, vec []float64) (err error)
	Update(uuid string, vec []float64) (err error)
	Delete(uuid string) (err error)
	GetObject(uuid string) (vec []float64, err error)
	CreateIndex(poolSize uint32) (err error)
	SaveIndex() (err error)
	Exists(string) (uint32, bool)
	CreateAndSaveIndex(poolSize uint32) (err error)
	ObjectCount() uint64
	Close()
}

type ngt struct {
	cflg uint32 // create index flag 0 or 1
	ic   uint64 // insert count
	kvs  kvs.BidiMap
	core core.NGT
}

func New(cfg *config.NGT) (nn NGT, err error) {
	n := new(ngt)
	opts := []core.Option{
		core.WithIndexPath(cfg.IndexPath),
		core.WithDimension(cfg.Dimension),
		core.WithDistanceTypeByString(cfg.DistanceType),
		core.WithObjectTypeByString(cfg.ObjectType),
		core.WithBulkInsertChunkSize(cfg.BulkInsertChunkSize),
		core.WithCreationEdgeSize(cfg.CreationEdgeSize),
		core.WithSearchEdgeSize(cfg.SearchEdgeSize),
	}

	n.kvs = kvs.New()

	if _, err = os.Stat(cfg.IndexPath); os.IsNotExist(err) {
		n.core, err = core.New(opts...)
	} else {
		n.core, err = core.Load(opts...)
	}

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *ngt) Search(vec []float64, size uint32, epsilon, radius float32) ([]model.Distance, error) {
	sr, err := n.core.Search(vec, int(size), epsilon, radius)
	if err != nil {
		return nil, err
	}

	var errs error
	ds := make([]model.Distance, 0, len(sr))

	for _, d := range sr {
		if err = d.Error; d.ID == 0 && err != nil {
			errs = errors.Wrap(errs, err.Error())
			continue
		}
		key, ok := n.kvs.GetInverse(d.ID)
		if ok {
			ds = append(ds, model.Distance{
				ID:       key,
				Distance: d.Distance,
			})
		} else {
			errs = errors.Wrap(errs, errors.ErrUUIDNotFound(d.ID).Error())
		}
	}

	return ds[:len(ds)], errs
}

func (n *ngt) SearchByID(uuid string, size uint32, epsilon, radius float32) ([]model.Distance, error) {
	oid, ok := n.kvs.Get(uuid)
	if !ok {
		return nil, errors.ErrObjectIDNotFound(uuid)
	}
	vec, err := n.core.GetVector(uint(oid))
	if err != nil {
		return nil, errors.ErrObjectNotFound(err, uuid)
	}

	return n.Search(vec, size, epsilon, radius)
}

func (n *ngt) Insert(uuid string, vec []float64) (err error) {
	if len(uuid) == 0 {
		err = errors.ErrUUIDNotFound(0)
		return err
	}

	id, ok := n.kvs.Get(uuid)
	oid := uint(id)
	if ok {
		err = errors.ErrUUIDAlreadyExists(uuid, oid)
		return err
	}

	oid, err = n.core.Insert(vec)
	if err != nil {
		return err
	}

	atomic.AddUint64(&n.ic, 1)
	n.kvs.Set(uuid, uint32(oid))
	return nil
}

func (n *ngt) Update(uuid string, vec []float64) (err error) {
	err = n.Delete(uuid)
	if err != nil {
		return err
	}
	oid, err := n.core.Insert(vec)
	if err != nil {
		return err
	}

	n.kvs.Set(uuid, uint32(oid))
	return nil
}

func (n *ngt) Delete(uuid string) (err error) {
	if len(uuid) == 0 {
		err = errors.ErrUUIDNotFound(0)
		return err
	}

	if ic := atomic.LoadUint64(&n.ic); ic > 0 {
		return errors.ErrUncommittedIndexExists(ic)
	}
	oid, ok := n.kvs.Get(uuid)
	if !ok {
		err = errors.ErrObjectIDNotFound(uuid)
		return err
	}
	err = n.core.Remove(uint(oid))
	if err != nil {
		return err
	}

	n.kvs.DeleteInverse(oid)

	return nil
}

func (n *ngt) GetObject(uuid string) (vec []float64, err error) {
	oid, ok := n.kvs.Get(uuid)
	if !ok {
		err = errors.ErrObjectIDNotFound(uuid)
		return nil, err
	}
	return n.core.GetVector(uint(oid))
}

func (n *ngt) CreateIndex(poolSize uint32) (err error) {
	if ic := atomic.LoadUint64(&n.ic); ic == 0 {
		return errors.ErrUncommittedIndexNotFound
	}
	defer func() {
		if err == nil {
			// TODO CreateIndexが失敗した場合にObjectSpaceがどうなるか岩崎さんに確認
			atomic.StoreUint64(&n.ic, 0)
			atomic.StoreUint32(&n.cflg, 1)
		}
	}()
	err = n.core.CreateIndex(poolSize)
	return
}

func (n *ngt) SaveIndex() (err error) {
	if atomic.LoadUint32(&n.cflg) == 0 {
		return errors.ErrUncommittedIndexExists(atomic.LoadUint64(&n.ic))
	}
	defer func() {
		if err == nil {
			atomic.StoreUint32(&n.cflg, 0)
		}
	}()
	return n.core.SaveIndex()
}

func (n *ngt) CreateAndSaveIndex(poolSize uint32) (err error) {
	err = n.CreateIndex(poolSize)
	if err != nil {
		return err
	}

	return n.SaveIndex()
}

func (n *ngt) Close() {
	n.core.Close()
}

func (n *ngt) Exists(uuid string) (uint32, bool) {
	oid, ok := n.kvs.Get(uuid)
	if !ok {
		return 0, false
	}

	return oid, true
}

func (n *ngt) ObjectCount() uint64 {
	return n.kvs.Len()
}
