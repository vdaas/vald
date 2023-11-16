// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package service

import (
	"context"
	"reflect"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

// Mirror represents an interface for managing mirroring operations.
// It provides methods for starting the mirroring service, connecting and disconnecting targets,
// checking the connectivity status of a given address, checking the existence of an address,
// retrieving all mirror targets, and iterating over all mirror addresses.
type Mirror interface {
	Start(ctx context.Context) <-chan error
	Connect(ctx context.Context, targets ...*payload.Mirror_Target) error
	Disconnect(ctx context.Context, targets ...*payload.Mirror_Target) error
	IsConnected(ctx context.Context, addr string) bool
	Exist(ctx context.Context, addr string) bool
	MirrorTargets() ([]*payload.Mirror_Target, error)
	RangeAllMirrorAddr(f func(addr string, _ any) bool)
}

type mirr struct {
	addrl         sync.Map[string, any]    // List of all connected addresses
	selfMirrTgts  []*payload.Mirror_Target // Targets of self mirror gateway
	selfMirrAddrl sync.Map[string, any]    // List of self Mirror gateway addresses
	gwAddrl       sync.Map[string, any]    // List of Vald Gateway addresses
	eg            errgroup.Group
	registerDur   time.Duration
	gateway       Gateway
}

// NewMirror creates the Mirror object with optional configuration options.
// It returns the initialized Mirror object and an error if the creation process fails.
func NewMirror(opts ...MirrorOption) (_ Mirror, err error) {
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

	m.selfMirrTgts = make([]*payload.Mirror_Target, 0)
	m.selfMirrAddrl.Range(func(addr string, _ any) bool {
		var (
			host string
			port uint16
		)
		host, port, err = net.SplitHostPort(addr)
		if err != nil {
			return false
		}
		m.selfMirrTgts = append(m.selfMirrTgts, &payload.Mirror_Target{
			Host: host,
			Port: uint32(port),
		})
		return true
	})
	return m, err
}

// Start starts the mirroring service.
// It returns a channel for receiving errors during the mirroring process.
func (m *mirr) Start(ctx context.Context) <-chan error {
	ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Mirror.Start")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	ech := make(chan error, 100)

	m.eg.Go(func() error {
		tic := time.NewTicker(m.registerDur)
		defer close(ech)
		defer tic.Stop()

		for {
			select {
			case <-ctx.Done():
			case <-tic.C:
				tgt, err := m.MirrorTargets()
				if err != nil {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case ech <- err:
						break
					}
				}

				resTgts, err := m.registers(ctx, &payload.Mirror_Targets{Targets: tgt})
				if err != nil || len(resTgts) == 0 {
					if !errors.Is(err, errors.ErrTargetNotFound) && len(resTgts) == 0 {
						err = errors.Join(err, errors.ErrTargetNotFound)
					} else if len(resTgts) == 0 {
						err = errors.ErrTargetNotFound
					}
					select {
					case <-ctx.Done():
						return ctx.Err()
					case ech <- err:
					}
				}
				if len(resTgts) > 0 {
					if err := m.Connect(ctx, resTgts...); err != nil {
						select {
						case <-ctx.Done():
							return ctx.Err()
						case ech <- err:
							break
						}
					}
				}
				log.Debugf("[mirror]: connected mirror gateway targets: %v", m.gateway.GRPCClient().ConnectedAddrs())
			}
		}
	})
	return ech
}

