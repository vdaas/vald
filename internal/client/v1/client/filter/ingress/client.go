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

package ingress

import (
	"context"
	"reflect"
	"slices"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/filter/ingress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/filter/internal/dyntarget"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/strings"
)

type Client interface {
	ingress.FilterClient
	Target(ctx context.Context, addrs ...string) (ingress.FilterClient, error)
	// Release ends a Target use of addrs; see client.Release.
	Release(ctx context.Context, addrs ...string)
	GRPCClient() grpc.Client
	Start(context.Context) (<-chan error, error)
	Stop(context.Context) error
}

type client struct {
	c       grpc.Client
	targets *dyntarget.Pool // request-derived target connections; see Target
	addrs   []string
	idle    time.Duration // how long an unused request-derived target stays connected
}

type specificAddrClient struct {
	c    grpc.Client
	addr string
}

type multipleAddrsClient struct {
	c     grpc.Client
	addrs []string
}

const (
	apiName = "vald/internal/client/v1/client/filter/ingress"

	// defaultIdleTimeout is used when the filter client's
	// dial_option.idle_timeout is unset; it matches grpc-go's default channel
	// idle timeout, after which gRPC itself would idle the channel.
	defaultIdleTimeout = 30 * time.Minute
)

func New(opts ...Option) (Client, error) {
	c := new(client)
	for _, opt := range append(defaultOptions, opts...) {
		err := opt(c)
		if err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	if c.c == nil {
		if c.addrs == nil {
			return nil, errors.ErrGRPCTargetAddrNotFound
		}
		c.c = grpc.New("Ingress Filter Client", grpc.WithAddrs(c.addrs...))
	}
	if c.idle <= 0 {
		c.idle = defaultIdleTimeout
	}
	c.targets = dyntarget.New(c.c, c.addrs, c.idle)
	return c, nil
}

func (c *client) Start(ctx context.Context) (<-chan error, error) {
	c.targets.Start(ctx)
	return c.c.StartConnectionMonitor(ctx)
}

func (c *client) Stop(ctx context.Context) error {
	c.targets.Stop()
	return c.c.Close(ctx)
}

func (c *client) GRPCClient() grpc.Client {
	return c.c
}

// Target connects to targets and returns a client for them. Every successful
// call must be paired with Release(ctx, targets...) once the returned client is
// no longer used.
func (c *client) Target(ctx context.Context, targets ...string) (ingress.FilterClient, error) {
	if len(targets) == 0 {
		return nil, errors.ErrTargetNotFound
	}
	if err := c.targets.Acquire(ctx, targets...); err != nil {
		return nil, err
	}
	if len(targets) == 1 {
		return &specificAddrClient{
			addr: targets[0],
			c:    c.c,
		}, nil
	}
	return &multipleAddrsClient{
		addrs: slices.Clone(targets),
		c:     c.c,
	}, nil
}

// Release ends a Target use of addrs. A request-derived addr (one that is not
// in the client's configured addrs) is disconnected once it has had no
// in-flight user for the idle timeout; configured addrs stay connected.
func (c *client) Release(ctx context.Context, addrs ...string) {
	c.targets.Release(ctx, addrs...)
}

func (c *client) GenVector(
	ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption,
) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.GenVector")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Vector, error) {
		return ingress.NewFilterClient(conn).GenVector(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *client) FilterVector(
	ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption,
) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.FilterVector")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	res, err = grpc.RoundRobin(ctx, c.c, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (*payload.Object_Vector, error) {
		return ingress.NewFilterClient(conn).FilterVector(ctx, in, append(copts, opts...)...)
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *specificAddrClient) GenVector(
	ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption,
) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.GenVector/"+s.addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = s.c.Do(ctx, s.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (any, error) {
		res, err = ingress.NewFilterClient(conn).GenVector(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *specificAddrClient) FilterVector(
	ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption,
) (res *payload.Object_Vector, err error) {
	ctx, span := trace.StartSpan(ctx, apiName+"/Client.FilterVector/"+s.addr)
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = s.c.Do(ctx, s.addr, func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (any, error) {
		res, err = ingress.NewFilterClient(conn).FilterVector(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *multipleAddrsClient) GenVector(
	ctx context.Context, in *payload.Object_Blob, opts ...grpc.CallOption,
) (res *payload.Object_Vector, err error) {
	// Joining the addresses is per call; skip it when spans are not recorded.
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(ctx, apiName+"/Client.GetVector/["+strings.Join(m.addrs, ",")+"]")
	}
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	_, err = m.c.Do(ctx, m.addrs[0], func(ctx context.Context,
		conn *grpc.ClientConn,
		copts ...grpc.CallOption,
	) (any, error) {
		res, err = ingress.NewFilterClient(conn).GenVector(ctx, in, append(copts, opts...)...)
		return nil, err
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *multipleAddrsClient) FilterVector(
	ctx context.Context, in *payload.Object_Vector, opts ...grpc.CallOption,
) (res *payload.Object_Vector, err error) {
	// Joining the addresses is per call; skip it when spans are not recorded.
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(ctx, apiName+"/Client.FilterVector/["+strings.Join(m.addrs, ",")+"]")
	}
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	// Chain with Do instead of OrderedRange: OrderedRange swallows per-address
	// errors, which would silently skip a failing filter or return (nil, nil).
	for _, addr := range m.addrs {
		_, err = m.c.Do(ctx, addr, func(ctx context.Context,
			conn *grpc.ClientConn,
			copts ...grpc.CallOption,
		) (any, error) {
			r, ierr := ingress.NewFilterClient(conn).FilterVector(ctx, in, append(copts, opts...)...)
			if ierr != nil {
				return nil, ierr
			}
			in = r
			return r, nil
		})
		if err != nil {
			return nil, err
		}
	}
	return in, nil
}
