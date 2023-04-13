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
package observability

import (
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/exporter"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/observability/trace"
)

type Option func(*observability) error

var defaultOpts = []Option{
	WithErrGroup(errgroup.Get()),
}

// WithErrGroup returns an option that sets the errgroup.
func WithErrGroup(eg errgroup.Group) Option {
	return func(o *observability) error {
		if eg == nil {
			return errors.NewErrInvalidOption("errgroup", eg)
		}
		o.eg = eg
		return nil
	}
}

// WithMetrics returns an option that sets the metrics.
func WithMetrics(ms ...metrics.Metric) Option {
	return func(o *observability) error {
		if len(ms) != 0 {
			if o.metrics == nil {
				o.metrics = ms
			} else {
				o.metrics = append(o.metrics, ms...)
			}
		}
		return nil
	}
}

// WithExporters returns an option that sets the exporters.
func WithExporters(exps ...exporter.Exporter) Option {
	return func(o *observability) error {
		if len(exps) != 0 {
			if o.exporters == nil {
				o.exporters = exps
			} else {
				o.exporters = append(o.exporters, exps...)
			}
		}
		return nil
	}
}

// WithTracer returns an option that sets the tracer.
func WithTracer(tr trace.Tracer) Option {
	return func(o *observability) error {
		if tr == nil {
			return errors.NewErrInvalidOption("tracer", tr)
		}
		o.tracer = tr
		return nil
	}
}
