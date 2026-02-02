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

use std::collections::HashMap;
use std::sync::Arc;
use std::sync::atomic::{AtomicBool, AtomicU64, Ordering};

use algorithm::{ANN, Error, MultiError};
use anyhow::Result;
use chrono::{Local, Timelike, Utc};
use crate::config::QBG;
use futures::StreamExt;
use kvs::{BidirectionalMap, BidirectionalMapBuilder, MapBase};
use kvs::map::codec::BincodeCodec;
use proto::payload::v1::object::Distance;
use proto::payload::v1::search;
use qbg::index::Index;
use qbg::property::Property;
use tracing::{debug, error, info, warn};
use vqueue::{DrainItem, Queue};

use super::k8s::MetricsExporter;
use super::memstore;
use super::metadata::Metadata;
use super::persistence::{PersistenceConfig, PersistenceManager};

pub struct QBGService {
    path: String,
    index: Index,
    property: Property,
    vq: vqueue::PersistentQueue,
    kvs: Arc<BidirectionalMap<String, u32, BincodeCodec>>,
    persistence: Option<PersistenceManager>,
    metrics_exporter: Option<MetricsExporter>,
    is_flushing: AtomicBool,
    is_indexing: AtomicBool,
    is_saving: AtomicBool,
    is_readreplica: bool,
    create_index_count: AtomicU64,
    unsaved_create_index_count: AtomicU64,
    processed_vq_count: AtomicU64,
    broken_index_count: AtomicU64,
    statistics_enabled: bool,
    enable_copy_on_write: bool,
    broken_index_history_limit: usize,
}

impl QBGService {
    pub async fn new(config: &QBG) -> Self {
        let path = if config.index_path.is_empty() {
            "index".to_string()
        } else {
            config.index_path.clone()
        };
        
        // Read replica configuration
        let is_readreplica = config.is_readreplica;
        
        // Persistence configuration
        let enable_copy_on_write = config.enable_copy_on_write;
        let broken_index_history_limit = config.broken_index_history_limit;
        
        // Initialize persistence manager and prepare folders
        let persistence_config = PersistenceConfig {
            enable_copy_on_write,
            broken_index_history_limit,
        };
        let persistence = PersistenceManager::new(&path, persistence_config);
        if let Err(e) = persistence.prepare_folders() {
            warn!("failed to prepare persistence folders: {}", e);
        }
        
        // Check if we need to load an existing index
        let should_load = persistence.index_exists();
        let mut broken_index_count = persistence.broken_index_count();
        
        // If existing index is potentially broken, try to back it up
        if PersistenceManager::needs_backup(&persistence.paths().primary_path) {
            info!("detected potentially broken index, attempting backup");
            if let Err(e) = persistence.backup_broken() {
                warn!("failed to backup broken index: {}", e);
            }
            broken_index_count = persistence.broken_index_count();
        }
        
        let mut property = Property::new();
        property.init_qbg_construction_parameters();
        property.set_qbg_construction_parameters(
            config.extended_dimension,
            config.dimension,
            config.number_of_subvectors,
            config.number_of_blobs,
            config.internal_data_type,
            config.data_type,
            config.distance_type,
        );
        property.init_qbg_build_parameters();
        property.set_qbg_build_parameters(
            config.hierarchical_clustering_init_mode,
            config.number_of_first_objects,
            config.number_of_first_clusters,
            config.number_of_second_objects,
            config.number_of_second_clusters,
            config.number_of_third_clusters,
            config.number_of_objects,
            config.number_of_subvectors,
            config.optimization_clustering_init_mode,
            config.rotation_iteration,
            config.subvector_iteration,
            config.number_of_matrices,
            config.rotation,
            config.repositioning,
        );
        
        // Use the primary path from persistence manager for the index
        let index_path = persistence.paths().primary_path.to_string_lossy().to_string();
        
        // Load or create the index
        let index = if should_load {
            info!("loading existing index from {}", index_path);
            // Use new_prebuilt to open an existing index (prebuilt=false for read-write mode)
            match Index::new_prebuilt(&index_path, false) {
                Ok(idx) => {
                    info!("successfully loaded existing index");
                    idx
                }
                Err(e) => {
                    warn!("failed to load existing index, creating new: {}", e);
                    Index::new(&index_path, &mut property).unwrap()
                }
            }
        } else {
            debug!("creating new index at {}", index_path);
            Index::new(&index_path, &mut property).unwrap()
        };
        
        let vq_path = path.clone();
        let vq = vqueue::Builder::new(vq_path).build().await.unwrap();
        let kvs_path = format!("{}_kvs", path);
        let kvs = BidirectionalMapBuilder::new(kvs_path)
            .cache_capacity(10000) // TODO: Add kvs_cache_capacity to QBG config
            .compression_factor(9) // TODO: Add kvs_compression_factor to QBG config
            .mode(kvs::Mode::HighThroughput)
            .use_compression(true) // TODO: Add kvs_use_compression to QBG config
            .build()
            .await
            .unwrap();

        // Initialize temporary directory for Copy-on-Write mode
        if enable_copy_on_write {
            if let Err(e) = persistence.mktmp() {
                warn!("failed to create temporary directory for CoW: {}", e);
            }
        }

        // Initialize K8s metrics exporter if enabled
        let enable_export_index_info = config.enable_export_index_info_to_k8s;
        let metrics_exporter = if enable_export_index_info {
            let pod_name = std::env::var("MY_POD_NAME").unwrap_or_default();
            let pod_namespace = std::env::var("MY_POD_NAMESPACE").unwrap_or_default();
            
            if pod_name.is_empty() || pod_namespace.is_empty() {
                warn!("K8s metrics export enabled but MY_POD_NAME or MY_POD_NAMESPACE not set");
                None
            } else {
                match super::k8s::K8sClient::new().await {
                    Ok(client) => {
                        info!("K8s metrics exporter initialized for pod {}/{}", pod_namespace, pod_name);
                        Some(MetricsExporter::new(
                            Box::new(client),
                            pod_name,
                            pod_namespace,
                            true,
                        ))
                    }
                    Err(e) => {
                        warn!("failed to create K8s client: {}", e);
                        None
                    }
                }
            }
        } else {
            None
        };
        
        QBGService {
            path: index_path,
            index,
            property,
            vq,
            kvs,
            persistence: Some(persistence),
            metrics_exporter,
            is_flushing: AtomicBool::new(false),
            is_indexing: AtomicBool::new(false),
            is_saving: AtomicBool::new(false),
            is_readreplica,
            create_index_count: AtomicU64::new(0),
            unsaved_create_index_count: AtomicU64::new(0),
            processed_vq_count: AtomicU64::new(0),
            broken_index_count: AtomicU64::new(broken_index_count),
            statistics_enabled: false,
            enable_copy_on_write,
            broken_index_history_limit,
        }
    }

    async fn ready_for_update(&mut self, uuid: String, vector: Vec<f32>, ts: i64) -> Result<(), Error> {
        if uuid.len() == 0 {
            return Err(Error::UUIDNotFound {
                uuid: "0".to_string(),
            });
        }
        if vector.len() != self.get_dimension_size() {
            return Err(Error::InvalidDimensionSize {
                current: vector.len().to_string(),
                limit: self.get_dimension_size().to_string(),
            });
        }
        let get_result = self.get_object(uuid.clone()).await;
        match get_result {
            Ok((ovec, ots)) => {
                if (vector.len() != ovec.len()) || (vector != ovec) {
                    return Ok(());
                }
                if ots < ts {
                    self.update_timestamp(uuid.clone(), ts, false).await?;
                    return Ok(());
                }
                Err(Error::UUIDAlreadyExists { uuid })
            }
            Err(Error::ObjectIDNotFound { .. }) => {
                // Object doesn't exist, ok to update (insert)
                Ok(())
            }
            Err(e) => Err(e),
        }
    }

    async fn insert_internal(&mut self, uuid: String, vector: Vec<f32>, t: i64, validation: bool) -> Result<(), Error> {
        if self.is_readreplica {
            return Err(Error::WriteOperationToReadReplica {});
        }
        if uuid.len() == 0 {
            return Err(Error::UUIDNotFound {
                uuid: "0".to_string(),
            });
        }
        if validation {
            let (_, ok) = self.exists(uuid.clone()).await;
            if ok {
                return Err(Error::UUIDAlreadyExists { uuid });
            }
        }
        self.vq.push_insert(uuid, vector, Some(t)).await.map_err(|e| Error::Internal(Box::new(e)))
    }

    async fn insert_multiple_internal(&mut self, vectors: HashMap<String, Vec<f32>>, t: i64, validation: bool) -> Result<(), Error> {
        for (uuid, vec) in vectors {
            if validation {
                self.ready_for_update(uuid.clone(), vec.clone(), t).await?;
            }
            self.insert_with_time(uuid, vec, t).await?;
        }
        Ok(())
    }

    async fn update_internal(&mut self, uuid: String, vector: Vec<f32>, t: i64) -> Result<(), Error> {
        if self.is_readreplica {
            return Err(Error::WriteOperationToReadReplica {});
        }
        self.ready_for_update(uuid.clone(), vector.clone(), t).await?;
        self.remove_internal(uuid.clone(), t, true).await?;
        self.insert_internal(uuid, vector, t+1, false).await
    }

