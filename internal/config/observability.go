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

// Package config providers configuration type and load configuration logic
package config

type Observability struct {
	Enabled    bool        `json:"enabled" yaml:"enabled"`
	Collector  *Collector  `json:"collector" yaml:"collector"`
	Trace      *Trace      `json:"trace" yaml:"trace"`
	Prometheus *Prometheus `json:"prometheus" yaml:"prometheus"`
	Jaeger     *Jaeger     `json:"jaeger" yaml:"jaeger"`
}

type Collector struct {
	Duration string   `json:"duration" yaml:"duration"`
	Metrics  *Metrics `json:"metrics" yaml:"metrics"`
}

type Trace struct {
	Enabled      bool    `json:"enabled" yaml:"enabled"`
	SamplingRate float64 `json:"sampling_rate" yaml:"sampling_rate"`
}

type Metrics struct {
	EnableVersionInfo bool `json:"enable_version_info" yaml:"enable_version_info"`
	EnableMemory      bool `json:"enable_memory" yaml:"enable_memory"`
	EnableGoroutine   bool `json:"enable_goroutine" yaml:"enable_goroutine"`
	EnableCGO         bool `json:"enable_cgo" yaml:"enable_cgo"`
}

type Prometheus struct {
	Enabled   bool   `json:"enabled" yaml:"enabled"`
	Endpoint  string `json:"endpoint" yaml:"endpoint"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

type Jaeger struct {
	Enabled bool `json:"enabled" yaml:"enabled"`

	CollectorEndpoint string `json:"collector_endpoint" yaml:"collector_endpoint"`
	AgentEndpoint     string `json:"agent_endpoint" yaml:"agent_endpoint"`

	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`

	ServiceName string `json:"service_name" yaml:"service_name"`

	BufferMaxCount int `json:"buffer_max_count" yaml:"buffer_max_count"`
}

func (o *Observability) Bind() *Observability {
	if o.Collector != nil {
		o.Collector.Duration = GetActualValue(o.Collector.Duration)
	} else {
		o.Collector = new(Collector)
	}

	if o.Trace == nil {
		o.Trace = new(Trace)
	}

	if o.Prometheus != nil {
		o.Prometheus.Endpoint = GetActualValue(o.Prometheus.Endpoint)
		o.Prometheus.Namespace = GetActualValue(o.Prometheus.Namespace)
	} else {
		o.Prometheus = new(Prometheus)
	}

	if o.Jaeger != nil {
		o.Jaeger.CollectorEndpoint = GetActualValue(o.Jaeger.CollectorEndpoint)
		o.Jaeger.AgentEndpoint = GetActualValue(o.Jaeger.AgentEndpoint)
		o.Jaeger.Username = GetActualValue(o.Jaeger.Username)
		o.Jaeger.Password = GetActualValue(o.Jaeger.Password)
		o.Jaeger.ServiceName = GetActualValue(o.Jaeger.ServiceName)
	} else {
		o.Jaeger = new(Jaeger)
	}

	return o
}
