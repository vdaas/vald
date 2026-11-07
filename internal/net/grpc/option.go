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

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/circuitbreaker"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/client/metric"
	"github.com/vdaas/vald/internal/net/grpc/interceptor/client/trace"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/timeutil"
	"github.com/vdaas/vald/internal/tls"
	"google.golang.org/grpc"
	gbackoff "google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Option func(*gRPCClient)

var defaultOptions = []Option{
	WithConnectionPoolSize(10),
	WithKeepaliveParams("30s", "10s", true),
	WithEnableConnectionPoolRebalance(false),
	WithConnectionPoolRebalanceDuration("1h"),
	WithErrGroup(errgroup.Get()),
	WithHealthCheckDuration("10s"),
	WithResolveDNS(true),
	WithBackoffMaxDelay(gbackoff.DefaultConfig.MaxDelay.String()),
	WithBackoffBaseDelay(gbackoff.DefaultConfig.BaseDelay.String()),
	WithBackoffMultiplier(gbackoff.DefaultConfig.Multiplier),
	WithBackoffJitter(gbackoff.DefaultConfig.Jitter),
	WithMinConnectTimeout("20s"),
}

func WithAddrs(addrs ...string) Option {
	return func(g *gRPCClient) {
		if len(addrs) == 0 {
			return
		}
		if g.addrs == nil {
			g.addrs = make(map[string]struct{})
		}
		for _, addr := range addrs {
			g.addrs[addr] = struct{}{}
		}
	}
}

func WithHealthCheckDuration(dur string) Option {
	return func(g *gRPCClient) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			log.Errorf("failed to parse health check duration: %v", err)
			return
		}
		if d <= 0 {
			log.Errorf("invalid health check duration: %d", d)
			return
		}
		g.hcDur = d
	}
}

func WithConnectionPoolRebalanceDuration(dur string) Option {
	return func(g *gRPCClient) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			log.Errorf("failed to parse connection pool rebalance duration: %v", err)
			return
		}
		if d <= 0 {
			log.Errorf("invalid connection pool rebalance duration: %d", d)
			return
		}
		g.prDur = d
	}
}

func WithResolveDNS(flg bool) Option {
	return func(g *gRPCClient) {
		g.resolveDNS = flg
	}
}

func WithEnableConnectionPoolRebalance(flg bool) Option {
	return func(g *gRPCClient) {
		g.enablePoolRebalance = flg
	}
}

func WithConnectionPoolSize(size int) Option {
	return func(g *gRPCClient) {
		if size >= 1 {
			g.poolSize = uint64(size)
		}
	}
}

func WithEnableConnectionPoolMetrics(flg bool) Option {
	return func(g *gRPCClient) {
		g.enablePoolMetrics = flg
	}
}

func WithBackoffMaxDelay(dur string) Option {
	return func(g *gRPCClient) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			log.Errorf("failed to parse backoff max delay: %v", err)
			return
		}
		if d <= 0 {
			log.Errorf("invalid backoff max delay: %d", d)
			return
		}
		g.gbo.MaxDelay = d
	}
}

func WithBackoffBaseDelay(dur string) Option {
	return func(g *gRPCClient) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			log.Errorf("failed to parse backoff base delay: %v", err)
			return
		}
		if d <= 0 {
			log.Errorf("invalid backoff base delay: %d", d)
			return
		}
		g.gbo.BaseDelay = d
	}
}

func WithBackoffMultiplier(m float64) Option {
	return func(g *gRPCClient) {
		g.gbo.Multiplier = m
	}
}

func WithBackoffJitter(j float64) Option {
	return func(g *gRPCClient) {
		g.gbo.Jitter = j
	}
}

func WithMinConnectTimeout(dur string) Option {
	return func(g *gRPCClient) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			log.Errorf("failed to parse minimum connection timeout: %v", err)
			return
		}
		if d <= 0 {
			log.Errorf("invalid minimum connection timeout: %d", d)
			return
		}
		g.mcd = d
	}
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(g *gRPCClient) {
		if eg != nil {
			g.eg = eg
		}
	}
}

