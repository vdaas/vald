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

package service

import (
	"context"
	"reflect"
	"strconv"

	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/manager/backup/cassandra/model"
)

const (
	uuidColumn   = "uuid"
	vectorColumn = "vector"
	ipsColumn    = "ips"
)

var columns = []string{uuidColumn, vectorColumn, ipsColumn}

type Cassandra interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	GetVector(ctx context.Context, uuid string) (*model.Vector, error)
	GetIPs(ctx context.Context, uuid string) ([]string, error)
	SetVector(ctx context.Context, vec *model.Vector) error
	SetVectors(ctx context.Context, vecs ...*model.Vector) error
	DeleteVector(ctx context.Context, uuid string) error
	DeleteVectors(ctx context.Context, uuids ...string) error
	SetIPs(ctx context.Context, uuid string, ips ...string) error
	RemoveIPs(ctx context.Context, ips ...string) error
}

type client struct {
	db        cassandra.Cassandra
	tableName string
}

func New(opts ...Option) (Cassandra, error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

func (c *client) Connect(ctx context.Context) error {
	return c.db.Open(ctx)
}

func (c *client) Close(ctx context.Context) error {
	return c.db.Close(ctx)
}

func (c *client) getVector(ctx context.Context, uuid string) (*model.Vector, error) {
	var vector model.Vector
	if err := c.db.Query(cassandra.Select(c.tableName,
		columns,
		cassandra.Eq(uuidColumn))).
		BindMap(map[string]interface{}{
			uuidColumn: uuid,
		}).GetRelease(&vector); err != nil {
		return nil, cassandra.WrapErrorWithKeys(err, uuid)
	}
	return &vector, nil
}

func (c *client) GetVector(ctx context.Context, uuid string) (*model.Vector, error) {
	return c.getVector(ctx, uuid)
}

func (c *client) GetIPs(ctx context.Context, uuid string) ([]string, error) {
	mv, err := c.getVector(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return mv.IPs, nil
}

func (c *client) SetVector(ctx context.Context, vec *model.Vector) error {
	stmt, names := cassandra.Insert(c.tableName, columns...).ToCql()
	return c.db.Query(stmt, names).BindStruct(vec).ExecRelease()
}

func (c *client) SetVectors(ctx context.Context, vecs ...*model.Vector) error {
	ib := cassandra.Insert(c.tableName, columns...)
	bt := cassandra.Batch()

	entities := make(map[string]interface{}, len(vecs)*3)
	for i, mv := range vecs {
		prefix := "p" + strconv.Itoa(i)
		bt = bt.AddWithPrefix(prefix, ib)
		entities[prefix+"."+uuidColumn] = mv.UUID
		entities[prefix+"."+vectorColumn] = mv.Vector
		entities[prefix+"."+ipsColumn] = mv.IPs
	}

	return c.db.Query(bt.ToCql()).BindMap(entities).ExecRelease()
}

func (c *client) DeleteVector(ctx context.Context, uuid string) error {
	return c.db.Query(cassandra.Delete(c.tableName,
		cassandra.Eq(uuidColumn)).ToCql()).
		BindMap(map[string]interface{}{uuidColumn: uuid}).
		ExecRelease()
}

func (c *client) DeleteVectors(ctx context.Context, uuids ...string) error {
	deleteBuilder := cassandra.Delete(c.tableName, cassandra.Eq(uuidColumn))
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
	return c.db.Query(cassandra.Update(c.tableName).
		AddNamed(ipsColumn, ipsColumn).
		Where(cassandra.Eq(uuidColumn)).ToCql()).
		BindMap(map[string]interface{}{
			uuidColumn: uuid,
			ipsColumn:  ips,
		}).ExecRelease()
}

func (c *client) RemoveIPs(ctx context.Context, ips ...string) error {
	var vectors []model.Vector

	for _, ip := range ips {
		err := c.db.Query(cassandra.Select(c.tableName,
			[]string{uuidColumn, ipsColumn},
			cassandra.Contains(ipsColumn))).
			BindMap(map[string]interface{}{ipsColumn: ip}).
			SelectRelease(&vectors)
		if err != nil {
			return err
		}

		for _, mv := range vectors {
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

			err = c.db.Query(cassandra.Update(c.tableName).Set(ipsColumn).
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
