//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"reflect"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/tls"
)

var (
	// ErrNotFound is a alias of gocql.ErrNotFound.
	ErrNotFound = gocql.ErrNotFound
	// ErrUnavailable is a alias of gocql.ErrUnavailable.
	ErrUnavailable = gocql.ErrUnavailable
	// ErrUnsupported is a alias of gocql.ErrUnsupported.
	ErrUnsupported = gocql.ErrUnsupported
	// ErrTooManyStmts is a alias of gocql.ErrTooManyStmts.
	ErrTooManyStmts = gocql.ErrTooManyStmts
	// ErrUseStmt is a alias of gocql.ErrUseStmt.
	ErrUseStmt = gocql.ErrUseStmt
	// ErrSessionClosed is a alias of gocql.ErrSessionClosed.
	ErrSessionClosed = gocql.ErrSessionClosed
	// ErrNoConnections is a alias of gocql.ErrNoConnections.
	ErrNoConnections = gocql.ErrNoConnections
	// ErrNoKeyspace is a alias of gocql.ErrNoKeyspace.
	ErrNoKeyspace = gocql.ErrNoKeyspace
	// ErrKeyspaceDoesNotExist is a alias of gocql.ErrKeyspaceDoesNotExist.
	ErrKeyspaceDoesNotExist = gocql.ErrKeyspaceDoesNotExist
	// ErrNoMetadata is a alias of gocql.ErrNoMetadata.
	ErrNoMetadata = gocql.ErrNoMetadata
	// ErrNoHosts is a alias of gocql.ErrNoHosts.
	ErrNoHosts = gocql.ErrNoHosts
	// ErrNoConnectionsStarted is a alias of gocql.ErrNoConnectionsStarted.
	ErrNoConnectionsStarted = gocql.ErrNoConnectionsStarted
	// ErrHostQueryFailed is a alias of gocql.ErrHostQueryFailed.
	ErrHostQueryFailed = gocql.ErrHostQueryFailed
)

// Cassandra represent an interface to query on cassandra.
type Cassandra interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	Query(stmt string, names []string) *Queryx
}

// ClusterConfig represent an interface of cassandra cluster configuation.
type ClusterConfig interface {
	CreateSession() (*gocql.Session, error)
}

type (
	// Session is a alias of gocql.Session.
	Session = gocql.Session
	// Cmp is a alias of qb.Cmp.
	Cmp = qb.Cmp
	// BatchBuilder is a alias of qb.BatchBuilder.
	BatchBuilder = qb.BatchBuilder
	// InsertBuilder is a alias of qb.InsertBuilder.
	InsertBuilder = qb.InsertBuilder
	// DeleteBuilder is a alias of qb.DeleteBuilder.
	DeleteBuilder = qb.DeleteBuilder
	// UpdateBuilder is a alias of qb.UpdateBuilder.
	UpdateBuilder = qb.UpdateBuilder
	// Queryx is a alias of gocqlx.Queryx.
	Queryx = gocqlx.Queryx
)

type (
	retryPolicy struct {
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
		enableTokenAwareHostPolicy     bool
	}
	hostFilter struct {
		dcHost    string
		whiteList []string
		enable    bool
	}
	// skipcq: SCC-U1000
	events struct {
		DisableNodeStatusEvents bool
		DisableTopologyEvents   bool
		DisableSchemaEvents     bool
	}
	client struct {
		cluster                  ClusterConfig
		dialer                   gocql.Dialer
		rawDialer                net.Dialer
		frameHeaderObserver      FrameHeaderObserver
		connectObserver          ConnectObserver
		batchObserver            BatchObserver
		queryObserver            QueryObserver
		compressor               gocql.Compressor
		tls                      *tls.Config
		session                  *gocql.Session
		authProvider             func(h *gocql.HostInfo) (gocql.Authenticator, error)
		tlsCAPath                string
		tlsCertPath              string
		cqlVersion               string
		keyspace                 string
		password                 string
		tlsKeyPath               string
		username                 string
		poolConfig               poolConfig
		hosts                    []string
		hostFilter               hostFilter
		retryPolicy              retryPolicy
		reconnectionPolicy       reconnectionPolicy
		reconnectInterval        time.Duration
		port                     int
		pageSize                 int
		socketKeepalive          time.Duration
		protoVersion             int
		maxRoutingKeyInfo        int
		maxWaitSchemaAgreement   time.Duration
		writeCoalesceWaitTime    time.Duration
		timeout                  time.Duration
		connectTimeout           time.Duration
		maxPreparedStmts         int
		numConns                 int
		consistency              gocql.Consistency
		serialConsistency        gocql.SerialConsistency
		disableSkipMetadata      bool
		disableSchemaEvents      bool
		disableTopologyEvents    bool
		defaultIdempotence       bool
		disableNodeStatusEvents  bool
		disableInitialHostLookup bool
		ignorePeerAddr           bool
		defaultTimestamp         bool
		enableHostVerification   bool
	}
)