func WithBackoff(bo backoff.Backoff) Option {
	return func(g *gRPCClient) {
		if bo != nil {
			g.bo = bo
		}
	}
}

func WithCircuitBreaker(cb circuitbreaker.CircuitBreaker) Option {
	return func(gr *gRPCClient) {
		if cb != nil {
			gr.cb = cb
		}
	}
}

/*
API References https://pkg.go.dev/google.golang.org/grpc#CallOption

1. Already Implemented APIs
- func CallContentSubtype(contentSubtype string) CallOption
- func MaxCallRecvMsgSize(bytes int) CallOption
- func MaxCallSendMsgSize(bytes int) CallOption
- func MaxRetryRPCBufferSize(bytes int) CallOption
- func WaitForReady(waitForReady bool) CallOption

2. Unnecessary for this package APIs
- func Header(md *metadata.MD) CallOption
- func Peer(p *peer.Peer) CallOption
- func PerRPCCredentials(creds credentials.PerRPCCredentials) CallOption
- func StaticMethod() CallOption
- func Trailer(md *metadata.MD) CallOption

3. Experimental APIs
- func ForceCodec(codec encoding.Codec) CallOption
- func ForceCodecV2(codec encoding.CodecV2) CallOption
- func OnFinish(onFinish func(err error)) CallOption
- func UseCompressor(name string) CallOption

4. Deprecated APIs
- func CallCustomCodec(codec Codec) CallOption
- func FailFast(failFast bool) CallOption.
*/
const defaultCallOptionLength = 5

// appendCallOpts appends opts to g.copts, preallocating defaultCallOptionLength
// capacity on first use. It never retains the caller's slice.
func (g *gRPCClient) appendCallOpts(opts ...grpc.CallOption) {
	if g.copts == nil {
		g.copts = make([]grpc.CallOption, 0, defaultCallOptionLength)
	}
	g.copts = append(g.copts, opts...)
}

func WithCallOptions(opts ...grpc.CallOption) Option {
	return func(g *gRPCClient) {
		if len(opts) != 0 {
			g.appendCallOpts(opts...)
		}
	}
}

func WithCallContentSubtype(contentSubtype string) Option {
	return func(g *gRPCClient) {
		g.appendCallOpts(grpc.CallContentSubtype(contentSubtype))
	}
}

func WithMaxRecvMsgSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.appendCallOpts(grpc.MaxCallRecvMsgSize(size))
		}
	}
}

func WithMaxSendMsgSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.appendCallOpts(grpc.MaxCallSendMsgSize(size))
		}
	}
}

func WithMaxRetryRPCBufferSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.appendCallOpts(grpc.MaxRetryRPCBufferSize(size))
		}
	}
}

func WithWaitForReady(flg bool) Option {
	return func(g *gRPCClient) {
		g.appendCallOpts(grpc.WaitForReady(flg))
	}
}

