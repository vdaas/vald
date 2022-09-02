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

// Package jaeger provides a jaeger exporter.
package jaeger

import (
	"context"
	"net/http"
	"time"

	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/observability/exporter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type Jaeger interface {
	exporter.Exporter
}

type export struct {
	tp                  *trace.TracerProvider
	exp                 *jaeger.Exporter
	collectorEndpoint   string
	client              *http.Client
	collectorPassword   string
	collectorUserName   string
	agentHost           string
	agentPort           string
	agentReconnInterval time.Duration
	agentMaxPacketSize  int
	serviceName         string
}

func New(opts ...Option) (j Jaeger, err error) {
	e := new(export)

	for _, opt := range append(jaegerDefaultOpts, opts...) {
		err = opt(e)
		if err != nil {
			return nil, err
		}
	}

	var eop jaeger.EndpointOption
	if len(e.agentHost) != 0 && len(e.agentPort) != 0 {
		// TODO: we can not get trace data, so we need to fix it later.
		eop = jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(e.agentHost),
			jaeger.WithAgentPort(e.agentPort),
			jaeger.WithAttemptReconnectingInterval(e.agentReconnInterval),
			jaeger.WithMaxPacketSize(e.agentMaxPacketSize))
	} else {
		eop = jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(e.collectorEndpoint),
			jaeger.WithHTTPClient(http.DefaultClient),
			// jaeger.WithDisableAttemptReconnecting(),
			jaeger.WithPassword(e.collectorPassword),
			jaeger.WithUsername(e.collectorUserName))
	}
	e.exp, err = jaeger.New(eop)
	if err != nil {
		return nil, err
	}
	e.tp = trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(e.exp), // TODO we should set batch option here. like below and get configuration from yaml
		// trace.WithBatcher(e.exp,
		// 	trace.WithBatchTimeout(time.Second*5),
		// 	trace.WithExportTimeout(time.Minute),
		// 	trace.WithMaxExportBatchSize(1024),
		// 	trace.WithMaxQueueSize(256),
		// 	// trace.WithBlocking(),
		// ),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(e.serviceName),
		)),
	)
	return e, nil
}

func (e *export) Start(ctx context.Context) (err error) {
	otel.SetTracerProvider(e.tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return nil
}

func (e *export) Stop(ctx context.Context) {
	var err error
	if e.tp != nil {
		err = e.tp.ForceFlush(ctx)
		if err != nil {
			log.Error(err)
		}
		err = e.tp.Shutdown(ctx)
		if err != nil {
			log.Error(err)
		}
	}
	if e.exp != nil {
		err = e.exp.Shutdown(ctx)
		if err != nil {
			log.Error(err)
		}
	}
}
