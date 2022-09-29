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
	// exporter *prometheus.Exporter
	views metrics.View

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

	// // Create controller for prometheus exporter.
	// controller := basic.New(
	// 	processor.NewFactory(
	// 		simple.NewWithHistogramDistribution(
	// 			histogram.WithExplicitBoundaries(e.histogramBoundarie),
	// 		),
	// 		aggregation.CumulativeTemporalitySelector(),
	// 		processor.WithMemory(e.inmemoryEnabled),
	// 	),
	// 	basic.WithCollectPeriod(e.collectInterval),
	// 	basic.WithCollectTimeout(e.collectTimeout),
	// 	basic.WithResource(resource.NewWithAttributes(
	// 		semconv.SchemaURL,
	// 		semconv.ServiceNamespaceKey.String(e.namespace),
	// 	)),
	// )

	// cfg := prometheus.Config{
	// 	DefaultHistogramBoundaries: e.histogramBoundarie,
	// }
	//
	// var err error
	// e.exporter, err = prometheus.New(cfg, nil)
	// if err != nil {
	// 	return nil, err
	// }

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
	// global.SetMeterProvider(e.exporter.MeterProvider())
	return nil
	// return e.exporter.Controller().Start(ctx)
}

func (e *exp) Stop(ctx context.Context) error {
	return nil
	// return e.exporter.Controller().Stop(ctx)
}

func (e *exp) NewHTTPHandler() http.Handler {
	mux := http.NewServeMux()
	// mux.Handle(e.endpoint, e.exporter)
	return mux
}

func Exporter() (Prometheus, error) {
	if instance == nil {
		return Init()
	}
	return instance, nil
}
