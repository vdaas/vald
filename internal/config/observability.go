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

// Package config providers configuration type and load configuration logic
package config

type Observability struct {
	Enabled     bool         `json:"enabled" yaml:"enabled"`
	Collector   *Collector   `json:"collector" yaml:"collector"`
	Trace       *Trace       `json:"trace" yaml:"trace"`
	Prometheus  *Prometheus  `json:"prometheus" yaml:"prometheus"`
	Jaeger      *Jaeger      `json:"jaeger" yaml:"jaeger"`
	Stackdriver *Stackdriver `json:"stackdriver" yaml:"stackdriver"`
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
	EnableVersionInfo bool     `json:"enable_version_info" yaml:"enable_version_info"`
	VersionInfoLabels []string `json:"version_info_labels" yaml:"version_info_labels"`
	EnableMemory      bool     `json:"enable_memory" yaml:"enable_memory"`
	EnableGoroutine   bool     `json:"enable_goroutine" yaml:"enable_goroutine"`
	EnableCGO         bool     `json:"enable_cgo" yaml:"enable_cgo"`
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

type Stackdriver struct {
	ProjectID string `json:"project_id" yaml:"project_id"`

	Client *StackdriverClient `json:"client" yaml:"client"`

	Exporter *StackdriverExporter `json:"exporter" yaml:"exporter"`
	Profiler *StackdriverProfiler `json:"profiler" yaml:"profiler"`
}

type StackdriverClient struct {
	APIKey                string   `json:"api_key" yaml:"api_key"`
	Audiences             []string `json:"audiences" yaml:"audiences"`
	CredentialsFile       string   `json:"credentials_file" yaml:"credentials_file"`
	CredentialsJSON       string   `json:"credentials_json" yaml:"credentials_json"`
	Endpoint              string   `json:"endpoint" yaml:"endpoint"`
	QuotaProject          string   `json:"quota_project" yaml:"quota_project"`
	RequestReason         string   `json:"request_reason" yaml:"request_reason"`
	Scopes                []string `json:"scopes" yaml:"scopes"`
	UserAgent             string   `json:"user_agent" yaml:"user_agent"`
	TelemetryEnabled      bool     `json:"telemetry_enabled" yaml:"telemetry_enabled"`
	AuthenticationEnabled bool     `json:"authentication_enabled" yaml:"authentication_enabled"`
}

type StackdriverExporter struct {
	MonitoringEnabled bool `json:"monitoring_enabled" yaml:"monitoring_enabled"`
	TracingEnabled    bool `json:"tracing_enabled" yaml:"tracing_enabled"`

	Location                 string `json:"location" yaml:"location"`
	BundleDelayThreshold     string `json:"bundle_delay_threshold" yaml:"bundle_delay_threshold"`
	BundleCountThreshold     int    `json:"bundle_count_threshold" yaml:"bundle_count_threshold"`
	TraceSpansBufferMaxBytes int    `json:"trace_spans_buffer_max_bytes" yaml:"trace_spans_buffer_max_bytes"`

	MetricPrefix string `json:"metric_prefix" yaml:"metric_prefix"`

	SkipCMD           bool   `json:"skip_cmd" yaml:"skip_cmd"`
	Timeout           string `json:"timeout" yaml:"timeout"`
	ReportingInterval string `json:"reporting_interval" yaml:"reporting_interval"`
	NumberOfWorkers   int    `json:"number_of_workers" yaml:"number_of_workers"`
}

type StackdriverProfiler struct {
	Enabled        bool   `json:"enabled" yaml:"enabled"`
	Service        string `json:"service" yaml:"service"`
	ServiceVersion string `json:"service_version" yaml:"service_version"`
	DebugLogging   bool   `json:"debug_logging" yaml:"debug_logging"`

	MutexProfiling     bool `json:"mutex_profiling" yaml:"mutex_profiling"`
	CPUProfiling       bool `json:"cpu_profiling" yaml:"cpu_profiling"`
	AllocProfiling     bool `json:"alloc_profiling" yaml:"alloc_profiling"`
	HeapProfiling      bool `json:"heap_profiling" yaml:"heap_profiling"`
	GoroutineProfiling bool `json:"goroutine_profiling" yaml:"goroutine_profiling"`

	AllocForceGC bool `json:"alloc_force_gc" yaml:"alloc_force_gc"`

	APIAddr string `json:"api_addr" yaml:"api_addr"`

	Instance string `json:"instance" yaml:"instance"`
	Zone     string `json:"zone" yaml:"zone"`
}

func (o *Observability) Bind() *Observability {
	if o.Collector != nil {
		o.Collector = o.Collector.Bind()
	} else {
		o.Collector = new(Collector)
		o.Collector.Metrics = new(Metrics)
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

	if o.Stackdriver != nil {
		o.Stackdriver = o.Stackdriver.Bind()
	} else {
		o.Stackdriver = new(Stackdriver)
		o.Stackdriver.Client = new(StackdriverClient)
		o.Stackdriver.Exporter = new(StackdriverExporter)
		o.Stackdriver.Profiler = new(StackdriverProfiler)
	}

	return o
}

func (c *Collector) Bind() *Collector {
	c.Duration = GetActualValue(c.Duration)

	if c.Metrics != nil {
		c.Metrics.VersionInfoLabels = GetActualValues(c.Metrics.VersionInfoLabels)
	} else {
		c.Metrics = new(Metrics)
	}

	return c
}

func (sd *Stackdriver) Bind() *Stackdriver {
	sd.ProjectID = GetActualValue(sd.ProjectID)

	if sd.Client != nil {
		sd.Client.APIKey = GetActualValue(sd.Client.APIKey)
		sd.Client.Audiences = GetActualValues(sd.Client.Audiences)
		sd.Client.CredentialsFile = GetActualValue(sd.Client.CredentialsFile)
		sd.Client.CredentialsJSON = GetActualValue(sd.Client.CredentialsJSON)
		sd.Client.Endpoint = GetActualValue(sd.Client.Endpoint)
		sd.Client.QuotaProject = GetActualValue(sd.Client.QuotaProject)
		sd.Client.RequestReason = GetActualValue(sd.Client.RequestReason)
		sd.Client.Scopes = GetActualValues(sd.Client.Scopes)
		sd.Client.UserAgent = GetActualValue(sd.Client.UserAgent)
	} else {
		sd.Client = new(StackdriverClient)
	}

	if sd.Exporter != nil {
		sd.Exporter.Location = GetActualValue(sd.Exporter.Location)
		sd.Exporter.BundleDelayThreshold = GetActualValue(sd.Exporter.BundleDelayThreshold)
		sd.Exporter.MetricPrefix = GetActualValue(sd.Exporter.MetricPrefix)
		sd.Exporter.Timeout = GetActualValue(sd.Exporter.Timeout)
		sd.Exporter.ReportingInterval = GetActualValue(sd.Exporter.ReportingInterval)
	} else {
		sd.Exporter = new(StackdriverExporter)
	}

	if sd.Profiler != nil {
		sd.Profiler.Service = GetActualValue(sd.Profiler.Service)
		sd.Profiler.ServiceVersion = GetActualValue(sd.Profiler.ServiceVersion)
		sd.Profiler.APIAddr = GetActualValue(sd.Profiler.APIAddr)
		sd.Profiler.Instance = GetActualValue(sd.Profiler.Instance)
		sd.Profiler.Zone = GetActualValue(sd.Profiler.Zone)
	} else {
		sd.Profiler = new(StackdriverProfiler)
	}

	return sd
}
