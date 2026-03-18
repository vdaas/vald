//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package config

// Observability represents the configuration for the observability.
type Observability struct {
	// OTLP represents the OTLP configuration.
	OTLP *OTLP `json:"otlp" yaml:"otlp"`
	// Metrics represents the metrics configuration.
	Metrics *Metrics `json:"metrics" yaml:"metrics"`
	// Trace represents the trace configuration.
	Trace *Trace `json:"trace" yaml:"trace"`
	// Enabled enables observability.
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// Bind binds the actual data from the Observability receiver fields.
func (o *Observability) Bind() *Observability {
	if o.OTLP == nil {
		o.OTLP = new(OTLP)
	}
	if o.OTLP != nil {
		o.OTLP.Bind()
	}

	if o.Metrics == nil {
		o.Metrics = new(Metrics)
	}
	if o.Metrics != nil {
		o.Metrics.Bind()
	}

	if o.Trace == nil {
		o.Trace = new(Trace)
	}
	if o.Trace != nil {
		o.Trace.Bind()
	}

	return o
}

// OTLP represents the configuration for the OTLP.
type OTLP struct {
	// Attribute represents the OTLP attribute configuration.
	Attribute *OTLPAttribute `json:"attribute" yaml:"attribute"`
	// CollectorEndpoint represents the collector endpoint.
	CollectorEndpoint string `json:"collector_endpoint" yaml:"collector_endpoint"`
	// TraceBatchTimeout represents the trace batch timeout.
	TraceBatchTimeout string `json:"trace_batch_timeout" yaml:"trace_batch_timeout"`
	// TraceExportTimeout represents the trace export timeout.
	TraceExportTimeout string `json:"trace_export_timeout" yaml:"trace_export_timeout"`
	// MetricsExportInterval represents the metrics export interval.
	MetricsExportInterval string `json:"metrics_export_interval" yaml:"metrics_export_interval"`
	// MetricsExportTimeout represents the metrics export timeout.
	MetricsExportTimeout string `json:"metrics_export_timeout" yaml:"metrics_export_timeout"`
	// TraceMaxExportBatchSize represents the trace max export batch size.
	TraceMaxExportBatchSize int `json:"trace_max_export_batch_size" yaml:"trace_max_export_batch_size"`
	// TraceMaxQueueSize represents the trace max queue size.
	TraceMaxQueueSize int `json:"trace_max_queue_size" yaml:"trace_max_queue_size"`
}

// Bind binds the actual data from the OTLP receiver fields.
func (o *OTLP) Bind() *OTLP {
	o.CollectorEndpoint = GetActualValue(o.CollectorEndpoint)
	o.TraceBatchTimeout = GetActualValue(o.TraceBatchTimeout)
	o.TraceExportTimeout = GetActualValue(o.TraceExportTimeout)
	o.MetricsExportInterval = GetActualValue(o.MetricsExportInterval)
	o.MetricsExportTimeout = GetActualValue(o.MetricsExportTimeout)
	if o.Attribute == nil {
		o.Attribute = new(OTLPAttribute)
	}
	if o.Attribute != nil {
		o.Attribute.Bind()
	}
	return o
}

// OTLPAttribute represents the configuration for the OTLP attribute.
type OTLPAttribute struct {
	// Namespace represents the namespace.
	Namespace string `json:"namespace" yaml:"namespace"`
	// PodName represents the pod name.
	PodName string `json:"pod_name" yaml:"pod_name"`
	// NodeName represents the node name.
	NodeName string `json:"node_name" yaml:"node_name"`
	// ServiceName represents the service name.
	ServiceName string `json:"service_name" yaml:"service_name"`
}

// Bind binds the actual data from the OTLPAttribute receiver fields.
func (o *OTLPAttribute) Bind() *OTLPAttribute {
	o.Namespace = GetActualValue(o.Namespace)
	o.PodName = GetActualValue(o.PodName)
	o.NodeName = GetActualValue(o.NodeName)
	o.ServiceName = GetActualValue(o.ServiceName)
	return o
}

// Trace represents the configuration for the trace.
type Trace struct {
	// Enabled enables trace.
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// Bind binds the actual data from the Trace receiver fields.
func (t *Trace) Bind() *Trace {
	// No fields to bind as per rules
	return t
}

// Metrics represents the configuration for the metrics.
type Metrics struct {
	// VersionInfoLabels represents the version info labels.
	VersionInfoLabels []string `json:"version_info_labels" yaml:"version_info_labels"`
	// EnableVersionInfo enables version info metrics.
	EnableVersionInfo bool `json:"enable_version_info" yaml:"enable_version_info"`
	// EnableMemory enables memory metrics.
	EnableMemory bool `json:"enable_memory" yaml:"enable_memory"`
	// EnableGoroutine enables goroutine metrics.
	EnableGoroutine bool `json:"enable_goroutine" yaml:"enable_goroutine"`
	// EnableCGO enables cgo metrics.
	EnableCGO bool `json:"enable_cgo" yaml:"enable_cgo"`
}

// Bind binds the actual data from the Metrics receiver fields.
func (m *Metrics) Bind() *Metrics {
	m.VersionInfoLabels = GetActualValues(m.VersionInfoLabels)
	return m
}
