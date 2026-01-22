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
pub mod meta_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value
    )]
    use tonic::codegen::http::Uri;
    use tonic::codegen::*;
    #[derive(Debug, Clone)]
    pub struct MetaClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl MetaClient<tonic::transport::Channel> {
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
    impl<T> MetaClient<T>
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
        pub fn with_interceptor<F>(inner: T, interceptor: F) -> MetaClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                    http::Request<tonic::body::Body>,
                    Response = http::Response<
                        <T as tonic::client::GrpcService<tonic::body::Body>>::ResponseBody,
                    >,
                >,
            <T as tonic::codegen::Service<http::Request<tonic::body::Body>>>::Error:
                Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            MetaClient::new(InterceptedService::new(inner, interceptor))
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
        pub async fn get(
            &mut self,
            request: impl tonic::IntoRequest<super::super::super::payload::v1::meta::Key>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::meta::Value>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::unknown(format!("Service was not ready: {}", e.into()))
            })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/meta.v1.Meta/Get");
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("meta.v1.Meta", "Get"));
            self.inner.unary(req, path, codec).await
        }
        ///
        pub async fn set(
            &mut self,
            request: impl tonic::IntoRequest<super::super::super::payload::v1::meta::KeyValue>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::Empty>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::unknown(format!("Service was not ready: {}", e.into()))
            })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/meta.v1.Meta/Set");
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("meta.v1.Meta", "Set"));
            self.inner.unary(req, path, codec).await
        }
        ///
        pub async fn delete(
            &mut self,
            request: impl tonic::IntoRequest<super::super::super::payload::v1::meta::Key>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::Empty>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::unknown(format!("Service was not ready: {}", e.into()))
            })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/meta.v1.Meta/Delete");
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("meta.v1.Meta", "Delete"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod meta_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with MetaServer.
    #[async_trait]
    pub trait Meta: std::marker::Send + std::marker::Sync + 'static {
        async fn get(
            &self,
            request: tonic::Request<super::super::super::payload::v1::meta::Key>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::meta::Value>,
            tonic::Status,
        >;
        ///
        async fn set(
            &self,
            request: tonic::Request<super::super::super::payload::v1::meta::KeyValue>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::Empty>,
            tonic::Status,
        >;
        ///
        async fn delete(
            &self,
            request: tonic::Request<super::super::super::payload::v1::meta::Key>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::Empty>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct MetaServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> MetaServer<T> {
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
        pub fn with_interceptor<F>(inner: T, interceptor: F) -> InterceptedService<Self, F>
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for MetaServer<T>
    where
        T: Meta,
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
                "/meta.v1.Meta/Get" => {
                    #[allow(non_camel_case_types)]
                    struct GetSvc<T: Meta>(pub Arc<T>);
                    impl<T: Meta>
                        tonic::server::UnaryService<super::super::super::payload::v1::meta::Key>
                        for GetSvc<T>
                    {
                        type Response = super::super::super::payload::v1::meta::Value;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::super::super::payload::v1::meta::Key>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move { <T as Meta>::get(&inner, request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = GetSvc(inner);
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
                "/meta.v1.Meta/Set" => {
                    #[allow(non_camel_case_types)]
                    struct SetSvc<T: Meta>(pub Arc<T>);
                    impl<T: Meta>
                        tonic::server::UnaryService<
                            super::super::super::payload::v1::meta::KeyValue,
                        > for SetSvc<T>
                    {
                        type Response = super::super::super::payload::v1::Empty;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::meta::KeyValue,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move { <T as Meta>::set(&inner, request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = SetSvc(inner);
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
                "/meta.v1.Meta/Delete" => {
                    #[allow(non_camel_case_types)]
                    struct DeleteSvc<T: Meta>(pub Arc<T>);
                    impl<T: Meta>
                        tonic::server::UnaryService<super::super::super::payload::v1::meta::Key>
                        for DeleteSvc<T>
                    {
                        type Response = super::super::super::payload::v1::Empty;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::super::super::payload::v1::meta::Key>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move { <T as Meta>::delete(&inner, request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = DeleteSvc(inner);
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
                _ => Box::pin(async move {
                    let mut response = http::Response::new(tonic::body::Body::default());
                    let headers = response.headers_mut();
                    headers.insert(
                        tonic::Status::GRPC_STATUS,
                        (tonic::Code::Unimplemented as i32).into(),
                    );
                    headers.insert(
                        http::header::CONTENT_TYPE,
                        tonic::metadata::GRPC_CONTENT_TYPE,
                    );
                    Ok(response)
                }),
            }
        }
    }
    impl<T> Clone for MetaServer<T> {
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
    pub const SERVICE_NAME: &str = "meta.v1.Meta";
    impl<T> tonic::server::NamedService for MetaServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod search_with_metadata_client {
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
    pub struct SearchWithMetadataClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl SearchWithMetadataClient<tonic::transport::Channel> {
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
    impl<T> SearchWithMetadataClient<T>
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
        ) -> SearchWithMetadataClient<InterceptedService<T, F>>
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
            SearchWithMetadataClient::new(InterceptedService::new(inner, interceptor))
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
 SearchWithMetadata RPC is the method to search vector(s) similar to the request vector and to get metadata(s).
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
        pub async fn search_with_metadata(
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
                "/meta.v1.SearchWithMetadata/SearchWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("meta.v1.SearchWithMetadata", "SearchWithMetadata"),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 SearchByIDWithMetadata RPC is the method to search similar vectors using a user-defined vector ID and to get metadata.<br>
 The vector with the same requested ID should be indexed into the `vald-lb-gateway` before searching.
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
        pub async fn search_by_id_with_metadata(
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
                "/meta.v1.SearchWithMetadata/SearchByIDWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.SearchWithMetadata",
                        "SearchByIDWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamSearchWithMetadata RPC is the method to search vectors and to get metadata with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        pub async fn stream_search_with_metadata(
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
                "/meta.v1.SearchWithMetadata/StreamSearchWithMetadata",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.SearchWithMetadata",
                        "StreamSearchWithMetadata",
                    ),
                );
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 StreamSearchByIDWithMetadata RPC is the method to search vectors and to get metadata with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        pub async fn stream_search_by_id_with_metadata(
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
                "/meta.v1.SearchWithMetadata/StreamSearchByIDWithMetadata",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.SearchWithMetadata",
                        "StreamSearchByIDWithMetadata",
                    ),
                );
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiSearchWithMetadata RPC is the method to search vectors and to get metadata with multiple vectors in **1** request.

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
        pub async fn multi_search_with_metadata(
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
                "/meta.v1.SearchWithMetadata/MultiSearchWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.SearchWithMetadata",
                        "MultiSearchWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 MultiSearchByIDWithMetadata RPC is the method to search vectors and to get metadata with multiple IDs in **1** request.

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
        pub async fn multi_search_by_id_with_metadata(
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
                "/meta.v1.SearchWithMetadata/MultiSearchByIDWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.SearchWithMetadata",
                        "MultiSearchByIDWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 LinearSearchWithMetadata RPC is the method to linear search vector(s) similar to the request vector and to get metadata.
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
        pub async fn linear_search_with_metadata(
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
                "/meta.v1.SearchWithMetadata/LinearSearchWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.SearchWithMetadata",
                        "LinearSearchWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 LinearSearchByIDWithMetadata RPC is the method to linear search similar vectors using a user-defined vector ID and to get metadata.<br>
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
        pub async fn linear_search_by_id_with_metadata(
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
                "/meta.v1.SearchWithMetadata/LinearSearchByIDWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.SearchWithMetadata",
                        "LinearSearchByIDWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamLinearSearchWithMetadata RPC is the method to linear search vectors and to get metadata with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        pub async fn stream_linear_search_with_metadata(
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
                "/meta.v1.SearchWithMetadata/StreamLinearSearchWithMetadata",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.SearchWithMetadata",
                        "StreamLinearSearchWithMetadata",
                    ),
                );
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
   StreamLinearSearchByIDWithMetadata RPC is the method to linear search vectors and to get metadata with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        pub async fn stream_linear_search_by_id_with_metadata(
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
                "/meta.v1.SearchWithMetadata/StreamLinearSearchByIDWithMetadata",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.SearchWithMetadata",
                        "StreamLinearSearchByIDWithMetadata",
                    ),
                );
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiLinearSearchWithMetadata RPC is the method to linear search vectors and to get metadata with multiple vectors in **1** request.

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
        pub async fn multi_linear_search_with_metadata(
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
                "/meta.v1.SearchWithMetadata/MultiLinearSearchWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.SearchWithMetadata",
                        "MultiLinearSearchWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 MultiLinearSearchByIDWithMetadata RPC is the method to linear search vectors and to get metadata with multiple IDs in **1** request.

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
        pub async fn multi_linear_search_by_id_with_metadata(
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
                "/meta.v1.SearchWithMetadata/MultiLinearSearchByIDWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.SearchWithMetadata",
                        "MultiLinearSearchByIDWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod search_with_metadata_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with SearchWithMetadataServer.
    #[async_trait]
    pub trait SearchWithMetadata: std::marker::Send + std::marker::Sync + 'static {
        /** Overview
 SearchWithMetadata RPC is the method to search vector(s) similar to the request vector and to get metadata(s).
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
        async fn search_with_metadata(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /** Overview
 SearchByIDWithMetadata RPC is the method to search similar vectors using a user-defined vector ID and to get metadata.<br>
 The vector with the same requested ID should be indexed into the `vald-lb-gateway` before searching.
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
        async fn search_by_id_with_metadata(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::IdRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamSearchWithMetadata method.
        type StreamSearchWithMetadataStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamSearchWithMetadata RPC is the method to search vectors and to get metadata with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        async fn stream_search_with_metadata(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::search::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamSearchWithMetadataStream>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamSearchByIDWithMetadata method.
        type StreamSearchByIDWithMetadataStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamSearchByIDWithMetadata RPC is the method to search vectors and to get metadata with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        async fn stream_search_by_id_with_metadata(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::search::IdRequest>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamSearchByIDWithMetadataStream>,
            tonic::Status,
        >;
        /** Overview
 MultiSearchWithMetadata RPC is the method to search vectors and to get metadata with multiple vectors in **1** request.

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
        async fn multi_search_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::search::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        >;
        /** Overview
 MultiSearchByIDWithMetadata RPC is the method to search vectors and to get metadata with multiple IDs in **1** request.

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
        async fn multi_search_by_id_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::search::MultiIdRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        >;
        /** Overview
 LinearSearchWithMetadata RPC is the method to linear search vector(s) similar to the request vector and to get metadata.
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
        async fn linear_search_with_metadata(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /** Overview
 LinearSearchByIDWithMetadata RPC is the method to linear search similar vectors using a user-defined vector ID and to get metadata.<br>
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
        async fn linear_search_by_id_with_metadata(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::IdRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamLinearSearchWithMetadata method.
        type StreamLinearSearchWithMetadataStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamLinearSearchWithMetadata RPC is the method to linear search vectors and to get metadata with multi queries(vectors) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        async fn stream_linear_search_with_metadata(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::search::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamLinearSearchWithMetadataStream>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamLinearSearchByIDWithMetadata method.
        type StreamLinearSearchByIDWithMetadataStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
   StreamLinearSearchByIDWithMetadata RPC is the method to linear search vectors and to get metadata with multi queries(IDs) using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        async fn stream_linear_search_by_id_with_metadata(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::search::IdRequest>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamLinearSearchByIDWithMetadataStream>,
            tonic::Status,
        >;
        /** Overview
 MultiLinearSearchWithMetadata RPC is the method to linear search vectors and to get metadata with multiple vectors in **1** request.

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
        async fn multi_linear_search_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::search::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        >;
        /** Overview
 MultiLinearSearchByIDWithMetadata RPC is the method to linear search vectors and to get metadata with multiple IDs in **1** request.

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
        async fn multi_linear_search_by_id_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::search::MultiIdRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Responses>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct SearchWithMetadataServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> SearchWithMetadataServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for SearchWithMetadataServer<T>
    where
        T: SearchWithMetadata,
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
                "/meta.v1.SearchWithMetadata/SearchWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct SearchWithMetadataSvc<T: SearchWithMetadata>(pub Arc<T>);
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::Request,
                    > for SearchWithMetadataSvc<T> {
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
                                <T as SearchWithMetadata>::search_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = SearchWithMetadataSvc(inner);
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
                "/meta.v1.SearchWithMetadata/SearchByIDWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct SearchByIDWithMetadataSvc<T: SearchWithMetadata>(pub Arc<T>);
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::IdRequest,
                    > for SearchByIDWithMetadataSvc<T> {
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
                                <T as SearchWithMetadata>::search_by_id_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = SearchByIDWithMetadataSvc(inner);
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
                "/meta.v1.SearchWithMetadata/StreamSearchWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct StreamSearchWithMetadataSvc<T: SearchWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::search::Request,
                    > for StreamSearchWithMetadataSvc<T> {
                        type Response = super::super::super::payload::v1::search::StreamResponse;
                        type ResponseStream = T::StreamSearchWithMetadataStream;
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
                                <T as SearchWithMetadata>::stream_search_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = StreamSearchWithMetadataSvc(inner);
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
                "/meta.v1.SearchWithMetadata/StreamSearchByIDWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct StreamSearchByIDWithMetadataSvc<T: SearchWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::search::IdRequest,
                    > for StreamSearchByIDWithMetadataSvc<T> {
                        type Response = super::super::super::payload::v1::search::StreamResponse;
                        type ResponseStream = T::StreamSearchByIDWithMetadataStream;
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
                                <T as SearchWithMetadata>::stream_search_by_id_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = StreamSearchByIDWithMetadataSvc(inner);
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
                "/meta.v1.SearchWithMetadata/MultiSearchWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct MultiSearchWithMetadataSvc<T: SearchWithMetadata>(pub Arc<T>);
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::MultiRequest,
                    > for MultiSearchWithMetadataSvc<T> {
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
                                <T as SearchWithMetadata>::multi_search_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = MultiSearchWithMetadataSvc(inner);
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
                "/meta.v1.SearchWithMetadata/MultiSearchByIDWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct MultiSearchByIDWithMetadataSvc<T: SearchWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::MultiIdRequest,
                    > for MultiSearchByIDWithMetadataSvc<T> {
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
                                <T as SearchWithMetadata>::multi_search_by_id_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = MultiSearchByIDWithMetadataSvc(inner);
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
                "/meta.v1.SearchWithMetadata/LinearSearchWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct LinearSearchWithMetadataSvc<T: SearchWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::Request,
                    > for LinearSearchWithMetadataSvc<T> {
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
                                <T as SearchWithMetadata>::linear_search_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = LinearSearchWithMetadataSvc(inner);
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
                "/meta.v1.SearchWithMetadata/LinearSearchByIDWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct LinearSearchByIDWithMetadataSvc<T: SearchWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::IdRequest,
                    > for LinearSearchByIDWithMetadataSvc<T> {
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
                                <T as SearchWithMetadata>::linear_search_by_id_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = LinearSearchByIDWithMetadataSvc(inner);
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
                "/meta.v1.SearchWithMetadata/StreamLinearSearchWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct StreamLinearSearchWithMetadataSvc<T: SearchWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::search::Request,
                    > for StreamLinearSearchWithMetadataSvc<T> {
                        type Response = super::super::super::payload::v1::search::StreamResponse;
                        type ResponseStream = T::StreamLinearSearchWithMetadataStream;
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
                                <T as SearchWithMetadata>::stream_linear_search_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = StreamLinearSearchWithMetadataSvc(inner);
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
                "/meta.v1.SearchWithMetadata/StreamLinearSearchByIDWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct StreamLinearSearchByIDWithMetadataSvc<T: SearchWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::search::IdRequest,
                    > for StreamLinearSearchByIDWithMetadataSvc<T> {
                        type Response = super::super::super::payload::v1::search::StreamResponse;
                        type ResponseStream = T::StreamLinearSearchByIDWithMetadataStream;
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
                                <T as SearchWithMetadata>::stream_linear_search_by_id_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = StreamLinearSearchByIDWithMetadataSvc(inner);
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
                "/meta.v1.SearchWithMetadata/MultiLinearSearchWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct MultiLinearSearchWithMetadataSvc<T: SearchWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::MultiRequest,
                    > for MultiLinearSearchWithMetadataSvc<T> {
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
                                <T as SearchWithMetadata>::multi_linear_search_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = MultiLinearSearchWithMetadataSvc(inner);
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
                "/meta.v1.SearchWithMetadata/MultiLinearSearchByIDWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct MultiLinearSearchByIDWithMetadataSvc<T: SearchWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: SearchWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::search::MultiIdRequest,
                    > for MultiLinearSearchByIDWithMetadataSvc<T> {
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
                                <T as SearchWithMetadata>::multi_linear_search_by_id_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = MultiLinearSearchByIDWithMetadataSvc(inner);
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
    impl<T> Clone for SearchWithMetadataServer<T> {
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
    pub const SERVICE_NAME: &str = "meta.v1.SearchWithMetadata";
    impl<T> tonic::server::NamedService for SearchWithMetadataServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod insert_with_metadata_client {
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
 InsertWithMetadata Service is responsible for inserting new vectors into the `vald-agent` and set metadata.
*/
    #[derive(Debug, Clone)]
    pub struct InsertWithMetadataClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl InsertWithMetadataClient<tonic::transport::Channel> {
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
    impl<T> InsertWithMetadataClient<T>
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
        ) -> InsertWithMetadataClient<InterceptedService<T, F>>
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
            InsertWithMetadataClient::new(InterceptedService::new(inner, interceptor))
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
 InsertWithMetadata RPC is the method to add a new single vector and metadata.
 ---
 Status Code
 | 0    | OK                |
 | 1    | CANCELLED         |
 | 3    | INVALID_ARGUMENT  |
 | 4    | DEADLINE_EXCEEDED |
 | 5    | NOT_FOUND         |
 | 13   | INTERNAL          |
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
        pub async fn insert_with_metadata(
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
            let path = http::uri::PathAndQuery::from_static(
                "/meta.v1.InsertWithMetadata/InsertWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("meta.v1.InsertWithMetadata", "InsertWithMetadata"),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamInsertWithMetadata RPC is the method to add new multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        pub async fn stream_insert_with_metadata(
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
                "/meta.v1.InsertWithMetadata/StreamInsertWithMetadata",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.InsertWithMetadata",
                        "StreamInsertWithMetadata",
                    ),
                );
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiInsertWithMetadata RPC is the method to add multiple new vectors and metadata in **1** request.

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
        pub async fn multi_insert_with_metadata(
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
                "/meta.v1.InsertWithMetadata/MultiInsertWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.InsertWithMetadata",
                        "MultiInsertWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod insert_with_metadata_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with InsertWithMetadataServer.
    #[async_trait]
    pub trait InsertWithMetadata: std::marker::Send + std::marker::Sync + 'static {
        /** Overview
 InsertWithMetadata RPC is the method to add a new single vector and metadata.
 ---
 Status Code
 | 0    | OK                |
 | 1    | CANCELLED         |
 | 3    | INVALID_ARGUMENT  |
 | 4    | DEADLINE_EXCEEDED |
 | 5    | NOT_FOUND         |
 | 13   | INTERNAL          |
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
        async fn insert_with_metadata(
            &self,
            request: tonic::Request<super::super::super::payload::v1::insert::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamInsertWithMetadata method.
        type StreamInsertWithMetadataStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamInsertWithMetadata RPC is the method to add new multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        async fn stream_insert_with_metadata(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::insert::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamInsertWithMetadataStream>,
            tonic::Status,
        >;
        /** Overview
 MultiInsertWithMetadata RPC is the method to add multiple new vectors and metadata in **1** request.

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
        async fn multi_insert_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::insert::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
    }
    /** Overview
 InsertWithMetadata Service is responsible for inserting new vectors into the `vald-agent` and set metadata.
*/
    #[derive(Debug)]
    pub struct InsertWithMetadataServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> InsertWithMetadataServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for InsertWithMetadataServer<T>
    where
        T: InsertWithMetadata,
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
                "/meta.v1.InsertWithMetadata/InsertWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct InsertWithMetadataSvc<T: InsertWithMetadata>(pub Arc<T>);
                    impl<
                        T: InsertWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::insert::Request,
                    > for InsertWithMetadataSvc<T> {
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
                                <T as InsertWithMetadata>::insert_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = InsertWithMetadataSvc(inner);
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
                "/meta.v1.InsertWithMetadata/StreamInsertWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct StreamInsertWithMetadataSvc<T: InsertWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: InsertWithMetadata,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::insert::Request,
                    > for StreamInsertWithMetadataSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamLocation;
                        type ResponseStream = T::StreamInsertWithMetadataStream;
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
                                <T as InsertWithMetadata>::stream_insert_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = StreamInsertWithMetadataSvc(inner);
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
                "/meta.v1.InsertWithMetadata/MultiInsertWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct MultiInsertWithMetadataSvc<T: InsertWithMetadata>(pub Arc<T>);
                    impl<
                        T: InsertWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::insert::MultiRequest,
                    > for MultiInsertWithMetadataSvc<T> {
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
                                <T as InsertWithMetadata>::multi_insert_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = MultiInsertWithMetadataSvc(inner);
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
    impl<T> Clone for InsertWithMetadataServer<T> {
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
    pub const SERVICE_NAME: &str = "meta.v1.InsertWithMetadata";
    impl<T> tonic::server::NamedService for InsertWithMetadataServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod object_with_metadata_client {
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
 ObjectWithMetadata Service is responsible for getting inserted vectors and metadata, and checking whether vectors are inserted into the `vald-agent`.
*/
    #[derive(Debug, Clone)]
    pub struct ObjectWithMetadataClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl ObjectWithMetadataClient<tonic::transport::Channel> {
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
    impl<T> ObjectWithMetadataClient<T>
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
        ) -> ObjectWithMetadataClient<InterceptedService<T, F>>
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
            ObjectWithMetadataClient::new(InterceptedService::new(inner, interceptor))
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
 GetObjectWithMetadata RPC is the method to get the metadata of a vector inserted into the `vald-agent` and metadata.
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
        pub async fn get_object_with_metadata(
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
            let path = http::uri::PathAndQuery::from_static(
                "/meta.v1.ObjectWithMetadata/GetObjectWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.ObjectWithMetadata",
                        "GetObjectWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamGetObjectWithMetadata RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        pub async fn stream_get_object_with_metadata(
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
                "/meta.v1.ObjectWithMetadata/StreamGetObjectWithMetadata",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.ObjectWithMetadata",
                        "StreamGetObjectWithMetadata",
                    ),
                );
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
        pub async fn stream_list_object_with_metadata(
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
                "/meta.v1.ObjectWithMetadata/StreamListObjectWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.ObjectWithMetadata",
                        "StreamListObjectWithMetadata",
                    ),
                );
            self.inner.server_streaming(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod object_with_metadata_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with ObjectWithMetadataServer.
    #[async_trait]
    pub trait ObjectWithMetadata: std::marker::Send + std::marker::Sync + 'static {
        /** Overview
 GetObjectWithMetadata RPC is the method to get the metadata of a vector inserted into the `vald-agent` and metadata.
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
        async fn get_object_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::object::VectorRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Vector>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamGetObjectWithMetadata method.
        type StreamGetObjectWithMetadataStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamVector,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamGetObjectWithMetadata RPC is the method to get the metadata of multiple existing vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        async fn stream_get_object_with_metadata(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::object::VectorRequest>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamGetObjectWithMetadataStream>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamListObjectWithMetadata method.
        type StreamListObjectWithMetadataStream: tonic::codegen::tokio_stream::Stream<
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
        async fn stream_list_object_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::object::list::Request,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamListObjectWithMetadataStream>,
            tonic::Status,
        >;
    }
    /** Overview
 ObjectWithMetadata Service is responsible for getting inserted vectors and metadata, and checking whether vectors are inserted into the `vald-agent`.
*/
    #[derive(Debug)]
    pub struct ObjectWithMetadataServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> ObjectWithMetadataServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for ObjectWithMetadataServer<T>
    where
        T: ObjectWithMetadata,
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
                "/meta.v1.ObjectWithMetadata/GetObjectWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct GetObjectWithMetadataSvc<T: ObjectWithMetadata>(pub Arc<T>);
                    impl<
                        T: ObjectWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::object::VectorRequest,
                    > for GetObjectWithMetadataSvc<T> {
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
                                <T as ObjectWithMetadata>::get_object_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = GetObjectWithMetadataSvc(inner);
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
                "/meta.v1.ObjectWithMetadata/StreamGetObjectWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct StreamGetObjectWithMetadataSvc<T: ObjectWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: ObjectWithMetadata,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::object::VectorRequest,
                    > for StreamGetObjectWithMetadataSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamVector;
                        type ResponseStream = T::StreamGetObjectWithMetadataStream;
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
                                <T as ObjectWithMetadata>::stream_get_object_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = StreamGetObjectWithMetadataSvc(inner);
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
                "/meta.v1.ObjectWithMetadata/StreamListObjectWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct StreamListObjectWithMetadataSvc<T: ObjectWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: ObjectWithMetadata,
                    > tonic::server::ServerStreamingService<
                        super::super::super::payload::v1::object::list::Request,
                    > for StreamListObjectWithMetadataSvc<T> {
                        type Response = super::super::super::payload::v1::object::list::Response;
                        type ResponseStream = T::StreamListObjectWithMetadataStream;
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
                                <T as ObjectWithMetadata>::stream_list_object_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = StreamListObjectWithMetadataSvc(inner);
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
    impl<T> Clone for ObjectWithMetadataServer<T> {
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
    pub const SERVICE_NAME: &str = "meta.v1.ObjectWithMetadata";
    impl<T> tonic::server::NamedService for ObjectWithMetadataServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod remove_with_metadata_client {
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
 RemoveWithMetadata Service is responsible for removing vectors indexed in the `vald-agent`.
*/
    #[derive(Debug, Clone)]
    pub struct RemoveWithMetadataClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl RemoveWithMetadataClient<tonic::transport::Channel> {
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
    impl<T> RemoveWithMetadataClient<T>
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
        ) -> RemoveWithMetadataClient<InterceptedService<T, F>>
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
            RemoveWithMetadataClient::new(InterceptedService::new(inner, interceptor))
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
 RemoveWithMetadata RPC is the method to remove a single vector and metadata.
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
        pub async fn remove_with_metadata(
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
            let path = http::uri::PathAndQuery::from_static(
                "/meta.v1.RemoveWithMetadata/RemoveWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("meta.v1.RemoveWithMetadata", "RemoveWithMetadata"),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 RemoveByTimestampWithMetadata RPC is the method to remove vectors and metadata based on timestamp.

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
        pub async fn remove_by_timestamp_with_metadata(
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
                "/meta.v1.RemoveWithMetadata/RemoveByTimestampWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.RemoveWithMetadata",
                        "RemoveByTimestampWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 A method to remove multiple with metadata indexed vectors and metadata by bidirectional streaming.

 StreamRemoveWithMetadata RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        pub async fn stream_remove_with_metadata(
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
                "/meta.v1.RemoveWithMetadata/StreamRemoveWithMetadata",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.RemoveWithMetadata",
                        "StreamRemoveWithMetadata",
                    ),
                );
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiRemoveWithMetadata is the method to remove multiple vectors and metadata in **1** request.

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
        pub async fn multi_remove_with_metadata(
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
                "/meta.v1.RemoveWithMetadata/MultiRemoveWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.RemoveWithMetadata",
                        "MultiRemoveWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod remove_with_metadata_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with RemoveWithMetadataServer.
    #[async_trait]
    pub trait RemoveWithMetadata: std::marker::Send + std::marker::Sync + 'static {
        /** Overview
 RemoveWithMetadata RPC is the method to remove a single vector and metadata.
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
        async fn remove_with_metadata(
            &self,
            request: tonic::Request<super::super::super::payload::v1::remove::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /** Overview
 RemoveByTimestampWithMetadata RPC is the method to remove vectors and metadata based on timestamp.

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
        async fn remove_by_timestamp_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::remove::TimestampRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamRemoveWithMetadata method.
        type StreamRemoveWithMetadataStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 A method to remove multiple with metadata indexed vectors and metadata by bidirectional streaming.

 StreamRemoveWithMetadata RPC is the method to remove multiple vectors using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        async fn stream_remove_with_metadata(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::remove::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamRemoveWithMetadataStream>,
            tonic::Status,
        >;
        /** Overview
 MultiRemoveWithMetadata is the method to remove multiple vectors and metadata in **1** request.

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
        async fn multi_remove_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::remove::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
    }
    /** Overview
 RemoveWithMetadata Service is responsible for removing vectors indexed in the `vald-agent`.
*/
    #[derive(Debug)]
    pub struct RemoveWithMetadataServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> RemoveWithMetadataServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for RemoveWithMetadataServer<T>
    where
        T: RemoveWithMetadata,
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
                "/meta.v1.RemoveWithMetadata/RemoveWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct RemoveWithMetadataSvc<T: RemoveWithMetadata>(pub Arc<T>);
                    impl<
                        T: RemoveWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::remove::Request,
                    > for RemoveWithMetadataSvc<T> {
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
                                <T as RemoveWithMetadata>::remove_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = RemoveWithMetadataSvc(inner);
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
                "/meta.v1.RemoveWithMetadata/RemoveByTimestampWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct RemoveByTimestampWithMetadataSvc<T: RemoveWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: RemoveWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::remove::TimestampRequest,
                    > for RemoveByTimestampWithMetadataSvc<T> {
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
                                <T as RemoveWithMetadata>::remove_by_timestamp_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = RemoveByTimestampWithMetadataSvc(inner);
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
                "/meta.v1.RemoveWithMetadata/StreamRemoveWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct StreamRemoveWithMetadataSvc<T: RemoveWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: RemoveWithMetadata,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::remove::Request,
                    > for StreamRemoveWithMetadataSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamLocation;
                        type ResponseStream = T::StreamRemoveWithMetadataStream;
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
                                <T as RemoveWithMetadata>::stream_remove_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = StreamRemoveWithMetadataSvc(inner);
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
                "/meta.v1.RemoveWithMetadata/MultiRemoveWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct MultiRemoveWithMetadataSvc<T: RemoveWithMetadata>(pub Arc<T>);
                    impl<
                        T: RemoveWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::remove::MultiRequest,
                    > for MultiRemoveWithMetadataSvc<T> {
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
                                <T as RemoveWithMetadata>::multi_remove_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = MultiRemoveWithMetadataSvc(inner);
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
    impl<T> Clone for RemoveWithMetadataServer<T> {
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
    pub const SERVICE_NAME: &str = "meta.v1.RemoveWithMetadata";
    impl<T> tonic::server::NamedService for RemoveWithMetadataServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod update_with_metadata_client {
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
 UpdateWithMetadata Service updates to new vector from inserted vector in the `vald-agent` components and updates to new metadata from set metadata.
*/
    #[derive(Debug, Clone)]
    pub struct UpdateWithMetadataClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl UpdateWithMetadataClient<tonic::transport::Channel> {
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
    impl<T> UpdateWithMetadataClient<T>
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
        ) -> UpdateWithMetadataClient<InterceptedService<T, F>>
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
            UpdateWithMetadataClient::new(InterceptedService::new(inner, interceptor))
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
 UpdateWithMetadata RPC is the method to update a single vector.
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
        pub async fn update_with_metadata(
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
            let path = http::uri::PathAndQuery::from_static(
                "/meta.v1.UpdateWithMetadata/UpdateWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("meta.v1.UpdateWithMetadata", "UpdateWithMetadata"),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamUpdateWithMetadata RPC is the method to update multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        pub async fn stream_update_with_metadata(
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
                "/meta.v1.UpdateWithMetadata/StreamUpdateWithMetadata",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.UpdateWithMetadata",
                        "StreamUpdateWithMetadata",
                    ),
                );
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiUpdateWithMetadata is the method to update multiple vectors and metadata in **1** request.

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
        pub async fn multi_update_with_metadata(
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
                "/meta.v1.UpdateWithMetadata/MultiUpdateWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.UpdateWithMetadata",
                        "MultiUpdateWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 A method to update timestamp an indexed vector and metadata.
 ---
 Status Code
 TODO
 ---
 Troubleshooting
 TODO
*/
        pub async fn update_timestamp_with_metadata(
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
                "/meta.v1.UpdateWithMetadata/UpdateTimestampWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.UpdateWithMetadata",
                        "UpdateTimestampWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod update_with_metadata_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with UpdateWithMetadataServer.
    #[async_trait]
    pub trait UpdateWithMetadata: std::marker::Send + std::marker::Sync + 'static {
        /** Overview
 UpdateWithMetadata RPC is the method to update a single vector.
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
        async fn update_with_metadata(
            &self,
            request: tonic::Request<super::super::super::payload::v1::update::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamUpdateWithMetadata method.
        type StreamUpdateWithMetadataStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamUpdateWithMetadata RPC is the method to update multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        async fn stream_update_with_metadata(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::update::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamUpdateWithMetadataStream>,
            tonic::Status,
        >;
        /** Overview
 MultiUpdateWithMetadata is the method to update multiple vectors and metadata in **1** request.

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
        async fn multi_update_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::update::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
        /** Overview
 A method to update timestamp an indexed vector and metadata.
 ---
 Status Code
 TODO
 ---
 Troubleshooting
 TODO
*/
        async fn update_timestamp_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::update::TimestampRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
    }
    /** Overview
 UpdateWithMetadata Service updates to new vector from inserted vector in the `vald-agent` components and updates to new metadata from set metadata.
*/
    #[derive(Debug)]
    pub struct UpdateWithMetadataServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> UpdateWithMetadataServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for UpdateWithMetadataServer<T>
    where
        T: UpdateWithMetadata,
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
                "/meta.v1.UpdateWithMetadata/UpdateWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct UpdateWithMetadataSvc<T: UpdateWithMetadata>(pub Arc<T>);
                    impl<
                        T: UpdateWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::update::Request,
                    > for UpdateWithMetadataSvc<T> {
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
                                <T as UpdateWithMetadata>::update_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = UpdateWithMetadataSvc(inner);
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
                "/meta.v1.UpdateWithMetadata/StreamUpdateWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct StreamUpdateWithMetadataSvc<T: UpdateWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: UpdateWithMetadata,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::update::Request,
                    > for StreamUpdateWithMetadataSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamLocation;
                        type ResponseStream = T::StreamUpdateWithMetadataStream;
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
                                <T as UpdateWithMetadata>::stream_update_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = StreamUpdateWithMetadataSvc(inner);
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
                "/meta.v1.UpdateWithMetadata/MultiUpdateWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct MultiUpdateWithMetadataSvc<T: UpdateWithMetadata>(pub Arc<T>);
                    impl<
                        T: UpdateWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::update::MultiRequest,
                    > for MultiUpdateWithMetadataSvc<T> {
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
                                <T as UpdateWithMetadata>::multi_update_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = MultiUpdateWithMetadataSvc(inner);
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
                "/meta.v1.UpdateWithMetadata/UpdateTimestampWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct UpdateTimestampWithMetadataSvc<T: UpdateWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: UpdateWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::update::TimestampRequest,
                    > for UpdateTimestampWithMetadataSvc<T> {
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
                                <T as UpdateWithMetadata>::update_timestamp_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = UpdateTimestampWithMetadataSvc(inner);
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
    impl<T> Clone for UpdateWithMetadataServer<T> {
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
    pub const SERVICE_NAME: &str = "meta.v1.UpdateWithMetadata";
    impl<T> tonic::server::NamedService for UpdateWithMetadataServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
/// Generated client implementations.
pub mod upsert_with_metadata_client {
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
 UpsertWithMetadata Service is responsible for updating existing vectors and metadata in the `vald-agent` or inserting new vectors into the `vald-agent` and metadata if the vector does not exist.
*/
    #[derive(Debug, Clone)]
    pub struct UpsertWithMetadataClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl UpsertWithMetadataClient<tonic::transport::Channel> {
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
    impl<T> UpsertWithMetadataClient<T>
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
        ) -> UpsertWithMetadataClient<InterceptedService<T, F>>
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
            UpsertWithMetadataClient::new(InterceptedService::new(inner, interceptor))
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
 UpsertWithMetadata RPC is the method to update the inserted vector and metadata to a new single vector and metadata or add a new single vector and metadata if not inserted before.
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
        pub async fn upsert_with_metadata(
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
            let path = http::uri::PathAndQuery::from_static(
                "/meta.v1.UpsertWithMetadata/UpsertWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("meta.v1.UpsertWithMetadata", "UpsertWithMetadata"),
                );
            self.inner.unary(req, path, codec).await
        }
        /** Overview
 StreamUpsertWithMetadata RPC is the method to update multiple existing vectors and metadata or add new multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        pub async fn stream_upsert_with_metadata(
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
                "/meta.v1.UpsertWithMetadata/StreamUpsertWithMetadata",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.UpsertWithMetadata",
                        "StreamUpsertWithMetadata",
                    ),
                );
            self.inner.streaming(req, path, codec).await
        }
        /** Overview
 MultiUpsertWithMetadata is the method to update existing multiple vectors and metadata and add new multiple vectors and metadata in **1** request.

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
        pub async fn multi_upsert_with_metadata(
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
                "/meta.v1.UpsertWithMetadata/MultiUpsertWithMetadata",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "meta.v1.UpsertWithMetadata",
                        "MultiUpsertWithMetadata",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod upsert_with_metadata_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with UpsertWithMetadataServer.
    #[async_trait]
    pub trait UpsertWithMetadata: std::marker::Send + std::marker::Sync + 'static {
        /** Overview
 UpsertWithMetadata RPC is the method to update the inserted vector and metadata to a new single vector and metadata or add a new single vector and metadata if not inserted before.
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
        async fn upsert_with_metadata(
            &self,
            request: tonic::Request<super::super::super::payload::v1::upsert::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamUpsertWithMetadata method.
        type StreamUpsertWithMetadataStream: tonic::codegen::tokio_stream::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + std::marker::Send
            + 'static;
        /** Overview
 StreamUpsertWithMetadata RPC is the method to update multiple existing vectors and metadata or add new multiple vectors and metadata using the [bidirectional streaming RPC](https://grpc.io/docs/what-is-grpc/core-concepts/#bidirectional-streaming-rpc).<br>
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
        async fn stream_upsert_with_metadata(
            &self,
            request: tonic::Request<
                tonic::Streaming<super::super::super::payload::v1::upsert::Request>,
            >,
        ) -> std::result::Result<
            tonic::Response<Self::StreamUpsertWithMetadataStream>,
            tonic::Status,
        >;
        /** Overview
 MultiUpsertWithMetadata is the method to update existing multiple vectors and metadata and add new multiple vectors and metadata in **1** request.

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
        async fn multi_upsert_with_metadata(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::upsert::MultiRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Locations>,
            tonic::Status,
        >;
    }
    /** Overview
 UpsertWithMetadata Service is responsible for updating existing vectors and metadata in the `vald-agent` or inserting new vectors into the `vald-agent` and metadata if the vector does not exist.
*/
    #[derive(Debug)]
    pub struct UpsertWithMetadataServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> UpsertWithMetadataServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for UpsertWithMetadataServer<T>
    where
        T: UpsertWithMetadata,
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
                "/meta.v1.UpsertWithMetadata/UpsertWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct UpsertWithMetadataSvc<T: UpsertWithMetadata>(pub Arc<T>);
                    impl<
                        T: UpsertWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::upsert::Request,
                    > for UpsertWithMetadataSvc<T> {
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
                                <T as UpsertWithMetadata>::upsert_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = UpsertWithMetadataSvc(inner);
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
                "/meta.v1.UpsertWithMetadata/StreamUpsertWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct StreamUpsertWithMetadataSvc<T: UpsertWithMetadata>(
                        pub Arc<T>,
                    );
                    impl<
                        T: UpsertWithMetadata,
                    > tonic::server::StreamingService<
                        super::super::super::payload::v1::upsert::Request,
                    > for StreamUpsertWithMetadataSvc<T> {
                        type Response = super::super::super::payload::v1::object::StreamLocation;
                        type ResponseStream = T::StreamUpsertWithMetadataStream;
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
                                <T as UpsertWithMetadata>::stream_upsert_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = StreamUpsertWithMetadataSvc(inner);
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
                "/meta.v1.UpsertWithMetadata/MultiUpsertWithMetadata" => {
                    #[allow(non_camel_case_types)]
                    struct MultiUpsertWithMetadataSvc<T: UpsertWithMetadata>(pub Arc<T>);
                    impl<
                        T: UpsertWithMetadata,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::upsert::MultiRequest,
                    > for MultiUpsertWithMetadataSvc<T> {
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
                                <T as UpsertWithMetadata>::multi_upsert_with_metadata(
                                        &inner,
                                        request,
                                    )
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
                        let method = MultiUpsertWithMetadataSvc(inner);
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
    impl<T> Clone for UpsertWithMetadataServer<T> {
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
    pub const SERVICE_NAME: &str = "meta.v1.UpsertWithMetadata";
    impl<T> tonic::server::NamedService for UpsertWithMetadataServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
