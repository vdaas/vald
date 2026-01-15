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

use bincode::{Decode, Encode};
use serde::{Serialize, de::DeserializeOwned};
use std::fmt::Debug;
use std::hash::Hash;

/// A trait that defines the requirements for a key in the key-value store.
///
/// This trait is a marker trait that bundles all the necessary traits for a type to be used as a key.
/// This includes serialization, deserialization, encoding, decoding, equality, hashing, cloning, and thread safety.
pub trait KeyType:
    Serialize
    + DeserializeOwned
    + Encode
    + Decode<()>
    + Eq
    + Hash
    + Clone
    + Send
    + Sync
    + Debug
    + 'static
{
}
impl<
    T: Serialize
        + DeserializeOwned
        + Encode
        + Decode<()>
        + Eq
        + Hash
        + Clone
        + Send
        + Sync
        + Debug
        + 'static,
> KeyType for T
{
}

/// A trait that defines the requirements for a value in the key-value store.
///
/// This trait is a marker trait that bundles all the necessary traits for a type to be used as a value.
/// This includes serialization, deserialization, encoding, decoding, equality, hashing, cloning, and thread safety.
pub trait ValueType:
    Serialize
    + DeserializeOwned
    + Encode
    + Decode<()>
    + Eq
    + Hash
    + Clone
    + Send
    + Sync
    + Debug
    + 'static
{
}
impl<
    T: Serialize
        + DeserializeOwned
        + Encode
        + Decode<()>
        + Eq
        + Hash
        + Clone
        + Send
        + Sync
        + Debug
        + 'static,
> ValueType for T
{
}
