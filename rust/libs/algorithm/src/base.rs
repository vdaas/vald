use anyhow::Result;
use serde::Deserialize;
use serde::Serialize;

pub mod algorithm1 {
    pub trait Base<T, U, Param: Deserialize, Response: Serialize> {
        fn search(&self, v: &Vec<T>, p: Option<&Param>) -> Result<&Response>;
        fn insert(&self, v: &Vec<T>, p: Option<&Param>) -> Result<&Response>;
        fn update(&self, id: &U, v: &Vec<T>, p: Option<&Param>) -> Result<&Response>;
        fn remove(&self, id: &U, p: Option<&Param>) -> Result<&Response>;
        fn commit(&self, p: Option<&Param>) -> Result<&Response>;

        fn new(p: Option<&Param>) -> Result<&Self>;
        fn open(p: &str) -> Result<&Self>;
        fn save(p: &str) -> Result<&Self>;
    }

    pub trait Param {
        type A;
        fn get_parameters() -> Result<&A>;
        fn set_parameters(a: &A) -> Result<()>;
    }

    pub trait Response {
        type Status;
        type Result;
        fn get_status() -> Status;
        fn get_result() -> Result;
    }
}

pub mod algorithm2 {
    pub trait Base<Query: Deserialize, Response: Serialize, Param: Deserialize> {
        fn search(&self, q: &Query) -> Result<&Response>;
        fn insert(&self, q: &Query) -> Result<&Response>;
        fn update(&self, q: &Query) -> Result<&Response>;
        fn remove(&self, q: &Query) -> Result<&Response>;
        fn commit(&self, q: &Query) -> Result<&Response>;

        fn new(p: Option<&Param>) -> Result<&Self>;
        fn open(p: &str) -> Result<&Self>;
        fn save(p: &str) -> Result<&Self>;
    }

    // example like NGT
    #[derive(Serialize, Deserialize, Debug)]
    struct SearchQuery {
        vector: Vec<f32>,
        size: u64,
        epsilon: f32,
        radius: f32,
    }

    #[derive(Serialize, Deserialize, Debug)]
    struct InsertQuery {
        vector: Vec<f32>,
    }
}
