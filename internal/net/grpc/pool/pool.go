//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// Package pool provied grpc connection pooling functionality for better performance.
package pool

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

// Alias types for clarity.
type (
	ClientConn = grpc.ClientConn
	DialOption = grpc.DialOption
)

// Conn defines the interface for a gRPC connection pool.
type Conn interface {
	// Connect establishes connections for all slots.
	Connect(context.Context) (Conn, error)
	// Disconnect gracefully closes all connections in the pool.
	Disconnect(context.Context) error
	// Do executes the provided function using a healthy connection.
	Do(context.Context, func(*ClientConn) error) error
	// Get returns a healthy connection from the pool, if available.
	Get(context.Context) (*ClientConn, bool)
	// IsHealthy checks the overall health of the pool.
	IsHealthy(context.Context) bool
	// IsIPConn indicates whether the pool is using direct IP connections.
	IsIPConn() bool
	// Len returns the number of connection slots.
	Len() uint64
	// Size returns the configured pool size.
	Size() uint64
	// Reconnect re-establishes connections if the pool is unhealthy or if forced.
	Reconnect(context.Context, bool) (Conn, error)
	// String returns a string representation of the pool's state.
	String() string
}

// poolConn wraps a single gRPC connection and its target address.
type poolConn struct {
	conn *ClientConn // Underlying gRPC connection.
	addr string      // Target address used for dialing this connection.
}

