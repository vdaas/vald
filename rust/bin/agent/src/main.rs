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

use algorithm::Error;
use anyhow::Result;
use proto::payload::v1::object::Distance;
use proto::payload::v1::search;
use qbg::index::Index;
use qbg::property::Property;

mod handler;

#[derive(Debug)]
struct _MockService {
    dim: usize,
}

struct QBGService {
    path: String,
    index: Index,
    property: Property,
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

    fn update(&mut self, _uuid: String, _vector: Vec<f32>, _ts: i64) -> Result<(), Error> {
        todo!()
    }

    fn remove(&mut self, _uuid: String, _ts: i64) -> Result<(), Error> {
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

    fn update(&mut self, uuid: String, vector: Vec<f32>, ts: i64) -> Result<(), Error> {
        self.remove(uuid.clone(), ts).unwrap();
        self.insert(uuid, vector, ts).unwrap();
        Ok(())
    }

    fn remove(&mut self, _uuid: String, _ts: i64) -> Result<(), Error> {
        // convert uuid to id
        let id = 1;
        self.index.remove(id).unwrap();
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
        self.index.get_dimension().unwrap()
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
    let path: String = "index".to_string();
    let mut p = Property::new();
    p.init_qbg_construction_parameters();
    p.set_dimension(128);
    p.set_number_of_subvectors(64);
    p.set_number_of_blobs(0);
    p.init_qbg_build_parameters();
    p.set_number_of_objects(500);
    let index = Index::new(&path, &mut p).unwrap();
    let service = QBGService {
        path: path,
        index: index,
        property: p,
    };
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
