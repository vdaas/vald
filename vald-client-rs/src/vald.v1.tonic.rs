// @generated
/// Generated client implementations.
pub mod filter_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
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
        ) -> FilterClient<InterceptedService<T, F>>
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/SearchObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "SearchObject"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to search multiple objects.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/MultiSearchObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "MultiSearchObject"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to search object by bidirectional streaming.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/StreamSearchObject",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "StreamSearchObject"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method insert object.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/InsertObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "InsertObject"));
            self.inner.unary(req, path, codec).await
        }
        /** Represent the streaming RPC to insert object by bidirectional streaming.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/StreamInsertObject",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "StreamInsertObject"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to insert multiple objects.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/MultiInsertObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "MultiInsertObject"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to update object.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/UpdateObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "UpdateObject"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to update object by bidirectional streaming.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/StreamUpdateObject",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "StreamUpdateObject"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to update multiple objects.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/MultiUpdateObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "MultiUpdateObject"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to upsert object.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/UpsertObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "UpsertObject"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to upsert object by bidirectional streaming.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Filter/StreamUpsertObject",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Filter", "StreamUpsertObject"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to upsert multiple objects.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
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
/// Generated client implementations.
pub mod insert_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
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
        ) -> InsertClient<InterceptedService<T, F>>
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Insert/Insert");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Insert", "Insert"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to add new multiple vectors by bidirectional streaming.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Insert/StreamInsert",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Insert", "StreamInsert"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to add new multiple vectors in a single request.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
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
/// Generated client implementations.
pub mod object_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
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
        ) -> ObjectClient<InterceptedService<T, F>>
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Object/Exists");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Object", "Exists"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to fetch a vector.
*/
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Object/GetObject");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Object", "GetObject"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to fetch vectors by bidirectional streaming.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Object/StreamGetObject",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Object", "StreamGetObject"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to get all the vectors with server streaming
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Object/StreamListObject",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Object", "StreamListObject"));
            self.inner.server_streaming(req, path, codec).await
        }
    }
}
/// Generated client implementations.
pub mod remove_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
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
        ) -> RemoveClient<InterceptedService<T, F>>
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Remove/Remove");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Remove", "Remove"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to remove an indexed vector based on timestamp.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Remove/RemoveByTimestamp",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Remove", "RemoveByTimestamp"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to remove multiple indexed vectors by bidirectional streaming.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Remove/StreamRemove",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Remove", "StreamRemove"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to remove multiple indexed vectors in a single request.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
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
/// Generated client implementations.
pub mod search_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
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
        ) -> SearchClient<InterceptedService<T, F>>
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Search/Search");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Search", "Search"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to search indexed vectors by ID.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/SearchByID",
            );
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Search", "SearchByID"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to search indexed vectors by multiple vectors.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/StreamSearch",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "StreamSearch"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to search indexed vectors by multiple IDs.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/StreamSearchByID",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "StreamSearchByID"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to search indexed vectors by multiple vectors in a single request.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/MultiSearch",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "MultiSearch"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to search indexed vectors by multiple IDs in a single request.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/MultiSearchByID",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "MultiSearchByID"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to linear search indexed vectors by a raw vector.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/LinearSearch",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "LinearSearch"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to linear search indexed vectors by ID.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/LinearSearchByID",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "LinearSearchByID"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to linear search indexed vectors by multiple vectors.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/StreamLinearSearch",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "StreamLinearSearch"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to linear search indexed vectors by multiple IDs.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/StreamLinearSearchByID",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "StreamLinearSearchByID"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to linear search indexed vectors by multiple vectors in a single
 request.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Search/MultiLinearSearch",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Search", "MultiLinearSearch"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to linear search indexed vectors by multiple IDs in a single
 request.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
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
/// Generated client implementations.
pub mod update_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
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
        ) -> UpdateClient<InterceptedService<T, F>>
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Update/Update");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Update", "Update"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to update multiple indexed vectors by bidirectional streaming.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Update/StreamUpdate",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Update", "StreamUpdate"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to update multiple indexed vectors in a single request.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Update/MultiUpdate",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Update", "MultiUpdate"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// Generated client implementations.
pub mod upsert_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
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
        ) -> UpsertClient<InterceptedService<T, F>>
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static("/vald.v1.Upsert/Upsert");
            let mut req = request.into_request();
            req.extensions_mut().insert(GrpcMethod::new("vald.v1.Upsert", "Upsert"));
            self.inner.unary(req, path, codec).await
        }
        /** A method to insert/update multiple vectors by bidirectional streaming.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/vald.v1.Upsert/StreamUpsert",
            );
            let mut req = request.into_streaming_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("vald.v1.Upsert", "StreamUpsert"));
            self.inner.streaming(req, path, codec).await
        }
        /** A method to insert/update multiple vectors in a single request.
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
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
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
