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

//! QBG (Quantized Blob Graph) ANN Algorithm Wrapper for Vald.
//!
//! This library provides a **Rust wrapper for the C++ QBG library**, enabling high-performance
//! approximate nearest neighbor (ANN) search with graph-based indexing. The wrapper abstracts
//! the complexity of C++ FFI while maintaining full access to QBG's performance optimizations
//! and advanced configuration options.
//!
//! # What is QBG?
//!
//! QBG (Quantized Blob Graph) is an efficient ANN algorithm that:
//! - Uses hierarchical clustering and graph-based indexing
//! - Supports multiple distance metrics (L1, L2, Hamming, Angle, Cosine)
//! - Handles various data types (uint8, float, float16)
//! - Provides fast approximate search with configurable accuracy/speed tradeoffs
//! - Optimizes for AVX-512 and AVX-2 CPU instructions for maximum performance
//!
//! # C++ Integration
//!
//! This crate wraps the C++ QBG implementation from the `qbg-sys` crate, which provides:
//! - Safe FFI bindings to the QBG C++ library
//! - Memory management and pointer handling
//! - Support for prebuilt and freshly created indexes
//! - Atomic operations for thread-safe updates
//!
//! # Core Components
//!
//! - **`Index`** - The main entry point for QBG operations (create, search, insert, etc.)
//! - **`Property`** - Configuration for index construction (dimension, clustering parameters, etc.)
//! - **`ObjectType` / `DataType` / `DistanceType`** - Enums for type safety and serialization
//! - **`Result`** - Error handling wrapper around QBG operations
//!
//! # Safety Considerations
//!
//! Every `unsafe` block in this library is documented with `// SAFETY:` comments explaining:
//! - Why unsafe code is necessary (C++ interop, memory management)
//! - How memory safety is guaranteed
//! - What invariants must be upheld
//!
//! # Example Usage
//!
//! ```ignore
//! use qbg::Index;
//! use qbg::Property;
//!
//! // Create or load an index
//! let mut property = Property::new();
//! property.set_qbg_construction_parameters(
//!     512,      // extended_dimension
//!     512,      // dimension
//!     8,        // number_of_subvectors
//!     10000,    // number_of_blobs
//!     ObjectType::Float,
//!     DataType::Float,
//!     DistanceType::L2,
//! );
//!
//! let index = Index::new("path/to/index", &mut property)?;
//!
//! // Insert vectors
//! let vector = vec![0.1, 0.2, 0.3, /* ... */];
//! index.insert(0, &vector)?;
//!
//! // Search
//! let results = index.search(&vector, 10)?;
//! ```

use serde::{Deserialize, Serialize};

/// Data type for internal vector representation in the index.
///
/// This enum specifies how vector components are represented in the quantized index structure.
/// It affects memory usage, precision, and computational efficiency.
///
/// # Variants
///
/// * `None` - Invalid type
/// * `Uint8` - 8-bit unsigned integer quantization. Provides maximum compression but lowest precision.
/// * `Float` - 32-bit floating-point. Full precision but higher memory usage.
/// * `Float16` - 16-bit half-precision floating-point. Good balance between precision and compression.
#[derive(Serialize, Deserialize, Debug, Clone, Copy, PartialEq, Eq)]
pub enum ObjectType {
    /// Invalid type.
    #[serde(rename = "None", alias = "none")]
    None,
    /// 8-bit unsigned integer quantization.
    #[serde(rename = "uint8", alias = "Uint8", alias = "u8", alias = "U8")]
    Uint8,
    /// 32-bit floating-point representation.
    #[serde(rename = "float", alias = "Float", alias = "f32", alias = "F32")]
    Float,
    /// 16-bit half-precision floating-point.
    #[serde(rename = "float16", alias = "Float16", alias = "f16", alias = "F16")]
    Float16,
}

impl From<ffi::ObjectType> for ObjectType {
    fn from(value: ffi::ObjectType) -> Self {
        match value {
            ffi::ObjectType::Uint8 => ObjectType::Uint8,
            ffi::ObjectType::Float => ObjectType::Float,
            ffi::ObjectType::Float16 => ObjectType::Float16,
            _ => ObjectType::None,
        }
    }
}

impl From<ObjectType> for ffi::ObjectType {
    fn from(value: ObjectType) -> Self {
        match value {
            ObjectType::Uint8 => ffi::ObjectType::Uint8,
            ObjectType::Float => ffi::ObjectType::Float,
            ObjectType::Float16 => ffi::ObjectType::Float16,
            _ => ffi::ObjectType::None,
        }
    }
}

