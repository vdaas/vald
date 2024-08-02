use std::{any::type_name, path::{Path, PathBuf}, sync::OnceLock};

use kvs::KVS;
use rand::{thread_rng, Rng};

static TIMESTAMP: OnceLock<i64> = OnceLock::new();

pub fn setup_kvs<T: KVS + Send + Sync + 'static>(param: &str, size: Option<usize>, kdim: Option<usize>, vdim: Option<usize>) -> (PathBuf, T) {
    let name = type_name::<T>().split("::").last().unwrap();
    let path = Path::new(&format!("/tmp/kvs_bench/{}", TIMESTAMP.get_or_init(|| chrono::Local::now().timestamp()))).join(format!("{}-{}", name, param));
    let db = T::new(path.to_str().unwrap()).unwrap();
    if let Some(size) = size {
        let kdim = match kdim {
            Some(d) => d,
            None => 1024,
        };
        let vdim = match vdim {
            Some(d) => d,
            None => 64,
        };
        for i in 0..size {
            let mut key = i.to_ne_bytes().to_vec();
            key.resize_with(kdim, Default::default);
            db.set(&key, &random_bytes(vdim)).unwrap();
        }
    }
    (path.clone(), db)
}

pub fn random_bytes(dim: usize) -> Vec<u8> {
    let mut buf = vec![0u8; dim];
    thread_rng().fill(&mut buf[..]);
    buf
}

