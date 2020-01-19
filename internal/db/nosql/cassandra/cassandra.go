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

package cassandra

import (
	"context"
	"crypto/tls"
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

var (
	ErrNotFound = gocql.ErrNotFound
)

type Cassandra interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	Lister
	Getter
	Setter
	Deleter
	Querier
}

type Session = gocql.Session
type Cmp = qb.Cmp
type BatchBuilder = qb.BatchBuilder
type InsertBuilder = qb.InsertBuilder
type DeleteBuilder = qb.DeleteBuilder
type UpdateBuilder = qb.UpdateBuilder
type Queryx = gocqlx.Queryx

type client struct {
	hosts          []string
	cqlVersion     string
	protoVersion   int
	timeout        time.Duration
	connectTimeout time.Duration
	port           int
	keyspace       string
	numConns       int
	consistency    gocql.Consistency
	compressor     gocql.Compressor
	username       string
	password       string
	authProvider   func(h *gocql.HostInfo) (gocql.Authenticator, error)
	retryPolicy    struct {
		numRetries  int
		minDuration time.Duration
		maxDuration time.Duration
	}
	reconnectionPolicy struct {
		initialInterval time.Duration
		maxRetries      int
	}
	poolConfig struct {
		dataCenterName                 string
		enableDCAwareRouting           bool
		enableShuffleReplicas          bool
		enableNonLocalReplicasFallback bool
	}
	socketKeepalive          time.Duration
	maxPreparedStmts         int
	maxRoutingKeyInfo        int
	pageSize                 int
	serialConsistency        gocql.SerialConsistency
	tls                      *tls.Config
	tlsCertPath              string
	tlsKeyPath               string
	tlsCAPath                string
	enableHostVerification   bool
	defaultTimestamp         bool
	reconnectInterval        time.Duration
	maxWaitSchemaAgreement   time.Duration
	ignorePeerAddr           bool
	disableInitialHostLookup bool
	disableNodeStatusEvents  bool
	disableTopologyEvents    bool
	disableSchemaEvents      bool
	disableSkipMetadata      bool
	defaultIdempotence       bool
	writeCoalesceWaitTime    time.Duration
	kvTable                  string
	vkTable                  string

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

	c.cluster = &gocql.ClusterConfig{
		Hosts:          c.hosts,
		CQLVersion:     c.cqlVersion,
		ProtoVersion:   c.protoVersion,
		Timeout:        c.timeout,
		ConnectTimeout: c.connectTimeout,
		Port:           c.port,
		Keyspace:       c.keyspace,
		NumConns:       c.numConns,
		Consistency:    c.consistency,
		Compressor:     c.compressor,
		Authenticator: &gocql.PasswordAuthenticator{
			Username: c.username,
			Password: c.password,
		},
		AuthProvider: c.authProvider,
		RetryPolicy: &gocql.ExponentialBackoffRetryPolicy{
			NumRetries: c.retryPolicy.numRetries,
			Min:        c.retryPolicy.minDuration,
			Max:        c.retryPolicy.maxDuration,
		},
		ConvictionPolicy: &gocql.SimpleConvictionPolicy{},
		ReconnectionPolicy: &gocql.ExponentialReconnectionPolicy{
			MaxRetries:      c.reconnectionPolicy.maxRetries,
			InitialInterval: c.reconnectionPolicy.initialInterval,
		},
		SocketKeepalive:   c.socketKeepalive,
		MaxPreparedStmts:  c.maxPreparedStmts,
		MaxRoutingKeyInfo: c.maxRoutingKeyInfo,
		PageSize:          c.pageSize,
		SerialConsistency: c.serialConsistency,
		DefaultTimestamp:  c.defaultTimestamp,
		PoolConfig: gocql.PoolConfig{
			HostSelectionPolicy: func() (hsp gocql.HostSelectionPolicy) {
				if c.poolConfig.enableDCAwareRouting && len(c.poolConfig.dataCenterName) != 0 {
					hsp = gocql.DCAwareRoundRobinPolicy(c.poolConfig.dataCenterName)
				} else {
					hsp = gocql.RoundRobinHostPolicy()
				}
				switch {
				case c.poolConfig.enableShuffleReplicas &&
					c.poolConfig.enableNonLocalReplicasFallback:
					return gocql.TokenAwareHostPolicy(hsp, gocql.ShuffleReplicas(), gocql.NonLocalReplicasFallback())
				case c.poolConfig.enableShuffleReplicas:
					return gocql.TokenAwareHostPolicy(hsp, gocql.ShuffleReplicas())
				case c.poolConfig.enableNonLocalReplicasFallback:
					return gocql.TokenAwareHostPolicy(hsp, gocql.NonLocalReplicasFallback())
				}
				return gocql.TokenAwareHostPolicy(hsp)
			}(),
		},
		ReconnectInterval:      c.reconnectInterval,
		MaxWaitSchemaAgreement: c.maxWaitSchemaAgreement,
		// HostFilter
		// AddressTranslator
		IgnorePeerAddr:           c.ignorePeerAddr,
		DisableInitialHostLookup: c.disableInitialHostLookup,
		Events: struct {
			DisableNodeStatusEvents bool
			DisableTopologyEvents   bool
			DisableSchemaEvents     bool
		}{
			DisableNodeStatusEvents: c.disableNodeStatusEvents,
			DisableTopologyEvents:   c.disableTopologyEvents,
			DisableSchemaEvents:     c.disableSchemaEvents,
		},
		DisableSkipMetadata: c.disableSkipMetadata,
		// QueryObserver
		// BatchObserver
		// ConnectObserver
		// FrameHeaderObserver
		DefaultIdempotence:    c.defaultIdempotence,
		WriteCoalesceWaitTime: c.writeCoalesceWaitTime,
	}

	if c.tls != nil {
		c.cluster.SslOpts = &gocql.SslOptions{
			Config:                 c.tls,
			CertPath:               c.tlsCertPath,
			KeyPath:                c.tlsKeyPath,
			CaPath:                 c.tlsCAPath,
			EnableHostVerification: c.enableHostVerification,
		}
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

func Select(table string, columns []string, cmps ...Cmp) (stmt string, names []string) {
	sb := qb.Select(table).Columns(columns...)
	for _, cmp := range cmps {
		sb = sb.Where(cmp)
	}
	return sb.ToCql()
}

func Delete(table string, cmps ...Cmp) *DeleteBuilder {
	db := qb.Delete(table)
	for _, cmp := range cmps {
		db = db.Where(cmp)
	}
	return db
}

func Insert(table string, columns ...string) *InsertBuilder {
	return qb.Insert(table).Columns(columns...)
}

func Update(table string) *UpdateBuilder {
	return qb.Update(table)
}

func Batch() *BatchBuilder {
	return qb.Batch()
}

func Eq(column string) Cmp {
	return qb.Eq(column)
}

func Contains(column string) Cmp {
	return qb.Contains(column)
}

func (c *client) Query(stmt string, names []string) *Queryx {
	return gocqlx.Query(c.session.Query(stmt), names)
}

func (c *client) GetValue(key string) (string, error) {
	var value string

	stmt, names := qb.Select(c.kvTable).Columns(metaColumn).Where(qb.Eq(uuidColumn)).ToCql()
	q := gocqlx.Query(c.session.Query(stmt), names).BindMap(qb.M{
		uuidColumn: key,
	})
	if err := q.GetRelease(&value); err != nil {
		switch err {
		case gocql.ErrNotFound:
			return "", errors.ErrCassandraNotFound(key)
		default:
			return "", err
		}
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
		switch err {
		case gocql.ErrNotFound:
			return "", errors.ErrCassandraNotFound(value)
		default:
			return "", err
		}
	}
	return key, nil
}

func (c *client) MultiGetValue(keys ...string) ([]string, error) {
	var keyvals []struct {
		UUID string
		Meta string
	}
	stmt, names := qb.Select(c.kvTable).Columns(uuidColumn, metaColumn).Where(qb.In(uuidColumn)).ToCql()
	q := gocqlx.Query(c.session.Query(stmt), names).BindMap(qb.M{
		uuidColumn: keys,
	})
	if err := q.SelectRelease(&keyvals); err != nil {
		return nil, err
	}

	kvs := make(map[string]string, len(keyvals))
	for _, keyval := range keyvals {
		kvs[keyval.UUID] = keyval.Meta
	}

	var errs error
	values := make([]string, 0, len(keyvals))
	for _, key := range keys {
		if kvs[key] == "" {
			if errs != nil {
				errs = errors.Wrap(errs, errors.ErrCassandraNotFound(key).Error())
			} else {
				errs = errors.ErrCassandraNotFound(key)
			}
			values = append(values, "")
			continue
		}
		values = append(values, kvs[key])
	}

	return values, errs
}

func (c *client) MultiGetKey(values ...string) ([]string, error) {
	var keyvals []struct {
		UUID string
		Meta string
	}
	stmt, names := qb.Select(c.vkTable).Columns(uuidColumn, metaColumn).Where(qb.In(metaColumn)).ToCql()
	q := gocqlx.Query(c.session.Query(stmt), names).BindMap(qb.M{
		metaColumn: values,
	})
	if err := q.SelectRelease(&keyvals); err != nil {
		return nil, err
	}

	kvs := make(map[string]string, len(keyvals))
	for _, keyval := range keyvals {
		kvs[keyval.Meta] = keyval.UUID
	}

	var errs error
	keys := make([]string, 0, len(keyvals))
	for _, value := range values {
		if kvs[value] == "" {
			if errs != nil {
				errs = errors.Wrap(errs, errors.ErrCassandraNotFound(value).Error())
			} else {
				errs = errors.ErrCassandraNotFound(value)
			}
			keys = append(keys, "")
			continue
		}
		keys = append(keys, kvs[value])
	}

	return keys, errs
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

func (c *client) DeleteByKeys(keys ...string) ([]string, error) {
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