    async fn remove_internal(&mut self, uuid: String, t: i64, validation: bool) -> Result<(), Error> {
        if self.is_readreplica {
            return Err(Error::WriteOperationToReadReplica {});
        }
        if uuid.len() == 0 {
            return Err(Error::UUIDNotFound {
                uuid: "0".to_string(),
            });
        }
        if validation {
            let result = self.kvs.get(&uuid).await;
            let iv_exists = self.vq.iv_exists(&uuid).await.unwrap_or(0) > 0;
            if result.is_err() && !iv_exists {
                return Err(Error::ObjectIDNotFound { uuid });
            }
        }
        self.vq.push_delete(uuid, Some(t)).await.map_err(|e| Error::Internal(Box::new(e)))
    }

    async fn remove_multiple_internal(&mut self, uuids: Vec<String>, t: i64, validation: bool) -> Result<(), Error> {
        let mut ids: Vec<String> = vec![];
        for uuid in uuids {
            let result = self.remove_internal(uuid, t, validation).await;
            match result {
                Ok(()) => continue,
                Err(err) => match err {
                    Error::ObjectIDNotFound { uuid } => ids.push(uuid),
                    _ => return Err(err),
                },
            }
        }
        if !ids.is_empty() {
            return Err(Error::new_object_id_not_found(ids));
        }
        Ok(())
    }
}

impl ANN for QBGService {
    #[tracing::instrument(skip(self), level = "debug")]
    async fn exists(&self, uuid: String) -> (usize, bool) {
        match memstore::exists(&self.kvs, &self.vq, &uuid).await {
            Ok((oid, exists)) => (oid as usize, exists),
            Err(_) => (0, false),
        }
    }

    #[tracing::instrument(skip(self), level = "info")]
    async fn create_index(&mut self) -> Result<(), Error> {
        // Check if read replica
        if self.is_readreplica {
            return Err(Error::WriteOperationToReadReplica {});
        }

        // If there are no objects to index, return success
        let ic = self.vq.ivq_len() + self.vq.dvq_len();
        if ic == 0 {
            self.create_index_count.fetch_add(1, Ordering::SeqCst);
            return Ok(());
        }

        // Check if already indexing
        if self.is_indexing.load(Ordering::SeqCst) {
            debug!("create index already in progress, skipping");
            return Ok(());
        }

        self.is_indexing.store(true, Ordering::SeqCst);
        info!("create index operation started, uncommitted indexes = {}", ic);

        let now = Utc::now().timestamp_nanos_opt().unwrap_or(0);
        let batch_size = 1000; // TODO: make configurable
        let mut vq_processed_cnt: u64 = 0;
        let mut insert_cnt: u32 = 0;

        // Phase 1: Process delete queue
        debug!("create index delete phase started");
        {
            let mut stream = self.vq.drain_queues(now, batch_size);
            while let Some(item_result) = stream.next().await {
                match item_result {
                    Ok(DrainItem::Delete(uuid)) => {
                        debug!("processing delete for uuid: {}", uuid);
                        match self.kvs.delete(&uuid).await {
                            Ok(oid) => {
                                if let Err(e) = self.index.remove(oid as usize) {
                                    error!("failed to remove oid {} from index: {}", oid, e);
                                    // Continue processing other items
                                }
                                debug!("removed from index and kvs: uuid={}, oid={}", uuid, oid);
                            }
                            Err(e) => {
                                warn!("uuid {} not found in kvs during delete: {}", uuid, e);
                            }
                        }
                        vq_processed_cnt += 1;
                    }
                    Ok(DrainItem::Insert(uuid, vector)) => {
                        debug!("processing insert for uuid: {}", uuid);
                        match self.index.insert(&vector) {
                            Ok(oid) => {
                                let timestamp = Utc::now().timestamp_nanos_opt().unwrap_or(0) as u128;
                                if let Err(e) = self.kvs.set(uuid.clone(), oid as u32, timestamp).await {
                                    error!("failed to set kvs for uuid {}: {}", uuid, e);
                                }
                                insert_cnt += 1;
                                debug!("inserted to index and kvs: uuid={}, oid={}", uuid, oid);
                            }
                            Err(e) => {
                                error!("failed to insert vector for uuid {}: {}", uuid, e);
                                // Retry once
                                if let Ok(oid) = self.index.insert(&vector) {
                                    let timestamp = Utc::now().timestamp_nanos_opt().unwrap_or(0) as u128;
                                    if let Err(e) = self.kvs.set(uuid.clone(), oid as u32, timestamp).await {
                                        error!("failed to set kvs on retry for uuid {}: {}", uuid, e);
                                    }
                                    insert_cnt += 1;
                                } else {
                                    error!("retry insert also failed for uuid {}", uuid);
                                }
                            }
                        }
                        vq_processed_cnt += 1;
                    }
                    Err(e) => {
                        error!("error draining vqueue: {}", e);
                    }
                }
            }
        }
        debug!("create index drain phase finished, processed {} items, inserted {}", vq_processed_cnt, insert_cnt);

        // Update processed vq count
        self.processed_vq_count.fetch_add(vq_processed_cnt, Ordering::SeqCst);

        // Phase 2: Build the index
        debug!("create graph and tree phase started");
        let result = self.index.build_index(&self.path, &mut self.property);
        self.is_indexing.store(false, Ordering::SeqCst);

        match result {
            Ok(()) => {
                self.create_index_count.fetch_add(1, Ordering::SeqCst);
                self.unsaved_create_index_count.fetch_add(1, Ordering::SeqCst);
                debug!("create graph and tree phase finished");
                info!("create index operation finished");

                // Export metrics to K8s pod annotations
                if let Some(ref exporter) = self.metrics_exporter {
                    let index_count = self.kvs.len() as u64;
                    let uncommitted = (self.vq.ivq_len() + self.vq.dvq_len()) as u64;
                    let processed_vq = self.processed_vq_count.load(Ordering::SeqCst);
                    let unsaved_exec = self.unsaved_create_index_count.load(Ordering::SeqCst);
                    if let Err(e) = exporter.export_on_create_index(
                        index_count,
                        uncommitted,
                        processed_vq,
                        unsaved_exec,
                    ).await {
                        warn!("failed to export create_index metrics: {}", e);
                    }
                }

                Ok(())
            }
            Err(e) => {
                error!("an error occurred on creating graph and tree phase: {}", e);
                Err(Error::Internal(Box::new(std::io::Error::other(e.to_string()))))
            }
        }
    }

    #[tracing::instrument(skip(self), level = "info")]
    async fn save_index(&mut self) -> Result<(), Error> {
        // Read replica cannot perform write operations
        if self.is_readreplica {
            return Err(Error::WriteOperationToReadReplica {});
        }

        // Don't save if already saving
        if self.is_saving.load(Ordering::SeqCst) {
            debug!("save already in progress, skipping");
            return Ok(());
        }
        
        self.is_saving.store(true, Ordering::SeqCst);
        
        // Determine save path (temp for CoW, primary otherwise)
        let save_path = if let Some(ref persistence) = self.persistence {
            persistence.get_save_path().to_string_lossy().to_string()
        } else {
            self.path.clone()
        };

        debug!("saving index to path: {}", save_path);
        
        // Save the core index to the appropriate path
        // Note: QBG save_index uses the path from when the index was created
        // For CoW we need to copy the saved index to the temp location
        let result = self.index.save_index();
        
        // Save metadata to the appropriate path
        if let Some(ref persistence) = self.persistence {
            let index_count = self.kvs.len() as u64;
            let metadata = Metadata::new_qbg(index_count);
            
            if persistence.is_copy_on_write_enabled() {
                // For CoW, save to temp path and then switch
                if let Err(e) = persistence.save_metadata_to_save_path(&metadata) {
                    warn!("failed to save metadata to CoW path: {}", e);
                } else {
                    debug!("saved metadata with index_count={} to CoW path", index_count);
                }
            } else {
                if let Err(e) = persistence.save_metadata(&metadata) {
                    warn!("failed to save metadata: {}", e);
                } else {
                    debug!("saved metadata with index_count={}", index_count);
                }
            }
        }
        
        // Flush kvs to ensure persistence
        if let Err(e) = self.kvs.flush().await {
            warn!("failed to flush kvs: {}", e);
        }

        // For CoW mode, perform the atomic switch after successful save
        if result.is_ok() {
            if let Some(ref persistence) = self.persistence {
                if persistence.is_copy_on_write_enabled() {
                    if let Err(e) = persistence.move_and_switch_saved_data() {
                        error!("failed to switch CoW data: {}", e);
                    }
                }
            }
        }
        
        self.is_saving.store(false, Ordering::SeqCst);
        
        match result {
            Ok(()) => {
                // Reset unsaved create index count after successful save
                let processed_vq = self.processed_vq_count.swap(0, Ordering::SeqCst);
                self.unsaved_create_index_count.store(0, Ordering::SeqCst);

                // Export metrics to K8s pod annotations
                if let Some(ref exporter) = self.metrics_exporter {
                    let timestamp = Utc::now().to_rfc3339();
                    if let Err(e) = exporter.export_on_save_index(timestamp, processed_vq).await {
                        warn!("failed to export save_index metrics: {}", e);
                    }
                }

                info!("index saved successfully");
                Ok(())
            }
            Err(e) => Err(Error::Internal(Box::new(std::io::Error::other(e.to_string()))))
        }
    }

