use criterion::{black_box, criterion_group, criterion_main, measurement::WallTime, BenchmarkGroup, Criterion};
use std::error::Error;
use kvs::KVS;
use rand::{thread_rng, Rng};

// Implement KVS for sled
struct Sled(sled::Db);
impl KVS for Sled {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> {
        Ok(Sled(sled::open(path)?))
    }
    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>, Box<dyn Error>> {
        Ok(self.0.get(key)?.map(|iv| iv.to_vec()))
    }
    fn set(&self, key: &[u8], value: &[u8]) -> Result<(), Box<dyn Error>> {
        self.0.insert(key, value)?;
        self.0.flush()?;
        Ok(())
    }
}

// Implement KVS for kv
struct Kv(kv::Store);
impl KVS for Kv {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> {
        let cfg = kv::Config::new(path);
        let store = kv::Store::new(cfg)?;
        Ok(Kv(store))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>, Box<dyn Error>> {
        let bucket = self.0.bucket::<kv::Raw, kv::Raw>(None)?;
        Ok(bucket.get(&kv::Raw::from(key))?.map(|v| v.to_vec()))
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<(), Box<dyn Error>> {
        let bucket = self.0.bucket::<kv::Raw, kv::Raw>(None)?;
        bucket.set(&kv::Raw::from(key), &kv::Raw::from(value))?;
        Ok(())
    }
}

// Implement KVS for rkv
struct Rkv(std::sync::Arc<std::sync::RwLock<rkv::Rkv<rkv::backend::SafeModeEnvironment>>>, rkv::SingleStore<rkv::backend::SafeModeDatabase>);
impl KVS for Rkv {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> {
        std::fs::create_dir_all(path)?;
        let mut manager = rkv::Manager::<rkv::backend::SafeModeEnvironment>::singleton().write()?;
        let created_arc = manager.get_or_create(std::path::Path::new(path), rkv::Rkv::new::<rkv::backend::SafeMode>)?;
        let store = created_arc.read().unwrap().open_single("mydb", rkv::StoreOptions::create())?;
        Ok(Rkv(created_arc, store))
    }
    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>, Box<dyn Error + '_>> {
        let env = self.0.read()?;
        let reader = env.read()?;
        Ok(Some(self.1.get(&reader, key)?.unwrap().to_bytes()?.to_vec()))
    }
    fn set(&self, key: &[u8], value: &[u8]) -> Result<(), Box<dyn Error + '_>> {
        let env = self.0.read()?;
        let mut writer = env.write()?;
        self.1.put(&mut writer, key, &rkv::Value::Blob(value))?;
        writer.commit()?;
        Ok(())
    }
}