/// Data type for the input vectors before quantization.
///
/// This enum specifies the original format of the vectors provided to the index.
/// The index will handle type conversion and quantization as needed.
///
/// # Variants
///
/// * `None` - Invalid type
/// * `Uint8` - 8-bit unsigned integer vectors. Useful for binary/categorical data.
/// * `Float` - 32-bit floating-point vectors. Standard format for most applications.
/// * `Float16` - 16-bit half-precision floating-point vectors.
/// * `Any` - Accept vectors in any supported format. Useful for flexible implementations.
#[derive(Serialize, Deserialize, Debug, Clone, Copy, PartialEq, Eq)]
pub enum DataType {
    /// Invalid type.
    #[serde(rename = "None", alias = "none")]
    None,
    /// 8-bit unsigned integer input vectors.
    #[serde(rename = "uint8", alias = "Uint8", alias = "u8", alias = "U8")]
    Uint8,
    /// 32-bit floating-point input vectors.
    #[serde(rename = "float", alias = "Float", alias = "f32", alias = "F32")]
    Float,
    /// 16-bit half-precision floating-point input vectors.
    #[serde(rename = "float16", alias = "Float16", alias = "f16", alias = "F16")]
    Float16,
    /// Accept vectors in any supported format.
    #[serde(rename = "any", alias = "Any")]
    Any,
}

impl From<ffi::DataType> for DataType {
    fn from(value: ffi::DataType) -> Self {
        match value {
            ffi::DataType::Uint8 => DataType::Uint8,
            ffi::DataType::Float => DataType::Float,
            ffi::DataType::Float16 => DataType::Float16,
            ffi::DataType::Any => DataType::Any,
            _ => DataType::None,
        }
    }
}

impl From<DataType> for ffi::DataType {
    fn from(value: DataType) -> Self {
        match value {
            DataType::Uint8 => ffi::DataType::Uint8,
            DataType::Float => ffi::DataType::Float,
            DataType::Float16 => ffi::DataType::Float16,
            DataType::Any => ffi::DataType::Any,
            _ => ffi::DataType::None,
        }
    }
}

