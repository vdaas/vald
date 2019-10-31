//
// Copyright (C) 2019 kpango (Yusuke Kato)
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
	"fmt"
	"reflect"

	dbr "github.com/gocraft/dbr/v2"
	"github.com/vdaas/vald/internal/errors"
)

type MySQL interface {
	Open() error
	Close() error
	Getter
	Setter
}

type mySQLClient struct {
	db        string
	host      string
	port      int
	user      string
	pass      string
	name      string
	session   *dbr.Session
	connected bool
}

func New(ctx context.Context, opts ...Option) (MySQL, error) {
	m := new(mySQLClient)
	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	err := m.Open()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *mySQLClient) Open() error {
	conn, err := dbr.Open(m.db,
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
			m.user, m.pass, m.host, m.port, m.name), nil)
	if err != nil {
		return err
	}

	m.session = conn.NewSession(nil)
	m.connected = true

	return nil
}

func (m *mySQLClient) Close() error {
	if m.connected {
		m.session.Close()
		m.connected = false
	}
	return nil
}

func (m *mySQLClient) GetMeta(uuid string) (MetaVector, error) {
	if !m.connected {
		return nil, errors.ErrMySQLConnectionClosed
	}

	var metas []meta
	_, err := m.session.Select("*").From("meta_vector").Where(dbr.Eq("uuid", uuid)).Load(&metas)
	if err != nil {
		return nil, err
	}

	if len(metas) > 0 {
		return nil, errors.ErrRequiredElementNotFoundByUUID(uuid)
	}

	var podIPs []podIP
	_, err = m.session.Select("*").From("pod_ip").Where(dbr.Eq("uuid", uuid)).Load(&podIPs)
	if err != nil {
		return nil, err
	}

	return &metaVector{
		meta:   metas[0],
		podIPs: podIPs,
	}, nil
}

func (m *mySQLClient) GetIPs(uuid string) ([]string, error) {
	if !m.connected {
		return nil, errors.ErrMySQLConnectionClosed
	}

	var podIPs []podIP
	_, err := m.session.Select("*").From("pod_ip").Where(dbr.Eq("uuid", uuid)).Load(&podIPs)
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0, len(podIPs))
	for _, podIP := range podIPs {
		ips = append(ips, podIP.IP)
	}

	return ips, nil
}

func setMetaWithTx(tx *dbr.Tx, meta MetaVector) error {
	_, err := tx.InsertBySql("INSERT INTO meta_vector(uuid, vector, meta) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE vector = ?, meta = ?",
		meta.GetUUID(), meta.GetVector(), meta.GetMeta(), meta.GetVector(), meta.GetMeta()).Exec()
	if err != nil {
		return err
	}

	_, err = tx.DeleteFrom("pod_ip").Where(dbr.Eq("uuid", meta.GetUUID())).Exec()
	if err != nil {
		return err
	}

	stmt := tx.InsertInto("pod_ip").Columns("uuid", "ip")
	for _, ip := range meta.GetIPs() {
		stmt.Record(ip)
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (m *mySQLClient) SetMeta(meta MetaVector) error {
	if !m.connected {
		return errors.ErrMySQLConnectionClosed
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = setMetaWithTx(tx, meta)

	return tx.Commit()
}

func (m *mySQLClient) SetMetas(metas ...MetaVector) error {
	if !m.connected {
		return errors.ErrMySQLConnectionClosed
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	for _, meta := range metas {
		err = setMetaWithTx(tx, meta)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func deleteMetaWithTx(tx *dbr.Tx, uuid string) error {
	_, err := tx.DeleteFrom("meta_vector").Where(dbr.Eq("uuid", uuid)).Exec()
	if err != nil {
		return err
	}

	_, err = tx.DeleteFrom("pod_ip").Where(dbr.Eq("uuid", uuid)).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (m *mySQLClient) DeleteMeta(uuid string) error {
	if !m.connected {
		return errors.ErrMySQLConnectionClosed
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	err = deleteMetaWithTx(tx, uuid)

	return tx.Commit()
}

func (m *mySQLClient) DeleteMetas(uuids ...string) error {
	if !m.connected {
		return errors.ErrMySQLConnectionClosed
	}

	tx, err := m.session.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	for _, uuid := range uuids {
		err = deleteMetaWithTx(tx, uuid)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
