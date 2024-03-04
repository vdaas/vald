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
        fn new_index(p: Pin<&mut Property>) -> UniquePtr<Index>;
        fn insert(self: Pin<&mut Index>, v: &Vec<f32>) -> u32;
        fn create_index(self: Pin<&mut Index>, pool_size: u32);
        fn remove(self: Pin<&mut Index>, id: u32);
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_works() {
        let result = add(2, 2);
        assert_eq!(result, 4);
    }
}
