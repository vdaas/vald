pub mod algorithm1 {
    use anyhow::Result;

    pub trait Base<T, U, SP, SR, IP, IR, UP, UR, RP, RR, CP, CR, P>: Send + Sync
    where
        T: num::Num,
        U: num::sign::Unsigned,
    {
        fn search(&self, v: &[T], p: Option<SP>) -> Result<SR>;
        fn insert(&self, v: &[T], id: &U, p: Option<IP>) -> Result<IR>;
        fn update(&self, v: &[T], id: &U, p: Option<UP>) -> Result<UR>;
        fn remove(&self, id: &U, p: Option<RP>) -> Result<RR>;
        fn commit(&self, p: Option<CP>) -> Result<CR>;
    }
}

pub mod algorithm2 {
    use anyhow::Result;

    pub trait Search<T: num::Num, P, R>: Send + Sync {
        fn search(&self, v: &[T], p: Option<P>) -> Result<R>;
    }

    pub trait Insert<T: num::Num, U: num::sign::Unsigned, P, R>: Send + Sync {
        fn insert(&self, v: &[T], id: &U, p: Option<P>) -> Result<R>;
    }

    pub trait Update<T: num::Num, U: num::sign::Unsigned, P, R>: Send + Sync {
        fn update(&self, v: &[T], id: &U, p: Option<P>) -> Result<R>;
    }

    pub trait Remove<U: num::sign::Unsigned, P, R>: Send + Sync {
        fn remove(&self, id: &U, p: Option<P>) -> Result<R>;
    }

    pub trait Commit<P, R>: Send + Sync {
        fn commit(&self, p: Option<P>) -> Result<R>;
    }
}

// inspired by https://github.com/jonluca/hora
pub mod algorithm3 {
    use anyhow::Result;

    pub trait Base<T, U, SP, SR, IP, UP, RP, CP>: Send + Sync
    where
        T: num::Num,
        U: num::sign::Unsigned,
    {
        fn search(&self, v: &[T], p: Option<SP>) -> Result<Vec<SR>>;
        fn insert(&mut self, v: &[T], id: &U, p: Option<IP>) -> Result<()>;
        fn update(&mut self, v: &[T], id: &U, p: Option<UP>) -> Result<()>;
        fn remove(&mut self, id: &U, p: Option<RP>) -> Result<()>;
        fn commit(&mut self, p: Option<CP>) -> Result<()>;
    }
}
