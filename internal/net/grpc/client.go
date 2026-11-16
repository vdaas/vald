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
)

type (
	// CallOption wraps grpc.CallOption for convenience.
	CallOption = grpc.CallOption
	// DialOption wraps pool.DialOption for convenience.
	DialOption = pool.DialOption
	// ClientConn wraps pool.ClientConn for convenience.
	ClientConn = pool.ClientConn
)

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
	// bo is the backoff strategy used for retrying failed operations.
	bo backoff.Backoff
	// addrs is a set of initial target addresses for the gRPC client.
	addrs map[string]struct{}
	// stopMonitor is the cancel function to stop the connection monitor goroutine.
	// It is guarded by monitorMu.
	stopMonitor context.CancelFunc
	// ech is the channel for reporting errors from the connection monitor.
	// It is guarded by monitorMu.
	ech <-chan error
	// crl is the connection request list used to manage pending reconnection requests.
	crl sync.Map[string, bool]
	// disableResolveDNSAddrs stores addresses for which DNS resolution should be explicitly disabled.
	disableResolveDNSAddrs sync.Map[string, bool]
	// conns stores the active connection pools for each target address.
	conns sync.Map[string, pool.Conn]
	// name is the logical name of the gRPC client, used for logging and tracing.
	name string
	// roccd represents the duration to wait before closing old connections during reconnection (Reconnection Old Connection Closing Duration).
	roccd string
	// dopts is the list of default dial options applied to new connections.
	dopts []DialOption
	// copts is the list of default call options applied to RPC calls.
	copts []CallOption
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
	// monitorRunning indicates whether the connection monitor goroutine is currently active.
	monitorRunning atomic.Bool
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
		reconnLimitDuration := time.Second

		if g.enablePoolRebalance && g.prDur.Nanoseconds() > 0 {
			err = safety.RecoverFunc(func() error {
				prTick = time.NewTicker(g.prDur) // pool rebalance ticker
				return nil
			})()
			reconnLimitDuration = time.Duration(int64(math.Min(float64(g.hcDur.Nanoseconds()), float64(g.prDur.Nanoseconds()))))
		} else {
			err = safety.RecoverFunc(func() error {
				prTick = time.NewTicker(g.hcDur) // pool rebalance ticker
				return nil
			})()
			reconnLimitDuration = g.hcDur
		}
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

		disconnectTargets := make([]string, 0, len(addrs))
		log.Debugf("connection monitor loop starting for gRPC %s addrs = %v", g.name, addrs)
		for {
			select {
			case <-ctx.Done():
				if err != nil {
					return errors.Join(ctx.Err(), err)
				}
				return ctx.Err()
			case <-prTick.C:
				if g.enablePoolRebalance {
					err = g.rangeConns("pool rebalance", func(addr string, p pool.Conn) bool {
						// if addr or pool is nil or empty the registration of conns is invalid let's disconnect them
						if addr == "" || p == nil {
							disconnectTargets = append(disconnectTargets, addr)
							return true
						}
						log.Debugf("rebalancing pool connection for gRPC %s addr: %s, detail: %s", g.name, addr, p.String())
						var err error
						// for rebalancing connection we don't need to check connection health
						p, err = p.Connect(ctx)
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
						// p is the pool already held in g.conns (Connect returns its
						// receiver); storing it again would only resurrect it if a
						// concurrent Disconnect has removed and closed it meanwhile.
						return true
					})
				}
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
					log.Debugf("unheallthy connection detected for gRPC addr: %s trying to reconnect. detail: %s", g.name, addr, p.String())
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
				disconnected   = make(map[string]bool, len(disconnectTargets))
			)
			if len(disconnectTargets) > 0 {
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
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.Range")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if g.conns.Len() == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	err = g.rangeConns("Range", func(addr string, p pool.Conn) bool {
		ssctx, sspan := startAddrSpan(sctx, apiName+"/Client.Range/", addr)
		defer func() {
			if sspan != nil {
				sspan.End()
			}
		}()
		_, err := g.do(ssctx, p, addr, func(ictx context.Context, conn *ClientConn, copts ...CallOption,
		) (any, error) {
			return nil, f(ictx, addr, conn, copts...)
		})
		recordSpanError(sspan, err)
		return true
	})
	if err != nil {
		recordSpanError(span, err)
		if errors.Is(err, errors.ErrGRPCClientConnNotFound("*")) {
			return err
		}
	}
	return nil
}

