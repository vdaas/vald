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

package grpc

import (
	"context"
	"maps"
	"math"
	"slices"
	"time"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/credentials"
	_ "github.com/vdaas/vald/internal/net/grpc/logger"
	"github.com/vdaas/vald/internal/net/grpc/pool"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/sync/singleflight"
	"google.golang.org/grpc"
	gbackoff "google.golang.org/grpc/backoff"
	gstatus "google.golang.org/grpc/status"
)

type (
	// CallOption wraps grpc.CallOption for convenience.
	CallOption = grpc.CallOption
	// DialOption wraps pool.DialOption for convenience.
	DialOption = pool.DialOption
	// ClientConn wraps pool.ClientConn for convenience.
	ClientConn = pool.ClientConn
)

// NewClient creates a new gRPC ClientConn for the given target and dial options.
func NewClient(target string, opts ...DialOption) (*ClientConn, error) {
	return grpc.NewClient(target, opts...)
}

// WithTransportCredentials returns a DialOption that configures a connection level security credentials.
func WithTransportCredentials(creds credentials.TransportCredentials) DialOption {
	return grpc.WithTransportCredentials(creds)
}

// WaitForReady returns a CallOption that configures whether an RPC waits for a ready connection.
func WaitForReady(waitForReady bool) CallOption {
	return grpc.WaitForReady(waitForReady)
}

