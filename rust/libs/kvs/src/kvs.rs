use anyhow::Result;

pub mod kvs {
    pub trait KVS<K, V> {
        fn set(k: &K, v: &V) -> Result<()>;
        fn get_value(k: &K) -> Result<&V>;
        fn get_key(v: &V) -> Result<&K>;

        fn new(p: &str) -> Result<&Self>;
        fn open(p: &str) -> Result<&Self>;
    }
}
