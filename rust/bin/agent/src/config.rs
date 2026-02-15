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

use serde::{Deserialize, Serialize};
use std::env;
use std::path::Path;

/// AgentConfig represents the global configuration for the agent
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct AgentConfig {
    #[serde(default)]
    pub logging: Logging,

    #[serde(default)]
    pub observability: Observability,

    #[serde(default)]
    pub server_config: ServerConfig,

    #[serde(default)]
    pub service: Service,

    #[serde(default)]
    pub daemon: Daemon,

    #[serde(default)]
    pub qbg: QBG,
}

impl AgentConfig {
    pub fn bind(&mut self) -> &mut Self {
        self.qbg.bind();
        self
    }

    pub fn validate(&self) -> Result<(), String> {
        self.qbg.validate()?;
        Ok(())
    }
}

/// Logging configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Logging {
    #[serde(default = "default_logging_level")]
    pub level: String,

    #[serde(default)]
    pub json: bool,
}

fn default_logging_level() -> String {
    "info".to_string()
}

impl Default for Logging {
    fn default() -> Self {
        Self {
            level: default_logging_level(),
            json: false,
        }
    }
}

/// Observability configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Observability {
    #[serde(default)]
    pub enabled: bool,

    #[serde(default)]
    pub endpoint: String,

    #[serde(default = "default_service_name")]
    pub service_name: String,

    #[serde(default)]
    pub tracer: Tracer,

    #[serde(default)]
    pub meter: Meter,
}

fn default_service_name() -> String {
    "vald-agent".to_string()
}

