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

/// LLVM version range to probe when searching for libomp.so.
/// Updated as new LLVM releases ship distro packages.
const LLVM_VERSION_MIN: u32 = 14;
const LLVM_VERSION_MAX: u32 = 25;

/// Emit the correct OpenMP linker flags for the given static library.
///
/// Ubuntu 25+ ships NGT compiled with clang's `-fopenmp=libomp`, which
/// emits `__kmpc_*` symbols. Earlier Ubuntu releases use GCC's libgomp
/// (GOMP_* symbols). We inspect the static archive with `nm` to decide
/// which runtime to link. If `nm` is unavailable we fall back to libgomp
/// because it is present on all supported build environments.
pub fn link_openmp(ngt_static_lib: &str) {
    // Detect whether NGT was compiled with LLVM OpenMP (libomp) or GCC OpenMP (libgomp).
    // Ubuntu 25+ with clang uses -fopenmp=libomp by default, generating __kmpc_* symbols.
    // Ubuntu 24 and earlier use GNU OpenMP (libgomp) with GOMP_* symbols.
    let uses_llvm_omp = std::process::Command::new("nm")
        .args(["-u", ngt_static_lib])
        .output()
        .ok()
        .map(|o| String::from_utf8_lossy(&o.stdout).contains("__kmpc_"))
        .unwrap_or(false);

    if uses_llvm_omp {
        // Search newest-first so we pick up the highest available LLVM version.
        for v in (LLVM_VERSION_MIN..=LLVM_VERSION_MAX).rev() {
            let dir = format!("/usr/lib/llvm-{}/lib", v);
            if std::path::Path::new(&format!("{}/libomp.so", dir)).exists() {
                println!("cargo:rustc-link-search=native={}", dir);
                break;
            }
        }
        println!("cargo:rustc-link-lib=dylib=omp");
    } else {
        println!("cargo:rustc-link-lib=dylib=gomp");
    }
}
