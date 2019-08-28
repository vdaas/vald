// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package grpc provides generic functionallity for grpc
package grpc

import (
	"context"
	"io"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func BidirectionalStream(stream grpc.ServerStream,
	f func(context.Context, interface{}) (interface{}, error)) error {
	eg, ctx := errgroup.WithContext(stream.Context())
	for {
		select {
		case <-ctx.Done():
			return eg.Wait()
		default:
			var data interface{}
			err := stream.RecvMsg(&data)
			if err != nil {
				if err == io.EOF {
					return eg.Wait()
				}
				return err
			}

			eg.Go(func() error {
				res, err := f(ctx, data)
				if err != nil {
					return err
				}
				return stream.SendMsg(res)
			})
		}
	}
}
