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
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Error2 {
    /// ErrorType type = 1;
    #[prost(string, tag="2")]
    pub message: ::prost::alloc::string::String,
}
impl ::prost::Name for Error2 {
const NAME: &'static str = "Error2";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.Error2".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.Error2".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct ResponseHeader {
    /// cluster_id is the ID of the cluster which sent the response.
    #[prost(uint64, tag="1")]
    pub cluster_id: u64,
    #[prost(message, optional, tag="2")]
    pub error: ::core::option::Option<Error2>,
}
impl ::prost::Name for ResponseHeader {
const NAME: &'static str = "ResponseHeader";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.ResponseHeader".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.ResponseHeader".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct GetClusterInfoRequest {
    #[prost(message, optional, tag="1")]
    pub header: ::core::option::Option<ResponseHeader>,
}
impl ::prost::Name for GetClusterInfoRequest {
const NAME: &'static str = "GetClusterInfoRequest";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.GetClusterInfoRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.GetClusterInfoRequest".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct GetClusterInfoResponse {
    #[prost(message, optional, tag="1")]
    pub header: ::core::option::Option<ResponseHeader>,
}
impl ::prost::Name for GetClusterInfoResponse {
const NAME: &'static str = "GetClusterInfoResponse";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.GetClusterInfoResponse".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.GetClusterInfoResponse".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RequestHeader {
    /// cluster_id is the ID of the cluster which be sent to.
    #[prost(uint64, tag="1")]
    pub cluster_id: u64,
    /// sender_id is the ID of the sender server, also member ID or etcd ID.
    /// sender_id is used in PD internal communication.
    #[prost(uint64, tag="2")]
    pub sender_id: u64,
    /// caller_id is the ID of the client which sends the request, such as tikv,
    /// tidb, cdc, etc.
    #[prost(string, tag="3")]
    pub caller_id: ::prost::alloc::string::String,
    /// caller_component is the component of the client which sends the request,
    /// such as ddl, optimizer, etc.
    #[prost(string, tag="4")]
    pub caller_component: ::prost::alloc::string::String,
}
impl ::prost::Name for RequestHeader {
const NAME: &'static str = "RequestHeader";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.RequestHeader".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.RequestHeader".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct GetAllStoresRequest {
    #[prost(message, optional, tag="1")]
    pub header: ::core::option::Option<RequestHeader>,
    /// Do NOT return tombstone stores if set to true.
    #[prost(bool, tag="2")]
    pub exclude_tombstone_stores: bool,
}
impl ::prost::Name for GetAllStoresRequest {
const NAME: &'static str = "GetAllStoresRequest";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.GetAllStoresRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.GetAllStoresRequest".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetAllStoresResponse {
    #[prost(message, optional, tag="1")]
    pub header: ::core::option::Option<ResponseHeader>,
    #[prost(message, repeated, tag="2")]
    pub stores: ::prost::alloc::vec::Vec<super::metapb::Store>,
}
impl ::prost::Name for GetAllStoresResponse {
const NAME: &'static str = "GetAllStoresResponse";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.GetAllStoresResponse".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.GetAllStoresResponse".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Region {
    #[prost(message, optional, tag="1")]
    pub region: ::core::option::Option<super::metapb::Region2>,
    #[prost(message, optional, tag="2")]
    pub leader: ::core::option::Option<super::metapb::Peer>,
    /// Leader considers that these peers are down.
    /// repeated PeerStats down_peers = 3;
    /// Pending peers are the peers that the leader can't consider as
    /// working followers.
    ///
    /// buckets isn't nil only when need_buckets is true.
    /// metapb.Buckets buckets = 5;
    #[prost(message, repeated, tag="4")]
    pub pending_peers: ::prost::alloc::vec::Vec<super::metapb::Peer>,
}
impl ::prost::Name for Region {
const NAME: &'static str = "Region";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.Region".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.Region".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct KeyRange {
    #[prost(bytes="vec", tag="1")]
    pub start_key: ::prost::alloc::vec::Vec<u8>,
    /// end_key is +inf when it is empty.
    #[prost(bytes="vec", tag="2")]
    pub end_key: ::prost::alloc::vec::Vec<u8>,
}
impl ::prost::Name for KeyRange {
const NAME: &'static str = "KeyRange";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.KeyRange".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.KeyRange".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchScanRegionsRequest {
    #[prost(message, optional, tag="1")]
    pub header: ::core::option::Option<RequestHeader>,
    #[prost(bool, tag="2")]
    pub need_buckets: bool,
    /// the given ranges must be in order.
    #[prost(message, repeated, tag="3")]
    pub ranges: ::prost::alloc::vec::Vec<KeyRange>,
    /// limit the total number of regions to scan.
    #[prost(int32, tag="4")]
    pub limit: i32,
    /// If contain_all_key_range is true, the output must contain all
    /// key ranges in the request.
    /// If the output does not contain all key ranges, the request is considered
    /// failed and returns an error(REGIONS_NOT_CONTAIN_ALL_KEY_RANGE).
    #[prost(bool, tag="5")]
    pub contain_all_key_range: bool,
}
impl ::prost::Name for BatchScanRegionsRequest {
const NAME: &'static str = "BatchScanRegionsRequest";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.BatchScanRegionsRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.BatchScanRegionsRequest".into() }}
/// ResponseHeader header = 1;
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct BatchScanRegionsResponse {
    /// the returned regions are flattened into a list, because the given ranges can located in the same range, we do not return duplicated regions then.
    #[prost(message, repeated, tag="2")]
    pub regions: ::prost::alloc::vec::Vec<Region>,
}
impl ::prost::Name for BatchScanRegionsResponse {
const NAME: &'static str = "BatchScanRegionsResponse";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.BatchScanRegionsResponse".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.BatchScanRegionsResponse".into() }}
include!("pdpb.serde.rs");
// @@protoc_insertion_point(module)
