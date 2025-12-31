//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Search {
}
/// Nested message and enum types in `Search`.
pub mod search {
    /// Represent a search request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Request {
        /// The vector to be searched.
        #[prost(float, repeated, packed="false", tag="1")]
        pub vector: ::prost::alloc::vec::Vec<f32>,
        /// The configuration of the search request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
impl ::prost::Name for Request {
const NAME: &'static str = "Request";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Search.Request".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Search.Request".into() }}
    /// Represent the multiple search request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiRequest {
        /// Represent the multiple search request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<Request>,
    }
impl ::prost::Name for MultiRequest {
const NAME: &'static str = "MultiRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Search.MultiRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Search.MultiRequest".into() }}
    /// Represent a search by ID request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct IdRequest {
        /// The vector ID to be searched.
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
        /// The configuration of the search request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
impl ::prost::Name for IdRequest {
const NAME: &'static str = "IDRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Search.IDRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Search.IDRequest".into() }}
    /// Represent the multiple search by ID request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiIdRequest {
        /// Represent the multiple search by ID request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<IdRequest>,
    }
impl ::prost::Name for MultiIdRequest {
const NAME: &'static str = "MultiIDRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Search.MultiIDRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Search.MultiIDRequest".into() }}
    /// Represent a search by binary object request.
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
impl ::prost::Name for ObjectRequest {
const NAME: &'static str = "ObjectRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Search.ObjectRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Search.ObjectRequest".into() }}
    /// Represent the multiple search by binary object request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiObjectRequest {
        /// Represent the multiple search by binary object request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<ObjectRequest>,
    }
impl ::prost::Name for MultiObjectRequest {
const NAME: &'static str = "MultiObjectRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Search.MultiObjectRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Search.MultiObjectRequest".into() }}
    /// Represent search configuration.
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
        #[prost(message, optional, tag="6")]
        pub ingress_filters: ::core::option::Option<super::filter::Config>,
        /// Egress filter configurations.
        #[prost(message, optional, tag="7")]
        pub egress_filters: ::core::option::Option<super::filter::Config>,
        /// Minimum number of result to be returned.
        #[prost(uint32, tag="8")]
        pub min_num: u32,
        /// Aggregation Algorithm
        #[prost(enumeration="AggregationAlgorithm", tag="9")]
        pub aggregation_algorithm: i32,
        /// Search ratio for agent return result number.
        #[prost(message, optional, tag="10")]
        pub ratio: ::core::option::Option<f32>,
        /// Search nprobe.
        #[prost(uint32, tag="11")]
        pub nprobe: u32,
    }
impl ::prost::Name for Config {
const NAME: &'static str = "Config";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Search.Config".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Search.Config".into() }}
    /// Represent a search response.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Response {
        /// The unique request ID.
        #[prost(string, tag="1")]
        pub request_id: ::prost::alloc::string::String,
        /// Search results.
        #[prost(message, repeated, tag="2")]
        pub results: ::prost::alloc::vec::Vec<super::object::Distance>,
    }
impl ::prost::Name for Response {
const NAME: &'static str = "Response";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Search.Response".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Search.Response".into() }}
    /// Represent multiple search responses.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Responses {
        /// Represent the multiple search response content.
        #[prost(message, repeated, tag="1")]
        pub responses: ::prost::alloc::vec::Vec<Response>,
    }
impl ::prost::Name for Responses {
const NAME: &'static str = "Responses";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Search.Responses".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Search.Responses".into() }}
    /// Represent stream search response.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct StreamResponse {
        #[prost(oneof="stream_response::Payload", tags="1, 2")]
        pub payload: ::core::option::Option<stream_response::Payload>,
    }
    /// Nested message and enum types in `StreamResponse`.
    pub mod stream_response {
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
impl ::prost::Name for StreamResponse {
const NAME: &'static str = "StreamResponse";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Search.StreamResponse".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Search.StreamResponse".into() }}
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
                Self::Unknown => "Unknown",
                Self::ConcurrentQueue => "ConcurrentQueue",
                Self::SortSlice => "SortSlice",
                Self::SortPoolSlice => "SortPoolSlice",
                Self::PairingHeap => "PairingHeap",
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
impl ::prost::Name for Search {
const NAME: &'static str = "Search";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Search".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Search".into() }}
/// Filter related messages.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Filter {
}
/// Nested message and enum types in `Filter`.
pub mod filter {
    /// Represent the target filter server.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct Target {
        /// The target hostname.
        #[prost(string, tag="1")]
        pub host: ::prost::alloc::string::String,
        /// The target port.
        #[prost(uint32, tag="2")]
        pub port: u32,
    }
impl ::prost::Name for Target {
const NAME: &'static str = "Target";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Filter.Target".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Filter.Target".into() }}
    /// Represent filter configuration.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Config {
        /// Represent the filter target configuration.
        #[prost(message, repeated, tag="1")]
        pub targets: ::prost::alloc::vec::Vec<Target>,
    }
impl ::prost::Name for Config {
const NAME: &'static str = "Config";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Filter.Config".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Filter.Config".into() }}
}
impl ::prost::Name for Filter {
const NAME: &'static str = "Filter";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Filter".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Filter".into() }}
/// Insert related messages.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Insert {
}
/// Nested message and enum types in `Insert`.
pub mod insert {
    /// Represent the insert request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Request {
        /// The vector to be inserted.
        #[prost(message, optional, tag="1")]
        pub vector: ::core::option::Option<super::object::Vector>,
        /// The configuration of the insert request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
impl ::prost::Name for Request {
const NAME: &'static str = "Request";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Insert.Request".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Insert.Request".into() }}
    /// Represent the multiple insert request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiRequest {
        /// Represent multiple insert request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<Request>,
    }
impl ::prost::Name for MultiRequest {
const NAME: &'static str = "MultiRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Insert.MultiRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Insert.MultiRequest".into() }}
    /// Represent the insert by binary object request.
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
impl ::prost::Name for ObjectRequest {
const NAME: &'static str = "ObjectRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Insert.ObjectRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Insert.ObjectRequest".into() }}
    /// Represent the multiple insert by binary object request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiObjectRequest {
        /// Represent multiple insert by object content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<ObjectRequest>,
    }
impl ::prost::Name for MultiObjectRequest {
const NAME: &'static str = "MultiObjectRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Insert.MultiObjectRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Insert.MultiObjectRequest".into() }}
    /// Represent insert configurations.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Config {
        /// A flag to skip exist check during insert operation.
        #[prost(bool, tag="1")]
        pub skip_strict_exist_check: bool,
        /// Filter configurations.
        #[prost(message, optional, tag="2")]
        pub filters: ::core::option::Option<super::filter::Config>,
        /// Insert timestamp.
        #[prost(int64, tag="3")]
        pub timestamp: i64,
    }
impl ::prost::Name for Config {
const NAME: &'static str = "Config";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Insert.Config".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Insert.Config".into() }}
}
impl ::prost::Name for Insert {
const NAME: &'static str = "Insert";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Insert".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Insert".into() }}
/// Update related messages
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Update {
}
/// Nested message and enum types in `Update`.
pub mod update {
    /// Represent the update request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Request {
        /// The vector to be updated.
        #[prost(message, optional, tag="1")]
        pub vector: ::core::option::Option<super::object::Vector>,
        /// The configuration of the update request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
impl ::prost::Name for Request {
const NAME: &'static str = "Request";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Update.Request".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Update.Request".into() }}
    /// Represent the multiple update request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiRequest {
        /// Represent the multiple update request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<Request>,
    }
impl ::prost::Name for MultiRequest {
const NAME: &'static str = "MultiRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Update.MultiRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Update.MultiRequest".into() }}
    /// Represent the update binary object request.
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
impl ::prost::Name for ObjectRequest {
const NAME: &'static str = "ObjectRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Update.ObjectRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Update.ObjectRequest".into() }}
    /// Represent the multiple update binary object request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiObjectRequest {
        /// Represent the multiple update object request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<ObjectRequest>,
    }
impl ::prost::Name for MultiObjectRequest {
const NAME: &'static str = "MultiObjectRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Update.MultiObjectRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Update.MultiObjectRequest".into() }}
    /// Represent a vector meta data.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct TimestampRequest {
        /// The vector ID.
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
        /// timestamp represents when this vector inserted.
        #[prost(int64, tag="2")]
        pub timestamp: i64,
        /// force represents forcefully update the timestamp.
        #[prost(bool, tag="3")]
        pub force: bool,
    }
impl ::prost::Name for TimestampRequest {
const NAME: &'static str = "TimestampRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Update.TimestampRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Update.TimestampRequest".into() }}
    /// Represent the update configuration.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Config {
        /// A flag to skip exist check during update operation.
        #[prost(bool, tag="1")]
        pub skip_strict_exist_check: bool,
        /// Filter configuration.
        #[prost(message, optional, tag="2")]
        pub filters: ::core::option::Option<super::filter::Config>,
        /// Update timestamp.
        #[prost(int64, tag="3")]
        pub timestamp: i64,
        /// A flag to disable balanced update (split remove -> insert operation)
        /// during update operation.
        #[prost(bool, tag="4")]
        pub disable_balanced_update: bool,
    }
impl ::prost::Name for Config {
const NAME: &'static str = "Config";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Update.Config".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Update.Config".into() }}
}
impl ::prost::Name for Update {
const NAME: &'static str = "Update";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Update".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Update".into() }}
/// Upsert related messages.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Upsert {
}
/// Nested message and enum types in `Upsert`.
pub mod upsert {
    /// Represent the upsert request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Request {
        /// The vector to be upserted.
        #[prost(message, optional, tag="1")]
        pub vector: ::core::option::Option<super::object::Vector>,
        /// The configuration of the upsert request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
impl ::prost::Name for Request {
const NAME: &'static str = "Request";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Upsert.Request".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Upsert.Request".into() }}
    /// Represent mthe ultiple upsert request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiRequest {
        /// Represent the multiple upsert request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<Request>,
    }
impl ::prost::Name for MultiRequest {
const NAME: &'static str = "MultiRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Upsert.MultiRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Upsert.MultiRequest".into() }}
    /// Represent the upsert binary object request.
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
impl ::prost::Name for ObjectRequest {
const NAME: &'static str = "ObjectRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Upsert.ObjectRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Upsert.ObjectRequest".into() }}
    /// Represent the multiple upsert binary object request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiObjectRequest {
        /// Represent the multiple upsert object request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<ObjectRequest>,
    }
impl ::prost::Name for MultiObjectRequest {
const NAME: &'static str = "MultiObjectRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Upsert.MultiObjectRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Upsert.MultiObjectRequest".into() }}
    /// Represent the upsert configuration.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Config {
        /// A flag to skip exist check during upsert operation.
        #[prost(bool, tag="1")]
        pub skip_strict_exist_check: bool,
        /// Filter configuration.
        #[prost(message, optional, tag="2")]
        pub filters: ::core::option::Option<super::filter::Config>,
        /// Upsert timestamp.
        #[prost(int64, tag="3")]
        pub timestamp: i64,
        /// A flag to disable balanced update (split remove -> insert operation)
        /// during update operation.
        #[prost(bool, tag="4")]
        pub disable_balanced_update: bool,
    }
impl ::prost::Name for Config {
const NAME: &'static str = "Config";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Upsert.Config".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Upsert.Config".into() }}
}
impl ::prost::Name for Upsert {
const NAME: &'static str = "Upsert";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Upsert".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Upsert".into() }}
/// Remove related messages.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Remove {
}
/// Nested message and enum types in `Remove`.
pub mod remove {
    /// Represent the remove request.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct Request {
        /// The object ID to be removed.
        #[prost(message, optional, tag="1")]
        pub id: ::core::option::Option<super::object::Id>,
        /// The configuration of the remove request.
        #[prost(message, optional, tag="2")]
        pub config: ::core::option::Option<Config>,
    }
impl ::prost::Name for Request {
const NAME: &'static str = "Request";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Remove.Request".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Remove.Request".into() }}
    /// Represent the multiple remove request.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct MultiRequest {
        /// Represent the multiple remove request content.
        #[prost(message, repeated, tag="1")]
        pub requests: ::prost::alloc::vec::Vec<Request>,
    }
impl ::prost::Name for MultiRequest {
const NAME: &'static str = "MultiRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Remove.MultiRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Remove.MultiRequest".into() }}
    /// Represent the remove request based on timestamp.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct TimestampRequest {
        /// The timestamp comparison list. If more than one is specified, the `AND`
        /// search is applied.
        #[prost(message, repeated, tag="1")]
        pub timestamps: ::prost::alloc::vec::Vec<Timestamp>,
    }
impl ::prost::Name for TimestampRequest {
const NAME: &'static str = "TimestampRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Remove.TimestampRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Remove.TimestampRequest".into() }}
    /// Represent the timestamp comparison.
    #[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
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
            /// The timestamp is greater than or equal to the specified value in the
            /// request.
            Ge = 2,
            /// The timestamp is greater than the specified value in the request.
            Gt = 3,
            /// The timestamp is less than or equal to the specified value in the
            /// request.
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
                    Self::Eq => "Eq",
                    Self::Ne => "Ne",
                    Self::Ge => "Ge",
                    Self::Gt => "Gt",
                    Self::Le => "Le",
                    Self::Lt => "Lt",
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
impl ::prost::Name for Timestamp {
const NAME: &'static str = "Timestamp";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Remove.Timestamp".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Remove.Timestamp".into() }}
    /// Represent the remove configuration.
    #[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct Config {
        /// A flag to skip exist check during upsert operation.
        #[prost(bool, tag="1")]
        pub skip_strict_exist_check: bool,
        /// Remove timestamp.
        #[prost(int64, tag="3")]
        pub timestamp: i64,
    }
impl ::prost::Name for Config {
const NAME: &'static str = "Config";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Remove.Config".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Remove.Config".into() }}
}
impl ::prost::Name for Remove {
const NAME: &'static str = "Remove";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Remove".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Remove".into() }}
/// Flush related messages.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Flush {
}
/// Nested message and enum types in `Flush`.
pub mod flush {
    #[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct Request {
    }
impl ::prost::Name for Request {
const NAME: &'static str = "Request";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Flush.Request".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Flush.Request".into() }}
}
impl ::prost::Name for Flush {
const NAME: &'static str = "Flush";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Flush".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Flush".into() }}
/// Common messages.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Object {
}
/// Nested message and enum types in `Object`.
pub mod object {
    /// Represent a request to fetch raw vector.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct VectorRequest {
        /// The vector ID to be fetched.
        #[prost(message, optional, tag="1")]
        pub id: ::core::option::Option<Id>,
        /// Filter configurations.
        #[prost(message, optional, tag="2")]
        pub filters: ::core::option::Option<super::filter::Config>,
    }
impl ::prost::Name for VectorRequest {
const NAME: &'static str = "VectorRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.VectorRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.VectorRequest".into() }}
    /// Represent the ID and distance pair.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Distance {
        /// The vector ID.
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
        /// The distance.
        #[prost(float, tag="2")]
        pub distance: f32,
    }
impl ::prost::Name for Distance {
const NAME: &'static str = "Distance";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.Distance".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.Distance".into() }}
    /// Represent stream response of distances.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct StreamDistance {
        #[prost(oneof="stream_distance::Payload", tags="1, 2")]
        pub payload: ::core::option::Option<stream_distance::Payload>,
    }
    /// Nested message and enum types in `StreamDistance`.
    pub mod stream_distance {
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
impl ::prost::Name for StreamDistance {
const NAME: &'static str = "StreamDistance";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.StreamDistance".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.StreamDistance".into() }}
    /// Represent the vector ID.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct Id {
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
    }
impl ::prost::Name for Id {
const NAME: &'static str = "ID";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.ID".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.ID".into() }}
    /// Represent multiple vector IDs.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct IDs {
        #[prost(string, repeated, tag="1")]
        pub ids: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    }
impl ::prost::Name for IDs {
const NAME: &'static str = "IDs";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.IDs".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.IDs".into() }}
    /// Represent a vector.
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
impl ::prost::Name for Vector {
const NAME: &'static str = "Vector";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.Vector".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.Vector".into() }}
    /// Represent a request to fetch vector meta data.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct TimestampRequest {
        /// The vector ID to be fetched.
        #[prost(message, optional, tag="1")]
        pub id: ::core::option::Option<Id>,
    }
impl ::prost::Name for TimestampRequest {
const NAME: &'static str = "TimestampRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.TimestampRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.TimestampRequest".into() }}
    /// Represent a vector meta data.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct Timestamp {
        /// The vector ID.
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
        /// timestamp represents when this vector inserted.
        #[prost(int64, tag="2")]
        pub timestamp: i64,
    }
impl ::prost::Name for Timestamp {
const NAME: &'static str = "Timestamp";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.Timestamp".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.Timestamp".into() }}
    /// Represent multiple vectors.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Vectors {
        #[prost(message, repeated, tag="1")]
        pub vectors: ::prost::alloc::vec::Vec<Vector>,
    }
impl ::prost::Name for Vectors {
const NAME: &'static str = "Vectors";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.Vectors".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.Vectors".into() }}
    /// Represent stream response of the vector.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct StreamVector {
        #[prost(oneof="stream_vector::Payload", tags="1, 2")]
        pub payload: ::core::option::Option<stream_vector::Payload>,
    }
    /// Nested message and enum types in `StreamVector`.
    pub mod stream_vector {
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
impl ::prost::Name for StreamVector {
const NAME: &'static str = "StreamVector";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.StreamVector".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.StreamVector".into() }}
    /// Represent reshape vector.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct ReshapeVector {
        /// The binary object.
        #[prost(bytes="vec", tag="1")]
        pub object: ::prost::alloc::vec::Vec<u8>,
        /// The new shape.
        #[prost(int32, repeated, tag="2")]
        pub shape: ::prost::alloc::vec::Vec<i32>,
    }
impl ::prost::Name for ReshapeVector {
const NAME: &'static str = "ReshapeVector";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.ReshapeVector".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.ReshapeVector".into() }}
    /// Represent the binary object.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct Blob {
        /// The object ID.
        #[prost(string, tag="1")]
        pub id: ::prost::alloc::string::String,
        /// The binary object.
        #[prost(bytes="vec", tag="2")]
        pub object: ::prost::alloc::vec::Vec<u8>,
    }
impl ::prost::Name for Blob {
const NAME: &'static str = "Blob";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.Blob".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.Blob".into() }}
    /// Represent stream response of binary objects.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct StreamBlob {
        #[prost(oneof="stream_blob::Payload", tags="1, 2")]
        pub payload: ::core::option::Option<stream_blob::Payload>,
    }
    /// Nested message and enum types in `StreamBlob`.
    pub mod stream_blob {
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
impl ::prost::Name for StreamBlob {
const NAME: &'static str = "StreamBlob";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.StreamBlob".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.StreamBlob".into() }}
    /// Represent the vector location.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
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
impl ::prost::Name for Location {
const NAME: &'static str = "Location";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.Location".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.Location".into() }}
    /// Represent the stream response of the vector location.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct StreamLocation {
        #[prost(oneof="stream_location::Payload", tags="1, 2")]
        pub payload: ::core::option::Option<stream_location::Payload>,
    }
    /// Nested message and enum types in `StreamLocation`.
    pub mod stream_location {
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
impl ::prost::Name for StreamLocation {
const NAME: &'static str = "StreamLocation";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.StreamLocation".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.StreamLocation".into() }}
    /// Represent multiple vector locations.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Locations {
        #[prost(message, repeated, tag="1")]
        pub locations: ::prost::alloc::vec::Vec<Location>,
    }
impl ::prost::Name for Locations {
const NAME: &'static str = "Locations";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.Locations".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.Locations".into() }}
    /// Represent the list object vector stream request and response.
    #[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct List {
    }
    /// Nested message and enum types in `List`.
    pub mod list {
        #[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
        pub struct Request {
        }
impl ::prost::Name for Request {
const NAME: &'static str = "Request";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.List.Request".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.List.Request".into() }}
        #[derive(Clone, PartialEq, ::prost::Message)]
        pub struct Response {
            #[prost(oneof="response::Payload", tags="1, 2")]
            pub payload: ::core::option::Option<response::Payload>,
        }
        /// Nested message and enum types in `Response`.
        pub mod response {
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
impl ::prost::Name for Response {
const NAME: &'static str = "Response";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.List.Response".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.List.Response".into() }}
    }
impl ::prost::Name for List {
const NAME: &'static str = "List";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object.List".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object.List".into() }}
}
impl ::prost::Name for Object {
const NAME: &'static str = "Object";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Object".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Object".into() }}
/// Control related messages.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Control {
}
/// Nested message and enum types in `Control`.
pub mod control {
    /// Represent the create index request.
    #[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct CreateIndexRequest {
        /// The pool size of the create index operation.
        #[prost(uint32, tag="1")]
        pub pool_size: u32,
    }
impl ::prost::Name for CreateIndexRequest {
const NAME: &'static str = "CreateIndexRequest";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Control.CreateIndexRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Control.CreateIndexRequest".into() }}
}
impl ::prost::Name for Control {
const NAME: &'static str = "Control";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Control".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Control".into() }}
/// Discoverer related messages.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Discoverer {
}
/// Nested message and enum types in `Discoverer`.
pub mod discoverer {
    /// Represent the dicoverer request.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
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
impl ::prost::Name for Request {
const NAME: &'static str = "Request";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Discoverer.Request".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Discoverer.Request".into() }}
}
impl ::prost::Name for Discoverer {
const NAME: &'static str = "Discoverer";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Discoverer".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Discoverer".into() }}
/// Info related messages.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Info {
}
/// Nested message and enum types in `Info`.
pub mod info {
    /// Represent the index information messages.
    #[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct Index {
    }
    /// Nested message and enum types in `Index`.
    pub mod index {
        /// Represent the index count message.
        #[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
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
impl ::prost::Name for Count {
const NAME: &'static str = "Count";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Index.Count".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Index.Count".into() }}
        /// Represent the index count for each Agents message.
        #[derive(Clone, PartialEq, ::prost::Message)]
        pub struct Detail {
            /// count infos for each agents
            #[prost(map="string, message", tag="1")]
            pub counts: ::std::collections::HashMap<::prost::alloc::string::String, Count>,
            /// index replica of vald cluster
            #[prost(uint32, tag="2")]
            pub replica: u32,
            /// live agent replica of vald cluster
            #[prost(uint32, tag="3")]
            pub live_agents: u32,
        }
impl ::prost::Name for Detail {
const NAME: &'static str = "Detail";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Index.Detail".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Index.Detail".into() }}
        /// Represent the UUID message.
        #[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
        pub struct Uuid {
        }
        /// Nested message and enum types in `UUID`.
        pub mod uuid {
            /// The committed UUID.
            #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
            pub struct Committed {
                #[prost(string, tag="1")]
                pub uuid: ::prost::alloc::string::String,
            }
impl ::prost::Name for Committed {
const NAME: &'static str = "Committed";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Index.UUID.Committed".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Index.UUID.Committed".into() }}
            /// The uncommitted UUID.
            #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
            pub struct Uncommitted {
                #[prost(string, tag="1")]
                pub uuid: ::prost::alloc::string::String,
            }
impl ::prost::Name for Uncommitted {
const NAME: &'static str = "Uncommitted";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Index.UUID.Uncommitted".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Index.UUID.Uncommitted".into() }}
        }
impl ::prost::Name for Uuid {
const NAME: &'static str = "UUID";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Index.UUID".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Index.UUID".into() }}
        /// Represents index Statistics
        #[derive(Clone, PartialEq, ::prost::Message)]
        pub struct Statistics {
            #[prost(bool, tag="1")]
            pub valid: bool,
            #[prost(int32, tag="2")]
            pub median_indegree: i32,
            #[prost(int32, tag="3")]
            pub median_outdegree: i32,
            #[prost(uint64, tag="4")]
            pub max_number_of_indegree: u64,
            #[prost(uint64, tag="5")]
            pub max_number_of_outdegree: u64,
            #[prost(uint64, tag="6")]
            pub min_number_of_indegree: u64,
            #[prost(uint64, tag="7")]
            pub min_number_of_outdegree: u64,
            #[prost(uint64, tag="8")]
            pub mode_indegree: u64,
            #[prost(uint64, tag="9")]
            pub mode_outdegree: u64,
            #[prost(uint64, tag="10")]
            pub nodes_skipped_for_10_edges: u64,
            #[prost(uint64, tag="11")]
            pub nodes_skipped_for_indegree_distance: u64,
            #[prost(uint64, tag="12")]
            pub number_of_edges: u64,
            #[prost(uint64, tag="13")]
            pub number_of_indexed_objects: u64,
            #[prost(uint64, tag="14")]
            pub number_of_nodes: u64,
            #[prost(uint64, tag="15")]
            pub number_of_nodes_without_edges: u64,
            #[prost(uint64, tag="16")]
            pub number_of_nodes_without_indegree: u64,
            #[prost(uint64, tag="17")]
            pub number_of_objects: u64,
            #[prost(uint64, tag="18")]
            pub number_of_removed_objects: u64,
            #[prost(uint64, tag="19")]
            pub size_of_object_repository: u64,
            #[prost(uint64, tag="20")]
            pub size_of_refinement_object_repository: u64,
            #[prost(double, tag="21")]
            pub variance_of_indegree: f64,
            #[prost(double, tag="22")]
            pub variance_of_outdegree: f64,
            #[prost(double, tag="23")]
            pub mean_edge_length: f64,
            #[prost(double, tag="24")]
            pub mean_edge_length_for_10_edges: f64,
            #[prost(double, tag="25")]
            pub mean_indegree_distance_for_10_edges: f64,
            #[prost(double, tag="26")]
            pub mean_number_of_edges_per_node: f64,
            #[prost(double, tag="27")]
            pub c1_indegree: f64,
            #[prost(double, tag="28")]
            pub c5_indegree: f64,
            #[prost(double, tag="29")]
            pub c95_outdegree: f64,
            #[prost(double, tag="30")]
            pub c99_outdegree: f64,
            #[prost(int64, repeated, tag="31")]
            pub indegree_count: ::prost::alloc::vec::Vec<i64>,
            #[prost(uint64, repeated, tag="32")]
            pub outdegree_histogram: ::prost::alloc::vec::Vec<u64>,
            #[prost(uint64, repeated, tag="33")]
            pub indegree_histogram: ::prost::alloc::vec::Vec<u64>,
        }
impl ::prost::Name for Statistics {
const NAME: &'static str = "Statistics";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Index.Statistics".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Index.Statistics".into() }}
        /// Represents index Statistics for each Agents
        #[derive(Clone, PartialEq, ::prost::Message)]
        pub struct StatisticsDetail {
            /// count infos for each agents
            #[prost(map="string, message", tag="1")]
            pub details: ::std::collections::HashMap<::prost::alloc::string::String, Statistics>,
        }
impl ::prost::Name for StatisticsDetail {
const NAME: &'static str = "StatisticsDetail";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Index.StatisticsDetail".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Index.StatisticsDetail".into() }}
        /// Represents index Property
        #[derive(Clone, PartialEq, ::prost::Message)]
        pub struct Property {
            #[prost(int32, tag="1")]
            pub dimension: i32,
            #[prost(int32, tag="2")]
            pub thread_pool_size: i32,
            #[prost(string, tag="3")]
            pub object_type: ::prost::alloc::string::String,
            #[prost(string, tag="4")]
            pub distance_type: ::prost::alloc::string::String,
            #[prost(string, tag="5")]
            pub index_type: ::prost::alloc::string::String,
            #[prost(string, tag="6")]
            pub database_type: ::prost::alloc::string::String,
            #[prost(string, tag="7")]
            pub object_alignment: ::prost::alloc::string::String,
            #[prost(int32, tag="8")]
            pub path_adjustment_interval: i32,
            #[prost(int32, tag="9")]
            pub graph_shared_memory_size: i32,
            #[prost(int32, tag="10")]
            pub tree_shared_memory_size: i32,
            #[prost(int32, tag="11")]
            pub object_shared_memory_size: i32,
            #[prost(int32, tag="12")]
            pub prefetch_offset: i32,
            #[prost(int32, tag="13")]
            pub prefetch_size: i32,
            #[prost(string, tag="14")]
            pub accuracy_table: ::prost::alloc::string::String,
            #[prost(string, tag="15")]
            pub search_type: ::prost::alloc::string::String,
            #[prost(float, tag="16")]
            pub max_magnitude: f32,
            #[prost(int32, tag="17")]
            pub n_of_neighbors_for_insertion_order: i32,
            #[prost(float, tag="18")]
            pub epsilon_for_insertion_order: f32,
            #[prost(string, tag="19")]
            pub refinement_object_type: ::prost::alloc::string::String,
            #[prost(int32, tag="20")]
            pub truncation_threshold: i32,
            #[prost(int32, tag="21")]
            pub edge_size_for_creation: i32,
            #[prost(int32, tag="22")]
            pub edge_size_for_search: i32,
            #[prost(int32, tag="23")]
            pub edge_size_limit_for_creation: i32,
            #[prost(double, tag="24")]
            pub insertion_radius_coefficient: f64,
            #[prost(int32, tag="25")]
            pub seed_size: i32,
            #[prost(string, tag="26")]
            pub seed_type: ::prost::alloc::string::String,
            #[prost(int32, tag="27")]
            pub truncation_thread_pool_size: i32,
            #[prost(int32, tag="28")]
            pub batch_size_for_creation: i32,
            #[prost(string, tag="29")]
            pub graph_type: ::prost::alloc::string::String,
            #[prost(int32, tag="30")]
            pub dynamic_edge_size_base: i32,
            #[prost(int32, tag="31")]
            pub dynamic_edge_size_rate: i32,
            #[prost(float, tag="32")]
            pub build_time_limit: f32,
            #[prost(int32, tag="33")]
            pub outgoing_edge: i32,
            #[prost(int32, tag="34")]
            pub incoming_edge: i32,
            #[prost(float, tag="35")]
            pub epsilon_for_creation: f32,
        }
impl ::prost::Name for Property {
const NAME: &'static str = "Property";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Index.Property".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Index.Property".into() }}
        /// Represents index Properties for each Agents
        #[derive(Clone, PartialEq, ::prost::Message)]
        pub struct PropertyDetail {
            #[prost(map="string, message", tag="1")]
            pub details: ::std::collections::HashMap<::prost::alloc::string::String, Property>,
        }
impl ::prost::Name for PropertyDetail {
const NAME: &'static str = "PropertyDetail";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Index.PropertyDetail".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Index.PropertyDetail".into() }}
    }
impl ::prost::Name for Index {
const NAME: &'static str = "Index";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Index".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Index".into() }}
    /// Represent the resource stats
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct ResourceStats {
        #[prost(string, tag="1")]
        pub name: ::prost::alloc::string::String,
        #[prost(string, tag="2")]
        pub ip: ::prost::alloc::string::String,
        /// Container resource usage statistics
        #[prost(message, optional, tag="3")]
        pub cgroup_stats: ::core::option::Option<CgroupStats>,
    }
impl ::prost::Name for ResourceStats {
const NAME: &'static str = "ResourceStats";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.ResourceStats".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.ResourceStats".into() }}
    #[derive(Clone, Copy, PartialEq, ::prost::Message)]
    pub struct CgroupStats {
        /// CPU cores available
        #[prost(double, tag="1")]
        pub cpu_limit_cores: f64,
        /// CPU usage in cores (not percentage)
        #[prost(double, tag="2")]
        pub cpu_usage_cores: f64,
        /// Memory limit in bytes
        #[prost(uint64, tag="3")]
        pub memory_limit_bytes: u64,
        /// Memory usage in bytes
        #[prost(uint64, tag="4")]
        pub memory_usage_bytes: u64,
    }
impl ::prost::Name for CgroupStats {
const NAME: &'static str = "CgroupStats";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.CgroupStats".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.CgroupStats".into() }}
    /// Represent the pod information message.
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
impl ::prost::Name for Pod {
const NAME: &'static str = "Pod";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Pod".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Pod".into() }}
    /// Represent the node information message.
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
impl ::prost::Name for Node {
const NAME: &'static str = "Node";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Node".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Node".into() }}
    /// Represent the service information message.
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
impl ::prost::Name for Service {
const NAME: &'static str = "Service";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Service".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Service".into() }}
    /// Represets the service port information message.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct ServicePort {
        /// The name of the port.
        #[prost(string, tag="1")]
        pub name: ::prost::alloc::string::String,
        /// The port number
        #[prost(int32, tag="2")]
        pub port: i32,
    }
impl ::prost::Name for ServicePort {
const NAME: &'static str = "ServicePort";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.ServicePort".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.ServicePort".into() }}
    /// Represent the kubernetes labels.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Labels {
        #[prost(map="string, string", tag="1")]
        pub labels: ::std::collections::HashMap<::prost::alloc::string::String, ::prost::alloc::string::String>,
    }
impl ::prost::Name for Labels {
const NAME: &'static str = "Labels";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Labels".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Labels".into() }}
    /// Represent the kubernetes annotations.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Annotations {
        #[prost(map="string, string", tag="1")]
        pub annotations: ::std::collections::HashMap<::prost::alloc::string::String, ::prost::alloc::string::String>,
    }
impl ::prost::Name for Annotations {
const NAME: &'static str = "Annotations";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Annotations".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Annotations".into() }}
    /// Represent the CPU information message.
    #[derive(Clone, Copy, PartialEq, ::prost::Message)]
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
impl ::prost::Name for Cpu {
const NAME: &'static str = "CPU";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.CPU".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.CPU".into() }}
    /// Represent the memory information message.
    #[derive(Clone, Copy, PartialEq, ::prost::Message)]
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
impl ::prost::Name for Memory {
const NAME: &'static str = "Memory";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Memory".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Memory".into() }}
    /// Represent the multiple pod information message.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Pods {
        /// The multiple pod information.
        #[prost(message, repeated, tag="1")]
        pub pods: ::prost::alloc::vec::Vec<Pod>,
    }
impl ::prost::Name for Pods {
const NAME: &'static str = "Pods";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Pods".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Pods".into() }}
    /// Represent the multiple node information message.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Nodes {
        /// The multiple node information.
        #[prost(message, repeated, tag="1")]
        pub nodes: ::prost::alloc::vec::Vec<Node>,
    }
impl ::prost::Name for Nodes {
const NAME: &'static str = "Nodes";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Nodes".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Nodes".into() }}
    /// Represent the multiple service information message.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Services {
        /// The multiple service information.
        #[prost(message, repeated, tag="1")]
        pub services: ::prost::alloc::vec::Vec<Service>,
    }
impl ::prost::Name for Services {
const NAME: &'static str = "Services";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.Services".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.Services".into() }}
    /// Represent the multiple IP message.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct IPs {
        #[prost(string, repeated, tag="1")]
        pub ip: ::prost::alloc::vec::Vec<::prost::alloc::string::String>,
    }
impl ::prost::Name for IPs {
const NAME: &'static str = "IPs";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info.IPs".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info.IPs".into() }}
}
impl ::prost::Name for Info {
const NAME: &'static str = "Info";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Info".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Info".into() }}
/// Mirror related messages.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Mirror {
}
/// Nested message and enum types in `Mirror`.
pub mod mirror {
    /// Represent server information.
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct Target {
        /// The target hostname.
        #[prost(string, tag="1")]
        pub host: ::prost::alloc::string::String,
        /// The target port.
        #[prost(uint32, tag="2")]
        pub port: u32,
    }
impl ::prost::Name for Target {
const NAME: &'static str = "Target";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Mirror.Target".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Mirror.Target".into() }}
    /// Represent the multiple Target message.
    #[derive(Clone, PartialEq, ::prost::Message)]
    pub struct Targets {
        /// The multiple target information.
        #[prost(message, repeated, tag="1")]
        pub targets: ::prost::alloc::vec::Vec<Target>,
    }
impl ::prost::Name for Targets {
const NAME: &'static str = "Targets";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Mirror.Targets".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Mirror.Targets".into() }}
}
impl ::prost::Name for Mirror {
const NAME: &'static str = "Mirror";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Mirror".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Mirror".into() }}
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Meta {
}
/// Nested message and enum types in `Meta`.
pub mod meta {
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct Key {
        #[prost(string, tag="1")]
        pub key: ::prost::alloc::string::String,
    }
impl ::prost::Name for Key {
const NAME: &'static str = "Key";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Meta.Key".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Meta.Key".into() }}
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct Value {
        #[prost(message, optional, tag="1")]
        pub value: ::core::option::Option<::prost_types::Any>,
    }
impl ::prost::Name for Value {
const NAME: &'static str = "Value";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Meta.Value".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Meta.Value".into() }}
    #[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
    pub struct KeyValue {
        #[prost(message, optional, tag="1")]
        pub key: ::core::option::Option<Key>,
        #[prost(message, optional, tag="2")]
        pub value: ::core::option::Option<Value>,
    }
impl ::prost::Name for KeyValue {
const NAME: &'static str = "KeyValue";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Meta.KeyValue".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Meta.KeyValue".into() }}
}
impl ::prost::Name for Meta {
const NAME: &'static str = "Meta";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Meta".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Meta".into() }}
/// Represent an empty message.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Empty {
}
impl ::prost::Name for Empty {
const NAME: &'static str = "Empty";
const PACKAGE: &'static str = "payload.v1";
fn full_name() -> ::prost::alloc::string::String { "payload.v1.Empty".into() }fn type_url() -> ::prost::alloc::string::String { "/payload.v1.Empty".into() }}
// @@protoc_insertion_point(module)