/// Distance metric for approximate nearest neighbor search.
///
/// This enum specifies the distance metric used to measure similarity between vectors.
/// Different metrics are appropriate for different types of data and use cases.
///
/// # Metrics
///
/// ## Euclidean and L-norms
/// * `L1` - Manhattan distance (sum of absolute differences)
/// * `L2` - Euclidean distance. Most common metric for continuous data.
/// * `NormalizedL2` - L2 distance normalized by vector magnitude
///
/// ## Angular distances
/// * `Angle` - Angular distance. Useful for directional similarity.
/// * `Cosine` - Cosine similarity distance. Good for high-dimensional data.
/// * `NormalizedAngle` - Normalized angular distance
/// * `NormalizedCosine` - Normalized cosine similarity distance
///
/// ## Hamming and Jaccard distances
/// * `Hamming` - Hamming distance for binary/categorical vectors.
/// * `Jaccard` - Jaccard distance for set similarity.
/// * `SparseJaccard` - Optimized Jaccard for sparse vectors.
///
/// ## Inner product
/// * `InnerProduct` - Inner product distance. Optimized for dot product similarity. Common aliases: `DotProduct`, `dp`.
///
/// ## Hyperbolic distances
/// * `Poincare` - Poincare distance for hyperbolic geometry
/// * `Lorentz` - Lorentz distance for Lorentz model
///
/// * `None` - Invalid or uninitialized state
#[derive(Serialize, Deserialize, Debug, Clone, Copy, PartialEq, Eq)]
pub enum DistanceType {
    /// Invalid or uninitialized distance metric.
    #[serde(rename = "None", alias = "none")]
    None,
    /// Manhattan distance (L1 norm).
    ///
    /// Calculated as the sum of absolute differences: $\sum |x_i - y_i|$
    /// Useful for sparse data and when different dimensions have different importance.
    #[serde(rename = "l1", alias = "L1")]
    L1,
    /// Euclidean distance (L2 norm).
    ///
    /// Calculated as: $\sqrt{\sum (x_i - y_i)^2}$
    /// The most commonly used distance metric. Suitable for most continuous data.
    #[serde(rename = "l2", alias = "L2")]
    L2,
    /// Hamming distance.
    ///
    /// Counts the number of positions where vector components differ.
    /// Useful for binary or categorical data encoded as bit vectors.
    #[serde(rename = "hamming", alias = "Hamming", alias = "ham")]
    Hamming,
    /// Angular distance.
    ///
    /// Measures the angle between vectors. Useful for directional similarity
    /// and when magnitude is irrelevant.
    #[serde(rename = "angle", alias = "Angle", alias = "ang")]
    Angle,
    /// Cosine similarity distance.
    ///
    /// Calculated as: $1 - \cos(\theta) = 1 - \frac{x \cdot y}{||x|| \cdot ||y||}$
    /// Excellent for high-dimensional sparse data, NLP applications, and when
    /// only direction matters, not magnitude.
    #[serde(rename = "cosine", alias = "Cosine", alias = "cos")]
    Cosine,
    /// Normalized angular distance.
    ///
    /// Angular distance normalized to a standard range.
    /// Useful when you need bounded values between 0 and 1.
    #[serde(
        rename = "normalizedangle",
        alias = "NormalizedAngle",
        alias = "normang",
        alias = "NormAng"
    )]
    NormalizedAngle,
    /// Normalized cosine similarity distance.
    ///
    /// Cosine similarity normalized to a standard range [0, 1].
    /// Provides the same properties as cosine distance but with normalized bounds.
    #[serde(
        rename = "normalizedcosine",
        alias = "NormalizedCosine",
        alias = "normcos",
        alias = "NormCos"
    )]
    NormalizedCosine,
    /// Jaccard distance for sets.
    ///
    /// Calculated as: $1 - \frac{|A \cap B|}{|A \cup B|}$
    /// Useful for set-based similarity, typical for categorical or presence/absence data.
    #[serde(rename = "jaccard", alias = "Jaccard", alias = "jac")]
    Jaccard,
    /// Jaccard distance optimized for sparse vectors.
    ///
    /// Optimized version of Jaccard distance for sparse vector representations.
    /// Better performance when vectors have many zero elements.
    #[serde(rename = "sparsejaccard", alias = "SparseJaccard", alias = "spjac")]
    SparseJaccard,
    /// L2 distance normalized by vector magnitude.
    ///
    /// Normalized version of L2 distance that accounts for vector length differences.
    /// Useful when you want Euclidean distance but normalized by magnitude.
    #[serde(rename = "normalizedl2", alias = "NormalizedL2", alias = "norml2")]
    NormalizedL2,
    /// Inner product distance.
    ///
    /// Calculated as: $x \cdot y = \sum x_i \cdot y_i$
    /// Note: Higher inner product = greater similarity (opposite of other metrics).
    /// Optimized for dot product similarity searches, common in recommendation systems.
    /// Aliases: `DotProduct`, `dp`
    #[serde(
        rename = "innerproduct",
        alias = "InnerProduct",
        alias = "ip",
        alias = "dotproduct",
        alias = "DotProduct",
        alias = "dp"
    )]
    InnerProduct,
    /// Poincaré distance for hyperbolic geometry.
    ///
    /// Distance metric in the Poincaré model of hyperbolic space.
    /// Useful for hierarchical data structures and tree-like relationships.
    #[serde(rename = "poincare", alias = "Poincare", alias = "poinc")]
    Poincare,
    /// Lorentz distance (Lorentz model of hyperbolic geometry).
    ///
    /// Alternative distance metric for hyperbolic space using the Lorentz model.
    /// Can be more efficient than Poincaré distance in some scenarios.
    #[serde(rename = "lorentz", alias = "Lorentz", alias = "loren")]
    Lorentz,
}

impl From<ffi::DistanceType> for DistanceType {
    fn from(value: ffi::DistanceType) -> Self {
        match value {
            ffi::DistanceType::L1 => DistanceType::L1,
            ffi::DistanceType::L2 => DistanceType::L2,
            ffi::DistanceType::Hamming => DistanceType::Hamming,
            ffi::DistanceType::Angle => DistanceType::Angle,
            ffi::DistanceType::Cosine => DistanceType::Cosine,
            ffi::DistanceType::NormalizedAngle => DistanceType::NormalizedAngle,
            ffi::DistanceType::NormalizedCosine => DistanceType::NormalizedCosine,
            ffi::DistanceType::Jaccard => DistanceType::Jaccard,
            ffi::DistanceType::SparseJaccard => DistanceType::SparseJaccard,
            ffi::DistanceType::NormalizedL2 => DistanceType::NormalizedL2,
            ffi::DistanceType::InnerProduct => DistanceType::InnerProduct,
            ffi::DistanceType::Poincare => DistanceType::Poincare,
            ffi::DistanceType::Lorentz => DistanceType::Lorentz,
            _ => DistanceType::None,
        }
    }
}

impl From<DistanceType> for ffi::DistanceType {
    fn from(value: DistanceType) -> Self {
        match value {
            DistanceType::L1 => ffi::DistanceType::L1,
            DistanceType::L2 => ffi::DistanceType::L2,
            DistanceType::Hamming => ffi::DistanceType::Hamming,
            DistanceType::Angle => ffi::DistanceType::Angle,
            DistanceType::Cosine => ffi::DistanceType::Cosine,
            DistanceType::NormalizedAngle => ffi::DistanceType::NormalizedAngle,
            DistanceType::NormalizedCosine => ffi::DistanceType::NormalizedCosine,
            DistanceType::Jaccard => ffi::DistanceType::Jaccard,
            DistanceType::SparseJaccard => ffi::DistanceType::SparseJaccard,
            DistanceType::NormalizedL2 => ffi::DistanceType::NormalizedL2,
            DistanceType::InnerProduct => ffi::DistanceType::InnerProduct,
            DistanceType::Poincare => ffi::DistanceType::Poincare,
            DistanceType::Lorentz => ffi::DistanceType::Lorentz,
            _ => ffi::DistanceType::None,
        }
    }
}

