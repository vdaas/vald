//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

use anyhow::*;
use kv::{Bincode, Bucket, Config, Integer, Store};
use serde::{Deserialize, Serialize};

pub trait BidiMap {
    fn get(&self, key: &str) -> Result<(u32, u128)>;
    fn get_inverse(&self, key: u32) -> Result<(String, u128)>;
    fn set(&mut self, key: &str, value: u32, timestamp: u128) -> Result<()>;
    fn delete(&mut self, key: &str) -> Result<u32>;
    fn delete_inverse(&mut self, key: u32) -> Result<String>;
    fn range<F: Fn(&str, u32, u128) -> Result<()>>(&self, f: F) -> Result<()>;
    fn len(&self) -> Result<usize>;
    fn close(&self) -> Result<()>;
}

struct Bidi<'a> {
    ou: Bucket<'a, Integer, Bincode<ValueOu>>,
    uo: Bucket<'a, &'a str, Bincode<ValueUo>>,
}

#[derive(Serialize, Deserialize)]
struct ValueOu {
    value: String,
    timestamp: u128,
}

#[derive(Serialize, Deserialize)]
struct ValueUo {
    value: u32,
    timestamp: u128,
}

pub fn new(path: &str) -> Result<impl BidiMap> {
    let s = Store::new(Config::new(path))?;
    Ok(Bidi {
        ou: s.bucket::<Integer, Bincode<ValueOu>>(Some("ou"))?,
        uo: s.bucket::<&str, Bincode<ValueUo>>(Some("uo"))?,
    })
}

