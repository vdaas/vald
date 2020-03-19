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

// Package service manages the main logic of server.
package service

import (
	"context"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/errors"
)

const (
	uuidColumn = "uuid"
	metaColumn = "meta"
)

type Cassandra interface {
	Connect(context.Context) error
	Close(context.Context) error
	Get(string) (string, error)
	GetMultiple(...string) ([]string, error)
	GetInverse(string) (string, error)
	GetInverseMultiple(...string) ([]string, error)
	Set(string, string) error
	SetMultiple(map[string]string) error
	Delete(string) (string, error)
	DeleteMultiple(...string) ([]string, error)
	DeleteInverse(string) (string, error)
	DeleteInverseMultiple(...string) ([]string, error)
}

type client struct {
	db      cassandra.Cassandra
	kvTable string
	vkTable string
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

	if cfg.KVTable == "" {
		cfg.KVTable = "kv"
	}
	if cfg.VKTable == "" {
		cfg.VKTable = "vk"
	}

	return &client{
		db:      db,
		kvTable: cfg.KVTable,
		vkTable: cfg.VKTable,
	}, nil
}

func (c *client) Connect(ctx context.Context) error {
	return c.db.Open(ctx)
}

func (c *client) Close(ctx context.Context) error {
	return c.db.Close(ctx)
}

func wrapErrorWithKeys(err error, keys ...string) error {
	switch err {
	case cassandra.ErrNotFound:
		return errors.ErrCassandraNotFound(keys...)
	case cassandra.ErrUnavailable:
		return errors.ErrCassandraUnavailable()
	default:
		return err
	}
}

func (c *client) Get(key string) (string, error) {
	var val string
	if err := c.db.Query(cassandra.Select(c.kvTable,
		[]string{metaColumn},
		cassandra.Eq(uuidColumn))).BindMap(map[string]interface{}{
		uuidColumn: key,
	}).GetRelease(&val); err != nil {
		return "", wrapErrorWithKeys(err, key)
	}
	return val, nil
}

func (c *client) GetMultiple(keys ...string) (vals []string, err error) {
	var keyvals []struct {
		UUID string
		Meta string
	}
	if err = c.db.Query(cassandra.Select(c.kvTable,
		[]string{uuidColumn, metaColumn},
		cassandra.In(uuidColumn))).BindMap(map[string]interface{}{
		uuidColumn: keys,
	}).SelectRelease(&keyvals); err != nil {
		return nil, wrapErrorWithKeys(err, keys...)
	}

	kvs := make(map[string]string, len(keyvals))
	for _, keyval := range keyvals {
		kvs[keyval.UUID] = keyval.Meta
	}

	vals = make([]string, 0, len(keyvals))
	for _, key := range keys {
		if kvs[key] == "" {
			if err != nil {
				err = errors.Wrap(err, errors.ErrCassandraNotFound(key).Error())
			} else {
				err = errors.ErrCassandraNotFound(key)
			}
			vals = append(vals, "")
			continue
		}
		vals = append(vals, kvs[key])
	}
	if err != nil {
		return nil, err
	}
	return vals, nil
}

func (c *client) GetInverse(val string) (string, error) {
	var key string
	if err := c.db.Query(cassandra.Select(c.vkTable,
		[]string{uuidColumn},
		cassandra.Eq(metaColumn))).BindMap(map[string]interface{}{
		metaColumn: val,
	}).GetRelease(&key); err != nil {
		return "", wrapErrorWithKeys(err, val)
	}
	return key, nil
}

