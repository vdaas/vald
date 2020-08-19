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

package service

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/db/rdb/mysql"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/model"
)

type MySQL interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	GetMeta(ctx context.Context, uuid string) (*model.MetaVector, error)
	GetIPs(ctx context.Context, uuid string) ([]string, error)
	SetMeta(ctx context.Context, meta *model.MetaVector) error
	SetMetas(ctx context.Context, metas ...*model.MetaVector) error
	DeleteMeta(ctx context.Context, uuid string) error
	DeleteMetas(ctx context.Context, uuids ...string) error
	SetIPs(ctx context.Context, uuid string, ips ...string) error
	RemoveIPs(ctx context.Context, ips ...string) error
}

type client struct {
	db mysql.MySQL
}

func New(opts ...Option) (ms MySQL, err error) {
	c := new(client)
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(c); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return c, nil
}

func (c *client) Connect(ctx context.Context) error {
	return c.db.Open(ctx)
}

func (c *client) Close(ctx context.Context) error {
	return c.db.Close(ctx)
}

func (c *client) GetMeta(ctx context.Context, uuid string) (*model.MetaVector, error) {
	res, err := c.db.GetMeta(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &model.MetaVector{
		UUID:   res.GetUUID(),
		Vector: res.GetVector(),
		Meta:   res.GetMeta(),
		IPs:    res.GetIPs(),
	}, err
}

func (c *client) GetIPs(ctx context.Context, uuid string) ([]string, error) {
	return c.db.GetIPs(ctx, uuid)
}

func (c *client) SetMeta(ctx context.Context, meta *model.MetaVector) error {
	return c.db.SetMeta(ctx, meta)
}

func (c *client) SetMetas(ctx context.Context, metas ...*model.MetaVector) error {
	ms := make([]mysql.MetaVector, 0, len(metas))
	for _, meta := range metas {
		m := meta
		ms = append(ms, m)
	}
	return c.db.SetMetas(ctx, ms...)
}

func (c *client) DeleteMeta(ctx context.Context, uuid string) error {
	return c.db.DeleteMeta(ctx, uuid)
}

func (c *client) DeleteMetas(ctx context.Context, uuids ...string) error {
	return c.db.DeleteMetas(ctx, uuids...)
}

func (c *client) SetIPs(ctx context.Context, uuid string, ips ...string) error {
	return c.db.SetIPs(ctx, uuid, ips...)
}

func (c *client) RemoveIPs(ctx context.Context, ips ...string) error {
	return c.db.RemoveIPs(ctx, ips...)
}
