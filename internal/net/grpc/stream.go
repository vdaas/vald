//
// Copyright (C) 2019 kpango (Yusuke Kato)
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

// Package grpc provides generic functionallity for grpc
package grpc

import (
	"context"
	"io"

	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/log"
	"google.golang.org/grpc"
)

type ServerStream grpc.ServerStream

func BidirectionalStream(stream grpc.ServerStream,
	newData func() interface{},
	f func(context.Context, interface{}) (interface{}, error)) (err error) {
	ctx, cancel := context.WithCancel(stream.Context())
	eg, ctx := errgroup.New(ctx)
	var count uint32
	for {
		select {
		case <-ctx.Done():
			log.Error("canceled")
			err = eg.Wait()
			log.Errorf("finished %s", err.Error())
			return err
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
				eg.Go(func() (err error) {
					var res interface{}
					res, err = f(ctx, data)
					if err != nil {
						cancel()
						// runtime.Gosched()
						// stream.SendMsg(nil)
						return err
					}
					return stream.SendMsg(res)
				})
			}
		}
	}
}
