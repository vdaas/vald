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
	"reflect"

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
	NextRotateOffset(total int) int
	BroadCastSlice(ctx context.Context, kind BroadCastKind, addrs []string, offset, count int,
		f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error) error
}

type BroadCastKind int

const (
	READ BroadCastKind = iota
	WRITE
)

type gateway struct {
	client  discoverer.Client
	vcCache vc.ConnCache
	_       [128]byte
	rrIdx   atomic.Uint64
}

func (g *gateway) getValdClient(conn *grpc.ClientConn) vald.Client {
	return g.vcCache.Get(conn)
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

const (
	defaultBroadcastRanCap    = 8
	maxBroadcastRanCap        = 256
	defaultDoMultiVisitedCap  = 16
	maxDoMultiVisitedSliceCap = 256
)

type broadcastCarrier struct {
	g      *gateway
	f      func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error
	cbFn   func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error
	ran    map[string]struct{}
	ferr   map[string]error
	kind   BroadCastKind
	hasErr atomic.Bool
	mu     sync.Mutex
}

func (c *broadcastCarrier) callback(
	ictx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
) error {
	if ictx.Err() != nil {
		return nil //nolint:nilerr // canceled worker returns nil so Range* collects healthy agents
	}
	cerr := c.f(ictx, addr, c.g.getValdClient(conn), copts...)
	if c.kind == READ && cerr == nil && !c.hasErr.Load() {
		return nil
	}
	c.mu.Lock()
	if c.kind == WRITE {
		c.ran[addr] = struct{}{}
	}
	if cerr != nil {
		if c.ferr == nil {
			c.ferr = make(map[string]error)
		}
		c.ferr[addr] = cerr
		if c.kind == READ {
			c.hasErr.Store(true)
		}
	} else if c.ferr != nil {
		delete(c.ferr, addr)
	}
	c.mu.Unlock()
	return cerr
}

type doMultiCarrier struct {
	g       *gateway
	f       func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error
	visitFn func(ctx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) error
	visited []string
	cur     atomic.Int64
	limit   int64
}

func (c *doMultiCarrier) visit(
	ictx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
) error {
	if c.cur.Load() >= c.limit {
		return nil
	}
	if cerr := c.f(ictx, addr, c.g.getValdClient(conn), copts...); cerr != nil {
		return cerr
	}
	if pos := int(c.cur.Add(1) - 1); pos < len(c.visited) {
		c.visited[pos] = addr
	}
	return nil
}

var (
	broadcastCarrierPool = sync.Pool[*broadcastCarrier]{ //nolint:gochecknoglobals // broadcast callback state is reused across BroadCast/BroadCastSlice calls
		New: func() *broadcastCarrier {
			c := &broadcastCarrier{
				ran: make(map[string]struct{}, defaultBroadcastRanCap),
			}
			c.cbFn = c.callback
			return c
		},
	}
	doMultiCarrierPool = sync.Pool[*doMultiCarrier]{ //nolint:gochecknoglobals // DoMulti callback state and visited slice are reused across DoMulti calls
		New: func() *doMultiCarrier {
			c := &doMultiCarrier{
				visited: make([]string, 0, defaultDoMultiVisitedCap),
			}
			c.visitFn = c.visit
			return c
		},
	}
)

func getBroadcastCarrier(
	g *gateway,
	kind BroadCastKind,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) *broadcastCarrier {
	c := broadcastCarrierPool.Get()
	c.g = g
	c.kind = kind
	c.f = f
	c.hasErr.Store(false)
	return c
}

func putBroadcastCarrier(c *broadcastCarrier) {
	if len(c.ran) > maxBroadcastRanCap {
		c.ran = make(map[string]struct{}, defaultBroadcastRanCap)
	} else {
		clear(c.ran)
	}
	if len(c.ferr) > maxBroadcastRanCap {
		c.ferr = nil
	} else if c.ferr != nil {
		clear(c.ferr)
	}
	c.g = nil
	c.f = nil
	broadcastCarrierPool.Put(c)
}

func getDoMultiCarrier(
	g *gateway,
	limitInt int,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) *doMultiCarrier {
	c := doMultiCarrierPool.Get()
	c.g = g
	c.f = f
	c.limit = int64(limitInt)
	c.cur.Store(0)
	if cap(c.visited) >= limitInt {
		c.visited = c.visited[:limitInt]
	} else {
		c.visited = make([]string, limitInt)
	}
	return c
}

func putDoMultiCarrier(c *doMultiCarrier) {
	if cap(c.visited) > maxDoMultiVisitedSliceCap {
		c.visited = make([]string, 0, defaultDoMultiVisitedCap)
	} else {
		clear(c.visited)
		c.visited = c.visited[:0]
	}
	c.g = nil
	c.f = nil
	doMultiCarrierPool.Put(c)
}

func (g *gateway) clientForKind(kind BroadCastKind) grpc.Client {
	switch kind {
	case READ:
		return g.client.GetReadClient()
	case WRITE:
		return g.client.GetClient()
	default:
		return g.client.GetClient()
	}
}

func (g *gateway) broadCastWith(
	fctx context.Context,
	client grpc.Client,
	kind BroadCastKind,
	addrs []string,
	offset, count int,
	isSlice bool,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) error {
	c := getBroadcastCarrier(g, kind, f)
	defer putBroadcastCarrier(c)

	var err error
	if isSlice {
		err = client.RangeSliceConcurrent(fctx, addrs, offset, count, -1, c.cbFn)
	} else {
		err = client.RangeConcurrent(fctx, -1, c.cbFn)
	}
	if err != nil {
		return err
	}
	if kind == READ && !c.hasErr.Load() {
		return nil
	}
	// RangeConcurrent/RangeSliceConcurrent has joined all workers before reaching
	// this point, so c.ran and c.ferr are read without acquiring c.mu.
	if kind == WRITE && fctx.Err() == nil {
		if isSlice {
			n := len(addrs)
			for i := range count {
				idx := offset + i
				if idx >= n {
					idx -= n
				}
				addr := addrs[idx]
				if _, ok := c.ran[addr]; !ok {
					if c.ferr == nil {
						c.ferr = make(map[string]error)
					}
					c.ferr[addr] = status.Errorf(codes.Unavailable,
						"agent %s was not reached: no healthy connection, or the request was rejected before it was sent", addr)
				}
			}
		} else {
			for _, addr := range addrs {
				if _, ok := c.ran[addr]; !ok {
					if c.ferr == nil {
						c.ferr = make(map[string]error)
					}
					c.ferr[addr] = status.Errorf(codes.Unavailable,
						"agent %s was not reached: no healthy connection, or the request was rejected before it was sent", addr)
				}
			}
		}
	}
	return errors.JoinMap(c.ferr)
}

func (g *gateway) BroadCast(
	ctx context.Context,
	kind BroadCastKind,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) (err error) {
	fctx := ctx
	var span trace.Span
	if trace.Enabled() {
		fctx, span = trace.StartSpan(ctx, "vald/gateway-lb/service/Gateway.BroadCast")
		defer trace.End(span)
	}

	client := g.clientForKind(kind)

	// The agents a write has to reach are fixed before the broadcast. A discoverer update
	// can change the list while it runs, and an agent that fell out of it in the meantime
	// must not make the write look complete: it can rejoin with the object this write was
	// meant to remove. The price is that an agent removed on purpose (a scale-down) during
	// the broadcast is reported as unreached too, which the caller can retry.
	var targets []string
	if kind == WRITE {
		targets = g.client.GetAddrs(fctx)
	}
	return g.broadCastWith(fctx, client, kind, targets, 0, 0, false, f)
}

// NextRotateOffset returns the next starting index in [0, total) for round-robin subset routing.
func (g *gateway) NextRotateOffset(total int) int {
	if total <= 1 {
		return 0
	}
	return int((g.rrIdx.Add(1) - 1) % uint64(total)) //nolint:gosec // modulo by total bounds result within [0, total)
}

// BroadCastSlice runs f against count agents starting at cyclic offset in addrs.
func (g *gateway) BroadCastSlice(
	ctx context.Context,
	kind BroadCastKind,
	addrs []string,
	offset int,
	count int,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) error {
	n := len(addrs)
	if n == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	if count <= 0 {
		return nil
	}
	if count > n {
		count = n
	}
	offset %= n
	if offset < 0 {
		offset += n
	}

	fctx := ctx
	var span trace.Span
	if trace.Enabled() {
		fctx, span = trace.StartSpan(ctx, "vald/gateway-lb/service/Gateway.BroadCastSlice")
		defer trace.End(span)
	}

	client := g.clientForKind(kind)
	if kind == READ && client != g.client.GetClient() {
		// ReadReplica clients key snapByAddr by service DNS names rather than primary pod IPs,
		// so fall back to full RangeConcurrent on the chosen read client instead of passing pod IPs.
		return g.broadCastWith(fctx, client, kind, nil, 0, 0, false, f)
	}
	return g.broadCastWith(fctx, client, kind, addrs, offset, count, true, f)
}

func (g *gateway) DoMulti(
	ctx context.Context,
	num int,
	f func(ctx context.Context, target string, ac vald.Client, copts ...grpc.CallOption) error,
) (err error) {
	sctx := ctx
	var span trace.Span
	if trace.Enabled() {
		sctx, span = trace.StartSpan(ctx, "vald/gateway-lb/service/Gateway.DoMulti")
		defer trace.End(span)
	}
	addrs := g.client.GetAddrs(sctx)
	if len(addrs) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	limitInt := min(max(num, 0), len(addrs))
	limit := int64(limitInt)
	if limitInt == 0 {
		return nil
	}
	offset := 0
	if len(addrs) > limitInt {
		offset = g.NextRotateOffset(len(addrs))
	}

	c := getDoMultiCarrier(g, limitInt, f)
	defer putDoMultiCarrier(c)

	client := g.client.GetClient()
	err = client.RangeSliceConcurrent(sctx, addrs, offset, limitInt, limitInt, c.visitFn)
	if err != nil || c.cur.Load() < limit {
		recorded := min(int(c.cur.Load()), len(c.visited))
		seen := make(map[string]struct{}, recorded)
		for _, addr := range c.visited[:recorded] {
			if addr != "" {
				seen[addr] = struct{}{}
			}
		}
		remaining := make([]string, 0, len(addrs))
		for _, addr := range addrs {
			if _, ok := seen[addr]; !ok {
				remaining = append(remaining, addr)
			}
		}
		if len(remaining) > 0 {
			err = client.OrderedRange(sctx, remaining, func(ictx context.Context,
				addr string,
				conn *grpc.ClientConn,
				copts ...grpc.CallOption,
			) error {
				if c.cur.Load() < limit {
					if _, ok := seen[addr]; !ok {
						if cerr := f(ictx, addr, g.getValdClient(conn), copts...); cerr != nil {
							return cerr
						}
						c.cur.Add(1)
						seen[addr] = struct{}{}
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
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
