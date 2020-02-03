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

// Package prometheus provides a prometheus exporter.
package prometheus

import (
	"sync"

	"contrib.go.opencensus.io/exporter/prometheus"
)

type prometheusOptions = prometheus.Options

var (
	exporter *prometheus.Exporter
	once     sync.Once
)

func Init(opts ...PrometheusOption) (err error) {
	po := new(prometheusOptions)

	for _, opt := range append(prometheusDefaultOpts, opts...) {
		err = opt(po)
		if err != nil {
			return err
		}
	}

	ex, err := prometheus.NewExporter(*po)
	if err != nil {
		return err
	}

	once.Do(func() {
		exporter = ex
	})
	return nil
}

func Exporter() *prometheus.Exporter {
	return exporter
}
