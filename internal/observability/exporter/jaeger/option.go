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
	"time"

	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/vdaas/vald/internal/errors"
)

type Option func(*export) error

var jaegerDefaultOpts = []Option{
	WithAgentMaxPacketSize(65000),
	WithAgentReconnectInterval("30s"),
	WithServiceName("vald"),
	WithHTTPClient(http.DefaultClient),
	WithBatchTimeout("5s"),
	WithExportTimeout("30s"),
	WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
	WithMaxQueueSize(trace.DefaultMaxQueueSize),
}

func WithAgentEndpoint(aep string) Option {
	return func(exp *export) error {
		if len(aep) == 0 {
			return errors.NewErrInvalidOption("agentEndpoint", aep)
		}
		host, port, err := net.SplitHostPort(aep)
		if err != nil {
			return errors.NewErrCriticalOption("agentEndpoint", aep, err)
		}
		exp.agentHost = host
		exp.agentPort = port
		return nil
	}
}

func WithAgentReconnectInterval(dur string) Option {
	return func(e *export) error {
		if len(dur) == 0 {
			return errors.NewErrInvalidOption("agentReconnectInterval", dur)
		}
		d, err := time.ParseDuration(dur)
		if err != nil {
			return errors.NewErrInvalidOption("agentReconnectInterval", dur, err)
		}
		e.agentReconnInterval = d
		return nil
	}
}

func WithAgentMaxPacketSize(cnt int) Option {
	return func(exp *export) error {
		if cnt <= 0 {
			return errors.NewErrInvalidOption("agentMaxPacketSize", cnt)
		}
		exp.agentMaxPacketSize = cnt
		return nil
	}
}

func WithCollectorEndpoint(cep string) Option {
	return func(exp *export) error {
		if len(cep) == 0 {
			return errors.NewErrInvalidOption("collectorEndpoint", cep)
		}
		exp.collectorEndpoint = cep
		return nil
	}
}

func WithHTTPClient(c *http.Client) Option {
	return func(exp *export) error {
		if c == nil {
			return errors.NewErrInvalidOption("httpClient", c)
		}
		exp.client = c
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
		if len(password) == 0 {
			return errors.NewErrInvalidOption("password", password)
		}
		exp.collectorPassword = password
		return nil
	}
}

func WithServiceName(serviceName string) Option {
	return func(exp *export) error {
		if len(serviceName) == 0 {
			return errors.NewErrInvalidOption("serviceName", serviceName)
		}
		exp.serviceName = serviceName
		return nil
	}
}

func WithBatchTimeout(dur string) Option {
	return func(e *export) error {
		if len(dur) == 0 {
			return errors.NewErrInvalidOption("batchTimeout", dur)
		}
		d, err := time.ParseDuration(dur)
		if err != nil {
			return errors.NewErrInvalidOption("batchTimeout", dur, err)
		}
		e.batchTimeout = d
		return nil
	}
}

func WithExportTimeout(dur string) Option {
	return func(e *export) error {
		if len(dur) == 0 {
			return errors.NewErrInvalidOption("exportTimeout", dur)
		}
		d, err := time.ParseDuration(dur)
		if err != nil {
			return errors.NewErrInvalidOption("exportTimeout", dur, err)
		}
		e.exportTimeout = d
		return nil
	}
}

func WithMaxExportBatchSize(size int) Option {
	return func(e *export) error {
		if size <= 0 {
			return errors.NewErrInvalidOption("maxExportBatchSize", size)
		}
		e.maxExportBatchSize = size
		return nil
	}
}

func WithMaxQueueSize(size int) Option {
	return func(e *export) error {
		if size <= 0 {
			return errors.NewErrInvalidOption("maxQueueSize", size)
		}
		e.maxQueueSize = size
		return nil
	}
}
