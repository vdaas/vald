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

package grpc

import (
	"context"
	"fmt"
	"runtime"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/errhandler"
	"github.com/vdaas/vald/internal/net/grpc/proto"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"google.golang.org/grpc"
)

type (
	ClientStream = grpc.ClientStream
	ServerStream = grpc.ServerStream

	TypedClientStream[Q, R proto.Message] interface {
		Send(Q) error
		Recv() (R, error)
		ClientStream
	}
	TypedServerStream[Q, R proto.Message] interface {
		Send(R) error
		Recv() (Q, error)
		ServerStream
	}
)

// BidirectionalStream represents gRPC bidirectional stream server handler.
// It receives messages from the stream, calls the function with the received message, and sends the returned message to the stream.
// It limits the number of concurrent calls to the function with the concurrency integer.
// It records errors and returns them as a single error.
func BidirectionalStream[Q, R proto.Message, S TypedServerStream[Q, R]](
	ctx context.Context, stream S, concurrency int, handle func(context.Context, Q) (R, error),
) (err error) {
	var span trace.Span
	if trace.Enabled() {
		ctx, span = trace.StartSpan(ctx, apiName+"/BidirectionalStream")
		defer trace.End(span)
	}
	if any(stream) == nil {
		return errors.ErrGRPCServerStreamNotFound
	}
	// Leave eg unbounded: concurrency == 1 dispatches inline, and concurrency > 1
	// bounds active goroutines via the worker pool below (avoiding errgroup's
	// SetLimit(1) inline deadlock on the persistent range worker and avoiding
	// semaphore acquisition failures while draining dc on canceled contexts).
	eg, ctx := errgroup.New(ctx)

	var (
		mu   sync.Mutex
		emu  sync.Mutex
		errs []error
	)

	finalize := func(workCh chan Q) (err error) {
		if workCh != nil {
			close(workCh)
		}
		err = eg.Wait()
		if err != nil {
			emu.Lock()
			errs = append(errs, err)
			emu.Unlock()
		}
		emu.Lock()
		nerrs := len(errs)
		if nerrs == 0 && span == nil {
			emu.Unlock()
			return nil
		}
		dedup := errors.RemoveDuplicates(errs)
		err = errors.Join(dedup...)
		emu.Unlock()
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse BidirectionalStream final gRPC error response")
		if span != nil && span.IsRecording() {
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(trace.StatusError, msg)
			}
		}
		return err
	}

	var cnt atomic.Uint64
	wctx, wcancel := context.WithCancel(ctx)
	defer wcancel()
	// One receiver goroutine serves the whole stream. dc is buffered with max(concurrency*2, 16)
	// so network read prefetching can overlap with worker dispatch.
	dc := make(chan Q, max(concurrency*2, 16))
	ec := make(chan error, 1)
	done := make(chan struct{})
	defer close(done)
	go func() {
		for {
			data, err := stream.Recv()
			if err != nil {
				ec <- err
				return
			}
			select {
			case dc <- data:
			case <-done:
				return
			}
		}
	}()

	processOne := func(data Q) {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("recovered from panic in BidirectionalStream: %v", r)
			}
		}()
		id := cnt.Add(1)
		sctx, sspan := ctx, trace.Span(nil)
		if trace.Enabled() {
			sctx, sspan = trace.StartSpan(ctx, fmt.Sprintf("%s/BidirectionalStream/stream-%020d", apiName, id))
		}
		res, herr := handle(sctx, data)
		if herr != nil {
			var (
				st   *status.Status
				perr error
			)
			if gst, ok := status.FromError(herr); ok && gst != nil {
				st, perr = gst, herr
			} else {
				var msg string
				st, msg, perr = status.ParseError(herr, codes.Internal, fmt.Sprintf("failed to parse BidirectionalStream id= %020d gRPC error response", id))
				_ = msg
			}
			if sspan != nil && sspan.IsRecording() {
				msg := perr.Error()
				if st != nil {
					msg = st.Message()
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				}
				sspan.RecordError(perr)
				sspan.SetStatus(trace.StatusError, msg)
			}
			if perr != nil && st != nil {
				code := st.Code()
				if code != codes.Canceled &&
					code != codes.DeadlineExceeded &&
					code != codes.InvalidArgument &&
					code != codes.NotFound &&
					code != codes.OK &&
					code != codes.Unimplemented {
					runtime.Gosched()
					log.Error(perr)
				}
			}
		}
		mu.Lock()
		serr := stream.Send(res)
		mu.Unlock()
		if serr != nil {
			runtime.Gosched()
			serr = errors.Wrapf(serr, "BidirectionalStream SendMsg returned error at stream-%020d", id)
			emu.Lock()
			errs = append(errs, serr)
			emu.Unlock()
			st, msg, perr := status.ParseError(serr, codes.Internal, fmt.Sprintf("failed to parse BidirectionalStream.SendMsg id= %020d gRPC error response", id),
				&errdetails.RequestInfo{
					RequestId:   fmt.Sprintf("%s/BidirectionalStream/stream-%020d/SendMsg", apiName, id),
					ServingData: errdetails.Serialize(res),
				})
			if sspan != nil && sspan.IsRecording() {
				sspan.RecordError(perr)
				sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
				sspan.SetStatus(trace.StatusError, msg)
			}
			emu.Lock()
			errs = append(errs, perr)
			emu.Unlock()
			wcancel()
		}
		if sspan != nil {
			sspan.End()
		}
	}

	var (
		workCh         chan Q
		workersStarted int
		workerFn       func() error
	)
	if concurrency > 1 {
		const (
			workChScale  = 2
			minWorkChCap = 16
		)
		workCh = make(chan Q, max(concurrency*workChScale, minWorkChCap))
		workerFn = func() error { //nolint:unparam // eg.Go requires func() error
			for data := range workCh {
				processOne(data)
			}
			return nil
		}
	}

	// dispatch executes inline when concurrency == 1, hands the message to a
	// bounded persistent worker pool when concurrency > 1, or spawns an
	// unbounded errgroup goroutine when concurrency <= 0.
	dispatch := func(data Q) {
		switch {
		case concurrency == 1:
			processOne(data)
		case concurrency > 1:
			if workersStarted < concurrency {
				workersStarted++
				eg.Go(workerFn)
			}
			workCh <- data
		default:
			eg.Go(safety.RecoverWithoutPanicFunc(func() error {
				processOne(data)
				return nil
			}))
		}
	}
	for {
		select {
		case <-ctx.Done():
			// A message may have landed in dc concurrently with the cancellation
			// (e.g. a sibling worker's Send failure self-cancelled this errgroup ctx
			// at the same instant a message was received). Drain dc non-blockingly and
			// dispatch before finalizing, so fully-received messages are not silently
			// dropped.
			for {
				select {
				case data := <-dc:
					dispatch(data)
				default:
					return finalize(workCh)
				}
			}
		case <-wctx.Done():
			for {
				select {
				case data := <-dc:
					dispatch(data)
				default:
					return finalize(workCh)
				}
			}
		case err := <-ec:
			for {
				select {
				case data := <-dc:
					dispatch(data)
				default:
					goto ecdrained
				}
			}
		ecdrained:
			if !errors.Is(err, io.EOF) {
				err = errors.Wrap(err, "BidirectionalStream Recv returned error")
				emu.Lock()
				errs = append(errs, err)
				emu.Unlock()
				log.Errorf("failed to receive stream message: %v", err)
			}
			return finalize(workCh)
		case data := <-dc:
			dispatch(data)
		}
	}
}

