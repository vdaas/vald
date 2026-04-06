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
impl serde::Serialize for Peer {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.id != 0 {
            len += 1;
        }
        if self.store_id != 0 {
            len += 1;
        }
        if self.is_witness {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("metapb.Peer", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if self.store_id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("storeId", ToString::to_string(&self.store_id).as_str())?;
        }
        if self.is_witness {
            struct_ser.serialize_field("isWitness", &self.is_witness)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Peer {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["id", "store_id", "storeId", "is_witness", "isWitness"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            StoreId,
            IsWitness,
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
                            "id" => Ok(GeneratedField::Id),
                            "storeId" | "store_id" => Ok(GeneratedField::StoreId),
                            "isWitness" | "is_witness" => Ok(GeneratedField::IsWitness),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Peer;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct metapb.Peer")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Peer, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut store_id__ = None;
                let mut is_witness__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::StoreId => {
                            if store_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("storeId"));
                            }
                            store_id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::IsWitness => {
                            if is_witness__.is_some() {
                                return Err(serde::de::Error::duplicate_field("isWitness"));
                            }
                            is_witness__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(Peer {
                    id: id__.unwrap_or_default(),
                    store_id: store_id__.unwrap_or_default(),
                    is_witness: is_witness__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("metapb.Peer", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Region2 {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.id != 0 {
            len += 1;
        }
        if !self.start_key.is_empty() {
            len += 1;
        }
        if !self.end_key.is_empty() {
            len += 1;
        }
        if self.region_epoch.is_some() {
            len += 1;
        }
        if !self.peers.is_empty() {
            len += 1;
        }
        if self.is_in_flashback {
            len += 1;
        }
        if self.flashback_start_ts != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("metapb.Region2", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
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
        if let Some(v) = self.region_epoch.as_ref() {
            struct_ser.serialize_field("regionEpoch", v)?;
        }
        if !self.peers.is_empty() {
            struct_ser.serialize_field("peers", &self.peers)?;
        }
        if self.is_in_flashback {
            struct_ser.serialize_field("isInFlashback", &self.is_in_flashback)?;
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
impl<'de> serde::Deserialize<'de> for Region2 {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "start_key",
            "startKey",
            "end_key",
            "endKey",
            "region_epoch",
            "regionEpoch",
            "peers",
            "is_in_flashback",
            "isInFlashback",
            "flashback_start_ts",
            "flashbackStartTs",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            StartKey,
            EndKey,
            RegionEpoch,
            Peers,
            IsInFlashback,
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
                            "id" => Ok(GeneratedField::Id),
                            "startKey" | "start_key" => Ok(GeneratedField::StartKey),
                            "endKey" | "end_key" => Ok(GeneratedField::EndKey),
                            "regionEpoch" | "region_epoch" => Ok(GeneratedField::RegionEpoch),
                            "peers" => Ok(GeneratedField::Peers),
                            "isInFlashback" | "is_in_flashback" => {
                                Ok(GeneratedField::IsInFlashback)
                            }
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
            type Value = Region2;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct metapb.Region2")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Region2, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut start_key__ = None;
                let mut end_key__ = None;
                let mut region_epoch__ = None;
                let mut peers__ = None;
                let mut is_in_flashback__ = None;
                let mut flashback_start_ts__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(
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
                        GeneratedField::RegionEpoch => {
                            if region_epoch__.is_some() {
                                return Err(serde::de::Error::duplicate_field("regionEpoch"));
                            }
                            region_epoch__ = map_.next_value()?;
                        }
                        GeneratedField::Peers => {
                            if peers__.is_some() {
                                return Err(serde::de::Error::duplicate_field("peers"));
                            }
                            peers__ = Some(map_.next_value()?);
                        }
                        GeneratedField::IsInFlashback => {
                            if is_in_flashback__.is_some() {
                                return Err(serde::de::Error::duplicate_field("isInFlashback"));
                            }
                            is_in_flashback__ = Some(map_.next_value()?);
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
                Ok(Region2 {
                    id: id__.unwrap_or_default(),
                    start_key: start_key__.unwrap_or_default(),
                    end_key: end_key__.unwrap_or_default(),
                    region_epoch: region_epoch__,
                    peers: peers__.unwrap_or_default(),
                    is_in_flashback: is_in_flashback__.unwrap_or_default(),
                    flashback_start_ts: flashback_start_ts__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("metapb.Region2", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for RegionEpoch {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.conf_ver != 0 {
            len += 1;
        }
        if self.version != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("metapb.RegionEpoch", len)?;
        if self.conf_ver != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("confVer", ToString::to_string(&self.conf_ver).as_str())?;
        }
        if self.version != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("version", ToString::to_string(&self.version).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for RegionEpoch {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["conf_ver", "confVer", "version"];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            ConfVer,
            Version,
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
                            "confVer" | "conf_ver" => Ok(GeneratedField::ConfVer),
                            "version" => Ok(GeneratedField::Version),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = RegionEpoch;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct metapb.RegionEpoch")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<RegionEpoch, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut conf_ver__ = None;
                let mut version__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::ConfVer => {
                            if conf_ver__.is_some() {
                                return Err(serde::de::Error::duplicate_field("confVer"));
                            }
                            conf_ver__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Version => {
                            if version__.is_some() {
                                return Err(serde::de::Error::duplicate_field("version"));
                            }
                            version__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                    }
                }
                Ok(RegionEpoch {
                    conf_ver: conf_ver__.unwrap_or_default(),
                    version: version__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("metapb.RegionEpoch", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Store {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.id != 0 {
            len += 1;
        }
        if !self.address.is_empty() {
            len += 1;
        }
        if self.state != 0 {
            len += 1;
        }
        if !self.version.is_empty() {
            len += 1;
        }
        if !self.peer_address.is_empty() {
            len += 1;
        }
        if !self.status_address.is_empty() {
            len += 1;
        }
        if !self.git_hash.is_empty() {
            len += 1;
        }
        if self.start_timestamp != 0 {
            len += 1;
        }
        if !self.deploy_path.is_empty() {
            len += 1;
        }
        if self.last_heartbeat != 0 {
            len += 1;
        }
        if self.physically_destroyed {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("metapb.Store", len)?;
        if self.id != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("id", ToString::to_string(&self.id).as_str())?;
        }
        if !self.address.is_empty() {
            struct_ser.serialize_field("address", &self.address)?;
        }
        if self.state != 0 {
            let v = StoreState::try_from(self.state).map_err(|_| {
                serde::ser::Error::custom(format!("Invalid variant {}", self.state))
            })?;
            struct_ser.serialize_field("state", &v)?;
        }
        if !self.version.is_empty() {
            struct_ser.serialize_field("version", &self.version)?;
        }
        if !self.peer_address.is_empty() {
            struct_ser.serialize_field("peerAddress", &self.peer_address)?;
        }
        if !self.status_address.is_empty() {
            struct_ser.serialize_field("statusAddress", &self.status_address)?;
        }
        if !self.git_hash.is_empty() {
            struct_ser.serialize_field("gitHash", &self.git_hash)?;
        }
        if self.start_timestamp != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "startTimestamp",
                ToString::to_string(&self.start_timestamp).as_str(),
            )?;
        }
        if !self.deploy_path.is_empty() {
            struct_ser.serialize_field("deployPath", &self.deploy_path)?;
        }
        if self.last_heartbeat != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field(
                "lastHeartbeat",
                ToString::to_string(&self.last_heartbeat).as_str(),
            )?;
        }
        if self.physically_destroyed {
            struct_ser.serialize_field("physicallyDestroyed", &self.physically_destroyed)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Store {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "address",
            "state",
            "version",
            "peer_address",
            "peerAddress",
            "status_address",
            "statusAddress",
            "git_hash",
            "gitHash",
            "start_timestamp",
            "startTimestamp",
            "deploy_path",
            "deployPath",
            "last_heartbeat",
            "lastHeartbeat",
            "physically_destroyed",
            "physicallyDestroyed",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Address,
            State,
            Version,
            PeerAddress,
            StatusAddress,
            GitHash,
            StartTimestamp,
            DeployPath,
            LastHeartbeat,
            PhysicallyDestroyed,
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
                            "id" => Ok(GeneratedField::Id),
                            "address" => Ok(GeneratedField::Address),
                            "state" => Ok(GeneratedField::State),
                            "version" => Ok(GeneratedField::Version),
                            "peerAddress" | "peer_address" => Ok(GeneratedField::PeerAddress),
                            "statusAddress" | "status_address" => Ok(GeneratedField::StatusAddress),
                            "gitHash" | "git_hash" => Ok(GeneratedField::GitHash),
                            "startTimestamp" | "start_timestamp" => {
                                Ok(GeneratedField::StartTimestamp)
                            }
                            "deployPath" | "deploy_path" => Ok(GeneratedField::DeployPath),
                            "lastHeartbeat" | "last_heartbeat" => Ok(GeneratedField::LastHeartbeat),
                            "physicallyDestroyed" | "physically_destroyed" => {
                                Ok(GeneratedField::PhysicallyDestroyed)
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
            type Value = Store;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct metapb.Store")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Store, V::Error>
            where
                V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut address__ = None;
                let mut state__ = None;
                let mut version__ = None;
                let mut peer_address__ = None;
                let mut status_address__ = None;
                let mut git_hash__ = None;
                let mut start_timestamp__ = None;
                let mut deploy_path__ = None;
                let mut last_heartbeat__ = None;
                let mut physically_destroyed__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::Address => {
                            if address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("address"));
                            }
                            address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::State => {
                            if state__.is_some() {
                                return Err(serde::de::Error::duplicate_field("state"));
                            }
                            state__ = Some(map_.next_value::<StoreState>()? as i32);
                        }
                        GeneratedField::Version => {
                            if version__.is_some() {
                                return Err(serde::de::Error::duplicate_field("version"));
                            }
                            version__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PeerAddress => {
                            if peer_address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("peerAddress"));
                            }
                            peer_address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::StatusAddress => {
                            if status_address__.is_some() {
                                return Err(serde::de::Error::duplicate_field("statusAddress"));
                            }
                            status_address__ = Some(map_.next_value()?);
                        }
                        GeneratedField::GitHash => {
                            if git_hash__.is_some() {
                                return Err(serde::de::Error::duplicate_field("gitHash"));
                            }
                            git_hash__ = Some(map_.next_value()?);
                        }
                        GeneratedField::StartTimestamp => {
                            if start_timestamp__.is_some() {
                                return Err(serde::de::Error::duplicate_field("startTimestamp"));
                            }
                            start_timestamp__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::DeployPath => {
                            if deploy_path__.is_some() {
                                return Err(serde::de::Error::duplicate_field("deployPath"));
                            }
                            deploy_path__ = Some(map_.next_value()?);
                        }
                        GeneratedField::LastHeartbeat => {
                            if last_heartbeat__.is_some() {
                                return Err(serde::de::Error::duplicate_field("lastHeartbeat"));
                            }
                            last_heartbeat__ = Some(
                                map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?
                                    .0,
                            );
                        }
                        GeneratedField::PhysicallyDestroyed => {
                            if physically_destroyed__.is_some() {
                                return Err(serde::de::Error::duplicate_field(
                                    "physicallyDestroyed",
                                ));
                            }
                            physically_destroyed__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(Store {
                    id: id__.unwrap_or_default(),
                    address: address__.unwrap_or_default(),
                    state: state__.unwrap_or_default(),
                    version: version__.unwrap_or_default(),
                    peer_address: peer_address__.unwrap_or_default(),
                    status_address: status_address__.unwrap_or_default(),
                    git_hash: git_hash__.unwrap_or_default(),
                    start_timestamp: start_timestamp__.unwrap_or_default(),
                    deploy_path: deploy_path__.unwrap_or_default(),
                    last_heartbeat: last_heartbeat__.unwrap_or_default(),
                    physically_destroyed: physically_destroyed__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("metapb.Store", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for StoreState {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Up => "Up",
            Self::Offline => "Offline",
            Self::Tombstone => "Tombstone",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for StoreState {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &["Up", "Offline", "Tombstone"];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = StoreState;

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
                    "Up" => Ok(StoreState::Up),
                    "Offline" => Ok(StoreState::Offline),
                    "Tombstone" => Ok(StoreState::Tombstone),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
