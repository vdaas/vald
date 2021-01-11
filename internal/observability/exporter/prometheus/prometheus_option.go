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

// Package prometheus provides a prometheus exporter.
package prometheus

import "github.com/vdaas/vald/internal/log"

type PrometheusOption func(*prometheusOptions) error

var prometheusDefaultOpts = []PrometheusOption{
	WithEndpoint("/metrics"),
	WithNamespace("vald"),
	WithOnErrorFunc(func(err error) {
		if err != nil {
			log.Warnf("Failed to export to Prometheus: %v", err)
		}
	}),
}

func WithEndpoint(ep string) PrometheusOption {
	return func(po *prometheusOptions) error {
		if ep != "" {
			po.endpoint = ep
		}
		return nil
	}
}

func WithNamespace(ns string) PrometheusOption {
	return func(po *prometheusOptions) error {
		if ns != "" {
			po.options.Namespace = ns
		}
		return nil
	}
}

func WithOnErrorFunc(f func(error)) PrometheusOption {
	return func(po *prometheusOptions) error {
		if f != nil {
			po.options.OnError = f
		}
		return nil
	}
}
