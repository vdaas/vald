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

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestObservability_Bind(t *testing.T) {
	type fields struct {
		Enabled    bool
		Collector  *Collector
		Trace      *Trace
		Prometheus *Prometheus
		Jaeger     *Jaeger
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
					Enabled: true,
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
				},
				want: want{
					want: &Observability{
						Enabled: false,
						Collector: &Collector{
							Duration: collectorDur,
							Metrics:  new(Metrics),
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

			envPrefix := "OBSERVABILITY_BIND_"
			m := map[string]string{
				envPrefix + "PROMETHEUS_ENDPOINT":       prometheusEndpoint,
				envPrefix + "PROMETHUS_NAMESPACE":       prometheusNamespace,
				envPrefix + "JAEGER_COLLECTOR_ENDPOINT": jaegerCollectorEndpoint,
				envPrefix + "JAEGER_AGENT_ENDPOINT":     jaegerAgentEndpoint,
				envPrefix + "JAEGER_USERNAME":           jaegerUsername,
				envPrefix + "JAEGER_PASSWORD":           jaegerPassword,
				envPrefix + "JAEGER_SERVICE_NAME":       jaegerServiceName,
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
						Endpoint:  "_" + envPrefix + "PROMETHEUS_ENDPOINT_",
						Namespace: "_" + envPrefix + "PROMETHUS_NAMESPACE_",
					},
					Jaeger: &Jaeger{
						CollectorEndpoint: "_" + envPrefix + "JAEGER_COLLECTOR_ENDPOINT_",
						AgentEndpoint:     "_" + envPrefix + "JAEGER_AGENT_ENDPOINT_",
						Username:          "_" + envPrefix + "JAEGER_USERNAME_",
						Password:          "_" + envPrefix + "JAEGER_PASSWORD_",
						ServiceName:       "_" + envPrefix + "JAEGER_SERVICE_NAME_",
					},
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						t.Setenv(k, v)
					}
				},
				want: want{
					want: &Observability{
						Enabled: false,
						Collector: &Collector{
							Duration: collectorDur,
							Metrics:  new(Metrics),
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
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			o := &Observability{
				Enabled:    test.fields.Enabled,
				Collector:  test.fields.Collector,
				Trace:      test.fields.Trace,
				Prometheus: test.fields.Prometheus,
				Jaeger:     test.fields.Jaeger,
			}

			got := o.Bind()
			if err := checkFunc(test.want, got); err != nil {
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

			envPrefix := "COLLECTOR_BIND_"
			m := map[string]string{
				envPrefix + "DURATION":                    duration,
				envPrefix + "METRICS_VERSION_INFO_LABELS": versionInfoLabels,
			}
			return test{
				name: "return Collector when the data is loaded from the environment variable",
				fields: fields{
					Duration: "_" + envPrefix + "DURATION_",
					Metrics: &Metrics{
						VersionInfoLabels: []string{
							"_" + envPrefix + "METRICS_VERSION_INFO_LABELS_",
						},
					},
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						t.Setenv(k, v)
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
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			c := &Collector{
				Duration: test.fields.Duration,
				Metrics:  test.fields.Metrics,
			}

			got := c.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