// Close gracefully closes the connection with the specified delay.
// It periodically checks the connection state until either the connection is closed or the delay elapses.
func (pc *poolConn) Close(ctx context.Context, delay time.Duration) error {
	// Determine the ticker interval (at least 5ms, at most 5s).
	interval := delay / 10
	if interval < 5*time.Millisecond {
		interval = 5 * time.Millisecond
	} else if interval > time.Minute {
		interval = 5 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Create a context with timeout to ensure closure does not hang indefinitely.
	ctx, cancel := context.WithTimeout(ctx, delay)
	defer cancel()

	log.Debugf("Closing connection for %s with delay %s", pc.addr, delay)
	for {
		switch pc.conn.GetState() {
		case connectivity.Idle, connectivity.Connecting, connectivity.Ready, connectivity.TransientFailure:
			err := pc.conn.Close()
			if err != nil {
				log.Errorf("failed to close gRPC pool connection for %s, error: %v", pc.addr, err)
			}
		case connectivity.Shutdown:
			return nil
		}
		select {
		case <-ctx.Done():
			err := ctx.Err()
			if errors.IsNot(err, context.DeadlineExceeded, context.Canceled) {
				return err
			}
			return nil
		case <-ticker.C:
		}
	}
}

// pool implements the Conn interface.
// It stores connection slots in a lock-free manner using an atomic.Value.
type pool struct {
	// connSlots holds a slice of atomic pointers to poolConn.
	connSlots atomic.Pointer[[]atomic.Pointer[poolConn]] // holds []atomic.Pointer[poolConn]

	// Configuration parameters.
	startPort       uint16 // Starting port for scanning if needed.
	endPort         uint16 // Ending port for scanning if needed.
	host            string // Target host.
	port            uint16 // Target port.
	addr            string // Complete address (host:port).
	isIPAddr        bool   // True if the target is an IP address.
	enableDNSLookup bool   // Whether to perform DNS resolution.
	enableMetrics   bool   // Whether to enable metrics collection.

	// Pool management fields.
	poolSize     atomic.Uint64 // Configured pool size.
	currentIndex atomic.Uint64 // Atomic counter for round-robin indexing.

	// gRPC dial options and timeouts.
	dialOpts          []DialOption
	dialTimeout       time.Duration // Timeout for dialing a connection.
	oldConnCloseDelay time.Duration // Delay before closing old connections.

	// Retry/backoff strategy.
	bo backoff.Backoff

	// Goroutine management.
	errGroup errgroup.Group

	// Used for DNS change detection during reconnection.
	dnsHash atomic.Pointer[string]

	// Flag indicating whether the pool is closing.
	closing atomic.Bool
}

// Default pool size.
const defaultPoolSize = uint64(4)

// Global metrics are stored in a sync.Map (key: address, value: healthy connection count).
var metrics sync.Map[string, uint64]

// New creates and initializes a connection pool for the configured target address and options.
// It applies default and provided options, parses and normalizes the target address (including host/port extraction with fallbacks and optional port scanning),
// allocates and initializes connection slots, and performs an initial dial check to validate reachability. On success it returns the ready pool; on failure it
// returns an error.
func New(ctx context.Context, opts ...Option) (c Conn, err error) {
	p := &pool{
		dialTimeout:       time.Second,
		oldConnCloseDelay: 2 * time.Minute,
		enableDNSLookup:   false,
	}
	// Apply default and user-specified options.
	for _, opt := range append(defaultOptions, opts...) {
		opt(p)
	}

	if p.addr == "" {
		return nil, errors.Errorf("target address is not provided")
	}

	// Initialize the connection slots.
	p.init()
	p.closing.Store(false)

	// Parse the address to extract host and port.
	var isIPv4, isIPv6 bool
	p.host, p.port, _, isIPv4, isIPv6, err = net.Parse(ctx, p.addr)
	p.isIPAddr = isIPv4 || isIPv6
	if err != nil {
		log.Warnf("failed to parse addr %s: %s", p.addr, err)
		// Fallback: split using Cut.
		if p.host == "" {
			p.host, p.port, err = net.SplitHostPort(p.addr)
			if err != nil {
				if host, portStr, ok := strings.Cut(p.addr, ":"); ok {
					p.host = host
					if portNum, err := strconv.ParseUint(portStr, 10, 16); err == nil {
						p.port = uint16(portNum)
					}
				} else {
					p.host = p.addr
				}
			}
		}
		// If port is still zero, attempt port scanning.
		if p.port == 0 {
			var port uint16
			if port, err = p.scanGRPCPort(ctx); err != nil {
				return nil, err
			}
			p.port = port
		}
		p.addr = net.JoinHostPort(p.host, p.port)
	}

	log.Debugf("Initial connection target: %s, host: %s, port: %d, isIP: %t", p.addr, p.host, p.port, p.isIPAddr)
	conn, err := p.dial(ctx, 0, p.addr)
	if err != nil {
		log.Warnf("Initial dial check to %s failed: %v", p.addr, err)
		var port uint16
		if port, err = p.scanGRPCPort(ctx); err != nil {
			return nil, err
		}
		p.port = port
		p.addr = net.JoinHostPort(p.host, p.port)
		log.Debugf("Fallback target: %s, host: %s, port: %d, isIP: %t", p.addr, p.host, p.port, p.isIPAddr)
		conn, err = p.dial(ctx, 0, p.addr)
		if err != nil {
			if conn != nil {
				cerr := conn.Close()
				if cerr != nil {
					return nil, errors.Join(err, cerr)
				}
			}
			return nil, err
		}
	}
	if conn != nil {
		err = conn.Close()
		if err != nil {
			return nil, err
		}
	}

	if p.errGroup == nil {
		p.errGroup = errgroup.Get()
	}

	return p, nil
}

// init initializes the connection slots slice using an atomic.Value.
func (p *pool) init() {
	size := p.Size()
	if size < 1 {
		size = defaultPoolSize
		p.poolSize.Store(size)
	}
	slots := make([]atomic.Pointer[poolConn], size)
	p.connSlots.Store(&slots)
}

// getSlots returns the current connection slots slice.
func (p *pool) getSlots() []atomic.Pointer[poolConn] {
	if v := p.connSlots.Load(); v != nil && len(*v) > 0 {
		return *v
	}
	return nil
}

// grow increases the number of connection slots if the new size is larger.
func (p *pool) grow(newSize uint64) {
	oldSlots := p.getSlots()
	newSlots := make([]atomic.Pointer[poolConn], newSize)
	if oldSlots == nil {
		p.connSlots.Store(&newSlots)
		p.poolSize.Store(newSize)
		return
	}
	currentLen := uint64(len(oldSlots))
	if currentLen >= newSize {
		return
	}
	copy(newSlots, oldSlots)
	p.connSlots.Store(&newSlots)
	p.poolSize.Store(newSize)
}

// load retrieves the poolConn and real index pos at the specified index.
func (p *pool) load(idx uint64) (ridx uint64, pc *poolConn) {
	if idx >= p.poolSize.Load() {
		return 0, nil
	}
	slots := p.getSlots()
	sz := uint64(len(slots))
	if slots != nil && sz != 0 {
		if sz < idx {
			return sz, slots[sz].Load()
		}
		return idx, slots[idx].Load()
	}
	return 0, nil
}

// store sets the poolConn at the specified index.
func (p *pool) store(idx uint64, pc *poolConn) {
	if idx >= p.poolSize.Load() {
		return
	}
	size := p.Size()
	if size <= idx {
		size = max(idx+1, defaultPoolSize)
		p.grow(size)
	}
	slots := p.getSlots()
	if slots == nil {
		slots = make([]atomic.Pointer[poolConn], size)
		slots[idx].Store(pc)
		p.connSlots.Store(&slots)
		return
	}
	slots[idx].Store(pc)
}

// loop iterates over each connection slot and applies the provided function.
func (p *pool) loop(
	ctx context.Context, fn func(ctx context.Context, idx uint64, pc *poolConn) bool,
) (err error) {
	var count uint64
	for idx := range p.poolSize.Load() {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			if errors.IsNot(err, context.DeadlineExceeded, context.Canceled) {
				return err
			}
			return nil
		default:
			count++
			ridx, pc := p.load(idx)
			if !fn(ctx, ridx, pc) {
				return nil
			}
		}
	}
	if count == 0 {
		return errors.ErrGRPCPoolConnectionNotFound
	}
	return nil
}

