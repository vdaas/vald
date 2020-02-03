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

// Package jaeger provides a jaeger exporter.
package jaeger

import (
	"sync"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

type jaegerOptions = jaeger.Options

var (
	exporter *jaeger.Exporter
	once     sync.Once
)

func Init(opts ...JaegerOption) (err error) {
	jo := new(jaegerOptions)

	for _, opt := range append(jaegerDefaultOpts, opts...) {
		err = opt(jo)
		if err != nil {
			return err
		}
	}

	ex, err := jaeger.NewExporter(*jo)
	if err != nil {
		return err
	}

	once.Do(func() {
		exporter = ex
		trace.RegisterExporter(ex)
	})
	return nil
}

func Exporter() *jaeger.Exporter {
	return exporter
}
