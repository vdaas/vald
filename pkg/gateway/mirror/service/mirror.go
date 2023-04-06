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
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	mclient "github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
)

// Mirror manages other mirror gateway connection.
// If there is a new Mirror Gateway components, registers new connection.
type Mirror interface {
	Start(ctx context.Context) (<-chan error, error)
	Connect(ctx context.Context, targets ...*payload.Mirror_Target) error
	MirrorTargets() ([]*payload.Mirror_Target, error)
}

type mirr struct {
	addrl         sync.Map // List of all connected addresses
	selfMirrAddrs []string // Address of self mirror gateway
	selfMirrAddrl sync.Map // List of self Mirror gateway addresses
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
	ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Mirror.startAdvertise")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ech := make(chan error, 100)

	tgts, err := m.mirrorAddrsToTargets(m.selfMirrAddrs...)
	if err != nil {
		close(ech)
		return nil, err
	}
	err = m.registers(ctx, tgts)
	if err != nil &&
		!errors.Is(err, errors.ErrTargetNotFound) &&
		!errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
		var attrs trace.Attributes

		switch {
		case errors.Is(err, context.Canceled):
			err = status.WrapWithCanceled(
				vald.InsertRPCName+" API canceld", err,
			)
			attrs = trace.StatusCodeCancelled(err.Error())
		case errors.Is(err, context.DeadlineExceeded):
			err = status.WrapWithDeadlineExceeded(
				vald.InsertRPCName+" API deadline exceeded", err,
			)
			attrs = trace.StatusCodeDeadlineExceeded(err.Error())
		default:
			var (
				st  *status.Status
				msg string
			)
			st, msg, err = status.ParseError(err, codes.Internal, "failed to parse "+mirror.RegisterRPCName+" gRPC error response")
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				span.SetStatus(trace.StatusError, err.Error())
			}
		}
		log.Warn(err)
		if span != nil {
			span.RecordError(err)
			span.SetAttributes(attrs...)
			span.SetStatus(trace.StatusError, err.Error())
		}
		close(ech)
		return nil, err
	}

	m.eg.Go(func() (err error) {
		tic := time.NewTicker(m.advertiseDur)
		defer close(ech)
		defer tic.Stop()

		for {
			select {
			case <-ctx.Done():
				return err
			case <-tic.C:
				tgts, err := m.mirrorAddrsToTargets(append(m.selfMirrAddrs, m.client.GRPCClient().ConnectedAddrs()...)...)
				if err != nil || len(tgts.GetTargets()) == 0 {
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
				resTgts, err := m.advertises(ctx, tgts)
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
				log.Debugf("[mirror]: connected mirror gateway targets: %v", m.client.GRPCClient().ConnectedAddrs())
			}
		}
	})
	return ech, nil
}

func (m *mirr) registers(ctx context.Context, tgts *payload.Mirror_Targets) error {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+mirror.MirrorRPCServiceName+"/"+mirror.RegisterRPCName), "vald/gateway/mirror/service/Mirror.registers")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	reqInfo := &errdetails.RequestInfo{
		ServingData: errdetails.Serialize(tgts),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + mirror.RegisterRPCName,
	}
	return m.broadCast(ctx,
		func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Mirror.registers/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			_, err := mirror.NewMirrorClient(conn).Register(ctx, tgts, copts...)
			if err != nil {
				var attrs trace.Attributes
				switch {
				case errors.Is(err, context.Canceled):
					err = status.WrapWithCanceled(
						mirror.RegisterRPCName+" API canceld", err, reqInfo, resInfo,
					)
					attrs = trace.StatusCodeCancelled(err.Error())
				case errors.Is(err, context.DeadlineExceeded):
					err = status.WrapWithCanceled(
						mirror.RegisterRPCName+" API deadline exceeded", err, reqInfo, resInfo,
					)
					attrs = trace.StatusCodeDeadlineExceeded(err.Error())
				case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
					err = status.WrapWithInternal(
						mirror.RegisterRPCName+" API connection not found", err, reqInfo, resInfo,
					)
					attrs = trace.StatusCodeInternal(err.Error())
				default:
					var (
						st  *status.Status
						msg string
					)
					st, msg, err = status.ParseError(err, codes.Internal,
						"failed to parse "+mirror.RegisterRPCName+" gRPC error response", reqInfo, resInfo,
					)
					attrs = trace.FromGRPCStatus(st.Code(), msg)
				}
				log.Error("failed to send Register API to %s\t: %v", target, err)
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(attrs...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return err
			}
			return nil
		},
	)
}

func (m *mirr) advertises(ctx context.Context, tgts *payload.Mirror_Targets) ([]*payload.Mirror_Target, error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+mirror.MirrorRPCServiceName+"/"+mirror.AdvertiseRPCName), "vald/gateway/mirror/service/Mirror.advertises")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	reqInfo := &errdetails.RequestInfo{
		ServingData: errdetails.Serialize(tgts),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + mirror.AdvertiseRPCName,
	}
	resTgts := make([]*payload.Mirror_Target, 0, len(tgts.GetTargets()))
	var mu sync.Mutex
	err := m.broadCast(ctx,
		func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
			ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Mirror.advertises/"+target)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			res, err := mirror.NewMirrorClient(conn).Advertise(ctx, tgts)
			if err != nil {
				var attrs trace.Attributes
				switch {
				case errors.Is(err, context.Canceled):
					err = status.WrapWithCanceled(
						mirror.AdvertiseRPCName+" API canceld", err, reqInfo, resInfo,
					)
					attrs = trace.StatusCodeCancelled(err.Error())
				case errors.Is(err, context.DeadlineExceeded):
					err = status.WrapWithCanceled(
						mirror.AdvertiseRPCName+" API deadline exceeded", err, reqInfo, resInfo,
					)
					attrs = trace.StatusCodeDeadlineExceeded(err.Error())
				case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
					err = status.WrapWithInternal(
						mirror.AdvertiseRPCName+" API connection not found", err, reqInfo, resInfo,
					)
					attrs = trace.StatusCodeInternal(err.Error())
				default:
					var (
						st  *status.Status
						msg string
					)
					st, msg, err = status.ParseError(err, codes.Internal,
						"failed to parse "+mirror.AdvertiseRPCName+" gRPC error response", reqInfo, resInfo,
					)
					attrs = trace.FromGRPCStatus(st.Code(), msg)
				}
				log.Errorf("failed to process advertise requst to %s\terror: %s", target, err.Error())
				if span != nil {
					span.RecordError(err)
					span.SetAttributes(attrs...)
					span.SetStatus(trace.StatusError, err.Error())
				}
				return err
			}
			if res != nil && len(res.GetTargets()) > 0 {
				mu.Lock()
				resTgts = append(resTgts, res.GetTargets()...)
				mu.Unlock()
			}
			return nil
		},
	)
	return resTgts, err
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
