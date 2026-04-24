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
pub mod tikv_client {
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
    pub struct TikvClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl TikvClient<tonic::transport::Channel> {
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
    impl<T> TikvClient<T>
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
        pub fn with_interceptor<F>(inner: T, interceptor: F) -> TikvClient<InterceptedService<T, F>>
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
            TikvClient::new(InterceptedService::new(inner, interceptor))
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
        ///
        pub async fn raw_get(
            &mut self,
            request: impl tonic::IntoRequest<super::super::tikv::RawGetRequest>,
        ) -> std::result::Result<tonic::Response<super::super::tikv::RawGetResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::unknown(format!("Service was not ready: {}", e.into()))
            })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/tikvpb.Tikv/RawGet");
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("tikvpb.Tikv", "RawGet"));
            self.inner.unary(req, path, codec).await
        }
        ///
        pub async fn raw_batch_get(
            &mut self,
            request: impl tonic::IntoRequest<super::super::tikv::RawBatchGetRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::tikv::RawBatchGetResponse>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::unknown(format!("Service was not ready: {}", e.into()))
            })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/tikvpb.Tikv/RawBatchGet");
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("tikvpb.Tikv", "RawBatchGet"));
            self.inner.unary(req, path, codec).await
        }
        ///
        pub async fn raw_put(
            &mut self,
            request: impl tonic::IntoRequest<super::super::tikv::RawPutRequest>,
        ) -> std::result::Result<tonic::Response<super::super::tikv::RawPutResponse>, tonic::Status>
        {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::unknown(format!("Service was not ready: {}", e.into()))
            })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/tikvpb.Tikv/RawPut");
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("tikvpb.Tikv", "RawPut"));
            self.inner.unary(req, path, codec).await
        }
        ///
        pub async fn raw_batch_put(
            &mut self,
            request: impl tonic::IntoRequest<super::super::tikv::RawBatchPutRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::tikv::RawBatchPutResponse>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::unknown(format!("Service was not ready: {}", e.into()))
            })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/tikvpb.Tikv/RawBatchPut");
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("tikvpb.Tikv", "RawBatchPut"));
            self.inner.unary(req, path, codec).await
        }
        ///
        pub async fn raw_delete(
            &mut self,
            request: impl tonic::IntoRequest<super::super::tikv::RawDeleteRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::tikv::RawDeleteResponse>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::unknown(format!("Service was not ready: {}", e.into()))
            })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/tikvpb.Tikv/RawDelete");
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("tikvpb.Tikv", "RawDelete"));
            self.inner.unary(req, path, codec).await
        }
        ///
        pub async fn raw_batch_delete(
            &mut self,
            request: impl tonic::IntoRequest<super::super::tikv::RawBatchDeleteRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::tikv::RawBatchDeleteResponse>,
            tonic::Status,
        > {
            self.inner.ready().await.map_err(|e| {
                tonic::Status::unknown(format!("Service was not ready: {}", e.into()))
            })?;
            let codec = tonic_prost::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/tikvpb.Tikv/RawBatchDelete");
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("tikvpb.Tikv", "RawBatchDelete"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod tikv_server {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value
    )]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with TikvServer.
    #[async_trait]
    pub trait Tikv: std::marker::Send + std::marker::Sync + 'static {
        ///
        async fn raw_get(
            &self,
            request: tonic::Request<super::super::tikv::RawGetRequest>,
        ) -> std::result::Result<tonic::Response<super::super::tikv::RawGetResponse>, tonic::Status>;
        ///
        async fn raw_batch_get(
            &self,
            request: tonic::Request<super::super::tikv::RawBatchGetRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::tikv::RawBatchGetResponse>,
            tonic::Status,
        >;
        ///
        async fn raw_put(
            &self,
            request: tonic::Request<super::super::tikv::RawPutRequest>,
        ) -> std::result::Result<tonic::Response<super::super::tikv::RawPutResponse>, tonic::Status>;
        ///
        async fn raw_batch_put(
            &self,
            request: tonic::Request<super::super::tikv::RawBatchPutRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::tikv::RawBatchPutResponse>,
            tonic::Status,
        >;
        ///
        async fn raw_delete(
            &self,
            request: tonic::Request<super::super::tikv::RawDeleteRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::tikv::RawDeleteResponse>,
            tonic::Status,
        >;
        ///
        async fn raw_batch_delete(
            &self,
            request: tonic::Request<super::super::tikv::RawBatchDeleteRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::tikv::RawBatchDeleteResponse>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct TikvServer<T> {
        inner: Arc<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    impl<T> TikvServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for TikvServer<T>
    where
        T: Tikv,
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
                "/tikvpb.Tikv/RawGet" => {
                    #[allow(non_camel_case_types)]
                    struct RawGetSvc<T: Tikv>(pub Arc<T>);
                    impl<T: Tikv> tonic::server::UnaryService<super::super::tikv::RawGetRequest> for RawGetSvc<T> {
                        type Response = super::super::tikv::RawGetResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::super::tikv::RawGetRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move { <T as Tikv>::raw_get(&inner, request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = RawGetSvc(inner);
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
                "/tikvpb.Tikv/RawBatchGet" => {
                    #[allow(non_camel_case_types)]
                    struct RawBatchGetSvc<T: Tikv>(pub Arc<T>);
                    impl<T: Tikv>
                        tonic::server::UnaryService<super::super::tikv::RawBatchGetRequest>
                        for RawBatchGetSvc<T>
                    {
                        type Response = super::super::tikv::RawBatchGetResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::super::tikv::RawBatchGetRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut =
                                async move { <T as Tikv>::raw_batch_get(&inner, request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = RawBatchGetSvc(inner);
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
                "/tikvpb.Tikv/RawPut" => {
                    #[allow(non_camel_case_types)]
                    struct RawPutSvc<T: Tikv>(pub Arc<T>);
                    impl<T: Tikv> tonic::server::UnaryService<super::super::tikv::RawPutRequest> for RawPutSvc<T> {
                        type Response = super::super::tikv::RawPutResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::super::tikv::RawPutRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move { <T as Tikv>::raw_put(&inner, request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = RawPutSvc(inner);
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
                "/tikvpb.Tikv/RawBatchPut" => {
                    #[allow(non_camel_case_types)]
                    struct RawBatchPutSvc<T: Tikv>(pub Arc<T>);
                    impl<T: Tikv>
                        tonic::server::UnaryService<super::super::tikv::RawBatchPutRequest>
                        for RawBatchPutSvc<T>
                    {
                        type Response = super::super::tikv::RawBatchPutResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::super::tikv::RawBatchPutRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut =
                                async move { <T as Tikv>::raw_batch_put(&inner, request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = RawBatchPutSvc(inner);
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
                "/tikvpb.Tikv/RawDelete" => {
                    #[allow(non_camel_case_types)]
                    struct RawDeleteSvc<T: Tikv>(pub Arc<T>);
                    impl<T: Tikv> tonic::server::UnaryService<super::super::tikv::RawDeleteRequest>
                        for RawDeleteSvc<T>
                    {
                        type Response = super::super::tikv::RawDeleteResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::super::tikv::RawDeleteRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move { <T as Tikv>::raw_delete(&inner, request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = RawDeleteSvc(inner);
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
                "/tikvpb.Tikv/RawBatchDelete" => {
                    #[allow(non_camel_case_types)]
                    struct RawBatchDeleteSvc<T: Tikv>(pub Arc<T>);
                    impl<T: Tikv>
                        tonic::server::UnaryService<super::super::tikv::RawBatchDeleteRequest>
                        for RawBatchDeleteSvc<T>
                    {
                        type Response = super::super::tikv::RawBatchDeleteResponse;
                        type Future = BoxFuture<tonic::Response<Self::Response>, tonic::Status>;
                        fn call(
                            &mut self,
                            request: tonic::Request<super::super::tikv::RawBatchDeleteRequest>,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut =
                                async move { <T as Tikv>::raw_batch_delete(&inner, request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let method = RawBatchDeleteSvc(inner);
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
    impl<T> Clone for TikvServer<T> {
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
    pub const SERVICE_NAME: &str = "tikvpb.Tikv";
    impl<T> tonic::server::NamedService for TikvServer<T> {
        const NAME: &'static str = SERVICE_NAME;
    }
}