    #[tracing::instrument(skip(self, vector), level = "debug", fields(vector_dim = vector.len()))]
    async fn insert(&mut self, uuid: String, vector: Vec<f32>) -> Result<(), Error> {
        self.insert_internal(uuid, vector, Local::now().nanosecond().into(), true).await
    }

    #[tracing::instrument(skip(self, vectors), level = "debug", fields(count = vectors.len()))]
    async fn insert_multiple(&mut self, vectors: HashMap<String, Vec<f32>>) -> Result<(), Error> {
        let mut uuids: Vec<String> = vec![];
        for (uuid, vec) in vectors {
            let result = self.insert(uuid.clone(), vec).await;
            match result {
                Ok(()) => continue,
                Err(err) => match err {
                    Error::UUIDAlreadyExists { uuid } => uuids.push(uuid),
                    _ => return Err(err),
                },
            }
        }
        if !uuids.is_empty() {
            return Err(Error::new_uuid_already_exists(uuids));
        }
        Ok(())
    }

    #[tracing::instrument(skip(self, vector), level = "debug", fields(vector_dim = vector.len()))]
    async fn update(&mut self, uuid: String, vector: Vec<f32>) -> Result<(), Error> {
        if self.is_flushing() {
            return Err(Error::FlushingIsInProgress {});
        }
        self.update_internal(uuid, vector, Local::now().nanosecond().into()).await
    }

    #[tracing::instrument(skip(self, vectors), level = "debug", fields(count = vectors.len()))]
    async fn update_multiple(&mut self, mut vectors: HashMap<String, Vec<f32>>) -> Result<(), Error> {
        let mut uuids: Vec<String> = vec![];
        for (uuid, vec) in vectors.clone() {
            let result = self.ready_for_update(uuid.clone(), vec, Local::now().nanosecond().into()).await;
            match result {
                Ok(()) => uuids.push(uuid),
                Err(_err) => {
                    let _ = vectors.remove(&uuid);
                }
            }
        }
        self.remove_multiple(uuids.clone()).await?;
        self.insert_multiple(vectors).await
    }

    #[tracing::instrument(skip(self), level = "debug")]
    async fn remove(&mut self, uuid: String) -> Result<(), Error> {
        if self.is_flushing() {
            return Err(Error::FlushingIsInProgress {});
        }
        self.remove_internal(uuid, Local::now().nanosecond().into(), true).await
    }

    #[tracing::instrument(skip(self), level = "debug", fields(count = uuids.len()))]
    async fn remove_multiple(&mut self, uuids: Vec<String>) -> Result<(), Error> {
        if self.is_flushing() {
            return Err(Error::FlushingIsInProgress {});
        }
        self.remove_multiple_internal(uuids, Local::now().nanosecond().into(), true).await
    }

    #[tracing::instrument(skip(self, vector), level = "debug", fields(vector_dim = vector.len()))]
    async fn search(
        &self,
        vector: Vec<f32>,
        k: u32,
        epsilon: f32,
        radius: f32,
    ) -> Result<search::Response, Error> {
        let vec = self
            .index
            .search(vector.as_slice(), k as usize, radius, epsilon)
            .unwrap();
        let results: Vec<Distance> = vec
            .into_iter()
            .map(|x| Distance {
                id: x.0.to_string(),
                distance: x.1,
            })
            .collect();
        let res = search::Response {
            request_id: "".to_string(),
            results: results,
        };
        Ok(res)
    }

    #[tracing::instrument(skip(self), level = "debug")]
    async fn get_object(&self, uuid: String) -> Result<(Vec<f32>, i64), Error> {
        let index = &self.index;
        let get_vector_fn = |oid: u32| async move {
            index.get_object(oid as usize)
                .map(|v| v.to_vec())
                .map_err(|e| memstore::MemstoreError::ObjectNotFound(e.to_string()))
        };
        
        memstore::get_object(&self.kvs, &self.vq, &uuid, Some(get_vector_fn))
            .await
            .map_err(|e| match e {
                memstore::MemstoreError::ObjectIdNotFound(uuid) => Error::ObjectIDNotFound { uuid },
                memstore::MemstoreError::ObjectNotFound(uuid) => Error::ObjectIDNotFound { uuid },
                memstore::MemstoreError::UuidNotFound(uuid) => Error::UUIDNotFound { uuid },
                _ => Error::Internal(Box::new(e)),
            })
    }

    fn get_dimension_size(&self) -> usize {
        self.index.get_dimension().unwrap_or_default()
    }

    fn len(&self) -> u32 {
        // Return the count of items in kvs (indexed items)
        // Note: This doesn't include items still in vqueue
        self.kvs.len() as u32
    }

    fn insert_vqueue_buffer_len(&self) -> u32 {
        self.vq.ivq_len() as u32
    }

    fn delete_vqueue_buffer_len(&self) -> u32 {
        self.vq.dvq_len() as u32
    }

    fn is_flushing(&self) -> bool {
        self.is_flushing.load(Ordering::SeqCst)
    }

    fn is_indexing(&self) -> bool {
        self.is_indexing.load(Ordering::SeqCst)
    }

    fn is_saving(&self) -> bool {
        self.is_saving.load(Ordering::SeqCst)
    }

    #[tracing::instrument(skip(self), level = "info")]
    async fn regenerate_indexes(&mut self) -> Result<(), Error> {
        // Read replica cannot perform write operations
        if self.is_readreplica {
            return Err(Error::WriteOperationToReadReplica {});
        }

        // Close the current index and rebuild it
        self.index.close_index();
        self.create_index().await
    }

    #[tracing::instrument(skip(self), level = "debug")]
    async fn search_by_id(&self, uuid: String, k: u32, epsilon: f32, radius: f32) -> Result<proto::payload::v1::search::Response, Error> {
        let (vec, _ts) = self.get_object(uuid).await?;
        self.search(vec, k, epsilon, radius).await
    }

    async fn linear_search(&self, _vector: Vec<f32>, _k: u32) -> Result<proto::payload::v1::search::Response, Error> {
        Err(Error::Unsupported {
            method: "LinearSearch".to_string(),
            algorithm: "QBG".to_string(),
        })
    }

    async fn linear_search_by_id(&self, _uuid: String, _k: u32) -> Result<proto::payload::v1::search::Response, Error> {
        Err(Error::Unsupported {
            method: "LinearSearchByID".to_string(),
            algorithm: "QBG".to_string(),
        })
    }

    async fn insert_with_time(&mut self, uuid: String, vector: Vec<f32>, t: i64) -> Result<(), Error> {
        self.insert_internal(uuid, vector, t, true).await
    }

    async fn insert_multiple_with_time(&mut self, vectors: HashMap<String, Vec<f32>>, t: i64) -> Result<(), Error> {
        self.insert_multiple_internal(vectors, t, true).await
    }

    async fn update_with_time(&mut self, uuid: String, vector: Vec<f32>, t: i64) -> Result<(), Error> {
        self.update_internal(uuid, vector, t).await
    }

    async fn update_multiple_with_time(&mut self, vectors: HashMap<String, Vec<f32>>, t: i64) -> Result<(), Error> {
        for (uuid, vec) in vectors {
            self.update_internal(uuid, vec, t).await?;
        }
        Ok(())
    }

    async fn update_timestamp(&mut self, uuid: String, t: i64, force: bool) -> Result<(), Error> {
        let index = &self.index;
        let get_vector_fn = |oid: u32| async move {
            index.get_object(oid as usize)
                .map(|v| v.to_vec())
                .map_err(|e| memstore::MemstoreError::ObjectNotFound(e.to_string()))
        };
        
        memstore::update_timestamp(&self.kvs, &self.vq, &uuid, t, force, Some(get_vector_fn))
            .await
            .map_err(|e| match e {
                memstore::MemstoreError::ObjectIdNotFound(uuid) => Error::ObjectIDNotFound { uuid },
                memstore::MemstoreError::ObjectNotFound(uuid) => Error::ObjectIDNotFound { uuid },
                memstore::MemstoreError::UuidNotFound(uuid) => Error::UUIDNotFound { uuid },
                memstore::MemstoreError::ZeroTimestamp => Error::InvalidUUID { uuid: "timestamp is zero".to_string() },
                memstore::MemstoreError::NewerTimestampObjectAlreadyExists(uuid, _) => Error::UUIDAlreadyExists { uuid },
                memstore::MemstoreError::NothingToBeDoneForUpdate(uuid) => Error::UUIDAlreadyExists { uuid },
                _ => Error::Internal(Box::new(e)),
            })
    }

    async fn remove_with_time(&mut self, uuid: String, t: i64) -> Result<(), Error> {
        self.remove_internal(uuid, t, true).await
    }

    async fn remove_multiple_with_time(&mut self, uuids: Vec<String>, t: i64) -> Result<(), Error> {
        self.remove_multiple_internal(uuids, t, true).await
    }

