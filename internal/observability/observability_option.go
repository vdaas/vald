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

// Package observability provides observability functions
package observability

import (
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/observability/collector"
	"github.com/vdaas/vald/internal/observability/exporter"
	"github.com/vdaas/vald/internal/observability/trace"
)

type Option func(*observability) error

var observabilityDefaultOpts = []Option{
	WithErrGroup(errgroup.Get()),
}

func WithErrGroup(eg errgroup.Group) Option {
	return func(o *observability) error {
		if eg != nil {
			o.eg = eg
		}
		return nil
	}
}

func WithCollector(c collector.Collector) Option {
	return func(o *observability) error {
		if c != nil {
			o.collector = c
		}
		return nil
	}
}

func WithTracer(t trace.Tracer) Option {
	return func(o *observability) error {
		if t != nil {
			o.tracer = t
		}
		return nil
	}
}

func WithExporters(exps ...exporter.Exporter) Option {
	return func(o *observability) error {
		if o.exporters == nil {
			o.exporters = exps
			return nil
		}

		o.exporters = append(o.exporters, exps...)

		return nil
	}
}