// New initialize and return the cassandra client, or any error occurred.
func New(opts ...Option) (Cassandra, error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(c); err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))

			e := new(errors.ErrCriticalOption)
			if errors.As(err, &e) {
				log.Error(werr)
				return nil, werr
			}
			log.Warn(werr)
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
		Authenticator: func() *gocql.PasswordAuthenticator {
			if len(c.username)+len(c.password) == 0 {
				return nil
			}
			return &gocql.PasswordAuthenticator{
				Username: c.username,
				Password: c.password,
			}
		}(),
		AuthProvider: c.authProvider,
		RetryPolicy: func() *gocql.ExponentialBackoffRetryPolicy {
			if c.retryPolicy.numRetries < 1 {
				return nil
			}
			return &gocql.ExponentialBackoffRetryPolicy{
				NumRetries: c.retryPolicy.numRetries,
				Min:        c.retryPolicy.minDuration,
				Max:        c.retryPolicy.maxDuration,
			}
		}(),
		ConvictionPolicy: NewConvictionPolicy(),
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
				if c.poolConfig.enableTokenAwareHostPolicy {
					return gocql.TokenAwareHostPolicy(hsp)
				}
				return hsp
			}(),
		},
		ReconnectInterval:      c.reconnectInterval,
		MaxWaitSchemaAgreement: c.maxWaitSchemaAgreement,
		HostFilter: func() (hf gocql.HostFilter) {
			if !c.hostFilter.enable {
				return nil
			}

			if len(c.hostFilter.dcHost) != 0 {
				hf = gocql.DataCentreHostFilter(c.hostFilter.dcHost)
			}
			if len(c.hostFilter.whiteList) != 0 {
				wlhf := gocql.WhiteListHostFilter(c.hostFilter.whiteList...)
				if hf == nil {
					hf = wlhf
				} else {
					hf = gocql.HostFilterFunc(func(host *gocql.HostInfo) bool {
						return hf.Accept(host) || wlhf.Accept(host)
					})
				}
			}
			return hf
		}(),
		// AddressTranslator
		IgnorePeerAddr:           c.ignorePeerAddr,
		DisableInitialHostLookup: c.disableInitialHostLookup,
		Events: events{
			DisableNodeStatusEvents: c.disableNodeStatusEvents,
			DisableTopologyEvents:   c.disableTopologyEvents,
			DisableSchemaEvents:     c.disableSchemaEvents,
		},
		DisableSkipMetadata:   c.disableSkipMetadata,
		QueryObserver:         c.queryObserver,
		BatchObserver:         c.batchObserver,
		ConnectObserver:       c.connectObserver,
		FrameHeaderObserver:   c.frameHeaderObserver,
		DefaultIdempotence:    c.defaultIdempotence,
		WriteCoalesceWaitTime: c.writeCoalesceWaitTime,
		SslOpts: func() *gocql.SslOptions {
			if c.tls != nil {
				return &gocql.SslOptions{
					Config:                 c.tls,
					CertPath:               c.tlsCertPath,
					KeyPath:                c.tlsKeyPath,
					CaPath:                 c.tlsCAPath,
					EnableHostVerification: c.enableHostVerification,
				}
			}
			return nil
		}(),
	}

	return c, nil
}

// Open creates a session to cassandra and return any error occurred.
func (c *client) Open(ctx context.Context) (err error) {
	if c.session, err = c.cluster.CreateSession(); err != nil {
		log.Errorf("failed to create session %#v", c)
		return errors.ErrCassandraFailedToCreateSession(err, c.hosts, c.port, c.cqlVersion)
	}
	if c.rawDialer != nil {
		c.rawDialer.StartDialerCache(ctx)
	}
	return nil
}

// Close closes the session to cassandra.
func (c *client) Close(context.Context) error {
	c.session.Close()
	return nil
}

// Query creates an query that can be executed on cassandra.
func (c *client) Query(stmt string, names []string) *Queryx {
	return gocqlx.Query(c.session.Query(stmt), names)
}

// Select build and returns the cql string and the named args.
func Select(table string, columns []string, cmps ...Cmp) (stmt string, names []string) {
	sb := qb.Select(table).Columns(columns...)
	for _, cmp := range cmps {
		sb = sb.Where(cmp)
	}
	return sb.ToCql()
}

// Delete returns the delete builder.
func Delete(table string, cmps ...Cmp) *DeleteBuilder {
	db := qb.Delete(table)
	for _, cmp := range cmps {
		db = db.Where(cmp)
	}
	return db
}

// Insert returns the insert builder.
func Insert(table string, columns ...string) *InsertBuilder {
	return qb.Insert(table).Columns(columns...)
}

// Update returns the update builder.
func Update(table string) *UpdateBuilder {
	return qb.Update(table)
}

// Batch returns the batch builder.
func Batch() *BatchBuilder {
	return qb.Batch()
}

// Eq returns the equal comparator.
func Eq(column string) Cmp {
	return qb.Eq(column)
}

// In returns the in comparator.
func In(column string) Cmp {
	return qb.In(column)
}

// Contains return the contains comparator.
func Contains(column string) Cmp {
	return qb.Contains(column)
}

// WrapErrorWithKeys wraps the cassandra error to Vald internal error.
func WrapErrorWithKeys(err error, keys ...string) error {
	switch {
	case errors.Is(err, ErrNotFound):
		return errors.ErrCassandraNotFound(keys...)
	case errors.Is(err, ErrUnavailable):
		return errors.ErrCassandraUnavailable
	case errors.Is(err, ErrUnsupported):
		return err
	case errors.Is(err, ErrTooManyStmts):
		return err
	case errors.Is(err, ErrUseStmt):
		return err
	case errors.Is(err, ErrSessionClosed):
		return err
	case errors.Is(err, ErrNoConnections):
		return err
	case errors.Is(err, ErrNoKeyspace):
		return err
	case errors.Is(err, ErrKeyspaceDoesNotExist):
		return err
	case errors.Is(err, ErrNoMetadata):
		return err
	case errors.Is(err, ErrNoHosts):
		return err
	case errors.Is(err, ErrNoConnectionsStarted):
		return err
	case errors.Is(err, ErrHostQueryFailed):
		return err
	default:
		return err
	}
}
