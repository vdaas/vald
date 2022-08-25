//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package trace provides trace functions.
package trace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	enabled bool

	BoolAttribute = func(key string, val bool) attribute.KeyValue {
		return attribute.Key(key).Bool(val)
	}
	Float64Attribute = func(key string, val float64) attribute.KeyValue {
		return attribute.Key(key).Float64(val)
	}
	Int64Attribute = func(key string, val int64) attribute.KeyValue {
		return attribute.Key(key).Int64(val)
	}
	StringAttribute = func(key, val string) attribute.KeyValue {
		return attribute.Key(key).String(val)
	}

	FromContext = trace.SpanFromContext
)

type Span = trace.Span

type Tracer interface {
	Start(ctx context.Context)
}

type tracer struct {
	samplingRate float64
}

func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	if !enabled {
		return ctx, nil
	}

	return otel.Tracer("component-main").Start(ctx, name, opts...)
}

func New(opts ...TraceOption) Tracer {
	t := new(tracer)

	for _, opt := range append(traceDefaultOpts, opts...) {
		opt(t)
	}

	enabled = true

	tr := otel.Tracer("component-main")

	_ = tr
	return t
}

func (t *tracer) Start(ctx context.Context) {
	// trace.ApplyConfig(
	// 	trace.Config{
	// 		DefaultSampler: trace.ProbabilitySampler(t.samplingRate),
	// 	},
	// )
}
