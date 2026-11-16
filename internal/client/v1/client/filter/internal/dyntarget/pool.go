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

// Package dyntarget tracks the filter targets named in requests for the
// ingress and egress filter clients, so that their connections are reference
// counted and disconnected once they stay unused for an idle timeout.
package dyntarget

import (
	"context"
	"slices"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync"
)

// Pool reference counts the connections of request-derived targets of c.
type Pool struct {
	c     grpc.Client
	dyn   map[string]*dynamicTarget // guarded by mu; see Acquire
	addrs []string
	idle  time.Duration // how long an unused request-derived target stays connected
	mu    sync.Mutex
	// closed is guarded by mu; once set (Stop or the Start ctx is done) no idle
	// disconnect is scheduled or run.
	closed bool
}

// dynamicTarget tracks an address that is not one of the pool's configured
// addrs (i.e. it comes from a request). refs, gen and evict are guarded by
// Pool.mu; mu serializes its Connect and Disconnect.
type dynamicTarget struct {
	evict *time.Timer // pending idle disconnect, nil while in use
	mu    sync.Mutex
	refs  int    // in-flight Target uses
	gen   uint64 // identifies the latest scheduled idle disconnect
}

// New returns a Pool for the connections of c. addrs are the configured
// targets, which are never disconnected by the Pool; idle must be positive.
func New(c grpc.Client, addrs []string, idle time.Duration) *Pool {
	return &Pool{
		c:     c,
		addrs: addrs,
		idle:  idle,
	}
}

// Start stops the idle disconnects once ctx is done.
func (p *Pool) Start(ctx context.Context) {
	context.AfterFunc(ctx, p.Stop)
}

// Stop cancels every pending idle disconnect and prevents new ones; the
// connections left are closed by the owner of the gRPC client.
func (p *Pool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	for _, dt := range p.dyn {
		if dt.evict != nil {
			dt.evict.Stop()
			dt.evict = nil
		}
	}
}

// Acquire connects targets. Every successful call must be paired with
// Release(ctx, targets...) once the targets are no longer used; on error the
// targets already acquired are released.
func (p *Pool) Acquire(ctx context.Context, targets ...string) error {
	for i, addr := range targets {
		if err := p.acquire(ctx, addr); err != nil {
			p.Release(ctx, targets[:i]...)
			return err
		}
	}
	return nil
}

// acquire connects addr if needed; a request-derived addr is reference counted
// so that Release can schedule its disconnect once its last user is done, and a
// pending disconnect is cancelled here.
func (p *Pool) acquire(ctx context.Context, addr string) error {
	var dt *dynamicTarget
	if !slices.Contains(p.addrs, addr) {
		p.mu.Lock()
		if p.dyn == nil {
			p.dyn = make(map[string]*dynamicTarget)
		}
		dt = p.dyn[addr]
		if dt == nil {
			dt = new(dynamicTarget)
			p.dyn[addr] = dt
		}
		dt.refs++
		if dt.evict != nil {
			// An already running callback sees refs > 0 and does nothing.
			dt.evict.Stop()
			dt.evict = nil
		}
		p.mu.Unlock()
		dt.mu.Lock()
		defer dt.mu.Unlock()
	}
	if p.c.IsConnected(ctx, addr) {
		return nil
	}
	if _, err := p.c.Connect(ctx, addr); err != nil {
		if dt != nil {
			// Release this use like a finished request: a pool left behind by the
			// failed Connect (or by an earlier use whose pending disconnect was
			// cancelled above) is then disconnected after the idle timeout
			// instead of being forgotten.
			p.release(ctx, addr)
		}
		return err
	}
	return nil
}

// Release ends an Acquire use of addrs. A request-derived addr (one that is not
// in the configured addrs) is disconnected once it has had no in-flight user
// for the idle timeout, so filter targets named in requests do not accumulate
// connection pools while targets in steady use keep theirs. Configured addrs
// stay connected.
func (p *Pool) Release(ctx context.Context, addrs ...string) {
	for _, addr := range addrs {
		p.release(ctx, addr)
	}
}

func (p *Pool) release(ctx context.Context, addr string) {
	if slices.Contains(p.addrs, addr) {
		return
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	dt, ok := p.dyn[addr]
	if !ok || dt.refs <= 0 {
		return
	}
	dt.refs--
	if dt.refs > 0 || p.closed {
		return
	}
	dt.gen++
	gen := dt.gen
	// The request is over by the time the disconnect runs.
	ctx = context.WithoutCancel(ctx)
	dt.evict = time.AfterFunc(p.idle, func() {
		p.disconnectIdle(ctx, addr, dt, gen)
	})
}

// disconnectIdle disconnects addr unless it was used again, or the pool was
// stopped, after the idle disconnect gen was scheduled.
func (p *Pool) disconnectIdle(ctx context.Context, addr string, dt *dynamicTarget, gen uint64) {
	dt.mu.Lock()
	defer dt.mu.Unlock()
	p.mu.Lock()
	if p.closed || dt.refs > 0 || dt.gen != gen || p.dyn[addr] != dt {
		p.mu.Unlock()
		return
	}
	dt.evict = nil
	p.mu.Unlock()
	if err := p.c.Disconnect(ctx, addr); err != nil &&
		errors.IsNot(err, errors.ErrGRPCClientConnNotFound(addr)) {
		log.Warnf("failed to disconnect request-derived filter target %s: %v", addr, err)
	}
	p.mu.Lock()
	if dt.refs == 0 && dt.gen == gen && p.dyn[addr] == dt {
		delete(p.dyn, addr)
	}
	p.mu.Unlock()
}
