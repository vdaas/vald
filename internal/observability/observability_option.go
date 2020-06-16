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

// Package observability provides observability functions
package observability

import (
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/observability/collector"
	"github.com/vdaas/vald/internal/observability/exporter/jaeger"
	"github.com/vdaas/vald/internal/observability/exporter/prometheus"
	"github.com/vdaas/vald/internal/observability/exporter/stackdriver"
	"github.com/vdaas/vald/internal/observability/trace"
)

type Option func(*observability) error

var (
	observabilityDefaultOpts = []Option{
		WithErrGroup(errgroup.Get()),
	}
)

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

func WithPrometheus(p prometheus.Prometheus) Option {
	return func(o *observability) error {
		if p != nil {
			o.prometheus = p
		}
		return nil
	}
}

func WithJaeger(j jaeger.Jaeger) Option {
	return func(o *observability) error {
		if j != nil {
			o.jaeger = j
		}
		return nil
	}
}

func WithStackdriver(sd stackdriver.Stackdriver) Option {
	return func(o *observability) error {
		if sd != nil {
			o.stackdriver = sd
		}
		return nil
	}
}