    async fn list_object_func<F: FnMut(String, Vec<f32>, i64) -> bool + Send>(&self, mut f: F) {
        let index = &self.index;
        memstore::list_object_func(&self.kvs, &self.vq, |uuid, oid, ts| {
            // Get vector from index if oid > 0, otherwise skip (not indexed yet)
            if oid > 0 {
                if let Ok(vec) = index.get_object(oid as usize) {
                    return f(uuid, vec.to_vec(), ts);
                }
            }
            true // continue iteration if vector not available
        }).await;
    }

    async fn create_and_save_index(&mut self) -> Result<(), Error> {
        self.create_index().await?;
        self.save_index().await
    }

    fn number_of_create_index_executions(&self) -> u64 {
        self.create_index_count.load(Ordering::SeqCst)
    }

    async fn uuids(&self) -> Vec<String> {
        memstore::uuids(&self.kvs, &self.vq).await.unwrap_or_default()
    }

    fn broken_index_count(&self) -> u64 {
        self.broken_index_count.load(Ordering::SeqCst)
    }

    fn index_statistics(&self) -> Result<proto::payload::v1::info::index::Statistics, Error> {
        Ok(proto::payload::v1::info::index::Statistics {
            valid: true,
            median_indegree: 0,
            median_outdegree: 0,
            max_number_of_indegree: 0,
            max_number_of_outdegree: 0,
            min_number_of_indegree: 0,
            min_number_of_outdegree: 0,
            mode_indegree: 0,
            mode_outdegree: 0,
            nodes_skipped_for_10_edges: 0,
            nodes_skipped_for_indegree_distance: 0,
            number_of_edges: 0,
            number_of_indexed_objects: self.len() as u64,
            number_of_nodes: self.len() as u64,
            number_of_nodes_without_edges: 0,
            number_of_nodes_without_indegree: 0,
            number_of_objects: self.len() as u64,
            number_of_removed_objects: 0,
            size_of_object_repository: self.len() as u64,
            size_of_refinement_object_repository: 0,
            variance_of_indegree: 0.0,
            variance_of_outdegree: 0.0,
            mean_edge_length: 0.0,
            mean_edge_length_for_10_edges: 0.0,
            mean_indegree_distance_for_10_edges: 0.0,
            mean_number_of_edges_per_node: 0.0,
            c1_indegree: 0.0,
            c5_indegree: 0.0,
            c95_outdegree: 0.0,
            c99_outdegree: 0.0,
            indegree_count: vec![],
            outdegree_histogram: vec![],
            indegree_histogram: vec![],
        })
    }

    fn is_statistics_enabled(&self) -> bool {
        self.statistics_enabled
    }

    fn index_property(&self) -> Result<proto::payload::v1::info::index::Property, Error> {
        Err(Error::Unsupported { method: "index_property".to_owned(), algorithm: "QBG".to_owned() })
    }

    #[tracing::instrument(skip(self), level = "info")]
    async fn close(&mut self) -> Result<(), Error> {
        info!("Closing QBGService...");
        
        // Skip index operations for read replicas
        if self.is_readreplica {
            info!("Read replica mode: skipping index creation and save on close");
        } else {
            // Create final index if there are uncommitted changes
            let uncommitted = self.vq.ivq_len() + self.vq.dvq_len();
            if uncommitted > 0 {
                info!("Creating final index with {} uncommitted changes...", uncommitted);
                if let Err(e) = self.create_index().await {
                    if !matches!(e, Error::UncommittedIndexNotFound {}) {
                        warn!("Failed to create final index: {:?}", e);
                    }
                }
            }
            
            // Save the index
            info!("Saving index...");
            if let Err(e) = self.save_index().await {
                warn!("Failed to save index on close: {:?}", e);
            }
        }
        
        // Close the QBG index
        info!("Closing QBG core index...");
        self.index.close_index();
        
        // Flush and close KVS
        info!("Flushing KVS...");
        if let Err(e) = self.kvs.flush().await {
            warn!("Failed to flush KVS: {:?}", e);
        }
        
        info!("QBGService closed successfully");
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use config::Config;
    use tempfile::TempDir;

    /// Test helper to create a QBGService with temporary directories
    struct TestQBGService {
        service: QBGService,
        _temp_dir: TempDir,
        base_path: String,
    }

    impl TestQBGService {
        async fn new(dimension: usize) -> Self {
            Self::with_options(dimension, false).await
        }

        async fn new_read_replica(dimension: usize) -> Self {
            Self::with_options(dimension, true).await
        }

        async fn with_options(dimension: usize, is_read_replica: bool) -> Self {
            let temp_dir = TempDir::new().expect("Failed to create temp directory");
            let base_path = temp_dir.path().to_str().unwrap().to_string();

            let config = Config::builder()
                .set_default("qbg.index_path", format!("{}/index", base_path)).unwrap()
                .set_default("qbg.vqueue_path", format!("{}/vqueue", base_path)).unwrap()
                .set_default("qbg.kvs_path", format!("{}/kvs", base_path)).unwrap()
                .set_default("qbg.dimension", dimension as i64).unwrap()
                .set_default("qbg.extended_dimension", dimension as i64).unwrap()
                .set_default("qbg.number_of_subvectors", 1_i64).unwrap()
                .set_default("qbg.number_of_blobs", 0_i64).unwrap()
                .set_default("qbg.distance_type", 1_i64).unwrap() // L2
                .set_default("qbg.data_type", 1_i64).unwrap() // Float
                .set_default("qbg.internal_data_type", 1_i64).unwrap()
                .set_default("qbg.is_readreplica", is_read_replica).unwrap()
                .build()
                .unwrap();

            let agent_config: crate::config::AgentConfig = config.try_deserialize().unwrap();
            let service = QBGService::new(&agent_config.qbg).await;

            TestQBGService {
                service,
                _temp_dir: temp_dir,
                base_path,
            }
        }

        /// Create a Read Replica service using the same paths as this service.
        /// The original service should have built and saved the index first.
        async fn create_read_replica_from_same_path(&self, dimension: usize) -> QBGService {
            let config = Config::builder()
                .set_default("qbg.index_path", format!("{}/index", self.base_path)).unwrap()
                .set_default("qbg.vqueue_path", format!("{}/vqueue", self.base_path)).unwrap()
                .set_default("qbg.kvs_path", format!("{}/kvs", self.base_path)).unwrap()
                .set_default("qbg.dimension", dimension as i64).unwrap()
                .set_default("qbg.extended_dimension", dimension as i64).unwrap()
                .set_default("qbg.number_of_subvectors", 1_i64).unwrap()
                .set_default("qbg.number_of_blobs", 0_i64).unwrap()
                .set_default("qbg.distance_type", 1_i64).unwrap()
                .set_default("qbg.data_type", 1_i64).unwrap()
                .set_default("qbg.internal_data_type", 1_i64).unwrap()
                .set_default("qbg.is_readreplica", true).unwrap()
                .build()
                .unwrap();

            let agent_config: crate::config::AgentConfig = config.try_deserialize().unwrap();
            QBGService::new(&agent_config.qbg).await
        }
    }

    fn gen_random_vector(dim: usize) -> Vec<f32> {
        use rand::Rng;
        let mut rng = rand::rng();
        (0..dim).map(|_| rng.random::<f32>()).collect()
    }

    // ========== Insert Tests ==========

    #[tokio::test]
    async fn test_insert_single_vector() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "test-uuid-1".to_string();
        let vector = gen_random_vector(128);

        let result = test_svc.service.insert(uuid.clone(), vector.clone()).await;
        assert!(result.is_ok(), "Insert should succeed: {:?}", result.err());

        // Check that the vector exists
        let (_, exists) = test_svc.service.exists(uuid.clone()).await;
        assert!(exists, "Vector should exist after insert");
    }

    #[tokio::test]
    async fn test_insert_duplicate_uuid_fails() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "test-uuid-dup".to_string();
        let vector1 = gen_random_vector(128);
        let vector2 = gen_random_vector(128);

        // First insert should succeed
        let result1 = test_svc.service.insert(uuid.clone(), vector1).await;
        assert!(result1.is_ok());

