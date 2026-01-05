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
use anyhow::{Ok, Result};
use opentelemetry::global;
use opentelemetry_otlp::{MetricExporter, SpanExporter, WithExportConfig};
use opentelemetry_sdk::metrics::{PeriodicReader, SdkMeterProvider};
use opentelemetry_sdk::propagation::TraceContextPropagator;
use opentelemetry_sdk::trace::{self, SdkTracerProvider};
use opentelemetry_sdk::Resource;
use url::Url;

use crate::config::Config;

pub const SERVICE_NAME: &str = opentelemetry_semantic_conventions::resource::SERVICE_NAME;

pub trait Observability {
    fn shutdown(&mut self) -> Result<()>;
}

pub struct ObservabilityImpl {
    config: Config,
    meter_provider: Option<SdkMeterProvider>,
    tracer_provider: Option<SdkTracerProvider>,
}

impl ObservabilityImpl {
    pub fn new(cfg: Config) -> Result<ObservabilityImpl, anyhow::Error> {
        let mut obj = ObservabilityImpl {
            config: cfg,
            meter_provider: None,
            tracer_provider: None,
        };

        if !obj.config.enabled {
            return Ok(obj);
        }

        if obj.config.meter.enabled {
            // NOTE: Since the agent implementation does not use views, we will use the simplest implementation for the current phase.
            // If we want flexibility and customization, use SdkMeterProvider::builder.
            let exporter = MetricExporter::builder()
                .with_tonic()
                .with_endpoint(
                    Url::parse(obj.config.endpoint.as_str())?
                        .join("/v1/metrics")?
                        .as_str(),
                )
                .with_timeout(obj.config.meter.export_timeout_duration)
                .build()?;
            let reader = PeriodicReader::builder(exporter)
                .with_interval(obj.config.meter.export_duration)
                .build();
            let provider = SdkMeterProvider::builder()
                .with_reader(reader)
                .with_resource(Resource::from(obj.config()))
                .build();
            obj.meter_provider = Some(provider.clone());
            global::set_meter_provider(provider.clone());
        }

        if obj.config.tracer.enabled {
            let exporter = SpanExporter::builder()
                .with_tonic()
                .with_endpoint(
                    Url::parse(obj.config.endpoint.as_str())?
                        .join("/v1/traces")?
                        .as_str(),
                )
                .build()?;
            let provider = SdkTracerProvider::builder()
                .with_batch_exporter(exporter)
                .with_sampler(trace::Sampler::AlwaysOn)
                .with_resource(Resource::from(obj.config()))
                .with_id_generator(trace::RandomIdGenerator::default())
                .build();
            obj.tracer_provider = Some(provider.clone());
            global::set_text_map_propagator(TraceContextPropagator::new());
            global::set_tracer_provider(provider.clone());
        }
        Ok(obj)
    }

    fn config(&self) -> &Config {
        &self.config
    }
}

impl Observability for ObservabilityImpl {
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

        if self.config.tracer.enabled {
            if let Some(ref provider) = self.tracer_provider {
                provider.force_flush()?;
                provider.shutdown()?;
            }
        }
        Ok(())
    }
}
