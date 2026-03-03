//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

// SearchWithMetadataClient is the client API for SearchWithMetadata service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchWithMetadataClient interface {
	// Overview
	// SearchWithMetadata RPC is the method to search vector(s) similar to the request vector and to get metadata(s).
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
	SearchWithMetadata(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error)
	// Overview
	// SearchByIDWithMetadata RPC is the method to search similar vectors using a user-defined vector ID and to get metadata.<br>
	// The vector with the same requested ID should be indexed into the `vald-lb-gateway` before searching.
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
	SearchByIDWithMetadata(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)
	// Overview
	// StreamSearchWithMetadata RPC is the method to search vectors and to get metadata with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
	StreamSearchWithMetadata(ctx context.Context, opts ...grpc.CallOption) (SearchWithMetadata_StreamSearchWithMetadataClient, error)
	// Overview
	// StreamSearchByIDWithMetadata RPC is the method to search vectors and to get metadata with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
	StreamSearchByIDWithMetadata(ctx context.Context, opts ...grpc.CallOption) (SearchWithMetadata_StreamSearchByIDWithMetadataClient, error)
	// Overview
	// MultiSearchWithMetadata RPC is the method to search vectors and to get metadata with multiple vectors in **1** request.
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
	MultiSearchWithMetadata(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	// Overview
	// MultiSearchByIDWithMetadata RPC is the method to search vectors and to get metadata with multiple IDs in **1** request.
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
	MultiSearchByIDWithMetadata(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	// Overview
	// LinearSearchWithMetadata RPC is the method to linear search vector(s) similar to the request vector and to get metadata.
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
	LinearSearchWithMetadata(ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption) (*payload.Search_Response, error)
	// Overview
	// LinearSearchByIDWithMetadata RPC is the method to linear search similar vectors using a user-defined vector ID and to get metadata.<br>
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
	LinearSearchByIDWithMetadata(ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption) (*payload.Search_Response, error)
	// Overview
	// StreamLinearSearchWithMetadata RPC is the method to linear search vectors and to get metadata with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
	StreamLinearSearchWithMetadata(ctx context.Context, opts ...grpc.CallOption) (SearchWithMetadata_StreamLinearSearchWithMetadataClient, error)
	// Overview
	//
	//	StreamLinearSearchByIDWithMetadata RPC is the method to linear search vectors and to get metadata with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
	StreamLinearSearchByIDWithMetadata(ctx context.Context, opts ...grpc.CallOption) (SearchWithMetadata_StreamLinearSearchByIDWithMetadataClient, error)
	// Overview
	// MultiLinearSearchWithMetadata RPC is the method to linear search vectors and to get metadata with multiple vectors in **1** request.
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
	MultiLinearSearchWithMetadata(ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
	// Overview
	// MultiLinearSearchByIDWithMetadata RPC is the method to linear search vectors and to get metadata with multiple IDs in **1** request.
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
	MultiLinearSearchByIDWithMetadata(ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption) (*payload.Search_Responses, error)
}

type searchWithMetadataClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchWithMetadataClient(cc grpc.ClientConnInterface) SearchWithMetadataClient {
	return &searchWithMetadataClient{cc}
}

func (c *searchWithMetadataClient) SearchWithMetadata(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/meta.v1.SearchWithMetadata/SearchWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchWithMetadataClient) SearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/meta.v1.SearchWithMetadata/SearchByIDWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchWithMetadataClient) StreamSearchWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (SearchWithMetadata_StreamSearchWithMetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &SearchWithMetadata_ServiceDesc.Streams[0], "/meta.v1.SearchWithMetadata/StreamSearchWithMetadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchWithMetadataStreamSearchWithMetadataClient{stream}
	return x, nil
}

type SearchWithMetadata_StreamSearchWithMetadataClient interface {
	Send(*payload.Search_Request) error
	Recv() (*payload.Search_StreamResponse, error)
	grpc.ClientStream
}

type searchWithMetadataStreamSearchWithMetadataClient struct {
	grpc.ClientStream
}

func (x *searchWithMetadataStreamSearchWithMetadataClient) Send(m *payload.Search_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchWithMetadataStreamSearchWithMetadataClient) Recv() (
	*payload.Search_StreamResponse,
	error,
) {
	m := new(payload.Search_StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchWithMetadataClient) StreamSearchByIDWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (SearchWithMetadata_StreamSearchByIDWithMetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &SearchWithMetadata_ServiceDesc.Streams[1], "/meta.v1.SearchWithMetadata/StreamSearchByIDWithMetadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchWithMetadataStreamSearchByIDWithMetadataClient{stream}
	return x, nil
}

type SearchWithMetadata_StreamSearchByIDWithMetadataClient interface {
	Send(*payload.Search_IDRequest) error
	Recv() (*payload.Search_StreamResponse, error)
	grpc.ClientStream
}

type searchWithMetadataStreamSearchByIDWithMetadataClient struct {
	grpc.ClientStream
}

func (x *searchWithMetadataStreamSearchByIDWithMetadataClient) Send(
	m *payload.Search_IDRequest,
) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchWithMetadataStreamSearchByIDWithMetadataClient) Recv() (
	*payload.Search_StreamResponse,
	error,
) {
	m := new(payload.Search_StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchWithMetadataClient) MultiSearchWithMetadata(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	out := new(payload.Search_Responses)
	err := c.cc.Invoke(ctx, "/meta.v1.SearchWithMetadata/MultiSearchWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchWithMetadataClient) MultiSearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	out := new(payload.Search_Responses)
	err := c.cc.Invoke(ctx, "/meta.v1.SearchWithMetadata/MultiSearchByIDWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchWithMetadataClient) LinearSearchWithMetadata(
	ctx context.Context, in *payload.Search_Request, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/meta.v1.SearchWithMetadata/LinearSearchWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchWithMetadataClient) LinearSearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_IDRequest, opts ...grpc.CallOption,
) (*payload.Search_Response, error) {
	out := new(payload.Search_Response)
	err := c.cc.Invoke(ctx, "/meta.v1.SearchWithMetadata/LinearSearchByIDWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchWithMetadataClient) StreamLinearSearchWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (SearchWithMetadata_StreamLinearSearchWithMetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &SearchWithMetadata_ServiceDesc.Streams[2], "/meta.v1.SearchWithMetadata/StreamLinearSearchWithMetadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchWithMetadataStreamLinearSearchWithMetadataClient{stream}
	return x, nil
}

type SearchWithMetadata_StreamLinearSearchWithMetadataClient interface {
	Send(*payload.Search_Request) error
	Recv() (*payload.Search_StreamResponse, error)
	grpc.ClientStream
}

type searchWithMetadataStreamLinearSearchWithMetadataClient struct {
	grpc.ClientStream
}

func (x *searchWithMetadataStreamLinearSearchWithMetadataClient) Send(
	m *payload.Search_Request,
) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchWithMetadataStreamLinearSearchWithMetadataClient) Recv() (
	*payload.Search_StreamResponse,
	error,
) {
	m := new(payload.Search_StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchWithMetadataClient) StreamLinearSearchByIDWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (SearchWithMetadata_StreamLinearSearchByIDWithMetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &SearchWithMetadata_ServiceDesc.Streams[3], "/meta.v1.SearchWithMetadata/StreamLinearSearchByIDWithMetadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchWithMetadataStreamLinearSearchByIDWithMetadataClient{stream}
	return x, nil
}

type SearchWithMetadata_StreamLinearSearchByIDWithMetadataClient interface {
	Send(*payload.Search_IDRequest) error
	Recv() (*payload.Search_StreamResponse, error)
	grpc.ClientStream
}

type searchWithMetadataStreamLinearSearchByIDWithMetadataClient struct {
	grpc.ClientStream
}

func (x *searchWithMetadataStreamLinearSearchByIDWithMetadataClient) Send(
	m *payload.Search_IDRequest,
) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchWithMetadataStreamLinearSearchByIDWithMetadataClient) Recv() (
	*payload.Search_StreamResponse,
	error,
) {
	m := new(payload.Search_StreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *searchWithMetadataClient) MultiLinearSearchWithMetadata(
	ctx context.Context, in *payload.Search_MultiRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	out := new(payload.Search_Responses)
	err := c.cc.Invoke(ctx, "/meta.v1.SearchWithMetadata/MultiLinearSearchWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *searchWithMetadataClient) MultiLinearSearchByIDWithMetadata(
	ctx context.Context, in *payload.Search_MultiIDRequest, opts ...grpc.CallOption,
) (*payload.Search_Responses, error) {
	out := new(payload.Search_Responses)
	err := c.cc.Invoke(ctx, "/meta.v1.SearchWithMetadata/MultiLinearSearchByIDWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SearchWithMetadataServer is the server API for SearchWithMetadata service.
// All implementations must embed UnimplementedSearchWithMetadataServer
// for forward compatibility
type SearchWithMetadataServer interface {
	// Overview
	// SearchWithMetadata RPC is the method to search vector(s) similar to the request vector and to get metadata(s).
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
	SearchWithMetadata(context.Context, *payload.Search_Request) (*payload.Search_Response, error)
	// Overview
	// SearchByIDWithMetadata RPC is the method to search similar vectors using a user-defined vector ID and to get metadata.<br>
	// The vector with the same requested ID should be indexed into the `vald-lb-gateway` before searching.
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
	SearchByIDWithMetadata(context.Context, *payload.Search_IDRequest) (*payload.Search_Response, error)
	// Overview
	// StreamSearchWithMetadata RPC is the method to search vectors and to get metadata with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
	StreamSearchWithMetadata(SearchWithMetadata_StreamSearchWithMetadataServer) error
	// Overview
	// StreamSearchByIDWithMetadata RPC is the method to search vectors and to get metadata with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
	StreamSearchByIDWithMetadata(SearchWithMetadata_StreamSearchByIDWithMetadataServer) error
	// Overview
	// MultiSearchWithMetadata RPC is the method to search vectors and to get metadata with multiple vectors in **1** request.
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
	MultiSearchWithMetadata(context.Context, *payload.Search_MultiRequest) (*payload.Search_Responses, error)
	// Overview
	// MultiSearchByIDWithMetadata RPC is the method to search vectors and to get metadata with multiple IDs in **1** request.
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
	MultiSearchByIDWithMetadata(context.Context, *payload.Search_MultiIDRequest) (*payload.Search_Responses, error)
	// Overview
	// LinearSearchWithMetadata RPC is the method to linear search vector(s) similar to the request vector and to get metadata.
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
	LinearSearchWithMetadata(context.Context, *payload.Search_Request) (*payload.Search_Response, error)
	// Overview
	// LinearSearchByIDWithMetadata RPC is the method to linear search similar vectors using a user-defined vector ID and to get metadata.<br>
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
	LinearSearchByIDWithMetadata(context.Context, *payload.Search_IDRequest) (*payload.Search_Response, error)
	// Overview
	// StreamLinearSearchWithMetadata RPC is the method to linear search vectors and to get metadata with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
	StreamLinearSearchWithMetadata(SearchWithMetadata_StreamLinearSearchWithMetadataServer) error
	// Overview
	//
	//	StreamLinearSearchByIDWithMetadata RPC is the method to linear search vectors and to get metadata with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
	StreamLinearSearchByIDWithMetadata(SearchWithMetadata_StreamLinearSearchByIDWithMetadataServer) error
	// Overview
	// MultiLinearSearchWithMetadata RPC is the method to linear search vectors and to get metadata with multiple vectors in **1** request.
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
	MultiLinearSearchWithMetadata(context.Context, *payload.Search_MultiRequest) (*payload.Search_Responses, error)
	// Overview
	// MultiLinearSearchByIDWithMetadata RPC is the method to linear search vectors and to get metadata with multiple IDs in **1** request.
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
	MultiLinearSearchByIDWithMetadata(context.Context, *payload.Search_MultiIDRequest) (*payload.Search_Responses, error)
	mustEmbedUnimplementedSearchWithMetadataServer()
}

// UnimplementedSearchWithMetadataServer must be embedded to have forward compatible implementations.
type UnimplementedSearchWithMetadataServer struct{}

func (UnimplementedSearchWithMetadataServer) SearchWithMetadata(
	context.Context, *payload.Search_Request,
) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchWithMetadata not implemented")
}

func (UnimplementedSearchWithMetadataServer) SearchByIDWithMetadata(
	context.Context, *payload.Search_IDRequest,
) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchByIDWithMetadata not implemented")
}

func (UnimplementedSearchWithMetadataServer) StreamSearchWithMetadata(
	SearchWithMetadata_StreamSearchWithMetadataServer,
) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearchWithMetadata not implemented")
}

func (UnimplementedSearchWithMetadataServer) StreamSearchByIDWithMetadata(
	SearchWithMetadata_StreamSearchByIDWithMetadataServer,
) error {
	return status.Errorf(codes.Unimplemented, "method StreamSearchByIDWithMetadata not implemented")
}

func (UnimplementedSearchWithMetadataServer) MultiSearchWithMetadata(
	context.Context, *payload.Search_MultiRequest,
) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiSearchWithMetadata not implemented")
}

func (UnimplementedSearchWithMetadataServer) MultiSearchByIDWithMetadata(
	context.Context, *payload.Search_MultiIDRequest,
) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiSearchByIDWithMetadata not implemented")
}

func (UnimplementedSearchWithMetadataServer) LinearSearchWithMetadata(
	context.Context, *payload.Search_Request,
) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LinearSearchWithMetadata not implemented")
}

func (UnimplementedSearchWithMetadataServer) LinearSearchByIDWithMetadata(
	context.Context, *payload.Search_IDRequest,
) (*payload.Search_Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LinearSearchByIDWithMetadata not implemented")
}

func (UnimplementedSearchWithMetadataServer) StreamLinearSearchWithMetadata(
	SearchWithMetadata_StreamLinearSearchWithMetadataServer,
) error {
	return status.Errorf(codes.Unimplemented, "method StreamLinearSearchWithMetadata not implemented")
}

func (UnimplementedSearchWithMetadataServer) StreamLinearSearchByIDWithMetadata(
	SearchWithMetadata_StreamLinearSearchByIDWithMetadataServer,
) error {
	return status.Errorf(codes.Unimplemented, "method StreamLinearSearchByIDWithMetadata not implemented")
}

func (UnimplementedSearchWithMetadataServer) MultiLinearSearchWithMetadata(
	context.Context, *payload.Search_MultiRequest,
) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiLinearSearchWithMetadata not implemented")
}

