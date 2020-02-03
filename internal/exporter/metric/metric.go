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

// Package metric provides meters.
package metric

import (
	"strings"
	"unsafe"

	"github.com/vdaas/vald/internal/log"
	"go.opentelemetry.io/otel/exporter/metric/stdout"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
)

type stdoutWriter struct {
	printer func(...interface{})
}

func (t *stdoutWriter) Write(p []byte) (int, error) {
	t.printer(strings.TrimSpace(*(*string)(unsafe.Pointer(&p))))
	return len(p), nil
}

func InitStdoutMeter() (*push.Controller, error) {
	return stdout.InstallNewPipeline(
		stdout.Config{
			Writer:      &stdoutWriter{printer: log.Info},
			Quantiles:   []float64{0.5, 0.9, 0.99},
			PrettyPrint: false,
		})
}
