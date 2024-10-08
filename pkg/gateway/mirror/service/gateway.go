// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/internal/client/v1/client/mirror"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

const (
	// forwardedContextKey is the key used to store forwarding-related information in a context.
	forwardedContextKey = "forwarded-for"
)

// Gateway represents an interface for interacting with gRPC clients.
type Gateway interface {
	FromForwardedContext(ctx context.Context) string
	BroadCast(ctx context.Context,
		f func(ctx context.Context, target string, vc MirrorClient, copts ...grpc.CallOption) error) error
	Do(ctx context.Context, target string,
		f func(ctx context.Context, target string, vc MirrorClient, copts ...grpc.CallOption) (any, error)) (any, error)
	DoMulti(ctx context.Context, targets []string,
		f func(ctx context.Context, target string, vc MirrorClient, copts ...grpc.CallOption) error) error
	GRPCClient() grpc.Client
}

type gateway struct {
	// client is the Mirror gateway client for other clusters and the Vald gateway (e.g. LB gateway) client for the own cluster.
	client  mirror.Client
	eg      errgroup.Group
	podName string
}

// NewGateway returns Gateway object if no error occurs.
func NewGateway(opts ...Option) (Gateway, error) {
	g := new(gateway)
	for _, opt := range append(defaultGatewayOpts, opts...) {
		if err := opt(g); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(err, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	return g, nil
}

// GRPCClient returns the underlying gRPC client associated with this object.
// It provides access to the low-level gRPC client for advanced use cases.
func (g *gateway) GRPCClient() grpc.Client {
	return g.client.GRPCClient()
}

// forwardedContext takes a context and a podName, returning a new context
// with additional information related to forwarding.
func (*gateway) forwardedContext(ctx context.Context, podName string) context.Context {
	return grpc.NewOutgoingContext(ctx, grpc.MD{
		forwardedContextKey: []string{
			podName,
		},
	})
}

// FromForwardedContext extracts information from the forwarded context
// and returns the podName associated with it.
func (*gateway) FromForwardedContext(ctx context.Context) string {
	md, ok := grpc.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	vals, ok := md[forwardedContextKey]
	if !ok {
		return ""
	}
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// BroadCast performs a broadcast operation using the provided function
// to interact with gRPC clients for multiple targets.
// The provided function should handle the communication logic for a target.
func (g *gateway) BroadCast(
	ctx context.Context,
	f func(ctx context.Context, target string, vc MirrorClient, copts ...grpc.CallOption) error,
) (err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Gateway.BroadCast")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	return g.client.GRPCClient().RangeConcurrent(g.forwardedContext(ctx, g.podName), -1, func(ictx context.Context,
		addr string, conn *grpc.ClientConn, copts ...grpc.CallOption,
	) (err error) {
		select {
		case <-ictx.Done():
			return nil
		default:
			return f(ictx, addr, NewMirrorClient(conn), copts...)
		}
	})
}

// Do performs a gRPC operation on a single target using the provided function.
// It returns the result of the operation and any associated error.
// The provided function should handle the communication logic for a target.
func (g *gateway) Do(
	ctx context.Context,
	target string,
	f func(ctx context.Context, addr string, vc MirrorClient, copts ...grpc.CallOption) (any, error),
) (res any, err error) {
	ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Gateway.Do")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if target == "" {
		return nil, errors.ErrTargetNotFound
	}
	fctx := g.forwardedContext(ctx, g.podName)
	res, err = g.client.GRPCClient().Do(fctx, target,
		func(ictx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
			return f(ictx, target, NewMirrorClient(conn), copts...)
		},
	)
	if err != nil {
		return g.client.GRPCClient().RoundRobin(fctx, func(ictx context.Context, conn *grpc.ClientConn, copts ...grpc.CallOption) (any, error) {
			if strings.EqualFold(conn.Target(), target) {
				return nil, errors.ErrTargetNotFound
			}
			return f(ictx, conn.Target(), NewMirrorClient(conn), copts...)
		})
	}
	return res, nil
}

// DoMulti performs a gRPC operation on multiple targets using the provided function.
// It returns an error if any of the operations fails.
// The provided function should handle the communication logic for a target.
func (g *gateway) DoMulti(
	ctx context.Context,
	targets []string,
	f func(ctx context.Context, target string, vc MirrorClient, copts ...grpc.CallOption) error,
) error {
	ctx, span := trace.StartSpan(ctx, "vald/gateway/mirror/service/Gateway.DoMulti")
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	if len(targets) == 0 {
		return errors.ErrTargetNotFound
	}
	return g.client.GRPCClient().OrderedRangeConcurrent(g.forwardedContext(ctx, g.podName), targets, -1,
		func(ictx context.Context, addr string, conn *grpc.ClientConn, copts ...grpc.CallOption) (err error) {
			select {
			case <-ictx.Done():
				return nil
			default:
				return f(ictx, addr, NewMirrorClient(conn), copts...)
			}
		},
	)
}
