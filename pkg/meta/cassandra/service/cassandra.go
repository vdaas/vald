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

// Package service manages the main logic of server.
package service

import (
	"context"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/nosql/cassandra"
	"github.com/vdaas/vald/internal/errors"
)

type Cassandra interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
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
	db cassandra.Cassandra
}

func New(cfg *config.Cassandra) (Cassandra, error) {
	opts := []cassandra.Option{
		cassandra.WithHosts(cfg.Hosts...),
		cassandra.WithCQLVersion(cfg.CQLVersion),
		cassandra.WithTimeout(cfg.Timeout),
		cassandra.WithConnectTimeout(cfg.ConnectTimeout),
		cassandra.WithPort(cfg.Port),
		cassandra.WithNumConns(cfg.NumConns),
		cassandra.WithConsistency(cfg.Consistency),
		cassandra.WithMaxPreparedStmts(cfg.MaxPreparedStmts),
		cassandra.WithMaxRoutingKeyInfo(cfg.MaxRoutingKeyInfo),
		cassandra.WithPageSize(cfg.PageSize),
		cassandra.WithDefaultTimestamp(cfg.DefaultTimestamp),
		cassandra.WithMaxWaitSchemaAgreement(cfg.MaxWaitSchemaAgreement),
		cassandra.WithReconnectInterval(cfg.ReconnectInterval),
		cassandra.WithReconnectionPolicyMaxRetries(cfg.ReconnectionPolicy.MaxRetries),
		cassandra.WithReconnectionPolicyInitialInterval(cfg.ReconnectionPolicy.InitialInterval),
		cassandra.WithWriteCoalesceWaitTime(cfg.WriteCoalesceWaitTime),
		cassandra.WithKeyspace(cfg.Keyspace),
		cassandra.WithKVTable(cfg.KVTable),
		cassandra.WithVKTable(cfg.VKTable),
		cassandra.WithUsername(cfg.Username),
		cassandra.WithPassword(cfg.Password),
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

func (c *client) Disconnect(ctx context.Context) error {
	return c.db.Close(ctx)
}

func (c *client) Get(key string) (string, error) {
	return c.db.GetValue(key)
}

func (c *client) GetMultiple(keys ...string) (vals []string, err error) {
	return c.db.MultiGetValue(keys...)
}

func (c *client) GetInverse(val string) (string, error) {
	return c.db.GetKey(val)
}

func (c *client) GetInverseMultiple(vals ...string) ([]string, error) {
	return c.db.MultiGetKey(vals...)
}

func (c *client) Set(key, val string) error {
	return c.db.Set(key, val)
}

func (c *client) SetMultiple(kvs map[string]string) (err error) {
	return c.db.MultiSet(kvs)
}

func (c *client) Delete(key string) (string, error) {
	vals, err := c.db.Delete(key)
	if err != nil {
		return "", err
	}

	if len(vals) != 1 {
		return "", errors.ErrCassandraDeleteOperationFailed(key, nil)
	}

	return vals[0], nil
}

func (c *client) DeleteMultiple(keys ...string) ([]string, error) {
	return c.db.Delete(keys...)
}

func (c *client) DeleteInverse(val string) (string, error) {
	keys, err := c.db.DeleteByValues(val)
	if err != nil {
		return "", err
	}

	if len(keys) != 1 {
		return "", errors.ErrCassandraDeleteOperationFailed(val, nil)
	}

	return keys[0], nil
}

func (c *client) DeleteInverseMultiple(vals ...string) ([]string, error) {
	return c.db.DeleteByValues(vals...)
}