/*
API References https://pkg.go.dev/google.golang.org/grpc#DialOption

1. Already Implemented APIs
- func WithAuthority(a string) DialOption
- func WithContextDialer(f func(context.Context, string) (net.Conn, error)) DialOption
- func WithDisableRetry() DialOption
- func WithIdleTimeout(d time.Duration) DialOption
- func WithInitialConnWindowSize(s int32) DialOption
- func WithInitialWindowSize(s int32) DialOption
- func WithKeepaliveParams(kp keepalive.ClientParameters) DialOption
- func WithMaxCallAttempts(n int) DialOption
- func WithMaxHeaderListSize(s uint32) DialOption
- func WithReadBufferSize(s int) DialOption
- func WithSharedWriteBuffer(val bool) DialOption
- func WithTransportCredentials(creds credentials.TransportCredentials) DialOption
- func WithUserAgent(s string) DialOption
- func WithWriteBufferSize(s int) DialOption

2. Unnecessary for this package APIs
- func WithChainStreamInterceptor(interceptors ...StreamClientInterceptor) DialOption
- func WithChainUnaryInterceptor(interceptors ...UnaryClientInterceptor) DialOption
- func WithConnectParams(p ConnectParams) DialOption
- func WithDefaultCallOptions(cos ...CallOption) DialOption
- func WithDefaultServiceConfig(s string) DialOption
- func WithDisableServiceConfig() DialOption
- func WithPerRPCCredentials(creds credentials.PerRPCCredentials) DialOption
- func WithStatsHandler(h stats.Handler) DialOption
- func WithStreamInterceptor(f StreamClientInterceptor) DialOption
- func WithUnaryInterceptor(f UnaryClientInterceptor) DialOption

3. Experimental APIs
- func WithChannelzParentID(c channelz.Identifier) DialOption
- func WithCredentialsBundle(b credentials.Bundle) DialOption
- func WithDisableHealthCheck() DialOption
- func WithNoProxy() DialOption
- func WithResolvers(rs ...resolver.Builder) DialOption

4. Deprecated APIs
- func FailOnNonTempDialError(f bool) DialOption
- func WithBackoffConfig(b BackoffConfig) DialOption
- func WithBackoffMaxDelay(md time.Duration) DialOption
- func WithBlock() DialOption
- func WithCodec(c Codec) DialOption
- func WithCompressor(cp Compressor) DialOption
- func WithDecompressor(dc Decompressor) DialOption
- func WithDialer(f func(string, time.Duration) (net.Conn, error)) DialOption
- func WithInsecure() DialOption
- func WithMaxMsgSize(s int) DialOption
- func WithReturnConnectionError() DialOption
- func WithTimeout(d time.Duration) DialOption
*/

const defaultDialOptionLength = 14

// appendDialOpts appends opts to g.dopts, preallocating defaultDialOptionLength
// capacity on first use. It never retains the caller's slice.
func (g *gRPCClient) appendDialOpts(opts ...grpc.DialOption) {
	if g.dopts == nil {
		g.dopts = make([]grpc.DialOption, 0, defaultDialOptionLength)
	}
	g.dopts = append(g.dopts, opts...)
}

func WithDialOptions(opts ...grpc.DialOption) Option {
	return func(g *gRPCClient) {
		if len(opts) != 0 {
			g.appendDialOpts(opts...)
		}
	}
}

func WithWriteBufferSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.appendDialOpts(
				grpc.WithWriteBufferSize(size),
			)
		}
	}
}

func WithReadBufferSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.appendDialOpts(
				grpc.WithReadBufferSize(size),
			)
		}
	}
}

func WithInitialWindowSize(size int32) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.appendDialOpts(
				grpc.WithInitialWindowSize(size),
			)
		}
	}
}

func WithInitialConnectionWindowSize(size int32) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.appendDialOpts(
				grpc.WithInitialConnWindowSize(size),
			)
		}
	}
}

func WithMaxMsgSize(size int) Option {
	return func(g *gRPCClient) {
		if size > 1 {
			g.appendDialOpts(
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(size)),
			)
		}
	}
}