// slotCount returns the number of connection slots.
func (p *pool) slotCount() uint64 {
	if p == nil {
		return 0
	}
	slots := p.getSlots()
	if slots == nil {
		return 0
	}
	return uint64(len(slots))
}

// flush clears the connection slots.
func (p *pool) flush() {
	slots := make([]atomic.Pointer[poolConn], p.Size())
	p.connSlots.Store(&slots)
}

// refreshConn checks if the connection at slot idx is healthy for the given address.
// If not, it dials a new connection and updates the slot atomically.
// It also schedules graceful closure of any existing (old) connection.
func (p *pool) refreshConn(ctx context.Context, idx uint64, pc *poolConn, addr string) {
	if pc != nil && pc.conn != nil {
		state := pc.conn.GetState()
		if pc.addr == addr {
			switch state {
			case connectivity.Ready,
				connectivity.Connecting:
				return
			case connectivity.Idle:
				pc.conn.Connect()
				return
			case connectivity.TransientFailure:
				log.Debugf("connection for %s pool %d/%d is in TransientFailure, keeping it to allow retry", addr, idx+1, p.Size())
				return
			case connectivity.Shutdown:
				log.Debugf("connection for %s pool %d/%d is Shutdown, will redial", addr, idx+1, p.Size())
			}
		} else {
			log.Debugf("connection for %s pool %d/%d addr mismatch (got %s), will redial", addr, idx+1, p.Size(), pc.addr)
		}
	} else {
		log.Debugf("connection pool %d/%d len %d is empty, establish new pool member connection to %s", idx+1, p.Size(), p.Len(), addr)
	}

	p.errGroup.Go(func() error {
		newConn, err := p.dial(ctx, idx, addr)
		if err != nil {
			if pc != nil && pc.conn != nil {
				state, healthy := p.isHealthy(idx, pc.conn)
				if healthy {
					return nil
				}
				log.Debugf("re-dialed connection for %s pool %d/%d len %d is still unhealthy (state: %s) going to close connection for %s",
					pc.addr, idx+1, p.Size(), p.Len(), state.String(), addr)

				if pc.conn != nil {
					log.Debugf("closing unhealthy connection pool %d/%d len %d for addr: %s", idx+1, p.Size(), p.Len(), pc.addr)
					err := pc.Close(ctx, p.oldConnCloseDelay)
					if err != nil {
						log.Errorf("failed to close connection pool %d/%d addr = %s\terror = %v", idx+1, p.Size(), pc.addr, err)
					}
					return nil
				}
			}

			if errors.IsNot(err, context.DeadlineExceeded, context.Canceled) {
				log.Warnf("failed to refresh connection for pool %d/%d addr: %s: %v", idx+1, p.Size(), addr, err)
			} else {
				log.Debugf("pool connect operation has been canceled for connection pool %d/%d, addr: %s,\terror: %v", idx+1, p.Size(), addr, err)
			}
			return nil // do not propagate dial error in async
		}

		p.store(idx, &poolConn{conn: newConn, addr: addr})

		if pc != nil && pc.conn != nil {
			log.Debugf("closing unhealthy connection pool %d/%d len %d for addr: %s", idx+1, p.Size(), p.Len(), pc.addr)
			err := pc.Close(ctx, p.oldConnCloseDelay)
			if err != nil {
				log.Errorf("failed to close connection pool %d/%d addr = %s\terror = %v", idx+1, p.Size(), pc.addr, err)
			}
		}
		return nil
	})
}

