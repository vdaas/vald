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
#[cxx::bridge]
pub mod ffi {
    #[repr(i32)]
    enum ObjectType {
        None = 0,
        Uint8 = 1,
        Float = 2,
        Float16 = 3,
    }

    #[repr(i32)]
    enum DistanceType {
        None = -1,
        L1 = 0,
        L2 = 1,
        Hamming = 2,
        Angle = 3,
        Cosine = 4,
        NormalizedAngle = 5,
        NormalizedCosine = 6,
        Jaccard = 7,
        SparseJaccard = 8,
        NormalizedL2 = 9,
        InnerProduct = 10,
        Poincare = 100,
        Lorentz = 101,
    }

    unsafe extern "C++" {
        include!("ngt-rs/src/input.h");

        type Property;
        fn new_property() -> UniquePtr<Property>;
        fn set_dimension(self: Pin<&mut Property>, dimension: i32);
        fn set_object_type(self: Pin<&mut Property>, t: ObjectType);
        fn set_distance_type(self: Pin<&mut Property>, t: DistanceType);

        type Index;
        fn open_index(path: &String, read_only: bool) -> Result<UniquePtr<Index>>;
        fn new_index(path: &String, p: Pin<&mut Property>) -> Result<UniquePtr<Index>>;
        fn new_index_in_memory(p: Pin<&mut Property>) -> Result<UniquePtr<Index>>;
        unsafe fn search(
            self: Pin<&mut Index>,
            v: &[f32],
            size: i32,
            epsilon: f32,
            radius: f32,
            edge_size: i32,
            ids: *mut i32,
            distances: *mut f32,
        ) -> Result<()>;
        unsafe fn linear_search(
            self: Pin<&mut Index>,
            v: &[f32],
            size: i32,
            edge_size: i32,
            ids: *mut i32,
            distances: *mut f32,
        ) -> Result<()>;
        fn insert(self: Pin<&mut Index>, v: &[f32]) -> Result<u32>;
        fn create_index(self: Pin<&mut Index>, pool_size: u32) -> Result<()>;
        fn remove(self: Pin<&mut Index>, id: u32) -> Result<()>;
        fn get_vector(self: Pin<&mut Index>, id: u32) -> Result<&[f32]>;
    }
}

#[cfg(test)]
mod tests {
    use std::vec;

    use anyhow::Result;
    use rand::distributions::Standard;
    use rand::prelude::*;

    use super::*;

    const DIMENSION: i32 = 128;
    const COUNT: u32 = 1000;
    const K: usize = 10;

    fn gen_random_vector(dim: i32) -> Vec<f32> {
        (0..dim)
            .map(|_| StdRng::from_entropy().sample(Standard))
            .collect()
    }

    #[test]
    fn test_ngt() -> Result<()> {
        let mut p = ffi::new_property();
        p.pin_mut().set_dimension(DIMENSION);
        p.pin_mut().set_distance_type(ffi::DistanceType::L2);
        p.pin_mut().set_object_type(ffi::ObjectType::Float);

        let mut index = ffi::new_index_in_memory(p.pin_mut())?;
        let vectors: Vec<Vec<f32>> = (0..COUNT).map(|_| gen_random_vector(DIMENSION)).collect();
        for (i, v) in vectors.iter().enumerate() {
            let id = index.pin_mut().insert(v.as_slice())?;
            assert_eq!(i + 1, id as usize);
        }
        index.pin_mut().create_index(4)?;

        for _ in 0..COUNT {
            let mut ids: Vec<i32> = vec![-1; K];
            let mut distances: Vec<f32> = vec![-1.0; K];
            unsafe {
                index.pin_mut().search(
                    gen_random_vector(DIMENSION).as_slice(),
                    K as i32,
                    0.05,
                    -1.0,
                    i32::MIN,
                    &mut ids[0] as *mut i32,
                    &mut distances[0] as *mut f32,
                )?
            };
            for i in 0..K {
                assert!(
                    1 <= ids[i] && ids[i] <= COUNT as i32,
                    "invalid id {}",
                    ids[i]
                );
                assert!(distances[i] >= 0.0, "invalid distance {}", distances[i]);
            }
        }

        for (i, v) in vectors.iter().enumerate() {
            let ret = index.pin_mut().get_vector((i + 1) as u32)?;
            assert_eq!(v.as_slice(), ret);
        }

        for i in 1..COUNT + 1 {
            index.pin_mut().remove(i)?;
        }
        Ok(())
    }
}
