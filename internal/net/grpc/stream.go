//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/errdetails"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/slices"
	"google.golang.org/grpc"
)

type (
	ClientStream = grpc.ClientStream
	ServerStream = grpc.ServerStream
)

// BidirectionalStream represents gRPC bidirectional stream server handler.
func BidirectionalStream[Q any, R any](ctx context.Context, stream ServerStream,
	concurrency int,
	f func(context.Context, *Q) (*R, error),
) (err error) {
	ctx, span := trace.StartSpan(stream.Context(), apiName+"/BidirectionalStream")
	defer func() {
		if span != nil {
			span.End()
		}
	}()
	eg, ctx := errgroup.New(ctx)
	if concurrency > 0 {
		eg.Limitation(concurrency)
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
		slices.RemoveDuplicates(errs, func(left, right error) bool {
			return left.Error() < right.Error()
		})
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
			data := new(Q)
			err = stream.RecvMsg(data)
			if err != nil {
				if err != io.EOF && !errors.Is(err, io.EOF) {
					err = errors.Wrap(err, "BidirectionalStream RecvMsg returned error")
					emu.Lock()
					errs = append(errs, err)
					emu.Unlock()
					log.Errorf("failed to receive stream message: %v", err)
				}
				return finalize()

			}
			if data != nil {
				eg.Go(safety.RecoverWithoutPanicFunc(func() (err error) {
					id := atomic.AddUint64(&cnt, 1)
					ctx, sspan := trace.StartSpan(ctx, fmt.Sprintf("%s/BidirectionalStream/stream-%020d", apiName, id))
					defer func() {
						if sspan != nil {
							sspan.End()
						}
					}()
					var res *R
					res, err = f(ctx, data)
					if err != nil {
						runtime.Gosched()
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
							log.Error(err)
						}
					}
					if res != nil {
						mu.Lock()
						err = stream.SendMsg(res)
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
					}
					return nil
				}))
			}
		}
	}
}

// BidirectionalStreamClient is gRPC client stream.
func BidirectionalStreamClient(stream ClientStream,
	dataProvider, newData func() interface{},
	f func(interface{}, error),
) (err error) {
	if stream == nil {
		return errors.ErrGRPCClientStreamNotFound
	}

	ctx, cancel := context.WithCancel(stream.Context())
	eg, ctx := errgroup.New(ctx)

	eg.Go(safety.RecoverFunc(func() (err error) {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				res := newData()
				err = stream.RecvMsg(res)
				if err == io.EOF || errors.Is(err, io.EOF) {
					cancel()
					return nil
				}
				f(res, err)
			}
		}
	}))

	defer func() {
		if err != nil {
			err = errors.Join(stream.CloseSend(), err)
		} else {
			err = stream.CloseSend()
		}
	}()

	return func() (err error) {
		for {
			select {
			case <-ctx.Done():
				return eg.Wait()
			default:
				data := dataProvider()
				if data == nil {
					err = stream.CloseSend()
					cancel()
					if err != nil {
						return errors.Join(eg.Wait(), err)
					}
					return eg.Wait()
				}

				err = stream.SendMsg(data)
				if err != nil {
					return err
				}
			}
		}
	}()
}
