//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
fn main() -> miette::Result<()> {
    let current_dir = std::env::current_dir().unwrap();
    println!("cargo:rustc-link-search=native={}", current_dir.display());

    cxx_build::bridge("src/lib.rs")
        .file("src/input.cpp")
        .flag_if_supported("-std=c++20")
        .flag_if_supported("-fopenmp")
        .flag_if_supported("-DNGT_BFLOAT_DISABLED")
        .compile("qbg-rs");

    println!("cargo:rustc-link-search=native=/usr/local/lib");
    println!("cargo:rustc-link-lib=static=ngt");
    println!("cargo:rustc-link-lib=blas");
    println!("cargo:rustc-link-lib=lapack");
    println!("cargo:rustc-link-lib=dylib=gomp");
    println!("cargo:rerun-if-changed=src/*");

    Ok(())
}
