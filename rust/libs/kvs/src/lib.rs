use std::error::Error;

// KVS trait for abstraction
pub trait KVS {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> where Self: Sized;
    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>, Box<dyn Error + '_>>;
    fn set(&self, key: &[u8], value: &[u8]) -> Result<(), Box<dyn Error + '_>>;
}
