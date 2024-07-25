use std::sync::Arc;

use anyhow::{Ok, Result};
use opentelemetry::global::shutdown_tracer_provider;
use opentelemetry_sdk::metrics::SdkMeterProvider;
use opentelemetry_sdk::trace::TracerProvider;

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
