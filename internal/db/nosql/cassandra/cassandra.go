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
	ErrNotFound             = gocql.ErrNotFound
	ErrUnavailable          = gocql.ErrUnavailable
	ErrUnsupported          = gocql.ErrUnsupported
	ErrTooManyStmts         = gocql.ErrTooManyStmts
	ErrUseStmt              = gocql.ErrUseStmt
	ErrSessionClosed        = gocql.ErrSessionClosed
	ErrNoConnections        = gocql.ErrNoConnections
	ErrNoKeyspace           = gocql.ErrNoKeyspace
	ErrKeyspaceDoesNotExist = gocql.ErrKeyspaceDoesNotExist
	ErrNoMetadata           = gocql.ErrNoMetadata
	ErrNoHosts              = gocql.ErrNoHosts
	ErrNoConnectionsStarted = gocql.ErrNoConnectionsStarted
	ErrHostQueryFailed      = gocql.ErrHostQueryFailed
)

type Cassandra interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	Query(stmt string, names []string) *Queryx
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
	dialer                   gocql.Dialer
	writeCoalesceWaitTime    time.Duration

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
		Dialer:            c.dialer,
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

func In(column string) Cmp {
	return qb.In(column)
}

func Contains(column string) Cmp {
	return qb.Contains(column)
}

func (c *client) Query(stmt string, names []string) *Queryx {
	return gocqlx.Query(c.session.Query(stmt), names)
}
