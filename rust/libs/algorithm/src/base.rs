pub mod algorithm {
    use anyhow::Result;

    pub trait Base<T, U, SP, SR, IP, IR, UP, UR, RP, RR, CP, CR> {
        fn search(&self, v: &Vec<T>, p: Option<&SP>) -> Result<&SR>;
        fn insert(&self, v: &Vec<T>, id: &U, p: Option<&IP>) -> Result<&IR>;
        fn update(&self, v: &Vec<T>, id: &U, p: Option<&UP>) -> Result<&UR>;
        fn remove(&self, id: &U, p: Option<&RP>) -> Result<&RR>;
        fn commit(&self, p: Option<&CP>) -> Result<&CR>;
    }
}