/// C++ Foreign Function Interface (FFI) bindings for QBG.
///
/// This module defines the low-level C++ FFI bindings using the `cxx` crate.
/// It provides direct mapping between Rust and C++ types and function calls.
///
/// # C++ Library Integration
///
/// The `ffi` module is generated from C++ code and provides:
/// - C++ type definitions (`Property`, `Index`) as opaque types
/// - C++ function wrappers (`new_index`, `new_prebuilt_index`, etc.)
/// - Enum mappings for data types and distance metrics
/// - Raw FFI calls that are wrapped by higher-level modules
///
/// # Safety
///
/// All items in this module should be considered `unsafe` to use directly.
/// Use the higher-level wrappers in `property` and `index` modules instead,
/// which provide safe abstractions and proper error handling.
///
/// # Memory Management
///
/// Objects like `Property` and `Index` are owned via `UniquePtr<T>`, which ensures
/// automatic deallocation when dropped, preventing memory leaks from C++ allocations.
#[cxx::bridge]
pub mod ffi {
    #[repr(i32)]
    enum ObjectType {
        Uint8 = 0,
        Float = 1,
        Float16 = 2,
        None = 99,
    }

    #[repr(i32)]
    enum DataType {
        Uint8 = 0,
        Float = 1,
        Float16 = 2,
        None = 99,
        Any = 100,
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
            internal_data_type: DataType,
            data_type: ObjectType,
            distance_type: DistanceType,
        );
        fn set_extended_dimension(self: Pin<&mut Property>, extended_dimension: usize);
        fn set_dimension(self: Pin<&mut Property>, dimension: usize);
        fn set_number_of_subvectors(self: Pin<&mut Property>, number_of_subvectors: usize);
        fn set_number_of_blobs(self: Pin<&mut Property>, number_of_blobs: usize);
        fn set_internal_data_type(self: Pin<&mut Property>, internal_data_type: DataType);
        fn set_data_type(self: Pin<&mut Property>, data_type: ObjectType);
        fn set_distance_type(self: Pin<&mut Property>, distance_type: DistanceType);
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

/// Configuration management for QBG index construction.
///
/// This module provides the `Property` struct, which wraps the C++ QBG property configuration.
/// It allows users to set construction parameters (dimension, clustering, quantization) and
/// build parameters (hierarchical clustering, optimization) before creating or modifying an index.
///
/// # C++ Binding
///
/// Property wraps `ffi::Property`, which is a UniquePtr to the underlying C++ property object.
/// All configuration is delegated directly to the C++ implementation for consistency.
///
/// # Usage Pattern
///
/// Properties must be configured before index creation:
/// 1. Create a Property instance via `Property::new()`
/// 2. Initialize construction parameters with `init_qbg_construction_parameters()`
/// 3. Set construction parameters with `set_qbg_construction_parameters()`
/// 4. Pass to `Index::new()` to create the index
pub mod property {
    use super::ffi;
    use cxx::UniquePtr;
    use std::pin::Pin;

    /// QBG index property configuration.
    ///
    /// `Property` encapsulates all configuration parameters needed to create or load a QBG index.
    /// It provides a type-safe interface to the underlying C++ property object, managing memory
    /// automatically through Rust's ownership system.
    ///
    /// # Usage
    ///
    /// Typically used in this pattern:
    /// 1. Create a new Property with `Property::new()`
    /// 2. Configure parameters using setter methods
    /// 3. Pass to `Index::new()` or `Index::open()` to create/open an index
    ///
    /// # Thread Safety
    ///
    /// A Property should not be shared across threads during configuration. Once created,
    /// pass it to an Index which manages thread safety.
    pub struct Property {
        /// The underlying C++ QBG Property object.
        ///
        /// Manages the lifetime and memory of the C++ property instance.
        /// Automatically cleaned up when Property is dropped.
        inner: UniquePtr<ffi::Property>,
    }

    impl Default for Property {
        fn default() -> Self {
            Property {
                inner: ffi::new_property()
            }
        }
    }

    impl Property {
        /// Creates a new Property instance with default C++ configuration.
        ///
        /// This initializes the underlying C++ property object which can be configured
        /// before using it to create or modify a QBG index.
        pub fn new() -> Self {
            Property {
                inner: ffi::new_property()
            }
        }

