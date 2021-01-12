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
	"github.com/vdaas/vald/internal/info"
	"github.com/vdaas/vald/internal/observability/client/google"
)

type Option func(p *prof) error

var defaultOptions = []Option{
	WithCPUProfiling(true),
	WithAllocProfiling(true),
	WithHeapProfiling(true),
	WithGoroutineProfiling(true),
	WithServiceVersion(info.Version),
}

func WithProjectID(pid string) Option {
	return func(p *prof) error {
		if pid != "" {
			p.ProjectID = pid
		}

		return nil
	}
}

func WithService(name string) Option {
	return func(p *prof) error {
		if name != "" {
			p.Service = name
		}

		return nil
	}
}

func WithServiceVersion(version string) Option {
	return func(p *prof) error {
		if version != "" {
			p.ServiceVersion = version
		}

		return nil
	}
}

func WithDebugLogging(enabled bool) Option {
	return func(p *prof) error {
		p.DebugLogging = enabled

		return nil
	}
}

func WithMutexProfiling(enabled bool) Option {
	return func(p *prof) error {
		p.MutexProfiling = enabled

		return nil
	}
}

func WithCPUProfiling(enabled bool) Option {
	return func(p *prof) error {
		p.NoCPUProfiling = !enabled

		return nil
	}
}

func WithAllocProfiling(enabled bool) Option {
	return func(p *prof) error {
		p.NoAllocProfiling = !enabled

		return nil
	}
}

func WithHeapProfiling(enabled bool) Option {
	return func(p *prof) error {
		p.NoHeapProfiling = !enabled

		return nil
	}
}

func WithGoroutineProfiling(enabled bool) Option {
	return func(p *prof) error {
		p.NoGoroutineProfiling = !enabled

		return nil
	}
}

func WithAllocForceGC(enabled bool) Option {
	return func(p *prof) error {
		p.AllocForceGC = enabled

		return nil
	}
}

func WithAPIAddr(addr string) Option {
	return func(p *prof) error {
		if addr != "" {
			p.APIAddr = addr
		}

		return nil
	}
}

func WithInstance(instance string) Option {
	return func(p *prof) error {
		if instance != "" {
			p.Instance = instance
		}

		return nil
	}
}

func WithZone(zone string) Option {
	return func(p *prof) error {
		if zone != "" {
			p.Zone = zone
		}

		return nil
	}
}

func WithClientOptions(copts ...google.Option) Option {
	return func(p *prof) error {
		opts := make([]google.Option, 0, len(copts))
		for _, opt := range copts {
			if opt != nil {
				opts = append(opts, opt)
			}
		}

		if p.clientOpts == nil {
			p.clientOpts = opts
			return nil
		}

		p.clientOpts = append(p.clientOpts, opts...)

		return nil
	}
}
