// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package grpc

import (
	"context"

	"github.com/vdaas/vald/internal/backoff"
	"github.com/vdaas/vald/internal/strings"
)

type contextKey string

// GRPCMethodContextKey represents a context key for gRPC method.
// This is exported only for testing.
const GRPCMethodContextKey contextKey = "grpc_method"

// WrapGRPCMethod returns a copy of parent in which the method associated with key (grpcMethodContextKey).
func WrapGRPCMethod(ctx context.Context, method string) context.Context {
	m := FromGRPCMethod(ctx)
	if m == "" {
		return context.WithValue(ctx, GRPCMethodContextKey, method)
	}
	if strings.HasSuffix(m, method) {
		return ctx
	}
	return context.WithValue(ctx, GRPCMethodContextKey, m+"/"+method)
}

// WithGRPCMethod returns a copy of parent in which the method associated with key (grpcMethodContextKey).
func WithGRPCMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, GRPCMethodContextKey, method)
}

// FromGRPCMethod returns the value associated with this context for key (grpcMethodContextKey).
func FromGRPCMethod(ctx context.Context) string {
	if v := ctx.Value(GRPCMethodContextKey); v != nil {
		if method, ok := v.(string); ok {
			return method
		}
	}
	return ""
}

// quietErrorLogContextKey marks calls whose per-call error log is demoted to
// Debug because the caller aggregates and reports the failures itself.
const quietErrorLogContextKey contextKey = "quiet_error_log"

// WithQuietErrorLog returns a copy of ctx whose RPC errors the client logs at
// Debug instead of the per-code level, for fan-out callers (e.g. a search
// broadcast) that tolerate per-target failures and log a per-request summary.
// The backoff retry loop's per-attempt error log is demoted as well.
func WithQuietErrorLog(ctx context.Context) context.Context {
	return context.WithValue(backoff.WithQuietErrorLog(ctx), quietErrorLogContextKey, true)
}

// isQuietErrorLog reports whether ctx was marked by WithQuietErrorLog.
func isQuietErrorLog(ctx context.Context) bool {
	quiet, ok := ctx.Value(quietErrorLogContextKey).(bool)
	return ok && quiet
}

// noResourceExhaustedRetryContextKey marks calls whose ResourceExhausted
// errors are returned at once instead of being retried with backoff.
const noResourceExhaustedRetryContextKey contextKey = "no_resource_exhausted_retry"

// WithoutResourceExhaustedRetry returns a copy of ctx whose RPCs are not
// retried on ResourceExhausted, for fan-out reads (e.g. a search broadcast)
// whose per-target failures are covered by the other replicas: retrying an
// overloaded agent inside the request deadline only adds load and latency.
// Write paths keep the default backoff-and-wait policy.
func WithoutResourceExhaustedRetry(ctx context.Context) context.Context {
	return context.WithValue(ctx, noResourceExhaustedRetryContextKey, true)
}

// isNoResourceExhaustedRetry reports whether ctx was marked by WithoutResourceExhaustedRetry.
func isNoResourceExhaustedRetry(ctx context.Context) bool {
	no, ok := ctx.Value(noResourceExhaustedRetryContextKey).(bool)
	return ok && no
}
