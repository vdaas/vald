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

import (
	"os"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestObservability_Bind(t *testing.T) {
	type fields struct {
		Enabled     bool
		Collector   *Collector
		Trace       *Trace
		Prometheus  *Prometheus
		Jaeger      *Jaeger
		Stackdriver *Stackdriver
	}
	type want struct {
		want *Observability
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Observability) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Observability) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return Observability when all object parameters are nil",
				fields: fields{
					Enabled: false,
				},
				want: want{
					want: &Observability{
						Enabled: true,
						Collector: &Collector{
							Metrics: new(Metrics),
						},
						Trace:      new(Trace),
						Prometheus: new(Prometheus),
						Jaeger:     new(Jaeger),
						Stackdriver: &Stackdriver{
							Client:   new(StackdriverClient),
							Exporter: new(StackdriverExporter),
							Profiler: new(StackdriverProfiler),
						},
					},
				},
			}
		}(),
		func() test {
			collectorDur := "5ms"
			prometheusEndpoint := "http://prometheus.kube-system.svc.cluster.local.:9090"
			prometheusNamespace := "monitoring"
			jaegerCollectorEndpoint := "http://jaeger-collector.monitoring.svc.cluster.local:14268/api/traces"
			jaegerAgentEndpoint := "jaeger-agent.monitoring.svc.cluster.local:6831"
			jaegerUsername := "username"
			jaegerPassword := "pass"
			jaegerServiceName := "jaeger"
			stackdriverProjectID := "vald"
			return test{
				name: "return Observability when all object parameters are not nil",
				fields: fields{
					Enabled: false,
					Collector: &Collector{
						Duration: collectorDur,
					},
					Trace: new(Trace),
					Prometheus: &Prometheus{
						Endpoint:  prometheusEndpoint,
						Namespace: prometheusNamespace,
					},
					Jaeger: &Jaeger{
						CollectorEndpoint: jaegerCollectorEndpoint,
						AgentEndpoint:     jaegerAgentEndpoint,
						Username:          jaegerUsername,
						Password:          jaegerPassword,
						ServiceName:       jaegerServiceName,
					},
					Stackdriver: &Stackdriver{
						ProjectID: stackdriverProjectID,
						Client:    new(StackdriverClient),
						Exporter:  new(StackdriverExporter),
						Profiler:  new(StackdriverProfiler),
					},
				},
				want: want{
					want: &Observability{
						Enabled: false,
						Collector: &Collector{
							Duration: collectorDur,
						},
						Trace: new(Trace),
						Prometheus: &Prometheus{
							Endpoint:  prometheusEndpoint,
							Namespace: prometheusNamespace,
						},
						Jaeger: &Jaeger{
							CollectorEndpoint: jaegerCollectorEndpoint,
							AgentEndpoint:     jaegerAgentEndpoint,
							Username:          jaegerUsername,
							Password:          jaegerPassword,
							ServiceName:       jaegerServiceName,
						},
						Stackdriver: &Stackdriver{
							ProjectID: stackdriverProjectID,
							Client:    new(StackdriverClient),
							Exporter:  new(StackdriverExporter),
							Profiler:  new(StackdriverProfiler),
						},
					},
				},
			}
		}(),
		func() test {
			collectorDur := "5ms"
			prometheusEndpoint := "http://prometheus.kube-system.svc.cluster.local.:9090"
			prometheusNamespace := "monitoring"
			jaegerCollectorEndpoint := "http://jaeger-collector.monitoring.svc.cluster.local:14268/api/traces"
			jaegerAgentEndpoint := "jaeger-agent.monitoring.svc.cluster.local:6831"
			jaegerUsername := "username"
			jaegerPassword := "pass"
			jaegerServiceName := "jaeger"
			stackdriverProjectID := "vald"
			m := map[string]string{
				"PROMETHEUS_ENDPOINT":       prometheusEndpoint,
				"PROMETHUS_NAMESPACE":       prometheusNamespace,
				"JAEGER_COLLECTOR_ENDPOINT": jaegerCollectorEndpoint,
				"JAEGER_AGENT_ENDPOINT":     jaegerAgentEndpoint,
				"JAEGER_USERNAME":           jaegerUsername,
				"JAEGER_PASSWORD":           jaegerPassword,
				"JAEGER_SERVICE_NAME":       jaegerServiceName,
				"STACKDRIVER_PROJECT_ID":    stackdriverProjectID,
			}
			return test{
				name: "return Observability when the data is loaded environment variable",
				fields: fields{
					Enabled: false,
					Collector: &Collector{
						Duration: collectorDur,
					},
					Trace: new(Trace),
					Prometheus: &Prometheus{
						Endpoint:  "_PROMETHEUS_ENDPOINT_",
						Namespace: "_PROMETHUS_NAMESPACE_",
					},
					Jaeger: &Jaeger{
						CollectorEndpoint: "_JAEGER_COLLECTOR_ENDPOINT_",
						AgentEndpoint:     "_JAEGER_AGENT_ENDPOINT_",
						Username:          "_JAEGER_USERNAME_",
						Password:          "_JAEGER_PASSWORD_",
						ServiceName:       "_JAEGER_SERVICE_NAME_",
					},
					Stackdriver: &Stackdriver{
						ProjectID: "_STACKDRIVER_PROJECT_ID_",
						Client:    new(StackdriverClient),
						Exporter:  new(StackdriverExporter),
						Profiler:  new(StackdriverProfiler),
					},
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						if err := os.Setenv(k, v); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for k := range m {
						if err := os.Unsetenv(k); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &Observability{
						Enabled: false,
						Collector: &Collector{
							Duration: collectorDur,
						},
						Trace: new(Trace),
						Prometheus: &Prometheus{
							Endpoint:  prometheusEndpoint,
							Namespace: prometheusNamespace,
						},
						Jaeger: &Jaeger{
							CollectorEndpoint: jaegerCollectorEndpoint,
							AgentEndpoint:     jaegerAgentEndpoint,
							Username:          jaegerUsername,
							Password:          jaegerPassword,
							ServiceName:       jaegerServiceName,
						},
						Stackdriver: &Stackdriver{
							ProjectID: stackdriverProjectID,
							Client:    new(StackdriverClient),
							Exporter:  new(StackdriverExporter),
							Profiler:  new(StackdriverProfiler),
						},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			o := &Observability{
				Enabled:     test.fields.Enabled,
				Collector:   test.fields.Collector,
				Trace:       test.fields.Trace,
				Prometheus:  test.fields.Prometheus,
				Jaeger:      test.fields.Jaeger,
				Stackdriver: test.fields.Stackdriver,
			}

			got := o.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestCollector_Bind(t *testing.T) {
	type fields struct {
		Duration string
		Metrics  *Metrics
	}
	type want struct {
		want *Collector
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Collector) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Collector) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return Collector when the Metrics is nil",
				fields: fields{
					Duration: "5ms",
				},
				want: want{
					want: &Collector{
						Duration: "5ms",
						Metrics:  new(Metrics),
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return Collector when the Metrics is not nil",
				fields: fields{
					Duration: "5ms",
					Metrics:  new(Metrics),
				},
				want: want{
					want: &Collector{
						Duration: "5ms",
						Metrics:  new(Metrics),
					},
				},
			}
		}(),
		func() test {
			duration := "5ms"
			versionInfoLabels := "vald_version"
			m := map[string]string{
				"DURATION":                    duration,
				"METRICS_VERSION_INFO_LABELS": versionInfoLabels,
			}
			return test{
				name: "return Collector when the data is loaded from the environment variable",
				fields: fields{
					Duration: "_DURATION_",
					Metrics: &Metrics{
						VersionInfoLabels: []string{
							"_METRICS_VERSION_INFO_LABELS_",
						},
					},
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						if err := os.Setenv(k, v); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for k := range m {
						if err := os.Unsetenv(k); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &Collector{
						Duration: duration,
						Metrics: &Metrics{
							VersionInfoLabels: []string{
								versionInfoLabels,
							},
						},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &Collector{
				Duration: test.fields.Duration,
				Metrics:  test.fields.Metrics,
			}

			got := c.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestStackdriver_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		ProjectID string
		Client    *StackdriverClient
		Exporter  *StackdriverExporter
		Profiler  *StackdriverProfiler
	}
	type want struct {
		want *Stackdriver
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Stackdriver) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *Stackdriver) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			projectID := "vald"
			return test{
				name: "return Stackdriver when the Client and the Exporter and the Profiler is nil",
				fields: fields{
					ProjectID: projectID,
				},
				want: want{
					want: &Stackdriver{
						ProjectID: projectID,
						Client:    new(StackdriverClient),
						Exporter:  new(StackdriverExporter),
						Profiler:  new(StackdriverProfiler),
					},
				},
			}
		}(),
		func() test {
			projectID := "vald"
			return test{
				name: "return Stackdriver when the Client and the Exporter and the Profiler is not nil",
				fields: fields{
					ProjectID: projectID,
					Client:    new(StackdriverClient),
					Exporter:  new(StackdriverExporter),
					Profiler:  new(StackdriverProfiler),
				},
				want: want{
					want: &Stackdriver{
						ProjectID: projectID,
						Client:    new(StackdriverClient),
						Exporter:  new(StackdriverExporter),
						Profiler:  new(StackdriverProfiler),
					},
				},
			}
		}(),
		func() test {
			projectID := "vdaas/vald"
			clientAPIKey := "api_key"
			exporterLocation := "asia-northeast1-a"
			profileService := "vald-service"
			m := map[string]string{
				"PROJECT_ID":        projectID,
				"CLIENT_API_KEY":    clientAPIKey,
				"EXPORTER_LOCATION": exporterLocation,
				"PROFILER_SERVICE":  profileService,
			}
			return test{
				name: "return Stackdriver when the data is loaded from the environment variable",
				fields: fields{
					ProjectID: "_PROJECT_ID_",
					Client: &StackdriverClient{
						APIKey: "_API_KEY_",
					},
					Exporter: &StackdriverExporter{
						Location: "_EXPORTER_LOCATION_",
					},
					Profiler: &StackdriverProfiler{
						Service: "_PROFILER_SERVICE_",
					},
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						if err := os.Setenv(k, v); err != nil {
							t.Fatal(err)
						}
					}
				},
				afterFunc: func(t *testing.T) {
					t.Helper()
					for k := range m {
						if err := os.Unsetenv(k); err != nil {
							t.Fatal(err)
						}
					}
				},
				want: want{
					want: &Stackdriver{
						ProjectID: projectID,
						Client: &StackdriverClient{
							APIKey: clientAPIKey,
						},
						Exporter: &StackdriverExporter{
							Location: exporterLocation,
						},
						Profiler: &StackdriverProfiler{
							Service: profileService,
						},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			sd := &Stackdriver{
				ProjectID: test.fields.ProjectID,
				Client:    test.fields.Client,
				Exporter:  test.fields.Exporter,
				Profiler:  test.fields.Profiler,
			}

			got := sd.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
