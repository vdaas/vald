// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/pkg/gateway/mirror/service"
)

type mockClient struct {
	vald.ClientWithMirror

	InsertFunc            func(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	UpdateFunc            func(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	UpsertFunc            func(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	RemoveFunc            func(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	RemoveByTimestampFunc func(ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

func (m *mockClient) Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return m.InsertFunc(ctx, in, opts...)
}

func (m *mockClient) Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return m.UpdateFunc(ctx, in, opts...)
}

func (m *mockClient) Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return m.UpsertFunc(ctx, in, opts...)
}

func (m *mockClient) Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return m.RemoveFunc(ctx, in, opts...)
}

func (m *mockClient) RemoveByTimestamp(ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	return m.RemoveByTimestampFunc(ctx, in, opts...)
}

type mockGateway struct {
	service.Gateway

	StartFunc                func(ctx context.Context) (<-chan error, error)
	ForwardedContextFunc     func(ctx context.Context, podName string) context.Context
	FromForwardedContextFunc func(ctx context.Context) string
	BroadCastFunc            func(ctx context.Context,
		f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error
	DoMultiFunc func(ctx context.Context, targets []string,
		f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error
}

func (m *mockGateway) ForwardedContext(ctx context.Context, podName string) context.Context {
	return m.ForwardedContextFunc(ctx, podName)
}

func (m *mockGateway) FromForwardedContext(ctx context.Context) string {
	return m.FromForwardedContextFunc(ctx)
}

func (m *mockGateway) BroadCast(ctx context.Context,
	f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error,
) error {
	return m.BroadCastFunc(ctx, f)
}

func (m *mockGateway) DoMulti(ctx context.Context, targets []string,
	f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error,
) error {
	return m.DoMultiFunc(ctx, targets, f)
}
