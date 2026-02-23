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
