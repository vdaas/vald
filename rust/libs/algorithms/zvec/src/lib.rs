// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

use algorithm::{ANN, Error, MultiError};
use prost::Message;
use proto::google::protobuf::Any;
use proto::payload::v1::{object::Distance, search, z_vec};
use std::collections::HashMap;
use std::path::Path;
use std::sync::OnceLock;
use zvec_rust::{
    Collection, CollectionSchema, DataType, Doc, ErrorCode, Fts, IndexParams, MetricType,
    MultiQuery, SearchQuery, SubQuery,
};

const VECTOR_FIELD: &str = "embedding";
const TIMESTAMP_FIELD: &str = "timestamp";
static INITIALIZED: OnceLock<Result<(), String>> = OnceLock::new();

pub struct Service {
    collection: Collection,
    dimension: usize,
    fields: Vec<String>,
}

impl Service {
    pub fn new(
        path: &str,
        dimension: usize,
        fields: Vec<String>,
        fts_fields: Vec<String>,
    ) -> Result<Self, Error> {
        if dimension < 2 {
            return Err(Error::InvalidDimensionSize {
                uuid: String::new(),
                current: dimension.to_string(),
                limit: "0".to_string(),
            });
        }
        match INITIALIZED.get_or_init(|| zvec_rust::initialize(None).map_err(|err| err.to_string()))
        {
            Ok(()) => {}
            Err(message) => {
                return Err(Error::Backend {
                    message: message.clone(),
                });
            }
        }

        let collection = if Path::new(path).exists() {
            Collection::open(path, None).map_err(backend_error)?
        } else {
            let mut schema = CollectionSchema::builder("vald_zvec")
                .add_field(
                    zvec_rust::FieldSchema::new(TIMESTAMP_FIELD, DataType::Int64, false, 0)
                        .map_err(backend_error)?,
                )
                .add_vector_field(
                    VECTOR_FIELD,
                    DataType::VectorFp32,
                    dimension as u32,
                    IndexParams::hnsw(MetricType::L2, 16, 200).map_err(backend_error)?,
                );
            for field in &fields {
                let mut field_schema =
                    zvec_rust::FieldSchema::new(field, DataType::String, true, 0)
                        .map_err(backend_error)?;
                field_schema
                    .set_index_params(&IndexParams::invert(false, false).map_err(backend_error)?)
                    .map_err(backend_error)?;
                schema = schema.add_field(field_schema);
            }
            for field in &fts_fields {
                let mut field_schema =
                    zvec_rust::FieldSchema::new(field, DataType::String, true, 0)
                        .map_err(backend_error)?;
                field_schema
                    .set_index_params(&IndexParams::fts(None, None, None).map_err(backend_error)?)
                    .map_err(backend_error)?;
                schema = schema.add_field(field_schema);
            }
            let schema = schema.build().map_err(backend_error)?;
            Collection::create_and_open(path, &schema, None).map_err(backend_error)?
        };
        Ok(Self {
            collection,
            dimension,
            fields: fields.into_iter().chain(fts_fields).collect(),
        })
    }

    fn document(
        &self,
        uuid: &str,
        vector: &[f32],
        timestamp: i64,
        options: &[Any],
    ) -> Result<Doc, Error> {
        if uuid.is_empty() {
            return Err(Error::InvalidUUID {
                uuid: uuid.to_string(),
            });
        }
        self.validate_dimension(uuid, vector)?;
        let options = decode_any::<z_vec::DocumentOptions>(options)?;
        let mut doc = Doc::new().map_err(backend_error)?;
        doc.set_pk(uuid);
        doc.add_i64(TIMESTAMP_FIELD, timestamp)
            .map_err(backend_error)?;
        doc.add_vector_f32(VECTOR_FIELD, vector)
            .map_err(backend_error)?;
        if let Some(options) = options {
            for (name, value) in options.fields {
                if self.fields.iter().any(|field| field == &name) {
                    doc.add_string(&name, &value).map_err(backend_error)?;
                }
            }
        }
        Ok(doc)
    }

