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
	OTLP       *OTLP       `json:"otlp"       yaml:"otlp"`
	Metrics    *Metrics    `json:"metrics"    yaml:"metrics"`
	Trace      *Trace      `json:"trace"      yaml:"trace"`
	Prometheus *Prometheus `json:"prometheus" yaml:"prometheus"`
	Jaeger     *Jaeger     `json:"jaeger"     yaml:"jaeger"`
}

type OTLP struct {
	CollectorEndpoint       string         `json:"collector_endpoint"          yaml:"collector_endpoint"`
	Attribute               *OTLPAttribute `json:"attribute"                   yaml:"attribute"`
	TraceBatchTimeout       string         `json:"trace_batch_timeout"         yaml:"trace_batch_timeout"`
	TraceExportTimeout      string         `json:"trace_export_timeout"        yaml:"trace_export_timeout"`
	TraceMaxExportBatchSize int            `json:"trace_max_export_batch_size" yaml:"trace_max_export_batch_size"`
	TraceMaxQueueSize       int            `json:"trace_max_queue_size"        yaml:"trace_max_queue_size"`
	MetricsExportInterval   string         `json:"metrics_export_interval"     yaml:"metrics_export_interval"`
	MetricsExportTimeout    string         `json:"metrics_export_timeout"      yaml:"metrics_export_timeout"`
}

type OTLPAttribute struct {
	Namespace   string `json:"namespace"    yaml:"namespace"`
	PodName     string `json:"pod_name"     yaml:"pod_name"`
	NodeName    string `json:"node_name"    yaml:"node_name"`
	ServiceName string `json:"service_name" yaml:"service_name"`
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
	Enabled   bool   `json:"enabled"   yaml:"enabled"`
	Endpoint  string `json:"endpoint"  yaml:"endpoint"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

// Jaeger represents the configuration for the jaeger.
type Jaeger struct {
	Enabled bool `json:"enabled" yaml:"enabled"`

	CollectorEndpoint string `json:"collector_endpoint" yaml:"collector_endpoint"`
	AgentEndpoint     string `json:"agent_endpoint"     yaml:"agent_endpoint"`

	AgentMaxPacketSize     int    `json:"agent_max_packet_size"    yaml:"agent_max_packet_size"`
	AgentReconnectInterval string `json:"agent_reconnect_interval" yaml:"agent_reconnect_interval"`

	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`

	ServiceName        string `json:"service_name"          yaml:"service_name"`
	BatchTimeout       string `json:"batch_timeout"         yaml:"batch_timeout"`
	ExportTimeout      string `json:"export_timeout"        yaml:"export_timeout"`
	MaxExportBatchSize int    `json:"max_export_batch_size" yaml:"max_export_batch_size"`
	MaxQueueSize       int    `json:"max_queue_size"        yaml:"max_queue_size"`
}

// Bind binds the actual data from the OTLPAttribute receiver fields.
func (o *OTLPAttribute) Bind() *OTLPAttribute {
	o.Namespace = GetActualValue(o.Namespace)
	o.PodName = GetActualValue(o.PodName)
	o.NodeName = GetActualValue(o.NodeName)
	o.ServiceName = GetActualValue(o.ServiceName)
	return o
}

// Bind binds the actual data from the Observability receiver fields.
func (o *Observability) Bind() *Observability {
	if o.OTLP != nil {
		o.OTLP.CollectorEndpoint = GetActualValue(o.OTLP.CollectorEndpoint)
		o.OTLP.TraceBatchTimeout = GetActualValue(o.OTLP.TraceBatchTimeout)
		o.OTLP.TraceExportTimeout = GetActualValue(o.OTLP.TraceExportTimeout)
		o.OTLP.MetricsExportInterval = GetActualValue(o.OTLP.MetricsExportInterval)
		o.OTLP.MetricsExportTimeout = GetActualValue(o.OTLP.MetricsExportTimeout)
	} else {
		o.OTLP = new(OTLP)
	}

	if o.OTLP.Attribute != nil {
		o.OTLP.Attribute.Bind()
	} else {
		o.OTLP.Attribute = new(OTLPAttribute)
	}

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
	} else {
		o.Prometheus = new(Prometheus)
	}

	if o.Jaeger != nil {
		o.Jaeger.AgentEndpoint = GetActualValue(o.Jaeger.AgentEndpoint)
		o.Jaeger.AgentReconnectInterval = GetActualValue(o.Jaeger.AgentReconnectInterval)
		o.Jaeger.CollectorEndpoint = GetActualValue(o.Jaeger.CollectorEndpoint)
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
