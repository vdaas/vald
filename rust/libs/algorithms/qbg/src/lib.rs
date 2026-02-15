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
#[cxx::bridge]
pub mod ffi {
    unsafe extern "C++" {
        include!("qbg/src/input.h");

        type Property;
        fn new_property() -> UniquePtr<Property>;
        fn init_qbg_construction_parameters(self: Pin<&mut Property>);
        fn set_qbg_construction_parameters(
            self: Pin<&mut Property>,
            extended_dimension: usize,
            dimension: usize,
            number_of_subvectors: usize,
            number_of_blobs: usize,
            internal_data_type: i32,
            data_type: i32,
            distance_type: i32,
        );
        fn set_extended_dimension(self: Pin<&mut Property>, extended_dimension: usize);
        fn set_dimension(self: Pin<&mut Property>, dimension: usize);
        fn set_number_of_subvectors(self: Pin<&mut Property>, number_of_subvectors: usize);
        fn set_number_of_blobs(self: Pin<&mut Property>, number_of_blobs: usize);
        fn set_internal_data_type(self: Pin<&mut Property>, internal_data_type: i32);
        fn set_data_type(self: Pin<&mut Property>, data_type: i32);
        fn set_distance_type(self: Pin<&mut Property>, distance_type: i32);
        fn init_qbg_build_parameters(self: Pin<&mut Property>);
        fn set_qbg_build_parameters(
            self: Pin<&mut Property>,
            hierarchical_clustering_init_mode: i32,
            number_of_first_objects: usize,
            number_of_first_clusters: usize,
            number_of_second_objects: usize,
            number_of_second_clusters: usize,
            number_of_third_clusters: usize,
            number_of_objects: usize,
            number_of_subvectors: usize,
            optimization_clustering_init_mode: i32,
            rotation_iteration: usize,
            subvector_iteration: usize,
            number_of_matrices: usize,
            rotation: bool,
            repositioning: bool,
        );
        fn set_hierarchical_clustering_init_mode(
            self: Pin<&mut Property>,
            hierarchical_clustering_init_mode: i32,
        );
        fn set_number_of_first_objects(self: Pin<&mut Property>, number_of_first_objects: usize);
        fn set_number_of_first_clusters(self: Pin<&mut Property>, number_of_first_clusters: usize);
        fn set_number_of_second_objects(self: Pin<&mut Property>, number_of_second_objects: usize);
        fn set_number_of_second_clusters(
            self: Pin<&mut Property>,
            number_of_second_clusters: usize,
        );
        fn set_number_of_third_clusters(self: Pin<&mut Property>, number_of_third_clusters: usize);
        fn set_number_of_objects(self: Pin<&mut Property>, number_of_objects: usize);
        fn set_number_of_subvectors_for_bp(self: Pin<&mut Property>, number_of_subvectors: usize);
        fn set_optimization_clustering_init_mode(
            self: Pin<&mut Property>,
            optimization_clustering_init_mode: i32,
        );
        fn set_rotation_iteration(self: Pin<&mut Property>, rotation_iteration: usize);
        fn set_subvector_iteration(self: Pin<&mut Property>, subvector_iteration: usize);
        fn set_number_of_matrices(self: Pin<&mut Property>, number_of_matrices: usize);
        fn set_rotation(self: Pin<&mut Property>, rotation: bool);
        fn set_repositioning(self: Pin<&mut Property>, repositioning: bool);

        type SearchResult;
        fn get_id(self: Pin<&mut SearchResult>) -> u32;
        fn get_distance(self: Pin<&mut SearchResult>) -> f32;

        type Index;
        fn new_index(path: &String, p: Pin<&mut Property>) -> Result<UniquePtr<Index>>;
        fn new_prebuilt_index(path: &String, p: bool) -> Result<UniquePtr<Index>>;
        fn open_index(self: Pin<&mut Index>, path: &String, prebuilt: bool) -> Result<()>;
        fn build_index(self: Pin<&mut Index>, path: &String, p: Pin<&mut Property>) -> Result<()>;
        fn save_index(self: Pin<&mut Index>) -> Result<()>;
        fn close_index(self: Pin<&mut Index>);
        fn append(self: Pin<&mut Index>, v: &[f32]) -> Result<i32>;
        fn insert(self: Pin<&mut Index>, v: &[f32]) -> Result<i32>;
        fn remove(self: Pin<&mut Index>, id: usize) -> Result<()>;
        fn search(
            self: &Index,
            v: &[f32],
            k: usize,
            radius: f32,
            epsilon: f32,
        ) -> Result<UniquePtr<CxxVector<SearchResult>>>;
        fn get_object(self: &Index, id: usize) -> Result<*mut f32>;
        fn get_dimension(self: &Index) -> Result<usize>;
    }
}

