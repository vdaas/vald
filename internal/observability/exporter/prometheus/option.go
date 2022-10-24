// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package prometheus

import (
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/observability/metrics"
)

type Option func(e *exp) error

var defaultOpts = []Option{
	WithEndpoint("/metrics"),
	WithNamespace("vald"),
}

func WithEndpoint(ep string) Option {
	return func(e *exp) error {
		if ep == "" {
			return errors.NewErrInvalidOption("endpoint", ep)
		}
		e.endpoint = ep
		return nil
	}
}

func WithNamespace(ns string) Option {
	return func(e *exp) error {
		if ns == "" {
			return errors.NewErrInvalidOption("namespace", ns)
		}
		e.namespace = ns
		return nil
	}
}

func WithView(viewers ...metrics.Viewer) Option {
	return func(e *exp) error {
		views := make([]metrics.View, 0, len(viewers))
		for _, viewer := range viewers {
			vs, err := viewer.View()
			if err != nil {
				return errors.NewErrCriticalOption("view", viewer, err)
			}
			for _, v := range vs {
				views = append(views, *v)
			}
		}
		e.views = append(e.views, views...)
		return nil
	}
}
