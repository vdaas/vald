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

package service

import (
	"context"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/tls"
	"github.com/vdaas/vald/pkg/manager/backup/cassandra/model"
)

type Cassandra interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	GetMeta(ctx context.Context, uuid string) (*model.MetaVector, error)
	GetIPs(ctx context.Context, uuid string) ([]string, error)
	SetMeta(ctx context.Context, meta model.MetaVector) error
	SetMetas(ctx context.Context, metas ...model.MetaVector) error
	DeleteMeta(ctx context.Context, uuid string) error
	DeleteMetas(ctx context.Context, uuids ...string) error
	SetIPs(ctx context.Context, uuid string, ips ...string) error
	RemoveIPs(ctx context.Context, ips ...string) error
}

type client struct {
	db          cassandra.Cassandra
	metaTable   string
	metaColumns []string
}

func New(cfg *config.Cassandra) (Cassandra, error) {
	opts := []cassandra.Option{
		cassandra.WithHosts(cfg.Hosts...),
		cassandra.WithCQLVersion(cfg.CQLVersion),
		cassandra.WithProtoVersion(cfg.ProtoVersion),
		cassandra.WithTimeout(cfg.Timeout),
		cassandra.WithConnectTimeout(cfg.ConnectTimeout),
		cassandra.WithPort(cfg.Port),
		cassandra.WithKeyspace(cfg.Keyspace),
		cassandra.WithNumConns(cfg.NumConns),
		cassandra.WithConsistency(cfg.Consistency),
		cassandra.WithUsername(cfg.Username),
		cassandra.WithPassword(cfg.Password),
		cassandra.WithRetryPolicyNumRetries(cfg.RetryPolicy.NumRetries),
		cassandra.WithRetryPolicyMinDuration(cfg.RetryPolicy.MinDuration),
		cassandra.WithRetryPolicyMaxDuration(cfg.RetryPolicy.MaxDuration),
		cassandra.WithReconnectionPolicyMaxRetries(cfg.ReconnectionPolicy.MaxRetries),
		cassandra.WithReconnectionPolicyInitialInterval(cfg.ReconnectionPolicy.InitialInterval),
		cassandra.WithSocketKeepalive(cfg.SocketKeepalive),
		cassandra.WithMaxPreparedStmts(cfg.MaxPreparedStmts),
		cassandra.WithMaxRoutingKeyInfo(cfg.MaxRoutingKeyInfo),
		cassandra.WithPageSize(cfg.PageSize),
		cassandra.WithEnableHostVerification(cfg.EnableHostVerification),
		cassandra.WithDefaultTimestamp(cfg.DefaultTimestamp),
		cassandra.WithReconnectInterval(cfg.ReconnectInterval),
		cassandra.WithMaxWaitSchemaAgreement(cfg.MaxWaitSchemaAgreement),
		cassandra.WithIgnorePeerAddr(cfg.IgnorePeerAddr),
		cassandra.WithDisableInitialHostLookup(cfg.DisableInitialHostLookup),
		cassandra.WithDisableNodeStatusEvents(cfg.DisableNodeStatusEvents),
		cassandra.WithDisableTopologyEvents(cfg.DisableTopologyEvents),
		cassandra.WithDisableSkipMetadata(cfg.DisableSkipMetadata),
		cassandra.WithDefaultIdempotence(cfg.DefaultIdempotence),
		cassandra.WithWriteCoalesceWaitTime(cfg.WriteCoalesceWaitTime),
		cassandra.WithKVTable(cfg.KVTable),
		cassandra.WithVKTable(cfg.VKTable),
	}

	if cfg.TLS != nil && cfg.TLS.Enabled {
		tcfg, err := tls.New(
			tls.WithCert(cfg.TLS.Cert),
			tls.WithKey(cfg.TLS.Key),
			tls.WithCa(cfg.TLS.CA),
		)
		if err != nil {
			return nil, err
		}

		opts = append(
			opts,
			cassandra.WithTLS(tcfg),
			cassandra.WithTLSCertPath(cfg.TLS.Cert),
			cassandra.WithTLSKeyPath(cfg.TLS.Key),
			cassandra.WithTLSCAPath(cfg.TLS.CA),
		)
	}

	db, err := cassandra.New(opts...)
	if err != nil {
		return nil, err
	}
	return &client{db: db}, nil
}

func (c *client) Connect(ctx context.Context) error {
	return c.db.Open(ctx)
}

func (c *client) Close(ctx context.Context) error {
	return c.db.Close(ctx)
}

func (c *client) getMetaVector(ctx context.Context, uuid string) (*model.MetaVector, error) {
	var metaVector model.MetaVector

	stmt, names := c.db.Select(c.metaTable, c.metaColumns, c.db.Eq("uuid"))
	q := c.db.Query(stmt, names).BindMap(map[string]interface{}{"uuid": uuid})

	if err := q.GetRelease(&metaVector); err != nil {
		switch err {
		case cassandra.ErrNotFound:
			return nil, errors.ErrCassandraNotFound(uuid)
		default:
			return nil, err
		}
	}

	return &metaVector, nil
}

func (c *client) GetMeta(ctx context.Context, uuid string) (*model.MetaVector, error) {
	return c.getMetaVector(ctx, uuid)
}

func (c *client) GetIPs(ctx context.Context, uuid string) ([]string, error) {
	mv, err := c.getMetaVector(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return mv.GetIPs(), nil
}

func (c *client) setMetaVectors(ctx context.Context, mvs ...model.MetaVector) error {
	ib := c.db.Insert(c.metaTable, c.metaColumns...)
	bt := c.db.Batch()

	entities := make([]interface{}, 0, len(mvs))
	for _, mv := range mvs {
		bt = bt.Add(ib)
		entities = append(entities, mv)
	}

	stmt, names := bt.ToCql()
	return c.db.Query(stmt, names).ExecRelease()
}

func (c *client) SetMeta(ctx context.Context, meta model.MetaVector) error {
	return c.setMetaVectors(ctx, meta)
}

func (c *client) SetMetas(ctx context.Context, metas ...model.MetaVector) error {
	return c.setMetaVectors(ctx, metas...)
}

func (c *client) deleteMetaVectors(ctx context.Context, uuids ...string) error {
	deleteBuilder := c.db.Delete(c.metaTable, c.db.Eq("uuid"))
	bt := c.db.Batch()
	bindUUIDs := make([]interface{}, 0, len(uuids))
	for _, uuid := range uuids {
		bt.Add(deleteBuilder)
		bindUUIDs = append(bindUUIDs, uuid)
	}

	stmt, names := bt.ToCql()
	return c.db.Query(stmt, names).Bind(bindUUIDs...).ExecRelease()
}

func (c *client) DeleteMeta(ctx context.Context, uuid string) error {
	return c.deleteMetaVectors(ctx, uuid)
}

func (c *client) DeleteMetas(ctx context.Context, uuids ...string) error {
	return c.deleteMetaVectors(ctx, uuids...)
}

func (c *client) SetIPs(ctx context.Context, uuid string, ips ...string) error {
	return nil
}

func (c *client) RemoveIPs(ctx context.Context, ips ...string) error {
	return nil
}
