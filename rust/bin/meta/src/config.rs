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

use std::collections::HashMap;
use std::error::Error;
use std::fs::File;
use std::time::Duration;

use serde::Deserialize;

#[derive(Debug, Clone, Deserialize, Default)]
pub struct DurationConfig {
    #[serde(default)]
    pub secs: u64,
    #[serde(default)]
    pub nanos: u32,
}

impl DurationConfig {
    pub fn to_duration(&self) -> Result<Duration, Box<dyn Error>> {
        if self.nanos >= 1_000_000_000 {
            return Err("nanos must be less than 1e9".into());
        }
        Ok(Duration::new(self.secs, self.nanos))
    }
}

#[derive(Debug, Clone, Deserialize)]
pub struct TracerConfig {
    #[serde(default)]
    pub enabled: bool,
}

impl Default for TracerConfig {
    fn default() -> Self {
        Self { enabled: false }
    }
}

#[derive(Debug, Clone, Deserialize)]
pub struct MeterConfig {
    #[serde(default)]
    pub enabled: bool,
    #[serde(default)]
    pub export_duration: Option<DurationConfig>,
    #[serde(default)]
    pub export_timeout_duration: Option<DurationConfig>,
}

impl Default for MeterConfig {
    fn default() -> Self {
        Self {
            enabled: false,
            export_duration: None,
            export_timeout_duration: None,
        }
    }
}

#[derive(Debug, Clone, Deserialize)]
pub struct ObservabilityConfig {
    #[serde(default)]
    pub enabled: bool,
    #[serde(default)]
    pub endpoint: String,
    #[serde(default)]
    pub attributes: HashMap<String, String>,
    #[serde(default)]
    pub tracer: TracerConfig,
    #[serde(default)]
    pub meter: MeterConfig,
}

impl Default for ObservabilityConfig {
    fn default() -> Self {
        Self {
            enabled: false,
            endpoint: "".to_string(),
            attributes: HashMap::new(),
            tracer: TracerConfig::default(),
            meter: MeterConfig::default(),
        }
    }
}

#[derive(Debug, Clone, Deserialize)]
pub struct Config {
    #[serde(default)]
    pub observability: ObservabilityConfig,
    #[serde(default = "default_server_addr")]
    pub server_addr: String,
    #[serde(default = "default_database_path")]
    pub database_path: String,
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
    pub fn load(path: Option<&str>) -> Result<Self, Box<dyn Error>> {
        match path {
            Some(p) => {
                let f = File::open(p)?;
                let cfg: Config = serde_yaml::from_reader(f)?;
                Ok(cfg)
            }
            None => Ok(Config::default()),
        }
    }
}

// Conversion to internal observability types
impl ObservabilityConfig {
    pub fn to_observability_config(&self) -> Result<observability::config::Config, Box<dyn Error>> {
        let mut cfg = observability::config::Config::new()
            .enabled(self.enabled)
            .endpoint(&self.endpoint)
            .attributes(self.attributes.clone());

        let mut tracer = observability::config::Tracer::new();
        tracer = tracer.enabled(self.tracer.enabled);
        cfg = cfg.tracer(tracer);

        let mut meter = observability::config::Meter::new();
        meter = meter.enabled(self.meter.enabled);
        if let Some(d) = &self.meter.export_duration {
            meter = meter.export_duration(d.to_duration()?);
        }
        if let Some(d) = &self.meter.export_timeout_duration {
            meter = meter.export_timeout_duration(d.to_duration()?);
        }
        cfg = cfg.meter(meter);

        Ok(cfg)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::io::Write;

    #[test]
    fn test_duration_config() {
        let d = DurationConfig { secs: 1, nanos: 0 };
        assert_eq!(d.to_duration().unwrap(), Duration::new(1, 0));

        let d = DurationConfig {
            secs: 1,
            nanos: 1_000_000_000,
        };
        assert!(d.to_duration().is_err());
    }

    #[test]
    fn test_load_default() {
        let cfg = Config::load(None).unwrap();
        assert_eq!(cfg.server_addr, default_server_addr());
        assert_eq!(cfg.database_path, default_database_path());
        assert!(!cfg.observability.enabled);
    }

    #[test]
    fn test_load_from_file() {
        let mut file = tempfile::NamedTempFile::new().unwrap();
        let yaml = r#"
server_addr: "127.0.0.1:9000"
database_path: "/var/lib/meta"
observability:
  enabled: true
  endpoint: "http://localhost:4317"
"#;
        write!(file, "{}", yaml).unwrap();
        // Keep the file alive
        let path = file.path().to_str().unwrap().to_string();

        let cfg = Config::load(Some(&path)).unwrap();
        assert_eq!(cfg.server_addr, "127.0.0.1:9000");
        assert_eq!(cfg.database_path, "/var/lib/meta");
        assert!(cfg.observability.enabled);
        assert_eq!(cfg.observability.endpoint, "http://localhost:4317");
    }
}
