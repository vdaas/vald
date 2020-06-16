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

// Package stackdriver provides a stackdriver exporter.
package stackdriver

import (
	"context"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
)

type Stackdriver interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context)
	Exporter() *stackdriver.Exporter
}

type exporter struct {
	exporter *stackdriver.Exporter
	*stackdriver.Options
}

func New(opts ...Option) (s Stackdriver, err error) {
	e := new(exporter)

	for _, opt := range append(defaultOpts, opts...) {
		err = opt(e)
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *exporter) Start(ctx context.Context) (err error) {
	e.Options.Context = ctx

	e.exporter, err = stackdriver.NewExporter(*e.Options)
	if err != nil {
		return err
	}

	err = e.exporter.StartMetricsExporter()
	if err != nil {
		return err
	}

	trace.RegisterExporter(e.exporter)

	return nil
}

func (e *exporter) Stop(ctx context.Context) {
	if e.exporter != nil {
		e.exporter.StopMetricsExporter()
		e.exporter.Flush()
	}
}

func (e *exporter) Exporter() *stackdriver.Exporter {
	return e.exporter
}
