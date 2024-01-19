// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

type gatewayMock struct {
	service.Gateway

	StartFunc                func(ctx context.Context) (<-chan error, error)
	ForwardedContextFunc     func(ctx context.Context, podName string) context.Context
	FromForwardedContextFunc func(ctx context.Context) string
	BroadCastFunc            func(ctx context.Context,
		f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error
	DoMultiFunc func(ctx context.Context, targets []string,
		f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error
}

func (gm *gatewayMock) ForwardedContext(ctx context.Context, podName string) context.Context {
	return gm.ForwardedContextFunc(ctx, podName)
}

func (gm *gatewayMock) FromForwardedContext(ctx context.Context) string {
	return gm.FromForwardedContextFunc(ctx)
}

func (gm *gatewayMock) BroadCast(ctx context.Context,
	f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error,
) error {
	return gm.BroadCastFunc(ctx, f)
}

func (gm *gatewayMock) DoMulti(ctx context.Context, targets []string,
	f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error,
) error {
	return gm.DoMultiFunc(ctx, targets, f)
}
