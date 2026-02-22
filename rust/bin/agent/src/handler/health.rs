use axum::{
    Json,
    extract::State,
    http::StatusCode,
    response::IntoResponse,
    routing::get,
    Router,
};
use serde_json::json;
use std::sync::Arc;

/// Health check handler
pub async fn liveness() -> impl IntoResponse {
    (StatusCode::OK, Json(json!({ "status": "ok", "mode": "liveness" })))
}

pub async fn readiness() -> impl IntoResponse {
    (StatusCode::OK, Json(json!({ "status": "ok", "mode": "readiness" })))
}

pub async fn startup() -> impl IntoResponse {
    (StatusCode::OK, Json(json!({ "status": "ok", "mode": "startup" })))
}

pub fn router() -> Router {
    Router::new()
        .route("/liveness", get(liveness))
        .route("/readiness", get(readiness))
        .route("/startup", get(startup))
}