        /// Gets a mutable reference to the underlying C++ Property object.
        ///
        /// This is used internally when passing the property to C++ functions.
        /// Users should typically use the typed setter methods instead.
        pub fn get_property(&mut self) -> Pin<&mut ffi::Property> {
            self.inner.pin_mut()
        }

        /// Initializes QBG construction parameters to default values.
        ///
        /// Must be called before setting construction parameters.
        pub fn init_qbg_construction_parameters(&mut self) {
            self.inner.pin_mut().init_qbg_construction_parameters()
        }

        /// Sets all QBG construction parameters at once.
        ///
        /// # Arguments
        ///
        /// * `extended_dimension` - The extended vector dimension (usually equal to or greater than dimension)
        /// * `dimension` - The actual vector dimension
        /// * `number_of_subvectors` - Number of subvectors for quantization (typically 8-256)
        /// * `number_of_blobs` - Number of blobs in the graph. 0 means automatic.
        /// * `internal_data_type` - Data type for internal index storage (Float, Uint8, Float16)
        /// * `data_type` - Input vector data type (ObjectType)
        /// * `distance_type` - Distance metric to use for similarity measurement
        pub fn set_qbg_construction_parameters(
            &mut self,
            extended_dimension: usize,
            dimension: usize,
            number_of_subvectors: usize,
            number_of_blobs: usize,
            internal_data_type: ffi::DataType,
            data_type: ffi::ObjectType,
            distance_type: ffi::DistanceType,
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

        /// Sets the extended vector dimension.
        ///
        /// The extended dimension is used for preprocessing and can be larger than
        /// the actual data dimension.
        pub fn set_extended_dimension(&mut self, extended_dimension: usize) {
            self.inner
                .pin_mut()
                .set_extended_dimension(extended_dimension)
        }

        /// Sets the actual vector dimension.
        ///
        /// This should typically equal or be less than extended_dimension.
        pub fn set_dimension(&mut self, dimension: usize) {
            self.inner.pin_mut().set_dimension(dimension)
        }

        /// Sets the number of subvectors for quantization.
        ///
        /// Higher values increase precision but also increase memory and computation.
        /// Typical values: 8, 16, 32, 64, 128, 256
        pub fn set_number_of_subvectors(&mut self, number_of_subvectors: usize) {
            self.inner
                .pin_mut()
                .set_number_of_subvectors(number_of_subvectors)
        }

        /// Sets the number of blobs in the graph structure.
        ///
        /// A blob is a cluster of vectors. 0 means automatic calculation.
        pub fn set_number_of_blobs(&mut self, number_of_blobs: usize) {
            self.inner.pin_mut().set_number_of_blobs(number_of_blobs)
        }

        /// Sets the internal data type for index storage.
        ///
        /// This determines how vectors are quantized and stored internally.
        pub fn set_internal_data_type(&mut self, internal_data_type: ffi::DataType) {
            self.inner
                .pin_mut()
                .set_internal_data_type(internal_data_type)
        }

        /// Sets the input vector data type.
        ///
        /// This specifies the format of vectors provided to the index.
        pub fn set_data_type(&mut self, data_type: ffi::ObjectType) {
            self.inner.pin_mut().set_data_type(data_type)
        }

        /// Sets the distance metric for similarity measurement.
        pub fn set_distance_type(&mut self, distance_type: ffi::DistanceType) {
            self.inner.pin_mut().set_distance_type(distance_type)
        }

        /// Initializes QBG build parameters to default values.
        ///
        /// Must be called before setting build parameters.
        pub fn init_qbg_build_parameters(&mut self) {
            self.inner.pin_mut().init_qbg_build_parameters()
        }

        /// Sets all QBG build parameters at once.
        ///
        /// Build parameters control the index construction process including clustering
        /// hierarchy and rotation/optimization settings.
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

        /// Sets the initialization mode for hierarchical clustering.
        pub fn set_hierarchical_clustering_init_mode(
            &mut self,
            hierarchical_clustering_init_mode: i32,
        ) {
            self.inner
                .pin_mut()
                .set_hierarchical_clustering_init_mode(hierarchical_clustering_init_mode)
        }

        /// Sets the number of objects in the first clustering level.
        pub fn set_number_of_first_objects(&mut self, number_of_first_objects: usize) {
            self.inner
                .pin_mut()
                .set_number_of_first_objects(number_of_first_objects)
        }

        /// Sets the number of clusters in the first clustering level.
        pub fn set_number_of_first_clusters(&mut self, number_of_first_clusters: usize) {
            self.inner
                .pin_mut()
                .set_number_of_first_clusters(number_of_first_clusters)
        }

        /// Sets the number of objects in the second clustering level.
        pub fn set_number_of_second_objects(&mut self, number_of_second_objects: usize) {
            self.inner
                .pin_mut()
                .set_number_of_second_objects(number_of_second_objects)
        }

        /// Sets the number of clusters in the second clustering level.
        pub fn set_number_of_second_clusters(&mut self, number_of_second_clusters: usize) {
            self.inner
                .pin_mut()
                .set_number_of_second_clusters(number_of_second_clusters)
        }

        /// Sets the number of clusters in the third clustering level.
        pub fn set_number_of_third_clusters(&mut self, number_of_third_clusters: usize) {
            self.inner
                .pin_mut()
                .set_number_of_third_clusters(number_of_third_clusters)
        }

        /// Sets the total number of objects to consider in clustering.
        pub fn set_number_of_objects(&mut self, number_of_objects: usize) {
            self.inner
                .pin_mut()
                .set_number_of_objects(number_of_objects)
        }

        /// Sets the number of subvectors for build parameters.
        pub fn set_number_of_subvectors_for_bp(&mut self, number_of_subvectors: usize) {
            self.inner
                .pin_mut()
                .set_number_of_subvectors_for_bp(number_of_subvectors)
        }

        /// Sets the initialization mode for optimization clustering.
        pub fn set_optimization_clustering_init_mode(
            &mut self,
            optimization_clustering_init_mode: i32,
        ) {
            self.inner
                .pin_mut()
                .set_optimization_clustering_init_mode(optimization_clustering_init_mode)
        }

        /// Sets the number of iterations for rotation optimization.
        ///
        /// More iterations increase rotation quality but also increase build time.
        pub fn set_rotation_iteration(&mut self, rotation_iteration: usize) {
            self.inner
                .pin_mut()
                .set_rotation_iteration(rotation_iteration)
        }

        /// Sets the number of iterations for subvector optimization.
        pub fn set_subvector_iteration(&mut self, subvector_iteration: usize) {
            self.inner
                .pin_mut()
                .set_subvector_iteration(subvector_iteration)
        }

        /// Sets the number of rotation matrices.
        pub fn set_number_of_matrices(&mut self, number_of_matrices: usize) {
            self.inner
                .pin_mut()
                .set_number_of_matrices(number_of_matrices)
        }

        /// Enables or disables rotation during index construction.
        ///
        /// Rotation can improve search quality for certain data distributions.
        pub fn set_rotation(&mut self, rotation: bool) {
            self.inner.pin_mut().set_rotation(rotation)
        }

        /// Enables or disables repositioning during index construction.
        pub fn set_repositioning(&mut self, repositioning: bool) {
            self.inner.pin_mut().set_repositioning(repositioning)
        }
    }
}

/// QBG Index operations and search functionality.
///
/// This module provides the `Index` struct, which is the main interface for all QBG operations.
/// It wraps the C++ QBG index implementation and provides safe Rust abstractions for:
/// - Creating new indexes
/// - Loading prebuilt indexes from disk
/// - Inserting, updating, and removing vectors
/// - Searching for approximate nearest neighbors
/// - Saving and closing indexes
///
/// # C++ Binding
///
/// Index wraps `ffi::Index`, which is a UniquePtr to the underlying C++ index object.
/// All heavy lifting is performed by the C++ implementation, which uses optimized SIMD
/// instructions (AVX-512/AVX-2) for maximum performance.
///
/// # Memory Safety
///
/// The Index holds ownership of the C++ index object via UniquePtr, ensuring automatic
/// cleanup when the Index is dropped. This prevents memory leaks and dangling pointers.
///
/// # Thread Safety
///
/// Index implements Send and Sync, but users must ensure proper synchronization when
/// sharing index access across threads, as the C++ implementation may not be internally
/// thread-safe for concurrent modifications.
pub mod index {
    use super::ffi;
    use super::property;
    use core::slice;
    use cxx::UniquePtr;

