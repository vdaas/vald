use criterion::{black_box, criterion_group, criterion_main, measurement::WallTime, BenchmarkGroup, Criterion};
use kvs::KVS;

mod kvs_impl;
use kvs_impl::*;

mod util;
use util::*;

const SIZE: usize = 1 << 10;
const DIM: usize = 1 << 7;

fn bench_get<T: KVS + Send + Sync + 'static>(group: &mut BenchmarkGroup<WallTime>, name: &str, keys: &Vec<Vec<u8>>) {
    let db = setup_kvs::<T>(SIZE, DIM);

    group.bench_function(&format!("{}", name), |b| {
        b.iter(|| {
            for key in keys {
                let _ = black_box(db.get(black_box(&key[..])));
            }
        })
    });
}

fn bench_set<T: KVS + Send + Sync + 'static>(group: &mut BenchmarkGroup<WallTime>, name: &str, keys: &Vec<Vec<u8>>, value: &[u8]) {
    let db = setup_kvs::<T>(SIZE, DIM);

    group.bench_function(&format!("{}", name), |b| {
        b.iter(|| {
            for key in keys.iter() {
                black_box(db.set(black_box(&key[..]), black_box(value))).unwrap();
            }
        })
    });
}

fn criterion_benchmark(c: &mut Criterion) {
    let keys = sequential_keys(SIZE, 0);
    let value = random_bytes(DIM);
    {
        let mut group = c.benchmark_group("set");
        bench_set::<Sled>(&mut group, "sled", &keys, &value);
        bench_set::<Kv>(&mut group, "kv", &keys, &value);
        bench_set::<Rkv>(&mut group, "rkv", &keys, &value);
        bench_set::<Redb>(&mut group, "redb", &keys, &value);
        bench_set::<Rocksdb>(&mut group, "rocksdb", &keys, &value);
        bench_set::<Persy>(&mut group, "persy", &keys, &value);
        group.finish();
    }
    {
        let mut group = c.benchmark_group("get");
        bench_get::<Sled>(&mut group, "sled", &keys);
        bench_get::<Kv>(&mut group, "kv", &keys);
        bench_get::<Rkv>(&mut group, "rkv", &keys);
        bench_get::<Redb>(&mut group, "redb", &keys);
        bench_get::<Rocksdb>(&mut group, "rocksdb", &keys);
        bench_get::<Persy>(&mut group, "persy", &keys);
        group.finish();
    }
    {
        let unused_keys = (SIZE..SIZE+SIZE).map(|i| -> Vec<u8> {i.to_ne_bytes().to_vec()}).collect();
        let mut group = c.benchmark_group("unused get");
        bench_get::<Sled>(&mut group, "sled", &unused_keys);
        bench_get::<Kv>(&mut group, "kv", &unused_keys);
        bench_get::<Rkv>(&mut group, "rkv", &unused_keys);
        bench_get::<Redb>(&mut group, "redb", &unused_keys);
        bench_get::<Rocksdb>(&mut group, "rocksdb", &unused_keys);
        bench_get::<Persy>(&mut group, "persy", &unused_keys);
        group.finish();
    }
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);