    fn validate_dimension(&self, uuid: &str, vector: &[f32]) -> Result<(), Error> {
        if vector.len() != self.dimension {
            return Err(Error::InvalidDimensionSize {
                uuid: uuid.to_string(),
                current: vector.len().to_string(),
                limit: self.dimension.to_string(),
            });
        }
        Ok(())
    }

    fn write_result(result: zvec_rust::WriteResult, uuids: &[&str]) -> Result<(), Error> {
        if result.error_count == 0 {
            return Ok(());
        }
        let (index, failed) = result
            .results
            .into_iter()
            .enumerate()
            .find(|(_, result)| !result.is_success())
            .ok_or_else(|| Error::Backend {
                message: "zvec write failed".to_string(),
            })?;
        let uuid = uuids.get(index).copied().unwrap_or_default().to_string();
        match failed.code {
            ErrorCode::AlreadyExists => Err(Error::UUIDAlreadyExists { uuid }),
            ErrorCode::NotFound => Err(Error::ObjectIDNotFound { uuid }),
            _ => Err(Error::Backend {
                message: failed.message,
            }),
        }
    }

    fn search_documents(
        &self,
        vector: &[f32],
        k: u32,
        options: Option<z_vec::SearchOptions>,
    ) -> Result<Vec<Doc>, Error> {
        let Some(options) = options else {
            let query = SearchQuery::new(VECTOR_FIELD, vector, k as i32).map_err(backend_error)?;
            return self.collection.query(&query).map_err(backend_error);
        };
        if options.hybrid_queries.is_empty() {
            let mut builder = SearchQuery::builder()
                .field_name(VECTOR_FIELD)
                .vector(vector)
                .topk(k as i32);
            if !options.pre_filter.is_empty() {
                builder = builder.filter(&options.pre_filter);
            }
            let query = builder.build().map_err(backend_error)?;
            return self.collection.query(&query).map_err(backend_error);
        }
        if options.hybrid_queries.len() == 1 {
            let item = &options.hybrid_queries[0];
            let mut query = match &item.query {
                Some(z_vec::query::Query::Vector(vector)) => {
                    SearchQuery::new(&item.field_name, &vector.values, k as i32)
                        .map_err(backend_error)?
                }
                Some(z_vec::query::Query::Fts(fts)) => {
                    let mut value = Fts::new().map_err(backend_error)?;
                    if !fts.match_string.is_empty() {
                        value
                            .set_match_string(&fts.match_string)
                            .map_err(backend_error)?;
                    }
                    if !fts.query_string.is_empty() {
                        value
                            .set_query_string(&fts.query_string)
                            .map_err(backend_error)?;
                    }
                    SearchQuery::fts(&item.field_name, &value, k as i32).map_err(backend_error)?
                }
                None => {
                    return Err(Error::Backend {
                        message: "Zvec hybrid query has no payload".to_string(),
                    });
                }
            };
            if !options.pre_filter.is_empty() {
                query
                    .set_filter(&options.pre_filter)
                    .map_err(backend_error)?;
            }
            return self.collection.query(&query).map_err(backend_error);
        }

        let mut query = MultiQuery::new().map_err(backend_error)?;
        query.set_topk(k as i32).map_err(backend_error)?;
        if !options.pre_filter.is_empty() {
            query
                .set_filter(&options.pre_filter)
                .map_err(backend_error)?;
        }
        for item in options.hybrid_queries {
            let mut sub = SubQuery::new().map_err(backend_error)?;
            sub.set_field_name(&item.field_name)
                .map_err(backend_error)?;
            sub.set_num_candidates((k.max(10) * 5) as i32)
                .map_err(backend_error)?;
            match item.query {
                Some(z_vec::query::Query::Vector(vector)) => sub
                    .set_query_vector(&vector.values)
                    .map_err(backend_error)?,
                Some(z_vec::query::Query::Fts(fts)) => {
                    let mut value = Fts::new().map_err(backend_error)?;
                    if !fts.match_string.is_empty() {
                        value
                            .set_match_string(&fts.match_string)
                            .map_err(backend_error)?;
                    }
                    if !fts.query_string.is_empty() {
                        value
                            .set_query_string(&fts.query_string)
                            .map_err(backend_error)?;
                    }
                    sub.set_fts(&value).map_err(backend_error)?;
                }
                None => continue,
            }
            query.add_sub_query(&sub).map_err(backend_error)?;
        }
        if !options.hybrid_weights.is_empty() {
            let weights = options
                .hybrid_weights
                .iter()
                .map(|weight| f64::from(*weight))
                .collect::<Vec<_>>();
            query.set_rerank_weighted(&weights).map_err(backend_error)?;
        } else {
            query.set_rerank_rrf(60).map_err(backend_error)?;
        }
        self.collection.multi_query(&query).map_err(backend_error)
    }
}

