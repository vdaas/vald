//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	GetVector(ctx context.Context, uuid string) (*model.Vector, error)
	GetIPs(ctx context.Context, uuid string) ([]string, error)
	SetVector(ctx context.Context, vector *model.Vector) error
	SetVectors(ctx context.Context, vectors ...*model.Vector) error
	DeleteVector(ctx context.Context, uuid string) error
	DeleteVectors(ctx context.Context, uuids ...string) error
	SetIPs(ctx context.Context, uuid string, ips ...string) error
	RemoveIPs(ctx context.Context, ips ...string) error
}

type client struct {
	db mysql.MySQL
}

func New(opts ...Option) (ms MySQL, err error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
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

func (c *client) GetVector(ctx context.Context, uuid string) (*model.Vector, error) {
	res, err := c.db.GetVector(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &model.Vector{
		UUID:   res.GetUUID(),
		Vector: res.GetVector(),
		IPs:    res.GetIPs(),
	}, err
}

func (c *client) GetIPs(ctx context.Context, uuid string) ([]string, error) {
	return c.db.GetIPs(ctx, uuid)
}

func (c *client) SetVector(ctx context.Context, vector *model.Vector) error {
	return c.db.SetVector(ctx, vector)
}

func (c *client) SetVectors(ctx context.Context, vectors ...*model.Vector) error {
	ms := make([]mysql.Vector, 0, len(vectors))
	for _, vector := range vectors {
		m := vector
		ms = append(ms, m)
	}
	return c.db.SetVectors(ctx, ms...)
}

func (c *client) DeleteVector(ctx context.Context, uuid string) error {
	return c.db.DeleteVector(ctx, uuid)
}

func (c *client) DeleteVectors(ctx context.Context, uuids ...string) error {
	return c.db.DeleteVectors(ctx, uuids...)
}

func (c *client) SetIPs(ctx context.Context, uuid string, ips ...string) error {
	return c.db.SetIPs(ctx, uuid, ips...)
}

func (c *client) RemoveIPs(ctx context.Context, ips ...string) error {
	return c.db.RemoveIPs(ctx, ips...)
}
