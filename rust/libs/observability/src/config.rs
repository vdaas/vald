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

use opentelemetry::KeyValue;
use opentelemetry_sdk::{self, Resource};

#[derive(Clone, Debug)]
pub struct Config {
    pub enabled: bool,
    pub endpoint: String,
    pub attributes: HashMap<String, String>,
    pub tracer: Tracer,
    pub meter: Meter,
}

#[derive(Clone, Debug, Default)]
pub struct Tracer {
    pub enabled: bool,
}

#[derive(Clone, Debug)]
pub struct Meter {
    pub enabled: bool,
    pub export_duration: Duration,
    pub export_timeout_duration: Duration,
}

impl Config {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn enabled(mut self, enabled: bool) -> Self {
        self.enabled = enabled;
        self
    }

    pub fn endpoint(mut self, endpoint: &str) -> Self {
        self.endpoint = endpoint.to_string();
        self
    }

    pub fn attributes(mut self, attrs: HashMap<String, String>) -> Self {
        self.attributes = attrs;
        self
    }

    pub fn attribute(mut self, key: &str, value: &str) -> Self {
        self.attributes.insert(key.to_string(), value.to_string());
        self
    }

    pub fn tracer(mut self, cfg: Tracer) -> Self {
        self.tracer = cfg;
        self
    }

    pub fn meter(mut self, cfg: Meter) -> Self {
        self.meter = cfg;
        self
    }
}

impl Default for Config {
    fn default() -> Self {
        Self {
            enabled: false,
            endpoint: "".to_string(),
            attributes: HashMap::new(),
            tracer: Tracer::default(),
            meter: Meter::default(),
        }
    }
}

impl From<&Config> for Resource {
    fn from(value: &Config) -> Self {
        let key_values: Vec<KeyValue> = value
            .attributes
            .iter()
            .map(|(key, val)| KeyValue::new(key.clone(), val.clone()))
            .collect();
        Resource::new(key_values)
    }
}

impl Tracer {
    pub fn new() -> Self {
        Tracer::default()
    }

    pub fn enabled(mut self, enabled: bool) -> Self {
        self.enabled = enabled;
        self
    }
}

impl Meter {
    pub fn new() -> Self {
        Meter::default()
    }

    pub fn enabled(mut self, enabled: bool) -> Self {
        self.enabled = enabled;
        self
    }

    pub fn export_duration(mut self, dur: Duration) -> Self {
        self.export_duration = dur;
        self
    }

    pub fn export_timeout_duration(mut self, dur: Duration) -> Self {
        self.export_timeout_duration = dur;
        self
    }
}

impl Default for Meter {
    fn default() -> Self {
        Self {
            enabled: false,
            export_duration: Duration::from_secs(1),
            export_timeout_duration: Duration::from_secs(5),
        }
    }
}