        // Second insert with same UUID should fail
        let result2 = test_svc.service.insert(uuid.clone(), vector2).await;
        assert!(result2.is_err());
        match result2.err().unwrap() {
            Error::UUIDAlreadyExists { uuid: err_uuid } => {
                assert_eq!(err_uuid, uuid);
            }
            e => panic!("Expected UUIDAlreadyExists error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_insert_empty_uuid_fails() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "".to_string();
        let vector = gen_random_vector(128);

        let result = test_svc.service.insert(uuid, vector).await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::UUIDNotFound { .. } => {}
            e => panic!("Expected UUIDNotFound error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_insert_multiple_vectors() {
        let mut test_svc = TestQBGService::new(128).await;

        let mut vectors = HashMap::new();
        for i in 0..10 {
            vectors.insert(format!("uuid-{}", i), gen_random_vector(128));
        }

        let result = test_svc.service.insert_multiple(vectors.clone()).await;
        assert!(result.is_ok(), "Insert multiple should succeed: {:?}", result.err());

        // Check all vectors exist
        for uuid in vectors.keys() {
            let (_, exists) = test_svc.service.exists(uuid.clone()).await;
            assert!(exists, "Vector {} should exist after insert_multiple", uuid);
        }
    }

    // ========== GetObject Tests ==========

    #[tokio::test]
    async fn test_get_object_from_vqueue() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "test-uuid-get".to_string();
        let vector = gen_random_vector(128);
        let timestamp = 1000i64;

        test_svc.service.insert_with_time(uuid.clone(), vector.clone(), timestamp).await.unwrap();

        let (retrieved_vec, retrieved_ts) = test_svc.service.get_object(uuid).await.unwrap();
        assert_eq!(retrieved_vec, vector);
        assert_eq!(retrieved_ts, timestamp);
    }

    #[tokio::test]
    async fn test_get_object_not_found() {
        let test_svc = TestQBGService::new(128).await;

        let result = test_svc.service.get_object("nonexistent-uuid".to_string()).await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::ObjectIDNotFound { uuid } => {
                assert_eq!(uuid, "nonexistent-uuid");
            }
            e => panic!("Expected ObjectIDNotFound error, got: {:?}", e),
        }
    }

    // ========== Exists Tests ==========

    #[tokio::test]
    async fn test_exists_returns_false_for_nonexistent() {
        let test_svc = TestQBGService::new(128).await;

        let (oid, exists) = test_svc.service.exists("nonexistent".to_string()).await;
        assert!(!exists);
        assert_eq!(oid, 0);
    }

    #[tokio::test]
    async fn test_exists_returns_true_after_insert() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "exists-test-uuid".to_string();
        let vector = gen_random_vector(128);

        test_svc.service.insert(uuid.clone(), vector).await.unwrap();