impl BidiMap for Bidi<'_> {
    fn get(&self,  key: &str) -> Result<(u32, u128)> {
        if let Some(value) = self.uo.get(&key)? {
            Ok((value.0.value, value.0.timestamp))
        } else {
            Err(anyhow!("key not found"))
        }
    }

    fn get_inverse(&self, key: u32) -> Result<(String, u128)> {
        if let Some(value) = self.ou.get(&Integer::from(key))? {
            Ok((value.0.value, value.0.timestamp))
        } else {
            Err(anyhow!("key not found"))
        }
    }

    fn set(&mut self, key: &str, value: u32, timestamp: u128) -> Result<()> {
        self.uo.set(&key, &Bincode(ValueUo{
            value,
            timestamp,
        }))?;
        self.ou.set(&Integer::from(value), &Bincode(ValueOu {
            value: key.to_owned(),
            timestamp,
        }))?;
        Ok(())
    }

    fn delete(&mut self, key: &str) -> Result<u32> {
        let value = self.uo.remove(&key)?;
        if let Some(value) = value {
            self.ou.remove(&Integer::from(value.0.value))?;
            Ok(value.0.value)
        } else {
            Err(anyhow!("key not found"))
        }
    }

    fn delete_inverse(&mut self, key: u32) -> Result<String> {
        if let Some(value) = self.ou.remove(&Integer::from(key))? {
            self.uo.remove(&value.0.value.as_str())?;
            Ok(value.0.value)
        } else {
            Err(anyhow!("key not found"))
        }
    }
    
    fn range<F: Fn(&str, u32, u128) -> Result<()>>(&self, f: F) -> Result<()> {
        for item in self.uo.iter() {
            let item = item?;
            let key = item.key::<&str>()?;
            let value = item.value::<Bincode<ValueUo>>()?;
            f(key, value.0.value, value.0.timestamp)?
        }
        Ok(())
    }

    fn len(&self) -> Result<usize> {
        let ou_len = self.ou.len();
        let uo_len = self.uo.len();
        ensure!(ou_len == uo_len, "mismatch length, ou:{}, uo:{}", ou_len, uo_len);
        Ok(ou_len)
    }

    fn close(&self) -> Result<()> {
        Ok(())
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    
    use std::{fs::remove_dir_all, sync::{atomic::{AtomicUsize, Ordering}, Arc}, time::SystemTime};
    use tokio::sync::RwLock;

    fn setup(path: &str) -> (String, impl Fn() -> Result<(), std::io::Error>) {
        let s = format!("./test/{}", path);
        let _ = remove_dir_all(&s);
        (s, || { remove_dir_all("./test") })
    }

    #[test]
    fn test_unary() -> Result<()> {
        let (path, teardown) = setup("unary");
        let mut db = new(&path)?;

        assert!(std::path::Path::new(path.as_str()).exists());
        assert_eq!(db.len()?, 0);

        let key1 = "test-1";
        let value1 = 10;
        let timestamp1 = SystemTime::now().duration_since(SystemTime::UNIX_EPOCH)?.as_nanos();
        db.set(key1, value1, timestamp1)?;
        let key2 = "test-2";
        let value2 = 20;
        let timestamp2 = SystemTime::now().duration_since(SystemTime::UNIX_EPOCH)?.as_nanos();
        db.set(key2, value2, timestamp2)?;

        assert_eq!(db.len()?, 2);

        let got = db.get(&key1)?;
        assert_eq!(got.0, value1);
        assert_eq!(got.1, timestamp1);

        let got = db.get_inverse(value1)?;
        assert_eq!(got.0, key1);
        assert_eq!(got.1, timestamp1);

        let del = db.delete(&key1)?;
        assert_eq!(del, value1);

        let del = db.delete_inverse(value2)?;
        assert_eq!(del, key2);

        assert_eq!(db.len()?, 0);

        db.close()?;

        teardown()?;

        Ok(())
    }

    #[test]
    fn test_range() -> Result<()> {
        let (path, teardown) = setup("range");
        let mut db = new(&path)?;

        assert!(std::path::Path::new(path.as_str()).exists());
        assert_eq!(db.len()?, 0);

        let mut inputs = vec![];
        for i in 0..10 {
            let key = format!("key-{}", i);
            let value = i;
            let timestamp = SystemTime::now().duration_since(SystemTime::UNIX_EPOCH)?.as_nanos();
            db.set(&key, value, timestamp)?;
            inputs.push((key, value, timestamp));
        };

        assert_eq!(db.len()?, 10);

        db.range(|key, value, timestamp| {
            static COUNT: AtomicUsize = AtomicUsize::new(0);
            let i = COUNT.fetch_add(1, Ordering::Relaxed);
            assert_eq!(inputs[i].0, key);
            assert_eq!(inputs[i].1, value);
            assert_eq!(inputs[i].2, timestamp);
            Ok(())
        })?;

        teardown()?;

        Ok(())
    }

    #[tokio::test]
    async fn test_multithread() {
        let (path, teardown) = setup("mt");
        let db = Arc::new(RwLock::new(new(&path).unwrap()));

        assert!(std::path::Path::new(path.as_str()).exists());
        assert_eq!(db.read().await.len().unwrap(), 0);

        let mut inputs = vec![];
        for i in 0..100 {
            let key = format!("key-{}", i);
            let value = i;
            let timestamp = SystemTime::now().duration_since(SystemTime::UNIX_EPOCH).unwrap().as_nanos();
            inputs.push((key, value, timestamp));
        }
        let inputs = Arc::new(RwLock::new(inputs));

        let mut handlers = vec![];
        {
            let db = Arc::clone(&db);
            let inputs = Arc::clone(&inputs);
            let handler = tokio::spawn(async move {
                let inputs = Arc::clone(&inputs);
                for input in inputs.read().await.iter() {
                    let _ = db.write().await.set(&input.0, input.1, input.2);
                }
            });
            handlers.push(handler);
        }

        {
            let db = Arc::clone(&db);
            let inputs = Arc::clone(&inputs);
            let handler = tokio::spawn(async move {
                let inputs = Arc::clone(&inputs);
                for input in inputs.read().await.iter().step_by(2) {
                    let value = db.read().await.get(&input.0).unwrap();
                    assert_eq!(value.0, input.1);
                    assert_eq!(value.1, input.2);
                }
            });
            handlers.push(handler);
        }

        {
            let db = Arc::clone(&db);
            let inputs = Arc::clone(&inputs);
            let handler = tokio::spawn(async move {
                let inputs = Arc::clone(&inputs);
                for input in inputs.read().await.iter().skip(1).step_by(2) {
                    let value = db.read().await.get_inverse(input.1).unwrap();
                    assert_eq!(value.0, input.0);
                    assert_eq!(value.1, input.2);
                }
            });
            handlers.push(handler);
        }

        for handler in handlers {
            handler.await.unwrap();
        }

        assert_eq!(db.read().await.len().unwrap(), 100);

        let mut handlers = vec![];
        {
            let db = Arc::clone(&db);
            let inputs = Arc::clone(&inputs);
            let handler = tokio::spawn(async move {
                let inputs = Arc::clone(&inputs);
                for input in inputs.read().await.iter().step_by(2) {
                    let value = db.write().await.delete(&input.0).unwrap();
                    assert_eq!(value, input.1);
                }
            });
            handlers.push(handler);
        }

        {
            let db = Arc::clone(&db);
            let inputs = Arc::clone(&inputs);
            let handler = tokio::spawn(async move {
                let inputs = Arc::clone(&inputs);
                for input in inputs.read().await.iter().skip(1).step_by(2) {
                    let value = db.write().await.delete_inverse(input.1).unwrap();
                    assert_eq!(value, input.0);
                }
            });
            handlers.push(handler);
        }

        for handler in handlers {
            handler.await.unwrap();
        }

        assert_eq!(db.read().await.len().unwrap(), 0);

        let _ = teardown();
    }
}