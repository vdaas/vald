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

package recover

import (
	"context"
	"runtime"
	"runtime/debug"

	"github.com/vdaas/vald/internal/conv"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/net/grpc"
)

// Interceptor returns a unary server interceptor with zero-alloc open-coded recover.
func Interceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = handlePanic(r, err)
			}
		}()
		return handler(ctx, req)
	}
}

// StreamInterceptor returns a stream server interceptor with zero-alloc open-coded recover.
func StreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv any,
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = handlePanic(r, err)
			}
		}()
		return handler(srv, ss)
	}
}

func handlePanic(r any, origErr error) error {
	ds := debug.Stack()
	infoStr := info.Get().String() + "\r\n" + conv.Btoa(ds)
	log.Warnf("gRPC server panic recovered: %#v\ninfo:\n%s", r, infoStr)

	var err error
	switch x := r.(type) {
	case runtime.Error:
		err = errors.ErrRuntimeError(origErr, x)
	case *string:
		if x != nil {
			err = errors.ErrPanicString(origErr, *x)
		} else {
			err = errors.ErrPanicRecovered(origErr, x)
		}
	case string:
		err = errors.ErrPanicString(origErr, x)
	case error:
		err = errors.Join(origErr, x)
	default:
		err = errors.ErrPanicRecovered(origErr, x)
	}
	if err != nil {
		log.Errorf("gRPC server panic recovered error:\t%v\ninfo:\n%s\nrecovered:\t%#v", err, infoStr, r)
	}
	return err
}
