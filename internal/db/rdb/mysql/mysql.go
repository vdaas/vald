//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

package mysql

import (
	"context"
	"crypto/tls"
	"fmt"
	"reflect"
	"sync/atomic"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/vdaas/vald/internal/db/rdb/mysql/dbr"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
)

const (
	vectorTableName  = "backup_vector"
	podIPTableName   = "pod_ip"
	idColumnName     = "id"
	uuidColumnName   = "uuid"
	vectorColumnName = "vector"
	ipColumnName     = "ip"
	asterisk         = "*"
)

// MySQL represents the interface to handle MySQL operation.
type MySQL interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	Getter
	Setter
}

type mySQLClient struct {
	db                   string
	network              string
	socketPath           string
	host                 string
	port                 uint16
	user                 string
	pass                 string
	name                 string
	charset              string
	timezone             string
	initialPingTimeLimit time.Duration
	initialPingDuration  time.Duration
	connMaxLifeTime      time.Duration
	dialer               net.Dialer
	dialerFunc           func(ctx context.Context, network, addr string) (net.Conn, error)
	tlsConfig            *tls.Config
	maxOpenConns         int
	maxIdleConns         int
	session              dbr.Session
	connected            atomic.Value
	eventReceiver        EventReceiver
	dbr                  dbr.DBR
}

// New creates the new mySQLClient with option.
// It will return error when set option is failed.
func New(opts ...Option) (MySQL, error) {
	m := &mySQLClient{
		dbr: dbr.New(),
	}
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(m); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return m, nil
}

// Open opens the connection with MySQL.
// It will return error when connecting to MySQL ends with fail.
func (m *mySQLClient) Open(ctx context.Context) (err error) {
	if m.dialer != nil {
		m.dialer.StartDialerCache(ctx)
		m.dialerFunc = m.dialer.GetDialer()
	}

	var addParam string
	if m.dialerFunc != nil {
		mysql.RegisterDialContext(m.network, func(ctx context.Context, addr string) (net.Conn, error) {
			return m.dialerFunc(ctx, m.network, addr)
		})
	}

	if m.tlsConfig != nil {
		tlsConfName := "tls"
		err = mysql.RegisterTLSConfig(tlsConfName, m.tlsConfig)
		if err != nil {
			return err
		}
		addParam += "&tls=" + tlsConfName
	}

	var addr, network string
	nt := net.NetworkTypeFromString(m.network)
	if len(m.socketPath) != 0 && (nt == net.UNIX ||
		nt == net.UNIXGRAM ||
		nt == net.UNIXPACKET) {
		network = net.UNIX.String()
		addr = m.socketPath
	} else {
		network = net.TCP.String()
		addr = net.JoinHostPort(m.host, m.port)
	}
	conn, err := m.dbr.Open(m.db, fmt.Sprintf(
		"%s:%s@%s(%s)/%s?charset=%s&parseTime=true&loc=%s%s",
		m.user, m.pass, network, addr, m.name, m.charset, m.timezone, addParam),
		m.eventReceiver)
	if err != nil {
		return err
	}

	conn.SetConnMaxLifetime(m.connMaxLifeTime)
	conn.SetMaxIdleConns(m.maxIdleConns)
	conn.SetMaxOpenConns(m.maxOpenConns)

	if m.session == nil {
		m.session = dbr.NewSession(conn, m.eventReceiver)
	}
	m.connected.Store(true)

	return m.Ping(ctx)
}

// Ping check the connection of MySQL database.
// If the connection is closed, it returns error.
func (m *mySQLClient) Ping(ctx context.Context) (err error) {
	if m.session == nil {
		err = errors.ErrMySQLSessionNil
		m.errorLog(err)
		return err
	}
	pctx, cancel := context.WithTimeout(ctx, m.initialPingTimeLimit)
	defer cancel()
	tick := time.NewTicker(m.initialPingDuration)
	for {
		select {
		case <-pctx.Done():
			if err != nil {
				err = errors.Wrap(errors.ErrMySQLConnectionPingFailed, err.Error())
			} else {
				err = errors.ErrMySQLConnectionPingFailed
			}
			cerr := pctx.Err()
			if cerr != nil {
				err = errors.Wrap(err, cerr.Error())
			}
			return err
		case <-tick.C:
			err = m.session.PingContext(pctx)
			if err == nil {
				return nil
			}
			log.Error(err)
		}
	}
}