    /// A QBG (Query-by-Graph) approximate nearest neighbor search index.
    ///
    /// `Index` is the core data structure for QBG-based vector search operations. It provides
    /// methods to create/load indexes, insert vectors, search for nearest neighbors, and optimize
    /// the index structure.
    ///
    /// # Creation and Loading
    ///
    /// - `Index::new()` - Create a new index from scratch with configuration from a Property
    /// - `Index::open()` - Load an existing index from disk
    ///
    /// # Operations
    ///
    /// - **Insert/Update**: `insert()` - Add or update vectors in the index
    /// - **Search**: `search()` - Find k nearest neighbors to a query vector
    /// - **Optimization**: `rebuild()` - Reconstruct and optimize the index structure
    /// - **Serialization**: `save()` - Persist index to disk
    ///
    /// # Thread Safety
    ///
    /// The underlying C++ QBG index supports concurrent read operations (searches) but
    /// write operations (insert, rebuild) may have synchronization overhead. The index
    /// should be accessed through proper synchronization primitives (Arc<Mutex<Index>>) in
    /// multi-threaded contexts.
    ///
    /// # Memory Management
    ///
    /// The Index automatically manages C++ memory through a UniquePtr. All vectors and
    /// indexes are cleaned up when the Index is dropped.
    pub struct Index {
        /// The underlying C++ QBG Index object.
        ///
        /// Manages the C++ index instance and its associated data structures.
        /// Automatically cleaned up when Index is dropped.
        inner: UniquePtr<ffi::Index>,
    }