func (UnimplementedSearchWithMetadataServer) MultiLinearSearchByIDWithMetadata(
	context.Context, *payload.Search_MultiIDRequest,
) (*payload.Search_Responses, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiLinearSearchByIDWithMetadata not implemented")
}
func (UnimplementedSearchWithMetadataServer) mustEmbedUnimplementedSearchWithMetadataServer() {}

// UnsafeSearchWithMetadataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchWithMetadataServer will
// result in compilation errors.
type UnsafeSearchWithMetadataServer interface {
	mustEmbedUnimplementedSearchWithMetadataServer()
}

func RegisterSearchWithMetadataServer(s grpc.ServiceRegistrar, srv SearchWithMetadataServer) {
	s.RegisterService(&SearchWithMetadata_ServiceDesc, srv)
}

func _SearchWithMetadata_SearchWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchWithMetadataServer).SearchWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.SearchWithMetadata/SearchWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchWithMetadataServer).SearchWithMetadata(ctx, req.(*payload.Search_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchWithMetadata_SearchByIDWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchWithMetadataServer).SearchByIDWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.SearchWithMetadata/SearchByIDWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchWithMetadataServer).SearchByIDWithMetadata(ctx, req.(*payload.Search_IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchWithMetadata_StreamSearchWithMetadata_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(SearchWithMetadataServer).StreamSearchWithMetadata(&searchWithMetadataStreamSearchWithMetadataServer{stream})
}

type SearchWithMetadata_StreamSearchWithMetadataServer interface {
	Send(*payload.Search_StreamResponse) error
	Recv() (*payload.Search_Request, error)
	grpc.ServerStream
}

type searchWithMetadataStreamSearchWithMetadataServer struct {
	grpc.ServerStream
}

func (x *searchWithMetadataStreamSearchWithMetadataServer) Send(
	m *payload.Search_StreamResponse,
) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchWithMetadataStreamSearchWithMetadataServer) Recv() (*payload.Search_Request, error) {
	m := new(payload.Search_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SearchWithMetadata_StreamSearchByIDWithMetadata_Handler(
	srv any, stream grpc.ServerStream,
) error {
	return srv.(SearchWithMetadataServer).StreamSearchByIDWithMetadata(&searchWithMetadataStreamSearchByIDWithMetadataServer{stream})
}

type SearchWithMetadata_StreamSearchByIDWithMetadataServer interface {
	Send(*payload.Search_StreamResponse) error
	Recv() (*payload.Search_IDRequest, error)
	grpc.ServerStream
}

type searchWithMetadataStreamSearchByIDWithMetadataServer struct {
	grpc.ServerStream
}

func (x *searchWithMetadataStreamSearchByIDWithMetadataServer) Send(
	m *payload.Search_StreamResponse,
) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchWithMetadataStreamSearchByIDWithMetadataServer) Recv() (
	*payload.Search_IDRequest,
	error,
) {
	m := new(payload.Search_IDRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SearchWithMetadata_MultiSearchWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchWithMetadataServer).MultiSearchWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.SearchWithMetadata/MultiSearchWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchWithMetadataServer).MultiSearchWithMetadata(ctx, req.(*payload.Search_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchWithMetadata_MultiSearchByIDWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_MultiIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchWithMetadataServer).MultiSearchByIDWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.SearchWithMetadata/MultiSearchByIDWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchWithMetadataServer).MultiSearchByIDWithMetadata(ctx, req.(*payload.Search_MultiIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchWithMetadata_LinearSearchWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchWithMetadataServer).LinearSearchWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.SearchWithMetadata/LinearSearchWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchWithMetadataServer).LinearSearchWithMetadata(ctx, req.(*payload.Search_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchWithMetadata_LinearSearchByIDWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_IDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchWithMetadataServer).LinearSearchByIDWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.SearchWithMetadata/LinearSearchByIDWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchWithMetadataServer).LinearSearchByIDWithMetadata(ctx, req.(*payload.Search_IDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchWithMetadata_StreamLinearSearchWithMetadata_Handler(
	srv any, stream grpc.ServerStream,
) error {
	return srv.(SearchWithMetadataServer).StreamLinearSearchWithMetadata(&searchWithMetadataStreamLinearSearchWithMetadataServer{stream})
}

type SearchWithMetadata_StreamLinearSearchWithMetadataServer interface {
	Send(*payload.Search_StreamResponse) error
	Recv() (*payload.Search_Request, error)
	grpc.ServerStream
}

type searchWithMetadataStreamLinearSearchWithMetadataServer struct {
	grpc.ServerStream
}

func (x *searchWithMetadataStreamLinearSearchWithMetadataServer) Send(
	m *payload.Search_StreamResponse,
) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchWithMetadataStreamLinearSearchWithMetadataServer) Recv() (
	*payload.Search_Request,
	error,
) {
	m := new(payload.Search_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SearchWithMetadata_StreamLinearSearchByIDWithMetadata_Handler(
	srv any, stream grpc.ServerStream,
) error {
	return srv.(SearchWithMetadataServer).StreamLinearSearchByIDWithMetadata(&searchWithMetadataStreamLinearSearchByIDWithMetadataServer{stream})
}

type SearchWithMetadata_StreamLinearSearchByIDWithMetadataServer interface {
	Send(*payload.Search_StreamResponse) error
	Recv() (*payload.Search_IDRequest, error)
	grpc.ServerStream
}

type searchWithMetadataStreamLinearSearchByIDWithMetadataServer struct {
	grpc.ServerStream
}

func (x *searchWithMetadataStreamLinearSearchByIDWithMetadataServer) Send(
	m *payload.Search_StreamResponse,
) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchWithMetadataStreamLinearSearchByIDWithMetadataServer) Recv() (
	*payload.Search_IDRequest,
	error,
) {
	m := new(payload.Search_IDRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SearchWithMetadata_MultiLinearSearchWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchWithMetadataServer).MultiLinearSearchWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.SearchWithMetadata/MultiLinearSearchWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchWithMetadataServer).MultiLinearSearchWithMetadata(ctx, req.(*payload.Search_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SearchWithMetadata_MultiLinearSearchByIDWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Search_MultiIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SearchWithMetadataServer).MultiLinearSearchByIDWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.SearchWithMetadata/MultiLinearSearchByIDWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(SearchWithMetadataServer).MultiLinearSearchByIDWithMetadata(ctx, req.(*payload.Search_MultiIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SearchWithMetadata_ServiceDesc is the grpc.ServiceDesc for SearchWithMetadata service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SearchWithMetadata_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "meta.v1.SearchWithMetadata",
	HandlerType: (*SearchWithMetadataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SearchWithMetadata",
			Handler:    _SearchWithMetadata_SearchWithMetadata_Handler,
		},
		{
			MethodName: "SearchByIDWithMetadata",
			Handler:    _SearchWithMetadata_SearchByIDWithMetadata_Handler,
		},
		{
			MethodName: "MultiSearchWithMetadata",
			Handler:    _SearchWithMetadata_MultiSearchWithMetadata_Handler,
		},
		{
			MethodName: "MultiSearchByIDWithMetadata",
			Handler:    _SearchWithMetadata_MultiSearchByIDWithMetadata_Handler,
		},
		{
			MethodName: "LinearSearchWithMetadata",
			Handler:    _SearchWithMetadata_LinearSearchWithMetadata_Handler,
		},
		{
			MethodName: "LinearSearchByIDWithMetadata",
			Handler:    _SearchWithMetadata_LinearSearchByIDWithMetadata_Handler,
		},
		{
			MethodName: "MultiLinearSearchWithMetadata",
			Handler:    _SearchWithMetadata_MultiLinearSearchWithMetadata_Handler,
		},
		{
			MethodName: "MultiLinearSearchByIDWithMetadata",
			Handler:    _SearchWithMetadata_MultiLinearSearchByIDWithMetadata_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamSearchWithMetadata",
			Handler:       _SearchWithMetadata_StreamSearchWithMetadata_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamSearchByIDWithMetadata",
			Handler:       _SearchWithMetadata_StreamSearchByIDWithMetadata_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamLinearSearchWithMetadata",
			Handler:       _SearchWithMetadata_StreamLinearSearchWithMetadata_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamLinearSearchByIDWithMetadata",
			Handler:       _SearchWithMetadata_StreamLinearSearchByIDWithMetadata_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "v1/vald/meta.proto",
}

// InsertWithMetadataClient is the client API for InsertWithMetadata service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InsertWithMetadataClient interface {
	// Overview
	// InsertWithMetadata RPC is the method to add a new single vector and metadata.
	// ---
	// Status Code
	// | 0    | OK                |
	// | 1    | CANCELLED         |
	// | 3    | INVALID_ARGUMENT  |
	// | 4    | DEADLINE_EXCEEDED |
	// | 5    | NOT_FOUND         |
	// | 13   | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	InsertWithMetadata(ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	// Overview
	// StreamInsertWithMetadata RPC is the method to add new multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
	// Each Insert request and response are independent.
	// It's the recommended method to insert a large number of vectors.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	StreamInsertWithMetadata(ctx context.Context, opts ...grpc.CallOption) (InsertWithMetadata_StreamInsertWithMetadataClient, error)
	// Overview
	// MultiInsertWithMetadata RPC is the method to add multiple new vectors and metadata in **1** request.
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
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	MultiInsertWithMetadata(ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type insertWithMetadataClient struct {
	cc grpc.ClientConnInterface
}

func NewInsertWithMetadataClient(cc grpc.ClientConnInterface) InsertWithMetadataClient {
	return &insertWithMetadataClient{cc}
}

func (c *insertWithMetadataClient) InsertWithMetadata(
	ctx context.Context, in *payload.Insert_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/meta.v1.InsertWithMetadata/InsertWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *insertWithMetadataClient) StreamInsertWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (InsertWithMetadata_StreamInsertWithMetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &InsertWithMetadata_ServiceDesc.Streams[0], "/meta.v1.InsertWithMetadata/StreamInsertWithMetadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &insertWithMetadataStreamInsertWithMetadataClient{stream}
	return x, nil
}

type InsertWithMetadata_StreamInsertWithMetadataClient interface {
	Send(*payload.Insert_Request) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type insertWithMetadataStreamInsertWithMetadataClient struct {
	grpc.ClientStream
}

func (x *insertWithMetadataStreamInsertWithMetadataClient) Send(m *payload.Insert_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *insertWithMetadataStreamInsertWithMetadataClient) Recv() (
	*payload.Object_StreamLocation,
	error,
) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *insertWithMetadataClient) MultiInsertWithMetadata(
	ctx context.Context, in *payload.Insert_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/meta.v1.InsertWithMetadata/MultiInsertWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InsertWithMetadataServer is the server API for InsertWithMetadata service.
// All implementations must embed UnimplementedInsertWithMetadataServer
// for forward compatibility
type InsertWithMetadataServer interface {
	// Overview
	// InsertWithMetadata RPC is the method to add a new single vector and metadata.
	// ---
	// Status Code
	// | 0    | OK                |
	// | 1    | CANCELLED         |
	// | 3    | INVALID_ARGUMENT  |
	// | 4    | DEADLINE_EXCEEDED |
	// | 5    | NOT_FOUND         |
	// | 13   | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	InsertWithMetadata(context.Context, *payload.Insert_Request) (*payload.Object_Location, error)
	// Overview
	// StreamInsertWithMetadata RPC is the method to add new multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
	// Each Insert request and response are independent.
	// It's the recommended method to insert a large number of vectors.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	StreamInsertWithMetadata(InsertWithMetadata_StreamInsertWithMetadataServer) error
	// Overview
	// MultiInsertWithMetadata RPC is the method to add multiple new vectors and metadata in **1** request.
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
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	MultiInsertWithMetadata(context.Context, *payload.Insert_MultiRequest) (*payload.Object_Locations, error)
	mustEmbedUnimplementedInsertWithMetadataServer()
}

// UnimplementedInsertWithMetadataServer must be embedded to have forward compatible implementations.
type UnimplementedInsertWithMetadataServer struct{}

func (UnimplementedInsertWithMetadataServer) InsertWithMetadata(
	context.Context, *payload.Insert_Request,
) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InsertWithMetadata not implemented")
}

func (UnimplementedInsertWithMetadataServer) StreamInsertWithMetadata(
	InsertWithMetadata_StreamInsertWithMetadataServer,
) error {
	return status.Errorf(codes.Unimplemented, "method StreamInsertWithMetadata not implemented")
}

func (UnimplementedInsertWithMetadataServer) MultiInsertWithMetadata(
	context.Context, *payload.Insert_MultiRequest,
) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiInsertWithMetadata not implemented")
}
func (UnimplementedInsertWithMetadataServer) mustEmbedUnimplementedInsertWithMetadataServer() {}

// UnsafeInsertWithMetadataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InsertWithMetadataServer will
// result in compilation errors.
type UnsafeInsertWithMetadataServer interface {
	mustEmbedUnimplementedInsertWithMetadataServer()
}

func RegisterInsertWithMetadataServer(s grpc.ServiceRegistrar, srv InsertWithMetadataServer) {
	s.RegisterService(&InsertWithMetadata_ServiceDesc, srv)
}

func _InsertWithMetadata_InsertWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Insert_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InsertWithMetadataServer).InsertWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.InsertWithMetadata/InsertWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(InsertWithMetadataServer).InsertWithMetadata(ctx, req.(*payload.Insert_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _InsertWithMetadata_StreamInsertWithMetadata_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(InsertWithMetadataServer).StreamInsertWithMetadata(&insertWithMetadataStreamInsertWithMetadataServer{stream})
}

type InsertWithMetadata_StreamInsertWithMetadataServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Insert_Request, error)
	grpc.ServerStream
}

type insertWithMetadataStreamInsertWithMetadataServer struct {
	grpc.ServerStream
}

func (x *insertWithMetadataStreamInsertWithMetadataServer) Send(
	m *payload.Object_StreamLocation,
) error {
	return x.ServerStream.SendMsg(m)
}

func (x *insertWithMetadataStreamInsertWithMetadataServer) Recv() (*payload.Insert_Request, error) {
	m := new(payload.Insert_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _InsertWithMetadata_MultiInsertWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Insert_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InsertWithMetadataServer).MultiInsertWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.InsertWithMetadata/MultiInsertWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(InsertWithMetadataServer).MultiInsertWithMetadata(ctx, req.(*payload.Insert_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// InsertWithMetadata_ServiceDesc is the grpc.ServiceDesc for InsertWithMetadata service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InsertWithMetadata_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "meta.v1.InsertWithMetadata",
	HandlerType: (*InsertWithMetadataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "InsertWithMetadata",
			Handler:    _InsertWithMetadata_InsertWithMetadata_Handler,
		},
		{
			MethodName: "MultiInsertWithMetadata",
			Handler:    _InsertWithMetadata_MultiInsertWithMetadata_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamInsertWithMetadata",
			Handler:       _InsertWithMetadata_StreamInsertWithMetadata_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "v1/vald/meta.proto",
}

// ObjectWithMetadataClient is the client API for ObjectWithMetadata service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ObjectWithMetadataClient interface {
	// Overview
	// GetObjectWithMetadata RPC is the method to get the metadata of a vector inserted into the `vald-agent` and metadata.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	GetObjectWithMetadata(ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption) (*payload.Object_Vector, error)
	// Overview
	// StreamGetObjectWithMetadata RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the GetObject request can be communicated in any order between client and server.
	// Each Upsert request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	StreamGetObjectWithMetadata(ctx context.Context, opts ...grpc.CallOption) (ObjectWithMetadata_StreamGetObjectWithMetadataClient, error)
	// Overview
	// A method to get all the vectors with server streaming
	// ---
	// Status Code
	// TODO
	// ---
	// Troubleshooting
	// TODO
	StreamListObjectWithMetadata(ctx context.Context, in *payload.Object_List_Request, opts ...grpc.CallOption) (ObjectWithMetadata_StreamListObjectWithMetadataClient, error)
}

type objectWithMetadataClient struct {
	cc grpc.ClientConnInterface
}

func NewObjectWithMetadataClient(cc grpc.ClientConnInterface) ObjectWithMetadataClient {
	return &objectWithMetadataClient{cc}
}

func (c *objectWithMetadataClient) GetObjectWithMetadata(
	ctx context.Context, in *payload.Object_VectorRequest, opts ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	out := new(payload.Object_Vector)
	err := c.cc.Invoke(ctx, "/meta.v1.ObjectWithMetadata/GetObjectWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *objectWithMetadataClient) StreamGetObjectWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (ObjectWithMetadata_StreamGetObjectWithMetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &ObjectWithMetadata_ServiceDesc.Streams[0], "/meta.v1.ObjectWithMetadata/StreamGetObjectWithMetadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &objectWithMetadataStreamGetObjectWithMetadataClient{stream}
	return x, nil
}

type ObjectWithMetadata_StreamGetObjectWithMetadataClient interface {
	Send(*payload.Object_VectorRequest) error
	Recv() (*payload.Object_StreamVector, error)
	grpc.ClientStream
}

type objectWithMetadataStreamGetObjectWithMetadataClient struct {
	grpc.ClientStream
}

func (x *objectWithMetadataStreamGetObjectWithMetadataClient) Send(
	m *payload.Object_VectorRequest,
) error {
	return x.ClientStream.SendMsg(m)
}

func (x *objectWithMetadataStreamGetObjectWithMetadataClient) Recv() (
	*payload.Object_StreamVector,
	error,
) {
	m := new(payload.Object_StreamVector)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *objectWithMetadataClient) StreamListObjectWithMetadata(
	ctx context.Context, in *payload.Object_List_Request, opts ...grpc.CallOption,
) (ObjectWithMetadata_StreamListObjectWithMetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &ObjectWithMetadata_ServiceDesc.Streams[1], "/meta.v1.ObjectWithMetadata/StreamListObjectWithMetadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &objectWithMetadataStreamListObjectWithMetadataClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ObjectWithMetadata_StreamListObjectWithMetadataClient interface {
	Recv() (*payload.Object_List_Response, error)
	grpc.ClientStream
}

type objectWithMetadataStreamListObjectWithMetadataClient struct {
	grpc.ClientStream
}

func (x *objectWithMetadataStreamListObjectWithMetadataClient) Recv() (
	*payload.Object_List_Response,
	error,
) {
	m := new(payload.Object_List_Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ObjectWithMetadataServer is the server API for ObjectWithMetadata service.
// All implementations must embed UnimplementedObjectWithMetadataServer
// for forward compatibility
type ObjectWithMetadataServer interface {
	// Overview
	// GetObjectWithMetadata RPC is the method to get the metadata of a vector inserted into the `vald-agent` and metadata.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	GetObjectWithMetadata(context.Context, *payload.Object_VectorRequest) (*payload.Object_Vector, error)
	// Overview
	// StreamGetObjectWithMetadata RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the GetObject request can be communicated in any order between client and server.
	// Each Upsert request and response are independent.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	StreamGetObjectWithMetadata(ObjectWithMetadata_StreamGetObjectWithMetadataServer) error
	// Overview
	// A method to get all the vectors with server streaming
	// ---
	// Status Code
	// TODO
	// ---
	// Troubleshooting
	// TODO
	StreamListObjectWithMetadata(*payload.Object_List_Request, ObjectWithMetadata_StreamListObjectWithMetadataServer) error
	mustEmbedUnimplementedObjectWithMetadataServer()
}

// UnimplementedObjectWithMetadataServer must be embedded to have forward compatible implementations.
type UnimplementedObjectWithMetadataServer struct{}

func (UnimplementedObjectWithMetadataServer) GetObjectWithMetadata(
	context.Context, *payload.Object_VectorRequest,
) (*payload.Object_Vector, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetObjectWithMetadata not implemented")
}

func (UnimplementedObjectWithMetadataServer) StreamGetObjectWithMetadata(
	ObjectWithMetadata_StreamGetObjectWithMetadataServer,
) error {
	return status.Errorf(codes.Unimplemented, "method StreamGetObjectWithMetadata not implemented")
}

func (UnimplementedObjectWithMetadataServer) StreamListObjectWithMetadata(
	*payload.Object_List_Request, ObjectWithMetadata_StreamListObjectWithMetadataServer,
) error {
	return status.Errorf(codes.Unimplemented, "method StreamListObjectWithMetadata not implemented")
}
func (UnimplementedObjectWithMetadataServer) mustEmbedUnimplementedObjectWithMetadataServer() {}

// UnsafeObjectWithMetadataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ObjectWithMetadataServer will
// result in compilation errors.
type UnsafeObjectWithMetadataServer interface {
	mustEmbedUnimplementedObjectWithMetadataServer()
}

func RegisterObjectWithMetadataServer(s grpc.ServiceRegistrar, srv ObjectWithMetadataServer) {
	s.RegisterService(&ObjectWithMetadata_ServiceDesc, srv)
}

func _ObjectWithMetadata_GetObjectWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Object_VectorRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ObjectWithMetadataServer).GetObjectWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.ObjectWithMetadata/GetObjectWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(ObjectWithMetadataServer).GetObjectWithMetadata(ctx, req.(*payload.Object_VectorRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ObjectWithMetadata_StreamGetObjectWithMetadata_Handler(
	srv any, stream grpc.ServerStream,
) error {
	return srv.(ObjectWithMetadataServer).StreamGetObjectWithMetadata(&objectWithMetadataStreamGetObjectWithMetadataServer{stream})
}

type ObjectWithMetadata_StreamGetObjectWithMetadataServer interface {
	Send(*payload.Object_StreamVector) error
	Recv() (*payload.Object_VectorRequest, error)
	grpc.ServerStream
}

type objectWithMetadataStreamGetObjectWithMetadataServer struct {
	grpc.ServerStream
}

func (x *objectWithMetadataStreamGetObjectWithMetadataServer) Send(
	m *payload.Object_StreamVector,
) error {
	return x.ServerStream.SendMsg(m)
}

func (x *objectWithMetadataStreamGetObjectWithMetadataServer) Recv() (
	*payload.Object_VectorRequest,
	error,
) {
	m := new(payload.Object_VectorRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _ObjectWithMetadata_StreamListObjectWithMetadata_Handler(
	srv any, stream grpc.ServerStream,
) error {
	m := new(payload.Object_List_Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ObjectWithMetadataServer).StreamListObjectWithMetadata(m, &objectWithMetadataStreamListObjectWithMetadataServer{stream})
}

type ObjectWithMetadata_StreamListObjectWithMetadataServer interface {
	Send(*payload.Object_List_Response) error
	grpc.ServerStream
}

type objectWithMetadataStreamListObjectWithMetadataServer struct {
	grpc.ServerStream
}

func (x *objectWithMetadataStreamListObjectWithMetadataServer) Send(
	m *payload.Object_List_Response,
) error {
	return x.ServerStream.SendMsg(m)
}

// ObjectWithMetadata_ServiceDesc is the grpc.ServiceDesc for ObjectWithMetadata service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ObjectWithMetadata_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "meta.v1.ObjectWithMetadata",
	HandlerType: (*ObjectWithMetadataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetObjectWithMetadata",
			Handler:    _ObjectWithMetadata_GetObjectWithMetadata_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamGetObjectWithMetadata",
			Handler:       _ObjectWithMetadata_StreamGetObjectWithMetadata_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamListObjectWithMetadata",
			Handler:       _ObjectWithMetadata_StreamListObjectWithMetadata_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "v1/vald/meta.proto",
}

// RemoveWithMetadataClient is the client API for RemoveWithMetadata service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RemoveWithMetadataClient interface {
	// Overview
	// RemoveWithMetadata RPC is the method to remove a single vector and metadata.
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
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	RemoveWithMetadata(ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	// Overview
	// RemoveByTimestampWithMetadata RPC is the method to remove vectors and metadata based on timestamp.
	//
	// <div class="notice">
	// In the TimestampRequest message, the 'timestamps' field is repeated, allowing the inclusion of multiple Timestamp.<br>
	// When multiple Timestamps are provided, it results in an `AND` condition, enabling the realization of deletions with specified ranges.<br>
	// This design allows for versatile deletion operations, facilitating tasks such as removing data within a specific time range.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                                                       |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.                              |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed.                             |
	// | NOT_FOUND         | No vectors in the system match the specified timestamp conditions.                              | Check whether vectors matching the specified timestamp conditions exist in the system, and fix conditions if needed. |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.
	RemoveByTimestampWithMetadata(ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
	// Overview
	// A method to remove multiple with metadata indexed vectors and metadata by bidirectional streaming.
	//
	// StreamRemoveWithMetadata RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the remove request can be communicated in any order between client and server.
	// Each Remove request and response are independent.
	// It's the recommended method to remove a large number of vectors.
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
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	StreamRemoveWithMetadata(ctx context.Context, opts ...grpc.CallOption) (RemoveWithMetadata_StreamRemoveWithMetadataClient, error)
	// Overview
	// MultiRemoveWithMetadata is the method to remove multiple vectors and metadata in **1** request.
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
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	MultiRemoveWithMetadata(ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type removeWithMetadataClient struct {
	cc grpc.ClientConnInterface
}

func NewRemoveWithMetadataClient(cc grpc.ClientConnInterface) RemoveWithMetadataClient {
	return &removeWithMetadataClient{cc}
}

func (c *removeWithMetadataClient) RemoveWithMetadata(
	ctx context.Context, in *payload.Remove_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/meta.v1.RemoveWithMetadata/RemoveWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *removeWithMetadataClient) RemoveByTimestampWithMetadata(
	ctx context.Context, in *payload.Remove_TimestampRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/meta.v1.RemoveWithMetadata/RemoveByTimestampWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *removeWithMetadataClient) StreamRemoveWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (RemoveWithMetadata_StreamRemoveWithMetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &RemoveWithMetadata_ServiceDesc.Streams[0], "/meta.v1.RemoveWithMetadata/StreamRemoveWithMetadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &removeWithMetadataStreamRemoveWithMetadataClient{stream}
	return x, nil
}

type RemoveWithMetadata_StreamRemoveWithMetadataClient interface {
	Send(*payload.Remove_Request) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type removeWithMetadataStreamRemoveWithMetadataClient struct {
	grpc.ClientStream
}

func (x *removeWithMetadataStreamRemoveWithMetadataClient) Send(m *payload.Remove_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *removeWithMetadataStreamRemoveWithMetadataClient) Recv() (
	*payload.Object_StreamLocation,
	error,
) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *removeWithMetadataClient) MultiRemoveWithMetadata(
	ctx context.Context, in *payload.Remove_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/meta.v1.RemoveWithMetadata/MultiRemoveWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveWithMetadataServer is the server API for RemoveWithMetadata service.
// All implementations must embed UnimplementedRemoveWithMetadataServer
// for forward compatibility
type RemoveWithMetadataServer interface {
	// Overview
	// RemoveWithMetadata RPC is the method to remove a single vector and metadata.
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
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	RemoveWithMetadata(context.Context, *payload.Remove_Request) (*payload.Object_Location, error)
	// Overview
	// RemoveByTimestampWithMetadata RPC is the method to remove vectors and metadata based on timestamp.
	//
	// <div class="notice">
	// In the TimestampRequest message, the 'timestamps' field is repeated, allowing the inclusion of multiple Timestamp.<br>
	// When multiple Timestamps are provided, it results in an `AND` condition, enabling the realization of deletions with specified ranges.<br>
	// This design allows for versatile deletion operations, facilitating tasks such as removing data within a specific time range.
	// </div>
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                   | how to resolve                                                                                                       |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.                              |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed.                             |
	// | NOT_FOUND         | No vectors in the system match the specified timestamp conditions.                              | Check whether vectors matching the specified timestamp conditions exist in the system, and fix conditions if needed. |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.
	RemoveByTimestampWithMetadata(context.Context, *payload.Remove_TimestampRequest) (*payload.Object_Locations, error)
	// Overview
	// A method to remove multiple with metadata indexed vectors and metadata by bidirectional streaming.
	//
	// StreamRemoveWithMetadata RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the remove request can be communicated in any order between client and server.
	// Each Remove request and response are independent.
	// It's the recommended method to remove a large number of vectors.
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
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	StreamRemoveWithMetadata(RemoveWithMetadata_StreamRemoveWithMetadataServer) error
	// Overview
	// MultiRemoveWithMetadata is the method to remove multiple vectors and metadata in **1** request.
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
	// | name              | common reason                                                                                   | how to resolve                                                                           |
	// | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
	MultiRemoveWithMetadata(context.Context, *payload.Remove_MultiRequest) (*payload.Object_Locations, error)
	mustEmbedUnimplementedRemoveWithMetadataServer()
}

// UnimplementedRemoveWithMetadataServer must be embedded to have forward compatible implementations.
type UnimplementedRemoveWithMetadataServer struct{}

func (UnimplementedRemoveWithMetadataServer) RemoveWithMetadata(
	context.Context, *payload.Remove_Request,
) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveWithMetadata not implemented")
}

func (UnimplementedRemoveWithMetadataServer) RemoveByTimestampWithMetadata(
	context.Context, *payload.Remove_TimestampRequest,
) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveByTimestampWithMetadata not implemented")
}

func (UnimplementedRemoveWithMetadataServer) StreamRemoveWithMetadata(
	RemoveWithMetadata_StreamRemoveWithMetadataServer,
) error {
	return status.Errorf(codes.Unimplemented, "method StreamRemoveWithMetadata not implemented")
}

func (UnimplementedRemoveWithMetadataServer) MultiRemoveWithMetadata(
	context.Context, *payload.Remove_MultiRequest,
) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiRemoveWithMetadata not implemented")
}
func (UnimplementedRemoveWithMetadataServer) mustEmbedUnimplementedRemoveWithMetadataServer() {}

// UnsafeRemoveWithMetadataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RemoveWithMetadataServer will
// result in compilation errors.
type UnsafeRemoveWithMetadataServer interface {
	mustEmbedUnimplementedRemoveWithMetadataServer()
}

func RegisterRemoveWithMetadataServer(s grpc.ServiceRegistrar, srv RemoveWithMetadataServer) {
	s.RegisterService(&RemoveWithMetadata_ServiceDesc, srv)
}

func _RemoveWithMetadata_RemoveWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Remove_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoveWithMetadataServer).RemoveWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.RemoveWithMetadata/RemoveWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(RemoveWithMetadataServer).RemoveWithMetadata(ctx, req.(*payload.Remove_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _RemoveWithMetadata_RemoveByTimestampWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Remove_TimestampRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoveWithMetadataServer).RemoveByTimestampWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.RemoveWithMetadata/RemoveByTimestampWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(RemoveWithMetadataServer).RemoveByTimestampWithMetadata(ctx, req.(*payload.Remove_TimestampRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RemoveWithMetadata_StreamRemoveWithMetadata_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(RemoveWithMetadataServer).StreamRemoveWithMetadata(&removeWithMetadataStreamRemoveWithMetadataServer{stream})
}

type RemoveWithMetadata_StreamRemoveWithMetadataServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Remove_Request, error)
	grpc.ServerStream
}

type removeWithMetadataStreamRemoveWithMetadataServer struct {
	grpc.ServerStream
}

func (x *removeWithMetadataStreamRemoveWithMetadataServer) Send(
	m *payload.Object_StreamLocation,
) error {
	return x.ServerStream.SendMsg(m)
}

func (x *removeWithMetadataStreamRemoveWithMetadataServer) Recv() (*payload.Remove_Request, error) {
	m := new(payload.Remove_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RemoveWithMetadata_MultiRemoveWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Remove_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemoveWithMetadataServer).MultiRemoveWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.RemoveWithMetadata/MultiRemoveWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(RemoveWithMetadataServer).MultiRemoveWithMetadata(ctx, req.(*payload.Remove_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RemoveWithMetadata_ServiceDesc is the grpc.ServiceDesc for RemoveWithMetadata service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RemoveWithMetadata_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "meta.v1.RemoveWithMetadata",
	HandlerType: (*RemoveWithMetadataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RemoveWithMetadata",
			Handler:    _RemoveWithMetadata_RemoveWithMetadata_Handler,
		},
		{
			MethodName: "RemoveByTimestampWithMetadata",
			Handler:    _RemoveWithMetadata_RemoveByTimestampWithMetadata_Handler,
		},
		{
			MethodName: "MultiRemoveWithMetadata",
			Handler:    _RemoveWithMetadata_MultiRemoveWithMetadata_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamRemoveWithMetadata",
			Handler:       _RemoveWithMetadata_StreamRemoveWithMetadata_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "v1/vald/meta.proto",
}

// UpdateWithMetadataClient is the client API for UpdateWithMetadata service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UpdateWithMetadataClient interface {
	// Overview
	// UpdateWithMetadata RPC is the method to update a single vector.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
	// | ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	UpdateWithMetadata(ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	// Overview
	// StreamUpdateWithMetadata RPC is the method to update multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the update request can be communicated in any order between client and server.
	// Each Update request and response are independent.
	// It's the recommended method to update the large amount of vectors.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
	// | ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	StreamUpdateWithMetadata(ctx context.Context, opts ...grpc.CallOption) (UpdateWithMetadata_StreamUpdateWithMetadataClient, error)
	// Overview
	// MultiUpdateWithMetadata is the method to update multiple vectors and metadata in **1** request.
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
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
	// | ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	MultiUpdateWithMetadata(ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
	// Overview
	// A method to update timestamp an indexed vector and metadata.
	// ---
	// Status Code
	// TODO
	// ---
	// Troubleshooting
	// TODO
	UpdateTimestampWithMetadata(ctx context.Context, in *payload.Update_TimestampRequest, opts ...grpc.CallOption) (*payload.Object_Location, error)
}

type updateWithMetadataClient struct {
	cc grpc.ClientConnInterface
}

func NewUpdateWithMetadataClient(cc grpc.ClientConnInterface) UpdateWithMetadataClient {
	return &updateWithMetadataClient{cc}
}

func (c *updateWithMetadataClient) UpdateWithMetadata(
	ctx context.Context, in *payload.Update_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/meta.v1.UpdateWithMetadata/UpdateWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *updateWithMetadataClient) StreamUpdateWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (UpdateWithMetadata_StreamUpdateWithMetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &UpdateWithMetadata_ServiceDesc.Streams[0], "/meta.v1.UpdateWithMetadata/StreamUpdateWithMetadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &updateWithMetadataStreamUpdateWithMetadataClient{stream}
	return x, nil
}

type UpdateWithMetadata_StreamUpdateWithMetadataClient interface {
	Send(*payload.Update_Request) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type updateWithMetadataStreamUpdateWithMetadataClient struct {
	grpc.ClientStream
}

func (x *updateWithMetadataStreamUpdateWithMetadataClient) Send(m *payload.Update_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *updateWithMetadataStreamUpdateWithMetadataClient) Recv() (
	*payload.Object_StreamLocation,
	error,
) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *updateWithMetadataClient) MultiUpdateWithMetadata(
	ctx context.Context, in *payload.Update_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/meta.v1.UpdateWithMetadata/MultiUpdateWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *updateWithMetadataClient) UpdateTimestampWithMetadata(
	ctx context.Context, in *payload.Update_TimestampRequest, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/meta.v1.UpdateWithMetadata/UpdateTimestampWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpdateWithMetadataServer is the server API for UpdateWithMetadata service.
// All implementations must embed UnimplementedUpdateWithMetadataServer
// for forward compatibility
type UpdateWithMetadataServer interface {
	// Overview
	// UpdateWithMetadata RPC is the method to update a single vector.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
	// | ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	UpdateWithMetadata(context.Context, *payload.Update_Request) (*payload.Object_Location, error)
	// Overview
	// StreamUpdateWithMetadata RPC is the method to update multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the update request can be communicated in any order between client and server.
	// Each Update request and response are independent.
	// It's the recommended method to update the large amount of vectors.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
	// | ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	StreamUpdateWithMetadata(UpdateWithMetadata_StreamUpdateWithMetadataServer) error
	// Overview
	// MultiUpdateWithMetadata is the method to update multiple vectors and metadata in **1** request.
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
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
	// | ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	MultiUpdateWithMetadata(context.Context, *payload.Update_MultiRequest) (*payload.Object_Locations, error)
	// Overview
	// A method to update timestamp an indexed vector and metadata.
	// ---
	// Status Code
	// TODO
	// ---
	// Troubleshooting
	// TODO
	UpdateTimestampWithMetadata(context.Context, *payload.Update_TimestampRequest) (*payload.Object_Location, error)
	mustEmbedUnimplementedUpdateWithMetadataServer()
}

// UnimplementedUpdateWithMetadataServer must be embedded to have forward compatible implementations.
type UnimplementedUpdateWithMetadataServer struct{}

func (UnimplementedUpdateWithMetadataServer) UpdateWithMetadata(
	context.Context, *payload.Update_Request,
) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateWithMetadata not implemented")
}

func (UnimplementedUpdateWithMetadataServer) StreamUpdateWithMetadata(
	UpdateWithMetadata_StreamUpdateWithMetadataServer,
) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpdateWithMetadata not implemented")
}

func (UnimplementedUpdateWithMetadataServer) MultiUpdateWithMetadata(
	context.Context, *payload.Update_MultiRequest,
) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpdateWithMetadata not implemented")
}

func (UnimplementedUpdateWithMetadataServer) UpdateTimestampWithMetadata(
	context.Context, *payload.Update_TimestampRequest,
) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTimestampWithMetadata not implemented")
}
func (UnimplementedUpdateWithMetadataServer) mustEmbedUnimplementedUpdateWithMetadataServer() {}

// UnsafeUpdateWithMetadataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UpdateWithMetadataServer will
// result in compilation errors.
type UnsafeUpdateWithMetadataServer interface {
	mustEmbedUnimplementedUpdateWithMetadataServer()
}

func RegisterUpdateWithMetadataServer(s grpc.ServiceRegistrar, srv UpdateWithMetadataServer) {
	s.RegisterService(&UpdateWithMetadata_ServiceDesc, srv)
}

func _UpdateWithMetadata_UpdateWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Update_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpdateWithMetadataServer).UpdateWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.UpdateWithMetadata/UpdateWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UpdateWithMetadataServer).UpdateWithMetadata(ctx, req.(*payload.Update_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _UpdateWithMetadata_StreamUpdateWithMetadata_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(UpdateWithMetadataServer).StreamUpdateWithMetadata(&updateWithMetadataStreamUpdateWithMetadataServer{stream})
}

type UpdateWithMetadata_StreamUpdateWithMetadataServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Update_Request, error)
	grpc.ServerStream
}

type updateWithMetadataStreamUpdateWithMetadataServer struct {
	grpc.ServerStream
}

func (x *updateWithMetadataStreamUpdateWithMetadataServer) Send(
	m *payload.Object_StreamLocation,
) error {
	return x.ServerStream.SendMsg(m)
}

func (x *updateWithMetadataStreamUpdateWithMetadataServer) Recv() (*payload.Update_Request, error) {
	m := new(payload.Update_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _UpdateWithMetadata_MultiUpdateWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Update_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpdateWithMetadataServer).MultiUpdateWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.UpdateWithMetadata/MultiUpdateWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UpdateWithMetadataServer).MultiUpdateWithMetadata(ctx, req.(*payload.Update_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UpdateWithMetadata_UpdateTimestampWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Update_TimestampRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpdateWithMetadataServer).UpdateTimestampWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.UpdateWithMetadata/UpdateTimestampWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UpdateWithMetadataServer).UpdateTimestampWithMetadata(ctx, req.(*payload.Update_TimestampRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UpdateWithMetadata_ServiceDesc is the grpc.ServiceDesc for UpdateWithMetadata service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UpdateWithMetadata_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "meta.v1.UpdateWithMetadata",
	HandlerType: (*UpdateWithMetadataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateWithMetadata",
			Handler:    _UpdateWithMetadata_UpdateWithMetadata_Handler,
		},
		{
			MethodName: "MultiUpdateWithMetadata",
			Handler:    _UpdateWithMetadata_MultiUpdateWithMetadata_Handler,
		},
		{
			MethodName: "UpdateTimestampWithMetadata",
			Handler:    _UpdateWithMetadata_UpdateTimestampWithMetadata_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamUpdateWithMetadata",
			Handler:       _UpdateWithMetadata_StreamUpdateWithMetadata_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "v1/vald/meta.proto",
}

// UpsertWithMetadataClient is the client API for UpsertWithMetadata service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UpsertWithMetadataClient interface {
	// Overview
	// UpsertWithMetadata RPC is the method to update the inserted vector and metadata to a new single vector and metadata or add a new single vector and metadata if not inserted before.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	UpsertWithMetadata(ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption) (*payload.Object_Location, error)
	// Overview
	// StreamUpsertWithMetadata RPC is the method to update multiple existing vectors and metadata or add new multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the upsert request can be communicated in any order between the client and server.
	// Each Upsert request and response are independent.
	// It’s the recommended method to upsert a large number of vectors.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	StreamUpsertWithMetadata(ctx context.Context, opts ...grpc.CallOption) (UpsertWithMetadata_StreamUpsertWithMetadataClient, error)
	// Overview
	// MultiUpsertWithMetadata is the method to update existing multiple vectors and metadata and add new multiple vectors and metadata in **1** request.
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
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	MultiUpsertWithMetadata(ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption) (*payload.Object_Locations, error)
}

type upsertWithMetadataClient struct {
	cc grpc.ClientConnInterface
}

func NewUpsertWithMetadataClient(cc grpc.ClientConnInterface) UpsertWithMetadataClient {
	return &upsertWithMetadataClient{cc}
}

func (c *upsertWithMetadataClient) UpsertWithMetadata(
	ctx context.Context, in *payload.Upsert_Request, opts ...grpc.CallOption,
) (*payload.Object_Location, error) {
	out := new(payload.Object_Location)
	err := c.cc.Invoke(ctx, "/meta.v1.UpsertWithMetadata/UpsertWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *upsertWithMetadataClient) StreamUpsertWithMetadata(
	ctx context.Context, opts ...grpc.CallOption,
) (UpsertWithMetadata_StreamUpsertWithMetadataClient, error) {
	stream, err := c.cc.NewStream(ctx, &UpsertWithMetadata_ServiceDesc.Streams[0], "/meta.v1.UpsertWithMetadata/StreamUpsertWithMetadata", opts...)
	if err != nil {
		return nil, err
	}
	x := &upsertWithMetadataStreamUpsertWithMetadataClient{stream}
	return x, nil
}

type UpsertWithMetadata_StreamUpsertWithMetadataClient interface {
	Send(*payload.Upsert_Request) error
	Recv() (*payload.Object_StreamLocation, error)
	grpc.ClientStream
}

type upsertWithMetadataStreamUpsertWithMetadataClient struct {
	grpc.ClientStream
}

func (x *upsertWithMetadataStreamUpsertWithMetadataClient) Send(m *payload.Upsert_Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *upsertWithMetadataStreamUpsertWithMetadataClient) Recv() (
	*payload.Object_StreamLocation,
	error,
) {
	m := new(payload.Object_StreamLocation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *upsertWithMetadataClient) MultiUpsertWithMetadata(
	ctx context.Context, in *payload.Upsert_MultiRequest, opts ...grpc.CallOption,
) (*payload.Object_Locations, error) {
	out := new(payload.Object_Locations)
	err := c.cc.Invoke(ctx, "/meta.v1.UpsertWithMetadata/MultiUpsertWithMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UpsertWithMetadataServer is the server API for UpsertWithMetadata service.
// All implementations must embed UnimplementedUpsertWithMetadataServer
// for forward compatibility
type UpsertWithMetadataServer interface {
	// Overview
	// UpsertWithMetadata RPC is the method to update the inserted vector and metadata to a new single vector and metadata or add a new single vector and metadata if not inserted before.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	UpsertWithMetadata(context.Context, *payload.Upsert_Request) (*payload.Object_Location, error)
	// Overview
	// StreamUpsertWithMetadata RPC is the method to update multiple existing vectors and metadata or add new multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
	// Using the bidirectional streaming RPC, the upsert request can be communicated in any order between the client and server.
	// Each Upsert request and response are independent.
	// It’s the recommended method to upsert a large number of vectors.
	// ---
	// Status Code
	// |  0   | OK                |
	// |  1   | CANCELLED         |
	// |  3   | INVALID_ARGUMENT  |
	// |  4   | DEADLINE_EXCEEDED |
	// |  5   | NOT_FOUND         |
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	StreamUpsertWithMetadata(UpsertWithMetadata_StreamUpsertWithMetadataServer) error
	// Overview
	// MultiUpsertWithMetadata is the method to update existing multiple vectors and metadata and add new multiple vectors and metadata in **1** request.
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
	// |  6   | ALREADY_EXISTS    |
	// |  10  | ABORTED           |
	// |  13  | INTERNAL          |
	// ---
	// Troubleshooting
	// The request process may not be completed when the response code is NOT `0 (OK)`.
	//
	// Here are some common reasons and how to resolve each error.
	//
	// | name              | common reason                                                                                                                                       | how to resolve                                                                           |
	// | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
	// | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
	// | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
	// | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
	// | ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
	// | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
	MultiUpsertWithMetadata(context.Context, *payload.Upsert_MultiRequest) (*payload.Object_Locations, error)
	mustEmbedUnimplementedUpsertWithMetadataServer()
}

// UnimplementedUpsertWithMetadataServer must be embedded to have forward compatible implementations.
type UnimplementedUpsertWithMetadataServer struct{}

func (UnimplementedUpsertWithMetadataServer) UpsertWithMetadata(
	context.Context, *payload.Upsert_Request,
) (*payload.Object_Location, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpsertWithMetadata not implemented")
}

func (UnimplementedUpsertWithMetadataServer) StreamUpsertWithMetadata(
	UpsertWithMetadata_StreamUpsertWithMetadataServer,
) error {
	return status.Errorf(codes.Unimplemented, "method StreamUpsertWithMetadata not implemented")
}

func (UnimplementedUpsertWithMetadataServer) MultiUpsertWithMetadata(
	context.Context, *payload.Upsert_MultiRequest,
) (*payload.Object_Locations, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiUpsertWithMetadata not implemented")
}
func (UnimplementedUpsertWithMetadataServer) mustEmbedUnimplementedUpsertWithMetadataServer() {}

// UnsafeUpsertWithMetadataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UpsertWithMetadataServer will
// result in compilation errors.
type UnsafeUpsertWithMetadataServer interface {
	mustEmbedUnimplementedUpsertWithMetadataServer()
}

func RegisterUpsertWithMetadataServer(s grpc.ServiceRegistrar, srv UpsertWithMetadataServer) {
	s.RegisterService(&UpsertWithMetadata_ServiceDesc, srv)
}

func _UpsertWithMetadata_UpsertWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Upsert_Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpsertWithMetadataServer).UpsertWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.UpsertWithMetadata/UpsertWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UpsertWithMetadataServer).UpsertWithMetadata(ctx, req.(*payload.Upsert_Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _UpsertWithMetadata_StreamUpsertWithMetadata_Handler(srv any, stream grpc.ServerStream) error {
	return srv.(UpsertWithMetadataServer).StreamUpsertWithMetadata(&upsertWithMetadataStreamUpsertWithMetadataServer{stream})
}

type UpsertWithMetadata_StreamUpsertWithMetadataServer interface {
	Send(*payload.Object_StreamLocation) error
	Recv() (*payload.Upsert_Request, error)
	grpc.ServerStream
}

type upsertWithMetadataStreamUpsertWithMetadataServer struct {
	grpc.ServerStream
}

func (x *upsertWithMetadataStreamUpsertWithMetadataServer) Send(
	m *payload.Object_StreamLocation,
) error {
	return x.ServerStream.SendMsg(m)
}

func (x *upsertWithMetadataStreamUpsertWithMetadataServer) Recv() (*payload.Upsert_Request, error) {
	m := new(payload.Upsert_Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _UpsertWithMetadata_MultiUpsertWithMetadata_Handler(
	srv any, ctx context.Context, dec func(any) error, interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(payload.Upsert_MultiRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UpsertWithMetadataServer).MultiUpsertWithMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/meta.v1.UpsertWithMetadata/MultiUpsertWithMetadata",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(UpsertWithMetadataServer).MultiUpsertWithMetadata(ctx, req.(*payload.Upsert_MultiRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UpsertWithMetadata_ServiceDesc is the grpc.ServiceDesc for UpsertWithMetadata service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UpsertWithMetadata_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "meta.v1.UpsertWithMetadata",
	HandlerType: (*UpsertWithMetadataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpsertWithMetadata",
			Handler:    _UpsertWithMetadata_UpsertWithMetadata_Handler,
		},
		{
			MethodName: "MultiUpsertWithMetadata",
			Handler:    _UpsertWithMetadata_MultiUpsertWithMetadata_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamUpsertWithMetadata",
			Handler:       _UpsertWithMetadata_StreamUpsertWithMetadata_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "v1/vald/meta.proto",
}
