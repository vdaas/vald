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
	"strconv"

	"github.com/kpango/gache"
	"github.com/vdaas/vald/internal/config"
	core "github.com/vdaas/vald/internal/core/ngt"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/pkg/manager/backup/model"
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
	Exists(string) (string, bool)
	CreateAndSaveIndex(poolSize uint32) (err error)
	Close()
}

type ngt struct {
	ou   gache.Gache // map[oid]uuid
	uo   gache.Gache // map[uuid]oid
	core core.NGT
}

func NewNGT(cfg *config.NGT) (NGT, error) {

	var (
		n   core.NGT
		err error

		opts = []core.Option{
			core.WithIndexPath(cfg.IndexPath),
			core.WithDimension(cfg.Dimension),
			core.WithDistanceTypeByString(cfg.DistanceType),
			core.WithObjectTypeByString(cfg.ObjectType),
			core.WithBulkInsertChunkSize(cfg.BulkInsertChunkSize),
			core.WithCreationEdgeSize(cfg.CreationEdgeSize),
			core.WithSearchEdgeSize(cfg.SearchEdgeSize),
		}

		ou = gache.New().
			SetDefaultExpire(0).
			DisableExpiredHook()

		uo = gache.New().
			SetDefaultExpire(0).
			DisableExpiredHook()
	)

	if _, err := os.Stat(cfg.IndexPath); os.IsNotExist(err) {
		n, err = core.New(opts...)
	} else {
		n, err = core.Load(opts...)
	}

	if err != nil {
		return nil, err
	}

	return &ngt{
		ou:   ou,
		uo:   uo,
		core: n,
	}, nil
}

func (n *ngt) Search(vec []float64, size uint32, epsilon, radius float32) ([]model.Distance, error) {
	sr, err := n.core.Search(vec, int(size), epsilon, radius)
	if err != nil {
		return nil, err
	}

	var (
		ds   = make([]model.Distance, 0, len(sr))
		errs error
	)

	for i, d := range sr {
		if err = d.Error; d.ID == 0 && err != nil {
			errs = errors.Wrap(errs, err.Error())
			continue
		}
		key, ok := n.ou.Get(strconv.FormatInt(int64(d.ID), 10))
		if ok {
			ds[i] = model.Distance{
				ID:       key.(string),
				Distance: d.Distance,
			}
		} else {
			log.Warn(errors.ErrUUIDNotFound(d.ID))
		}
	}

	return ds[:len(ds)], errs
}

func (n *ngt) SearchByID(uuid string, size uint32, epsilon, radius float32) ([]model.Distance, error) {
	oid, ok := n.uo.Get(uuid)
	if !ok {
		return nil, errors.ErrObjectIDNotFound(uuid)
	}
	vec, err := n.core.GetVector(oid.(uint))
	if err != nil {
		return nil, errors.ErrObjectNotFound(err, uuid)
	}

	return n.Search(vec, size, epsilon, radius)
}

func (n *ngt) Insert(uuid string, vec []float64) (err error) {
	i, ok := n.uo.Get(uuid)
	if ok && i != 0 {
		err = errors.ErrUUIDAlreadyExists(uuid, i.(uint32))
		return err
	}

	oid, err := n.core.Insert(vec)
	if err != nil {
		return err
	}

	n.uo.SetWithExpire(uuid, oid, 0)
	n.ou.SetWithExpire(strconv.FormatInt(int64(oid), 10), uuid, 0)

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

	n.uo.SetWithExpire(uuid, oid, 0)
	n.ou.SetWithExpire(strconv.FormatInt(int64(oid), 10), uuid, 0)

	return nil
}

func (n *ngt) Delete(uuid string) (err error) {
	i, ok := n.uo.Get(uuid)
	oid := i.(uint)
	if !ok || oid == 0 {
		err = errors.ErrObjectIDNotFound(uuid)
		return err
	}
	err = n.core.Remove(oid)
	if err != nil {
		return err
	}

	n.uo.Delete(uuid)
	n.ou.Delete(strconv.FormatInt(int64(oid), 10))

	return nil
}

func (n *ngt) GetObject(uuid string) (vec []float64, err error) {
	i, ok := n.uo.Get(uuid)

	if !ok || i == 0 {
		err = errors.ErrObjectIDNotFound(uuid)
		return nil, err
	}

	return n.core.GetVector(i.(uint))
}

func (n *ngt) CreateIndex(poolSize uint32) (err error) {
	return n.core.CreateIndex(poolSize)
}

func (n *ngt) SaveIndex() (err error) {
	return n.core.SaveIndex()
}

func (n *ngt) CreateAndSaveIndex(poolSize uint32) (err error) {
	return n.core.CreateAndSaveIndex(poolSize)
}

func (n *ngt) Close() {
	n.core.Close()
	n.ou.Stop()
	n.ou.Clear()
	n.uo.Stop()
	n.uo.Clear()
}

func (n *ngt) Exists(uuid string) (string, bool) {
	oid, ok := n.uo.Get(uuid)
	if !ok {
		return "", false
	}

	return oid.(string), true
}