    impl Index {
        /// Creates a new QBG index at the specified path.
        ///
        /// This constructs a new index from scratch using the provided property configuration.
        /// The index is built using the parameters specified in the Property object.
        ///
        /// # Arguments
        ///
        /// * `path` - File system path where the index will be stored
        /// * `p` - Property object containing index configuration
        ///
        /// # Returns
        ///
        /// A new Index instance or an error if index creation fails.
        pub fn new(path: &String, p: &mut property::Property) -> Result<Self, cxx::Exception> {
            let inner = ffi::new_index(path, p.get_property())?;
            Ok(Index { inner })
        }

        /// Opens a prebuilt index from disk.
        ///
        /// This loads an existing index that was previously saved. Use this when you have
        /// an index file already built and want to perform search operations.
        ///
        /// # Arguments
        ///
        /// * `path` - File system path to the existing index
        /// * `p` - Whether the index is prebuilt (typically true for loading existing indexes)
        ///
        /// # Returns
        ///
        /// An Index instance wrapping the loaded index, or an error if loading fails.
        pub fn new_prebuilt(path: &String, p: bool) -> Result<Self, cxx::Exception> {
            let inner = ffi::new_prebuilt_index(path, p)?;
            Ok(Index { inner })
        }

        /// Opens or reopens an index from disk.
        ///
        /// This allows switching which index file is being used by the current Index instance.
        ///
        /// # Arguments
        ///
        /// * `path` - File system path to the index
        /// * `prebuilt` - Whether the index should be treated as prebuilt
        pub fn open_index(&mut self, path: &String, prebuilt: bool) -> Result<(), cxx::Exception> {
            self.inner.pin_mut().open_index(path, prebuilt)
        }

        /// Rebuilds the index with new parameters.
        ///
        /// This is useful when you want to recreate or optimize an index with different
        /// clustering or optimization parameters.
        ///
        /// # Arguments
        ///
        /// * `path` - File system path for the rebuilt index
        /// * `p` - Property object with new construction parameters
        pub fn build_index(
            &mut self,
            path: &String,
            p: &mut property::Property,
        ) -> Result<(), cxx::Exception> {
            self.inner.pin_mut().build_index(path, p.get_property())
        }

        /// Saves the current index state to disk.
        ///
        /// This persists all vectors and internal structures to the index file.
        /// Should be called after performing insert/update/delete operations to ensure
        /// changes are not lost.
        pub fn save_index(&mut self) -> Result<(), cxx::Exception> {
            self.inner.pin_mut().save_index()
        }

        /// Closes the index and frees associated resources.
        ///
        /// After calling this, the Index should not be used for further operations.
        pub fn close_index(&mut self) {
            self.inner.pin_mut().close_index()
        }

        /// Appends a vector to the index and returns its assigned ID.
        ///
        /// This assigns a new sequential ID to the vector. Use this when you want
        /// the system to assign IDs automatically.
        ///
        /// # Arguments
        ///
        /// * `v` - Vector data with dimension matching the index configuration
        ///
        /// # Returns
        ///
        /// The auto-assigned object ID or an error if the operation fails.
        pub fn append(&mut self, v: &[f32]) -> Result<i32, cxx::Exception> {
            self.inner.pin_mut().append(v)
        }

        /// Inserts a vector into the index.
        ///
        /// Similar to append but may have different semantics depending on the C++ implementation.
        ///
        /// # Arguments
        ///
        /// * `v` - Vector data with dimension matching the index configuration
        ///
        /// # Returns
        ///
        /// The assigned object ID or an error if the operation fails.
        pub fn insert(&mut self, v: &[f32]) -> Result<i32, cxx::Exception> {
            self.inner.pin_mut().insert(v)
        }

        /// Removes a vector from the index by its object ID.
        ///
        /// This marks the vector as deleted and removes it from search results.
        ///
        /// # Arguments
        ///
        /// * `id` - Object ID of the vector to remove
        pub fn remove(&mut self, id: usize) -> Result<(), cxx::Exception> {
            self.inner.pin_mut().remove(id)
        }

