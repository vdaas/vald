use proto::{
    payload::v1::{insert, object},
    vald::v1::insert_server,
};
#[tonic::async_trait]
impl insert_server::Insert for super::Agent {
    async fn insert(
        &self,
        request: tonic::Request<insert::Request>,
    ) -> std::result::Result<tonic::Response<object::Location>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamInsert method."]
    type StreamInsertStream = crate::stream_type!(object::StreamLocation);

    #[doc = " A method to add new multiple vectors by bidirectional streaming.\n"]
    async fn stream_insert(
        &self,
        request: tonic::Request<tonic::Streaming<insert::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamInsertStream>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to add new multiple vectors in a single request.\n"]
    async fn multi_insert(
        &self,
        request: tonic::Request<insert::MultiRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        todo!()
    }
}
