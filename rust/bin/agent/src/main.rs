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

use algorithm::{Error, MultiError};
use anyhow::Result;
use chrono::{Local, Timelike};
use config::Config;
use proto::payload::v1::object::Distance;
use proto::payload::v1::search;
use qbg::index::Index;
use qbg::property::Property;
use std::collections::HashMap;

mod handler;

#[derive(Debug)]
struct _MockService {
    dim: usize,
}

impl algorithm::ANN for _MockService {
    fn exists(&self, _uuid: String) -> bool {
        todo!()
    }

    fn create_index(&mut self) -> Result<(), Error> {
        todo!()
    }

    fn save_index(&mut self) -> Result<(), Error> {
        todo!()
    }

    fn insert(&mut self, _uuid: String, _vector: Vec<f32>, _ts: i64) -> Result<(), Error> {
        todo!()
    }

    fn insert_multiple(&mut self, _vectors: HashMap<String, Vec<f32>>) -> Result<(), Error> {
        todo!()
    }

    fn update(&mut self, _uuid: String, _vector: Vec<f32>, _ts: i64) -> Result<(), Error> {
        todo!()
    }

    fn update_multiple(&mut self, _vectors: HashMap<String, Vec<f32>>) -> Result<(), Error> {
        todo!()
    }

    fn ready_for_update(
        &mut self,
        _uuid: String,
        _vector: Vec<f32>,
        _ts: i64,
    ) -> Result<(), Error> {
        todo!()
    }

    fn remove(&mut self, _uuid: String, _ts: i64) -> Result<(), Error> {
        todo!()
    }

    fn remove_multiple(&mut self, _uuids: Vec<String>) -> Result<(), Error> {
        todo!()
    }

    fn search(
        &self,
        vector: Vec<f32>,
        _k: u32,
        _epsilon: f32,
        _radius: f32,
    ) -> Result<search::Response, Error> {
        Err(Error::IncompatibleDimensionSize {
            got: vector.len() as usize,
            want: self.dim,
        }
        .into())
    }

    fn get_object(&self, _uuid: String) -> Result<(Vec<f32>, i64), Error> {
        todo!()
    }

    fn get_dimension_size(&self) -> usize {
        self.dim
    }

    fn len(&self) -> u32 {
        todo!()
    }

    fn insert_vqueue_buffer_len(&self) -> u32 {
        todo!()
    }

    fn delete_vqueue_buffer_len(&self) -> u32 {
        todo!()
    }

    fn is_indexing(&self) -> bool {
        todo!()
    }

    fn is_saving(&self) -> bool {
        todo!()
    }
}

struct QBGService {
    path: String,
    index: Index,
    property: Property,
}

impl QBGService {
    fn new(settings: Config) -> Self {
        let path = settings
            .get::<String>("qbg.index_path")
            .unwrap_or("index".to_string());
        let mut property = Property::new();
        property.init_qbg_construction_parameters();
        property.set_qbg_construction_parameters(
            settings.get::<usize>("qbg.extended_dimension").unwrap_or(0),
            settings.get::<usize>("qbg.dimension").unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_subvectors")
                .unwrap_or(1),
            settings.get::<usize>("qbg.number_of_blobs").unwrap_or(0),
            settings.get::<i32>("qbg.internal_data_type").unwrap_or(1),
            settings.get::<i32>("qbg.data_type").unwrap_or(1),
            settings.get::<i32>("qbg.distance_type").unwrap_or(1),
        );
        property.init_qbg_build_parameters();
        property.set_qbg_build_parameters(
            settings
                .get::<i32>("qbg.hierarchical_clustering_init_mode")
                .unwrap_or(2),
            settings
                .get::<usize>("qbg.number_of_first_objects")
                .unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_first_clusters")
                .unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_second_objects")
                .unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_second_clusters")
                .unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_third_clusters")
                .unwrap_or(0),
            settings
                .get::<usize>("qbg.number_of_objects")
                .unwrap_or(1000),
            settings
                .get::<usize>("qbg.number_of_subvectors")
                .unwrap_or(1),
            settings
                .get::<i32>("qbg.optimization_clustering_init_mode")
                .unwrap_or(2),
            settings
                .get::<usize>("qbg.rotation_iteration")
                .unwrap_or(2000),
            settings
                .get::<usize>("qbg.subvector_iteration")
                .unwrap_or(400),
            settings.get::<usize>("qbg.number_of_matrices").unwrap_or(3),
            settings.get::<bool>("qbg.rotation").unwrap_or(true),
            settings.get::<bool>("qbg.repositioning").unwrap_or(false),
        );
        let index = Index::new(&path, &mut property).unwrap();
        QBGService {
            path,
            index,
            property,
        }
    }
}

impl algorithm::ANN for QBGService {
    fn exists(&self, _uuid: String) -> bool {
        // convert uuid to id
        let id = 1;
        let result = self.index.get_object(id);
        match result {
            Ok(_vec) => true,
            Err(_err) => false,
        }
    }

    fn create_index(&mut self) -> Result<(), Error> {
        self.index
            .build_index(&self.path, &mut self.property)
            .unwrap();
        Ok(())
    }

    fn save_index(&mut self) -> Result<(), Error> {
        self.index.save_index().unwrap();
        Ok(())
    }