        /// Searches for approximate nearest neighbors.
        ///
        /// Performs an ANN search and returns the k nearest neighbors within the search radius.
        ///
        /// # Arguments
        ///
        /// * `v` - Query vector with dimension matching the index configuration
        /// * `k` - Number of nearest neighbors to return
        /// * `radius` - Maximum search radius (0.0 means no radius limit)
        /// * `epsilon` - Search accuracy parameter (higher values = faster but less accurate)
        ///
        /// # Returns
        ///
        /// A vector of (object_id, distance) tuples for the found neighbors.
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

        /// Retrieves a vector from the index by its object ID.
        ///
        /// # Arguments
        ///
        /// * `id` - Object ID of the vector to retrieve
        ///
        /// # Returns
        ///
        /// A slice containing the vector data with dimension matching the index configuration.
        pub fn get_object(&self, id: usize) -> Result<&[f32], cxx::Exception> {
            let dim = self.inner.get_dimension()?;
            match self.inner.get_object(id) {
                Ok(v) => Ok(unsafe { slice::from_raw_parts(v, dim) }),
                Err(e) => Err(e),
            }
        }

        /// Returns the vector dimension configured for this index.
        ///
        /// # Returns
        ///
        /// The dimension size (number of elements per vector) or an error if the query fails.
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
        p.pin_mut().set_internal_data_type(ffi::DataType::Float);
        p.pin_mut().set_data_type(ffi::ObjectType::Float);
        p.pin_mut().set_distance_type(ffi::DistanceType::L2);
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
        let mut search_results = index.pin_mut().search(vec.as_slice(), K, RADIUS, EPSILON)?;
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
        let mut search_results = index.pin_mut().search(vec.as_slice(), K, RADIUS, EPSILON)?;
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
        let mut search_results = index.pin_mut().search(vec.as_slice(), K, RADIUS, EPSILON)?;
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
        let mut search_results = index.pin_mut().search(vec.as_slice(), K, RADIUS, EPSILON)?;
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
        p.set_qbg_construction_parameters(
            1,
            1,
            1,
            1,
            ffi::DataType::Float,
            ffi::ObjectType::Float,
            ffi::DistanceType::L2,
        );
        p.set_extended_dimension(1);
        p.set_dimension(1);
        p.set_number_of_subvectors(1);
        p.set_number_of_blobs(1);
        p.set_internal_data_type(ffi::DataType::Float);
        p.set_data_type(ffi::ObjectType::Float);
        p.set_distance_type(ffi::DistanceType::L2);
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
            let res = index.append(vec.as_slice());
            assert!(res.is_ok(), "append failed: {:?}", res.err());
            assert_eq!((i + 1) as i32, res.unwrap())
        }

        // Build
        println!("building the index...");
        let res = index.build_index(&path, &mut p);
        assert!(res.is_ok(), "build_index failed: {:?}", res.err());
        let res = index.open_index(&path, true);
        assert!(res.is_ok(), "open_index failed: {:?}", res.err());

        // Insert
        for i in 0..100 {
            let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|x| (x + i) as f32).collect();
            let res = index.insert(vec.as_slice());
            assert!(res.is_ok(), "insert failed: {:?}", res.err());
            assert_eq!((i + 1 + 100) as i32, res.unwrap());
        }

        // Get Object
        let res = index.get_object(1);
        assert!(res.is_ok(), "get_object failed: {:?}", res.err());

        // Get Dimension
        let res = index.get_dimension();
        assert!(res.is_ok(), "get_dimension failed: {:?}", res.err());
        assert!(res.unwrap() > 0, "dimension should be greater than 0");

        // Search
        println!("search the index for the specified query...");
        let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|i| i as f32).collect();
        let res = index.search(vec.as_slice(), K, RADIUS, EPSILON);
        assert!(res.is_ok(), "search failed: {:?}", res.err());
        let search_results = res.unwrap();
        let ids: Vec<u32> = search_results.iter().map(|s| s.0).collect();
        let distances: Vec<f32> = search_results.iter().map(|s| s.1).collect();
        println!("search results:\n\t{:?}", search_results);
        println!("ids:\n\t{:?}", ids);
        println!("distances:\n\t{:?}", distances);

        // Remove
        let res = index.remove(1);
        assert!(res.is_ok(), "remove failed: {:?}", res.err());
        let vec: Vec<f32> = (0..DIMENSION).into_iter().map(|i| i as f32).collect();
        let res = index.search(vec.as_slice(), K, RADIUS, EPSILON);
        assert!(res.is_ok(), "search failed: {:?}", res.err());
        let search_results = res.unwrap();
        assert!(
            !search_results.is_empty(),
            "search results should not be empty"
        );

        index.close_index();

        Ok(())
    }
}