unsafe impl Sync for ffi::Property {}
unsafe impl Send for ffi::Property {}
unsafe impl Sync for ffi::Index {}
unsafe impl Send for ffi::Index {}

pub mod property {
    use super::ffi;
    use cxx::UniquePtr;
    use std::pin::Pin;

    pub struct Property {
        inner: UniquePtr<ffi::Property>,
    }

    impl Property {
        pub fn new() -> Self {
            let inner = ffi::new_property();
            Property { inner }
        }

        pub fn get_property(&mut self) -> Pin<&mut ffi::Property> {
            self.inner.pin_mut()
        }

        pub fn init_qbg_construction_parameters(&mut self) {
            self.inner.pin_mut().init_qbg_construction_parameters()
        }

        pub fn set_qbg_construction_parameters(
            &mut self,
            extended_dimension: usize,
            dimension: usize,
            number_of_subvectors: usize,
            number_of_blobs: usize,
            internal_data_type: i32,
            data_type: i32,
            distance_type: i32,
        ) {
            self.inner.pin_mut().set_qbg_construction_parameters(
                extended_dimension,
                dimension,
                number_of_subvectors,
                number_of_blobs,
                internal_data_type,
                data_type,
                distance_type,
            )
        }

        pub fn set_extended_dimension(&mut self, extended_dimension: usize) {
            self.inner
                .pin_mut()
                .set_extended_dimension(extended_dimension)
        }

        pub fn set_dimension(&mut self, dimension: usize) {
            self.inner.pin_mut().set_dimension(dimension)
        }

        pub fn set_number_of_subvectors(&mut self, number_of_subvectors: usize) {
            self.inner
                .pin_mut()
                .set_number_of_subvectors(number_of_subvectors)
        }

        pub fn set_number_of_blobs(&mut self, number_of_blobs: usize) {
            self.inner.pin_mut().set_number_of_blobs(number_of_blobs)
        }

        pub fn set_internal_data_type(&mut self, internal_data_type: i32) {
            self.inner
                .pin_mut()
                .set_internal_data_type(internal_data_type)
        }

        pub fn set_data_type(&mut self, data_type: i32) {
            self.inner.pin_mut().set_data_type(data_type)
        }

        pub fn set_distance_type(&mut self, distance_type: i32) {
            self.inner.pin_mut().set_distance_type(distance_type)
        }

        pub fn init_qbg_build_parameters(&mut self) {
            self.inner.pin_mut().init_qbg_build_parameters()
        }

        pub fn set_qbg_build_parameters(
            &mut self,
            hierarchical_clustering_init_mode: i32,
            number_of_first_objects: usize,
            number_of_first_clusters: usize,
            number_of_second_objects: usize,
            number_of_second_clusters: usize,
            number_of_third_clusters: usize,
            number_of_objects: usize,
            number_of_subvectors: usize,
            optimization_clustering_init_mode: i32,
            rotation_iteration: usize,
            subvector_iteration: usize,
            number_of_matrices: usize,
            rotation: bool,
            repositioning: bool,
        ) {
            self.inner.pin_mut().set_qbg_build_parameters(
                hierarchical_clustering_init_mode,
                number_of_first_objects,
                number_of_first_clusters,
                number_of_second_objects,
                number_of_second_clusters,
                number_of_third_clusters,
                number_of_objects,
                number_of_subvectors,
                optimization_clustering_init_mode,
                rotation_iteration,
                subvector_iteration,
                number_of_matrices,
                rotation,
                repositioning,
            )
        }

