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

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/plugin/grpctrace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ClientInterceptor interface {
	GetClientInterceptor() UnaryClientInterceptor
}

type clientInterceptor struct {
	tracerName string
	spanName   string
}

func NewClientInterceptor(opts ...ClientOption) ClientInterceptor {
	i := new(clientInterceptor)

	for _, opt := range append(clientDefaultOpts, opts...) {
		opt(i)
	}

	return i
}

func (i *clientInterceptor) GetClientInterceptor() UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		reqMeta, _ := metadata.FromOutgoingContext(ctx)
		metaCopy := reqMeta.Copy()

		tr := global.TraceProvider().Tracer(i.tracerName)
		err := tr.WithSpan(ctx, i.spanName,
			func(ctx context.Context) error {
				grpctrace.Inject(ctx, &metaCopy)
				ctx = metadata.NewOutgoingContext(ctx, metaCopy)

				err := invoker(ctx, method, req, reply, cc, opts...)
				// if err != nil {
				// 	s, _ := status.FromError(err)
				// 	trace.SpanFromContext(ctx).SetStatus(s.Code())
				// } else {
				// 	trace.SpanFromContext(ctx).SetStatus(codes.OK)
				// }
				return err
			})
		return err
	}
}
