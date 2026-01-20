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

pub mod service;
pub use service::qbg::QBGService;

#[cfg(test)]
mod tests {
#[derive(Debug)]
struct _MockService {
    dim: usize,
}

impl algorithm::ANN for _MockService {
    fn search(&self, vector: Vec<f32>, k: u32, epsilon: f32, radius: f32) -> Result<search::Response, Error> {
        Err(Error::IncompatibleDimensionSize {
            got: vector.len() as usize,
            want: self.dim,
        }
        .into())
    }

    fn search_by_id(&self, uuid: String, k: u32, epsilon: f32, radius: f32) -> Result<search::Response, Error> {
        todo!()
    }

    fn linear_search(&self, vector: Vec<f32>, k: u32) -> Result<search::Response, Error> {
        todo!()
    }

    fn linear_search_by_id(&self, uuid: String, k: u32) -> Result<search::Response, Error> {
        todo!()
    }

    fn insert(&mut self, uuid: String, vector: Vec<f32>) -> Result<(), Error> {
        todo!()
    }

    fn insert_with_time(&mut self, uuid: String, vector: Vec<f32>, t: i64) -> Result<(), Error> {
        todo!()
    }

    fn insert_multiple(&mut self, vectors: HashMap<String, Vec<f32>>) -> Result<(), Error> {
        todo!()
    }

    fn insert_multiple_with_time(&mut self, vectors: HashMap<String, Vec<f32>>, t: i64) -> Result<(), Error> {
        todo!()
    }

    fn update(&mut self, uuid: String, vector: Vec<f32>, ts: i64) -> Result<(), Error> {
        todo!()
    }

    fn update_with_time(&mut self, uuid: String, vector: Vec<f32>, t: i64) -> Result<(), Error> {
        todo!()
    }

    fn update_multiple(&mut self, vectors: HashMap<String, Vec<f32>>) -> Result<(), Error> {
        todo!()
    }

    fn update_multiple_with_time(&mut self, vectors: HashMap<String, Vec<f32>>, t: i64) -> Result<(), Error> {
        todo!()
    }

    fn remove(&mut self, uuid: String, ts: i64) -> Result<(), Error> {
        todo!()
    }

    fn remove_with_time(&mut self, uuid: String, t: i64) -> Result<(), Error> {
        todo!()
    }

    fn remove_multiple(&mut self, uuids: Vec<String>) -> Result<(), Error> {
        todo!()
    }

    fn remove_multiple_with_time(&mut self, uuids: Vec<String>, t: i64) -> Result<(), Error> {
        todo!()
    }

    fn regenerate_indexes(&mut self) -> Result<(), Error> {
        todo!()
    }

    fn get_object(&self, uuid: String) -> Result<(Vec<f32>, i64), Error> {
        todo!()
    }

    fn list_object_func<F: Fn(String, Vec<f32>, i64) -> bool>(&self, f: F) {
        todo!()
    }

    fn exists(&self, uuid: String) -> (usize, bool) {
        todo!()
    }

    fn create_index(&mut self) -> Result<(), Error> {
        todo!()
    }

    fn save_index(&mut self) -> Result<(), Error> {
        todo!()
    }

    fn create_and_save_index(&mut self) -> Result<(), Error> {
        todo!()
    }

    fn is_indexing(&self) -> bool {
        todo!()
    }

    fn is_flushing(&self) -> bool {
        todo!()
    }

    fn is_saving(&self) -> bool {
        todo!()
    }

    fn len(&self) -> u32 {
        todo!()
    }

    fn number_of_create_index_executions(&self) -> u64 {
        todo!()
    }

    fn uuids(&self) -> Vec<String> {
        todo!()
    }

    fn insert_vqueue_buffer_len(&self) -> u32 {
        todo!()
    }

    fn delete_vqueue_buffer_len(&self) -> u32 {
        todo!()
    }

    fn get_dimension_size(&self) -> i32 {
        todo!()
    }

    fn broken_index_count(&self) -> u64 {
        todo!()
    }

    fn index_statistics(&self) -> Result<info::index::Statistics, Error> {
        todo!()
    }

    fn is_statistics_enabled(&self) -> bool {
        todo!()
    }

    fn index_property(&self) -> Result<info::index::Property, Error> {
        todo!()
    }

    fn close(&mut self) -> Result<(), Error> {
        todo!()
    }
}
}
