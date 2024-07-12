use std::{
    path::PathBuf, sync::{
        atomic::{AtomicBool, AtomicUsize, Ordering},
        Arc,
    }, thread
};

use kvs::KVS;
use rand::{thread_rng, Rng};

mod kvs_impl;
use kvs_impl::*;

mod util;
use util::*;

static TOTAL: AtomicUsize = AtomicUsize::new(0);
static GET: AtomicUsize = AtomicUsize::new(0);
static SET: AtomicUsize = AtomicUsize::new(0);
static DEL: AtomicUsize = AtomicUsize::new(0);

const USAGE: &str = "
Usage: stress [--threads=<#>] [--burn-in] [--duration=<s>] \
    [--key-len=<l>] [--val-len=<l>] \
    [--get-prop=<p>] \
    [--set-prop=<p>] \
    [--del-prop=<p>] \
    [--entries=<n>] \
    [--sequential] \
    [--total-ops=<n>] \
    [--flush-every=<ms>]

Options:
    --threads=<#>      Number of threads [default: 4].
    --burn-in          Don't halt until we receive a signal.
    --duration=<s>     Seconds to run for [default: 10].
    --key-len=<l>      The length of keys [default: 10].
    --val-len=<l>      The length of values [default: 100].
    --get-prop=<p>     The relative proportion of get requests [default: 94].
    --set-prop=<p>     The relative proportion of set requests [default: 2].
    --del-prop=<p>     The relative proportion of del requests [default: 1].
    --entries=<n>      The total keyspace [default: 100000].
    --sequential       Run the test in sequential mode instead of random.
    --total-ops=<n>    Stop test after executing a total number of operations.
    --flush-every=<m>  Flush and sync the database every ms [default: 200].
    --cache-mb=<mb>    Size of the page cache in megabytes [default: 1024].
";

#[derive(Debug, Clone, Copy)]
struct Args {
    threads: usize,
    burn_in: bool,
    duration: u64,
    key_len: usize,
    val_len: usize,
    get_prop: usize,
    set_prop: usize,
    del_prop: usize,
    entries: usize,
    sequential: bool,
    total_ops: Option<usize>,
    flush_every: u64,
    cache_mb: usize,
}

impl Default for Args {
    fn default() -> Args {
        Args {
            threads: 4,
            burn_in: false,
            duration: 10,
            key_len: 10,
            val_len: 100,
            get_prop: 80,
            set_prop: 20,
            del_prop: 0,
            entries: 100000,
            sequential: false,
            total_ops: None,
            flush_every: 200,
            cache_mb: 1024,
        }
    }
}

fn parse<'a, I, T>(mut iter: I) -> T
where
    I: Iterator<Item = &'a str>,
    T: std::str::FromStr,
    <T as std::str::FromStr>::Err: std::fmt::Debug,
{
    iter.next().expect(USAGE).parse().expect(USAGE)
}

impl Args {
    fn parse() -> Args {
        let mut args = Args::default();
        for raw_arg in std::env::args().skip(1) {
            let mut splits = raw_arg[2..].split('=');
            match splits.next().unwrap() {
                "threads" => args.threads = parse(&mut splits),
                "burn-in" => args.burn_in = true,
                "duration" => args.duration = parse(&mut splits),
                "key-len" => args.key_len = parse(&mut splits),
                "val-len" => args.val_len = parse(&mut splits),
                "get-prop" => args.get_prop = parse(&mut splits),
                "set-prop" => args.set_prop = parse(&mut splits),
                "del-prop" => args.del_prop = parse(&mut splits),
                "entries" => args.entries = parse(&mut splits),
                "sequential" => args.sequential = true,
                "total-ops" => args.total_ops = Some(parse(&mut splits)),
                "flush-every" => args.flush_every = parse(&mut splits),
                "cache-mb" => args.cache_mb = parse(&mut splits),
                _ => {},
            }
        }
        args
    }
}

fn report(shutdown: Arc<AtomicBool>, path: Arc<PathBuf>) {
    let mut last = 0;
    let mut last_get = 0;
    let mut last_set = 0;
    let mut last_del = 0;
    let me = procfs::process::Process::myself().unwrap();

    while !shutdown.load(Ordering::Relaxed) {
        thread::sleep(std::time::Duration::from_secs(1));
        let total = TOTAL.load(Ordering::Acquire);
        let cur_get = GET.load(Ordering::Acquire);
        let cur_set = SET.load(Ordering::Acquire);
        let cur_del = DEL.load(Ordering::Acquire);

        println!(
            "did {}({}/{}/{}) ops, {}mb RSS, {}mb File",
            (total - last), (cur_get - last_get), (cur_set - last_set), (cur_del - last_del),
            me.status().unwrap().vmrss.unwrap() / 1024,
            fs_extra::dir::get_size(path.as_path()).unwrap() / 1024 / 1024
        );

        last = total;
        last_get = cur_get;
        last_set = cur_set;
        last_del = cur_del;
    }
}

fn run<T: KVS>(args: Args, db: Arc<T>, shutdown: Arc<AtomicBool>) {
    let get_max = args.get_prop;
    let set_max = get_max + args.set_prop;
    let del_max = set_max + args.del_prop;

    let keygen = |len| -> Vec<u8> {
        static SEQ: AtomicUsize = AtomicUsize::new(0);
        let i = if args.sequential {
            SEQ.fetch_add(1, Ordering::Relaxed)
        } else {
            thread_rng().gen::<usize>()
        } % args.entries;

        let mut arg = i.to_ne_bytes().to_vec();
        arg.resize_with(len, Default::default);

        arg
    };

    let mut rng = thread_rng();

    while !shutdown.load(Ordering::Relaxed) {
        let _op = TOTAL.fetch_add(1, Ordering::Release);
        let key = keygen(args.key_len);
        let choice = rng.gen_range(0..del_max + 1);

        match choice {
            v if v <= get_max => {
                GET.fetch_add(1, Ordering::Release);
                db.get(&key).unwrap();
            }
            v if v > get_max && v <= set_max => {
                SET.fetch_add(1, Ordering::Release);
                let value = random_bytes(args.val_len);
                db.set(&key, &value).unwrap();
            }
            v if v > set_max && v <= del_max => {
                DEL.fetch_add(1, Ordering::Release);
                db.del(&key).unwrap();
            }
            _ => {
            }
        }
    }
}

fn benchmark_main<T: KVS + 'static>(args: Args) {
    let shutdown = Arc::new(AtomicBool::new(false));

    dbg!(args);

    let (path, db) = setup_kvs::<T>(args.key_len, args.val_len);
    let db = Arc::new(db);
    let path = Arc::new(path);

    let mut threads = vec![];

    let now = std::time::Instant::now();

    let n_threads = args.threads;

    for i in 0..=n_threads {
        let db = db.clone();
        let shutdown = shutdown.clone();

        let t = if i == 0 {
            let path = path.clone();
            thread::Builder::new()
                .name("reporter".into())
                .spawn(move || report(shutdown, path))
                .unwrap()
        } else {
            thread::spawn(move || run(args, db, shutdown))
        };

        threads.push(t);
    }

    if let Some(ops) = args.total_ops {
        assert!(!args.burn_in, "don't set both --burn-in and --total-ops");
        while TOTAL.load(Ordering::Relaxed) < ops {
            thread::sleep(std::time::Duration::from_millis(50));
        }
        shutdown.store(true, Ordering::SeqCst);
    } else if !args.burn_in {
        thread::sleep(std::time::Duration::from_secs(args.duration));
        shutdown.store(true, Ordering::SeqCst);
    }

    for t in threads.into_iter() {
        t.join().unwrap();
    }
    let ops = TOTAL.load(Ordering::SeqCst);
    let time = now.elapsed().as_secs() as usize;

    println!(
        "did {} total ops in {} seconds. {} ops/s",
        ops,
        time,
        ((ops * 1_000) / (time * 1_000))
    );
}

fn main() {
    let args = Args::parse();

    benchmark_main::<Rkv>(args);
}
