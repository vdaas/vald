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

use kv::*;
use defer::defer;
use opentelemetry::{trace::{Tracer, TraceContextExt}, KeyValue, Context};
use observability::{ctx_span, tracer};
use proto::{meta::v1::meta_server, payload::v1::{meta, Empty}};

#[tonic::async_trait]
impl meta_server::Meta for super::Meta {
    async fn get(
        &self,
        request: tonic::Request<meta::Key>,
    ) -> std::result::Result<tonic::Response<meta::Value>, tonic::Status> {
        let parent_cx = request.extensions().get::<Context>().cloned().unwrap_or_else(Context::new);
        let ctx = ctx_span!(&parent_cx, "Meta::get");
        defer!(ctx.span().end());

        let key = request.into_inner().key;
        let raw_key = Raw::from(key.as_bytes());

        match self.bucket.get(&raw_key) {
            Ok(Some(value_bytes)) => {
                ctx.span().add_event("Key found", vec![KeyValue::new("key", key.clone())]);

                let any_value = prost_types::Any {
                    type_url: "type.googleapis.com/your.package.MessageType".to_string(),
                    value: value_bytes.to_vec(),
                };
                let response = meta::Value {
                    value: Some(any_value),
                };

                Ok(tonic::Response::new(response))
            },
            Ok(None) => {
                ctx.span().add_event("Key not found", vec![KeyValue::new("key", key)]);
                Err(tonic::Status::not_found("Key not found"))
            }
            Err(e) => {
                ctx.span().add_event("Database error", vec![KeyValue::new("error", e.to_string())]);
                Err(tonic::Status::internal(format!("Database error: {}", e)))
            }
        }
    }

    async fn set(
        &self,
        request: tonic::Request<meta::KeyValue>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        let parent_cx = request.extensions().get::<Context>().cloned().unwrap_or_else(Context::new);
        let ctx = ctx_span!(&parent_cx, "Meta::set");
        defer!(ctx.span().end());

        let key_value = request.into_inner();

        let key = match key_value.key {
            Some(k) => k.key,
            None => {
                ctx.span().add_event("Invalid argument", vec![KeyValue::new("error", "Key is missing")]);
                return Err(tonic::Status::invalid_argument("Key is missing"));
            }
        };

        let value = match key_value.value {
            Some(v) => match v.value {
                Some(any_value) => any_value.value,
                None => {
                    ctx.span().add_event("Invalid argument", vec![KeyValue::new("error", "Value is missing")]);
                    return Err(tonic::Status::invalid_argument("Value is missing"));
                }
            },
            None => {
                ctx.span().add_event("Invalid argument", vec![KeyValue::new("error", "Value is missing")]);
                return Err(tonic::Status::invalid_argument("Value is missing"));
            }
        };

        let raw_key = Raw::from(key.as_bytes());
        let raw_value = sled::IVec::from(value);

        match self.bucket.set(&raw_key, &raw_value) {
            Ok(_) => {
                ctx.span().add_event("Value set successfully", vec![KeyValue::new("key", key)]);
                Ok(tonic::Response::new(Empty {}))
            },
            Err(e) => {
                ctx.span().add_event("Failed to set value", vec![KeyValue::new("error", e.to_string())]);
                Err(tonic::Status::internal(format!("Failed to set value: {}", e)))
            }
        }
    }

    async fn delete(
        &self,
        request: tonic::Request<meta::Key>,
    ) -> std::result::Result<tonic::Response<Empty>, tonic::Status> {
        let parent_cx = request.extensions().get::<Context>().cloned().unwrap_or_else(Context::new);
        let ctx = ctx_span!(&parent_cx, "Meta::delete");
        defer!(ctx.span().end());

        let key = request.into_inner().key;
        let raw_key = Raw::from(key.as_bytes());

        match self.bucket.remove(&raw_key) {
            Ok(_) => {
                ctx.span().add_event("Key deleted successfully", vec![KeyValue::new("key", key)]);
                Ok(tonic::Response::new(Empty {}))
            },
            Err(e) => {
                ctx.span().add_event("Failed to delete key", vec![KeyValue::new("error", e.to_string())]);
                Err(tonic::Status::internal(format!("Failed to delete key: {}", e)))
            }
        }
    }
}
