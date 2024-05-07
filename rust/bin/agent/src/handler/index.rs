use proto::{
    core::v1::agent_server,
    payload::v1::{control, info, object, Empty},
};

#[tonic::async_trait]
impl agent_server::Agent for super::Agent {
    async fn create_index(
        &self,
        request: tonic::Request<control::CreateIndexRequest>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        todo!()
    }

    async fn save_index(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        todo!()
    }

    #[doc = " Represent the creating and saving index RPC.\n"]
    async fn create_and_save_index(
        &self,
        request: tonic::Request<control::CreateIndexRequest>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        todo!()
    }

    #[doc = " Represent the RPC to get the agent index information.\n"]
    async fn index_info(
        &self,
        request: tonic::Request<Empty>,
    ) -> std::result::Result<tonic::Response<info::index::Count>, tonic::Status> {
        todo!()
    }

    #[doc = " Represent the RPC to get the vector metadata. This RPC is mainly used for index correction process\n"]
    async fn get_timestamp(
        &self,
        request: tonic::Request<object::GetTimestampRequest>,
    ) -> std::result::Result<tonic::Response<object::Timestamp>, tonic::Status> {
        todo!()
    }
}
