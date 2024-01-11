// @generated
/// Generated client implementations.
pub mod agent_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    #[derive(Debug, Clone)]
    pub struct AgentClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl AgentClient<tonic::transport::Channel> {
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
    impl<T> AgentClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + Send,
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
        ) -> AgentClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
            >>::Error: Into<StdError> + Send + Sync,
        {
            AgentClient::new(InterceptedService::new(inner, interceptor))
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
        pub async fn create_index(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::control::CreateIndexRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::Empty>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/core.v1.Agent/CreateIndex",
            );
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("core.v1.Agent", "CreateIndex"));
            self.inner.unary(req, path, codec).await
        }
        pub async fn save_index(
            &mut self,
            request: impl tonic::IntoRequest<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::Empty>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/core.v1.Agent/SaveIndex");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("core.v1.Agent", "SaveIndex"));
            self.inner.unary(req, path, codec).await
        }
        /** Represent the creating and saving index RPC.
*/
        pub async fn create_and_save_index(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::control::CreateIndexRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::Empty>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/core.v1.Agent/CreateAndSaveIndex",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("core.v1.Agent", "CreateAndSaveIndex"));
            self.inner.unary(req, path, codec).await
        }
        /** Represent the RPC to get the agent index information.
*/
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/core.v1.Agent/IndexInfo");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("core.v1.Agent", "IndexInfo"));
            self.inner.unary(req, path, codec).await
        }
        /** Represent the RPC to get the vector metadata. This RPC is mainly used for index correction process
*/
        pub async fn get_timestamp(
            &mut self,
            request: impl tonic::IntoRequest<
                super::super::super::payload::v1::object::GetTimestampRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Timestamp>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/core.v1.Agent/GetTimestamp",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("core.v1.Agent", "GetTimestamp"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated server implementations.
pub mod agent_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with AgentServer.
    #[async_trait]
    pub trait Agent: Send + Sync + 'static {
        async fn create_index(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::control::CreateIndexRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::Empty>,
            tonic::Status,
        >;
        async fn save_index(
            &self,
            request: tonic::Request<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::Empty>,
            tonic::Status,
        >;
        /** Represent the creating and saving index RPC.
*/
        async fn create_and_save_index(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::control::CreateIndexRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::Empty>,
            tonic::Status,
        >;
        /** Represent the RPC to get the agent index information.
*/
        async fn index_info(
            &self,
            request: tonic::Request<super::super::super::payload::v1::Empty>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::info::index::Count>,
            tonic::Status,
        >;
        /** Represent the RPC to get the vector metadata. This RPC is mainly used for index correction process
*/
        async fn get_timestamp(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::object::GetTimestampRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Timestamp>,
            tonic::Status,
        >;
    }
    #[derive(Debug)]
    pub struct AgentServer<T: Agent> {
        inner: _Inner<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    struct _Inner<T>(Arc<T>);
    impl<T: Agent> AgentServer<T> {
        pub fn new(inner: T) -> Self {
            Self::from_arc(Arc::new(inner))
        }
        pub fn from_arc(inner: Arc<T>) -> Self {
            let inner = _Inner(inner);
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for AgentServer<T>
    where
        T: Agent,
        B: Body + Send + 'static,
        B::Error: Into<StdError> + Send + 'static,
    {
        type Response = http::Response<tonic::body::BoxBody>;
        type Error = std::convert::Infallible;
        type Future = BoxFuture<Self::Response, Self::Error>;
        fn poll_ready(
            &mut self,
            _cx: &mut Context<'_>,
        ) -> Poll<std::result::Result<(), Self::Error>> {
            Poll::Ready(Ok(()))
        }
        fn call(&mut self, req: http::Request<B>) -> Self::Future {
            let inner = self.inner.clone();
            match req.uri().path() {
                "/core.v1.Agent/CreateIndex" => {
                    #[allow(non_camel_case_types)]
                    struct CreateIndexSvc<T: Agent>(pub Arc<T>);
                    impl<
                        T: Agent,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::control::CreateIndexRequest,
                    > for CreateIndexSvc<T> {
                        type Response = super::super::super::payload::v1::Empty;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::control::CreateIndexRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                (*inner).create_index(request).await
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
                        let inner = inner.0;
                        let method = CreateIndexSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
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
                "/core.v1.Agent/SaveIndex" => {
                    #[allow(non_camel_case_types)]
                    struct SaveIndexSvc<T: Agent>(pub Arc<T>);
                    impl<
                        T: Agent,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::Empty,
                    > for SaveIndexSvc<T> {
                        type Response = super::super::super::payload::v1::Empty;
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
                            let fut = async move { (*inner).save_index(request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let inner = inner.0;
                        let method = SaveIndexSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
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
                "/core.v1.Agent/CreateAndSaveIndex" => {
                    #[allow(non_camel_case_types)]
                    struct CreateAndSaveIndexSvc<T: Agent>(pub Arc<T>);
                    impl<
                        T: Agent,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::control::CreateIndexRequest,
                    > for CreateAndSaveIndexSvc<T> {
                        type Response = super::super::super::payload::v1::Empty;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::control::CreateIndexRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                (*inner).create_and_save_index(request).await
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
                        let inner = inner.0;
                        let method = CreateAndSaveIndexSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
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
                "/core.v1.Agent/IndexInfo" => {
                    #[allow(non_camel_case_types)]
                    struct IndexInfoSvc<T: Agent>(pub Arc<T>);
                    impl<
                        T: Agent,
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
                            let fut = async move { (*inner).index_info(request).await };
                            Box::pin(fut)
                        }
                    }
                    let accept_compression_encodings = self.accept_compression_encodings;
                    let send_compression_encodings = self.send_compression_encodings;
                    let max_decoding_message_size = self.max_decoding_message_size;
                    let max_encoding_message_size = self.max_encoding_message_size;
                    let inner = self.inner.clone();
                    let fut = async move {
                        let inner = inner.0;
                        let method = IndexInfoSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
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
                "/core.v1.Agent/GetTimestamp" => {
                    #[allow(non_camel_case_types)]
                    struct GetTimestampSvc<T: Agent>(pub Arc<T>);
                    impl<
                        T: Agent,
                    > tonic::server::UnaryService<
                        super::super::super::payload::v1::object::GetTimestampRequest,
                    > for GetTimestampSvc<T> {
                        type Response = super::super::super::payload::v1::object::Timestamp;
                        type Future = BoxFuture<
                            tonic::Response<Self::Response>,
                            tonic::Status,
                        >;
                        fn call(
                            &mut self,
                            request: tonic::Request<
                                super::super::super::payload::v1::object::GetTimestampRequest,
                            >,
                        ) -> Self::Future {
                            let inner = Arc::clone(&self.0);
                            let fut = async move {
                                (*inner).get_timestamp(request).await
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
                        let inner = inner.0;
                        let method = GetTimestampSvc(inner);
                        let codec = tonic::codec::ProstCodec::default();
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
                        Ok(
                            http::Response::builder()
                                .status(200)
                                .header("grpc-status", "12")
                                .header("content-type", "application/grpc")
                                .body(empty_body())
                                .unwrap(),
                        )
                    })
                }
            }
        }
    }
    impl<T: Agent> Clone for AgentServer<T> {
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
    impl<T: Agent> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(Arc::clone(&self.0))
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: Agent> tonic::server::NamedService for AgentServer<T> {
        const NAME: &'static str = "core.v1.Agent";
    }
}
