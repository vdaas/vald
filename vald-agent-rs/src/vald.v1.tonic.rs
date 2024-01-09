// @generated
/// Generated server implementations.
pub mod filter_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with FilterServer.
    #[async_trait]
    pub trait Filter: Send + Sync + 'static {
        async fn search_object(
            &self,
            request: tonic::Request<
                super::super::super::payload::v1::search::ObjectRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /** A method to search multiple objects.
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
        type StreamSearchObjectStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to search object by bidirectional streaming.
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
        /** A method insert object.
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
        type StreamInsertObjectStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** Represent the streaming RPC to insert object by bidirectional streaming.
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
        /** A method to insert multiple objects.
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
        /** A method to update object.
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
        type StreamUpdateObjectStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to update object by bidirectional streaming.
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
        /** A method to update multiple objects.
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
        /** A method to upsert object.
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
        type StreamUpsertObjectStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to upsert object by bidirectional streaming.
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
        /** A method to upsert multiple objects.
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
    #[derive(Debug)]
    pub struct FilterServer<T: Filter> {
        inner: _Inner<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    struct _Inner<T>(Arc<T>);
    impl<T: Filter> FilterServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for FilterServer<T>
    where
        T: Filter,
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
                                (*inner).search_object(request).await
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
                        let method = SearchObjectSvc(inner);
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
                                (*inner).multi_search_object(request).await
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
                        let method = MultiSearchObjectSvc(inner);
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
                                (*inner).stream_search_object(request).await
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
                        let method = StreamSearchObjectSvc(inner);
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
                                (*inner).insert_object(request).await
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
                        let method = InsertObjectSvc(inner);
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
                                (*inner).stream_insert_object(request).await
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
                        let method = StreamInsertObjectSvc(inner);
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
                                (*inner).multi_insert_object(request).await
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
                        let method = MultiInsertObjectSvc(inner);
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
                                (*inner).update_object(request).await
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
                        let method = UpdateObjectSvc(inner);
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
                                (*inner).stream_update_object(request).await
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
                        let method = StreamUpdateObjectSvc(inner);
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
                                (*inner).multi_update_object(request).await
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
                        let method = MultiUpdateObjectSvc(inner);
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
                                (*inner).upsert_object(request).await
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
                        let method = UpsertObjectSvc(inner);
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
                                (*inner).stream_upsert_object(request).await
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
                        let method = StreamUpsertObjectSvc(inner);
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
                                (*inner).multi_upsert_object(request).await
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
                        let method = MultiUpsertObjectSvc(inner);
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
    impl<T: Filter> Clone for FilterServer<T> {
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
    impl<T: Filter> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(Arc::clone(&self.0))
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: Filter> tonic::server::NamedService for FilterServer<T> {
        const NAME: &'static str = "vald.v1.Filter";
    }
}
/// Generated server implementations.
pub mod insert_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with InsertServer.
    #[async_trait]
    pub trait Insert: Send + Sync + 'static {
        async fn insert(
            &self,
            request: tonic::Request<super::super::super::payload::v1::insert::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamInsert method.
        type StreamInsertStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to add new multiple vectors by bidirectional streaming.
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
        /** A method to add new multiple vectors in a single request.
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
    pub struct InsertServer<T: Insert> {
        inner: _Inner<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    struct _Inner<T>(Arc<T>);
    impl<T: Insert> InsertServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for InsertServer<T>
    where
        T: Insert,
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
                            let fut = async move { (*inner).insert(request).await };
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
                        let method = InsertSvc(inner);
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
                                (*inner).stream_insert(request).await
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
                        let method = StreamInsertSvc(inner);
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
                                (*inner).multi_insert(request).await
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
                        let method = MultiInsertSvc(inner);
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
    impl<T: Insert> Clone for InsertServer<T> {
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
    impl<T: Insert> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(Arc::clone(&self.0))
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: Insert> tonic::server::NamedService for InsertServer<T> {
        const NAME: &'static str = "vald.v1.Insert";
    }
}
/// Generated server implementations.
pub mod object_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with ObjectServer.
    #[async_trait]
    pub trait Object: Send + Sync + 'static {
        async fn exists(
            &self,
            request: tonic::Request<super::super::super::payload::v1::object::Id>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Id>,
            tonic::Status,
        >;
        /** A method to fetch a vector.
*/
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
        type StreamGetObjectStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamVector,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to fetch vectors by bidirectional streaming.
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
        type StreamListObjectStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::list::Response,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to get all the vectors with server streaming
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
    }
    #[derive(Debug)]
    pub struct ObjectServer<T: Object> {
        inner: _Inner<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    struct _Inner<T>(Arc<T>);
    impl<T: Object> ObjectServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for ObjectServer<T>
    where
        T: Object,
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
                            let fut = async move { (*inner).exists(request).await };
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
                        let method = ExistsSvc(inner);
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
                            let fut = async move { (*inner).get_object(request).await };
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
                        let method = GetObjectSvc(inner);
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
                                (*inner).stream_get_object(request).await
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
                        let method = StreamGetObjectSvc(inner);
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
                                (*inner).stream_list_object(request).await
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
                        let method = StreamListObjectSvc(inner);
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
                        let res = grpc.server_streaming(method, req).await;
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
    impl<T: Object> Clone for ObjectServer<T> {
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
    impl<T: Object> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(Arc::clone(&self.0))
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: Object> tonic::server::NamedService for ObjectServer<T> {
        const NAME: &'static str = "vald.v1.Object";
    }
}
/// Generated server implementations.
pub mod remove_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with RemoveServer.
    #[async_trait]
    pub trait Remove: Send + Sync + 'static {
        async fn remove(
            &self,
            request: tonic::Request<super::super::super::payload::v1::remove::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /** A method to remove an indexed vector based on timestamp.
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
        type StreamRemoveStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to remove multiple indexed vectors by bidirectional streaming.
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
        /** A method to remove multiple indexed vectors in a single request.
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
    pub struct RemoveServer<T: Remove> {
        inner: _Inner<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    struct _Inner<T>(Arc<T>);
    impl<T: Remove> RemoveServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for RemoveServer<T>
    where
        T: Remove,
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
                            let fut = async move { (*inner).remove(request).await };
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
                        let method = RemoveSvc(inner);
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
                                (*inner).remove_by_timestamp(request).await
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
                        let method = RemoveByTimestampSvc(inner);
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
                                (*inner).stream_remove(request).await
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
                        let method = StreamRemoveSvc(inner);
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
                                (*inner).multi_remove(request).await
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
                        let method = MultiRemoveSvc(inner);
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
    impl<T: Remove> Clone for RemoveServer<T> {
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
    impl<T: Remove> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(Arc::clone(&self.0))
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: Remove> tonic::server::NamedService for RemoveServer<T> {
        const NAME: &'static str = "vald.v1.Remove";
    }
}
/// Generated server implementations.
pub mod search_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with SearchServer.
    #[async_trait]
    pub trait Search: Send + Sync + 'static {
        async fn search(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /** A method to search indexed vectors by ID.
*/
        async fn search_by_id(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::IdRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamSearch method.
        type StreamSearchStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to search indexed vectors by multiple vectors.
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
        type StreamSearchByIDStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to search indexed vectors by multiple IDs.
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
        /** A method to search indexed vectors by multiple vectors in a single request.
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
        /** A method to search indexed vectors by multiple IDs in a single request.
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
        /** A method to linear search indexed vectors by a raw vector.
*/
        async fn linear_search(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /** A method to linear search indexed vectors by ID.
*/
        async fn linear_search_by_id(
            &self,
            request: tonic::Request<super::super::super::payload::v1::search::IdRequest>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::search::Response>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamLinearSearch method.
        type StreamLinearSearchStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to linear search indexed vectors by multiple vectors.
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
        type StreamLinearSearchByIDStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::search::StreamResponse,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to linear search indexed vectors by multiple IDs.
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
        /** A method to linear search indexed vectors by multiple vectors in a single
 request.
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
        /** A method to linear search indexed vectors by multiple IDs in a single
 request.
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
    #[derive(Debug)]
    pub struct SearchServer<T: Search> {
        inner: _Inner<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    struct _Inner<T>(Arc<T>);
    impl<T: Search> SearchServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for SearchServer<T>
    where
        T: Search,
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
                            let fut = async move { (*inner).search(request).await };
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
                        let method = SearchSvc(inner);
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
                                (*inner).search_by_id(request).await
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
                        let method = SearchByIDSvc(inner);
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
                                (*inner).stream_search(request).await
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
                        let method = StreamSearchSvc(inner);
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
                                (*inner).stream_search_by_id(request).await
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
                        let method = StreamSearchByIDSvc(inner);
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
                                (*inner).multi_search(request).await
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
                        let method = MultiSearchSvc(inner);
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
                                (*inner).multi_search_by_id(request).await
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
                        let method = MultiSearchByIDSvc(inner);
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
                                (*inner).linear_search(request).await
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
                        let method = LinearSearchSvc(inner);
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
                                (*inner).linear_search_by_id(request).await
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
                        let method = LinearSearchByIDSvc(inner);
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
                                (*inner).stream_linear_search(request).await
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
                        let method = StreamLinearSearchSvc(inner);
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
                                (*inner).stream_linear_search_by_id(request).await
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
                        let method = StreamLinearSearchByIDSvc(inner);
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
                                (*inner).multi_linear_search(request).await
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
                        let method = MultiLinearSearchSvc(inner);
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
                                (*inner).multi_linear_search_by_id(request).await
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
                        let method = MultiLinearSearchByIDSvc(inner);
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
    impl<T: Search> Clone for SearchServer<T> {
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
    impl<T: Search> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(Arc::clone(&self.0))
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: Search> tonic::server::NamedService for SearchServer<T> {
        const NAME: &'static str = "vald.v1.Search";
    }
}
/// Generated server implementations.
pub mod update_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with UpdateServer.
    #[async_trait]
    pub trait Update: Send + Sync + 'static {
        async fn update(
            &self,
            request: tonic::Request<super::super::super::payload::v1::update::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamUpdate method.
        type StreamUpdateStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to update multiple indexed vectors by bidirectional streaming.
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
        /** A method to update multiple indexed vectors in a single request.
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
    }
    #[derive(Debug)]
    pub struct UpdateServer<T: Update> {
        inner: _Inner<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    struct _Inner<T>(Arc<T>);
    impl<T: Update> UpdateServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for UpdateServer<T>
    where
        T: Update,
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
                            let fut = async move { (*inner).update(request).await };
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
                        let method = UpdateSvc(inner);
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
                                (*inner).stream_update(request).await
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
                        let method = StreamUpdateSvc(inner);
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
                                (*inner).multi_update(request).await
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
                        let method = MultiUpdateSvc(inner);
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
    impl<T: Update> Clone for UpdateServer<T> {
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
    impl<T: Update> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(Arc::clone(&self.0))
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: Update> tonic::server::NamedService for UpdateServer<T> {
        const NAME: &'static str = "vald.v1.Update";
    }
}
/// Generated server implementations.
pub mod upsert_server {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    /// Generated trait containing gRPC methods that should be implemented for use with UpsertServer.
    #[async_trait]
    pub trait Upsert: Send + Sync + 'static {
        async fn upsert(
            &self,
            request: tonic::Request<super::super::super::payload::v1::upsert::Request>,
        ) -> std::result::Result<
            tonic::Response<super::super::super::payload::v1::object::Location>,
            tonic::Status,
        >;
        /// Server streaming response type for the StreamUpsert method.
        type StreamUpsertStream: futures_core::Stream<
                Item = std::result::Result<
                    super::super::super::payload::v1::object::StreamLocation,
                    tonic::Status,
                >,
            >
            + Send
            + 'static;
        /** A method to insert/update multiple vectors by bidirectional streaming.
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
        /** A method to insert/update multiple vectors in a single request.
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
    pub struct UpsertServer<T: Upsert> {
        inner: _Inner<T>,
        accept_compression_encodings: EnabledCompressionEncodings,
        send_compression_encodings: EnabledCompressionEncodings,
        max_decoding_message_size: Option<usize>,
        max_encoding_message_size: Option<usize>,
    }
    struct _Inner<T>(Arc<T>);
    impl<T: Upsert> UpsertServer<T> {
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
    impl<T, B> tonic::codegen::Service<http::Request<B>> for UpsertServer<T>
    where
        T: Upsert,
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
                            let fut = async move { (*inner).upsert(request).await };
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
                        let method = UpsertSvc(inner);
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
                                (*inner).stream_upsert(request).await
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
                        let method = StreamUpsertSvc(inner);
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
                                (*inner).multi_upsert(request).await
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
                        let method = MultiUpsertSvc(inner);
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
    impl<T: Upsert> Clone for UpsertServer<T> {
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
    impl<T: Upsert> Clone for _Inner<T> {
        fn clone(&self) -> Self {
            Self(Arc::clone(&self.0))
        }
    }
    impl<T: std::fmt::Debug> std::fmt::Debug for _Inner<T> {
        fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
            write!(f, "{:?}", self.0)
        }
    }
    impl<T: Upsert> tonic::server::NamedService for UpsertServer<T> {
        const NAME: &'static str = "vald.v1.Upsert";
    }
}
