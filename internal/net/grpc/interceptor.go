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

	"github.com/vdaas/vald/internal/safety"
	"google.golang.org/grpc"
)

type (
	UnaryServerInterceptor  = grpc.UnaryServerInterceptor
	StreamServerInterceptor = grpc.StreamServerInterceptor
)

var (
	UnaryInterceptor       = grpc.UnaryInterceptor
	ChainUnaryInterceptor  = grpc.ChainUnaryInterceptor
	StreamInterceptor      = grpc.StreamInterceptor
	ChainStreamInterceptor = grpc.ChainStreamInterceptor
)

func RecoverInterceptor() UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		err = safety.RecoverWithoutPanicFunc(func() (err error) {
			resp, err = handler(ctx, req)
			return err
		})()
		return resp, err
	}
}

func RecoverStreamInterceptor() StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		return safety.RecoverWithoutPanicFunc(func() (err error) {
			return handler(srv, ss)
		})()
	}
}
