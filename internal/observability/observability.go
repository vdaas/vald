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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/client/google"
	"github.com/vdaas/vald/internal/observability/collector"
	"github.com/vdaas/vald/internal/observability/exporter"
	"github.com/vdaas/vald/internal/observability/exporter/jaeger"
	"github.com/vdaas/vald/internal/observability/exporter/prometheus"
	"github.com/vdaas/vald/internal/observability/exporter/stackdriver"
	"github.com/vdaas/vald/internal/observability/metrics"
	"github.com/vdaas/vald/internal/observability/metrics/grpc"
	"github.com/vdaas/vald/internal/observability/profiler"
	sdprof "github.com/vdaas/vald/internal/observability/profiler/stackdriver"
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
	profilers []profiler.Profiler
}

func NewWithConfig(cfg *config.Observability, metrics ...metrics.Metric) (Observability, error) {
	opts := make([]Option, 0)
	traceOpts := make([]trace.Option, 0)
	exps := make([]exporter.Exporter, 0)
	profs := make([]profiler.Profiler, 0)

	col, err := newCollector(cfg.Collector, metrics...)
	if err != nil {
		return nil, err
	}
	opts = append(opts, WithCollector(col))

	if cfg.Prometheus.Enabled {
		prom, err := newPrometheus(cfg.Prometheus)
		if err != nil {
			return nil, err
		}

		exps = append(exps, prom)
	}

	if cfg.Jaeger.Enabled {
		jae, err := newJaeger(cfg.Jaeger)
		if err != nil {
			return nil, err
		}

		exps = append(exps, jae)
		traceOpts = append(traceOpts, trace.WithSyncer(jae.Exporter()))
	}

	if cfg.Stackdriver.Exporter.MonitoringEnabled || cfg.Stackdriver.Exporter.TracingEnabled {
		sdex, err := newStackdriverExporter(cfg.Stackdriver)
		if err != nil {
			return nil, err
		}

		exps = append(exps, sdex)
	}

	if cfg.Stackdriver.Profiler.Enabled {
		sdp, err := newStackdriverProfiler(cfg.Stackdriver)
		if err != nil {
			return nil, err
		}

		profs = append(profs, sdp)
	}

	if cfg.Trace.Enabled {
		tr, err := trace.New(
			append(
				traceOpts,
				trace.WithSamplingRate(cfg.Trace.SamplingRate),
			)...,
		)
		if err != nil {
			return nil, err
		}

		opts = append(opts, WithTracer(tr))
	}

	opts = append(
		opts,
		WithExporters(exps...),
		WithProfilers(profs...),
	)

	return New(opts...)
}

func newCollector(cfg *config.Collector, metrics ...metrics.Metric) (collector.Collector, error) {
	return collector.New(
		collector.WithDuration(cfg.Duration),
		collector.WithVersionInfo(
			cfg.Metrics.EnableVersionInfo,
			cfg.Metrics.VersionInfoLabels...,
		),
		collector.WithMemoryMetrics(cfg.Metrics.EnableMemory),
		collector.WithGoroutineMetrics(cfg.Metrics.EnableGoroutine),
		collector.WithCGOMetrics(cfg.Metrics.EnableCGO),
		collector.WithMetrics(metrics...),
	)
}

func newPrometheus(cfg *config.Prometheus) (prometheus.Prometheus, error) {
	return prometheus.New(
		prometheus.WithEndpoint(cfg.Endpoint),
		prometheus.WithNamespace(cfg.Namespace),
	)
}

func newJaeger(cfg *config.Jaeger) (jaeger.Jaeger, error) {
	return jaeger.New(
		jaeger.WithCollectorEndpoint(cfg.CollectorEndpoint),
		jaeger.WithAgentEndpoint(cfg.AgentEndpoint),
		jaeger.WithUsername(cfg.Username),
		jaeger.WithPassword(cfg.Password),
		jaeger.WithServiceName(cfg.ServiceName),
		jaeger.WithBufferMaxCount(cfg.BufferMaxCount),
	)
}

