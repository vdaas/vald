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

// Package config providers configuration type and load configuration logic
package config

// Observability represents the configuration for the observability.
type Observability struct {
	Enabled    bool        `json:"enabled"    yaml:"enabled"`
	Metrics    *Metrics    `json:"metrics"    yaml:"metrics"`
	Trace      *Trace      `json:"trace"      yaml:"trace"`
	Prometheus *Prometheus `json:"prometheus" yaml:"prometheus"`
	Jaeger     *Jaeger     `json:"jaeger"     yaml:"jaeger"`
}

// Collector represents the configuration for the collector.
type Collector struct {
	Duration string   `json:"duration" yaml:"duration"`
	Metrics  *Metrics `json:"metrics"  yaml:"metrics"`
}

// Trace represents the configuration for the trace.
type Trace struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// Metrics represents the configuration for the metrics.
type Metrics struct {
	EnableVersionInfo bool     `json:"enable_version_info" yaml:"enable_version_info"`
	VersionInfoLabels []string `json:"version_info_labels" yaml:"version_info_labels"`
	EnableMemory      bool     `json:"enable_memory"       yaml:"enable_memory"`
	EnableGoroutine   bool     `json:"enable_goroutine"    yaml:"enable_goroutine"`
	EnableCGO         bool     `json:"enable_cgo"          yaml:"enable_cgo"`
}

// Prometheus represents the configuration for the prometheus.
type Prometheus struct {
	Enabled            bool   `json:"enabled"               yaml:"enabled"`
	Endpoint           string `json:"endpoint"              yaml:"endpoint"`
	Namespace          string `json:"namespace"             yaml:"namespace"`
	CollectInterval    string `json:"collect_interval"      yaml:"collect_interval"`
	CollectTimeout     string `json:"collect_timeout"       yaml:"collect_timeout"`
	EnableInMemoryMode bool   `json:"enable_in_memory_mode" yaml:"enable_in_memory_mode"`
}

// Jaeger represents the configuration for the jaeger.
type Jaeger struct {
	Enabled bool `json:"enabled" yaml:"enabled"`

	CollectorEndpoint string `json:"collector_endpoint" yaml:"collector_endpoint"`
	AgentEndpoint     string `json:"agent_endpoint"     yaml:"agent_endpoint"`

	AgentMaxPacketSize int `json:"agent_max_packet_size" yaml:"agent_max_packet_size"`

	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`

	ServiceName        string `json:"service_name" yaml:"service_name"`
	BatchTimeout       string `json:"batch_timeout" yaml:"batch_timeout"`
	ExportTimeout      string `json:"export_timeout" yaml:"export_timeout"`
	MaxExportBatchSize int    `json:"max_export_batch_size" yaml:"max_export_batch_size"`
	MaxQueueSize       int    `json:"max_queue_size" yaml:"max_queue_size"`
}

// Bind binds the actual data from the Observability receiver fields.
func (o *Observability) Bind() *Observability {
	if o.Metrics != nil {
		o.Metrics.VersionInfoLabels = GetActualValues(o.Metrics.VersionInfoLabels)
	} else {
		o.Metrics = new(Metrics)
	}

	if o.Trace == nil {
		o.Trace = new(Trace)
	}

	if o.Prometheus != nil {
		o.Prometheus.Endpoint = GetActualValue(o.Prometheus.Endpoint)
		o.Prometheus.Namespace = GetActualValue(o.Prometheus.Namespace)
		o.Prometheus.CollectInterval = GetActualValue(o.Prometheus.CollectInterval)
		o.Prometheus.CollectTimeout = GetActualValue(o.Prometheus.CollectTimeout)
	} else {
		o.Prometheus = new(Prometheus)
	}

	if o.Jaeger != nil {
		o.Jaeger.CollectorEndpoint = GetActualValue(o.Jaeger.CollectorEndpoint)
		o.Jaeger.AgentEndpoint = GetActualValue(o.Jaeger.AgentEndpoint)
		o.Jaeger.Username = GetActualValue(o.Jaeger.Username)
		o.Jaeger.Password = GetActualValue(o.Jaeger.Password)
		o.Jaeger.ServiceName = GetActualValue(o.Jaeger.ServiceName)
		o.Jaeger.BatchTimeout = GetActualValue(o.Jaeger.BatchTimeout)
		o.Jaeger.ExportTimeout = GetActualValue(o.Jaeger.ExportTimeout)
	} else {
		o.Jaeger = new(Jaeger)
	}

	return o
}

// Bind binds the actual data from the Collector receiver fields.
func (c *Collector) Bind() *Collector {
	c.Duration = GetActualValue(c.Duration)

	if c.Metrics != nil {
		c.Metrics.VersionInfoLabels = GetActualValues(c.Metrics.VersionInfoLabels)
	} else {
		c.Metrics = new(Metrics)
	}

	return c
}
