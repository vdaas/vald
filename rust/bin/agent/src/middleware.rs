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

#[derive(Debug, Clone, Default)]
pub struct AccessLogMiddlewareLayer {}

impl<S> Layer<S> for AccessLogMiddlewareLayer {
    type Service = AccessLogMiddleware<S>;

    fn layer(&self, service: S) -> Self::Service {
        AccessLogMiddleware { inner: service }
    }
}

#[derive(Debug, Clone)]
pub struct AccessLogMiddleware<S> {
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
            let response = match result {
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
                    res
                }
                Err(e) => {
                    warn!("{}, {:?}, {:?}", RPC_FAILED_MESSAGE, entity, e);
                    return Err(e);
                }
            };
            Ok(response)
        })
    }
}
