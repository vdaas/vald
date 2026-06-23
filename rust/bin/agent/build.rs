//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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
fn main() {
    // Compile glibc compatibility shims so the binary runs on glibc < 2.38.
    // GCC 13 + glibc 2.38 headers emit __isoc23_strtoX references even for
    // C++17 builds; defining them here prevents runtime failures on Debian 12.
    cc::Build::new()
        .file("src/glibc_compat.c")
        .compile("glibc_compat");

    println!("cargo:rustc-link-lib=static=stdc++");
    println!("cargo:rerun-if-changed=build.rs");
    println!("cargo:rerun-if-changed=src/glibc_compat.c");
}