    fn insert(&mut self, _uuid: String, vector: Vec<f32>, _ts: i64) -> Result<(), Error> {
        let _i = self.index.append(vector.as_slice()).unwrap();
        Ok(())
    }

    fn insert_multiple(&mut self, vectors: HashMap<String, Vec<f32>>) -> Result<(), Error> {
        let mut uuids: Vec<String> = vec![];
        for (uuid, vec) in vectors {
            let result = self.insert(uuid, vec, Local::now().nanosecond().into());
            match result {
                Ok(()) => continue,
                Err(err) => match err {
                    Error::UUIDAlreadyExists { uuid } => uuids.push(uuid),
                    _ => return Err(err),
                },
            }
        }
        if !uuids.is_empty() {
            return Err(Error::new_uuid_already_exists(uuids));
        }
        Ok(())
    }

    fn update(&mut self, uuid: String, vector: Vec<f32>, ts: i64) -> Result<(), Error> {
        self.remove(uuid.clone(), ts)?;
        self.insert(uuid, vector, ts)?;
        Ok(())
    }

    fn update_multiple(&mut self, mut vectors: HashMap<String, Vec<f32>>) -> Result<(), Error> {
        let mut uuids: Vec<String> = vec![];
        for (uuid, vec) in vectors.clone() {
            let result = self.ready_for_update(uuid.clone(), vec, Local::now().nanosecond().into());
            match result {
                Ok(()) => uuids.push(uuid),
                Err(_err) => {
                    let _ = vectors.remove(&uuid);
                }
            }
        }
        self.remove_multiple(uuids.clone())?;
        self.insert_multiple(vectors)
    }

    fn ready_for_update(&mut self, uuid: String, vector: Vec<f32>, ts: i64) -> Result<(), Error> {
        if uuid.len() == 0 {
            return Err(Error::UUIDNotFound {
                uuid: "0".to_string(),
            });
        }
        if vector.len() != self.get_dimension_size() {
            return Err(Error::InvalidDimensionSize {
                uuid: uuid,
                current: vector.len().to_string(),
                limit: self.get_dimension_size().to_string(),
            });
        }
        let (ovec, ots) = self.get_object(uuid.clone())?;
        if (vector.len() != ovec.len()) || (vector != ovec) {
            return Ok(());
        }
        if ots < ts {
            self.update(uuid.clone(), vector, ts)?;
            return Ok(());
        }
        Err(Error::UUIDAlreadyExists { uuid })
    }

    fn remove(&mut self, _uuid: String, _ts: i64) -> Result<(), Error> {
        // convert uuid to id
        let id = 1;
        self.index.remove(id).unwrap();
        Ok(())
    }

    fn remove_multiple(&mut self, uuids: Vec<String>) -> Result<(), Error> {
        let mut ids: Vec<String> = vec![];
        for uuid in uuids {
            let result = self.remove(uuid, Local::now().nanosecond().into());
            match result {
                Ok(()) => continue,
                Err(err) => match err {
                    Error::ObjectIDNotFound { uuid } => ids.push(uuid),
                    _ => return Err(err),
                },
            }
        }
        if !ids.is_empty() {
            return Err(Error::new_object_id_not_found(ids));
        }
        Ok(())
    }

    fn search(
        &self,
        vector: Vec<f32>,
        k: u32,
        epsilon: f32,
        radius: f32,
    ) -> Result<search::Response, Error> {
        let vec = self
            .index
            .search(vector.as_slice(), k as usize, radius, epsilon)
            .unwrap();
        let results: Vec<Distance> = vec
            .into_iter()
            .map(|x| Distance {
                id: x.0.to_string(),
                distance: x.1,
            })
            .collect();
        let res = search::Response {
            request_id: "".to_string(),
            results: results,
        };
        Ok(res)
    }

    fn get_object(&self, _uuid: String) -> Result<(Vec<f32>, i64), Error> {
        // convert uuid to id
        let id = 1;
        let vec = self.index.get_object(id).unwrap();
        // get timestamp
        let ts: i64 = 0;
        Ok((vec.to_vec(), ts))
    }

    fn get_dimension_size(&self) -> usize {
        self.index.get_dimension().unwrap_or_default()
    }

    fn len(&self) -> u32 {
        todo!()
    }

    fn insert_vqueue_buffer_len(&self) -> u32 {
        todo!()
    }

    fn delete_vqueue_buffer_len(&self) -> u32 {
        todo!()
    }

    fn is_indexing(&self) -> bool {
        todo!()
    }

    fn is_saving(&self) -> bool {
        todo!()
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "0.0.0.0:8081".parse()?;
    let settings = Config::builder()
        .add_source(config::File::with_name("/etc/server/config.yaml"))
        .build()
        .unwrap();
    let _logger =
        flexi_logger::Logger::try_with_str(settings.get::<String>("logging.level")?)?.start()?;
    let service = QBGService::new(settings);
    let agent = handler::Agent::new(
        service,
        "agent-qbg",
        "127.0.0.1",
        "vald/internal/core/algorithm",
        "vald-agent",
    );

    tonic::transport::Server::builder()
        .add_service(proto::core::v1::agent_server::AgentServer::new(agent))
        .serve(addr)
        .await?;

    Ok(())
}