// Implement KVS for redb
struct Redb(redb::Database, redb::TableDefinition<'static, &'static [u8], &'static [u8]>);
impl KVS for Redb {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> {
        let db = redb::Database::create(path)?;
        let def = redb::TableDefinition::new("my_table");
        Ok(Redb(db, def))
    }
    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>, Box<dyn Error>> {
        let txn = self.0.begin_read()?;
        let table = txn.open_table(self.1)?;
        Ok(Some(table.get(key)?.unwrap().value().to_vec()))
    }
    fn set(&self, key: &[u8], value: &[u8]) -> Result<(), Box<dyn Error>> {
        let txn = self.0.begin_write()?;
        {
            let mut table = txn.open_table(self.1)?;
            table.insert(key, value)?;
        }
        txn.commit()?;
        Ok(())
    }
}

// Implement KVS for pickledb
struct Pickledb(std::cell::RefCell<pickledb::PickleDb>);
impl KVS for Pickledb {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> {
        let db = pickledb::PickleDb::new(path, pickledb::PickleDbDumpPolicy::AutoDump, pickledb::SerializationMethod::Bin);
        Ok(Pickledb(std::cell::RefCell::new(db)))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>, Box<dyn Error>> {
        Ok(self.0.borrow().get::<Vec<u8>>(std::str::from_utf8(key)?))
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<(), Box<dyn Error>> {
        self.0.borrow_mut().set(std::str::from_utf8(key)?, &value.to_vec())?;
        Ok(())
    }
}

// Implement KVS for rocksdb
struct Rocksdb(rocksdb::DB);
impl KVS for Rocksdb {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> {
        let db = rocksdb::DB::open_default(path)?;
        Ok(Rocksdb(db))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>, Box<dyn Error>> {
        Ok(self.0.get(key)?)
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<(), Box<dyn Error>> {
        self.0.put(key, value)?;
        Ok(())
    }
}

// Implement KVS for persy
struct Persy(persy::Persy);
impl KVS for Persy {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> {
        let persy = persy::Persy::open_or_create_with(path, persy::Config::new(), |persy| {
            let mut tx = persy.begin()?;
            tx.create_segment("main")?;
            tx.create_index::<persy::ByteVec, persy::ByteVec>("index", persy::ValueMode::Exclusive)?;
            tx.prepare()?.commit()?;
            Ok(())
        })?;
        Ok(Persy(persy))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>, Box<dyn Error>> {
        let mut tx = self.0.begin()?;
        let result = tx.one::<persy::ByteVec, persy::ByteVec>("index", &persy::ByteVec::new(key.to_vec()))?.map(|v| v.to_vec());
        Ok(result)
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<(), Box<dyn Error>> {
        let mut tx = self.0.begin()?;
        tx.put::<persy::ByteVec, persy::ByteVec>("index", persy::ByteVec::new(key.to_vec()), persy::ByteVec::new(value.to_vec()))?;
        tx.prepare()?.commit()?;
        Ok(())
    }
}

fn gen_random_bytes(size: usize, dim: usize) -> Vec<Vec<u8>> {
    let mut rng = thread_rng();
    (0..size).map(|_| {
        let mut buf = vec![0u8; dim];
        rng.fill(&mut buf[..]);
        buf
    }).collect()
}

fn setup_dataset(size: usize, key_dim: usize, val_dim: usize) -> (Vec<Vec<u8>>, Vec<Vec<u8>>) {
    (gen_random_bytes(size, key_dim), gen_random_bytes(size, val_dim))
}

fn bench_get<T: KVS>(group: &mut BenchmarkGroup<WallTime>, name: &str) {
    let db = T::new(&format!("/tmp/bench_db/{}", name)).unwrap();
    let size = 1000;
    let (keys, values) = setup_dataset(size, 64, 256);
    for i in 0..size {
        db.set(&keys[i][..], &values[i][..]).unwrap();
    }

    group.bench_function(&format!("{}", name), |b| {
        b.iter(|| {
            for i in 0..size {
                black_box(db.get(black_box(&keys[i][..])).unwrap());
            }
        })
    });
}

fn bench_set<T: KVS>(group: &mut BenchmarkGroup<WallTime>, name: &str) {
    let db = T::new(&format!("/tmp/bench_db/{}", name)).unwrap();
    let size = 1000;
    let (keys, values) = setup_dataset(size, 64, 256);

    group.bench_function(&format!("{}", name), |b| {
        b.iter(|| {
            for i in 0..size {
                black_box(db.set(black_box(&keys[i][..]), black_box(&values[i][..])).unwrap());
            }
        })
    });
}

fn criterion_benchmark(c: &mut Criterion) {
    {
        let mut group = c.benchmark_group("get");
        bench_get::<Sled>(&mut group, "sled");
        bench_get::<Kv>(&mut group, "kv");
        bench_get::<Rkv>(&mut group, "rkv");
        bench_get::<Redb>(&mut group, "redb");
        //bench_get::<Pickledb>(&mut group, "pickledb"); // only use string key
        bench_get::<Rocksdb>(&mut group, "rocksdb");
        bench_get::<Persy>(&mut group, "persy");    
        group.finish();
    }
    {
        let mut group = c.benchmark_group("set");
        bench_set::<Sled>(&mut group, "sled");
        bench_set::<Kv>(&mut group, "kv");
        bench_set::<Rkv>(&mut group, "rkv");
        bench_set::<Redb>(&mut group, "redb");
        //bench_set::<Pickledb>(&mut group, "pickledb"); // only use string key
        bench_set::<Rocksdb>(&mut group, "rocksdb");
        bench_set::<Persy>(&mut group, "persy");
        group.finish();
    }
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);
