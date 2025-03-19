//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

// Package pool provides a lock-free gRPC connection pool client.
// This re-implementation maintains the public Conn interface unchanged while
// using atomic operations for efficient, lock-free connection management.
// Additional features such as DNS lookup, port scanning, and metrics collection are incorporated.
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
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
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
		err := pc.conn.Close()
		if err != nil && !status.Is(err, codes.Canceled) {
			return err
		}
		select {
		case <-ctx.Done():
			if ctx.Err() != nil &&
				!errors.Is(ctx.Err(), context.DeadlineExceeded) &&
				!errors.Is(ctx.Err(), context.Canceled) {
				return ctx.Err()
			}
			return nil
		case <-ticker.C:
			switch pc.conn.GetState() {
			case connectivity.Idle, connectivity.Connecting, connectivity.TransientFailure:
				err := pc.conn.Close()
				if err != nil && !status.Is(err, codes.Canceled) {
					return err
				}
				return nil
			case connectivity.Shutdown:
				return nil
			}
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

// New creates a new connection pool with the provided options.
// It parses the target address, initializes the connection slots, and performs an initial dial check.
func New(ctx context.Context, opts ...Option) (Conn, error) {
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
	var err error
	var isIPv4, isIPv6 bool
	p.host, p.port, _, isIPv4, isIPv6, err = net.Parse(p.addr)
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
	// Perform an initial dial check.
	conn, err := p.dial(ctx, p.addr)
	if err != nil {
		log.Warnf("Initial dial check to %s failed: %v", p.addr, err)
		var port uint16
		if port, err = p.scanGRPCPort(ctx); err != nil {
			return nil, err
		}
		p.port = port
		p.addr = net.JoinHostPort(p.host, p.port)
		log.Debugf("Fallback target: %s, host: %s, port: %d, isIP: %t", p.addr, p.host, p.port, p.isIPAddr)
		conn, err = p.dial(ctx, p.addr)
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
func (p *pool) getSlots() *[]atomic.Pointer[poolConn] {
	if v := p.connSlots.Load(); v != nil {
		return v
	}
	return nil
}

// grow increases the number of connection slots if the new size is larger.
func (p *pool) grow(newSize uint64) {
	oldSlots := *p.getSlots()
	currentLen := uint64(len(oldSlots))
	if currentLen >= newSize {
		return
	}
	newSlots := make([]atomic.Pointer[poolConn], newSize)
	copy(newSlots, oldSlots)
	p.connSlots.Store(&newSlots)
	p.poolSize.Store(newSize)
}

// load retrieves the poolConn at the specified index.
func (p *pool) load(idx uint64) *poolConn {
	slots := *p.getSlots()
	if slots == nil || idx < 0 || idx >= p.slotCount() {
		return nil
	}
	return slots[idx].Load()
}

// store sets the poolConn at the specified index.
func (p *pool) store(idx uint64, pc *poolConn) {
	slots := *p.getSlots()
	if slots == nil || idx < 0 || idx >= p.slotCount() {
		return
	}
	slots[idx].Store(pc)
}

// loop iterates over each connection slot and applies the provided function.
func (p *pool) loop(
	ctx context.Context, fn func(ctx context.Context, idx uint64, pc *poolConn) bool,
) error {
	slots := *p.getSlots()
	if slots == nil {
		return errors.Errorf("connection slots not initialized")
	}
	var count uint64
	for idx := range slots {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			count++
			if !fn(ctx, uint64(idx), slots[idx].Load()) {
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
	return uint64(len(*p.getSlots()))
}

// flush clears the connection slots.
func (p *pool) flush() {
	p.connSlots.Store(nil)
}

// refreshConn checks if the connection at slot idx is healthy for the given address.
// If not, it dials a new connection and updates the slot atomically.
// It also schedules graceful closure of any existing (old) connection.
func (p *pool) refreshConn(ctx context.Context, idx uint64, pc *poolConn, addr string) error {
	if pc != nil && pc.addr == addr && p.isHealthy(ctx, pc.conn) {
		return nil
	}
	if pc != nil {
		log.Debugf("connection for %s pool %d/%d is unhealthy trying to establish new pool member connection to %s", pc.addr, idx+1, p.Size(), addr)
	} else {
		log.Debugf("connection pool %d/%d is empty, establish new pool member connection to %s", idx+1, p.Size(), addr)
	}
	newConn, err := p.dial(ctx, addr)
	if err != nil {
		if pc != nil && p.isHealthy(ctx, pc.conn) {
			return nil
		}
		if pc != nil && pc.conn != nil {
			p.errGroup.Go(func() error {
				log.Debugf("waiting for invalid connection to %s to be closed...", pc.addr)
				err := pc.Close(ctx, p.oldConnCloseDelay)
				if err != nil {
					log.Debugf("failed to close connection pool addr = %s\terror = %v", pc.addr, err)
				}
				return nil
			})
		}
		return errors.Join(err, errors.ErrInvalidGRPCClientConn(addr))
	}

	p.store(idx, &poolConn{conn: newConn, addr: addr})

	if pc != nil && pc.conn != nil {
		p.errGroup.Go(func() error {
			log.Debugf("waiting for old connection to %s to be closed...", pc.addr)
			err := pc.Close(ctx, p.oldConnCloseDelay)
			if err != nil {
				log.Debugf("failed to close connection pool addr = %s\terror = %v", pc.addr, err)
			}
			return nil
		})
	}
	return nil
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
		target := net.JoinHostPort(ips[0], p.port)
		return p.singleTargetConnect(ctx, target)
	}
	return p.connect(ctx, ips...)
}

// connect establishes connections using multiple IP addresses.
func (p *pool) connect(ctx context.Context, ips ...string) (c Conn, err error) {
	if uint64(len(ips)) > p.Size() {
		p.grow(uint64(len(ips)))
	}
	log.Debugf("Connecting to multiple IPs: %v on port %d", ips, p.port)
	err = p.loop(ctx, func(ctx context.Context, idx uint64, pc *poolConn) bool {
		target := net.JoinHostPort(ips[idx%uint64(len(ips))], p.port)
		if err = p.refreshConn(ctx, idx, pc, target); err != nil {
			if !errors.Is(err, context.DeadlineExceeded) &&
				!errors.Is(err, context.Canceled) {
				log.Warnf("An error occurred while dialing pool slot %d connection to %s,\terror: %v", idx, target, err)
			} else {
				log.Debugf("Connect loop operation canceled while dialing pool slot %d connection to %s,\terror: %v", idx, target, err)
				return false
			}
		}
		return true
	})
	if err != nil && !errors.Is(err, context.Canceled) &&
		!errors.Is(err, context.DeadlineExceeded) {
		return p, err
	}
	hash := strings.Join(ips, "-")
	p.dnsHash.Store(&hash)
	return p, err
}

// singleTargetConnect connects every slot to a single target address.
func (p *pool) singleTargetConnect(ctx context.Context, addr string) (Conn, error) {
	log.Debugf("Connecting to single target: %s", addr)
	failCount := uint64(0)
	err := p.loop(ctx, func(ctx context.Context, idx uint64, pc *poolConn) bool {
		if err := p.refreshConn(ctx, idx, pc, addr); err != nil {
			if !errors.Is(err, context.DeadlineExceeded) &&
				!errors.Is(err, context.Canceled) {
				log.Warnf("An error occurred while dialing pool member connection to %s,\terror: %v", p.addr, err)
				failCount++
				if p.isIPAddr && (p.slotCount() <= 2 || failCount >= p.slotCount()/3) {
					return false
				}
			} else {
				log.Debugf("Connect loop operation canceled while dialing pool member connection to %s,\terror: %v", p.addr, err)
				return false
			}
		}
		return true
	})
	p.dnsHash.Store(&p.host)
	return p, err
}

// Reconnect re-establishes connections if the pool is unhealthy or if forced.
func (p *pool) Reconnect(ctx context.Context, force bool) (Conn, error) {
	if p.closing.Load() {
		return p, nil
	}
	hash := p.dnsHash.Load()
	if force || hash == nil || *hash == "" {
		return p.Connect(ctx)
	}
	if p.IsHealthy(ctx) {
		if !p.isIPAddr && p.enableDNSLookup && hash != nil && *hash != "" {
			ips, err := p.lookupIPAddr(ctx)
			if err != nil {
				return p, nil
			}
			if len(ips) == 1 {
				target := net.JoinHostPort(ips[0], p.port)
				return p.singleTargetConnect(ctx, target)
			}
			if *hash != strings.Join(ips, "-") {
				return p.connect(ctx, ips...)
			}
		}
		return p, nil
	}
	return p.Connect(ctx)
}

// Disconnect gracefully closes all connections in the pool.
func (p *pool) Disconnect(ctx context.Context) (err error) {
	log.Debug("Disconnecting pool...")
	p.closing.Store(true)
	defer p.closing.Store(false)
	emap := make(map[string]error, p.Size())
	err = p.loop(ctx, func(ctx context.Context, idx uint64, pc *poolConn) bool {
		if pc != nil && pc.conn != nil {
			log.Debugf("Closing slot %d (addr: %s)", idx, pc.addr)
			if err := pc.Close(ctx, p.oldConnCloseDelay); err != nil {
				if !errors.Is(err, context.DeadlineExceeded) &&
					!errors.Is(err, context.Canceled) {
					log.Debugf("failed to close connection pool addr = %s\terror = %v", pc.addr, err)
					emap[err.Error()] = err
				} else {
					log.Debugf("Disconnect loop operation canceled while closing pool member connection to %s,\terror: %v", pc.addr, err)
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
func (p *pool) dial(ctx context.Context, addr string) (*ClientConn, error) {
	dialFunc := func(ctx context.Context) (*ClientConn, error) {
		ctx, cancel := context.WithTimeout(ctx, p.dialTimeout)
		defer cancel()
		log.Debugf("Dialing %s with timeout %s", addr, p.dialTimeout)
		conn, err := grpc.NewClient(addr, p.dialOpts...)
		if err != nil {
			if conn != nil {
				_ = conn.Close()
			}
			return nil, err
		}
		if !p.isHealthy(ctx, conn) {
			if conn != nil {
				err = conn.Close()
				if err != nil {
					err = errors.Join(errors.ErrGRPCClientConnNotFound(addr), err)
				} else {
					err = errors.ErrGRPCClientConnNotFound(addr)
				}
			}
			log.Debugf("connection for %s is unhealthy: %v", addr, err)
			return nil, errors.Wrapf(err, "connection to %s is unhealthy", addr)
		}
		return conn, nil
	}
	if p.bo != nil {
		var conn *ClientConn
		_, err := p.bo.Do(ctx, func(ctx context.Context) (interface{}, bool, error) {
			var err error
			conn, err = dialFunc(ctx)
			return conn, err != nil, err
		})
		if err != nil && conn != nil {
			_ = conn.Close()
			return nil, errors.Join(err, conn.Close())
		}
		return conn, nil
	}
	return dialFunc(ctx)
}

// getHealthyConn retrieves a healthy connection from the pool using round-robin indexing.
// It attempts up to poolSize times.
func (p *pool) getHealthyConn(ctx context.Context) (idx uint64, pc *poolConn, ok bool) {
	if p == nil || p.closing.Load() {
		return 0, nil, false
	}
	sz := p.Size()
	if sz == 0 {
		return 0, nil, false
	}
	for i := uint64(0); i < sz; i++ {
		idx = p.currentIndex.Add(1) % sz
		pc = p.load(idx)
		if pc != nil && p.isHealthy(ctx, pc.conn) {
			return idx, pc, true
		}
		if err := p.refreshConn(ctx, idx, pc, p.addr); err == nil {
			if pc = p.load(idx); pc != nil && p.isHealthy(ctx, pc.conn) {
				return idx, pc, true
			}
		}
	}
	return 0, nil, false
}

// Do executes the provided function using a healthy connection.
// If an error indicating a closed connection is returned, it attempts to refresh the connection and retries.
func (p *pool) Do(ctx context.Context, f func(conn *ClientConn) error) error {
	if p == nil {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	_, pc, ok := p.getHealthyConn(ctx)
	if !ok || pc == nil || pc.conn == nil {
		return errors.ErrGRPCClientConnNotFound(p.addr)
	}
	return f(pc.conn)
}

// Get returns a healthy connection from the pool, if available.
func (p *pool) Get(ctx context.Context) (conn *ClientConn, ok bool) {
	_, pc, ok := p.getHealthyConn(ctx)
	if ok && pc != nil {
		return pc.conn, true
	}
	return nil, false
}

// IsHealthy checks the overall health of the pool.
// For IP-based connections, all slots must be healthy; otherwise, at least one healthy slot is acceptable.
// Global metrics are updated accordingly.
func (p *pool) IsHealthy(ctx context.Context) bool {
	sz := p.slotCount()
	if sz == 0 {
		return false
	}
	healthyCount := uint64(0)
	err := p.loop(ctx, func(ctx context.Context, _ uint64, pc *poolConn) bool {
		if pc != nil && p.isHealthy(ctx, pc.conn) {
			healthyCount++
		}
		return true
	})
	metrics.Store(p.addr, healthyCount)
	if err != nil {
		log.Debugf("health check loop for addr=%s returned error: %v", p.addr, err)
	}
	if healthyCount == 0 {
		log.Debugf("no connection pool member is healthy for addr=%s", p.addr)
		return false
	}
	if p.isIPAddr {
		return healthyCount == uint64(sz)
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
		p.addr, p.host, p.port, p.isIPAddr, p.enableDNSLookup, hash, p.slotCount(), p.Size(), p.currentIndex.Load(),
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
			if err != nil && !errors.Is(err, context.Canceled) {
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
	log.Debugf("Scanning ports: %v", ports)
	var conn *ClientConn
	for _, port := range ports {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			conn, err = grpc.NewClient(net.JoinHostPort(p.host, port), p.dialOpts...)
			if err == nil && p.isHealthy(ctx, conn) {
				_ = conn.Close()
				log.Debugf("Found valid gRPC port: %d", port)
				return port, nil
			}
			if conn != nil {
				_ = conn.Close()
			}
		}
	}
	return 0, errors.ErrInvalidGRPCPort(p.addr, p.host, p.port)
}

// Metrics returns a map of healthy connection counts per target address.
func Metrics(ctx context.Context) map[string]uint64 {
	result := make(map[string]uint64)
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
func (p *pool) isHealthy(ctx context.Context, conn *ClientConn) bool {
	if conn == nil {
		log.Warn("gRPC connection is nil")
		return false
	}
	state := conn.GetState()
	switch state {
	case connectivity.Ready:
		return true
	case connectivity.Connecting:
		return true
	case connectivity.Idle:
		// Trigger connection if idle.
		p.errGroup.Go(func() error {
			conn.Connect()
			return nil
		})
		if conn.WaitForStateChange(ctx, state) {
			return p.isHealthy(ctx, conn)
		}
		log.Errorf("Connection %s did not recover from idle", conn.Target())
		return false
	case connectivity.Shutdown, connectivity.TransientFailure:
		log.Errorf("Connection %s is unhealthy (state: %s)", conn.Target(), state.String())
		return false
	default:
		log.Errorf("Connection %s has unknown state: %s", conn.Target(), state.String())
		return false
	}
}