func (c *client) GetInverseMultiple(vals ...string) (keys []string, err error) {
	var keyvals []struct {
		UUID string
		Meta string
	}
	if err = c.db.Query(cassandra.Select(c.vkTable,
		[]string{uuidColumn, metaColumn},
		cassandra.In(metaColumn))).BindMap(map[string]interface{}{
		metaColumn: vals,
	}).SelectRelease(&keyvals); err != nil {
		return nil, wrapErrorWithKeys(err, vals...)
	}

	kvs := make(map[string]string, len(keyvals))
	for _, keyval := range keyvals {
		kvs[keyval.Meta] = keyval.UUID
	}

	keys = make([]string, 0, len(keyvals))
	for _, val := range vals {
		if kvs[val] == "" {
			if err != nil {
				err = errors.Wrap(err, errors.ErrCassandraNotFound(val).Error())
			} else {
				err = errors.ErrCassandraNotFound(val)
			}
			keys = append(keys, "")
			continue
		}
		keys = append(keys, kvs[val])
	}
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (c *client) Set(key, val string) error {
	return c.db.Query(
		cassandra.Batch().Add(
			cassandra.Insert(c.kvTable).Columns(uuidColumn, metaColumn),
		).Add(
			cassandra.Insert(c.vkTable).Columns(uuidColumn, metaColumn),
		).ToCql(),
	).Bind(
		[]interface{}{key, val, key, val}...,
	).ExecRelease()
}

func (c *client) SetMultiple(kvs map[string]string) (err error) {
	kvi := cassandra.Insert(c.kvTable).Columns(uuidColumn, metaColumn)
	vki := cassandra.Insert(c.vkTable).Columns(uuidColumn, metaColumn)

	bt := cassandra.Batch()
	entities := make([]interface{}, 0, len(kvs)*4)
	for key, val := range kvs {
		bt = bt.Add(kvi).Add(vki)
		entities = append(entities, key, val, key, val)
	}

	return c.db.Query(bt.ToCql()).Bind(entities...).ExecRelease()
}

func (c *client) deleteByKeys(keys ...string) ([]string, error) {
	vals, err := c.GetMultiple(keys...)
	if err != nil {
		return nil, err
	}
	kvd := cassandra.Delete(c.kvTable).Where(cassandra.Eq(uuidColumn))
	vkd := cassandra.Delete(c.vkTable).Where(cassandra.Eq(metaColumn))

	bt := cassandra.Batch()
	uuids := make([]interface{}, 0, len(keys)*2)
	for i, key := range keys {
		bt = bt.Add(kvd).Add(vkd)
		uuids = append(uuids, key, vals[i])
	}

	err = c.db.Query(bt.ToCql()).Bind(uuids...).ExecRelease()
	if err != nil {
		return nil, err
	}
	return vals, nil
}

func (c *client) Delete(key string) (string, error) {
	vals, err := c.deleteByKeys(key)
	if err != nil {
		return "", err
	}

	if len(vals) != 1 {
		return "", errors.ErrCassandraDeleteOperationFailed(key, nil)
	}

	return vals[0], nil
}

func (c *client) DeleteMultiple(keys ...string) ([]string, error) {
	return c.deleteByKeys(keys...)
}

func (c *client) deleteByValues(vals ...string) ([]string, error) {
	keys, err := c.GetInverseMultiple(vals...)
	if err != nil {
		return nil, err
	}
	kvd := cassandra.Delete(c.kvTable).Where(cassandra.Eq(uuidColumn))
	vkd := cassandra.Delete(c.vkTable).Where(cassandra.Eq(metaColumn))

	bt := cassandra.Batch()
	metas := make([]interface{}, 0, len(vals)*2)
	for i, val := range vals {
		bt = bt.Add(kvd).Add(vkd)
		metas = append(metas, keys[i], val)
	}

	err = c.db.Query(bt.ToCql()).Bind(metas...).ExecRelease()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (c *client) DeleteInverse(val string) (string, error) {
	keys, err := c.deleteByValues(val)
	if err != nil {
		return "", err
	}

	if len(keys) != 1 {
		return "", errors.ErrCassandraDeleteOperationFailed(val, nil)
	}

	return keys[0], nil
}

func (c *client) DeleteInverseMultiple(vals ...string) ([]string, error) {
	return c.deleteByValues(vals...)
}
