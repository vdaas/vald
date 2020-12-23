//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/gogo/protobuf/jsonpb"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/observability/trace"
	"google.golang.org/grpc"
)

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
		span.AddAttributes(
			trace.StringAttribute("grpc.kind", "unary"),
			trace.StringAttribute("grpc.service", service),
			trace.StringAttribute("grpc.method", method),
		)

		if reqj := marshalJSON(req); reqj != "" {
			span.AddAttributes(
				trace.StringAttribute("grpc.request.payload", reqj),
			)
		}

		resp, err = handler(ctx, req)

		if resj := marshalJSON(resp); resj != "" {
			span.AddAttributes(
				trace.StringAttribute("grpc.response.payload", resj),
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
		span.AddAttributes(
			trace.StringAttribute("grpc.kind", "stream"),
			trace.StringAttribute("grpc.service", service),
			trace.StringAttribute("grpc.method", method),
		)

		tss := &tracingServerStream{
			ServerStream: ss,
		}

		err := handler(srv, tss)

		span.AddAttributes(
			trace.StringAttribute("grpc.request.payload.first", tss.request),
			trace.StringAttribute("grpc.response.payload.first", tss.response),
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
	p, ok := pbMsg.(proto.Message)
	if !ok {
		return ""
	}

	b := &bytes.Buffer{}
	marshaler := &jsonpb.Marshaler{}

	if err := marshaler.Marshal(b, p); err != nil {
		return ""
	}

	return b.String()
}
