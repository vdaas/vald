use anyhow::Result;

pub mod algorithm {
    pub trait Base<T, U, Param: Deserialize, Response: Serialize> {
        fn search(&self, v: &Vec<T>, p: Option<&Param>) -> Result<&Response>;
        fn insert(&self, id: &U, v: &Vec<T>, p: Option<&Param>) -> Result<&Response>;
        fn update(&self, id: &U, v: &Vec<T>, p: Option<&Param>) -> Result<&Response>;
        fn remove(&self, id: &U, p: Option<&Param>) -> Result<&Response>;
        fn commit(&self, p: Option<&Param>) -> Result<&Response>;

        fn new(p: Option<&Param>) -> Result<&Self>;
        fn open(p: &str) -> Result<&Self>;
        fn save(p: &str) -> Result<&Self>;
    }
}