impl ANN for Service {
    fn exists(&self, uuid: String) -> bool {
        self.collection
            .fetch_with_options(&[uuid.as_str()], None, false)
            .map(|docs| !docs.is_empty())
            .unwrap_or(false)
    }

    fn create_index(&mut self) -> Result<(), Error> {
        self.collection.optimize().map_err(backend_error)
    }

    fn save_index(&mut self) -> Result<(), Error> {
        self.collection.flush().map_err(backend_error)
    }

    fn insert(&mut self, uuid: String, vector: Vec<f32>, ts: i64) -> Result<(), Error> {
        self.insert_with_options(uuid, vector, ts, &[])
    }

    fn insert_with_options(
        &mut self,
        uuid: String,
        vector: Vec<f32>,
        ts: i64,
        options: &[Any],
    ) -> Result<(), Error> {
        let doc = self.document(&uuid, &vector, ts, options)?;
        Self::write_result(
            self.collection.insert(&[&doc]).map_err(backend_error)?,
            &[uuid.as_str()],
        )
    }

    fn insert_multiple(&mut self, vectors: HashMap<String, Vec<f32>>) -> Result<(), Error> {
        let mut duplicated = Vec::new();
        for (uuid, vector) in vectors {
            if let Err(err) = self.insert(uuid, vector, 0) {
                match err {
                    Error::UUIDAlreadyExists { uuid } => duplicated.push(uuid),
                    _ => return Err(err),
                }
            }
        }
        if duplicated.is_empty() {
            Ok(())
        } else {
            Err(Error::new_uuid_already_exists(duplicated))
        }
    }

    fn update(&mut self, uuid: String, vector: Vec<f32>, ts: i64) -> Result<(), Error> {
        let doc = self.document(&uuid, &vector, ts, &[])?;
        Self::write_result(
            self.collection.update(&[&doc]).map_err(backend_error)?,
            &[uuid.as_str()],
        )
    }

    fn update_multiple(&mut self, vectors: HashMap<String, Vec<f32>>) -> Result<(), Error> {
        for (uuid, vector) in vectors {
            self.update(uuid, vector, 0)?;
        }
        Ok(())
    }

    fn ready_for_update(&mut self, uuid: String, vector: Vec<f32>, ts: i64) -> Result<(), Error> {
        self.validate_dimension(&uuid, &vector)?;
        if !self.exists(uuid.clone()) {
            return Err(Error::UUIDNotFound { uuid });
        }
        self.update(uuid, vector, ts)
    }

    fn remove(&mut self, uuid: String, _ts: i64) -> Result<(), Error> {
        Self::write_result(
            self.collection
                .delete(&[uuid.as_str()])
                .map_err(backend_error)?,
            &[uuid.as_str()],
        )
    }

    fn remove_multiple(&mut self, uuids: Vec<String>) -> Result<(), Error> {
        let values = uuids.iter().map(String::as_str).collect::<Vec<_>>();
        Self::write_result(
            self.collection.delete(&values).map_err(backend_error)?,
            &values,
        )
    }

