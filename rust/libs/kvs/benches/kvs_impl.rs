use std::{error::Error, fs, path::Path};
use kvs::KVS;

// Implement KVS for sled
pub struct Sled(sled::Db);
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
pub struct Kv(kv::Store);
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
pub struct Rkv(rkv::Rkv<rkv::backend::SafeModeEnvironment>, rkv::SingleStore<rkv::backend::SafeModeDatabase>);
impl KVS for Rkv {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> {
        fs::create_dir_all(path)?;
        let db = rkv::Rkv::new::<rkv::backend::SafeMode>(Path::new(path))?;
        let store = db.open_single("mydb", rkv::StoreOptions::create())?;
        Ok(Rkv(db, store))
    }

    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>, Box<dyn Error + '_>> {
        let reader = self.0.read()?;
        if let Some(val) = self.1.get(&reader, key)? {
            Ok(Some(val.to_bytes()?.to_vec()))
        } else {
            Ok(None)
        }
        
    }

    fn set(&self, key: &[u8], value: &[u8]) -> Result<(), Box<dyn Error + '_>> {
        let mut writer = self.0.write()?;
        self.1.put(&mut writer, key, &rkv::Value::Blob(value))?;
        writer.commit()?;
        Ok(())
    }
}

// Implement KVS for redb
pub struct Redb(redb::Database, redb::TableDefinition<'static, &'static [u8], &'static [u8]>);
impl KVS for Redb {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> {
        fs::create_dir_all(path)?;
        let db = redb::Database::create(Path::new(path).join("db"))?;
        let def = redb::TableDefinition::new("x");
        Ok(Redb(db, def))
    }
    fn get(&self, key: &[u8]) -> Result<Option<Vec<u8>>, Box<dyn Error>> {
        let txn = self.0.begin_read()?;
        let table = txn.open_table(self.1)?;
        if let Some(val) = table.get(key)? {
            Ok(Some(val.value().to_vec()))
        } else {
            Ok(None)
        }
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

// Implement KVS for rocksdb
pub struct Rocksdb(rocksdb::DB);
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
pub struct Persy(persy::Persy);
impl KVS for Persy {
    fn new(path: &str) -> Result<Self, Box<dyn Error>> {
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