// Client is an interface for a gRPC client that manages connections to multiple servers.
type Client interface {
	// StartConnectionMonitor starts a background goroutine to monitor and maintain connections.
	// It returns a channel that reports errors during connection management.
	StartConnectionMonitor(ctx context.Context) (<-chan error, error)
	// Connect establishes a connection to the specified address.
	Connect(ctx context.Context, addr string, dopts ...DialOption) (pool.Conn, error)
	// IsConnected checks if a connection to the specified address is healthy.
	IsConnected(ctx context.Context, addr string) bool
	// Disconnect closes the connection to the specified address.
	Disconnect(ctx context.Context, addr string) error
	// Range iterates over all healthy connections and executes the given function.
	Range(ctx context.Context,
		f func(ctx context.Context,
			addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	// RangeConcurrent iterates over all healthy connections and executes the given function concurrently.
	RangeConcurrent(ctx context.Context,
		concurrency int,
		f func(ctx context.Context,
			addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	// OrderedRange iterates over connections in a specified order and executes the given function.
	OrderedRange(ctx context.Context,
		order []string,
		f func(ctx context.Context,
			addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	// OrderedRangeConcurrent iterates over connections in a specified order and executes the given function concurrently.
	OrderedRangeConcurrent(ctx context.Context,
		order []string,
		concurrency int,
		f func(ctx context.Context,
			addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	// RangeSliceConcurrent executes f concurrently over count addresses starting at offset in addrs.
	RangeSliceConcurrent(ctx context.Context,
		addrs []string,
		offset int,
		count int,
		concurrency int,
		f func(ctx context.Context,
			addr string,
			conn *ClientConn,
			copts ...CallOption) error) error
	// Do executes a function on a specific connection.
	Do(ctx context.Context, addr string,
		f func(ctx context.Context,
			conn *ClientConn,
			copts ...CallOption) (any, error)) (any, error)
	// RoundRobin executes a function on connections in a round-robin manner until one succeeds.
	RoundRobin(ctx context.Context, f func(ctx context.Context,
		conn *ClientConn,
		copts ...CallOption) (any, error)) (any, error)
	// GetDialOption returns the default dial options for the client.
	GetDialOption() []DialOption
	// GetCallOption returns the default call options for the client.
	GetCallOption() []CallOption
	// GetBackoff returns the backoff strategy used for retries.
	GetBackoff() backoff.Backoff
	// SetDisableResolveDNSAddr sets whether to disable DNS resolution for a specific address.
	SetDisableResolveDNSAddr(addr string, disabled bool)
	// ConnectedAddrs returns a slice of currently connected addresses.
	ConnectedAddrs(context.Context) []string
	// Close closes all connections and stops the connection monitor.
	Close(ctx context.Context) error
}

type cbResolver interface {
	Resolve(key string) (*circuitbreaker.Handle, error)
}

type boRetryer interface {
	RetryVoid(ctx context.Context, firstErr error, f func(ctx context.Context) (retryable bool, err error)) error
}

type breakerKeyEntry struct {
	h    *circuitbreaker.Handle
	name string
}

// gRPCClient is an implementation of the Client interface.
type gRPCClient struct {
	// dialer is the custom dialer interface for establishing network connections.
	dialer net.Dialer
	// eg is the error group for managing the synchronization of background goroutines.
	eg errgroup.Group
	// group is singlefight group to prevent duplicated connection control.
	group singleflight.Group[pool.Conn]
	// cb is the circuit breaker used to prevent cascading failures.
	cb circuitbreaker.CircuitBreaker
	// cbResolver caches the optional *circuitbreaker.Handle resolver for fast-path RPCs.
	cbResolver cbResolver
	// bo is the backoff strategy used for retrying failed operations.
	bo backoff.Backoff
	// boRetryer caches the optional RetryVoid implementation for fast-path RPCs.
	boRetryer boRetryer
	// copts is the list of default call options applied to RPC calls.
	copts []CallOption
	// name is the logical name of the gRPC client, used for logging and tracing.
	name string
	// addrs is a set of initial target addresses for the gRPC client.
	addrs map[string]struct{}
	// conns stores the active connection pools for each target address.
	conns sync.Map[string, pool.Conn]
	// targetsSnap holds an immutable snapshot of active connection targets, updated under snapMu.
	targetsSnap atomic.Pointer[connSnapshot]
	// crl is the connection request list used to manage pending reconnection requests.
	crl sync.Map[string, bool]
	// disableResolveDNSAddrs stores addresses for which DNS resolution should be explicitly disabled.
	disableResolveDNSAddrs sync.Map[string, bool]
	// breakerKeys caches (method, addr) -> method+"/"+addr breaker keys and handles to eliminate per-RPC string allocations.
	breakerKeys sync.Map[breakerKeyPair, breakerKeyEntry]
	// roccd represents the duration to wait before closing old connections during reconnection (Reconnection Old Connection Closing Duration).
	roccd string
	// dopts is the list of default dial options applied to new connections.
	dopts []DialOption
	// stopMonitor is the cancel function to stop the connection monitor goroutine.
	// It is guarded by monitorMu.
	stopMonitor context.CancelFunc
	// ech is the channel for reporting errors from the connection monitor.
	// It is guarded by monitorMu.
	ech <-chan error
	// workers is placed last among pointer-bearing fields (after cold fields such as
	// dopts/stopMonitor/ech) so gRPCClient minimizes GC pointerBytes (govet
	// fieldalignment) and keeps hot read-only fields on separate cache lines from
	// workers' mutable shards.
	workers errgroup.WorkerPool
	// gbo is the configuration for gRPC specific backoff strategy.
	gbo gbackoff.Config
	// hcDur is the interval duration for periodic health checks of connections.
	hcDur time.Duration
	// mcd is the minimum duration for connection timeouts (Minimum Connection Duration).
	mcd time.Duration
	// prDur is the interval duration for periodic pool rebalancing.
	prDur time.Duration
	// clientCount tracks the total number of active connections managed by this client.
	clientCount uint64
	// poolSize specifies the number of connections per pool (per address).
	poolSize uint64
	// disconnectGen is incremented whenever a connection pool is removed or replaced in conns.
	disconnectGen atomic.Uint64
	_             [120]byte
	// rrIdx selects the starting target index in RoundRobin so requests distribute evenly across targets.
	rrIdx atomic.Uint64
	// monitorRunning indicates whether the connection monitor goroutine is currently active.
	monitorRunning atomic.Bool
	// snapMu serializes updates to targetsSnap on Connect/Disconnect/storeConn.
	snapMu sync.Mutex
	// monitorMu guards stopMonitor and ech.
	monitorMu sync.Mutex
	// startMu serializes StartConnectionMonitor so that only one caller performs the
	// initial connect at a time. Callers waiting behind a successful start return the
	// running monitor's ech; if the start failed, monitorRunning is reset and each
	// waiting caller re-runs the initial connect itself in turn.
	startMu sync.Mutex
	// resolveDNS indicates whether to resolve DNS for target addresses by default.
	resolveDNS bool
	// enablePoolRebalance indicates whether periodic rebalancing of the connection pool is enabled.
	enablePoolRebalance bool
	// enablePoolMetrics enables healthy connection count metrics of connection pools.
	enablePoolMetrics bool
}

type breakerKeyPair struct {
	method string
	addr   string
}

type connTarget struct {
	p      pool.Conn
	addrBr atomic.Pointer[circuitbreaker.Handle]
	addr   string
}

type connSnapshot struct {
	byAddr  map[string]*connTarget
	targets []*connTarget
	gen     uint64
}

type fanOutMode uint8

const (
	fanOutModeTargets fanOutMode = iota
	fanOutModeOrdered
	fanOutModeSlice
)

// fanOutBatch is the pooled carrier for a concurrent fan-out. It borrows an
// immutable target snapshot (or order slice) and exposes a pre-bound runFn for WorkerPool.FanOut.
type fanOutBatch struct {
	g           *gRPCClient
	ctx         context.Context //nolint:containedctx // request-scoped batch carrier bounded to synchronous FanOut
	f           func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
	runFn       func(int)
	snapTargets []*connTarget
	snapByAddr  map[string]*connTarget
	prefix      string
	orders      []string
	connGen     uint64
	offset      int
	last        int
	mode        fanOutMode
}

// Cold misses are constructed by getFanOutBatch to avoid an initializer cycle.
var fanOutBatchPool sync.Pool[*fanOutBatch] //nolint:gochecknoglobals // batches are reused across fan-out calls

func newFanOutBatch() *fanOutBatch {
	b := &fanOutBatch{}
	b.runFn = b.runAt
	return b
}

func getFanOutBatch() *fanOutBatch {
	if b := fanOutBatchPool.Get(); b != nil {
		return b
	}
	return newFanOutBatch()
}

func putFanOutBatch(b *fanOutBatch) {
	b.g = nil
	b.ctx = nil
	b.f = nil
	b.snapTargets = nil
	b.snapByAddr = nil
	b.prefix = ""
	b.orders = nil
	b.connGen = 0
	b.offset = 0
	b.last = 0
	b.mode = fanOutModeTargets
	fanOutBatchPool.Put(b)
}

func (c *fanOutBatch) runAt(i int) {
	defer func() {
		if r := recover(); r != nil {
			switch {
			case i != c.last:
				log.Errorf("recovered from panic in grpc fan-out task: %v", r)
			case c.mode == fanOutModeOrdered:
				log.Errorf("recovered from panic in grpc OrderedRangeConcurrent: %v", r)
			case c.mode == fanOutModeSlice:
				log.Errorf("recovered from panic in grpc RangeSliceConcurrent: %v", r)
			default:
				log.Errorf("recovered from panic in grpc RangeConcurrent: %v", r)
			}
		}
	}()
	var (
		t    *connTarget
		addr string
	)
	switch c.mode {
	case fanOutModeTargets:
		t = c.snapTargets[i]
		addr = t.addr
	case fanOutModeOrdered:
		addr = c.orders[i]
		if c.snapByAddr != nil {
			t = c.snapByAddr[addr]
		}
	case fanOutModeSlice:
		idx := c.offset + i
		if idx >= len(c.orders) {
			idx -= len(c.orders)
		}
		addr = c.orders[idx]
		if c.snapByAddr != nil {
			t = c.snapByAddr[addr]
		}
	}
	if t == nil || t.p == nil {
		c.g.crl.Store(addr, true)
		op := "OrderedRangeConcurrent"
		if c.mode == fanOutModeSlice {
			op = "RangeSliceConcurrent"
		}
		log.Warnf("gRPCClient.%s operation failed, gRPC connection pool for %s is invalid,\terror: %v", op, addr, errors.ErrGRPCClientConnNotFound(addr))
		return
	}
	ssctx, sspan := startAddrSpan(c.ctx, c.prefix, addr)
	err := c.g.doVoidPolicyTarget(ssctx, t, c.connGen, sspan != nil, c.f)
	if err != nil && sspan != nil {
		recordSpanError(sspan, err)
	}
	if sspan != nil {
		sspan.End()
	}
}

const (
	apiName                    = "vald/internal/net/grpc"
	defaultHealthCheckDuration = 10 * time.Second
)

// New creates a new gRPC client with the given options.
func New(name string, opts ...Option) (c Client) {
	g := &gRPCClient{
		name:  name,
		addrs: make(map[string]struct{}),
		group: singleflight.New[pool.Conn](),
	}

	for _, opt := range append(defaultOptions, opts...) {
		opt(g)
	}
	if g.cb != nil {
		if r, ok := g.cb.(cbResolver); ok {
			g.cbResolver = r
		}
	}
	if g.bo != nil {
		if r, ok := g.bo.(boRetryer); ok {
			g.boRetryer = r
		}
	}
	g.dopts = append(g.dopts, grpc.WithConnectParams(
		grpc.ConnectParams{
			Backoff: gbackoff.Config{
				MaxDelay:   g.gbo.MaxDelay,
				BaseDelay:  g.gbo.BaseDelay,
				Multiplier: g.gbo.Multiplier,
				Jitter:     g.gbo.Jitter,
			},
			MinConnectTimeout: g.mcd,
		},
	))
	if len(g.copts) != 0 {
		g.dopts = append(g.dopts, grpc.WithDefaultCallOptions(g.copts...))
	}
	g.monitorRunning.Store(false)
	return g
}

// StartConnectionMonitor starts a background goroutine to monitor connections,
// perform health checks, and handle reconnections.
func (g *gRPCClient) StartConnectionMonitor(ctx context.Context) (<-chan error, error) {
	g.startMu.Lock()
	defer g.startMu.Unlock()
	if g.monitorRunning.Load() {
		g.monitorMu.Lock()
		ech := g.ech
		g.monitorMu.Unlock()
		return ech, nil
	}
	g.monitorRunning.Store(true)

	addrs := slices.Collect(maps.Keys(g.addrs))
	slices.Sort(addrs)
	if g.dialer != nil {
		g.dialer.StartDialerCache(ctx)
	}

	log.Debugf("gRPC %s connection monitor started for %v", g.name, addrs)

	ech := make(chan error, len(addrs))
	for _, addr := range addrs {
		if addr != "" {
			_, err := g.Connect(ctx, addr)
			if err != nil {
				if errors.IsNot(err, context.Canceled,
					context.DeadlineExceeded,
					errors.ErrCircuitBreakerOpenState,
					errors.ErrGRPCClientConnNotFound("*"),
					errors.ErrGRPCClientConnNotFound(addr),
					errors.ErrGRPCClientNotFound) {
					log.Errorf("failed to initial gRPC %s connection to %s,\terror: %v", g.name, addr, err)
					ech <- err
				} else {
					log.Warn(err)
				}
			}
		}
	}

	if len(addrs) != 0 && atomic.LoadUint64(&g.clientCount) == 0 {
		err := errors.ErrGRPCClientConnNotFound(strings.Join(addrs, ",\t"))
		log.Error(err)
		// the monitor has not been started, allow the next call to retry.
		g.monitorRunning.Store(false)
		return nil, err
	}

	log.Debugf("initial connection succeeded for gRPC %s addrs = %v, client_count: %d", g.name, addrs, atomic.LoadUint64(&g.clientCount))

	var stop context.CancelFunc
	ctx, stop = context.WithCancel(ctx)
	g.monitorMu.Lock()
	g.stopMonitor = stop
	g.monitorMu.Unlock()
	g.eg.Go(safety.RecoverFunc(func() (err error) {
		defer g.monitorRunning.Store(false)
		defer close(ech)
		defer func() {
			if err := g.Close(context.Background()); err != nil {
				log.Error(err)
			}
		}()

		var hcTick, prTick *time.Ticker
		var prCh <-chan time.Time

		if g.hcDur.Nanoseconds() <= 0 {
			g.hcDur = defaultHealthCheckDuration
		}

		err = safety.RecoverFunc(func() error {
			hcTick = time.NewTicker(g.hcDur) // health check ticker
			return nil
		})()
		if err != nil || hcTick == nil {
			select {
			case <-ctx.Done():
				cerr := ctx.Err()
				if errors.IsNot(cerr, context.Canceled, context.DeadlineExceeded) {
					return errors.Join(err, cerr)
				}
			case ech <- err:
				return err
			}
		}
		defer hcTick.Stop()

		// this duration is for timeout to prevent blocking health check loop, which should be minimum duration of hcDur and prDur
		reconnLimitDuration := g.hcDur

		if g.enablePoolRebalance && g.prDur.Nanoseconds() > 0 {
			err = safety.RecoverFunc(func() error {
				prTick = time.NewTicker(g.prDur) // pool rebalance ticker
				return nil
			})()
			if err != nil || prTick == nil {
				select {
				case <-ctx.Done():
					cerr := ctx.Err()
					if errors.IsNot(cerr, context.Canceled, context.DeadlineExceeded) {
						return errors.Join(err, cerr)
					}
				case ech <- err:
					return err
				}
			}
			defer prTick.Stop()
			prCh = prTick.C
			reconnLimitDuration = time.Duration(int64(math.Min(float64(g.hcDur.Nanoseconds()), float64(g.prDur.Nanoseconds()))))
		}

		disconnectTargets := make([]string, 0, len(addrs))
		log.Debugf("connection monitor loop starting for gRPC %s addrs = %v", g.name, addrs)
		for {
			select {
			case <-ctx.Done():
				if err != nil {
					return errors.Join(ctx.Err(), err)
				}
				return ctx.Err()
			case <-prCh:
				err = g.rangeConns("pool rebalance", func(addr string, p pool.Conn) bool {
					// if addr or pool is nil or empty the registration of conns is invalid let's disconnect them
					if addr == "" || p == nil {
						disconnectTargets = append(disconnectTargets, addr)
						return true
					}
					log.Debugf("rebalancing pool connection for gRPC %s addr: %s, detail: %s", g.name, addr, p.String())
					var perr error
					// for rebalancing connection we don't need to check connection health
					p, perr = p.Connect(ctx)
					if errors.IsNot(perr, context.Canceled,
						context.DeadlineExceeded,
						errors.ErrCircuitBreakerOpenState,
						errors.ErrGRPCClientConnNotFound("*"),
						errors.ErrGRPCClientConnNotFound(addr),
						errors.ErrGRPCClientNotFound) {
						select {
						case <-ctx.Done():
							cerr := ctx.Err()
							if errors.IsNot(cerr, context.Canceled, context.DeadlineExceeded) {
								log.Error(errors.Join(perr, cerr))
							}
						case ech <- perr:
						}
					}
					// if rebalanced connection pool is nil even error is nil we should disconnect and delete it
					if perr == nil && p == nil {
						disconnectTargets = append(disconnectTargets, addr)
						return true
					}
					// if connection pool could not recover we should try next connection loop
					if perr != nil || p == nil || !p.IsHealthy(ctx) {
						_, _ = g.crl.LoadOrStore(addr, false)
						return true
					}
					// p is the pool already held in g.conns (Connect returns its
					// receiver); storing it again would only resurrect it if a
					// concurrent Disconnect has removed and closed it meanwhile.
					return true
				})
			case <-hcTick.C:
				err = g.rangeConns("health check", func(addr string, p pool.Conn) bool {
					// if addr or pool is nil or empty the registration of conns is invalid let's disconnect them
					if addr == "" || p == nil {
						disconnectTargets = append(disconnectTargets, addr)
						return true
					}
					// for health check we don't need to reconnect when connection is healthy
					if p.IsHealthy(ctx) {
						return true
					}
					log.Debugf("unhealthy connection detected for gRPC %s addr: %s trying to reconnect. detail: %s", g.name, addr, p.String())
					// if connection is not ip direct or unhealthy let's re-connect
					var err error
					// if not healthy we should try reconnect
					p, err = p.Reconnect(ctx, false)
					if errors.IsNot(err, context.Canceled,
						context.DeadlineExceeded,
						errors.ErrCircuitBreakerOpenState,
						errors.ErrGRPCClientConnNotFound("*"),
						errors.ErrGRPCClientConnNotFound(addr),
						errors.ErrGRPCClientNotFound) {
						select {
						case <-ctx.Done():
							cerr := ctx.Err()
							if errors.IsNot(cerr, context.Canceled, context.DeadlineExceeded) {
								log.Error(errors.Join(err, cerr))
							}
						case ech <- err:
						}
					}
					// if rebalanced connection pool is nil even error is nil we should disconnect and delete it
					if err == nil && p == nil {
						disconnectTargets = append(disconnectTargets, addr)
						return true
					}
					// if connection pool could not recover we should try next connection loop
					if err != nil || p == nil || !p.IsHealthy(ctx) {
						_, _ = g.crl.LoadOrStore(addr, false)
						return true
					}
					// p is the pool already held in g.conns (Reconnect returns its
					// receiver); storing it again would only resurrect it if a
					// concurrent Disconnect has removed and closed it meanwhile.
					return true
				})
			}
			if err != nil && errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) && len(addrs) != 0 {
				for _, addr := range addrs {
					if addr != "" {
						log.Debugf("connection for gRPC %s addr = %s not found in connection map will re-connect soon", g.name, addr)
						_, _ = g.crl.LoadOrStore(addr, false)
					}
				}
			}

			clctx, cancel := context.WithTimeout(ctx, reconnLimitDuration)
			g.crl.Range(func(addr string, enabled bool) bool {
				select {
				case <-clctx.Done():
					return false
				default:
					defer g.crl.Delete(addr)
					var p pool.Conn
					if enabled && g.bo != nil {
						_, err = g.bo.Do(clctx, func(ictx context.Context) (r any, ret bool, err error) {
							p, err = g.Connect(ictx, addr)
							return nil, err != nil, err
						})
					} else {
						p, err = g.Connect(clctx, addr)
					}
					if err != nil || p == nil || !p.IsHealthy(clctx) {
						log.Debugf("connection for gRPC %s addr = %s is not healthy will delete soon,\terror: %v,\tpool: [%v]", g.name, addr, err, p)
						disconnectTargets = append(disconnectTargets, addr)
					}
					return true
				}
			})
			cancel()
			var (
				disconnectFlag bool
				host           string
				port           uint16
			)
			if len(disconnectTargets) > 0 {
				disconnected := make(map[string]bool, len(disconnectTargets))
				log.Debugf("starting to bulk disconnection for gRPC %s addrs %v", g.name, disconnectTargets)
				for _, addr := range disconnectTargets {
					var isLocal, isIPv4, isIPv6 bool
					host, port, isLocal, isIPv4, isIPv6, err = net.Parse(ctx, addr)
					_ = isLocal
					_ = isIPv4
					_ = isIPv6
					disconnectFlag = net.IsIPLiteral(host) // Disconnect only if the connection is a direct IP connection; do not delete connections via DNS due to retry.
					if err != nil {
						log.Warnf("failed to parse addr %s for disconnection checking, will disconnect soon: host: %s, port %d, err: %v", addr, host, port, err)
						disconnectFlag = true // Disconnect if the address connected to is not parseable.
					}
					log.Debugf("disconnection target is addr: %s, host: %s, port: %d, disconnectFlag: %t, disconnected: %v", addr, host, port, disconnectFlag, disconnected)
					if disconnectFlag &&
						!disconnected[addr] {
						log.Debugf("bulk part addr %s disconnecting", addr)
						err = g.Disconnect(ctx, addr)
						if errors.IsNot(err, context.Canceled,
							context.DeadlineExceeded,
							errors.ErrCircuitBreakerOpenState,
							errors.ErrGRPCClientConnNotFound("*"),
							errors.ErrGRPCClientConnNotFound(addr),
							errors.ErrGRPCClientNotFound) {
							select {
							case <-ctx.Done():
								cerr := ctx.Err()
								if errors.IsNot(cerr, context.Canceled, context.DeadlineExceeded) {
									log.Error(errors.Join(err, cerr))
								}
							case ech <- err:
							}
						}
						disconnected[addr] = true
					}
				}
				disconnectTargets = disconnectTargets[:0]
			}
		}
	}))
	g.monitorMu.Lock()
	g.ech = ech
	g.monitorMu.Unlock()

	log.Debugf("connection monitor successfully started for %s", g.name)

	return ech, nil
}

// Range executes a function f for each healthy connection.
// It returns an error if no connections are available or if any of the function calls fail.
func (g *gRPCClient) Range(
	ctx context.Context,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (err error) {
	var span trace.Span
	sctx := ctx
	if trace.Enabled() {
		sctx, span = trace.StartSpan(ctx, apiName+"/Client.Range")
		defer trace.End(span)
	}
	snap := g.snapshot()
	if snap == nil || len(snap.targets) == 0 {
		log.Warnf("%s rangeConns for %s client conn Not Found Error at ending, len: %d,\tsize: %d,\taddrs: %v", g.name, "Range", g.conns.Len(), g.conns.Size(), g.addrs)
		err = errors.ErrGRPCClientConnNotFound("*")
		recordSpanError(span, err)
		return err
	}
	for _, t := range snap.targets {
		ssctx, sspan := startAddrSpan(sctx, apiName+"/Client.Range/", t.addr)
		rerr := g.doVoidPolicyTarget(ssctx, t, snap.gen, sspan != nil, f)
		if sspan != nil {
			recordSpanError(sspan, rerr)
			sspan.End()
		}
	}
	return nil
}

func (g *gRPCClient) runFanOutBatch(
	ctx context.Context,
	span trace.Span,
	b *fanOutBatch,
	count int,
	prefix string,
	mode fanOutMode,
	connGen uint64,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) error {
	b.g = g
	b.ctx = ctx
	b.f = f
	b.prefix = prefix
	b.connGen = connGen
	b.last = count - 1
	b.mode = mode
	err := g.workers.FanOut(g.eg, count, b.runFn)
	putFanOutBatch(b)
	if err != nil {
		recordSpanError(span, err)
		return err
	}
	return nil
}

func (g *gRPCClient) lookupOrderedTarget(snap *connSnapshot, addr, op string) *connTarget {
	t := snap.byAddr[addr]
	if t == nil || t.p == nil {
		g.crl.Store(addr, true)
		log.Warnf("gRPCClient.%s operation failed, gRPC connection pool for %s is invalid,\terror: %v", op, addr, errors.ErrGRPCClientConnNotFound(addr))
		return nil
	}
	return t
}

func (g *gRPCClient) execOrderedTarget(
	ctx context.Context,
	snap *connSnapshot,
	addr, op, prefix string,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) error {
	t := g.lookupOrderedTarget(snap, addr, op)
	if t == nil {
		return nil
	}
	ssctx, sspan := startAddrSpan(ctx, prefix, addr)
	err := g.doVoidPolicyTarget(ssctx, t, snap.gen, sspan != nil, f)
	if sspan != nil {
		recordSpanError(sspan, err)
		sspan.End()
	}
	return nil
}

func (g *gRPCClient) runTargetSingle(
	ctx context.Context,
	t *connTarget,
	connGen uint64,
	prefix string,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("recovered from panic in grpc %s: %v", prefix, r)
		}
	}()
	ssctx, sspan := startAddrSpan(ctx, prefix, t.addr)
	derr := g.doVoidPolicyTarget(ssctx, t, connGen, sspan != nil, f)
	if sspan != nil {
		if derr != nil {
			recordSpanError(sspan, derr)
		}
		sspan.End()
	}
}

// RangeConcurrent executes a function f for each healthy connection concurrently.
// The concurrency level is controlled by the `concurrency` parameter.
func (g *gRPCClient) RangeConcurrent(
	ctx context.Context,
	concurrency int,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (err error) {
	snap := g.snapshot()
	if snap == nil || len(snap.targets) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	var span trace.Span
	sctx := ctx
	if trace.Enabled() {
		sctx, span = trace.StartSpan(ctx, apiName+"/Client.RangeConcurrent")
		defer trace.End(span)
	}
	if concurrency == 1 {
		return g.Range(sctx, f)
	}
	targets := snap.targets
	connGen := snap.gen
	if len(targets) == 1 {
		g.runTargetSingle(sctx, targets[0], connGen, apiName+"/Client.RangeConcurrent/", f)
		return nil
	}
	if concurrency <= 0 || concurrency >= len(targets) {
		b := getFanOutBatch()
		b.snapTargets = targets
		return g.runFanOutBatch(sctx, span, b, len(targets), apiName+"/Client.RangeConcurrent/", fanOutModeTargets, connGen, f)
	}
	eg, egctx := errgroup.New(sctx)
	eg.SetLimit(concurrency)
	for _, t := range targets {
		target := t
		eg.Go(safety.RecoverFunc(func() (err error) {
			ssctx, sspan := startAddrSpan(egctx, apiName+"/Client.RangeConcurrent/", target.addr)
			err = g.doVoidPolicyTarget(ssctx, target, connGen, sspan != nil, f)
			if sspan != nil {
				recordSpanError(sspan, err)
				sspan.End()
			}
			if err != nil && errors.IsAny(err, context.Canceled, context.DeadlineExceeded) {
				return err
			}
			return nil
		}))
	}
	err = eg.Wait()
	if err != nil {
		recordSpanError(span, err)
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			return err
		}
	}
	return nil
}

// OrderedRange executes a function f for each healthy connection in the specified order.
func (g *gRPCClient) OrderedRange(
	ctx context.Context,
	orders []string,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (err error) {
	var span trace.Span
	sctx := ctx
	if trace.Enabled() {
		sctx, span = trace.StartSpan(ctx, apiName+"/Client.OrderedRange")
		defer trace.End(span)
	}
	if len(orders) == 0 {
		log.Warn("no order found for OrderedRange")
		return g.Range(sctx, f)
	}
	snap := g.snapshot()
	if snap == nil || len(snap.targets) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	var cnt int
	for _, addr := range orders {
		t := g.lookupOrderedTarget(snap, addr, "OrderedRange")
		if t == nil {
			continue
		}
		cnt++
		ssctx, sspan := startAddrSpan(sctx, apiName+"/Client.OrderedRange/", addr)
		ierr := g.doVoidPolicyTarget(ssctx, t, snap.gen, sspan != nil, f)
		if ierr != nil && (sspan != nil || errors.Is(ierr, errors.ErrGRPCClientConnNotFound("*"))) {
			err = errors.Join(err, ierr)
		}
		// End the per-address span here: a deferred End would keep every span open
		// until OrderedRange returns and stack one closure per address.
		if sspan != nil {
			recordSpanError(sspan, ierr)
			sspan.End()
		}
	}
	if cnt == 0 {
		err = errors.ErrGRPCClientConnNotFound("*")
	}
	if err != nil {
		recordSpanError(span, err)
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			return err
		}
	}
	return nil
}

func (g *gRPCClient) runOrderedSingle(
	ctx context.Context,
	snap *connSnapshot,
	addr string,
	op string,
	prefix string,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) {
	t := g.lookupOrderedTarget(snap, addr, op)
	if t == nil {
		return
	}
	g.runTargetSingle(ctx, t, snap.gen, prefix, f)
}

// OrderedRangeConcurrent executes a function f for each healthy connection in the specified order concurrently.
func (g *gRPCClient) OrderedRangeConcurrent(
	ctx context.Context,
	orders []string,
	concurrency int,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (err error) {
	var span trace.Span
	sctx := ctx
	if trace.Enabled() {
		sctx, span = trace.StartSpan(ctx, apiName+"/Client.OrderedRangeConcurrent")
		defer trace.End(span)
	}
	if len(orders) == 0 {
		log.Warn("no order found for OrderedRangeConcurrent")
		return g.RangeConcurrent(sctx, concurrency, f)
	}
	snap := g.snapshot()
	if snap == nil || len(snap.targets) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	if concurrency == 1 {
		return g.OrderedRange(sctx, orders, f)
	}
	if concurrency <= 0 || concurrency >= len(orders) {
		if len(orders) == 1 {
			g.runOrderedSingle(sctx, snap, orders[0], "OrderedRangeConcurrent", apiName+"/Client.OrderedRangeConcurrent/", f)
			return nil
		}
		b := getFanOutBatch()
		b.orders = orders
		b.snapByAddr = snap.byAddr
		return g.runFanOutBatch(sctx, span, b, len(orders), apiName+"/Client.OrderedRangeConcurrent/", fanOutModeOrdered, snap.gen, f)
	}
	eg, egctx := errgroup.New(sctx)
	eg.SetLimit(concurrency)
	for _, order := range orders {
		addr := order
		eg.Go(safety.RecoverFunc(func() error {
			return g.execOrderedTarget(egctx, snap, addr, "OrderedRangeConcurrent", apiName+"/Client.OrderedRangeConcurrent/", f)
		}))
	}
	err = eg.Wait()
	recordSpanError(span, err)
	return nil
}

// RangeSliceConcurrent executes f concurrently over count addresses from addrs
// starting at cyclic offset without allocating a rotated slice.
func (g *gRPCClient) RangeSliceConcurrent(
	ctx context.Context,
	addrs []string,
	offset int,
	count int,
	concurrency int,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (err error) {
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
	var span trace.Span
	sctx := ctx
	if trace.Enabled() {
		sctx, span = trace.StartSpan(ctx, apiName+"/Client.RangeSliceConcurrent")
		defer trace.End(span)
	}
	snap := g.snapshot()
	if snap == nil || len(snap.targets) == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	if concurrency == 1 {
		var cnt int
		for i := range count {
			idx := offset + i
			if idx >= n {
				idx -= n
			}
			addr := addrs[idx]
			t := g.lookupOrderedTarget(snap, addr, "RangeSliceConcurrent")
			if t == nil {
				continue
			}
			cnt++
			ssctx, sspan := startAddrSpan(sctx, apiName+"/Client.RangeSliceConcurrent/", addr)
			ierr := g.doVoidPolicyTarget(ssctx, t, snap.gen, sspan != nil, f)
			if ierr != nil && (sspan != nil || errors.Is(ierr, errors.ErrGRPCClientConnNotFound("*"))) {
				err = errors.Join(err, ierr)
			}
			if sspan != nil {
				recordSpanError(sspan, ierr)
				sspan.End()
			}
		}
		if cnt == 0 {
			err = errors.ErrGRPCClientConnNotFound("*")
		}
		if err != nil {
			recordSpanError(span, err)
			if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
				return err
			}
		}
		return nil
	}
	if concurrency <= 0 || concurrency >= count {
		if count == 1 {
			g.runOrderedSingle(sctx, snap, addrs[offset], "RangeSliceConcurrent", apiName+"/Client.RangeSliceConcurrent/", f)
			return nil
		}
		b := getFanOutBatch()
		b.orders = addrs
		b.offset = offset
		b.snapByAddr = snap.byAddr
		return g.runFanOutBatch(sctx, span, b, count, apiName+"/Client.RangeSliceConcurrent/", fanOutModeSlice, snap.gen, f)
	}
	eg, egctx := errgroup.New(sctx)
	eg.SetLimit(concurrency)
	for i := range count {
		idx := offset + i
		if idx >= n {
			idx -= n
		}
		addr := addrs[idx]
		eg.Go(safety.RecoverFunc(func() error {
			return g.execOrderedTarget(egctx, snap, addr, "RangeSliceConcurrent", apiName+"/Client.RangeSliceConcurrent/", f)
		}))
	}
	err = eg.Wait()
	recordSpanError(span, err)
	return nil
}

// RoundRobin is a generic function that executes a gRPC call in a round-robin fashion.
func RoundRobin[R any](
	ctx context.Context,
	c Client,
	f func(ctx context.Context,
		conn *grpc.ClientConn, copts ...grpc.CallOption) (R, error),
) (data R, err error) {
	res, err := c.RoundRobin(ctx, func(ctx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
		return f(ctx, conn, copts...)
	})
	// data is zero (nil) value of R
	if err != nil {
		return data, err
	}
	var ok bool
	if data, ok = res.(R); ok {
		return data, nil
	}
	return data, errors.UnexpectedProtoMessageType(data, res)
}

// RoundRobin executes a function f on available connections in a round-robin manner
// until one of the calls succeeds.
func (g *gRPCClient) RoundRobin(
	ctx context.Context,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (any, error),
) (data any, err error) {
	var span trace.Span
	sctx := ctx
	if trace.Enabled() {
		sctx, span = trace.StartSpan(ctx, apiName+"/Client.RoundRobin")
		defer trace.End(span)
	}
	snap := g.snapshot()
	if snap == nil || len(snap.targets) == 0 {
		return nil, errors.ErrGRPCClientConnNotFound("*")
	}

	if g.canFastPolicy() {
		var ret bool
		data, ret, err = g.roundRobinOnce(sctx, snap, f)
		if err == nil {
			return data, nil
		}
		if !ret || g.boRetryer == nil {
			return nil, err
		}
		return g.retryRoundRobin(sctx, err, f)
	}

	var boName string
	if boName = FromGRPCMethod(sctx); boName != "" {
		sctx = backoff.WithBackoffName(sctx, boName)
	}
	exec := func(ictx context.Context) (data any, ret bool, err error) {
		curSnap := g.snapshot()
		if curSnap == nil || len(curSnap.targets) == 0 {
			return nil, false, errors.ErrGRPCClientConnNotFound("*")
		}
		connGen := curSnap.gen
		n := len(curSnap.targets)
		start := 0
		if n > 1 {
			start = int((g.rrIdx.Add(1) - 1) % uint64(n)) //nolint:gosec // modulo bounds within [0, n)
		}
		for off := range n {
			idx := start + off
			if idx >= n {
				idx -= n
			}
			t := curSnap.targets[idx]
			addr, p := t.addr, t.p
			if p == nil {
				g.crl.Store(addr, true)
				continue
			}
			ret = false
			tctx, sspan := startAddrSpan(ictx, apiName+"/Client.RoundRobin/", addr)
			var boName string
			tctx = WrapGRPCMethod(tctx, addr)
			if boName = FromGRPCMethod(tctx); boName != "" {
				tctx = backoff.WithBackoffName(tctx, boName)
			}

			if boName != "" && g.useBreaker(addr, p, true, connGen) && g.cb != nil {
				data, err = g.cb.Do(tctx, boName, func(cbctx context.Context) (any, error) {
					data, ret, err = g.executeRPC(cbctx, p, addr, f)
					if err != nil && !ret {
						return data, errors.NewErrCircuitBreakerIgnorable(err)
					}
					return data, err
				})
			} else {
				data, ret, err = g.executeRPC(tctx, p, addr, f)
			}

			if err != nil {
				if sspan != nil {
					recordSpanError(sspan, err)
					sspan.End()
				}
				continue
			}
			if sspan != nil {
				sspan.End()
			}
			return data, ret, nil
		}
		return data, ret, err
	}

	if g.bo != nil {
		return g.bo.Do(sctx, exec)
	}
	data, _, err = exec(sctx)
	return data, err
}

func (g *gRPCClient) roundRobinOnce(
	ctx context.Context,
	snap *connSnapshot,
	f func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error),
) (data any, ret bool, err error) {
	connGen := snap.gen
	n := len(snap.targets)
	start := 0
	if n > 1 {
		start = int((g.rrIdx.Add(1) - 1) % uint64(n)) //nolint:gosec // modulo bounds within [0, n)
	}
	for off := range n {
		idx := start + off
		if idx >= n {
			idx -= n
		}
		t := snap.targets[idx]
		addr, p := t.addr, t.p
		if p == nil {
			g.crl.Store(addr, true)
			continue
		}
		ret = false
		tctx, sspan := startAddrSpan(ctx, apiName+"/Client.RoundRobin/", addr)
		var br *circuitbreaker.Handle
		if FromGRPCMethod(ctx) != "" {
			br = g.resolveBreakerHandle(ctx, t, p, addr, true, connGen)
		}
		if br != nil {
			if berr := br.Allow(); berr != nil {
				err = berr
				if sspan != nil {
					recordSpanError(sspan, err)
					sspan.End()
				}
				continue
			}
		}
		var callErr error
		data, ret, callErr = g.executeRPC(tctx, p, addr, f)
		if callErr == nil {
			if br != nil {
				br.FinishSuccess()
			}
			if sspan != nil {
				sspan.End()
			}
			return data, false, nil
		}
		if br != nil {
			err = br.FinishError(callErr, ret)
		} else {
			err = callErr
		}
		if sspan != nil {
			recordSpanError(sspan, err)
			sspan.End()
		}
	}
	return nil, ret, err
}

// Keep retry closure in a separate noinline frame so healthy RoundRobin calls never allocate.
//
//go:noinline

func (g *gRPCClient) retryRoundRobin(
	ctx context.Context,
	firstErr error,
	f func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error),
) (data any, err error) {
	if backoff.FromBackoffName(ctx) == "" {
		if m := FromGRPCMethod(ctx); m != "" {
			ctx = backoff.WithBackoffName(ctx, m)
		}
	}
	rerr := g.boRetryer.RetryVoid(ctx, firstErr, func(ictx context.Context) (bool, error) {
		curSnap := g.snapshot()
		if curSnap == nil || len(curSnap.targets) == 0 {
			return false, errors.ErrGRPCClientConnNotFound("*")
		}
		var ret bool
		var passErr error
		data, ret, passErr = g.roundRobinOnce(ictx, curSnap, f)
		return ret, passErr
	})
	if rerr != nil {
		return nil, rerr
	}
	return data, nil
}

// Do executes a function f on a specific connection address.
func (g *gRPCClient) Do(
	ctx context.Context,
	addr string,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (any, error),
) (data any, err error) {
	snap := g.snapshot()
	var t *connTarget
	if snap != nil {
		t = snap.byAddr[addr]
	}
	if t == nil || t.p == nil {
		g.crl.Store(addr, true)
		err = errors.ErrGRPCClientConnNotFound(addr)
		log.Warnf("gRPCClient.Do operation failed, gRPC connection pool for %s is invalid,\terror: %v", addr, err)
		return nil, err
	}
	data, err = g.doPolicyConn(ctx, t, t.p, t.addr, true, snap.gen, f)
	return data, err
}

// executeRPC executes the given function f with a connection from the pool.
// It returns the result, a boolean indicating if the call is retryable, and an error.
func (g *gRPCClient) executeRPC(
	ctx context.Context,
	p pool.Conn,
	addr string,
	f func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error),
) (res any, retryable bool, err error) {
	return g.callRPC(ctx, p, addr, f, nil)
}

// callRPC runs an RPC on a pooled conn and classifies its error via
// handleExecuteErr. It is the generic body shared by executeRPC and the rpcCall
// callbacks (exec/execCB); the unconfigured void fast path is specialized in
// callVoid instead of going through here.
func (g *gRPCClient) callRPC(
	ctx context.Context,
	p pool.Conn,
	addr string,
	f func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error),
	fVoid func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (res any, retryable bool, err error) {
	conn, ok := p.Get(ctx)
	if !ok || conn == nil {
		return g.handleExecuteErr(ctx, p, addr, errors.ErrGRPCClientConnNotFound(addr))
	}
	// g.copts are already the default call options of every pooled conn
	// (New installs WithDefaultCallOptions into g.dopts, which every pool
	// dials with), so passing them per call again would only make grpc
	// allocate a combined slice and apply each option twice per RPC.
	if fVoid != nil {
		err = fVoid(ctx, addr, conn)
	} else {
		res, err = f(ctx, conn)
	}
	if err != nil {
		return g.handleExecuteErr(ctx, p, addr, err)
	}
	return res, false, nil
}

type rpcCallFailedError struct {
	err  error
	addr string
}

func (e *rpcCallFailedError) Error() string {
	if e.err == nil {
		return "addr: " + e.addr
	}
	return "addr: " + e.addr + ": " + e.err.Error()
}

func (e *rpcCallFailedError) Unwrap() error {
	return e.err
}

type rpcCall struct {
	p            pool.Conn
	g            *gRPCClient
	f            func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error)
	fVoid        func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error
	execFn       func(context.Context) (any, bool, error)
	execCBFn     func(context.Context) (any, error)
	addr         string
	boName       string
	connGen      uint64
	connVerified bool
	ret          bool
}

// callVoid is the void-RPC fast path of the unconfigured (no backoff, no
// circuit breaker) callers. It deliberately repeats callRPC's void branch so
// the void path carries no result/void dispatch of its own; this is an intentional
// specialization of one hot path, not speculative duplication. Error classification stays in
// the single handleExecuteErr, so only the two-step get/call is repeated.
func (g *gRPCClient) callVoid(
	ctx context.Context,
	p pool.Conn,
	addr string,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) error {
	conn, ok := p.Get(ctx)
	if !ok || conn == nil {
		_, _, err := g.handleExecuteErr(ctx, p, addr, errors.ErrGRPCClientConnNotFound(addr))
		return err
	}
	if err := f(ctx, addr, conn); err != nil {
		_, _, err = g.handleExecuteErr(ctx, p, addr, err)
		return err
	}
	return nil
}

func (rc *rpcCall) exec(ictx context.Context) (r any, ret bool, err error) {
	rc.ret = false
	if rc.boName != "" && rc.g.useBreaker(rc.addr, rc.p, rc.connVerified, rc.connGen) {
		r, err = rc.g.cb.Do(ictx, rc.boName, rc.execCBFn)
		return r, rc.ret, err
	}
	return rc.g.callRPC(ictx, rc.p, rc.addr, rc.f, rc.fVoid)
}

func (rc *rpcCall) execCB(ictx context.Context) (any, error) {
	r, ret, err := rc.g.callRPC(ictx, rc.p, rc.addr, rc.f, rc.fVoid)
	rc.ret = ret
	if err != nil && !ret {
		return r, errors.NewErrCircuitBreakerIgnorable(err)
	}
	return r, err
}

func (g *gRPCClient) handleExecuteErr(
	ctx context.Context, p pool.Conn, addr string, err error,
) (any, bool, error) {
	//nolint:errorlint // fast path for direct context sentinels before errors.IsAny fallback
	if err == context.Canceled || err == context.DeadlineExceeded {
		return nil, false, err
	}
	type grpcStatusError interface {
		GRPCStatus() *gstatus.Status
	}
	if se, ok := err.(grpcStatusError); ok {
		if st := se.GRPCStatus(); st != nil {
			code := st.Code()
			if code != codes.Unknown && code != codes.Canceled && code != codes.DeadlineExceeded {
				if backoff.IsQuietErrorLog(ctx) {
					if code != codes.NotFound &&
						code != codes.AlreadyExists &&
						code != codes.Aborted &&
						code != codes.Canceled &&
						code != codes.DeadlineExceeded &&
						code != codes.OK &&
						code != codes.Unimplemented {
						log.Debug(err)
					}
				} else {
					status.Log(code, err)
				}
				if isRetryableCode(ctx, code) {
					return nil, p.IsHealthy(ctx), err
				}
				return nil, false, err
			}
		}
	}
	if errors.IsAny(err,
		context.Canceled,
		context.DeadlineExceeded,
		errors.ErrCircuitBreakerOpenState,
		errors.ErrGRPCClientNotFound,
		errors.ErrGRPCClientConnNotFound("*"),
	) || (p.IsIPConn() && errors.Is(err, errors.ErrGRPCClientConnNotFound(addr))) {
		return nil, false, err
	}
	st, ok := status.FromError(err)
	if !ok || st == nil {
		return nil, p.IsHealthy(ctx), err
	}
	if backoff.IsQuietErrorLog(ctx) {
		log.Debug(err)
	} else {
		status.Log(st.Code(), err)
	}
	if isRetryableCode(ctx, st.Code()) {
		return nil, p.IsHealthy(ctx), err
	}
	return nil, false, err
}

func isRetryableCode(ctx context.Context, code codes.Code) bool {
	if code == codes.ResourceExhausted {
		return !isNoResourceExhaustedRetry(ctx)
	}
	return code == codes.Internal || code == codes.Unavailable
}

func (g *gRPCClient) handleNilPool(addr string) error {
	g.crl.Store(addr, true)
	err := errors.ErrGRPCClientConnNotFound(addr)
	log.Warnf("gRPCClient.do operation failed, gRPC connection pool for %s is invalid,\terror: %v", addr, err)
	return err
}

func (g *gRPCClient) execCallVoid(
	ctx context.Context,
	p pool.Conn,
	addr string,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) error {
	if err := g.callVoid(ctx, p, addr, f); err != nil {
		return &rpcCallFailedError{addr: addr, err: err}
	}
	return nil
}

func (g *gRPCClient) doVoid(
	ctx context.Context,
	p pool.Conn,
	addr string,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) error {
	return g.doVoidPolicyConn(ctx, nil, p, addr, false, 0, true, f)
}

func (g *gRPCClient) canFastPolicy() bool {
	return (g.cb == nil || g.cbResolver != nil) && (g.bo == nil || g.boRetryer != nil)
}

func (g *gRPCClient) doVoidPolicyTarget(
	ctx context.Context,
	t *connTarget,
	connGen uint64,
	needWrappedErr bool,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) error {
	return g.doVoidPolicyConn(ctx, t, t.p, t.addr, true, connGen, needWrappedErr, f)
}

func (g *gRPCClient) doVoidPolicyConn(
	ctx context.Context,
	t *connTarget,
	p pool.Conn,
	addr string,
	connVerified bool,
	connGen uint64,
	needWrappedErr bool,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) error {
	if p == nil {
		return g.handleNilPool(addr)
	}
	if !trace.Enabled() {
		return g.dispatchVoid(ctx, t, p, addr, connVerified, connGen, needWrappedErr, f)
	}
	return g.doVoidSlow(ctx, t, p, addr, connVerified, connGen, needWrappedErr, f)
}

func (g *gRPCClient) doVoidSlow(
	ctx context.Context,
	t *connTarget,
	p pool.Conn,
	addr string,
	connVerified bool,
	connGen uint64,
	needWrappedErr bool,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) error {
	sctx, span := startAddrSpan(ctx, apiName+"/Client.do/", addr)
	defer trace.End(span)
	return g.dispatchVoid(sctx, t, p, addr, connVerified, connGen, needWrappedErr, f)
}

func (g *gRPCClient) dispatchVoid(
	ctx context.Context,
	t *connTarget,
	p pool.Conn,
	addr string,
	connVerified bool,
	connGen uint64,
	needWrappedErr bool,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) error {
	if g.bo == nil && g.cb == nil {
		return g.execCallVoid(ctx, p, addr, f)
	}
	if g.canFastPolicy() {
		return g.execVoidFast(ctx, t, p, addr, connVerified, connGen, needWrappedErr, f)
	}
	_, err := g.execWithPolicy(ctx, p, addr, connVerified, connGen, nil, f)
	return err
}

func (g *gRPCClient) resolveBreakerHandle(
	ctx context.Context, t *connTarget, p pool.Conn, addr string, connVerified bool, connGen uint64,
) *circuitbreaker.Handle {
	if g.cbResolver == nil {
		return nil
	}
	if !connVerified || g.disconnectGen.Load() != connGen {
		snap := g.snapshot()
		if snap == nil {
			return nil
		}
		cur := snap.byAddr[addr]
		if cur == nil || cur.p != p {
			return nil
		}
		if t == nil {
			t = cur
		}
	}
	m := FromGRPCMethod(ctx)
	if m == "" || m == addr {
		if t != nil {
			if br := t.addrBr.Load(); br != nil && !br.IsEvicted() {
				return br
			}
		}
		br, err := g.cbResolver.Resolve(addr)
		if err != nil || br == nil {
			return nil
		}
		if t != nil {
			t.addrBr.Store(br)
		}
		return br
	}
	var key string
	if hasAddrSuffix(m, addr) {
		key = m
	} else {
		kp := breakerKeyPair{method: m, addr: addr}
		if entry, ok := g.breakerKeys.Load(kp); ok {
			if entry.h != nil && !entry.h.IsEvicted() {
				return entry.h
			}
			key = entry.name
		} else {
			key = m + "/" + addr
		}
		br, err := g.cbResolver.Resolve(key)
		if err != nil || br == nil {
			return nil
		}
		g.breakerKeys.Store(kp, breakerKeyEntry{h: br, name: key})
		return br
	}
	br, err := g.cbResolver.Resolve(key)
	if err != nil {
		return nil
	}
	return br
}

func (g *gRPCClient) execVoidFast(
	ctx context.Context,
	t *connTarget,
	p pool.Conn,
	addr string,
	connVerified bool,
	connGen uint64,
	needWrappedErr bool,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) error {
	br := g.resolveBreakerHandle(ctx, t, p, addr, connVerified, connGen)
	if br != nil {
		if berr := br.Allow(); berr != nil {
			if needWrappedErr {
				return &rpcCallFailedError{addr: addr, err: berr}
			}
			return berr
		}
	}
	conn, ok := p.Get(ctx)
	var err error
	if !ok || conn == nil {
		err = errors.ErrGRPCClientConnNotFound(addr)
	} else {
		err = f(ctx, addr, conn)
	}
	if err == nil {
		if br != nil {
			br.FinishSuccess()
		}
		return nil
	}
	_, ret, herr := g.handleExecuteErr(ctx, p, addr, err)
	if !ret || g.boRetryer == nil {
		if br != nil {
			if needWrappedErr {
				herr = br.FinishError(herr, ret)
			} else {
				br.RecordError(herr, ret)
			}
		}
		if needWrappedErr {
			return &rpcCallFailedError{addr: addr, err: herr}
		}
		return herr
	}
	if br != nil {
		herr = br.FinishError(herr, ret)
	}
	_, err = g.retryFast(ctx, br, p, addr, connVerified, connGen, herr, nil, f)
	return err
}

// Keep retry closure in a separate noinline frame so healthy calls never heap-escape captured variables.
//
//go:noinline

func (g *gRPCClient) retryFast(
	ctx context.Context,
	br *circuitbreaker.Handle,
	p pool.Conn,
	addr string,
	connVerified bool,
	connGen uint64,
	firstErr error,
	f func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error),
	fVoid func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (data any, err error) {
	if backoff.FromBackoffName(ctx) == "" {
		ctx = backoff.WithBackoffName(ctx, g.cachedBreakerName(ctx, addr))
	}
	rerr := g.boRetryer.RetryVoid(ctx, firstErr, func(ictx context.Context) (bool, error) {
		if br != nil && !g.useBreaker(addr, p, connVerified, connGen) {
			br = nil
		}
		if br != nil {
			if berr := br.Allow(); berr != nil {
				return false, berr
			}
		}
		rconn, rok := p.Get(ictx)
		var callErr error
		if !rok || rconn == nil {
			callErr = errors.ErrGRPCClientConnNotFound(addr)
		} else if fVoid != nil {
			callErr = fVoid(ictx, addr, rconn)
		} else {
			data, callErr = f(ictx, rconn)
		}
		if callErr == nil {
			if br != nil {
				br.FinishSuccess()
			}
			return false, nil
		}
		_, rret, rherr := g.handleExecuteErr(ictx, p, addr, callErr)
		if br != nil {
			rherr = br.FinishError(rherr, rret)
		}
		return rret, rherr
	})
	if rerr != nil {
		return nil, &rpcCallFailedError{addr: addr, err: rerr}
	}
	return data, nil
}

func (g *gRPCClient) execWithPolicy(
	sctx context.Context,
	p pool.Conn,
	addr string,
	connVerified bool,
	connGen uint64,
	f func(ctx context.Context, conn *ClientConn, copts ...CallOption) (any, error),
	fVoid func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (data any, err error) {
	rc := getRPCCall()
	rc.g = g
	rc.p = p
	rc.addr = addr
	rc.connVerified = connVerified
	rc.connGen = connGen
	rc.boName = ""
	if g.cb != nil {
		rc.boName = g.cachedBreakerName(sctx, addr)
	}
	rc.f = f
	rc.fVoid = fVoid
	rc.ret = false
	if g.bo != nil {
		if backoff.FromBackoffName(sctx) == "" {
			name := rc.boName
			if name == "" {
				name = g.cachedBreakerName(sctx, addr)
			}
			sctx = backoff.WithBackoffName(sctx, name)
		}
		data, err = g.bo.Do(sctx, rc.execFn)
	} else {
		data, _, err = rc.exec(sctx)
	}
	rc.g = nil
	rc.p = nil
	rc.addr = ""
	rc.boName = ""
	rc.connVerified = false
	rc.connGen = 0
	rc.f = nil
	rc.fVoid = nil
	rpcCallPool.Put(rc)
	if err != nil {
		return nil, &rpcCallFailedError{addr: addr, err: err}
	}
	return data, nil
}

func hasAddrSuffix(m, addr string) bool {
	return m == addr || (len(m) > len(addr) && m[len(m)-len(addr)-1] == '/' && strings.HasSuffix(m, addr))
}

// cachedBreakerName returns the circuit breaker key using g.breakerKeys to avoid
// per-RPC string concatenation allocations when the method prefix is not yet
// suffixed by "/"+addr.
func (g *gRPCClient) cachedBreakerName(ctx context.Context, addr string) string {
	m := FromGRPCMethod(ctx)
	switch {
	case m == "":
		return addr
	case hasAddrSuffix(m, addr):
		return m
	}
	kp := breakerKeyPair{method: m, addr: addr}
	if entry, ok := g.breakerKeys.Load(kp); ok {
		return entry.name
	}
	k := m + "/" + addr
	actual, _ := g.breakerKeys.LoadOrStore(kp, breakerKeyEntry{name: k})
	return actual.name
}

func newRPCCall() *rpcCall {
	rc := &rpcCall{}
	rc.execFn = rc.exec
	rc.execCBFn = rc.execCB
	return rc
}

func getRPCCall() *rpcCall {
	if rc := rpcCallPool.Get(); rc != nil {
		return rc
	}
	return newRPCCall()
}

var rpcCallPool = sync.Pool[*rpcCall]{New: newRPCCall} //nolint:gochecknoglobals // RPC call carriers are reused across gRPCClient.do/doVoid calls

// do is an internal method that handles the execution of a function on a connection,
// including backoff and circuit breaker logic.
func (g *gRPCClient) do(
	ctx context.Context,
	p pool.Conn,
	addr string,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (any, error),
) (data any, err error) {
	return g.doPolicyConn(ctx, nil, p, addr, false, 0, f)
}

func (g *gRPCClient) doPolicyConn(
	ctx context.Context,
	t *connTarget,
	p pool.Conn,
	addr string,
	connVerified bool,
	connGen uint64,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (any, error),
) (data any, err error) {
	if p == nil {
		return nil, g.handleNilPool(addr)
	}
	if !trace.Enabled() {
		return g.dispatchAny(ctx, t, p, addr, connVerified, connGen, f)
	}
	return g.doSlow(ctx, t, p, addr, connVerified, connGen, f)
}

func (g *gRPCClient) doSlow(
	ctx context.Context,
	t *connTarget,
	p pool.Conn,
	addr string,
	connVerified bool,
	connGen uint64,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (any, error),
) (data any, err error) {
	sctx, span := startAddrSpan(ctx, apiName+"/Client.do/", addr)
	defer trace.End(span)
	return g.dispatchAny(sctx, t, p, addr, connVerified, connGen, f)
}

func (g *gRPCClient) dispatchAny(
	ctx context.Context,
	t *connTarget,
	p pool.Conn,
	addr string,
	connVerified bool,
	connGen uint64,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (any, error),
) (data any, err error) {
	if g.bo == nil && g.cb == nil {
		if data, _, err = g.callRPC(ctx, p, addr, f, nil); err != nil {
			return nil, &rpcCallFailedError{addr: addr, err: err}
		}
		return data, nil
	}
	if g.canFastPolicy() {
		br := g.resolveBreakerHandle(ctx, t, p, addr, connVerified, connGen)
		if br != nil {
			if berr := br.Allow(); berr != nil {
				return nil, &rpcCallFailedError{addr: addr, err: berr}
			}
		}
		conn, ok := p.Get(ctx)
		if !ok || conn == nil {
			err = errors.ErrGRPCClientConnNotFound(addr)
		} else {
			data, err = f(ctx, conn)
		}
		if err == nil {
			if br != nil {
				br.FinishSuccess()
			}
			return data, nil
		}
		_, ret, herr := g.handleExecuteErr(ctx, p, addr, err)
		if br != nil {
			herr = br.FinishError(herr, ret)
		}
		if !ret || g.boRetryer == nil {
			return nil, &rpcCallFailedError{addr: addr, err: herr}
		}
		return g.retryFast(ctx, br, p, addr, connVerified, connGen, herr, f, nil)
	}
	return g.execWithPolicy(ctx, p, addr, connVerified, connGen, f, nil)
}

// GetDialOption returns the default dial options.
func (g *gRPCClient) GetDialOption() []DialOption {
	return g.dopts
}

// GetCallOption returns the default call options.
func (g *gRPCClient) GetCallOption() []CallOption {
	return g.copts
}

// GetBackoff returns the backoff strategy.
func (g *gRPCClient) GetBackoff() backoff.Backoff {
	return g.bo
}

// SetDisableResolveDNSAddr dynamically disables DNS resolution for a specific address.
func (g *gRPCClient) SetDisableResolveDNSAddr(addr string, disabled bool) {
	// NOTE: When connecting to multiple locations, it was necessary to switch dynamically, so implementation was added.
	// There is no setting for disable on the helm chart side, so I used this implementation.
	g.disableResolveDNSAddrs.Store(addr, disabled)
}

// mergedDialOptions returns g.dopts followed by dopts in a newly allocated slice.
// g.dopts is shared by concurrent Connect calls and may have spare capacity,
// so appending to it directly would write into its shared backing array.
func (g *gRPCClient) mergedDialOptions(dopts []DialOption) []DialOption {
	return slices.Concat(g.dopts, dopts)
}

// Connect establishes or retrieves a connection to the given address.
// It handles connection pooling and retries.
func (g *gRPCClient) Connect(
	ctx context.Context, addr string, dopts ...DialOption,
) (conn pool.Conn, err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(ctx, apiName+"/Client.Connect/"+addr)
		defer trace.End(span)
	}

	conn, _, err = g.group.Do(ctx, "connect-"+addr, func(ctx context.Context) (pool.Conn, error) {
		handleError := func(err error) error {
			if err != nil {
				recordSpanError(span, err)
				return err
			}
			return nil
		}
		var (
			ok      bool
			oldConn pool.Conn
		)
		loadedConn, ok := g.conns.Load(addr)
		if ok && loadedConn != nil {
			if loadedConn.IsHealthy(ctx) {
				return loadedConn, nil
			}
			oldConn = loadedConn
			log.Debugf("%s connecting unhealthy pool addr= %s, conn=[%s]", g.name, addr, oldConn.String())
			conn, err = oldConn.Connect(ctx)
			if err == nil && conn != nil && conn.IsHealthy(ctx) {
				// conn is oldConn, already in g.conns: do not Store it again, which
				// would resurrect it after a concurrent Disconnect closed it.
				return conn, nil
			}
			if err != nil {
				log.Warnf("%s failed to reconnect unhealthy pool conn=[%s] for addr %s\terror= %v\t trying to disconnect", g.name, oldConn.String(), addr, err)
			}
		}
		log.Warnf("%s creating new connection pool (size: %d) for addr = %s", g.name, g.poolSize, addr)
		opts := []pool.Option{
			pool.WithAddr(addr),
			pool.WithSize(g.poolSize),
			pool.WithDialOptions(g.mergedDialOptions(dopts)...),
			pool.WithOldConnCloseDelay(g.roccd),
			pool.WithResolveDNS(func() bool {
				disabled, ok := g.disableResolveDNSAddrs.Load(addr)
				if ok && disabled {
					return false
				}
				return g.resolveDNS
			}()),
			pool.WithEnableMetrics(g.enablePoolMetrics),
			pool.WithErrGroup(g.eg),
		}
		if g.bo != nil {
			opts = append(opts, pool.WithBackoff(g.bo))
		}
		conn, err = pool.New(ctx, opts...)
		if err != nil || conn == nil {
			derr := g.Disconnect(ctx, addr)
			if derr != nil && errors.IsNot(derr, errors.ErrGRPCClientConnNotFound(addr)) {
				log.Warnf("%s failed to disconnect unhealthy pool addr= %s\terror= %s", g.name, addr, derr.Error())
				if err != nil {
					err = errors.Join(err, derr)
				} else {
					err = derr
				}
			}
			return nil, handleError(err)
		}
		log.Warnf("%s connecting to new connection pool for addr= %s, conn=[%s]", g.name, addr, conn.String())
		// newConn is never stored in g.conns on the error paths below, so it
		// must be disconnected there or its dialed slots leak.
		newConn := conn
		conn, err = conn.Connect(ctx)
		if err != nil {
			log.Error(err)
			derr := g.Disconnect(ctx, addr)
			if errors.IsNot(derr, errors.ErrGRPCClientConnNotFound(addr)) {
				var connStr string
				if conn != nil {
					connStr = conn.String()
				}
				log.Warnf("%s failed to disconnect unhealthy pool addr= %s\terror= %s\tconn=%s", g.name, addr, derr.Error(), connStr)
				err = errors.Join(err, derr)
			}
			if derr := newConn.Disconnect(ctx); derr != nil {
				err = errors.Join(err, derr)
			}
			return nil, handleError(err)
		}
		if conn == nil || !conn.IsHealthy(ctx) {
			if conn != nil {
				log.Debugf("%s connection to %s is unhealthy, conn=%s", g.name, addr, conn.String())
			} else {
				log.Debugf("%s connection to %s is nil", g.name, addr)
			}
			err = errors.ErrGRPCClientConnNotFound(addr)
			if derr := newConn.Disconnect(ctx); derr != nil {
				err = errors.Join(err, derr)
			}
			return nil, handleError(err)
		}

		// Count from the map transition, not the Load snapshot above: a
		// concurrent Disconnect may have removed (and uncounted) oldConn since.
		g.snapMu.Lock()
		if _, loaded := g.conns.Swap(addr, conn); !loaded {
			atomic.AddUint64(&g.clientCount, 1)
		} else {
			g.disconnectGen.Add(1)
		}
		g.rebuildSnapshotLocked()
		g.snapMu.Unlock()
		if oldConn != nil {
			dErr := oldConn.Disconnect(ctx)
			if dErr != nil {
				log.Warnf("failed to disconnect old connection for addr %s: %v", addr, dErr)
				err = errors.Join(err, dErr)
			}
		}
		return conn, handleError(err)
	})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (g *gRPCClient) IsConnected(ctx context.Context, addr string) bool {
	snap := g.snapshot()
	if snap == nil {
		return false
	}
	t := snap.byAddr[addr]
	if t == nil || t.p == nil {
		return false
	}
	return t.p.IsHealthy(ctx)
}

// Disconnect closes the connection to the given address and removes it from the pool.
func (g *gRPCClient) Disconnect(ctx context.Context, addr string) (err error) {
	log.Warnf("Disconnecting %s client connection for %s", g.name, addr)
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(ctx, apiName+"/Client.Disconnect/"+addr)
		defer trace.End(span)
	}

	_, _, err = g.group.Do(ctx, "disconnect-"+addr, func(ctx context.Context) (pool.Conn, error) {
		g.snapMu.Lock()
		p, ok := g.conns.LoadAndDelete(addr)
		if !ok || p == nil {
			g.snapMu.Unlock()
			log.Debugf("gRPC %s connection pool addr = %s is already unavailable or deleted", g.name, addr)
			err := errors.ErrGRPCClientConnNotFound(addr)
			if span != nil {
				span.RecordError(err)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		atomic.AddUint64(&g.clientCount, ^uint64(0))
		g.disconnectGen.Add(1)
		g.rebuildSnapshotLocked()
		g.snapMu.Unlock()
		g.evictAddr(addr)
		log.Debugf("gRPC %s connection pool addr = %s will disconnect soon...", g.name, addr)
		err := p.Disconnect(ctx)
		// RPCs still in flight on p re-create breaker/metrics keys when they
		// finish; once p is closed they fail fast, so evict again afterwards.
		// A call that obtained p before LoadAndDelete and starts after this
		// point skips the breaker (see useBreaker).
		g.evictAddr(addr)
		if err != nil {
			recordSpanError(span, err)
			return nil, err
		}

		return nil, nil
	})
	return err
}

// startAddrSpan starts the span prefix+addr. It skips building the name when
// tracing is disabled, where StartSpan would return (ctx, nil) anyway.
func startAddrSpan(ctx context.Context, prefix, addr string) (context.Context, trace.Span) {
	if !trace.Enabled() {
		return ctx, nil
	}
	return trace.StartSpan(ctx, prefix+addr)
}

// recordSpanError records a non-nil err on a non-nil span, adding the gRPC
// status attributes only when err carries a gRPC status.
func recordSpanError(span trace.Span, err error) {
	if span == nil || err == nil || !span.IsRecording() {
		return
	}
	msg := err.Error()
	span.RecordError(err)
	if st, ok := status.FromError(err); ok && st != nil {
		span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
	}
	span.SetStatus(trace.StatusError, msg)
}

// useBreaker reports whether a call to addr on p goes through the circuit
// breaker. A call that obtained p before Disconnect removed it and evicted the
// breakers of addr would otherwise re-create a breaker that nothing evicts
// again, so calls on a pool that is no longer connected skip the breaker.
// A Disconnect that completes between this check and the breaker's store can
// still leave one; that window is not closed.
func (g *gRPCClient) useBreaker(addr string, p pool.Conn, verified bool, gen uint64) bool {
	if g.cb == nil {
		return false
	}
	if verified && g.disconnectGen.Load() == gen {
		return true
	}
	if snap := g.snapshot(); snap != nil {
		t := snap.byAddr[addr]
		return t != nil && t.p == p
	}
	cur, ok := g.conns.Load(addr)
	return ok && cur == p
}

// addrEvicter is implemented by the circuit breaker manager and the dialer,
// which both keep state per address.
type addrEvicter interface {
	EvictAddr(addr string)
}

// evictAddr drops the per-address state that the breaker and the dialer keep
// for addr once its pool is removed. Both re-create it lazily, so a pool that
// connects to addr again later only starts from a fresh state.
func (g *gRPCClient) evictAddr(addr string) {
	g.breakerKeys.Range(func(kp breakerKeyPair, _ breakerKeyEntry) bool {
		if kp.addr == addr {
			g.breakerKeys.Delete(kp)
		}
		return true
	})
	if e, ok := g.cb.(addrEvicter); ok {
		e.EvictAddr(addr)
	}
	if e, ok := g.dialer.(addrEvicter); ok {
		e.EvictAddr(addr)
	}
}

func (g *gRPCClient) ConnectedAddrs(_ context.Context) (addrs []string) {
	snap := g.snapshot()
	if snap == nil || len(snap.targets) == 0 {
		return nil
	}
	addrs = make([]string, 0, len(snap.targets))
	for _, t := range snap.targets {
		addrs = append(addrs, t.addr)
	}
	return addrs
}

func (g *gRPCClient) snapshot() *connSnapshot {
	if snap := g.targetsSnap.Load(); snap != nil {
		return snap
	}
	if g.conns.Len() == 0 {
		return nil
	}
	g.snapMu.Lock()
	defer g.snapMu.Unlock()
	if snap := g.targetsSnap.Load(); snap != nil {
		return snap
	}
	g.rebuildSnapshotLocked()
	return g.targetsSnap.Load()
}

func (g *gRPCClient) storeConn(addr string, p pool.Conn) {
	g.snapMu.Lock()
	defer g.snapMu.Unlock()
	g.conns.Store(addr, p)
	g.rebuildSnapshotLocked()
}

func (g *gRPCClient) rebuildSnapshotLocked() {
	old := g.targetsSnap.Load()
	n := g.conns.Len()
	if n == 0 {
		g.targetsSnap.Store(&connSnapshot{
			gen: g.disconnectGen.Load(),
		})
		return
	}
	targets := make([]*connTarget, 0, n)
	byAddr := make(map[string]*connTarget, n)
	g.conns.Range(func(addr string, p pool.Conn) bool {
		if p == nil {
			return true
		}
		var t *connTarget
		if old != nil && old.byAddr != nil {
			if prev, ok := old.byAddr[addr]; ok && prev.p == p {
				t = prev
			}
		}
		if t == nil {
			t = &connTarget{p: p, addr: addr}
		}
		targets = append(targets, t)
		byAddr[addr] = t
		return true
	})
	g.targetsSnap.Store(&connSnapshot{
		targets: targets,
		byAddr:  byAddr,
		gen:     g.disconnectGen.Load(),
	})
}

// Close gracefully closes all connections and stops the connection monitor.
func (g *gRPCClient) Close(ctx context.Context) (err error) {
	g.monitorMu.Lock()
	stop := g.stopMonitor
	g.monitorMu.Unlock()
	if stop != nil {
		stop()
	}
	defer g.workers.Close()
	for _, addr := range g.ConnectedAddrs(ctx) {
		derr := g.Disconnect(ctx, addr)
		if errors.IsNot(derr, errors.ErrGRPCClientConnNotFound(addr)) {
			err = errors.Join(err, derr)
		}
	}
	return err
}

// rangeConns is a helper function to iterate over the connection map.
func (g *gRPCClient) rangeConns(action string, fn func(addr string, p pool.Conn) bool) (err error) {
	if g.conns.Len() == 0 {
		log.Warnf("%s rangeConns for %s client conn Not Found Error at ending, len: %d,\tsize: %d,\taddrs: %v", g.name, action, g.conns.Len(), g.conns.Size(), g.addrs)
		return errors.ErrGRPCClientConnNotFound("*")
	}
	g.conns.Range(fn)
	return nil
}
