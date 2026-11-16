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
	ctx, span := trace.StartSpan(ctx, apiName+"/BidirectionalStream")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	if any(stream) == nil {
		return errors.ErrGRPCServerStreamNotFound
	}
	eg, ctx := errgroup.New(ctx)
	if concurrency > 0 {
		eg.SetLimit(concurrency)
	}

	var (
		mu   sync.Mutex
		emu  sync.Mutex
		errs = make([]error, 0, concurrency*2) // concurrency * recv+send
	)

	finalize := func() (err error) {
		err = eg.Wait()
		if err != nil {
			emu.Lock()
			errs = append(errs, err)
			emu.Unlock()
		}
		errs := errors.RemoveDuplicates(errs)
		emu.Lock()
		err = errors.Join(errs...)
		emu.Unlock()
		st, msg, err := status.ParseError(err, codes.Internal, "failed to parse BidirectionalStream final gRPC error response")
		if span != nil {
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(trace.StatusError, msg)
			}
		}
		return err
	}

	var cnt atomic.Uint64
	// Run Recv in a helper goroutine so a ctx cancellation actually interrupts the
	// receive loop. ctx here is a child of stream.Context() (errgroup.New above);
	// grpc's stream.Recv observes only the parent stream.Context(), so the previous
	// `default: stream.Recv()` could block forever once entered — e.g. after a
	// worker's stream.Send fails non-transport-fatally (MaxSendMsgSize/marshal) and
	// self-cancels this errgroup ctx — deadlocking finalize()'s eg.Wait against the
	// stuck receiver (the server-side counterpart of the BidirectionalStreamClient fix).
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
	// dispatch hands a received message to a worker in the errgroup. It is a named
	// closure (not inlined into the receive case) so the ctx.Done() branch can also
	// drain and dispatch an already-received message instead of dropping it.
	dispatch := func(data Q) {
		eg.Go(safety.RecoverWithoutPanicFunc(func() (err error) {
			id := cnt.Add(1)
			sctx, sspan := ctx, trace.Span(nil)
			if trace.Enabled() {
				sctx, sspan = trace.StartSpan(ctx, fmt.Sprintf("%s/BidirectionalStream/stream-%020d", apiName, id))
			}
			defer func() {
				if sspan != nil {
					sspan.End()
				}
			}()
			res, err := handle(sctx, data)
			if err != nil {
				st, msg, perr := status.ParseError(err, codes.Internal, fmt.Sprintf("failed to parse BidirectionalStream id= %020d gRPC error response", id))
				if sspan != nil {
					sspan.RecordError(perr)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, msg)
				}
				code := st.Code()
				if perr != nil && st != nil &&
					code != codes.Canceled &&
					code != codes.DeadlineExceeded &&
					code != codes.InvalidArgument &&
					code != codes.NotFound &&
					code != codes.OK &&
					code != codes.Unimplemented {
					runtime.Gosched()
					log.Error(perr)
				}
			}
			mu.Lock()
			err = stream.Send(res)
			mu.Unlock()
			if err != nil {
				runtime.Gosched()
				err = errors.Wrapf(err, "BidirectionalStream SendMsg returned error at stream-%020d", id)
				emu.Lock()
				errs = append(errs, err)
				emu.Unlock()
				st, msg, perr := status.ParseError(err, codes.Internal, fmt.Sprintf("failed to parse BidirectionalStream.SendMsg id= %020d gRPC error response", id),
					&errdetails.RequestInfo{
						RequestId:   fmt.Sprintf("%s/BidirectionalStream/stream-%020d/SendMsg", apiName, id),
						ServingData: errdetails.Serialize(res),
					})
				if sspan != nil {
					sspan.RecordError(perr)
					sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
					sspan.SetStatus(trace.StatusError, msg)
				}
				return perr
			}
			return nil
		}))
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
					return finalize()
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
			return finalize()
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