func stackdriverClientOpts(cfg *config.Stackdriver) []google.Option {
	return []google.Option{
		google.WithAPIKey(cfg.Client.APIKey),
		google.WithAudiences(cfg.Client.Audiences...),
		google.WithCredentialsFile(cfg.Client.CredentialsFile),
		google.WithCredentialsJSON(cfg.Client.CredentialsJSON),
		google.WithEndpoint(cfg.Client.Endpoint),
		google.WithQuotaProject(cfg.Client.QuotaProject),
		google.WithRequestReason(cfg.Client.RequestReason),
		google.WithScopes(cfg.Client.Scopes...),
		google.WithUserAgent(cfg.Client.UserAgent),
		google.WithTelemetry(cfg.Client.TelemetryEnabled),
		google.WithAuthentication(cfg.Client.AuthenticationEnabled),
	}
}

func newStackdriverExporter(cfg *config.Stackdriver) (stackdriver.Stackdriver, error) {
	clientOpts := stackdriverClientOpts(cfg)
	return stackdriver.New(
		stackdriver.WithProjectID(cfg.ProjectID),
		stackdriver.WithMonitoring(cfg.Exporter.MonitoringEnabled),
		stackdriver.WithTracing(cfg.Exporter.TracingEnabled),
		stackdriver.WithLocation(cfg.Exporter.Location),
		stackdriver.WithBundleDelayThreshold(cfg.Exporter.BundleDelayThreshold),
		stackdriver.WithBundleCountThreshold(cfg.Exporter.BundleCountThreshold),
		stackdriver.WithTraceSpansBufferMaxBytes(cfg.Exporter.TraceSpansBufferMaxBytes),
		stackdriver.WithMetricPrefix(cfg.Exporter.MetricPrefix),
		stackdriver.WithSkipCMD(cfg.Exporter.SkipCMD),
		stackdriver.WithTimeout(cfg.Exporter.Timeout),
		stackdriver.WithReportingInterval(cfg.Exporter.ReportingInterval),
		stackdriver.WithNumberOfWorkers(cfg.Exporter.NumberOfWorkers),
		stackdriver.WithMonitoringClientOptions(clientOpts...),
		stackdriver.WithTraceClientOptions(clientOpts...),
	)
}

func newStackdriverProfiler(cfg *config.Stackdriver) (sdprof.Stackdriver, error) {
	clientOpts := stackdriverClientOpts(cfg)
	return sdprof.New(
		sdprof.WithProjectID(cfg.ProjectID),
		sdprof.WithService(cfg.Profiler.Service),
		sdprof.WithServiceVersion(cfg.Profiler.ServiceVersion),
		sdprof.WithDebugLogging(cfg.Profiler.DebugLogging),
		sdprof.WithMutexProfiling(cfg.Profiler.MutexProfiling),
		sdprof.WithCPUProfiling(cfg.Profiler.CPUProfiling),
		sdprof.WithAllocProfiling(cfg.Profiler.AllocProfiling),
		sdprof.WithHeapProfiling(cfg.Profiler.HeapProfiling),
		sdprof.WithGoroutineProfiling(cfg.Profiler.GoroutineProfiling),
		sdprof.WithAllocForceGC(cfg.Profiler.AllocForceGC),
		sdprof.WithAPIAddr(cfg.Profiler.APIAddr),
		sdprof.WithInstance(cfg.Profiler.Instance),
		sdprof.WithZone(cfg.Profiler.Zone),
		sdprof.WithClientOptions(clientOpts...),
	)
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

	for i, prof := range o.profilers {
		err = prof.Start(ctx)
		if err != nil {
			for _, ex := range o.exporters {
				ex.Stop(ctx)
			}
			for _, prof = range o.profilers[:i] {
				prof.Stop(ctx)
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

	for _, prof := range o.profilers {
		prof.Stop(ctx)
	}
}
