//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
pub struct Error {
    #[prost(enumeration="ErrorType", tag="1")]
    pub r#type: i32,
    #[prost(string, tag="2")]
    pub message: ::prost::alloc::string::String,
}
impl ::prost::Name for Error {
const NAME: &'static str = "Error";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.Error".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.Error".into() }}
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct GetAllStoresRequest {
    /// RequestHeader header = 1;
    /// Do NOT return tombstone stores if set to true.
    #[prost(bool, tag="2")]
    pub exclude_tombstone_stores: bool,
}
impl ::prost::Name for GetAllStoresRequest {
const NAME: &'static str = "GetAllStoresRequest";
const PACKAGE: &'static str = "pdpb";
fn full_name() -> ::prost::alloc::string::String { "pdpb.GetAllStoresRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/pdpb.GetAllStoresRequest".into() }}
/// ResponseHeader header = 1;
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GetAllStoresResponse {
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
    pub region: ::core::option::Option<super::metapb::Region>,
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
    /// RequestHeader header = 1;
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
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum ErrorType {
    Ok = 0,
    Unknown = 1,
    NotBootstrapped = 2,
    StoreTombstone = 3,
    AlreadyBootstrapped = 4,
    IncompatibleVersion = 5,
    RegionNotFound = 6,
    GlobalConfigNotFound = 7,
    DuplicatedEntry = 8,
    EntryNotFound = 9,
    InvalidValue = 10,
    /// required watch revision is smaller than current compact/min revision.
    DataCompacted = 11,
    RegionsNotContainAllKeyRange = 12,
}
impl ErrorType {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            Self::Ok => "OK",
            Self::Unknown => "UNKNOWN",
            Self::NotBootstrapped => "NOT_BOOTSTRAPPED",
            Self::StoreTombstone => "STORE_TOMBSTONE",
            Self::AlreadyBootstrapped => "ALREADY_BOOTSTRAPPED",
            Self::IncompatibleVersion => "INCOMPATIBLE_VERSION",
            Self::RegionNotFound => "REGION_NOT_FOUND",
            Self::GlobalConfigNotFound => "GLOBAL_CONFIG_NOT_FOUND",
            Self::DuplicatedEntry => "DUPLICATED_ENTRY",
            Self::EntryNotFound => "ENTRY_NOT_FOUND",
            Self::InvalidValue => "INVALID_VALUE",
            Self::DataCompacted => "DATA_COMPACTED",
            Self::RegionsNotContainAllKeyRange => "REGIONS_NOT_CONTAIN_ALL_KEY_RANGE",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "OK" => Some(Self::Ok),
            "UNKNOWN" => Some(Self::Unknown),
            "NOT_BOOTSTRAPPED" => Some(Self::NotBootstrapped),
            "STORE_TOMBSTONE" => Some(Self::StoreTombstone),
            "ALREADY_BOOTSTRAPPED" => Some(Self::AlreadyBootstrapped),
            "INCOMPATIBLE_VERSION" => Some(Self::IncompatibleVersion),
            "REGION_NOT_FOUND" => Some(Self::RegionNotFound),
            "GLOBAL_CONFIG_NOT_FOUND" => Some(Self::GlobalConfigNotFound),
            "DUPLICATED_ENTRY" => Some(Self::DuplicatedEntry),
            "ENTRY_NOT_FOUND" => Some(Self::EntryNotFound),
            "INVALID_VALUE" => Some(Self::InvalidValue),
            "DATA_COMPACTED" => Some(Self::DataCompacted),
            "REGIONS_NOT_CONTAIN_ALL_KEY_RANGE" => Some(Self::RegionsNotContainAllKeyRange),
            _ => None,
        }
    }
}
include!("pdpb.serde.rs");
// @@protoc_insertion_point(module)
