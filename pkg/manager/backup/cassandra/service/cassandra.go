//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package service

import (
	"context"
	"strconv"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/manager/backup/cassandra/model"
)

const (
	uuidColumn   = "uuid"
	vectorColumn = "vector"
	metaColumn   = "meta"
	ipsColumn    = "ips"
)

var (
	metaColumns = []string{uuidColumn, vectorColumn, metaColumn, ipsColumn}
)

type Cassandra interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	GetMeta(ctx context.Context, uuid string) (*model.MetaVector, error)
	GetIPs(ctx context.Context, uuid string) ([]string, error)
	SetMeta(ctx context.Context, meta *model.MetaVector) error
	SetMetas(ctx context.Context, metas ...*model.MetaVector) error
	DeleteMeta(ctx context.Context, uuid string) error
	DeleteMetas(ctx context.Context, uuids ...string) error
	SetIPs(ctx context.Context, uuid string, ips ...string) error
	RemoveIPs(ctx context.Context, ips ...string) error
}

type client struct {
	db        cassandra.Cassandra
	metaTable string
}

func New(cfg *config.Cassandra) (Cassandra, error) {
	opts, err := cfg.Opts()
	if err != nil {
		return nil, err
	}

	db, err := cassandra.New(opts...)
	if err != nil {
		return nil, err
	}

	if cfg.MetaTable == "" {
		cfg.MetaTable = "meta_vector"
	}

	return &client{
		db:        db,
		metaTable: cfg.MetaTable,
	}, nil
}

func (c *client) Connect(ctx context.Context) error {
	return c.db.Open(ctx)
}

func (c *client) Close(ctx context.Context) error {
	return c.db.Close(ctx)
}

func (c *client) getMetaVector(ctx context.Context, uuid string) (*model.MetaVector, error) {
	var metaVector model.MetaVector
	switch err := c.db.Query(cassandra.Select(c.metaTable,
		metaColumns,
		cassandra.Eq(uuidColumn))).
		BindMap(map[string]interface{}{
			uuidColumn: uuid,
		}).GetRelease(&metaVector); err {
	case cassandra.ErrNotFound:
		return nil, errors.ErrCassandraNotFound(uuid)
	case nil:
		return &metaVector, nil
	default:
		return nil, err
	}
}

func (c *client) GetMeta(ctx context.Context, uuid string) (*model.MetaVector, error) {
	return c.getMetaVector(ctx, uuid)
}

func (c *client) GetIPs(ctx context.Context, uuid string) ([]string, error) {
	mv, err := c.getMetaVector(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return mv.IPs, nil
}

func (c *client) SetMeta(ctx context.Context, meta *model.MetaVector) error {
	stmt, names := cassandra.Insert(c.metaTable, metaColumns...).ToCql()
	return c.db.Query(stmt, names).BindStruct(meta).ExecRelease()
}

func (c *client) SetMetas(ctx context.Context, metas ...*model.MetaVector) error {
	ib := cassandra.Insert(c.metaTable, metaColumns...)
	bt := cassandra.Batch()

	entities := make(map[string]interface{}, len(metas)*4)
	for i, mv := range metas {
		prefix := "p" + strconv.Itoa(i)
		bt = bt.AddWithPrefix(prefix, ib)
		entities[prefix+"."+uuidColumn] = mv.UUID
		entities[prefix+"."+vectorColumn] = mv.Vector
		entities[prefix+"."+metaColumn] = mv.Meta
		entities[prefix+"."+ipsColumn] = mv.IPs
	}

	return c.db.Query(bt.ToCql()).BindMap(entities).ExecRelease()
}

func (c *client) DeleteMeta(ctx context.Context, uuid string) error {
	return c.db.Query(cassandra.Delete(c.metaTable,
		cassandra.Eq(uuidColumn)).ToCql()).
		BindMap(map[string]interface{}{uuidColumn: uuid}).
		ExecRelease()
}

func (c *client) DeleteMetas(ctx context.Context, uuids ...string) error {
	deleteBuilder := cassandra.Delete(c.metaTable, cassandra.Eq(uuidColumn))
	bt := cassandra.Batch()
	bindUUIDs := make(map[string]interface{}, len(uuids))
	for i, uuid := range uuids {
		prefix := "p" + strconv.Itoa(i)
		bt.AddWithPrefix(prefix, deleteBuilder)
		bindUUIDs[prefix+"."+uuidColumn] = uuid
	}

	return c.db.Query(bt.ToCql()).BindMap(bindUUIDs).ExecRelease()
}

func (c *client) SetIPs(ctx context.Context, uuid string, ips ...string) error {
	return c.db.Query(cassandra.Update(c.metaTable).
		AddNamed(ipsColumn, ipsColumn).
		Where(cassandra.Eq(uuidColumn)).ToCql()).
		BindMap(map[string]interface{}{
			uuidColumn: uuid,
			ipsColumn:  ips,
		}).ExecRelease()
}

func (c *client) RemoveIPs(ctx context.Context, ips ...string) error {
	var metaVectors []model.MetaVector

	for _, ip := range ips {
		err := c.db.Query(cassandra.Select(c.metaTable,
			[]string{uuidColumn, ipsColumn},
			cassandra.Contains(ipsColumn))).
			BindMap(map[string]interface{}{ipsColumn: ip}).
			SelectRelease(&metaVectors)
		if err != nil {
			return err
		}

		for _, mv := range metaVectors {
			currentIPs := mv.IPs
			newIPs := make([]string, 0, len(currentIPs)-1)
			for i, cIP := range currentIPs {
				if cIP == ip {
					if i != len(currentIPs) {
						newIPs = append(newIPs, currentIPs[i+1:]...)
					}
					break
				}
				newIPs = append(newIPs, cIP)
			}

			err = c.db.Query(cassandra.Update(c.metaTable).Set(ipsColumn).
				Where(cassandra.Eq(uuidColumn)).ToCql()).
				BindMap(map[string]interface{}{
					uuidColumn: mv.UUID,
					ipsColumn:  newIPs,
				}).ExecRelease()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
