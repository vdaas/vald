use std::error::Error;

// KVS trait for abstraction
pub trait KVS: Send + Sync {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> where Self: Sized + Send + Sync;
    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>, Box<dyn Error + '_>>;
    fn set(&self, key: &[u8], value: &[u8]) -> Result<(), Box<dyn Error + '_>>;
    fn del(&self, key: &[u8]) -> Result<(), Box<dyn Error + '_>>;
}
