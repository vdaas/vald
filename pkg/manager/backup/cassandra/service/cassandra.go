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
	"strconv"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/tls"
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
	SetMeta(ctx context.Context, meta model.MetaVector) error
	SetMetas(ctx context.Context, metas ...model.MetaVector) error
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

	stmt, names := cassandra.Select(c.metaTable, metaColumns, cassandra.Eq(uuidColumn))
	err := c.db.Query(stmt, names).BindMap(map[string]interface{}{uuidColumn: uuid}).GetRelease(&metaVector)

	if err != nil {
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

	return mv.IPs, nil
}

func (c *client) SetMeta(ctx context.Context, meta model.MetaVector) error {
	stmt, names := cassandra.Insert(c.metaTable, metaColumns...).ToCql()
	return c.db.Query(stmt, names).BindStruct(meta).ExecRelease()
}

func (c *client) SetMetas(ctx context.Context, metas ...model.MetaVector) error {
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

	stmt, names := bt.ToCql()
	return c.db.Query(stmt, names).BindMap(entities).ExecRelease()
}

func (c *client) DeleteMeta(ctx context.Context, uuid string) error {
	stmt, names := cassandra.Delete(c.metaTable, cassandra.Eq(uuidColumn)).ToCql()
	return c.db.Query(stmt, names).BindMap(map[string]interface{}{uuidColumn: uuid}).ExecRelease()
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

	stmt, names := bt.ToCql()
	return c.db.Query(stmt, names).BindMap(bindUUIDs).ExecRelease()
}

func (c *client) SetIPs(ctx context.Context, uuid string, ips ...string) error {
	stmt, names := cassandra.Update(c.metaTable).AddNamed(ipsColumn, ipsColumn).Where(cassandra.Eq(uuidColumn)).ToCql()
	return c.db.Query(stmt, names).BindMap(map[string]interface{}{uuidColumn: uuid, ipsColumn: ips}).ExecRelease()
}

func (c *client) RemoveIPs(ctx context.Context, ips ...string) error {
	var metaVectors []model.MetaVector

	for _, ip := range ips {
		stmt, names := cassandra.Select(c.metaTable, []string{uuidColumn, ipsColumn}, cassandra.Contains(ipsColumn))
		err := c.db.Query(stmt, names).BindMap(map[string]interface{}{ipsColumn: ip}).SelectRelease(&metaVectors)
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

			stmt, names = cassandra.Update(c.metaTable).Set(ipsColumn).Where(cassandra.Eq(uuidColumn)).ToCql()
			err = c.db.Query(stmt, names).BindMap(map[string]interface{}{uuidColumn: mv.UUID, ipsColumn: newIPs}).ExecRelease()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
