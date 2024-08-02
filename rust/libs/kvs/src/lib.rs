use anyhow::Result;

// KVS trait for abstraction
pub trait KVS: Send + Sync {
    fn new(path: &str) -> Result<Self> where Self: Sized + Send + Sync;
    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>>;
    fn range(&self, f: impl Fn(&[u8], &[u8]) -> Result<()>) -> Vec<Result<()>>;
    fn set(&self, key: &[u8], value: &[u8]) -> Result<()>;
    fn del(&self, key: &[u8]) -> Result<()>;
}
