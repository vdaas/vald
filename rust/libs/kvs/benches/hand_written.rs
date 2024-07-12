use std::{any::type_name, path::Path, time::{Duration, Instant}};

use kvs::KVS;

mod kvs_impl;
use kvs_impl::*;

mod util;
use procfs::process::Status;
use util::{random_bytes, setup_kvs};

fn bencher(name: &str, operation: &str, path: &Path, size: usize, kdim: usize, vdim: usize, f: impl Fn(&[u8])) -> Vec<(usize, Duration, Status, u64)> {
    let me = procfs::process::Process::myself().unwrap();
    let mut progress = vec![];
    let begin = Instant::now();
    for i in 0..size {
        let mut arg = i.to_ne_bytes().to_vec();
        arg.resize_with(kdim, Default::default);
        f(&arg);
        if i & (i + 1) == 0 && i + 1 >= 1 << 14 {
            let p = begin.elapsed();
            let st = me.status().unwrap();
            let dir_size = fs_extra::dir::get_size(path).unwrap();
            println!("{},{},{},{},{},{},{},{}", name, operation, kdim, vdim, i+1, p.as_nanos(), st.vmrss.unwrap(), dir_size/1024);
            progress.push((i+1, p, st, dir_size));
        }
    }
    progress
}

fn monotonic_benchmark<T: KVS + 'static>(size: usize, kdims: &[usize], vdims: &[usize]) {
    let name = type_name::<T>().split("::").last().unwrap();
    for &kdim in kdims {
        for &vdim in vdims {
            let (path, db) = setup_kvs::<T>(kdim, vdim);
            bencher(name, "set", path.as_path(), size, kdim, vdim, |k| { db.set(k, &random_bytes(vdim)).unwrap(); });
            bencher(name, "get", path.as_path(), size, kdim, vdim, |k| { db.get(k).unwrap(); });
        }
    }
}

fn main() {
    let size = 1 << 16;
    let kdims: &[usize] = &[1 << 3, 1 << 10];
    let vdims: &[usize] = &[1 << 8, 1 << 12];
    println!("name,operation,key size(B),value size(B),insert size,time(ns),vmrss(KB),file size(B)");
    monotonic_benchmark::<Kv>(size, kdims, vdims);
    monotonic_benchmark::<Persy>(size, kdims, vdims);
    monotonic_benchmark::<Redb>(size, kdims, vdims);
    monotonic_benchmark::<Rkv>(size, kdims, vdims);
    monotonic_benchmark::<Rocksdb>(size, kdims, vdims);
    monotonic_benchmark::<Sled>(size, kdims, vdims);
}