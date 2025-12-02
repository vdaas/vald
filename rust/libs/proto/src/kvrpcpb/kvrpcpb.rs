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
pub struct RawGetRequest {
    #[prost(message, optional, tag="1")]
    pub context: ::core::option::Option<Context>,
    #[prost(bytes="vec", tag="2")]
    pub key: ::prost::alloc::vec::Vec<u8>,
    #[prost(string, tag="3")]
    pub cf: ::prost::alloc::string::String,
}
impl ::prost::Name for RawGetRequest {
const NAME: &'static str = "RawGetRequest";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawGetRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawGetRequest".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RawGetResponse {
    #[prost(message, optional, tag="1")]
    pub region_error: ::core::option::Option<super::errorpb::Error>,
    #[prost(string, tag="2")]
    pub error: ::prost::alloc::string::String,
    #[prost(bytes="vec", tag="3")]
    pub value: ::prost::alloc::vec::Vec<u8>,
    #[prost(bool, tag="4")]
    pub not_found: bool,
}
impl ::prost::Name for RawGetResponse {
const NAME: &'static str = "RawGetResponse";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawGetResponse".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawGetResponse".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RawBatchGetRequest {
    #[prost(message, optional, tag="1")]
    pub context: ::core::option::Option<Context>,
    #[prost(bytes="vec", repeated, tag="2")]
    pub keys: ::prost::alloc::vec::Vec<::prost::alloc::vec::Vec<u8>>,
    #[prost(string, tag="3")]
    pub cf: ::prost::alloc::string::String,
}
impl ::prost::Name for RawBatchGetRequest {
const NAME: &'static str = "RawBatchGetRequest";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawBatchGetRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawBatchGetRequest".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RawBatchGetResponse {
    #[prost(message, optional, tag="1")]
    pub region_error: ::core::option::Option<super::errorpb::Error>,
    #[prost(message, repeated, tag="2")]
    pub pairs: ::prost::alloc::vec::Vec<KvPair>,
}
impl ::prost::Name for RawBatchGetResponse {
const NAME: &'static str = "RawBatchGetResponse";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawBatchGetResponse".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawBatchGetResponse".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RawPutRequest {
    #[prost(message, optional, tag="1")]
    pub context: ::core::option::Option<Context>,
    #[prost(bytes="vec", tag="2")]
    pub key: ::prost::alloc::vec::Vec<u8>,
    #[prost(bytes="vec", tag="3")]
    pub value: ::prost::alloc::vec::Vec<u8>,
    #[prost(string, tag="4")]
    pub cf: ::prost::alloc::string::String,
    #[prost(uint64, tag="5")]
    pub ttl: u64,
    #[prost(bool, tag="6")]
    pub for_cas: bool,
}
impl ::prost::Name for RawPutRequest {
const NAME: &'static str = "RawPutRequest";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawPutRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawPutRequest".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RawPutResponse {
    #[prost(message, optional, tag="1")]
    pub region_error: ::core::option::Option<super::errorpb::Error>,
    #[prost(string, tag="2")]
    pub error: ::prost::alloc::string::String,
}
impl ::prost::Name for RawPutResponse {
const NAME: &'static str = "RawPutResponse";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawPutResponse".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawPutResponse".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct RawBatchPutRequest {
    #[prost(message, optional, tag="1")]
    pub context: ::core::option::Option<Context>,
    #[prost(message, repeated, tag="2")]
    pub pairs: ::prost::alloc::vec::Vec<KvPair>,
    #[prost(string, tag="3")]
    pub cf: ::prost::alloc::string::String,
    #[deprecated]
    #[prost(uint64, tag="4")]
    pub ttl: u64,
    #[prost(bool, tag="5")]
    pub for_cas: bool,
    /// The time-to-live for each keys in seconds, and if the length of `ttls`
    /// is exactly one, the ttl will be applied to all keys. Otherwise, the length
    /// mismatch between `ttls` and `pairs` will return an error.
    #[prost(uint64, repeated, tag="6")]
    pub ttls: ::prost::alloc::vec::Vec<u64>,
}
impl ::prost::Name for RawBatchPutRequest {
const NAME: &'static str = "RawBatchPutRequest";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawBatchPutRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawBatchPutRequest".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RawBatchPutResponse {
    #[prost(message, optional, tag="1")]
    pub region_error: ::core::option::Option<super::errorpb::Error>,
    #[prost(string, tag="2")]
    pub error: ::prost::alloc::string::String,
}
impl ::prost::Name for RawBatchPutResponse {
const NAME: &'static str = "RawBatchPutResponse";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawBatchPutResponse".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawBatchPutResponse".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RawDeleteRequest {
    #[prost(message, optional, tag="1")]
    pub context: ::core::option::Option<Context>,
    #[prost(bytes="vec", tag="2")]
    pub key: ::prost::alloc::vec::Vec<u8>,
    #[prost(string, tag="3")]
    pub cf: ::prost::alloc::string::String,
    #[prost(bool, tag="4")]
    pub for_cas: bool,
}
impl ::prost::Name for RawDeleteRequest {
const NAME: &'static str = "RawDeleteRequest";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawDeleteRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawDeleteRequest".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RawDeleteResponse {
    #[prost(message, optional, tag="1")]
    pub region_error: ::core::option::Option<super::errorpb::Error>,
    #[prost(string, tag="2")]
    pub error: ::prost::alloc::string::String,
}
impl ::prost::Name for RawDeleteResponse {
const NAME: &'static str = "RawDeleteResponse";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawDeleteResponse".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawDeleteResponse".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RawBatchDeleteRequest {
    #[prost(message, optional, tag="1")]
    pub context: ::core::option::Option<Context>,
    #[prost(bytes="vec", repeated, tag="2")]
    pub keys: ::prost::alloc::vec::Vec<::prost::alloc::vec::Vec<u8>>,
    #[prost(string, tag="3")]
    pub cf: ::prost::alloc::string::String,
    #[prost(bool, tag="4")]
    pub for_cas: bool,
}
impl ::prost::Name for RawBatchDeleteRequest {
const NAME: &'static str = "RawBatchDeleteRequest";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawBatchDeleteRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawBatchDeleteRequest".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RawBatchDeleteResponse {
    #[prost(message, optional, tag="1")]
    pub region_error: ::core::option::Option<super::errorpb::Error>,
    #[prost(string, tag="2")]
    pub error: ::prost::alloc::string::String,
}
impl ::prost::Name for RawBatchDeleteResponse {
const NAME: &'static str = "RawBatchDeleteResponse";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.RawBatchDeleteResponse".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.RawBatchDeleteResponse".into() }}
// Helper messages.

/// Miscellaneous metadata attached to most requests.
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Context {
    #[prost(uint64, tag="1")]
    pub region_id: u64,
    /// metapb.RegionEpoch region_epoch = 2;
    /// metapb.Peer peer = 3;
    #[prost(uint64, tag="5")]
    pub term: u64,
    #[prost(enumeration="CommandPri", tag="6")]
    pub priority: i32,
    #[prost(enumeration="IsolationLevel", tag="7")]
    pub isolation_level: i32,
    #[prost(bool, tag="8")]
    pub not_fill_cache: bool,
    #[prost(bool, tag="9")]
    pub sync_log: bool,
    /// True means execution time statistics should be recorded and returned.
    #[prost(bool, tag="10")]
    pub record_time_stat: bool,
    /// True means RocksDB scan statistics should be recorded and returned.
    #[prost(bool, tag="11")]
    pub record_scan_stat: bool,
    #[prost(bool, tag="12")]
    pub replica_read: bool,
    /// Read requests can ignore locks belonging to these transactions because either
    /// these transactions are rolled back or theirs commit_ts > read request's start_ts.
    #[prost(uint64, repeated, tag="13")]
    pub resolved_locks: ::prost::alloc::vec::Vec<u64>,
    #[prost(uint64, tag="14")]
    pub max_execution_duration_ms: u64,
    /// After a region applies to `applied_index`, we can get a
    /// snapshot for the region even if the peer is a follower.
    #[prost(uint64, tag="15")]
    pub applied_index: u64,
    /// A hint for TiKV to schedule tasks more fairly. Query with same task ID
    /// may share same priority and resource quota.
    #[prost(uint64, tag="16")]
    pub task_id: u64,
    /// Not required to read the most up-to-date data, replicas with `safe_ts` >= `start_ts`
    /// can handle read request directly
    #[prost(bool, tag="17")]
    pub stale_read: bool,
    /// Any additional serialized information about the request.
    #[prost(bytes="vec", tag="18")]
    pub resource_group_tag: ::prost::alloc::vec::Vec<u8>,
    /// Used to tell TiKV whether operations are allowed or not on different disk usages.
    #[prost(enumeration="DiskFullOpt", tag="19")]
    pub disk_full_opt: i32,
    /// Indicates the request is a retry request and the same request may have been sent before.
    #[prost(bool, tag="20")]
    pub is_retry_request: bool,
    /// API version implies the encode of the key and value.
    #[prost(enumeration="ApiVersion", tag="21")]
    pub api_version: i32,
    /// Read request should read through locks belonging to these transactions because these
    /// transactions are committed and theirs commit_ts <= read request's start_ts.
    #[prost(uint64, repeated, tag="22")]
    pub committed_locks: ::prost::alloc::vec::Vec<u64>,
    // // The informantion to trace a request sent to TiKV.
    // tracepb.TraceContext trace_context = 23;

    /// The source of the request, will be used as the tag of the metrics reporting.
    /// This field can be set for any requests that require to report metrics with any extra labels.
    #[prost(string, tag="24")]
    pub request_source: ::prost::alloc::string::String,
    /// The source of the current transaction.
    #[prost(uint64, tag="25")]
    pub txn_source: u64,
    /// If `busy_threshold_ms` is given, TiKV can reject the request and return a `ServerIsBusy`
    /// error before processing if the estimated waiting duration exceeds the threshold.
    #[prost(uint32, tag="27")]
    pub busy_threshold_ms: u32,
    /// Some information used for resource control.
    #[prost(message, optional, tag="28")]
    pub resource_control_context: ::core::option::Option<ResourceControlContext>,
    /// The keyspace that the request is sent to.
    /// NOTE: This field is only meaningful while the api_version is V2.
    #[prost(string, tag="31")]
    pub keyspace_name: ::prost::alloc::string::String,
    /// The keyspace that the request is sent to.
    /// NOTE: This field is only meaningful while the api_version is V2.
    #[prost(uint32, tag="32")]
    pub keyspace_id: u32,
    /// The buckets version that the request is sent to.
    /// NOTE: This field is only meaningful while enable buckets.
    #[prost(uint64, tag="33")]
    pub buckets_version: u64,
    /// It tells us where the request comes from in TiDB. If it isn't from TiDB, leave it blank.
    /// This is for tests only and thus can be safely changed/removed without affecting compatibility.
    #[prost(message, optional, tag="34")]
    pub source_stmt: ::core::option::Option<SourceStmt>,
    /// The cluster id of the request
    #[prost(uint64, tag="35")]
    pub cluster_id: u64,
    /// The trace id of the request, will be used for tracing the request's execution's inner steps.
    #[prost(bytes="vec", tag="36")]
    pub trace_id: ::prost::alloc::vec::Vec<u8>,
    /// Control flags for trace logging behavior.
    /// Bit 0: immediate_log - Force immediate logging without buffering
    /// Bit 1: category_req_resp - Enable request/response tracing
    /// Bit 2: category_write_details - Enable detailed write tracing
    /// Bit 3: category_read_details - Enable detailed read tracing
    /// Bits 4-63: Reserved for future use
    /// This field is set by client-go based on an extractor function provided by TiDB.
    #[prost(uint64, tag="37")]
    pub trace_control_flags: u64,
}
impl ::prost::Name for Context {
const NAME: &'static str = "Context";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.Context".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.Context".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct ResourceControlContext {
    /// It's used to identify which resource group the request belongs to.
    #[prost(string, tag="1")]
    pub resource_group_name: ::prost::alloc::string::String,
    // // The resource consumption of the resource group that have completed at all TiKVs between the previous request to this TiKV and current request.
    // // It's used as penalty to make the local resource scheduling on one TiKV takes the gloabl resource consumption into consideration.
    // resource_manager.Consumption penalty = 2;

