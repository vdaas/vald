// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

// Package service
package service

import (
	"context"
	"maps"
	"reflect"
	"slices"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/client/v1/client/discoverer"
	vc "github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
)

type Gateway interface {
	Start(ctx context.Context) (<-chan error, error)
	GetAgentCount(ctx context.Context) int
	Addrs(ctx context.Context) []string
	DoMulti(ctx context.Context, num int,
		f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error) error
	// BroadCast runs f against every agent. The agents' errors are not dropped: the error of f for
	// each failing agent is returned joined (sorted by target), and the caller decides what a failure
	// means for its RPC (for example search tolerates some of them, Remove and Flush do not).
	// For WRITE, an agent of the discoverer's list that was known when the call started and whose
	// callback did not run (no healthy connection, circuit breaker open) is reported as Unavailable,
	// while every reachable agent is still written. READ does not report unreached agents.
	BroadCast(ctx context.Context, kind BroadCastKind,
		f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error) error
}

type BroadCastKind int

const (
	READ BroadCastKind = iota
	WRITE
)

type gateway struct {
	client discoverer.Client
}

func NewGateway(opts ...Option) (gw Gateway, err error) {
	g := new(gateway)
	for _, opt := range opts {
		if err := opt(g); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return g, nil
}

func (g *gateway) Start(ctx context.Context) (<-chan error, error) {
	return g.client.Start(ctx)
}

func (g *gateway) BroadCast(
	ctx context.Context,
	kind BroadCastKind,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) (err error) {
	fctx, span := trace.StartSpan(ctx, "vald/gateway-lb/service/Gateway.BroadCast")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	var client grpc.Client
	switch kind {
	case READ:
		client = g.client.GetReadClient()
	case WRITE:
		client = g.client.GetClient()
	}

	// The agents a write has to reach are fixed before the broadcast. A discoverer update
	// can change the list while it runs, and an agent that fell out of it in the meantime
	// must not make the write look complete: it can rejoin with the object this write was
	// meant to remove. The price is that an agent removed on purpose (a scale-down) during
	// the broadcast is reported as unreached too, which the caller can retry.
	var targets []string
	if kind == WRITE {
		targets = g.client.GetAddrs(fctx)
	}

	// RangeConcurrent drops every per-target error except a local ctx error, so
	// the callback errors are collected here; otherwise a failing agent would look
	// to the handlers like one that has nothing to report. Only the last result per
	// target is kept because the client backoff may re-invoke f for that target.
	var (
		mu   sync.Mutex
		ferr = make(map[string]error)
		// ran holds the targets whose callback was invoked.
		ran = make(map[string]struct{})
	)
	err = client.RangeConcurrent(fctx, -1, func(ictx context.Context,
		addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (err error) {
		select {
		case <-ictx.Done():
			return nil
		default:
			err = f(ictx, addr, vc.NewFromConn(conn), copts...)
			mu.Lock()
			ran[addr] = struct{}{}
			if err != nil {
				ferr[addr] = err
			} else {
				delete(ferr, addr)
			}
			mu.Unlock()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	mu.Lock()
	defer mu.Unlock()
	// RangeConcurrent also drops a target whose connection could not be acquired
	// (or whose circuit breaker is open) before the callback ran. A write is still
	// sent to every reachable agent, but the unreached ones must be reported: a
	// Remove or Flush that skipped an agent is not a success. Reads keep serving
	// partial results, so they are left as they are.
	if kind == WRITE && fctx.Err() == nil {
		for _, addr := range targets {
			if _, ok := ran[addr]; !ok {
				ferr[addr] = status.Errorf(codes.Unavailable,
					"agent %s was not reached: no healthy connection, or the request was rejected before it was sent", addr)
			}
		}
	}
	if len(ferr) == 0 {
		return nil
	}
	// sorted by target so the leading error (whose code a joined error reports)
	// does not depend on map order.
	errs := make([]error, 0, len(ferr))
	for _, addr := range slices.Sorted(maps.Keys(ferr)) {
		errs = append(errs, ferr[addr])
	}
	return errors.Join(errs...)
}

func (g *gateway) DoMulti(
	ctx context.Context,
	num int,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) (err error) {
	sctx, span := trace.StartSpan(ctx, "vald/gateway-lb/service/Gateway.DoMulti")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	var cur uint32 = 0
	addrs := g.client.GetAddrs(sctx)
	var limit uint32
	if len(addrs) < num {
		limit = uint32(len(addrs))
	} else {
		limit = uint32(num)
	}
	var visited sync.Map[string, any]
	err = g.client.GetClient().OrderedRange(sctx, addrs, func(ictx context.Context,
		addr string,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (err error) {
		if atomic.LoadUint32(&cur) < limit {
			err = f(ictx, addr, vc.NewFromConn(conn), copts...)
			if err != nil {
				return err
			}
			atomic.AddUint32(&cur, 1)
			visited.Store(addr, struct{}{})
		}
		return nil
	})
	if err != nil || cur < limit {
		err = g.client.GetClient().OrderedRange(sctx, addrs, func(ictx context.Context,
			addr string,
			conn *grpc.ClientConn,
			copts ...grpc.CallOption,
		) (err error) {
			if atomic.LoadUint32(&cur) < limit {
				_, ok := visited.Load(addr)
				if !ok {
					err = f(ictx, addr, vc.NewFromConn(conn), copts...)
					if err != nil {
						return err
					}
					atomic.AddUint32(&cur, 1)
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *gateway) GetAgentCount(ctx context.Context) int {
	return len(g.Addrs(ctx))
}

func (g *gateway) Addrs(ctx context.Context) []string {
	return g.client.GetAddrs(ctx)
}
