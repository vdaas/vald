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

use log::{debug, warn};
use std::{
    pin::Pin,
    task::{Context, Poll},
    time::{SystemTime, UNIX_EPOCH},
};
use tower::{Layer, Service};

const GRPC_KIND_UNARY: &str = "unary";
const GRPC_KIND_STREAM: &str = "stream";
const RPC_COMPLETED_MESSAGE: &str = "rpc completed";
const RPC_FAILED_MESSAGE: &str = "rpc failed";
const VALD_ORG: &str = "vald.vdaas.org";
const LATENCY_METRICS_NAME: &str = "server_latency";
const COMPLETED_RPCS_METRICS_NAME: &str = "server_completed_rpcs";
const MILLISECONDS: &str = "ms";
const GRPCMETHOD_KEY_NAME: &str = "grpc_server_method";
const GRPCSTATUS: &str = "grpc_server_status";

#[derive(Debug)]
struct AccessLogEntity {
    grpc: AccessLogGRPCEntity,
    start_time: i64,
    end_time: i64,
    latency: i64,
    trace_id: String,
}

#[derive(Debug)]
struct AccessLogGRPCEntity {
    kind: String,
    service: String,
    method: String,
}

/// Layer that wraps services with access logging middleware.
#[derive(Debug, Clone, Default)]
pub struct AccessLogMiddlewareLayer {}

/// Layer that wraps services with metrics recording middleware.
#[derive(Debug, Clone, Default)]
pub struct MetricMiddlewareLayer {}

impl<S> Layer<S> for AccessLogMiddlewareLayer {
    type Service = AccessLogMiddleware<S>;

    fn layer(&self, service: S) -> Self::Service {
        AccessLogMiddleware { inner: service }
    }
}

impl<S> Layer<S> for MetricMiddlewareLayer {
    type Service = MetricMiddleware<S>;

    fn layer(&self, service: S) -> Self::Service {
        MetricMiddleware { inner: service }
    }
}

/// Service wrapper that logs access information for each request.
#[derive(Debug, Clone)]
pub struct AccessLogMiddleware<S> {
    inner: S,
}

/// Service wrapper that records metrics for each request.
#[derive(Debug, Clone)]
pub struct MetricMiddleware<S> {
    inner: S,
}

type BoxFuture<'a, T> = Pin<Box<dyn std::future::Future<Output = T> + Send + 'a>>;

impl<S, ReqBody, ResBody> Service<http::Request<ReqBody>> for AccessLogMiddleware<S>
where
    S: Service<http::Request<ReqBody>, Response = http::Response<ResBody>> + Clone + Send + 'static,
    S::Future: Send + 'static,
    S::Error: std::fmt::Debug,
    ReqBody: Send + 'static,
{
    type Response = S::Response;
    type Error = S::Error;
    type Future = BoxFuture<'static, Result<Self::Response, Self::Error>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }

    fn call(&mut self, req: http::Request<ReqBody>) -> Self::Future {
        // See: https://docs.rs/tower/latest/tower/trait.Service.html#be-careful-when-cloning-inner-services
        let clone = self.inner.clone();
        let mut inner = std::mem::replace(&mut self.inner, clone);

        let path = req.uri().path().to_string().clone();
        Box::pin(async move {
            // Do extra async work here...
            let (service, method) = if let Some(pos) = path.rfind('/') {
                let service_part = &path[..pos];
                let method_part = &path[(pos + 1)..];
                if let Some(service_pos) = service_part.rfind('/') {
                    (&service_part[(service_pos + 1)..], method_part)
                } else {
                    (service_part.trim_start_matches('/'), method_part)
                }
            } else {
                ("", path.trim_start_matches('/'))
            };
            let kind = if method.to_lowercase().contains("stream") {
                GRPC_KIND_STREAM
            } else {
                GRPC_KIND_UNARY
            };

            let start = SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap_or_default();
            let result = inner.call(req).await;
            let end = SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap_or_default();
            let start_nanos = start.as_secs() as i64 * 1_000_000_000 + start.subsec_nanos() as i64;
            let end_nanos = end.as_secs() as i64 * 1_000_000_000 + end.subsec_nanos() as i64;
            let entity = AccessLogEntity {
                grpc: AccessLogGRPCEntity {
                    kind: kind.to_string(),
                    service: service.to_string(),
                    method: method.to_string(),
                },
                start_time: start_nanos,
                end_time: end_nanos,
                latency: end_nanos - start_nanos,
                trace_id: "".to_string(),
            };
            match result {
                Ok(res) => {
                    let status = res.headers().get("grpc-status");
                    if status.is_none() {
                        debug!("{}, {:?}", RPC_COMPLETED_MESSAGE, entity);
                    } else {
                        let message = res
                            .headers()
                            .get("grpc-message")
                            .and_then(|v| v.to_str().ok())
                            .unwrap_or("internal error");
                        warn!("{}, {:?}, {:?}", RPC_FAILED_MESSAGE, entity, message);
                    }
                    return Ok(res);
                }
                Err(e) => {
                    warn!("{}, {:?}, {:?}", RPC_FAILED_MESSAGE, entity, e);
                    Err(e)
                }
            }
        })
    }
}