// BidirectionalStreamClient is gRPC client stream.
func BidirectionalStreamClient[Q, R proto.Message, S TypedClientStream[Q, R]](
	stream S, concurrency int, sendDataProvider func() (Q, bool), callBack func(R, error) bool,
) (err error) {
	if any(stream) == nil {
		return errors.ErrGRPCClientStreamNotFound
	}
	ctx, cancel := context.WithCancel(stream.Context())
	defer cancel()

	pipelined := concurrency > 1
	bufSize := 1
	if pipelined {
		bufSize = max(concurrency*2, 16)
	}

	recvDone := make(chan error, 1)
	go func() {
		recvDone <- safety.RecoverFunc(func() (err error) {
			resCh := make(chan R, bufSize)
			errCh := make(chan error, 1)
			var next chan struct{}
			if !pipelined {
				next = make(chan struct{})
			}
			go func() {
				for {
					res, rerr := stream.Recv()
					if rerr != nil {
						select {
						case errCh <- rerr:
						case <-ctx.Done():
						}
						return
					}
					select {
					case resCh <- res:
					case <-ctx.Done():
						return
					}
					if !pipelined {
						select {
						case <-next:
						case <-ctx.Done():
							return
						}
					}
				}
			}()
			for {
				select {
				case <-ctx.Done():
					return nil
				case res := <-resCh:
					if !callBack(res, nil) {
						cancel()
						return nil
					}
					if !pipelined {
						select {
						case next <- struct{}{}:
						case <-ctx.Done():
							return nil
						}
					}
				case rerr := <-errCh:
					for {
						select {
						case res := <-resCh:
							if !callBack(res, nil) {
								cancel()
								if errors.IsAny(rerr, io.EOF, context.Canceled, context.DeadlineExceeded) {
									return nil
								}
								return rerr
							}
						default:
							goto drained
						}
					}
				drained:
					cancel()
					if errors.IsAny(rerr, io.EOF, context.Canceled, context.DeadlineExceeded) {
						return nil
					}
					var zero R
					callBack(zero, rerr)
					return rerr
				}
			}
		})()
	}()

	return func() (err error) {
		finalize := func(sendErr error) error {
			if errors.IsAny(sendErr, io.EOF, context.Canceled, context.DeadlineExceeded) {
				sendErr = nil
			}
			serr := stream.CloseSend()
			if errors.IsNot(serr, io.EOF, context.Canceled, context.DeadlineExceeded) {
				sendErr = errors.Join(sendErr, serr)
			}
			select {
			case rerr := <-recvDone:
				if errors.IsNot(rerr, io.EOF, context.Canceled, context.DeadlineExceeded) {
					sendErr = errors.Join(sendErr, rerr)
				}
			case <-ctx.Done():
				cancel()
				rerr := <-recvDone
				if errors.IsNot(rerr, io.EOF, context.Canceled, context.DeadlineExceeded) {
					sendErr = errors.Join(sendErr, rerr)
				}
			}
			return sendErr
		}

		for {
			select {
			case <-ctx.Done():
				return finalize(ctx.Err())
			default:
				data, ok := sendDataProvider()
				if !ok {
					return finalize(nil)
				}
				err := stream.Send(data)
				if err != nil {
					return finalize(err)
				}
			}
		}
	}()
}

