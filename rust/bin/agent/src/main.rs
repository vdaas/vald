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
#![cfg_attr(test, allow(missing_docs))]

use clap::Parser;

mod version;

#[derive(Parser, Debug)]
#[command(name = "agent")]
#[command(about = "Vald Agent - Vector Search Engine", long_about = None)]
struct Args {
    /// Print version information
    #[arg(short, long)]
    version: bool,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let raw_args: Vec<String> = std::env::args().collect();
    if version::is_version_request(&raw_args) {
        version::print_version_info();
        return Ok(());
    }

    let args = Args::parse();

    if args.version {
        version::print_version_info();
        return Ok(());
    }

    let settings = ::config::Config::builder()
        .add_source(::config::File::with_name("/etc/server/config.yaml"))
        .build()?;

    let mut config: agent::config::AgentConfig = settings.try_deserialize()?;
    config.bind();
    config.validate()?;

    agent::serve(config).await
}

