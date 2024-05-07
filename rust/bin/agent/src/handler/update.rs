use proto::{
    payload::v1::{object, update},
    vald::v1::update_server,
};

#[tonic::async_trait]
impl update_server::Update for super::Agent {
    async fn update(
        &self,
        request: tonic::Request<update::Request>,
    ) -> std::result::Result<tonic::Response<object::Location>, tonic::Status> {
        todo!()
    }

    #[doc = " Server streaming response type for the StreamUpdate method."]
    type StreamUpdateStream = crate::stream_type!(object::StreamLocation);

    #[doc = " A method to update multiple indexed vectors by bidirectional streaming.\n"]
    async fn stream_update(
        &self,
        request: tonic::Request<tonic::Streaming<update::Request>>,
    ) -> std::result::Result<tonic::Response<Self::StreamUpdateStream>, tonic::Status> {
        todo!()
    }

    #[doc = " A method to update multiple indexed vectors in a single request.\n"]
    async fn multi_update(
        &self,
        request: tonic::Request<update::MultiRequest>,
    ) -> std::result::Result<tonic::Response<object::Locations>, tonic::Status> {
        todo!()
    }
}