    fn search(
        &self,
        vector: Vec<f32>,
        k: u32,
        epsilon: f32,
        radius: f32,
    ) -> Result<search::Response, Error> {
        self.search_with_options(vector, k, epsilon, radius, &[])
    }

    fn search_with_options(
        &self,
        vector: Vec<f32>,
        k: u32,
        _epsilon: f32,
        _radius: f32,
        options: &[Any],
    ) -> Result<search::Response, Error> {
        self.validate_dimension(String::new().as_str(), &vector)?;
        let options = decode_any::<z_vec::SearchOptions>(options)?;
        let docs = self.search_documents(&vector, k, options)?;
        Ok(search::Response {
            request_id: String::new(),
            results: docs
                .into_iter()
                .filter_map(|doc| {
                    doc.get_pk().map(|id| Distance {
                        id: id.to_string(),
                        distance: doc.get_score(),
                    })
                })
                .collect(),
        })
    }

    fn get_object(&self, uuid: String) -> Result<(Vec<f32>, i64), Error> {
        let mut docs = self
            .collection
            .fetch(&[uuid.as_str()])
            .map_err(backend_error)?;
        let doc = docs.pop().ok_or(Error::ObjectIDNotFound { uuid })?;
        let vector = doc
            .get_vector_f32(VECTOR_FIELD)
            .map_err(backend_error)?
            .ok_or_else(|| Error::Backend {
                message: "zvec document has no embedding".to_string(),
            })?;
        let timestamp = doc
            .get_i64(TIMESTAMP_FIELD)
            .map_err(backend_error)?
            .unwrap_or_default();
        Ok((vector, timestamp))
    }

    fn get_dimension_size(&self) -> usize {
        self.dimension
    }

    fn len(&self) -> u32 {
        self.collection
            .stats()
            .map(|stats| stats.doc_count.min(u64::from(u32::MAX)) as u32)
            .unwrap_or_default()
    }

    fn insert_vqueue_buffer_len(&self) -> u32 {
        0
    }

    fn delete_vqueue_buffer_len(&self) -> u32 {
        0
    }

    fn is_indexing(&self) -> bool {
        false
    }

    fn is_saving(&self) -> bool {
        false
    }
}

fn decode_any<T>(options: &[Any]) -> Result<Option<T>, Error>
where
    T: Message + Default + prost::Name,
{
    let expected = format!("/{}", T::full_name());
    for option in options {
        if option.type_url.ends_with(&expected) {
            return T::decode(option.value.as_slice())
                .map(Some)
                .map_err(|err| Error::Backend {
                    message: err.to_string(),
                });
        }
    }
    Ok(None)
}

