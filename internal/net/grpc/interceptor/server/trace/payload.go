//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package trace provides gRPC interceptors for traces
package trace

import (
	"bytes"
	"context"
	"path"
	"sync"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
)

const (
	grpcKindUnary  = "unary"
	grpcKindStream = "stream"

	traceAttrGRPCKind    = "grpc.kind"
	traceAttrGRPCService = "grpc.service"
	traceAttrGRPCMethod  = "grpc.method"

	traceAttrGRPCRequestPayload  = "grpc.request.payload"
	traceAttrGRPCResponsePayload = "grpc.response.payload"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func TracePayloadInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		span := trace.FromContext(ctx)
		if span == nil {
			return handler(ctx, req)
		}

		service, method := parseMethod(info.FullMethod)
		span.SetAttributes(
			trace.StringAttribute(traceAttrGRPCKind, grpcKindUnary),
			trace.StringAttribute(traceAttrGRPCService, service),
			trace.StringAttribute(traceAttrGRPCMethod, method),
		)

		if reqj := marshalJSON(req); reqj != "" {
			span.SetAttributes(
				trace.StringAttribute(traceAttrGRPCRequestPayload, reqj),
			)
		}

		resp, err = handler(ctx, req)

		if resj := marshalJSON(resp); resj != "" {
			span.SetAttributes(
				trace.StringAttribute(traceAttrGRPCResponsePayload, resj),
			)
		}

		return resp, err
	}
}

func TracePayloadStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		span := trace.FromContext(ss.Context())
		if span == nil {
			return handler(srv, ss)
		}

		service, method := parseMethod(info.FullMethod)
		span.SetAttributes(
			trace.StringAttribute(traceAttrGRPCKind, grpcKindStream),
			trace.StringAttribute(traceAttrGRPCService, service),
			trace.StringAttribute(traceAttrGRPCMethod, method),
		)

		tss := &tracingServerStream{
			ServerStream: ss,
		}

		err := handler(srv, tss)

		span.SetAttributes(
			trace.StringAttribute(traceAttrGRPCRequestPayload, tss.request),
			trace.StringAttribute(traceAttrGRPCResponsePayload, tss.response),
		)

		return err
	}
}

type tracingServerStream struct {
	grpc.ServerStream
	request  string
	response string
}

func (tss *tracingServerStream) RecvMsg(m interface{}) error {
	err := tss.ServerStream.RecvMsg(m)
	if err == nil && tss.request == "" {
		if reqj := marshalJSON(m); reqj != "" {
			tss.request = reqj
		}
	}

	return err
}

func (tss *tracingServerStream) SendMsg(m interface{}) error {
	err := tss.ServerStream.SendMsg(m)
	if err == nil && tss.response == "" {
		if resj := marshalJSON(m); resj != "" {
			tss.response = resj
		}
	}

	return err
}

func parseMethod(fullMethod string) (service, method string) {
	service = path.Dir(fullMethod)[1:]
	method = path.Base(fullMethod)

	return service, method
}

func marshalJSON(pbMsg interface{}) string {
	b := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(b)
	defer b.Reset()

	err := json.Encode(b, pbMsg)
	if err != nil {
		return ""
	}

	return conv.Btoa(b.Bytes())
}
