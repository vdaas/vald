//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kou-m, rinx )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/safety"
	"google.golang.org/grpc"
)

type ServerStream grpc.ServerStream

func BidirectionalStream(stream ServerStream,
	newData func() interface{},
	f func(context.Context, interface{}) (interface{}, error)) (err error) {
	ctx := stream.Context()
	eg, ctx := errgroup.New(stream.Context())
	eg.Limitation(10)
	for {
		select {
		case <-ctx.Done():
			return eg.Wait()
		default:
			data := newData()
			err = stream.RecvMsg(data)
			if err != nil {
				if err == io.EOF {
					return eg.Wait()
				}
				return err
			}
			if data != nil {
				eg.Go(safety.RecoverFunc(func() (err error) {
					var res interface{}
					res, err = f(ctx, data)
					if err != nil {
						runtime.Gosched()
						return err
					}
					if res != nil {
						err = stream.SendMsg(res)
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
