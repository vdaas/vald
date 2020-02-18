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
	"context"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/observability/collector"
	"github.com/vdaas/vald/internal/observability/exporter/jaeger"
	"github.com/vdaas/vald/internal/observability/exporter/prometheus"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/observability/metrics/cpu"
	"github.com/vdaas/vald/internal/observability/metrics/grpc"
	"github.com/vdaas/vald/internal/observability/metrics/mem"
	"github.com/vdaas/vald/internal/observability/metrics/runtime"
	"github.com/vdaas/vald/internal/safety"
)

type Observability interface {
	PreStart(ctx context.Context) error
	Start(ctx context.Context) <-chan error
	Stop(ctx context.Context)
}

type observability struct {
	jaeger     jaeger.Jaeger
	prometheus prometheus.Prometheus
	collector  collector.Collector
	eg         errgroup.Group
}

func New(cfg *config.Observability) (Observability, error) {
	o := new(observability)
	if cfg != nil {
		if cfg.Collector != nil {
			cpuMetric, err := cpu.NewMetric()
			if err != nil {
				return nil, err
			}

			col, err := collector.New(
				collector.WithDuration(cfg.Collector.Duration),
				collector.WithMetrics(
					cpuMetric,
					mem.NewMetric(),
					runtime.NewNumberOfGoroutines(),
					runtime.NewNumberOfCGOCall(),
				),
			)
			if err != nil {
				return nil, err
			}

			o.collector = col
		}
		if cfg.Prometheus != nil && cfg.Prometheus.Enabled {
			prom, err := prometheus.New(
				prometheus.WithNamespace(cfg.Prometheus.Namespace),
			)
			if err != nil {
				return nil, err
			}

			o.prometheus = prom
		}

		if cfg.Jaeger != nil && cfg.Jaeger.Enabled {
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

			o.jaeger = jae
		}
	}

	o.eg = errgroup.Get()

	return o, nil
}

func (o *observability) PreStart(ctx context.Context) (err error) {
	err = metrics.RegisterView(grpc.DefaultServerViews...)
	if err != nil {
		return err
	}
	err = o.collector.PreStart(ctx)
	if err != nil {
		return err
	}
	if o.prometheus != nil {
		err = o.prometheus.Start(ctx)
		if err != nil {
			return err
		}
	}
	if o.jaeger != nil {
		err = o.jaeger.Start(ctx)
		if err != nil {
			return err
		}
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
	if o.prometheus != nil {
		o.prometheus.Stop(ctx)
	}
	if o.jaeger != nil {
		o.jaeger.Stop(ctx)
	}
}
