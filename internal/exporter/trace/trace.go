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

// Package trace provides tracers.
package trace

import (
	"strings"
	"unsafe"

	"github.com/vdaas/vald/internal/log"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporter/trace/stdout"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type stdoutWriter struct {
	printer func(...interface{})
}

func (t *stdoutWriter) Write(p []byte) (int, error) {
	t.printer(strings.TrimSpace(*(*string)(unsafe.Pointer(&p))))
	return len(p), nil
}

func InitStdoutTracer() error {
	exporter, err := stdout.NewExporter(
		stdout.Options{
			Writer:      &stdoutWriter{printer: log.Info},
			PrettyPrint: false,
		})
	if err != nil {
		return err
	}

	tp, err := sdktrace.NewProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter),
	)
	if err != nil {
		return err
	}

	global.SetTraceProvider(tp)
	return nil
}