        pub fn set_hierarchical_clustering_init_mode(
            &mut self,
            hierarchical_clustering_init_mode: i32,
        ) {
            self.inner
                .pin_mut()
                .set_hierarchical_clustering_init_mode(hierarchical_clustering_init_mode)
        }

        pub fn set_number_of_first_objects(&mut self, number_of_first_objects: usize) {
            self.inner
                .pin_mut()
                .set_number_of_first_objects(number_of_first_objects)
        }

        pub fn set_number_of_first_clusters(&mut self, number_of_first_clusters: usize) {
            self.inner
                .pin_mut()
                .set_number_of_first_clusters(number_of_first_clusters)
        }

        pub fn set_number_of_second_objects(&mut self, number_of_second_objects: usize) {
            self.inner
                .pin_mut()
                .set_number_of_second_objects(number_of_second_objects)
        }

        pub fn set_number_of_second_clusters(&mut self, number_of_second_clusters: usize) {
            self.inner
                .pin_mut()
                .set_number_of_second_clusters(number_of_second_clusters)
        }

        pub fn set_number_of_third_clusters(&mut self, number_of_third_clusters: usize) {
            self.inner
                .pin_mut()
                .set_number_of_third_clusters(number_of_third_clusters)
        }

        pub fn set_number_of_objects(&mut self, number_of_objects: usize) {
            self.inner
                .pin_mut()
                .set_number_of_objects(number_of_objects)
        }

        pub fn set_number_of_subvectors_for_bp(&mut self, number_of_subvectors: usize) {
            self.inner
                .pin_mut()
                .set_number_of_subvectors_for_bp(number_of_subvectors)
        }

        pub fn set_optimization_clustering_init_mode(
            &mut self,
            optimization_clustering_init_mode: i32,
        ) {
            self.inner
                .pin_mut()
                .set_optimization_clustering_init_mode(optimization_clustering_init_mode)
        }

        pub fn set_rotation_iteration(&mut self, rotation_iteration: usize) {
            self.inner
                .pin_mut()
                .set_rotation_iteration(rotation_iteration)
        }

        pub fn set_subvector_iteration(&mut self, subvector_iteration: usize) {
            self.inner
                .pin_mut()
                .set_subvector_iteration(subvector_iteration)
        }

        pub fn set_number_of_matrices(&mut self, number_of_matrices: usize) {
            self.inner
                .pin_mut()
                .set_number_of_matrices(number_of_matrices)
        }

        pub fn set_rotation(&mut self, rotation: bool) {
            self.inner.pin_mut().set_rotation(rotation)
        }

        pub fn set_repositioning(&mut self, repositioning: bool) {
            self.inner.pin_mut().set_repositioning(repositioning)
        }
    }
}

pub mod index {
    use super::ffi;
    use super::property;
    use core::slice;
    use cxx::UniquePtr;

    pub struct Index {
        inner: UniquePtr<ffi::Index>,
    }

    impl Index {
        pub fn new(path: &String, p: &mut property::Property) -> Result<Self, cxx::Exception> {
            let inner = ffi::new_index(path, p.get_property())?;
            Ok(Index { inner })
        }

        pub fn new_prebuilt(path: &String, p: bool) -> Result<Self, cxx::Exception> {
            let inner = ffi::new_prebuilt_index(path, p)?;
            Ok(Index { inner })
        }

        pub fn open_index(&mut self, path: &String, prebuilt: bool) -> Result<(), cxx::Exception> {
            self.inner.pin_mut().open_index(path, prebuilt)
        }

        pub fn build_index(
            &mut self,
            path: &String,
            p: &mut property::Property,
        ) -> Result<(), cxx::Exception> {
            self.inner.pin_mut().build_index(path, p.get_property())
        }

        pub fn save_index(&mut self) -> Result<(), cxx::Exception> {
            self.inner.pin_mut().save_index()
        }

        pub fn close_index(&mut self) {
            self.inner.pin_mut().close_index()
        }

        pub fn append(&mut self, v: &[f32]) -> Result<i32, cxx::Exception> {
            self.inner.pin_mut().append(v)
        }

