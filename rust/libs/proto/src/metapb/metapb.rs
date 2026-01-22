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
pub struct Peer {
    #[prost(uint64, tag="1")]
    pub id: u64,
    #[prost(uint64, tag="2")]
    pub store_id: u64,
    /// PeerRole role = 3;
    #[prost(bool, tag="4")]
    pub is_witness: bool,
}
impl ::prost::Name for Peer {
const NAME: &'static str = "Peer";
const PACKAGE: &'static str = "metapb";
fn full_name() -> ::prost::alloc::string::String { "metapb.Peer".into() }fn type_url() -> ::prost::alloc::string::String { "/metapb.Peer".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Store {
    #[prost(uint64, tag="1")]
    pub id: u64,
    /// Address to handle client requests (kv, cop, etc.)
    #[prost(string, tag="2")]
    pub address: ::prost::alloc::string::String,
    #[prost(enumeration="StoreState", tag="3")]
    pub state: i32,
    /// repeated StoreLabel labels = 4;
    #[prost(string, tag="5")]
    pub version: ::prost::alloc::string::String,
    /// Address to handle peer requests (raft messages from other store).
    /// Empty means same as address.
    #[prost(string, tag="6")]
    pub peer_address: ::prost::alloc::string::String,
    /// Status address provides the HTTP service for external components
    #[prost(string, tag="7")]
    pub status_address: ::prost::alloc::string::String,
    #[prost(string, tag="8")]
    pub git_hash: ::prost::alloc::string::String,
    /// The start timestamp of the current store
    #[prost(int64, tag="9")]
    pub start_timestamp: i64,
    #[prost(string, tag="10")]
    pub deploy_path: ::prost::alloc::string::String,
    /// The last heartbeat timestamp of the store.
    #[prost(int64, tag="11")]
    pub last_heartbeat: i64,
    /// If the store is physically destroyed, which means it can never up again.
    ///
    /// NodeState is used to replace StoreState which will be deprecated in the future.
    /// NodeState node_state = 13;
    #[prost(bool, tag="12")]
    pub physically_destroyed: bool,
}
impl ::prost::Name for Store {
const NAME: &'static str = "Store";
const PACKAGE: &'static str = "metapb";
fn full_name() -> ::prost::alloc::string::String { "metapb.Store".into() }fn type_url() -> ::prost::alloc::string::String { "/metapb.Store".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Region2 {
    #[prost(uint64, tag="1")]
    pub id: u64,
    /// Region key range [start_key, end_key).
    #[prost(bytes="vec", tag="2")]
    pub start_key: ::prost::alloc::vec::Vec<u8>,
    #[prost(bytes="vec", tag="3")]
    pub end_key: ::prost::alloc::vec::Vec<u8>,
    #[prost(message, optional, tag="4")]
    pub region_epoch: ::core::option::Option<RegionEpoch>,
    #[prost(message, repeated, tag="5")]
    pub peers: ::prost::alloc::vec::Vec<Peer>,
    /// Encryption metadata for start_key and end_key. encryption_meta.iv is IV for start_key.
    /// IV for end_key is calculated from (encryption_meta.iv + len(start_key)).
    /// The field is only used by PD and should be ignored otherwise.
    /// If encryption_meta is empty (i.e. nil), it means start_key and end_key are unencrypted.
    /// encryptionpb.EncryptionMeta encryption_meta = 6;
    /// The flashback state indicates whether this region is in the flashback state.
    /// TODO: only check by `flashback_start_ts` in the future. Keep for compatibility now.
    #[prost(bool, tag="7")]
    pub is_in_flashback: bool,
    /// The start_ts that the current flashback progress is using.
    #[prost(uint64, tag="8")]
    pub flashback_start_ts: u64,
}
impl ::prost::Name for Region2 {
const NAME: &'static str = "Region2";
const PACKAGE: &'static str = "metapb";
fn full_name() -> ::prost::alloc::string::String { "metapb.Region2".into() }fn type_url() -> ::prost::alloc::string::String { "/metapb.Region2".into() }}
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RegionEpoch {
    /// Conf change version, auto increment when add or remove peer
    #[prost(uint64, tag="1")]
    pub conf_ver: u64,
    /// Region version, auto increment when split or merge
    #[prost(uint64, tag="2")]
    pub version: u64,
}
impl ::prost::Name for RegionEpoch {
const NAME: &'static str = "RegionEpoch";
const PACKAGE: &'static str = "metapb";
fn full_name() -> ::prost::alloc::string::String { "metapb.RegionEpoch".into() }fn type_url() -> ::prost::alloc::string::String { "/metapb.RegionEpoch".into() }}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum StoreState {
    Up = 0,
    Offline = 1,
    Tombstone = 2,
}
impl StoreState {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            Self::Up => "Up",
            Self::Offline => "Offline",
            Self::Tombstone => "Tombstone",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "Up" => Some(Self::Up),
            "Offline" => Some(Self::Offline),
            "Tombstone" => Some(Self::Tombstone),
            _ => None,
        }
    }
}
include!("metapb.serde.rs");
// @@protoc_insertion_point(module)
