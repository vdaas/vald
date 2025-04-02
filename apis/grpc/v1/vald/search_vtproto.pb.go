//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package vald

import (
	context "context"

	payload "github.com/vdaas/vald/apis/grpc/v1/payload"
	codes "github.com/vdaas/vald/internal/net/grpc/codes"
	status "github.com/vdaas/vald/internal/net/grpc/status"
	grpc "google.golang.org/grpc"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SearchClient is the client API for Search service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchClient interface {
	// Overview
	// Search RPC is the method to search vector(s) similar to the request vector.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	Search(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error)
	// Overview
	// SearchByID RPC is the method to search similar vectors using a user-defined vector ID.<br>
	// The vector with the same requested ID should be indexed into the `vald-agent` before searching.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	SearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)
	// Overview
	// StreamSearch RPC is the method to search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server.
	// Each Search request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	StreamSearch(ctx context.Context, opts ...grpc.CallOption) (Search_StreamSearchClient, error)
	// Overview
	// StreamSearchByID RPC is the method to search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server.
	// Each SearchByID request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	StreamSearchByID(ctx context.Context, opts ...grpc.CallOption) (Search_StreamSearchByIDClient, error)
	// Overview
	// MultiSearch RPC is the method to search vectors with multiple vectors in **1** request.
	//
	// <div class="notice">
	// gRPC has a message size limitation.<br>
	// Please be careful that the size of the request exceeds the limit.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	//
	//	The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	MultiSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	// Overview
	// MultiSearchByID RPC is the method to search vectors with multiple IDs in **1** request.
	//
	// <div class="notice">
	// gRPC has a message size limitation.<br>
	// Please be careful that the size of the request exceeds the limit.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	MultiSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	// Overview
	// LinearSearch RPC is the method to linear search vector(s) similar to the request vector.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	LinearSearch(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error)
	// Overview
	// LinearSearchByID RPC is the method to linear search similar vectors using a user-defined vector ID.<br>
	// The vector with the same requested ID should be indexed into the `vald-agent` before searching.
	// You will get a `NOT_FOUND` error if the vector isn't stored.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	LinearSearchByID(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)
	// Overview
	// StreamLinearSearch RPC is the method to linear search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server.
	// Each LinearSearch request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	StreamLinearSearch(ctx context.Context, opts ...grpc.CallOption) (Search_StreamLinearSearchClient, error)
	// Overview
	//
	//	StreamLinearSearchByID RPC is the method to linear search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	//
	// Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server.
	// Each LinearSearchByID request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	StreamLinearSearchByID(ctx context.Context, opts ...grpc.CallOption) (Search_StreamLinearSearchByIDClient, error)
	// Overview
	// MultiLinearSearch RPC is the method to linear search vectors with multiple vectors in **1** request.
	//
	// <div class="notice">
	// gRPC has a message size limitation.<br>
	// Please be careful that the size of the request exceeds the limit.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	//
	//	The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	MultiLinearSearch(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	// Overview
	// MultiLinearSearchByID RPC is the method to linear search vectors with multiple IDs in **1** request.
	//
	// <div class="notice">
	// gRPC has a message size limitation.<br>
	// Please be careful that the size of the request exceeds the limit.
	// </div>
	// // ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	MultiLinearSearchByID(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
}

type searchClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchClient(cc grpc.ClientConnInterface) SearchClient {
	return &searchClient{cc}
}

func (c *searchClient) Search(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/Search", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) SearchByID(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/SearchByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) StreamSearch(
	ctx context.Context, opts ...grpc.CallOption,
) (Search_StreamSearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Search_ServiceDesc.Streams[0], "/vald.v1.Search/StreamSearch", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchStreamSearchClient{stream}
	return x, nil
}

type Search_StreamSearchClient interface {
	Send(*payload.Search_Request) error
	Recv() (*payload.Search_StreamResponse, error)
	grpc.ClientStream
}

type searchStreamSearchClient struct {
	grpc.ClientStream
}

func (x *searchStreamSearchClient) Send(m *payload.Search_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchStreamSearchClient) Recv() (*payload.Search_StreamResponse, error) {
	m := new(payload.Search_StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) StreamSearchByID(
	ctx context.Context, opts ...grpc.CallOption,
) (Search_StreamSearchByIDClient, error) {
	stream, err := c.cc.NewStream(ctx, &Search_ServiceDesc.Streams[1], "/vald.v1.Search/StreamSearchByID", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchStreamSearchByIDClient{stream}
	return x, nil
}

type Search_StreamSearchByIDClient interface {
	Send(*payload.Search_IDRequest) error
	Recv() (*payload.Search_StreamResponse, error)
	grpc.ClientStream
}

type searchStreamSearchByIDClient struct {
	grpc.ClientStream
}

func (x *searchStreamSearchByIDClient) Send(m *payload.Search_IDRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchStreamSearchByIDClient) Recv() (*payload.Search_StreamResponse, error) {
	m := new(payload.Search_StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) MultiSearch(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	out := new(payload.Search_Responses)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/MultiSearch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) MultiSearchByID(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	out := new(payload.Search_Responses)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/MultiSearchByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) LinearSearch(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/LinearSearch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) LinearSearchByID(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/LinearSearchByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) StreamLinearSearch(
	ctx context.Context, opts ...grpc.CallOption,
) (Search_StreamLinearSearchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Search_ServiceDesc.Streams[2], "/vald.v1.Search/StreamLinearSearch", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchStreamLinearSearchClient{stream}
	return x, nil
}

type Search_StreamLinearSearchClient interface {
	Send(*payload.Search_Request) error
	Recv() (*payload.Search_StreamResponse, error)
	grpc.ClientStream
}

type searchStreamLinearSearchClient struct {
	grpc.ClientStream
}

func (x *searchStreamLinearSearchClient) Send(m *payload.Search_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchStreamLinearSearchClient) Recv() (*payload.Search_StreamResponse, error) {
	m := new(payload.Search_StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) StreamLinearSearchByID(
	ctx context.Context, opts ...grpc.CallOption,
) (Search_StreamLinearSearchByIDClient, error) {
	stream, err := c.cc.NewStream(ctx, &Search_ServiceDesc.Streams[3], "/vald.v1.Search/StreamLinearSearchByID", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchStreamLinearSearchByIDClient{stream}
	return x, nil
}

type Search_StreamLinearSearchByIDClient interface {
	Send(*payload.Search_IDRequest) error
	Recv() (*payload.Search_StreamResponse, error)
	grpc.ClientStream
}

type searchStreamLinearSearchByIDClient struct {
	grpc.ClientStream
}

func (x *searchStreamLinearSearchByIDClient) Send(m *payload.Search_IDRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchStreamLinearSearchByIDClient) Recv() (*payload.Search_StreamResponse, error) {
	m := new(payload.Search_StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchClient) MultiLinearSearch(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	out := new(payload.Search_Responses)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/MultiLinearSearch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchClient) MultiLinearSearchByID(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	out := new(payload.Search_Responses)
	err := c.cc.Invoke(ctx, "/vald.v1.Search/MultiLinearSearchByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchServer is the server API for Search service.
// All implementations must embed UnimplementedSearchServer
// for forward compatibility
type SearchServer interface {
	// Overview
	// Search RPC is the method to search vector(s) similar to the request vector.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	Search(context.Context, *payload.Search_Request) (*payload.Search_Response, error)
	// Overview
	// SearchByID RPC is the method to search similar vectors using a user-defined vector ID.<br>
	// The vector with the same requested ID should be indexed into the `vald-agent` before searching.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	SearchByID(context.Context, *payload.Search_IDRequest) (*payload.Search_Response, error)
	// Overview
	// StreamSearch RPC is the method to search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server.
	// Each Search request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	StreamSearch(Search_StreamSearchServer) error
	// Overview
	// StreamSearchByID RPC is the method to search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server.
	// Each SearchByID request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	StreamSearchByID(Search_StreamSearchByIDServer) error
	// Overview
	// MultiSearch RPC is the method to search vectors with multiple vectors in **1** request.
	//
	// <div class="notice">
	// gRPC has a message size limitation.<br>
	// Please be careful that the size of the request exceeds the limit.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	//
	//	The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	MultiSearch(context.Context, *payload.Search_MultiRequest) (*payload.Search_Responses, error)
	// Overview
	// MultiSearchByID RPC is the method to search vectors with multiple IDs in **1** request.
	//
	// <div class="notice">
	// gRPC has a message size limitation.<br>
	// Please be careful that the size of the request exceeds the limit.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	MultiSearchByID(context.Context, *payload.Search_MultiIDRequest) (*payload.Search_Responses, error)
	// Overview
	// LinearSearch RPC is the method to linear search vector(s) similar to the request vector.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	LinearSearch(context.Context, *payload.Search_Request) (*payload.Search_Response, error)
	// Overview
	// LinearSearchByID RPC is the method to linear search similar vectors using a user-defined vector ID.<br>
	// The vector with the same requested ID should be indexed into the `vald-agent` before searching.
	// You will get a `NOT_FOUND` error if the vector isn't stored.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	LinearSearchByID(context.Context, *payload.Search_IDRequest) (*payload.Search_Response, error)
	// Overview
	// StreamLinearSearch RPC is the method to linear search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server.
	// Each LinearSearch request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	StreamLinearSearch(Search_StreamLinearSearchServer) error
	// Overview
	//
	//	StreamLinearSearchByID RPC is the method to linear search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	//
	// Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server.
	// Each LinearSearchByID request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	StreamLinearSearchByID(Search_StreamLinearSearchByIDServer) error
	// Overview
	// MultiLinearSearch RPC is the method to linear search vectors with multiple vectors in **1** request.
	//
	// <div class="notice">
	// gRPC has a message size limitation.<br>
	// Please be careful that the size of the request exceeds the limit.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	//
	//	The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                   | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
	MultiLinearSearch(context.Context, *payload.Search_MultiRequest) (*payload.Search_Responses, error)
	// Overview
	// MultiLinearSearchByID RPC is the method to linear search vectors with multiple IDs in **1** request.
	//
	// <div class="notice">
	// gRPC has a message size limitation.<br>
	// Please be careful that the size of the request exceeds the limit.
	// </div>
	// // ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                    | how to resolve                                                                           |
	// | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
	MultiLinearSearchByID(context.Context, *payload.Search_MultiIDRequest) (*payload.Search_Responses, error)
	mustEmbedUnimplementedSearchServer()
}

// UnimplementedSearchServer must be embedded to have forward compatible implementations.
type UnimplementedSearchServer struct{}

func (UnimplementedSearchServer) Search(
	context.Context, *payload.Search_Request,
) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Search not implemented")
}

func (UnimplementedSearchServer) SearchByID(
	context.Context, *payload.Search_IDRequest,
) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchByID not implemented")
}

func (UnimplementedSearchServer) StreamSearch(Search_StreamSearchServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearch not implemented")
}

func (UnimplementedSearchServer) StreamSearchByID(Search_StreamSearchByIDServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearchByID not implemented")
}

func (UnimplementedSearchServer) MultiSearch(
	context.Context, *payload.Search_MultiRequest,
) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiSearch not implemented")
}

func (UnimplementedSearchServer) MultiSearchByID(
	context.Context, *payload.Search_MultiIDRequest,
) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiSearchByID not implemented")
}

func (UnimplementedSearchServer) LinearSearch(
	context.Context, *payload.Search_Request,
) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LinearSearch not implemented")
}

func (UnimplementedSearchServer) LinearSearchByID(
	context.Context, *payload.Search_IDRequest,
) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LinearSearchByID not implemented")
}

func (UnimplementedSearchServer) StreamLinearSearch(Search_StreamLinearSearchServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamLinearSearch not implemented")
}

func (UnimplementedSearchServer) StreamLinearSearchByID(Search_StreamLinearSearchByIDServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamLinearSearchByID not implemented")
}

func (UnimplementedSearchServer) MultiLinearSearch(
	context.Context, *payload.Search_MultiRequest,
) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiLinearSearch not implemented")
}

func (UnimplementedSearchServer) MultiLinearSearchByID(
	context.Context, *payload.Search_MultiIDRequest,
) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiLinearSearchByID not implemented")
}
func (UnimplementedSearchServer) mustEmbedUnimplementedSearchServer() {}

// UnsafeSearchServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchServer will
// result in compilation errors.
type UnsafeSearchServer interface {
	mustEmbedUnimplementedSearchServer()
}

func RegisterSearchServer(s grpc.ServiceRegistrar, srv SearchServer) {
	s.RegisterService(&Search_ServiceDesc, srv)
}

func _Search_Search_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).Search(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/Search",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchServer).Search(ctx, req.(*payload.Search_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_SearchByID_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).SearchByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/SearchByID",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchServer).SearchByID(ctx, req.(*payload.Search_IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_StreamSearch_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(SearchServer).StreamSearch(&searchStreamSearchServer{stream})
}

type Search_StreamSearchServer interface {
	Send(*payload.Search_StreamResponse) error
	Recv() (*payload.Search_Request, error)
	grpc.ServerStream
}

type searchStreamSearchServer struct {
	grpc.ServerStream
}

func (x *searchStreamSearchServer) Send(m *payload.Search_StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchStreamSearchServer) Recv() (*payload.Search_Request, error) {
	m := new(payload.Search_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Search_StreamSearchByID_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(SearchServer).StreamSearchByID(&searchStreamSearchByIDServer{stream})
}

type Search_StreamSearchByIDServer interface {
	Send(*payload.Search_StreamResponse) error
	Recv() (*payload.Search_IDRequest, error)
	grpc.ServerStream
}

type searchStreamSearchByIDServer struct {
	grpc.ServerStream
}

func (x *searchStreamSearchByIDServer) Send(m *payload.Search_StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchStreamSearchByIDServer) Recv() (*payload.Search_IDRequest, error) {
	m := new(payload.Search_IDRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Search_MultiSearch_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).MultiSearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/MultiSearch",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchServer).MultiSearch(ctx, req.(*payload.Search_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_MultiSearchByID_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_MultiIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).MultiSearchByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/MultiSearchByID",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchServer).MultiSearchByID(ctx, req.(*payload.Search_MultiIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_LinearSearch_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).LinearSearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/LinearSearch",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchServer).LinearSearch(ctx, req.(*payload.Search_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_LinearSearchByID_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).LinearSearchByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/LinearSearchByID",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchServer).LinearSearchByID(ctx, req.(*payload.Search_IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_StreamLinearSearch_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(SearchServer).StreamLinearSearch(&searchStreamLinearSearchServer{stream})
}

type Search_StreamLinearSearchServer interface {
	Send(*payload.Search_StreamResponse) error
	Recv() (*payload.Search_Request, error)
	grpc.ServerStream
}

type searchStreamLinearSearchServer struct {
	grpc.ServerStream
}

func (x *searchStreamLinearSearchServer) Send(m *payload.Search_StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchStreamLinearSearchServer) Recv() (*payload.Search_Request, error) {
	m := new(payload.Search_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Search_StreamLinearSearchByID_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(SearchServer).StreamLinearSearchByID(&searchStreamLinearSearchByIDServer{stream})
}

type Search_StreamLinearSearchByIDServer interface {
	Send(*payload.Search_StreamResponse) error
	Recv() (*payload.Search_IDRequest, error)
	grpc.ServerStream
}

type searchStreamLinearSearchByIDServer struct {
	grpc.ServerStream
}

func (x *searchStreamLinearSearchByIDServer) Send(m *payload.Search_StreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchStreamLinearSearchByIDServer) Recv() (*payload.Search_IDRequest, error) {
	m := new(payload.Search_IDRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Search_MultiLinearSearch_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).MultiLinearSearch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/MultiLinearSearch",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchServer).MultiLinearSearch(ctx, req.(*payload.Search_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Search_MultiLinearSearchByID_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_MultiIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchServer).MultiLinearSearchByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/vald.v1.Search/MultiLinearSearchByID",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchServer).MultiLinearSearchByID(ctx, req.(*payload.Search_MultiIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Search_ServiceDesc is the grpc.ServiceDesc for Search service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Search_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "vald.v1.Search",
	HandlerType: (*SearchServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Search",
			Handler:    _Search_Search_Handler,
		},
		{
			MethodName: "SearchByID",
			Handler:    _Search_SearchByID_Handler,
		},
		{
			MethodName: "MultiSearch",
			Handler:    _Search_MultiSearch_Handler,
		},
		{
			MethodName: "MultiSearchByID",
			Handler:    _Search_MultiSearchByID_Handler,
		},
		{
			MethodName: "LinearSearch",
			Handler:    _Search_LinearSearch_Handler,
		},
		{
			MethodName: "LinearSearchByID",
			Handler:    _Search_LinearSearchByID_Handler,
		},
		{
			MethodName: "MultiLinearSearch",
			Handler:    _Search_MultiLinearSearch_Handler,
		},
		{
			MethodName: "MultiLinearSearchByID",
			Handler:    _Search_MultiLinearSearchByID_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamSearch",
			Handler:       _Search_StreamSearch_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamSearchByID",
			Handler:       _Search_StreamSearchByID_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamLinearSearch",
			Handler:       _Search_StreamLinearSearch_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamLinearSearchByID",
			Handler:       _Search_StreamLinearSearchByID_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "v1/vald/search.proto",
}
