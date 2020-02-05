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
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/observability/exporter/metric/prometheus"
	"github.com/vdaas/vald/internal/observability/exporter/trace/jaeger"
)

func New(cfg *config.Observability) error {
	if cfg != nil {
		if cfg.Prometheus != nil && cfg.Prometheus.Enabled {
			err := prometheus.Init(
				prometheus.WithNamespace(cfg.Prometheus.Namespace),
			)
			if err != nil {
				return err
			}
		}

		if cfg.Jaeger != nil && cfg.Jaeger.Enabled {
			err := jaeger.Init(
				jaeger.WithCollectorEndpoint(cfg.Jaeger.CollectorEndpoint),
				jaeger.WithAgentEndpoint(cfg.Jaeger.AgentEndpoint),
				jaeger.WithUsername(cfg.Jaeger.Username),
				jaeger.WithPassword(cfg.Jaeger.Password),
				jaeger.WithServiceName(cfg.Jaeger.ServiceName),
				jaeger.WithBufferMaxCount(cfg.Jaeger.BufferMaxCount),
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Stop() error {
	jaeger.Flush()
	return nil
}
