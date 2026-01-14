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
pub mod filter_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    /** Overview
 Filter Server is responsible for providing insert, update, upsert and search interface for `Vald Filter Gateway`.

 Vald Filter Gateway forward user request to user-defined ingress/egress filter components allowing user to run custom logic.
*/
    #[derive(Debug, Clone)]
    pub struct FilterClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl FilterClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> FilterClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::Body>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> FilterClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::Body>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::Body>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::Body>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            FilterClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        /** Overview
 SearchObject RPC is the method to search object(s) similar to request object.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn search_object(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::search::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/SearchObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "SearchObject"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamSearchObject RPC is the method to search vectors with multi queries(objects) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 By using the bidirectional streaming RPC, the search request can be communicated in any order between client and server.
 Each Search request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn multi_search_object(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::search::MultiObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/MultiSearchObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "MultiSearchObject"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 MultiSearchObject RPC is the method to search objects with multiple objects in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn stream_search_object(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::search::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::search::StreamResponse,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/StreamSearchObject",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "StreamSearchObject"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 InsertObject RPC is the method to insert object through Vald Filter Gateway.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn insert_object(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::insert::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/InsertObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "InsertObject"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamInsertObject RPC is the method to add new multiple object using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).

 By using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
 Each Insert request and response are independent.
 It's the recommended method to insert a large number of objects.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn stream_insert_object(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::insert::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::object::StreamLocation,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/StreamInsertObject",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "StreamInsertObject"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiInsertObject RPC is the method to add multiple new objects in **1** request.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn multi_insert_object(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::insert::MultiObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/MultiInsertObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "MultiInsertObject"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 UpdateObject RPC is the method to update a single vector.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn update_object(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::update::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/UpdateObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "UpdateObject"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamUpdateObject RPC is the method to update multiple objects using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 By using the bidirectional streaming RPC, the update request can be communicated in any order between client and server.
 Each Update request and response are independent.
 It's the recommended method to update the large amount of objects.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn stream_update_object(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::update::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::object::StreamLocation,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/StreamUpdateObject",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "StreamUpdateObject"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiUpdateObject is the method to update multiple objects in **1** request.

 <div class="notice">
 gRPC has the message size limitation.<br>
 Please be careful that the size of the request exceed the limit.
 </div>
 ---
 Status Code

 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn multi_update_object(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::update::MultiObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/MultiUpdateObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "MultiUpdateObject"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 UpsertObject RPC is the method to update a single object and add a new single object.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn upsert_object(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::upsert::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/UpsertObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "UpsertObject"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 UpsertObject RPC is the method to update a single object and add a new single object.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn stream_upsert_object(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::upsert::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::object::StreamLocation,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/StreamUpsertObject",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "StreamUpsertObject"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiUpsertObject is the method to update existing multiple objects and add new multiple objects in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        pub async fn multi_upsert_object(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::upsert::MultiObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/MultiUpsertObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "MultiUpsertObject"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod filter_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with FilterServer.
    #[async_trait]
    pub trait Filter: std::marker::Send + std::marker::Sync + 'static {
        /** Overview
 SearchObject RPC is the method to search object(s) similar to request object.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn search_object(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::search::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /** Overview
 StreamSearchObject RPC is the method to search vectors with multi queries(objects) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 By using the bidirectional streaming RPC, the search request can be communicated in any order between client and server.
 Each Search request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn multi_search_object(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::search::MultiObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamSearchObject method.
        type StreamSearchObjectStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 MultiSearchObject RPC is the method to search objects with multiple objects in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn stream_search_object(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::search::ObjectRequest>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamSearchObjectStream>,
            tonic::Status,
        >;
        /** Overview
 InsertObject RPC is the method to insert object through Vald Filter Gateway.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn insert_object(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::insert::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamInsertObject method.
        type StreamInsertObjectStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamInsertObject RPC is the method to add new multiple object using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).

 By using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
 Each Insert request and response are independent.
 It's the recommended method to insert a large number of objects.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn stream_insert_object(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::insert::ObjectRequest>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamInsertObjectStream>,
            tonic::Status,
        >;
        /** Overview
 MultiInsertObject RPC is the method to add multiple new objects in **1** request.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn multi_insert_object(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::insert::MultiObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
        /** Overview
 UpdateObject RPC is the method to update a single vector.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn update_object(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::update::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamUpdateObject method.
        type StreamUpdateObjectStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamUpdateObject RPC is the method to update multiple objects using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 By using the bidirectional streaming RPC, the update request can be communicated in any order between client and server.
 Each Update request and response are independent.
 It's the recommended method to update the large amount of objects.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn stream_update_object(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::update::ObjectRequest>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamUpdateObjectStream>,
            tonic::Status,
        >;
        /** Overview
 MultiUpdateObject is the method to update multiple objects in **1** request.

 <div class="notice">
 gRPC has the message size limitation.<br>
 Please be careful that the size of the request exceed the limit.
 </div>
 ---
 Status Code

 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn multi_update_object(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::update::MultiObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
        /** Overview
 UpsertObject RPC is the method to update a single object and add a new single object.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn upsert_object(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::upsert::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamUpsertObject method.
        type StreamUpsertObjectStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 UpsertObject RPC is the method to update a single object and add a new single object.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn stream_upsert_object(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::upsert::ObjectRequest>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamUpsertObjectStream>,
            tonic::Status,
        >;
        /** Overview
 MultiUpsertObject is the method to update existing multiple objects and add new multiple objects in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  13  | INTERNAL          |
*/
        async fn multi_upsert_object(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::upsert::MultiObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
    }
    /** Overview
 Filter Server is responsible for providing insert, update, upsert and search interface for `Vald Filter Gateway`.

 Vald Filter Gateway forward user request to user-defined ingress/egress filter components allowing user to run custom logic.
*/
    #[derive(Debug)]
    pub struct FilterServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> FilterServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for FilterServer<T>
    where
        T: Filter,
        B: Body + std::marker::Send + 'static,
        B::Error: Into<StdError> + std::marker::Send + 'static,
    {
        type Response = http::Response<tonic::body::Body>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/vald.v1.Filter/SearchObject" => {
                    #[allow(non_camel_case_types)]
                    struct SearchObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::ObjectRequest,
                    > for SearchObjectSvc<T> {
                        type Response = super::super::super::payload::v1::search::Response;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::search::ObjectRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::search_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = SearchObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Filter/MultiSearchObject" => {
                    #[allow(non_camel_case_types)]
                    struct MultiSearchObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::MultiObjectRequest,
                    > for MultiSearchObjectSvc<T> {
                        type Response = super::super::super::payload::v1::search::Responses;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::search::MultiObjectRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::multi_search_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiSearchObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Filter/StreamSearchObject" => {
                    #[allow(non_camel_case_types)]
                    struct StreamSearchObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::search::ObjectRequest,
                    > for StreamSearchObjectSvc<T> {
                        type Response = super::super::super::payload::v1::search::StreamResponse;
                        type ResponseStream = T::StreamSearchObjectStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::search::ObjectRequest,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::stream_search_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamSearchObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Filter/InsertObject" => {
                    #[allow(non_camel_case_types)]
                    struct InsertObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::insert::ObjectRequest,
                    > for InsertObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::Location;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::insert::ObjectRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::insert_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = InsertObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Filter/StreamInsertObject" => {
                    #[allow(non_camel_case_types)]
                    struct StreamInsertObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::insert::ObjectRequest,
                    > for StreamInsertObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamLocation;
                        type ResponseStream = T::StreamInsertObjectStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::insert::ObjectRequest,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::stream_insert_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamInsertObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Filter/MultiInsertObject" => {
                    #[allow(non_camel_case_types)]
                    struct MultiInsertObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::insert::MultiObjectRequest,
                    > for MultiInsertObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::Locations;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::insert::MultiObjectRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::multi_insert_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiInsertObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Filter/UpdateObject" => {
                    #[allow(non_camel_case_types)]
                    struct UpdateObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::update::ObjectRequest,
                    > for UpdateObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::Location;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::update::ObjectRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::update_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = UpdateObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Filter/StreamUpdateObject" => {
                    #[allow(non_camel_case_types)]
                    struct StreamUpdateObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::update::ObjectRequest,
                    > for StreamUpdateObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamLocation;
                        type ResponseStream = T::StreamUpdateObjectStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::update::ObjectRequest,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::stream_update_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamUpdateObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Filter/MultiUpdateObject" => {
                    #[allow(non_camel_case_types)]
                    struct MultiUpdateObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::update::MultiObjectRequest,
                    > for MultiUpdateObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::Locations;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::update::MultiObjectRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::multi_update_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiUpdateObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Filter/UpsertObject" => {
                    #[allow(non_camel_case_types)]
                    struct UpsertObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::upsert::ObjectRequest,
                    > for UpsertObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::Location;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::upsert::ObjectRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::upsert_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = UpsertObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Filter/StreamUpsertObject" => {
                    #[allow(non_camel_case_types)]
                    struct StreamUpsertObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::upsert::ObjectRequest,
                    > for StreamUpsertObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamLocation;
                        type ResponseStream = T::StreamUpsertObjectStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::upsert::ObjectRequest,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::stream_upsert_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamUpsertObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Filter/MultiUpsertObject" => {
                    #[allow(non_camel_case_types)]
                    struct MultiUpsertObjectSvc<T: Filter>(pub Arc<T>);
                    impl<
                        T: Filter,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::upsert::MultiObjectRequest,
                    > for MultiUpsertObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::Locations;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::upsert::MultiObjectRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Filter>::multi_upsert_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiUpsertObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        let mut response = http::Response::new(
                            tonic::body::Body::default(),
                        );
                        let headers = response.headers_mut();
                        headers
                            .insert(
                                tonic::Status::GRPC_STATUS,
                                (tonic::Code::Unimplemented as i32).into(),
                            );
                        headers
                            .insert(
                                http::header::CONTENT_TYPE,
                                tonic::metadata::GRPC_CONTENT_TYPE,
                            );
                        Ok(response)
                    })
                }
            }
        }
    }
    impl<T> Clone for FilterServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    /// Generated gRPC service name
    pub const SERVICE_NAME: &str = "vald.v1.Filter";
    impl<T> tonic::server::NamedService for FilterServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod flush_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct FlushClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl FlushClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> FlushClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::Body>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> FlushClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::Body>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::Body>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::Body>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            FlushClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn flush(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::flush::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::info::index::Count>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Flush/Flush");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Flush", "Flush"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod flush_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with FlushServer.
    #[async_trait]
    pub trait Flush: std::marker::Send + std::marker::Sync + 'static {
        async fn flush(
            &self,
            request: tonic::Request<super::super::super::payload::v1::flush::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::info::index::Count>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct FlushServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> FlushServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for FlushServer<T>
    where
        T: Flush,
        B: Body + std::marker::Send + 'static,
        B::Error: Into<StdError> + std::marker::Send + 'static,
    {
        type Response = http::Response<tonic::body::Body>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/vald.v1.Flush/Flush" => {
                    #[allow(non_camel_case_types)]
                    struct FlushSvc<T: Flush>(pub Arc<T>);
                    impl<
                        T: Flush,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::flush::Request,
                    > for FlushSvc<T> {
                        type Response = super::super::super::payload::v1::info::index::Count;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::flush::Request,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Flush>::flush(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = FlushSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        let mut response = http::Response::new(
                            tonic::body::Body::default(),
                        );
                        let headers = response.headers_mut();
                        headers
                            .insert(
                                tonic::Status::GRPC_STATUS,
                                (tonic::Code::Unimplemented as i32).into(),
                            );
                        headers
                            .insert(
                                http::header::CONTENT_TYPE,
                                tonic::metadata::GRPC_CONTENT_TYPE,
                            );
                        Ok(response)
                    })
                }
            }
        }
    }
    impl<T> Clone for FlushServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    /// Generated gRPC service name
    pub const SERVICE_NAME: &str = "vald.v1.Flush";
    impl<T> tonic::server::NamedService for FlushServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod index_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct IndexClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl IndexClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> IndexClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::Body>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> IndexClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::Body>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::Body>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::Body>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            IndexClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn index_info(
            &mut self,
            request: impl tonic::IntoRequest<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::info::index::Count>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Index/IndexInfo");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Index", "IndexInfo"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn index_detail(
            &mut self,
            request: impl tonic::IntoRequest<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::info::index::Detail>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Index/IndexDetail",
            );
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Index", "IndexDetail"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 Represent the RPC to get the index statistics.
*/
        pub async fn index_statistics(
            &mut self,
            request: impl tonic::IntoRequest<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::info::index::Statistics>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Index/IndexStatistics",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Index", "IndexStatistics"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 Represent the RPC to get the index statistics for each agents.
*/
        pub async fn index_statistics_detail(
            &mut self,
            request: impl tonic::IntoRequest<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<
                super::super::super::payload::v1::info::index::StatisticsDetail,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Index/IndexStatisticsDetail",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Index", "IndexStatisticsDetail"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 Represent the RPC to get the index property.
*/
        pub async fn index_property(
            &mut self,
            request: impl tonic::IntoRequest<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<
                super::super::super::payload::v1::info::index::PropertyDetail,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Index/IndexProperty",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Index", "IndexProperty"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod index_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with IndexServer.
    #[async_trait]
    pub trait Index: std::marker::Send + std::marker::Sync + 'static {
        async fn index_info(
            &self,
            request: tonic::Request<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::info::index::Count>,
            tonic::Status,
        >;
        async fn index_detail(
            &self,
            request: tonic::Request<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::info::index::Detail>,
            tonic::Status,
        >;
        /** Overview
 Represent the RPC to get the index statistics.
*/
        async fn index_statistics(
            &self,
            request: tonic::Request<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::info::index::Statistics>,
            tonic::Status,
        >;
        /** Overview
 Represent the RPC to get the index statistics for each agents.
*/
        async fn index_statistics_detail(
            &self,
            request: tonic::Request<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<
                super::super::super::payload::v1::info::index::StatisticsDetail,
            >,
            tonic::Status,
        >;
        /** Overview
 Represent the RPC to get the index property.
*/
        async fn index_property(
            &self,
            request: tonic::Request<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<
                super::super::super::payload::v1::info::index::PropertyDetail,
            >,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct IndexServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> IndexServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for IndexServer<T>
    where
        T: Index,
        B: Body + std::marker::Send + 'static,
        B::Error: Into<StdError> + std::marker::Send + 'static,
    {
        type Response = http::Response<tonic::body::Body>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/vald.v1.Index/IndexInfo" => {
                    #[allow(non_camel_case_types)]
                    struct IndexInfoSvc<T: Index>(pub Arc<T>);
                    impl<
                        T: Index,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::Empty,
                    > for IndexInfoSvc<T> {
                        type Response = super::super::super::payload::v1::info::index::Count;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::Empty,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Index>::index_info(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = IndexInfoSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Index/IndexDetail" => {
                    #[allow(non_camel_case_types)]
                    struct IndexDetailSvc<T: Index>(pub Arc<T>);
                    impl<
                        T: Index,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::Empty,
                    > for IndexDetailSvc<T> {
                        type Response = super::super::super::payload::v1::info::index::Detail;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::Empty,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Index>::index_detail(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = IndexDetailSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Index/IndexStatistics" => {
                    #[allow(non_camel_case_types)]
                    struct IndexStatisticsSvc<T: Index>(pub Arc<T>);
                    impl<
                        T: Index,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::Empty,
                    > for IndexStatisticsSvc<T> {
                        type Response = super::super::super::payload::v1::info::index::Statistics;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::Empty,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Index>::index_statistics(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = IndexStatisticsSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Index/IndexStatisticsDetail" => {
                    #[allow(non_camel_case_types)]
                    struct IndexStatisticsDetailSvc<T: Index>(pub Arc<T>);
                    impl<
                        T: Index,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::Empty,
                    > for IndexStatisticsDetailSvc<T> {
                        type Response = super::super::super::payload::v1::info::index::StatisticsDetail;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::Empty,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Index>::index_statistics_detail(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = IndexStatisticsDetailSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Index/IndexProperty" => {
                    #[allow(non_camel_case_types)]
                    struct IndexPropertySvc<T: Index>(pub Arc<T>);
                    impl<
                        T: Index,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::Empty,
                    > for IndexPropertySvc<T> {
                        type Response = super::super::super::payload::v1::info::index::PropertyDetail;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::Empty,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Index>::index_property(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = IndexPropertySvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        let mut response = http::Response::new(
                            tonic::body::Body::default(),
                        );
                        let headers = response.headers_mut();
                        headers
                            .insert(
                                tonic::Status::GRPC_STATUS,
                                (tonic::Code::Unimplemented as i32).into(),
                            );
                        headers
                            .insert(
                                http::header::CONTENT_TYPE,
                                tonic::metadata::GRPC_CONTENT_TYPE,
                            );
                        Ok(response)
                    })
                }
            }
        }
    }
    impl<T> Clone for IndexServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    /// Generated gRPC service name
    pub const SERVICE_NAME: &str = "vald.v1.Index";
    impl<T> tonic::server::NamedService for IndexServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod insert_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct InsertClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl InsertClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> InsertClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::Body>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InsertClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::Body>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::Body>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::Body>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            InsertClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn insert(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::insert::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Insert/Insert");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Insert", "Insert"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamInsert RPC is the method to add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
 Each Insert request and response are independent.
 It's the recommended method to insert a large number of vectors.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn stream_insert(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::insert::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::object::StreamLocation,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Insert/StreamInsert",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Insert", "StreamInsert"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiInsert RPC is the method to add multiple new vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn multi_insert(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::insert::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Insert/MultiInsert",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Insert", "MultiInsert"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod insert_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with InsertServer.
    #[async_trait]
    pub trait Insert: std::marker::Send + std::marker::Sync + 'static {
        async fn insert(
            &self,
            request: tonic::Request<super::super::super::payload::v1::insert::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamInsert method.
        type StreamInsertStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamInsert RPC is the method to add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the insert request can be communicated in any order between client and server.
 Each Insert request and response are independent.
 It's the recommended method to insert a large number of vectors.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn stream_insert(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::insert::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamInsertStream>,
            tonic::Status,
        >;
        /** Overview
 MultiInsert RPC is the method to add multiple new vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | ALREADY_EXISTS    | Request ID is already inserted.                                                                                                                     | Change request ID.                                                                       |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn multi_insert(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::insert::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct InsertServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> InsertServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for InsertServer<T>
    where
        T: Insert,
        B: Body + std::marker::Send + 'static,
        B::Error: Into<StdError> + std::marker::Send + 'static,
    {
        type Response = http::Response<tonic::body::Body>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/vald.v1.Insert/Insert" => {
                    #[allow(non_camel_case_types)]
                    struct InsertSvc<T: Insert>(pub Arc<T>);
                    impl<
                        T: Insert,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::insert::Request,
                    > for InsertSvc<T> {
                        type Response = super::super::super::payload::v1::object::Location;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::insert::Request,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Insert>::insert(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = InsertSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Insert/StreamInsert" => {
                    #[allow(non_camel_case_types)]
                    struct StreamInsertSvc<T: Insert>(pub Arc<T>);
                    impl<
                        T: Insert,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::insert::Request,
                    > for StreamInsertSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamLocation;
                        type ResponseStream = T::StreamInsertStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::insert::Request,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Insert>::stream_insert(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamInsertSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Insert/MultiInsert" => {
                    #[allow(non_camel_case_types)]
                    struct MultiInsertSvc<T: Insert>(pub Arc<T>);
                    impl<
                        T: Insert,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::insert::MultiRequest,
                    > for MultiInsertSvc<T> {
                        type Response = super::super::super::payload::v1::object::Locations;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::insert::MultiRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Insert>::multi_insert(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiInsertSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        let mut response = http::Response::new(
                            tonic::body::Body::default(),
                        );
                        let headers = response.headers_mut();
                        headers
                            .insert(
                                tonic::Status::GRPC_STATUS,
                                (tonic::Code::Unimplemented as i32).into(),
                            );
                        headers
                            .insert(
                                http::header::CONTENT_TYPE,
                                tonic::metadata::GRPC_CONTENT_TYPE,
                            );
                        Ok(response)
                    })
                }
            }
        }
    }
    impl<T> Clone for InsertServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    /// Generated gRPC service name
    pub const SERVICE_NAME: &str = "vald.v1.Insert";
    impl<T> tonic::server::NamedService for InsertServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod object_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct ObjectClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl ObjectClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> ObjectClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::Body>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> ObjectClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::Body>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::Body>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::Body>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            ObjectClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn exists(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::object::Id,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Id>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Object/Exists");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Object", "Exists"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn get_object(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::object::VectorRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Vector>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Object/GetObject");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Object", "GetObject"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamGetObject RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the GetObject request can be communicated in any order between client and server.
 Each Upsert request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                   | how to resolve                                                                           |
 | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn stream_get_object(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::object::VectorRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::object::StreamVector,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Object/StreamGetObject",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Object", "StreamGetObject"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 A method to get all the vectors with server streaming
 ---
 Status Code
 TODO
 ---
 Troubleshooting
 TODO
*/
        pub async fn stream_list_object(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::object::list::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::object::list::Response,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Object/StreamListObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Object", "StreamListObject"));
            self.inner.server_streaming(req, path, codec).await
        }
        /** Overview
 Represent the RPC to get the vector metadata. This RPC is mainly used for index correction process
 ---
 Status Code
 TODO
 ---
 Troubleshooting
 TODO
*/
        pub async fn get_timestamp(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::object::TimestampRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Timestamp>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Object/GetTimestamp",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Object", "GetTimestamp"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod object_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with ObjectServer.
    #[async_trait]
    pub trait Object: std::marker::Send + std::marker::Sync + 'static {
        async fn exists(
            &self,
            request: tonic::Request<super::super::super::payload::v1::object::Id>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Id>,
            tonic::Status,
        >;
        async fn get_object(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::object::VectorRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Vector>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamGetObject method.
        type StreamGetObjectStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamVector,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamGetObject RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the GetObject request can be communicated in any order between client and server.
 Each Upsert request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                   | how to resolve                                                                           |
 | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn stream_get_object(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::object::VectorRequest>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamGetObjectStream>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamListObject method.
        type StreamListObjectStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::list::Response,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 A method to get all the vectors with server streaming
 ---
 Status Code
 TODO
 ---
 Troubleshooting
 TODO
*/
        async fn stream_list_object(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::object::list::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamListObjectStream>,
            tonic::Status,
        >;
        /** Overview
 Represent the RPC to get the vector metadata. This RPC is mainly used for index correction process
 ---
 Status Code
 TODO
 ---
 Troubleshooting
 TODO
*/
        async fn get_timestamp(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::object::TimestampRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Timestamp>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct ObjectServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> ObjectServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for ObjectServer<T>
    where
        T: Object,
        B: Body + std::marker::Send + 'static,
        B::Error: Into<StdError> + std::marker::Send + 'static,
    {
        type Response = http::Response<tonic::body::Body>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/vald.v1.Object/Exists" => {
                    #[allow(non_camel_case_types)]
                    struct ExistsSvc<T: Object>(pub Arc<T>);
                    impl<
                        T: Object,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::object::Id,
                    > for ExistsSvc<T> {
                        type Response = super::super::super::payload::v1::object::Id;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::object::Id,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Object>::exists(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = ExistsSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Object/GetObject" => {
                    #[allow(non_camel_case_types)]
                    struct GetObjectSvc<T: Object>(pub Arc<T>);
                    impl<
                        T: Object,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::object::VectorRequest,
                    > for GetObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::Vector;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::object::VectorRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Object>::get_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Object/StreamGetObject" => {
                    #[allow(non_camel_case_types)]
                    struct StreamGetObjectSvc<T: Object>(pub Arc<T>);
                    impl<
                        T: Object,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::object::VectorRequest,
                    > for StreamGetObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamVector;
                        type ResponseStream = T::StreamGetObjectStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::object::VectorRequest,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Object>::stream_get_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamGetObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Object/StreamListObject" => {
                    #[allow(non_camel_case_types)]
                    struct StreamListObjectSvc<T: Object>(pub Arc<T>);
                    impl<
                        T: Object,
                    > tonic::server::ServerStreamingService<
                        super::super::super::payload::v1::object::list::Request,
                    > for StreamListObjectSvc<T> {
                        type Response = super::super::super::payload::v1::object::list::Response;
                        type ResponseStream = T::StreamListObjectStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::object::list::Request,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Object>::stream_list_object(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamListObjectSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.server_streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Object/GetTimestamp" => {
                    #[allow(non_camel_case_types)]
                    struct GetTimestampSvc<T: Object>(pub Arc<T>);
                    impl<
                        T: Object,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::object::TimestampRequest,
                    > for GetTimestampSvc<T> {
                        type Response = super::super::super::payload::v1::object::Timestamp;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::object::TimestampRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Object>::get_timestamp(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetTimestampSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        let mut response = http::Response::new(
                            tonic::body::Body::default(),
                        );
                        let headers = response.headers_mut();
                        headers
                            .insert(
                                tonic::Status::GRPC_STATUS,
                                (tonic::Code::Unimplemented as i32).into(),
                            );
                        headers
                            .insert(
                                http::header::CONTENT_TYPE,
                                tonic::metadata::GRPC_CONTENT_TYPE,
                            );
                        Ok(response)
                    })
                }
            }
        }
    }
    impl<T> Clone for ObjectServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    /// Generated gRPC service name
    pub const SERVICE_NAME: &str = "vald.v1.Object";
    impl<T> tonic::server::NamedService for ObjectServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod remove_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct RemoveClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl RemoveClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> RemoveClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::Body>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> RemoveClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::Body>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::Body>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::Body>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            RemoveClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn remove(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::remove::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Remove/Remove");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Remove", "Remove"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 RemoveByTimestamp RPC is the method to remove vectors based on timestamp.

 <div class="notice">
 In the TimestampRequest message, the 'timestamps' field is repeated, allowing the inclusion of multiple Timestamp.<br>
 When multiple Timestamps are provided, it results in an `AND` condition, enabling the realization of deletions with specified ranges.<br>
 This design allows for versatile deletion operations, facilitating tasks such as removing data within a specific time range.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                   | how to resolve                                                                                                       |
 | :---------------- | :---------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.                              |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed.                             |
 | NOT_FOUND         | No vectors in the system match the specified timestamp conditions.                              | Check whether vectors matching the specified timestamp conditions exist in the system, and fix conditions if needed. |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.
*/
        pub async fn remove_by_timestamp(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::remove::TimestampRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Remove/RemoveByTimestamp",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Remove", "RemoveByTimestamp"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 A method to remove multiple indexed vectors by bidirectional streaming.

 StreamRemove RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the remove request can be communicated in any order between client and server.
 Each Remove request and response are independent.
 It's the recommended method to remove a large number of vectors.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
   The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                   | how to resolve                                                                           |
 | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn stream_remove(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::remove::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::object::StreamLocation,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Remove/StreamRemove",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Remove", "StreamRemove"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiRemove is the method to remove multiple vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                   | how to resolve                                                                           |
 | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn multi_remove(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::remove::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Remove/MultiRemove",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Remove", "MultiRemove"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod remove_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with RemoveServer.
    #[async_trait]
    pub trait Remove: std::marker::Send + std::marker::Sync + 'static {
        async fn remove(
            &self,
            request: tonic::Request<super::super::super::payload::v1::remove::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /** Overview
 RemoveByTimestamp RPC is the method to remove vectors based on timestamp.

 <div class="notice">
 In the TimestampRequest message, the 'timestamps' field is repeated, allowing the inclusion of multiple Timestamp.<br>
 When multiple Timestamps are provided, it results in an `AND` condition, enabling the realization of deletions with specified ranges.<br>
 This design allows for versatile deletion operations, facilitating tasks such as removing data within a specific time range.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                   | how to resolve                                                                                                       |
 | :---------------- | :---------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.                              |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed.                             |
 | NOT_FOUND         | No vectors in the system match the specified timestamp conditions.                              | Check whether vectors matching the specified timestamp conditions exist in the system, and fix conditions if needed. |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.
*/
        async fn remove_by_timestamp(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::remove::TimestampRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamRemove method.
        type StreamRemoveStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 A method to remove multiple indexed vectors by bidirectional streaming.

 StreamRemove RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the remove request can be communicated in any order between client and server.
 Each Remove request and response are independent.
 It's the recommended method to remove a large number of vectors.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
   The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                   | how to resolve                                                                           |
 | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn stream_remove(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::remove::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamRemoveStream>,
            tonic::Status,
        >;
        /** Overview
 MultiRemove is the method to remove multiple vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                   | how to resolve                                                                           |
 | :---------------- | :---------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server. | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                         | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Requested ID is NOT inserted.                                                                   | Send a request with an ID that is already inserted.                                      |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn multi_remove(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::remove::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct RemoveServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> RemoveServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for RemoveServer<T>
    where
        T: Remove,
        B: Body + std::marker::Send + 'static,
        B::Error: Into<StdError> + std::marker::Send + 'static,
    {
        type Response = http::Response<tonic::body::Body>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/vald.v1.Remove/Remove" => {
                    #[allow(non_camel_case_types)]
                    struct RemoveSvc<T: Remove>(pub Arc<T>);
                    impl<
                        T: Remove,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::remove::Request,
                    > for RemoveSvc<T> {
                        type Response = super::super::super::payload::v1::object::Location;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::remove::Request,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Remove>::remove(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = RemoveSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Remove/RemoveByTimestamp" => {
                    #[allow(non_camel_case_types)]
                    struct RemoveByTimestampSvc<T: Remove>(pub Arc<T>);
                    impl<
                        T: Remove,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::remove::TimestampRequest,
                    > for RemoveByTimestampSvc<T> {
                        type Response = super::super::super::payload::v1::object::Locations;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::remove::TimestampRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Remove>::remove_by_timestamp(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = RemoveByTimestampSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Remove/StreamRemove" => {
                    #[allow(non_camel_case_types)]
                    struct StreamRemoveSvc<T: Remove>(pub Arc<T>);
                    impl<
                        T: Remove,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::remove::Request,
                    > for StreamRemoveSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamLocation;
                        type ResponseStream = T::StreamRemoveStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::remove::Request,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Remove>::stream_remove(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamRemoveSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Remove/MultiRemove" => {
                    #[allow(non_camel_case_types)]
                    struct MultiRemoveSvc<T: Remove>(pub Arc<T>);
                    impl<
                        T: Remove,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::remove::MultiRequest,
                    > for MultiRemoveSvc<T> {
                        type Response = super::super::super::payload::v1::object::Locations;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::remove::MultiRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Remove>::multi_remove(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiRemoveSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        let mut response = http::Response::new(
                            tonic::body::Body::default(),
                        );
                        let headers = response.headers_mut();
                        headers
                            .insert(
                                tonic::Status::GRPC_STATUS,
                                (tonic::Code::Unimplemented as i32).into(),
                            );
                        headers
                            .insert(
                                http::header::CONTENT_TYPE,
                                tonic::metadata::GRPC_CONTENT_TYPE,
                            );
                        Ok(response)
                    })
                }
            }
        }
    }
    impl<T> Clone for RemoveServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    /// Generated gRPC service name
    pub const SERVICE_NAME: &str = "vald.v1.Remove";
    impl<T> tonic::server::NamedService for RemoveServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod search_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    /** Overview
 Search Service is responsible for searching vectors similar to the user request vector from `vald-agent`.
*/
    #[derive(Debug, Clone)]
    pub struct SearchClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl SearchClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> SearchClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::Body>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> SearchClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::Body>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::Body>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::Body>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            SearchClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        /** Overview
 Search RPC is the method to search vector(s) similar to the request vector.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn search(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::search::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Search/Search");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Search", "Search"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 SearchByID RPC is the method to search similar vectors using a user-defined vector ID.<br>
 The vector with the same requested ID should be indexed into the `vald-agent` before searching.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn search_by_id(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::search::IdRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/SearchByID",
            );
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Search", "SearchByID"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamSearch RPC is the method to search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server.
 Each Search request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn stream_search(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::search::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::search::StreamResponse,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/StreamSearch",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "StreamSearch"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 StreamSearchByID RPC is the method to search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server.
 Each SearchByID request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn stream_search_by_id(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::search::IdRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::search::StreamResponse,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/StreamSearchByID",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "StreamSearchByID"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiSearch RPC is the method to search vectors with multiple vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
   The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn multi_search(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::search::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/MultiSearch",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "MultiSearch"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 MultiSearchByID RPC is the method to search vectors with multiple IDs in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn multi_search_by_id(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::search::MultiIdRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/MultiSearchByID",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "MultiSearchByID"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 LinearSearch RPC is the method to linear search vector(s) similar to the request vector.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn linear_search(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::search::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/LinearSearch",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "LinearSearch"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 LinearSearchByID RPC is the method to linear search similar vectors using a user-defined vector ID.<br>
 The vector with the same requested ID should be indexed into the `vald-agent` before searching.
 You will get a `NOT_FOUND` error if the vector isn't stored.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn linear_search_by_id(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::search::IdRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/LinearSearchByID",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "LinearSearchByID"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamLinearSearch RPC is the method to linear search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server.
 Each LinearSearch request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn stream_linear_search(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::search::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::search::StreamResponse,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/StreamLinearSearch",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "StreamLinearSearch"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
   StreamLinearSearchByID RPC is the method to linear search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server.
 Each LinearSearchByID request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn stream_linear_search_by_id(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::search::IdRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::search::StreamResponse,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/StreamLinearSearchByID",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "StreamLinearSearchByID"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiLinearSearch RPC is the method to linear search vectors with multiple vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
   The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn multi_linear_search(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::search::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/MultiLinearSearch",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "MultiLinearSearch"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 MultiLinearSearchByID RPC is the method to linear search vectors with multiple IDs in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 // ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn multi_linear_search_by_id(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::search::MultiIdRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/MultiLinearSearchByID",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "MultiLinearSearchByID"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod search_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with SearchServer.
    #[async_trait]
    pub trait Search: std::marker::Send + std::marker::Sync + 'static {
        /** Overview
 Search RPC is the method to search vector(s) similar to the request vector.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn search(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /** Overview
 SearchByID RPC is the method to search similar vectors using a user-defined vector ID.<br>
 The vector with the same requested ID should be indexed into the `vald-agent` before searching.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn search_by_id(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::IdRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamSearch method.
        type StreamSearchStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamSearch RPC is the method to search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server.
 Each Search request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn stream_search(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::search::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamSearchStream>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamSearchByID method.
        type StreamSearchByIDStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamSearchByID RPC is the method to search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the search request can be communicated in any order between the client and server.
 Each SearchByID request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn stream_search_by_id(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::search::IdRequest>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamSearchByIDStream>,
            tonic::Status,
        >;
        /** Overview
 MultiSearch RPC is the method to search vectors with multiple vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
   The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn multi_search(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::search::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        >;
        /** Overview
 MultiSearchByID RPC is the method to search vectors with multiple IDs in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn multi_search_by_id(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::search::MultiIdRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        >;
        /** Overview
 LinearSearch RPC is the method to linear search vector(s) similar to the request vector.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn linear_search(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /** Overview
 LinearSearchByID RPC is the method to linear search similar vectors using a user-defined vector ID.<br>
 The vector with the same requested ID should be indexed into the `vald-agent` before searching.
 You will get a `NOT_FOUND` error if the vector isn't stored.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn linear_search_by_id(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::IdRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamLinearSearch method.
        type StreamLinearSearchStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamLinearSearch RPC is the method to linear search vectors with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server.
 Each LinearSearch request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn stream_linear_search(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::search::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamLinearSearchStream>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamLinearSearchByID method.
        type StreamLinearSearchByIDStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
   StreamLinearSearchByID RPC is the method to linear search vectors with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the linear search request can be communicated in any order between the client and server.
 Each LinearSearchByID request and response are independent.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn stream_linear_search_by_id(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::search::IdRequest>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamLinearSearchByIDStream>,
            tonic::Status,
        >;
        /** Overview
 MultiLinearSearch RPC is the method to linear search vectors with multiple vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
   The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                   | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                 | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                 | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Search result is empty or insufficient to request result length.                                                | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                   | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn multi_linear_search(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::search::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        >;
        /** Overview
 MultiLinearSearchByID RPC is the method to linear search vectors with multiple IDs in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 // ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                    | how to resolve                                                                           |
 | :---------------- | :------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                  | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Requested vector's ID is empty, or some request payload is invalid.                                                          | Check request payload and fix request payload.                                           |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                  | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | The Requested ID is not inserted on the target Vald cluster, or the search result is insufficient to the required result length. | Send a request with another vector or set min_num to a smaller value.                    |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                    | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn multi_linear_search_by_id(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::search::MultiIdRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        >;
    }
    /** Overview
 Search Service is responsible for searching vectors similar to the user request vector from `vald-agent`.
*/
    #[derive(Debug)]
    pub struct SearchServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> SearchServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for SearchServer<T>
    where
        T: Search,
        B: Body + std::marker::Send + 'static,
        B::Error: Into<StdError> + std::marker::Send + 'static,
    {
        type Response = http::Response<tonic::body::Body>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/vald.v1.Search/Search" => {
                    #[allow(non_camel_case_types)]
                    struct SearchSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::Request,
                    > for SearchSvc<T> {
                        type Response = super::super::super::payload::v1::search::Response;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::search::Request,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::search(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = SearchSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Search/SearchByID" => {
                    #[allow(non_camel_case_types)]
                    struct SearchByIDSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::IdRequest,
                    > for SearchByIDSvc<T> {
                        type Response = super::super::super::payload::v1::search::Response;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::search::IdRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::search_by_id(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = SearchByIDSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Search/StreamSearch" => {
                    #[allow(non_camel_case_types)]
                    struct StreamSearchSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::search::Request,
                    > for StreamSearchSvc<T> {
                        type Response = super::super::super::payload::v1::search::StreamResponse;
                        type ResponseStream = T::StreamSearchStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::search::Request,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::stream_search(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamSearchSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Search/StreamSearchByID" => {
                    #[allow(non_camel_case_types)]
                    struct StreamSearchByIDSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::search::IdRequest,
                    > for StreamSearchByIDSvc<T> {
                        type Response = super::super::super::payload::v1::search::StreamResponse;
                        type ResponseStream = T::StreamSearchByIDStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::search::IdRequest,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::stream_search_by_id(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamSearchByIDSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Search/MultiSearch" => {
                    #[allow(non_camel_case_types)]
                    struct MultiSearchSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::MultiRequest,
                    > for MultiSearchSvc<T> {
                        type Response = super::super::super::payload::v1::search::Responses;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::search::MultiRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::multi_search(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiSearchSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Search/MultiSearchByID" => {
                    #[allow(non_camel_case_types)]
                    struct MultiSearchByIDSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::MultiIdRequest,
                    > for MultiSearchByIDSvc<T> {
                        type Response = super::super::super::payload::v1::search::Responses;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::search::MultiIdRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::multi_search_by_id(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiSearchByIDSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Search/LinearSearch" => {
                    #[allow(non_camel_case_types)]
                    struct LinearSearchSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::Request,
                    > for LinearSearchSvc<T> {
                        type Response = super::super::super::payload::v1::search::Response;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::search::Request,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::linear_search(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = LinearSearchSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Search/LinearSearchByID" => {
                    #[allow(non_camel_case_types)]
                    struct LinearSearchByIDSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::IdRequest,
                    > for LinearSearchByIDSvc<T> {
                        type Response = super::super::super::payload::v1::search::Response;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::search::IdRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::linear_search_by_id(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = LinearSearchByIDSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Search/StreamLinearSearch" => {
                    #[allow(non_camel_case_types)]
                    struct StreamLinearSearchSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::search::Request,
                    > for StreamLinearSearchSvc<T> {
                        type Response = super::super::super::payload::v1::search::StreamResponse;
                        type ResponseStream = T::StreamLinearSearchStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::search::Request,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::stream_linear_search(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamLinearSearchSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Search/StreamLinearSearchByID" => {
                    #[allow(non_camel_case_types)]
                    struct StreamLinearSearchByIDSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::search::IdRequest,
                    > for StreamLinearSearchByIDSvc<T> {
                        type Response = super::super::super::payload::v1::search::StreamResponse;
                        type ResponseStream = T::StreamLinearSearchByIDStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::search::IdRequest,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::stream_linear_search_by_id(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamLinearSearchByIDSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Search/MultiLinearSearch" => {
                    #[allow(non_camel_case_types)]
                    struct MultiLinearSearchSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::MultiRequest,
                    > for MultiLinearSearchSvc<T> {
                        type Response = super::super::super::payload::v1::search::Responses;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::search::MultiRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::multi_linear_search(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiLinearSearchSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Search/MultiLinearSearchByID" => {
                    #[allow(non_camel_case_types)]
                    struct MultiLinearSearchByIDSvc<T: Search>(pub Arc<T>);
                    impl<
                        T: Search,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::MultiIdRequest,
                    > for MultiLinearSearchByIDSvc<T> {
                        type Response = super::super::super::payload::v1::search::Responses;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::search::MultiIdRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Search>::multi_linear_search_by_id(&inner, request)
                                    .await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiLinearSearchByIDSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        let mut response = http::Response::new(
                            tonic::body::Body::default(),
                        );
                        let headers = response.headers_mut();
                        headers
                            .insert(
                                tonic::Status::GRPC_STATUS,
                                (tonic::Code::Unimplemented as i32).into(),
                            );
                        headers
                            .insert(
                                http::header::CONTENT_TYPE,
                                tonic::metadata::GRPC_CONTENT_TYPE,
                            );
                        Ok(response)
                    })
                }
            }
        }
    }
    impl<T> Clone for SearchServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    /// Generated gRPC service name
    pub const SERVICE_NAME: &str = "vald.v1.Search";
    impl<T> tonic::server::NamedService for SearchServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod update_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct UpdateClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl UpdateClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> UpdateClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::Body>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> UpdateClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::Body>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::Body>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::Body>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            UpdateClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn update(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::update::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Update/Update");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Update", "Update"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamUpdate RPC is the method to update multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the update request can be communicated in any order between client and server.
 Each Update request and response are independent.
 It's the recommended method to update the large amount of vectors.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
 | ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn stream_update(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::update::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::object::StreamLocation,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Update/StreamUpdate",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Update", "StreamUpdate"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiUpdate is the method to update multiple vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
 | ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn multi_update(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::update::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Update/MultiUpdate",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Update", "MultiUpdate"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 A method to update timestamp an indexed vector.
 ---
 Status Code
 TODO
 ---
 Troubleshooting
 TODO
*/
        pub async fn update_timestamp(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::update::TimestampRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Update/UpdateTimestamp",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Update", "UpdateTimestamp"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod update_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with UpdateServer.
    #[async_trait]
    pub trait Update: std::marker::Send + std::marker::Sync + 'static {
        async fn update(
            &self,
            request: tonic::Request<super::super::super::payload::v1::update::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamUpdate method.
        type StreamUpdateStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamUpdate RPC is the method to update multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the update request can be communicated in any order between client and server.
 Each Update request and response are independent.
 It's the recommended method to update the large amount of vectors.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
 | ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn stream_update(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::update::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamUpdateStream>,
            tonic::Status,
        >;
        /** Overview
 MultiUpdate is the method to update multiple vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | NOT_FOUND         | Requested ID is NOT inserted.                                                                                                                       | Send a request with an ID that is already inserted.                                      |
 | ALREADY_EXISTS    | Request pair of ID and vector is already inserted.                                                                                                  | Change request ID.                                                                       |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn multi_update(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::update::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
        /** Overview
 A method to update timestamp an indexed vector.
 ---
 Status Code
 TODO
 ---
 Troubleshooting
 TODO
*/
        async fn update_timestamp(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::update::TimestampRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct UpdateServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> UpdateServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for UpdateServer<T>
    where
        T: Update,
        B: Body + std::marker::Send + 'static,
        B::Error: Into<StdError> + std::marker::Send + 'static,
    {
        type Response = http::Response<tonic::body::Body>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/vald.v1.Update/Update" => {
                    #[allow(non_camel_case_types)]
                    struct UpdateSvc<T: Update>(pub Arc<T>);
                    impl<
                        T: Update,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::update::Request,
                    > for UpdateSvc<T> {
                        type Response = super::super::super::payload::v1::object::Location;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::update::Request,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Update>::update(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = UpdateSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Update/StreamUpdate" => {
                    #[allow(non_camel_case_types)]
                    struct StreamUpdateSvc<T: Update>(pub Arc<T>);
                    impl<
                        T: Update,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::update::Request,
                    > for StreamUpdateSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamLocation;
                        type ResponseStream = T::StreamUpdateStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::update::Request,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Update>::stream_update(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamUpdateSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Update/MultiUpdate" => {
                    #[allow(non_camel_case_types)]
                    struct MultiUpdateSvc<T: Update>(pub Arc<T>);
                    impl<
                        T: Update,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::update::MultiRequest,
                    > for MultiUpdateSvc<T> {
                        type Response = super::super::super::payload::v1::object::Locations;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::update::MultiRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Update>::multi_update(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiUpdateSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Update/UpdateTimestamp" => {
                    #[allow(non_camel_case_types)]
                    struct UpdateTimestampSvc<T: Update>(pub Arc<T>);
                    impl<
                        T: Update,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::update::TimestampRequest,
                    > for UpdateTimestampSvc<T> {
                        type Response = super::super::super::payload::v1::object::Location;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::update::TimestampRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Update>::update_timestamp(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = UpdateTimestampSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        let mut response = http::Response::new(
                            tonic::body::Body::default(),
                        );
                        let headers = response.headers_mut();
                        headers
                            .insert(
                                tonic::Status::GRPC_STATUS,
                                (tonic::Code::Unimplemented as i32).into(),
                            );
                        headers
                            .insert(
                                http::header::CONTENT_TYPE,
                                tonic::metadata::GRPC_CONTENT_TYPE,
                            );
                        Ok(response)
                    })
                }
            }
        }
    }
    impl<T> Clone for UpdateServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    /// Generated gRPC service name
    pub const SERVICE_NAME: &str = "vald.v1.Update";
    impl<T> tonic::server::NamedService for UpdateServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod upsert_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct UpsertClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl UpsertClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> UpsertClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::Body>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> UpsertClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::Body>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::Body>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::Body>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            UpsertClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        pub async fn upsert(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::upsert::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Upsert/Upsert");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Upsert", "Upsert"));
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamUpsert RPC is the method to update multiple existing vectors or add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the upsert request can be communicated in any order between the client and server.
 Each Upsert request and response are independent.
 It’s the recommended method to upsert a large number of vectors.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn stream_upsert(
            &mut self,
            request: impl tonic::IntoStreamingRequest<
                Message = super::super::super::payload::v1::upsert::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<
                tonic::codec::Streaming<
                    super::super::super::payload::v1::object::StreamLocation,
                >,
            >,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Upsert/StreamUpsert",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Upsert", "StreamUpsert"));
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiUpsert is the method to update existing multiple vectors and add new multiple vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        pub async fn multi_upsert(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::upsert::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Upsert/MultiUpsert",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Upsert", "MultiUpsert"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod upsert_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with UpsertServer.
    #[async_trait]
    pub trait Upsert: std::marker::Send + std::marker::Sync + 'static {
        async fn upsert(
            &self,
            request: tonic::Request<super::super::super::payload::v1::upsert::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamUpsert method.
        type StreamUpsertStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamUpsert RPC is the method to update multiple existing vectors or add new multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
 Using the bidirectional streaming RPC, the upsert request can be communicated in any order between the client and server.
 Each Upsert request and response are independent.
 It’s the recommended method to upsert a large number of vectors.
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn stream_upsert(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::upsert::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamUpsertStream>,
            tonic::Status,
        >;
        /** Overview
 MultiUpsert is the method to update existing multiple vectors and add new multiple vectors in **1** request.

 <div class="notice">
 gRPC has a message size limitation.<br>
 Please be careful that the size of the request exceeds the limit.
 </div>
 ---
 Status Code
 |  0   | OK                |
 |  1   | CANCELLED         |
 |  3   | INVALID_ARGUMENT  |
 |  4   | DEADLINE_EXCEEDED |
 |  5   | NOT_FOUND         |
 |  6   | ALREADY_EXISTS    |
 |  10  | ABORTED           |
 |  13  | INTERNAL          |
 ---
 Troubleshooting
 The request process may not be completed when the response code is NOT `0 (OK)`.

 Here are some common reasons and how to resolve each error.

 | name              | common reason                                                                                                                                       | how to resolve                                                                           |
 | :---------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------- |
 | CANCELLED         | Executed cancel() of rpc from client/server-side or network problems between client and server.                                                     | Check the code, especially around timeout and connection management, and fix if needed.  |
 | INVALID_ARGUMENT  | The Dimension of the request vector is NOT the same as Vald Agent's config, the requested vector's ID is empty, or some request payload is invalid. | Check Agent config, request payload, and fix request payload or Agent config.            |
 | DEADLINE_EXCEEDED | The RPC timeout setting is too short on the client/server side.                                                                                     | Check the gRPC timeout setting on both the client and server sides and fix it if needed. |
 | ALREADY_EXISTS    | Requested pair of ID and vector is already inserted                                                                                                 | Change request payload or nothing to do if update is unnecessary.                        |
 | INTERNAL          | Target Vald cluster or network route has some critical error.                                                                                       | Check target Vald cluster first and check network route including ingress as second.     |
*/
        async fn multi_upsert(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::upsert::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct UpsertServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> UpsertServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            Self {
                inner,
                accept_compression_encodings: Default::default(),
                send_compression_encodings: Default::default(),
                max_decoding_message_size: None,
                max_encoding_message_size: None,
            }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> InterceptedService<Self, F>
        where
            F: tonic::service::Interceptor,
        {
            InterceptedService::new(Self::new(inner), interceptor)
        }
        /// Enable decompressing requests with the given encoding.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.accept_compression_encodings.enable(encoding);
            self
        }
        /// Compress responses with the given encoding, if the client supports it.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.send_compression_encodings.enable(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.max_decoding_message_size = Some(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.max_encoding_message_size = Some(limit);
            self
        }
    }
    impl<T, B> tonic::codegen::Service<http::Request<B>> for UpsertServer<T>
    where
        T: Upsert,
        B: Body + std::marker::Send + 'static,
        B::Error: Into<StdError> + std::marker::Send + 'static,
    {
        type Response = http::Response<tonic::body::Body>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            match req.uri().path() {
                "/vald.v1.Upsert/Upsert" => {
                    #[allow(non_camel_case_types)]
                    struct UpsertSvc<T: Upsert>(pub Arc<T>);
                    impl<
                        T: Upsert,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::upsert::Request,
                    > for UpsertSvc<T> {
                        type Response = super::super::super::payload::v1::object::Location;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::upsert::Request,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Upsert>::upsert(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = UpsertSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Upsert/StreamUpsert" => {
                    #[allow(non_camel_case_types)]
                    struct StreamUpsertSvc<T: Upsert>(pub Arc<T>);
                    impl<
                        T: Upsert,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::upsert::Request,
                    > for StreamUpsertSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamLocation;
                        type ResponseStream = T::StreamUpsertStream;
                        type Future = BoxFuture<
                            tonic::Response<Self::ResponseStream>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                tonic::Streaming<
                                    super::super::super::payload::v1::upsert::Request,
                                >,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Upsert>::stream_upsert(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = StreamUpsertSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.streaming(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                "/vald.v1.Upsert/MultiUpsert" => {
                    #[allow(non_camel_case_types)]
                    struct MultiUpsertSvc<T: Upsert>(pub Arc<T>);
                    impl<
                        T: Upsert,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::upsert::MultiRequest,
                    > for MultiUpsertSvc<T> {
                        type Response = super::super::super::payload::v1::object::Locations;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::upsert::MultiRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                <T as Upsert>::multi_upsert(&inner, request).await
                            };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = MultiUpsertSvc(inner);
                        let codec = tonic_prost::ProstCodec::default();
                        let mut grpc = tonic::server::Grpc::new(codec)
                            .apply_compression_config(
                                accept_compression_encodings,
                                send_compression_encodings,
                            )
                            .apply_max_message_size_config(
                                max_decoding_message_size,
                                max_encoding_message_size,
                            );
                        let res = grpc.unary(method, req).await;
                        Ok(res)
                    };
                    Box::pin(fut)
                }
                _ => {
                    Box::pin(async move {
                        let mut response = http::Response::new(
                            tonic::body::Body::default(),
                        );
                        let headers = response.headers_mut();
                        headers
                            .insert(
                                tonic::Status::GRPC_STATUS,
                                (tonic::Code::Unimplemented as i32).into(),
                            );
                        headers
                            .insert(
                                http::header::CONTENT_TYPE,
                                tonic::metadata::GRPC_CONTENT_TYPE,
                            );
                        Ok(response)
                    })
                }
            }
        }
    }
    impl<T> Clone for UpsertServer<T> {
        fn clone(&self) -> Self {
            let inner = self.inner.clone();
            Self {
                inner,
                accept_compression_encodings: self.accept_compression_encodings,
                send_compression_encodings: self.send_compression_encodings,
                max_decoding_message_size: self.max_decoding_message_size,
                max_encoding_message_size: self.max_encoding_message_size,
            }
        }
    }
    /// Generated gRPC service name
    pub const SERVICE_NAME: &str = "vald.v1.Upsert";
    impl<T> tonic::server::NamedService for UpsertServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