// Connect establishes connections for all slots.
// It uses DNS lookup if enabled; otherwise, it connects to the single target address.
func (p *pool) Connect(ctx context.Context) (Conn, error) {
	if p.closing.Load() {
		return p, nil
	}
	log.Debugf("Connecting: addr=%s, host=%s, port=%d, isIP=%t, enableDNS=%t",
		p.addr, p.host, p.port, p.isIPAddr, p.enableDNSLookup)

	if p.isIPAddr || !p.enableDNSLookup {
		return p.singleTargetConnect(ctx, p.addr)
	}
	ips, err := p.lookupIPAddr(ctx)
	if err != nil {
		return p.singleTargetConnect(ctx, p.addr)
	}
	if len(ips) == 1 {
		return p.singleTargetConnect(ctx, net.JoinHostPort(ips[0], p.port))
	}
	return p.connect(ctx, ips...)
}

// connect establishes connections using multiple IP addresses.
func (p *pool) connect(ctx context.Context, ips ...string) (c Conn, err error) {
	if p == nil || p.closing.Load() {
		return p, nil
	}

	if uint64(len(ips)) > p.Size() {
		p.grow(uint64(len(ips)))
	}
	log.Debugf("Connecting to %s multiple IPs: %v on port %d", p.addr, ips, p.port)
	err = p.loop(ctx, func(ctx context.Context, idx uint64, pc *poolConn) bool {
		p.refreshConn(ctx, idx, pc, net.JoinHostPort(ips[idx%uint64(len(ips))], p.port))
		return true
	})
	if errors.IsNot(err, context.DeadlineExceeded, context.Canceled) {
		return p, err
	}
	hash := strings.Join(ips, "-")
	p.dnsHash.Store(&hash)
	return p, nil
}

// singleTargetConnect connects every slot to a single target address.
func (p *pool) singleTargetConnect(ctx context.Context, addr string) (Conn, error) {
	if p == nil || p.closing.Load() {
		return p, nil
	}
	log.Debugf("Connecting to single target: %s", addr)
	err := p.loop(ctx, func(ctx context.Context, idx uint64, pc *poolConn) bool {
		p.refreshConn(ctx, idx, pc, addr)
		return true
	})
	if errors.IsNot(err, context.DeadlineExceeded, context.Canceled) {
		return p, err
	}
	p.dnsHash.Store(&p.host)
	return p, nil
}

// Reconnect re-establishes connections if the pool is unhealthy or if forced.
func (p *pool) Reconnect(ctx context.Context, force bool) (Conn, error) {
	if p == nil || p.closing.Load() {
		return p, nil
	}
	hash := p.dnsHash.Load()
	if force || hash == nil || *hash == "" {
		return p.Connect(ctx)
	}
	if p.IsHealthy(ctx) {
		if !p.isIPAddr && p.enableDNSLookup {
			if *hash != "" {
				ips, err := p.lookupIPAddr(ctx)
				if err != nil {
					return p, nil
				}
				if len(ips) == 1 {
					return p.singleTargetConnect(ctx, net.JoinHostPort(ips[0], p.port))
				}
				if *hash != strings.Join(ips, "-") {
					return p.connect(ctx, ips...)
				}
			}
		} else {
			return p.singleTargetConnect(ctx, p.addr)
		}
		return p, nil
	}
	return p.Connect(ctx)
}

