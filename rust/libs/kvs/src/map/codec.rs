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

use crate::map::error::Error;

/// A trait for defining custom serialization and deserialization logic.
///
/// This allows `Map` to be generic over the data format, enabling users to
/// plug in their preferred serialization framework (e.g., Wincode, JSON, Protobuf).
pub trait Codec: Send + Sync + 'static {
    /// Serializes a given value into a byte vector.
    fn encode<T: serde::Serialize + wincode::SchemaWrite<Src = T> + ?Sized>(
        &self,
        v: &T,
    ) -> Result<Vec<u8>, Error>;
    /// Deserializes a byte slice into a value of a specific type.
    fn decode<T: serde::de::DeserializeOwned + for<'de> wincode::SchemaRead<'de, Dst = T>>(
        &self,
        bytes: &[u8],
    ) -> Result<T, Error>;
}

/// The default codec implementation using `bincode`.
#[derive(Clone, Debug, Default)]
pub struct BincodeCodec;

impl Codec for BincodeCodec {
    fn encode<T: serde::Serialize + wincode::SchemaWrite<Src = T> + ?Sized>(
        &self,
        v: &T,
    ) -> Result<Vec<u8>, Error> {
        wincode::serialize(v).map_err(|e| Error::Codec {
            source: Box::new(e),
        })
    }

    fn decode<T: serde::de::DeserializeOwned + for<'de> wincode::SchemaRead<'de, Dst = T>>(
        &self,
        bytes: &[u8],
    ) -> Result<T, Error> {
        wincode::deserialize(bytes)
            .map(|decoded| decoded)
            .map_err(|e| Error::Codec {
                source: Box::new(e),
            })
    }
}