// Close closes the connection of MySQL database.
// If the connection is already closed or closing connection is failed, it returns error.
func (m *mySQLClient) Close(context.Context) (err error) {
	if m.session == nil {
		err = errors.ErrMySQLSessionNil
		m.errorLog(err)
		m.connected.Store(false)
		return err
	}

	if m.connected.Load().(bool) {
		err = m.session.Close()
		if err == nil {
			m.connected.Store(false)
		}
	}
	return nil
}

// GetVector gets the vector data and podIPs which have index of vector.
func (m *mySQLClient) GetVector(ctx context.Context, uuid string) (Vector, error) {
	if !m.connected.Load().(bool) {
		return nil, errors.ErrMySQLConnectionClosed
	}

	if m.session == nil {
		err := errors.ErrMySQLSessionNil
		m.errorLog(err)
		return nil, err
	}

	var data *data
	_, err := m.session.Select(asterisk).From(vectorTableName).Where(m.dbr.Eq(uuidColumnName, uuid)).Limit(1).LoadContext(ctx, &data)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.ErrRequiredElementNotFoundByUUID(uuid)
	}

	var podIPs []podIP
	_, err = m.session.Select(asterisk).From(podIPTableName).Where(m.dbr.Eq(idColumnName, data.ID)).LoadContext(ctx, &podIPs)
	if err != nil {
		return nil, err
	}

	return &vector{
		data:   *data,
		podIPs: podIPs,
	}, nil
}

// GetIPs gets the pod ips which have index of requested uuids' vector data's vector.
func (m *mySQLClient) GetIPs(ctx context.Context, uuid string) ([]string, error) {
	if !m.connected.Load().(bool) {
		return nil, errors.ErrMySQLConnectionClosed
	}

	if m.session == nil {
		err := errors.ErrMySQLSessionNil
		m.errorLog(err)
		return nil, err
	}

	var id int64
	_, err := m.session.Select(idColumnName).From(vectorTableName).Where(m.dbr.Eq(uuidColumnName, uuid)).Limit(1).LoadContext(ctx, &id)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, errors.ErrRequiredElementNotFoundByUUID(uuid)
	}

	var podIPs []podIP
	_, err = m.session.Select(asterisk).From(podIPTableName).Where(m.dbr.Eq(idColumnName, id)).LoadContext(ctx, &podIPs)
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0, len(podIPs))
	for _, podIP := range podIPs {
		ips = append(ips, podIP.IP)
	}

	return ips, nil
}

func validateVector(vec Vector) error {
	if len(vec.GetVector()) == 0 {
		return errors.ErrRequiredMemberNotFilled("vector")
	}
	return nil
}

