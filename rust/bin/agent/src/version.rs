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

use backtrace::Backtrace;
use chrono::Local;
use std::collections::BTreeMap;
use std::env;

const SERVER_NAME: &str = "agent qbg";
const STACK_TRACE_LIMIT: usize = 4;

pub fn is_version_request(args: &[String]) -> bool {
    args.iter()
        .skip(1)
        .any(|arg| matches!(arg.as_str(), "-version" | "--version" | "-v" | "-V"))
}

pub fn print_version_info() {
    println!("{}", build_version_output());
}

fn build_version_output() -> String {
    let mut info = BTreeMap::new();

    insert_value(
        &mut info,
        "algorithm info",
        option_env!("VALD_ALGORITHM_INFO"),
    );
    insert_value_owned(
        &mut info,
        "build cpu info flags",
        option_env!("BUILD_CPU_INFO_FLAGS").and_then(format_cpu_flags),
    );
    insert_value(&mut info, "build time", option_env!("BUILD_TIME"));
    insert_value(&mut info, "cgo call", option_env!("CGO_CALL"));
    insert_value(&mut info, "cgo enabled", option_env!("CGO_ENABLED"));
    insert_value(&mut info, "git commit", option_env!("GIT_COMMIT"));
    insert_value(&mut info, "go arch", Some(env::consts::ARCH));
    insert_value_owned(
        &mut info,
        "go max procs",
        Some(available_parallelism().to_string()),
    );
    insert_value(&mut info, "go os", Some(env::consts::OS));
    insert_value(&mut info, "go version", option_env!("RUSTC_VERSION"));
    insert_value(&mut info, "goroutine count", Some("1"));
    insert_value_owned(
        &mut info,
        "runtime cpu cores",
        Some(available_parallelism().to_string()),
    );
    insert_value(&mut info, "server name", Some(SERVER_NAME));
    insert_value(
        &mut info,
        "vald version",
        option_env!("VALD_VERSION").or(Some(env!("CARGO_PKG_VERSION"))),
    );

    for (index, trace) in collect_stack_traces().into_iter().enumerate() {
        let key = format!("stack trace-{:03}", index);
        let value = format!(
            "{}\t{}#L{}\t{}",
            trace.url, trace.file, trace.line, trace.func_name
        );
        info.insert(key, value);
    }

    let width = info.keys().map(|k| k.len()).max().unwrap_or(0);
    let mut lines = Vec::with_capacity(info.len());
    for (key, value) in info {
        if !value.is_empty() {
            lines.push(format!("{key:<width$} ->\t{value}", width = width));
        }
    }

    let now = Local::now().format("%Y-%m-%d %H:%M:%S");
    format!("{}     [INFO]:\n{}", now, lines.join("\n"))
}

fn insert_value(map: &mut BTreeMap<String, String>, key: &str, value: Option<&str>) {
    if let Some(value) = value {
        let value = value.trim();
        if !value.is_empty() {
            map.insert(key.to_string(), value.to_string());
        }
    }
}

fn insert_value_owned(map: &mut BTreeMap<String, String>, key: &str, value: Option<String>) {
    if let Some(value) = value {
        let value = value.trim();
        if !value.is_empty() {
            map.insert(key.to_string(), value.to_string());
        }
    }
}

fn available_parallelism() -> usize {
    std::thread::available_parallelism()
        .map(|n| n.get())
        .unwrap_or(1)
}

fn format_cpu_flags(flags: &str) -> Option<String> {
    let flags = flags
        .split_whitespace()
        .filter(|flag| !flag.is_empty())
        .collect::<Vec<_>>();
    if flags.is_empty() {
        None
    } else {
        Some(format!("[{}]", flags.join(" ")))
    }
}

struct StackTraceEntry {
    url: String,
    file: String,
    line: u32,
    func_name: String,
}

fn collect_stack_traces() -> Vec<StackTraceEntry> {
    let bt = Backtrace::new();
    let mut traces = Vec::new();

    for frame in bt.frames() {
        for symbol in frame.symbols() {
            let file = match symbol.filename() {
                Some(file) => file.display().to_string(),
                None => continue,
            };
            let line = match symbol.lineno() {
                Some(line) => line as u32,
                None => continue,
            };
            let func_name = symbol
                .name()
                .map(|name| name.to_string())
                .unwrap_or_else(|| "unknown".to_string());
            if should_skip_frame(&func_name) {
                continue;
            }

            let url = build_stack_url(&file, line);
            traces.push(StackTraceEntry {
                url,
                file,
                line,
                func_name,
            });
            if traces.len() >= STACK_TRACE_LIMIT {
                return traces;
            }
        }
    }

    traces
}

fn should_skip_frame(func_name: &str) -> bool {
    func_name.contains("version::") || func_name.contains("print_version_info")
}

fn build_stack_url(file: &str, line: u32) -> String {
    let repo_root = option_env!("VALD_REPO_ROOT").unwrap_or("");
    let git_commit = option_env!("GIT_COMMIT").unwrap_or("main");

    if !repo_root.is_empty() {
        let repo_root = repo_root.replace('\\', "/");
        let file_norm = file.replace('\\', "/");
        if let Some(relative) = file_norm.strip_prefix(&repo_root) {
            let relative = relative.trim_start_matches('/');
            return format!(
                "https://github.com/vdaas/vald/blob/{}/{}#L{}",
                git_commit, relative, line
            );
        }
    }

    format!("{}#L{}", file, line)
}
