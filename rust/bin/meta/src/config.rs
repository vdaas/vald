//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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
use std::collections::HashMap;
use std::time::Duration;

use observability::config as obs_config;
use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Config {
    #[serde(default)]
    pub observability: ObservabilityConfig,

    #[serde(default = "default_server_addr")]
    pub server_addr: String,

    #[serde(default = "default_database_path")]
    pub database_path: String,
}

#[derive(Debug, Default, Deserialize)]
pub struct ObservabilityConfig {
    #[serde(default)]
    pub enabled: bool,
    #[serde(default)]
    pub endpoint: String,
    #[serde(default)]
    pub attributes: HashMap<String, String>,
    #[serde(default)]
    pub tracer: Tracer,
    #[serde(default)]
    pub meter: Meter,
}

#[derive(Debug, Default, Deserialize)]
pub struct Tracer {
    #[serde(default)]
    pub enabled: bool,
}

#[derive(Debug, Default, Deserialize)]
pub struct Meter {
    #[serde(default)]
    pub enabled: bool,
    pub export_duration: Option<DurationConfig>,
    pub export_timeout_duration: Option<DurationConfig>,
}

#[derive(Debug, Deserialize)]
pub struct DurationConfig {
    pub secs: u64,
    pub nanos: u32,
}

fn default_server_addr() -> String {
    "[::1]:8095".to_string()
}

fn default_database_path() -> String {
    "/tmp/meta/database".to_string()
}

impl Default for Config {
    fn default() -> Self {
        Self {
            observability: ObservabilityConfig::default(),
            server_addr: default_server_addr(),
            database_path: default_database_path(),
        }
    }
}

impl Config {
    pub fn load(path: &str) -> Result<Self, Box<dyn std::error::Error>> {
        let f = std::fs::File::open(path)?;
        let cfg: Config = serde_yaml::from_reader(f)?;
        Ok(cfg)
    }
}

impl ObservabilityConfig {
    pub fn to_observability_config(self) -> obs_config::Config {
        obs_config::Config::new()
            .enabled(self.enabled)
            .endpoint(&self.endpoint)
            .attributes(self.attributes)
            .tracer(self.tracer.into())
            .meter(self.meter.into())
    }
}

impl From<Tracer> for obs_config::Tracer {
    fn from(t: Tracer) -> Self {
        obs_config::Tracer::new().enabled(t.enabled)
    }
}

impl From<Meter> for obs_config::Meter {
    fn from(m: Meter) -> Self {
        let mut meter = obs_config::Meter::new().enabled(m.enabled);
        if let Some(d) = m.export_duration {
            meter = meter.export_duration(Duration::new(d.secs, d.nanos));
        }
        if let Some(d) = m.export_timeout_duration {
            meter = meter.export_timeout_duration(Duration::new(d.secs, d.nanos));
        }
        meter
    }
}