func WithInsecure(flg bool) Option {
	return func(g *gRPCClient) {
		if flg {
			g.appendDialOpts(
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
		}
	}
}

func WithKeepaliveParams(t, to string, permitWithoutStream bool) Option {
	return func(g *gRPCClient) {
		if len(t) == 0 || len(to) == 0 {
			return
		}
		td, err := timeutil.Parse(t)
		if err != nil {
			log.Errorf("failed to parse grpc keepalive time: %s,\t%v", t, err)
			return
		}
		if td <= 0 {
			log.Errorf("invalid grpc keepalive time: %d", td)
			return
		}
		tod, err := timeutil.Parse(to)
		if err != nil {
			log.Errorf("failed to parse grpc keepalive timeout: %s,\t%v", to, err)
			return
		}
		if tod <= 0 {
			log.Errorf("invalid grpc keepalive timeout: %d", tod)
			return
		}
		g.appendDialOpts(
			grpc.WithKeepaliveParams(
				keepalive.ClientParameters{
					Time:                td,
					Timeout:             tod,
					PermitWithoutStream: permitWithoutStream,
				},
			),
		)
	}
}

func WithDialer(network string, der net.Dialer) Option {
	return func(g *gRPCClient) {
		if der != nil {
			g.dialer = der
			nt := net.NetworkTypeFromString(network)
			switch nt {
			case net.UDP, net.UDP4, net.UDP6:
				nt = net.UDP
			case net.UNIX, net.UNIXGRAM, net.UNIXPACKET:
				nt = net.UNIX
			default:
				nt = net.TCP
			}
			g.appendDialOpts(
				grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
					log.Debugf("gRPC context Dialer for network %s, addr is %s", nt.String(), addr)
					return g.dialer.GetDialer()(ctx, nt.String(), addr)
				}),
			)
		}
	}
}

func WithTLSConfig(cfg *tls.Config) Option {
	return func(g *gRPCClient) {
		if cfg != nil {
			g.appendDialOpts(
				grpc.WithTransportCredentials(credentials.NewTLS(cfg)),
			)
		}
	}
}

func WithAuthority(a string) Option {
	return func(g *gRPCClient) {
		if a != "" {
			g.appendDialOpts(
				grpc.WithAuthority(a),
			)
		}
	}
}

func WithDisableRetry(disable bool) Option {
	return func(g *gRPCClient) {
		if disable {
			g.appendDialOpts(
				grpc.WithDisableRetry(),
			)
		}
	}
}

func WithIdleTimeout(dur string) Option {
	return func(g *gRPCClient) {
		if len(dur) == 0 {
			return
		}
		d, err := timeutil.Parse(dur)
		if err != nil {
			log.Errorf("failed to parse idle timeout duration: %v", err)
			return
		}
		if d <= 0 {
			log.Errorf("invalid idle timeout duration: %d", d)
			return
		}
		g.appendDialOpts(
			grpc.WithIdleTimeout(d),
		)
	}
}

func WithMaxCallAttempts(n int) Option {
	return func(g *gRPCClient) {
		if n > 2 {
			g.appendDialOpts(
				grpc.WithMaxCallAttempts(n),
			)
		}
	}
}

func WithMaxHeaderListSize(size uint32) Option {
	return func(g *gRPCClient) {
		if size > 0 {
			g.appendDialOpts(
				grpc.WithMaxHeaderListSize(size),
			)
		}
	}
}

func WithSharedWriteBuffer(enable bool) Option {
	return func(g *gRPCClient) {
		if enable {
			g.appendDialOpts(
				grpc.WithSharedWriteBuffer(enable),
			)
		}
	}
}

func WithUserAgent(ua string) Option {
	return func(g *gRPCClient) {
		if ua != "" {
			g.appendDialOpts(
				grpc.WithUserAgent(ua),
			)
		}
	}
}

func WithClientInterceptors(names ...string) Option {
	return func(g *gRPCClient) {
		for _, name := range names {
			switch strings.ToLower(name) {
			case "traceinterceptor", "trace":
				g.appendDialOpts(
					WithStatsHandler(trace.NewStatsHandler()),
				)
			case "metricinterceptor", "metric":
				uci, sci, err := metric.ClientInterceptors()
				if err != nil {
					lerr := errors.NewErrCriticalOption("gRPCInterceptors", "metric", errors.Wrap(err, "failed to create interceptor"))
					log.Warn(lerr.Error())
					continue
				}
				g.appendDialOpts(
					grpc.WithUnaryInterceptor(uci),
					grpc.WithStreamInterceptor(sci),
				)
			default:
			}
		}
	}
}

func WithOldConnCloseDelay(dur string) Option {
	return func(g *gRPCClient) {
		if len(dur) == 0 {
			return
		}
		g.roccd = dur
	}
}
