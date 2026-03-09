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
use axum::{Json, Router, http::StatusCode, response::IntoResponse, routing::get};
use serde_json::json;

/// Health check handler
pub async fn liveness() -> impl IntoResponse {
    (
        StatusCode::OK,
        Json(json!({ "status": "ok", "mode": "liveness" })),
    )
}

/// Readiness check handler
pub async fn readiness() -> impl IntoResponse {
    (
        StatusCode::OK,
        Json(json!({ "status": "ok", "mode": "readiness" })),
    )
}

/// Startup check handler
pub async fn startup() -> impl IntoResponse {
    (
        StatusCode::OK,
        Json(json!({ "status": "ok", "mode": "startup" })),
    )
}

/// Create and configure the health check router
pub fn router() -> Router {
    Router::new()
        .route("/liveness", get(liveness))
        .route("/readiness", get(readiness))
        .route("/startup", get(startup))
}
