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
package service

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
)

// GatewayMock represents mock struct for Gateway.
type GatewayMock struct {
	Gateway
	ForwardedContextFunc     func(ctx context.Context, podName string) context.Context
	FromForwardedContextFunc func(ctx context.Context) string
	BroadCastFunc            func(ctx context.Context,
		f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error
	DoFunc func(ctx context.Context, target string,
		f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error)) (interface{}, error)
	DoMultiFunc func(ctx context.Context, targets []string,
		f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error
	GRPCClientFunc func() grpc.Client
}

// ForwardedContext calls ForwardedContextFunc object.
func (gm *GatewayMock) ForwardedContext(ctx context.Context, podName string) context.Context {
	return gm.ForwardedContextFunc(ctx, podName)
}

// FromForwardedContext calls FromForwardedContextFunc object.
func (gm *GatewayMock) FromForwardedContext(ctx context.Context) string {
	return gm.FromForwardedContextFunc(ctx)
}

// BroadCast calls BroadCastFunc object.
func (gm *GatewayMock) BroadCast(ctx context.Context,
	f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) error,
) error {
	return gm.BroadCastFunc(ctx, f)
}

// Do calls DoFunc object.
func (gm *GatewayMock) Do(ctx context.Context, target string,
	f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) (interface{}, error),
) (interface{}, error) {
	return gm.DoFunc(ctx, target, f)
}

// DoMulti calls DoMultiFunc object.
func (gm *GatewayMock) DoMulti(ctx context.Context, targets []string,
	f func(_ context.Context, _ string, _ vald.ClientWithMirror, _ ...grpc.CallOption) error,
) error {
	return gm.DoMultiFunc(ctx, targets, f)
}

// GRPCClient calls GRPCClientFunc object.
func (gm *GatewayMock) GRPCClient() grpc.Client {
	return gm.GRPCClientFunc()
}
