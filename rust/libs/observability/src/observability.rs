//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
use std::sync::Arc;

use anyhow::{Ok, Result};
use opentelemetry::global::{self, shutdown_tracer_provider};
use opentelemetry_otlp::WithExportConfig;
use opentelemetry_sdk::metrics::SdkMeterProvider;
use opentelemetry_sdk::propagation::TraceContextPropagator;
use opentelemetry_sdk::trace::TracerProvider;
use opentelemetry_sdk::Resource;

use crate::config::Config;

pub trait Observability {
    fn build(&mut self) -> Result<()>;
    fn shutdown(&mut self) -> Result<()>;
}

pub struct ObservabilityImpl {
    config: Arc<Config>,
    meter_provider: Option<SdkMeterProvider>,
    tracer_provider: Option<TracerProvider>,
}

impl ObservabilityImpl {
    fn new(cfg: Arc<Config>) -> ObservabilityImpl {
        ObservabilityImpl {
            config: cfg,
            meter_provider: None,
            tracer_provider: None,
        }
    }

    fn config(&self) -> &Config {
        &self.config
    }
}

impl Observability for ObservabilityImpl {
    fn build(&mut self) -> Result<()> {
        if !self.config.enabled {
            return Ok(());
        }

        if self.config.meter.enabled {
            let provider = opentelemetry_otlp::new_pipeline()
                .metrics(opentelemetry_sdk::runtime::Tokio)
                .with_period(self.config.meter.export_duration)
                .with_resource(Resource::from(self.config()))
                .with_exporter(
                    opentelemetry_otlp::new_exporter()
                        .http()
                        .with_endpoint(self.config.meter.endpoint.as_str()),
                )
                .build()?;
            self.meter_provider = Some(provider.clone());
            global::set_meter_provider(provider.clone());
        }

        if self.config.tracer.enabled {
            let tracer = opentelemetry_otlp::new_pipeline()
                .tracing()
                .with_exporter(
                    opentelemetry_otlp::new_exporter()
                        .http()
                        .with_endpoint(self.config.tracer.endpoint.as_str()),
                )
                .with_trace_config(
                    opentelemetry_sdk::trace::config()
                        .with_sampler(opentelemetry_sdk::trace::Sampler::AlwaysOn)
                        .with_resource(Resource::from(self.config()))
                        .with_id_generator(opentelemetry_sdk::trace::RandomIdGenerator::default()),
                )
                .install_batch(opentelemetry_sdk::runtime::Tokio)?;
            global::set_text_map_propagator(TraceContextPropagator::new());
            global::set_tracer_provider(tracer.provider().unwrap());
        }
        Ok(())
    }

    fn shutdown(&mut self) -> Result<()> {
        if !self.config.enabled {
            return Ok(());
        }

        if self.config.meter.enabled {
            if let Some(ref provider) = self.meter_provider {
                provider.force_flush()?;
                provider.shutdown()?;
            }
        }

        if self.config.meter.enabled {
            if let Some(ref provider) = self.tracer_provider {
                for result in provider.force_flush() {
                    result?;
                }
                shutdown_tracer_provider();
            }
        }
        Ok(())
    }
}
