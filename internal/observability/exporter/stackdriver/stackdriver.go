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

// Package stackdriver provides a stackdriver exporter.
package stackdriver

import (
	"context"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"github.com/vdaas/vald/internal/observability/exporter"
	"go.opencensus.io/trace"
)

type Stackdriver interface {
	exporter.Exporter
}

type exp struct {
	exporter *stackdriver.Exporter

	monitoringEnabled bool
	tracingEnabled    bool

	*stackdriver.Options
}

func New(opts ...Option) (s Stackdriver, err error) {
	e := new(exp)
	e.Options = new(stackdriver.Options)

	for _, opt := range append(defaultOptions, opts...) {
		err = opt(e)
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *exp) Start(ctx context.Context) (err error) {
	e.Options.Context = ctx

	e.exporter, err = stackdriver.NewExporter(*e.Options)
	if err != nil {
		return err
	}

	if e.monitoringEnabled {
		err = e.exporter.StartMetricsExporter()
		if err != nil {
			return err
		}
	}

	if e.tracingEnabled {
		trace.RegisterExporter(e.exporter)
	}

	return nil
}

func (e *exp) Stop(ctx context.Context) {
	if e.exporter != nil {
		if e.monitoringEnabled {
			e.exporter.StopMetricsExporter()
		}

		e.exporter.Flush()
	}
}