impl Default for Observability {
    fn default() -> Self {
        Self {
            enabled: false,
            endpoint: String::new(),
            service_name: default_service_name(),
            tracer: Tracer::default(),
            meter: Meter::default(),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct Tracer {
    #[serde(default)]
    pub enabled: bool,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Meter {
    #[serde(default)]
    pub enabled: bool,

    #[serde(default = "default_meter_export_duration_secs")]
    pub export_duration_secs: u64,

    #[serde(default = "default_meter_export_timeout_secs")]
    pub export_timeout_secs: u64,
}

fn default_meter_export_duration_secs() -> u64 {
    1
}

fn default_meter_export_timeout_secs() -> u64 {
    5
}

impl Default for Meter {
    fn default() -> Self {
        Self {
            enabled: false,
            export_duration_secs: default_meter_export_duration_secs(),
            export_timeout_secs: default_meter_export_timeout_secs(),
        }
    }
}

/// Server configuration
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct ServerConfig {
    #[serde(default)]
    pub servers: Vec<Server>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
#[derive(Default)]
pub struct Server {
    #[serde(default)]
    pub name: String,

    #[serde(default)]
    pub host: String,

    #[serde(default)]
    pub port: u16,

    #[serde(default)]
    pub grpc: GrpcServerConfig,
}


#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GrpcServerConfig {
    #[serde(default)]
    pub max_receive_message_size: usize,

    #[serde(default)]
    pub max_send_message_size: usize,

    #[serde(default)]
    pub initial_window_size: u32,

    #[serde(default)]
    pub initial_conn_window_size: u32,

    #[serde(default)]
    pub max_header_list_size: u32,

    #[serde(default)]
    pub max_concurrent_streams: u32,

    #[serde(default)]
    pub connection_timeout: String,

    #[serde(default)]
    pub keepalive: Keepalive,

    #[serde(default)]
    pub interceptors: Vec<String>,
}

impl Default for GrpcServerConfig {
    fn default() -> Self {
        Self {
            max_receive_message_size: 4 * 1024 * 1024,
            max_send_message_size: 4 * 1024 * 1024,
            initial_window_size: 65535,
            initial_conn_window_size: 65535,
            max_header_list_size: 8192,
            max_concurrent_streams: 100,
            connection_timeout: String::new(),
            keepalive: Keepalive::default(),
            interceptors: Vec::new(),
        }
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct Keepalive {
    #[serde(default)]
    pub max_conn_age: String,

    #[serde(default)]
    pub time: String,

    #[serde(default)]
    pub timeout: String,
}

/// Service configuration
#[derive(Debug, Clone, Serialize, Deserialize, Default)]
pub struct Service {
    #[serde(rename = "type")]
    #[serde(default)]
    pub type_: String,
}

/// Daemon configuration
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Daemon {
    #[serde(default = "default_daemon_auto_index_check_duration_ms")]
    pub auto_index_check_duration_ms: u64,

    #[serde(default = "default_daemon_auto_save_index_duration_ms")]
    pub auto_save_index_duration_ms: u64,

    #[serde(default = "default_daemon_auto_index_limit_ms")]
    pub auto_index_limit_ms: u64,

    #[serde(default = "default_daemon_auto_index_length")]
    pub auto_index_length: usize,

    #[serde(default = "default_daemon_pool_size")]
    pub pool_size: u32,

    #[serde(default = "default_daemon_initial_delay_ms")]
    pub initial_delay_ms: u64,

    #[serde(default)]
    pub enable_proactive_gc: bool,
}

fn default_daemon_auto_index_check_duration_ms() -> u64 {
    1000
}

fn default_daemon_auto_save_index_duration_ms() -> u64 {
    60000
}

fn default_daemon_auto_index_limit_ms() -> u64 {
    3600000
}

fn default_daemon_auto_index_length() -> usize {
    100
}

fn default_daemon_pool_size() -> u32 {
    10000
}

fn default_daemon_initial_delay_ms() -> u64 {
    0
}

impl Default for Daemon {
    fn default() -> Self {
        Self {
            auto_index_check_duration_ms: default_daemon_auto_index_check_duration_ms(),
            auto_save_index_duration_ms: default_daemon_auto_save_index_duration_ms(),
            auto_index_limit_ms: default_daemon_auto_index_limit_ms(),
            auto_index_length: default_daemon_auto_index_length(),
            pool_size: default_daemon_pool_size(),
            initial_delay_ms: default_daemon_initial_delay_ms(),
            enable_proactive_gc: false,
        }
    }
}

/// VQueue configuration for vector queue buffer sizes
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VQueue {
    // ... existing code ...
    /// InsertBufferPoolSize represents insert time ordered slice buffer size
    #[serde(default = "default_insert_buffer_pool_size")]
    pub insert_buffer_pool_size: usize,

    /// DeleteBufferPoolSize represents delete time ordered slice buffer size
    #[serde(default = "default_delete_buffer_pool_size")]
    pub delete_buffer_pool_size: usize,
}

fn default_insert_buffer_pool_size() -> usize {
    1000
}

fn default_delete_buffer_pool_size() -> usize {
    1000
}

impl VQueue {
    pub fn new() -> Self {
        Self {
            insert_buffer_pool_size: default_insert_buffer_pool_size(),
            delete_buffer_pool_size: default_delete_buffer_pool_size(),
        }
    }

    pub fn bind(&mut self) -> &mut Self {
        self
    }
}

impl Default for VQueue {
    fn default() -> Self {
        Self::new()
    }
}

/// KVSDB configuration for bidirectional kv store
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct KVSDB {
    /// Concurrency represents kvsdb range loop processing concurrency
    #[serde(default = "default_kvsdb_concurrency")]
    pub concurrency: usize,

    /// CacheCapacity represents kvsdb cache capacity
    #[serde(default = "default_kvsdb_cache_capacity")]
    pub cache_capacity: usize,

    /// CompressionFactor represents kvsdb compression factor
    #[serde(default = "default_kvsdb_compression_factor")]
    pub compression_factor: i32,

    /// UseCompression represents kvsdb compression usage
    #[serde(default = "default_kvsdb_use_compression")]
    pub use_compression: bool,
}

fn default_kvsdb_concurrency() -> usize {
    10
}

fn default_kvsdb_cache_capacity() -> usize {
    10000
}

fn default_kvsdb_compression_factor() -> i32 {
    9
}

fn default_kvsdb_use_compression() -> bool {
    true
}

impl KVSDB {
    pub fn new() -> Self {
        Self {
            concurrency: default_kvsdb_concurrency(),
            cache_capacity: default_kvsdb_cache_capacity(),
            compression_factor: default_kvsdb_compression_factor(),
            use_compression: default_kvsdb_use_compression(),
        }
    }

    pub fn bind(&mut self) -> &mut Self {
        self
    }
}

impl Default for KVSDB {
    fn default() -> Self {
        Self::new()
    }
}

/// QBG configuration structure
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct QBG {
    /// PodName represent the pod name
    #[serde(default)]
    pub pod_name: String,

    /// PodNamespace represent the pod namespace
    #[serde(default)]
    pub namespace: String,

    /// IndexPath represent the qbg index file path
    #[serde(default)]
    pub index_path: String,

    /// Dimension represent the qbg index dimension
    #[serde(default)]
    pub dimension: usize,

    /// ExtendedDimension represent the qbg extended dimension
    #[serde(default)]
    pub extended_dimension: usize,

    /// NumberOfSubvectors represent the number of subvectors
    #[serde(default = "default_number_of_subvectors")]
    pub number_of_subvectors: usize,

    /// NumberOfBlobs represent the number of blobs
    #[serde(default)]
    pub number_of_blobs: usize,

    /// InternalDataType represent the internal data type (1 for float32, 2 for uint8)
    #[serde(default = "default_internal_data_type")]
    pub internal_data_type: i32,

    /// DataType represent the data type (1 for float32, 2 for uint8)
    #[serde(default = "default_data_type")]
    pub data_type: i32,

    /// DistanceType represent the distance type
    #[serde(default = "default_distance_type")]
    pub distance_type: i32,

    /// HierarchicalClusteringInitMode represent hierarchical clustering init mode
    #[serde(default = "default_hierarchical_clustering_init_mode")]
    pub hierarchical_clustering_init_mode: i32,

    /// NumberOfFirstObjects represent number of first objects
    #[serde(default)]
    pub number_of_first_objects: usize,

    /// NumberOfFirstClusters represent number of first clusters
    #[serde(default)]
    pub number_of_first_clusters: usize,

    /// NumberOfSecondObjects represent number of second objects
    #[serde(default)]
    pub number_of_second_objects: usize,

    /// NumberOfSecondClusters represent number of second clusters
    #[serde(default)]
    pub number_of_second_clusters: usize,

    /// NumberOfThirdClusters represent number of third clusters
    #[serde(default)]
    pub number_of_third_clusters: usize,

    /// NumberOfObjects represent total number of objects
    #[serde(default = "default_number_of_objects")]
    pub number_of_objects: usize,

    /// OptimizationClusteringInitMode represent optimization clustering init mode
    #[serde(default = "default_optimization_clustering_init_mode")]
    pub optimization_clustering_init_mode: i32,

    /// RotationIteration represent rotation iteration count
    #[serde(default = "default_rotation_iteration")]
    pub rotation_iteration: usize,

    /// SubvectorIteration represent subvector iteration count
    #[serde(default = "default_subvector_iteration")]
    pub subvector_iteration: usize,

    /// NumberOfMatrices represent number of matrices
    #[serde(default = "default_number_of_matrices")]
    pub number_of_matrices: usize,

    /// Rotation enable rotation
    #[serde(default = "default_rotation")]
    pub rotation: bool,

    /// Repositioning enable repositioning
    #[serde(default)]
    pub repositioning: bool,

    /// BulkInsertChunkSize represent the bulk insert chunk size
    #[serde(default = "default_bulk_insert_chunk_size")]
    pub bulk_insert_chunk_size: usize,

    /// DefaultPoolSize represent default create index batch pool size
    #[serde(default = "default_pool_size")]
    pub default_pool_size: u32,

    /// DefaultRadius represent default radius used for search
    #[serde(default = "default_radius")]
    pub default_radius: f32,

    /// DefaultEpsilon represent default epsilon used for search
    #[serde(default = "default_epsilon")]
    pub default_epsilon: f32,

    /// AutoIndexDurationLimit represents auto indexing duration limit
    #[serde(default)]
    pub auto_index_duration_limit: String,

    /// AutoIndexCheckDuration represent checking loop duration about auto indexing execution
    #[serde(default)]
    pub auto_index_check_duration: String,

    /// AutoSaveIndexDuration represent checking loop duration about auto save index execution
    #[serde(default)]
    pub auto_save_index_duration: String,

    /// AutoIndexLength represent auto index length limit
    #[serde(default)]
    pub auto_index_length: usize,

    /// InitialDelayMaxDuration represent maximum duration for initial delay
    #[serde(default)]
    pub initial_delay_max_duration: String,

    /// EnableInMemoryMode enables on memory qbg indexing mode
    #[serde(default)]
    pub enable_in_memory_mode: bool,

    /// EnableCopyOnWrite enables copy on write saving
    #[serde(default)]
    pub enable_copy_on_write: bool,

    /// VQueue represent the qbg vector queue buffer size
    #[serde(default)]
    pub vqueue: Option<VQueue>,

    /// KVSDB represent the qbg bidirectional kv store configuration
    #[serde(default)]
    pub kvsdb: Option<KVSDB>,

    /// BrokenIndexHistoryLimit represents the maximum number of broken index generations
    #[serde(default = "default_broken_index_history_limit")]
    pub broken_index_history_limit: usize,

    /// ErrorBufferLimit represents the maximum number of core qbg error buffer pool size limit
    #[serde(default)]
    pub error_buffer_limit: u64,

    /// IsReadReplica represents whether the qbg is read replica or not
    #[serde(default)]
    pub is_readreplica: bool,

    /// EnableExportIndexInfoToK8s represents whether the qbg index info is exported to k8s or not
    #[serde(default)]
    pub enable_export_index_info_to_k8s: bool,

    /// ExportIndexInfoDuration represents the duration of exporting index info to k8s
    #[serde(default)]
    pub export_index_info_duration: String,

    /// EnableStatistics represents whether the qbg index statistics load or not
    #[serde(default)]
    pub enable_statistics: bool,
}

// Default value functions
fn default_number_of_subvectors() -> usize {
    1
}

fn default_internal_data_type() -> i32 {
    1 // float32
}

fn default_data_type() -> i32 {
    1 // float32
}

fn default_distance_type() -> i32 {
    1 // L2
}

fn default_hierarchical_clustering_init_mode() -> i32 {
    2
}

fn default_optimization_clustering_init_mode() -> i32 {
    2
}

fn default_number_of_objects() -> usize {
    1000
}

fn default_rotation_iteration() -> usize {
    2000
}

fn default_subvector_iteration() -> usize {
    400
}

fn default_number_of_matrices() -> usize {
    3
}

fn default_rotation() -> bool {
    true
}

fn default_bulk_insert_chunk_size() -> usize {
    100
}

fn default_pool_size() -> u32 {
    10
}

fn default_radius() -> f32 {
    -1.0
}

fn default_epsilon() -> f32 {
    0.1
}

fn default_broken_index_history_limit() -> usize {
    3
}

impl QBG {
    /// Create a new QBG configuration with default values
    pub fn new() -> Self {
        Self {
            pod_name: String::new(),
            namespace: String::new(),
            index_path: String::new(),
            dimension: 0,
            extended_dimension: 0,
            number_of_subvectors: default_number_of_subvectors(),
            number_of_blobs: 0,
            internal_data_type: default_internal_data_type(),
            data_type: default_data_type(),
            distance_type: default_distance_type(),
            hierarchical_clustering_init_mode: default_hierarchical_clustering_init_mode(),
            number_of_first_objects: 0,
            number_of_first_clusters: 0,
            number_of_second_objects: 0,
            number_of_second_clusters: 0,
            number_of_third_clusters: 0,
            number_of_objects: default_number_of_objects(),
            optimization_clustering_init_mode: default_optimization_clustering_init_mode(),
            rotation_iteration: default_rotation_iteration(),
            subvector_iteration: default_subvector_iteration(),
            number_of_matrices: default_number_of_matrices(),
            rotation: default_rotation(),
            repositioning: false,
            bulk_insert_chunk_size: default_bulk_insert_chunk_size(),
            default_pool_size: default_pool_size(),
            default_radius: default_radius(),
            default_epsilon: default_epsilon(),
            auto_index_duration_limit: String::new(),
            auto_index_check_duration: String::new(),
            auto_save_index_duration: String::new(),
            auto_index_length: 0,
            initial_delay_max_duration: String::new(),
            enable_in_memory_mode: false,
            enable_copy_on_write: false,
            vqueue: None,
            kvsdb: None,
            broken_index_history_limit: default_broken_index_history_limit(),
            error_buffer_limit: 0,
            is_readreplica: false,
            enable_export_index_info_to_k8s: false,
            export_index_info_duration: String::new(),
            enable_statistics: false,
        }
    }

    /// Bind applies environment variable expansion to string fields
    pub fn bind(&mut self) -> &mut Self {
        self.pod_name = get_actual_value(&self.pod_name);
        self.namespace = get_actual_value(&self.namespace);
        self.index_path = get_actual_value(&self.index_path);
        self.auto_index_check_duration = get_actual_value(&self.auto_index_check_duration);
        self.auto_index_duration_limit = get_actual_value(&self.auto_index_duration_limit);
        self.auto_save_index_duration = get_actual_value(&self.auto_save_index_duration);
        self.initial_delay_max_duration = get_actual_value(&self.initial_delay_max_duration);
        self.export_index_info_duration = get_actual_value(&self.export_index_info_duration);

        if let Some(ref mut vq) = self.vqueue {
            vq.bind();
        } else {
            self.vqueue = Some(VQueue::default());
        }

        if let Some(ref mut kvs) = self.kvsdb {
            kvs.bind();
        } else {
            self.kvsdb = Some(KVSDB::default());
        }

        self
    }

    /// Validate configuration values
    pub fn validate(&self) -> Result<(), String> {
        if self.dimension == 0 {
            return Err("dimension must be greater than 0".to_string());
        }

        if self.index_path.is_empty() {
            return Err("index_path must not be empty".to_string());
        }

        if self.bulk_insert_chunk_size == 0 {
            return Err("bulk_insert_chunk_size must be greater than 0".to_string());
        }

        if self.number_of_subvectors == 0 {
            return Err("number_of_subvectors must be greater than 0".to_string());
        }

        // Validate data types (1 for float32, 2 for uint8)
        if !(self.internal_data_type == 1 || self.internal_data_type == 2) {
            return Err(format!(
                "invalid internal_data_type: {} (must be 1 or 2)",
                self.internal_data_type
            ));
        }

        if !(self.data_type == 1 || self.data_type == 2) {
            return Err(format!(
                "invalid data_type: {} (must be 1 or 2)",
                self.data_type
            ));
        }

        Ok(())
    }
}

impl Default for QBG {
    fn default() -> Self {
        Self::new()
    }
}

/// Get actual value by expanding environment variables
/// If value starts with ${, it attempts to resolve from environment variables
fn get_actual_value(value: &str) -> String {
    if value.starts_with("${") && value.ends_with("}") {
        let env_var = &value[2..value.len() - 1];
        if let Some(idx) = env_var.find(':') {
            let (var_name, default_val) = env_var.split_at(idx);
            env::var(var_name).unwrap_or_else(|_| default_val[1..].to_string())
        } else {
            env::var(env_var).unwrap_or_else(|_| value.to_string())
        }
    } else {
        value.to_string()
    }
}

/// Load configuration from YAML file
pub fn load_config_from_file<P: AsRef<Path>>(path: P) -> Result<QBG, Box<dyn std::error::Error>> {
    let content = std::fs::read_to_string(path)?;
    let mut config: QBG = serde_yaml::from_str(&content)?;
    config.bind();
    config.validate()?;
    Ok(config)
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::io::Write;
    use tempfile::NamedTempFile;

    #[test]
    fn test_vqueue_new() {
        let vq = VQueue::new();
        assert_eq!(vq.insert_buffer_pool_size, 1000);
        assert_eq!(vq.delete_buffer_pool_size, 1000);
    }

    #[test]
    fn test_vqueue_default() {
        let vq = VQueue::default();
        assert_eq!(vq.insert_buffer_pool_size, 1000);
        assert_eq!(vq.delete_buffer_pool_size, 1000);
    }

    #[test]
    fn test_kvsdb_new() {
        let kvs = KVSDB::new();
        assert_eq!(kvs.concurrency, 10);
        assert_eq!(kvs.cache_capacity, 10000);
        assert_eq!(kvs.compression_factor, 9);
        assert!(kvs.use_compression);
    }

    #[test]
    fn test_kvsdb_default() {
        let kvs = KVSDB::default();
        assert_eq!(kvs.concurrency, 10);
        assert_eq!(kvs.cache_capacity, 10000);
        assert_eq!(kvs.compression_factor, 9);
        assert!(kvs.use_compression);
    }

    #[test]
    fn test_qbg_new() {
        let qbg = QBG::new();
        assert_eq!(qbg.dimension, 0);
        assert_eq!(qbg.extended_dimension, 0);
        assert_eq!(qbg.number_of_subvectors, 1);
        assert_eq!(qbg.internal_data_type, 1);
        assert_eq!(qbg.data_type, 1);
        assert_eq!(qbg.distance_type, 1);
        assert_eq!(qbg.number_of_objects, 1000);
        assert_eq!(qbg.rotation_iteration, 2000);
        assert_eq!(qbg.subvector_iteration, 400);
        assert_eq!(qbg.number_of_matrices, 3);
        assert!(qbg.rotation);
        assert!(!qbg.repositioning);
        assert_eq!(qbg.bulk_insert_chunk_size, 100);
        assert_eq!(qbg.default_pool_size, 10);
        assert_eq!(qbg.default_radius, -1.0);
        assert_eq!(qbg.default_epsilon, 0.1);
        assert_eq!(qbg.broken_index_history_limit, 3);
        assert!(!qbg.enable_in_memory_mode);
        assert!(!qbg.enable_copy_on_write);
        assert!(!qbg.is_readreplica);
        assert!(!qbg.enable_export_index_info_to_k8s);
        assert!(!qbg.enable_statistics);
    }

    #[test]
    fn test_qbg_default() {
        let qbg = QBG::default();
        assert_eq!(qbg.dimension, 0);
        assert_eq!(qbg.number_of_subvectors, 1);
    }

    #[test]
    fn test_qbg_bind_with_vqueue_kvsdb() {
        let mut qbg = QBG {
            pod_name: "test-pod".to_string(),
            namespace: "test-ns".to_string(),
            index_path: "/tmp/index".to_string(),
            dimension: 128,
            vqueue: None,
            kvsdb: None,
            ..QBG::new()
        };

        qbg.bind();

        assert!(qbg.vqueue.is_some());
        assert!(qbg.kvsdb.is_some());
        assert_eq!(qbg.vqueue.as_ref().unwrap().insert_buffer_pool_size, 1000);
        assert_eq!(qbg.kvsdb.as_ref().unwrap().concurrency, 10);
    }

    #[test]
    fn test_qbg_validate_valid() {
        let qbg = QBG {
            dimension: 128,
            index_path: "/tmp/index".to_string(),
            bulk_insert_chunk_size: 100,
            number_of_subvectors: 1,
            ..QBG::new()
        };

        assert!(qbg.validate().is_ok());
    }

    #[test]
    fn test_qbg_validate_zero_dimension() {
        let qbg = QBG {
            dimension: 0,
            index_path: "/tmp/index".to_string(),
            ..QBG::new()
        };

        let result = qbg.validate();
        assert!(result.is_err());
        assert_eq!(result.unwrap_err(), "dimension must be greater than 0");
    }

    #[test]
    fn test_qbg_validate_empty_index_path() {
        let qbg = QBG {
            dimension: 128,
            index_path: String::new(),
            ..QBG::new()
        };

        let result = qbg.validate();
        assert!(result.is_err());
        assert_eq!(result.unwrap_err(), "index_path must not be empty");
    }

    #[test]
    fn test_qbg_validate_zero_bulk_insert_chunk_size() {
        let qbg = QBG {
            dimension: 128,
            index_path: "/tmp/index".to_string(),
            bulk_insert_chunk_size: 0,
            ..QBG::new()
        };

        let result = qbg.validate();
        assert!(result.is_err());
        assert_eq!(
            result.unwrap_err(),
            "bulk_insert_chunk_size must be greater than 0"
        );
    }

    #[test]
    fn test_qbg_validate_zero_number_of_subvectors() {
        let qbg = QBG {
            dimension: 128,
            index_path: "/tmp/index".to_string(),
            number_of_subvectors: 0,
            ..QBG::new()
        };

        let result = qbg.validate();
        assert!(result.is_err());
        assert_eq!(
            result.unwrap_err(),
            "number_of_subvectors must be greater than 0"
        );
    }

    #[test]
    fn test_qbg_validate_invalid_internal_data_type() {
        let qbg = QBG {
            dimension: 128,
            index_path: "/tmp/index".to_string(),
            internal_data_type: 3,
            ..QBG::new()
        };

        let result = qbg.validate();
        assert!(result.is_err());
        assert!(result.unwrap_err().contains("invalid internal_data_type"));
    }

    #[test]
    fn test_qbg_validate_invalid_data_type() {
        let qbg = QBG {
            dimension: 128,
            index_path: "/tmp/index".to_string(),
            data_type: 99,
            ..QBG::new()
        };

        let result = qbg.validate();
        assert!(result.is_err());
        assert!(result.unwrap_err().contains("invalid data_type"));
    }

    #[test]
    fn test_get_actual_value_no_env_var() {
        let value = "simple_value";
        let result = get_actual_value(value);
        assert_eq!(result, "simple_value");
    }

    #[test]
    fn test_get_actual_value_with_env_var() {
        let existing = match env::var("HOME") {
            Ok(value) => value,
            Err(_) => return,
        };
        let value = "${HOME}";
        let result = get_actual_value(value);
        assert_eq!(result, existing);
    }

    #[test]
    fn test_get_actual_value_with_env_var_and_default() {
        let value = "${NONEXISTENT_VAR:default_value}";
        let result = get_actual_value(value);
        assert_eq!(result, "default_value");
    }

    #[test]
    fn test_deserialize_from_yaml_string() {
        let yaml_str = r#"
pod_name: test-pod
namespace: test-namespace
index_path: /tmp/test_index
dimension: 256
extended_dimension: 512
number_of_subvectors: 4
number_of_blobs: 8
internal_data_type: 1
data_type: 1
distance_type: 1
bulk_insert_chunk_size: 50
rotation_iteration: 3000
subvector_iteration: 500
number_of_matrices: 4
rotation: true
repositioning: false
vqueue:
  insert_buffer_pool_size: 2000
  delete_buffer_pool_size: 2000
kvsdb:
  concurrency: 20
enable_copy_on_write: true
enable_in_memory_mode: true
is_readreplica: false
"#;

        let qbg: QBG = serde_yaml::from_str(yaml_str).expect("Failed to deserialize");
        assert_eq!(qbg.pod_name, "test-pod");
        assert_eq!(qbg.namespace, "test-namespace");
        assert_eq!(qbg.index_path, "/tmp/test_index");
        assert_eq!(qbg.dimension, 256);
        assert_eq!(qbg.extended_dimension, 512);
        assert_eq!(qbg.number_of_subvectors, 4);
        assert_eq!(qbg.number_of_blobs, 8);
        assert_eq!(qbg.internal_data_type, 1);
        assert_eq!(qbg.data_type, 1);
        assert_eq!(qbg.distance_type, 1);
        assert_eq!(qbg.bulk_insert_chunk_size, 50);
        assert_eq!(qbg.rotation_iteration, 3000);
        assert_eq!(qbg.subvector_iteration, 500);
        assert_eq!(qbg.number_of_matrices, 4);
        assert!(qbg.rotation);
        assert!(!qbg.repositioning);
        assert_eq!(qbg.vqueue.as_ref().unwrap().insert_buffer_pool_size, 2000);
        assert_eq!(qbg.vqueue.as_ref().unwrap().delete_buffer_pool_size, 2000);
        assert_eq!(qbg.kvsdb.as_ref().unwrap().concurrency, 20);
        assert!(qbg.enable_copy_on_write);
        assert!(qbg.enable_in_memory_mode);
        assert!(!qbg.is_readreplica);
    }

    #[test]
    fn test_load_config_from_file() {
        let mut file = NamedTempFile::new().expect("Failed to create temp file");
        let yaml_str = "\
index_path: /tmp/test_index
dimension: 128
";
        file.write_all(yaml_str.as_bytes())
            .expect("Failed to write config file");

        let cfg = load_config_from_file(file.path()).expect("Failed to load config");
        assert_eq!(cfg.index_path, "/tmp/test_index");
        assert_eq!(cfg.dimension, 128);
    }

    #[test]
    fn test_qbg_serialization_round_trip() {
        let qbg = QBG {
            pod_name: "test-pod".to_string(),
            namespace: "test-ns".to_string(),
            index_path: "/tmp/index".to_string(),
            dimension: 128,
            extended_dimension: 256,
            number_of_subvectors: 4,
            number_of_blobs: 8,
            vqueue: Some(VQueue {
                insert_buffer_pool_size: 2000,
                delete_buffer_pool_size: 1500,
            }),
            kvsdb: Some(KVSDB {
                concurrency: 15,
                cache_capacity: 10000,
                compression_factor: 9,
                use_compression: true,
            }),
            ..QBG::new()
        };

        let yaml_str = serde_yaml::to_string(&qbg).expect("Failed to serialize");
        let deserialized: QBG = serde_yaml::from_str(&yaml_str).expect("Failed to deserialize");

        assert_eq!(qbg.pod_name, deserialized.pod_name);
        assert_eq!(qbg.namespace, deserialized.namespace);
        assert_eq!(qbg.index_path, deserialized.index_path);
        assert_eq!(qbg.dimension, deserialized.dimension);
        assert_eq!(qbg.extended_dimension, deserialized.extended_dimension);
        assert_eq!(qbg.number_of_subvectors, deserialized.number_of_subvectors);
    }

    #[test]
    fn test_qbg_validate_data_types() {
        // Valid data types
        for dt in &[1, 2] {
            let qbg = QBG {
                dimension: 128,
                index_path: "/tmp/index".to_string(),
                data_type: *dt,
                internal_data_type: *dt,
                ..QBG::new()
            };
            assert!(qbg.validate().is_ok(), "Failed for data_type: {}", dt);
        }
    }
}
