use anyhow::{Result, anyhow};
use std::{fs, path::Path};
use kvs::KVS;

// Implement KVS for sled
pub struct Sled(sled::Db);
impl KVS for Sled {
    fn new(path: &str) -> Result<Self> {
        Ok(Sled(sled::open(path)?))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>> {
        Ok(self.0.get(key)?.map(|iv| iv.to_vec()))
    }

    fn range(&self, f: impl Fn(&[u8], &[u8]) -> Result<()>) -> Vec<Result<()>> {
        self.0.iter().map(|result| {
            match result {
                Ok(r) => f(r.0.as_ref(), r.1.as_ref()),
                Err(e) => Err(anyhow!(e))
            }
        }).collect()
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<()> {
        self.0.insert(key, value)?;
        self.0.flush()?;
        Ok(())
    }

    fn del(&self, _key: &[u8]) -> Result<()> {
        todo!()
    }
}

// Implement KVS for kv
pub struct Kv(kv::Store);
impl KVS for Kv {
    fn new(path: &str) -> Result<Self> {
        let cfg = kv::Config::new(path);
        let store = kv::Store::new(cfg)?;
        Ok(Kv(store))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>> {
        let bucket = self.0.bucket::<kv::Raw, kv::Raw>(None)?;
        Ok(bucket.get(&kv::Raw::from(key))?.map(|v| v.to_vec()))
    }

    fn range(&self, f: impl Fn(&[u8], &[u8]) -> Result<()>) -> Vec<Result<()>> {
        let bucket: kv::Bucket<kv::Raw, kv::Raw> = self.0.bucket::<kv::Raw, kv::Raw>(None).unwrap();
        bucket.iter().map(|result| {
            match result {
                Ok(r) => f(r.key::<kv::Raw>().unwrap().as_ref(), r.value::<kv::Raw>().unwrap().as_ref()),
                Err(e) => Err(anyhow!(e)),
            }
        }).collect()
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<()> {
        let bucket = self.0.bucket::<kv::Raw, kv::Raw>(None)?;
        bucket.set(&kv::Raw::from(key), &kv::Raw::from(value))?;
        Ok(())
    }

    fn del(&self, _key: &[u8]) -> Result<()> {
        todo!()
    }
}

pub struct Kv2(kv::Bucket<'static, kv::Raw, kv::Raw>);
impl KVS for Kv2 {
    fn new(path: &str) -> Result<Self> {
        let cfg = kv::Config::new(path);
        let store = kv::Store::new(cfg)?;
        let bucket = store.bucket::<kv::Raw, kv::Raw>(Some("root"))?;
        Ok(Kv2(bucket))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>> {
        Ok(self.0.get(&kv::Raw::from(key))?.map(|v| v.to_vec()))
    }

    fn range(&self, f: impl Fn(&[u8], &[u8]) -> Result<()>) -> Vec<Result<()>> {
        self.0.iter().map(|result| {
            match result {
                Ok(r) => f(r.key::<kv::Raw>().unwrap().as_ref(), r.value::<kv::Raw>().unwrap().as_ref()),
                Err(e) => Err(anyhow!(e)),
            }
        }).collect()
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<()> {
        self.0.set(&kv::Raw::from(key), &kv::Raw::from(value))?;
        Ok(())
    }

    fn del(&self, _key: &[u8]) -> Result<()> {
        todo!()
    }
}

// Implement KVS for rkv
pub struct Rkv(rkv::Rkv<rkv::backend::SafeModeEnvironment>, rkv::SingleStore<rkv::backend::SafeModeDatabase>);
impl KVS for Rkv {
    fn new(path: &str) -> Result<Self> {
        fs::create_dir_all(path)?;
        let db = rkv::Rkv::new::<rkv::backend::SafeMode>(Path::new(path))?;
        let store = db.open_single("mydb", rkv::StoreOptions::create())?;
        Ok(Rkv(db, store))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>> {
        let reader = self.0.read()?;
        if let Some(val) = self.1.get(&reader, key)? {
            Ok(Some(val.to_bytes()?.to_vec()))
        } else {
            Ok(None)
        }
    }

    fn range(&self, f: impl Fn(&[u8], &[u8]) -> Result<()>) -> Vec<Result<()>> {
        todo!()
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<()> {
        let mut writer = self.0.write()?;
        self.1.put(&mut writer, key, &rkv::Value::Blob(value))?;
        writer.commit()?;
        Ok(())
    }

    fn del(&self, _key: &[u8]) -> Result<()> {
        todo!()
    }
}

// Implement KVS for redb
pub struct Redb(redb::Database, redb::TableDefinition<'static, &'static [u8], &'static [u8]>);
impl KVS for Redb {
    fn new(path: &str) -> Result<Self> {
        fs::create_dir_all(path)?;
        let db = redb::Database::create(Path::new(path).join("db"))?;
        let def = redb::TableDefinition::new("x");
        let txn = db.begin_write()?;
        {
            txn.open_table(def)?;
        }
        Ok(Redb(db, def))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>> {
        let txn = self.0.begin_read()?;
        if let Ok(table) = txn.open_table(self.1) {
            if let Some(val) = table.get(key)? {
                return Ok(Some(val.value().to_vec()));
            }
        }
        Ok(None)
    }

    fn range(&self, f: impl Fn(&[u8], &[u8]) -> Result<()>) -> Vec<Result<()>> {
        todo!()
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<()> {
        let txn = self.0.begin_write()?;
        {
            let mut table = txn.open_table(self.1)?;
            table.insert(key, value)?;
        }
        txn.commit()?;
        Ok(())
    }

    fn del(&self, _key: &[u8]) -> Result<()> {
        todo!()
    }
}

// Implement KVS for rocksdb
pub struct Rocksdb(rocksdb::DB);
impl KVS for Rocksdb {
    fn new(path: &str) -> Result<Self> {
        let db = rocksdb::DB::open_default(path)?;
        Ok(Rocksdb(db))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>> {
        Ok(self.0.get(key)?)
    }

    fn range(&self, f: impl Fn(&[u8], &[u8]) -> Result<()>) -> Vec<Result<()>> {
        self.0.iterator(rocksdb::IteratorMode::Start).map(|v| {
            match v {
                Ok(r) => f(&r.0, &r.1),
                Err(err) => Err(anyhow!(err)),
            }
        }).collect()
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<()> {
        self.0.put(key, value)?;
        Ok(())
    }

    fn del(&self, _key: &[u8]) -> Result<()> {
        todo!()
    }
}

// Implement KVS for persy
pub struct Persy(persy::Persy);
impl KVS for Persy {
    fn new(path: &str) -> Result<Self> {
        fs::create_dir_all(path)?;
        let persy = persy::Persy::open_or_create_with(Path::new(path).join("db"), persy::Config::new(), |persy| {
            let mut tx = persy.begin()?;
            tx.create_segment("main")?;
            tx.create_index::<persy::ByteVec, persy::ByteVec>("index", persy::ValueMode::Replace)?;
            tx.prepare()?.commit()?;
            Ok(())
        })?;
        Ok(Persy(persy))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>> {
        let mut tx = self.0.begin()?;
        let result = tx.one::<persy::ByteVec, persy::ByteVec>("index", &persy::ByteVec::new(key.to_vec()))?.map(|v| v.to_vec());
        Ok(result)
    }

    fn range(&self, f: impl Fn(&[u8], &[u8]) -> Result<()>) -> Vec<Result<()>> {
        todo!()
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<()> {
        let mut tx = self.0.begin()?;
        tx.put::<persy::ByteVec, persy::ByteVec>("index", persy::ByteVec::new(key.to_vec()), persy::ByteVec::new(value.to_vec()))?;
        tx.prepare()?.commit()?;
        Ok(())
    }

    fn del(&self, _key: &[u8]) -> Result<()> {
        todo!()
    }
}