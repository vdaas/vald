use std::{any::type_name, path::{Path, PathBuf}, sync::OnceLock};

use kvs::KVS;
use rand::{thread_rng, Rng};

static TIMESTAMP: OnceLock<i64> = OnceLock::new();

pub fn setup_kvs<T: KVS + Send + Sync + 'static>(param: &str) -> (PathBuf, T) {
    let name = type_name::<T>().split("::").last().unwrap();
    let path = Path::new(&format!("/tmp/kvs_bench/{}", TIMESTAMP.get_or_init(|| chrono::Local::now().timestamp()))).join(format!("{}-{}", name, param));
    (path.clone(), T::new(path.to_str().unwrap()).unwrap())
}

pub fn random_bytes(dim: usize) -> Vec<u8> {
    let mut buf = vec![0u8; dim];
    thread_rng().fill(&mut buf[..]);
    buf
}

