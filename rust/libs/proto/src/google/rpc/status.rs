#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Status {
    /// The status code, which should be an enum value of [google.rpc.Code][google.rpc.Code].
    #[prost(int32, tag="1")]
    pub code: i32,
    /// A developer-facing error message, which should be in English. Any
    /// user-facing error message should be localized and sent in the
    /// [google.rpc.Status.details][google.rpc.Status.details] field, or localized by the client.
    #[prost(string, tag="2")]
    pub message: ::prost::alloc::string::String,
    /// A list of messages that carry the error details.  There is a common set of
    /// message types for APIs to use.
    #[prost(message, repeated, tag="3")]
    pub details: ::prost::alloc::vec::Vec<::prost_types::Any>,
}
