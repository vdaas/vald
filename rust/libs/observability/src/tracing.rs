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

//! Tracing integration module for OpenTelemetry.
//!
//! This module provides integration between the `tracing` crate and OpenTelemetry,
//! allowing spans and events from `tracing` to be exported to OpenTelemetry backends.

use opentelemetry::global;
use opentelemetry::trace::TracerProvider;
use opentelemetry_otlp::{SpanExporter, WithExportConfig};
use opentelemetry_sdk::Resource;
use opentelemetry_sdk::propagation::TraceContextPropagator;
use opentelemetry_sdk::trace::{self, SdkTracerProvider};
use tracing_opentelemetry::OpenTelemetryLayer;
use tracing_subscriber::EnvFilter;
use tracing_subscriber::layer::SubscriberExt;
use tracing_subscriber::util::SubscriberInitExt;
use url::Url;

use crate::config::Config;
use crate::error::{ObservabilityError, Result};

/// Configuration for tracing initialization.
#[derive(Clone, Debug)]
pub struct TracingConfig {
    /// Enable tracing output to stdout/stderr.
    pub enable_stdout: bool,
    /// Enable JSON format for stdout output.
    pub enable_json: bool,
    /// Enable OpenTelemetry export.
    pub enable_otel: bool,
    /// Log level filter (e.g., "info", "debug", "trace").
    pub level: String,
    /// Service name for tracing.
    pub service_name: String,
}

impl Default for TracingConfig {
    fn default() -> Self {
        Self {
            enable_stdout: true,
            enable_json: false,
            enable_otel: false,
            level: "info".to_string(),
            service_name: "vald-agent".to_string(),
        }
    }
}

impl TracingConfig {
    /// Creates a tracing configuration with defaults.
    pub fn new() -> Self {
        Self::default()
    }

    /// Enables or disables stdout/stderr output.
    pub fn enable_stdout(mut self, enable: bool) -> Self {
        self.enable_stdout = enable;
        self
    }

    /// Enables or disables JSON output formatting.
    pub fn enable_json(mut self, enable: bool) -> Self {
        self.enable_json = enable;
        self
    }

    /// Enables or disables OpenTelemetry export.
    pub fn enable_otel(mut self, enable: bool) -> Self {
        self.enable_otel = enable;
        self
    }

    /// Sets the log level filter.
    pub fn level(mut self, level: &str) -> Self {
        self.level = level.to_string();
        self
    }

    /// Sets the service name used in tracing.
    pub fn service_name(mut self, name: &str) -> Self {
        self.service_name = name.to_string();
        self
    }
}

/// Initialize tracing with the given configuration.
///
/// This sets up a tracing subscriber with optional layers:
/// - Stdout/stderr output (with optional JSON formatting)
/// - OpenTelemetry export (if otel_config is provided)
///
/// # Arguments
/// * `tracing_config` - Configuration for tracing behavior
/// * `otel_config` - Optional OpenTelemetry configuration for exporting traces
///
/// # Returns
/// * `Ok(Option<SdkTracerProvider>)` - The tracer provider if OpenTelemetry is enabled
pub fn init_tracing(
    tracing_config: &TracingConfig,
    otel_config: Option<&Config>,
) -> Result<Option<SdkTracerProvider>> {
    let env_filter =
        EnvFilter::try_from_default_env().unwrap_or_else(|_| EnvFilter::new(&tracing_config.level));

    // Initialize OpenTelemetry tracer if enabled
    let tracer_provider = if tracing_config.enable_otel {
        if let Some(cfg) = otel_config {
            if cfg.enabled && cfg.tracer.enabled {
                Some(init_otel_tracer(cfg)?)
            } else {
                None
            }
        } else {
            None
        }
    } else {
        None
    };

    // Build subscriber based on configuration
    // Note: We use separate match branches to avoid complex type combinations
    match (
        tracing_config.enable_stdout,
        tracing_config.enable_json,
        &tracer_provider,
    ) {
        // stdout + json + otel
        (true, true, Some(provider)) => {
            let tracer = provider.tracer(tracing_config.service_name.clone());
            tracing_subscriber::registry()
                .with(env_filter)
                .with(tracing_subscriber::fmt::layer().json())
                .with(OpenTelemetryLayer::new(tracer))
                .try_init()
                .map_err(|e| ObservabilityError::TracingInit(Box::new(e)))?
        }
        // stdout + json (no otel)
        (true, true, None) => {
            tracing_subscriber::registry()
                .with(env_filter)
                .with(tracing_subscriber::fmt::layer().json())
                .try_init()
                .map_err(|e| ObservabilityError::TracingInit(Box::new(e)))?
        }
        // stdout + text + otel
        (true, false, Some(provider)) => {
            let tracer = provider.tracer(tracing_config.service_name.clone());
            tracing_subscriber::registry()
                .with(env_filter)
                .with(tracing_subscriber::fmt::layer())
                .with(OpenTelemetryLayer::new(tracer))
                .try_init()
                .map_err(|e| ObservabilityError::TracingInit(Box::new(e)))?
        }
        // stdout + text (no otel)
        (true, false, None) => {
            tracing_subscriber::registry()
                .with(env_filter)
                .with(tracing_subscriber::fmt::layer())
                .try_init()
                .map_err(|e| ObservabilityError::TracingInit(Box::new(e)))?
        }
        // no stdout + otel only
        (false, _, Some(provider)) => {
            let tracer = provider.tracer(tracing_config.service_name.clone());
            tracing_subscriber::registry()
                .with(env_filter)
                .with(OpenTelemetryLayer::new(tracer))
                .try_init()
                .map_err(|e| ObservabilityError::TracingInit(Box::new(e)))?
        }
        // no output at all
        (false, _, None) => {
            tracing_subscriber::registry()
                .with(env_filter)
                .try_init()
                .map_err(|e| ObservabilityError::TracingInit(Box::new(e)))?
        }
    }

    if let Some(provider) = &tracer_provider {
        global::set_text_map_propagator(TraceContextPropagator::new());
        global::set_tracer_provider(provider.clone());
    }
    
    Ok(tracer_provider)
}

/// Initialize OpenTelemetry tracer provider.
fn init_otel_tracer(cfg: &Config) -> Result<SdkTracerProvider> {
    let exporter = SpanExporter::builder()
        .with_tonic()
        .with_endpoint(
            Url::parse(cfg.endpoint.as_str())?
                .join("/v1/traces")?
                .as_str(),
        )
        .build()?;

    let provider = SdkTracerProvider::builder()
        .with_batch_exporter(exporter)
        .with_sampler(trace::Sampler::AlwaysOn)
        .with_resource(Resource::from(cfg))
        .with_id_generator(trace::RandomIdGenerator::default())
        .build();

    Ok(provider)
}

/// Shutdown tracing and flush any pending spans.
pub fn shutdown_tracing(provider: Option<SdkTracerProvider>) -> Result<()> {
    if let Some(provider) = provider {
        provider.force_flush()?;
        provider.shutdown()?;
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_tracing_config_default() {
        let config = TracingConfig::default();
        assert!(config.enable_stdout);
        assert!(!config.enable_json);
        assert!(!config.enable_otel);
        assert_eq!(config.level, "info");
    }

    #[test]
    fn test_tracing_config_builder() {
        let config = TracingConfig::new()
            .enable_stdout(false)
            .enable_json(true)
            .enable_otel(true)
            .level("debug")
            .service_name("test-service");

        assert!(!config.enable_stdout);
        assert!(config.enable_json);
        assert!(config.enable_otel);
        assert_eq!(config.level, "debug");
        assert_eq!(config.service_name, "test-service");
    }
}
