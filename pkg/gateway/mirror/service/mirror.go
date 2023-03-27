//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package service
package service

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	mclient "github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

// Mirror manages other mirror gateway connection.
// If there is a new Mirror Gateway components, mirror create register new connection.
type Mirror interface {
	Start(ctx context.Context) (<-chan error, error)
	Connect(ctx context.Context, targets ...*payload.Mirror_Target) error
	MirrorTargets() ([]*payload.Mirror_Target, error)
}

type mirr struct {
	addrl         sync.Map // List of all connected addresses
	selfMirrAddrs []string // Address of my mirror gateway
	selfMirrAddrl sync.Map // List of my Mirror gateway addresses
	gwAddrs       []string // Address of Vald Gateway (LB gateway)
	gwAddrl       sync.Map // List of Vald Gateway addresses
	client        mclient.Client
	eg            errgroup.Group
	advertiseDur  time.Duration
}

func NewMirror(opts ...MirrorOption) (Mirror, error) {
	m := new(mirr)
	for _, opt := range append(defaultMirrOpts, opts...) {
		if err := opt(m); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(err, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	for _, addr := range m.selfMirrAddrs {
		m.selfMirrAddrl.Store(addr, struct{}{})
	}
	for _, addr := range m.gwAddrs {
		m.gwAddrl.Store(addr, struct{}{})
	}
	return m, nil
}

func (m *mirr) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 100)

	aech, err := m.startAdvertise(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	m.eg.Go(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-aech:
			}
			if err != nil {
				select {
				case <-ctx.Done():
				case ech <- err:
				}
				err = nil
			}
		}
	})
	return ech, nil
}

func (m *mirr) startAdvertise(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 100)

	req, err := m.mirrorAddrsToTargets(m.selfMirrAddrs...)
	if err != nil {
		close(ech)
		return nil, err
	}
	err = m.broadCast(ctx,
		func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Mirror.startAdvertise/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			_, err := mirror.NewMirrorClient(conn).Register(ctx, req)
			if err != nil {
				st, msg, err := status.ParseError(err, codes.Internal,
					"failed to parse Register API gRPC error response",
				)
				log.Errorf("failed to process register requst to %s\terror: %s", target, err.Error())
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return err
			}
			return nil
		},
	)
	if err != nil && !errors.Is(err, errors.ErrTargetNotFound) {
		close(ech)
		return nil, err
	}

	m.eg.Go(func() (err error) {
		tic := time.NewTicker(m.advertiseDur)
		mutex := new(sync.Mutex)
		defer close(ech)
		defer tic.Stop()

		for {
			select {
			case <-ctx.Done():
				return err
			case <-tic.C:
				req, err := m.mirrorAddrsToTargets(append(m.selfMirrAddrs, m.client.GRPCClient().ConnectedAddrs()...)...)
				if err != nil || len(req.GetTargets()) == 0 {
					if err == nil {
						err = errors.ErrTargetNotFound
					}
					select {
					case <-ctx.Done():
						return ctx.Err()
					case ech <- err:
					}
					continue
				}
				resTgts := make([]*payload.Mirror_Target, 0, len(req.GetTargets()))
				err = m.broadCast(ctx,
					func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
						ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Mirror.startAdvertise/"+target)
						defer func() {
							if span != nil {
								span.End()
							}
						}()
						res, err := mirror.NewMirrorClient(conn).Advertise(ctx, req)
						if err != nil {
							st, msg, err := status.ParseError(err, codes.Internal,
								"failed to parse Advertise API gRPC error response",
							)
							log.Errorf("failed to process advertise requst to %s\terror: %s", target, err.Error())
							if span != nil {
								span.RecordError(err)
								span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
								span.SetStatus(trace.StatusError, err.Error())
							}
							return err
						}
						if res != nil && len(res.GetTargets()) > 0 {
							mutex.Lock()
							resTgts = append(resTgts, res.GetTargets()...)
							mutex.Unlock()
						}
						return nil
					},
				)
				if err != nil || len(resTgts) == 0 {
					if err == nil {
						err = errors.ErrTargetNotFound
					}
					select {
					case <-ctx.Done():
						return ctx.Err()
					case ech <- err:
					}
					continue
				}
				if err = m.Connect(ctx, resTgts...); err != nil {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case ech <- err:
					}
				}
				log.Infof("[mirror]: connected mirror gateway targets: %v", m.client.GRPCClient().ConnectedAddrs())
			}
		}
	})
	return ech, nil
}

func (m *mirr) Connect(ctx context.Context, targets ...*payload.Mirror_Target) error {
	ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Mirror.Connect")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(targets) == 0 {
		return errors.ErrTargetNotFound
	}
	for _, target := range targets {
		addr := net.JoinHostPort(target.GetIp(), uint16(target.GetPort())) // addr: host:port
		if !m.isSelfMirrorAddr(addr) && !m.isGatewayAddr(addr) {
			_, ok := m.addrl.Load(addr)
			if !ok || !m.client.GRPCClient().IsConnected(ctx, addr) {
				_, err := m.client.GRPCClient().Connect(ctx, addr)
				if err != nil {
					m.addrl.Delete(addr)
					return err
				}
			}
			m.addrl.Store(addr, struct{}{})
		}
	}
	return nil
}

func (m *mirr) MirrorTargets() ([]*payload.Mirror_Target, error) {
	addrs := append(m.selfMirrAddrs, m.client.GRPCClient().ConnectedAddrs()...)
	tgts := make([]*payload.Mirror_Target, 0, len(addrs))
	for _, addr := range addrs {
		if !m.isGatewayAddr(addr) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			tgts = append(tgts, &payload.Mirror_Target{
				Ip:   host,
				Port: uint32(port),
			})
		}
	}
	return tgts, nil
}

func (m *mirr) broadCast(ctx context.Context,
	f func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) (err error) {
	fctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Mirror.broadCast")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	connectedAddrs := m.client.GRPCClient().ConnectedAddrs()
	addrs := make([]string, 0, len(connectedAddrs))
	for _, addr := range connectedAddrs {
		if !m.isSelfMirrorAddr(addr) && !m.isGatewayAddr(addr) {
			addrs = append(addrs, addr)
		}
	}
	if len(addrs) == 0 {
		return errors.ErrTargetNotFound
	}

	return m.client.GRPCClient().OrderedRangeConcurrent(fctx, addrs, -1, func(ictx context.Context,
		addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
	) error {
		select {
		case <-ictx.Done():
			return nil
		default:
			err = f(ictx, addr, conn, copts...)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *mirr) isSelfMirrorAddr(addr string) bool {
	if _, ok := m.selfMirrAddrl.Load(addr); ok {
		return true
	}
	return false
}

func (m *mirr) isGatewayAddr(addr string) bool {
	if _, ok := m.gwAddrl.Load(addr); ok {
		return true
	}
	return false
}

func (m *mirr) mirrorAddrsToTargets(addrs ...string) (*payload.Mirror_Targets, error) {
	tgts := make([]*payload.Mirror_Target, 0, len(addrs))
	for _, addr := range addrs {
		if ok := m.isGatewayAddr(addr); !ok {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			tgts = append(tgts, &payload.Mirror_Target{
				Ip:   host,
				Port: uint32(port),
			})
		}
	}
	return &payload.Mirror_Targets{
		Targets: tgts,
	}, nil
}
