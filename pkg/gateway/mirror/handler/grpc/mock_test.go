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
	InsertFunc       func(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamInsertFunc func(ctx context.Context, opts ...grpc.CallOption) (vald.Insert_StreamInsertClient, error)
	MultiInsertFunc  func(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)

	UpdateFunc       func(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamUpdateFunc func(ctx context.Context, opts ...grpc.CallOption) (vald.Update_StreamUpdateClient, error)
	MultiUpdateFunc  func(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)

	UpsertFunc       func(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	StreamUpsertFunc func(ctx context.Context, opts ...grpc.CallOption) (vald.Upsert_StreamUpsertClient, error)
	MultiUpsertFunc  func(ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)

	SearchFunc                 func(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error)
	SearchByIDFunc             func(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)
	StreamSearchFunc           func(ctx context.Context, opts ...grpc.CallOption) (vald.Search_StreamSearchClient, error)
	StreamSearchByIDFunc       func(ctx context.Context, opts ...grpc.CallOption) (vald.Search_StreamSearchByIDClient, error)
	MultiSearchFunc            func(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	MultiSearchByIDFunc        func(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	LinearSearchFunc           func(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error)
	LinearSearchByIDFunc       func(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)
	StreamLinearSearchFunc     func(ctx context.Context, opts ...grpc.CallOption) (vald.Search_StreamLinearSearchClient, error)
	StreamLinearSearchByIDFunc func(ctx context.Context, opts ...grpc.CallOption) (vald.Search_StreamLinearSearchByIDClient, error)
	MultiLinearSearchFunc      func(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	MultiLinearSearchByIDFunc  func(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)

	RemoveFunc            func(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	RemoveByTimestampFunc func(ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
	StreamRemoveFunc      func(ctx context.Context, opts ...grpc.CallOption) (vald.Remove_StreamRemoveClient, error)
	MultiRemoveFunc       func(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)

	ExistsFunc           func(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_ID, error)
	GetObjectFunc        func(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (*payload.Object_Vector, error)
	StreamGetObjectFunc  func(ctx context.Context, opts ...grpc.CallOption) (vald.Object_StreamGetObjectClient, error)
	StreamListObjectFunc func(ctx context.Context, opts ...grpc.CallOption) (vald.Object_StreamListObjectClient, error)

	RegisterFunc  func(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (*payload.Mirror_Targets, error)
	AdvertiseFunc func(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (*payload.Mirror_Targets, error)
}

func (m *mockClient) Insert(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return m.InsertFunc(ctx, in, opts...)
}

func (m *mockClient) StreamInsert(ctx context.Context, opts ...grpc.CallOption) (vald.Insert_StreamInsertClient, error) {
	return m.StreamInsertFunc(ctx, opts...)
}

func (m *mockClient) MultiInsert(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	return m.MultiInsertFunc(ctx, in, opts...)
}

func (m *mockClient) Update(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return m.UpdateFunc(ctx, in, opts...)
}

func (m *mockClient) StreamUpdate(ctx context.Context, opts ...grpc.CallOption) (vald.Update_StreamUpdateClient, error) {
	return m.StreamUpdateFunc(ctx, opts...)
}

func (m *mockClient) MultiUpdate(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	return m.MultiUpdateFunc(ctx, in, opts...)
}

func (m *mockClient) Upsert(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return m.UpsertFunc(ctx, in, opts...)
}

func (m *mockClient) StreamUpsert(ctx context.Context, opts ...grpc.CallOption) (vald.Upsert_StreamUpsertClient, error) {
	return m.StreamUpsertFunc(ctx, opts...)
}

func (m *mockClient) MultiUpsert(ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	return m.MultiUpsertFunc(ctx, in, opts...)
}

func (m *mockClient) Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	return m.SearchFunc(ctx, in, opts...)
}

func (m *mockClient) SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	return m.SearchByIDFunc(ctx, in, opts...)
}

func (m *mockClient) StreamSearch(ctx context.Context, opts ...grpc.CallOption) (vald.Search_StreamSearchClient, error) {
	return m.StreamSearchFunc(ctx, opts...)
}

func (m *mockClient) StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (vald.Search_StreamSearchByIDClient, error) {
	return m.StreamSearchByIDFunc(ctx, opts...)
}

func (m *mockClient) MultiSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error) {
	return m.MultiSearchFunc(ctx, in, opts...)
}

func (m *mockClient) MultiSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error) {
	return m.MultiSearchByIDFunc(ctx, in, opts...)
}

func (m *mockClient) LinearSearch(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	return m.LinearSearchFunc(ctx, in, opts...)
}

func (m *mockClient) LinearSearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error) {
	return m.LinearSearchByIDFunc(ctx, in, opts...)
}

func (m *mockClient) StreamLinearSearch(ctx context.Context, opts ...grpc.CallOption) (vald.Search_StreamLinearSearchClient, error) {
	return m.StreamLinearSearchFunc(ctx, opts...)
}

func (m *mockClient) StreamLinearSearchByID(ctx context.Context, opts ...grpc.CallOption) (vald.Search_StreamLinearSearchByIDClient, error) {
	return m.StreamLinearSearchByIDFunc(ctx, opts...)
}

func (m *mockClient) MultiLinearSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error) {
	return m.MultiLinearSearchFunc(ctx, in, opts...)
}

func (m *mockClient) MultiLinearSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error) {
	return m.MultiLinearSearchByIDFunc(ctx, in, opts...)
}

func (m *mockClient) Remove(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error) {
	return m.RemoveFunc(ctx, in, opts...)
}

func (m *mockClient) RemoveByTimestamp(ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	return m.RemoveByTimestamp(ctx, in, opts...)
}

func (m *mockClient) StreamRemove(ctx context.Context, opts ...grpc.CallOption) (vald.Remove_StreamRemoveClient, error) {
	return m.StreamRemoveFunc(ctx, opts...)
}

func (m *mockClient) MultiRemove(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error) {
	return m.MultiRemoveFunc(ctx, in, opts...)
}

func (m *mockClient) Exists(ctx context.Context, in *payload.Object_ID, opts ...grpc.CallOption) (*payload.Object_ID, error) {
	return m.ExistsFunc(ctx, in, opts...)
}

func (m *mockClient) GetObject(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (*payload.Object_Vector, error) {
	return m.GetObjectFunc(ctx, in, opts...)
}

func (m *mockClient) StreamGetObject(ctx context.Context, opts ...grpc.CallOption) (vald.Object_StreamGetObjectClient, error) {
	return m.StreamGetObjectFunc(ctx, opts...)
}

func (m *mockClient) StreamListObject(ctx context.Context, in *payload.Object_List_Request, opts ...grpc.CallOption) (vald.Object_StreamListObjectClient, error) {
	return m.StreamListObjectFunc(ctx, opts...)
}

func (m *mockClient) Register(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (*payload.Mirror_Targets, error) {
	return m.RegisterFunc(ctx, in)
}

func (m *mockClient) Advertise(ctx context.Context, in *payload.Mirror_Targets, opts ...grpc.CallOption) (*payload.Mirror_Targets, error) {
	return m.AdvertiseFunc(ctx, in)
}

var _ vald.ClientWithMirror = (*mockClient)(nil)

type mockGateway struct {
	StartFunc                func(ctx context.Context) (<-chan error, error)
	ForwardedContextFunc     func(ctx context.Context, podName string) context.Context
	FromForwardedContextFunc func(ctx context.Context) string
	BroadCastFunc            func(ctx context.Context,
		f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error
	DoFunc func(ctx context.Context, target string,
		f func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error)) (interface{}, error)
	DoMultiFunc func(ctx context.Context, targets []string,
		f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error) error
	GRPCClientFunc func() grpc.Client
}

func (m *mockGateway) Start(ctx context.Context) (<-chan error, error) {
	return m.StartFunc(ctx)
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

func (m *mockGateway) Do(ctx context.Context, target string,
	f func(ctx context.Context, vc vald.ClientWithMirror, copts ...grpc.CallOption) (interface{}, error),
) (interface{}, error) {
	return m.DoFunc(ctx, target, f)
}

func (m *mockGateway) DoMulti(ctx context.Context, targets []string,
	f func(ctx context.Context, target string, vc vald.ClientWithMirror, copts ...grpc.CallOption) error,
) error {
	return m.DoMultiFunc(ctx, targets, f)
}

func (m *mockGateway) GRPCClient() grpc.Client {
	return m.GRPCClientFunc()
}

var _ service.Gateway = (*mockGateway)(nil)
