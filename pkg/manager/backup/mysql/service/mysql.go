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

	"github.com/kpango/gache"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/db/rdb/mysql"
	"github.com/vdaas/vald/internal/net/tcp"
	"github.com/vdaas/vald/internal/tls"
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
	db  mysql.MySQL
	der tcp.Dialer
}

func New(cfg *config.MySQL) (MySQL, error) {
	c := new(client)

	opts := append(make([]mysql.Option, 0, 13),
		mysql.WithDB(cfg.DB),
		mysql.WithHost(cfg.Host),
		mysql.WithPort(cfg.Port),
		mysql.WithUser(cfg.User),
		mysql.WithPass(cfg.Pass),
		mysql.WithName(cfg.Name),
		mysql.WithCharset(cfg.Charset),
		mysql.WithTimezone(cfg.Timezone),
		mysql.WithInitialPingTimeLimit(cfg.InitialPingTimeLimit),
		mysql.WithInitialPingDuration(cfg.InitialPingDuration),
		mysql.WithConnectionLifeTimeLimit(cfg.ConnMaxLifeTime),
		mysql.WithMaxIdleConns(cfg.MaxIdleConns),
		mysql.WithMaxOpenConns(cfg.MaxOpenConns),
	)

	if cfg.TLS != nil && cfg.TLS.Enabled {
		tcfg, err := tls.New(
			tls.WithCert(cfg.TLS.Cert),
			tls.WithKey(cfg.TLS.Key),
			tls.WithCa(cfg.TLS.CA),
		)
		if err != nil {
			return nil, err
		}
		opts = append(opts, mysql.WithTLSConfig(tcfg))
	}

	if cfg.TCP != nil {
		topts := make([]tcp.DialerOption, 0, 8)
		if cfg.TCP.DNS != nil && cfg.TCP.DNS.CacheEnabled {
			topts = append(topts,
				tcp.WithCache(gache.New()),
				tcp.WithEnableDNSCache(),
				tcp.WithDNSCacheExpiration(cfg.TCP.DNS.CacheExpiration),
				tcp.WithDNSRefreshDuration(cfg.TCP.DNS.RefreshDuration),
			)
		}

		if cfg.TCP.Dialer != nil && cfg.TCP.Dialer.DualStackEnabled {
			topts = append(topts, tcp.WithEnableDialerDualStack())
		} else {
			topts = append(topts, tcp.WithDisableDialerDualStack())
		}

		if cfg.TCP.TLS != nil && cfg.TCP.TLS.Enabled {
			tcfg, err := tls.New(
				tls.WithCert(cfg.TCP.TLS.Cert),
				tls.WithKey(cfg.TCP.TLS.Key),
				tls.WithCa(cfg.TCP.TLS.CA),
			)
			if err != nil {
				return nil, err
			}
			topts = append(topts, tcp.WithTLS(tcfg))
		}
		c.der = tcp.NewDialer(append(topts,
			tcp.WithDialerKeepAlive(cfg.TCP.Dialer.KeepAlive),
			tcp.WithDialerTimeout(cfg.TCP.Dialer.Timeout),
		)...)
		opts = append(opts, mysql.WithDialer(c.der.GetDialer()))
	}

	m, err := mysql.New(opts...)

	if err != nil {
		return nil, err
	}

	c.db = m

	return c, nil
}

func (c *client) Connect(ctx context.Context) error {
	if c.der != nil {
		c.der.StartDialerCache(ctx)
	}
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