// Disconnect gracefully closes all connections in the pool.
func (p *pool) Disconnect(ctx context.Context) (err error) {
	log.Warn("Disconnecting pool...")
	p.closing.Store(true)
	defer p.closing.Store(false)
	emap := make(map[string]error, p.Size())
	err = p.loop(ctx, func(ctx context.Context, idx uint64, pc *poolConn) bool {
		if pc != nil && pc.conn != nil {
			log.Debugf("Closing pool connection %d/%d for %s", idx+1, p.Size(), pc.addr)
			if err := pc.Close(ctx, p.oldConnCloseDelay); err != nil {
				if errors.IsNot(err, context.DeadlineExceeded, context.Canceled) {
					log.Errorf("failed to close connection pool %d/%d addr = %s\terror = %v", idx+1, p.Size(), pc.addr, err)
					emap[err.Error()] = err
				} else {
					log.Debugf("Disconnect loop operation canceled while closing connection pool %d/%d addr = %s\terror = %v", idx+1, p.Size(), pc.addr, err)
					return false
				}
			}
		}
		return true
	})
	p.flush()
	for _, e := range emap {
		err = errors.Join(err, e)
	}
	return err
}

// dial creates a new gRPC connection to the specified address.
// It uses a dial timeout and, if configured, a backoff strategy.
func (p *pool) dial(ctx context.Context, idx uint64, addr string) (*ClientConn, error) {
	dialFunc := func() (*ClientConn, error) {
		log.Debugf("Dialing pool connection %d/%d to %s", idx+1, p.Size(), addr)
		conn, err := grpc.NewClient(addr, p.dialOpts...)
		if err != nil {
			if conn != nil {
				return nil, errors.Join(err, conn.Close())
			}
			return nil, err
		}

		// Don't close immediately if unhealthy, check specific states
		_, healthy := p.isHealthy(idx, conn)
		if !healthy {
			if conn != nil {
				// If connection is in TransientFailure, retain and let it retry; otherwise close.
				state := conn.GetState()
				if state == connectivity.TransientFailure {
					return conn, nil
				}
				err = conn.Close()
				if err != nil {
					err = errors.Join(errors.ErrGRPCClientConnNotFound(addr), err)
				} else {
					err = errors.ErrGRPCClientConnNotFound(addr)
				}
			}
			err = errors.Wrapf(err, "pool connection %d/%d for %s is unhealthy", idx+1, p.Size(), addr)
			log.Debug(err)
			return nil, err
		}
		return conn, nil
	}
	if p.bo != nil {
		var conn *ClientConn
		_, err := p.bo.Do(ctx, func(_ context.Context) (any, bool, error) {
			var err error
			conn, err = dialFunc()
			return conn, err != nil, err
		})
		if err != nil && conn != nil {
			return nil, errors.Join(err, conn.Close())
		}
		return conn, nil
	}
	return dialFunc()
}

// getHealthyConn retrieves a healthy connection from the pool using round-robin indexing.
// It attempts up to poolSize times.
func (p *pool) getHealthyConn(ctx context.Context) (pc *poolConn, ok bool) {
	if p == nil || p.closing.Load() {
		return nil, false
	}
	sz := p.Size()
	if sz == 0 {
		return nil, false
	}
	start := p.currentIndex.Add(1)
	for i := range sz {
		idx := (start + i) % sz
		_, pc = p.load(idx)

		if pc != nil && pc.conn != nil {
			state := pc.conn.GetState()
			switch state {
			case connectivity.Ready,
				connectivity.Connecting:
				return pc, true
			case connectivity.Idle:
				// If the connection is idle, trigger a connection attempt but continue searching.
				// If we don't find a Ready one, we might return this one.
				pc.conn.Connect()
				return pc, true
			case connectivity.TransientFailure:
				// Log debug info but continue searching.
				// We prefer Ready/Idle/Connecting.
				log.Debugf("connection for %s pool %d/%d (state: %s) is in transient failure", pc.addr, idx+1, sz, state.String())
			case connectivity.Shutdown:
				// Skip
			}
		}
	}

	rc := make([]uint64, 0, sz) // refreshedConnection list
	// Second pass: if no healthy connection found, try to refresh empty slots or shutdown connections.
	for i := range sz {
		idx := (start + i) % sz
		_, pc = p.load(idx)
		if pc == nil || pc.conn == nil || pc.conn.GetState() == connectivity.Shutdown {
			p.refreshConn(ctx, idx, pc, p.addr)
			rc = append(rc, idx)
		} else if pc.conn.GetState() == connectivity.TransientFailure {
			// If TransientFailure, return it as a last resort.
			return pc, true
		}
	}

	// Third pass return connected pool conn
	for _, idx := range rc {
		_, pc = p.load(idx)
		if pc != nil && pc.conn != nil && pc.conn.GetState() != connectivity.Shutdown {
			return pc, true
		}
	}

	return nil, false
}

