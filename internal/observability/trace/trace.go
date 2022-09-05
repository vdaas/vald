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
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	tracerName = "vdaas/vald"
)

var (
	enabled bool

	FromContext = trace.SpanFromContext
)

type Span = trace.Span

type Tracer interface {
	Start(ctx context.Context) error
}

type tracer struct {
}

func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	if !enabled {
		return ctx, nil
	}

	return otel.Tracer(tracerName).Start(ctx, name, opts...)
}

func New(opts ...TraceOption) (Tracer, error) {
	t := &tracer{}

	for _, opt := range append(traceDefaultOpts, opts...) {
		if err := opt(t); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(err))

			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}

	enabled = true

	return t, nil
}

func (t *tracer) Start(ctx context.Context) error {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}
