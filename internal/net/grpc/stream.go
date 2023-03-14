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

	var mu sync.Mutex

	errMap := sync.Map{}

	finalize := func() error {
		var errs error
		err = eg.Wait()
		errMap.Range(func(_, e interface{}) bool {
			err, ok := e.(error)
			if !ok || err == nil {
				return true
			}
			if errs == nil {
				errs = err
			} else {
				errs = errors.Wrap(err, errs.Error())
			}
			return true
		})
		if errs == nil {
			return nil
		}
		st, msg, err := status.ParseError(errs, codes.Internal, "failed to parse BidirectionalStream final gRPC error response")
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
					log.Errorf("failed to receive stream message: %v", err)
					errMap.Store(err.Error(), err)
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
						errMap.Store(err.Error(), err)
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
			err = errors.Wrap(stream.CloseSend(), err.Error())
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
						return errors.Wrap(eg.Wait(), err.Error())
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
