//
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
//

// Package metric provides gRPC client interceptors for client metric
package metric

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"google.golang.org/grpc"
)

const (
	latencyMetricsName       = "client_latency"
	completedRPCsMetricsName = "client_completed_rpcs"

	gRPCMethodKeyName = "grpc_client_method"
	gRPCStatus        = "grpc_client_status"
)

func ClientMetricInterceptors() (grpc.UnaryClientInterceptor, grpc.StreamClientInterceptor, error) {
	meter := metrics.GetMeter()

	latencyHistgram, err := meter.Float64Histogram(
		latencyMetricsName,
		metrics.WithDescription("Client latency in milliseconds, by method"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create latency metric")
	}

	completedRPCCnt, err := meter.Int64Counter(
		completedRPCsMetricsName,
		metrics.WithDescription("Count of RPCs by method and status"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create completedRPCs metric")
	}

	record := func(ctx context.Context, method string, err error, latency float64) {
		attrs := attributesFromError(method, err)
		latencyHistgram.Record(ctx, latency, metrics.WithAttributes(attrs...))
		completedRPCCnt.Add(ctx, 1, metrics.WithAttributes(attrs...))
	}
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			now := time.Now()
			err := invoker(ctx, method, req, reply, cc, opts...)
			elapsedTime := time.Since(now)
			record(ctx, method, err, float64(elapsedTime)/float64(time.Millisecond))
			return err
		}, func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			now := time.Now()
			_, err := streamer(ctx, desc, cc, method, opts...)
			elapsedTime := time.Since(now)
			record(ctx, method, err, float64(elapsedTime)/float64(time.Millisecond))
			return nil, nil
		}, nil
}

func attributesFromError(method string, err error) []attribute.KeyValue {
	code := codes.OK // default error is success when error is nil
	if err != nil {
		st, _ := status.FromError(err)
		if st != nil {
			code = st.Code()
		}
	}
	return []attribute.KeyValue{
		attribute.String(gRPCMethodKeyName, method),
		attribute.String(gRPCStatus, code.String()),
	}
}
