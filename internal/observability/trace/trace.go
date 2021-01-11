//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

	"go.opencensus.io/trace"
)

var (
	enabled bool

	BoolAttribute    = trace.BoolAttribute
	Float64Attribute = trace.Float64Attribute
	Int64Attribute   = trace.Int64Attribute
	StringAttribute  = trace.StringAttribute

	FromContext = trace.FromContext
)

type Span = trace.Span

type Tracer interface {
	Start(ctx context.Context)
}

type tracer struct {
	samplingRate float64
}

func StartSpan(ctx context.Context, name string, opts ...trace.StartOption) (context.Context, *Span) {
	if !enabled {
		return ctx, nil
	}

	return trace.StartSpan(ctx, name, opts...)
}

func New(opts ...TraceOption) Tracer {
	t := new(tracer)

	for _, opt := range append(traceDefaultOpts, opts...) {
		opt(t)
	}

	enabled = true

	return t
}

func (t *tracer) Start(ctx context.Context) {
	trace.ApplyConfig(
		trace.Config{
			DefaultSampler: trace.ProbabilitySampler(t.samplingRate),
		},
	)
}