fn grpc_status_to_string(grpc_status: &str) -> &str {
    match grpc_status {
        "0" => "OK",
        "1" => "Canceled",
        "2" => "Unknown",
        "3" => "InvalidArgument",
        "4" => "DeadlineExceeded",
        "5" => "NotFound",
        "6" => "AlreadyExists",
        "7" => "PermissionDenied",
        "8" => "ResourceExhausted",
        "9" => "FailedPrecondition",
        "10" => "Aborted",
        "11" => "OutOfRange",
        "12" => "Unimplemented",
        "13" => "Internal",
        "14" => "Unavailable",
        "15" => "DataLoss",
        "16" => "Unauthenticated",
        _ => "InvalidStatus",
    }
}

impl<S, ReqBody, ResBody> Service<http::Request<ReqBody>> for MetricMiddleware<S>
where
    S: Service<http::Request<ReqBody>, Response = http::Response<ResBody>> + Clone + Send + 'static,
    S::Future: Send + 'static,
    S::Error: std::fmt::Debug,
    ReqBody: Send + 'static,
{
    type Response = S::Response;
    type Error = S::Error;
    type Future = BoxFuture<'static, Result<Self::Response, Self::Error>>;

    fn poll_ready(&mut self, cx: &mut Context<'_>) -> Poll<Result<(), Self::Error>> {
        self.inner.poll_ready(cx)
    }

    fn call(&mut self, req: http::Request<ReqBody>) -> Self::Future {
        let clone = self.inner.clone();
        let mut inner = std::mem::replace(&mut self.inner, clone);

        let path = req.uri().path().to_string().clone();
        Box::pin(async move {
            let meter = opentelemetry::global::meter(VALD_ORG);
            let latency_histogram = meter
                .f64_histogram(LATENCY_METRICS_NAME)
                .with_description("Server latency in milliseconds, by method")
                .with_unit(MILLISECONDS)
                .build();
            let completed_rpc_cnt = meter
                .u64_counter(COMPLETED_RPCS_METRICS_NAME)
                .with_description("Count of RPCs by method and status")
                .with_unit(MILLISECONDS)
                .build();

            let start = SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap_or_default();
            let result = inner.call(req).await;
            let end = SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap_or_default();
            let start_nanos =
                (start.as_secs() as i64 * 1_000_000_000 + start.subsec_nanos() as i64) as f64;
            let end_nanos =
                (end.as_secs() as i64 * 1_000_000_000 + end.subsec_nanos() as i64) as f64;
            match result {
                Ok(res) => {
                    let status = res.headers().get("grpc-status");
                    let grpc_status_str = status.and_then(|v| v.to_str().ok()).unwrap_or("0");
                    let code = grpc_status_to_string(grpc_status_str).to_string();
                    let attributes = [
                        opentelemetry::KeyValue::new(GRPCMETHOD_KEY_NAME, path),
                        opentelemetry::KeyValue::new(GRPCSTATUS, code),
                    ];
                    latency_histogram
                        .record((end_nanos - start_nanos) / 1_000_000_f64, &attributes);
                    completed_rpc_cnt.add(1, &attributes);
                    return Ok(res);
                }
                Err(e) => {
                    let attributes = [
                        opentelemetry::KeyValue::new(GRPCMETHOD_KEY_NAME, path),
                        opentelemetry::KeyValue::new(
                            GRPCSTATUS,
                            grpc_status_to_string("13").to_string(),
                        ),
                    ];
                    latency_histogram
                        .record((end_nanos - start_nanos) / 1_000_000_f64, &attributes);
                    completed_rpc_cnt.add(1, &attributes);
                    Err(e)
                }
            }
        })
    }
}
