// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package backoff

import "context"

type contextKey string

const (
	backoffNameContextKey      contextKey = "backoff_name"
	backoffRunningContextKey   contextKey = "backoff_running"
	backoffRunningContextValue string     = string(backoffRunningContextKey)
)

// WithBackoffName returns a copy of parent in which the method associated with key (backoffNameContextKey).
func WithBackoffName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, backoffNameContextKey, name)
}

// FromBackoffName returns the value associated with this context for key (backoffNameContextKey).
func FromBackoffName(ctx context.Context) string {
	if val := ctx.Value(backoffNameContextKey); val != nil {
		if name, ok := val.(string); ok {
			return name
		}
	}
	return ""
}

// isRunning returns context and is running or not for preventing duplicated backoff execution
func isRunning(ctx context.Context) (context.Context, bool) {
	if val := ctx.Value(backoffRunningContextKey); val != nil {
		if v, ok := val.(string); ok && v == backoffRunningContextValue {
			return ctx, true
		}
	}
	return context.WithValue(ctx, backoffRunningContextKey, backoffRunningContextValue), false
}