// Do executes the provided function using a healthy connection.
// If an error indicating a closed connection is returned, it attempts to refresh the connection and retries.
func (p *pool) Do(ctx context.Context, f func(conn *ClientConn) error) (err error) {
	if p == nil {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	pc, ok := p.getHealthyConn(ctx)
	if !ok || pc == nil || pc.conn == nil {
		return errors.ErrGRPCClientConnNotFound(p.addr)
	}
	return f(pc.conn)
}

// Get returns a healthy connection from the pool, if available.
func (p *pool) Get(ctx context.Context) (conn *ClientConn, ok bool) {
	pc, ok := p.getHealthyConn(ctx)
	if ok && pc != nil {
		return pc.conn, true
	}
	return nil, false
}

// IsHealthy checks the overall health of the pool.
// For IP-based connections, all slots must be healthy; otherwise, at least one healthy slot is acceptable.
// Global metrics are updated accordingly.
func (p *pool) IsHealthy(ctx context.Context) bool {
	healthyCount := uint64(0)
	err := p.loop(ctx, func(ctx context.Context, idx uint64, pc *poolConn) bool {
		if pc != nil && pc.conn != nil {
			state, healthy := p.isHealthy(idx, pc.conn)
			if healthy {
				healthyCount++
				if p.enableMetrics {
					cnt, ok := metrics.Load(pc.addr)
					if ok {
						metrics.Store(pc.addr, cnt+1)
					} else {
						metrics.Store(pc.addr, 1)
					}
				}
			} else {
				log.Debugf("unhealthy connection detected for %s pool %d/%d len %d is unhealthy (state: %s)",
					pc.addr, idx+1, p.Size(), p.Len(), state.String())
			}
		} else {
			// Do not block health check with dial
			p.refreshConn(ctx, idx, pc, p.addr)
			// Assume unhealthy if nil for now
			log.Debugf("nil pool connection detected for %s pool %d/%d len %d is unhealthy", p.addr, idx+1, p.Size(), p.Len())
		}
		return true
	})
	if p.enableMetrics {
		metrics.Store(p.addr, healthyCount)
	}
	if err != nil {
		log.Debugf("health check loop for addr=%s returned error: %v", p.addr, err)
	}
	if healthyCount == 0 {
		log.Warnf("no connection pool member is healthy for addr=%s size=%d, len=%d", p.addr, p.Size(), p.Len())
		return false
	}
	if p.isIPAddr {
		return healthyCount == p.slotCount()
	}
	return healthyCount > 0
}

// Len returns the number of connection slots.
func (p *pool) Len() uint64 {
	return p.slotCount()
}

// Size returns the configured pool size.
func (p *pool) Size() uint64 {
	return p.poolSize.Load()
}

// IsIPConn indicates whether the pool is using direct IP connections.
func (p *pool) IsIPConn() bool {
	return p.isIPAddr
}

// String returns a string representation of the pool's state.
func (p *pool) String() string {
	hash := ""
	if rh := p.dnsHash.Load(); rh != nil {
		hash = *rh
	}
	return fmt.Sprintf("addr: %s, host: %s, port: %d, isIP: %t, enableDNS: %t, dnsHash: %s, slotCount: %d, poolSize: %d, currentIndex: %d, dialTimeout: %s, oldConnCloseDelay: %s, closing: %t",
		p.addr, p.host, p.port, p.isIPAddr, p.enableDNSLookup, hash, p.Len(), p.Size(), p.currentIndex.Load(),
		p.dialTimeout.String(), p.oldConnCloseDelay.String(), p.closing.Load())
}

// lookupIPAddr performs DNS lookup for the host and returns a list of reachable IP addresses.
// It also attempts a short TCP dial ("ping") for each IP.
func (p *pool) lookupIPAddr(ctx context.Context) ([]string, error) {
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, p.host)
	if err != nil {
		log.Debugf("Failed to lookup IP addresses for %s: %s", p.addr, err.Error())
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, errors.ErrGRPCLookupIPAddrNotFound(p.host)
	}
	var ips []string
	for _, ip := range addrs {
		ipStr := ip.String()
		target := net.JoinHostPort(ipStr, p.port)
		pingCtx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		conn, err := net.DialContext(pingCtx, net.TCP.String(), target)
		cancel()
		if err == nil {
			ips = append(ips, ipStr)
		} else {
			log.Warnf("Failed to ping %s: %s", target, err.Error())
		}
		if conn != nil {
			err = conn.Close()
			if errors.IsNot(err, context.DeadlineExceeded, context.Canceled) {
				log.Warn("failed to close connection:", err)
			}
		}
	}
	if len(ips) == 0 {
		return nil, errors.ErrGRPCLookupIPAddrNotFound(p.host)
	}
	// Sorting can be added here if needed.
	return ips, nil
}

