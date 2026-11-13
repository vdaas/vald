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

package metric

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/attribute"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/sync"
	"google.golang.org/grpc"
)

const (
	latencyMetricsName       = "client_latency"
	completedRPCsMetricsName = "client_completed_rpcs"

	gRPCMethodKeyName = "grpc_client_method"
	gRPCStatus        = "grpc_client_status"
)

func ClientInterceptors() (grpc.UnaryClientInterceptor, grpc.StreamClientInterceptor, error) {
	return clientInterceptors(metrics.GetMeter())
}

func clientInterceptors(
	meter metrics.Meter,
) (grpc.UnaryClientInterceptor, grpc.StreamClientInterceptor, error) {
	latencyHistogram, err := meter.Float64Histogram(
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
		code := codes.OK
		if err != nil {
			st, _ := status.FromError(err)
			if st != nil {
				code = st.Code()
			}
		}
		opt := getMeasurementOption(method, code)
		latencyHistogram.Record(ctx, latency, opt)
		completedRPCCnt.Add(ctx, 1, opt)
	}
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			now := time.Now()
			err := invoker(ctx, method, req, reply, cc, opts...)
			elapsedTime := time.Since(now)
			record(ctx, method, err, float64(elapsedTime)/float64(time.Millisecond))
			return err
		}, func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			now := time.Now()
			stream, err := streamer(ctx, desc, cc, method, opts...)
			if err != nil {
				record(ctx, method, err, float64(time.Since(now))/float64(time.Millisecond))
				// Return the real stream and error from streamer: discarding them
				// (returning nil, nil) hands the caller a nil ClientStream with a
				// nil error, so its first Send/Recv/CloseSend nil-panics and any
				// stream-creation error is silently swallowed.
				return stream, err
			}
			// streamer only sets the stream up; the RPC's final status surfaces
			// later, so record once when the stream actually ends.
			cs := &clientStream{
				ClientStream:  stream,
				serverStreams: desc != nil && desc.ServerStreams,
			}
			cs.finish = func(err error) {
				cs.once.Do(func() {
					record(ctx, method, err, float64(time.Since(now))/float64(time.Millisecond))
				})
			}
			// Mirror grpc-go, which also ends the stream when ctx is done, so a
			// stream abandoned without a terminal RecvMsg is still counted.
			cs.stop = context.AfterFunc(ctx, func() { cs.finish(ctx.Err()) })
			return cs, nil
		}, nil
}

// clientStream records the stream's metrics once, when the RPC ends.
// It deliberately does not record on CloseSend: that is only a half-close,
// and the final status still arrives through RecvMsg.
type clientStream struct {
	grpc.ClientStream
	finish        func(error)
	stop          func() bool
	once          sync.Once
	serverStreams bool
}

func (cs *clientStream) end(err error) {
	cs.stop()
	cs.finish(err)
}

func (cs *clientStream) RecvMsg(m any) error {
	err := cs.ClientStream.RecvMsg(m)
	switch {
	case errors.Is(err, io.EOF):
		cs.end(nil)
	case err != nil:
		cs.end(err)
	case !cs.serverStreams:
		// A non-server-streaming RPC is complete once its single response
		// arrives; callers (e.g. CloseAndRecv) do not read on to io.EOF.
		cs.end(nil)
	}
	return err
}

func (cs *clientStream) SendMsg(m any) error {
	err := cs.ClientStream.SendMsg(m)
	// grpc-go ends the RPC on SendMsg errors other than io.EOF; io.EOF means
	// the real status will be returned by RecvMsg.
	if err != nil && !errors.Is(err, io.EOF) {
		cs.end(err)
	}
	return err
}

type metricOptKey struct {
	method string
	code   codes.Code
}

var optCache sync.Map[metricOptKey, metrics.MeasurementOption]

func getMeasurementOption(method string, code codes.Code) metrics.MeasurementOption {
	k := metricOptKey{
		method: method,
		code:   code,
	}
	if opt, ok := optCache.Load(k); ok {
		return opt
	}
	opt := metrics.WithAttributes(
		attribute.String(gRPCMethodKeyName, method),
		attribute.String(gRPCStatus, code.String()),
	)
	actual, _ := optCache.LoadOrStore(k, opt)
	return actual
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
