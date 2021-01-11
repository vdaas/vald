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

// Package compressor represents compressor client
package compressor

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/apis/grpc/v1/manager/compressor"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

type Client interface {
	Start(ctx context.Context) (<-chan error, error)
	Stop(context.Context) error
	GetVector(ctx context.Context, uuid string) (*payload.Backup_Vector, error)
	GetLocation(ctx context.Context, uuid string) ([]string, error)
	Register(ctx context.Context, vec *payload.Backup_Vector) error
	RegisterMultiple(ctx context.Context, vecs *payload.Backup_Vectors) error
	Remove(ctx context.Context, uuid string) error
	RemoveMultiple(ctx context.Context, uuids ...string) error
	RegisterIPs(ctx context.Context, ips []string) error
	RemoveIPs(ctx context.Context, ips []string) error
}

type client struct {
	client grpc.Client
}

func New(opts ...Option) (c Client, err error) {
	cc := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(cc); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}

	return cc, nil
}

func (c *client) Start(ctx context.Context) (<-chan error, error) {
	return c.client.StartConnectionMonitor(ctx)
}

func (c *client) Stop(ctx context.Context) error {
	return c.client.Close(ctx)
}

func (c *client) GetVector(ctx context.Context, uuid string) (vec *payload.Backup_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/internal/client/v1/client/compressor/Client.GetVector")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.client.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (i interface{}, err error) {
		vec, err = compressor.NewBackupClient(conn).GetVector(ctx, &payload.Backup_GetVector_Request{
			Uuid: uuid,
		}, copts...)
		if err != nil {
			return nil, err
		}
		return
	})
	return
}

func (c *client) GetLocation(ctx context.Context, uuid string) (ipList []string, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/internal/client/v1/client/compressor/Client.GetLocation")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.client.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (i interface{}, err error) {
		ips, err := compressor.NewBackupClient(conn).Locations(ctx, &payload.Backup_Locations_Request{
			Uuid: uuid,
		}, copts...)
		if err != nil {
			return nil, err
		}
		ipList = ips.GetIp()
		return
	})
	return
}

func (c *client) Register(ctx context.Context, vec *payload.Backup_Vector) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/internal/client/v1/client/compressor/Client.Register")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.client.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (i interface{}, err error) {
		_, err = compressor.NewBackupClient(conn).Register(ctx, vec, copts...)
		if err != nil {
			return nil, err
		}
		return
	})
	return
}

func (c *client) RegisterMultiple(ctx context.Context, vecs *payload.Backup_Vectors) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/internal/client/v1/client/compressor/Client.RegisterMultiple")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.client.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (i interface{}, err error) {
		_, err = compressor.NewBackupClient(conn).RegisterMulti(ctx, vecs, copts...)
		if err != nil {
			return nil, err
		}
		return
	})
	return
}

func (c *client) Remove(ctx context.Context, uuid string) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/internal/client/v1/client/compressor/Client.Remove")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = c.client.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (i interface{}, err error) {
		_, err = compressor.NewBackupClient(conn).Remove(ctx, &payload.Backup_Remove_Request{
			Uuid: uuid,
		}, copts...)
		if err != nil {
			return nil, err
		}
		return
	})
	return
}

func (c *client) RemoveMultiple(ctx context.Context, uuids ...string) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/internal/client/v1/client/compressor/Client.RemoveMultiple")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	req := new(payload.Backup_Remove_RequestMulti)
	req.Uuids = uuids
	_, err = c.client.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (i interface{}, err error) {
		_, err = compressor.NewBackupClient(conn).RemoveMulti(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return
	})
	return
}

func (c *client) RegisterIPs(ctx context.Context, ips []string) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/internal/client/v1/client/compressor/Client.RegisterIPs")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	req := new(payload.Backup_IP_Register_Request)
	req.Ips = ips
	_, err = c.client.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (i interface{}, err error) {
		_, err = compressor.NewBackupClient(conn).RegisterIPs(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return
	})
	return
}

func (c *client) RemoveIPs(ctx context.Context, ips []string) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/internal/client/v1/client/compressor/Client.RemoveIPs")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	req := new(payload.Backup_IP_Remove_Request)
	req.Ips = ips
	_, err = c.client.RoundRobin(ctx, func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (i interface{}, err error) {
		_, err = compressor.NewBackupClient(conn).RemoveIPs(ctx, req, copts...)
		if err != nil {
			return nil, err
		}
		return
	})
	return
}
