//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"sync"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	sdk "go.opentelemetry.io/otel/sdk/trace"
)

var (
	enabled bool

	FromContext = trace.SpanFromContext

	instance *tracer
	once     sync.Once
)

type Span = trace.Span

type Tracer interface {
	Start(ctx context.Context)
	StartSpan(ctx context.Context, name string, opts ...trace.StartOption) (context.Context, *Span)
}

type tracer struct {
	name         string
	tracer       trace.Tracer
	samplingRate float64
}

func StartSpan(ctx context.Context, name string, opts ...trace.StartOption) (context.Context, *Span) {
	if !enabled || instance == nil {
		return ctx, nil
	}

	return tracer.StartSpan(ctx, name, opts...)
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
	global.SetTraceProvider(
		sdk.NewProvider(
			sdk.WithConfig(
				sdk.Config{
					DefaultSampler: sdk.ProbabilitySampler(t.samplingRate),
				},
			),
		),
	)

	t.tracer = global.Tracer(t.name)

	once.Do(func() {
		instance = t
	})
}

func (t *tracer) StartSpan(ctx context.Context, name string, opts ...trace.StartOption) (context.Context, *Span) {
	return t.tracer.Start(ctx, name, opts...)
}
