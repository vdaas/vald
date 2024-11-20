//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
pub mod property {
    use super::ffi;
    use cxx::UniquePtr;

    pub struct Property {
        inner: UniquePtr<ffi::Property>,
    }

    impl Property {
        pub fn new() -> Self {
            let inner = ffi::new_property();
            Property { inner }
        }

        pub fn init_construction_parameters(&mut self) {
            self.inner.pin_mut().init_qbg_construction_parameters();
        }

        pub fn set_extended_dimension(&mut self, extended_dimension: usize) {
            self.inner.pin_mut().set_extended_dimension(extended_dimension);
        }

        pub fn set_dimension(&mut self, dimension: usize) {
            self.inner.pin_mut().set_dimension(dimension);
        }

        pub fn set_number_of_subvectors(&mut self, number_of_subvectors: usize) {
            self.inner.pin_mut().set_number_of_subvectors(number_of_subvectors);
        }

        pub fn set_number_of_blobs(&mut self, number_of_blobs: usize) {
            self.inner.pin_mut().set_number_of_blobs(number_of_blobs);
        }

        pub fn set_internal_data_type(&mut self, internal_data_type: i32) {
            self.inner.pin_mut().set_internal_data_type(internal_data_type);
        }

        pub fn set_data_type(&mut self, data_type: i32) {
            self.inner.pin_mut().set_data_type(data_type);
        }

        pub fn set_distance_type(&mut self, distance_type: i32) {
            self.inner.pin_mut().set_distance_type(distance_type);
        }

        pub fn init_build_parameters(&mut self) {
            self.inner.pin_mut().init_qbg_build_parameters();
        }

        pub fn set_number_of_objects(&mut self, number_of_objects: usize) {
            self.inner.pin_mut().set_number_of_objects(number_of_objects);
        }
    }
}

pub mod index {
    use super::ffi;
    use super::property::Property;
    use cxx::UniquePtr;

    pub struct Index {
        inner: UniquePtr<ffi::Index>,
    }

    impl Index {
        pub fn new(path: &str, property: &mut Property) -> Result<Self, cxx::Exception> {
            let inner = ffi::new_index(&path.to_string(), property.inner.pin_mut())?;
            Ok(Index { inner })
        }

        pub fn append(&mut self, vector: &[f32]) -> Result<i32, cxx::Exception> {
            self.inner.pin_mut().append(vector)
        }

        pub fn build(&mut self, path: &str, property: &mut Property) {
            self.inner
                .pin_mut()
                .build_index(&path.to_string(), property.inner.pin_mut());
        }

        pub fn save(&mut self) {
            self.inner.pin_mut().save_index();
        }

        pub fn close(&mut self) {
            self.inner.pin_mut().close_index();
        }

        pub fn search(
            &mut self,
            query: &[f32],
            k: usize,
        ) -> Result<(Vec<i32>, Vec<f32>), cxx::Exception> {
            let mut ids = vec![0; k];
            let mut distances = vec![0.0; k];
            unsafe {
                self.inner.pin_mut().search(
                    query,
                    k,
                    ids.as_mut_ptr(),
                    distances.as_mut_ptr(),
                );
            }
            Ok((ids, distances))
        }
    }
}

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
            hierarchical_clustering_init_mode: i16,
        );
        fn set_number_of_first_objects(self: Pin<&mut Property>, number_of_first_objects: usize);
        fn set_number_of_first_clusters(self: Pin<&mut Property>, number_of_first_clusters: usize);
        fn set_number_of_second_objects(self: Pin<&mut Property>, number_of_second_objects: u32);
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

        type Index;
        fn new_index(path: &String, p: Pin<&mut Property>) -> Result<UniquePtr<Index>>;
        fn new_prebuilt_index(path: &String, p: bool) -> Result<UniquePtr<Index>>;
        fn open_index(self: Pin<&mut Index>, path: &String, prebuilt: bool);
        fn build_index(self: Pin<&mut Index>, path: &String, p: Pin<&mut Property>);
        fn save_index(self: Pin<&mut Index>);
        fn close_index(self: Pin<&mut Index>);
        fn append(self: Pin<&mut Index>, v: &[f32]) -> Result<i32>;
        unsafe fn search(
            self: Pin<&mut Index>,
            v: &[f32],
            k: usize,
            ids: *mut i32,
            distances: *mut f32,
        );
    }
}

