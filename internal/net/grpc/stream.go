//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"io"
	"runtime"
	"sync"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"google.golang.org/grpc"
)

type ClientStream = grpc.ClientStream

// BidirectionalStream represents gRPC bidirectional stream server handler.
func BidirectionalStream(ctx context.Context, stream grpc.ServerStream,
	concurrency int,
	newData func() interface{},
	f func(context.Context, interface{}) (interface{}, error)) (err error) {
	eg, ctx := errgroup.New(ctx)
	if concurrency > 0 {
		eg.Limitation(concurrency)
	}

	var mu sync.Mutex

	errMap := sync.Map{}

	for {
		select {
		case <-ctx.Done():
			err = eg.Wait()
			if err != nil {
				log.Error(err)
				return err
			}
			return nil
		default:
			data := newData()
			err = stream.RecvMsg(data)
			if err != nil {
				if err == io.EOF {
					err = eg.Wait()
					if err != nil {
						log.Error(err)
						return err
					}
					var errs error
					errMap.Range(func(_, err interface{}) bool {
						errs = errors.Wrap(errs, err.(error).Error())
						return true
					})
					return errs
				}
				log.Error(err)
				return err
			}
			if data != nil {
				eg.Go(safety.RecoverWithoutPanicFunc(func() (err error) {
					var res interface{}
					res, err = f(ctx, data)
					if err != nil {
						runtime.Gosched()
						errMap.Store(err.Error(), err)
						return nil
					}
					if res != nil {
						mu.Lock()
						err = stream.SendMsg(res)
						mu.Unlock()
						if err != nil {
							runtime.Gosched()
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
func BidirectionalStreamClient(stream grpc.ClientStream,
	dataProvider, newData func() interface{},
	f func(interface{}, error)) (err error) {
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
				if err == io.EOF {
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
