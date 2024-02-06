pub mod algorithm {
    pub trait Base<T, U, Param, Response> {
        fn search(&self, v: &Vec<T>, p: Option<&Param>) -> Result<&Response>;
        fn insert(&self, v: &Vec<T>, p: Option<&Param>) -> Result<()>;
        fn update(&self, id: &U, v: &Vec<T>, p: Option<&Param>) -> Result<()>;
        fn remove(&self, id: &U, p: Option<&Param>) -> Result<()>;
        fn commit(&self, p: Option<&Param>) -> Result<()>;

        fn new(p: Option<&Param>) -> Result<&Self>;
        fn open(p: &str) -> Result<&Self>;
        fn save(p: &str) -> Result<&Self>;
    }

    pub trait Param {
        type A;
        fn get_parameter() -> Result<&A>;
        fn set_parameter(a: &A) -> Result<()>;
    }
}
