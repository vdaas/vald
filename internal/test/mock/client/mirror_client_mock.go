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
package client

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/v1/mirror"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/internal/net/grpc"
)

type MirrorClientMock struct {
	vald.Client
	mirror.MirrorClient

	InsertFunc            func(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	UpdateFunc            func(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	UpsertFunc            func(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	RemoveFunc            func(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	RemoveByTimestampFunc func(ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

func (mc *MirrorClientMock) Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return mc.InsertFunc(ctx, in, opts...)
}

func (mc *MirrorClientMock) Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return mc.UpdateFunc(ctx, in, opts...)
}

func (mc *MirrorClientMock) Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return mc.UpsertFunc(ctx, in, opts...)
}

func (mc *MirrorClientMock) Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return mc.RemoveFunc(ctx, in, opts...)
}

func (mc *MirrorClientMock) RemoveByTimestamp(ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	return mc.RemoveByTimestampFunc(ctx, in, opts...)
}
