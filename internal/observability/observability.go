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
	"context"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/collector"
	"github.com/vdaas/vald/internal/observability/exporter"
	"github.com/vdaas/vald/internal/observability/exporter/jaeger"
	"github.com/vdaas/vald/internal/observability/exporter/prometheus"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/observability/metrics/grpc"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/safety"
)

type Observability interface {
	PreStart(ctx context.Context) error
	Start(ctx context.Context) <-chan error
	Stop(ctx context.Context)
}

type observability struct {
	eg        errgroup.Group
	collector collector.Collector
	tracer    trace.Tracer

	exporters []exporter.Exporter
}

func NewWithConfig(cfg *config.Observability, metrics ...metrics.Metric) (Observability, error) {
	opts := make([]Option, 0)
	exps := make([]exporter.Exporter, 0)

	col, err := collector.New(
		collector.WithDuration(cfg.Collector.Duration),
		collector.WithVersionInfo(
			cfg.Collector.Metrics.EnableVersionInfo,
			cfg.Collector.Metrics.VersionInfoLabels...,
		),
		collector.WithMemoryMetrics(cfg.Collector.Metrics.EnableMemory),
		collector.WithGoroutineMetrics(cfg.Collector.Metrics.EnableGoroutine),
		collector.WithCGOMetrics(cfg.Collector.Metrics.EnableCGO),
		collector.WithMetrics(metrics...),
	)
	if err != nil {
		return nil, err
	}
	opts = append(opts, WithCollector(col))

	if cfg.Trace.Enabled {
		opts = append(opts,
			WithTracer(trace.New(trace.WithSamplingRate(cfg.Trace.SamplingRate))),
		)
	}

	if cfg.Prometheus.Enabled {
		prom, err := prometheus.New(
			prometheus.WithEndpoint(cfg.Prometheus.Endpoint),
			prometheus.WithNamespace(cfg.Prometheus.Namespace),
		)
		if err != nil {
			return nil, err
		}

		exps = append(exps, prom)
	}

	if cfg.Jaeger.Enabled {
		jae, err := jaeger.New(
			jaeger.WithCollectorEndpoint(cfg.Jaeger.CollectorEndpoint),
			jaeger.WithAgentEndpoint(cfg.Jaeger.AgentEndpoint),
			jaeger.WithUsername(cfg.Jaeger.Username),
			jaeger.WithPassword(cfg.Jaeger.Password),
			jaeger.WithServiceName(cfg.Jaeger.ServiceName),
			jaeger.WithBufferMaxCount(cfg.Jaeger.BufferMaxCount),
		)
		if err != nil {
			return nil, err
		}

		exps = append(exps, jae)
	}

	opts = append(
		opts,
		WithExporters(exps...),
	)

	return New(opts...)
}

func New(opts ...Option) (Observability, error) {
	o := new(observability)

	for _, opt := range append(observabilityDefaultOpts, opts...) {
		err := opt(o)
		if err != nil {
			return nil, err
		}
	}

	return o, nil
}

func (o *observability) PreStart(ctx context.Context) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	err = metrics.RegisterView(grpc.DefaultServerViews...)
	if err != nil {
		return err
	}

	if o.collector == nil {
		return errors.ErrCollectorNotFound()
	}
	err = o.collector.PreStart(ctx)
	if err != nil {
		return err
	}

	for i, ex := range o.exporters {
		err = ex.Start(ctx)
		if err != nil {
			for _, ex = range o.exporters[:i] {
				ex.Stop(ctx)
			}
			return err
		}
	}

	if o.tracer != nil {
		o.tracer.Start(ctx)
	}

	return nil
}

func (o *observability) Start(ctx context.Context) <-chan error {
	ech := make(chan error, 2)

	o.eg.Go(safety.RecoverFunc(func() (err error) {
		defer close(ech)
		cech := o.collector.Start(ctx)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case err = <-cech:
			}
			if err != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case ech <- err:
				}
			}
		}
	}))

	return ech
}

func (o *observability) Stop(ctx context.Context) {
	o.collector.Stop(ctx)

	for _, ex := range o.exporters {
		ex.Stop(ctx)
	}
}
