use divan::{black_box, AllocProfiler};
use kvs::KVS;
use std::sync::{Arc, OnceLock};
use std::thread;

mod kvs_impl;
use kvs_impl::{random_bytes, sequential_keys, setup_kvs, Kv, Persy, Redb, Rkv, Rocksdb, Sled};

const SIZE: usize = 1 << 10;
const DIM: usize = 1 << 7;

static KEYS: OnceLock<Vec<Vec<u8>>> = OnceLock::new();
static VALUE: OnceLock<Vec<u8>> = OnceLock::new();

#[global_allocator]
static ALLOC: AllocProfiler = AllocProfiler::system();

#[divan::bench(types = [
    Sled, Kv, Rkv, Redb, Rocksdb, Persy
])]
fn bench_set<T: KVS + 'static>(bencher: divan::Bencher) {
    let db = setup_kvs::<T>();
    bencher.bench(||{
        for key in KEYS.get().unwrap().iter() {
            black_box(db.set(key, VALUE.get().unwrap()).unwrap());
        }
    });
}

#[divan::bench(types = [
    Sled, Kv, Rkv, Redb, Rocksdb, Persy
])]
fn bench_get<T: KVS + 'static>(bencher: divan::Bencher) {
    let db = setup_kvs::<T>();
    for key in KEYS.get().unwrap().iter() {
        db.set(key, VALUE.get().unwrap()).unwrap();
    }

    bencher.bench(||{
        for key in KEYS.get().unwrap().iter() {
            black_box(db.get(key)).unwrap();
        }    
    });
}

#[divan::bench(types = [
    Sled, Kv, Rkv, Redb, Rocksdb, Persy
])]
fn bench_parallel_get_set<T: KVS + Send + Sync + 'static>(bencher: divan::Bencher) {
    let db = Arc::new(setup_kvs::<T>());
    for key in KEYS.get().unwrap().iter() {
        db.set(key, VALUE.get().unwrap()).unwrap();
    }

    bencher.bench(||{
        let mut handles = vec![];
        {
            let db = Arc::clone(&db);
            let handle = thread::spawn(move || {
                for key in KEYS.get().unwrap().iter() {
                    black_box(db.set(key, VALUE.get().unwrap()).unwrap())
                }
            });
            handles.push(handle);
        }
        {
            let db = Arc::clone(&db);
            let handle = thread::spawn(move || {
                for key in KEYS.get().unwrap().iter() {
                    black_box(db.get(key).unwrap());
                    ()
                }
            });
            handles.push(handle);
        }
        for handle in handles {
            handle.join().unwrap();
        }
    });
}

fn main() {
    KEYS.get_or_init(||{sequential_keys(SIZE, 0)});
    VALUE.get_or_init(||{random_bytes(DIM)});
    divan::main();
}
