//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/vdaas/vald/internal/log"
)

type JaegerOption func(*jaegerOptions) error

var jaegerDefaultOpts = []JaegerOption{
	WithServiceName("vald"),
	WithOnErrorFunc(func(err error) {
		if err != nil {
			log.Warnf("Error when uploading spans to Jaeger: %v", err)
		}
	}),
}

func WithCollectorEndpoint(cep string) JaegerOption {
	return func(jo *jaegerOptions) error {
		if cep != "" {
			jo.CollectorEndpoint = cep
		}
		return nil
	}
}

func WithAgentEndpoint(aep string) JaegerOption {
	return func(jo *jaegerOptions) error {
		if aep != "" {
			jo.AgentEndpoint = aep
		}
		return nil
	}
}

func WithUsername(username string) JaegerOption {
	return func(jo *jaegerOptions) error {
		if username != "" {
			jo.Username = username
		}
		return nil
	}
}

func WithPassword(password string) JaegerOption {
	return func(jo *jaegerOptions) error {
		if password != "" {
			jo.Password = password
		}
		return nil
	}
}

func WithServiceName(serviceName string) JaegerOption {
	return func(jo *jaegerOptions) error {
		if serviceName != "" {
			jo.Process = jaeger.Process{ServiceName: serviceName}
		}
		return nil
	}
}

func WithBufferMaxCount(cnt int) JaegerOption {
	return func(jo *jaegerOptions) error {
		jo.BufferMaxCount = cnt
		return nil
	}
}

func WithOnErrorFunc(f func(error)) JaegerOption {
	return func(jo *jaegerOptions) error {
		if f != nil {
			jo.OnError = f
		}
		return nil
	}
}
