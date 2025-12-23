//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
impl serde::Serialize for BatchScanRegionsRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.need_buckets {
            len += 1;
        }
        if !self.ranges.is_empty() {
            len += 1;
        }
        if self.limit != 0 {
            len += 1;
        }
        if self.contain_all_key_range {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pdpb.BatchScanRegionsRequest", len)?;
        if self.need_buckets {
            struct_ser.serialize_field("needBuckets", &self.need_buckets)?;
        }
        if !self.ranges.is_empty() {
            struct_ser.serialize_field("ranges", &self.ranges)?;
        }
        if self.limit != 0 {
            struct_ser.serialize_field("limit", &self.limit)?;
        }
        if self.contain_all_key_range {
            struct_ser.serialize_field("containAllKeyRange", &self.contain_all_key_range)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for BatchScanRegionsRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "need_buckets",
            "needBuckets",
            "ranges",
            "limit",
            "contain_all_key_range",
            "containAllKeyRange",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            NeedBuckets,
            Ranges,
            Limit,
            ContainAllKeyRange,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "needBuckets" | "need_buckets" => Ok(GeneratedField::NeedBuckets),
                            "ranges" => Ok(GeneratedField::Ranges),
                            "limit" => Ok(GeneratedField::Limit),
                            "containAllKeyRange" | "contain_all_key_range" => Ok(GeneratedField::ContainAllKeyRange),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = BatchScanRegionsRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pdpb.BatchScanRegionsRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<BatchScanRegionsRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut need_buckets__ = None;
                let mut ranges__ = None;
                let mut limit__ = None;
                let mut contain_all_key_range__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::NeedBuckets => {
                            if need_buckets__.is_some() {
                                return Err(serde::de::Error::duplicate_field("needBuckets"));
                            }
                            need_buckets__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ranges => {
                            if ranges__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ranges"));
                            }
                            ranges__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Limit => {
                            if limit__.is_some() {
                                return Err(serde::de::Error::duplicate_field("limit"));
                            }
                            limit__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ContainAllKeyRange => {
                            if contain_all_key_range__.is_some() {
                                return Err(serde::de::Error::duplicate_field("containAllKeyRange"));
                            }
                            contain_all_key_range__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(BatchScanRegionsRequest {
                    need_buckets: need_buckets__.unwrap_or_default(),
                    ranges: ranges__.unwrap_or_default(),
                    limit: limit__.unwrap_or_default(),
                    contain_all_key_range: contain_all_key_range__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pdpb.BatchScanRegionsRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for BatchScanRegionsResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.regions.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pdpb.BatchScanRegionsResponse", len)?;
        if !self.regions.is_empty() {
            struct_ser.serialize_field("regions", &self.regions)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for BatchScanRegionsResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "regions",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Regions,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regions" => Ok(GeneratedField::Regions),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = BatchScanRegionsResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pdpb.BatchScanRegionsResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<BatchScanRegionsResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut regions__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Regions => {
                            if regions__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regions"));
                            }
                            regions__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(BatchScanRegionsResponse {
                    regions: regions__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pdpb.BatchScanRegionsResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Error {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.r#type != 0 {
            len += 1;
        }
        if !self.message.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pdpb.Error", len)?;
        if self.r#type != 0 {
            let v = ErrorType::try_from(self.r#type)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.r#type)))?;
            struct_ser.serialize_field("type", &v)?;
        }
        if !self.message.is_empty() {
            struct_ser.serialize_field("message", &self.message)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Error {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "type",
            "message",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Type,
            Message,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "type" => Ok(GeneratedField::Type),
                            "message" => Ok(GeneratedField::Message),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Error;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pdpb.Error")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Error, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut r#type__ = None;
                let mut message__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Type => {
                            if r#type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("type"));
                            }
                            r#type__ = Some(map_.next_value::<ErrorType>()? as i32);
                        }
                        GeneratedField::Message => {
                            if message__.is_some() {
                                return Err(serde::de::Error::duplicate_field("message"));
                            }
                            message__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(Error {
                    r#type: r#type__.unwrap_or_default(),
                    message: message__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pdpb.Error", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ErrorType {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Ok => "OK",
            Self::Unknown => "UNKNOWN",
            Self::NotBootstrapped => "NOT_BOOTSTRAPPED",
            Self::StoreTombstone => "STORE_TOMBSTONE",
            Self::AlreadyBootstrapped => "ALREADY_BOOTSTRAPPED",
            Self::IncompatibleVersion => "INCOMPATIBLE_VERSION",
            Self::RegionNotFound => "REGION_NOT_FOUND",
            Self::GlobalConfigNotFound => "GLOBAL_CONFIG_NOT_FOUND",
            Self::DuplicatedEntry => "DUPLICATED_ENTRY",
            Self::EntryNotFound => "ENTRY_NOT_FOUND",
            Self::InvalidValue => "INVALID_VALUE",
            Self::DataCompacted => "DATA_COMPACTED",
            Self::RegionsNotContainAllKeyRange => "REGIONS_NOT_CONTAIN_ALL_KEY_RANGE",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for ErrorType {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "OK",
            "UNKNOWN",
            "NOT_BOOTSTRAPPED",
            "STORE_TOMBSTONE",
            "ALREADY_BOOTSTRAPPED",
            "INCOMPATIBLE_VERSION",
            "REGION_NOT_FOUND",
            "GLOBAL_CONFIG_NOT_FOUND",
            "DUPLICATED_ENTRY",
            "ENTRY_NOT_FOUND",
            "INVALID_VALUE",
            "DATA_COMPACTED",
            "REGIONS_NOT_CONTAIN_ALL_KEY_RANGE",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ErrorType;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                write!(formatter, "expected one of: {:?}", &FIELDS)
            }

            fn visit_i64<E>(self, v: i64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Signed(v), &self)
                    })
            }

            fn visit_u64<E>(self, v: u64) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                i32::try_from(v)
                    .ok()
                    .and_then(|x| x.try_into().ok())
                    .ok_or_else(|| {
                        serde::de::Error::invalid_value(serde::de::Unexpected::Unsigned(v), &self)
                    })
            }

            fn visit_str<E>(self, value: &str) -> std::result::Result<Self::Value, E>
            where
                E: serde::de::Error,
            {
                match value {
                    "OK" => Ok(ErrorType::Ok),
                    "UNKNOWN" => Ok(ErrorType::Unknown),
                    "NOT_BOOTSTRAPPED" => Ok(ErrorType::NotBootstrapped),
                    "STORE_TOMBSTONE" => Ok(ErrorType::StoreTombstone),
                    "ALREADY_BOOTSTRAPPED" => Ok(ErrorType::AlreadyBootstrapped),
                    "INCOMPATIBLE_VERSION" => Ok(ErrorType::IncompatibleVersion),
                    "REGION_NOT_FOUND" => Ok(ErrorType::RegionNotFound),
                    "GLOBAL_CONFIG_NOT_FOUND" => Ok(ErrorType::GlobalConfigNotFound),
                    "DUPLICATED_ENTRY" => Ok(ErrorType::DuplicatedEntry),
                    "ENTRY_NOT_FOUND" => Ok(ErrorType::EntryNotFound),
                    "INVALID_VALUE" => Ok(ErrorType::InvalidValue),
                    "DATA_COMPACTED" => Ok(ErrorType::DataCompacted),
                    "REGIONS_NOT_CONTAIN_ALL_KEY_RANGE" => Ok(ErrorType::RegionsNotContainAllKeyRange),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for GetAllStoresRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.exclude_tombstone_stores {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pdpb.GetAllStoresRequest", len)?;
        if self.exclude_tombstone_stores {
            struct_ser.serialize_field("excludeTombstoneStores", &self.exclude_tombstone_stores)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetAllStoresRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "exclude_tombstone_stores",
            "excludeTombstoneStores",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            ExcludeTombstoneStores,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "excludeTombstoneStores" | "exclude_tombstone_stores" => Ok(GeneratedField::ExcludeTombstoneStores),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetAllStoresRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pdpb.GetAllStoresRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetAllStoresRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut exclude_tombstone_stores__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::ExcludeTombstoneStores => {
                            if exclude_tombstone_stores__.is_some() {
                                return Err(serde::de::Error::duplicate_field("excludeTombstoneStores"));
                            }
                            exclude_tombstone_stores__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetAllStoresRequest {
                    exclude_tombstone_stores: exclude_tombstone_stores__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pdpb.GetAllStoresRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for GetAllStoresResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.stores.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pdpb.GetAllStoresResponse", len)?;
        if !self.stores.is_empty() {
            struct_ser.serialize_field("stores", &self.stores)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for GetAllStoresResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "stores",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Stores,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "stores" => Ok(GeneratedField::Stores),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = GetAllStoresResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pdpb.GetAllStoresResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<GetAllStoresResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut stores__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Stores => {
                            if stores__.is_some() {
                                return Err(serde::de::Error::duplicate_field("stores"));
                            }
                            stores__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(GetAllStoresResponse {
                    stores: stores__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pdpb.GetAllStoresResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for KeyRange {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.start_key.is_empty() {
            len += 1;
        }
        if !self.end_key.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pdpb.KeyRange", len)?;
        if !self.start_key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("startKey", pbjson::private::base64::encode(&self.start_key).as_str())?;
        }
        if !self.end_key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("endKey", pbjson::private::base64::encode(&self.end_key).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for KeyRange {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "start_key",
            "startKey",
            "end_key",
            "endKey",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            StartKey,
            EndKey,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "startKey" | "start_key" => Ok(GeneratedField::StartKey),
                            "endKey" | "end_key" => Ok(GeneratedField::EndKey),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = KeyRange;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pdpb.KeyRange")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<KeyRange, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut start_key__ = None;
                let mut end_key__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::StartKey => {
                            if start_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startKey"));
                            }
                            start_key__ = 
                                Some(map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::EndKey => {
                            if end_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("endKey"));
                            }
                            end_key__ = 
                                Some(map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(KeyRange {
                    start_key: start_key__.unwrap_or_default(),
                    end_key: end_key__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pdpb.KeyRange", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Region {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region.is_some() {
            len += 1;
        }
        if self.leader.is_some() {
            len += 1;
        }
        if !self.pending_peers.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("pdpb.Region", len)?;
        if let Some(v) = self.region.as_ref() {
            struct_ser.serialize_field("region", v)?;
        }
        if let Some(v) = self.leader.as_ref() {
            struct_ser.serialize_field("leader", v)?;
        }
        if !self.pending_peers.is_empty() {
            struct_ser.serialize_field("pendingPeers", &self.pending_peers)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Region {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "region",
            "leader",
            "pending_peers",
            "pendingPeers",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Region,
            Leader,
            PendingPeers,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "region" => Ok(GeneratedField::Region),
                            "leader" => Ok(GeneratedField::Leader),
                            "pendingPeers" | "pending_peers" => Ok(GeneratedField::PendingPeers),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Region;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct pdpb.Region")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Region, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut region__ = None;
                let mut leader__ = None;
                let mut pending_peers__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Region => {
                            if region__.is_some() {
                                return Err(serde::de::Error::duplicate_field("region"));
                            }
                            region__ = map_.next_value()?;
                        }
                        GeneratedField::Leader => {
                            if leader__.is_some() {
                                return Err(serde::de::Error::duplicate_field("leader"));
                            }
                            leader__ = map_.next_value()?;
                        }
                        GeneratedField::PendingPeers => {
                            if pending_peers__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pendingPeers"));
                            }
                            pending_peers__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(Region {
                    region: region__,
                    leader: leader__,
                    pending_peers: pending_peers__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("pdpb.Region", FIELDS, GeneratedVisitor)
    }
}
