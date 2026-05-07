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
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Text {
    #[prost(string, tag="1")]
    pub text: ::prost::alloc::string::String,
}
impl ::prost::Name for Text {
const NAME: &'static str = "Text";
const PACKAGE: &'static str = "embedder.v1";
fn full_name() -> ::prost::alloc::string::String { "embedder.v1.Text".into() }fn type_url() -> ::prost::alloc::string::String { "/embedder.v1.Text".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct Document {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(string, tag="2")]
    pub text: ::prost::alloc::string::String,
    #[prost(int64, tag="3")]
    pub timestamp: i64,
}
impl ::prost::Name for Document {
const NAME: &'static str = "Document";
const PACKAGE: &'static str = "embedder.v1";
fn full_name() -> ::prost::alloc::string::String { "embedder.v1.Document".into() }fn type_url() -> ::prost::alloc::string::String { "/embedder.v1.Document".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SearchRequest {
    #[prost(string, tag="1")]
    pub text: ::prost::alloc::string::String,
    #[prost(message, optional, tag="2")]
    pub config: ::core::option::Option<super::super::payload::v1::search::Config>,
}
impl ::prost::Name for SearchRequest {
const NAME: &'static str = "SearchRequest";
const PACKAGE: &'static str = "embedder.v1";
fn full_name() -> ::prost::alloc::string::String { "embedder.v1.SearchRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/embedder.v1.SearchRequest".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct InsertRequest {
    #[prost(message, optional, tag="1")]
    pub document: ::core::option::Option<Document>,
    #[prost(message, optional, tag="2")]
    pub config: ::core::option::Option<super::super::payload::v1::insert::Config>,
}
impl ::prost::Name for InsertRequest {
const NAME: &'static str = "InsertRequest";
const PACKAGE: &'static str = "embedder.v1";
fn full_name() -> ::prost::alloc::string::String { "embedder.v1.InsertRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/embedder.v1.InsertRequest".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct InsertWithMetadataRequest {
    #[prost(message, optional, tag="1")]
    pub request: ::core::option::Option<InsertRequest>,
    #[prost(message, optional, tag="2")]
    pub metadata: ::core::option::Option<super::super::payload::v1::meta::Value>,
}
impl ::prost::Name for InsertWithMetadataRequest {
const NAME: &'static str = "InsertWithMetadataRequest";
const PACKAGE: &'static str = "embedder.v1";
fn full_name() -> ::prost::alloc::string::String { "embedder.v1.InsertWithMetadataRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/embedder.v1.InsertWithMetadataRequest".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpdateRequest {
    #[prost(message, optional, tag="1")]
    pub document: ::core::option::Option<Document>,
    #[prost(message, optional, tag="2")]
    pub config: ::core::option::Option<super::super::payload::v1::update::Config>,
}
impl ::prost::Name for UpdateRequest {
const NAME: &'static str = "UpdateRequest";
const PACKAGE: &'static str = "embedder.v1";
fn full_name() -> ::prost::alloc::string::String { "embedder.v1.UpdateRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/embedder.v1.UpdateRequest".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpdateWithMetadataRequest {
    #[prost(message, optional, tag="1")]
    pub request: ::core::option::Option<UpdateRequest>,
    #[prost(message, optional, tag="2")]
    pub metadata: ::core::option::Option<super::super::payload::v1::meta::Value>,
}
impl ::prost::Name for UpdateWithMetadataRequest {
const NAME: &'static str = "UpdateWithMetadataRequest";
const PACKAGE: &'static str = "embedder.v1";
fn full_name() -> ::prost::alloc::string::String { "embedder.v1.UpdateWithMetadataRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/embedder.v1.UpdateWithMetadataRequest".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpsertRequest {
    #[prost(message, optional, tag="1")]
    pub document: ::core::option::Option<Document>,
    #[prost(message, optional, tag="2")]
    pub config: ::core::option::Option<super::super::payload::v1::upsert::Config>,
}
impl ::prost::Name for UpsertRequest {
const NAME: &'static str = "UpsertRequest";
const PACKAGE: &'static str = "embedder.v1";
fn full_name() -> ::prost::alloc::string::String { "embedder.v1.UpsertRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/embedder.v1.UpsertRequest".into() }}
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct UpsertWithMetadataRequest {
    #[prost(message, optional, tag="1")]
    pub request: ::core::option::Option<UpsertRequest>,
    #[prost(message, optional, tag="2")]
    pub metadata: ::core::option::Option<super::super::payload::v1::meta::Value>,
}
impl ::prost::Name for UpsertWithMetadataRequest {
const NAME: &'static str = "UpsertWithMetadataRequest";
const PACKAGE: &'static str = "embedder.v1";
fn full_name() -> ::prost::alloc::string::String { "embedder.v1.UpsertWithMetadataRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/embedder.v1.UpsertWithMetadataRequest".into() }}
#[derive(Clone, PartialEq, Eq, Hash, ::prost::Message)]
pub struct RemoveRequest {
    #[prost(string, tag="1")]
    pub id: ::prost::alloc::string::String,
    #[prost(message, optional, tag="2")]
    pub config: ::core::option::Option<super::super::payload::v1::remove::Config>,
}
impl ::prost::Name for RemoveRequest {
const NAME: &'static str = "RemoveRequest";
const PACKAGE: &'static str = "embedder.v1";
fn full_name() -> ::prost::alloc::string::String { "embedder.v1.RemoveRequest".into() }fn type_url() -> ::prost::alloc::string::String { "/embedder.v1.RemoveRequest".into() }}
include!("embedder.v1.serde.rs");
// @@protoc_insertion_point(module)