// SetVector records vector data at backup_vector table and set of (podIP, uuid) at podIPtable through same transaction.
// If error occurs it will rollback by defer function.
func (m *mySQLClient) SetVector(ctx context.Context, vec Vector) error {
	if !m.connected.Load().(bool) {
		return errors.ErrMySQLConnectionClosed
	}

	if m.session == nil {
		err := errors.ErrMySQLSessionNil
		m.errorLog(err)
		return err
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = validateVector(vec)
	if err != nil {
		return err
	}

	_, err = tx.InsertBySql("INSERT INTO "+vectorTableName+"(uuid, vector) VALUES (?, ?) ON DUPLICATE KEY UPDATE vector = ?",
		vec.GetUUID(),
		vec.GetVector(),
		vec.GetVector()).ExecContext(ctx)
	if err != nil {
		return err
	}

	var id int64
	_, err = tx.Select(idColumnName).From(vectorTableName).Where(m.dbr.Eq(uuidColumnName, vec.GetUUID())).Limit(1).LoadContext(ctx, &id)
	if err != nil {
		return err
	}
	if id == 0 {
		return errors.ErrRequiredElementNotFoundByUUID(vec.GetUUID())
	}

	_, err = tx.DeleteFrom(podIPTableName).Where(m.dbr.Eq(idColumnName, id)).ExecContext(ctx)
	if err != nil {
		return err
	}

	stmt := tx.InsertInto(podIPTableName).Columns(idColumnName, ipColumnName)
	for _, ip := range vec.GetIPs() {
		stmt.Record(&podIP{ID: id, IP: ip})
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// SetVectors records multiple vector data like as SetVector().
func (m *mySQLClient) SetVectors(ctx context.Context, vecs ...Vector) error {
	if !m.connected.Load().(bool) {
		return errors.ErrMySQLConnectionClosed
	}

	if m.session == nil {
		err := errors.ErrMySQLSessionNil
		m.errorLog(err)
		return err
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	for _, vec := range vecs {
		err = validateVector(vec)
		if err != nil {
			return err
		}

		_, err = tx.InsertBySql("INSERT INTO "+vectorTableName+"(uuid, vector) VALUES (?, ?) ON DUPLICATE KEY UPDATE vector = ?",
			vec.GetUUID(),
			vec.GetVector(),
			vec.GetVector()).ExecContext(ctx)
		if err != nil {
			return err
		}
	}

	for _, vec := range vecs {
		var id int64
		_, err = tx.Select(idColumnName).From(vectorTableName).Where(m.dbr.Eq(uuidColumnName, vec.GetUUID())).Limit(1).LoadContext(ctx, &id)
		if err != nil {
			return err
		}
		if id == 0 {
			return errors.ErrRequiredElementNotFoundByUUID(vec.GetUUID())
		}

		_, err = tx.DeleteFrom(podIPTableName).Where(m.dbr.Eq(idColumnName, id)).ExecContext(ctx)
		if err != nil {
			return err
		}

		stmt := tx.InsertInto(podIPTableName).Columns(idColumnName, ipColumnName)
		for _, ip := range vec.GetIPs() {
			stmt.Record(&podIP{ID: id, IP: ip})
		}
		_, err = stmt.ExecContext(ctx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteVector deletes vector data from backup_vector table and podIPs from pod_ip table using vector's uuid.
func (m *mySQLClient) DeleteVector(ctx context.Context, val string) error {
	if !m.connected.Load().(bool) {
		return errors.ErrMySQLConnectionClosed
	}

	if m.session == nil {
		err := errors.ErrMySQLSessionNil
		m.errorLog(err)
		return err
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	if tx == nil {
		return errors.ErrMySQLTransactionNotCreated
	}
	defer tx.RollbackUnlessCommitted()

	var id int64
	_, err = tx.Select(idColumnName).From(vectorTableName).Where(m.dbr.Eq(uuidColumnName, val)).Limit(1).LoadContext(ctx, &id)
	if err != nil {
		return err
	}
	if id == 0 {
		return errors.ErrRequiredElementNotFoundByUUID(val)
	}

	_, err = tx.DeleteFrom(vectorTableName).Where(m.dbr.Eq(uuidColumnName, val)).ExecContext(ctx)
	if err != nil {
		return err
	}

	_, err = tx.DeleteFrom(podIPTableName).Where(m.dbr.Eq(idColumnName, id)).ExecContext(ctx)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// DeleteVectors is the same as DeleteVector() but it deletes multiple records.
func (m *mySQLClient) DeleteVectors(ctx context.Context, uuids ...string) (err error) {
	for _, uuid := range uuids {
		err = m.DeleteVector(ctx, uuid)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetIPs insert the vector's uuid and the podIPs into database.
func (m *mySQLClient) SetIPs(ctx context.Context, uuid string, ips ...string) error {
	if !m.connected.Load().(bool) {
		return errors.ErrMySQLConnectionClosed
	}

	if m.session == nil {
		err := errors.ErrMySQLSessionNil
		m.errorLog(err)
		return err
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	var id int64
	_, err = tx.Select(idColumnName).From(vectorTableName).Where(m.dbr.Eq(uuidColumnName, uuid)).Limit(1).LoadContext(ctx, &id)
	if err != nil {
		return err
	}
	if id == 0 {
		return errors.ErrRequiredElementNotFoundByUUID(uuid)
	}

	stmt := tx.InsertInto(podIPTableName).Columns(idColumnName, ipColumnName)
	for _, ip := range ips {
		stmt.Record(&podIP{ID: id, IP: ip})
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RemoveIPs delete the podIPs from database by podIPs.
func (m *mySQLClient) RemoveIPs(ctx context.Context, ips ...string) error {
	if !m.connected.Load().(bool) {
		return errors.ErrMySQLConnectionClosed
	}

	if m.session == nil {
		err := errors.ErrMySQLSessionNil
		m.errorLog(err)
		return err
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	_, err = tx.DeleteFrom(podIPTableName).Where(m.dbr.Eq(ipColumnName, ips)).ExecContext(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (m *mySQLClient) errorLog(err error) {
	log.Errorf(
		"err: %v, { host: %s, port: %d, user: %s, name: %s, db: %s, charset: %s, socketPath: %s, network: %s} ",
		err, m.host, m.port, m.user, m.name, m.db, m.charset, m.socketPath, m.network,
	)
}