func (m *mirr) registers(ctx context.Context, tgts *payload.Mirror_Targets) ([]*payload.Mirror_Target, error) {
	ctx, span := trace.StartSpan(grpc.WithGRPCMethod(ctx, vald.PackageName+"."+vald.MirrorRPCServiceName+"/"+vald.RegisterRPCName), "vald/gateway/mirror/service/Mirror.registers")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	reqInfo := &errdetails.RequestInfo{
		ServingData: errdetails.Serialize(tgts),
	}
	resInfo := &errdetails.ResourceInfo{
		ResourceType: errdetails.ValdGRPCResourceTypePrefix + "/vald.v1." + vald.RegisterRPCName,
	}
	resTgts := make([]*payload.Mirror_Target, 0, len(tgts.GetTargets()))
	exists := make(map[string]bool)
	var mu sync.Mutex

	err := m.gateway.DoMulti(ctx, m.connectedMirrorAddrs(), func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error {
		ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Mirror.registers/"+target)
		defer func() {
			if span != nil {
				span.End()
			}
		}()

		res, err := vc.Register(ctx, tgts, copts...)
		if err != nil {
			var attrs trace.Attributes
			switch {
			case errors.Is(err, context.Canceled):
				err = status.WrapWithCanceled(
					vald.RegisterRPCName+" API canceld", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeCancelled(err.Error())
			case errors.Is(err, context.DeadlineExceeded):
				err = status.WrapWithCanceled(
					vald.RegisterRPCName+" API deadline exceeded", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeDeadlineExceeded(err.Error())
			case errors.Is(err, errors.ErrGRPCClientConnNotFound("*")):
				err = status.WrapWithInternal(
					vald.RegisterRPCName+" API connection not found", err, reqInfo, resInfo,
				)
				attrs = trace.StatusCodeInternal(err.Error())
			default:
				var (
					st  *status.Status
					msg string
				)
				st, msg, err = status.ParseError(err, codes.Internal,
					"failed to parse "+vald.RegisterRPCName+" gRPC error response", reqInfo, resInfo,
				)
				attrs = trace.FromGRPCStatus(st.Code(), msg)

				// When the ingress resource is deleted, the controller's default backend results(Unimplemented error) are returned so that the connection should be disconnected.
				// If it is a different namespace on the same cluster, the connection is automatically disconnected because the net.grpc health check fails.
				if st != nil && st.Code() == codes.Unimplemented {
					host, port, err := net.SplitHostPort(target)
					if err != nil {
						log.Warn(err)
					} else {
						if err := m.Disconnect(ctx, &payload.Mirror_Target{
							Host: host,
							Port: uint32(port),
						}); err != nil {
							log.Errorf("failed to disconnect %s, err: %v", target, err)
						}
					}
				}
			}
			log.Error("failed to send Register API to %s\t: %v", target, err)
			if span != nil {
				span.RecordError(err)
				span.SetAttributes(attrs...)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return err
		}
		if res != nil && len(res.GetTargets()) > 0 {
			for _, tgt := range res.GetTargets() {
				addr := net.JoinHostPort(tgt.Host, uint16(tgt.Port))
				mu.Lock()
				if !exists[addr] {
					exists[addr] = true
					resTgts = append(resTgts, res.GetTargets()...)
				}
				mu.Unlock()
			}
		}
		return nil
	})
	return resTgts, err
}

// Connect establishes gRPC connections to the specified Mirror targets, excluding this gateway and the LB Gateway.
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
		addr := net.JoinHostPort(target.GetHost(), uint16(target.GetPort())) // addr: host:port
		if !m.isSelfMirrorAddr(addr) && !m.isGatewayAddr(addr) {
			_, ok := m.addrl.Load(addr)
			if !ok || !m.IsConnected(ctx, addr) {
				_, err := m.gateway.GRPCClient().Connect(ctx, addr)
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

// Disconnect terminates gRPC connections to the specified Mirror targets.
func (m *mirr) Disconnect(ctx context.Context, targets ...*payload.Mirror_Target) error {
	ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Mirror.Disconnect")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(targets) == 0 {
		return errors.ErrTargetNotFound
	}
	for _, target := range targets {
		addr := net.JoinHostPort(target.GetHost(), uint16(target.GetPort()))
		if _, ok := m.gwAddrl.Load(addr); !ok {
			_, ok := m.addrl.Load(addr)
			if ok || m.IsConnected(ctx, addr) {
				if err := m.gateway.GRPCClient().Disconnect(ctx, addr); err != nil &&
					!errors.Is(err, errors.ErrGRPCClientConnNotFound(addr)) {
					return err
				}
				m.addrl.Delete(addr)
			}
		}
	}
	return nil
}

// IsConnected checks if the gRPC connection to the given address is connected.
func (m *mirr) IsConnected(ctx context.Context, addr string) bool {
	return m.gateway.GRPCClient().IsConnected(ctx, addr)
}

// Exist checks if the given address exists in the Mmirror.
func (m *mirr) Exist(_ context.Context, addr string) bool {
	_, ok := m.addrl.Load(addr)
	return ok
}

// MirrorTargets returns the Mirror targets, including the address of this gateway and the addresses of other Mirror Gateways
// to which this gateway is currently connected.
func (m *mirr) MirrorTargets() ([]*payload.Mirror_Target, error) {
	addrs := m.gateway.GRPCClient().ConnectedAddrs()
	tgts := make([]*payload.Mirror_Target, 0, len(addrs)+1)
	tgts = append(tgts, m.selfMirrTgts...)
	for _, addr := range addrs {
		if !m.isGatewayAddr(addr) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			tgts = append(tgts, &payload.Mirror_Target{
				Host: host,
				Port: uint32(port),
			})
		}
	}
	return tgts, nil
}

func (m *mirr) isSelfMirrorAddr(addr string) bool {
	_, ok := m.selfMirrAddrl.Load(addr)
	return ok
}

func (m *mirr) isGatewayAddr(addr string) bool {
	_, ok := m.gwAddrl.Load(addr)
	return ok
}

// connectedMirrorAddrs returns the addresses of other Mirror Gateways to which this gateway is currently connected.
func (m *mirr) connectedMirrorAddrs() []string {
	connectedAddrs := m.gateway.GRPCClient().ConnectedAddrs()
	addrs := make([]string, 0, len(connectedAddrs))
	for _, addr := range connectedAddrs {
		if !m.isSelfMirrorAddr(addr) &&
			!m.isGatewayAddr(addr) {
			addrs = append(addrs, addr)
		}
	}
	return addrs
}

func (m *mirr) RangeAllMirrorAddr(f func(addr string, _ any) bool) {
	m.addrl.Range(func(addr string, value any) bool {
		if !m.isGatewayAddr(addr) {
			if !f(addr, value) {
				return false
			}
		}
		return true
	})
}

func (m *mirr) toMirrorTargets(addrs ...string) ([]*payload.Mirror_Target, error) {
	tgts := make([]*payload.Mirror_Target, 0, len(addrs))
	for _, addr := range addrs {
		if ok := m.isGatewayAddr(addr); !ok {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			tgts = append(tgts, &payload.Mirror_Target{
				Host: host,
				Port: uint32(port),
			})
		}
	}
	return tgts, nil
}
