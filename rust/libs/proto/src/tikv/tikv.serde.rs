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
impl serde::Serialize for ApiVersion {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::V1 => "V1",
            Self::V1ttl => "V1TTL",
            Self::V2 => "V2",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for ApiVersion {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["V1", "V1TTL", "V2"];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ApiVersion;

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
                    "V1" => Ok(ApiVersion::V1),
                    "V1TTL" => Ok(ApiVersion::V1ttl),
                    "V2" => Ok(ApiVersion::V2),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for AlreadyExist {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.key.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.AlreadyExist", len)?;
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for AlreadyExist {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["key"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Key,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "key" => Ok(GeneratedField::Key),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = AlreadyExist;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.AlreadyExist")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<AlreadyExist, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut key__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(AlreadyExist {
                    key: key__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.AlreadyExist", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Assertion {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::None => "None",
            Self::Exist => "Exist",
            Self::NotExist => "NotExist",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for Assertion {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["None", "Exist", "NotExist"];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Assertion;

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
                    "None" => Ok(Assertion::None),
                    "Exist" => Ok(Assertion::Exist),
                    "NotExist" => Ok(Assertion::NotExist),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for AssertionFailed {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.start_ts != 0 {
            len += 1;
        }
        if !self.key.is_empty() {
            len += 1;
        }
        if self.assertion != 0 {
            len += 1;
        }
        if self.existing_start_ts != 0 {
            len += 1;
        }
        if self.existing_commit_ts != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.AssertionFailed", len)?;
        if self.start_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("startTs", ToString::to_string(&self.start_ts).as_str())?;
        }
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        if self.assertion != 0 {
            let v = Assertion::try_from(self.assertion).map_err(|_| {
                serde::ser::Error::custom(format!("Invalid variant {}", self.assertion))
            })?;
            struct_ser.serialize_field("assertion", &v)?;
        }
        if self.existing_start_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "existingStartTs",
                ToString::to_string(&self.existing_start_ts).as_str(),
            )?;
        }
        if self.existing_commit_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "existingCommitTs",
                ToString::to_string(&self.existing_commit_ts).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for AssertionFailed {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "start_ts",
            "startTs",
            "key",
            "assertion",
            "existing_start_ts",
            "existingStartTs",
            "existing_commit_ts",
            "existingCommitTs",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            StartTs,
            Key,
            Assertion,
            ExistingStartTs,
            ExistingCommitTs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "startTs" | "start_ts" => Ok(GeneratedField::StartTs),
                            "key" => Ok(GeneratedField::Key),
                            "assertion" => Ok(GeneratedField::Assertion),
                            "existingStartTs" | "existing_start_ts" => {
                                Ok(GeneratedField::ExistingStartTs)
                            }
                            "existingCommitTs" | "existing_commit_ts" => {
                                Ok(GeneratedField::ExistingCommitTs)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = AssertionFailed;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.AssertionFailed")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<AssertionFailed, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut start_ts__ = None;
                let mut key__ = None;
                let mut assertion__ = None;
                let mut existing_start_ts__ = None;
                let mut existing_commit_ts__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::StartTs => {
                            if start_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startTs"));
                            }
                            start_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Assertion => {
                            if assertion__.is_some() {
                                return Err(serde::de::Error::duplicate_field("assertion"));
                            }
                            assertion__ = Some(map_.next_value::<Assertion>()? as i32);
                        }
                        GeneratedField::ExistingStartTs => {
                            if existing_start_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("existingStartTs"));
                            }
                            existing_start_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ExistingCommitTs => {
                            if existing_commit_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("existingCommitTs"));
                            }
                            existing_commit_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(AssertionFailed {
                    start_ts: start_ts__.unwrap_or_default(),
                    key: key__.unwrap_or_default(),
                    assertion: assertion__.unwrap_or_default(),
                    existing_start_ts: existing_start_ts__.unwrap_or_default(),
                    existing_commit_ts: existing_commit_ts__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.AssertionFailed", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for BucketVersionNotMatch {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.version != 0 {
            len += 1;
        }
        if !self.keys.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.BucketVersionNotMatch", len)?;
        if self.version != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("version", ToString::to_string(&self.version).as_str())?;
        }
        if !self.keys.is_empty() {
            struct_ser.serialize_field(
                "keys",
                &self
                    .keys
                    .iter()
                    .map(pbjson::private::base64::encode)
                    .collect::<Vec<_>>(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for BucketVersionNotMatch {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["version", "keys"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Version,
            Keys,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "version" => Ok(GeneratedField::Version),
                            "keys" => Ok(GeneratedField::Keys),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = BucketVersionNotMatch;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.BucketVersionNotMatch")
            }

            fn visit_map<V>(
                self,
                mut map_: V,
            ) -> std::result::Result<BucketVersionNotMatch, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut version__ = None;
                let mut keys__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Version => {
                            if version__.is_some() {
                                return Err(serde::de::Error::duplicate_field("version"));
                            }
                            version__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Keys => {
                            if keys__.is_some() {
                                return Err(serde::de::Error::duplicate_field("keys"));
                            }
                            keys__ = Some(
                                map_.next_value::<Vec<::pbjson::private::BytesDeserialize<_>>>()?
                                    .into_iter()
                                    .map(|x| x.0)
                                    .collect(),
                            );
                        }
                    }
                }
                Ok(BucketVersionNotMatch {
                    version: version__.unwrap_or_default(),
                    keys: keys__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.BucketVersionNotMatch", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CommandPri {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Normal => "Normal",
            Self::Low => "Low",
            Self::High => "High",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for CommandPri {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["Normal", "Low", "High"];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CommandPri;

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
                    "Normal" => Ok(CommandPri::Normal),
                    "Low" => Ok(CommandPri::Low),
                    "High" => Ok(CommandPri::High),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for CommitTsExpired {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.start_ts != 0 {
            len += 1;
        }
        if self.attempted_commit_ts != 0 {
            len += 1;
        }
        if !self.key.is_empty() {
            len += 1;
        }
        if self.min_commit_ts != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.CommitTsExpired", len)?;
        if self.start_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("startTs", ToString::to_string(&self.start_ts).as_str())?;
        }
        if self.attempted_commit_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "attemptedCommitTs",
                ToString::to_string(&self.attempted_commit_ts).as_str(),
            )?;
        }
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        if self.min_commit_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "minCommitTs",
                ToString::to_string(&self.min_commit_ts).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CommitTsExpired {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "start_ts",
            "startTs",
            "attempted_commit_ts",
            "attemptedCommitTs",
            "key",
            "min_commit_ts",
            "minCommitTs",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            StartTs,
            AttemptedCommitTs,
            Key,
            MinCommitTs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "startTs" | "start_ts" => Ok(GeneratedField::StartTs),
                            "attemptedCommitTs" | "attempted_commit_ts" => {
                                Ok(GeneratedField::AttemptedCommitTs)
                            }
                            "key" => Ok(GeneratedField::Key),
                            "minCommitTs" | "min_commit_ts" => Ok(GeneratedField::MinCommitTs),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CommitTsExpired;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.CommitTsExpired")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CommitTsExpired, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut start_ts__ = None;
                let mut attempted_commit_ts__ = None;
                let mut key__ = None;
                let mut min_commit_ts__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::StartTs => {
                            if start_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startTs"));
                            }
                            start_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::AttemptedCommitTs => {
                            if attempted_commit_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("attemptedCommitTs"));
                            }
                            attempted_commit_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::MinCommitTs => {
                            if min_commit_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("minCommitTs"));
                            }
                            min_commit_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(CommitTsExpired {
                    start_ts: start_ts__.unwrap_or_default(),
                    attempted_commit_ts: attempted_commit_ts__.unwrap_or_default(),
                    key: key__.unwrap_or_default(),
                    min_commit_ts: min_commit_ts__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.CommitTsExpired", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for CommitTsTooLarge {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.commit_ts != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.CommitTsTooLarge", len)?;
        if self.commit_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("commitTs", ToString::to_string(&self.commit_ts).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for CommitTsTooLarge {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["commit_ts", "commitTs"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            CommitTs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "commitTs" | "commit_ts" => Ok(GeneratedField::CommitTs),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = CommitTsTooLarge;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.CommitTsTooLarge")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<CommitTsTooLarge, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut commit_ts__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::CommitTs => {
                            if commit_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("commitTs"));
                            }
                            commit_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(CommitTsTooLarge {
                    commit_ts: commit_ts__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.CommitTsTooLarge", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Context {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_id != 0 {
            len += 1;
        }
        if self.region_epoch.is_some() {
            len += 1;
        }
        if self.peer.is_some() {
            len += 1;
        }
        if self.term != 0 {
            len += 1;
        }
        if self.priority != 0 {
            len += 1;
        }
        if self.isolation_level != 0 {
            len += 1;
        }
        if self.not_fill_cache {
            len += 1;
        }
        if self.sync_log {
            len += 1;
        }
        if self.record_time_stat {
            len += 1;
        }
        if self.record_scan_stat {
            len += 1;
        }
        if self.replica_read {
            len += 1;
        }
        if !self.resolved_locks.is_empty() {
            len += 1;
        }
        if self.max_execution_duration_ms != 0 {
            len += 1;
        }
        if self.applied_index != 0 {
            len += 1;
        }
        if self.task_id != 0 {
            len += 1;
        }
        if self.stale_read {
            len += 1;
        }
        if !self.resource_group_tag.is_empty() {
            len += 1;
        }
        if self.disk_full_opt != 0 {
            len += 1;
        }
        if self.is_retry_request {
            len += 1;
        }
        if self.api_version != 0 {
            len += 1;
        }
        if !self.committed_locks.is_empty() {
            len += 1;
        }
        if !self.request_source.is_empty() {
            len += 1;
        }
        if self.txn_source != 0 {
            len += 1;
        }
        if self.busy_threshold_ms != 0 {
            len += 1;
        }
        if self.resource_control_context.is_some() {
            len += 1;
        }
        if !self.keyspace_name.is_empty() {
            len += 1;
        }
        if self.keyspace_id != 0 {
            len += 1;
        }
        if self.buckets_version != 0 {
            len += 1;
        }
        if self.source_stmt.is_some() {
            len += 1;
        }
        if self.cluster_id != 0 {
            len += 1;
        }
        if !self.trace_id.is_empty() {
            len += 1;
        }
        if self.trace_control_flags != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.Context", len)?;
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        if let Some(v) = self.region_epoch.as_ref() {
            struct_ser.serialize_field("regionEpoch", v)?;
        }
        if let Some(v) = self.peer.as_ref() {
            struct_ser.serialize_field("peer", v)?;
        }
        if self.term != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("term", ToString::to_string(&self.term).as_str())?;
        }
        if self.priority != 0 {
            let v = CommandPri::try_from(self.priority).map_err(|_| {
                serde::ser::Error::custom(format!("Invalid variant {}", self.priority))
            })?;
            struct_ser.serialize_field("priority", &v)?;
        }
        if self.isolation_level != 0 {
            let v = IsolationLevel::try_from(self.isolation_level).map_err(|_| {
                serde::ser::Error::custom(format!("Invalid variant {}", self.isolation_level))
            })?;
            struct_ser.serialize_field("isolationLevel", &v)?;
        }
        if self.not_fill_cache {
            struct_ser.serialize_field("notFillCache", &self.not_fill_cache)?;
        }
        if self.sync_log {
            struct_ser.serialize_field("syncLog", &self.sync_log)?;
        }
        if self.record_time_stat {
            struct_ser.serialize_field("recordTimeStat", &self.record_time_stat)?;
        }
        if self.record_scan_stat {
            struct_ser.serialize_field("recordScanStat", &self.record_scan_stat)?;
        }
        if self.replica_read {
            struct_ser.serialize_field("replicaRead", &self.replica_read)?;
        }
        if !self.resolved_locks.is_empty() {
            struct_ser.serialize_field(
                "resolvedLocks",
                &self
                    .resolved_locks
                    .iter()
                    .map(ToString::to_string)
                    .collect::<Vec<_>>(),
            )?;
        }
        if self.max_execution_duration_ms != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "maxExecutionDurationMs",
                ToString::to_string(&self.max_execution_duration_ms).as_str(),
            )?;
        }
        if self.applied_index != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "appliedIndex",
                ToString::to_string(&self.applied_index).as_str(),
            )?;
        }
        if self.task_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("taskId", ToString::to_string(&self.task_id).as_str())?;
        }
        if self.stale_read {
            struct_ser.serialize_field("staleRead", &self.stale_read)?;
        }
        if !self.resource_group_tag.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "resourceGroupTag",
                pbjson::private::base64::encode(&self.resource_group_tag).as_str(),
            )?;
        }
        if self.disk_full_opt != 0 {
            let v = DiskFullOpt::try_from(self.disk_full_opt).map_err(|_| {
                serde::ser::Error::custom(format!("Invalid variant {}", self.disk_full_opt))
            })?;
            struct_ser.serialize_field("diskFullOpt", &v)?;
        }
        if self.is_retry_request {
            struct_ser.serialize_field("isRetryRequest", &self.is_retry_request)?;
        }
        if self.api_version != 0 {
            let v = ApiVersion::try_from(self.api_version).map_err(|_| {
                serde::ser::Error::custom(format!("Invalid variant {}", self.api_version))
            })?;
            struct_ser.serialize_field("apiVersion", &v)?;
        }
        if !self.committed_locks.is_empty() {
            struct_ser.serialize_field(
                "committedLocks",
                &self
                    .committed_locks
                    .iter()
                    .map(ToString::to_string)
                    .collect::<Vec<_>>(),
            )?;
        }
        if !self.request_source.is_empty() {
            struct_ser.serialize_field("requestSource", &self.request_source)?;
        }
        if self.txn_source != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("txnSource", ToString::to_string(&self.txn_source).as_str())?;
        }
        if self.busy_threshold_ms != 0 {
            struct_ser.serialize_field("busyThresholdMs", &self.busy_threshold_ms)?;
        }
        if let Some(v) = self.resource_control_context.as_ref() {
            struct_ser.serialize_field("resourceControlContext", v)?;
        }
        if !self.keyspace_name.is_empty() {
            struct_ser.serialize_field("keyspaceName", &self.keyspace_name)?;
        }
        if self.keyspace_id != 0 {
            struct_ser.serialize_field("keyspaceId", &self.keyspace_id)?;
        }
        if self.buckets_version != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "bucketsVersion",
                ToString::to_string(&self.buckets_version).as_str(),
            )?;
        }
        if let Some(v) = self.source_stmt.as_ref() {
            struct_ser.serialize_field("sourceStmt", v)?;
        }
        if self.cluster_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("clusterId", ToString::to_string(&self.cluster_id).as_str())?;
        }
        if !self.trace_id.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "traceId",
                pbjson::private::base64::encode(&self.trace_id).as_str(),
            )?;
        }
        if self.trace_control_flags != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "traceControlFlags",
                ToString::to_string(&self.trace_control_flags).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Context {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "region_id",
            "regionId",
            "region_epoch",
            "regionEpoch",
            "peer",
            "term",
            "priority",
            "isolation_level",
            "isolationLevel",
            "not_fill_cache",
            "notFillCache",
            "sync_log",
            "syncLog",
            "record_time_stat",
            "recordTimeStat",
            "record_scan_stat",
            "recordScanStat",
            "replica_read",
            "replicaRead",
            "resolved_locks",
            "resolvedLocks",
            "max_execution_duration_ms",
            "maxExecutionDurationMs",
            "applied_index",
            "appliedIndex",
            "task_id",
            "taskId",
            "stale_read",
            "staleRead",
            "resource_group_tag",
            "resourceGroupTag",
            "disk_full_opt",
            "diskFullOpt",
            "is_retry_request",
            "isRetryRequest",
            "api_version",
            "apiVersion",
            "committed_locks",
            "committedLocks",
            "request_source",
            "requestSource",
            "txn_source",
            "txnSource",
            "busy_threshold_ms",
            "busyThresholdMs",
            "resource_control_context",
            "resourceControlContext",
            "keyspace_name",
            "keyspaceName",
            "keyspace_id",
            "keyspaceId",
            "buckets_version",
            "bucketsVersion",
            "source_stmt",
            "sourceStmt",
            "cluster_id",
            "clusterId",
            "trace_id",
            "traceId",
            "trace_control_flags",
            "traceControlFlags",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionId,
            RegionEpoch,
            Peer,
            Term,
            Priority,
            IsolationLevel,
            NotFillCache,
            SyncLog,
            RecordTimeStat,
            RecordScanStat,
            ReplicaRead,
            ResolvedLocks,
            MaxExecutionDurationMs,
            AppliedIndex,
            TaskId,
            StaleRead,
            ResourceGroupTag,
            DiskFullOpt,
            IsRetryRequest,
            ApiVersion,
            CommittedLocks,
            RequestSource,
            TxnSource,
            BusyThresholdMs,
            ResourceControlContext,
            KeyspaceName,
            KeyspaceId,
            BucketsVersion,
            SourceStmt,
            ClusterId,
            TraceId,
            TraceControlFlags,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            "regionEpoch" | "region_epoch" => Ok(GeneratedField::RegionEpoch),
                            "peer" => Ok(GeneratedField::Peer),
                            "term" => Ok(GeneratedField::Term),
                            "priority" => Ok(GeneratedField::Priority),
                            "isolationLevel" | "isolation_level" => {
                                Ok(GeneratedField::IsolationLevel)
                            }
                            "notFillCache" | "not_fill_cache" => Ok(GeneratedField::NotFillCache),
                            "syncLog" | "sync_log" => Ok(GeneratedField::SyncLog),
                            "recordTimeStat" | "record_time_stat" => {
                                Ok(GeneratedField::RecordTimeStat)
                            }
                            "recordScanStat" | "record_scan_stat" => {
                                Ok(GeneratedField::RecordScanStat)
                            }
                            "replicaRead" | "replica_read" => Ok(GeneratedField::ReplicaRead),
                            "resolvedLocks" | "resolved_locks" => Ok(GeneratedField::ResolvedLocks),
                            "maxExecutionDurationMs" | "max_execution_duration_ms" => {
                                Ok(GeneratedField::MaxExecutionDurationMs)
                            }
                            "appliedIndex" | "applied_index" => Ok(GeneratedField::AppliedIndex),
                            "taskId" | "task_id" => Ok(GeneratedField::TaskId),
                            "staleRead" | "stale_read" => Ok(GeneratedField::StaleRead),
                            "resourceGroupTag" | "resource_group_tag" => {
                                Ok(GeneratedField::ResourceGroupTag)
                            }
                            "diskFullOpt" | "disk_full_opt" => Ok(GeneratedField::DiskFullOpt),
                            "isRetryRequest" | "is_retry_request" => {
                                Ok(GeneratedField::IsRetryRequest)
                            }
                            "apiVersion" | "api_version" => Ok(GeneratedField::ApiVersion),
                            "committedLocks" | "committed_locks" => {
                                Ok(GeneratedField::CommittedLocks)
                            }
                            "requestSource" | "request_source" => Ok(GeneratedField::RequestSource),
                            "txnSource" | "txn_source" => Ok(GeneratedField::TxnSource),
                            "busyThresholdMs" | "busy_threshold_ms" => {
                                Ok(GeneratedField::BusyThresholdMs)
                            }
                            "resourceControlContext" | "resource_control_context" => {
                                Ok(GeneratedField::ResourceControlContext)
                            }
                            "keyspaceName" | "keyspace_name" => Ok(GeneratedField::KeyspaceName),
                            "keyspaceId" | "keyspace_id" => Ok(GeneratedField::KeyspaceId),
                            "bucketsVersion" | "buckets_version" => {
                                Ok(GeneratedField::BucketsVersion)
                            }
                            "sourceStmt" | "source_stmt" => Ok(GeneratedField::SourceStmt),
                            "clusterId" | "cluster_id" => Ok(GeneratedField::ClusterId),
                            "traceId" | "trace_id" => Ok(GeneratedField::TraceId),
                            "traceControlFlags" | "trace_control_flags" => {
                                Ok(GeneratedField::TraceControlFlags)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Context;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.Context")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Context, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_id__ = None;
                let mut region_epoch__ = None;
                let mut peer__ = None;
                let mut term__ = None;
                let mut priority__ = None;
                let mut isolation_level__ = None;
                let mut not_fill_cache__ = None;
                let mut sync_log__ = None;
                let mut record_time_stat__ = None;
                let mut record_scan_stat__ = None;
                let mut replica_read__ = None;
                let mut resolved_locks__ = None;
                let mut max_execution_duration_ms__ = None;
                let mut applied_index__ = None;
                let mut task_id__ = None;
                let mut stale_read__ = None;
                let mut resource_group_tag__ = None;
                let mut disk_full_opt__ = None;
                let mut is_retry_request__ = None;
                let mut api_version__ = None;
                let mut committed_locks__ = None;
                let mut request_source__ = None;
                let mut txn_source__ = None;
                let mut busy_threshold_ms__ = None;
                let mut resource_control_context__ = None;
                let mut keyspace_name__ = None;
                let mut keyspace_id__ = None;
                let mut buckets_version__ = None;
                let mut source_stmt__ = None;
                let mut cluster_id__ = None;
                let mut trace_id__ = None;
                let mut trace_control_flags__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::RegionEpoch => {
                            if region_epoch__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionEpoch"));
                            }
                            region_epoch__ = map_.next_value()?;
                        }
                        GeneratedField::Peer => {
                            if peer__.is_some() {
                                return Err(serde::de::Error::duplicate_field("peer"));
                            }
                            peer__ = map_.next_value()?;
                        }
                        GeneratedField::Term => {
                            if term__.is_some() {
                                return Err(serde::de::Error::duplicate_field("term"));
                            }
                            term__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Priority => {
                            if priority__.is_some() {
                                return Err(serde::de::Error::duplicate_field("priority"));
                            }
                            priority__ = Some(map_.next_value::<CommandPri>()? as i32);
                        }
                        GeneratedField::IsolationLevel => {
                            if isolation_level__.is_some() {
                                return Err(serde::de::Error::duplicate_field("isolationLevel"));
                            }
                            isolation_level__ = Some(map_.next_value::<IsolationLevel>()? as i32);
                        }
                        GeneratedField::NotFillCache => {
                            if not_fill_cache__.is_some() {
                                return Err(serde::de::Error::duplicate_field("notFillCache"));
                            }
                            not_fill_cache__ = Some(map_.next_value()?);
                        }
                        GeneratedField::SyncLog => {
                            if sync_log__.is_some() {
                                return Err(serde::de::Error::duplicate_field("syncLog"));
                            }
                            sync_log__ = Some(map_.next_value()?);
                        }
                        GeneratedField::RecordTimeStat => {
                            if record_time_stat__.is_some() {
                                return Err(serde::de::Error::duplicate_field("recordTimeStat"));
                            }
                            record_time_stat__ = Some(map_.next_value()?);
                        }
                        GeneratedField::RecordScanStat => {
                            if record_scan_stat__.is_some() {
                                return Err(serde::de::Error::duplicate_field("recordScanStat"));
                            }
                            record_scan_stat__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ReplicaRead => {
                            if replica_read__.is_some() {
                                return Err(serde::de::Error::duplicate_field("replicaRead"));
                            }
                            replica_read__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ResolvedLocks => {
                            if resolved_locks__.is_some() {
                                return Err(serde::de::Error::duplicate_field("resolvedLocks"));
                            }
                            resolved_locks__ = Some(
                                map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter()
                                    .map(|x| x.0)
                                    .collect(),
                            );
                        }
                        GeneratedField::MaxExecutionDurationMs => {
                            if max_execution_duration_ms__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "maxExecutionDurationMs",
                                ));
                            }
                            max_execution_duration_ms__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::AppliedIndex => {
                            if applied_index__.is_some() {
                                return Err(serde::de::Error::duplicate_field("appliedIndex"));
                            }
                            applied_index__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::TaskId => {
                            if task_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("taskId"));
                            }
                            task_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::StaleRead => {
                            if stale_read__.is_some() {
                                return Err(serde::de::Error::duplicate_field("staleRead"));
                            }
                            stale_read__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ResourceGroupTag => {
                            if resource_group_tag__.is_some() {
                                return Err(serde::de::Error::duplicate_field("resourceGroupTag"));
                            }
                            resource_group_tag__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::DiskFullOpt => {
                            if disk_full_opt__.is_some() {
                                return Err(serde::de::Error::duplicate_field("diskFullOpt"));
                            }
                            disk_full_opt__ = Some(map_.next_value::<DiskFullOpt>()? as i32);
                        }
                        GeneratedField::IsRetryRequest => {
                            if is_retry_request__.is_some() {
                                return Err(serde::de::Error::duplicate_field("isRetryRequest"));
                            }
                            is_retry_request__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ApiVersion => {
                            if api_version__.is_some() {
                                return Err(serde::de::Error::duplicate_field("apiVersion"));
                            }
                            api_version__ = Some(map_.next_value::<ApiVersion>()? as i32);
                        }
                        GeneratedField::CommittedLocks => {
                            if committed_locks__.is_some() {
                                return Err(serde::de::Error::duplicate_field("committedLocks"));
                            }
                            committed_locks__ = Some(
                                map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter()
                                    .map(|x| x.0)
                                    .collect(),
                            );
                        }
                        GeneratedField::RequestSource => {
                            if request_source__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requestSource"));
                            }
                            request_source__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TxnSource => {
                            if txn_source__.is_some() {
                                return Err(serde::de::Error::duplicate_field("txnSource"));
                            }
                            txn_source__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::BusyThresholdMs => {
                            if busy_threshold_ms__.is_some() {
                                return Err(serde::de::Error::duplicate_field("busyThresholdMs"));
                            }
                            busy_threshold_ms__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ResourceControlContext => {
                            if resource_control_context__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "resourceControlContext",
                                ));
                            }
                            resource_control_context__ = map_.next_value()?;
                        }
                        GeneratedField::KeyspaceName => {
                            if keyspace_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("keyspaceName"));
                            }
                            keyspace_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::KeyspaceId => {
                            if keyspace_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("keyspaceId"));
                            }
                            keyspace_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::BucketsVersion => {
                            if buckets_version__.is_some() {
                                return Err(serde::de::Error::duplicate_field("bucketsVersion"));
                            }
                            buckets_version__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::SourceStmt => {
                            if source_stmt__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sourceStmt"));
                            }
                            source_stmt__ = map_.next_value()?;
                        }
                        GeneratedField::ClusterId => {
                            if cluster_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("clusterId"));
                            }
                            cluster_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::TraceId => {
                            if trace_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("traceId"));
                            }
                            trace_id__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::TraceControlFlags => {
                            if trace_control_flags__.is_some() {
                                return Err(serde::de::Error::duplicate_field("traceControlFlags"));
                            }
                            trace_control_flags__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(Context {
                    region_id: region_id__.unwrap_or_default(),
                    region_epoch: region_epoch__,
                    peer: peer__,
                    term: term__.unwrap_or_default(),
                    priority: priority__.unwrap_or_default(),
                    isolation_level: isolation_level__.unwrap_or_default(),
                    not_fill_cache: not_fill_cache__.unwrap_or_default(),
                    sync_log: sync_log__.unwrap_or_default(),
                    record_time_stat: record_time_stat__.unwrap_or_default(),
                    record_scan_stat: record_scan_stat__.unwrap_or_default(),
                    replica_read: replica_read__.unwrap_or_default(),
                    resolved_locks: resolved_locks__.unwrap_or_default(),
                    max_execution_duration_ms: max_execution_duration_ms__.unwrap_or_default(),
                    applied_index: applied_index__.unwrap_or_default(),
                    task_id: task_id__.unwrap_or_default(),
                    stale_read: stale_read__.unwrap_or_default(),
                    resource_group_tag: resource_group_tag__.unwrap_or_default(),
                    disk_full_opt: disk_full_opt__.unwrap_or_default(),
                    is_retry_request: is_retry_request__.unwrap_or_default(),
                    api_version: api_version__.unwrap_or_default(),
                    committed_locks: committed_locks__.unwrap_or_default(),
                    request_source: request_source__.unwrap_or_default(),
                    txn_source: txn_source__.unwrap_or_default(),
                    busy_threshold_ms: busy_threshold_ms__.unwrap_or_default(),
                    resource_control_context: resource_control_context__,
                    keyspace_name: keyspace_name__.unwrap_or_default(),
                    keyspace_id: keyspace_id__.unwrap_or_default(),
                    buckets_version: buckets_version__.unwrap_or_default(),
                    source_stmt: source_stmt__,
                    cluster_id: cluster_id__.unwrap_or_default(),
                    trace_id: trace_id__.unwrap_or_default(),
                    trace_control_flags: trace_control_flags__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.Context", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for DataIsNotReady {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_id != 0 {
            len += 1;
        }
        if self.peer_id != 0 {
            len += 1;
        }
        if self.safe_ts != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.DataIsNotReady", len)?;
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        if self.peer_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("peerId", ToString::to_string(&self.peer_id).as_str())?;
        }
        if self.safe_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("safeTs", ToString::to_string(&self.safe_ts).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for DataIsNotReady {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "region_id",
            "regionId",
            "peer_id",
            "peerId",
            "safe_ts",
            "safeTs",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionId,
            PeerId,
            SafeTs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            "peerId" | "peer_id" => Ok(GeneratedField::PeerId),
                            "safeTs" | "safe_ts" => Ok(GeneratedField::SafeTs),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = DataIsNotReady;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.DataIsNotReady")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<DataIsNotReady, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_id__ = None;
                let mut peer_id__ = None;
                let mut safe_ts__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::PeerId => {
                            if peer_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("peerId"));
                            }
                            peer_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::SafeTs => {
                            if safe_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("safeTs"));
                            }
                            safe_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(DataIsNotReady {
                    region_id: region_id__.unwrap_or_default(),
                    peer_id: peer_id__.unwrap_or_default(),
                    safe_ts: safe_ts__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.DataIsNotReady", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for DebugInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.mvcc_info.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.DebugInfo", len)?;
        if !self.mvcc_info.is_empty() {
            struct_ser.serialize_field("mvccInfo", &self.mvcc_info)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for DebugInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["mvcc_info", "mvccInfo"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            MvccInfo,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "mvccInfo" | "mvcc_info" => Ok(GeneratedField::MvccInfo),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = DebugInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.DebugInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<DebugInfo, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut mvcc_info__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::MvccInfo => {
                            if mvcc_info__.is_some() {
                                return Err(serde::de::Error::duplicate_field("mvccInfo"));
                            }
                            mvcc_info__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(DebugInfo {
                    mvcc_info: mvcc_info__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.DebugInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for DiskFull {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.store_id.is_empty() {
            len += 1;
        }
        if !self.reason.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.DiskFull", len)?;
        if !self.store_id.is_empty() {
            struct_ser.serialize_field(
                "storeId",
                &self
                    .store_id
                    .iter()
                    .map(ToString::to_string)
                    .collect::<Vec<_>>(),
            )?;
        }
        if !self.reason.is_empty() {
            struct_ser.serialize_field("reason", &self.reason)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for DiskFull {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["store_id", "storeId", "reason"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            StoreId,
            Reason,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "storeId" | "store_id" => Ok(GeneratedField::StoreId),
                            "reason" => Ok(GeneratedField::Reason),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = DiskFull;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.DiskFull")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<DiskFull, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut store_id__ = None;
                let mut reason__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::StoreId => {
                            if store_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("storeId"));
                            }
                            store_id__ = Some(
                                map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter()
                                    .map(|x| x.0)
                                    .collect(),
                            );
                        }
                        GeneratedField::Reason => {
                            if reason__.is_some() {
                                return Err(serde::de::Error::duplicate_field("reason"));
                            }
                            reason__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(DiskFull {
                    store_id: store_id__.unwrap_or_default(),
                    reason: reason__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.DiskFull", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for DiskFullOpt {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::NotAllowedOnFull => "NotAllowedOnFull",
            Self::AllowedOnAlmostFull => "AllowedOnAlmostFull",
            Self::AllowedOnAlreadyFull => "AllowedOnAlreadyFull",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for DiskFullOpt {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "NotAllowedOnFull",
            "AllowedOnAlmostFull",
            "AllowedOnAlreadyFull",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = DiskFullOpt;

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
                    "NotAllowedOnFull" => Ok(DiskFullOpt::NotAllowedOnFull),
                    "AllowedOnAlmostFull" => Ok(DiskFullOpt::AllowedOnAlmostFull),
                    "AllowedOnAlreadyFull" => Ok(DiskFullOpt::AllowedOnAlreadyFull),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for EpochNotMatch {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.current_regions.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.EpochNotMatch", len)?;
        if !self.current_regions.is_empty() {
            struct_ser.serialize_field("currentRegions", &self.current_regions)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for EpochNotMatch {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["current_regions", "currentRegions"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            CurrentRegions,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "currentRegions" | "current_regions" => {
                                Ok(GeneratedField::CurrentRegions)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = EpochNotMatch;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.EpochNotMatch")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<EpochNotMatch, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut current_regions__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::CurrentRegions => {
                            if current_regions__.is_some() {
                                return Err(serde::de::Error::duplicate_field("currentRegions"));
                            }
                            current_regions__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(EpochNotMatch {
                    current_regions: current_regions__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.EpochNotMatch", FIELDS, GeneratedVisitor)
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
        if !self.message.is_empty() {
            len += 1;
        }
        if self.not_leader.is_some() {
            len += 1;
        }
        if self.region_not_found.is_some() {
            len += 1;
        }
        if self.key_not_in_region.is_some() {
            len += 1;
        }
        if self.epoch_not_match.is_some() {
            len += 1;
        }
        if self.server_is_busy.is_some() {
            len += 1;
        }
        if self.stale_command.is_some() {
            len += 1;
        }
        if self.store_not_match.is_some() {
            len += 1;
        }
        if self.raft_entry_too_large.is_some() {
            len += 1;
        }
        if self.max_timestamp_not_synced.is_some() {
            len += 1;
        }
        if self.read_index_not_ready.is_some() {
            len += 1;
        }
        if self.proposal_in_merging_mode.is_some() {
            len += 1;
        }
        if self.data_is_not_ready.is_some() {
            len += 1;
        }
        if self.region_not_initialized.is_some() {
            len += 1;
        }
        if self.disk_full.is_some() {
            len += 1;
        }
        if self.recovery_in_progress.is_some() {
            len += 1;
        }
        if self.flashback_in_progress.is_some() {
            len += 1;
        }
        if self.flashback_not_prepared.is_some() {
            len += 1;
        }
        if self.is_witness.is_some() {
            len += 1;
        }
        if self.mismatch_peer_id.is_some() {
            len += 1;
        }
        if self.bucket_version_not_match.is_some() {
            len += 1;
        }
        if self.undetermined_result.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.Error", len)?;
        if !self.message.is_empty() {
            struct_ser.serialize_field("message", &self.message)?;
        }
        if let Some(v) = self.not_leader.as_ref() {
            struct_ser.serialize_field("notLeader", v)?;
        }
        if let Some(v) = self.region_not_found.as_ref() {
            struct_ser.serialize_field("regionNotFound", v)?;
        }
        if let Some(v) = self.key_not_in_region.as_ref() {
            struct_ser.serialize_field("keyNotInRegion", v)?;
        }
        if let Some(v) = self.epoch_not_match.as_ref() {
            struct_ser.serialize_field("epochNotMatch", v)?;
        }
        if let Some(v) = self.server_is_busy.as_ref() {
            struct_ser.serialize_field("serverIsBusy", v)?;
        }
        if let Some(v) = self.stale_command.as_ref() {
            struct_ser.serialize_field("staleCommand", v)?;
        }
        if let Some(v) = self.store_not_match.as_ref() {
            struct_ser.serialize_field("storeNotMatch", v)?;
        }
        if let Some(v) = self.raft_entry_too_large.as_ref() {
            struct_ser.serialize_field("raftEntryTooLarge", v)?;
        }
        if let Some(v) = self.max_timestamp_not_synced.as_ref() {
            struct_ser.serialize_field("maxTimestampNotSynced", v)?;
        }
        if let Some(v) = self.read_index_not_ready.as_ref() {
            struct_ser.serialize_field("readIndexNotReady", v)?;
        }
        if let Some(v) = self.proposal_in_merging_mode.as_ref() {
            struct_ser.serialize_field("proposalInMergingMode", v)?;
        }
        if let Some(v) = self.data_is_not_ready.as_ref() {
            struct_ser.serialize_field("dataIsNotReady", v)?;
        }
        if let Some(v) = self.region_not_initialized.as_ref() {
            struct_ser.serialize_field("regionNotInitialized", v)?;
        }
        if let Some(v) = self.disk_full.as_ref() {
            struct_ser.serialize_field("diskFull", v)?;
        }
        if let Some(v) = self.recovery_in_progress.as_ref() {
            struct_ser.serialize_field("RecoveryInProgress", v)?;
        }
        if let Some(v) = self.flashback_in_progress.as_ref() {
            struct_ser.serialize_field("FlashbackInProgress", v)?;
        }
        if let Some(v) = self.flashback_not_prepared.as_ref() {
            struct_ser.serialize_field("FlashbackNotPrepared", v)?;
        }
        if let Some(v) = self.is_witness.as_ref() {
            struct_ser.serialize_field("isWitness", v)?;
        }
        if let Some(v) = self.mismatch_peer_id.as_ref() {
            struct_ser.serialize_field("mismatchPeerId", v)?;
        }
        if let Some(v) = self.bucket_version_not_match.as_ref() {
            struct_ser.serialize_field("bucketVersionNotMatch", v)?;
        }
        if let Some(v) = self.undetermined_result.as_ref() {
            struct_ser.serialize_field("undeterminedResult", v)?;
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
            "message",
            "not_leader",
            "notLeader",
            "region_not_found",
            "regionNotFound",
            "key_not_in_region",
            "keyNotInRegion",
            "epoch_not_match",
            "epochNotMatch",
            "server_is_busy",
            "serverIsBusy",
            "stale_command",
            "staleCommand",
            "store_not_match",
            "storeNotMatch",
            "raft_entry_too_large",
            "raftEntryTooLarge",
            "max_timestamp_not_synced",
            "maxTimestampNotSynced",
            "read_index_not_ready",
            "readIndexNotReady",
            "proposal_in_merging_mode",
            "proposalInMergingMode",
            "data_is_not_ready",
            "dataIsNotReady",
            "region_not_initialized",
            "regionNotInitialized",
            "disk_full",
            "diskFull",
            "RecoveryInProgress",
            "FlashbackInProgress",
            "FlashbackNotPrepared",
            "is_witness",
            "isWitness",
            "mismatch_peer_id",
            "mismatchPeerId",
            "bucket_version_not_match",
            "bucketVersionNotMatch",
            "undetermined_result",
            "undeterminedResult",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Message,
            NotLeader,
            RegionNotFound,
            KeyNotInRegion,
            EpochNotMatch,
            ServerIsBusy,
            StaleCommand,
            StoreNotMatch,
            RaftEntryTooLarge,
            MaxTimestampNotSynced,
            ReadIndexNotReady,
            ProposalInMergingMode,
            DataIsNotReady,
            RegionNotInitialized,
            DiskFull,
            RecoveryInProgress,
            FlashbackInProgress,
            FlashbackNotPrepared,
            IsWitness,
            MismatchPeerId,
            BucketVersionNotMatch,
            UndeterminedResult,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "message" => Ok(GeneratedField::Message),
                            "notLeader" | "not_leader" => Ok(GeneratedField::NotLeader),
                            "regionNotFound" | "region_not_found" => {
                                Ok(GeneratedField::RegionNotFound)
                            }
                            "keyNotInRegion" | "key_not_in_region" => {
                                Ok(GeneratedField::KeyNotInRegion)
                            }
                            "epochNotMatch" | "epoch_not_match" => {
                                Ok(GeneratedField::EpochNotMatch)
                            }
                            "serverIsBusy" | "server_is_busy" => Ok(GeneratedField::ServerIsBusy),
                            "staleCommand" | "stale_command" => Ok(GeneratedField::StaleCommand),
                            "storeNotMatch" | "store_not_match" => {
                                Ok(GeneratedField::StoreNotMatch)
                            }
                            "raftEntryTooLarge" | "raft_entry_too_large" => {
                                Ok(GeneratedField::RaftEntryTooLarge)
                            }
                            "maxTimestampNotSynced" | "max_timestamp_not_synced" => {
                                Ok(GeneratedField::MaxTimestampNotSynced)
                            }
                            "readIndexNotReady" | "read_index_not_ready" => {
                                Ok(GeneratedField::ReadIndexNotReady)
                            }
                            "proposalInMergingMode" | "proposal_in_merging_mode" => {
                                Ok(GeneratedField::ProposalInMergingMode)
                            }
                            "dataIsNotReady" | "data_is_not_ready" => {
                                Ok(GeneratedField::DataIsNotReady)
                            }
                            "regionNotInitialized" | "region_not_initialized" => {
                                Ok(GeneratedField::RegionNotInitialized)
                            }
                            "diskFull" | "disk_full" => Ok(GeneratedField::DiskFull),
                            "RecoveryInProgress" => Ok(GeneratedField::RecoveryInProgress),
                            "FlashbackInProgress" => Ok(GeneratedField::FlashbackInProgress),
                            "FlashbackNotPrepared" => Ok(GeneratedField::FlashbackNotPrepared),
                            "isWitness" | "is_witness" => Ok(GeneratedField::IsWitness),
                            "mismatchPeerId" | "mismatch_peer_id" => {
                                Ok(GeneratedField::MismatchPeerId)
                            }
                            "bucketVersionNotMatch" | "bucket_version_not_match" => {
                                Ok(GeneratedField::BucketVersionNotMatch)
                            }
                            "undeterminedResult" | "undetermined_result" => {
                                Ok(GeneratedField::UndeterminedResult)
                            }
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
                formatter.write_str("struct tikv.Error")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Error, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut message__ = None;
                let mut not_leader__ = None;
                let mut region_not_found__ = None;
                let mut key_not_in_region__ = None;
                let mut epoch_not_match__ = None;
                let mut server_is_busy__ = None;
                let mut stale_command__ = None;
                let mut store_not_match__ = None;
                let mut raft_entry_too_large__ = None;
                let mut max_timestamp_not_synced__ = None;
                let mut read_index_not_ready__ = None;
                let mut proposal_in_merging_mode__ = None;
                let mut data_is_not_ready__ = None;
                let mut region_not_initialized__ = None;
                let mut disk_full__ = None;
                let mut recovery_in_progress__ = None;
                let mut flashback_in_progress__ = None;
                let mut flashback_not_prepared__ = None;
                let mut is_witness__ = None;
                let mut mismatch_peer_id__ = None;
                let mut bucket_version_not_match__ = None;
                let mut undetermined_result__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Message => {
                            if message__.is_some() {
                                return Err(serde::de::Error::duplicate_field("message"));
                            }
                            message__ = Some(map_.next_value()?);
                        }
                        GeneratedField::NotLeader => {
                            if not_leader__.is_some() {
                                return Err(serde::de::Error::duplicate_field("notLeader"));
                            }
                            not_leader__ = map_.next_value()?;
                        }
                        GeneratedField::RegionNotFound => {
                            if region_not_found__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionNotFound"));
                            }
                            region_not_found__ = map_.next_value()?;
                        }
                        GeneratedField::KeyNotInRegion => {
                            if key_not_in_region__.is_some() {
                                return Err(serde::de::Error::duplicate_field("keyNotInRegion"));
                            }
                            key_not_in_region__ = map_.next_value()?;
                        }
                        GeneratedField::EpochNotMatch => {
                            if epoch_not_match__.is_some() {
                                return Err(serde::de::Error::duplicate_field("epochNotMatch"));
                            }
                            epoch_not_match__ = map_.next_value()?;
                        }
                        GeneratedField::ServerIsBusy => {
                            if server_is_busy__.is_some() {
                                return Err(serde::de::Error::duplicate_field("serverIsBusy"));
                            }
                            server_is_busy__ = map_.next_value()?;
                        }
                        GeneratedField::StaleCommand => {
                            if stale_command__.is_some() {
                                return Err(serde::de::Error::duplicate_field("staleCommand"));
                            }
                            stale_command__ = map_.next_value()?;
                        }
                        GeneratedField::StoreNotMatch => {
                            if store_not_match__.is_some() {
                                return Err(serde::de::Error::duplicate_field("storeNotMatch"));
                            }
                            store_not_match__ = map_.next_value()?;
                        }
                        GeneratedField::RaftEntryTooLarge => {
                            if raft_entry_too_large__.is_some() {
                                return Err(serde::de::Error::duplicate_field("raftEntryTooLarge"));
                            }
                            raft_entry_too_large__ = map_.next_value()?;
                        }
                        GeneratedField::MaxTimestampNotSynced => {
                            if max_timestamp_not_synced__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "maxTimestampNotSynced",
                                ));
                            }
                            max_timestamp_not_synced__ = map_.next_value()?;
                        }
                        GeneratedField::ReadIndexNotReady => {
                            if read_index_not_ready__.is_some() {
                                return Err(serde::de::Error::duplicate_field("readIndexNotReady"));
                            }
                            read_index_not_ready__ = map_.next_value()?;
                        }
                        GeneratedField::ProposalInMergingMode => {
                            if proposal_in_merging_mode__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "proposalInMergingMode",
                                ));
                            }
                            proposal_in_merging_mode__ = map_.next_value()?;
                        }
                        GeneratedField::DataIsNotReady => {
                            if data_is_not_ready__.is_some() {
                                return Err(serde::de::Error::duplicate_field("dataIsNotReady"));
                            }
                            data_is_not_ready__ = map_.next_value()?;
                        }
                        GeneratedField::RegionNotInitialized => {
                            if region_not_initialized__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "regionNotInitialized",
                                ));
                            }
                            region_not_initialized__ = map_.next_value()?;
                        }
                        GeneratedField::DiskFull => {
                            if disk_full__.is_some() {
                                return Err(serde::de::Error::duplicate_field("diskFull"));
                            }
                            disk_full__ = map_.next_value()?;
                        }
                        GeneratedField::RecoveryInProgress => {
                            if recovery_in_progress__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "RecoveryInProgress",
                                ));
                            }
                            recovery_in_progress__ = map_.next_value()?;
                        }
                        GeneratedField::FlashbackInProgress => {
                            if flashback_in_progress__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "FlashbackInProgress",
                                ));
                            }
                            flashback_in_progress__ = map_.next_value()?;
                        }
                        GeneratedField::FlashbackNotPrepared => {
                            if flashback_not_prepared__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "FlashbackNotPrepared",
                                ));
                            }
                            flashback_not_prepared__ = map_.next_value()?;
                        }
                        GeneratedField::IsWitness => {
                            if is_witness__.is_some() {
                                return Err(serde::de::Error::duplicate_field("isWitness"));
                            }
                            is_witness__ = map_.next_value()?;
                        }
                        GeneratedField::MismatchPeerId => {
                            if mismatch_peer_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("mismatchPeerId"));
                            }
                            mismatch_peer_id__ = map_.next_value()?;
                        }
                        GeneratedField::BucketVersionNotMatch => {
                            if bucket_version_not_match__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "bucketVersionNotMatch",
                                ));
                            }
                            bucket_version_not_match__ = map_.next_value()?;
                        }
                        GeneratedField::UndeterminedResult => {
                            if undetermined_result__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "undeterminedResult",
                                ));
                            }
                            undetermined_result__ = map_.next_value()?;
                        }
                    }
                }
                Ok(Error {
                    message: message__.unwrap_or_default(),
                    not_leader: not_leader__,
                    region_not_found: region_not_found__,
                    key_not_in_region: key_not_in_region__,
                    epoch_not_match: epoch_not_match__,
                    server_is_busy: server_is_busy__,
                    stale_command: stale_command__,
                    store_not_match: store_not_match__,
                    raft_entry_too_large: raft_entry_too_large__,
                    max_timestamp_not_synced: max_timestamp_not_synced__,
                    read_index_not_ready: read_index_not_ready__,
                    proposal_in_merging_mode: proposal_in_merging_mode__,
                    data_is_not_ready: data_is_not_ready__,
                    region_not_initialized: region_not_initialized__,
                    disk_full: disk_full__,
                    recovery_in_progress: recovery_in_progress__,
                    flashback_in_progress: flashback_in_progress__,
                    flashback_not_prepared: flashback_not_prepared__,
                    is_witness: is_witness__,
                    mismatch_peer_id: mismatch_peer_id__,
                    bucket_version_not_match: bucket_version_not_match__,
                    undetermined_result: undetermined_result__,
                })
            }
        }
        deserializer.deserialize_struct("tikv.Error", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ExecDetails {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.time_detail.is_some() {
            len += 1;
        }
        if self.scan_detail.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.ExecDetails", len)?;
        if let Some(v) = self.time_detail.as_ref() {
            struct_ser.serialize_field("timeDetail", v)?;
        }
        if let Some(v) = self.scan_detail.as_ref() {
            struct_ser.serialize_field("scanDetail", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ExecDetails {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["time_detail", "timeDetail", "scan_detail", "scanDetail"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            TimeDetail,
            ScanDetail,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "timeDetail" | "time_detail" => Ok(GeneratedField::TimeDetail),
                            "scanDetail" | "scan_detail" => Ok(GeneratedField::ScanDetail),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ExecDetails;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.ExecDetails")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ExecDetails, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut time_detail__ = None;
                let mut scan_detail__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::TimeDetail => {
                            if time_detail__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timeDetail"));
                            }
                            time_detail__ = map_.next_value()?;
                        }
                        GeneratedField::ScanDetail => {
                            if scan_detail__.is_some() {
                                return Err(serde::de::Error::duplicate_field("scanDetail"));
                            }
                            scan_detail__ = map_.next_value()?;
                        }
                    }
                }
                Ok(ExecDetails {
                    time_detail: time_detail__,
                    scan_detail: scan_detail__,
                })
            }
        }
        deserializer.deserialize_struct("tikv.ExecDetails", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ExecDetailsV2 {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.time_detail.is_some() {
            len += 1;
        }
        if self.scan_detail_v2.is_some() {
            len += 1;
        }
        if self.write_detail.is_some() {
            len += 1;
        }
        if self.time_detail_v2.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.ExecDetailsV2", len)?;
        if let Some(v) = self.time_detail.as_ref() {
            struct_ser.serialize_field("timeDetail", v)?;
        }
        if let Some(v) = self.scan_detail_v2.as_ref() {
            struct_ser.serialize_field("scanDetailV2", v)?;
        }
        if let Some(v) = self.write_detail.as_ref() {
            struct_ser.serialize_field("writeDetail", v)?;
        }
        if let Some(v) = self.time_detail_v2.as_ref() {
            struct_ser.serialize_field("timeDetailV2", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ExecDetailsV2 {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "time_detail",
            "timeDetail",
            "scan_detail_v2",
            "scanDetailV2",
            "write_detail",
            "writeDetail",
            "time_detail_v2",
            "timeDetailV2",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            TimeDetail,
            ScanDetailV2,
            WriteDetail,
            TimeDetailV2,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "timeDetail" | "time_detail" => Ok(GeneratedField::TimeDetail),
                            "scanDetailV2" | "scan_detail_v2" => Ok(GeneratedField::ScanDetailV2),
                            "writeDetail" | "write_detail" => Ok(GeneratedField::WriteDetail),
                            "timeDetailV2" | "time_detail_v2" => Ok(GeneratedField::TimeDetailV2),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ExecDetailsV2;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.ExecDetailsV2")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ExecDetailsV2, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut time_detail__ = None;
                let mut scan_detail_v2__ = None;
                let mut write_detail__ = None;
                let mut time_detail_v2__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::TimeDetail => {
                            if time_detail__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timeDetail"));
                            }
                            time_detail__ = map_.next_value()?;
                        }
                        GeneratedField::ScanDetailV2 => {
                            if scan_detail_v2__.is_some() {
                                return Err(serde::de::Error::duplicate_field("scanDetailV2"));
                            }
                            scan_detail_v2__ = map_.next_value()?;
                        }
                        GeneratedField::WriteDetail => {
                            if write_detail__.is_some() {
                                return Err(serde::de::Error::duplicate_field("writeDetail"));
                            }
                            write_detail__ = map_.next_value()?;
                        }
                        GeneratedField::TimeDetailV2 => {
                            if time_detail_v2__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timeDetailV2"));
                            }
                            time_detail_v2__ = map_.next_value()?;
                        }
                    }
                }
                Ok(ExecDetailsV2 {
                    time_detail: time_detail__,
                    scan_detail_v2: scan_detail_v2__,
                    write_detail: write_detail__,
                    time_detail_v2: time_detail_v2__,
                })
            }
        }
        deserializer.deserialize_struct("tikv.ExecDetailsV2", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for FlashbackInProgress {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_id != 0 {
            len += 1;
        }
        if self.flashback_start_ts != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.FlashbackInProgress", len)?;
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        if self.flashback_start_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "flashbackStartTs",
                ToString::to_string(&self.flashback_start_ts).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for FlashbackInProgress {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "region_id",
            "regionId",
            "flashback_start_ts",
            "flashbackStartTs",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionId,
            FlashbackStartTs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            "flashbackStartTs" | "flashback_start_ts" => {
                                Ok(GeneratedField::FlashbackStartTs)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = FlashbackInProgress;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.FlashbackInProgress")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<FlashbackInProgress, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_id__ = None;
                let mut flashback_start_ts__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::FlashbackStartTs => {
                            if flashback_start_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("flashbackStartTs"));
                            }
                            flashback_start_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(FlashbackInProgress {
                    region_id: region_id__.unwrap_or_default(),
                    flashback_start_ts: flashback_start_ts__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.FlashbackInProgress", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for FlashbackNotPrepared {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_id != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.FlashbackNotPrepared", len)?;
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for FlashbackNotPrepared {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_id", "regionId"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionId,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = FlashbackNotPrepared;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.FlashbackNotPrepared")
            }

            fn visit_map<V>(
                self,
                mut map_: V,
            ) -> std::result::Result<FlashbackNotPrepared, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(FlashbackNotPrepared {
                    region_id: region_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.FlashbackNotPrepared", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for IsWitness {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_id != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.IsWitness", len)?;
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for IsWitness {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_id", "regionId"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionId,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = IsWitness;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.IsWitness")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<IsWitness, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(IsWitness {
                    region_id: region_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.IsWitness", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for IsolationLevel {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Si => "SI",
            Self::Rc => "RC",
            Self::RcCheckTs => "RCCheckTS",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for IsolationLevel {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["SI", "RC", "RCCheckTS"];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = IsolationLevel;

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
                    "SI" => Ok(IsolationLevel::Si),
                    "RC" => Ok(IsolationLevel::Rc),
                    "RCCheckTS" => Ok(IsolationLevel::RcCheckTs),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for KeyError {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.locked.is_some() {
            len += 1;
        }
        if !self.retryable.is_empty() {
            len += 1;
        }
        if !self.abort.is_empty() {
            len += 1;
        }
        if self.conflict.is_some() {
            len += 1;
        }
        if self.already_exist.is_some() {
            len += 1;
        }
        if self.commit_ts_expired.is_some() {
            len += 1;
        }
        if self.txn_not_found.is_some() {
            len += 1;
        }
        if self.commit_ts_too_large.is_some() {
            len += 1;
        }
        if self.assertion_failed.is_some() {
            len += 1;
        }
        if self.primary_mismatch.is_some() {
            len += 1;
        }
        if self.txn_lock_not_found.is_some() {
            len += 1;
        }
        if self.debug_info.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.KeyError", len)?;
        if let Some(v) = self.locked.as_ref() {
            struct_ser.serialize_field("locked", v)?;
        }
        if !self.retryable.is_empty() {
            struct_ser.serialize_field("retryable", &self.retryable)?;
        }
        if !self.abort.is_empty() {
            struct_ser.serialize_field("abort", &self.abort)?;
        }
        if let Some(v) = self.conflict.as_ref() {
            struct_ser.serialize_field("conflict", v)?;
        }
        if let Some(v) = self.already_exist.as_ref() {
            struct_ser.serialize_field("alreadyExist", v)?;
        }
        if let Some(v) = self.commit_ts_expired.as_ref() {
            struct_ser.serialize_field("commitTsExpired", v)?;
        }
        if let Some(v) = self.txn_not_found.as_ref() {
            struct_ser.serialize_field("txnNotFound", v)?;
        }
        if let Some(v) = self.commit_ts_too_large.as_ref() {
            struct_ser.serialize_field("commitTsTooLarge", v)?;
        }
        if let Some(v) = self.assertion_failed.as_ref() {
            struct_ser.serialize_field("assertionFailed", v)?;
        }
        if let Some(v) = self.primary_mismatch.as_ref() {
            struct_ser.serialize_field("primaryMismatch", v)?;
        }
        if let Some(v) = self.txn_lock_not_found.as_ref() {
            struct_ser.serialize_field("txnLockNotFound", v)?;
        }
        if let Some(v) = self.debug_info.as_ref() {
            struct_ser.serialize_field("debugInfo", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for KeyError {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "locked",
            "retryable",
            "abort",
            "conflict",
            "already_exist",
            "alreadyExist",
            "commit_ts_expired",
            "commitTsExpired",
            "txn_not_found",
            "txnNotFound",
            "commit_ts_too_large",
            "commitTsTooLarge",
            "assertion_failed",
            "assertionFailed",
            "primary_mismatch",
            "primaryMismatch",
            "txn_lock_not_found",
            "txnLockNotFound",
            "debug_info",
            "debugInfo",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Locked,
            Retryable,
            Abort,
            Conflict,
            AlreadyExist,
            CommitTsExpired,
            TxnNotFound,
            CommitTsTooLarge,
            AssertionFailed,
            PrimaryMismatch,
            TxnLockNotFound,
            DebugInfo,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "locked" => Ok(GeneratedField::Locked),
                            "retryable" => Ok(GeneratedField::Retryable),
                            "abort" => Ok(GeneratedField::Abort),
                            "conflict" => Ok(GeneratedField::Conflict),
                            "alreadyExist" | "already_exist" => Ok(GeneratedField::AlreadyExist),
                            "commitTsExpired" | "commit_ts_expired" => {
                                Ok(GeneratedField::CommitTsExpired)
                            }
                            "txnNotFound" | "txn_not_found" => Ok(GeneratedField::TxnNotFound),
                            "commitTsTooLarge" | "commit_ts_too_large" => {
                                Ok(GeneratedField::CommitTsTooLarge)
                            }
                            "assertionFailed" | "assertion_failed" => {
                                Ok(GeneratedField::AssertionFailed)
                            }
                            "primaryMismatch" | "primary_mismatch" => {
                                Ok(GeneratedField::PrimaryMismatch)
                            }
                            "txnLockNotFound" | "txn_lock_not_found" => {
                                Ok(GeneratedField::TxnLockNotFound)
                            }
                            "debugInfo" | "debug_info" => Ok(GeneratedField::DebugInfo),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = KeyError;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.KeyError")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<KeyError, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut locked__ = None;
                let mut retryable__ = None;
                let mut abort__ = None;
                let mut conflict__ = None;
                let mut already_exist__ = None;
                let mut commit_ts_expired__ = None;
                let mut txn_not_found__ = None;
                let mut commit_ts_too_large__ = None;
                let mut assertion_failed__ = None;
                let mut primary_mismatch__ = None;
                let mut txn_lock_not_found__ = None;
                let mut debug_info__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Locked => {
                            if locked__.is_some() {
                                return Err(serde::de::Error::duplicate_field("locked"));
                            }
                            locked__ = map_.next_value()?;
                        }
                        GeneratedField::Retryable => {
                            if retryable__.is_some() {
                                return Err(serde::de::Error::duplicate_field("retryable"));
                            }
                            retryable__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Abort => {
                            if abort__.is_some() {
                                return Err(serde::de::Error::duplicate_field("abort"));
                            }
                            abort__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Conflict => {
                            if conflict__.is_some() {
                                return Err(serde::de::Error::duplicate_field("conflict"));
                            }
                            conflict__ = map_.next_value()?;
                        }
                        GeneratedField::AlreadyExist => {
                            if already_exist__.is_some() {
                                return Err(serde::de::Error::duplicate_field("alreadyExist"));
                            }
                            already_exist__ = map_.next_value()?;
                        }
                        GeneratedField::CommitTsExpired => {
                            if commit_ts_expired__.is_some() {
                                return Err(serde::de::Error::duplicate_field("commitTsExpired"));
                            }
                            commit_ts_expired__ = map_.next_value()?;
                        }
                        GeneratedField::TxnNotFound => {
                            if txn_not_found__.is_some() {
                                return Err(serde::de::Error::duplicate_field("txnNotFound"));
                            }
                            txn_not_found__ = map_.next_value()?;
                        }
                        GeneratedField::CommitTsTooLarge => {
                            if commit_ts_too_large__.is_some() {
                                return Err(serde::de::Error::duplicate_field("commitTsTooLarge"));
                            }
                            commit_ts_too_large__ = map_.next_value()?;
                        }
                        GeneratedField::AssertionFailed => {
                            if assertion_failed__.is_some() {
                                return Err(serde::de::Error::duplicate_field("assertionFailed"));
                            }
                            assertion_failed__ = map_.next_value()?;
                        }
                        GeneratedField::PrimaryMismatch => {
                            if primary_mismatch__.is_some() {
                                return Err(serde::de::Error::duplicate_field("primaryMismatch"));
                            }
                            primary_mismatch__ = map_.next_value()?;
                        }
                        GeneratedField::TxnLockNotFound => {
                            if txn_lock_not_found__.is_some() {
                                return Err(serde::de::Error::duplicate_field("txnLockNotFound"));
                            }
                            txn_lock_not_found__ = map_.next_value()?;
                        }
                        GeneratedField::DebugInfo => {
                            if debug_info__.is_some() {
                                return Err(serde::de::Error::duplicate_field("debugInfo"));
                            }
                            debug_info__ = map_.next_value()?;
                        }
                    }
                }
                Ok(KeyError {
                    locked: locked__,
                    retryable: retryable__.unwrap_or_default(),
                    abort: abort__.unwrap_or_default(),
                    conflict: conflict__,
                    already_exist: already_exist__,
                    commit_ts_expired: commit_ts_expired__,
                    txn_not_found: txn_not_found__,
                    commit_ts_too_large: commit_ts_too_large__,
                    assertion_failed: assertion_failed__,
                    primary_mismatch: primary_mismatch__,
                    txn_lock_not_found: txn_lock_not_found__,
                    debug_info: debug_info__,
                })
            }
        }
        deserializer.deserialize_struct("tikv.KeyError", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for KeyNotInRegion {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.key.is_empty() {
            len += 1;
        }
        if self.region_id != 0 {
            len += 1;
        }
        if !self.start_key.is_empty() {
            len += 1;
        }
        if !self.end_key.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.KeyNotInRegion", len)?;
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        if !self.start_key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "startKey",
                pbjson::private::base64::encode(&self.start_key).as_str(),
            )?;
        }
        if !self.end_key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "endKey",
                pbjson::private::base64::encode(&self.end_key).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for KeyNotInRegion {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "key",
            "region_id",
            "regionId",
            "start_key",
            "startKey",
            "end_key",
            "endKey",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Key,
            RegionId,
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

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "key" => Ok(GeneratedField::Key),
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
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
            type Value = KeyNotInRegion;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.KeyNotInRegion")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<KeyNotInRegion, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut key__ = None;
                let mut region_id__ = None;
                let mut start_key__ = None;
                let mut end_key__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::StartKey => {
                            if start_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startKey"));
                            }
                            start_key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::EndKey => {
                            if end_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("endKey"));
                            }
                            end_key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(KeyNotInRegion {
                    key: key__.unwrap_or_default(),
                    region_id: region_id__.unwrap_or_default(),
                    start_key: start_key__.unwrap_or_default(),
                    end_key: end_key__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.KeyNotInRegion", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for KvPair {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.error.is_some() {
            len += 1;
        }
        if !self.key.is_empty() {
            len += 1;
        }
        if !self.value.is_empty() {
            len += 1;
        }
        if self.commit_ts != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.KvPair", len)?;
        if let Some(v) = self.error.as_ref() {
            struct_ser.serialize_field("error", v)?;
        }
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        if !self.value.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "value",
                pbjson::private::base64::encode(&self.value).as_str(),
            )?;
        }
        if self.commit_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("commitTs", ToString::to_string(&self.commit_ts).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for KvPair {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["error", "key", "value", "commit_ts", "commitTs"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Error,
            Key,
            Value,
            CommitTs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "error" => Ok(GeneratedField::Error),
                            "key" => Ok(GeneratedField::Key),
                            "value" => Ok(GeneratedField::Value),
                            "commitTs" | "commit_ts" => Ok(GeneratedField::CommitTs),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = KvPair;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.KvPair")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<KvPair, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut error__ = None;
                let mut key__ = None;
                let mut value__ = None;
                let mut commit_ts__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Error => {
                            if error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("error"));
                            }
                            error__ = map_.next_value()?;
                        }
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Value => {
                            if value__.is_some() {
                                return Err(serde::de::Error::duplicate_field("value"));
                            }
                            value__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::CommitTs => {
                            if commit_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("commitTs"));
                            }
                            commit_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(KvPair {
                    error: error__,
                    key: key__.unwrap_or_default(),
                    value: value__.unwrap_or_default(),
                    commit_ts: commit_ts__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.KvPair", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for LockInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.primary_lock.is_empty() {
            len += 1;
        }
        if self.lock_version != 0 {
            len += 1;
        }
        if !self.key.is_empty() {
            len += 1;
        }
        if self.lock_ttl != 0 {
            len += 1;
        }
        if self.txn_size != 0 {
            len += 1;
        }
        if self.lock_type != 0 {
            len += 1;
        }
        if self.lock_for_update_ts != 0 {
            len += 1;
        }
        if self.use_async_commit {
            len += 1;
        }
        if self.min_commit_ts != 0 {
            len += 1;
        }
        if !self.secondaries.is_empty() {
            len += 1;
        }
        if self.duration_to_last_update_ms != 0 {
            len += 1;
        }
        if self.is_txn_file {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.LockInfo", len)?;
        if !self.primary_lock.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "primaryLock",
                pbjson::private::base64::encode(&self.primary_lock).as_str(),
            )?;
        }
        if self.lock_version != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "lockVersion",
                ToString::to_string(&self.lock_version).as_str(),
            )?;
        }
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        if self.lock_ttl != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("lockTtl", ToString::to_string(&self.lock_ttl).as_str())?;
        }
        if self.txn_size != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("txnSize", ToString::to_string(&self.txn_size).as_str())?;
        }
        if self.lock_type != 0 {
            let v = Op::try_from(self.lock_type).map_err(|_| {
                serde::ser::Error::custom(format!("Invalid variant {}", self.lock_type))
            })?;
            struct_ser.serialize_field("lockType", &v)?;
        }
        if self.lock_for_update_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "lockForUpdateTs",
                ToString::to_string(&self.lock_for_update_ts).as_str(),
            )?;
        }
        if self.use_async_commit {
            struct_ser.serialize_field("useAsyncCommit", &self.use_async_commit)?;
        }
        if self.min_commit_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "minCommitTs",
                ToString::to_string(&self.min_commit_ts).as_str(),
            )?;
        }
        if !self.secondaries.is_empty() {
            struct_ser.serialize_field(
                "secondaries",
                &self
                    .secondaries
                    .iter()
                    .map(pbjson::private::base64::encode)
                    .collect::<Vec<_>>(),
            )?;
        }
        if self.duration_to_last_update_ms != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "durationToLastUpdateMs",
                ToString::to_string(&self.duration_to_last_update_ms).as_str(),
            )?;
        }
        if self.is_txn_file {
            struct_ser.serialize_field("isTxnFile", &self.is_txn_file)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for LockInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "primary_lock",
            "primaryLock",
            "lock_version",
            "lockVersion",
            "key",
            "lock_ttl",
            "lockTtl",
            "txn_size",
            "txnSize",
            "lock_type",
            "lockType",
            "lock_for_update_ts",
            "lockForUpdateTs",
            "use_async_commit",
            "useAsyncCommit",
            "min_commit_ts",
            "minCommitTs",
            "secondaries",
            "duration_to_last_update_ms",
            "durationToLastUpdateMs",
            "is_txn_file",
            "isTxnFile",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PrimaryLock,
            LockVersion,
            Key,
            LockTtl,
            TxnSize,
            LockType,
            LockForUpdateTs,
            UseAsyncCommit,
            MinCommitTs,
            Secondaries,
            DurationToLastUpdateMs,
            IsTxnFile,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "primaryLock" | "primary_lock" => Ok(GeneratedField::PrimaryLock),
                            "lockVersion" | "lock_version" => Ok(GeneratedField::LockVersion),
                            "key" => Ok(GeneratedField::Key),
                            "lockTtl" | "lock_ttl" => Ok(GeneratedField::LockTtl),
                            "txnSize" | "txn_size" => Ok(GeneratedField::TxnSize),
                            "lockType" | "lock_type" => Ok(GeneratedField::LockType),
                            "lockForUpdateTs" | "lock_for_update_ts" => {
                                Ok(GeneratedField::LockForUpdateTs)
                            }
                            "useAsyncCommit" | "use_async_commit" => {
                                Ok(GeneratedField::UseAsyncCommit)
                            }
                            "minCommitTs" | "min_commit_ts" => Ok(GeneratedField::MinCommitTs),
                            "secondaries" => Ok(GeneratedField::Secondaries),
                            "durationToLastUpdateMs" | "duration_to_last_update_ms" => {
                                Ok(GeneratedField::DurationToLastUpdateMs)
                            }
                            "isTxnFile" | "is_txn_file" => Ok(GeneratedField::IsTxnFile),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = LockInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.LockInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<LockInfo, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut primary_lock__ = None;
                let mut lock_version__ = None;
                let mut key__ = None;
                let mut lock_ttl__ = None;
                let mut txn_size__ = None;
                let mut lock_type__ = None;
                let mut lock_for_update_ts__ = None;
                let mut use_async_commit__ = None;
                let mut min_commit_ts__ = None;
                let mut secondaries__ = None;
                let mut duration_to_last_update_ms__ = None;
                let mut is_txn_file__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::PrimaryLock => {
                            if primary_lock__.is_some() {
                                return Err(serde::de::Error::duplicate_field("primaryLock"));
                            }
                            primary_lock__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::LockVersion => {
                            if lock_version__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lockVersion"));
                            }
                            lock_version__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::LockTtl => {
                            if lock_ttl__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lockTtl"));
                            }
                            lock_ttl__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::TxnSize => {
                            if txn_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("txnSize"));
                            }
                            txn_size__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::LockType => {
                            if lock_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lockType"));
                            }
                            lock_type__ = Some(map_.next_value::<Op>()? as i32);
                        }
                        GeneratedField::LockForUpdateTs => {
                            if lock_for_update_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lockForUpdateTs"));
                            }
                            lock_for_update_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::UseAsyncCommit => {
                            if use_async_commit__.is_some() {
                                return Err(serde::de::Error::duplicate_field("useAsyncCommit"));
                            }
                            use_async_commit__ = Some(map_.next_value()?);
                        }
                        GeneratedField::MinCommitTs => {
                            if min_commit_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("minCommitTs"));
                            }
                            min_commit_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Secondaries => {
                            if secondaries__.is_some() {
                                return Err(serde::de::Error::duplicate_field("secondaries"));
                            }
                            secondaries__ = Some(
                                map_.next_value::<Vec<::pbjson::private::BytesDeserialize<_>>>()?
                                    .into_iter()
                                    .map(|x| x.0)
                                    .collect(),
                            );
                        }
                        GeneratedField::DurationToLastUpdateMs => {
                            if duration_to_last_update_ms__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "durationToLastUpdateMs",
                                ));
                            }
                            duration_to_last_update_ms__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::IsTxnFile => {
                            if is_txn_file__.is_some() {
                                return Err(serde::de::Error::duplicate_field("isTxnFile"));
                            }
                            is_txn_file__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(LockInfo {
                    primary_lock: primary_lock__.unwrap_or_default(),
                    lock_version: lock_version__.unwrap_or_default(),
                    key: key__.unwrap_or_default(),
                    lock_ttl: lock_ttl__.unwrap_or_default(),
                    txn_size: txn_size__.unwrap_or_default(),
                    lock_type: lock_type__.unwrap_or_default(),
                    lock_for_update_ts: lock_for_update_ts__.unwrap_or_default(),
                    use_async_commit: use_async_commit__.unwrap_or_default(),
                    min_commit_ts: min_commit_ts__.unwrap_or_default(),
                    secondaries: secondaries__.unwrap_or_default(),
                    duration_to_last_update_ms: duration_to_last_update_ms__.unwrap_or_default(),
                    is_txn_file: is_txn_file__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.LockInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for MaxTimestampNotSynced {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("tikv.MaxTimestampNotSynced", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for MaxTimestampNotSynced {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {}
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = MaxTimestampNotSynced;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.MaxTimestampNotSynced")
            }

            fn visit_map<V>(
                self,
                mut map_: V,
            ) -> std::result::Result<MaxTimestampNotSynced, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(MaxTimestampNotSynced {})
            }
        }
        deserializer.deserialize_struct("tikv.MaxTimestampNotSynced", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for MismatchPeerId {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.request_peer_id != 0 {
            len += 1;
        }
        if self.store_peer_id != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.MismatchPeerId", len)?;
        if self.request_peer_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "requestPeerId",
                ToString::to_string(&self.request_peer_id).as_str(),
            )?;
        }
        if self.store_peer_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "storePeerId",
                ToString::to_string(&self.store_peer_id).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for MismatchPeerId {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "request_peer_id",
            "requestPeerId",
            "store_peer_id",
            "storePeerId",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RequestPeerId,
            StorePeerId,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "requestPeerId" | "request_peer_id" => {
                                Ok(GeneratedField::RequestPeerId)
                            }
                            "storePeerId" | "store_peer_id" => Ok(GeneratedField::StorePeerId),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = MismatchPeerId;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.MismatchPeerId")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<MismatchPeerId, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut request_peer_id__ = None;
                let mut store_peer_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RequestPeerId => {
                            if request_peer_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requestPeerId"));
                            }
                            request_peer_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::StorePeerId => {
                            if store_peer_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("storePeerId"));
                            }
                            store_peer_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(MismatchPeerId {
                    request_peer_id: request_peer_id__.unwrap_or_default(),
                    store_peer_id: store_peer_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.MismatchPeerId", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for MvccDebugInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.key.is_empty() {
            len += 1;
        }
        if self.mvcc.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.MvccDebugInfo", len)?;
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        if let Some(v) = self.mvcc.as_ref() {
            struct_ser.serialize_field("mvcc", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for MvccDebugInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["key", "mvcc"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Key,
            Mvcc,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "key" => Ok(GeneratedField::Key),
                            "mvcc" => Ok(GeneratedField::Mvcc),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = MvccDebugInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.MvccDebugInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<MvccDebugInfo, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut key__ = None;
                let mut mvcc__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Mvcc => {
                            if mvcc__.is_some() {
                                return Err(serde::de::Error::duplicate_field("mvcc"));
                            }
                            mvcc__ = map_.next_value()?;
                        }
                    }
                }
                Ok(MvccDebugInfo {
                    key: key__.unwrap_or_default(),
                    mvcc: mvcc__,
                })
            }
        }
        deserializer.deserialize_struct("tikv.MvccDebugInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for MvccInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.lock.is_some() {
            len += 1;
        }
        if !self.writes.is_empty() {
            len += 1;
        }
        if !self.values.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.MvccInfo", len)?;
        if let Some(v) = self.lock.as_ref() {
            struct_ser.serialize_field("lock", v)?;
        }
        if !self.writes.is_empty() {
            struct_ser.serialize_field("writes", &self.writes)?;
        }
        if !self.values.is_empty() {
            struct_ser.serialize_field("values", &self.values)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for MvccInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["lock", "writes", "values"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Lock,
            Writes,
            Values,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "lock" => Ok(GeneratedField::Lock),
                            "writes" => Ok(GeneratedField::Writes),
                            "values" => Ok(GeneratedField::Values),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = MvccInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.MvccInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<MvccInfo, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut lock__ = None;
                let mut writes__ = None;
                let mut values__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Lock => {
                            if lock__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lock"));
                            }
                            lock__ = map_.next_value()?;
                        }
                        GeneratedField::Writes => {
                            if writes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("writes"));
                            }
                            writes__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Values => {
                            if values__.is_some() {
                                return Err(serde::de::Error::duplicate_field("values"));
                            }
                            values__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(MvccInfo {
                    lock: lock__,
                    writes: writes__.unwrap_or_default(),
                    values: values__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.MvccInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for MvccLock {
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
        if self.start_ts != 0 {
            len += 1;
        }
        if !self.primary.is_empty() {
            len += 1;
        }
        if !self.short_value.is_empty() {
            len += 1;
        }
        if self.ttl != 0 {
            len += 1;
        }
        if self.for_update_ts != 0 {
            len += 1;
        }
        if self.txn_size != 0 {
            len += 1;
        }
        if self.use_async_commit {
            len += 1;
        }
        if !self.secondaries.is_empty() {
            len += 1;
        }
        if !self.rollback_ts.is_empty() {
            len += 1;
        }
        if self.last_change_ts != 0 {
            len += 1;
        }
        if self.versions_to_last_change != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.MvccLock", len)?;
        if self.r#type != 0 {
            let v = Op::try_from(self.r#type).map_err(|_| {
                serde::ser::Error::custom(format!("Invalid variant {}", self.r#type))
            })?;
            struct_ser.serialize_field("type", &v)?;
        }
        if self.start_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("startTs", ToString::to_string(&self.start_ts).as_str())?;
        }
        if !self.primary.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "primary",
                pbjson::private::base64::encode(&self.primary).as_str(),
            )?;
        }
        if !self.short_value.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "shortValue",
                pbjson::private::base64::encode(&self.short_value).as_str(),
            )?;
        }
        if self.ttl != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("ttl", ToString::to_string(&self.ttl).as_str())?;
        }
        if self.for_update_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "forUpdateTs",
                ToString::to_string(&self.for_update_ts).as_str(),
            )?;
        }
        if self.txn_size != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("txnSize", ToString::to_string(&self.txn_size).as_str())?;
        }
        if self.use_async_commit {
            struct_ser.serialize_field("useAsyncCommit", &self.use_async_commit)?;
        }
        if !self.secondaries.is_empty() {
            struct_ser.serialize_field(
                "secondaries",
                &self
                    .secondaries
                    .iter()
                    .map(pbjson::private::base64::encode)
                    .collect::<Vec<_>>(),
            )?;
        }
        if !self.rollback_ts.is_empty() {
            struct_ser.serialize_field(
                "rollbackTs",
                &self
                    .rollback_ts
                    .iter()
                    .map(ToString::to_string)
                    .collect::<Vec<_>>(),
            )?;
        }
        if self.last_change_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "lastChangeTs",
                ToString::to_string(&self.last_change_ts).as_str(),
            )?;
        }
        if self.versions_to_last_change != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "versionsToLastChange",
                ToString::to_string(&self.versions_to_last_change).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for MvccLock {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "type",
            "start_ts",
            "startTs",
            "primary",
            "short_value",
            "shortValue",
            "ttl",
            "for_update_ts",
            "forUpdateTs",
            "txn_size",
            "txnSize",
            "use_async_commit",
            "useAsyncCommit",
            "secondaries",
            "rollback_ts",
            "rollbackTs",
            "last_change_ts",
            "lastChangeTs",
            "versions_to_last_change",
            "versionsToLastChange",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Type,
            StartTs,
            Primary,
            ShortValue,
            Ttl,
            ForUpdateTs,
            TxnSize,
            UseAsyncCommit,
            Secondaries,
            RollbackTs,
            LastChangeTs,
            VersionsToLastChange,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "type" => Ok(GeneratedField::Type),
                            "startTs" | "start_ts" => Ok(GeneratedField::StartTs),
                            "primary" => Ok(GeneratedField::Primary),
                            "shortValue" | "short_value" => Ok(GeneratedField::ShortValue),
                            "ttl" => Ok(GeneratedField::Ttl),
                            "forUpdateTs" | "for_update_ts" => Ok(GeneratedField::ForUpdateTs),
                            "txnSize" | "txn_size" => Ok(GeneratedField::TxnSize),
                            "useAsyncCommit" | "use_async_commit" => {
                                Ok(GeneratedField::UseAsyncCommit)
                            }
                            "secondaries" => Ok(GeneratedField::Secondaries),
                            "rollbackTs" | "rollback_ts" => Ok(GeneratedField::RollbackTs),
                            "lastChangeTs" | "last_change_ts" => Ok(GeneratedField::LastChangeTs),
                            "versionsToLastChange" | "versions_to_last_change" => {
                                Ok(GeneratedField::VersionsToLastChange)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = MvccLock;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.MvccLock")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<MvccLock, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut r#type__ = None;
                let mut start_ts__ = None;
                let mut primary__ = None;
                let mut short_value__ = None;
                let mut ttl__ = None;
                let mut for_update_ts__ = None;
                let mut txn_size__ = None;
                let mut use_async_commit__ = None;
                let mut secondaries__ = None;
                let mut rollback_ts__ = None;
                let mut last_change_ts__ = None;
                let mut versions_to_last_change__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Type => {
                            if r#type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("type"));
                            }
                            r#type__ = Some(map_.next_value::<Op>()? as i32);
                        }
                        GeneratedField::StartTs => {
                            if start_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startTs"));
                            }
                            start_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Primary => {
                            if primary__.is_some() {
                                return Err(serde::de::Error::duplicate_field("primary"));
                            }
                            primary__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ShortValue => {
                            if short_value__.is_some() {
                                return Err(serde::de::Error::duplicate_field("shortValue"));
                            }
                            short_value__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Ttl => {
                            if ttl__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ttl"));
                            }
                            ttl__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ForUpdateTs => {
                            if for_update_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("forUpdateTs"));
                            }
                            for_update_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::TxnSize => {
                            if txn_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("txnSize"));
                            }
                            txn_size__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::UseAsyncCommit => {
                            if use_async_commit__.is_some() {
                                return Err(serde::de::Error::duplicate_field("useAsyncCommit"));
                            }
                            use_async_commit__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Secondaries => {
                            if secondaries__.is_some() {
                                return Err(serde::de::Error::duplicate_field("secondaries"));
                            }
                            secondaries__ = Some(
                                map_.next_value::<Vec<::pbjson::private::BytesDeserialize<_>>>()?
                                    .into_iter()
                                    .map(|x| x.0)
                                    .collect(),
                            );
                        }
                        GeneratedField::RollbackTs => {
                            if rollback_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("rollbackTs"));
                            }
                            rollback_ts__ = Some(
                                map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter()
                                    .map(|x| x.0)
                                    .collect(),
                            );
                        }
                        GeneratedField::LastChangeTs => {
                            if last_change_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastChangeTs"));
                            }
                            last_change_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::VersionsToLastChange => {
                            if versions_to_last_change__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "versionsToLastChange",
                                ));
                            }
                            versions_to_last_change__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(MvccLock {
                    r#type: r#type__.unwrap_or_default(),
                    start_ts: start_ts__.unwrap_or_default(),
                    primary: primary__.unwrap_or_default(),
                    short_value: short_value__.unwrap_or_default(),
                    ttl: ttl__.unwrap_or_default(),
                    for_update_ts: for_update_ts__.unwrap_or_default(),
                    txn_size: txn_size__.unwrap_or_default(),
                    use_async_commit: use_async_commit__.unwrap_or_default(),
                    secondaries: secondaries__.unwrap_or_default(),
                    rollback_ts: rollback_ts__.unwrap_or_default(),
                    last_change_ts: last_change_ts__.unwrap_or_default(),
                    versions_to_last_change: versions_to_last_change__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.MvccLock", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for MvccValue {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.start_ts != 0 {
            len += 1;
        }
        if !self.value.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.MvccValue", len)?;
        if self.start_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("startTs", ToString::to_string(&self.start_ts).as_str())?;
        }
        if !self.value.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "value",
                pbjson::private::base64::encode(&self.value).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for MvccValue {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["start_ts", "startTs", "value"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            StartTs,
            Value,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "startTs" | "start_ts" => Ok(GeneratedField::StartTs),
                            "value" => Ok(GeneratedField::Value),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = MvccValue;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.MvccValue")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<MvccValue, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut start_ts__ = None;
                let mut value__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::StartTs => {
                            if start_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startTs"));
                            }
                            start_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Value => {
                            if value__.is_some() {
                                return Err(serde::de::Error::duplicate_field("value"));
                            }
                            value__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(MvccValue {
                    start_ts: start_ts__.unwrap_or_default(),
                    value: value__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.MvccValue", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for MvccWrite {
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
        if self.start_ts != 0 {
            len += 1;
        }
        if self.commit_ts != 0 {
            len += 1;
        }
        if !self.short_value.is_empty() {
            len += 1;
        }
        if self.has_overlapped_rollback {
            len += 1;
        }
        if self.has_gc_fence {
            len += 1;
        }
        if self.gc_fence != 0 {
            len += 1;
        }
        if self.last_change_ts != 0 {
            len += 1;
        }
        if self.versions_to_last_change != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.MvccWrite", len)?;
        if self.r#type != 0 {
            let v = Op::try_from(self.r#type).map_err(|_| {
                serde::ser::Error::custom(format!("Invalid variant {}", self.r#type))
            })?;
            struct_ser.serialize_field("type", &v)?;
        }
        if self.start_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("startTs", ToString::to_string(&self.start_ts).as_str())?;
        }
        if self.commit_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("commitTs", ToString::to_string(&self.commit_ts).as_str())?;
        }
        if !self.short_value.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "shortValue",
                pbjson::private::base64::encode(&self.short_value).as_str(),
            )?;
        }
        if self.has_overlapped_rollback {
            struct_ser.serialize_field("hasOverlappedRollback", &self.has_overlapped_rollback)?;
        }
        if self.has_gc_fence {
            struct_ser.serialize_field("hasGcFence", &self.has_gc_fence)?;
        }
        if self.gc_fence != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("gcFence", ToString::to_string(&self.gc_fence).as_str())?;
        }
        if self.last_change_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "lastChangeTs",
                ToString::to_string(&self.last_change_ts).as_str(),
            )?;
        }
        if self.versions_to_last_change != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "versionsToLastChange",
                ToString::to_string(&self.versions_to_last_change).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for MvccWrite {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "type",
            "start_ts",
            "startTs",
            "commit_ts",
            "commitTs",
            "short_value",
            "shortValue",
            "has_overlapped_rollback",
            "hasOverlappedRollback",
            "has_gc_fence",
            "hasGcFence",
            "gc_fence",
            "gcFence",
            "last_change_ts",
            "lastChangeTs",
            "versions_to_last_change",
            "versionsToLastChange",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Type,
            StartTs,
            CommitTs,
            ShortValue,
            HasOverlappedRollback,
            HasGcFence,
            GcFence,
            LastChangeTs,
            VersionsToLastChange,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "type" => Ok(GeneratedField::Type),
                            "startTs" | "start_ts" => Ok(GeneratedField::StartTs),
                            "commitTs" | "commit_ts" => Ok(GeneratedField::CommitTs),
                            "shortValue" | "short_value" => Ok(GeneratedField::ShortValue),
                            "hasOverlappedRollback" | "has_overlapped_rollback" => {
                                Ok(GeneratedField::HasOverlappedRollback)
                            }
                            "hasGcFence" | "has_gc_fence" => Ok(GeneratedField::HasGcFence),
                            "gcFence" | "gc_fence" => Ok(GeneratedField::GcFence),
                            "lastChangeTs" | "last_change_ts" => Ok(GeneratedField::LastChangeTs),
                            "versionsToLastChange" | "versions_to_last_change" => {
                                Ok(GeneratedField::VersionsToLastChange)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = MvccWrite;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.MvccWrite")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<MvccWrite, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut r#type__ = None;
                let mut start_ts__ = None;
                let mut commit_ts__ = None;
                let mut short_value__ = None;
                let mut has_overlapped_rollback__ = None;
                let mut has_gc_fence__ = None;
                let mut gc_fence__ = None;
                let mut last_change_ts__ = None;
                let mut versions_to_last_change__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Type => {
                            if r#type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("type"));
                            }
                            r#type__ = Some(map_.next_value::<Op>()? as i32);
                        }
                        GeneratedField::StartTs => {
                            if start_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startTs"));
                            }
                            start_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::CommitTs => {
                            if commit_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("commitTs"));
                            }
                            commit_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ShortValue => {
                            if short_value__.is_some() {
                                return Err(serde::de::Error::duplicate_field("shortValue"));
                            }
                            short_value__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::HasOverlappedRollback => {
                            if has_overlapped_rollback__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "hasOverlappedRollback",
                                ));
                            }
                            has_overlapped_rollback__ = Some(map_.next_value()?);
                        }
                        GeneratedField::HasGcFence => {
                            if has_gc_fence__.is_some() {
                                return Err(serde::de::Error::duplicate_field("hasGcFence"));
                            }
                            has_gc_fence__ = Some(map_.next_value()?);
                        }
                        GeneratedField::GcFence => {
                            if gc_fence__.is_some() {
                                return Err(serde::de::Error::duplicate_field("gcFence"));
                            }
                            gc_fence__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::LastChangeTs => {
                            if last_change_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastChangeTs"));
                            }
                            last_change_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::VersionsToLastChange => {
                            if versions_to_last_change__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "versionsToLastChange",
                                ));
                            }
                            versions_to_last_change__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(MvccWrite {
                    r#type: r#type__.unwrap_or_default(),
                    start_ts: start_ts__.unwrap_or_default(),
                    commit_ts: commit_ts__.unwrap_or_default(),
                    short_value: short_value__.unwrap_or_default(),
                    has_overlapped_rollback: has_overlapped_rollback__.unwrap_or_default(),
                    has_gc_fence: has_gc_fence__.unwrap_or_default(),
                    gc_fence: gc_fence__.unwrap_or_default(),
                    last_change_ts: last_change_ts__.unwrap_or_default(),
                    versions_to_last_change: versions_to_last_change__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.MvccWrite", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for NotLeader {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_id != 0 {
            len += 1;
        }
        if self.leader.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.NotLeader", len)?;
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        if let Some(v) = self.leader.as_ref() {
            struct_ser.serialize_field("leader", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for NotLeader {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_id", "regionId", "leader"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionId,
            Leader,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            "leader" => Ok(GeneratedField::Leader),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = NotLeader;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.NotLeader")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<NotLeader, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_id__ = None;
                let mut leader__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Leader => {
                            if leader__.is_some() {
                                return Err(serde::de::Error::duplicate_field("leader"));
                            }
                            leader__ = map_.next_value()?;
                        }
                    }
                }
                Ok(NotLeader {
                    region_id: region_id__.unwrap_or_default(),
                    leader: leader__,
                })
            }
        }
        deserializer.deserialize_struct("tikv.NotLeader", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Op {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Put => "Put",
            Self::Del => "Del",
            Self::Lock => "Lock",
            Self::Rollback => "Rollback",
            Self::Insert => "Insert",
            Self::PessimisticLock => "PessimisticLock",
            Self::CheckNotExists => "CheckNotExists",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for Op {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "Put",
            "Del",
            "Lock",
            "Rollback",
            "Insert",
            "PessimisticLock",
            "CheckNotExists",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Op;

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
                    "Put" => Ok(Op::Put),
                    "Del" => Ok(Op::Del),
                    "Lock" => Ok(Op::Lock),
                    "Rollback" => Ok(Op::Rollback),
                    "Insert" => Ok(Op::Insert),
                    "PessimisticLock" => Ok(Op::PessimisticLock),
                    "CheckNotExists" => Ok(Op::CheckNotExists),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for PrimaryMismatch {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.lock_info.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.PrimaryMismatch", len)?;
        if let Some(v) = self.lock_info.as_ref() {
            struct_ser.serialize_field("lockInfo", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for PrimaryMismatch {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["lock_info", "lockInfo"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            LockInfo,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "lockInfo" | "lock_info" => Ok(GeneratedField::LockInfo),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = PrimaryMismatch;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.PrimaryMismatch")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<PrimaryMismatch, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut lock_info__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::LockInfo => {
                            if lock_info__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lockInfo"));
                            }
                            lock_info__ = map_.next_value()?;
                        }
                    }
                }
                Ok(PrimaryMismatch {
                    lock_info: lock_info__,
                })
            }
        }
        deserializer.deserialize_struct("tikv.PrimaryMismatch", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ProposalInMergingMode {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_id != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.ProposalInMergingMode", len)?;
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ProposalInMergingMode {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_id", "regionId"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionId,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ProposalInMergingMode;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.ProposalInMergingMode")
            }

            fn visit_map<V>(
                self,
                mut map_: V,
            ) -> std::result::Result<ProposalInMergingMode, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(ProposalInMergingMode {
                    region_id: region_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.ProposalInMergingMode", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RaftEntryTooLarge {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_id != 0 {
            len += 1;
        }
        if self.entry_size != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RaftEntryTooLarge", len)?;
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        if self.entry_size != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("entrySize", ToString::to_string(&self.entry_size).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RaftEntryTooLarge {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_id", "regionId", "entry_size", "entrySize"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionId,
            EntrySize,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            "entrySize" | "entry_size" => Ok(GeneratedField::EntrySize),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RaftEntryTooLarge;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RaftEntryTooLarge")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RaftEntryTooLarge, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_id__ = None;
                let mut entry_size__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::EntrySize => {
                            if entry_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("entrySize"));
                            }
                            entry_size__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(RaftEntryTooLarge {
                    region_id: region_id__.unwrap_or_default(),
                    entry_size: entry_size__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RaftEntryTooLarge", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawBatchDeleteRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.context.is_some() {
            len += 1;
        }
        if !self.keys.is_empty() {
            len += 1;
        }
        if !self.cf.is_empty() {
            len += 1;
        }
        if self.for_cas {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawBatchDeleteRequest", len)?;
        if let Some(v) = self.context.as_ref() {
            struct_ser.serialize_field("context", v)?;
        }
        if !self.keys.is_empty() {
            struct_ser.serialize_field(
                "keys",
                &self
                    .keys
                    .iter()
                    .map(pbjson::private::base64::encode)
                    .collect::<Vec<_>>(),
            )?;
        }
        if !self.cf.is_empty() {
            struct_ser.serialize_field("cf", &self.cf)?;
        }
        if self.for_cas {
            struct_ser.serialize_field("forCas", &self.for_cas)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawBatchDeleteRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["context", "keys", "cf", "for_cas", "forCas"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Context,
            Keys,
            Cf,
            ForCas,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "context" => Ok(GeneratedField::Context),
                            "keys" => Ok(GeneratedField::Keys),
                            "cf" => Ok(GeneratedField::Cf),
                            "forCas" | "for_cas" => Ok(GeneratedField::ForCas),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawBatchDeleteRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawBatchDeleteRequest")
            }

            fn visit_map<V>(
                self,
                mut map_: V,
            ) -> std::result::Result<RawBatchDeleteRequest, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut context__ = None;
                let mut keys__ = None;
                let mut cf__ = None;
                let mut for_cas__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Context => {
                            if context__.is_some() {
                                return Err(serde::de::Error::duplicate_field("context"));
                            }
                            context__ = map_.next_value()?;
                        }
                        GeneratedField::Keys => {
                            if keys__.is_some() {
                                return Err(serde::de::Error::duplicate_field("keys"));
                            }
                            keys__ = Some(
                                map_.next_value::<Vec<::pbjson::private::BytesDeserialize<_>>>()?
                                    .into_iter()
                                    .map(|x| x.0)
                                    .collect(),
                            );
                        }
                        GeneratedField::Cf => {
                            if cf__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cf"));
                            }
                            cf__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ForCas => {
                            if for_cas__.is_some() {
                                return Err(serde::de::Error::duplicate_field("forCas"));
                            }
                            for_cas__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RawBatchDeleteRequest {
                    context: context__,
                    keys: keys__.unwrap_or_default(),
                    cf: cf__.unwrap_or_default(),
                    for_cas: for_cas__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawBatchDeleteRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawBatchDeleteResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_error.is_some() {
            len += 1;
        }
        if !self.error.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawBatchDeleteResponse", len)?;
        if let Some(v) = self.region_error.as_ref() {
            struct_ser.serialize_field("regionError", v)?;
        }
        if !self.error.is_empty() {
            struct_ser.serialize_field("error", &self.error)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawBatchDeleteResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_error", "regionError", "error"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionError,
            Error,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionError" | "region_error" => Ok(GeneratedField::RegionError),
                            "error" => Ok(GeneratedField::Error),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawBatchDeleteResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawBatchDeleteResponse")
            }

            fn visit_map<V>(
                self,
                mut map_: V,
            ) -> std::result::Result<RawBatchDeleteResponse, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_error__ = None;
                let mut error__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionError => {
                            if region_error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionError"));
                            }
                            region_error__ = map_.next_value()?;
                        }
                        GeneratedField::Error => {
                            if error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("error"));
                            }
                            error__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RawBatchDeleteResponse {
                    region_error: region_error__,
                    error: error__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawBatchDeleteResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawBatchGetRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.context.is_some() {
            len += 1;
        }
        if !self.keys.is_empty() {
            len += 1;
        }
        if !self.cf.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawBatchGetRequest", len)?;
        if let Some(v) = self.context.as_ref() {
            struct_ser.serialize_field("context", v)?;
        }
        if !self.keys.is_empty() {
            struct_ser.serialize_field(
                "keys",
                &self
                    .keys
                    .iter()
                    .map(pbjson::private::base64::encode)
                    .collect::<Vec<_>>(),
            )?;
        }
        if !self.cf.is_empty() {
            struct_ser.serialize_field("cf", &self.cf)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawBatchGetRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["context", "keys", "cf"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Context,
            Keys,
            Cf,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "context" => Ok(GeneratedField::Context),
                            "keys" => Ok(GeneratedField::Keys),
                            "cf" => Ok(GeneratedField::Cf),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawBatchGetRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawBatchGetRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RawBatchGetRequest, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut context__ = None;
                let mut keys__ = None;
                let mut cf__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Context => {
                            if context__.is_some() {
                                return Err(serde::de::Error::duplicate_field("context"));
                            }
                            context__ = map_.next_value()?;
                        }
                        GeneratedField::Keys => {
                            if keys__.is_some() {
                                return Err(serde::de::Error::duplicate_field("keys"));
                            }
                            keys__ = Some(
                                map_.next_value::<Vec<::pbjson::private::BytesDeserialize<_>>>()?
                                    .into_iter()
                                    .map(|x| x.0)
                                    .collect(),
                            );
                        }
                        GeneratedField::Cf => {
                            if cf__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cf"));
                            }
                            cf__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RawBatchGetRequest {
                    context: context__,
                    keys: keys__.unwrap_or_default(),
                    cf: cf__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawBatchGetRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawBatchGetResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_error.is_some() {
            len += 1;
        }
        if !self.pairs.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawBatchGetResponse", len)?;
        if let Some(v) = self.region_error.as_ref() {
            struct_ser.serialize_field("regionError", v)?;
        }
        if !self.pairs.is_empty() {
            struct_ser.serialize_field("pairs", &self.pairs)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawBatchGetResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_error", "regionError", "pairs"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionError,
            Pairs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionError" | "region_error" => Ok(GeneratedField::RegionError),
                            "pairs" => Ok(GeneratedField::Pairs),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawBatchGetResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawBatchGetResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RawBatchGetResponse, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_error__ = None;
                let mut pairs__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionError => {
                            if region_error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionError"));
                            }
                            region_error__ = map_.next_value()?;
                        }
                        GeneratedField::Pairs => {
                            if pairs__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pairs"));
                            }
                            pairs__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RawBatchGetResponse {
                    region_error: region_error__,
                    pairs: pairs__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawBatchGetResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawBatchPutRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.context.is_some() {
            len += 1;
        }
        if !self.pairs.is_empty() {
            len += 1;
        }
        if !self.cf.is_empty() {
            len += 1;
        }
        if self.ttl != 0 {
            len += 1;
        }
        if self.for_cas {
            len += 1;
        }
        if !self.ttls.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawBatchPutRequest", len)?;
        if let Some(v) = self.context.as_ref() {
            struct_ser.serialize_field("context", v)?;
        }
        if !self.pairs.is_empty() {
            struct_ser.serialize_field("pairs", &self.pairs)?;
        }
        if !self.cf.is_empty() {
            struct_ser.serialize_field("cf", &self.cf)?;
        }
        if self.ttl != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("ttl", ToString::to_string(&self.ttl).as_str())?;
        }
        if self.for_cas {
            struct_ser.serialize_field("forCas", &self.for_cas)?;
        }
        if !self.ttls.is_empty() {
            struct_ser.serialize_field(
                "ttls",
                &self
                    .ttls
                    .iter()
                    .map(ToString::to_string)
                    .collect::<Vec<_>>(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawBatchPutRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["context", "pairs", "cf", "ttl", "for_cas", "forCas", "ttls"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Context,
            Pairs,
            Cf,
            Ttl,
            ForCas,
            Ttls,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "context" => Ok(GeneratedField::Context),
                            "pairs" => Ok(GeneratedField::Pairs),
                            "cf" => Ok(GeneratedField::Cf),
                            "ttl" => Ok(GeneratedField::Ttl),
                            "forCas" | "for_cas" => Ok(GeneratedField::ForCas),
                            "ttls" => Ok(GeneratedField::Ttls),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawBatchPutRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawBatchPutRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RawBatchPutRequest, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut context__ = None;
                let mut pairs__ = None;
                let mut cf__ = None;
                let mut ttl__ = None;
                let mut for_cas__ = None;
                let mut ttls__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Context => {
                            if context__.is_some() {
                                return Err(serde::de::Error::duplicate_field("context"));
                            }
                            context__ = map_.next_value()?;
                        }
                        GeneratedField::Pairs => {
                            if pairs__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pairs"));
                            }
                            pairs__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Cf => {
                            if cf__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cf"));
                            }
                            cf__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ttl => {
                            if ttl__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ttl"));
                            }
                            ttl__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ForCas => {
                            if for_cas__.is_some() {
                                return Err(serde::de::Error::duplicate_field("forCas"));
                            }
                            for_cas__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ttls => {
                            if ttls__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ttls"));
                            }
                            ttls__ = Some(
                                map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter()
                                    .map(|x| x.0)
                                    .collect(),
                            );
                        }
                    }
                }
                Ok(RawBatchPutRequest {
                    context: context__,
                    pairs: pairs__.unwrap_or_default(),
                    cf: cf__.unwrap_or_default(),
                    ttl: ttl__.unwrap_or_default(),
                    for_cas: for_cas__.unwrap_or_default(),
                    ttls: ttls__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawBatchPutRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawBatchPutResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_error.is_some() {
            len += 1;
        }
        if !self.error.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawBatchPutResponse", len)?;
        if let Some(v) = self.region_error.as_ref() {
            struct_ser.serialize_field("regionError", v)?;
        }
        if !self.error.is_empty() {
            struct_ser.serialize_field("error", &self.error)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawBatchPutResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_error", "regionError", "error"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionError,
            Error,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionError" | "region_error" => Ok(GeneratedField::RegionError),
                            "error" => Ok(GeneratedField::Error),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawBatchPutResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawBatchPutResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RawBatchPutResponse, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_error__ = None;
                let mut error__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionError => {
                            if region_error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionError"));
                            }
                            region_error__ = map_.next_value()?;
                        }
                        GeneratedField::Error => {
                            if error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("error"));
                            }
                            error__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RawBatchPutResponse {
                    region_error: region_error__,
                    error: error__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawBatchPutResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawDeleteRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.context.is_some() {
            len += 1;
        }
        if !self.key.is_empty() {
            len += 1;
        }
        if !self.cf.is_empty() {
            len += 1;
        }
        if self.for_cas {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawDeleteRequest", len)?;
        if let Some(v) = self.context.as_ref() {
            struct_ser.serialize_field("context", v)?;
        }
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        if !self.cf.is_empty() {
            struct_ser.serialize_field("cf", &self.cf)?;
        }
        if self.for_cas {
            struct_ser.serialize_field("forCas", &self.for_cas)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawDeleteRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["context", "key", "cf", "for_cas", "forCas"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Context,
            Key,
            Cf,
            ForCas,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "context" => Ok(GeneratedField::Context),
                            "key" => Ok(GeneratedField::Key),
                            "cf" => Ok(GeneratedField::Cf),
                            "forCas" | "for_cas" => Ok(GeneratedField::ForCas),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawDeleteRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawDeleteRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RawDeleteRequest, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut context__ = None;
                let mut key__ = None;
                let mut cf__ = None;
                let mut for_cas__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Context => {
                            if context__.is_some() {
                                return Err(serde::de::Error::duplicate_field("context"));
                            }
                            context__ = map_.next_value()?;
                        }
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Cf => {
                            if cf__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cf"));
                            }
                            cf__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ForCas => {
                            if for_cas__.is_some() {
                                return Err(serde::de::Error::duplicate_field("forCas"));
                            }
                            for_cas__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RawDeleteRequest {
                    context: context__,
                    key: key__.unwrap_or_default(),
                    cf: cf__.unwrap_or_default(),
                    for_cas: for_cas__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawDeleteRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawDeleteResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_error.is_some() {
            len += 1;
        }
        if !self.error.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawDeleteResponse", len)?;
        if let Some(v) = self.region_error.as_ref() {
            struct_ser.serialize_field("regionError", v)?;
        }
        if !self.error.is_empty() {
            struct_ser.serialize_field("error", &self.error)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawDeleteResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_error", "regionError", "error"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionError,
            Error,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionError" | "region_error" => Ok(GeneratedField::RegionError),
                            "error" => Ok(GeneratedField::Error),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawDeleteResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawDeleteResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RawDeleteResponse, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_error__ = None;
                let mut error__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionError => {
                            if region_error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionError"));
                            }
                            region_error__ = map_.next_value()?;
                        }
                        GeneratedField::Error => {
                            if error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("error"));
                            }
                            error__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RawDeleteResponse {
                    region_error: region_error__,
                    error: error__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawDeleteResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawGetRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.context.is_some() {
            len += 1;
        }
        if !self.key.is_empty() {
            len += 1;
        }
        if !self.cf.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawGetRequest", len)?;
        if let Some(v) = self.context.as_ref() {
            struct_ser.serialize_field("context", v)?;
        }
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        if !self.cf.is_empty() {
            struct_ser.serialize_field("cf", &self.cf)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawGetRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["context", "key", "cf"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Context,
            Key,
            Cf,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "context" => Ok(GeneratedField::Context),
                            "key" => Ok(GeneratedField::Key),
                            "cf" => Ok(GeneratedField::Cf),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawGetRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawGetRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RawGetRequest, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut context__ = None;
                let mut key__ = None;
                let mut cf__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Context => {
                            if context__.is_some() {
                                return Err(serde::de::Error::duplicate_field("context"));
                            }
                            context__ = map_.next_value()?;
                        }
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Cf => {
                            if cf__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cf"));
                            }
                            cf__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RawGetRequest {
                    context: context__,
                    key: key__.unwrap_or_default(),
                    cf: cf__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawGetRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawGetResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_error.is_some() {
            len += 1;
        }
        if !self.error.is_empty() {
            len += 1;
        }
        if !self.value.is_empty() {
            len += 1;
        }
        if self.not_found {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawGetResponse", len)?;
        if let Some(v) = self.region_error.as_ref() {
            struct_ser.serialize_field("regionError", v)?;
        }
        if !self.error.is_empty() {
            struct_ser.serialize_field("error", &self.error)?;
        }
        if !self.value.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "value",
                pbjson::private::base64::encode(&self.value).as_str(),
            )?;
        }
        if self.not_found {
            struct_ser.serialize_field("notFound", &self.not_found)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawGetResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "region_error",
            "regionError",
            "error",
            "value",
            "not_found",
            "notFound",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionError,
            Error,
            Value,
            NotFound,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionError" | "region_error" => Ok(GeneratedField::RegionError),
                            "error" => Ok(GeneratedField::Error),
                            "value" => Ok(GeneratedField::Value),
                            "notFound" | "not_found" => Ok(GeneratedField::NotFound),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawGetResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawGetResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RawGetResponse, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_error__ = None;
                let mut error__ = None;
                let mut value__ = None;
                let mut not_found__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionError => {
                            if region_error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionError"));
                            }
                            region_error__ = map_.next_value()?;
                        }
                        GeneratedField::Error => {
                            if error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("error"));
                            }
                            error__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Value => {
                            if value__.is_some() {
                                return Err(serde::de::Error::duplicate_field("value"));
                            }
                            value__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::NotFound => {
                            if not_found__.is_some() {
                                return Err(serde::de::Error::duplicate_field("notFound"));
                            }
                            not_found__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RawGetResponse {
                    region_error: region_error__,
                    error: error__.unwrap_or_default(),
                    value: value__.unwrap_or_default(),
                    not_found: not_found__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawGetResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawPutRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.context.is_some() {
            len += 1;
        }
        if !self.key.is_empty() {
            len += 1;
        }
        if !self.value.is_empty() {
            len += 1;
        }
        if !self.cf.is_empty() {
            len += 1;
        }
        if self.ttl != 0 {
            len += 1;
        }
        if self.for_cas {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawPutRequest", len)?;
        if let Some(v) = self.context.as_ref() {
            struct_ser.serialize_field("context", v)?;
        }
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        if !self.value.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "value",
                pbjson::private::base64::encode(&self.value).as_str(),
            )?;
        }
        if !self.cf.is_empty() {
            struct_ser.serialize_field("cf", &self.cf)?;
        }
        if self.ttl != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("ttl", ToString::to_string(&self.ttl).as_str())?;
        }
        if self.for_cas {
            struct_ser.serialize_field("forCas", &self.for_cas)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawPutRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["context", "key", "value", "cf", "ttl", "for_cas", "forCas"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Context,
            Key,
            Value,
            Cf,
            Ttl,
            ForCas,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "context" => Ok(GeneratedField::Context),
                            "key" => Ok(GeneratedField::Key),
                            "value" => Ok(GeneratedField::Value),
                            "cf" => Ok(GeneratedField::Cf),
                            "ttl" => Ok(GeneratedField::Ttl),
                            "forCas" | "for_cas" => Ok(GeneratedField::ForCas),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawPutRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawPutRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RawPutRequest, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut context__ = None;
                let mut key__ = None;
                let mut value__ = None;
                let mut cf__ = None;
                let mut ttl__ = None;
                let mut for_cas__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Context => {
                            if context__.is_some() {
                                return Err(serde::de::Error::duplicate_field("context"));
                            }
                            context__ = map_.next_value()?;
                        }
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Value => {
                            if value__.is_some() {
                                return Err(serde::de::Error::duplicate_field("value"));
                            }
                            value__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Cf => {
                            if cf__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cf"));
                            }
                            cf__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ttl => {
                            if ttl__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ttl"));
                            }
                            ttl__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ForCas => {
                            if for_cas__.is_some() {
                                return Err(serde::de::Error::duplicate_field("forCas"));
                            }
                            for_cas__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RawPutRequest {
                    context: context__,
                    key: key__.unwrap_or_default(),
                    value: value__.unwrap_or_default(),
                    cf: cf__.unwrap_or_default(),
                    ttl: ttl__.unwrap_or_default(),
                    for_cas: for_cas__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawPutRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RawPutResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_error.is_some() {
            len += 1;
        }
        if !self.error.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RawPutResponse", len)?;
        if let Some(v) = self.region_error.as_ref() {
            struct_ser.serialize_field("regionError", v)?;
        }
        if !self.error.is_empty() {
            struct_ser.serialize_field("error", &self.error)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RawPutResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_error", "regionError", "error"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionError,
            Error,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionError" | "region_error" => Ok(GeneratedField::RegionError),
                            "error" => Ok(GeneratedField::Error),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RawPutResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RawPutResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RawPutResponse, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_error__ = None;
                let mut error__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionError => {
                            if region_error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionError"));
                            }
                            region_error__ = map_.next_value()?;
                        }
                        GeneratedField::Error => {
                            if error__.is_some() {
                                return Err(serde::de::Error::duplicate_field("error"));
                            }
                            error__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(RawPutResponse {
                    region_error: region_error__,
                    error: error__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RawPutResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ReadIndexNotReady {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.reason.is_empty() {
            len += 1;
        }
        if self.region_id != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.ReadIndexNotReady", len)?;
        if !self.reason.is_empty() {
            struct_ser.serialize_field("reason", &self.reason)?;
        }
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ReadIndexNotReady {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["reason", "region_id", "regionId"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Reason,
            RegionId,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "reason" => Ok(GeneratedField::Reason),
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ReadIndexNotReady;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.ReadIndexNotReady")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ReadIndexNotReady, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut reason__ = None;
                let mut region_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Reason => {
                            if reason__.is_some() {
                                return Err(serde::de::Error::duplicate_field("reason"));
                            }
                            reason__ = Some(map_.next_value()?);
                        }
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(ReadIndexNotReady {
                    reason: reason__.unwrap_or_default(),
                    region_id: region_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.ReadIndexNotReady", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RecoveryInProgress {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_id != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RecoveryInProgress", len)?;
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RecoveryInProgress {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_id", "regionId"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionId,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RecoveryInProgress;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RecoveryInProgress")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RecoveryInProgress, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(RecoveryInProgress {
                    region_id: region_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RecoveryInProgress", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RegionNotFound {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_id != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RegionNotFound", len)?;
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RegionNotFound {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_id", "regionId"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionId,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RegionNotFound;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RegionNotFound")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RegionNotFound, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(RegionNotFound {
                    region_id: region_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RegionNotFound", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RegionNotInitialized {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.region_id != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.RegionNotInitialized", len)?;
        if self.region_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("regionId", ToString::to_string(&self.region_id).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RegionNotInitialized {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["region_id", "regionId"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RegionId,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "regionId" | "region_id" => Ok(GeneratedField::RegionId),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RegionNotInitialized;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.RegionNotInitialized")
            }

            fn visit_map<V>(
                self,
                mut map_: V,
            ) -> std::result::Result<RegionNotInitialized, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut region_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RegionId => {
                            if region_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionId"));
                            }
                            region_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(RegionNotInitialized {
                    region_id: region_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.RegionNotInitialized", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ResourceControlContext {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.resource_group_name.is_empty() {
            len += 1;
        }
        if self.override_priority != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.ResourceControlContext", len)?;
        if !self.resource_group_name.is_empty() {
            struct_ser.serialize_field("resourceGroupName", &self.resource_group_name)?;
        }
        if self.override_priority != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "overridePriority",
                ToString::to_string(&self.override_priority).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ResourceControlContext {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "resource_group_name",
            "resourceGroupName",
            "override_priority",
            "overridePriority",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            ResourceGroupName,
            OverridePriority,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "resourceGroupName" | "resource_group_name" => {
                                Ok(GeneratedField::ResourceGroupName)
                            }
                            "overridePriority" | "override_priority" => {
                                Ok(GeneratedField::OverridePriority)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ResourceControlContext;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.ResourceControlContext")
            }

            fn visit_map<V>(
                self,
                mut map_: V,
            ) -> std::result::Result<ResourceControlContext, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut resource_group_name__ = None;
                let mut override_priority__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::ResourceGroupName => {
                            if resource_group_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("resourceGroupName"));
                            }
                            resource_group_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::OverridePriority => {
                            if override_priority__.is_some() {
                                return Err(serde::de::Error::duplicate_field("overridePriority"));
                            }
                            override_priority__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(ResourceControlContext {
                    resource_group_name: resource_group_name__.unwrap_or_default(),
                    override_priority: override_priority__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.ResourceControlContext", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ScanDetail {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.write.is_some() {
            len += 1;
        }
        if self.lock.is_some() {
            len += 1;
        }
        if self.data.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.ScanDetail", len)?;
        if let Some(v) = self.write.as_ref() {
            struct_ser.serialize_field("write", v)?;
        }
        if let Some(v) = self.lock.as_ref() {
            struct_ser.serialize_field("lock", v)?;
        }
        if let Some(v) = self.data.as_ref() {
            struct_ser.serialize_field("data", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ScanDetail {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["write", "lock", "data"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Write,
            Lock,
            Data,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "write" => Ok(GeneratedField::Write),
                            "lock" => Ok(GeneratedField::Lock),
                            "data" => Ok(GeneratedField::Data),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ScanDetail;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.ScanDetail")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ScanDetail, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut write__ = None;
                let mut lock__ = None;
                let mut data__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Write => {
                            if write__.is_some() {
                                return Err(serde::de::Error::duplicate_field("write"));
                            }
                            write__ = map_.next_value()?;
                        }
                        GeneratedField::Lock => {
                            if lock__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lock"));
                            }
                            lock__ = map_.next_value()?;
                        }
                        GeneratedField::Data => {
                            if data__.is_some() {
                                return Err(serde::de::Error::duplicate_field("data"));
                            }
                            data__ = map_.next_value()?;
                        }
                    }
                }
                Ok(ScanDetail {
                    write: write__,
                    lock: lock__,
                    data: data__,
                })
            }
        }
        deserializer.deserialize_struct("tikv.ScanDetail", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ScanDetailV2 {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.processed_versions != 0 {
            len += 1;
        }
        if self.processed_versions_size != 0 {
            len += 1;
        }
        if self.total_versions != 0 {
            len += 1;
        }
        if self.rocksdb_delete_skipped_count != 0 {
            len += 1;
        }
        if self.rocksdb_key_skipped_count != 0 {
            len += 1;
        }
        if self.rocksdb_block_cache_hit_count != 0 {
            len += 1;
        }
        if self.rocksdb_block_read_count != 0 {
            len += 1;
        }
        if self.rocksdb_block_read_byte != 0 {
            len += 1;
        }
        if self.rocksdb_block_read_nanos != 0 {
            len += 1;
        }
        if self.get_snapshot_nanos != 0 {
            len += 1;
        }
        if self.read_index_propose_wait_nanos != 0 {
            len += 1;
        }
        if self.read_index_confirm_wait_nanos != 0 {
            len += 1;
        }
        if self.read_pool_schedule_wait_nanos != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.ScanDetailV2", len)?;
        if self.processed_versions != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "processedVersions",
                ToString::to_string(&self.processed_versions).as_str(),
            )?;
        }
        if self.processed_versions_size != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "processedVersionsSize",
                ToString::to_string(&self.processed_versions_size).as_str(),
            )?;
        }
        if self.total_versions != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "totalVersions",
                ToString::to_string(&self.total_versions).as_str(),
            )?;
        }
        if self.rocksdb_delete_skipped_count != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "rocksdbDeleteSkippedCount",
                ToString::to_string(&self.rocksdb_delete_skipped_count).as_str(),
            )?;
        }
        if self.rocksdb_key_skipped_count != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "rocksdbKeySkippedCount",
                ToString::to_string(&self.rocksdb_key_skipped_count).as_str(),
            )?;
        }
        if self.rocksdb_block_cache_hit_count != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "rocksdbBlockCacheHitCount",
                ToString::to_string(&self.rocksdb_block_cache_hit_count).as_str(),
            )?;
        }
        if self.rocksdb_block_read_count != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "rocksdbBlockReadCount",
                ToString::to_string(&self.rocksdb_block_read_count).as_str(),
            )?;
        }
        if self.rocksdb_block_read_byte != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "rocksdbBlockReadByte",
                ToString::to_string(&self.rocksdb_block_read_byte).as_str(),
            )?;
        }
        if self.rocksdb_block_read_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "rocksdbBlockReadNanos",
                ToString::to_string(&self.rocksdb_block_read_nanos).as_str(),
            )?;
        }
        if self.get_snapshot_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "getSnapshotNanos",
                ToString::to_string(&self.get_snapshot_nanos).as_str(),
            )?;
        }
        if self.read_index_propose_wait_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "readIndexProposeWaitNanos",
                ToString::to_string(&self.read_index_propose_wait_nanos).as_str(),
            )?;
        }
        if self.read_index_confirm_wait_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "readIndexConfirmWaitNanos",
                ToString::to_string(&self.read_index_confirm_wait_nanos).as_str(),
            )?;
        }
        if self.read_pool_schedule_wait_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "readPoolScheduleWaitNanos",
                ToString::to_string(&self.read_pool_schedule_wait_nanos).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ScanDetailV2 {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "processed_versions",
            "processedVersions",
            "processed_versions_size",
            "processedVersionsSize",
            "total_versions",
            "totalVersions",
            "rocksdb_delete_skipped_count",
            "rocksdbDeleteSkippedCount",
            "rocksdb_key_skipped_count",
            "rocksdbKeySkippedCount",
            "rocksdb_block_cache_hit_count",
            "rocksdbBlockCacheHitCount",
            "rocksdb_block_read_count",
            "rocksdbBlockReadCount",
            "rocksdb_block_read_byte",
            "rocksdbBlockReadByte",
            "rocksdb_block_read_nanos",
            "rocksdbBlockReadNanos",
            "get_snapshot_nanos",
            "getSnapshotNanos",
            "read_index_propose_wait_nanos",
            "readIndexProposeWaitNanos",
            "read_index_confirm_wait_nanos",
            "readIndexConfirmWaitNanos",
            "read_pool_schedule_wait_nanos",
            "readPoolScheduleWaitNanos",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            ProcessedVersions,
            ProcessedVersionsSize,
            TotalVersions,
            RocksdbDeleteSkippedCount,
            RocksdbKeySkippedCount,
            RocksdbBlockCacheHitCount,
            RocksdbBlockReadCount,
            RocksdbBlockReadByte,
            RocksdbBlockReadNanos,
            GetSnapshotNanos,
            ReadIndexProposeWaitNanos,
            ReadIndexConfirmWaitNanos,
            ReadPoolScheduleWaitNanos,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "processedVersions" | "processed_versions" => {
                                Ok(GeneratedField::ProcessedVersions)
                            }
                            "processedVersionsSize" | "processed_versions_size" => {
                                Ok(GeneratedField::ProcessedVersionsSize)
                            }
                            "totalVersions" | "total_versions" => Ok(GeneratedField::TotalVersions),
                            "rocksdbDeleteSkippedCount" | "rocksdb_delete_skipped_count" => {
                                Ok(GeneratedField::RocksdbDeleteSkippedCount)
                            }
                            "rocksdbKeySkippedCount" | "rocksdb_key_skipped_count" => {
                                Ok(GeneratedField::RocksdbKeySkippedCount)
                            }
                            "rocksdbBlockCacheHitCount" | "rocksdb_block_cache_hit_count" => {
                                Ok(GeneratedField::RocksdbBlockCacheHitCount)
                            }
                            "rocksdbBlockReadCount" | "rocksdb_block_read_count" => {
                                Ok(GeneratedField::RocksdbBlockReadCount)
                            }
                            "rocksdbBlockReadByte" | "rocksdb_block_read_byte" => {
                                Ok(GeneratedField::RocksdbBlockReadByte)
                            }
                            "rocksdbBlockReadNanos" | "rocksdb_block_read_nanos" => {
                                Ok(GeneratedField::RocksdbBlockReadNanos)
                            }
                            "getSnapshotNanos" | "get_snapshot_nanos" => {
                                Ok(GeneratedField::GetSnapshotNanos)
                            }
                            "readIndexProposeWaitNanos" | "read_index_propose_wait_nanos" => {
                                Ok(GeneratedField::ReadIndexProposeWaitNanos)
                            }
                            "readIndexConfirmWaitNanos" | "read_index_confirm_wait_nanos" => {
                                Ok(GeneratedField::ReadIndexConfirmWaitNanos)
                            }
                            "readPoolScheduleWaitNanos" | "read_pool_schedule_wait_nanos" => {
                                Ok(GeneratedField::ReadPoolScheduleWaitNanos)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ScanDetailV2;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.ScanDetailV2")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ScanDetailV2, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut processed_versions__ = None;
                let mut processed_versions_size__ = None;
                let mut total_versions__ = None;
                let mut rocksdb_delete_skipped_count__ = None;
                let mut rocksdb_key_skipped_count__ = None;
                let mut rocksdb_block_cache_hit_count__ = None;
                let mut rocksdb_block_read_count__ = None;
                let mut rocksdb_block_read_byte__ = None;
                let mut rocksdb_block_read_nanos__ = None;
                let mut get_snapshot_nanos__ = None;
                let mut read_index_propose_wait_nanos__ = None;
                let mut read_index_confirm_wait_nanos__ = None;
                let mut read_pool_schedule_wait_nanos__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::ProcessedVersions => {
                            if processed_versions__.is_some() {
                                return Err(serde::de::Error::duplicate_field("processedVersions"));
                            }
                            processed_versions__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ProcessedVersionsSize => {
                            if processed_versions_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "processedVersionsSize",
                                ));
                            }
                            processed_versions_size__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::TotalVersions => {
                            if total_versions__.is_some() {
                                return Err(serde::de::Error::duplicate_field("totalVersions"));
                            }
                            total_versions__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::RocksdbDeleteSkippedCount => {
                            if rocksdb_delete_skipped_count__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "rocksdbDeleteSkippedCount",
                                ));
                            }
                            rocksdb_delete_skipped_count__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::RocksdbKeySkippedCount => {
                            if rocksdb_key_skipped_count__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "rocksdbKeySkippedCount",
                                ));
                            }
                            rocksdb_key_skipped_count__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::RocksdbBlockCacheHitCount => {
                            if rocksdb_block_cache_hit_count__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "rocksdbBlockCacheHitCount",
                                ));
                            }
                            rocksdb_block_cache_hit_count__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::RocksdbBlockReadCount => {
                            if rocksdb_block_read_count__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "rocksdbBlockReadCount",
                                ));
                            }
                            rocksdb_block_read_count__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::RocksdbBlockReadByte => {
                            if rocksdb_block_read_byte__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "rocksdbBlockReadByte",
                                ));
                            }
                            rocksdb_block_read_byte__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::RocksdbBlockReadNanos => {
                            if rocksdb_block_read_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "rocksdbBlockReadNanos",
                                ));
                            }
                            rocksdb_block_read_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::GetSnapshotNanos => {
                            if get_snapshot_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field("getSnapshotNanos"));
                            }
                            get_snapshot_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ReadIndexProposeWaitNanos => {
                            if read_index_propose_wait_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "readIndexProposeWaitNanos",
                                ));
                            }
                            read_index_propose_wait_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ReadIndexConfirmWaitNanos => {
                            if read_index_confirm_wait_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "readIndexConfirmWaitNanos",
                                ));
                            }
                            read_index_confirm_wait_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ReadPoolScheduleWaitNanos => {
                            if read_pool_schedule_wait_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "readPoolScheduleWaitNanos",
                                ));
                            }
                            read_pool_schedule_wait_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(ScanDetailV2 {
                    processed_versions: processed_versions__.unwrap_or_default(),
                    processed_versions_size: processed_versions_size__.unwrap_or_default(),
                    total_versions: total_versions__.unwrap_or_default(),
                    rocksdb_delete_skipped_count: rocksdb_delete_skipped_count__
                        .unwrap_or_default(),
                    rocksdb_key_skipped_count: rocksdb_key_skipped_count__.unwrap_or_default(),
                    rocksdb_block_cache_hit_count: rocksdb_block_cache_hit_count__
                        .unwrap_or_default(),
                    rocksdb_block_read_count: rocksdb_block_read_count__.unwrap_or_default(),
                    rocksdb_block_read_byte: rocksdb_block_read_byte__.unwrap_or_default(),
                    rocksdb_block_read_nanos: rocksdb_block_read_nanos__.unwrap_or_default(),
                    get_snapshot_nanos: get_snapshot_nanos__.unwrap_or_default(),
                    read_index_propose_wait_nanos: read_index_propose_wait_nanos__
                        .unwrap_or_default(),
                    read_index_confirm_wait_nanos: read_index_confirm_wait_nanos__
                        .unwrap_or_default(),
                    read_pool_schedule_wait_nanos: read_pool_schedule_wait_nanos__
                        .unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.ScanDetailV2", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ScanInfo {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.total != 0 {
            len += 1;
        }
        if self.processed != 0 {
            len += 1;
        }
        if self.read_bytes != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.ScanInfo", len)?;
        if self.total != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("total", ToString::to_string(&self.total).as_str())?;
        }
        if self.processed != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("processed", ToString::to_string(&self.processed).as_str())?;
        }
        if self.read_bytes != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("readBytes", ToString::to_string(&self.read_bytes).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ScanInfo {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["total", "processed", "read_bytes", "readBytes"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Total,
            Processed,
            ReadBytes,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "total" => Ok(GeneratedField::Total),
                            "processed" => Ok(GeneratedField::Processed),
                            "readBytes" | "read_bytes" => Ok(GeneratedField::ReadBytes),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ScanInfo;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.ScanInfo")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ScanInfo, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut total__ = None;
                let mut processed__ = None;
                let mut read_bytes__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Total => {
                            if total__.is_some() {
                                return Err(serde::de::Error::duplicate_field("total"));
                            }
                            total__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Processed => {
                            if processed__.is_some() {
                                return Err(serde::de::Error::duplicate_field("processed"));
                            }
                            processed__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ReadBytes => {
                            if read_bytes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("readBytes"));
                            }
                            read_bytes__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(ScanInfo {
                    total: total__.unwrap_or_default(),
                    processed: processed__.unwrap_or_default(),
                    read_bytes: read_bytes__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.ScanInfo", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for ServerIsBusy {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.reason.is_empty() {
            len += 1;
        }
        if self.backoff_ms != 0 {
            len += 1;
        }
        if self.estimated_wait_ms != 0 {
            len += 1;
        }
        if self.applied_index != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.ServerIsBusy", len)?;
        if !self.reason.is_empty() {
            struct_ser.serialize_field("reason", &self.reason)?;
        }
        if self.backoff_ms != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("backoffMs", ToString::to_string(&self.backoff_ms).as_str())?;
        }
        if self.estimated_wait_ms != 0 {
            struct_ser.serialize_field("estimatedWaitMs", &self.estimated_wait_ms)?;
        }
        if self.applied_index != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "appliedIndex",
                ToString::to_string(&self.applied_index).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for ServerIsBusy {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "reason",
            "backoff_ms",
            "backoffMs",
            "estimated_wait_ms",
            "estimatedWaitMs",
            "applied_index",
            "appliedIndex",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Reason,
            BackoffMs,
            EstimatedWaitMs,
            AppliedIndex,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "reason" => Ok(GeneratedField::Reason),
                            "backoffMs" | "backoff_ms" => Ok(GeneratedField::BackoffMs),
                            "estimatedWaitMs" | "estimated_wait_ms" => {
                                Ok(GeneratedField::EstimatedWaitMs)
                            }
                            "appliedIndex" | "applied_index" => Ok(GeneratedField::AppliedIndex),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = ServerIsBusy;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.ServerIsBusy")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<ServerIsBusy, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut reason__ = None;
                let mut backoff_ms__ = None;
                let mut estimated_wait_ms__ = None;
                let mut applied_index__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Reason => {
                            if reason__.is_some() {
                                return Err(serde::de::Error::duplicate_field("reason"));
                            }
                            reason__ = Some(map_.next_value()?);
                        }
                        GeneratedField::BackoffMs => {
                            if backoff_ms__.is_some() {
                                return Err(serde::de::Error::duplicate_field("backoffMs"));
                            }
                            backoff_ms__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::EstimatedWaitMs => {
                            if estimated_wait_ms__.is_some() {
                                return Err(serde::de::Error::duplicate_field("estimatedWaitMs"));
                            }
                            estimated_wait_ms__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::AppliedIndex => {
                            if applied_index__.is_some() {
                                return Err(serde::de::Error::duplicate_field("appliedIndex"));
                            }
                            applied_index__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(ServerIsBusy {
                    reason: reason__.unwrap_or_default(),
                    backoff_ms: backoff_ms__.unwrap_or_default(),
                    estimated_wait_ms: estimated_wait_ms__.unwrap_or_default(),
                    applied_index: applied_index__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.ServerIsBusy", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for SourceStmt {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.start_ts != 0 {
            len += 1;
        }
        if self.connection_id != 0 {
            len += 1;
        }
        if self.stmt_id != 0 {
            len += 1;
        }
        if !self.session_alias.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.SourceStmt", len)?;
        if self.start_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("startTs", ToString::to_string(&self.start_ts).as_str())?;
        }
        if self.connection_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "connectionId",
                ToString::to_string(&self.connection_id).as_str(),
            )?;
        }
        if self.stmt_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("stmtId", ToString::to_string(&self.stmt_id).as_str())?;
        }
        if !self.session_alias.is_empty() {
            struct_ser.serialize_field("sessionAlias", &self.session_alias)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for SourceStmt {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "start_ts",
            "startTs",
            "connection_id",
            "connectionId",
            "stmt_id",
            "stmtId",
            "session_alias",
            "sessionAlias",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            StartTs,
            ConnectionId,
            StmtId,
            SessionAlias,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "startTs" | "start_ts" => Ok(GeneratedField::StartTs),
                            "connectionId" | "connection_id" => Ok(GeneratedField::ConnectionId),
                            "stmtId" | "stmt_id" => Ok(GeneratedField::StmtId),
                            "sessionAlias" | "session_alias" => Ok(GeneratedField::SessionAlias),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = SourceStmt;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.SourceStmt")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<SourceStmt, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut start_ts__ = None;
                let mut connection_id__ = None;
                let mut stmt_id__ = None;
                let mut session_alias__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::StartTs => {
                            if start_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startTs"));
                            }
                            start_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ConnectionId => {
                            if connection_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("connectionId"));
                            }
                            connection_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::StmtId => {
                            if stmt_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("stmtId"));
                            }
                            stmt_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::SessionAlias => {
                            if session_alias__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sessionAlias"));
                            }
                            session_alias__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(SourceStmt {
                    start_ts: start_ts__.unwrap_or_default(),
                    connection_id: connection_id__.unwrap_or_default(),
                    stmt_id: stmt_id__.unwrap_or_default(),
                    session_alias: session_alias__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.SourceStmt", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for StaleCommand {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("tikv.StaleCommand", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for StaleCommand {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {}
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = StaleCommand;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.StaleCommand")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<StaleCommand, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(StaleCommand {})
            }
        }
        deserializer.deserialize_struct("tikv.StaleCommand", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for StoreNotMatch {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.request_store_id != 0 {
            len += 1;
        }
        if self.actual_store_id != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.StoreNotMatch", len)?;
        if self.request_store_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "requestStoreId",
                ToString::to_string(&self.request_store_id).as_str(),
            )?;
        }
        if self.actual_store_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "actualStoreId",
                ToString::to_string(&self.actual_store_id).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for StoreNotMatch {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "request_store_id",
            "requestStoreId",
            "actual_store_id",
            "actualStoreId",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RequestStoreId,
            ActualStoreId,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "requestStoreId" | "request_store_id" => {
                                Ok(GeneratedField::RequestStoreId)
                            }
                            "actualStoreId" | "actual_store_id" => {
                                Ok(GeneratedField::ActualStoreId)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = StoreNotMatch;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.StoreNotMatch")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<StoreNotMatch, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut request_store_id__ = None;
                let mut actual_store_id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RequestStoreId => {
                            if request_store_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requestStoreId"));
                            }
                            request_store_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ActualStoreId => {
                            if actual_store_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("actualStoreId"));
                            }
                            actual_store_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(StoreNotMatch {
                    request_store_id: request_store_id__.unwrap_or_default(),
                    actual_store_id: actual_store_id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.StoreNotMatch", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for TimeDetail {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.wait_wall_time_ms != 0 {
            len += 1;
        }
        if self.process_wall_time_ms != 0 {
            len += 1;
        }
        if self.kv_read_wall_time_ms != 0 {
            len += 1;
        }
        if self.total_rpc_wall_time_ns != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.TimeDetail", len)?;
        if self.wait_wall_time_ms != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "waitWallTimeMs",
                ToString::to_string(&self.wait_wall_time_ms).as_str(),
            )?;
        }
        if self.process_wall_time_ms != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "processWallTimeMs",
                ToString::to_string(&self.process_wall_time_ms).as_str(),
            )?;
        }
        if self.kv_read_wall_time_ms != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "kvReadWallTimeMs",
                ToString::to_string(&self.kv_read_wall_time_ms).as_str(),
            )?;
        }
        if self.total_rpc_wall_time_ns != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "totalRpcWallTimeNs",
                ToString::to_string(&self.total_rpc_wall_time_ns).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for TimeDetail {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wait_wall_time_ms",
            "waitWallTimeMs",
            "process_wall_time_ms",
            "processWallTimeMs",
            "kv_read_wall_time_ms",
            "kvReadWallTimeMs",
            "total_rpc_wall_time_ns",
            "totalRpcWallTimeNs",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WaitWallTimeMs,
            ProcessWallTimeMs,
            KvReadWallTimeMs,
            TotalRpcWallTimeNs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "waitWallTimeMs" | "wait_wall_time_ms" => {
                                Ok(GeneratedField::WaitWallTimeMs)
                            }
                            "processWallTimeMs" | "process_wall_time_ms" => {
                                Ok(GeneratedField::ProcessWallTimeMs)
                            }
                            "kvReadWallTimeMs" | "kv_read_wall_time_ms" => {
                                Ok(GeneratedField::KvReadWallTimeMs)
                            }
                            "totalRpcWallTimeNs" | "total_rpc_wall_time_ns" => {
                                Ok(GeneratedField::TotalRpcWallTimeNs)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = TimeDetail;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.TimeDetail")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<TimeDetail, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut wait_wall_time_ms__ = None;
                let mut process_wall_time_ms__ = None;
                let mut kv_read_wall_time_ms__ = None;
                let mut total_rpc_wall_time_ns__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WaitWallTimeMs => {
                            if wait_wall_time_ms__.is_some() {
                                return Err(serde::de::Error::duplicate_field("waitWallTimeMs"));
                            }
                            wait_wall_time_ms__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ProcessWallTimeMs => {
                            if process_wall_time_ms__.is_some() {
                                return Err(serde::de::Error::duplicate_field("processWallTimeMs"));
                            }
                            process_wall_time_ms__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::KvReadWallTimeMs => {
                            if kv_read_wall_time_ms__.is_some() {
                                return Err(serde::de::Error::duplicate_field("kvReadWallTimeMs"));
                            }
                            kv_read_wall_time_ms__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::TotalRpcWallTimeNs => {
                            if total_rpc_wall_time_ns__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "totalRpcWallTimeNs",
                                ));
                            }
                            total_rpc_wall_time_ns__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(TimeDetail {
                    wait_wall_time_ms: wait_wall_time_ms__.unwrap_or_default(),
                    process_wall_time_ms: process_wall_time_ms__.unwrap_or_default(),
                    kv_read_wall_time_ms: kv_read_wall_time_ms__.unwrap_or_default(),
                    total_rpc_wall_time_ns: total_rpc_wall_time_ns__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.TimeDetail", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for TimeDetailV2 {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.wait_wall_time_ns != 0 {
            len += 1;
        }
        if self.process_wall_time_ns != 0 {
            len += 1;
        }
        if self.process_suspend_wall_time_ns != 0 {
            len += 1;
        }
        if self.kv_read_wall_time_ns != 0 {
            len += 1;
        }
        if self.total_rpc_wall_time_ns != 0 {
            len += 1;
        }
        if self.kv_grpc_process_time_ns != 0 {
            len += 1;
        }
        if self.kv_grpc_wait_time_ns != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.TimeDetailV2", len)?;
        if self.wait_wall_time_ns != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "waitWallTimeNs",
                ToString::to_string(&self.wait_wall_time_ns).as_str(),
            )?;
        }
        if self.process_wall_time_ns != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "processWallTimeNs",
                ToString::to_string(&self.process_wall_time_ns).as_str(),
            )?;
        }
        if self.process_suspend_wall_time_ns != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "processSuspendWallTimeNs",
                ToString::to_string(&self.process_suspend_wall_time_ns).as_str(),
            )?;
        }
        if self.kv_read_wall_time_ns != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "kvReadWallTimeNs",
                ToString::to_string(&self.kv_read_wall_time_ns).as_str(),
            )?;
        }
        if self.total_rpc_wall_time_ns != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "totalRpcWallTimeNs",
                ToString::to_string(&self.total_rpc_wall_time_ns).as_str(),
            )?;
        }
        if self.kv_grpc_process_time_ns != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "kvGrpcProcessTimeNs",
                ToString::to_string(&self.kv_grpc_process_time_ns).as_str(),
            )?;
        }
        if self.kv_grpc_wait_time_ns != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "kvGrpcWaitTimeNs",
                ToString::to_string(&self.kv_grpc_wait_time_ns).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for TimeDetailV2 {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "wait_wall_time_ns",
            "waitWallTimeNs",
            "process_wall_time_ns",
            "processWallTimeNs",
            "process_suspend_wall_time_ns",
            "processSuspendWallTimeNs",
            "kv_read_wall_time_ns",
            "kvReadWallTimeNs",
            "total_rpc_wall_time_ns",
            "totalRpcWallTimeNs",
            "kv_grpc_process_time_ns",
            "kvGrpcProcessTimeNs",
            "kv_grpc_wait_time_ns",
            "kvGrpcWaitTimeNs",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            WaitWallTimeNs,
            ProcessWallTimeNs,
            ProcessSuspendWallTimeNs,
            KvReadWallTimeNs,
            TotalRpcWallTimeNs,
            KvGrpcProcessTimeNs,
            KvGrpcWaitTimeNs,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "waitWallTimeNs" | "wait_wall_time_ns" => {
                                Ok(GeneratedField::WaitWallTimeNs)
                            }
                            "processWallTimeNs" | "process_wall_time_ns" => {
                                Ok(GeneratedField::ProcessWallTimeNs)
                            }
                            "processSuspendWallTimeNs" | "process_suspend_wall_time_ns" => {
                                Ok(GeneratedField::ProcessSuspendWallTimeNs)
                            }
                            "kvReadWallTimeNs" | "kv_read_wall_time_ns" => {
                                Ok(GeneratedField::KvReadWallTimeNs)
                            }
                            "totalRpcWallTimeNs" | "total_rpc_wall_time_ns" => {
                                Ok(GeneratedField::TotalRpcWallTimeNs)
                            }
                            "kvGrpcProcessTimeNs" | "kv_grpc_process_time_ns" => {
                                Ok(GeneratedField::KvGrpcProcessTimeNs)
                            }
                            "kvGrpcWaitTimeNs" | "kv_grpc_wait_time_ns" => {
                                Ok(GeneratedField::KvGrpcWaitTimeNs)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = TimeDetailV2;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.TimeDetailV2")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<TimeDetailV2, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut wait_wall_time_ns__ = None;
                let mut process_wall_time_ns__ = None;
                let mut process_suspend_wall_time_ns__ = None;
                let mut kv_read_wall_time_ns__ = None;
                let mut total_rpc_wall_time_ns__ = None;
                let mut kv_grpc_process_time_ns__ = None;
                let mut kv_grpc_wait_time_ns__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::WaitWallTimeNs => {
                            if wait_wall_time_ns__.is_some() {
                                return Err(serde::de::Error::duplicate_field("waitWallTimeNs"));
                            }
                            wait_wall_time_ns__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ProcessWallTimeNs => {
                            if process_wall_time_ns__.is_some() {
                                return Err(serde::de::Error::duplicate_field("processWallTimeNs"));
                            }
                            process_wall_time_ns__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ProcessSuspendWallTimeNs => {
                            if process_suspend_wall_time_ns__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "processSuspendWallTimeNs",
                                ));
                            }
                            process_suspend_wall_time_ns__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::KvReadWallTimeNs => {
                            if kv_read_wall_time_ns__.is_some() {
                                return Err(serde::de::Error::duplicate_field("kvReadWallTimeNs"));
                            }
                            kv_read_wall_time_ns__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::TotalRpcWallTimeNs => {
                            if total_rpc_wall_time_ns__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "totalRpcWallTimeNs",
                                ));
                            }
                            total_rpc_wall_time_ns__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::KvGrpcProcessTimeNs => {
                            if kv_grpc_process_time_ns__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "kvGrpcProcessTimeNs",
                                ));
                            }
                            kv_grpc_process_time_ns__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::KvGrpcWaitTimeNs => {
                            if kv_grpc_wait_time_ns__.is_some() {
                                return Err(serde::de::Error::duplicate_field("kvGrpcWaitTimeNs"));
                            }
                            kv_grpc_wait_time_ns__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(TimeDetailV2 {
                    wait_wall_time_ns: wait_wall_time_ns__.unwrap_or_default(),
                    process_wall_time_ns: process_wall_time_ns__.unwrap_or_default(),
                    process_suspend_wall_time_ns: process_suspend_wall_time_ns__
                        .unwrap_or_default(),
                    kv_read_wall_time_ns: kv_read_wall_time_ns__.unwrap_or_default(),
                    total_rpc_wall_time_ns: total_rpc_wall_time_ns__.unwrap_or_default(),
                    kv_grpc_process_time_ns: kv_grpc_process_time_ns__.unwrap_or_default(),
                    kv_grpc_wait_time_ns: kv_grpc_wait_time_ns__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.TimeDetailV2", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for TxnLockNotFound {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.key.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.TxnLockNotFound", len)?;
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for TxnLockNotFound {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["key"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Key,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "key" => Ok(GeneratedField::Key),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = TxnLockNotFound;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.TxnLockNotFound")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<TxnLockNotFound, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut key__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(TxnLockNotFound {
                    key: key__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.TxnLockNotFound", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for TxnNotFound {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.start_ts != 0 {
            len += 1;
        }
        if !self.primary_key.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.TxnNotFound", len)?;
        if self.start_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("startTs", ToString::to_string(&self.start_ts).as_str())?;
        }
        if !self.primary_key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "primaryKey",
                pbjson::private::base64::encode(&self.primary_key).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for TxnNotFound {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["start_ts", "startTs", "primary_key", "primaryKey"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            StartTs,
            PrimaryKey,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "startTs" | "start_ts" => Ok(GeneratedField::StartTs),
                            "primaryKey" | "primary_key" => Ok(GeneratedField::PrimaryKey),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = TxnNotFound;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.TxnNotFound")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<TxnNotFound, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut start_ts__ = None;
                let mut primary_key__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::StartTs => {
                            if start_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startTs"));
                            }
                            start_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::PrimaryKey => {
                            if primary_key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("primaryKey"));
                            }
                            primary_key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(TxnNotFound {
                    start_ts: start_ts__.unwrap_or_default(),
                    primary_key: primary_key__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.TxnNotFound", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for UndeterminedResult {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.message.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.UndeterminedResult", len)?;
        if !self.message.is_empty() {
            struct_ser.serialize_field("message", &self.message)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for UndeterminedResult {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["message"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
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
            type Value = UndeterminedResult;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.UndeterminedResult")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<UndeterminedResult, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut message__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Message => {
                            if message__.is_some() {
                                return Err(serde::de::Error::duplicate_field("message"));
                            }
                            message__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(UndeterminedResult {
                    message: message__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.UndeterminedResult", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for WriteConflict {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.start_ts != 0 {
            len += 1;
        }
        if self.conflict_ts != 0 {
            len += 1;
        }
        if !self.key.is_empty() {
            len += 1;
        }
        if !self.primary.is_empty() {
            len += 1;
        }
        if self.conflict_commit_ts != 0 {
            len += 1;
        }
        if self.reason != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.WriteConflict", len)?;
        if self.start_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("startTs", ToString::to_string(&self.start_ts).as_str())?;
        }
        if self.conflict_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "conflictTs",
                ToString::to_string(&self.conflict_ts).as_str(),
            )?;
        }
        if !self.key.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser
                .serialize_field("key", pbjson::private::base64::encode(&self.key).as_str())?;
        }
        if !self.primary.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "primary",
                pbjson::private::base64::encode(&self.primary).as_str(),
            )?;
        }
        if self.conflict_commit_ts != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "conflictCommitTs",
                ToString::to_string(&self.conflict_commit_ts).as_str(),
            )?;
        }
        if self.reason != 0 {
            let v = write_conflict::Reason::try_from(self.reason).map_err(|_| {
                serde::ser::Error::custom(format!("Invalid variant {}", self.reason))
            })?;
            struct_ser.serialize_field("reason", &v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for WriteConflict {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "start_ts",
            "startTs",
            "conflict_ts",
            "conflictTs",
            "key",
            "primary",
            "conflict_commit_ts",
            "conflictCommitTs",
            "reason",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            StartTs,
            ConflictTs,
            Key,
            Primary,
            ConflictCommitTs,
            Reason,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "startTs" | "start_ts" => Ok(GeneratedField::StartTs),
                            "conflictTs" | "conflict_ts" => Ok(GeneratedField::ConflictTs),
                            "key" => Ok(GeneratedField::Key),
                            "primary" => Ok(GeneratedField::Primary),
                            "conflictCommitTs" | "conflict_commit_ts" => {
                                Ok(GeneratedField::ConflictCommitTs)
                            }
                            "reason" => Ok(GeneratedField::Reason),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = WriteConflict;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.WriteConflict")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<WriteConflict, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut start_ts__ = None;
                let mut conflict_ts__ = None;
                let mut key__ = None;
                let mut primary__ = None;
                let mut conflict_commit_ts__ = None;
                let mut reason__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::StartTs => {
                            if start_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startTs"));
                            }
                            start_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ConflictTs => {
                            if conflict_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("conflictTs"));
                            }
                            conflict_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Primary => {
                            if primary__.is_some() {
                                return Err(serde::de::Error::duplicate_field("primary"));
                            }
                            primary__ = Some(
                                map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ConflictCommitTs => {
                            if conflict_commit_ts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("conflictCommitTs"));
                            }
                            conflict_commit_ts__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Reason => {
                            if reason__.is_some() {
                                return Err(serde::de::Error::duplicate_field("reason"));
                            }
                            reason__ = Some(map_.next_value::<write_conflict::Reason>()? as i32);
                        }
                    }
                }
                Ok(WriteConflict {
                    start_ts: start_ts__.unwrap_or_default(),
                    conflict_ts: conflict_ts__.unwrap_or_default(),
                    key: key__.unwrap_or_default(),
                    primary: primary__.unwrap_or_default(),
                    conflict_commit_ts: conflict_commit_ts__.unwrap_or_default(),
                    reason: reason__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.WriteConflict", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for write_conflict::Reason {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unknown => "Unknown",
            Self::Optimistic => "Optimistic",
            Self::PessimisticRetry => "PessimisticRetry",
            Self::SelfRolledBack => "SelfRolledBack",
            Self::RcCheckTs => "RcCheckTs",
            Self::LazyUniquenessCheck => "LazyUniquenessCheck",
            Self::NotLockedKeyConflict => "NotLockedKeyConflict",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for write_conflict::Reason {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "Unknown",
            "Optimistic",
            "PessimisticRetry",
            "SelfRolledBack",
            "RcCheckTs",
            "LazyUniquenessCheck",
            "NotLockedKeyConflict",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = write_conflict::Reason;

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
                    "Unknown" => Ok(write_conflict::Reason::Unknown),
                    "Optimistic" => Ok(write_conflict::Reason::Optimistic),
                    "PessimisticRetry" => Ok(write_conflict::Reason::PessimisticRetry),
                    "SelfRolledBack" => Ok(write_conflict::Reason::SelfRolledBack),
                    "RcCheckTs" => Ok(write_conflict::Reason::RcCheckTs),
                    "LazyUniquenessCheck" => Ok(write_conflict::Reason::LazyUniquenessCheck),
                    "NotLockedKeyConflict" => Ok(write_conflict::Reason::NotLockedKeyConflict),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for WriteDetail {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.store_batch_wait_nanos != 0 {
            len += 1;
        }
        if self.propose_send_wait_nanos != 0 {
            len += 1;
        }
        if self.persist_log_nanos != 0 {
            len += 1;
        }
        if self.raft_db_write_leader_wait_nanos != 0 {
            len += 1;
        }
        if self.raft_db_sync_log_nanos != 0 {
            len += 1;
        }
        if self.raft_db_write_memtable_nanos != 0 {
            len += 1;
        }
        if self.commit_log_nanos != 0 {
            len += 1;
        }
        if self.apply_batch_wait_nanos != 0 {
            len += 1;
        }
        if self.apply_log_nanos != 0 {
            len += 1;
        }
        if self.apply_mutex_lock_nanos != 0 {
            len += 1;
        }
        if self.apply_write_leader_wait_nanos != 0 {
            len += 1;
        }
        if self.apply_write_wal_nanos != 0 {
            len += 1;
        }
        if self.apply_write_memtable_nanos != 0 {
            len += 1;
        }
        if self.latch_wait_nanos != 0 {
            len += 1;
        }
        if self.process_nanos != 0 {
            len += 1;
        }
        if self.throttle_nanos != 0 {
            len += 1;
        }
        if self.pessimistic_lock_wait_nanos != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("tikv.WriteDetail", len)?;
        if self.store_batch_wait_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "storeBatchWaitNanos",
                ToString::to_string(&self.store_batch_wait_nanos).as_str(),
            )?;
        }
        if self.propose_send_wait_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "proposeSendWaitNanos",
                ToString::to_string(&self.propose_send_wait_nanos).as_str(),
            )?;
        }
        if self.persist_log_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "persistLogNanos",
                ToString::to_string(&self.persist_log_nanos).as_str(),
            )?;
        }
        if self.raft_db_write_leader_wait_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "raftDbWriteLeaderWaitNanos",
                ToString::to_string(&self.raft_db_write_leader_wait_nanos).as_str(),
            )?;
        }
        if self.raft_db_sync_log_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "raftDbSyncLogNanos",
                ToString::to_string(&self.raft_db_sync_log_nanos).as_str(),
            )?;
        }
        if self.raft_db_write_memtable_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "raftDbWriteMemtableNanos",
                ToString::to_string(&self.raft_db_write_memtable_nanos).as_str(),
            )?;
        }
        if self.commit_log_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "commitLogNanos",
                ToString::to_string(&self.commit_log_nanos).as_str(),
            )?;
        }
        if self.apply_batch_wait_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "applyBatchWaitNanos",
                ToString::to_string(&self.apply_batch_wait_nanos).as_str(),
            )?;
        }
        if self.apply_log_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "applyLogNanos",
                ToString::to_string(&self.apply_log_nanos).as_str(),
            )?;
        }
        if self.apply_mutex_lock_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "applyMutexLockNanos",
                ToString::to_string(&self.apply_mutex_lock_nanos).as_str(),
            )?;
        }
        if self.apply_write_leader_wait_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "applyWriteLeaderWaitNanos",
                ToString::to_string(&self.apply_write_leader_wait_nanos).as_str(),
            )?;
        }
        if self.apply_write_wal_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "applyWriteWalNanos",
                ToString::to_string(&self.apply_write_wal_nanos).as_str(),
            )?;
        }
        if self.apply_write_memtable_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "applyWriteMemtableNanos",
                ToString::to_string(&self.apply_write_memtable_nanos).as_str(),
            )?;
        }
        if self.latch_wait_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "latchWaitNanos",
                ToString::to_string(&self.latch_wait_nanos).as_str(),
            )?;
        }
        if self.process_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "processNanos",
                ToString::to_string(&self.process_nanos).as_str(),
            )?;
        }
        if self.throttle_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "throttleNanos",
                ToString::to_string(&self.throttle_nanos).as_str(),
            )?;
        }
        if self.pessimistic_lock_wait_nanos != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "pessimisticLockWaitNanos",
                ToString::to_string(&self.pessimistic_lock_wait_nanos).as_str(),
            )?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for WriteDetail {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "store_batch_wait_nanos",
            "storeBatchWaitNanos",
            "propose_send_wait_nanos",
            "proposeSendWaitNanos",
            "persist_log_nanos",
            "persistLogNanos",
            "raft_db_write_leader_wait_nanos",
            "raftDbWriteLeaderWaitNanos",
            "raft_db_sync_log_nanos",
            "raftDbSyncLogNanos",
            "raft_db_write_memtable_nanos",
            "raftDbWriteMemtableNanos",
            "commit_log_nanos",
            "commitLogNanos",
            "apply_batch_wait_nanos",
            "applyBatchWaitNanos",
            "apply_log_nanos",
            "applyLogNanos",
            "apply_mutex_lock_nanos",
            "applyMutexLockNanos",
            "apply_write_leader_wait_nanos",
            "applyWriteLeaderWaitNanos",
            "apply_write_wal_nanos",
            "applyWriteWalNanos",
            "apply_write_memtable_nanos",
            "applyWriteMemtableNanos",
            "latch_wait_nanos",
            "latchWaitNanos",
            "process_nanos",
            "processNanos",
            "throttle_nanos",
            "throttleNanos",
            "pessimistic_lock_wait_nanos",
            "pessimisticLockWaitNanos",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            StoreBatchWaitNanos,
            ProposeSendWaitNanos,
            PersistLogNanos,
            RaftDbWriteLeaderWaitNanos,
            RaftDbSyncLogNanos,
            RaftDbWriteMemtableNanos,
            CommitLogNanos,
            ApplyBatchWaitNanos,
            ApplyLogNanos,
            ApplyMutexLockNanos,
            ApplyWriteLeaderWaitNanos,
            ApplyWriteWalNanos,
            ApplyWriteMemtableNanos,
            LatchWaitNanos,
            ProcessNanos,
            ThrottleNanos,
            PessimisticLockWaitNanos,
        }
        impl<'de> serde::Deserialize<'de> for GeneratedField {
            fn deserialize<D>(deserializer: D) -> std::result::Result<GeneratedField, D::Error>
            where
                D: serde::Deserializer<'de>,
            {
                struct GeneratedVisitor;

                impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
                    type Value = GeneratedField;

                    fn expecting(
                        &self,
                        formatter: &mut std::fmt::Formatter<'_>,
                    ) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "storeBatchWaitNanos" | "store_batch_wait_nanos" => {
                                Ok(GeneratedField::StoreBatchWaitNanos)
                            }
                            "proposeSendWaitNanos" | "propose_send_wait_nanos" => {
                                Ok(GeneratedField::ProposeSendWaitNanos)
                            }
                            "persistLogNanos" | "persist_log_nanos" => {
                                Ok(GeneratedField::PersistLogNanos)
                            }
                            "raftDbWriteLeaderWaitNanos" | "raft_db_write_leader_wait_nanos" => {
                                Ok(GeneratedField::RaftDbWriteLeaderWaitNanos)
                            }
                            "raftDbSyncLogNanos" | "raft_db_sync_log_nanos" => {
                                Ok(GeneratedField::RaftDbSyncLogNanos)
                            }
                            "raftDbWriteMemtableNanos" | "raft_db_write_memtable_nanos" => {
                                Ok(GeneratedField::RaftDbWriteMemtableNanos)
                            }
                            "commitLogNanos" | "commit_log_nanos" => {
                                Ok(GeneratedField::CommitLogNanos)
                            }
                            "applyBatchWaitNanos" | "apply_batch_wait_nanos" => {
                                Ok(GeneratedField::ApplyBatchWaitNanos)
                            }
                            "applyLogNanos" | "apply_log_nanos" => {
                                Ok(GeneratedField::ApplyLogNanos)
                            }
                            "applyMutexLockNanos" | "apply_mutex_lock_nanos" => {
                                Ok(GeneratedField::ApplyMutexLockNanos)
                            }
                            "applyWriteLeaderWaitNanos" | "apply_write_leader_wait_nanos" => {
                                Ok(GeneratedField::ApplyWriteLeaderWaitNanos)
                            }
                            "applyWriteWalNanos" | "apply_write_wal_nanos" => {
                                Ok(GeneratedField::ApplyWriteWalNanos)
                            }
                            "applyWriteMemtableNanos" | "apply_write_memtable_nanos" => {
                                Ok(GeneratedField::ApplyWriteMemtableNanos)
                            }
                            "latchWaitNanos" | "latch_wait_nanos" => {
                                Ok(GeneratedField::LatchWaitNanos)
                            }
                            "processNanos" | "process_nanos" => Ok(GeneratedField::ProcessNanos),
                            "throttleNanos" | "throttle_nanos" => Ok(GeneratedField::ThrottleNanos),
                            "pessimisticLockWaitNanos" | "pessimistic_lock_wait_nanos" => {
                                Ok(GeneratedField::PessimisticLockWaitNanos)
                            }
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = WriteDetail;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct tikv.WriteDetail")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<WriteDetail, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut store_batch_wait_nanos__ = None;
                let mut propose_send_wait_nanos__ = None;
                let mut persist_log_nanos__ = None;
                let mut raft_db_write_leader_wait_nanos__ = None;
                let mut raft_db_sync_log_nanos__ = None;
                let mut raft_db_write_memtable_nanos__ = None;
                let mut commit_log_nanos__ = None;
                let mut apply_batch_wait_nanos__ = None;
                let mut apply_log_nanos__ = None;
                let mut apply_mutex_lock_nanos__ = None;
                let mut apply_write_leader_wait_nanos__ = None;
                let mut apply_write_wal_nanos__ = None;
                let mut apply_write_memtable_nanos__ = None;
                let mut latch_wait_nanos__ = None;
                let mut process_nanos__ = None;
                let mut throttle_nanos__ = None;
                let mut pessimistic_lock_wait_nanos__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::StoreBatchWaitNanos => {
                            if store_batch_wait_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "storeBatchWaitNanos",
                                ));
                            }
                            store_batch_wait_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ProposeSendWaitNanos => {
                            if propose_send_wait_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "proposeSendWaitNanos",
                                ));
                            }
                            propose_send_wait_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::PersistLogNanos => {
                            if persist_log_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field("persistLogNanos"));
                            }
                            persist_log_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::RaftDbWriteLeaderWaitNanos => {
                            if raft_db_write_leader_wait_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "raftDbWriteLeaderWaitNanos",
                                ));
                            }
                            raft_db_write_leader_wait_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::RaftDbSyncLogNanos => {
                            if raft_db_sync_log_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "raftDbSyncLogNanos",
                                ));
                            }
                            raft_db_sync_log_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::RaftDbWriteMemtableNanos => {
                            if raft_db_write_memtable_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "raftDbWriteMemtableNanos",
                                ));
                            }
                            raft_db_write_memtable_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::CommitLogNanos => {
                            if commit_log_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field("commitLogNanos"));
                            }
                            commit_log_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ApplyBatchWaitNanos => {
                            if apply_batch_wait_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "applyBatchWaitNanos",
                                ));
                            }
                            apply_batch_wait_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ApplyLogNanos => {
                            if apply_log_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field("applyLogNanos"));
                            }
                            apply_log_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ApplyMutexLockNanos => {
                            if apply_mutex_lock_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "applyMutexLockNanos",
                                ));
                            }
                            apply_mutex_lock_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ApplyWriteLeaderWaitNanos => {
                            if apply_write_leader_wait_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "applyWriteLeaderWaitNanos",
                                ));
                            }
                            apply_write_leader_wait_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ApplyWriteWalNanos => {
                            if apply_write_wal_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "applyWriteWalNanos",
                                ));
                            }
                            apply_write_wal_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ApplyWriteMemtableNanos => {
                            if apply_write_memtable_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "applyWriteMemtableNanos",
                                ));
                            }
                            apply_write_memtable_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::LatchWaitNanos => {
                            if latch_wait_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field("latchWaitNanos"));
                            }
                            latch_wait_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ProcessNanos => {
                            if process_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field("processNanos"));
                            }
                            process_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::ThrottleNanos => {
                            if throttle_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field("throttleNanos"));
                            }
                            throttle_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::PessimisticLockWaitNanos => {
                            if pessimistic_lock_wait_nanos__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "pessimisticLockWaitNanos",
                                ));
                            }
                            pessimistic_lock_wait_nanos__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(WriteDetail {
                    store_batch_wait_nanos: store_batch_wait_nanos__.unwrap_or_default(),
                    propose_send_wait_nanos: propose_send_wait_nanos__.unwrap_or_default(),
                    persist_log_nanos: persist_log_nanos__.unwrap_or_default(),
                    raft_db_write_leader_wait_nanos: raft_db_write_leader_wait_nanos__
                        .unwrap_or_default(),
                    raft_db_sync_log_nanos: raft_db_sync_log_nanos__.unwrap_or_default(),
                    raft_db_write_memtable_nanos: raft_db_write_memtable_nanos__
                        .unwrap_or_default(),
                    commit_log_nanos: commit_log_nanos__.unwrap_or_default(),
                    apply_batch_wait_nanos: apply_batch_wait_nanos__.unwrap_or_default(),
                    apply_log_nanos: apply_log_nanos__.unwrap_or_default(),
                    apply_mutex_lock_nanos: apply_mutex_lock_nanos__.unwrap_or_default(),
                    apply_write_leader_wait_nanos: apply_write_leader_wait_nanos__
                        .unwrap_or_default(),
                    apply_write_wal_nanos: apply_write_wal_nanos__.unwrap_or_default(),
                    apply_write_memtable_nanos: apply_write_memtable_nanos__.unwrap_or_default(),
                    latch_wait_nanos: latch_wait_nanos__.unwrap_or_default(),
                    process_nanos: process_nanos__.unwrap_or_default(),
                    throttle_nanos: throttle_nanos__.unwrap_or_default(),
                    pessimistic_lock_wait_nanos: pessimistic_lock_wait_nanos__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("tikv.WriteDetail", FIELDS, GeneratedVisitor)
    }
}
