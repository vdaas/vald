// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package metric

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
)

const (
	latencyMetricsName       = "server_latency"
	completedRPCsMetricsName = "server_completed_rpcs"

	gRPCMethodKeyName = "grpc_server_method"
	gRPCStatus        = "grpc_server_status"
)

func MetricInterceptor() (grpc.UnaryServerInterceptor, error) {
	meter := metrics.GetMeter()

	latencyHistgram, err := meter.SyncFloat64().Histogram(
		latencyMetricsName,
		metrics.WithDescription("Server latency in milliseconds, by method"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create latency metric")
	}

	completedRPCCnt, err := meter.SyncFloat64().Counter(
		completedRPCsMetricsName,
		metrics.WithDescription("Count of RPCs by method and status"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create completedRPCs metric")
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		now := time.Now()

		resp, err = handler(ctx, req)

		elapsedTime := time.Since(now)

		code := codes.Unknown.String()
		st, _ := status.FromError(err)
		if st != nil {
			code = st.Code().String()
		}

		latency := float64(elapsedTime) / float64(time.Millisecond)

		attrs := []attribute.KeyValue{
			attribute.String(gRPCMethodKeyName, info.FullMethod),
			attribute.String(gRPCStatus, code),
		}
		latencyHistgram.Record(ctx, latency, attrs...)
		completedRPCCnt.Add(ctx, latency, attrs...)

		return resp, err
	}, nil
}

func MetricStreamInterceptor() (grpc.StreamServerInterceptor, error) {
	meter := metrics.GetMeter()

	latencyHistgram, err := meter.SyncFloat64().Histogram(
		latencyMetricsName,
		metrics.WithDescription("Server latency in milliseconds, by method"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create latency metric")
	}

	completedRPCCnt, err := meter.SyncFloat64().Counter(
		completedRPCsMetricsName,
		metrics.WithDescription("Count of RPCs by method and status"),
		metrics.WithUnit(metrics.Milliseconds),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create completedRPCs metric")
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		now := time.Now()

		err = handler(srv, ss)

		elapsedTime := time.Since(now)

		code := codes.Unknown.String()
		st, _ := status.FromError(err)
		if st != nil {
			code = st.Code().String()
		}

		latency := float64(elapsedTime) / float64(time.Millisecond)

		attrs := []attribute.KeyValue{
			attribute.String(gRPCMethodKeyName, info.FullMethod),
			attribute.String(gRPCStatus, code),
		}
		latencyHistgram.Record(ss.Context(), latency, attrs...)
		completedRPCCnt.Add(ss.Context(), latency, attrs...)

		return err
	}, nil
}
