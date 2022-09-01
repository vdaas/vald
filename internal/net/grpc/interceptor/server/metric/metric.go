package metric

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability-v2/attribute"
	"github.com/vdaas/vald/internal/observability-v2/metrics"
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
		gRPCStatus,
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
