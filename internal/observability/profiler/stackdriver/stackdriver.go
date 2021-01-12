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

// Package stackdriver provides a stackdriver exporter.
package stackdriver

import (
	"context"

	"cloud.google.com/go/profiler"
	"github.com/vdaas/vald/internal/observability/client/google"
	vprof "github.com/vdaas/vald/internal/observability/profiler"
)

type Stackdriver interface {
	vprof.Profiler
}

type prof struct {
	*profiler.Config
	clientOpts []google.Option
}

func New(opts ...Option) (s Stackdriver, err error) {
	p := new(prof)
	p.Config = new(profiler.Config)

	for _, opt := range append(defaultOptions, opts...) {
		err = opt(p)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (p *prof) Start(ctx context.Context) (err error) {
	return profiler.Start(*p.Config, p.clientOpts...)
}

func (p *prof) Stop(ctx context.Context) {
}
