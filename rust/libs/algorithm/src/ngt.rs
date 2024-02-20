pub mod ngt {
    use crate::algorithm::Base;
    use anyhow::{anyhow, Ok, Result};

    #[derive(Debug)]
    struct Index {}

    #[derive(Debug)]
    pub struct NGT {
        index: Index,
    }
    #[derive(Debug)]
    pub struct SearchResult {
        id: String,
        distance: f32,
    }

    #[derive(Debug)]
    pub struct SearchResponse {
        results: Vec<SearchResult>,
    }

    #[derive(Debug)]
    pub struct SearchParam {
        epsilon: f32,
        radius: f32,
        num: u32,
    }

    #[derive(Debug)]    
    pub struct InsertParam {}
    
    #[derive(Debug)]
    pub struct UpdateParam {}

    #[derive(Debug)]
    pub struct RemoveParam {}

    #[derive(Debug)]
    pub struct CommitParam{}

    #[derive(Debug)]
    pub struct NGTParam {}

    impl Base<f32, String, SearchParam, SearchResponse, InsertParam, (), UpdateParam, (), RemoveParam, (), CommitParam, ()> for NGT {
        fn search()
    }
}