#[cfg(test)]
mod tests {
    use std::vec;

    use anyhow::Result;

    use crate::ffi;

    #[test]
    fn test_qbg() -> Result<()> {
        let dimension = 128;
        let path = "index".to_string();
        let k = 30;

        let mut p = ffi::new_property();
        //// Test Setter ////
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
        /////////////////////
        p.pin_mut().init_qbg_construction_parameters();
        p.pin_mut().set_dimension(dimension);
        p.pin_mut().set_number_of_subvectors(64);
        p.pin_mut().set_number_of_blobs(0);
        p.pin_mut().init_qbg_build_parameters();
        p.pin_mut().set_number_of_objects(500);

        // New
        println!("create an empty index...");
        let mut index = match ffi::new_index(&path, p.pin_mut()) {
            Ok(index) => index,
            Err(e) => panic!("{}", e),
        };

        // Append & Build
        println!("append objects...");
        for i in 0..100 {
            let mut vec: Vec<f32> = vec![0.0; dimension];
            for d in 0..dimension {
                vec[d] = (i + d) as f32
            }

            let result = index.pin_mut().append(vec.as_slice());
            match result {
                Ok(v) => assert_eq!((i + 1) as i32, v),
                Err(e) => panic!("{}", e),
            }
        }
        index.pin_mut().save_index();
        index.pin_mut().close_index();
        println!("building the index...");
        index.pin_mut().build_index(&path, p.pin_mut());
        index.pin_mut().open_index(&path, true);

        // Search
        let mut ids: Vec<i32> = vec![0; k];
        let mut distances: Vec<f32> = vec![0.0; k];
        let mut vec = vec![0.0; dimension];
        for i in 0..dimension {
            vec[i] = i as f32
        }
        unsafe {
            index.pin_mut().search(
                vec.as_slice(),
                k,
                &mut ids[0] as *mut i32,
                &mut distances[0] as *mut f32,
            )
        };
        println!("ids:\n{:?}", ids);
        println!("distances:\n{:?}", distances);
        index.pin_mut().close_index();

        return Ok(());
    }

    #[test]
    fn test_qbg_reopen() -> Result<()> {
        let dimension = 128;
        let path = "index".to_string();
        let k = 30;

        // New
        let mut index = match ffi::new_prebuilt_index(&path, true) {
            Ok(index) => index,
            Err(e) => panic!("{}", e),
        };

        // Search
        let mut ids: Vec<i32> = vec![0; k];
        let mut distances: Vec<f32> = vec![0.0; k];
        let mut vec = vec![0.0; dimension];
        for i in 0..dimension {
            vec[i] = i as f32
        }
        unsafe {
            index.pin_mut().search(
                vec.as_slice(),
                k,
                &mut ids[0] as *mut i32,
                &mut distances[0] as *mut f32,
            )
        };
        println!("ids:\n{:?}", ids);
        println!("distances:\n{:?}", distances);
        index.pin_mut().close_index();

        return Ok(());
    }

    use crate::{index::Index, property::Property};

    #[test]
    fn test_index_creation() {
        let dimension = 128;
        let path = "index";

        let mut property = Property::new();
        property.init_construction_parameters();
        property.set_dimension(dimension);
        property.set_number_of_subvectors(64);
        property.set_number_of_objects(500);

        let mut index = Index::new(path, &mut property).expect("Failed to create index");

        for i in 0..100 {
            let vec: Vec<f32> = (0..dimension).map(|d| (i + d) as f32).collect();
            index.append(&vec).expect("Failed to append vector");
        }

        index.save();
        index.close();
    }
}
