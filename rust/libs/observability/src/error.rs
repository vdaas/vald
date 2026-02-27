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

//! Error types for the observability crate.

use thiserror::Error;

/// Result type alias using ObservabilityError.
pub type Result<T> = std::result::Result<T, ObservabilityError>;

/// Error types for observability operations.
#[derive(Error, Debug)]
pub enum ObservabilityError {
    /// OpenTelemetry trace error.
    #[error("trace error: {0}")]
    Trace(#[from] opentelemetry_sdk::trace::TraceError),

    /// OpenTelemetry exporter build error.
    #[error("exporter build error: {0}")]
    ExporterBuild(#[from] opentelemetry_otlp::ExporterBuildError),

    /// OpenTelemetry SDK error.
    #[error("OTel SDK error: {0}")]
    OTelSdk(#[from] opentelemetry_sdk::error::OTelSdkError),

    /// Tracing subscriber initialization error.
    #[error("failed to initialize tracing subscriber: {0}")]
    TracingInit(Box<dyn std::error::Error + Send + Sync>),

    /// URL parsing error.
    #[error("invalid URL: {0}")]
    Url(#[from] url::ParseError),

    /// Generic error with string message.
    #[error("{0}")]
    Other(String),
}
