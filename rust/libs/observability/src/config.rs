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
use std::time::Duration;

use opentelemetry::KeyValue;
use opentelemetry_sdk::{self, Resource};

/// OpenTelemetry configuration for tracing and metrics.
#[derive(Clone, Debug)]
pub struct Config {
    /// Enables OpenTelemetry export.
    pub enabled: bool,
    /// OTLP endpoint for trace/metric export.
    pub endpoint: String,
    /// Resource attributes applied to all telemetry.
    pub attributes: HashMap<String, String>,
    /// Tracing configuration.
    pub tracer: Tracer,
    /// Metrics configuration.
    pub meter: Meter,
}

/// Tracing configuration settings.
#[derive(Clone, Debug, Default)]
pub struct Tracer {
    /// Enables tracing export.
    pub enabled: bool,
}

/// Metrics configuration settings.
#[derive(Clone, Debug)]
pub struct Meter {
    /// Enables metrics export.
    pub enabled: bool,
    /// Metric export interval.
    pub export_duration: Duration,
    /// Metric export timeout.
    pub export_timeout_duration: Duration,
}

impl Config {
    /// Creates a configuration with default values.
    pub fn new() -> Self {
        Self::default()
    }

    /// Sets whether OpenTelemetry export is enabled.
    pub fn enabled(mut self, enabled: bool) -> Self {
        self.enabled = enabled;
        self
    }

    /// Sets the OTLP endpoint.
    pub fn endpoint(mut self, endpoint: &str) -> Self {
        self.endpoint = endpoint.to_string();
        self
    }

    /// Sets resource attributes for exporters.
    pub fn attributes(mut self, attrs: HashMap<String, String>) -> Self {
        self.attributes = attrs;
        self
    }

    /// Adds a single resource attribute.
    pub fn attribute(mut self, key: &str, value: &str) -> Self {
        self.attributes.insert(key.to_string(), value.to_string());
        self
    }

    /// Sets the tracing configuration.
    pub fn tracer(mut self, cfg: Tracer) -> Self {
        self.tracer = cfg;
        self
    }

    /// Sets the metrics configuration.
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
        Resource::builder().with_attributes(key_values).build()
    }
}

impl Tracer {
    /// Creates a tracing configuration with default values.
    pub fn new() -> Self {
        Tracer::default()
    }

    /// Enables or disables tracing export.
    pub fn enabled(mut self, enabled: bool) -> Self {
        self.enabled = enabled;
        self
    }
}

impl Meter {
    /// Creates a metrics configuration with default values.
    pub fn new() -> Self {
        Meter::default()
    }

    /// Enables or disables metrics export.
    pub fn enabled(mut self, enabled: bool) -> Self {
        self.enabled = enabled;
        self
    }

    /// Sets the metrics export interval.
    pub fn export_duration(mut self, dur: Duration) -> Self {
        self.export_duration = dur;
        self
    }

    /// Sets the metrics export timeout.
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
