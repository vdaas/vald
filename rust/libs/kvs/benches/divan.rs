use divan::{bench, black_box, Bencher};
use kvs::KVS;
use std::any::type_name;
use std::collections::HashMap;
use std::sync::{Arc, Mutex, OnceLock};
use std::thread;

mod kvs_impl;
use kvs_impl::*;

mod util;
use util::*;

const SIZES: &[usize] = &[1 << 14, 1 << 15];
const DIM: usize = 1 << 8;

//#[global_allocator]
//static ALLOC: AllocProfiler = AllocProfiler::system();

#[bench(
    types = [Kv, Persy, Redb, Rkv, Rocksdb, Sled],
    args = SIZES,
    sample_size = 1,
    sample_count = 1,
)]
fn bench_1_set<T>(bencher: Bencher, size: usize)
where
    T: KVS + 'static,
{
    let (_, db) = setup_kvs::<T>(size, DIM);

    bencher.bench_local(|| {
        for i in 0..size {
            black_box(db.set(&i.to_ne_bytes().to_vec(), &random_bytes(DIM)).unwrap());
        }    
    });
}

#[bench(
    types = [Kv, Persy, Redb, Rkv, Rocksdb, Sled],
    args = SIZES,
    sample_size = 1,
    sample_count = 1,
)]
fn bench_2_get<T>(bencher: Bencher, size: usize)
where
    T: KVS + 'static,
{
    let (_, db)  = setup_kvs::<T>(size, DIM);

    bencher.bench_local(|| {
        for i in 0..size {
            black_box(db.get(&i.to_ne_bytes().to_vec()).unwrap());
        }    
    });
}

#[bench(
    types = [Kv, Persy, Redb, Rkv, Rocksdb, Sled],
    consts = SIZES,
    sample_size = 1,
    sample_count = 1,
)]
#[ignore]
fn bench_parallel_get_set<T, const N: usize>() -> Arc<T>
where
    T: KVS + 'static,
{
    let (_, db)  = setup_kvs::<T>(N, DIM);
    let db = Arc::new(db);

    let mut handles = vec![];
    {
        let db = Arc::clone(&db);
        let handle = thread::spawn(move || {
            for i in N..N*2 {
                black_box(db.set(&i.to_ne_bytes().to_vec(), &random_bytes(DIM)).unwrap());
            }
        });
        handles.push(handle);
    }
    {
        let db = Arc::clone(&db);
        let handle = thread::spawn(move || {
            for i in 0..N {
                black_box(db.get(&i.to_ne_bytes().to_vec()).unwrap());
            }
        });
        handles.push(handle);
    }
    for handle in handles {
        handle.join().unwrap();
    }

    db
}

fn main() {
    divan::main();
}
