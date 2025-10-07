//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	for {
		select {
		case <-ctx.Done():
			return finalize()
		default:
			data, err := stream.Recv()
			if err != nil {
				if err != io.EOF && !errors.Is(err, io.EOF) {
					err = errors.Wrap(err, "BidirectionalStream Recv returned error")
					emu.Lock()
					errs = append(errs, err)
					emu.Unlock()
					log.Errorf("failed to receive stream message: %v", err)
				}
				return finalize()
			}
			eg.Go(safety.RecoverWithoutPanicFunc(func() (err error) {
				id := atomic.AddUint64(&cnt, 1)
				ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s/BidirectionalStream/stream-%020d", apiName, id))
				defer func() {
					if sspan != nil {
						sspan.End()
					}
				}()
				res, err := handle(ctx, data)
				if err != nil {
					st, msg, err := status.ParseError(err, codes.Internal, fmt.Sprintf("failed to parse BidirectionalStream id= %020d gRPC error response", id))
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
						sspan.SetStatus(trace.StatusError, msg)
					}
					code := st.Code()
					if err != nil && st != nil &&
						code != codes.Canceled &&
						code != codes.DeadlineExceeded &&
						code != codes.InvalidArgument &&
						code != codes.NotFound &&
						code != codes.OK &&
						code != codes.Unimplemented {
						runtime.Gosched()
						log.Error(err)
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
					st, msg, err := status.ParseError(err, codes.Internal, fmt.Sprintf("failed to parse BidirectionalStream.SendMsg id= %020d gRPC error response", id),
						&errdetails.RequestInfo{
							RequestId:   fmt.Sprintf("%s/BidirectionalStream/stream-%020d/SendMsg", apiName, id),
							ServingData: errdetails.Serialize(res),
						})
					if sspan != nil {
						sspan.RecordError(err)
						sspan.SetAttributes(trace.FromGRPCStatus(st.Code(), msg)...)
						sspan.SetStatus(trace.StatusError, msg)
					}
					return err
				}
				return nil
			}))
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
	eg, ctx := errgroup.New(ctx)
	if concurrency > 0 {
		eg.SetLimit(concurrency)
	}
	eg.Go(safety.RecoverFunc(func() (err error) {
		for {
			select {
			case <-ctx.Done():
				err = ctx.Err()
				if errors.IsNot(err, io.EOF, context.Canceled, context.DeadlineExceeded) {
					return err
				}
				return nil
			default:
				res, err := stream.Recv()
				if errors.IsAny(err, io.EOF, context.Canceled, context.DeadlineExceeded) {
					cancel()
					return nil
				}
				if !callBack(res, err) {
					cancel()
					return nil
				}
			}
		}
	}))

	return func() (err error) {
		var mu sync.Mutex
		ech := make(chan error, concurrency)
		finalize := func(err error) (serr error) {
			if errors.IsAny(err, io.EOF, context.Canceled, context.DeadlineExceeded) {
				err = nil
			}
			cancel()
			err = errors.Join(err, eg.Wait())
			close(ech)
			for e := range ech {
				if errors.IsNot(e, io.EOF, context.Canceled, context.DeadlineExceeded) {
					err = errors.Join(err, e)
				}
			}
			mu.Lock()
			serr = stream.CloseSend()
			mu.Unlock()
			if errors.IsAny(serr, io.EOF, context.Canceled, context.DeadlineExceeded) {
				serr = nil
			}
			if err != nil {
				return errors.Join(err, serr)
			}
			return serr
		}
		for {
			select {
			case <-ctx.Done():
				return finalize(ctx.Err())
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
						case <-ctx.Done():
						case ech <- err:
						}
					}
					return nil
				}))
			}
		}
	}()
}
