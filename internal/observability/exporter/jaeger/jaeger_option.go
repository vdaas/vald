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
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"go.opentelemetry.io/otel/exporter/trace/jaeger"
)

type Option func(*exp) error

var (
	defaultOpts = []Option{
		WithServiceName("vald"),
		WithOnErrorFunc(func(err error) {
			if err != nil {
				log.Warnf("Error when uploading spans to Jaeger: %v", err)
			}
		}),
	}
)

func WithCollectorEndpoint(endpoint string) Option {
	return func(e *exp) error {
		if endpoint == "" {
			return nil
		}

		e.collectorEndpoint = endpoint

		return nil
	}
}

func WithAgentEndpoint(endpoint string) Option {
	return func(e *exp) error {
		if endpoint == "" {
			return nil
		}

		e.agentEndpoint = endpoint

		return nil
	}
}

func WithUsername(username string) Option {
	return func(e *exp) error {
		if username == "" {
			return nil
		}

		if e.collectorOpts == nil {
			e.collectorOpts = []jaeger.CollectorEndpointOption{
				jaeger.WithUsername(username),
			}
		} else {
			e.collectorOpts = append(
				e.collectorOpts,
				jaeger.WithUsername(username),
			)
		}

		return nil
	}
}

func WithPassword(password string) Option {
	return func(e *exp) error {
		if password == "" {
			return nil
		}

		if e.collectorOpts == nil {
			e.collectorOpts = []jaeger.CollectorEndpointOption{
				jaeger.WithPassword(password),
			}
		} else {
			e.collectorOpts = append(
				e.collectorOpts,
				jaeger.WithPassword(password),
			)
		}

		return nil
	}
}

func WithServiceName(name string) Option {
	return func(e *exp) error {
		if name != "" {
			e.name = name
		}
		return nil
	}
}

func WithBufferMaxCount(count int) Option {
	return func(e *exp) error {
		if count < 0 {
			return errors.NewErrInvalidOption("bufferMaxCount", count)
		}

		if e.options == nil {
			e.options = []jaeger.Option{
				jaeger.WithBufferMaxCount(count),
			}
		} else {
			e.options = append(
				e.options,
				jaeger.WithBufferMaxCount(count),
			)
		}

		return nil
	}
}

func WithOnErrorFunc(f func(error)) Option {
	return func(e *exp) error {
		if f == nil {
			return errors.NewErrInvalidOption("onErrorFunc", f)
		}

		if e.options == nil {
			e.options = []jaeger.Option{
				jaeger.WithOnError(f),
			}
		} else {
			e.options = append(
				e.options,
				jaeger.WithOnError(f),
			)
		}

		return nil
	}
}