    /// This priority would override the original priority of the resource group for the request.
    /// Used to deprioritize the runaway queries.
    #[prost(uint64, tag="3")]
    pub override_priority: u64,
}
impl ::prost::Name for ResourceControlContext {
const NAME: &'static str = "ResourceControlContext";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.ResourceControlContext".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.ResourceControlContext".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct SourceStmt {
    #[prost(uint64, tag="1")]
    pub start_ts: u64,
    #[prost(uint64, tag="2")]
    pub connection_id: u64,
    #[prost(uint64, tag="3")]
    pub stmt_id: u64,
    /// session alias set by user
    #[prost(string, tag="4")]
    pub session_alias: ::prost::alloc::string::String,
}
impl ::prost::Name for SourceStmt {
const NAME: &'static str = "SourceStmt";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.SourceStmt".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.SourceStmt".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct LockInfo {
    #[prost(bytes="vec", tag="1")]
    pub primary_lock: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag="2")]
    pub lock_version: u64,
    #[prost(bytes="vec", tag="3")]
    pub key: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag="4")]
    pub lock_ttl: u64,
    /// How many keys this transaction involves in this region.
    #[prost(uint64, tag="5")]
    pub txn_size: u64,
    #[prost(enumeration="Op", tag="6")]
    pub lock_type: i32,
    #[prost(uint64, tag="7")]
    pub lock_for_update_ts: u64,
    /// Fields for transactions that are using Async Commit.
    #[prost(bool, tag="8")]
    pub use_async_commit: bool,
    #[prost(uint64, tag="9")]
    pub min_commit_ts: u64,
    #[prost(bytes="vec", repeated, tag="10")]
    pub secondaries: ::prost::alloc::vec::Vec<::prost::alloc::vec::Vec<u8>>,
    /// The time elapsed since last update of lock wait info when waiting.
    /// It's used in timeout errors. 0 means unknown or not applicable.
    /// It can be used to help the client decide whether to try resolving the lock.
    #[prost(uint64, tag="11")]
    pub duration_to_last_update_ms: u64,
    /// Reserved for file based transaction.
    #[prost(bool, tag="100")]
    pub is_txn_file: bool,
}
impl ::prost::Name for LockInfo {
const NAME: &'static str = "LockInfo";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.LockInfo".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.LockInfo".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct KeyError {
    /// Client should backoff or cleanup the lock then retry.
    #[prost(message, optional, tag="1")]
    pub locked: ::core::option::Option<LockInfo>,
    /// Client may restart the txn. e.g write conflict.
    #[prost(string, tag="2")]
    pub retryable: ::prost::alloc::string::String,
    /// Client should abort the txn.
    #[prost(string, tag="3")]
    pub abort: ::prost::alloc::string::String,
    /// Write conflict is moved from retryable to here.
    #[prost(message, optional, tag="4")]
    pub conflict: ::core::option::Option<WriteConflict>,
    /// Key already exists
    #[prost(message, optional, tag="5")]
    pub already_exist: ::core::option::Option<AlreadyExist>,
    /// Deadlock deadlock = 6; // Deadlock is used in pessimistic transaction for single statement rollback.
    ///
    /// Commit ts is earlier than min commit ts of a transaction.
    #[prost(message, optional, tag="7")]
    pub commit_ts_expired: ::core::option::Option<CommitTsExpired>,
    /// Txn not found when checking txn status.
    #[prost(message, optional, tag="8")]
    pub txn_not_found: ::core::option::Option<TxnNotFound>,
    /// Calculated commit TS exceeds the limit given by the user.
    #[prost(message, optional, tag="9")]
    pub commit_ts_too_large: ::core::option::Option<CommitTsTooLarge>,
    /// Assertion of a `Mutation` is evaluated as a failure.
    #[prost(message, optional, tag="10")]
    pub assertion_failed: ::core::option::Option<AssertionFailed>,
    /// CheckTxnStatus is sent to a lock that's not the primary.
    #[prost(message, optional, tag="11")]
    pub primary_mismatch: ::core::option::Option<PrimaryMismatch>,
    /// TxnLockNotFound indicates the txn lock is not found.
    #[prost(message, optional, tag="12")]
    pub txn_lock_not_found: ::core::option::Option<TxnLockNotFound>,
    /// Extra information for error debugging
    #[prost(message, optional, tag="100")]
    pub debug_info: ::core::option::Option<DebugInfo>,
}
impl ::prost::Name for KeyError {
const NAME: &'static str = "KeyError";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.KeyError".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.KeyError".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct WriteConflict {
    #[prost(uint64, tag="1")]
    pub start_ts: u64,
    #[prost(uint64, tag="2")]
    pub conflict_ts: u64,
    #[prost(bytes="vec", tag="3")]
    pub key: ::prost::alloc::vec::Vec<u8>,
    #[prost(bytes="vec", tag="4")]
    pub primary: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag="5")]
    pub conflict_commit_ts: u64,
    #[prost(enumeration="write_conflict::Reason", tag="6")]
    pub reason: i32,
}
/// Nested message and enum types in `WriteConflict`.
pub mod write_conflict {
    #[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
    #[repr(i32)]
    pub enum Reason {
        Unknown = 0,
        /// in optimistic transactions.
        Optimistic = 1,
        /// a lock acquisition request waits for a lock and awakes, or meets a newer version of data, let TiDB retry.
        PessimisticRetry = 2,
        /// the transaction itself has been rolled back when it tries to prewrite.
        SelfRolledBack = 3,
        /// RcCheckTs failure by meeting a newer version, let TiDB retry.
        RcCheckTs = 4,
        /// write conflict found when deferring constraint checks in pessimistic transactions. Deprecated in next-gen (cloud-storage-engine).
        LazyUniquenessCheck = 5,
        /// write conflict found on keys that do not acquire pessimistic locks in pessimistic transactions.
        NotLockedKeyConflict = 6,
    }
    impl Reason {
        /// String value of the enum field names used in the ProtoBuf definition.
        ///
        /// The values are not transformed in any way and thus are considered stable
        /// (if the ProtoBuf definition does not change) and safe for programmatic use.
        pub fn as_str_name(&self) -> &'static str {
            match self {
                Self::Unknown => "Unknown",
                Self::Optimistic => "Optimistic",
                Self::PessimisticRetry => "PessimisticRetry",
                Self::SelfRolledBack => "SelfRolledBack",
                Self::RcCheckTs => "RcCheckTs",
                Self::LazyUniquenessCheck => "LazyUniquenessCheck",
                Self::NotLockedKeyConflict => "NotLockedKeyConflict",
            }
        }
        /// Creates an enum from field names used in the ProtoBuf definition.
        pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
            match value {
                "Unknown" => Some(Self::Unknown),
                "Optimistic" => Some(Self::Optimistic),
                "PessimisticRetry" => Some(Self::PessimisticRetry),
                "SelfRolledBack" => Some(Self::SelfRolledBack),
                "RcCheckTs" => Some(Self::RcCheckTs),
                "LazyUniquenessCheck" => Some(Self::LazyUniquenessCheck),
                "NotLockedKeyConflict" => Some(Self::NotLockedKeyConflict),
                _ => None,
            }
        }
    }
}
impl ::prost::Name for WriteConflict {
const NAME: &'static str = "WriteConflict";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.WriteConflict".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.WriteConflict".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct AlreadyExist {
    #[prost(bytes="vec", tag="1")]
    pub key: ::prost::alloc::vec::Vec<u8>,
}
impl ::prost::Name for AlreadyExist {
const NAME: &'static str = "AlreadyExist";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.AlreadyExist".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.AlreadyExist".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct CommitTsExpired {
    #[prost(uint64, tag="1")]
    pub start_ts: u64,
    #[prost(uint64, tag="2")]
    pub attempted_commit_ts: u64,
    #[prost(bytes="vec", tag="3")]
    pub key: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag="4")]
    pub min_commit_ts: u64,
}
impl ::prost::Name for CommitTsExpired {
const NAME: &'static str = "CommitTsExpired";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.CommitTsExpired".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.CommitTsExpired".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct TxnNotFound {
    #[prost(uint64, tag="1")]
    pub start_ts: u64,
    #[prost(bytes="vec", tag="2")]
    pub primary_key: ::prost::alloc::vec::Vec<u8>,
}
impl ::prost::Name for TxnNotFound {
const NAME: &'static str = "TxnNotFound";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.TxnNotFound".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.TxnNotFound".into() }}
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct CommitTsTooLarge {
    /// The calculated commit TS.
    #[prost(uint64, tag="1")]
    pub commit_ts: u64,
}
impl ::prost::Name for CommitTsTooLarge {
const NAME: &'static str = "CommitTsTooLarge";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.CommitTsTooLarge".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.CommitTsTooLarge".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct AssertionFailed {
    #[prost(uint64, tag="1")]
    pub start_ts: u64,
    #[prost(bytes="vec", tag="2")]
    pub key: ::prost::alloc::vec::Vec<u8>,
    #[prost(enumeration="Assertion", tag="3")]
    pub assertion: i32,
    #[prost(uint64, tag="4")]
    pub existing_start_ts: u64,
    #[prost(uint64, tag="5")]
    pub existing_commit_ts: u64,
}
impl ::prost::Name for AssertionFailed {
const NAME: &'static str = "AssertionFailed";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.AssertionFailed".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.AssertionFailed".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct PrimaryMismatch {
    #[prost(message, optional, tag="1")]
    pub lock_info: ::core::option::Option<LockInfo>,
}
impl ::prost::Name for PrimaryMismatch {
const NAME: &'static str = "PrimaryMismatch";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.PrimaryMismatch".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.PrimaryMismatch".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct TxnLockNotFound {
    #[prost(bytes="vec", tag="1")]
    pub key: ::prost::alloc::vec::Vec<u8>,
}
impl ::prost::Name for TxnLockNotFound {
const NAME: &'static str = "TxnLockNotFound";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.TxnLockNotFound".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.TxnLockNotFound".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MvccDebugInfo {
    #[prost(bytes="vec", tag="1")]
    pub key: ::prost::alloc::vec::Vec<u8>,
    #[prost(message, optional, tag="2")]
    pub mvcc: ::core::option::Option<MvccInfo>,
}
impl ::prost::Name for MvccDebugInfo {
const NAME: &'static str = "MvccDebugInfo";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.MvccDebugInfo".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.MvccDebugInfo".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DebugInfo {
    #[prost(message, repeated, tag="1")]
    pub mvcc_info: ::prost::alloc::vec::Vec<MvccDebugInfo>,
}
impl ::prost::Name for DebugInfo {
const NAME: &'static str = "DebugInfo";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.DebugInfo".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.DebugInfo".into() }}
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct TimeDetail {
    /// Off-cpu wall time elapsed in TiKV side. Usually this includes queue waiting time and
    /// other kind of waitings in series. (Wait time in the raftstore is not included.)
    #[prost(uint64, tag="1")]
    pub wait_wall_time_ms: u64,
    /// Off-cpu and on-cpu wall time elapsed to actually process the request payload. It does not
    /// include `wait_wall_time`.
    /// This field is very close to the CPU time in most cases. Some wait time spend in RocksDB
    /// cannot be excluded for now, like Mutex wait time, which is included in this field, so that
    /// this field is called wall time instead of CPU time.
    #[prost(uint64, tag="2")]
    pub process_wall_time_ms: u64,
    /// KV read wall Time means the time used in key/value scan and get.
    #[prost(uint64, tag="3")]
    pub kv_read_wall_time_ms: u64,
    /// Total wall clock time spent on this RPC in TiKV .
    #[prost(uint64, tag="4")]
    pub total_rpc_wall_time_ns: u64,
}
impl ::prost::Name for TimeDetail {
const NAME: &'static str = "TimeDetail";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.TimeDetail".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.TimeDetail".into() }}
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct TimeDetailV2 {
    /// Off-cpu wall time elapsed in TiKV side. Usually this includes queue waiting time and
    /// other kind of waitings in series. (Wait time in the raftstore is not included.)
    #[prost(uint64, tag="1")]
    pub wait_wall_time_ns: u64,
    /// Off-cpu and on-cpu wall time elapsed to actually process the request payload. It does not
    /// include `wait_wall_time` and `suspend_wall_time`.
    /// This field is very close to the CPU time in most cases. Some wait time spend in RocksDB
    /// cannot be excluded for now, like Mutex wait time, which is included in this field, so that
    /// this field is called wall time instead of CPU time.
    #[prost(uint64, tag="2")]
    pub process_wall_time_ns: u64,
    /// Cpu wall time elapsed that task is waiting in queue.
    #[prost(uint64, tag="3")]
    pub process_suspend_wall_time_ns: u64,
    /// KV read wall Time means the time used in key/value scan and get.
    #[prost(uint64, tag="4")]
    pub kv_read_wall_time_ns: u64,
    /// Total wall clock time spent on this RPC in TiKV .
    #[prost(uint64, tag="5")]
    pub total_rpc_wall_time_ns: u64,
    /// Time spent on the gRPC layer.
    #[prost(uint64, tag="6")]
    pub kv_grpc_process_time_ns: u64,
    /// Time spent on waiting for run again in grpc pool from other executor pool.
    #[prost(uint64, tag="7")]
    pub kv_grpc_wait_time_ns: u64,
}
impl ::prost::Name for TimeDetailV2 {
const NAME: &'static str = "TimeDetailV2";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.TimeDetailV2".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.TimeDetailV2".into() }}
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct ScanInfo {
    #[prost(int64, tag="1")]
    pub total: i64,
    #[prost(int64, tag="2")]
    pub processed: i64,
    #[prost(int64, tag="3")]
    pub read_bytes: i64,
}
impl ::prost::Name for ScanInfo {
const NAME: &'static str = "ScanInfo";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.ScanInfo".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.ScanInfo".into() }}
/// Only reserved for compatibility.
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct ScanDetail {
    #[prost(message, optional, tag="1")]
    pub write: ::core::option::Option<ScanInfo>,
    #[prost(message, optional, tag="2")]
    pub lock: ::core::option::Option<ScanInfo>,
    #[prost(message, optional, tag="3")]
    pub data: ::core::option::Option<ScanInfo>,
}
impl ::prost::Name for ScanDetail {
const NAME: &'static str = "ScanDetail";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.ScanDetail".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.ScanDetail".into() }}
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct ScanDetailV2 {
    /// Number of user keys scanned from the storage.
    /// It does not include deleted version or RocksDB tombstone keys.
    /// For Coprocessor requests, it includes keys that has been filtered out by
    /// Selection.
    #[prost(uint64, tag="1")]
    pub processed_versions: u64,
    /// Number of bytes of user key-value pairs scanned from the storage, i.e.
    /// total size of data returned from MVCC layer.
    #[prost(uint64, tag="8")]
    pub processed_versions_size: u64,
    /// Approximate number of MVCC keys meet during scanning. It includes
    /// deleted versions, but does not include RocksDB tombstone keys.
    ///
    /// When this field is notably larger than `processed_versions`, it means
    /// there are a lot of deleted MVCC keys.
    #[prost(uint64, tag="2")]
    pub total_versions: u64,
    /// Total number of deletes and single deletes skipped over during
    /// iteration, i.e. how many RocksDB tombstones are skipped.
    #[prost(uint64, tag="3")]
    pub rocksdb_delete_skipped_count: u64,
    /// Total number of internal keys skipped over during iteration.
    /// See <https://github.com/facebook/rocksdb/blob/9f1c84ca471d8b1ad7be9f3eebfc2c7e07dfd7a7/include/rocksdb/perf_context.h#L84> for details.
    #[prost(uint64, tag="4")]
    pub rocksdb_key_skipped_count: u64,
    /// Total number of RocksDB block cache hits.
    #[prost(uint64, tag="5")]
    pub rocksdb_block_cache_hit_count: u64,
    /// Total number of block reads (with IO).
    #[prost(uint64, tag="6")]
    pub rocksdb_block_read_count: u64,
    /// Total number of bytes from block reads.
    #[prost(uint64, tag="7")]
    pub rocksdb_block_read_byte: u64,
    /// Total time used for block reads.
    #[prost(uint64, tag="9")]
    pub rocksdb_block_read_nanos: u64,
    /// Time used for getting a raftstore snapshot (including proposing read index, leader confirmation and getting the RocksDB snapshot).
    #[prost(uint64, tag="10")]
    pub get_snapshot_nanos: u64,
    /// Time used for proposing read index from read pool to store pool, equals 0 when performing lease read.
    #[prost(uint64, tag="11")]
    pub read_index_propose_wait_nanos: u64,
    /// Time used for leader confirmation, equals 0 when performing lease read.
    #[prost(uint64, tag="12")]
    pub read_index_confirm_wait_nanos: u64,
    /// Time used for read pool scheduling.
    #[prost(uint64, tag="13")]
    pub read_pool_schedule_wait_nanos: u64,
}
impl ::prost::Name for ScanDetailV2 {
const NAME: &'static str = "ScanDetailV2";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.ScanDetailV2".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.ScanDetailV2".into() }}
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct ExecDetails {
    /// Available when ctx.record_time_stat = true or meet slow query.
    #[prost(message, optional, tag="1")]
    pub time_detail: ::core::option::Option<TimeDetail>,
    /// Available when ctx.record_scan_stat = true or meet slow query.
    #[prost(message, optional, tag="2")]
    pub scan_detail: ::core::option::Option<ScanDetail>,
}
impl ::prost::Name for ExecDetails {
const NAME: &'static str = "ExecDetails";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.ExecDetails".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.ExecDetails".into() }}
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct ExecDetailsV2 {
    /// Available when ctx.record_time_stat = true or meet slow query.
    /// deprecated. Should use `time_detail_v2` instead.
    #[prost(message, optional, tag="1")]
    pub time_detail: ::core::option::Option<TimeDetail>,
    /// Available when ctx.record_scan_stat = true or meet slow query.
    #[prost(message, optional, tag="2")]
    pub scan_detail_v2: ::core::option::Option<ScanDetailV2>,
    /// Raftstore writing durations of the request. Only available for some write requests.
    #[prost(message, optional, tag="3")]
    pub write_detail: ::core::option::Option<WriteDetail>,
    /// Available when ctx.record_time_stat = true or meet slow query.
    #[prost(message, optional, tag="4")]
    pub time_detail_v2: ::core::option::Option<TimeDetailV2>,
}
impl ::prost::Name for ExecDetailsV2 {
const NAME: &'static str = "ExecDetailsV2";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.ExecDetailsV2".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.ExecDetailsV2".into() }}
#[derive(Clone, Copy, PartialEq, Eq, Hash, ::prost::Message)]
pub struct WriteDetail {
    /// Wait duration in the store loop.
    #[prost(uint64, tag="1")]
    pub store_batch_wait_nanos: u64,
    /// Wait duration before sending proposal to peers.
    #[prost(uint64, tag="2")]
    pub propose_send_wait_nanos: u64,
    /// Total time spent on persisting the log.
    #[prost(uint64, tag="3")]
    pub persist_log_nanos: u64,
    /// Wait time until the Raft log write leader begins to write.
    #[prost(uint64, tag="4")]
    pub raft_db_write_leader_wait_nanos: u64,
    /// Time spent on synchronizing the Raft log to the disk.
    #[prost(uint64, tag="5")]
    pub raft_db_sync_log_nanos: u64,
    /// Time spent on writing the Raft log to the Raft memtable.
    #[prost(uint64, tag="6")]
    pub raft_db_write_memtable_nanos: u64,
    /// Time waiting for peers to confirm the proposal (counting from the instant when the leader sends the proposal message).
    #[prost(uint64, tag="7")]
    pub commit_log_nanos: u64,
    /// Wait duration in the apply loop.
    #[prost(uint64, tag="8")]
    pub apply_batch_wait_nanos: u64,
    /// Total time spend to applying the log.
    #[prost(uint64, tag="9")]
    pub apply_log_nanos: u64,
    /// Wait time until the KV RocksDB lock is acquired.
    #[prost(uint64, tag="10")]
    pub apply_mutex_lock_nanos: u64,
    /// Wait time until becoming the KV RocksDB write leader.
    #[prost(uint64, tag="11")]
    pub apply_write_leader_wait_nanos: u64,
    /// Time spent on writing the KV DB WAL to the disk.
    #[prost(uint64, tag="12")]
    pub apply_write_wal_nanos: u64,
    /// Time spent on writing to the memtable of the KV RocksDB.
    #[prost(uint64, tag="13")]
    pub apply_write_memtable_nanos: u64,
    /// Time spent on waiting in the latch.
    #[prost(uint64, tag="14")]
    pub latch_wait_nanos: u64,
    /// Processing time in the transaction layer.
    #[prost(uint64, tag="15")]
    pub process_nanos: u64,
    /// Wait time because of the scheduler flow control or quota limiter throttling.
    #[prost(uint64, tag="16")]
    pub throttle_nanos: u64,
    /// Wait time in the waiter manager for pessimistic locking.
    #[prost(uint64, tag="17")]
    pub pessimistic_lock_wait_nanos: u64,
}
impl ::prost::Name for WriteDetail {
const NAME: &'static str = "WriteDetail";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.WriteDetail".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.WriteDetail".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct KvPair {
    #[prost(message, optional, tag="1")]
    pub error: ::core::option::Option<KeyError>,
    #[prost(bytes="vec", tag="2")]
    pub key: ::prost::alloc::vec::Vec<u8>,
    #[prost(bytes="vec", tag="3")]
    pub value: ::prost::alloc::vec::Vec<u8>,
    /// The commit timestamp of the key.
    /// If it is zero, it means the commit timestamp is unknown.
    #[prost(uint64, tag="4")]
    pub commit_ts: u64,
}
impl ::prost::Name for KvPair {
const NAME: &'static str = "KvPair";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.KvPair".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.KvPair".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct MvccWrite {
    #[prost(enumeration="Op", tag="1")]
    pub r#type: i32,
    #[prost(uint64, tag="2")]
    pub start_ts: u64,
    #[prost(uint64, tag="3")]
    pub commit_ts: u64,
    #[prost(bytes="vec", tag="4")]
    pub short_value: ::prost::alloc::vec::Vec<u8>,
    #[prost(bool, tag="5")]
    pub has_overlapped_rollback: bool,
    #[prost(bool, tag="6")]
    pub has_gc_fence: bool,
    #[prost(uint64, tag="7")]
    pub gc_fence: u64,
    #[prost(uint64, tag="8")]
    pub last_change_ts: u64,
    #[prost(uint64, tag="9")]
    pub versions_to_last_change: u64,
}
impl ::prost::Name for MvccWrite {
const NAME: &'static str = "MvccWrite";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.MvccWrite".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.MvccWrite".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct MvccValue {
    #[prost(uint64, tag="1")]
    pub start_ts: u64,
    #[prost(bytes="vec", tag="2")]
    pub value: ::prost::alloc::vec::Vec<u8>,
}
impl ::prost::Name for MvccValue {
const NAME: &'static str = "MvccValue";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.MvccValue".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.MvccValue".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct MvccLock {
    #[prost(enumeration="Op", tag="1")]
    pub r#type: i32,
    #[prost(uint64, tag="2")]
    pub start_ts: u64,
    #[prost(bytes="vec", tag="3")]
    pub primary: ::prost::alloc::vec::Vec<u8>,
    #[prost(bytes="vec", tag="4")]
    pub short_value: ::prost::alloc::vec::Vec<u8>,
    #[prost(uint64, tag="5")]
    pub ttl: u64,
    #[prost(uint64, tag="6")]
    pub for_update_ts: u64,
    #[prost(uint64, tag="7")]
    pub txn_size: u64,
    #[prost(bool, tag="8")]
    pub use_async_commit: bool,
    #[prost(bytes="vec", repeated, tag="9")]
    pub secondaries: ::prost::alloc::vec::Vec<::prost::alloc::vec::Vec<u8>>,
    #[prost(uint64, repeated, tag="10")]
    pub rollback_ts: ::prost::alloc::vec::Vec<u64>,
    #[prost(uint64, tag="11")]
    pub last_change_ts: u64,
    #[prost(uint64, tag="12")]
    pub versions_to_last_change: u64,
}
impl ::prost::Name for MvccLock {
const NAME: &'static str = "MvccLock";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.MvccLock".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.MvccLock".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MvccInfo {
    #[prost(message, optional, tag="1")]
    pub lock: ::core::option::Option<MvccLock>,
    #[prost(message, repeated, tag="2")]
    pub writes: ::prost::alloc::vec::Vec<MvccWrite>,
    #[prost(message, repeated, tag="3")]
    pub values: ::prost::alloc::vec::Vec<MvccValue>,
}
impl ::prost::Name for MvccInfo {
const NAME: &'static str = "MvccInfo";
const PACKAGE: &'static str = "kvrpcpb";
fn full_name() -> ::prost::alloc::string::String { "kvrpcpb.MvccInfo".into() }fn type_url() -> ::prost::alloc::string::String { "/kvrpcpb.MvccInfo".into() }}
/// The API version the server and the client is using.
/// See more details in <https://github.com/tikv/rfcs/blob/master/text/0069-api-v2.md.>
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum ApiVersion {
    /// `V1` is mainly for TiDB & TxnKV, and is not safe to use RawKV along with the others.
    /// V1 server only accepts V1 requests. V1 raw requests with TTL will be rejected.
    V1 = 0,
    ///
    /// `V1TTL` is only available to RawKV, and 8 bytes representing the unix timestamp in
    /// seconds for expiring time will be append to the value of all RawKV entries. For example:
    /// ------------------------------------------------------------
    /// | User value     | Expire Ts                               |
    /// ------------------------------------------------------------
    /// | 0x12 0x34 0x56 | 0x00 0x00 0x00 0x00 0x00 0x00 0xff 0xff |
    /// ------------------------------------------------------------
    /// V1TTL server only accepts V1 raw requests.
    /// V1 client should not use `V1TTL` in request. V1 client should always send `V1`.
    V1ttl = 1,
    ///
    /// `V2` use new encoding for RawKV & TxnKV to support more features.
    ///
    /// Key Encoding:
    ///   TiDB: start with `m` or `t`, the same as `V1`.
    ///   TxnKV: prefix with `x`, encoded as `MCE( x{keyspace id} + {user key} ) + timestamp`.
    ///   RawKV: prefix with `r`, encoded as `MCE( r{keyspace id} + {user key} ) + timestamp`.
    ///   Where the `{keyspace id}` is fixed-length of 3 bytes in network byte order.
    ///   Besides, RawKV entires must be in `default` CF.
    ///
    /// Value Encoding:
    ///   TiDB & TxnKV: the same as `V1`.
    ///   RawKV: `{user value} + {optional fields} + {meta flag}`. The last byte in the
    ///   raw value must be meta flags. For example:
    ///   --------------------------------------
    ///   | User value     | Meta flags        |
    ///   --------------------------------------
    ///   | 0x12 0x34 0x56 | 0x00 (0b00000000) |
    ///   --------------------------------------
    ///   Bit 0 of meta flags is for TTL. If set, the value contains 8 bytes expiring time as
    ///   unix timestamp in seconds at the very left to the meta flags.
    ///   --------------------------------------------------------------------------------
    ///   | User value     | Expiring time                           | Meta flags        |
    ///   --------------------------------------------------------------------------------
    ///   | 0x12 0x34 0x56 | 0x00 0x00 0x00 0x00 0x00 0x00 0xff 0xff | 0x01 (0b00000001) |
    ///   --------------------------------------------------------------------------------
    ///   Bit 1 is for deletion. If set, the entry is logical deleted.
    ///   ---------------------
    ///   | Meta flags        |
    ///   ---------------------
    ///   | 0x02 (0b00000010) |
    ///   ---------------------
    ///
    /// V2 server accpets V2 requests and V1 transactional requests that statrts with TiDB key
    /// prefix (`m` and `t`).
    V2 = 2,
}
impl ApiVersion {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            Self::V1 => "V1",
            Self::V1ttl => "V1TTL",
            Self::V2 => "V2",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "V1" => Some(Self::V1),
            "V1TTL" => Some(Self::V1ttl),
            "V2" => Some(Self::V2),
            _ => None,
        }
    }
}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum CommandPri {
    /// Normal is the default value.
    Normal = 0,
    Low = 1,
    High = 2,
}
impl CommandPri {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            Self::Normal => "Normal",
            Self::Low => "Low",
            Self::High => "High",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "Normal" => Some(Self::Normal),
            "Low" => Some(Self::Low),
            "High" => Some(Self::High),
            _ => None,
        }
    }
}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum IsolationLevel {
    /// SI = snapshot isolation
    Si = 0,
    /// RC = read committed
    Rc = 1,
    /// RC read and it's needed to check if there exists more recent versions.
    RcCheckTs = 2,
}
impl IsolationLevel {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            Self::Si => "SI",
            Self::Rc => "RC",
            Self::RcCheckTs => "RCCheckTS",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "SI" => Some(Self::Si),
            "RC" => Some(Self::Rc),
            "RCCheckTS" => Some(Self::RcCheckTs),
            _ => None,
        }
    }
}
/// Operation allowed info during each TiKV storage threshold.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum DiskFullOpt {
    /// The default value, means operations are not allowed either under almost full or already full.
    NotAllowedOnFull = 0,
    /// Means operations will be allowed when disk is almost full.
    AllowedOnAlmostFull = 1,
    /// Means operations will be allowed when disk is already full.
    AllowedOnAlreadyFull = 2,
}
impl DiskFullOpt {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            Self::NotAllowedOnFull => "NotAllowedOnFull",
            Self::AllowedOnAlmostFull => "AllowedOnAlmostFull",
            Self::AllowedOnAlreadyFull => "AllowedOnAlreadyFull",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "NotAllowedOnFull" => Some(Self::NotAllowedOnFull),
            "AllowedOnAlmostFull" => Some(Self::AllowedOnAlmostFull),
            "AllowedOnAlreadyFull" => Some(Self::AllowedOnAlreadyFull),
            _ => None,
        }
    }
}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum Op {
    Put = 0,
    Del = 1,
    Lock = 2,
    Rollback = 3,
    /// insert operation has a constraint that key should not exist before.
    Insert = 4,
    PessimisticLock = 5,
    CheckNotExists = 6,
}
impl Op {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            Self::Put => "Put",
            Self::Del => "Del",
            Self::Lock => "Lock",
            Self::Rollback => "Rollback",
            Self::Insert => "Insert",
            Self::PessimisticLock => "PessimisticLock",
            Self::CheckNotExists => "CheckNotExists",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "Put" => Some(Self::Put),
            "Del" => Some(Self::Del),
            "Lock" => Some(Self::Lock),
            "Rollback" => Some(Self::Rollback),
            "Insert" => Some(Self::Insert),
            "PessimisticLock" => Some(Self::PessimisticLock),
            "CheckNotExists" => Some(Self::CheckNotExists),
            _ => None,
        }
    }
}
#[derive(Clone, Copy, Debug, PartialEq, Eq, Hash, PartialOrd, Ord, ::prost::Enumeration)]
#[repr(i32)]
pub enum Assertion {
    None = 0,
    Exist = 1,
    NotExist = 2,
}
impl Assertion {
    /// String value of the enum field names used in the ProtoBuf definition.
    ///
    /// The values are not transformed in any way and thus are considered stable
    /// (if the ProtoBuf definition does not change) and safe for programmatic use.
    pub fn as_str_name(&self) -> &'static str {
        match self {
            Self::None => "None",
            Self::Exist => "Exist",
            Self::NotExist => "NotExist",
        }
    }
    /// Creates an enum from field names used in the ProtoBuf definition.
    pub fn from_str_name(value: &str) -> ::core::option::Option<Self> {
        match value {
            "None" => Some(Self::None),
            "Exist" => Some(Self::Exist),
            "NotExist" => Some(Self::NotExist),
            _ => None,
        }
    }
}
// @@protoc_insertion_point(module)
