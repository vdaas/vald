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

package service

import (
	"context"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/rdb/mysql"
	"github.com/vdaas/vald/pkg/manager/backup/model"
)

type Mysql interface {
	Connect(ctx context.Context) error
	Close() error
	GetMeta(uuid string) (*model.MetaVector, error)
	GetIPs(uuid string) ([]string, error)
	SetMeta(meta model.MetaVector) error
	SetMetas(metas ...model.MetaVector) error
	DeleteMeta(uuid string) error
	DeleteMetas(uuids ...string) error
}

type client struct {
	db   mysql.Mysql
	opts []mysql.Option
}

func NewMysql(cfg *config.Mysql) (Mysql, error) {
	c := &client{}

	c.opts = make([]mysql.Option, 0, 6)
	c.opts = append(c.opts,
		mysql.WithDB(cfg.DB),
		mysql.WithHost(cfg.Host),
		mysql.WithPort(cfg.Port),
		mysql.WithUser(cfg.User),
		mysql.WithPass(cfg.Pass),
		mysql.WithName(cfg.Name))

	return c, nil
}

func (c *client) Connect(ctx context.Context) error {
	m, err := mysql.New(ctx, c.opts...)
	if err != nil {
		return err
	}

	c.db = m

	return nil
}

func (c *client) Close() error {
	return c.db.Close()
}

func (c *client) GetMeta(uuid string) (*model.MetaVector, error) {
	res, err := c.db.GetMeta(uuid)
	if err != nil {
		return nil, err
	}

	return &model.MetaVector{
		UUID:     res.GetUUID(),
		ObjectID: res.GetObjectID(),
		Vector:   res.GetVector(),
		Meta:     res.GetMeta(),
		IPs:      res.GetIPs(),
	}, err
}

func (c *client) GetIPs(uuid string) ([]string, error) {
	return c.db.GetIPs(uuid)
}

func (c *client) SetMeta(meta model.MetaVector) error {
	return c.db.SetMeta(&meta)
}

func (c *client) SetMetas(metas ...model.MetaVector) error {
	ms := make([]mysql.MetaVector, 0, len(metas))
	for _, meta := range metas {
		ms = append(ms, &meta)
	}
	return c.db.SetMetas(ms...)
}

func (c *client) DeleteMeta(uuid string) error {
	return c.db.DeleteMeta(uuid)
}

func (c *client) DeleteMetas(uuids ...string) error {
	return c.db.DeleteMetas(uuids...)
}
