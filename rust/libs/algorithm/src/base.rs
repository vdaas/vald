pub mod algorithm {
    use anyhow::Result;

    pub trait Base<T, U, SP, SR, IP, IR, UP, UR, RP, RR, CP, CR, P> {
        fn search(&self, v: &Vec<T>, p: Option<SP>) -> Result<SR>;
        fn insert(&self, v: &Vec<T>, id: &U, p: Option<IP>) -> Result<IR>;
        fn update(&self, v: &Vec<T>, id: &U, p: Option<UP>) -> Result<UR>;
        fn remove(&self, id: &U, p: Option<RP>) -> Result<RR>;
        fn commit(&self, p: Option<CP>) -> Result<CR>;

        fn new(p: Option<P>) -> Result<Self>;
    }
}

pub mod algorithm2 {
    use anyhow::Result;

    pub trait Search<T, P, R> {
        fn search(&self, v: &[T], p: Option<P>) -> Result<R>;
    }

    pub trait Insert<T, U, P, R> {
        fn insert(&self, v: &[T], id: &U, p: Option<P>) -> Result<R>;
    }

    pub trait Update<T, U, P, R> {
        fn update(&self, v: &[T], id: &U, p: Option<P>) -> Result<R>;
    }

    pub trait Remove<U, P, R> {
        fn remove(&self, id: &U, p: Option<P>) -> Result<R>;
    }

    pub trait Commit<P, R> {
        fn commit(&self, p: Option<P>) -> Result<R>;
    }
}
