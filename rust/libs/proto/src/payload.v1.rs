//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Search {
}
/// Nested message and enum types in `Search`.
pub mod search {
    /// Represent a search request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Request {
        /// The vector to be searched.
        #[prost(float, repeated, packed="false", tag="1")]
        pub vector: ::prost::alloc::vec::Vec<f32>,
        /// The configuration of the search request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
    /// Represent the multiple search request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiRequest {
        /// Represent the multiple search request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<Request>,
    }
    /// Represent a search by ID request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct IdRequest {
        /// The vector ID to be searched.
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
        /// The configuration of the search request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
    /// Represent the multiple search by ID request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiIdRequest {
        /// Represent the multiple search by ID request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<IdRequest>,
    }
    /// Represent a search by binary object request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct ObjectRequest {
        /// The binary object to be searched.
        #[prost(bytes="vec", tag="1")]
        pub object: ::prost::alloc::vec::Vec<u8>,
        /// The configuration of the search request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
        /// Filter configuration.
        #[prost(message, optional, tag="3")]
        pub vectorizer: ::core::option::Option<super::filter::Target>,
    }
    /// Represent the multiple search by binary object request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiObjectRequest {
        /// Represent the multiple search by binary object request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<ObjectRequest>,
    }
    /// Represent search configuration.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Config {
        /// Unique request ID.
        #[prost(string, tag="1")]
        pub request_id: ::prost::alloc::string::String,
        /// Maximum number of result to be returned.
        #[prost(uint32, tag="2")]
        pub num: u32,
        /// Search radius.
        #[prost(float, tag="3")]
        pub radius: f32,
        /// Search coefficient.
        #[prost(float, tag="4")]
        pub epsilon: f32,
        /// Search timeout in nanoseconds.
        #[prost(int64, tag="5")]
        pub timeout: i64,
        /// Ingress filter configurations.
        #[prost(message, repeated, tag="6")]
        pub ingress_filters: ::prost::alloc::vec::Vec<super::filter::Config>,
        /// Egress filter configurations.
        #[prost(message, repeated, tag="7")]
        pub egress_filters: ::prost::alloc::vec::Vec<super::filter::Config>,
        /// Minimum number of result to be returned.
        #[prost(uint32, tag="8")]
        pub min_num: u32,
        /// Aggregation Algorithm
        #[prost(enumeration="AggregationAlgorithm", tag="9")]
        pub aggregation_algorithm: i32,
    }
    /// Represent a search response.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Response {
        /// The unique request ID.
        #[prost(string, tag="1")]
        pub request_id: ::prost::alloc::string::String,
        /// Search results.
        #[prost(message, repeated, tag="2")]
        pub results: ::prost::alloc::vec::Vec<super::object::Distance>,
    }
    /// Represent multiple search responses.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Responses {
        /// Represent the multiple search response content.
        #[prost(message, repeated, tag="1")]
        pub responses: ::prost::alloc::vec::Vec<Response>,
    }
    /// Represent stream search response.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct StreamResponse {
        #[prost(oneof="stream_response::Payload", tags="1, 2")]
        pub payload: ::core::option::Option<stream_response::Payload>,
    }
    /// Nested message and enum types in `StreamResponse`.
    pub mod stream_response {
        #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Oneof)]
        pub enum Payload {
            /// Represent the search response.
            #[prost(message, tag="1")]
            Response(super::Response),
            /// The RPC error status.
            #[prost(message, tag="2")]
            Status(super::super::super::super::google::rpc::Status),
        }
    }
    /// AggregationAlgorithm is enum of each aggregation algorithms
    #[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
    #[repr(i32)]
    pub enum AggregationAlgorithm {
        Unknown = 0,
        ConcurrentQueue = 1,
        SortSlice = 2,
        SortPoolSlice = 3,
        PairingHeap = 4,
    }
    impl AggregationAlgorithm {
        /// String value of the enum field names used in the ProtoBuf definition.
        ///
        /// The values are not transformed in any way and thus are considered stable
        /// (if the ProtoBuf definition does not change) and safe for programmatic use.
        pub fn as_str_name(&self) -> &'static str {
            match self {
                AggregationAlgorithm::Unknown => "Unknown",
                AggregationAlgorithm::ConcurrentQueue => "ConcurrentQueue",
                AggregationAlgorithm::SortSlice => "SortSlice",
                AggregationAlgorithm::SortPoolSlice => "SortPoolSlice",
                AggregationAlgorithm::PairingHeap => "PairingHeap",
            }
        }
        /// Creates an enum from field names used in the ProtoBuf definition.
        pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
            match value {
                "Unknown" => Some(Self::Unknown),
                "ConcurrentQueue" => Some(Self::ConcurrentQueue),
                "SortSlice" => Some(Self::SortSlice),
                "SortPoolSlice" => Some(Self::SortPoolSlice),
                "PairingHeap" => Some(Self::PairingHeap),
                _ => None,
            }
        }
    }
}
/// Filter related messages.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Filter {
}
/// Nested message and enum types in `Filter`.
pub mod filter {
    /// Represent the target filter server.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Target {
        /// The target hostname.
        #[prost(string, tag="1")]
        pub host: ::prost::alloc::string::String,
        /// The target port.
        #[prost(uint32, tag="2")]
        pub port: u32,
    }
    /// Represent the filter query.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Query {
        /// The raw query string.
        #[prost(string, tag="1")]
        pub query: ::prost::alloc::string::String,
    }
    /// Represent filter configuration.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Config {
        /// Represent the filter target configuration.
        #[prost(message, optional, tag="1")]
        pub target: ::core::option::Option<Target>,
        /// The target query.
        #[prost(message, optional, tag="2")]
        pub query: ::core::option::Option<Query>,
    }
    /// Represent the ID and distance pair.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct DistanceRequest {
        /// Distance
        #[prost(message, repeated, tag="1")]
        pub distance: ::prost::alloc::vec::Vec<super::object::Distance>,
        /// Query
        #[prost(message, optional, tag="2")]
        pub query: ::core::option::Option<Query>,
    }
    /// Represent the ID and distance pair.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct DistanceResponse {
        /// Distance
        #[prost(message, repeated, tag="1")]
        pub distance: ::prost::alloc::vec::Vec<super::object::Distance>,
    }
    /// Represent the ID and vector pair.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct VectorRequest {
        /// Vector
        #[prost(message, optional, tag="1")]
        pub vector: ::core::option::Option<super::object::Vector>,
        /// Query
        #[prost(message, optional, tag="2")]
        pub query: ::core::option::Option<Query>,
    }
    /// Represent the ID and vector pair.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct VectorResponse {
        /// Distance
        #[prost(message, optional, tag="1")]
        pub vector: ::core::option::Option<super::object::Vector>,
    }
}
/// Insert related messages.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Insert {
}
/// Nested message and enum types in `Insert`.
pub mod insert {
    /// Represent the insert request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Request {
        /// The vector to be inserted.
        #[prost(message, optional, tag="1")]
        pub vector: ::core::option::Option<super::object::Vector>,
        /// The configuration of the insert request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
    /// Represent the multiple insert request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiRequest {
        /// Represent multiple insert request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<Request>,
    }
    /// Represent the insert by binary object request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct ObjectRequest {
        /// The binary object to be inserted.
        #[prost(message, optional, tag="1")]
        pub object: ::core::option::Option<super::object::Blob>,
        /// The configuration of the insert request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
        /// Filter configurations.
        #[prost(message, optional, tag="3")]
        pub vectorizer: ::core::option::Option<super::filter::Target>,
    }
    /// Represent the multiple insert by binary object request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiObjectRequest {
        /// Represent multiple insert by object content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<ObjectRequest>,
    }
    /// Represent insert configurations.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Config {
        /// A flag to skip exist check during insert operation.
        #[prost(bool, tag="1")]
        pub skip_strict_exist_check: bool,
        /// Filter configurations.
        #[prost(message, repeated, tag="2")]
        pub filters: ::prost::alloc::vec::Vec<super::filter::Config>,
        /// Insert timestamp.
        #[prost(int64, tag="3")]
        pub timestamp: i64,
    }
}
/// Update related messages
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Update {
}
/// Nested message and enum types in `Update`.
pub mod update {
    /// Represent the update request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Request {
        /// The vector to be updated.
        #[prost(message, optional, tag="1")]
        pub vector: ::core::option::Option<super::object::Vector>,
        /// The configuration of the update request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
    /// Represent the multiple update request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiRequest {
        /// Represent the multiple update request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<Request>,
    }
    /// Represent the update binary object request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct ObjectRequest {
        /// The binary object to be updated.
        #[prost(message, optional, tag="1")]
        pub object: ::core::option::Option<super::object::Blob>,
        /// The configuration of the update request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
        /// Filter target.
        #[prost(message, optional, tag="3")]
        pub vectorizer: ::core::option::Option<super::filter::Target>,
    }
    /// Represent the multiple update binary object request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiObjectRequest {
        /// Represent the multiple update object request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<ObjectRequest>,
    }
    /// Represent the update configuration.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Config {
        /// A flag to skip exist check during update operation.
        #[prost(bool, tag="1")]
        pub skip_strict_exist_check: bool,
        /// Filter configuration.
        #[prost(message, repeated, tag="2")]
        pub filters: ::prost::alloc::vec::Vec<super::filter::Config>,
        /// Update timestamp.
        #[prost(int64, tag="3")]
        pub timestamp: i64,
        /// A flag to disable balanced update (split remove -> insert operation)
        /// during update operation.
        #[prost(bool, tag="4")]
        pub disable_balanced_update: bool,
    }
}
/// Upsert related messages.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Upsert {
}
/// Nested message and enum types in `Upsert`.
pub mod upsert {
    /// Represent the upsert request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Request {
        /// The vector to be upserted.
        #[prost(message, optional, tag="1")]
        pub vector: ::core::option::Option<super::object::Vector>,
        /// The configuration of the upsert request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
    /// Represent mthe ultiple upsert request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiRequest {
        /// Represent the multiple upsert request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<Request>,
    }
    /// Represent the upsert binary object request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct ObjectRequest {
        /// The binary object to be upserted.
        #[prost(message, optional, tag="1")]
        pub object: ::core::option::Option<super::object::Blob>,
        /// The configuration of the upsert request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
        /// Filter target.
        #[prost(message, optional, tag="3")]
        pub vectorizer: ::core::option::Option<super::filter::Target>,
    }
    /// Represent the multiple upsert binary object request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiObjectRequest {
        /// Represent the multiple upsert object request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<ObjectRequest>,
    }
    /// Represent the upsert configuration.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Config {
        /// A flag to skip exist check during upsert operation.
        #[prost(bool, tag="1")]
        pub skip_strict_exist_check: bool,
        /// Filter configuration.
        #[prost(message, repeated, tag="2")]
        pub filters: ::prost::alloc::vec::Vec<super::filter::Config>,
        /// Upsert timestamp.
        #[prost(int64, tag="3")]
        pub timestamp: i64,
        /// A flag to disable balanced update (split remove -> insert operation)
        /// during update operation.
        #[prost(bool, tag="4")]
        pub disable_balanced_update: bool,
    }
}
/// Remove related messages.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Remove {
}
/// Nested message and enum types in `Remove`.
pub mod remove {
    /// Represent the remove request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Request {
        /// The object ID to be removed.
        #[prost(message, optional, tag="1")]
        pub id: ::core::option::Option<super::object::Id>,
        /// The configuration of the remove request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
    /// Represent the multiple remove request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiRequest {
        /// Represent the multiple remove request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<Request>,
    }
    /// Represent the remove request based on timestamp.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct TimestampRequest {
        /// The timestamp comparison list. If more than one is specified, the `AND` search is applied.
        #[prost(message, repeated, tag="1")]
        pub timestamps: ::prost::alloc::vec::Vec<Timestamp>,
    }
    /// Represent the timestamp comparison.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Timestamp {
        /// The timestamp.
        #[prost(int64, tag="1")]
        pub timestamp: i64,
        /// The conditional operator.
        #[prost(enumeration="timestamp::Operator", tag="2")]
        pub operator: i32,
    }
    /// Nested message and enum types in `Timestamp`.
    pub mod timestamp {
        /// Operator is enum of each conditional operator.
        #[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
        #[repr(i32)]
        pub enum Operator {
            /// The timestamp is equal to the specified value in the request.
            Eq = 0,
            /// The timestamp is not equal to the specified value in the request.
            Ne = 1,
            /// The timestamp is greater than or equal to the specified value in the request.
            Ge = 2,
            /// The timestamp is greater than the specified value in the request.
            Gt = 3,
            /// The timestamp is less than or equal to the specified value in the request.
            Le = 4,
            /// The timestamp is less than the specified value in the request.
            Lt = 5,
        }
        impl Operator {
            /// String value of the enum field names used in the ProtoBuf definition.
            ///
            /// The values are not transformed in any way and thus are considered stable
            /// (if the ProtoBuf definition does not change) and safe for programmatic use.
            pub fn as_str_name(&self) -> &'static str {
                match self {
                    Operator::Eq => "Eq",
                    Operator::Ne => "Ne",
                    Operator::Ge => "Ge",
                    Operator::Gt => "Gt",
                    Operator::Le => "Le",
                    Operator::Lt => "Lt",
                }
            }
            /// Creates an enum from field names used in the ProtoBuf definition.
            pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
                match value {
                    "Eq" => Some(Self::Eq),
                    "Ne" => Some(Self::Ne),
                    "Ge" => Some(Self::Ge),
                    "Gt" => Some(Self::Gt),
                    "Le" => Some(Self::Le),
                    "Lt" => Some(Self::Lt),
                    _ => None,
                }
            }
        }
    }
    /// Represent the remove configuration.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Config {
        /// A flag to skip exist check during upsert operation.
        #[prost(bool, tag="1")]
        pub skip_strict_exist_check: bool,
        /// Remove timestamp.
        #[prost(int64, tag="3")]
        pub timestamp: i64,
    }
}
/// Common messages.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Object {
}
/// Nested message and enum types in `Object`.
pub mod object {
    /// Represent a request to fetch raw vector.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct VectorRequest {
        /// The vector ID to be fetched.
        #[prost(message, optional, tag="1")]
        pub id: ::core::option::Option<Id>,
        /// Filter configurations.
        #[prost(message, repeated, tag="2")]
        pub filters: ::prost::alloc::vec::Vec<super::filter::Config>,
    }
    /// Represent the ID and distance pair.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Distance {
        /// The vector ID.
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
        /// The distance.
        #[prost(float, tag="2")]
        pub distance: f32,
    }
    /// Represent stream response of distances.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct StreamDistance {
        #[prost(oneof="stream_distance::Payload", tags="1, 2")]
        pub payload: ::core::option::Option<stream_distance::Payload>,
    }
    /// Nested message and enum types in `StreamDistance`.
    pub mod stream_distance {
        #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Oneof)]
        pub enum Payload {
            /// The distance.
            #[prost(message, tag="1")]
            Distance(super::Distance),
            /// The RPC error status.
            #[prost(message, tag="2")]
            Status(super::super::super::super::google::rpc::Status),
        }
    }
    /// Represent the vector ID.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Id {
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
    }
    /// Represent multiple vector IDs.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct IDs {
        #[prost(string, repeated, tag="1")]
        pub ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    }
    /// Represent a vector.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Vector {
        /// The vector ID.
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
        /// The vector.
        #[prost(float, repeated, packed="false", tag="2")]
        pub vector: ::prost::alloc::vec::Vec<f32>,
        /// timestamp represents when this vector inserted.
        #[prost(int64, tag="3")]
        pub timestamp: i64,
    }
    /// Represent a request to fetch vector meta data.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct GetTimestampRequest {
        /// The vector ID to be fetched.
        #[prost(message, optional, tag="1")]
        pub id: ::core::option::Option<Id>,
    }
    /// Represent a vector meta data.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Timestamp {
        /// The vector ID.
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
        /// timestamp represents when this vector inserted.
        #[prost(int64, tag="2")]
        pub timestamp: i64,
    }
    /// Represent multiple vectors.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Vectors {
        #[prost(message, repeated, tag="1")]
        pub vectors: ::prost::alloc::vec::Vec<Vector>,
    }
    /// Represent stream response of the vector.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct StreamVector {
        #[prost(oneof="stream_vector::Payload", tags="1, 2")]
        pub payload: ::core::option::Option<stream_vector::Payload>,
    }
    /// Nested message and enum types in `StreamVector`.
    pub mod stream_vector {
        #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Oneof)]
        pub enum Payload {
            /// The vector.
            #[prost(message, tag="1")]
            Vector(super::Vector),
            /// The RPC error status.
            #[prost(message, tag="2")]
            Status(super::super::super::super::google::rpc::Status),
        }
    }
    /// Represent reshape vector.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct ReshapeVector {
        /// The binary object.
        #[prost(bytes="vec", tag="1")]
        pub object: ::prost::alloc::vec::Vec<u8>,
        /// The new shape.
        #[prost(int32, repeated, tag="2")]
        pub shape: ::prost::alloc::vec::Vec<i32>,
    }
    /// Represent the binary object.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Blob {
        /// The object ID.
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
        /// The binary object.
        #[prost(bytes="vec", tag="2")]
        pub object: ::prost::alloc::vec::Vec<u8>,
    }
    /// Represent stream response of binary objects.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct StreamBlob {
        #[prost(oneof="stream_blob::Payload", tags="1, 2")]
        pub payload: ::core::option::Option<stream_blob::Payload>,
    }
    /// Nested message and enum types in `StreamBlob`.
    pub mod stream_blob {
        #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Oneof)]
        pub enum Payload {
            /// The binary object.
            #[prost(message, tag="1")]
            Blob(super::Blob),
            /// The RPC error status.
            #[prost(message, tag="2")]
            Status(super::super::super::super::google::rpc::Status),
        }
    }
    /// Represent the vector location.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Location {
        /// The name of the location.
        #[prost(string, tag="1")]
        pub name: ::prost::alloc::string::String,
        /// The UUID of the vector.
        #[prost(string, tag="2")]
        pub uuid: ::prost::alloc::string::String,
        /// The IP list.
        #[prost(string, repeated, tag="3")]
        pub ips: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    }
    /// Represent the stream response of the vector location.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct StreamLocation {
        #[prost(oneof="stream_location::Payload", tags="1, 2")]
        pub payload: ::core::option::Option<stream_location::Payload>,
    }
    /// Nested message and enum types in `StreamLocation`.
    pub mod stream_location {
        #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Oneof)]
        pub enum Payload {
            /// The vector location.
            #[prost(message, tag="1")]
            Location(super::Location),
            /// The RPC error status.
            #[prost(message, tag="2")]
            Status(super::super::super::super::google::rpc::Status),
        }
    }
    /// Represent multiple vector locations.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Locations {
        #[prost(message, repeated, tag="1")]
        pub locations: ::prost::alloc::vec::Vec<Location>,
    }
    /// Represent the list object vector stream request and response.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct List {
    }
    /// Nested message and enum types in `List`.
    pub mod list {
        #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
        pub struct Request {
        }
        #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
        pub struct Response {
            #[prost(oneof="response::Payload", tags="1, 2")]
            pub payload: ::core::option::Option<response::Payload>,
        }
        /// Nested message and enum types in `Response`.
        pub mod response {
            #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Oneof)]
            pub enum Payload {
                /// The vector
                #[prost(message, tag="1")]
                Vector(super::super::Vector),
                /// The RPC error status.
                #[prost(message, tag="2")]
                Status(super::super::super::super::super::google::rpc::Status),
            }
        }
    }
}
/// Control related messages.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Control {
}
/// Nested message and enum types in `Control`.
pub mod control {
    /// Represent the create index request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct CreateIndexRequest {
        /// The pool size of the create index operation.
        #[prost(uint32, tag="1")]
        pub pool_size: u32,
    }
}
/// Discoverer related messages.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Discoverer {
}
/// Nested message and enum types in `Discoverer`.
pub mod discoverer {
    /// Represent the dicoverer request.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Request {
        /// The agent name to be discovered.
        #[prost(string, tag="1")]
        pub name: ::prost::alloc::string::String,
        /// The namespace to be discovered.
        #[prost(string, tag="2")]
        pub namespace: ::prost::alloc::string::String,
        /// The node to be discovered.
        #[prost(string, tag="3")]
        pub node: ::prost::alloc::string::String,
    }
}
/// Info related messages.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Info {
}
/// Nested message and enum types in `Info`.
pub mod info {
    /// Represent the index information messages.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Index {
    }
    /// Nested message and enum types in `Index`.
    pub mod index {
        /// Represent the index count message.
        #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
        pub struct Count {
            /// The stored index count.
            #[prost(uint32, tag="1")]
            pub stored: u32,
            /// The uncommitted index count.
            #[prost(uint32, tag="2")]
            pub uncommitted: u32,
            /// The indexing index count.
            #[prost(bool, tag="3")]
            pub indexing: bool,
            /// The saving index count.
            #[prost(bool, tag="4")]
            pub saving: bool,
        }
        /// Represent the UUID message.
        #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
        pub struct Uuid {
        }
        /// Nested message and enum types in `UUID`.
        pub mod uuid {
            /// The committed UUID.
            #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
            pub struct Committed {
                #[prost(string, tag="1")]
                pub uuid: ::prost::alloc::string::String,
            }
            /// The uncommitted UUID.
            #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
            pub struct Uncommitted {
                #[prost(string, tag="1")]
                pub uuid: ::prost::alloc::string::String,
            }
        }
    }
    /// Represent the pod information message.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Pod {
        /// The app name of the pod on the label.
        #[prost(string, tag="1")]
        pub app_name: ::prost::alloc::string::String,
        /// The name of the pod.
        #[prost(string, tag="2")]
        pub name: ::prost::alloc::string::String,
        /// The namespace of the pod.
        #[prost(string, tag="3")]
        pub namespace: ::prost::alloc::string::String,
        /// The IP of the pod.
        #[prost(string, tag="4")]
        pub ip: ::prost::alloc::string::String,
        /// The CPU information of the pod.
        #[prost(message, optional, tag="5")]
        pub cpu: ::core::option::Option<Cpu>,
        /// The memory information of the pod.
        #[prost(message, optional, tag="6")]
        pub memory: ::core::option::Option<Memory>,
        /// The node information of the pod.
        #[prost(message, optional, tag="7")]
        pub node: ::core::option::Option<Node>,
    }
    /// Represent the node information message.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Node {
        /// The name of the node.
        #[prost(string, tag="1")]
        pub name: ::prost::alloc::string::String,
        /// The internal IP address of the node.
        #[prost(string, tag="2")]
        pub internal_addr: ::prost::alloc::string::String,
        /// The external IP address of the node.
        #[prost(string, tag="3")]
        pub external_addr: ::prost::alloc::string::String,
        /// The CPU information of the node.
        #[prost(message, optional, tag="4")]
        pub cpu: ::core::option::Option<Cpu>,
        /// The memory information of the node.
        #[prost(message, optional, tag="5")]
        pub memory: ::core::option::Option<Memory>,
        /// The pod information of the node.
        #[prost(message, optional, tag="6")]
        pub pods: ::core::option::Option<Pods>,
    }
    /// Represent the service information message.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Service {
        /// The name of the svc.
        #[prost(string, tag="1")]
        pub name: ::prost::alloc::string::String,
        /// The cluster ip of the svc.
        #[prost(string, tag="2")]
        pub cluster_ip: ::prost::alloc::string::String,
        /// The cluster ips of the svc.
        #[prost(string, repeated, tag="3")]
        pub cluster_ips: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
        /// The port of the svc.
        #[prost(message, repeated, tag="4")]
        pub ports: ::prost::alloc::vec::Vec<ServicePort>,
        /// The labels of the service.
        #[prost(message, optional, tag="5")]
        pub labels: ::core::option::Option<Labels>,
        /// The annotations of the service.
        #[prost(message, optional, tag="6")]
        pub annotations: ::core::option::Option<Annotations>,
    }
    /// Represets the service port information message.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct ServicePort {
        /// The name of the port.
        #[prost(string, tag="1")]
        pub name: ::prost::alloc::string::String,
        /// The port number
        #[prost(int32, tag="2")]
        pub port: i32,
    }
    /// Represent the kubernetes labels.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Labels {
        #[prost(map="string, string", tag="1")]
        pub labels: ::std::collections::HashMap<::prost::alloc::string::String, ::prost::alloc::string::String>,
    }
    /// Represent the kubernetes annotations.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Annotations {
        #[prost(map="string, string", tag="1")]
        pub annotations: ::std::collections::HashMap<::prost::alloc::string::String, ::prost::alloc::string::String>,
    }
    /// Represent the CPU information message.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Cpu {
        /// The CPU resource limit.
        #[prost(double, tag="1")]
        pub limit: f64,
        /// The CPU resource requested.
        #[prost(double, tag="2")]
        pub request: f64,
        /// The CPU usage.
        #[prost(double, tag="3")]
        pub usage: f64,
    }
    /// Represent the memory information message.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Memory {
        /// The memory limit.
        #[prost(double, tag="1")]
        pub limit: f64,
        /// The memory requested.
        #[prost(double, tag="2")]
        pub request: f64,
        /// The memory usage.
        #[prost(double, tag="3")]
        pub usage: f64,
    }
    /// Represent the multiple pod information message.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Pods {
        /// The multiple pod information.
        #[prost(message, repeated, tag="1")]
        pub pods: ::prost::alloc::vec::Vec<Pod>,
    }
    /// Represent the multiple node information message.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Nodes {
        /// The multiple node information.
        #[prost(message, repeated, tag="1")]
        pub nodes: ::prost::alloc::vec::Vec<Node>,
    }
    /// Represent the multiple service information message.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Services {
        /// The multiple service information.
        #[prost(message, repeated, tag="1")]
        pub services: ::prost::alloc::vec::Vec<Service>,
    }
    /// Represent the multiple IP message.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct IPs {
        #[prost(string, repeated, tag="1")]
        pub ip: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    }
}
/// Mirror related messages.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Mirror {
}
/// Nested message and enum types in `Mirror`.
pub mod mirror {
    /// Represent server information.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Target {
        /// The target hostname.
        #[prost(string, tag="1")]
        pub host: ::prost::alloc::string::String,
        /// The target port.
        #[prost(uint32, tag="2")]
        pub port: u32,
    }
    /// Represent the multiple Target message.
    #[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Targets {
        /// The multiple target information.
        #[prost(message, repeated, tag="1")]
        pub targets: ::prost::alloc::vec::Vec<Target>,
    }
}
/// Represent an empty message.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Empty {
}
// @@protoc_insertion_point(module)
