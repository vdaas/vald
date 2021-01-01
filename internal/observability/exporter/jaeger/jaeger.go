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

// Package jaeger provides a jaeger exporter.
package jaeger

import (
	"context"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/vdaas/vald/internal/observability/exporter"
	"go.opencensus.io/trace"
)

type jaegerOptions = jaeger.Options

type Jaeger interface {
	exporter.Exporter
}

type exp struct {
	exporter *jaeger.Exporter
	options  jaegerOptions
}

func New(opts ...JaegerOption) (j Jaeger, err error) {
	jo := new(jaegerOptions)

	for _, opt := range append(jaegerDefaultOpts, opts...) {
		err = opt(jo)
		if err != nil {
			return nil, err
		}
	}

	return &exp{
		options: *jo,
	}, nil
}

func (e *exp) Start(ctx context.Context) (err error) {
	e.exporter, err = jaeger.NewExporter(e.options)
	if err != nil {
		return err
	}

	trace.RegisterExporter(e.exporter)

	return nil
}

func (e *exp) Stop(ctx context.Context) {
	if e.exporter != nil {
		e.exporter.Flush()
	}
}
