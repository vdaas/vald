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
		Enabled bool
		OTLP    *OTLP
		Metrics *Metrics
		Trace   *Trace
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
						OTLP: &OTLP{
							Attribute: new(OTLPAttribute),
						},
						Metrics: new(Metrics),
						Trace:   new(Trace),
					},
				},
			}
		}(),
		func() test {
			collectorEndpoint := "collector.monitoring.svc.cluster.local:6831"
			traceMaxExportBatchSize := 256
			traceMaxQueueSize := 100
			traceBatchTimeout := "1m"
			traceExportTimeout := "1s"
			metricsExportInterval := "10ms"
			metricsExportTimeout := "30s"
			return test{
				name: "return Observability when all object parameters are not nil",
				fields: fields{
					Enabled: false,
					Metrics: new(Metrics),
					Trace:   new(Trace),
					OTLP: &OTLP{
						CollectorEndpoint:       collectorEndpoint,
						TraceBatchTimeout:       traceBatchTimeout,
						TraceExportTimeout:      traceExportTimeout,
						TraceMaxExportBatchSize: traceMaxExportBatchSize,
						TraceMaxQueueSize:       traceMaxQueueSize,
						MetricsExportInterval:   metricsExportInterval,
						MetricsExportTimeout:    metricsExportTimeout,
					},
				},
				want: want{
					want: &Observability{
						Enabled: false,
						Metrics: new(Metrics),
						Trace:   new(Trace),
						OTLP: &OTLP{
							CollectorEndpoint:       collectorEndpoint,
							TraceBatchTimeout:       traceBatchTimeout,
							TraceExportTimeout:      traceExportTimeout,
							TraceMaxExportBatchSize: traceMaxExportBatchSize,
							TraceMaxQueueSize:       traceMaxQueueSize,
							MetricsExportInterval:   metricsExportInterval,
							MetricsExportTimeout:    metricsExportTimeout,
							Attribute:               new(OTLPAttribute),
						},
					},
				},
			}
		}(),
		func() test {
			collectorEndpoint := "collector.monitoring.svc.cluster.local:6831"
			traceBatchTimeout := "1m"
			traceExportTimeout := "1s"
			metricsExportInterval := "10ms"
			metricsExportTimeout := "30s"

			envPrefix := "OBSERVABILITY_BIND_"
			m := map[string]string{
				envPrefix + "COLLECTOR_ENDPOINT":      collectorEndpoint,
				envPrefix + "TRACE_BATCH_TIMEOUT":     traceBatchTimeout,
				envPrefix + "TRACE_EXPORT_TIMEOUT":    traceExportTimeout,
				envPrefix + "METRICS_EXPORT_INTERVAL": metricsExportInterval,
				envPrefix + "METRICS_EXPORT_TIMEOUT":  metricsExportTimeout,
			}
			return test{
				name: "return Observability when the data is loaded environment variable",
				fields: fields{
					Enabled: false,
					Metrics: new(Metrics),
					Trace:   new(Trace),
					OTLP: &OTLP{
						CollectorEndpoint:     "_" + envPrefix + "COLLECTOR_ENDPOINT_",
						TraceBatchTimeout:     "_" + envPrefix + "TRACE_BATCH_TIMEOUT_",
						TraceExportTimeout:    "_" + envPrefix + "TRACE_EXPORT_TIMEOUT_",
						MetricsExportInterval: "_" + envPrefix + "METRICS_EXPORT_INTERVAL_",
						MetricsExportTimeout:  "_" + envPrefix + "METRICS_EXPORT_TIMEOUT_",
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
						Metrics: new(Metrics),
						Trace:   new(Trace),
						OTLP: &OTLP{
							CollectorEndpoint:     collectorEndpoint,
							TraceBatchTimeout:     traceBatchTimeout,
							TraceExportTimeout:    traceExportTimeout,
							MetricsExportInterval: metricsExportInterval,
							MetricsExportTimeout:  metricsExportTimeout,
							Attribute:             new(OTLPAttribute),
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
				Enabled: test.fields.Enabled,
				Trace:   test.fields.Trace,
				OTLP:    test.fields.OTLP,
			}

			got := o.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestOTLPAttribute_Bind(t *testing.T) {
	type fields struct {
		Namespace   string
		PodName     string
		NodeName    string
		ServiceName string
	}
	type want struct {
		want *OTLPAttribute
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *OTLPAttribute) error
		beforeFunc func(*testing.T)
		afterFunc  func(*testing.T)
	}
	defaultCheckFunc := func(w want, got *OTLPAttribute) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name:   "return OTLPAttribute when all object parameters are nil",
				fields: fields{},
				want: want{
					want: new(OTLPAttribute),
				},
			}
		}(),
		func() test {
			namespace := "monitoring"
			podName := "vald-agent-ngt-0"
			nodeName := "k3d-vald-cluster-agent-0"
			servicename := "vald-agent-ngt"
			return test{
				name: "return OTLPAttribute when all object parameters are not nil",
				fields: fields{
					Namespace:   namespace,
					PodName:     podName,
					NodeName:    nodeName,
					ServiceName: servicename,
				},
				want: want{
					want: &OTLPAttribute{
						Namespace:   namespace,
						PodName:     podName,
						NodeName:    nodeName,
						ServiceName: servicename,
					},
				},
			}
		}(),
		func() test {
			namespace := "monitoring"
			podName := "vald-agent-ngt-0"
			nodeName := "k3d-vald-cluster-agent-0"
			servicename := "vald-agent-ngt"

			envPrefix := "OTLPAttribute_BIND_"
			m := map[string]string{
				envPrefix + "NAMESPACE":    namespace,
				envPrefix + "POD_NAME":     podName,
				envPrefix + "NODE_NAME":    nodeName,
				envPrefix + "SERVICE_NAME": servicename,
			}

			return test{
				name: "return OTLPAttribute when the data is loaded environment variable",
				fields: fields{
					Namespace:   "_" + envPrefix + "NAMESPACE_",
					PodName:     "_" + envPrefix + "POD_NAME_",
					NodeName:    "_" + envPrefix + "NODE_NAME_",
					ServiceName: "_" + envPrefix + "SERVICE_NAME_",
				},
				beforeFunc: func(t *testing.T) {
					t.Helper()
					for k, v := range m {
						t.Setenv(k, v)
					}
				},
				want: want{
					want: &OTLPAttribute{
						Namespace:   namespace,
						PodName:     podName,
						NodeName:    nodeName,
						ServiceName: servicename,
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
			o := &OTLPAttribute{
				Namespace:   test.fields.Namespace,
				PodName:     test.fields.PodName,
				NodeName:    test.fields.NodeName,
				ServiceName: test.fields.ServiceName,
			}

			got := o.Bind()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