// scanGRPCPort scans ports from startPort to endPort for a valid gRPC endpoint.
func (p *pool) scanGRPCPort(ctx context.Context) (port uint16, err error) {
	ports, err := net.ScanPorts(ctx, p.startPort, p.endPort, p.host)
	if err != nil {
		return 0, err
	}
	log.Debugf("Scanning available gRPC ports: %v", ports)
	var conn *ClientConn
	for _, port := range ports {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			if errors.IsNot(err, context.DeadlineExceeded, context.Canceled) {
				return 0, err
			}
			return 0, nil
		default:
			conn, err = grpc.NewClient(net.JoinHostPort(p.host, port), p.dialOpts...)
			if err == nil && conn != nil {
				_, healthy := p.isHealthy(0, conn)
				if healthy {
					log.Debugf("Found valid gRPC port: %d", port)
					err = conn.Close()
					if err != nil {
						log.Warnf("Failed to close connection for port %d: %s", port, err.Error())
					}
					return port, nil
				}
			}
			if conn != nil {
				_ = conn.Close()
			}
		}
	}
	return 0, errors.ErrInvalidGRPCPort(p.addr, p.host, p.port)
}

// Metrics returns a map of healthy connection counts per target address.
func Metrics(_ context.Context) map[string]uint64 {
	result := make(map[string]uint64, metrics.Len())
	metrics.Range(func(addr string, count uint64) bool {
		if addr != "" {
			result[addr] = count
		}
		return true
	})
	if len(result) == 0 {
		return nil
	}
	return result
}

// p.isHealthy checks whether a given gRPC connection is healthy by examining its connectivity state.
func (p *pool) isHealthy(idx uint64, conn *ClientConn) (state connectivity.State, healthy bool) {
	if conn == nil {
		log.Warnf("gRPC target %s's pool connection %d/%d is nil", p.addr, idx+1, p.Size())
		return connectivity.State(-1), false
	}
	state = conn.GetState()
	switch state {
	case connectivity.Ready,
		connectivity.Connecting:
		return state, true
	case connectivity.Idle:
		log.Debugf("gRPC target %s's pool connection %d/%d status is Idle\tstate: %s\twill re-connect...", conn.Target(), idx+1, p.Size(), state.String())
		conn.Connect()
		return state, true
	case connectivity.TransientFailure:
		return state, true
	case connectivity.Shutdown:
		log.Debugf("gRPC target %s's pool connection %d/%d is unhealthy (state: %s)", conn.Target(), idx+1, p.Size(), state.String())
		return state, false
	default:
		log.Errorf("gRPC target %s's pool connection %d/%d has unknown state: %s", conn.Target(), idx+1, p.Size(), state.String())
		return state, false
	}
}
