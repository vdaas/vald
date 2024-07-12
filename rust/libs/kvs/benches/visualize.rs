use std::env::current_dir;

use anyhow::Result;
use polars::prelude::*;

fn main() -> Result<()> {
    let path = current_dir()?;
    let df = CsvReadOptions::default()
        .with_has_header(true)
        .try_into_reader_with_file_path(Some(path.join("benchmark.csv").into()))?
        .finish()?;

    println!("{:?}", df);

    

    Ok(())
}