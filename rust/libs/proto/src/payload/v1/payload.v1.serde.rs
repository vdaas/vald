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
impl serde::Serialize for Control {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Control", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Control {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Control;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Control")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Control, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Control {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Control", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for control::CreateIndexRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.pool_size != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Control.CreateIndexRequest", len)?;
        if self.pool_size != 0 {
            struct_ser.serialize_field("poolSize", &self.pool_size)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for control::CreateIndexRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "pool_size",
            "poolSize",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            PoolSize,
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
                            "poolSize" | "pool_size" => Ok(GeneratedField::PoolSize),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = control::CreateIndexRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Control.CreateIndexRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<control::CreateIndexRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut pool_size__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::PoolSize => {
                            if pool_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("poolSize"));
                            }
                            pool_size__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(control::CreateIndexRequest {
                    pool_size: pool_size__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Control.CreateIndexRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Discoverer {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Discoverer", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Discoverer {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Discoverer;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Discoverer")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Discoverer, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Discoverer {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Discoverer", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for discoverer::Request {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.name.is_empty() {
            len += 1;
        }
        if !self.namespace.is_empty() {
            len += 1;
        }
        if !self.node.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Discoverer.Request", len)?;
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.namespace.is_empty() {
            struct_ser.serialize_field("namespace", &self.namespace)?;
        }
        if !self.node.is_empty() {
            struct_ser.serialize_field("node", &self.node)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for discoverer::Request {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "name",
            "namespace",
            "node",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Name,
            Namespace,
            Node,
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
                            "name" => Ok(GeneratedField::Name),
                            "namespace" => Ok(GeneratedField::Namespace),
                            "node" => Ok(GeneratedField::Node),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = discoverer::Request;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Discoverer.Request")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<discoverer::Request, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut name__ = None;
                let mut namespace__ = None;
                let mut node__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Namespace => {
                            if namespace__.is_some() {
                                return Err(serde::de::Error::duplicate_field("namespace"));
                            }
                            namespace__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Node => {
                            if node__.is_some() {
                                return Err(serde::de::Error::duplicate_field("node"));
                            }
                            node__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(discoverer::Request {
                    name: name__.unwrap_or_default(),
                    namespace: namespace__.unwrap_or_default(),
                    node: node__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Discoverer.Request", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Empty {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Empty", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Empty {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Empty;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Empty")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Empty, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Empty {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Empty", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Filter {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Filter", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Filter {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Filter;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Filter")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Filter, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Filter {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Filter", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for filter::Config {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.targets.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Filter.Config", len)?;
        if !self.targets.is_empty() {
            struct_ser.serialize_field("targets", &self.targets)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for filter::Config {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "targets",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Targets,
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
                            "targets" => Ok(GeneratedField::Targets),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = filter::Config;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Filter.Config")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<filter::Config, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut targets__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Targets => {
                            if targets__.is_some() {
                                return Err(serde::de::Error::duplicate_field("targets"));
                            }
                            targets__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(filter::Config {
                    targets: targets__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Filter.Config", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for filter::Target {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.host.is_empty() {
            len += 1;
        }
        if self.port != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Filter.Target", len)?;
        if !self.host.is_empty() {
            struct_ser.serialize_field("host", &self.host)?;
        }
        if self.port != 0 {
            struct_ser.serialize_field("port", &self.port)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for filter::Target {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "host",
            "port",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Host,
            Port,
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
                            "host" => Ok(GeneratedField::Host),
                            "port" => Ok(GeneratedField::Port),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = filter::Target;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Filter.Target")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<filter::Target, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut host__ = None;
                let mut port__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Host => {
                            if host__.is_some() {
                                return Err(serde::de::Error::duplicate_field("host"));
                            }
                            host__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Port => {
                            if port__.is_some() {
                                return Err(serde::de::Error::duplicate_field("port"));
                            }
                            port__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(filter::Target {
                    host: host__.unwrap_or_default(),
                    port: port__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Filter.Target", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Flush {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Flush", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Flush {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Flush;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Flush")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Flush, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Flush {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Flush", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for flush::Request {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Flush.Request", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for flush::Request {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = flush::Request;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Flush.Request")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<flush::Request, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(flush::Request {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Flush.Request", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Info {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Info", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Info {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Info;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Info, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Info {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::Annotations {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.annotations.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Annotations", len)?;
        if !self.annotations.is_empty() {
            struct_ser.serialize_field("annotations", &self.annotations)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::Annotations {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "annotations",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Annotations,
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
                            "annotations" => Ok(GeneratedField::Annotations),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::Annotations;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Annotations")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::Annotations, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut annotations__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Annotations => {
                            if annotations__.is_some() {
                                return Err(serde::de::Error::duplicate_field("annotations"));
                            }
                            annotations__ = Some(
                                map_.next_value::<std::collections::HashMap<_, _>>()?
                            );
                        }
                    }
                }
                Ok(info::Annotations {
                    annotations: annotations__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Annotations", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::Cpu {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.limit != 0. {
            len += 1;
        }
        if self.request != 0. {
            len += 1;
        }
        if self.usage != 0. {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.CPU", len)?;
        if self.limit != 0. {
            struct_ser.serialize_field("limit", &self.limit)?;
        }
        if self.request != 0. {
            struct_ser.serialize_field("request", &self.request)?;
        }
        if self.usage != 0. {
            struct_ser.serialize_field("usage", &self.usage)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::Cpu {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "limit",
            "request",
            "usage",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Limit,
            Request,
            Usage,
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
                            "limit" => Ok(GeneratedField::Limit),
                            "request" => Ok(GeneratedField::Request),
                            "usage" => Ok(GeneratedField::Usage),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::Cpu;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.CPU")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::Cpu, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut limit__ = None;
                let mut request__ = None;
                let mut usage__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Limit => {
                            if limit__.is_some() {
                                return Err(serde::de::Error::duplicate_field("limit"));
                            }
                            limit__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Request => {
                            if request__.is_some() {
                                return Err(serde::de::Error::duplicate_field("request"));
                            }
                            request__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Usage => {
                            if usage__.is_some() {
                                return Err(serde::de::Error::duplicate_field("usage"));
                            }
                            usage__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(info::Cpu {
                    limit: limit__.unwrap_or_default(),
                    request: request__.unwrap_or_default(),
                    usage: usage__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.CPU", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::CgroupStats {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.cpu_limit_cores != 0. {
            len += 1;
        }
        if self.cpu_usage_cores != 0. {
            len += 1;
        }
        if self.memory_limit_bytes != 0 {
            len += 1;
        }
        if self.memory_usage_bytes != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.CgroupStats", len)?;
        if self.cpu_limit_cores != 0. {
            struct_ser.serialize_field("cpuLimitCores", &self.cpu_limit_cores)?;
        }
        if self.cpu_usage_cores != 0. {
            struct_ser.serialize_field("cpuUsageCores", &self.cpu_usage_cores)?;
        }
        if self.memory_limit_bytes != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("memoryLimitBytes", ToString::to_string(&self.memory_limit_bytes).as_str())?;
        }
        if self.memory_usage_bytes != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("memoryUsageBytes", ToString::to_string(&self.memory_usage_bytes).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::CgroupStats {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "cpu_limit_cores",
            "cpuLimitCores",
            "cpu_usage_cores",
            "cpuUsageCores",
            "memory_limit_bytes",
            "memoryLimitBytes",
            "memory_usage_bytes",
            "memoryUsageBytes",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            CpuLimitCores,
            CpuUsageCores,
            MemoryLimitBytes,
            MemoryUsageBytes,
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
                            "cpuLimitCores" | "cpu_limit_cores" => Ok(GeneratedField::CpuLimitCores),
                            "cpuUsageCores" | "cpu_usage_cores" => Ok(GeneratedField::CpuUsageCores),
                            "memoryLimitBytes" | "memory_limit_bytes" => Ok(GeneratedField::MemoryLimitBytes),
                            "memoryUsageBytes" | "memory_usage_bytes" => Ok(GeneratedField::MemoryUsageBytes),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::CgroupStats;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.CgroupStats")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::CgroupStats, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut cpu_limit_cores__ = None;
                let mut cpu_usage_cores__ = None;
                let mut memory_limit_bytes__ = None;
                let mut memory_usage_bytes__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::CpuLimitCores => {
                            if cpu_limit_cores__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cpuLimitCores"));
                            }
                            cpu_limit_cores__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::CpuUsageCores => {
                            if cpu_usage_cores__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cpuUsageCores"));
                            }
                            cpu_usage_cores__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MemoryLimitBytes => {
                            if memory_limit_bytes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("memoryLimitBytes"));
                            }
                            memory_limit_bytes__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MemoryUsageBytes => {
                            if memory_usage_bytes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("memoryUsageBytes"));
                            }
                            memory_usage_bytes__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(info::CgroupStats {
                    cpu_limit_cores: cpu_limit_cores__.unwrap_or_default(),
                    cpu_usage_cores: cpu_usage_cores__.unwrap_or_default(),
                    memory_limit_bytes: memory_limit_bytes__.unwrap_or_default(),
                    memory_usage_bytes: memory_usage_bytes__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.CgroupStats", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::IPs {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.ip.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.IPs", len)?;
        if !self.ip.is_empty() {
            struct_ser.serialize_field("ip", &self.ip)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::IPs {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "ip",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Ip,
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
                            "ip" => Ok(GeneratedField::Ip),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::IPs;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.IPs")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::IPs, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut ip__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Ip => {
                            if ip__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ip"));
                            }
                            ip__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(info::IPs {
                    ip: ip__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.IPs", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::Index {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Info.Index", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::Index {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::Index;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Index")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::Index, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(info::Index {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Index", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::index::Count {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.stored != 0 {
            len += 1;
        }
        if self.uncommitted != 0 {
            len += 1;
        }
        if self.indexing {
            len += 1;
        }
        if self.saving {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Index.Count", len)?;
        if self.stored != 0 {
            struct_ser.serialize_field("stored", &self.stored)?;
        }
        if self.uncommitted != 0 {
            struct_ser.serialize_field("uncommitted", &self.uncommitted)?;
        }
        if self.indexing {
            struct_ser.serialize_field("indexing", &self.indexing)?;
        }
        if self.saving {
            struct_ser.serialize_field("saving", &self.saving)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::index::Count {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "stored",
            "uncommitted",
            "indexing",
            "saving",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Stored,
            Uncommitted,
            Indexing,
            Saving,
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
                            "stored" => Ok(GeneratedField::Stored),
                            "uncommitted" => Ok(GeneratedField::Uncommitted),
                            "indexing" => Ok(GeneratedField::Indexing),
                            "saving" => Ok(GeneratedField::Saving),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::index::Count;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Index.Count")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::index::Count, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut stored__ = None;
                let mut uncommitted__ = None;
                let mut indexing__ = None;
                let mut saving__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Stored => {
                            if stored__.is_some() {
                                return Err(serde::de::Error::duplicate_field("stored"));
                            }
                            stored__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Uncommitted => {
                            if uncommitted__.is_some() {
                                return Err(serde::de::Error::duplicate_field("uncommitted"));
                            }
                            uncommitted__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Indexing => {
                            if indexing__.is_some() {
                                return Err(serde::de::Error::duplicate_field("indexing"));
                            }
                            indexing__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Saving => {
                            if saving__.is_some() {
                                return Err(serde::de::Error::duplicate_field("saving"));
                            }
                            saving__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(info::index::Count {
                    stored: stored__.unwrap_or_default(),
                    uncommitted: uncommitted__.unwrap_or_default(),
                    indexing: indexing__.unwrap_or_default(),
                    saving: saving__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Index.Count", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::index::Detail {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.counts.is_empty() {
            len += 1;
        }
        if self.replica != 0 {
            len += 1;
        }
        if self.live_agents != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Index.Detail", len)?;
        if !self.counts.is_empty() {
            struct_ser.serialize_field("counts", &self.counts)?;
        }
        if self.replica != 0 {
            struct_ser.serialize_field("replica", &self.replica)?;
        }
        if self.live_agents != 0 {
            struct_ser.serialize_field("liveAgents", &self.live_agents)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::index::Detail {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "counts",
            "replica",
            "live_agents",
            "liveAgents",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Counts,
            Replica,
            LiveAgents,
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
                            "counts" => Ok(GeneratedField::Counts),
                            "replica" => Ok(GeneratedField::Replica),
                            "liveAgents" | "live_agents" => Ok(GeneratedField::LiveAgents),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::index::Detail;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Index.Detail")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::index::Detail, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut counts__ = None;
                let mut replica__ = None;
                let mut live_agents__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Counts => {
                            if counts__.is_some() {
                                return Err(serde::de::Error::duplicate_field("counts"));
                            }
                            counts__ = Some(
                                map_.next_value::<std::collections::HashMap<_, _>>()?
                            );
                        }
                        GeneratedField::Replica => {
                            if replica__.is_some() {
                                return Err(serde::de::Error::duplicate_field("replica"));
                            }
                            replica__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::LiveAgents => {
                            if live_agents__.is_some() {
                                return Err(serde::de::Error::duplicate_field("liveAgents"));
                            }
                            live_agents__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(info::index::Detail {
                    counts: counts__.unwrap_or_default(),
                    replica: replica__.unwrap_or_default(),
                    live_agents: live_agents__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Index.Detail", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::index::Property {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.dimension != 0 {
            len += 1;
        }
        if self.thread_pool_size != 0 {
            len += 1;
        }
        if !self.object_type.is_empty() {
            len += 1;
        }
        if !self.distance_type.is_empty() {
            len += 1;
        }
        if !self.index_type.is_empty() {
            len += 1;
        }
        if !self.database_type.is_empty() {
            len += 1;
        }
        if !self.object_alignment.is_empty() {
            len += 1;
        }
        if self.path_adjustment_interval != 0 {
            len += 1;
        }
        if self.graph_shared_memory_size != 0 {
            len += 1;
        }
        if self.tree_shared_memory_size != 0 {
            len += 1;
        }
        if self.object_shared_memory_size != 0 {
            len += 1;
        }
        if self.prefetch_offset != 0 {
            len += 1;
        }
        if self.prefetch_size != 0 {
            len += 1;
        }
        if !self.accuracy_table.is_empty() {
            len += 1;
        }
        if !self.search_type.is_empty() {
            len += 1;
        }
        if self.max_magnitude != 0. {
            len += 1;
        }
        if self.n_of_neighbors_for_insertion_order != 0 {
            len += 1;
        }
        if self.epsilon_for_insertion_order != 0. {
            len += 1;
        }
        if !self.refinement_object_type.is_empty() {
            len += 1;
        }
        if self.truncation_threshold != 0 {
            len += 1;
        }
        if self.edge_size_for_creation != 0 {
            len += 1;
        }
        if self.edge_size_for_search != 0 {
            len += 1;
        }
        if self.edge_size_limit_for_creation != 0 {
            len += 1;
        }
        if self.insertion_radius_coefficient != 0. {
            len += 1;
        }
        if self.seed_size != 0 {
            len += 1;
        }
        if !self.seed_type.is_empty() {
            len += 1;
        }
        if self.truncation_thread_pool_size != 0 {
            len += 1;
        }
        if self.batch_size_for_creation != 0 {
            len += 1;
        }
        if !self.graph_type.is_empty() {
            len += 1;
        }
        if self.dynamic_edge_size_base != 0 {
            len += 1;
        }
        if self.dynamic_edge_size_rate != 0 {
            len += 1;
        }
        if self.build_time_limit != 0. {
            len += 1;
        }
        if self.outgoing_edge != 0 {
            len += 1;
        }
        if self.incoming_edge != 0 {
            len += 1;
        }
        if self.epsilon_for_creation != 0. {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Index.Property", len)?;
        if self.dimension != 0 {
            struct_ser.serialize_field("dimension", &self.dimension)?;
        }
        if self.thread_pool_size != 0 {
            struct_ser.serialize_field("threadPoolSize", &self.thread_pool_size)?;
        }
        if !self.object_type.is_empty() {
            struct_ser.serialize_field("objectType", &self.object_type)?;
        }
        if !self.distance_type.is_empty() {
            struct_ser.serialize_field("distanceType", &self.distance_type)?;
        }
        if !self.index_type.is_empty() {
            struct_ser.serialize_field("indexType", &self.index_type)?;
        }
        if !self.database_type.is_empty() {
            struct_ser.serialize_field("databaseType", &self.database_type)?;
        }
        if !self.object_alignment.is_empty() {
            struct_ser.serialize_field("objectAlignment", &self.object_alignment)?;
        }
        if self.path_adjustment_interval != 0 {
            struct_ser.serialize_field("pathAdjustmentInterval", &self.path_adjustment_interval)?;
        }
        if self.graph_shared_memory_size != 0 {
            struct_ser.serialize_field("graphSharedMemorySize", &self.graph_shared_memory_size)?;
        }
        if self.tree_shared_memory_size != 0 {
            struct_ser.serialize_field("treeSharedMemorySize", &self.tree_shared_memory_size)?;
        }
        if self.object_shared_memory_size != 0 {
            struct_ser.serialize_field("objectSharedMemorySize", &self.object_shared_memory_size)?;
        }
        if self.prefetch_offset != 0 {
            struct_ser.serialize_field("prefetchOffset", &self.prefetch_offset)?;
        }
        if self.prefetch_size != 0 {
            struct_ser.serialize_field("prefetchSize", &self.prefetch_size)?;
        }
        if !self.accuracy_table.is_empty() {
            struct_ser.serialize_field("accuracyTable", &self.accuracy_table)?;
        }
        if !self.search_type.is_empty() {
            struct_ser.serialize_field("searchType", &self.search_type)?;
        }
        if self.max_magnitude != 0. {
            struct_ser.serialize_field("maxMagnitude", &self.max_magnitude)?;
        }
        if self.n_of_neighbors_for_insertion_order != 0 {
            struct_ser.serialize_field("nOfNeighborsForInsertionOrder", &self.n_of_neighbors_for_insertion_order)?;
        }
        if self.epsilon_for_insertion_order != 0. {
            struct_ser.serialize_field("epsilonForInsertionOrder", &self.epsilon_for_insertion_order)?;
        }
        if !self.refinement_object_type.is_empty() {
            struct_ser.serialize_field("refinementObjectType", &self.refinement_object_type)?;
        }
        if self.truncation_threshold != 0 {
            struct_ser.serialize_field("truncationThreshold", &self.truncation_threshold)?;
        }
        if self.edge_size_for_creation != 0 {
            struct_ser.serialize_field("edgeSizeForCreation", &self.edge_size_for_creation)?;
        }
        if self.edge_size_for_search != 0 {
            struct_ser.serialize_field("edgeSizeForSearch", &self.edge_size_for_search)?;
        }
        if self.edge_size_limit_for_creation != 0 {
            struct_ser.serialize_field("edgeSizeLimitForCreation", &self.edge_size_limit_for_creation)?;
        }
        if self.insertion_radius_coefficient != 0. {
            struct_ser.serialize_field("insertionRadiusCoefficient", &self.insertion_radius_coefficient)?;
        }
        if self.seed_size != 0 {
            struct_ser.serialize_field("seedSize", &self.seed_size)?;
        }
        if !self.seed_type.is_empty() {
            struct_ser.serialize_field("seedType", &self.seed_type)?;
        }
        if self.truncation_thread_pool_size != 0 {
            struct_ser.serialize_field("truncationThreadPoolSize", &self.truncation_thread_pool_size)?;
        }
        if self.batch_size_for_creation != 0 {
            struct_ser.serialize_field("batchSizeForCreation", &self.batch_size_for_creation)?;
        }
        if !self.graph_type.is_empty() {
            struct_ser.serialize_field("graphType", &self.graph_type)?;
        }
        if self.dynamic_edge_size_base != 0 {
            struct_ser.serialize_field("dynamicEdgeSizeBase", &self.dynamic_edge_size_base)?;
        }
        if self.dynamic_edge_size_rate != 0 {
            struct_ser.serialize_field("dynamicEdgeSizeRate", &self.dynamic_edge_size_rate)?;
        }
        if self.build_time_limit != 0. {
            struct_ser.serialize_field("buildTimeLimit", &self.build_time_limit)?;
        }
        if self.outgoing_edge != 0 {
            struct_ser.serialize_field("outgoingEdge", &self.outgoing_edge)?;
        }
        if self.incoming_edge != 0 {
            struct_ser.serialize_field("incomingEdge", &self.incoming_edge)?;
        }
        if self.epsilon_for_creation != 0. {
            struct_ser.serialize_field("epsilonForCreation", &self.epsilon_for_creation)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::index::Property {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "dimension",
            "thread_pool_size",
            "threadPoolSize",
            "object_type",
            "objectType",
            "distance_type",
            "distanceType",
            "index_type",
            "indexType",
            "database_type",
            "databaseType",
            "object_alignment",
            "objectAlignment",
            "path_adjustment_interval",
            "pathAdjustmentInterval",
            "graph_shared_memory_size",
            "graphSharedMemorySize",
            "tree_shared_memory_size",
            "treeSharedMemorySize",
            "object_shared_memory_size",
            "objectSharedMemorySize",
            "prefetch_offset",
            "prefetchOffset",
            "prefetch_size",
            "prefetchSize",
            "accuracy_table",
            "accuracyTable",
            "search_type",
            "searchType",
            "max_magnitude",
            "maxMagnitude",
            "n_of_neighbors_for_insertion_order",
            "nOfNeighborsForInsertionOrder",
            "epsilon_for_insertion_order",
            "epsilonForInsertionOrder",
            "refinement_object_type",
            "refinementObjectType",
            "truncation_threshold",
            "truncationThreshold",
            "edge_size_for_creation",
            "edgeSizeForCreation",
            "edge_size_for_search",
            "edgeSizeForSearch",
            "edge_size_limit_for_creation",
            "edgeSizeLimitForCreation",
            "insertion_radius_coefficient",
            "insertionRadiusCoefficient",
            "seed_size",
            "seedSize",
            "seed_type",
            "seedType",
            "truncation_thread_pool_size",
            "truncationThreadPoolSize",
            "batch_size_for_creation",
            "batchSizeForCreation",
            "graph_type",
            "graphType",
            "dynamic_edge_size_base",
            "dynamicEdgeSizeBase",
            "dynamic_edge_size_rate",
            "dynamicEdgeSizeRate",
            "build_time_limit",
            "buildTimeLimit",
            "outgoing_edge",
            "outgoingEdge",
            "incoming_edge",
            "incomingEdge",
            "epsilon_for_creation",
            "epsilonForCreation",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Dimension,
            ThreadPoolSize,
            ObjectType,
            DistanceType,
            IndexType,
            DatabaseType,
            ObjectAlignment,
            PathAdjustmentInterval,
            GraphSharedMemorySize,
            TreeSharedMemorySize,
            ObjectSharedMemorySize,
            PrefetchOffset,
            PrefetchSize,
            AccuracyTable,
            SearchType,
            MaxMagnitude,
            NOfNeighborsForInsertionOrder,
            EpsilonForInsertionOrder,
            RefinementObjectType,
            TruncationThreshold,
            EdgeSizeForCreation,
            EdgeSizeForSearch,
            EdgeSizeLimitForCreation,
            InsertionRadiusCoefficient,
            SeedSize,
            SeedType,
            TruncationThreadPoolSize,
            BatchSizeForCreation,
            GraphType,
            DynamicEdgeSizeBase,
            DynamicEdgeSizeRate,
            BuildTimeLimit,
            OutgoingEdge,
            IncomingEdge,
            EpsilonForCreation,
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
                            "dimension" => Ok(GeneratedField::Dimension),
                            "threadPoolSize" | "thread_pool_size" => Ok(GeneratedField::ThreadPoolSize),
                            "objectType" | "object_type" => Ok(GeneratedField::ObjectType),
                            "distanceType" | "distance_type" => Ok(GeneratedField::DistanceType),
                            "indexType" | "index_type" => Ok(GeneratedField::IndexType),
                            "databaseType" | "database_type" => Ok(GeneratedField::DatabaseType),
                            "objectAlignment" | "object_alignment" => Ok(GeneratedField::ObjectAlignment),
                            "pathAdjustmentInterval" | "path_adjustment_interval" => Ok(GeneratedField::PathAdjustmentInterval),
                            "graphSharedMemorySize" | "graph_shared_memory_size" => Ok(GeneratedField::GraphSharedMemorySize),
                            "treeSharedMemorySize" | "tree_shared_memory_size" => Ok(GeneratedField::TreeSharedMemorySize),
                            "objectSharedMemorySize" | "object_shared_memory_size" => Ok(GeneratedField::ObjectSharedMemorySize),
                            "prefetchOffset" | "prefetch_offset" => Ok(GeneratedField::PrefetchOffset),
                            "prefetchSize" | "prefetch_size" => Ok(GeneratedField::PrefetchSize),
                            "accuracyTable" | "accuracy_table" => Ok(GeneratedField::AccuracyTable),
                            "searchType" | "search_type" => Ok(GeneratedField::SearchType),
                            "maxMagnitude" | "max_magnitude" => Ok(GeneratedField::MaxMagnitude),
                            "nOfNeighborsForInsertionOrder" | "n_of_neighbors_for_insertion_order" => Ok(GeneratedField::NOfNeighborsForInsertionOrder),
                            "epsilonForInsertionOrder" | "epsilon_for_insertion_order" => Ok(GeneratedField::EpsilonForInsertionOrder),
                            "refinementObjectType" | "refinement_object_type" => Ok(GeneratedField::RefinementObjectType),
                            "truncationThreshold" | "truncation_threshold" => Ok(GeneratedField::TruncationThreshold),
                            "edgeSizeForCreation" | "edge_size_for_creation" => Ok(GeneratedField::EdgeSizeForCreation),
                            "edgeSizeForSearch" | "edge_size_for_search" => Ok(GeneratedField::EdgeSizeForSearch),
                            "edgeSizeLimitForCreation" | "edge_size_limit_for_creation" => Ok(GeneratedField::EdgeSizeLimitForCreation),
                            "insertionRadiusCoefficient" | "insertion_radius_coefficient" => Ok(GeneratedField::InsertionRadiusCoefficient),
                            "seedSize" | "seed_size" => Ok(GeneratedField::SeedSize),
                            "seedType" | "seed_type" => Ok(GeneratedField::SeedType),
                            "truncationThreadPoolSize" | "truncation_thread_pool_size" => Ok(GeneratedField::TruncationThreadPoolSize),
                            "batchSizeForCreation" | "batch_size_for_creation" => Ok(GeneratedField::BatchSizeForCreation),
                            "graphType" | "graph_type" => Ok(GeneratedField::GraphType),
                            "dynamicEdgeSizeBase" | "dynamic_edge_size_base" => Ok(GeneratedField::DynamicEdgeSizeBase),
                            "dynamicEdgeSizeRate" | "dynamic_edge_size_rate" => Ok(GeneratedField::DynamicEdgeSizeRate),
                            "buildTimeLimit" | "build_time_limit" => Ok(GeneratedField::BuildTimeLimit),
                            "outgoingEdge" | "outgoing_edge" => Ok(GeneratedField::OutgoingEdge),
                            "incomingEdge" | "incoming_edge" => Ok(GeneratedField::IncomingEdge),
                            "epsilonForCreation" | "epsilon_for_creation" => Ok(GeneratedField::EpsilonForCreation),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::index::Property;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Index.Property")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::index::Property, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut dimension__ = None;
                let mut thread_pool_size__ = None;
                let mut object_type__ = None;
                let mut distance_type__ = None;
                let mut index_type__ = None;
                let mut database_type__ = None;
                let mut object_alignment__ = None;
                let mut path_adjustment_interval__ = None;
                let mut graph_shared_memory_size__ = None;
                let mut tree_shared_memory_size__ = None;
                let mut object_shared_memory_size__ = None;
                let mut prefetch_offset__ = None;
                let mut prefetch_size__ = None;
                let mut accuracy_table__ = None;
                let mut search_type__ = None;
                let mut max_magnitude__ = None;
                let mut n_of_neighbors_for_insertion_order__ = None;
                let mut epsilon_for_insertion_order__ = None;
                let mut refinement_object_type__ = None;
                let mut truncation_threshold__ = None;
                let mut edge_size_for_creation__ = None;
                let mut edge_size_for_search__ = None;
                let mut edge_size_limit_for_creation__ = None;
                let mut insertion_radius_coefficient__ = None;
                let mut seed_size__ = None;
                let mut seed_type__ = None;
                let mut truncation_thread_pool_size__ = None;
                let mut batch_size_for_creation__ = None;
                let mut graph_type__ = None;
                let mut dynamic_edge_size_base__ = None;
                let mut dynamic_edge_size_rate__ = None;
                let mut build_time_limit__ = None;
                let mut outgoing_edge__ = None;
                let mut incoming_edge__ = None;
                let mut epsilon_for_creation__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Dimension => {
                            if dimension__.is_some() {
                                return Err(serde::de::Error::duplicate_field("dimension"));
                            }
                            dimension__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ThreadPoolSize => {
                            if thread_pool_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("threadPoolSize"));
                            }
                            thread_pool_size__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ObjectType => {
                            if object_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("objectType"));
                            }
                            object_type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::DistanceType => {
                            if distance_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("distanceType"));
                            }
                            distance_type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::IndexType => {
                            if index_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("indexType"));
                            }
                            index_type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::DatabaseType => {
                            if database_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("databaseType"));
                            }
                            database_type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ObjectAlignment => {
                            if object_alignment__.is_some() {
                                return Err(serde::de::Error::duplicate_field("objectAlignment"));
                            }
                            object_alignment__ = Some(map_.next_value()?);
                        }
                        GeneratedField::PathAdjustmentInterval => {
                            if path_adjustment_interval__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pathAdjustmentInterval"));
                            }
                            path_adjustment_interval__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::GraphSharedMemorySize => {
                            if graph_shared_memory_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("graphSharedMemorySize"));
                            }
                            graph_shared_memory_size__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::TreeSharedMemorySize => {
                            if tree_shared_memory_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("treeSharedMemorySize"));
                            }
                            tree_shared_memory_size__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ObjectSharedMemorySize => {
                            if object_shared_memory_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("objectSharedMemorySize"));
                            }
                            object_shared_memory_size__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::PrefetchOffset => {
                            if prefetch_offset__.is_some() {
                                return Err(serde::de::Error::duplicate_field("prefetchOffset"));
                            }
                            prefetch_offset__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::PrefetchSize => {
                            if prefetch_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("prefetchSize"));
                            }
                            prefetch_size__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::AccuracyTable => {
                            if accuracy_table__.is_some() {
                                return Err(serde::de::Error::duplicate_field("accuracyTable"));
                            }
                            accuracy_table__ = Some(map_.next_value()?);
                        }
                        GeneratedField::SearchType => {
                            if search_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("searchType"));
                            }
                            search_type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::MaxMagnitude => {
                            if max_magnitude__.is_some() {
                                return Err(serde::de::Error::duplicate_field("maxMagnitude"));
                            }
                            max_magnitude__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::NOfNeighborsForInsertionOrder => {
                            if n_of_neighbors_for_insertion_order__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nOfNeighborsForInsertionOrder"));
                            }
                            n_of_neighbors_for_insertion_order__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::EpsilonForInsertionOrder => {
                            if epsilon_for_insertion_order__.is_some() {
                                return Err(serde::de::Error::duplicate_field("epsilonForInsertionOrder"));
                            }
                            epsilon_for_insertion_order__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::RefinementObjectType => {
                            if refinement_object_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("refinementObjectType"));
                            }
                            refinement_object_type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TruncationThreshold => {
                            if truncation_threshold__.is_some() {
                                return Err(serde::de::Error::duplicate_field("truncationThreshold"));
                            }
                            truncation_threshold__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::EdgeSizeForCreation => {
                            if edge_size_for_creation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("edgeSizeForCreation"));
                            }
                            edge_size_for_creation__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::EdgeSizeForSearch => {
                            if edge_size_for_search__.is_some() {
                                return Err(serde::de::Error::duplicate_field("edgeSizeForSearch"));
                            }
                            edge_size_for_search__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::EdgeSizeLimitForCreation => {
                            if edge_size_limit_for_creation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("edgeSizeLimitForCreation"));
                            }
                            edge_size_limit_for_creation__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::InsertionRadiusCoefficient => {
                            if insertion_radius_coefficient__.is_some() {
                                return Err(serde::de::Error::duplicate_field("insertionRadiusCoefficient"));
                            }
                            insertion_radius_coefficient__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::SeedSize => {
                            if seed_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("seedSize"));
                            }
                            seed_size__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::SeedType => {
                            if seed_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("seedType"));
                            }
                            seed_type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::TruncationThreadPoolSize => {
                            if truncation_thread_pool_size__.is_some() {
                                return Err(serde::de::Error::duplicate_field("truncationThreadPoolSize"));
                            }
                            truncation_thread_pool_size__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::BatchSizeForCreation => {
                            if batch_size_for_creation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("batchSizeForCreation"));
                            }
                            batch_size_for_creation__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::GraphType => {
                            if graph_type__.is_some() {
                                return Err(serde::de::Error::duplicate_field("graphType"));
                            }
                            graph_type__ = Some(map_.next_value()?);
                        }
                        GeneratedField::DynamicEdgeSizeBase => {
                            if dynamic_edge_size_base__.is_some() {
                                return Err(serde::de::Error::duplicate_field("dynamicEdgeSizeBase"));
                            }
                            dynamic_edge_size_base__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::DynamicEdgeSizeRate => {
                            if dynamic_edge_size_rate__.is_some() {
                                return Err(serde::de::Error::duplicate_field("dynamicEdgeSizeRate"));
                            }
                            dynamic_edge_size_rate__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::BuildTimeLimit => {
                            if build_time_limit__.is_some() {
                                return Err(serde::de::Error::duplicate_field("buildTimeLimit"));
                            }
                            build_time_limit__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::OutgoingEdge => {
                            if outgoing_edge__.is_some() {
                                return Err(serde::de::Error::duplicate_field("outgoingEdge"));
                            }
                            outgoing_edge__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::IncomingEdge => {
                            if incoming_edge__.is_some() {
                                return Err(serde::de::Error::duplicate_field("incomingEdge"));
                            }
                            incoming_edge__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::EpsilonForCreation => {
                            if epsilon_for_creation__.is_some() {
                                return Err(serde::de::Error::duplicate_field("epsilonForCreation"));
                            }
                            epsilon_for_creation__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(info::index::Property {
                    dimension: dimension__.unwrap_or_default(),
                    thread_pool_size: thread_pool_size__.unwrap_or_default(),
                    object_type: object_type__.unwrap_or_default(),
                    distance_type: distance_type__.unwrap_or_default(),
                    index_type: index_type__.unwrap_or_default(),
                    database_type: database_type__.unwrap_or_default(),
                    object_alignment: object_alignment__.unwrap_or_default(),
                    path_adjustment_interval: path_adjustment_interval__.unwrap_or_default(),
                    graph_shared_memory_size: graph_shared_memory_size__.unwrap_or_default(),
                    tree_shared_memory_size: tree_shared_memory_size__.unwrap_or_default(),
                    object_shared_memory_size: object_shared_memory_size__.unwrap_or_default(),
                    prefetch_offset: prefetch_offset__.unwrap_or_default(),
                    prefetch_size: prefetch_size__.unwrap_or_default(),
                    accuracy_table: accuracy_table__.unwrap_or_default(),
                    search_type: search_type__.unwrap_or_default(),
                    max_magnitude: max_magnitude__.unwrap_or_default(),
                    n_of_neighbors_for_insertion_order: n_of_neighbors_for_insertion_order__.unwrap_or_default(),
                    epsilon_for_insertion_order: epsilon_for_insertion_order__.unwrap_or_default(),
                    refinement_object_type: refinement_object_type__.unwrap_or_default(),
                    truncation_threshold: truncation_threshold__.unwrap_or_default(),
                    edge_size_for_creation: edge_size_for_creation__.unwrap_or_default(),
                    edge_size_for_search: edge_size_for_search__.unwrap_or_default(),
                    edge_size_limit_for_creation: edge_size_limit_for_creation__.unwrap_or_default(),
                    insertion_radius_coefficient: insertion_radius_coefficient__.unwrap_or_default(),
                    seed_size: seed_size__.unwrap_or_default(),
                    seed_type: seed_type__.unwrap_or_default(),
                    truncation_thread_pool_size: truncation_thread_pool_size__.unwrap_or_default(),
                    batch_size_for_creation: batch_size_for_creation__.unwrap_or_default(),
                    graph_type: graph_type__.unwrap_or_default(),
                    dynamic_edge_size_base: dynamic_edge_size_base__.unwrap_or_default(),
                    dynamic_edge_size_rate: dynamic_edge_size_rate__.unwrap_or_default(),
                    build_time_limit: build_time_limit__.unwrap_or_default(),
                    outgoing_edge: outgoing_edge__.unwrap_or_default(),
                    incoming_edge: incoming_edge__.unwrap_or_default(),
                    epsilon_for_creation: epsilon_for_creation__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Index.Property", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::index::PropertyDetail {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.details.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Index.PropertyDetail", len)?;
        if !self.details.is_empty() {
            struct_ser.serialize_field("details", &self.details)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::index::PropertyDetail {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "details",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Details,
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
                            "details" => Ok(GeneratedField::Details),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::index::PropertyDetail;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Index.PropertyDetail")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::index::PropertyDetail, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut details__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Details => {
                            if details__.is_some() {
                                return Err(serde::de::Error::duplicate_field("details"));
                            }
                            details__ = Some(
                                map_.next_value::<std::collections::HashMap<_, _>>()?
                            );
                        }
                    }
                }
                Ok(info::index::PropertyDetail {
                    details: details__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Index.PropertyDetail", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::index::Statistics {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.valid {
            len += 1;
        }
        if self.median_indegree != 0 {
            len += 1;
        }
        if self.median_outdegree != 0 {
            len += 1;
        }
        if self.max_number_of_indegree != 0 {
            len += 1;
        }
        if self.max_number_of_outdegree != 0 {
            len += 1;
        }
        if self.min_number_of_indegree != 0 {
            len += 1;
        }
        if self.min_number_of_outdegree != 0 {
            len += 1;
        }
        if self.mode_indegree != 0 {
            len += 1;
        }
        if self.mode_outdegree != 0 {
            len += 1;
        }
        if self.nodes_skipped_for_10_edges != 0 {
            len += 1;
        }
        if self.nodes_skipped_for_indegree_distance != 0 {
            len += 1;
        }
        if self.number_of_edges != 0 {
            len += 1;
        }
        if self.number_of_indexed_objects != 0 {
            len += 1;
        }
        if self.number_of_nodes != 0 {
            len += 1;
        }
        if self.number_of_nodes_without_edges != 0 {
            len += 1;
        }
        if self.number_of_nodes_without_indegree != 0 {
            len += 1;
        }
        if self.number_of_objects != 0 {
            len += 1;
        }
        if self.number_of_removed_objects != 0 {
            len += 1;
        }
        if self.size_of_object_repository != 0 {
            len += 1;
        }
        if self.size_of_refinement_object_repository != 0 {
            len += 1;
        }
        if self.variance_of_indegree != 0. {
            len += 1;
        }
        if self.variance_of_outdegree != 0. {
            len += 1;
        }
        if self.mean_edge_length != 0. {
            len += 1;
        }
        if self.mean_edge_length_for_10_edges != 0. {
            len += 1;
        }
        if self.mean_indegree_distance_for_10_edges != 0. {
            len += 1;
        }
        if self.mean_number_of_edges_per_node != 0. {
            len += 1;
        }
        if self.c1_indegree != 0. {
            len += 1;
        }
        if self.c5_indegree != 0. {
            len += 1;
        }
        if self.c95_outdegree != 0. {
            len += 1;
        }
        if self.c99_outdegree != 0. {
            len += 1;
        }
        if !self.indegree_count.is_empty() {
            len += 1;
        }
        if !self.outdegree_histogram.is_empty() {
            len += 1;
        }
        if !self.indegree_histogram.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Index.Statistics", len)?;
        if self.valid {
            struct_ser.serialize_field("valid", &self.valid)?;
        }
        if self.median_indegree != 0 {
            struct_ser.serialize_field("medianIndegree", &self.median_indegree)?;
        }
        if self.median_outdegree != 0 {
            struct_ser.serialize_field("medianOutdegree", &self.median_outdegree)?;
        }
        if self.max_number_of_indegree != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("maxNumberOfIndegree", ToString::to_string(&self.max_number_of_indegree).as_str())?;
        }
        if self.max_number_of_outdegree != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("maxNumberOfOutdegree", ToString::to_string(&self.max_number_of_outdegree).as_str())?;
        }
        if self.min_number_of_indegree != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("minNumberOfIndegree", ToString::to_string(&self.min_number_of_indegree).as_str())?;
        }
        if self.min_number_of_outdegree != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("minNumberOfOutdegree", ToString::to_string(&self.min_number_of_outdegree).as_str())?;
        }
        if self.mode_indegree != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("modeIndegree", ToString::to_string(&self.mode_indegree).as_str())?;
        }
        if self.mode_outdegree != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("modeOutdegree", ToString::to_string(&self.mode_outdegree).as_str())?;
        }
        if self.nodes_skipped_for_10_edges != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("nodesSkippedFor10Edges", ToString::to_string(&self.nodes_skipped_for_10_edges).as_str())?;
        }
        if self.nodes_skipped_for_indegree_distance != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("nodesSkippedForIndegreeDistance", ToString::to_string(&self.nodes_skipped_for_indegree_distance).as_str())?;
        }
        if self.number_of_edges != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("numberOfEdges", ToString::to_string(&self.number_of_edges).as_str())?;
        }
        if self.number_of_indexed_objects != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("numberOfIndexedObjects", ToString::to_string(&self.number_of_indexed_objects).as_str())?;
        }
        if self.number_of_nodes != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("numberOfNodes", ToString::to_string(&self.number_of_nodes).as_str())?;
        }
        if self.number_of_nodes_without_edges != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("numberOfNodesWithoutEdges", ToString::to_string(&self.number_of_nodes_without_edges).as_str())?;
        }
        if self.number_of_nodes_without_indegree != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("numberOfNodesWithoutIndegree", ToString::to_string(&self.number_of_nodes_without_indegree).as_str())?;
        }
        if self.number_of_objects != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("numberOfObjects", ToString::to_string(&self.number_of_objects).as_str())?;
        }
        if self.number_of_removed_objects != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("numberOfRemovedObjects", ToString::to_string(&self.number_of_removed_objects).as_str())?;
        }
        if self.size_of_object_repository != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("sizeOfObjectRepository", ToString::to_string(&self.size_of_object_repository).as_str())?;
        }
        if self.size_of_refinement_object_repository != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("sizeOfRefinementObjectRepository", ToString::to_string(&self.size_of_refinement_object_repository).as_str())?;
        }
        if self.variance_of_indegree != 0. {
            struct_ser.serialize_field("varianceOfIndegree", &self.variance_of_indegree)?;
        }
        if self.variance_of_outdegree != 0. {
            struct_ser.serialize_field("varianceOfOutdegree", &self.variance_of_outdegree)?;
        }
        if self.mean_edge_length != 0. {
            struct_ser.serialize_field("meanEdgeLength", &self.mean_edge_length)?;
        }
        if self.mean_edge_length_for_10_edges != 0. {
            struct_ser.serialize_field("meanEdgeLengthFor10Edges", &self.mean_edge_length_for_10_edges)?;
        }
        if self.mean_indegree_distance_for_10_edges != 0. {
            struct_ser.serialize_field("meanIndegreeDistanceFor10Edges", &self.mean_indegree_distance_for_10_edges)?;
        }
        if self.mean_number_of_edges_per_node != 0. {
            struct_ser.serialize_field("meanNumberOfEdgesPerNode", &self.mean_number_of_edges_per_node)?;
        }
        if self.c1_indegree != 0. {
            struct_ser.serialize_field("c1Indegree", &self.c1_indegree)?;
        }
        if self.c5_indegree != 0. {
            struct_ser.serialize_field("c5Indegree", &self.c5_indegree)?;
        }
        if self.c95_outdegree != 0. {
            struct_ser.serialize_field("c95Outdegree", &self.c95_outdegree)?;
        }
        if self.c99_outdegree != 0. {
            struct_ser.serialize_field("c99Outdegree", &self.c99_outdegree)?;
        }
        if !self.indegree_count.is_empty() {
            struct_ser.serialize_field("indegreeCount", &self.indegree_count.iter().map(ToString::to_string).collect::<Vec<_>>())?;
        }
        if !self.outdegree_histogram.is_empty() {
            struct_ser.serialize_field("outdegreeHistogram", &self.outdegree_histogram.iter().map(ToString::to_string).collect::<Vec<_>>())?;
        }
        if !self.indegree_histogram.is_empty() {
            struct_ser.serialize_field("indegreeHistogram", &self.indegree_histogram.iter().map(ToString::to_string).collect::<Vec<_>>())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::index::Statistics {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "valid",
            "median_indegree",
            "medianIndegree",
            "median_outdegree",
            "medianOutdegree",
            "max_number_of_indegree",
            "maxNumberOfIndegree",
            "max_number_of_outdegree",
            "maxNumberOfOutdegree",
            "min_number_of_indegree",
            "minNumberOfIndegree",
            "min_number_of_outdegree",
            "minNumberOfOutdegree",
            "mode_indegree",
            "modeIndegree",
            "mode_outdegree",
            "modeOutdegree",
            "nodes_skipped_for_10_edges",
            "nodesSkippedFor10Edges",
            "nodes_skipped_for_indegree_distance",
            "nodesSkippedForIndegreeDistance",
            "number_of_edges",
            "numberOfEdges",
            "number_of_indexed_objects",
            "numberOfIndexedObjects",
            "number_of_nodes",
            "numberOfNodes",
            "number_of_nodes_without_edges",
            "numberOfNodesWithoutEdges",
            "number_of_nodes_without_indegree",
            "numberOfNodesWithoutIndegree",
            "number_of_objects",
            "numberOfObjects",
            "number_of_removed_objects",
            "numberOfRemovedObjects",
            "size_of_object_repository",
            "sizeOfObjectRepository",
            "size_of_refinement_object_repository",
            "sizeOfRefinementObjectRepository",
            "variance_of_indegree",
            "varianceOfIndegree",
            "variance_of_outdegree",
            "varianceOfOutdegree",
            "mean_edge_length",
            "meanEdgeLength",
            "mean_edge_length_for_10_edges",
            "meanEdgeLengthFor10Edges",
            "mean_indegree_distance_for_10_edges",
            "meanIndegreeDistanceFor10Edges",
            "mean_number_of_edges_per_node",
            "meanNumberOfEdgesPerNode",
            "c1_indegree",
            "c1Indegree",
            "c5_indegree",
            "c5Indegree",
            "c95_outdegree",
            "c95Outdegree",
            "c99_outdegree",
            "c99Outdegree",
            "indegree_count",
            "indegreeCount",
            "outdegree_histogram",
            "outdegreeHistogram",
            "indegree_histogram",
            "indegreeHistogram",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Valid,
            MedianIndegree,
            MedianOutdegree,
            MaxNumberOfIndegree,
            MaxNumberOfOutdegree,
            MinNumberOfIndegree,
            MinNumberOfOutdegree,
            ModeIndegree,
            ModeOutdegree,
            NodesSkippedFor10Edges,
            NodesSkippedForIndegreeDistance,
            NumberOfEdges,
            NumberOfIndexedObjects,
            NumberOfNodes,
            NumberOfNodesWithoutEdges,
            NumberOfNodesWithoutIndegree,
            NumberOfObjects,
            NumberOfRemovedObjects,
            SizeOfObjectRepository,
            SizeOfRefinementObjectRepository,
            VarianceOfIndegree,
            VarianceOfOutdegree,
            MeanEdgeLength,
            MeanEdgeLengthFor10Edges,
            MeanIndegreeDistanceFor10Edges,
            MeanNumberOfEdgesPerNode,
            C1Indegree,
            C5Indegree,
            C95Outdegree,
            C99Outdegree,
            IndegreeCount,
            OutdegreeHistogram,
            IndegreeHistogram,
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
                            "valid" => Ok(GeneratedField::Valid),
                            "medianIndegree" | "median_indegree" => Ok(GeneratedField::MedianIndegree),
                            "medianOutdegree" | "median_outdegree" => Ok(GeneratedField::MedianOutdegree),
                            "maxNumberOfIndegree" | "max_number_of_indegree" => Ok(GeneratedField::MaxNumberOfIndegree),
                            "maxNumberOfOutdegree" | "max_number_of_outdegree" => Ok(GeneratedField::MaxNumberOfOutdegree),
                            "minNumberOfIndegree" | "min_number_of_indegree" => Ok(GeneratedField::MinNumberOfIndegree),
                            "minNumberOfOutdegree" | "min_number_of_outdegree" => Ok(GeneratedField::MinNumberOfOutdegree),
                            "modeIndegree" | "mode_indegree" => Ok(GeneratedField::ModeIndegree),
                            "modeOutdegree" | "mode_outdegree" => Ok(GeneratedField::ModeOutdegree),
                            "nodesSkippedFor10Edges" | "nodes_skipped_for_10_edges" => Ok(GeneratedField::NodesSkippedFor10Edges),
                            "nodesSkippedForIndegreeDistance" | "nodes_skipped_for_indegree_distance" => Ok(GeneratedField::NodesSkippedForIndegreeDistance),
                            "numberOfEdges" | "number_of_edges" => Ok(GeneratedField::NumberOfEdges),
                            "numberOfIndexedObjects" | "number_of_indexed_objects" => Ok(GeneratedField::NumberOfIndexedObjects),
                            "numberOfNodes" | "number_of_nodes" => Ok(GeneratedField::NumberOfNodes),
                            "numberOfNodesWithoutEdges" | "number_of_nodes_without_edges" => Ok(GeneratedField::NumberOfNodesWithoutEdges),
                            "numberOfNodesWithoutIndegree" | "number_of_nodes_without_indegree" => Ok(GeneratedField::NumberOfNodesWithoutIndegree),
                            "numberOfObjects" | "number_of_objects" => Ok(GeneratedField::NumberOfObjects),
                            "numberOfRemovedObjects" | "number_of_removed_objects" => Ok(GeneratedField::NumberOfRemovedObjects),
                            "sizeOfObjectRepository" | "size_of_object_repository" => Ok(GeneratedField::SizeOfObjectRepository),
                            "sizeOfRefinementObjectRepository" | "size_of_refinement_object_repository" => Ok(GeneratedField::SizeOfRefinementObjectRepository),
                            "varianceOfIndegree" | "variance_of_indegree" => Ok(GeneratedField::VarianceOfIndegree),
                            "varianceOfOutdegree" | "variance_of_outdegree" => Ok(GeneratedField::VarianceOfOutdegree),
                            "meanEdgeLength" | "mean_edge_length" => Ok(GeneratedField::MeanEdgeLength),
                            "meanEdgeLengthFor10Edges" | "mean_edge_length_for_10_edges" => Ok(GeneratedField::MeanEdgeLengthFor10Edges),
                            "meanIndegreeDistanceFor10Edges" | "mean_indegree_distance_for_10_edges" => Ok(GeneratedField::MeanIndegreeDistanceFor10Edges),
                            "meanNumberOfEdgesPerNode" | "mean_number_of_edges_per_node" => Ok(GeneratedField::MeanNumberOfEdgesPerNode),
                            "c1Indegree" | "c1_indegree" => Ok(GeneratedField::C1Indegree),
                            "c5Indegree" | "c5_indegree" => Ok(GeneratedField::C5Indegree),
                            "c95Outdegree" | "c95_outdegree" => Ok(GeneratedField::C95Outdegree),
                            "c99Outdegree" | "c99_outdegree" => Ok(GeneratedField::C99Outdegree),
                            "indegreeCount" | "indegree_count" => Ok(GeneratedField::IndegreeCount),
                            "outdegreeHistogram" | "outdegree_histogram" => Ok(GeneratedField::OutdegreeHistogram),
                            "indegreeHistogram" | "indegree_histogram" => Ok(GeneratedField::IndegreeHistogram),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::index::Statistics;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Index.Statistics")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::index::Statistics, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut valid__ = None;
                let mut median_indegree__ = None;
                let mut median_outdegree__ = None;
                let mut max_number_of_indegree__ = None;
                let mut max_number_of_outdegree__ = None;
                let mut min_number_of_indegree__ = None;
                let mut min_number_of_outdegree__ = None;
                let mut mode_indegree__ = None;
                let mut mode_outdegree__ = None;
                let mut nodes_skipped_for_10_edges__ = None;
                let mut nodes_skipped_for_indegree_distance__ = None;
                let mut number_of_edges__ = None;
                let mut number_of_indexed_objects__ = None;
                let mut number_of_nodes__ = None;
                let mut number_of_nodes_without_edges__ = None;
                let mut number_of_nodes_without_indegree__ = None;
                let mut number_of_objects__ = None;
                let mut number_of_removed_objects__ = None;
                let mut size_of_object_repository__ = None;
                let mut size_of_refinement_object_repository__ = None;
                let mut variance_of_indegree__ = None;
                let mut variance_of_outdegree__ = None;
                let mut mean_edge_length__ = None;
                let mut mean_edge_length_for_10_edges__ = None;
                let mut mean_indegree_distance_for_10_edges__ = None;
                let mut mean_number_of_edges_per_node__ = None;
                let mut c1_indegree__ = None;
                let mut c5_indegree__ = None;
                let mut c95_outdegree__ = None;
                let mut c99_outdegree__ = None;
                let mut indegree_count__ = None;
                let mut outdegree_histogram__ = None;
                let mut indegree_histogram__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Valid => {
                            if valid__.is_some() {
                                return Err(serde::de::Error::duplicate_field("valid"));
                            }
                            valid__ = Some(map_.next_value()?);
                        }
                        GeneratedField::MedianIndegree => {
                            if median_indegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("medianIndegree"));
                            }
                            median_indegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MedianOutdegree => {
                            if median_outdegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("medianOutdegree"));
                            }
                            median_outdegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MaxNumberOfIndegree => {
                            if max_number_of_indegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("maxNumberOfIndegree"));
                            }
                            max_number_of_indegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MaxNumberOfOutdegree => {
                            if max_number_of_outdegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("maxNumberOfOutdegree"));
                            }
                            max_number_of_outdegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MinNumberOfIndegree => {
                            if min_number_of_indegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("minNumberOfIndegree"));
                            }
                            min_number_of_indegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MinNumberOfOutdegree => {
                            if min_number_of_outdegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("minNumberOfOutdegree"));
                            }
                            min_number_of_outdegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ModeIndegree => {
                            if mode_indegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("modeIndegree"));
                            }
                            mode_indegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::ModeOutdegree => {
                            if mode_outdegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("modeOutdegree"));
                            }
                            mode_outdegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::NodesSkippedFor10Edges => {
                            if nodes_skipped_for_10_edges__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nodesSkippedFor10Edges"));
                            }
                            nodes_skipped_for_10_edges__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::NodesSkippedForIndegreeDistance => {
                            if nodes_skipped_for_indegree_distance__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nodesSkippedForIndegreeDistance"));
                            }
                            nodes_skipped_for_indegree_distance__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::NumberOfEdges => {
                            if number_of_edges__.is_some() {
                                return Err(serde::de::Error::duplicate_field("numberOfEdges"));
                            }
                            number_of_edges__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::NumberOfIndexedObjects => {
                            if number_of_indexed_objects__.is_some() {
                                return Err(serde::de::Error::duplicate_field("numberOfIndexedObjects"));
                            }
                            number_of_indexed_objects__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::NumberOfNodes => {
                            if number_of_nodes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("numberOfNodes"));
                            }
                            number_of_nodes__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::NumberOfNodesWithoutEdges => {
                            if number_of_nodes_without_edges__.is_some() {
                                return Err(serde::de::Error::duplicate_field("numberOfNodesWithoutEdges"));
                            }
                            number_of_nodes_without_edges__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::NumberOfNodesWithoutIndegree => {
                            if number_of_nodes_without_indegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("numberOfNodesWithoutIndegree"));
                            }
                            number_of_nodes_without_indegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::NumberOfObjects => {
                            if number_of_objects__.is_some() {
                                return Err(serde::de::Error::duplicate_field("numberOfObjects"));
                            }
                            number_of_objects__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::NumberOfRemovedObjects => {
                            if number_of_removed_objects__.is_some() {
                                return Err(serde::de::Error::duplicate_field("numberOfRemovedObjects"));
                            }
                            number_of_removed_objects__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::SizeOfObjectRepository => {
                            if size_of_object_repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sizeOfObjectRepository"));
                            }
                            size_of_object_repository__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::SizeOfRefinementObjectRepository => {
                            if size_of_refinement_object_repository__.is_some() {
                                return Err(serde::de::Error::duplicate_field("sizeOfRefinementObjectRepository"));
                            }
                            size_of_refinement_object_repository__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::VarianceOfIndegree => {
                            if variance_of_indegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("varianceOfIndegree"));
                            }
                            variance_of_indegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::VarianceOfOutdegree => {
                            if variance_of_outdegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("varianceOfOutdegree"));
                            }
                            variance_of_outdegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MeanEdgeLength => {
                            if mean_edge_length__.is_some() {
                                return Err(serde::de::Error::duplicate_field("meanEdgeLength"));
                            }
                            mean_edge_length__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MeanEdgeLengthFor10Edges => {
                            if mean_edge_length_for_10_edges__.is_some() {
                                return Err(serde::de::Error::duplicate_field("meanEdgeLengthFor10Edges"));
                            }
                            mean_edge_length_for_10_edges__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MeanIndegreeDistanceFor10Edges => {
                            if mean_indegree_distance_for_10_edges__.is_some() {
                                return Err(serde::de::Error::duplicate_field("meanIndegreeDistanceFor10Edges"));
                            }
                            mean_indegree_distance_for_10_edges__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::MeanNumberOfEdgesPerNode => {
                            if mean_number_of_edges_per_node__.is_some() {
                                return Err(serde::de::Error::duplicate_field("meanNumberOfEdgesPerNode"));
                            }
                            mean_number_of_edges_per_node__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::C1Indegree => {
                            if c1_indegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("c1Indegree"));
                            }
                            c1_indegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::C5Indegree => {
                            if c5_indegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("c5Indegree"));
                            }
                            c5_indegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::C95Outdegree => {
                            if c95_outdegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("c95Outdegree"));
                            }
                            c95_outdegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::C99Outdegree => {
                            if c99_outdegree__.is_some() {
                                return Err(serde::de::Error::duplicate_field("c99Outdegree"));
                            }
                            c99_outdegree__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::IndegreeCount => {
                            if indegree_count__.is_some() {
                                return Err(serde::de::Error::duplicate_field("indegreeCount"));
                            }
                            indegree_count__ = 
                                Some(map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter().map(|x| x.0).collect())
                            ;
                        }
                        GeneratedField::OutdegreeHistogram => {
                            if outdegree_histogram__.is_some() {
                                return Err(serde::de::Error::duplicate_field("outdegreeHistogram"));
                            }
                            outdegree_histogram__ = 
                                Some(map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter().map(|x| x.0).collect())
                            ;
                        }
                        GeneratedField::IndegreeHistogram => {
                            if indegree_histogram__.is_some() {
                                return Err(serde::de::Error::duplicate_field("indegreeHistogram"));
                            }
                            indegree_histogram__ = 
                                Some(map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter().map(|x| x.0).collect())
                            ;
                        }
                    }
                }
                Ok(info::index::Statistics {
                    valid: valid__.unwrap_or_default(),
                    median_indegree: median_indegree__.unwrap_or_default(),
                    median_outdegree: median_outdegree__.unwrap_or_default(),
                    max_number_of_indegree: max_number_of_indegree__.unwrap_or_default(),
                    max_number_of_outdegree: max_number_of_outdegree__.unwrap_or_default(),
                    min_number_of_indegree: min_number_of_indegree__.unwrap_or_default(),
                    min_number_of_outdegree: min_number_of_outdegree__.unwrap_or_default(),
                    mode_indegree: mode_indegree__.unwrap_or_default(),
                    mode_outdegree: mode_outdegree__.unwrap_or_default(),
                    nodes_skipped_for_10_edges: nodes_skipped_for_10_edges__.unwrap_or_default(),
                    nodes_skipped_for_indegree_distance: nodes_skipped_for_indegree_distance__.unwrap_or_default(),
                    number_of_edges: number_of_edges__.unwrap_or_default(),
                    number_of_indexed_objects: number_of_indexed_objects__.unwrap_or_default(),
                    number_of_nodes: number_of_nodes__.unwrap_or_default(),
                    number_of_nodes_without_edges: number_of_nodes_without_edges__.unwrap_or_default(),
                    number_of_nodes_without_indegree: number_of_nodes_without_indegree__.unwrap_or_default(),
                    number_of_objects: number_of_objects__.unwrap_or_default(),
                    number_of_removed_objects: number_of_removed_objects__.unwrap_or_default(),
                    size_of_object_repository: size_of_object_repository__.unwrap_or_default(),
                    size_of_refinement_object_repository: size_of_refinement_object_repository__.unwrap_or_default(),
                    variance_of_indegree: variance_of_indegree__.unwrap_or_default(),
                    variance_of_outdegree: variance_of_outdegree__.unwrap_or_default(),
                    mean_edge_length: mean_edge_length__.unwrap_or_default(),
                    mean_edge_length_for_10_edges: mean_edge_length_for_10_edges__.unwrap_or_default(),
                    mean_indegree_distance_for_10_edges: mean_indegree_distance_for_10_edges__.unwrap_or_default(),
                    mean_number_of_edges_per_node: mean_number_of_edges_per_node__.unwrap_or_default(),
                    c1_indegree: c1_indegree__.unwrap_or_default(),
                    c5_indegree: c5_indegree__.unwrap_or_default(),
                    c95_outdegree: c95_outdegree__.unwrap_or_default(),
                    c99_outdegree: c99_outdegree__.unwrap_or_default(),
                    indegree_count: indegree_count__.unwrap_or_default(),
                    outdegree_histogram: outdegree_histogram__.unwrap_or_default(),
                    indegree_histogram: indegree_histogram__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Index.Statistics", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::index::StatisticsDetail {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.details.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Index.StatisticsDetail", len)?;
        if !self.details.is_empty() {
            struct_ser.serialize_field("details", &self.details)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::index::StatisticsDetail {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "details",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Details,
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
                            "details" => Ok(GeneratedField::Details),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::index::StatisticsDetail;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Index.StatisticsDetail")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::index::StatisticsDetail, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut details__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Details => {
                            if details__.is_some() {
                                return Err(serde::de::Error::duplicate_field("details"));
                            }
                            details__ = Some(
                                map_.next_value::<std::collections::HashMap<_, _>>()?
                            );
                        }
                    }
                }
                Ok(info::index::StatisticsDetail {
                    details: details__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Index.StatisticsDetail", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::index::Uuid {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Info.Index.UUID", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::index::Uuid {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::index::Uuid;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Index.UUID")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::index::Uuid, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(info::index::Uuid {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Index.UUID", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::index::uuid::Committed {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.uuid.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Index.UUID.Committed", len)?;
        if !self.uuid.is_empty() {
            struct_ser.serialize_field("uuid", &self.uuid)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::index::uuid::Committed {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "uuid",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Uuid,
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
                            "uuid" => Ok(GeneratedField::Uuid),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::index::uuid::Committed;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Index.UUID.Committed")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::index::uuid::Committed, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut uuid__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Uuid => {
                            if uuid__.is_some() {
                                return Err(serde::de::Error::duplicate_field("uuid"));
                            }
                            uuid__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(info::index::uuid::Committed {
                    uuid: uuid__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Index.UUID.Committed", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::index::uuid::Uncommitted {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.uuid.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Index.UUID.Uncommitted", len)?;
        if !self.uuid.is_empty() {
            struct_ser.serialize_field("uuid", &self.uuid)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::index::uuid::Uncommitted {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "uuid",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Uuid,
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
                            "uuid" => Ok(GeneratedField::Uuid),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::index::uuid::Uncommitted;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Index.UUID.Uncommitted")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::index::uuid::Uncommitted, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut uuid__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Uuid => {
                            if uuid__.is_some() {
                                return Err(serde::de::Error::duplicate_field("uuid"));
                            }
                            uuid__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(info::index::uuid::Uncommitted {
                    uuid: uuid__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Index.UUID.Uncommitted", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::Labels {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.labels.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Labels", len)?;
        if !self.labels.is_empty() {
            struct_ser.serialize_field("labels", &self.labels)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::Labels {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "labels",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Labels,
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
                            "labels" => Ok(GeneratedField::Labels),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::Labels;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Labels")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::Labels, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut labels__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Labels => {
                            if labels__.is_some() {
                                return Err(serde::de::Error::duplicate_field("labels"));
                            }
                            labels__ = Some(
                                map_.next_value::<std::collections::HashMap<_, _>>()?
                            );
                        }
                    }
                }
                Ok(info::Labels {
                    labels: labels__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Labels", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::Memory {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.limit != 0. {
            len += 1;
        }
        if self.request != 0. {
            len += 1;
        }
        if self.usage != 0. {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Memory", len)?;
        if self.limit != 0. {
            struct_ser.serialize_field("limit", &self.limit)?;
        }
        if self.request != 0. {
            struct_ser.serialize_field("request", &self.request)?;
        }
        if self.usage != 0. {
            struct_ser.serialize_field("usage", &self.usage)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::Memory {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "limit",
            "request",
            "usage",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Limit,
            Request,
            Usage,
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
                            "limit" => Ok(GeneratedField::Limit),
                            "request" => Ok(GeneratedField::Request),
                            "usage" => Ok(GeneratedField::Usage),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::Memory;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Memory")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::Memory, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut limit__ = None;
                let mut request__ = None;
                let mut usage__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Limit => {
                            if limit__.is_some() {
                                return Err(serde::de::Error::duplicate_field("limit"));
                            }
                            limit__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Request => {
                            if request__.is_some() {
                                return Err(serde::de::Error::duplicate_field("request"));
                            }
                            request__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Usage => {
                            if usage__.is_some() {
                                return Err(serde::de::Error::duplicate_field("usage"));
                            }
                            usage__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(info::Memory {
                    limit: limit__.unwrap_or_default(),
                    request: request__.unwrap_or_default(),
                    usage: usage__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Memory", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::Node {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.name.is_empty() {
            len += 1;
        }
        if !self.internal_addr.is_empty() {
            len += 1;
        }
        if !self.external_addr.is_empty() {
            len += 1;
        }
        if self.cpu.is_some() {
            len += 1;
        }
        if self.memory.is_some() {
            len += 1;
        }
        if self.pods.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Node", len)?;
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.internal_addr.is_empty() {
            struct_ser.serialize_field("internalAddr", &self.internal_addr)?;
        }
        if !self.external_addr.is_empty() {
            struct_ser.serialize_field("externalAddr", &self.external_addr)?;
        }
        if let Some(v) = self.cpu.as_ref() {
            struct_ser.serialize_field("cpu", v)?;
        }
        if let Some(v) = self.memory.as_ref() {
            struct_ser.serialize_field("memory", v)?;
        }
        if let Some(v) = self.pods.as_ref() {
            struct_ser.serialize_field("Pods", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::Node {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "name",
            "internal_addr",
            "internalAddr",
            "external_addr",
            "externalAddr",
            "cpu",
            "memory",
            "Pods",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Name,
            InternalAddr,
            ExternalAddr,
            Cpu,
            Memory,
            Pods,
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
                            "name" => Ok(GeneratedField::Name),
                            "internalAddr" | "internal_addr" => Ok(GeneratedField::InternalAddr),
                            "externalAddr" | "external_addr" => Ok(GeneratedField::ExternalAddr),
                            "cpu" => Ok(GeneratedField::Cpu),
                            "memory" => Ok(GeneratedField::Memory),
                            "Pods" => Ok(GeneratedField::Pods),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::Node;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Node")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::Node, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut name__ = None;
                let mut internal_addr__ = None;
                let mut external_addr__ = None;
                let mut cpu__ = None;
                let mut memory__ = None;
                let mut pods__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::InternalAddr => {
                            if internal_addr__.is_some() {
                                return Err(serde::de::Error::duplicate_field("internalAddr"));
                            }
                            internal_addr__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ExternalAddr => {
                            if external_addr__.is_some() {
                                return Err(serde::de::Error::duplicate_field("externalAddr"));
                            }
                            external_addr__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Cpu => {
                            if cpu__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cpu"));
                            }
                            cpu__ = map_.next_value()?;
                        }
                        GeneratedField::Memory => {
                            if memory__.is_some() {
                                return Err(serde::de::Error::duplicate_field("memory"));
                            }
                            memory__ = map_.next_value()?;
                        }
                        GeneratedField::Pods => {
                            if pods__.is_some() {
                                return Err(serde::de::Error::duplicate_field("Pods"));
                            }
                            pods__ = map_.next_value()?;
                        }
                    }
                }
                Ok(info::Node {
                    name: name__.unwrap_or_default(),
                    internal_addr: internal_addr__.unwrap_or_default(),
                    external_addr: external_addr__.unwrap_or_default(),
                    cpu: cpu__,
                    memory: memory__,
                    pods: pods__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Node", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::Nodes {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.nodes.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Nodes", len)?;
        if !self.nodes.is_empty() {
            struct_ser.serialize_field("nodes", &self.nodes)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::Nodes {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "nodes",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Nodes,
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
                            "nodes" => Ok(GeneratedField::Nodes),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::Nodes;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Nodes")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::Nodes, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut nodes__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Nodes => {
                            if nodes__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nodes"));
                            }
                            nodes__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(info::Nodes {
                    nodes: nodes__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Nodes", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::Pod {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.app_name.is_empty() {
            len += 1;
        }
        if !self.name.is_empty() {
            len += 1;
        }
        if !self.namespace.is_empty() {
            len += 1;
        }
        if !self.ip.is_empty() {
            len += 1;
        }
        if self.cpu.is_some() {
            len += 1;
        }
        if self.memory.is_some() {
            len += 1;
        }
        if self.node.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Pod", len)?;
        if !self.app_name.is_empty() {
            struct_ser.serialize_field("appName", &self.app_name)?;
        }
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.namespace.is_empty() {
            struct_ser.serialize_field("namespace", &self.namespace)?;
        }
        if !self.ip.is_empty() {
            struct_ser.serialize_field("ip", &self.ip)?;
        }
        if let Some(v) = self.cpu.as_ref() {
            struct_ser.serialize_field("cpu", v)?;
        }
        if let Some(v) = self.memory.as_ref() {
            struct_ser.serialize_field("memory", v)?;
        }
        if let Some(v) = self.node.as_ref() {
            struct_ser.serialize_field("node", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::Pod {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "app_name",
            "appName",
            "name",
            "namespace",
            "ip",
            "cpu",
            "memory",
            "node",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            AppName,
            Name,
            Namespace,
            Ip,
            Cpu,
            Memory,
            Node,
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
                            "appName" | "app_name" => Ok(GeneratedField::AppName),
                            "name" => Ok(GeneratedField::Name),
                            "namespace" => Ok(GeneratedField::Namespace),
                            "ip" => Ok(GeneratedField::Ip),
                            "cpu" => Ok(GeneratedField::Cpu),
                            "memory" => Ok(GeneratedField::Memory),
                            "node" => Ok(GeneratedField::Node),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::Pod;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Pod")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::Pod, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut app_name__ = None;
                let mut name__ = None;
                let mut namespace__ = None;
                let mut ip__ = None;
                let mut cpu__ = None;
                let mut memory__ = None;
                let mut node__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::AppName => {
                            if app_name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("appName"));
                            }
                            app_name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Namespace => {
                            if namespace__.is_some() {
                                return Err(serde::de::Error::duplicate_field("namespace"));
                            }
                            namespace__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ip => {
                            if ip__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ip"));
                            }
                            ip__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Cpu => {
                            if cpu__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cpu"));
                            }
                            cpu__ = map_.next_value()?;
                        }
                        GeneratedField::Memory => {
                            if memory__.is_some() {
                                return Err(serde::de::Error::duplicate_field("memory"));
                            }
                            memory__ = map_.next_value()?;
                        }
                        GeneratedField::Node => {
                            if node__.is_some() {
                                return Err(serde::de::Error::duplicate_field("node"));
                            }
                            node__ = map_.next_value()?;
                        }
                    }
                }
                Ok(info::Pod {
                    app_name: app_name__.unwrap_or_default(),
                    name: name__.unwrap_or_default(),
                    namespace: namespace__.unwrap_or_default(),
                    ip: ip__.unwrap_or_default(),
                    cpu: cpu__,
                    memory: memory__,
                    node: node__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Pod", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::Pods {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.pods.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Pods", len)?;
        if !self.pods.is_empty() {
            struct_ser.serialize_field("pods", &self.pods)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::Pods {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "pods",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Pods,
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
                            "pods" => Ok(GeneratedField::Pods),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::Pods;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Pods")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::Pods, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut pods__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Pods => {
                            if pods__.is_some() {
                                return Err(serde::de::Error::duplicate_field("pods"));
                            }
                            pods__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(info::Pods {
                    pods: pods__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Pods", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::ResourceStats {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.name.is_empty() {
            len += 1;
        }
        if !self.ip.is_empty() {
            len += 1;
        }
        if self.cgroup_stats.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.ResourceStats", len)?;
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.ip.is_empty() {
            struct_ser.serialize_field("ip", &self.ip)?;
        }
        if let Some(v) = self.cgroup_stats.as_ref() {
            struct_ser.serialize_field("cgroupStats", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::ResourceStats {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "name",
            "ip",
            "cgroup_stats",
            "cgroupStats",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Name,
            Ip,
            CgroupStats,
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
                            "name" => Ok(GeneratedField::Name),
                            "ip" => Ok(GeneratedField::Ip),
                            "cgroupStats" | "cgroup_stats" => Ok(GeneratedField::CgroupStats),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::ResourceStats;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.ResourceStats")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::ResourceStats, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut name__ = None;
                let mut ip__ = None;
                let mut cgroup_stats__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ip => {
                            if ip__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ip"));
                            }
                            ip__ = Some(map_.next_value()?);
                        }
                        GeneratedField::CgroupStats => {
                            if cgroup_stats__.is_some() {
                                return Err(serde::de::Error::duplicate_field("cgroupStats"));
                            }
                            cgroup_stats__ = map_.next_value()?;
                        }
                    }
                }
                Ok(info::ResourceStats {
                    name: name__.unwrap_or_default(),
                    ip: ip__.unwrap_or_default(),
                    cgroup_stats: cgroup_stats__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.ResourceStats", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::Service {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.name.is_empty() {
            len += 1;
        }
        if !self.cluster_ip.is_empty() {
            len += 1;
        }
        if !self.cluster_ips.is_empty() {
            len += 1;
        }
        if !self.ports.is_empty() {
            len += 1;
        }
        if self.labels.is_some() {
            len += 1;
        }
        if self.annotations.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Service", len)?;
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.cluster_ip.is_empty() {
            struct_ser.serialize_field("clusterIp", &self.cluster_ip)?;
        }
        if !self.cluster_ips.is_empty() {
            struct_ser.serialize_field("clusterIps", &self.cluster_ips)?;
        }
        if !self.ports.is_empty() {
            struct_ser.serialize_field("ports", &self.ports)?;
        }
        if let Some(v) = self.labels.as_ref() {
            struct_ser.serialize_field("labels", v)?;
        }
        if let Some(v) = self.annotations.as_ref() {
            struct_ser.serialize_field("annotations", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::Service {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "name",
            "cluster_ip",
            "clusterIp",
            "cluster_ips",
            "clusterIps",
            "ports",
            "labels",
            "annotations",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Name,
            ClusterIp,
            ClusterIps,
            Ports,
            Labels,
            Annotations,
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
                            "name" => Ok(GeneratedField::Name),
                            "clusterIp" | "cluster_ip" => Ok(GeneratedField::ClusterIp),
                            "clusterIps" | "cluster_ips" => Ok(GeneratedField::ClusterIps),
                            "ports" => Ok(GeneratedField::Ports),
                            "labels" => Ok(GeneratedField::Labels),
                            "annotations" => Ok(GeneratedField::Annotations),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::Service;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Service")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::Service, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut name__ = None;
                let mut cluster_ip__ = None;
                let mut cluster_ips__ = None;
                let mut ports__ = None;
                let mut labels__ = None;
                let mut annotations__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ClusterIp => {
                            if cluster_ip__.is_some() {
                                return Err(serde::de::Error::duplicate_field("clusterIp"));
                            }
                            cluster_ip__ = Some(map_.next_value()?);
                        }
                        GeneratedField::ClusterIps => {
                            if cluster_ips__.is_some() {
                                return Err(serde::de::Error::duplicate_field("clusterIps"));
                            }
                            cluster_ips__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ports => {
                            if ports__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ports"));
                            }
                            ports__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Labels => {
                            if labels__.is_some() {
                                return Err(serde::de::Error::duplicate_field("labels"));
                            }
                            labels__ = map_.next_value()?;
                        }
                        GeneratedField::Annotations => {
                            if annotations__.is_some() {
                                return Err(serde::de::Error::duplicate_field("annotations"));
                            }
                            annotations__ = map_.next_value()?;
                        }
                    }
                }
                Ok(info::Service {
                    name: name__.unwrap_or_default(),
                    cluster_ip: cluster_ip__.unwrap_or_default(),
                    cluster_ips: cluster_ips__.unwrap_or_default(),
                    ports: ports__.unwrap_or_default(),
                    labels: labels__,
                    annotations: annotations__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Service", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::ServicePort {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.name.is_empty() {
            len += 1;
        }
        if self.port != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.ServicePort", len)?;
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if self.port != 0 {
            struct_ser.serialize_field("port", &self.port)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::ServicePort {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "name",
            "port",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Name,
            Port,
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
                            "name" => Ok(GeneratedField::Name),
                            "port" => Ok(GeneratedField::Port),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::ServicePort;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.ServicePort")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::ServicePort, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut name__ = None;
                let mut port__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Port => {
                            if port__.is_some() {
                                return Err(serde::de::Error::duplicate_field("port"));
                            }
                            port__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(info::ServicePort {
                    name: name__.unwrap_or_default(),
                    port: port__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.ServicePort", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for info::Services {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.services.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Info.Services", len)?;
        if !self.services.is_empty() {
            struct_ser.serialize_field("services", &self.services)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for info::Services {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "services",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Services,
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
                            "services" => Ok(GeneratedField::Services),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = info::Services;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Info.Services")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<info::Services, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut services__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Services => {
                            if services__.is_some() {
                                return Err(serde::de::Error::duplicate_field("services"));
                            }
                            services__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(info::Services {
                    services: services__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Info.Services", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Insert {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Insert", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Insert {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Insert;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Insert")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Insert, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Insert {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Insert", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for insert::Config {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.skip_strict_exist_check {
            len += 1;
        }
        if self.filters.is_some() {
            len += 1;
        }
        if self.timestamp != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Insert.Config", len)?;
        if self.skip_strict_exist_check {
            struct_ser.serialize_field("skipStrictExistCheck", &self.skip_strict_exist_check)?;
        }
        if let Some(v) = self.filters.as_ref() {
            struct_ser.serialize_field("filters", v)?;
        }
        if self.timestamp != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("timestamp", ToString::to_string(&self.timestamp).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for insert::Config {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "skip_strict_exist_check",
            "skipStrictExistCheck",
            "filters",
            "timestamp",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SkipStrictExistCheck,
            Filters,
            Timestamp,
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
                            "skipStrictExistCheck" | "skip_strict_exist_check" => Ok(GeneratedField::SkipStrictExistCheck),
                            "filters" => Ok(GeneratedField::Filters),
                            "timestamp" => Ok(GeneratedField::Timestamp),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = insert::Config;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Insert.Config")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<insert::Config, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut skip_strict_exist_check__ = None;
                let mut filters__ = None;
                let mut timestamp__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SkipStrictExistCheck => {
                            if skip_strict_exist_check__.is_some() {
                                return Err(serde::de::Error::duplicate_field("skipStrictExistCheck"));
                            }
                            skip_strict_exist_check__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Filters => {
                            if filters__.is_some() {
                                return Err(serde::de::Error::duplicate_field("filters"));
                            }
                            filters__ = map_.next_value()?;
                        }
                        GeneratedField::Timestamp => {
                            if timestamp__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timestamp"));
                            }
                            timestamp__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(insert::Config {
                    skip_strict_exist_check: skip_strict_exist_check__.unwrap_or_default(),
                    filters: filters__,
                    timestamp: timestamp__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Insert.Config", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for insert::MultiObjectRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.requests.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Insert.MultiObjectRequest", len)?;
        if !self.requests.is_empty() {
            struct_ser.serialize_field("requests", &self.requests)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for insert::MultiObjectRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "requests",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Requests,
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
                            "requests" => Ok(GeneratedField::Requests),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = insert::MultiObjectRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Insert.MultiObjectRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<insert::MultiObjectRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut requests__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Requests => {
                            if requests__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requests"));
                            }
                            requests__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(insert::MultiObjectRequest {
                    requests: requests__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Insert.MultiObjectRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for insert::MultiRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.requests.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Insert.MultiRequest", len)?;
        if !self.requests.is_empty() {
            struct_ser.serialize_field("requests", &self.requests)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for insert::MultiRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "requests",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Requests,
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
                            "requests" => Ok(GeneratedField::Requests),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = insert::MultiRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Insert.MultiRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<insert::MultiRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut requests__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Requests => {
                            if requests__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requests"));
                            }
                            requests__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(insert::MultiRequest {
                    requests: requests__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Insert.MultiRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for insert::ObjectRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.object.is_some() {
            len += 1;
        }
        if self.config.is_some() {
            len += 1;
        }
        if self.vectorizer.is_some() {
            len += 1;
        }
        if self.metadata.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Insert.ObjectRequest", len)?;
        if let Some(v) = self.object.as_ref() {
            struct_ser.serialize_field("object", v)?;
        }
        if let Some(v) = self.config.as_ref() {
            struct_ser.serialize_field("config", v)?;
        }
        if let Some(v) = self.vectorizer.as_ref() {
            struct_ser.serialize_field("vectorizer", v)?;
        }
        if let Some(v) = self.metadata.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("metadata", pbjson::private::base64::encode(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for insert::ObjectRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "object",
            "config",
            "vectorizer",
            "metadata",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Object,
            Config,
            Vectorizer,
            Metadata,
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
                            "object" => Ok(GeneratedField::Object),
                            "config" => Ok(GeneratedField::Config),
                            "vectorizer" => Ok(GeneratedField::Vectorizer),
                            "metadata" => Ok(GeneratedField::Metadata),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = insert::ObjectRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Insert.ObjectRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<insert::ObjectRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut object__ = None;
                let mut config__ = None;
                let mut vectorizer__ = None;
                let mut metadata__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Object => {
                            if object__.is_some() {
                                return Err(serde::de::Error::duplicate_field("object"));
                            }
                            object__ = map_.next_value()?;
                        }
                        GeneratedField::Config => {
                            if config__.is_some() {
                                return Err(serde::de::Error::duplicate_field("config"));
                            }
                            config__ = map_.next_value()?;
                        }
                        GeneratedField::Vectorizer => {
                            if vectorizer__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vectorizer"));
                            }
                            vectorizer__ = map_.next_value()?;
                        }
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::BytesDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                    }
                }
                Ok(insert::ObjectRequest {
                    object: object__,
                    config: config__,
                    vectorizer: vectorizer__,
                    metadata: metadata__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Insert.ObjectRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for insert::Request {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.vector.is_some() {
            len += 1;
        }
        if self.config.is_some() {
            len += 1;
        }
        if self.metadata.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Insert.Request", len)?;
        if let Some(v) = self.vector.as_ref() {
            struct_ser.serialize_field("vector", v)?;
        }
        if let Some(v) = self.config.as_ref() {
            struct_ser.serialize_field("config", v)?;
        }
        if let Some(v) = self.metadata.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("metadata", pbjson::private::base64::encode(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for insert::Request {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "vector",
            "config",
            "metadata",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Vector,
            Config,
            Metadata,
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
                            "vector" => Ok(GeneratedField::Vector),
                            "config" => Ok(GeneratedField::Config),
                            "metadata" => Ok(GeneratedField::Metadata),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = insert::Request;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Insert.Request")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<insert::Request, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut vector__ = None;
                let mut config__ = None;
                let mut metadata__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Vector => {
                            if vector__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vector"));
                            }
                            vector__ = map_.next_value()?;
                        }
                        GeneratedField::Config => {
                            if config__.is_some() {
                                return Err(serde::de::Error::duplicate_field("config"));
                            }
                            config__ = map_.next_value()?;
                        }
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::BytesDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                    }
                }
                Ok(insert::Request {
                    vector: vector__,
                    config: config__,
                    metadata: metadata__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Insert.Request", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Meta {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Meta", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Meta {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Meta;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Meta")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Meta, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Meta {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Meta", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for meta::Key {
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
        let mut struct_ser = serializer.serialize_struct("payload.v1.Meta.Key", len)?;
        if !self.key.is_empty() {
            struct_ser.serialize_field("key", &self.key)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for meta::Key {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "key",
        ];

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

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
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
            type Value = meta::Key;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Meta.Key")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<meta::Key, V::Error>
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
                            key__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(meta::Key {
                    key: key__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Meta.Key", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for meta::KeyValue {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.key.is_some() {
            len += 1;
        }
        if self.value.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Meta.KeyValue", len)?;
        if let Some(v) = self.key.as_ref() {
            struct_ser.serialize_field("key", v)?;
        }
        if let Some(v) = self.value.as_ref() {
            struct_ser.serialize_field("value", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for meta::KeyValue {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "key",
            "value",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Key,
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

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
                            "key" => Ok(GeneratedField::Key),
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
            type Value = meta::KeyValue;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Meta.KeyValue")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<meta::KeyValue, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut key__ = None;
                let mut value__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Key => {
                            if key__.is_some() {
                                return Err(serde::de::Error::duplicate_field("key"));
                            }
                            key__ = map_.next_value()?;
                        }
                        GeneratedField::Value => {
                            if value__.is_some() {
                                return Err(serde::de::Error::duplicate_field("value"));
                            }
                            value__ = map_.next_value()?;
                        }
                    }
                }
                Ok(meta::KeyValue {
                    key: key__,
                    value: value__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Meta.KeyValue", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for meta::Value {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.value.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Meta.Value", len)?;
        if let Some(v) = self.value.as_ref() {
            struct_ser.serialize_field("value", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for meta::Value {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "value",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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

                    fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                        write!(formatter, "expected one of: {:?}", &FIELDS)
                    }

                    #[allow(unused_variables)]
                    fn visit_str<E>(self, value: &str) -> std::result::Result<GeneratedField, E>
                    where
                        E: serde::de::Error,
                    {
                        match value {
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
            type Value = meta::Value;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Meta.Value")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<meta::Value, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut value__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Value => {
                            if value__.is_some() {
                                return Err(serde::de::Error::duplicate_field("value"));
                            }
                            value__ = map_.next_value()?;
                        }
                    }
                }
                Ok(meta::Value {
                    value: value__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Meta.Value", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Mirror {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Mirror", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Mirror {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Mirror;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Mirror")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Mirror, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Mirror {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Mirror", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for mirror::Target {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.host.is_empty() {
            len += 1;
        }
        if self.port != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Mirror.Target", len)?;
        if !self.host.is_empty() {
            struct_ser.serialize_field("host", &self.host)?;
        }
        if self.port != 0 {
            struct_ser.serialize_field("port", &self.port)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for mirror::Target {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "host",
            "port",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Host,
            Port,
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
                            "host" => Ok(GeneratedField::Host),
                            "port" => Ok(GeneratedField::Port),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = mirror::Target;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Mirror.Target")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<mirror::Target, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut host__ = None;
                let mut port__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Host => {
                            if host__.is_some() {
                                return Err(serde::de::Error::duplicate_field("host"));
                            }
                            host__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Port => {
                            if port__.is_some() {
                                return Err(serde::de::Error::duplicate_field("port"));
                            }
                            port__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(mirror::Target {
                    host: host__.unwrap_or_default(),
                    port: port__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Mirror.Target", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for mirror::Targets {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.targets.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Mirror.Targets", len)?;
        if !self.targets.is_empty() {
            struct_ser.serialize_field("targets", &self.targets)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for mirror::Targets {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "targets",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Targets,
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
                            "targets" => Ok(GeneratedField::Targets),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = mirror::Targets;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Mirror.Targets")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<mirror::Targets, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut targets__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Targets => {
                            if targets__.is_some() {
                                return Err(serde::de::Error::duplicate_field("targets"));
                            }
                            targets__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(mirror::Targets {
                    targets: targets__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Mirror.Targets", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Object {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Object", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Object {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Object;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Object, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Object {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::Blob {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.id.is_empty() {
            len += 1;
        }
        if !self.object.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.Blob", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.object.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("object", pbjson::private::base64::encode(&self.object).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::Blob {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "object",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Object,
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
                            "id" => Ok(GeneratedField::Id),
                            "object" => Ok(GeneratedField::Object),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::Blob;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.Blob")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::Blob, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut object__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Object => {
                            if object__.is_some() {
                                return Err(serde::de::Error::duplicate_field("object"));
                            }
                            object__ = 
                                Some(map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(object::Blob {
                    id: id__.unwrap_or_default(),
                    object: object__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.Blob", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::Distance {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.id.is_empty() {
            len += 1;
        }
        if self.distance != 0. {
            len += 1;
        }
        if self.metadata.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.Distance", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if self.distance != 0. {
            struct_ser.serialize_field("distance", &self.distance)?;
        }
        if let Some(v) = self.metadata.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("metadata", pbjson::private::base64::encode(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::Distance {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "distance",
            "metadata",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Distance,
            Metadata,
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
                            "id" => Ok(GeneratedField::Id),
                            "distance" => Ok(GeneratedField::Distance),
                            "metadata" => Ok(GeneratedField::Metadata),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::Distance;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.Distance")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::Distance, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut distance__ = None;
                let mut metadata__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Distance => {
                            if distance__.is_some() {
                                return Err(serde::de::Error::duplicate_field("distance"));
                            }
                            distance__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::BytesDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                    }
                }
                Ok(object::Distance {
                    id: id__.unwrap_or_default(),
                    distance: distance__.unwrap_or_default(),
                    metadata: metadata__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.Distance", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::Id {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.id.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.ID", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::Id {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
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
                            "id" => Ok(GeneratedField::Id),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::Id;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.ID")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::Id, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(object::Id {
                    id: id__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.ID", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::IDs {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.ids.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.IDs", len)?;
        if !self.ids.is_empty() {
            struct_ser.serialize_field("ids", &self.ids)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::IDs {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "ids",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Ids,
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
                            "ids" => Ok(GeneratedField::Ids),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::IDs;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.IDs")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::IDs, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut ids__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Ids => {
                            if ids__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ids"));
                            }
                            ids__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(object::IDs {
                    ids: ids__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.IDs", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::List {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Object.List", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::List {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::List;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.List")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::List, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(object::List {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.List", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::list::Request {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Object.List.Request", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::list::Request {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::list::Request;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.List.Request")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::list::Request, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(object::list::Request {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.List.Request", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::list::Response {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.payload.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.List.Response", len)?;
        if let Some(v) = self.payload.as_ref() {
            match v {
                object::list::response::Payload::Vector(v) => {
                    struct_ser.serialize_field("vector", v)?;
                }
                object::list::response::Payload::Status(v) => {
                    struct_ser.serialize_field("status", v)?;
                }
            }
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::list::Response {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "vector",
            "status",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Vector,
            Status,
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
                            "vector" => Ok(GeneratedField::Vector),
                            "status" => Ok(GeneratedField::Status),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::list::Response;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.List.Response")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::list::Response, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut payload__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Vector => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vector"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(object::list::response::Payload::Vector)
;
                        }
                        GeneratedField::Status => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(object::list::response::Payload::Status)
;
                        }
                    }
                }
                Ok(object::list::Response {
                    payload: payload__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.List.Response", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::Location {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.name.is_empty() {
            len += 1;
        }
        if !self.uuid.is_empty() {
            len += 1;
        }
        if !self.ips.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.Location", len)?;
        if !self.name.is_empty() {
            struct_ser.serialize_field("name", &self.name)?;
        }
        if !self.uuid.is_empty() {
            struct_ser.serialize_field("uuid", &self.uuid)?;
        }
        if !self.ips.is_empty() {
            struct_ser.serialize_field("ips", &self.ips)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::Location {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "name",
            "uuid",
            "ips",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Name,
            Uuid,
            Ips,
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
                            "name" => Ok(GeneratedField::Name),
                            "uuid" => Ok(GeneratedField::Uuid),
                            "ips" => Ok(GeneratedField::Ips),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::Location;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.Location")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::Location, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut name__ = None;
                let mut uuid__ = None;
                let mut ips__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Name => {
                            if name__.is_some() {
                                return Err(serde::de::Error::duplicate_field("name"));
                            }
                            name__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Uuid => {
                            if uuid__.is_some() {
                                return Err(serde::de::Error::duplicate_field("uuid"));
                            }
                            uuid__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Ips => {
                            if ips__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ips"));
                            }
                            ips__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(object::Location {
                    name: name__.unwrap_or_default(),
                    uuid: uuid__.unwrap_or_default(),
                    ips: ips__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.Location", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::Locations {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.locations.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.Locations", len)?;
        if !self.locations.is_empty() {
            struct_ser.serialize_field("locations", &self.locations)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::Locations {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "locations",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Locations,
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
                            "locations" => Ok(GeneratedField::Locations),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::Locations;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.Locations")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::Locations, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut locations__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Locations => {
                            if locations__.is_some() {
                                return Err(serde::de::Error::duplicate_field("locations"));
                            }
                            locations__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(object::Locations {
                    locations: locations__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.Locations", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::ReshapeVector {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.object.is_empty() {
            len += 1;
        }
        if !self.shape.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.ReshapeVector", len)?;
        if !self.object.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("object", pbjson::private::base64::encode(&self.object).as_str())?;
        }
        if !self.shape.is_empty() {
            struct_ser.serialize_field("shape", &self.shape)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::ReshapeVector {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "object",
            "shape",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Object,
            Shape,
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
                            "object" => Ok(GeneratedField::Object),
                            "shape" => Ok(GeneratedField::Shape),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::ReshapeVector;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.ReshapeVector")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::ReshapeVector, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut object__ = None;
                let mut shape__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Object => {
                            if object__.is_some() {
                                return Err(serde::de::Error::duplicate_field("object"));
                            }
                            object__ = 
                                Some(map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Shape => {
                            if shape__.is_some() {
                                return Err(serde::de::Error::duplicate_field("shape"));
                            }
                            shape__ = 
                                Some(map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter().map(|x| x.0).collect())
                            ;
                        }
                    }
                }
                Ok(object::ReshapeVector {
                    object: object__.unwrap_or_default(),
                    shape: shape__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.ReshapeVector", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::StreamBlob {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.payload.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.StreamBlob", len)?;
        if let Some(v) = self.payload.as_ref() {
            match v {
                object::stream_blob::Payload::Blob(v) => {
                    struct_ser.serialize_field("blob", v)?;
                }
                object::stream_blob::Payload::Status(v) => {
                    struct_ser.serialize_field("status", v)?;
                }
            }
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::StreamBlob {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "blob",
            "status",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Blob,
            Status,
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
                            "blob" => Ok(GeneratedField::Blob),
                            "status" => Ok(GeneratedField::Status),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::StreamBlob;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.StreamBlob")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::StreamBlob, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut payload__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Blob => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("blob"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(object::stream_blob::Payload::Blob)
;
                        }
                        GeneratedField::Status => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(object::stream_blob::Payload::Status)
;
                        }
                    }
                }
                Ok(object::StreamBlob {
                    payload: payload__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.StreamBlob", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::StreamDistance {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.payload.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.StreamDistance", len)?;
        if let Some(v) = self.payload.as_ref() {
            match v {
                object::stream_distance::Payload::Distance(v) => {
                    struct_ser.serialize_field("distance", v)?;
                }
                object::stream_distance::Payload::Status(v) => {
                    struct_ser.serialize_field("status", v)?;
                }
            }
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::StreamDistance {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "distance",
            "status",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Distance,
            Status,
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
                            "distance" => Ok(GeneratedField::Distance),
                            "status" => Ok(GeneratedField::Status),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::StreamDistance;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.StreamDistance")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::StreamDistance, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut payload__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Distance => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("distance"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(object::stream_distance::Payload::Distance)
;
                        }
                        GeneratedField::Status => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(object::stream_distance::Payload::Status)
;
                        }
                    }
                }
                Ok(object::StreamDistance {
                    payload: payload__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.StreamDistance", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::StreamLocation {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.payload.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.StreamLocation", len)?;
        if let Some(v) = self.payload.as_ref() {
            match v {
                object::stream_location::Payload::Location(v) => {
                    struct_ser.serialize_field("location", v)?;
                }
                object::stream_location::Payload::Status(v) => {
                    struct_ser.serialize_field("status", v)?;
                }
            }
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::StreamLocation {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "location",
            "status",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Location,
            Status,
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
                            "location" => Ok(GeneratedField::Location),
                            "status" => Ok(GeneratedField::Status),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::StreamLocation;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.StreamLocation")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::StreamLocation, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut payload__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Location => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("location"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(object::stream_location::Payload::Location)
;
                        }
                        GeneratedField::Status => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(object::stream_location::Payload::Status)
;
                        }
                    }
                }
                Ok(object::StreamLocation {
                    payload: payload__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.StreamLocation", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::StreamVector {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.payload.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.StreamVector", len)?;
        if let Some(v) = self.payload.as_ref() {
            match v {
                object::stream_vector::Payload::Vector(v) => {
                    struct_ser.serialize_field("vector", v)?;
                }
                object::stream_vector::Payload::Status(v) => {
                    struct_ser.serialize_field("status", v)?;
                }
            }
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::StreamVector {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "vector",
            "status",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Vector,
            Status,
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
                            "vector" => Ok(GeneratedField::Vector),
                            "status" => Ok(GeneratedField::Status),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::StreamVector;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.StreamVector")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::StreamVector, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut payload__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Vector => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vector"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(object::stream_vector::Payload::Vector)
;
                        }
                        GeneratedField::Status => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(object::stream_vector::Payload::Status)
;
                        }
                    }
                }
                Ok(object::StreamVector {
                    payload: payload__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.StreamVector", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::Timestamp {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.id.is_empty() {
            len += 1;
        }
        if self.timestamp != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.Timestamp", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if self.timestamp != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("timestamp", ToString::to_string(&self.timestamp).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::Timestamp {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "timestamp",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Timestamp,
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
                            "id" => Ok(GeneratedField::Id),
                            "timestamp" => Ok(GeneratedField::Timestamp),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::Timestamp;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.Timestamp")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::Timestamp, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut timestamp__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Timestamp => {
                            if timestamp__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timestamp"));
                            }
                            timestamp__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(object::Timestamp {
                    id: id__.unwrap_or_default(),
                    timestamp: timestamp__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.Timestamp", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::TimestampRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.id.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.TimestampRequest", len)?;
        if let Some(v) = self.id.as_ref() {
            struct_ser.serialize_field("id", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::TimestampRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
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
                            "id" => Ok(GeneratedField::Id),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::TimestampRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.TimestampRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::TimestampRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = map_.next_value()?;
                        }
                    }
                }
                Ok(object::TimestampRequest {
                    id: id__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.TimestampRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::Vector {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.id.is_empty() {
            len += 1;
        }
        if !self.vector.is_empty() {
            len += 1;
        }
        if self.timestamp != 0 {
            len += 1;
        }
        if self.metadata.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.Vector", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if !self.vector.is_empty() {
            struct_ser.serialize_field("vector", &self.vector)?;
        }
        if self.timestamp != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("timestamp", ToString::to_string(&self.timestamp).as_str())?;
        }
        if let Some(v) = self.metadata.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("metadata", pbjson::private::base64::encode(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::Vector {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "vector",
            "timestamp",
            "metadata",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Vector,
            Timestamp,
            Metadata,
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
                            "id" => Ok(GeneratedField::Id),
                            "vector" => Ok(GeneratedField::Vector),
                            "timestamp" => Ok(GeneratedField::Timestamp),
                            "metadata" => Ok(GeneratedField::Metadata),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::Vector;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.Vector")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::Vector, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut vector__ = None;
                let mut timestamp__ = None;
                let mut metadata__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Vector => {
                            if vector__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vector"));
                            }
                            vector__ = 
                                Some(map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter().map(|x| x.0).collect())
                            ;
                        }
                        GeneratedField::Timestamp => {
                            if timestamp__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timestamp"));
                            }
                            timestamp__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::BytesDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                    }
                }
                Ok(object::Vector {
                    id: id__.unwrap_or_default(),
                    vector: vector__.unwrap_or_default(),
                    timestamp: timestamp__.unwrap_or_default(),
                    metadata: metadata__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.Vector", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::VectorRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.id.is_some() {
            len += 1;
        }
        if self.filters.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.VectorRequest", len)?;
        if let Some(v) = self.id.as_ref() {
            struct_ser.serialize_field("id", v)?;
        }
        if let Some(v) = self.filters.as_ref() {
            struct_ser.serialize_field("filters", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::VectorRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "filters",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Filters,
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
                            "id" => Ok(GeneratedField::Id),
                            "filters" => Ok(GeneratedField::Filters),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::VectorRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.VectorRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::VectorRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut filters__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = map_.next_value()?;
                        }
                        GeneratedField::Filters => {
                            if filters__.is_some() {
                                return Err(serde::de::Error::duplicate_field("filters"));
                            }
                            filters__ = map_.next_value()?;
                        }
                    }
                }
                Ok(object::VectorRequest {
                    id: id__,
                    filters: filters__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.VectorRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for object::Vectors {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.vectors.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Object.Vectors", len)?;
        if !self.vectors.is_empty() {
            struct_ser.serialize_field("vectors", &self.vectors)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for object::Vectors {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "vectors",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Vectors,
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
                            "vectors" => Ok(GeneratedField::Vectors),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = object::Vectors;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Object.Vectors")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<object::Vectors, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut vectors__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Vectors => {
                            if vectors__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vectors"));
                            }
                            vectors__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(object::Vectors {
                    vectors: vectors__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Object.Vectors", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Remove {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Remove", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Remove {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Remove;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Remove")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Remove, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Remove {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Remove", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for remove::Config {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.skip_strict_exist_check {
            len += 1;
        }
        if self.timestamp != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Remove.Config", len)?;
        if self.skip_strict_exist_check {
            struct_ser.serialize_field("skipStrictExistCheck", &self.skip_strict_exist_check)?;
        }
        if self.timestamp != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("timestamp", ToString::to_string(&self.timestamp).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for remove::Config {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "skip_strict_exist_check",
            "skipStrictExistCheck",
            "timestamp",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SkipStrictExistCheck,
            Timestamp,
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
                            "skipStrictExistCheck" | "skip_strict_exist_check" => Ok(GeneratedField::SkipStrictExistCheck),
                            "timestamp" => Ok(GeneratedField::Timestamp),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = remove::Config;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Remove.Config")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<remove::Config, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut skip_strict_exist_check__ = None;
                let mut timestamp__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SkipStrictExistCheck => {
                            if skip_strict_exist_check__.is_some() {
                                return Err(serde::de::Error::duplicate_field("skipStrictExistCheck"));
                            }
                            skip_strict_exist_check__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Timestamp => {
                            if timestamp__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timestamp"));
                            }
                            timestamp__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(remove::Config {
                    skip_strict_exist_check: skip_strict_exist_check__.unwrap_or_default(),
                    timestamp: timestamp__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Remove.Config", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for remove::MultiRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.requests.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Remove.MultiRequest", len)?;
        if !self.requests.is_empty() {
            struct_ser.serialize_field("requests", &self.requests)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for remove::MultiRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "requests",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Requests,
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
                            "requests" => Ok(GeneratedField::Requests),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = remove::MultiRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Remove.MultiRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<remove::MultiRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut requests__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Requests => {
                            if requests__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requests"));
                            }
                            requests__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(remove::MultiRequest {
                    requests: requests__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Remove.MultiRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for remove::Request {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.id.is_some() {
            len += 1;
        }
        if self.config.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Remove.Request", len)?;
        if let Some(v) = self.id.as_ref() {
            struct_ser.serialize_field("id", v)?;
        }
        if let Some(v) = self.config.as_ref() {
            struct_ser.serialize_field("config", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for remove::Request {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "config",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Config,
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
                            "id" => Ok(GeneratedField::Id),
                            "config" => Ok(GeneratedField::Config),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = remove::Request;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Remove.Request")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<remove::Request, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut config__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = map_.next_value()?;
                        }
                        GeneratedField::Config => {
                            if config__.is_some() {
                                return Err(serde::de::Error::duplicate_field("config"));
                            }
                            config__ = map_.next_value()?;
                        }
                    }
                }
                Ok(remove::Request {
                    id: id__,
                    config: config__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Remove.Request", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for remove::Timestamp {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.timestamp != 0 {
            len += 1;
        }
        if self.operator != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Remove.Timestamp", len)?;
        if self.timestamp != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("timestamp", ToString::to_string(&self.timestamp).as_str())?;
        }
        if self.operator != 0 {
            let v = remove::timestamp::Operator::try_from(self.operator)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.operator)))?;
            struct_ser.serialize_field("operator", &v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for remove::Timestamp {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "timestamp",
            "operator",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Timestamp,
            Operator,
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
                            "timestamp" => Ok(GeneratedField::Timestamp),
                            "operator" => Ok(GeneratedField::Operator),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = remove::Timestamp;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Remove.Timestamp")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<remove::Timestamp, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut timestamp__ = None;
                let mut operator__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Timestamp => {
                            if timestamp__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timestamp"));
                            }
                            timestamp__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Operator => {
                            if operator__.is_some() {
                                return Err(serde::de::Error::duplicate_field("operator"));
                            }
                            operator__ = Some(map_.next_value::<remove::timestamp::Operator>()? as i32);
                        }
                    }
                }
                Ok(remove::Timestamp {
                    timestamp: timestamp__.unwrap_or_default(),
                    operator: operator__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Remove.Timestamp", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for remove::timestamp::Operator {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Eq => "Eq",
            Self::Ne => "Ne",
            Self::Ge => "Ge",
            Self::Gt => "Gt",
            Self::Le => "Le",
            Self::Lt => "Lt",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for remove::timestamp::Operator {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "Eq",
            "Ne",
            "Ge",
            "Gt",
            "Le",
            "Lt",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = remove::timestamp::Operator;

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
                    "Eq" => Ok(remove::timestamp::Operator::Eq),
                    "Ne" => Ok(remove::timestamp::Operator::Ne),
                    "Ge" => Ok(remove::timestamp::Operator::Ge),
                    "Gt" => Ok(remove::timestamp::Operator::Gt),
                    "Le" => Ok(remove::timestamp::Operator::Le),
                    "Lt" => Ok(remove::timestamp::Operator::Lt),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for remove::TimestampRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.timestamps.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Remove.TimestampRequest", len)?;
        if !self.timestamps.is_empty() {
            struct_ser.serialize_field("timestamps", &self.timestamps)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for remove::TimestampRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "timestamps",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Timestamps,
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
                            "timestamps" => Ok(GeneratedField::Timestamps),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = remove::TimestampRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Remove.TimestampRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<remove::TimestampRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut timestamps__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Timestamps => {
                            if timestamps__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timestamps"));
                            }
                            timestamps__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(remove::TimestampRequest {
                    timestamps: timestamps__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Remove.TimestampRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Search {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Search", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Search {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Search;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Search")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Search, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Search {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Search", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for search::AggregationAlgorithm {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        let variant = match self {
            Self::Unknown => "Unknown",
            Self::ConcurrentQueue => "ConcurrentQueue",
            Self::SortSlice => "SortSlice",
            Self::SortPoolSlice => "SortPoolSlice",
            Self::PairingHeap => "PairingHeap",
        };
        serializer.serialize_str(variant)
    }
}
impl<'de> serde::Deserialize<'de> for search::AggregationAlgorithm {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "Unknown",
            "ConcurrentQueue",
            "SortSlice",
            "SortPoolSlice",
            "PairingHeap",
        ];

        struct GeneratedVisitor;

        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = search::AggregationAlgorithm;

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
                    "Unknown" => Ok(search::AggregationAlgorithm::Unknown),
                    "ConcurrentQueue" => Ok(search::AggregationAlgorithm::ConcurrentQueue),
                    "SortSlice" => Ok(search::AggregationAlgorithm::SortSlice),
                    "SortPoolSlice" => Ok(search::AggregationAlgorithm::SortPoolSlice),
                    "PairingHeap" => Ok(search::AggregationAlgorithm::PairingHeap),
                    _ => Err(serde::de::Error::unknown_variant(value, FIELDS)),
                }
            }
        }
        deserializer.deserialize_any(GeneratedVisitor)
    }
}
impl serde::Serialize for search::Config {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.request_id.is_empty() {
            len += 1;
        }
        if self.num != 0 {
            len += 1;
        }
        if self.radius != 0. {
            len += 1;
        }
        if self.epsilon != 0. {
            len += 1;
        }
        if self.timeout != 0 {
            len += 1;
        }
        if self.ingress_filters.is_some() {
            len += 1;
        }
        if self.egress_filters.is_some() {
            len += 1;
        }
        if self.min_num != 0 {
            len += 1;
        }
        if self.aggregation_algorithm != 0 {
            len += 1;
        }
        if self.ratio.is_some() {
            len += 1;
        }
        if self.nprobe != 0 {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Search.Config", len)?;
        if !self.request_id.is_empty() {
            struct_ser.serialize_field("requestId", &self.request_id)?;
        }
        if self.num != 0 {
            struct_ser.serialize_field("num", &self.num)?;
        }
        if self.radius != 0. {
            struct_ser.serialize_field("radius", &self.radius)?;
        }
        if self.epsilon != 0. {
            struct_ser.serialize_field("epsilon", &self.epsilon)?;
        }
        if self.timeout != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("timeout", ToString::to_string(&self.timeout).as_str())?;
        }
        if let Some(v) = self.ingress_filters.as_ref() {
            struct_ser.serialize_field("ingressFilters", v)?;
        }
        if let Some(v) = self.egress_filters.as_ref() {
            struct_ser.serialize_field("egressFilters", v)?;
        }
        if self.min_num != 0 {
            struct_ser.serialize_field("minNum", &self.min_num)?;
        }
        if self.aggregation_algorithm != 0 {
            let v = search::AggregationAlgorithm::try_from(self.aggregation_algorithm)
                .map_err(|_| serde::ser::Error::custom(format!("Invalid variant {}", self.aggregation_algorithm)))?;
            struct_ser.serialize_field("aggregationAlgorithm", &v)?;
        }
        if let Some(v) = self.ratio.as_ref() {
            struct_ser.serialize_field("ratio", v)?;
        }
        if self.nprobe != 0 {
            struct_ser.serialize_field("nprobe", &self.nprobe)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for search::Config {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "request_id",
            "requestId",
            "num",
            "radius",
            "epsilon",
            "timeout",
            "ingress_filters",
            "ingressFilters",
            "egress_filters",
            "egressFilters",
            "min_num",
            "minNum",
            "aggregation_algorithm",
            "aggregationAlgorithm",
            "ratio",
            "nprobe",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RequestId,
            Num,
            Radius,
            Epsilon,
            Timeout,
            IngressFilters,
            EgressFilters,
            MinNum,
            AggregationAlgorithm,
            Ratio,
            Nprobe,
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
                            "requestId" | "request_id" => Ok(GeneratedField::RequestId),
                            "num" => Ok(GeneratedField::Num),
                            "radius" => Ok(GeneratedField::Radius),
                            "epsilon" => Ok(GeneratedField::Epsilon),
                            "timeout" => Ok(GeneratedField::Timeout),
                            "ingressFilters" | "ingress_filters" => Ok(GeneratedField::IngressFilters),
                            "egressFilters" | "egress_filters" => Ok(GeneratedField::EgressFilters),
                            "minNum" | "min_num" => Ok(GeneratedField::MinNum),
                            "aggregationAlgorithm" | "aggregation_algorithm" => Ok(GeneratedField::AggregationAlgorithm),
                            "ratio" => Ok(GeneratedField::Ratio),
                            "nprobe" => Ok(GeneratedField::Nprobe),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = search::Config;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Search.Config")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<search::Config, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut request_id__ = None;
                let mut num__ = None;
                let mut radius__ = None;
                let mut epsilon__ = None;
                let mut timeout__ = None;
                let mut ingress_filters__ = None;
                let mut egress_filters__ = None;
                let mut min_num__ = None;
                let mut aggregation_algorithm__ = None;
                let mut ratio__ = None;
                let mut nprobe__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RequestId => {
                            if request_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requestId"));
                            }
                            request_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Num => {
                            if num__.is_some() {
                                return Err(serde::de::Error::duplicate_field("num"));
                            }
                            num__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Radius => {
                            if radius__.is_some() {
                                return Err(serde::de::Error::duplicate_field("radius"));
                            }
                            radius__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Epsilon => {
                            if epsilon__.is_some() {
                                return Err(serde::de::Error::duplicate_field("epsilon"));
                            }
                            epsilon__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Timeout => {
                            if timeout__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timeout"));
                            }
                            timeout__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::IngressFilters => {
                            if ingress_filters__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ingressFilters"));
                            }
                            ingress_filters__ = map_.next_value()?;
                        }
                        GeneratedField::EgressFilters => {
                            if egress_filters__.is_some() {
                                return Err(serde::de::Error::duplicate_field("egressFilters"));
                            }
                            egress_filters__ = map_.next_value()?;
                        }
                        GeneratedField::MinNum => {
                            if min_num__.is_some() {
                                return Err(serde::de::Error::duplicate_field("minNum"));
                            }
                            min_num__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::AggregationAlgorithm => {
                            if aggregation_algorithm__.is_some() {
                                return Err(serde::de::Error::duplicate_field("aggregationAlgorithm"));
                            }
                            aggregation_algorithm__ = Some(map_.next_value::<search::AggregationAlgorithm>()? as i32);
                        }
                        GeneratedField::Ratio => {
                            if ratio__.is_some() {
                                return Err(serde::de::Error::duplicate_field("ratio"));
                            }
                            ratio__ = map_.next_value()?;
                        }
                        GeneratedField::Nprobe => {
                            if nprobe__.is_some() {
                                return Err(serde::de::Error::duplicate_field("nprobe"));
                            }
                            nprobe__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                    }
                }
                Ok(search::Config {
                    request_id: request_id__.unwrap_or_default(),
                    num: num__.unwrap_or_default(),
                    radius: radius__.unwrap_or_default(),
                    epsilon: epsilon__.unwrap_or_default(),
                    timeout: timeout__.unwrap_or_default(),
                    ingress_filters: ingress_filters__,
                    egress_filters: egress_filters__,
                    min_num: min_num__.unwrap_or_default(),
                    aggregation_algorithm: aggregation_algorithm__.unwrap_or_default(),
                    ratio: ratio__,
                    nprobe: nprobe__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Search.Config", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for search::IdRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.id.is_empty() {
            len += 1;
        }
        if self.config.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Search.IDRequest", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if let Some(v) = self.config.as_ref() {
            struct_ser.serialize_field("config", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for search::IdRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "config",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Config,
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
                            "id" => Ok(GeneratedField::Id),
                            "config" => Ok(GeneratedField::Config),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = search::IdRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Search.IDRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<search::IdRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut config__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Config => {
                            if config__.is_some() {
                                return Err(serde::de::Error::duplicate_field("config"));
                            }
                            config__ = map_.next_value()?;
                        }
                    }
                }
                Ok(search::IdRequest {
                    id: id__.unwrap_or_default(),
                    config: config__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Search.IDRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for search::MultiIdRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.requests.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Search.MultiIDRequest", len)?;
        if !self.requests.is_empty() {
            struct_ser.serialize_field("requests", &self.requests)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for search::MultiIdRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "requests",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Requests,
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
                            "requests" => Ok(GeneratedField::Requests),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = search::MultiIdRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Search.MultiIDRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<search::MultiIdRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut requests__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Requests => {
                            if requests__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requests"));
                            }
                            requests__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(search::MultiIdRequest {
                    requests: requests__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Search.MultiIDRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for search::MultiObjectRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.requests.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Search.MultiObjectRequest", len)?;
        if !self.requests.is_empty() {
            struct_ser.serialize_field("requests", &self.requests)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for search::MultiObjectRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "requests",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Requests,
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
                            "requests" => Ok(GeneratedField::Requests),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = search::MultiObjectRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Search.MultiObjectRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<search::MultiObjectRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut requests__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Requests => {
                            if requests__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requests"));
                            }
                            requests__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(search::MultiObjectRequest {
                    requests: requests__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Search.MultiObjectRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for search::MultiRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.requests.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Search.MultiRequest", len)?;
        if !self.requests.is_empty() {
            struct_ser.serialize_field("requests", &self.requests)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for search::MultiRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "requests",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Requests,
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
                            "requests" => Ok(GeneratedField::Requests),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = search::MultiRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Search.MultiRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<search::MultiRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut requests__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Requests => {
                            if requests__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requests"));
                            }
                            requests__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(search::MultiRequest {
                    requests: requests__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Search.MultiRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for search::ObjectRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.object.is_empty() {
            len += 1;
        }
        if self.config.is_some() {
            len += 1;
        }
        if self.vectorizer.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Search.ObjectRequest", len)?;
        if !self.object.is_empty() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("object", pbjson::private::base64::encode(&self.object).as_str())?;
        }
        if let Some(v) = self.config.as_ref() {
            struct_ser.serialize_field("config", v)?;
        }
        if let Some(v) = self.vectorizer.as_ref() {
            struct_ser.serialize_field("vectorizer", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for search::ObjectRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "object",
            "config",
            "vectorizer",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Object,
            Config,
            Vectorizer,
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
                            "object" => Ok(GeneratedField::Object),
                            "config" => Ok(GeneratedField::Config),
                            "vectorizer" => Ok(GeneratedField::Vectorizer),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = search::ObjectRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Search.ObjectRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<search::ObjectRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut object__ = None;
                let mut config__ = None;
                let mut vectorizer__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Object => {
                            if object__.is_some() {
                                return Err(serde::de::Error::duplicate_field("object"));
                            }
                            object__ = 
                                Some(map_.next_value::<::pbjson::private::BytesDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Config => {
                            if config__.is_some() {
                                return Err(serde::de::Error::duplicate_field("config"));
                            }
                            config__ = map_.next_value()?;
                        }
                        GeneratedField::Vectorizer => {
                            if vectorizer__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vectorizer"));
                            }
                            vectorizer__ = map_.next_value()?;
                        }
                    }
                }
                Ok(search::ObjectRequest {
                    object: object__.unwrap_or_default(),
                    config: config__,
                    vectorizer: vectorizer__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Search.ObjectRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for search::Request {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.vector.is_empty() {
            len += 1;
        }
        if self.config.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Search.Request", len)?;
        if !self.vector.is_empty() {
            struct_ser.serialize_field("vector", &self.vector)?;
        }
        if let Some(v) = self.config.as_ref() {
            struct_ser.serialize_field("config", v)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for search::Request {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "vector",
            "config",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Vector,
            Config,
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
                            "vector" => Ok(GeneratedField::Vector),
                            "config" => Ok(GeneratedField::Config),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = search::Request;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Search.Request")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<search::Request, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut vector__ = None;
                let mut config__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Vector => {
                            if vector__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vector"));
                            }
                            vector__ = 
                                Some(map_.next_value::<Vec<::pbjson::private::NumberDeserialize<_>>>()?
                                    .into_iter().map(|x| x.0).collect())
                            ;
                        }
                        GeneratedField::Config => {
                            if config__.is_some() {
                                return Err(serde::de::Error::duplicate_field("config"));
                            }
                            config__ = map_.next_value()?;
                        }
                    }
                }
                Ok(search::Request {
                    vector: vector__.unwrap_or_default(),
                    config: config__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Search.Request", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for search::Response {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.request_id.is_empty() {
            len += 1;
        }
        if !self.results.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Search.Response", len)?;
        if !self.request_id.is_empty() {
            struct_ser.serialize_field("requestId", &self.request_id)?;
        }
        if !self.results.is_empty() {
            struct_ser.serialize_field("results", &self.results)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for search::Response {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "request_id",
            "requestId",
            "results",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            RequestId,
            Results,
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
                            "requestId" | "request_id" => Ok(GeneratedField::RequestId),
                            "results" => Ok(GeneratedField::Results),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = search::Response;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Search.Response")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<search::Response, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut request_id__ = None;
                let mut results__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::RequestId => {
                            if request_id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requestId"));
                            }
                            request_id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Results => {
                            if results__.is_some() {
                                return Err(serde::de::Error::duplicate_field("results"));
                            }
                            results__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(search::Response {
                    request_id: request_id__.unwrap_or_default(),
                    results: results__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Search.Response", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for search::Responses {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.responses.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Search.Responses", len)?;
        if !self.responses.is_empty() {
            struct_ser.serialize_field("responses", &self.responses)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for search::Responses {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "responses",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Responses,
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
                            "responses" => Ok(GeneratedField::Responses),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = search::Responses;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Search.Responses")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<search::Responses, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut responses__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Responses => {
                            if responses__.is_some() {
                                return Err(serde::de::Error::duplicate_field("responses"));
                            }
                            responses__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(search::Responses {
                    responses: responses__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Search.Responses", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for search::StreamResponse {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.payload.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Search.StreamResponse", len)?;
        if let Some(v) = self.payload.as_ref() {
            match v {
                search::stream_response::Payload::Response(v) => {
                    struct_ser.serialize_field("response", v)?;
                }
                search::stream_response::Payload::Status(v) => {
                    struct_ser.serialize_field("status", v)?;
                }
            }
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for search::StreamResponse {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "response",
            "status",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Response,
            Status,
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
                            "response" => Ok(GeneratedField::Response),
                            "status" => Ok(GeneratedField::Status),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = search::StreamResponse;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Search.StreamResponse")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<search::StreamResponse, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut payload__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Response => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("response"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(search::stream_response::Payload::Response)
;
                        }
                        GeneratedField::Status => {
                            if payload__.is_some() {
                                return Err(serde::de::Error::duplicate_field("status"));
                            }
                            payload__ = map_.next_value::<::std::option::Option<_>>()?.map(search::stream_response::Payload::Status)
;
                        }
                    }
                }
                Ok(search::StreamResponse {
                    payload: payload__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Search.StreamResponse", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Update {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Update", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Update {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Update;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Update")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Update, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Update {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Update", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for update::Config {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.skip_strict_exist_check {
            len += 1;
        }
        if self.filters.is_some() {
            len += 1;
        }
        if self.timestamp != 0 {
            len += 1;
        }
        if self.disable_balanced_update {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Update.Config", len)?;
        if self.skip_strict_exist_check {
            struct_ser.serialize_field("skipStrictExistCheck", &self.skip_strict_exist_check)?;
        }
        if let Some(v) = self.filters.as_ref() {
            struct_ser.serialize_field("filters", v)?;
        }
        if self.timestamp != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("timestamp", ToString::to_string(&self.timestamp).as_str())?;
        }
        if self.disable_balanced_update {
            struct_ser.serialize_field("disableBalancedUpdate", &self.disable_balanced_update)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for update::Config {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "skip_strict_exist_check",
            "skipStrictExistCheck",
            "filters",
            "timestamp",
            "disable_balanced_update",
            "disableBalancedUpdate",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SkipStrictExistCheck,
            Filters,
            Timestamp,
            DisableBalancedUpdate,
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
                            "skipStrictExistCheck" | "skip_strict_exist_check" => Ok(GeneratedField::SkipStrictExistCheck),
                            "filters" => Ok(GeneratedField::Filters),
                            "timestamp" => Ok(GeneratedField::Timestamp),
                            "disableBalancedUpdate" | "disable_balanced_update" => Ok(GeneratedField::DisableBalancedUpdate),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = update::Config;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Update.Config")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<update::Config, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut skip_strict_exist_check__ = None;
                let mut filters__ = None;
                let mut timestamp__ = None;
                let mut disable_balanced_update__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SkipStrictExistCheck => {
                            if skip_strict_exist_check__.is_some() {
                                return Err(serde::de::Error::duplicate_field("skipStrictExistCheck"));
                            }
                            skip_strict_exist_check__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Filters => {
                            if filters__.is_some() {
                                return Err(serde::de::Error::duplicate_field("filters"));
                            }
                            filters__ = map_.next_value()?;
                        }
                        GeneratedField::Timestamp => {
                            if timestamp__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timestamp"));
                            }
                            timestamp__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::DisableBalancedUpdate => {
                            if disable_balanced_update__.is_some() {
                                return Err(serde::de::Error::duplicate_field("disableBalancedUpdate"));
                            }
                            disable_balanced_update__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(update::Config {
                    skip_strict_exist_check: skip_strict_exist_check__.unwrap_or_default(),
                    filters: filters__,
                    timestamp: timestamp__.unwrap_or_default(),
                    disable_balanced_update: disable_balanced_update__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Update.Config", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for update::MultiObjectRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.requests.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Update.MultiObjectRequest", len)?;
        if !self.requests.is_empty() {
            struct_ser.serialize_field("requests", &self.requests)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for update::MultiObjectRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "requests",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Requests,
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
                            "requests" => Ok(GeneratedField::Requests),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = update::MultiObjectRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Update.MultiObjectRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<update::MultiObjectRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut requests__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Requests => {
                            if requests__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requests"));
                            }
                            requests__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(update::MultiObjectRequest {
                    requests: requests__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Update.MultiObjectRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for update::MultiRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.requests.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Update.MultiRequest", len)?;
        if !self.requests.is_empty() {
            struct_ser.serialize_field("requests", &self.requests)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for update::MultiRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "requests",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Requests,
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
                            "requests" => Ok(GeneratedField::Requests),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = update::MultiRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Update.MultiRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<update::MultiRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut requests__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Requests => {
                            if requests__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requests"));
                            }
                            requests__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(update::MultiRequest {
                    requests: requests__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Update.MultiRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for update::ObjectRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.object.is_some() {
            len += 1;
        }
        if self.config.is_some() {
            len += 1;
        }
        if self.vectorizer.is_some() {
            len += 1;
        }
        if self.metadata.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Update.ObjectRequest", len)?;
        if let Some(v) = self.object.as_ref() {
            struct_ser.serialize_field("object", v)?;
        }
        if let Some(v) = self.config.as_ref() {
            struct_ser.serialize_field("config", v)?;
        }
        if let Some(v) = self.vectorizer.as_ref() {
            struct_ser.serialize_field("vectorizer", v)?;
        }
        if let Some(v) = self.metadata.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("metadata", pbjson::private::base64::encode(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for update::ObjectRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "object",
            "config",
            "vectorizer",
            "metadata",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Object,
            Config,
            Vectorizer,
            Metadata,
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
                            "object" => Ok(GeneratedField::Object),
                            "config" => Ok(GeneratedField::Config),
                            "vectorizer" => Ok(GeneratedField::Vectorizer),
                            "metadata" => Ok(GeneratedField::Metadata),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = update::ObjectRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Update.ObjectRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<update::ObjectRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut object__ = None;
                let mut config__ = None;
                let mut vectorizer__ = None;
                let mut metadata__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Object => {
                            if object__.is_some() {
                                return Err(serde::de::Error::duplicate_field("object"));
                            }
                            object__ = map_.next_value()?;
                        }
                        GeneratedField::Config => {
                            if config__.is_some() {
                                return Err(serde::de::Error::duplicate_field("config"));
                            }
                            config__ = map_.next_value()?;
                        }
                        GeneratedField::Vectorizer => {
                            if vectorizer__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vectorizer"));
                            }
                            vectorizer__ = map_.next_value()?;
                        }
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::BytesDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                    }
                }
                Ok(update::ObjectRequest {
                    object: object__,
                    config: config__,
                    vectorizer: vectorizer__,
                    metadata: metadata__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Update.ObjectRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for update::Request {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.vector.is_some() {
            len += 1;
        }
        if self.config.is_some() {
            len += 1;
        }
        if self.metadata.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Update.Request", len)?;
        if let Some(v) = self.vector.as_ref() {
            struct_ser.serialize_field("vector", v)?;
        }
        if let Some(v) = self.config.as_ref() {
            struct_ser.serialize_field("config", v)?;
        }
        if let Some(v) = self.metadata.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("metadata", pbjson::private::base64::encode(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for update::Request {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "vector",
            "config",
            "metadata",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Vector,
            Config,
            Metadata,
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
                            "vector" => Ok(GeneratedField::Vector),
                            "config" => Ok(GeneratedField::Config),
                            "metadata" => Ok(GeneratedField::Metadata),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = update::Request;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Update.Request")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<update::Request, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut vector__ = None;
                let mut config__ = None;
                let mut metadata__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Vector => {
                            if vector__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vector"));
                            }
                            vector__ = map_.next_value()?;
                        }
                        GeneratedField::Config => {
                            if config__.is_some() {
                                return Err(serde::de::Error::duplicate_field("config"));
                            }
                            config__ = map_.next_value()?;
                        }
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::BytesDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                    }
                }
                Ok(update::Request {
                    vector: vector__,
                    config: config__,
                    metadata: metadata__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Update.Request", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for update::TimestampRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.id.is_empty() {
            len += 1;
        }
        if self.timestamp != 0 {
            len += 1;
        }
        if self.force {
            len += 1;
        }
        if self.metadata.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Update.TimestampRequest", len)?;
        if !self.id.is_empty() {
            struct_ser.serialize_field("id", &self.id)?;
        }
        if self.timestamp != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("timestamp", ToString::to_string(&self.timestamp).as_str())?;
        }
        if self.force {
            struct_ser.serialize_field("force", &self.force)?;
        }
        if let Some(v) = self.metadata.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("metadata", pbjson::private::base64::encode(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for update::TimestampRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "id",
            "timestamp",
            "force",
            "metadata",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Id,
            Timestamp,
            Force,
            Metadata,
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
                            "id" => Ok(GeneratedField::Id),
                            "timestamp" => Ok(GeneratedField::Timestamp),
                            "force" => Ok(GeneratedField::Force),
                            "metadata" => Ok(GeneratedField::Metadata),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = update::TimestampRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Update.TimestampRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<update::TimestampRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut id__ = None;
                let mut timestamp__ = None;
                let mut force__ = None;
                let mut metadata__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Id => {
                            if id__.is_some() {
                                return Err(serde::de::Error::duplicate_field("id"));
                            }
                            id__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Timestamp => {
                            if timestamp__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timestamp"));
                            }
                            timestamp__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::Force => {
                            if force__.is_some() {
                                return Err(serde::de::Error::duplicate_field("force"));
                            }
                            force__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::BytesDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                    }
                }
                Ok(update::TimestampRequest {
                    id: id__.unwrap_or_default(),
                    timestamp: timestamp__.unwrap_or_default(),
                    force: force__.unwrap_or_default(),
                    metadata: metadata__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Update.TimestampRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for Upsert {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let len = 0;
        let struct_ser = serializer.serialize_struct("payload.v1.Upsert", len)?;
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for Upsert {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
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
                            Err(serde::de::Error::unknown_field(value, FIELDS))
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = Upsert;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Upsert")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<Upsert, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                while map_.next_key::<GeneratedField>()?.is_some() {
                    let _ = map_.next_value::<serde::de::IgnoredAny>()?;
                }
                Ok(Upsert {
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Upsert", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for upsert::Config {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.skip_strict_exist_check {
            len += 1;
        }
        if self.filters.is_some() {
            len += 1;
        }
        if self.timestamp != 0 {
            len += 1;
        }
        if self.disable_balanced_update {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Upsert.Config", len)?;
        if self.skip_strict_exist_check {
            struct_ser.serialize_field("skipStrictExistCheck", &self.skip_strict_exist_check)?;
        }
        if let Some(v) = self.filters.as_ref() {
            struct_ser.serialize_field("filters", v)?;
        }
        if self.timestamp != 0 {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("timestamp", ToString::to_string(&self.timestamp).as_str())?;
        }
        if self.disable_balanced_update {
            struct_ser.serialize_field("disableBalancedUpdate", &self.disable_balanced_update)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for upsert::Config {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "skip_strict_exist_check",
            "skipStrictExistCheck",
            "filters",
            "timestamp",
            "disable_balanced_update",
            "disableBalancedUpdate",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            SkipStrictExistCheck,
            Filters,
            Timestamp,
            DisableBalancedUpdate,
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
                            "skipStrictExistCheck" | "skip_strict_exist_check" => Ok(GeneratedField::SkipStrictExistCheck),
                            "filters" => Ok(GeneratedField::Filters),
                            "timestamp" => Ok(GeneratedField::Timestamp),
                            "disableBalancedUpdate" | "disable_balanced_update" => Ok(GeneratedField::DisableBalancedUpdate),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = upsert::Config;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Upsert.Config")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<upsert::Config, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut skip_strict_exist_check__ = None;
                let mut filters__ = None;
                let mut timestamp__ = None;
                let mut disable_balanced_update__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::SkipStrictExistCheck => {
                            if skip_strict_exist_check__.is_some() {
                                return Err(serde::de::Error::duplicate_field("skipStrictExistCheck"));
                            }
                            skip_strict_exist_check__ = Some(map_.next_value()?);
                        }
                        GeneratedField::Filters => {
                            if filters__.is_some() {
                                return Err(serde::de::Error::duplicate_field("filters"));
                            }
                            filters__ = map_.next_value()?;
                        }
                        GeneratedField::Timestamp => {
                            if timestamp__.is_some() {
                                return Err(serde::de::Error::duplicate_field("timestamp"));
                            }
                            timestamp__ = 
                                Some(map_.next_value::<::pbjson::private::NumberDeserialize<_>>()?.0)
                            ;
                        }
                        GeneratedField::DisableBalancedUpdate => {
                            if disable_balanced_update__.is_some() {
                                return Err(serde::de::Error::duplicate_field("disableBalancedUpdate"));
                            }
                            disable_balanced_update__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(upsert::Config {
                    skip_strict_exist_check: skip_strict_exist_check__.unwrap_or_default(),
                    filters: filters__,
                    timestamp: timestamp__.unwrap_or_default(),
                    disable_balanced_update: disable_balanced_update__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Upsert.Config", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for upsert::MultiObjectRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.requests.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Upsert.MultiObjectRequest", len)?;
        if !self.requests.is_empty() {
            struct_ser.serialize_field("requests", &self.requests)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for upsert::MultiObjectRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "requests",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Requests,
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
                            "requests" => Ok(GeneratedField::Requests),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = upsert::MultiObjectRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Upsert.MultiObjectRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<upsert::MultiObjectRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut requests__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Requests => {
                            if requests__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requests"));
                            }
                            requests__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(upsert::MultiObjectRequest {
                    requests: requests__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Upsert.MultiObjectRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for upsert::MultiRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if !self.requests.is_empty() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Upsert.MultiRequest", len)?;
        if !self.requests.is_empty() {
            struct_ser.serialize_field("requests", &self.requests)?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for upsert::MultiRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "requests",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Requests,
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
                            "requests" => Ok(GeneratedField::Requests),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = upsert::MultiRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Upsert.MultiRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<upsert::MultiRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut requests__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Requests => {
                            if requests__.is_some() {
                                return Err(serde::de::Error::duplicate_field("requests"));
                            }
                            requests__ = Some(map_.next_value()?);
                        }
                    }
                }
                Ok(upsert::MultiRequest {
                    requests: requests__.unwrap_or_default(),
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Upsert.MultiRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for upsert::ObjectRequest {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.object.is_some() {
            len += 1;
        }
        if self.config.is_some() {
            len += 1;
        }
        if self.vectorizer.is_some() {
            len += 1;
        }
        if self.metadata.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Upsert.ObjectRequest", len)?;
        if let Some(v) = self.object.as_ref() {
            struct_ser.serialize_field("object", v)?;
        }
        if let Some(v) = self.config.as_ref() {
            struct_ser.serialize_field("config", v)?;
        }
        if let Some(v) = self.vectorizer.as_ref() {
            struct_ser.serialize_field("vectorizer", v)?;
        }
        if let Some(v) = self.metadata.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("metadata", pbjson::private::base64::encode(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for upsert::ObjectRequest {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "object",
            "config",
            "vectorizer",
            "metadata",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Object,
            Config,
            Vectorizer,
            Metadata,
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
                            "object" => Ok(GeneratedField::Object),
                            "config" => Ok(GeneratedField::Config),
                            "vectorizer" => Ok(GeneratedField::Vectorizer),
                            "metadata" => Ok(GeneratedField::Metadata),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = upsert::ObjectRequest;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Upsert.ObjectRequest")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<upsert::ObjectRequest, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut object__ = None;
                let mut config__ = None;
                let mut vectorizer__ = None;
                let mut metadata__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Object => {
                            if object__.is_some() {
                                return Err(serde::de::Error::duplicate_field("object"));
                            }
                            object__ = map_.next_value()?;
                        }
                        GeneratedField::Config => {
                            if config__.is_some() {
                                return Err(serde::de::Error::duplicate_field("config"));
                            }
                            config__ = map_.next_value()?;
                        }
                        GeneratedField::Vectorizer => {
                            if vectorizer__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vectorizer"));
                            }
                            vectorizer__ = map_.next_value()?;
                        }
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::BytesDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                    }
                }
                Ok(upsert::ObjectRequest {
                    object: object__,
                    config: config__,
                    vectorizer: vectorizer__,
                    metadata: metadata__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Upsert.ObjectRequest", FIELDS, GeneratedVisitor)
    }
}
impl serde::Serialize for upsert::Request {
    #[allow(deprecated)]
    fn serialize<S>(&self, serializer: S) -> std::result::Result<S::Ok, S::Error>
    where
        S: serde::Serializer,
    {
        use serde::ser::SerializeStruct;
        let mut len = 0;
        if self.vector.is_some() {
            len += 1;
        }
        if self.config.is_some() {
            len += 1;
        }
        if self.metadata.is_some() {
            len += 1;
        }
        let mut struct_ser = serializer.serialize_struct("payload.v1.Upsert.Request", len)?;
        if let Some(v) = self.vector.as_ref() {
            struct_ser.serialize_field("vector", v)?;
        }
        if let Some(v) = self.config.as_ref() {
            struct_ser.serialize_field("config", v)?;
        }
        if let Some(v) = self.metadata.as_ref() {
            #[allow(clippy::needless_borrow)]
            #[allow(clippy::needless_borrows_for_generic_args)]
            struct_ser.serialize_field("metadata", pbjson::private::base64::encode(&v).as_str())?;
        }
        struct_ser.end()
    }
}
impl<'de> serde::Deserialize<'de> for upsert::Request {
    #[allow(deprecated)]
    fn deserialize<D>(deserializer: D) -> std::result::Result<Self, D::Error>
    where
        D: serde::Deserializer<'de>,
    {
        const FIELDS: &[&str] = &[
            "vector",
            "config",
            "metadata",
        ];

        #[allow(clippy::enum_variant_names)]
        enum GeneratedField {
            Vector,
            Config,
            Metadata,
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
                            "vector" => Ok(GeneratedField::Vector),
                            "config" => Ok(GeneratedField::Config),
                            "metadata" => Ok(GeneratedField::Metadata),
                            _ => Err(serde::de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }
                deserializer.deserialize_identifier(GeneratedVisitor)
            }
        }
        struct GeneratedVisitor;
        impl<'de> serde::de::Visitor<'de> for GeneratedVisitor {
            type Value = upsert::Request;

            fn expecting(&self, formatter: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
                formatter.write_str("struct payload.v1.Upsert.Request")
            }

            fn visit_map<V>(self, mut map_: V) -> std::result::Result<upsert::Request, V::Error>
                where
                    V: serde::de::MapAccess<'de>,
            {
                let mut vector__ = None;
                let mut config__ = None;
                let mut metadata__ = None;
                while let Some(k) = map_.next_key()? {
                    match k {
                        GeneratedField::Vector => {
                            if vector__.is_some() {
                                return Err(serde::de::Error::duplicate_field("vector"));
                            }
                            vector__ = map_.next_value()?;
                        }
                        GeneratedField::Config => {
                            if config__.is_some() {
                                return Err(serde::de::Error::duplicate_field("config"));
                            }
                            config__ = map_.next_value()?;
                        }
                        GeneratedField::Metadata => {
                            if metadata__.is_some() {
                                return Err(serde::de::Error::duplicate_field("metadata"));
                            }
                            metadata__ = 
                                map_.next_value::<::std::option::Option<::pbjson::private::BytesDeserialize<_>>>()?.map(|x| x.0)
                            ;
                        }
                    }
                }
                Ok(upsert::Request {
                    vector: vector__,
                    config: config__,
                    metadata: metadata__,
                })
            }
        }
        deserializer.deserialize_struct("payload.v1.Upsert.Request", FIELDS, GeneratedVisitor)
    }
}
