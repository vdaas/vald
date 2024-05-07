use proto::{
    payload::v1::{object, upsert},
    vald::v1::upsert_server,
};

#[tonic::async_trait]
impl upsert_server::Upsert for super::Agent {
    async fn upsert(
        &self,
        request: tonic::Request<upsert::Request>,
    ) -> std::result::Result<tonic::Response<object::Location>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamUpsert method."]
    type StreamUpsertStream = crate::stream_type!(object::StreamLocation);

    #[doc = " A method to insert/update multiple vectors by bidirectional streaming.\n"]
    async fn stream_upsert(
        &self,
        request: tonic::Request<tonic::Streaming<upsert::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamUpsertStream>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to insert/update multiple vectors in a single request.\n"]
    async fn multi_upsert(
        &self,
        request: tonic::Request<upsert::MultiRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        todo!()
    }
}
