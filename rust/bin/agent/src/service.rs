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

pub mod daemon;
pub mod k8s;
pub mod memstore;
pub mod metadata;
pub mod persistence;
mod qbg;
pub use daemon::{DaemonConfig, DaemonHandle, start as start_daemon};
pub use k8s::{IndexMetrics, K8sClient, MetricsExporter, Patcher};
pub use metadata::Metadata;
pub use persistence::{IndexPaths, PersistenceConfig, PersistenceManager};
pub use qbg::QBGService;

#[cfg(test)]
mod tests {
    use std::collections::HashMap;

    use algorithm::Error;
    use proto::payload::v1::{info, search};

    #[derive(Debug)]
    struct _MockService {
        dim: usize,
    }

    impl algorithm::ANN for _MockService {
        // Async search operations
        async fn search(
            &self,
            vector: Vec<f32>,
            _k: u32,
            _epsilon: f32,
            _radius: f32,
        ) -> Result<search::Response, Error> {
            Err(Error::IncompatibleDimensionSize {
                got: vector.len() as usize,
                want: self.dim,
            })
        }

        async fn search_by_id(
            &self,
            _uuid: String,
            _k: u32,
            _epsilon: f32,
            _radius: f32,
        ) -> Result<search::Response, Error> {
            todo!()
        }

        async fn linear_search(
            &self,
            _vector: Vec<f32>,
            _k: u32,
        ) -> Result<search::Response, Error> {
            todo!()
        }

        async fn linear_search_by_id(
            &self,
            _uuid: String,
            _k: u32,
        ) -> Result<search::Response, Error> {
            todo!()
        }

        // Async insert operations
        async fn insert(&mut self, _uuid: String, _vector: Vec<f32>) -> Result<(), Error> {
            todo!()
        }

        async fn insert_with_time(
            &mut self,
            _uuid: String,
            _vector: Vec<f32>,
            _t: i64,
        ) -> Result<(), Error> {
            todo!()
        }

        async fn insert_multiple(
            &mut self,
            _vectors: HashMap<String, Vec<f32>>,
        ) -> Result<(), Error> {
            todo!()
        }

        async fn insert_multiple_with_time(
            &mut self,
            _vectors: HashMap<String, Vec<f32>>,
            _t: i64,
        ) -> Result<(), Error> {
            todo!()
        }

        // Async update operations
        async fn update(&mut self, _uuid: String, _vector: Vec<f32>) -> Result<(), Error> {
            todo!()
        }

        async fn update_with_time(
            &mut self,
            _uuid: String,
            _vector: Vec<f32>,
            _t: i64,
        ) -> Result<(), Error> {
            todo!()
        }

        async fn update_multiple(
            &mut self,
            _vectors: HashMap<String, Vec<f32>>,
        ) -> Result<(), Error> {
            todo!()
        }

        async fn update_multiple_with_time(
            &mut self,
            _vectors: HashMap<String, Vec<f32>>,
            _t: i64,
        ) -> Result<(), Error> {
            todo!()
        }

        async fn update_timestamp(
            &mut self,
            _uuid: String,
            _t: i64,
            _force: bool,
        ) -> Result<(), Error> {
            todo!()
        }

        // Async remove operations
        async fn remove(&mut self, _uuid: String) -> Result<(), Error> {
            todo!()
        }

        async fn remove_with_time(&mut self, _uuid: String, _t: i64) -> Result<(), Error> {
            todo!()
        }

        async fn remove_multiple(&mut self, _uuids: Vec<String>) -> Result<(), Error> {
            todo!()
        }

        async fn remove_multiple_with_time(
            &mut self,
            _uuids: Vec<String>,
            _t: i64,
        ) -> Result<(), Error> {
            todo!()
        }

        // Async index management
        async fn regenerate_indexes(&mut self) -> Result<(), Error> {
            todo!()
        }

        async fn create_index(&mut self) -> Result<(), Error> {
            todo!()
        }

        async fn save_index(&mut self) -> Result<(), Error> {
            todo!()
        }

        async fn create_and_save_index(&mut self) -> Result<(), Error> {
            todo!()
        }

        // Async object retrieval
        async fn get_object(&self, _uuid: String) -> Result<(Vec<f32>, i64), Error> {
            todo!()
        }

        async fn exists(&self, _uuid: String) -> (usize, bool) {
            todo!()
        }

        async fn uuids(&self) -> Vec<String> {
            todo!()
        }

        async fn list_object_func<F: FnMut(String, Vec<f32>, i64) -> bool + Send>(&self, _f: F) {
            todo!()
        }

        async fn close(&mut self) -> Result<(), Error> {
            todo!()
        }

        // Sync status methods
        fn is_indexing(&self) -> bool {
            false
        }

        fn is_flushing(&self) -> bool {
            false
        }

        fn is_saving(&self) -> bool {
            false
        }

        fn len(&self) -> u32 {
            0
        }

        fn number_of_create_index_executions(&self) -> u64 {
            0
        }

        fn insert_vqueue_buffer_len(&self) -> u32 {
            0
        }

        fn delete_vqueue_buffer_len(&self) -> u32 {
            0
        }

        fn get_dimension_size(&self) -> usize {
            self.dim
        }

        fn broken_index_count(&self) -> u64 {
            0
        }

        fn index_statistics(&self) -> Result<info::index::Statistics, Error> {
            todo!()
        }

        fn is_statistics_enabled(&self) -> bool {
            false
        }

        fn index_property(&self) -> Result<info::index::Property, Error> {
            todo!()
        }
    }
}
