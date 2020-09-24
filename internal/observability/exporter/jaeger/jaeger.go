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

// Package jaeger provides a jaeger exporter.
package jaeger

import (
	"context"
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/exporter"
	"go.opentelemetry.io/otel/exporter/trace/jaeger"
)

type Jaeger interface {
	exporter.Exporter
	Exporter() *jaeger.Exporter
}

type exp struct {
	exporter *jaeger.Exporter
	name     string

	collectorEndpoint string
	agentEndpoint     string

	options       []jaeger.Option
	collectorOpts []jaeger.CollectorEndpointOption
}

func New(opts ...Option) (j Jaeger, err error) {
	e := new(exp)

	for _, opt := range append(defaultOpts, opts...) {
		if err = opt(e); err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))

			er := new(errors.ErrCriticalOption)
			if errors.As(err, &er) {
				log.Error(werr)
				return nil, werr
			}
			log.Warn(werr)
		}
	}

	return e, nil
}

func (e *exp) Start(ctx context.Context) (err error) {
	e.exporter, err = jaeger.NewExporter(e.endpoint(), e.options...)

	return err
}

func (e *exp) endpoint() jaeger.EndpointOption {
	if e.collectorEndpoint != "" {
		return jaeger.WithCollectorEndpoint(e.collectorEndpoint, e.collectorOpts...)
	}

	return jaeger.WithAgentEndpoint(e.agentEndpoint)
}

func (e *exp) Stop(ctx context.Context) {
	if e.exporter != nil {
		e.exporter.Flush()
	}
}

func (e *exp) Exporter() *jaeger.Exporter {
	return e.exporter
}
