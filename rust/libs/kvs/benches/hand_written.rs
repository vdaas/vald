use std::{any::type_name, path::Path, sync::{atomic::{AtomicBool, AtomicUsize, Ordering}, Arc}, thread, time::{Duration, Instant}};

use kvs::KVS;

mod kvs_impl;
use kvs_impl::*;

mod util;
use procfs::process::Status;
use rand::{thread_rng, Rng};
use util::{random_bytes, setup_kvs};

fn bencher<T: KVS + 'static>(name: &str, path: &Path, size: usize, kdim: usize, vdim: usize, db: T, ratio: f64, thread: usize, interval: u64, timer: u64) -> Vec<(usize, Duration, Status, u64)> {
    let me = Arc::new(procfs::process::Process::myself().unwrap());
    let shutdown = Arc::new(AtomicBool::new(false));
    let begin = Arc::new(Instant::now());
    let db = Arc::new(db);
    let get_count = Arc::new(AtomicUsize::new(0));
    let set_count = Arc::new(AtomicUsize::new(0));
    let mut threads = vec![];
    for _ in 0..thread {
        let db = db.clone();
        let shutdown = shutdown.clone();
        let get_count = get_count.clone();
        let set_count = set_count.clone();
        let t = thread::spawn(move || {
            let mut rng = thread_rng();
            while !shutdown.load(Ordering::Relaxed) {
                let choice: f64 = rng.gen_range(0.0..1.0);
                match choice {
                    v if v <= ratio => {
                        let mut key = set_count.fetch_add(1, Ordering::Release).to_ne_bytes().to_vec();
                        key.resize_with(kdim, Default::default);
                        db.set(&key, &random_bytes(vdim)).unwrap_or_default();
                    }
                    _ => {
                        get_count.fetch_add(1, Ordering::Release);
                        let mut key = rng.gen_range(0..size).to_ne_bytes().to_vec();
                        key.resize_with(kdim, Default::default);
                        db.get(&key).unwrap_or_default();
                    }
                }
            }
        });
        threads.push(t);
    }

    {
        let shutdown = shutdown.clone();
        thread::spawn(move || {
            thread::sleep(Duration::from_secs(timer));
            shutdown.store(true, Ordering::SeqCst);
        });    
    }

    let mut progress = vec![];
    while {
        thread::sleep(Duration::from_secs(interval));
        let p = begin.elapsed();
        let st = me.status().unwrap();
        let dir_size = fs_extra::dir::get_size(path).unwrap_or_default();
        let set_count = set_count.load(Ordering::Relaxed);
        println!("{},set,{},{},{},{},{},{},{}", name, kdim, vdim, thread, set_count, p.as_nanos(), st.vmrss.unwrap(), dir_size/1024);
        let get_count = get_count.load(Ordering::Relaxed);
        println!("{},get,{},{},{},{},{},{},{}", name, kdim, vdim, thread, get_count, p.as_nanos(), st.vmrss.unwrap(), dir_size/1024);
        progress.push((set_count, p, st, dir_size));
        
        !shutdown.load(Ordering::Relaxed) 
    } {}

    for t in threads {
        t.join().unwrap();
    }

    progress
}

fn benchmark<T: KVS + 'static>(size: usize, kdims: &[usize], vdims: &[usize], threads: &[usize], ratio: f64, interval: u64, timer: u64) {
    let name = type_name::<T>().split("::").last().unwrap();
    for &kdim in kdims {
        for &vdim in vdims {
            for &thread in threads {
                let (path, db) = setup_kvs::<T>(format!("{}-{}-{}", kdim, vdim, thread).as_str());
                bencher(name, path.as_path(), size, kdim, vdim, db, ratio, thread, interval, timer);
            }
        }
    }
}

fn main() {
    let size = 1 << 24;
    let kdims: &[usize] = &[1 << 2, 1 << 4, 1 << 6, 1 << 8, 1 << 10, 1 << 12, 1 << 14];
    let vdims: &[usize] = &[1 << 10];
    let threads: &[usize] = &[16];
    let ratio = 0.5;
    let interval = 5;
    let timer = 15;
    println!("name,operation,key size(B),value size(B),thread,operation count,time(ns),vmrss(KB),file size(B)");
    benchmark::<Kv>(size, kdims, vdims, threads, ratio, interval, timer);
    benchmark::<Kv2>(size, kdims, vdims, threads, ratio, interval, timer);
    //parallel_benchmark::<Persy>(size, kdims, vdims, threads, ratio, interval, timer);
    //parallel_benchmark::<Redb>(size, kdims, vdims, threads, ratio, interval, timer);
    //parallel_benchmark::<Rkv>(size, kdims, vdims, threads);
    benchmark::<Rocksdb>(size, kdims, vdims, threads, ratio, interval, timer);
    benchmark::<Sled>(size, kdims, vdims, threads, ratio, interval, timer);
}
