//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package interceptor provides interceptors for grpc
package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/distributedcontext"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/key"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/plugin/grpctrace"
)

type ServerInterceptor interface {
	GetServerInterceptor() UnaryServerInterceptor
}

type serverInterceptor struct {
	tracerName         string
	grpcServerKeyName  string
	grpcServerKeyValue string
}

func NewServerInterceptor(opts ...ServerOption) ServerInterceptor {
	i := new(serverInterceptor)

	for _, opt := range append(serverDefaultOpts, opts...) {
		opt(i)
	}

	return i
}

func (i *serverInterceptor) GetServerInterceptor() UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		reqMeta, _ := metadata.FromIncomingContext(ctx)
		metaCopy := reqMeta.Copy()

		entries, spanCtx := grpctrace.Extract(ctx, &metaCopy)
		ctx = distributedcontext.WithMap(ctx, distributedcontext.NewMap(distributedcontext.MapUpdate{
			MultiKV: entries,
		}))

		grpcServerKey := key.New(i.grpcServerKeyName)
		serverSpanAttrs := []core.KeyValue{
			grpcServerKey.String(i.grpcServerKeyValue),
		}

		tr := global.TraceProvider().Tracer(i.tracerName)
		ctx, span := tr.Start(
			ctx,
			info.FullMethod,
			trace.WithAttributes(serverSpanAttrs...),
			trace.ChildOf(spanCtx),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		return handler(ctx, req)
	}
}
