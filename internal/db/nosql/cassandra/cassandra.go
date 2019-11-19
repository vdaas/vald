//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
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

package cassandra

import (
	"context"
	"reflect"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/vdaas/vald/internal/errors"
)

const (
	uuidColumn = "uuid"
	metaColumn = "meta"
)

type Cassandra interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	Lister
	Getter
	Setter
	Deleter
}

type client struct {
	hosts                  []string
	cqlVersion             string
	timeout                time.Duration
	connectTimeout         time.Duration
	port                   int
	numConns               int
	consistency            gocql.Consistency
	maxPreparedStmts       int
	maxRoutingKeyInfo      int
	pageSize               int
	defaultTimestamp       bool
	maxWaitSchemaAgreement time.Duration
	reconnectInterval      time.Duration
	reconnectionPolicy     struct {
		initialInterval time.Duration
		maxRetries      int
	}
	writeCoalesceWaitTime time.Duration
	keyspace              string
	kvTable               string
	vkTable               string

	username string
	password string

	cluster *gocql.ClusterConfig
	session *gocql.Session
}

func New(opts ...Option) (Cassandra, error) {
	c := new(client)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	// FIXME: curretly only support fields in https://github.com/gocql/gocql/blob/ae2f7fc85f32248f9341a280ccdad16b44581f36/cluster.go#L162-L177
	c.cluster = &gocql.ClusterConfig{
		Hosts:                  c.hosts,
		CQLVersion:             c.cqlVersion,
		Timeout:                c.timeout,
		ConnectTimeout:         c.connectTimeout,
		Port:                   c.port,
		NumConns:               c.numConns,
		Consistency:            c.consistency,
		MaxPreparedStmts:       c.maxPreparedStmts,
		MaxRoutingKeyInfo:      c.maxRoutingKeyInfo,
		PageSize:               c.pageSize,
		DefaultTimestamp:       c.defaultTimestamp,
		MaxWaitSchemaAgreement: c.maxWaitSchemaAgreement,
		ReconnectInterval:      c.reconnectInterval,
		ConvictionPolicy:       &gocql.SimpleConvictionPolicy{},
		ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{
			MaxRetries:      c.reconnectionPolicy.maxRetries,
			InitialInterval: c.reconnectionPolicy.initialInterval,
		},
		WriteCoalesceWaitTime: c.writeCoalesceWaitTime,
		Keyspace:              c.keyspace,
		Authenticator: &gocql.PasswordAuthenticator{
			Username: c.username,
			Password: c.password,
		},
	}

	return c, nil
}

func (c *client) Open(ctx context.Context) error {
	session, err := c.cluster.CreateSession()
	if err != nil {
		return err
	}

	c.session = session

	return nil
}

func (c *client) Close(ctx context.Context) error {
	c.session.Close()
	return nil
}

func (c *client) GetValue(key string) (string, error) {
	var value string

	stmt, names := qb.Select(c.kvTable).Columns(metaColumn).Where(qb.Eq(uuidColumn)).ToCql()
	q := gocqlx.Query(c.session.Query(stmt), names).BindMap(qb.M{
		uuidColumn: key,
	})
	if err := q.GetRelease(&value); err != nil {
		return "", err
	}
	return value, nil
}

func (c *client) GetKey(value string) (string, error) {
	var key string

	stmt, names := qb.Select(c.vkTable).Columns(uuidColumn).Where(qb.Eq(metaColumn)).ToCql()
	q := gocqlx.Query(c.session.Query(stmt), names).BindMap(qb.M{
		metaColumn: value,
	})
	if err := q.GetRelease(&key); err != nil {
		return "", err
	}
	return key, nil
}

func (c *client) MultiGetValue(keys ...string) ([]string, error) {
	var values []string
	stmt, names := qb.Select(c.kvTable).Columns(metaColumn).Where(qb.In(uuidColumn)).ToCql()
	q := gocqlx.Query(c.session.Query(stmt), names).BindMap(qb.M{
		uuidColumn: keys,
	})
	if err := q.SelectRelease(&values); err != nil {
		return nil, err
	}

	return values, nil
}

func (c *client) MultiGetKey(values ...string) ([]string, error) {
	var keys []string
	stmt, names := qb.Select(c.vkTable).Columns(uuidColumn).Where(qb.In(metaColumn)).ToCql()
	q := gocqlx.Query(c.session.Query(stmt), names).BindMap(qb.M{
		metaColumn: values,
	})
	if err := q.SelectRelease(&keys); err != nil {
		return nil, err
	}

	return keys, nil
}

func (c *client) Set(key, value string) error {
	stmt, names := qb.Insert(c.kvTable).Columns(uuidColumn, metaColumn).ToCql()
	q := gocqlx.Query(c.session.Query(stmt), names).BindMap(qb.M{
		uuidColumn: key,
		metaColumn: value,
	})

	if err := q.ExecRelease(); err != nil {
		return err
	}

	stmt, names = qb.Insert(c.vkTable).Columns(metaColumn, uuidColumn).ToCql()
	q = gocqlx.Query(c.session.Query(stmt), names).BindMap(qb.M{
		metaColumn: value,
		uuidColumn: key,
	})

	return q.ExecRelease()
}

func (c *client) MultiSet(keyvals map[string]string) error {

	kvi := qb.Insert(c.kvTable).Columns(uuidColumn, metaColumn)
	vki := qb.Insert(c.vkTable).Columns(uuidColumn, metaColumn)

	bt := qb.Batch()
	entities := make([]interface{}, 0, len(keyvals)*4)
	for key, val := range keyvals {
		bt = bt.Add(kvi).Add(vki)
		entities = append(entities, key, val, key, val)
	}

	stmt, names := bt.ToCql()
	q := gocqlx.Query(c.session.Query(stmt), names).Bind(entities...)

	return q.ExecRelease()
}

func (c *client) Delete(keys ...string) ([]string, error) {
	vals, err := c.MultiGetValue(keys...)
	if err != nil {
		return nil, err
	}
	kvd := qb.Delete(c.kvTable).Where(qb.Eq(uuidColumn))
	vkd := qb.Delete(c.vkTable).Where(qb.Eq(metaColumn))

	bt := qb.Batch()
	uuids := make([]interface{}, 0, len(keys)*2)
	for i, key := range keys {
		bt = bt.Add(kvd).Add(vkd)
		uuids = append(uuids, key, vals[i])
	}

	stmt, names := bt.ToCql()
	q := gocqlx.Query(c.session.Query(stmt), names).Bind(uuids...)
	err = q.ExecRelease()
	if err != nil {
		return nil, err
	}

	return vals, nil
}

func (c *client) DeleteByValues(values ...string) ([]string, error) {
	keys, err := c.MultiGetKey(values...)
	if err != nil {
		return nil, err
	}
	kvd := qb.Delete(c.kvTable).Where(qb.Eq(uuidColumn))
	vkd := qb.Delete(c.vkTable).Where(qb.Eq(metaColumn))

	bt := qb.Batch()
	metas := make([]interface{}, 0, len(values)*2)
	for i, value := range values {
		bt = bt.Add(kvd).Add(vkd)
		metas = append(metas, keys[i], value)
	}

	stmt, names := bt.ToCql()
	q := gocqlx.Query(c.session.Query(stmt), names).Bind(metas...)
	err = q.ExecRelease()
	if err != nil {
		return nil, err
	}

	return keys, nil
}
