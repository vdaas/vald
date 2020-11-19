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

// Package grpc provides generic functionality for grpc
package grpc

import (
	"context"
	"path"
	"time"

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"google.golang.org/grpc"
)

type (
	UnaryServerInterceptor  = grpc.UnaryServerInterceptor
	StreamServerInterceptor = grpc.StreamServerInterceptor
)

var (
	UnaryInterceptor       = grpc.UnaryInterceptor
	ChainUnaryInterceptor  = grpc.ChainUnaryInterceptor
	StreamInterceptor      = grpc.StreamInterceptor
	ChainStreamInterceptor = grpc.ChainStreamInterceptor
)

func RecoverInterceptor() UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		err = safety.RecoverWithoutPanicFunc(func() (err error) {
			resp, err = handler(ctx, req)
			return err
		})()
		return resp, err
	}
}

func RecoverStreamInterceptor() StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		return safety.RecoverWithoutPanicFunc(func() (err error) {
			return handler(srv, ss)
		})()
	}
}

type AccessLogEntity struct {
	Grpc      *AccessLogGrpcEntity `json:"grpc,omitempty" yaml:"grpc"`
	StartTime float64              `json:"startTime,omitempty" yaml:"startTime"`
	Latency   float64              `json:"latency,omitempty" yaml:"latency"`
	TraceID   string               `json:"traceID,omitempty" yaml:"traceID"`
	Error     error                `json:"error,omitempty" yaml:"error"`
}

type AccessLogGrpcEntity struct {
	Kind    string `json:"kind,omitempty" yaml:"kind"`
	Service string `json:"service,omitempty" yaml:"service"`
	Method  string `json:"method,omitempty" yaml:"method"`
}

func AccessLogInterceptor() UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		var traceID string

		span := trace.FromContext(ctx)
		if span != nil {
			traceID = span.SpanContext().TraceID.String()
		}

		start := time.Now()

		resp, err = handler(ctx, req)

		latency := float64(time.Since(start)) / float64(time.Second)
		startTime := float64(start.UnixNano()) / float64(time.Second)

		service, method := parseMethod(info.FullMethod)

		entity := &AccessLogEntity{
			Grpc: &AccessLogGrpcEntity{
				Kind:    "unary",
				Service: service,
				Method:  method,
			},
			StartTime: startTime,
			Latency:   latency,
		}

		if traceID != "" {
			entity.TraceID = traceID
		}

		if err != nil {
			entity.Error = err
			log.Error("rpc completed", entity)
		} else {
			log.Info("rpc completed", entity)
		}

		return resp, err
	}
}

func parseMethod(fullMethod string) (service, method string) {
	service = path.Dir(fullMethod)[1:]
	method = path.Base(fullMethod)

	return service, method
}
