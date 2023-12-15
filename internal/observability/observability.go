// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/exporter"
	"github.com/vdaas/vald/internal/observability/exporter/otlp"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/observability/metrics/grpc"
	"github.com/vdaas/vald/internal/observability/metrics/mem"
	"github.com/vdaas/vald/internal/observability/metrics/runtime/cgo"
	"github.com/vdaas/vald/internal/observability/metrics/runtime/goroutine"
	"github.com/vdaas/vald/internal/observability/metrics/version"
	"github.com/vdaas/vald/internal/observability/trace"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

type Observability interface {
	PreStart(ctx context.Context) error
	Start(ctx context.Context) <-chan error
	Stop(ctx context.Context) error
}

type observability struct {
	eg        errgroup.Group
	exporters []exporter.Exporter
	tracer    trace.Tracer
	metrics   []metrics.Metric
}

func NewWithConfig(cfg *config.Observability, ms ...metrics.Metric) (Observability, error) {
	opts := make([]Option, 0)
	exps := make([]exporter.Exporter, 0)

	if cfg.Metrics != nil {
		ms = append(ms, grpc.New())
		if cfg.Metrics.EnableCGO {
			ms = append(ms, cgo.New())
		}
		if cfg.Metrics.EnableGoroutine {
			ms = append(ms, goroutine.New())
		}
		if cfg.Metrics.EnableMemory {
			ms = append(ms, mem.New())
		}
		if cfg.Metrics.EnableVersionInfo {
			ms = append(ms, version.New(cfg.Metrics.VersionInfoLabels...))
		}
	}

	if cfg.Trace.Enabled {
		tr, err := trace.New()
		if err != nil {
			return nil, err
		}
		opts = append(opts, WithTracer(tr))
	}

	if cfg.OTLP != nil {
		e, err := otlp.New(
			otlp.WithCollectorEndpoint(cfg.OTLP.CollectorEndpoint),
			otlp.WithTraceBatchTimeout(cfg.OTLP.TraceBatchTimeout),
			otlp.WithTraceExportTimeout(cfg.OTLP.TraceExportTimeout),
			otlp.WithTraceMaxExportBatchSize(cfg.OTLP.TraceMaxExportBatchSize),
			otlp.WithTraceMaxQueueSize(cfg.OTLP.TraceMaxQueueSize),
			otlp.WithMetricsExportInterval(cfg.OTLP.MetricsExportInterval),
			otlp.WithMetricsExportTimeout(cfg.OTLP.MetricsExportTimeout),
			otlp.WithAttributes(
				otlp.ServiceNameKey.String(cfg.OTLP.Attribute.ServiceName),
				otlp.NamespaceKey.String(cfg.OTLP.Attribute.Namespace),
				otlp.TargetPodNameKey.String(cfg.OTLP.Attribute.PodName),
				otlp.TargetNodeNameKey.String(cfg.OTLP.Attribute.NodeName),
				otlp.AppNameKey.String(cfg.OTLP.Attribute.ServiceName),
			),
		)
		if err != nil {
			return nil, err
		}
		exps = append(exps, e)
	}

	opts = append(
		opts,
		WithExporters(exps...),
		WithMetrics(ms...),
	)

	return New(opts...)
}

func New(opts ...Option) (Observability, error) {
	o := &observability{}
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(o); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	return o, nil
}

func (o *observability) PreStart(ctx context.Context) error {
	for i, ex := range o.exporters {
		if err := ex.Start(ctx); err != nil {
			for _, ex := range o.exporters[:i] {
				if err := ex.Stop(ctx); err != nil {
					log.Error(err)
				}
			}
			return err
		}
	}

	meter := metrics.GetMeter()
	for _, m := range o.metrics {
		if err := m.Register(meter); err != nil {
			return err
		}
	}

	if o.tracer != nil {
		if err := o.tracer.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (*observability) Start(context.Context) <-chan error {
	return nil
}

func (o *observability) Stop(ctx context.Context) (werr error) {
	for _, ex := range o.exporters {
		if err := ex.Stop(ctx); err != nil {
			log.Error(err)
			werr = errors.Join(werr, err)
		}
	}
	return werr
}
