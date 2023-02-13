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
	"fmt"
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
	"github.com/vdaas/vald/internal/observability/trace"
)

const (
	forwardedContextKey   = "forwarded-for"
	forwardedContextValue = "gateway mirror"
)

type Gateway interface {
	Start(ctx context.Context) (<-chan error, error)
	MirrorTargets() ([]*payload.Mirror_Target, error)
	OtherMirrorAddrs() []string
	Connect(ctx context.Context, targets ...*payload.Mirror_Target) ([]*payload.Mirror_Target, error)
	ForwardedContext(ctx context.Context, podName string) context.Context
	FromForwardedContext(ctx context.Context) string
	BroadCast(ctx context.Context,
		f func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error) error
}

type gateway struct {
	client       mclient.Client // Mirror Gateway client for other cluster.
	iclient      mclient.Client // Mirror Gateway client for the same cluster.
	eg           errgroup.Group
	advertiseDur time.Duration
	podName      string
}

func NewGateway(opts ...Option) (gw Gateway, err error) {
	g := new(gateway)
	for _, opt := range append(defaultGWOpts, opts...) {
		if err := opt(g); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(err, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return g, nil
}

func (g *gateway) Start(ctx context.Context) (<-chan error, error) {
	ech := make(chan error, 100)

	cech, err := g.client.Start(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}
	icech, err := g.iclient.Start(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}
	aech, err := g.startAdvertise(ctx)
	if err != nil {
		close(ech)
		return nil, err
	}

	g.eg.Go(func() (err error) {
		defer close(ech)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-cech:
			case err = <-icech:
			case err = <-aech:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ech <- err:
				}
				err = nil
			}
		}
	})
	return ech, nil
}

func (g *gateway) startAdvertise(ctx context.Context) (<-chan error, error) {
	tic := time.NewTicker(g.advertiseDur)

	tgts, err := g.selfMirrorTargets()
	if err != nil {
		return nil, err
	}
	req := &payload.Mirror_Targets{
		Targets: tgts,
	}
	err = g.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error {
		_, err := mirror.NewMirrorClient(conn).Register(ctx, req, copts...)
		return err
	})
	if err != nil {
		return nil, err
	}

	ech := make(chan error, 100)
	g.eg.Go(func() error {
		defer close(ech)
		defer tic.Stop()
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-tic.C:
				tgts, err := g.MirrorTargets()
				if err != nil || len(tgts) == 0 {
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
				req := &payload.Mirror_Targets{
					Targets: tgts,
				}
				resTgts := make([]*payload.Mirror_Target, 0, len(tgts))
				mutex := sync.Mutex{}
				err = g.BroadCast(ctx, func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) (err error) {
					res, err := mirror.NewMirrorClient(conn).Advertise(ctx, req, copts...)
					if err != nil {
						return err
					}
					mutex.Lock()
					resTgts = append(resTgts, res.GetTargets()...)
					mutex.Unlock()
					return nil
				})
				if err != nil {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case ech <- err:
					}
				} else {
					resTgts, err = g.Connect(ctx, resTgts...)
					if err != nil {
						select {
						case <-ctx.Done():
							return ctx.Err()
						case ech <- err:
						}
					}
					log.Infof("connected mirror gateway targets: %v", resTgts)
				}
			}
		}
	})
	return ech, nil
}

func (g *gateway) ForwardedContext(ctx context.Context, podName string) context.Context {
	return grpc.NewOutgoingContext(ctx, grpc.MD{
		forwardedContextKey: []string{
			podName,
		},
	})
}

func (g *gateway) FromForwardedContext(ctx context.Context) string {
	md, ok := grpc.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	vals, ok := md[forwardedContextKey]
	if !ok {
		return ""
	}
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// BroadCast executes a broadcast operation to the Mirror Gateway on other clusters.
func (g *gateway) BroadCast(ctx context.Context,
	f func(ctx context.Context, target string, conn *grpc.ClientConn, copts ...grpc.CallOption) error,
) (err error) {
	fctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Gateway.BroadCast")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	return g.client.GRPCClient().RangeConcurrent(g.ForwardedContext(fctx, g.podName), -1, func(ictx context.Context,
		addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (err error) {
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

// selfAddrs returns the addresses of Mirror Gateway on the same cluster.
func (g *gateway) selfMirrorAddrs() (m map[string]struct{}) {
	m = make(map[string]struct{})
	for _, addr := range g.iclient.GRPCClient().ConnectedAddrs() {
		m[addr] = struct{}{}
	}
	return m
}

// SelfAddrs returns the Targets of Mirror Gateway on the same cluster.
func (g *gateway) selfMirrorTargets() (tgts []*payload.Mirror_Target, err error) {
	tgts, err = g.addrsToTargets(g.iclient.GRPCClient().ConnectedAddrs())
	if err != nil {
		return nil, err
	}
	return tgts, nil
}

// OtherMirrorAddrs returns the addresses of Mirror Gateway on the other clusters.
func (g *gateway) OtherMirrorAddrs() (oaddrs []string) {
	selfAddrs := g.selfMirrorAddrs()
	addrs := g.client.GRPCClient().ConnectedAddrs()
	oaddrs = make([]string, 0, len(addrs))
	for _, addr := range g.client.GRPCClient().ConnectedAddrs() {
		if _, ok := selfAddrs[addr]; !ok {
			oaddrs = append(oaddrs, addr)
		}
	}
	return oaddrs
}

// MirrorTargets returns the connected Mirror Gateway targets.
func (g *gateway) MirrorTargets() (tgts []*payload.Mirror_Target, err error) {
	tgts, err = g.addrsToTargets(g.client.GRPCClient().ConnectedAddrs())
	if err != nil {
		return nil, err
	}
	return tgts, nil
}

func (g *gateway) addrsToTargets(addrs []string) ([]*payload.Mirror_Target, error) {
	targets := make([]*payload.Mirror_Target, 0, len(addrs))
	for _, addr := range addrs {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		targets = append(targets, &payload.Mirror_Target{
			Ip:   host,
			Port: uint32(port),
		})
	}
	return targets, nil
}

func (g *gateway) Connect(ctx context.Context, targets ...*payload.Mirror_Target) ([]*payload.Mirror_Target, error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Gateway.Connect")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	eg, egctx := errgroup.New(ctx)
	selfAddrs := g.selfMirrorAddrs()

	for _, target := range targets {
		if target == nil {
			continue
		}
		addr := fmt.Sprintf("%s:%d", target.GetIp(), target.GetPort())
		if _, ok := selfAddrs[addr]; ok {
			continue
		}
		eg.Go(func() error {
			if !g.client.GRPCClient().IsConnected(egctx, addr) {
				_, err := g.client.GRPCClient().Connect(egctx, addr, g.client.GRPCClient().GetDialOption()...)
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	tgts, err := g.MirrorTargets()
	if err != nil {
		return nil, err
	}
	return tgts, nil
}
