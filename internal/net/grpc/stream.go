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
	"sync/atomic"

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
			span.RecordError(err)
			span.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
			span.SetStatus(trace.StatusError, msg)
		}
		return err
	}

	var cnt uint64
	// Run Recv in a helper goroutine so a ctx cancellation actually interrupts the
	// receive loop. ctx here is a child of stream.Context() (errgroup.New above);
	// grpc's stream.Recv observes only the parent stream.Context(), so the previous
	// `default: stream.Recv()` could block forever once entered — e.g. after a
	// worker's stream.Send fails non-transport-fatally (MaxSendMsgSize/marshal) and
	// self-cancels this errgroup ctx — deadlocking finalize()'s eg.Wait against the
	// stuck receiver (the server-side counterpart of the BidirectionalStreamClient fix).
	recv := func() (<-chan Q, <-chan error) {
		dc := make(chan Q, 1)
		ec := make(chan error, 1)
		go func() {
			data, err := stream.Recv()
			if err != nil {
				ec <- err
				return
			}
			dc <- data
		}()
		return dc, ec
	}
	// dispatch hands a received message to a worker in the errgroup. It is a named
	// closure (not inlined into the receive case) so the ctx.Done() branch can also
	// drain and dispatch an already-received message instead of dropping it.
	dispatch := func(data Q) {
		eg.Go(safety.RecoverWithoutPanicFunc(func() (err error) {
			id := atomic.AddUint64(&cnt, 1)
			sctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s/BidirectionalStream/stream-%020d", apiName, id))
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
	dc, ec := recv()
	for {
		select {
		case <-ctx.Done():
			// A message may have landed in dc concurrently with the cancellation
			// (e.g. a sibling worker's Send failure self-cancelled this errgroup ctx
			// at the same instant a message was received). dc is buffered (cap 1) so
			// the receiver never blocks handing it off; drain it non-blockingly and
			// dispatch before finalizing, so a fully-received message is not silently
			// dropped by select choosing ctx.Done() over dc. This preserves the old
			// default-branch behavior that always handed a received message to eg.Go.
			select {
			case data := <-dc:
				dispatch(data)
			default:
			}
			return finalize()
		case err := <-ec:
			if !errors.Is(err, io.EOF) {
				err = errors.Wrap(err, "BidirectionalStream Recv returned error")
				emu.Lock()
				errs = append(errs, err)
				emu.Unlock()
				log.Errorf("failed to receive stream message: %v", err)
			}
			return finalize()
		case data := <-dc:
			dc, ec = recv()
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

	// The receiver runs on its own lifecycle, NOT inside the send errgroup, so the
	// send side can CloseSend first and only then wait for the receiver to drain to
	// EOF. Previously the receiver shared the errgroup that finalize waited on
	// *before* CloseSend, which deadlocked: CloseSend was never reached, so the
	// server kept blocking in its own Recv waiting for the client's half-close,
	// while this receiver blocked in stream.Recv waiting for a response that the
	// (idle) server would never send. cancel() also could not break it because it
	// cancels a context derived from stream.Context(), which grpc's stream.Recv
	// does not observe.
	recvDone := make(chan error, 1)
	go func() {
		recvDone <- safety.RecoverFunc(func() (err error) {
			for {
				// Run Recv in a helper goroutine so a cancelled ctx aborts the loop
				// even though grpc's stream.Recv only observes the stream's own
				// context, not the derived ctx that cancel() controls. This makes
				// cancellation effective (fix 2) instead of a no-op.
				resCh := make(chan R, 1)
				errCh := make(chan error, 1)
				go func() {
					res, rerr := stream.Recv()
					if rerr != nil {
						errCh <- rerr
						return
					}
					resCh <- res
				}()
				select {
				case <-ctx.Done():
					return nil
				case rerr := <-errCh:
					if errors.IsAny(rerr, io.EOF, context.Canceled, context.DeadlineExceeded) {
						cancel()
						return nil
					}
					var zero R
					if !callBack(zero, rerr) {
						cancel()
						return nil
					}
				case res := <-resCh:
					if !callBack(res, nil) {
						cancel()
						return nil
					}
				}
			}
		})()
	}()

	eg, egctx := errgroup.New(ctx)
	if concurrency > 0 {
		eg.SetLimit(concurrency)
	}

	return func() (err error) {
		var mu sync.Mutex
		ech := make(chan error, concurrency)
		finalize := func(err error) error {
			if errors.IsAny(err, io.EOF, context.Canceled, context.DeadlineExceeded) {
				err = nil
			}
			// 1. Wait for all in-flight sends to finish before half-closing.
			err = errors.Join(err, eg.Wait())
			close(ech)
			for e := range ech {
				if errors.IsNot(e, io.EOF, context.Canceled, context.DeadlineExceeded) {
					err = errors.Join(err, e)
				}
			}
			// 2. Half-close the send direction (fix 1). The server observes EOF on
			//    its Recv, finishes, and closes the response stream, which makes the
			//    receiver's stream.Recv return io.EOF so it can drain cleanly.
			mu.Lock()
			serr := stream.CloseSend()
			mu.Unlock()
			if errors.IsNot(serr, io.EOF, context.Canceled, context.DeadlineExceeded) {
				err = errors.Join(err, serr)
			}
			// 3. Wait for the receiver to drain to EOF. A cancelled ctx (caller
			//    deadline or cancel()) is the backstop so this cannot block forever.
			select {
			case rerr := <-recvDone:
				if errors.IsNot(rerr, io.EOF, context.Canceled, context.DeadlineExceeded) {
					err = errors.Join(err, rerr)
				}
			case <-ctx.Done():
				cancel()
				<-recvDone
			}
			return err
		}
		for {
			select {
			case <-egctx.Done():
				return finalize(egctx.Err())
			case err = <-ech:
				return finalize(err)
			default:
				data, ok := sendDataProvider()
				if !ok {
					return finalize(nil)
				}
				eg.Go(safety.RecoverFunc(func() (err error) {
					mu.Lock()
					err = stream.Send(data)
					mu.Unlock()
					if err != nil {
						select {
						case <-egctx.Done():
						case ech <- err:
						}
					}
					return nil
				}))
			}
		}
	}()
}