// RangeConcurrent executes a function f for each healthy connection concurrently.
// The concurrency level is controlled by the `concurrency` parameter.
func (g *gRPCClient) RangeConcurrent(
	ctx context.Context,
	concurrency int,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (err error) {
	if g.conns.Len() == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.RangeConcurrent")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if concurrency == 1 {
		return g.Range(sctx, f)
	}
	eg, egctx := errgroup.New(sctx)
	if concurrency > 1 {
		eg.SetLimit(concurrency)
	}
	err = g.rangeConns("RangeConcurrent", func(addr string, p pool.Conn) bool {
		eg.Go(safety.RecoverFunc(func() (err error) {
			ssctx, sspan := startAddrSpan(egctx, apiName+"/Client.RangeConcurrent/", addr)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err = g.do(ssctx, p, addr, func(ictx context.Context,
				conn *ClientConn, copts ...CallOption,
			) (any, error) {
				return nil, f(ictx, addr, conn, copts...)
			})
			if err != nil {
				recordSpanError(sspan, err)
				// Decide independently of tracing: only a local ctx error stops
				// the fan-out, whether or not a span is recorded.
				if errors.IsAny(err, context.Canceled, context.DeadlineExceeded) {
					return err
				}
			}
			return nil
		}))
		return true
	})
	err = errors.Join(err, eg.Wait())
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
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.OrderedRange")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(orders) == 0 {
		log.Warn("no order found for OrderedRange")
		return g.Range(sctx, f)
	}
	if g.conns.Len() == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	var cnt int
	for _, addr := range orders {
		p, ok := g.conns.Load(addr)
		if !ok || p == nil {
			g.crl.Store(addr, true)
			log.Warnf("gRPCClient.OrderedRange operation failed, gRPC connection pool for %s is invalid,\terror: %v", addr, errors.ErrGRPCClientConnNotFound(addr))
			continue
		}
		cnt++
		ssctx, sspan := startAddrSpan(sctx, apiName+"/Client.OrderedRange/", addr)
		_, ierr := g.do(ssctx, p, addr, func(ictx context.Context,
			conn *ClientConn, copts ...CallOption,
		) (any, error) {
			return nil, f(ictx, addr, conn, copts...)
		})
		if ierr != nil {
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

// OrderedRangeConcurrent executes a function f for each healthy connection in the specified order concurrently.
func (g *gRPCClient) OrderedRangeConcurrent(
	ctx context.Context,
	orders []string,
	concurrency int,
	f func(ctx context.Context, addr string, conn *ClientConn, copts ...CallOption) error,
) (err error) {
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.OrderedRangeConcurrent")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if len(orders) == 0 {
		log.Warn("no order found for OrderedRangeConcurrent")
		return g.RangeConcurrent(sctx, concurrency, f)
	}
	if g.conns.Len() == 0 {
		return errors.ErrGRPCClientConnNotFound("*")
	}
	if concurrency == 1 {
		return g.OrderedRange(sctx, orders, f)
	}
	eg, egctx := errgroup.New(sctx)
	if concurrency > 1 {
		eg.SetLimit(concurrency)
	}
	for _, order := range orders {
		addr := order
		eg.Go(safety.RecoverFunc(func() (err error) {
			p, ok := g.conns.Load(addr)
			if !ok || p == nil {
				g.crl.Store(addr, true)
				log.Warnf("gRPCClient.OrderedRangeConcurrent operation failed, gRPC connection pool for %s is invalid,\terror: %v", addr, errors.ErrGRPCClientConnNotFound(addr))
				return nil
			}
			ssctx, sspan := startAddrSpan(egctx, apiName+"/Client.OrderedRangeConcurrent/", addr)
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			_, err = g.do(ssctx, p, addr, func(ictx context.Context,
				conn *ClientConn, copts ...CallOption,
			) (any, error) {
				return nil, f(ictx, addr, conn, copts...)
			})
			recordSpanError(sspan, err)
			return nil
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
	sctx, span := trace.StartSpan(ctx, apiName+"/Client.RoundRobin")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if g.conns.Len() == 0 {
		return nil, errors.ErrGRPCClientConnNotFound("*")
	}

	var boName string
	if boName = FromGRPCMethod(sctx); boName != "" {
		sctx = backoff.WithBackoffName(sctx, boName)
	}

	exec := func(ictx context.Context) (data any, ret bool, err error) {
		cerr := g.rangeConns("RoundRobin", func(addr string, p pool.Conn) bool {
			if p == nil {
				g.crl.Store(addr, true)
				return true
			}
			tctx, span := startAddrSpan(ictx, apiName+"/Client.RoundRobin/", addr)
			defer func() {
				if span != nil {
					span.End()
				}
			}()
			var boName string
			tctx = WrapGRPCMethod(tctx, addr)
			if boName = FromGRPCMethod(tctx); boName != "" {
				tctx = backoff.WithBackoffName(tctx, boName)
			}

			if boName != "" && g.useBreaker(addr, p) {
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
				recordSpanError(span, err)
				return true
			}
			return false
		})
		if cerr != nil {
			return nil, false, cerr
		}
		return data, ret, err
	}

	if g.bo != nil {
		return g.bo.Do(sctx, exec)
	}
	data, _, err = exec(sctx)
	return data, err
}

// Do executes a function f on a specific connection address.
func (g *gRPCClient) Do(
	ctx context.Context,
	addr string,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (any, error),
) (data any, err error) {
	p, ok := g.conns.Load(addr)
	if !ok || p == nil {
		g.crl.Store(addr, true)
		err = errors.ErrGRPCClientConnNotFound(addr)
		log.Warnf("gRPCClient.Do operation failed, gRPC connection pool for %s is invalid,\terror: %v", addr, err)
		return nil, err
	}
	data, err = g.do(ctx, p, addr, f)
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
	conn, ok := p.Get(ctx)
	if !ok || conn == nil {
		err = errors.ErrGRPCClientConnNotFound(addr)
	} else {
		// g.copts are already the default call options of every pooled conn
		// (New installs WithDefaultCallOptions into g.dopts, which every pool
		// dials with), so passing them per call again would only make grpc
		// allocate a combined slice and apply each option twice per RPC.
		res, err = f(ctx, conn)
	}

	if err != nil {
		// not-retryable errors
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
			// Canceled/DeadlineExceeded were already returned above.
			return nil, p.IsHealthy(ctx), err
		}
		if isQuietErrorLog(ctx) {
			log.Debug(err)
		} else {
			status.Log(st.Code(), err)
		}
		switch st.Code() {
		case codes.ResourceExhausted:
			if isNoResourceExhaustedRetry(ctx) {
				return nil, false, err
			}
			return nil, p.IsHealthy(ctx), err
		// retryable errors
		case codes.Internal,
			codes.Unavailable:
			return nil, p.IsHealthy(ctx), err
		}
		// not-retryable errors
		return nil, false, err
	}
	// success
	return res, false, nil
}

// do is an internal method that handles the execution of a function on a connection,
// including backoff and circuit breaker logic.
func (g *gRPCClient) do(
	ctx context.Context,
	p pool.Conn,
	addr string,
	f func(ctx context.Context,
		conn *ClientConn, copts ...CallOption) (any, error),
) (data any, err error) {
	if p == nil {
		g.crl.Store(addr, true)
		err = errors.ErrGRPCClientConnNotFound(addr)
		log.Warnf("gRPCClient.do operation failed, gRPC connection pool for %s is invalid,\terror: %v", addr, err)
		return nil, err
	}
	sctx := ctx
	var span trace.Span
	sctx, span = startAddrSpan(ctx, apiName+"/Client.do/", addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	sctx = WrapGRPCMethod(sctx, addr)
	boName := FromGRPCMethod(sctx)
	// exec applies the circuit breaker with or without backoff, as RoundRobin does.
	exec := func(ictx context.Context) (r any, ret bool, err error) {
		if boName != "" && g.useBreaker(addr, p) {
			r, err = g.cb.Do(ictx, boName, func(ictx context.Context) (any, error) {
				r, ret, err = g.executeRPC(ictx, p, addr, f)
				if err != nil && !ret {
					return r, errors.NewErrCircuitBreakerIgnorable(err)
				}
				return r, err
			})
			return r, ret, err
		}
		return g.executeRPC(ictx, p, addr, f)
	}
	if g.bo != nil {
		if boName != "" {
			sctx = backoff.WithBackoffName(sctx, boName)
		}
		data, err = g.bo.Do(sctx, exec)
	} else {
		data, _, err = exec(sctx)
	}
	if err != nil {
		return nil, errors.ErrRPCCallFailed(addr, err)
	}
	return data, nil
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
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Connect/"+addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

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
				log.Warnf("%s failed to disconnect unhealthy pool addr= %s\terror= %s\tconn=%s", g.name, addr, err.Error(), conn.String())
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
		if _, loaded := g.conns.Swap(addr, conn); !loaded {
			atomic.AddUint64(&g.clientCount, 1)
		}
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
	p, ok := g.conns.Load(addr)
	if !ok || p == nil {
		return false
	}
	return p.IsHealthy(ctx)
}

// Disconnect closes the connection to the given address and removes it from the pool.
func (g *gRPCClient) Disconnect(ctx context.Context, addr string) (err error) {
	log.Warnf("Disconnecting %s client connection for %s", g.name, addr)
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.Disconnect/"+addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	_, _, err = g.group.Do(ctx, "disconnect-"+addr, func(ctx context.Context) (pool.Conn, error) {
		p, ok := g.conns.LoadAndDelete(addr)
		if !ok || p == nil {
			log.Debugf("gRPC %s connection pool addr = %s is already unavailable or deleted", g.name, addr)
			err := errors.ErrGRPCClientConnNotFound(addr)
			if span != nil {
				span.RecordError(err)
				span.SetStatus(trace.StatusError, err.Error())
			}
			return nil, err
		}
		atomic.AddUint64(&g.clientCount, ^uint64(0))
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
	if span == nil || err == nil {
		return
	}
	span.RecordError(err)
	if st, ok := status.FromError(err); ok && st != nil {
		span.SetAttributes(trace.FromGRPCStatus(st.Code(), err.Error())...)
	}
	span.SetStatus(trace.StatusError, err.Error())
}

// useBreaker reports whether a call to addr on p goes through the circuit
// breaker. A call that obtained p before Disconnect removed it and evicted the
// breakers of addr would otherwise re-create a breaker that nothing evicts
// again, so calls on a pool that is no longer connected skip the breaker.
// A Disconnect that completes between this check and the breaker's store can
// still leave one; that window is not closed.
func (g *gRPCClient) useBreaker(addr string, p pool.Conn) bool {
	if g.cb == nil {
		return false
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
	if e, ok := g.cb.(addrEvicter); ok {
		e.EvictAddr(addr)
	}
	if e, ok := g.dialer.(addrEvicter); ok {
		e.EvictAddr(addr)
	}
}

func (g *gRPCClient) ConnectedAddrs(_ context.Context) (addrs []string) {
	addrs = make([]string, 0, g.conns.Len())
	if err := g.rangeConns("ConnectedAddrs", func(addr string, _ pool.Conn) bool {
		addrs = append(addrs, addr)
		return true
	}); err != nil {
		return nil
	}
	return addrs
}

// Close gracefully closes all connections and stops the connection monitor.
func (g *gRPCClient) Close(ctx context.Context) (err error) {
	g.monitorMu.Lock()
	stop := g.stopMonitor
	g.monitorMu.Unlock()
	if stop != nil {
		stop()
	}
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
