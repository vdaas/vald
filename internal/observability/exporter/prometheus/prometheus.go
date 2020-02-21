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

// Package prometheus provides a prometheus exporter.
package prometheus

import (
	"context"
	"net/http"
	"sync"

	"contrib.go.opencensus.io/exporter/prometheus"
)

var (
	instance *exporter
	once     sync.Once
)

type prometheusOptions struct {
	endpoint string
	options  *prometheus.Options
}

type Prometheus interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context)
	Exporter() *prometheus.Exporter
}

type exporter struct {
	exporter *prometheus.Exporter
	options  prometheusOptions
}

func New(opts ...PrometheusOption) (Prometheus, error) {
	po := new(prometheusOptions)
	po.options = new(prometheus.Options)

	for _, opt := range append(prometheusDefaultOpts, opts...) {
		err := opt(po)
		if err != nil {
			return nil, err
		}
	}

	ex, err := prometheus.NewExporter(*po.options)
	if err != nil {
		return nil, err
	}

	e := exporter{
		exporter: ex,
		options:  *po,
	}

	once.Do(func() {
		instance = &e
	})

	return &e, nil
}

func (e *exporter) Start(ctx context.Context) error {
	return nil
}

func (e *exporter) Stop(ctx context.Context) {
}

func (e *exporter) Exporter() *prometheus.Exporter {
	return e.exporter
}

func Exporter() *prometheus.Exporter {
	return instance.exporter
}

func NewHTTPHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/metrics", instance.exporter)
	return mux
}
