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
	"net"
	"net/http"

	"go.opentelemetry.io/otel/sdk/trace"
)

type Option func(*export) error

var jaegerDefaultOpts = []Option{
	WithServiceName("vald"),
	WithHTTPClient(http.DefaultClient),
	WithBatchTimeout("5s"),
	WithExportTimeout("30s"),
	WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
	WithMaxQueueSize(trace.DefaultMaxQueueSize),
}

func WithAgentEndpoint(aep string) Option {
	return func(exp *export) error {
		if aep != "" {
			host, port, err := net.SplitHostPort(aep)
			if err != nil {
				return err
			}
			exp.agentHost = host
			exp.agentPort = port
		}
		return nil
	}
}

func WithAgentReconnectInterval(dur string) Option {
	return func(e *export) error {
		return nil
	}
}

func WithAgentMaxPacketSize(cnt int) Option {
	return func(exp *export) error {
		exp.agentMaxPacketSize = cnt
		return nil
	}
}

func WithCollectorEndpoint(cep string) Option {
	return func(exp *export) error {
		if cep != "" {
			exp.collectorEndpoint = cep
		}
		return nil
	}
}

func WithHTTPClient(c *http.Client) Option {
	return func(exp *export) error {
		if c != nil {
			exp.client = c
		}
		return nil
	}
}

func WithUsername(username string) Option {
	return func(exp *export) error {
		if username != "" {
			exp.collectorUserName = username
		}
		return nil
	}
}

func WithPassword(password string) Option {
	return func(exp *export) error {
		if password != "" {
			exp.collectorPassword = password
		}
		return nil
	}
}

func WithServiceName(serviceName string) Option {
	return func(exp *export) error {
		if serviceName != "" {
			exp.serviceName = serviceName
		}
		return nil
	}
}

func WithBatchTimeout(dur string) Option {
	return func(e *export) error {
		return nil
	}
}

func WithExportTimeout(dur string) Option {
	return func(e *export) error {
		return nil
	}
}

func WithMaxExportBatchSize(size int) Option {
	return func(e *export) error {
		return nil
	}
}

func WithMaxQueueSize(size int) Option {
	return func(e *export) error {
		return nil
	}
}