fn backend_error(error: zvec_rust::Error) -> Error {
    Error::Backend {
        message: error.to_string(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use prost::Name;
    use std::{
        fs,
        path::PathBuf,
        sync::atomic::{AtomicU64, Ordering},
    };

    static NEXT_TEST_ID: AtomicU64 = AtomicU64::new(0);

    struct TestDir(PathBuf);

    impl TestDir {
        fn new() -> Self {
            let id = NEXT_TEST_ID.fetch_add(1, Ordering::Relaxed);
            Self(std::env::temp_dir().join(format!("vald-zvec-{}-{id}", std::process::id())))
        }
    }

    impl Drop for TestDir {
        fn drop(&mut self) {
            if let Err(err) = fs::remove_dir_all(&self.0) {
                if err.kind() != std::io::ErrorKind::NotFound {
                    eprintln!("failed to remove Zvec test directory: {err}");
                }
            }
        }
    }

    fn document_options(fields: &[(&str, &str)]) -> Any {
        let value = z_vec::DocumentOptions {
            fields: fields
                .iter()
                .map(|(name, value)| ((*name).to_string(), (*value).to_string()))
                .collect(),
        };
        Any {
            type_url: format!(
                "type.googleapis.com/{}",
                z_vec::DocumentOptions::full_name()
            ),
            value: value.encode_to_vec(),
        }
    }

    fn search_options(value: z_vec::SearchOptions) -> Any {
        Any {
            type_url: format!("type.googleapis.com/{}", z_vec::SearchOptions::full_name()),
            value: value.encode_to_vec(),
        }
    }

    #[test]
    fn decode_any_options() {
        let tests = [
            (
                "matching type",
                vec![document_options(&[("category", "news")])],
                Some("news"),
                false,
            ),
            (
                "unrelated type",
                vec![Any {
                    type_url: "type.googleapis.com/example.Unrelated".to_string(),
                    value: Vec::new(),
                }],
                None,
                false,
            ),
            (
                "malformed payload",
                vec![Any {
                    type_url: format!(
                        "type.googleapis.com/{}",
                        z_vec::DocumentOptions::full_name()
                    ),
                    value: vec![0xff],
                }],
                None,
                true,
            ),
        ];

        for (name, options, expected, want_error) in tests {
            let result = decode_any::<z_vec::DocumentOptions>(&options);
            assert_eq!(result.is_err(), want_error, "{name}");
            if let Ok(value) = result {
                assert_eq!(
                    value.and_then(|value| value.fields.get("category").cloned()),
                    expected.map(str::to_string),
                    "{name}"
                );
            }
        }
    }

    #[test]
    fn collection_crud_and_search() {
        let directory = TestDir::new();
        {
            let mut service = Service::new(
                directory.0.to_str().expect("test path must be UTF-8"),
                3,
                vec!["category".to_string()],
                vec!["content".to_string()],
            )
            .expect("service creation must succeed");
            let options = [document_options(&[
                ("category", "news"),
                ("content", "vector search news"),
            ])];

            service
                .insert_with_options("first".to_string(), vec![0.0, 0.0, 0.0], 42, &options)
                .expect("insert must succeed");
            service
                .insert("second".to_string(), vec![1.0, 1.0, 1.0], 43)
                .expect("second insert must succeed");

            assert!(service.exists("first".to_string()));
            assert_eq!(service.len(), 2);
            assert_eq!(
                service
                    .get_object("first".to_string())
                    .expect("fetch must succeed"),
                (vec![0.0, 0.0, 0.0], 42)
            );

            let response = service
                .search(vec![0.0, 0.0, 0.0], 2, 0.0, -1.0)
                .expect("search must succeed");
            assert_eq!(response.results.len(), 2);
            assert_eq!(response.results[0].id, "first");
            assert!(response.results[0].distance <= response.results[1].distance);

            let filtered = service
                .search_with_options(
                    vec![0.0, 0.0, 0.0],
                    2,
                    0.0,
                    -1.0,
                    &[search_options(z_vec::SearchOptions {
                        pre_filter: "category = 'news'".to_string(),
                        ..Default::default()
                    })],
                )
                .expect("filtered search must succeed");
            assert_eq!(filtered.results.len(), 1);
            assert_eq!(filtered.results[0].id, "first");

            let full_text = service
                .search_with_options(
                    vec![0.0, 0.0, 0.0],
                    2,
                    0.0,
                    -1.0,
                    &[search_options(z_vec::SearchOptions {
                        hybrid_queries: vec![z_vec::Query {
                            field_name: "content".to_string(),
                            query: Some(z_vec::query::Query::Fts(z_vec::Fts {
                                match_string: "news".to_string(),
                                query_string: String::new(),
                            })),
                        }],
                        ..Default::default()
                    })],
                )
                .expect("full-text search must succeed");
            assert_eq!(full_text.results.len(), 1);
            assert_eq!(full_text.results[0].id, "first");

            assert!(matches!(
                service.insert("first".to_string(), vec![0.0, 0.0, 0.0], 44),
                Err(Error::UUIDAlreadyExists { uuid }) if uuid == "first"
            ));
            assert!(matches!(
                service.insert("bad-dimension".to_string(), vec![0.0, 0.0], 0),
                Err(Error::InvalidDimensionSize { .. })
            ));

            service
                .remove("first".to_string(), 0)
                .expect("remove must succeed");
            assert!(!service.exists("first".to_string()));
        }
    }
}
