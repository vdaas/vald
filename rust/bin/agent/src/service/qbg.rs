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
use chrono::{Local, Timelike};
use config::Config;
use kvs::{BidirectionalMap, BidirectionalMapBuilder, MapBase};
use kvs::map::codec::BincodeCodec;
use proto::payload::v1::object::Distance;
use proto::payload::v1::search;
use qbg::index::Index;
use qbg::property::Property;
use vqueue::Queue;

use super::memstore;

pub struct QBGService {
    path: String,
    index: Index,
    property: Property,
    vq: vqueue::PersistentQueue,
    kvs: Arc<BidirectionalMap<String, u32, BincodeCodec>>,
    is_flushing: AtomicBool,
    is_indexing: AtomicBool,
    is_saving: AtomicBool,
    create_index_count: AtomicU64,
    broken_index_count: AtomicU64,
    statistics_enabled: bool,
}

impl QBGService {
    pub async fn new(settings: Config) -> Self {
        let path = settings
            .get::<String>("qbg.index_path")
            .unwrap_or("index".to_string());
        let mut property = Property::new();
        property.init_qbg_construction_parameters();
        property.set_qbg_construction_parameters(
            settings.get::<usize>("qbg.extended_dimension").unwrap_or(0),
            settings.get::<usize>("qbg.dimension").unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_subvectors")
                .unwrap_or(1),
            settings.get::<usize>("qbg.number_of_blobs").unwrap_or(0),
            settings.get::<i32>("qbg.internal_data_type").unwrap_or(1),
            settings.get::<i32>("qbg.data_type").unwrap_or(1),
            settings.get::<i32>("qbg.distance_type").unwrap_or(1),
        );
        property.init_qbg_build_parameters();
        property.set_qbg_build_parameters(
            settings
                .get::<i32>("qbg.hierarchical_clustering_init_mode")
                .unwrap_or(2),
            settings
                .get::<usize>("qbg.number_of_first_objects")
                .unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_first_clusters")
                .unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_second_objects")
                .unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_second_clusters")
                .unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_third_clusters")
                .unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_objects")
                .unwrap_or(1000),
            settings
                .get::<usize>("qbg.number_of_subvectors")
                .unwrap_or(1),
            settings
                .get::<i32>("qbg.optimization_clustering_init_mode")
                .unwrap_or(2),
            settings
                .get::<usize>("qbg.rotation_iteration")
                .unwrap_or(2000),
            settings
                .get::<usize>("qbg.subvector_iteration")
                .unwrap_or(400),
            settings.get::<usize>("qbg.number_of_matrices").unwrap_or(3),
            settings.get::<bool>("qbg.rotation").unwrap_or(true),
            settings.get::<bool>("qbg.repositioning").unwrap_or(false),
        );
        let index = Index::new(&path, &mut property).unwrap();
        let vq_path = settings
            .get::<String>("qbg.vqueue_path")
            .unwrap_or("index".to_string());
        let vq = vqueue::Builder::new(vq_path).build().await.unwrap();
        let kvs_path = settings
            .get::<String>("qbg.kvs_path")
            .unwrap_or("kvs".to_string());
        let kvs = BidirectionalMapBuilder::new(kvs_path)
            .cache_capacity(settings.get::<u64>("qbg.kvs_cache_capacity").unwrap_or(10000))
            .compression_factor(settings.get::<i32>("qbg.kvs_compression_factor").unwrap_or(9))
            .mode(kvs::Mode::HighThroughput)
            .use_compression(settings.get::<bool>("qbg.kvs_use_compression").unwrap_or(true))
            .build()
            .await
            .unwrap();
        QBGService {
            path,
            index,
            property,
            vq,
            kvs,
            is_flushing: AtomicBool::new(false),
            is_indexing: AtomicBool::new(false),
            is_saving: AtomicBool::new(false),
            create_index_count: AtomicU64::new(0),
            broken_index_count: AtomicU64::new(0),
            statistics_enabled: false,
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
        self.ready_for_update(uuid.clone(), vector.clone(), t).await?;
        self.remove_internal(uuid.clone(), t, true).await?;
        self.insert_internal(uuid, vector, t+1, false).await
    }

    async fn remove_internal(&mut self, uuid: String, t: i64, validation: bool) -> Result<(), Error> {
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
    async fn exists(&self, uuid: String) -> (usize, bool) {
        match memstore::exists(&self.kvs, &self.vq, &uuid).await {
            Ok((oid, exists)) => (oid as usize, exists),
            Err(_) => (0, false),
        }
    }

    async fn create_index(&mut self) -> Result<(), Error> {
        // If there are no objects to index, return success
        if self.vq.ivq_len() == 0 {
            self.create_index_count.fetch_add(1, Ordering::SeqCst);
            return Ok(());
        }
        
        self.is_indexing.store(true, Ordering::SeqCst);
        let result = self.index
            .build_index(&self.path, &mut self.property);
        self.is_indexing.store(false, Ordering::SeqCst);
        match result {
            Ok(()) => {
                self.create_index_count.fetch_add(1, Ordering::SeqCst);
                Ok(())
            }
            Err(e) => Err(Error::Internal(Box::new(std::io::Error::other(e.to_string()))))
        }
    }

    async fn save_index(&mut self) -> Result<(), Error> {
        self.is_saving.store(true, Ordering::SeqCst);
        let result = self.index.save_index();
        self.is_saving.store(false, Ordering::SeqCst);
        match result {
            Ok(()) => Ok(()),
            Err(e) => Err(Error::Internal(Box::new(std::io::Error::other(e.to_string()))))
        }
    }

    async fn insert(&mut self, uuid: String, vector: Vec<f32>) -> Result<(), Error> {
        self.insert_internal(uuid, vector, Local::now().nanosecond().into(), true).await
    }

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

    async fn update(&mut self, uuid: String, vector: Vec<f32>) -> Result<(), Error> {
        if self.is_flushing() {
            return Err(Error::FlushingIsInProgress {});
        }
        self.update_internal(uuid, vector, Local::now().nanosecond().into()).await
    }

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

    async fn remove(&mut self, uuid: String) -> Result<(), Error> {
        if self.is_flushing() {
            return Err(Error::FlushingIsInProgress {});
        }
        self.remove_internal(uuid, Local::now().nanosecond().into(), true).await
    }

    async fn remove_multiple(&mut self, uuids: Vec<String>) -> Result<(), Error> {
        if self.is_flushing() {
            return Err(Error::FlushingIsInProgress {});
        }
        self.remove_multiple_internal(uuids, Local::now().nanosecond().into(), true).await
    }

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

    async fn regenerate_indexes(&mut self) -> Result<(), Error> {
        // Close the current index and rebuild it
        self.index.close_index();
        self.create_index().await
    }

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

    async fn close(&mut self) -> Result<(), Error> {
        // Close the QBG index
        self.index.close_index();
        // VQueue and KVS will be cleaned up when dropped
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::TempDir;

    /// Test helper to create a QBGService with temporary directories
    struct TestQBGService {
        service: QBGService,
        _temp_dir: TempDir,
    }

    impl TestQBGService {
        async fn new(dimension: usize) -> Self {
            let temp_dir = TempDir::new().expect("Failed to create temp directory");
            let base_path = temp_dir.path().to_str().unwrap().to_string();

            let settings = Config::builder()
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
                .build()
                .unwrap();

            let service = QBGService::new(settings).await;

            TestQBGService {
                service,
                _temp_dir: temp_dir,
            }
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

        // Create index to move vectors from vqueue to index
        test_svc.service.create_index().await.unwrap();

        assert_eq!(test_svc.service.len(), 10);
    }

    // ========== Create/Save Index Tests ==========

    #[tokio::test]
    async fn test_create_and_save_index() {
        let mut test_svc = TestQBGService::new(128).await;

        // Insert some vectors (QBG needs enough objects)
        for i in 0..100 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }

        // Create and save index
        let result = test_svc.service.create_and_save_index().await;
        assert!(result.is_ok(), "create_and_save_index should succeed: {:?}", result.err());
    }

    // ========== Search By ID Tests ==========

    #[tokio::test]
    async fn test_search_by_id() {
        let mut test_svc = TestQBGService::new(128).await;

        let uuid = "search-by-id-uuid".to_string();
        let vector = gen_random_vector(128);

        test_svc.service.insert(uuid.clone(), vector).await.unwrap();
        test_svc.service.create_index().await.unwrap();

        let result = test_svc.service.search_by_id(uuid, 5, 0.1, -1.0).await;
        assert!(result.is_ok(), "search_by_id should succeed: {:?}", result.err());
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

        // Insert and create index first
        for i in 0..5 {
            test_svc.service.insert(format!("uuid-{}", i), gen_random_vector(128)).await.unwrap();
        }
        test_svc.service.create_index().await.unwrap();

        // Regenerate indexes
        let result = test_svc.service.regenerate_indexes().await;
        assert!(result.is_ok(), "regenerate_indexes should succeed: {:?}", result.err());
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

        // Note: uuids() only returns items that are committed to kvs,
        // not items still in vqueue
        let mut uuids = test_svc.service.uuids().await;
        uuids.sort();

        let mut expected_sorted = expected_uuids.clone();
        expected_sorted.sort();

        assert_eq!(uuids, expected_sorted);
    }

    // ========== Number of Create Index Executions Tests ==========

    #[tokio::test]
    async fn test_number_of_create_index_executions() {
        let mut test_svc = TestQBGService::new(128).await;

        assert_eq!(test_svc.service.number_of_create_index_executions(), 0);

        // Insert and create index
        test_svc.service.insert("uuid-1".to_string(), gen_random_vector(128)).await.unwrap();
        test_svc.service.create_index().await.unwrap();

        assert_eq!(test_svc.service.number_of_create_index_executions(), 1);
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
        let enabled = test_svc.service.is_statistics_enabled();
        assert!(enabled);
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

        assert_eq!(count.load(Ordering::SeqCst), 3);
    }
}
