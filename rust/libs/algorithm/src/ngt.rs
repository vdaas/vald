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
    pub struct CommitParam {}

    #[derive(Debug)]
    pub struct NGTParam {}

    impl
        Base<
            f32,
            String,
            SearchParam,
            SearchResponse,
            InsertParam,
            (),
            UpdateParam,
            (),
            RemoveParam,
            (),
            CommitParam,
            (),
            NGTParam,
        > for NGT
    {
        fn search(&self, v: &Vec<T>, p: Option<SearchParam>) -> Result<SearchResponse>;
        fn insert(&self, v: &Vec<T>, id: &U, p: Option<InsertParam>) -> Result<()>;
        fn update(&self, v: &Vec<T>, id: &U, p: Option<UpdateParam>) -> Result<()>;
        fn remove(&self, id: &U, p: Option<RemoveParam>) -> Result<()>;
        fn commit(&self, p: Option<CommitParam>) -> Result<()>;

        fn new(p: Option<NGTParam>) -> Result<Self>;
    }
}
