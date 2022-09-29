// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package prometheus

import (
	"context"
	"net/http"
	"reflect"
	"sync"
	"time"

	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/view"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/exporter"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type Prometheus interface {
	exporter.Exporter
	NewHTTPHandler() http.Handler
}

type exp struct {
	exporter otelprom.Exporter
	viewers  []metrics.Viewer
	registry *prometheus.Registry

	namespace          string
	endpoint           string
	collectInterval    time.Duration
	collectTimeout     time.Duration
	inmemoryEnabled    bool
	histogramBoundarie []float64
}

var (
	instance Prometheus
	once     sync.Once
)

func New(opts ...Option) (Prometheus, error) {
	e := &exp{}
	for _, opt := range append(defaultOpts, opts...) {
		if err := opt(e); err != nil {
			oerr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := &errors.ErrCriticalOption{}
			if errors.As(oerr, &e) {
				log.Error(oerr)
				return nil, oerr
			}
			log.Warn(oerr)
		}
	}
	e.exporter = otelprom.New()
	e.registry = prometheus.NewRegistry()

	return e, nil
}

func Init(opts ...Option) (Prometheus, error) {
	var err error
	once.Do(func() {
		instance, err = New(opts...)
	})
	if err != nil {
		once = sync.Once{}
	}
	return instance, err
}

func (e *exp) Start(ctx context.Context) error {
	otlViews := make([]view.View, 0, len(e.viewers))
	for _, viewer := range e.viewers {
		views, err := viewer.View()
		if err != nil {
			return err
		}
		for _, v := range views {
			otlViews = append(otlViews, *v)
		}
	}

	provider := metric.NewMeterProvider(metric.WithReader(
		e.exporter,
		otlViews...,
	))
	global.SetMeterProvider(provider)

	registry := prometheus.NewRegistry()
	if err := registry.Register(e.exporter.Collector); err != nil {
		return err
	}

	return nil
}

func (e *exp) Stop(ctx context.Context) error {
	return e.exporter.Shutdown(ctx)
}

func (e *exp) NewHTTPHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(e.endpoint, promhttp.HandlerFor(e.registry, promhttp.HandlerOpts{}))
	return mux
}

func Exporter() (Prometheus, error) {
	if instance == nil {
		return Init()
	}
	return instance, nil
}
