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

use std::env;
use std::fs;
use std::path::PathBuf;
use std::process::Command;

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let manifest_dir = PathBuf::from(env::var("CARGO_MANIFEST_DIR")?);
    let repo_root = manifest_dir
        .join("../../..")
        .canonicalize()
        .unwrap_or(manifest_dir.clone());

    println!(
        "cargo:rerun-if-changed={}",
        repo_root.join("versions/NGT_VERSION").display()
    );
    println!(
        "cargo:rerun-if-changed={}",
        repo_root.join("versions/VALD_VERSION").display()
    );

    println!("cargo:rustc-env=VALD_REPO_ROOT={}", repo_root.display());

    if let Ok(ngt_version) = fs::read_to_string(repo_root.join("versions/NGT_VERSION")) {
        let ngt_version = ngt_version.trim();
        if !ngt_version.is_empty() {
            println!("cargo:rustc-env=VALD_ALGORITHM_INFO=NGT-{}", ngt_version);
        }
    }

    if let Ok(vald_version) = fs::read_to_string(repo_root.join("versions/VALD_VERSION")) {
        let vald_version = vald_version.trim();
        if !vald_version.is_empty() {
            println!("cargo:rustc-env=VALD_VERSION={}", vald_version);
        }
    }

    let build_time = chrono::Utc::now().format("%Y/%m/%d_%H:%M:%S%z").to_string();
    println!("cargo:rustc-env=BUILD_TIME={}", build_time);

    if let Some(git_commit) = command_output("git", &["rev-parse", "HEAD"], &repo_root) {
        println!("cargo:rustc-env=GIT_COMMIT={}", git_commit);
    }

    if let Some(rustc_version) = command_output("rustc", &["--version"], &repo_root) {
        println!("cargo:rustc-env=RUSTC_VERSION={}", rustc_version);
    }

    if let Some(cpu_flags) = read_cpu_flags("/proc/cpuinfo") {
        println!("cargo:rustc-env=BUILD_CPU_INFO_FLAGS={}", cpu_flags);
    }

    println!("cargo:rustc-env=CGO_ENABLED=true");
    println!("cargo:rustc-env=CGO_CALL=1");

    Ok(())
}

fn command_output(cmd: &str, args: &[&str], current_dir: &PathBuf) -> Option<String> {
    let output = Command::new(cmd)
        .args(args)
        .current_dir(current_dir)
        .output()
        .ok()?;
    if !output.status.success() {
        return None;
    }
    let value = String::from_utf8_lossy(&output.stdout).trim().to_string();
    if value.is_empty() { None } else { Some(value) }
}

fn read_cpu_flags(path: &str) -> Option<String> {
    let contents = fs::read_to_string(path).ok()?;
    for line in contents.lines() {
        if let Some(rest) = line.strip_prefix("flags") {
            let mut parts = rest.splitn(2, ':');
            let _ = parts.next();
            let flags = parts.next()?.trim();
            if !flags.is_empty() {
                return Some(flags.to_string());
            }
        }
    }
    None
}