        pub fn insert(&mut self, v: &[f32]) -> Result<i32, cxx::Exception> {
            self.inner.pin_mut().insert(v)
        }

        pub fn remove(&mut self, id: usize) -> Result<(), cxx::Exception> {
            self.inner.pin_mut().remove(id)
        }

        pub fn search(
            &self,
            v: &[f32],
            k: usize,
            radius: f32,
            epsilon: f32,
        ) -> Result<Vec<(u32, f32)>, cxx::Exception> {
            let index = self.inner.as_ref().unwrap();
            let mut search_results = index.search(v, k, radius, epsilon)?;
            Ok(search_results
                .pin_mut()
                .into_iter()
                .map(|mut s| (s.as_mut().get_id(), s.as_mut().get_distance()))
                .collect())
        }

        pub fn get_object(&self, id: usize) -> Result<&[f32], cxx::Exception> {
            let dim = self.inner.get_dimension()?;
            match self.inner.get_object(id) {
                Ok(v) => Ok(unsafe { slice::from_raw_parts(v, dim) }),
                Err(e) => Err(e),
            }
        }

        pub fn get_dimension(&self) -> Result<usize, cxx::Exception> {
            let index = self.inner.as_ref().unwrap();
            index.get_dimension()
        }
    }
}

#[cfg(test)]
mod tests {
    use crate::{ffi, index::Index, property::Property};
    use anyhow::Result;
    use tempfile::tempdir;

    const DIMENSION: usize = 128;
    const K: usize = 30;
    const RADIUS: f32 = 0.0;
    const EPSILON: f32 = 0.1;

    #[test]
    fn test_ffi_qbg() -> Result<()> {
        // New
        println!("create an empty index...");
        let temp_dir = tempdir()?;
        let path = temp_dir.path().join("index").to_string_lossy().to_string();
        let mut p = ffi::new_property();
        ////////// Test Setter //////////
        p.pin_mut().set_extended_dimension(1);
        p.pin_mut().set_dimension(1);
        p.pin_mut().set_number_of_subvectors(1);
        p.pin_mut().set_number_of_blobs(1);
        p.pin_mut().set_internal_data_type(1);
        p.pin_mut().set_data_type(1);
        p.pin_mut().set_distance_type(1);
        p.pin_mut().set_hierarchical_clustering_init_mode(1);
        p.pin_mut().set_number_of_first_objects(1);
        p.pin_mut().set_number_of_first_clusters(1);
        p.pin_mut().set_number_of_second_objects(1);
        p.pin_mut().set_number_of_second_clusters(1);
        p.pin_mut().set_number_of_third_clusters(1);
        p.pin_mut().set_number_of_objects(1);
        p.pin_mut().set_number_of_subvectors_for_bp(1);
        p.pin_mut().set_optimization_clustering_init_mode(1);
        p.pin_mut().set_rotation_iteration(1);
        p.pin_mut().set_subvector_iteration(1);
        p.pin_mut().set_number_of_matrices(1);
        p.pin_mut().set_rotation(false);
        p.pin_mut().set_repositioning(false);
        ////////// /////////// //////////
        p.pin_mut().init_qbg_construction_parameters();
        p.pin_mut().set_dimension(DIMENSION);
        p.pin_mut().set_number_of_subvectors(64);
        p.pin_mut().set_number_of_blobs(0);
        p.pin_mut().init_qbg_build_parameters();
        p.pin_mut().set_number_of_objects(500);
        let mut index = ffi::new_index(&path, p.pin_mut()).unwrap();

        // Append
        println!("append objects...");
        for i in 0..100 {
            let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|x| (x + i) as f32).collect();
            let id = index.pin_mut().append(vec.as_slice()).unwrap();
            assert_eq!((i + 1) as i32, id)
        }
        index.pin_mut().save_index().unwrap();
        index.pin_mut().close_index();

        // Build
        println!("building the index...");
        index.pin_mut().build_index(&path, p.pin_mut()).unwrap();
        index.pin_mut().open_index(&path, true).unwrap();