// HandleBidirectionalStream executes BidirectionalStream with standardized top-level
// and per-item tracing spans, status mapping, and span error recording.
// When grpcMethod is non-empty, it attaches gRPC method metadata to the parent
// span context and wraps per-item spans as "BidirectionalStream".
func HandleBidirectionalStream[Q, R, W proto.Message, S TypedServerStream[Q, W]](
	stream S,
	concurrency int,
	grpcMethod string,
	spanPrefix string,
	itemSpanSuffix func(Q) string,
	handler func(context.Context, Q) (R, error),
	wrapOK func(R) W,
	wrapErr func(*status.Status) W,
) error {
	ctx := stream.Context()
	var span trace.Span
	if trace.Enabled() {
		sctx := ctx
		if grpcMethod != "" {
			sctx = WithGRPCMethod(sctx, grpcMethod)
		}
		ctx, span = trace.StartSpan(sctx, spanPrefix)
		defer trace.End(span)
	}
	err := BidirectionalStream(ctx, stream, concurrency,
		func(sctx context.Context, req Q) (W, error) {
			var sspan trace.Span
			if trace.Enabled() {
				spanCtx := sctx
				if grpcMethod != "" {
					spanCtx = WrapGRPCMethod(spanCtx, "BidirectionalStream")
				}
				sctx, sspan = trace.StartSpan(spanCtx, spanPrefix+"/"+itemSpanSuffix(req))
				defer trace.End(sspan)
			}
			res, serr := handler(sctx, req)
			if serr != nil {
				st := errhandler.RecordSpanFromError(sspan, serr)
				return wrapErr(st), serr
			}
			return wrapOK(res), nil
		})
	if err != nil {
		if span != nil {
			errhandler.RecordSpanFromError(span, err)
		}
		return err
	}
	return nil
}