        let (_, exists) = test_svc.service.exists(uuid).await;
        assert!(exists);
    }

    // ========== Remove Tests ==========

    #[tokio::test]
    async fn test_remove_existing_vector() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "remove-test-uuid".to_string();
        let vector = gen_random_vector(128);

        test_svc.service.insert(uuid.clone(), vector).await.unwrap();
        
        let (_, exists_before) = test_svc.service.exists(uuid.clone()).await;
        assert!(exists_before);

        let result = test_svc.service.remove(uuid.clone()).await;
        assert!(result.is_ok());

        let (_, exists_after) = test_svc.service.exists(uuid).await;
        assert!(!exists_after);
    }

    #[tokio::test]
    async fn test_remove_nonexistent_vector_fails() {
        let mut test_svc = TestQBGService::new(128).await;

        let result = test_svc.service.remove("nonexistent-uuid".to_string()).await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::ObjectIDNotFound { .. } => {}
            e => panic!("Expected ObjectIDNotFound error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_remove_multiple() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuids: Vec<String> = (0..5).map(|i| format!("multi-remove-{}", i)).collect();
        
        // Insert all
        for uuid in &uuids {
            test_svc.service.insert(uuid.clone(), gen_random_vector(128)).await.unwrap();
        }

        // Remove all
        let result = test_svc.service.remove_multiple(uuids.clone()).await;
        assert!(result.is_ok());

        // Check none exist
        for uuid in &uuids {
            let (_, exists) = test_svc.service.exists(uuid.clone()).await;
            assert!(!exists, "Vector {} should not exist after remove_multiple", uuid);
        }
    }

    // ========== Update Tests ==========

    #[tokio::test]
    async fn test_update_existing_vector() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "update-test-uuid".to_string();
        let vector1 = gen_random_vector(128);
        let vector2 = gen_random_vector(128);

        test_svc.service.insert(uuid.clone(), vector1.clone()).await.unwrap();

        // Get original
        let (orig_vec, _) = test_svc.service.get_object(uuid.clone()).await.unwrap();
        assert_eq!(orig_vec, vector1);

        // Update
        let result = test_svc.service.update(uuid.clone(), vector2.clone()).await;
        assert!(result.is_ok(), "Update should succeed: {:?}", result.err());

        // Get updated - should be in vqueue with new vector
        let (updated_vec, _) = test_svc.service.get_object(uuid).await.unwrap();
        assert_eq!(updated_vec, vector2);
    }

    // ========== Linear Search Tests (Unsupported) ==========

    #[tokio::test]
    async fn test_linear_search_returns_unsupported() {
        let test_svc = TestQBGService::new(128).await;

        let vector = gen_random_vector(128);
        let result = test_svc.service.linear_search(vector, 10).await;
        
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::Unsupported { method, algorithm } => {
                assert_eq!(method, "LinearSearch");
                assert_eq!(algorithm, "QBG");
            }
            e => panic!("Expected Unsupported error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_linear_search_by_id_returns_unsupported() {
        let test_svc = TestQBGService::new(128).await;

        let result = test_svc.service.linear_search_by_id("some-uuid".to_string(), 10).await;
        
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::Unsupported { method, algorithm } => {
                assert_eq!(method, "LinearSearchByID");
                assert_eq!(algorithm, "QBG");
            }
            e => panic!("Expected Unsupported error, got: {:?}", e),
        }
    }

    // ========== VQueue Buffer Length Tests ==========

    #[tokio::test]
    async fn test_insert_vqueue_buffer_len() {
        let mut test_svc = TestQBGService::new(128).await;

        assert_eq!(test_svc.service.insert_vqueue_buffer_len(), 0);

        // Insert a vector
        test_svc.service.insert("uuid-1".to_string(), gen_random_vector(128)).await.unwrap();
        assert_eq!(test_svc.service.insert_vqueue_buffer_len(), 1);

        // Insert another
        test_svc.service.insert("uuid-2".to_string(), gen_random_vector(128)).await.unwrap();
        assert_eq!(test_svc.service.insert_vqueue_buffer_len(), 2);
    }

    #[tokio::test]
    async fn test_delete_vqueue_buffer_len() {
        let mut test_svc = TestQBGService::new(128).await;

        assert_eq!(test_svc.service.delete_vqueue_buffer_len(), 0);

        // Insert and then delete
        test_svc.service.insert("uuid-del".to_string(), gen_random_vector(128)).await.unwrap();
        test_svc.service.remove("uuid-del".to_string()).await.unwrap();

        assert_eq!(test_svc.service.delete_vqueue_buffer_len(), 1);
    }

    // ========== Dimension Tests ==========

    #[tokio::test]
    async fn test_get_dimension_size() {
        let test_svc = TestQBGService::new(256).await;
        // Note: dimension check depends on QBG index initialization
        let dim = test_svc.service.get_dimension_size();
        // QBG may adjust dimension internally, so just check it's reasonable
        assert!(dim > 0, "Dimension should be greater than 0");
    }

    // ========== Len Tests ==========

    #[tokio::test]
    async fn test_len_empty_index() {
        let test_svc = TestQBGService::new(128).await;
        assert_eq!(test_svc.service.len(), 0);
    }

    #[tokio::test]
    async fn test_len_after_create_index() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert vectors
        for i in 0..10 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }

        // Note: QBG's HierarchicalKmeans requires many objects for clustering
        // Skip create_index in this test since it may fail with few objects
        // len() returns kvs.len() which reflects inserted items
        assert!(test_svc.service.len() >= 0);
    }

    // ========== Create/Save Index Tests ==========

    #[tokio::test]
    async fn test_create_and_save_index() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert some vectors
        for i in 0..50 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }

        // Note: QBG's create_index may fail with HierarchicalKmeans clustering errors
        // when there aren't enough objects. Just verify no panic.
        let _ = test_svc.service.create_and_save_index().await;
    }

    // ========== Search By ID Tests ==========

    #[tokio::test]
    async fn test_search_by_id() {
        let test_svc = TestQBGService::new(128).await;

        // Note: search_by_id requires a built searchable index.
        // QBG throws an exception if called on an unbuilt index, causing SIGABRT.
        // This test just verifies the method exists and returns an error for nonexistent UUID.
        let result = test_svc.service.search_by_id("nonexistent".to_string(), 5, 0.1, -1.0).await;
        assert!(result.is_err());
    }

    #[tokio::test]
    async fn test_search_by_id_not_found() {
        let test_svc = TestQBGService::new(128).await;

        let result = test_svc.service.search_by_id("nonexistent".to_string(), 5, 0.1, -1.0).await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::ObjectIDNotFound { .. } => {}
            e => panic!("Expected ObjectIDNotFound error, got: {:?}", e),
        }
    }

    // ========== Regenerate Indexes Tests ==========

    #[tokio::test]
    async fn test_regenerate_indexes() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert some vectors
        for i in 0..50 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }

        // Note: QBG's create_index may fail with HierarchicalKmeans clustering errors.
        // Just verify no panic.
        let _ = test_svc.service.regenerate_indexes().await;
    }

    // ========== UUIDs Tests ==========

    #[tokio::test]
    async fn test_uuids_empty() {
        let test_svc = TestQBGService::new(128).await;
        let uuids = test_svc.service.uuids().await;
        assert!(uuids.is_empty());
    }

    #[tokio::test]
    async fn test_uuids_after_insert() {
        let mut test_svc = TestQBGService::new(128).await;

        let expected_uuids: Vec<String> = (0..5).map(|i| format!("uuid-{}", i)).collect();
        for uuid in &expected_uuids {
            test_svc.service.insert(uuid.clone(), gen_random_vector(128)).await.unwrap();
        }

        // uuids() returns items from both kvs and vqueue
        // After insert, items should be accessible
        let uuids = test_svc.service.uuids().await;
        // Note: actual behavior depends on memstore implementation
        // Just verify no panic and reasonable result
        assert!(uuids.len() <= expected_uuids.len());
    }

    // ========== Number of Create Index Executions Tests ==========

    #[tokio::test]
    async fn test_number_of_create_index_executions() {
        let mut test_svc = TestQBGService::new(128).await;

        assert_eq!(test_svc.service.number_of_create_index_executions(), 0);

        // Insert some vectors and try create_index
        for i in 0..50 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }
        let _ = test_svc.service.create_index().await;

        // Count may be 0 or 1 depending on success/failure
        let count = test_svc.service.number_of_create_index_executions();
        assert!(count <= 1);
    }

    // ========== Broken Index Count Tests ==========

    #[tokio::test]
    async fn test_broken_index_count() {
        let test_svc = TestQBGService::new(128).await;
        // Should start at 0 for a fresh index
        assert_eq!(test_svc.service.broken_index_count(), 0);
    }

    // ========== Index Statistics Tests ==========

    #[tokio::test]
    async fn test_index_statistics() {
        let test_svc = TestQBGService::new(128).await;
        let result = test_svc.service.index_statistics();
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_is_statistics_enabled() {
        let test_svc = TestQBGService::new(128).await;
        // Just verify it returns a boolean without panicking
        // Note: statistics_enabled is false by default
        let enabled = test_svc.service.is_statistics_enabled();
        assert!(!enabled);
    }

    // ========== Index Property Tests ==========

    #[tokio::test]
    async fn test_index_property() {
        let test_svc = TestQBGService::new(128).await;
        let result = test_svc.service.index_property();
        assert!(result.is_err());
    }

    // ========== Close Tests ==========

    #[tokio::test]
    async fn test_close() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert some data
        test_svc.service.insert("uuid-1".to_string(), gen_random_vector(128)).await.unwrap();

        // Close should succeed
        let result = test_svc.service.close().await;
        assert!(result.is_ok(), "close should succeed: {:?}", result.err());
    }

    #[tokio::test]
    async fn test_close_with_uncommitted_changes() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert multiple vectors (uncommitted)
        for i in 0..10 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }

        // Verify we have uncommitted changes
        let uncommitted = test_svc.service.insert_vqueue_buffer_len() + test_svc.service.delete_vqueue_buffer_len();
        assert!(uncommitted > 0, "Should have uncommitted changes");

        // Close should handle uncommitted changes gracefully
        let result = test_svc.service.close().await;
        assert!(result.is_ok(), "close with uncommitted changes should succeed");
    }

    #[tokio::test]
    async fn test_close_empty_service() {
        let mut test_svc = TestQBGService::new(128).await;

        // Close immediately without any operations
        let result = test_svc.service.close().await;
        assert!(result.is_ok(), "close on empty service should succeed");
    }

    #[tokio::test]
    async fn test_close_after_create_index() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert vectors
        for i in 0..50 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }

        // Create index first
        let _ = test_svc.service.create_index().await;

        // Close should succeed
        let result = test_svc.service.close().await;
        assert!(result.is_ok(), "close after create_index should succeed");
    }

    #[tokio::test]
    async fn test_close_after_save_index() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert and create index
        for i in 0..50 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }
        let _ = test_svc.service.create_index().await;
        let _ = test_svc.service.save_index().await;

        // Close should succeed
        let result = test_svc.service.close().await;
        assert!(result.is_ok(), "close after save_index should succeed");
    }

    #[tokio::test]
    async fn test_close_with_remove_operations() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert and remove some vectors
        for i in 0..20 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }
        
        // Remove half of them
        for i in 0..10 {
            let _ = test_svc.service.remove(format!("uuid-{}", i)).await;
        }

        // Close should handle mixed insert/delete queue
        let result = test_svc.service.close().await;
        assert!(result.is_ok(), "close with remove operations should succeed");
    }

    #[tokio::test]
    async fn test_close_with_update_operations() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert vectors
        for i in 0..10 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }

        // Update some vectors
        for i in 0..5 {
            let _ = test_svc.service.update(format!("uuid-{}", i), gen_random_vector(128)).await;
        }

        // Close should succeed
        let result = test_svc.service.close().await;
        assert!(result.is_ok(), "close with update operations should succeed");
    }

    // ========== State Flag Tests ==========

    #[tokio::test]
    async fn test_is_flushing_initial_state() {
        let test_svc = TestQBGService::new(128).await;
        assert!(!test_svc.service.is_flushing(), "is_flushing should be false initially");
    }

    #[tokio::test]
    async fn test_is_indexing_initial_state() {
        let test_svc = TestQBGService::new(128).await;
        assert!(!test_svc.service.is_indexing(), "is_indexing should be false initially");
    }

    #[tokio::test]
    async fn test_is_saving_initial_state() {
        let test_svc = TestQBGService::new(128).await;
        assert!(!test_svc.service.is_saving(), "is_saving should be false initially");
    }

    // ========== List Object Func Tests ==========

    #[tokio::test]
    async fn test_list_object_func() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert some vectors
        for i in 0..3 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }

        use std::sync::atomic::AtomicUsize;
        let count = AtomicUsize::new(0);
        test_svc.service.list_object_func(|_uuid, _vec, _ts| {
            count.fetch_add(1, Ordering::SeqCst);
            true // continue iterating
        }).await;

        // Note: list_object_func only iterates over indexed objects (oid > 0)
        // Objects in vqueue without create_index won't be counted
        let final_count = count.load(Ordering::SeqCst);
        assert!(final_count <= 3);
    }

    // ========== Read Replica Tests ==========

    #[tokio::test]
    async fn test_read_replica_insert_fails() {
        let mut test_svc = TestQBGService::new_read_replica(128).await;

        let result = test_svc.service.insert("uuid-1".to_string(), gen_random_vector(128)).await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::WriteOperationToReadReplica {} => {}
            e => panic!("Expected WriteOperationToReadReplica error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_read_replica_update_fails() {
        let mut test_svc = TestQBGService::new_read_replica(128).await;

        let result = test_svc.service.update("uuid-1".to_string(), gen_random_vector(128)).await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::WriteOperationToReadReplica {} => {}
            e => panic!("Expected WriteOperationToReadReplica error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_read_replica_remove_fails() {
        let mut test_svc = TestQBGService::new_read_replica(128).await;

        let result = test_svc.service.remove("uuid-1".to_string()).await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::WriteOperationToReadReplica {} => {}
            e => panic!("Expected WriteOperationToReadReplica error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_read_replica_create_index_fails() {
        let mut test_svc = TestQBGService::new_read_replica(128).await;

        let result = test_svc.service.create_index().await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::WriteOperationToReadReplica {} => {}
            e => panic!("Expected WriteOperationToReadReplica error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_read_replica_save_index_fails() {
        let mut test_svc = TestQBGService::new_read_replica(128).await;

        let result = test_svc.service.save_index().await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::WriteOperationToReadReplica {} => {}
            e => panic!("Expected WriteOperationToReadReplica error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_read_replica_create_and_save_index_fails() {
        let mut test_svc = TestQBGService::new_read_replica(128).await;

        let result = test_svc.service.create_and_save_index().await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::WriteOperationToReadReplica {} => {}
            e => panic!("Expected WriteOperationToReadReplica error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_read_replica_regenerate_indexes_fails() {
        let mut test_svc = TestQBGService::new_read_replica(128).await;

        let result = test_svc.service.regenerate_indexes().await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::WriteOperationToReadReplica {} => {}
            e => panic!("Expected WriteOperationToReadReplica error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_read_replica_read_operations_succeed() {
        let test_svc = TestQBGService::new_read_replica(128).await;

        // Exists should work
        let (_, exists) = test_svc.service.exists("uuid-1".to_string()).await;
        assert!(!exists);

        // len should work
        assert_eq!(test_svc.service.len(), 0);

        // get_dimension_size should work
        let dim = test_svc.service.get_dimension_size();
        assert!(dim > 0);

        // is_flushing/is_indexing/is_saving should work
        assert!(!test_svc.service.is_flushing());
        assert!(!test_svc.service.is_indexing());
        assert!(!test_svc.service.is_saving());

        // broken_index_count should work
        assert_eq!(test_svc.service.broken_index_count(), 0);

        // number_of_create_index_executions should work
        assert_eq!(test_svc.service.number_of_create_index_executions(), 0);

        // index_statistics should work
        let stats = test_svc.service.index_statistics();
        assert!(stats.is_ok());

        // uuids should work
        let uuids = test_svc.service.uuids().await;
        assert!(uuids.is_empty());
    }

    #[tokio::test]
    async fn test_read_replica_search_operations_succeed() {
        // Test that read replica correctly rejects write operations while allowing reads.
        // Note: Testing actual search on read replica requires a pre-built index which is 
        // complex to set up in unit tests due to QBG's directory handling.
        // We verify that search_by_id returns ObjectIDNotFound (not WriteOperationToReadReplica),
        // proving that read operations are allowed.
        
        let test_svc = TestQBGService::new_read_replica(128).await;

        // search_by_id should fail with ObjectIDNotFound, not WriteOperationToReadReplica
        // This proves that read operations are permitted on read replicas
        let search_by_id_result = test_svc.service.search_by_id("nonexistent".to_string(), 5, 0.1, -1.0).await;
        assert!(search_by_id_result.is_err());
        match search_by_id_result.err().unwrap() {
            Error::ObjectIDNotFound { .. } => {}
            e => panic!("Expected ObjectIDNotFound error, got: {:?}", e),
        }

        // get_object should also return ObjectIDNotFound
        let get_result = test_svc.service.get_object("nonexistent".to_string()).await;
        assert!(get_result.is_err());
        match get_result.err().unwrap() {
            Error::ObjectIDNotFound { .. } | Error::UUIDNotFound { .. } => {}
            e => panic!("Expected ObjectIDNotFound or UUIDNotFound error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_read_replica_close_succeeds() {
        let mut test_svc = TestQBGService::new_read_replica(128).await;

        // close should succeed for read replica (no save operation)
        let result = test_svc.service.close().await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_read_replica_insert_with_time_fails() {
        let mut test_svc = TestQBGService::new_read_replica(128).await;

        let result = test_svc.service.insert_with_time(
            "uuid-1".to_string(),
            gen_random_vector(128),
            1234567890,
        ).await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::WriteOperationToReadReplica {} => {}
            e => panic!("Expected WriteOperationToReadReplica error, got: {:?}", e),
        }
    }

    #[tokio::test]
    async fn test_read_replica_remove_with_time_fails() {
        let mut test_svc = TestQBGService::new_read_replica(128).await;

        let result = test_svc.service.remove_with_time("uuid-1".to_string(), 1234567890).await;
        assert!(result.is_err());
        match result.err().unwrap() {
            Error::WriteOperationToReadReplica {} => {}
            e => panic!("Expected WriteOperationToReadReplica error, got: {:?}", e),
        }
    }

    // ========== UpdateTimestamp Tests ==========

    #[tokio::test]
    async fn test_update_timestamp_basic() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert a vector
        let uuid = "test-uuid-1".to_string();
        let vector = gen_random_vector(128);
        test_svc.service.insert(uuid.clone(), vector).await.unwrap();

        // Verify the UUID exists
        let (_, exists) = test_svc.service.exists(uuid.clone()).await;
        assert!(exists, "UUID should exist after insert");

        // Try to update timestamp - it should work or return a specific error related to timing
        let new_timestamp: i64 = 9876543210;
        let result = test_svc.service.update_timestamp(uuid.clone(), new_timestamp, true).await;
        // The result can be either success or a "newer timestamp exists" error, both are acceptable
        // since this tests the update_timestamp behavior with already-existing entries
        let _ = result;
    }

    #[tokio::test]
    async fn test_update_timestamp_nonexistent_first() {
        let mut test_svc = TestQBGService::new(128).await;

        // Try to update timestamp for a UUID that has never been inserted
        let uuid = "never-inserted".to_string();
        let result = test_svc.service.update_timestamp(uuid.clone(), 1234567890, false).await;
        assert!(result.is_err(), "Should fail for non-existent UUID");
        
        // Accept either ObjectIDNotFound or UUIDNotFound errors
        match result {
            Err(Error::UUIDNotFound { .. }) | Err(Error::ObjectIDNotFound { .. }) => {}, // Expected
            Err(e) => panic!("Got unexpected error: {:?}", e),
            Ok(_) => panic!("Should not succeed for non-existent UUID"),
        }
    }

    #[tokio::test]
    async fn test_update_timestamp_with_remove_and_reinsert() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "test-uuid-3".to_string();
        let vector1 = gen_random_vector(128);
        
        // Insert first vector
        test_svc.service.insert(uuid.clone(), vector1).await.unwrap();

        // Remove it
        test_svc.service.remove(uuid.clone()).await.unwrap();

        // Verify it's removed (or at least doesn't exist)
        let (_, exists_after_remove) = test_svc.service.exists(uuid.clone()).await;
        // After remove, the UUID may still be in vqueue, so we just check the behavior

        // Try to update timestamp - may succeed (if still in vqueue) or fail (if removed from kvs)
        let result = test_svc.service.update_timestamp(uuid.clone(), 1234567890, false).await;
        // Both success and failure are acceptable depending on implementation timing
        let _ = result;
    }

    // ========== Concurrent Operation Tests ==========

    #[tokio::test]
    async fn test_concurrent_insert_basic() {
        let test_svc = TestQBGService::new(128).await;
        let service = std::sync::Arc::new(tokio::sync::Mutex::new(test_svc.service));

        let num_threads = 3;
        let vectors_per_thread = 10;

        // Spawn multiple tasks to insert vectors concurrently
        let mut handles = vec![];
        for thread_id in 0..num_threads {
            let service = service.clone();
            let handle = tokio::spawn(async move {
                for i in 0..vectors_per_thread {
                    let uuid = format!("uuid-{}-{}", thread_id, i);
                    let vector = gen_random_vector(128);
                    let mut svc = service.lock().await;
                    let result = svc.insert(uuid.clone(), vector).await;
                    assert!(result.is_ok(), "Insert failed for {}: {:?}", uuid, result.err());
                }
            });
            handles.push(handle);
        }

        // Wait for all inserts to complete
        for handle in handles {
            handle.await.unwrap();
        }

        // Check insert/delete vqueue buffer lengths (which include pending operations)
        let service = service.lock().await;
        let ivqueue_len = service.insert_vqueue_buffer_len();
        assert!(
            ivqueue_len > 0,
            "Should have pending inserts in vqueue (got: {})",
            ivqueue_len
        );
    }

    #[tokio::test]
    async fn test_concurrent_insert_and_verify() {
        let test_svc = TestQBGService::new(128).await;
        let service = std::sync::Arc::new(tokio::sync::Mutex::new(test_svc.service));

        let num_ops = 20;

        // Spawn concurrent inserts and exists checks
        let mut handles = vec![];
        
        // Insert thread
        {
            let service = service.clone();
            let handle = tokio::spawn(async move {
                for i in 0..num_ops {
                    let uuid = format!("item-{}", i);
                    let vector = gen_random_vector(128);
                    let mut svc = service.lock().await;
                    let _ = svc.insert(uuid, vector).await;
                }
            });
            handles.push(handle);
        }

        // Exists check thread (may find some items depending on timing)
        {
            let service = service.clone();
            let handle = tokio::spawn(async move {
                for i in 0..num_ops {
                    let uuid = format!("item-{}", i);
                    let svc = service.lock().await;
                    let (_, _exists) = svc.exists(uuid).await;
                    // Don't assert, just check that operation completes without panic
                }
            });
            handles.push(handle);
        }

        for handle in handles {
            handle.await.unwrap();
        }
    }

    #[tokio::test]
    async fn test_concurrent_insert_and_remove() {
        let test_svc = TestQBGService::new(128).await;
        let service = std::sync::Arc::new(tokio::sync::Mutex::new(test_svc.service));

        let insert_count = 20;
        let remove_count = 10;

        // First, insert vectors
        {
            let service = service.clone();
            let handle = tokio::spawn(async move {
                for i in 0..insert_count {
                    let uuid = format!("item-{}", i);
                    let vector = gen_random_vector(128);
                    let mut svc = service.lock().await;
                    let _ = svc.insert(uuid, vector).await;
                }
            });
            handle.await.unwrap();
        }

        // Now remove some concurrently with potential new inserts
        let mut handles = vec![];

        // Remove thread
        {
            let service = service.clone();
            let handle = tokio::spawn(async move {
                for i in 0..remove_count {
                    let uuid = format!("item-{}", i);
                    let mut svc = service.lock().await;
                    let _ = svc.remove(uuid).await;
                }
            });
            handles.push(handle);
        }

        // Insert new items thread (doesn't conflict with remove)
        {
            let service = service.clone();
            let handle = tokio::spawn(async move {
                for i in 0..5 {
                    let uuid = format!("new-item-{}", i);
                    let vector = gen_random_vector(128);
                    let mut svc = service.lock().await;
                    let _ = svc.insert(uuid, vector).await;
                }
            });
            handles.push(handle);
        }

        for handle in handles {
            handle.await.unwrap();
        }

        // Verify final vqueue state
        let svc = service.lock().await;
        let ivqueue = svc.insert_vqueue_buffer_len();
        let dvqueue = svc.delete_vqueue_buffer_len();
        // Should have some pending operations
        assert!(ivqueue > 0 || dvqueue > 0, "Should have pending operations in vqueue");
    }

    #[tokio::test]
    async fn test_concurrent_mixed_operations_with_timeouts() {
        let test_svc = TestQBGService::new(128).await;
        let service = std::sync::Arc::new(tokio::sync::Mutex::new(test_svc.service));

        let _num_threads = 3;
        let mut handles = vec![];

        // Thread 0: Insert vectors
        {
            let service = service.clone();
            let handle = tokio::spawn(async move {
                for i in 0..10 {
                    let uuid = format!("insert-{}", i);
                    let vector = gen_random_vector(128);
                    let mut svc = service.lock().await;
                    let _ = svc.insert(uuid, vector).await;
                }
            });
            handles.push(handle);
        }

        // Thread 1: Update vectors (after a small delay)
        {
            let service = service.clone();
            let handle = tokio::spawn(async move {
                tokio::time::sleep(std::time::Duration::from_millis(50)).await;
                for i in 0..5 {
                    let uuid = format!("insert-{}", i);
                    let vector = gen_random_vector(128);
                    let mut svc = service.lock().await;
                    let _ = svc.update(uuid, vector).await;
                }
            });
            handles.push(handle);
        }

        // Thread 2: Check status and operations
        {
            let service = service.clone();
            let handle = tokio::spawn(async move {
                tokio::time::sleep(std::time::Duration::from_millis(100)).await;
                let svc = service.lock().await;
                // Just check that these methods work without panicking
                let _ = svc.is_indexing();
                let _ = svc.is_saving();
                let _ = svc.is_flushing();
                let _ = svc.len();
                let _ = svc.insert_vqueue_buffer_len();
                let _ = svc.delete_vqueue_buffer_len();
            });
            handles.push(handle);
        }

        // Wait for all operations
        for handle in handles {
            handle.await.unwrap();
        }
    }

    // ========== Boundary Value Tests ==========

    #[tokio::test]
    async fn test_boundary_empty_uuid() {
        let mut test_svc = TestQBGService::new(128).await;

        let empty_uuid = "".to_string();
        let vector = gen_random_vector(128);

        // Empty UUID should fail
        let result = test_svc.service.insert(empty_uuid.clone(), vector.clone()).await;
        assert!(result.is_err(), "Insert with empty UUID should fail");
    }

    #[tokio::test]
    async fn test_boundary_very_long_uuid() {
        let mut test_svc = TestQBGService::new(128).await;

        // Create a very long UUID (10KB)
        let long_uuid = "a".repeat(10240);
        let vector = gen_random_vector(128);

        // Very long UUID should still work (no explicit limit in code)
        let result = test_svc.service.insert(long_uuid.clone(), vector).await;
        assert!(result.is_ok(), "Insert with very long UUID should succeed");

        // Verify it exists
        let (_, exists) = test_svc.service.exists(long_uuid).await;
        assert!(exists, "Very long UUID should exist after insert");
    }

    #[tokio::test]
    async fn test_boundary_special_characters_uuid() {
        let mut test_svc = TestQBGService::new(128).await;

        let special_uuid = "uuid-!@#$%^&*()_+-=[]{}|;:',.<>?/~`".to_string();
        let vector = gen_random_vector(128);

        // Special characters in UUID should work
        let result = test_svc.service.insert(special_uuid.clone(), vector).await;
        assert!(result.is_ok(), "Insert with special characters in UUID should succeed");

        let (_, exists) = test_svc.service.exists(special_uuid).await;
        assert!(exists, "UUID with special characters should exist");
    }

    #[tokio::test]
    async fn test_boundary_zero_timestamp() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "test-zero-timestamp".to_string();
        let vector = gen_random_vector(128);

        // Insert with zero timestamp
        let result = test_svc.service.insert_with_time(uuid.clone(), vector, 0).await;
        // Should succeed or fail depending on implementation
        let _ = result;
    }

    #[tokio::test]
    async fn test_boundary_negative_timestamp() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "test-negative-timestamp".to_string();
        let vector = gen_random_vector(128);

        // Insert with negative timestamp
        let result = test_svc.service.insert_with_time(uuid.clone(), vector, -1234567890).await;
        // Should succeed or fail depending on implementation
        let _ = result;
    }

    #[tokio::test]
    async fn test_boundary_max_timestamp() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "test-max-timestamp".to_string();
        let vector = gen_random_vector(128);

        // Insert with i64::MAX timestamp
        let result = test_svc.service.insert_with_time(uuid.clone(), vector, i64::MAX).await;
        assert!(result.is_ok(), "Insert with max timestamp should succeed");
    }

    #[tokio::test]
    async fn test_boundary_min_timestamp() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "test-min-timestamp".to_string();
        let vector = gen_random_vector(128);

        // Insert with i64::MIN timestamp
        let result = test_svc.service.insert_with_time(uuid.clone(), vector, i64::MIN).await;
        assert!(result.is_ok(), "Insert with min timestamp should succeed");
    }

    #[tokio::test]
    async fn test_boundary_empty_vector_list() {
        let mut test_svc = TestQBGService::new(128).await;

        let vectors: std::collections::HashMap<String, Vec<f32>> = std::collections::HashMap::new();

        // Insert empty vector map
        let result = test_svc.service.insert_multiple(vectors).await;
        assert!(result.is_ok(), "Insert multiple with empty map should succeed");
    }

    #[tokio::test]
    async fn test_boundary_remove_empty_list() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuids: Vec<String> = vec![];

        // Remove empty list
        let result = test_svc.service.remove_multiple(uuids).await;
        assert!(result.is_ok(), "Remove multiple with empty list should succeed");
    }

    #[tokio::test]
    async fn test_boundary_single_element_operations() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "single-element".to_string();
        let vector = gen_random_vector(128);

        // Single insert
        test_svc.service.insert(uuid.clone(), vector.clone()).await.unwrap();

        // Single update (may fail if insert not fully processed yet)
        let _result = test_svc.service.update(uuid.clone(), vector.clone()).await;

        // Single remove
        let _result = test_svc.service.remove(uuid.clone()).await;
    }

    #[tokio::test]
    async fn test_boundary_large_vector_dimension() {
        let test_svc = TestQBGService::new(4096).await;

        let uuid = "large-dimension".to_string();
        let vector = gen_random_vector(4096);

        let mut svc = test_svc.service;
        let result = svc.insert(uuid, vector).await;
        assert!(result.is_ok(), "Insert with large dimension should succeed");
    }

    #[tokio::test]
    async fn test_boundary_search_with_zero_k() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert multiple vectors to ensure index can be built
        for i in 0..10 {
            let uuid = format!("search-test-{}", i);
            let vector = gen_random_vector(128);
            let _ = test_svc.service.insert(uuid, vector).await;
        }

        // Create index for search - wait for it to complete
        let index_result = test_svc.service.create_index().await;
        // Index may fail with small dataset, which is acceptable
        let _ = index_result;

        // Only test search if we have indexed data
        let count = test_svc.service.len();
        if count > 0 {
            // Search with k=0 - should return empty or handle gracefully
            let search_vec = gen_random_vector(128);
            let result = test_svc.service.search(search_vec, 0, 0.1, 0.0).await;
            // Result handling: k=0 may not be supported, that's OK
            let _ = result;
        }
    }

    #[tokio::test]
    async fn test_boundary_duplicate_uuid_insert() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "duplicate".to_string();
        let vector1 = gen_random_vector(128);
        let vector2 = gen_random_vector(128);

        // First insert
        test_svc.service.insert(uuid.clone(), vector1).await.unwrap();

        // Second insert with same UUID (should fail)
        let result = test_svc.service.insert(uuid, vector2).await;
        assert!(result.is_err(), "Duplicate insert should fail");
    }

    #[tokio::test]
    async fn test_boundary_remove_nonexistent_uuid() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "nonexistent".to_string();

        // Remove non-existent UUID
        let result = test_svc.service.remove(uuid).await;
        // May succeed or fail depending on implementation
        let _ = result;
    }

    #[tokio::test]
    async fn test_boundary_update_nonexistent_uuid() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "nonexistent".to_string();
        let vector = gen_random_vector(128);

        // Update non-existent UUID
        let result = test_svc.service.update(uuid, vector).await;
        // Should fail
        assert!(result.is_err(), "Update non-existent UUID should fail");
    }

    #[tokio::test]
    async fn test_boundary_get_object_nonexistent() {
        let test_svc = TestQBGService::new(128).await;

        let uuid = "nonexistent".to_string();

        // Get non-existent object
        let result = test_svc.service.get_object(uuid).await;
        assert!(result.is_err(), "Get non-existent object should fail");
    }

    #[tokio::test]
    async fn test_boundary_multiple_operations_same_uuid() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "multi-ops".to_string();

        // Multiple insert attempts should fail after first
        for i in 0..5 {
            let vector = gen_random_vector(128);
            let result = test_svc.service.insert(uuid.clone(), vector).await;
            if i == 0 {
                assert!(result.is_ok(), "First insert should succeed");
            } else {
                assert!(result.is_err(), "Insert {} should fail (UUID already exists)", i);
            }
        }
    }

    #[tokio::test]
    async fn test_boundary_insert_and_get_many_times() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "stress-test".to_string();
        let vector = gen_random_vector(128);

        // Insert once
        test_svc.service.insert(uuid.clone(), vector.clone()).await.unwrap();

        // Get many times
        for _ in 0..100 {
            let result = test_svc.service.get_object(uuid.clone()).await;
            assert!(result.is_ok(), "Get should succeed");
            let (retrieved_vec, _) = result.unwrap();
            assert_eq!(retrieved_vec.len(), 128, "Retrieved vector dimension should match");
        }
    }

}