        // Insert
        for i in 0..100 {
            let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|x| (x + i) as f32).collect();
            let id = index.pin_mut().insert(vec.as_slice()).unwrap();
            assert_eq!((i + 1 + 100) as i32, id)
        }

        // Get Object
        let vec = index.pin_mut().get_object(1).unwrap();
        println!("vec:\n\t{:?}", vec);

        // Get Dimension
        let dim = index.pin_mut().get_dimension().unwrap();
        println!("dimension:\n\t{:?}", dim);

        // Search
        println!("search the index for the specified query...");
        let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|i| i as f32).collect();
        let mut search_results = index.pin_mut().search(vec.as_slice(), K, RADIUS, EPSILON);
        let ids: Vec<u32> = search_results
            .pin_mut()
            .into_iter()
            .map(|s| s.get_id())
            .collect();
        let distances: Vec<f32> = search_results
            .pin_mut()
            .into_iter()
            .map(|s| s.get_distance())
            .collect();
        println!("ids:\n\t{:?}", ids);
        println!("distances:\n\t{:?}", distances);

        // Remove
        index.pin_mut().remove(1).unwrap();
        let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|i| i as f32).collect();
        let mut search_results = index.pin_mut().search(vec.as_slice(), K, RADIUS, EPSILON);
        let ids: Vec<u32> = search_results
            .pin_mut()
            .into_iter()
            .map(|s| s.get_id())
            .collect();
        let distances: Vec<f32> = search_results
            .pin_mut()
            .into_iter()
            .map(|s| s.get_distance())
            .collect();
        println!("ids:\n\t{:?}", ids);
        println!("distances:\n\t{:?}", distances);

        index.pin_mut().close_index();

        Ok(())
    }

    #[test]
    fn test_ffi_qbg_prebuilt() -> Result<()> {
        // First create an index for this test
        let temp_dir = tempdir()?;
        let path = temp_dir.path().join("index").to_string_lossy().to_string();

        // Create and build a fresh index
        let mut p = ffi::new_property();
        p.pin_mut().init_qbg_construction_parameters();
        p.pin_mut().set_dimension(DIMENSION);
        p.pin_mut().set_number_of_subvectors(64);
        p.pin_mut().set_number_of_blobs(0);
        p.pin_mut().init_qbg_build_parameters();
        p.pin_mut().set_number_of_objects(500);
        let mut index = ffi::new_index(&path, p.pin_mut())?;

        // Append some objects
        for i in 0..100 {
            let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|x| (x + i) as f32).collect();
            index.pin_mut().append(vec.as_slice())?;
        }
        index.pin_mut().save_index()?;
        index.pin_mut().close_index();

        // Build the index
        index.pin_mut().build_index(&path, p.pin_mut())?;

        // Now test with prebuilt index
        let mut index = ffi::new_prebuilt_index(&path, true).unwrap();

        // Insert
        for i in 0..100 {
            let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|x| (x + i) as f32).collect();
            let id = index.pin_mut().insert(vec.as_slice()).unwrap();
            assert_eq!((i + 1 + 100) as i32, id)
        }

        // Get Object
        let vec = index.pin_mut().get_object(1).unwrap();
        println!("vec:\n\t{:?}", vec);

        // Get Dimension
        let dim = index.pin_mut().get_dimension().unwrap();
        println!("dimension:\n\t{:?}", dim);

        // Search
        let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|i| i as f32).collect();
        let mut search_results = index.pin_mut().search(vec.as_slice(), K, RADIUS, EPSILON);
        let ids: Vec<u32> = search_results
            .pin_mut()
            .into_iter()
            .map(|s| s.get_id())
            .collect();
        let distances: Vec<f32> = search_results
            .pin_mut()
            .into_iter()
            .map(|s| s.get_distance())
            .collect();
        println!("ids:\n\t{:?}", ids);
        println!("distances:\n\t{:?}", distances);

        // Remove
        index.pin_mut().remove(1).unwrap();
        let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|i| i as f32).collect();
        let mut search_results = index.pin_mut().search(vec.as_slice(), K, RADIUS, EPSILON);
        let ids: Vec<u32> = search_results
            .pin_mut()
            .into_iter()
            .map(|s| s.get_id())
            .collect();
        let distances: Vec<f32> = search_results
            .pin_mut()
            .into_iter()
            .map(|s| s.get_distance())
            .collect();
        println!("ids:\n\t{:?}", ids);
        println!("distances:\n\t{:?}", distances);

        index.pin_mut().close_index();

        Ok(())
    }

    #[test]
    fn test_property() -> Result<()> {
        let mut p = Property::new();
        p.init_qbg_construction_parameters();
        p.set_qbg_construction_parameters(1, 1, 1, 1, 1, 1, 1);
        p.set_extended_dimension(1);
        p.set_dimension(1);
        p.set_number_of_subvectors(1);
        p.set_number_of_blobs(1);
        p.set_internal_data_type(1);
        p.set_data_type(1);
        p.set_distance_type(1);
        p.init_qbg_build_parameters();
        p.set_qbg_build_parameters(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, true, false);
        p.set_hierarchical_clustering_init_mode(1);
        p.set_number_of_first_objects(1);
        p.set_number_of_first_clusters(1);
        p.set_number_of_second_objects(1);
        p.set_number_of_second_clusters(1);
        p.set_number_of_third_clusters(1);
        p.set_number_of_objects(1);
        p.set_number_of_subvectors_for_bp(1);
        p.set_optimization_clustering_init_mode(1);
        p.set_rotation_iteration(1);
        p.set_subvector_iteration(1);
        p.set_number_of_matrices(1);
        p.set_rotation(false);
        p.set_repositioning(false);

        Ok(())
    }

    #[test]
    fn test_index() -> Result<()> {
        // New
        println!("create an empty index...");
        let temp_dir = tempdir()?;
        let path = temp_dir.path().join("index").to_string_lossy().to_string();
        let mut p = Property::new();
        p.init_qbg_construction_parameters();
        p.set_dimension(DIMENSION);
        p.set_number_of_subvectors(64);
        p.set_number_of_blobs(0);
        p.init_qbg_build_parameters();
        p.set_number_of_objects(500);
        let mut index = Index::new(&path, &mut p).unwrap();

        // Append
        println!("append objects...");
        for i in 0..100 {
            let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|x| (x + i) as f32).collect();
            let id = index.append(vec.as_slice()).unwrap();
            assert_eq!((i + 1) as i32, id)
        }
        index.save_index().unwrap();
        index.close_index();

        // Build
        println!("building the index...");
        index.build_index(&path, &mut p).unwrap();
        index.open_index(&path, true).unwrap();

        // Insert
        for i in 0..100 {
            let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|x| (x + i) as f32).collect();
            let id = index.insert(vec.as_slice()).unwrap();
            assert_eq!((i + 1 + 100) as i32, id)
        }

        // Get Object
        let vec = index.get_object(1).unwrap();
        println!("vec:\n\t{:?}", vec);

        // Get Dimension
        let dim = index.get_dimension().unwrap();
        println!("dimension:\n\t{:?}", dim);

        // Search
        println!("search the index for the specified query...");
        let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|i| i as f32).collect();
        let search_results = index.search(vec.as_slice(), K, RADIUS, EPSILON).unwrap();
        let ids: Vec<u32> = search_results.iter().map(|s| s.0).collect();
        let distances: Vec<f32> = search_results.iter().map(|s| s.1).collect();
        println!("search results:\n\t{:?}", search_results);
        println!("ids:\n\t{:?}", ids);
        println!("distances:\n\t{:?}", distances);

        // Remove
        index.remove(1).unwrap();
        let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|i| i as f32).collect();
        let search_results = index.search(vec.as_slice(), K, RADIUS, EPSILON).unwrap();
        let ids: Vec<u32> = search_results.iter().map(|s| s.0).collect();
        let distances: Vec<f32> = search_results.iter().map(|s| s.1).collect();
        println!("search results:\n\t{:?}", search_results);
        println!("ids:\n\t{:?}", ids);
        println!("distances:\n\t{:?}", distances);

        index.close_index();

        Ok(())
    }
}
