#[macro_export]
macro_rules! stream_type {
    ($t:ty) => {
        std::pin::Pin<Box<dyn tokio_stream::Stream<Item = std::result::Result<$t, tonic::Status>> + Send>>
    };
}
