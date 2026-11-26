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
	if m == method || (len(m) > len(method) && m[len(m)-len(method)-1] == '/' && strings.HasSuffix(m, method)) {
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

// WithQuietErrorLog returns a copy of ctx whose RPC errors the client logs at
// Debug instead of the per-code level, for fan-out callers (e.g. a search
// broadcast) that tolerate per-target failures and log a per-request summary.
// The backoff retry loop's per-attempt error log is demoted as well.
func WithQuietErrorLog(ctx context.Context) context.Context {
	if _, ok := ctx.Value(backoff.QuietErrorLogContextKey()).(searchBroadcastPolicy); ok {
		return ctx
	}
	return backoff.WithQuietErrorLog(ctx)
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

// isNoResourceExhaustedRetry reports whether ctx was marked by WithoutResourceExhaustedRetry
// or WithSearchBroadcastPolicy.
func isNoResourceExhaustedRetry(ctx context.Context) bool {
	if _, ok := ctx.Value(backoff.QuietErrorLogContextKey()).(searchBroadcastPolicy); ok {
		return true
	}
	no, ok := ctx.Value(noResourceExhaustedRetryContextKey).(bool)
	return ok && no
}

type searchBroadcastPolicy uint8

const searchBroadcastPolicyEnabled searchBroadcastPolicy = 1

func (p searchBroadcastPolicy) QuietErrorLog() bool {
	return p == searchBroadcastPolicyEnabled
}

// WithSearchBroadcastPolicy returns a copy of ctx carrying both the quiet error
// log policy and the no-ResourceExhausted-retry policy in a single context value
// allocation.
func WithSearchBroadcastPolicy(ctx context.Context) context.Context {
	return context.WithValue(ctx, backoff.QuietErrorLogContextKey(), searchBroadcastPolicyEnabled)
}
