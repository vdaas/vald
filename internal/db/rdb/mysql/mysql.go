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

package mysql

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/go-sql-driver/mysql"
	dbr "github.com/gocraft/dbr/v2"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
)

const (
	metaVectorTableName = "meta_vector"
	podIPTableName      = "pod_ip"
	uuidColumnName      = "uuid"
	ipColumnName        = "ip"
	asterisk            = "*"
)

type MySQL interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	Getter
	Setter
}

type mySQLClient struct {
	db                   string
	host                 string
	port                 int
	user                 string
	pass                 string
	name                 string
	initialPingTimeLimit time.Duration
	initialPingDuration  time.Duration
	connMaxLifeTime      time.Duration
	dialer               func(ctx context.Context, network, addr string) (net.Conn, error)
	tlsConfig            *tls.Config
	maxOpenConns         int
	maxIdleConns         int
	session              *dbr.Session
	connected            atomic.Value
}

func New(opts ...Option) (MySQL, error) {
	m := new(mySQLClient)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(m); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return m, nil
}

func (m *mySQLClient) Open(ctx context.Context) error {
	var addParam string
	network := "tcp"
	if m.dialer != nil {
		mysql.RegisterDialContext(network, func(ctx context.Context, addr string) (net.Conn, error) {
			return m.dialer(ctx, network, addr)
		})
	}
	if m.tlsConfig != nil {
		tlsConfName := "tls"
		mysql.RegisterTLSConfig(tlsConfName, m.tlsConfig)
		addParam += "&tls=" + tlsConfName
	}

	conn, err := dbr.Open(m.db,
		fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local%s",
			m.user, m.pass, network, m.host, m.port, m.name, addParam), nil)
	if err != nil {
		return err
	}
	conn.SetConnMaxLifetime(m.connMaxLifeTime)
	conn.SetMaxIdleConns(m.maxIdleConns)
	conn.SetMaxOpenConns(m.maxOpenConns)

	m.session = conn.NewSession(nil)
	m.connected.Store(true)

	return m.Ping(ctx)
}

func (m *mySQLClient) Ping(ctx context.Context) (err error) {
	pctx, cancel := context.WithTimeout(ctx, m.initialPingTimeLimit)
	defer cancel()
	tick := time.NewTicker(m.initialPingDuration)
	for {
		select {
		case <-pctx.Done():
			return errors.Wrap(errors.Wrap(err, errors.ErrMySQLConnectionPingFailed.Error()), ctx.Err().Error())
		case <-tick.C:
			err = m.session.PingContext(ctx)
			if err == nil {
				return nil
			}
			log.Error(err)
		}
	}
}

func (m *mySQLClient) Close(ctx context.Context) error {
	if m.connected.Load().(bool) {
		m.session.Close()
		m.connected.Store(false)
	}
	return nil
}

func (m *mySQLClient) GetMeta(ctx context.Context, uuid string) (MetaVector, error) {
	if !m.connected.Load().(bool) {
		return nil, errors.ErrMySQLConnectionClosed
	}

	var meta *meta
	_, err := m.session.Select(asterisk).From(metaVectorTableName).Where(dbr.Eq(uuidColumnName, uuid)).Limit(1).LoadContext(ctx, &meta)
	if err != nil {
		return nil, err
	}
	if meta == nil {
		return nil, errors.ErrRequiredElementNotFoundByUUID(uuid)
	}

	var podIPs []podIP
	_, err = m.session.Select(asterisk).From(podIPTableName).Where(dbr.Eq(uuidColumnName, uuid)).LoadContext(ctx, &podIPs)
	if err != nil {
		return nil, err
	}

	return &metaVector{
		meta:   *meta,
		podIPs: podIPs,
	}, nil
}

func (m *mySQLClient) GetIPs(ctx context.Context, uuid string) ([]string, error) {
	if !m.connected.Load().(bool) {
		return nil, errors.ErrMySQLConnectionClosed
	}

	var podIPs []podIP
	_, err := m.session.Select(asterisk).From(podIPTableName).Where(dbr.Eq(uuidColumnName, uuid)).LoadContext(ctx, &podIPs)
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0, len(podIPs))
	for _, podIP := range podIPs {
		ips = append(ips, podIP.IP)
	}

	return ips, nil
}

func validateMeta(meta MetaVector) error {
	if meta.GetVectorString() == "" {
		return errors.ErrRequiredMemberNotFilled("vector")
	}
	return nil
}

func setMetaWithTx(ctx context.Context, tx *dbr.Tx, meta MetaVector) error {
	err := validateMeta(meta)
	if err != nil {
		return err
	}

	_, err = tx.InsertBySql("INSERT INTO meta_vector(uuid, vector, meta) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE vector = ?, meta = ?",
		meta.GetUUID(),
		meta.GetVectorString(),
		meta.GetMeta(),
		meta.GetVectorString(),
		meta.GetMeta()).ExecContext(ctx)
	if err != nil {
		return err
	}

	_, err = tx.DeleteFrom(podIPTableName).Where(dbr.Eq(uuidColumnName, meta.GetUUID())).ExecContext(ctx)
	if err != nil {
		return err
	}

	stmt := tx.InsertInto(podIPTableName).Columns(uuidColumnName, ipColumnName)
	for _, ip := range meta.GetIPs() {
		stmt.Record(&podIP{UUID: meta.GetUUID(), IP: ip})
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m *mySQLClient) SetMeta(ctx context.Context, meta MetaVector) error {
	if !m.connected.Load().(bool) {
		return errors.ErrMySQLConnectionClosed
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = setMetaWithTx(ctx, tx, meta)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *mySQLClient) SetMetas(ctx context.Context, metas ...MetaVector) error {
	if !m.connected.Load().(bool) {
		return errors.ErrMySQLConnectionClosed
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	for _, meta := range metas {
		err = setMetaWithTx(ctx, tx, meta)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func deleteMetaWithTx(ctx context.Context, tx *dbr.Tx, uuid string) error {
	_, err := tx.DeleteFrom(metaVectorTableName).Where(dbr.Eq(uuidColumnName, uuid)).ExecContext(ctx)
	if err != nil {
		return err
	}

	_, err = tx.DeleteFrom(podIPTableName).Where(dbr.Eq(uuidColumnName, uuid)).ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (m *mySQLClient) DeleteMeta(ctx context.Context, uuid string) error {
	if !m.connected.Load().(bool) {
		return errors.ErrMySQLConnectionClosed
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = deleteMetaWithTx(ctx, tx, uuid)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *mySQLClient) DeleteMetas(ctx context.Context, uuids ...string) error {
	if !m.connected.Load().(bool) {
		return errors.ErrMySQLConnectionClosed
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	for _, uuid := range uuids {
		err = deleteMetaWithTx(ctx, tx, uuid)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (m *mySQLClient) SetIPs(ctx context.Context, uuid string, ips ...string) error {
	if !m.connected.Load().(bool) {
		return errors.ErrMySQLConnectionClosed
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	stmt := tx.InsertInto(podIPTableName).Columns(uuidColumnName, ipColumnName)
	for _, ip := range ips {
		stmt.Record(&podIP{UUID: uuid, IP: ip})
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *mySQLClient) RemoveIPs(ctx context.Context, ips ...string) error {
	if !m.connected.Load().(bool) {
		return errors.ErrMySQLConnectionClosed
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	for _, ip := range ips {
		_, err = tx.DeleteFrom(podIPTableName).Where(dbr.Eq(ipColumnName, ip)).ExecContext(ctx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
