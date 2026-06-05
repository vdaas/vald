// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
fn gcc_lib_dir() -> Option<String> {
    let arch = std::env::var("CARGO_CFG_TARGET_ARCH").unwrap_or_default();
    let triple = match arch.as_str() {
        "x86_64" => "x86_64-linux-gnu",
        "aarch64" => "aarch64-linux-gnu",
        "arm" => "arm-linux-gnueabihf",
        _ => return None,
    };
    for ver in &["15", "14", "13", "12", "11", "10"] {
        let path = format!("/usr/lib/gcc/{}/{}", triple, ver);
        if std::path::Path::new(&path).exists() {
            return Some(path);
        }
    }
    None
}

fn main() -> miette::Result<()> {
    let current_dir = std::env::current_dir().unwrap();
    println!("cargo:rustc-link-search=native={}", current_dir.display());

    cxx_build::bridge("src/lib.rs")
        .file("src/input.cpp")
        .flag_if_supported("-std=c++20")
        .flag_if_supported("-fopenmp")
        .flag_if_supported("-flto=thin")
        .flag_if_supported("-mavx2")
        .flag_if_supported("-mno-avx512f")
        .flag_if_supported("-mno-avx512dq")
        .flag_if_supported("-mno-avx512cd")
        .flag_if_supported("-mno-avx512bw")
        .flag_if_supported("-mno-avx512vl")
        .flag_if_supported("-DNGT_BFLOAT_DISABLED")
        .flag_if_supported("-DNGT_LARGE_DATASET")
        .compile("qbg-rs");

    println!("cargo:rustc-link-search=native=/usr/local/lib");

    let arch = std::env::var("CARGO_CFG_TARGET_ARCH").unwrap_or_default();
    if arch == "x86_64" {
        println!("cargo:rustc-link-search=native=/usr/lib/x86_64-linux-gnu");
    } else if arch == "aarch64" {
        println!("cargo:rustc-link-search=native=/usr/lib/aarch64-linux-gnu");
    }

    if let Some(gcc_dir) = gcc_lib_dir() {
        println!("cargo:rustc-link-search=native={}", gcc_dir);
    }

    println!("cargo:rustc-link-lib=static:+whole-archive=ngt");
    println!("cargo:rustc-link-lib=static=blas");
    println!("cargo:rustc-link-lib=static=lapack");
    println!("cargo:rustc-link-lib=static=gfortran");
    // libquadmath is only available on x86/x86_64
    if arch == "x86_64" {
        println!("cargo:rustc-link-lib=static=quadmath");
    }
    println!("cargo:rustc-link-lib=static=z");
    // NGT is built with GCC/libgomp; use gomp (not omp/libomp)
    println!("cargo:rustc-link-lib=static=gomp");
    println!("cargo:rerun-if-changed=src/*");

    Ok(())
}
