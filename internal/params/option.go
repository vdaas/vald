//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package params

import (
	"os"
	"path/filepath"
)

type Option func(*parser)

var defaultOptions = []Option{
	WithName(filepath.Base(os.Args[0])),
	WithConfigFilePathKeys("f", "file", "c", "config"),
	WithConfigFilePathDefault("/etc/server/config.yaml"),
	WithConfigFileDescription("config file path"),
	WithVersionKeys("v", "ver", "version"),
	WithVersionFlagDefault(false),
	WithVersionDescription("show server version"),
	WithOverrideDefault(false),
}

// WithName returns Option that sets name.
func WithName(name string) Option {
	return func(p *parser) {
		if name != "" {
			p.name = name
		}
	}
}

// WithConfigFilePathKeys returns Option that sets filePath.keys.
func WithConfigFilePathKeys(keys ...string) Option {
	return func(p *parser) {
		p.filePath.keys = append(p.filePath.keys, keys...)
	}
}

// WithConfigFilePathDefault returns Option that sets filePath.defaultPath.
func WithConfigFilePathDefault(path string) Option {
	return func(p *parser) {
		p.filePath.defaultPath = path
	}
}

// WithConfigFileDescription returns Option that sets filePath.description.
func WithConfigFileDescription(desc string) Option {
	return func(p *parser) {
		p.filePath.description = desc
	}
}

// WithVersionKeys returns Option that sets version.keys.
func WithVersionKeys(keys ...string) Option {
	return func(p *parser) {
		p.version.keys = append(p.version.keys, keys...)
	}
}

// WithVersionFlagDefault returns Option that sets version.defaultFlag.
func WithVersionFlagDefault(flag bool) Option {
	return func(p *parser) {
		p.version.defaultFlag = flag
	}
}

// WithVersionDescription returns Option that sets version.description.
func WithVersionDescription(desc string) Option {
	return func(p *parser) {
		p.version.description = desc
	}
}

// WithOverrideDefault returns Option that overrides default flag.CommandLine.
func WithOverrideDefault(flag bool) Option {
	return func(p *parser) {
		p.overrideDefault = flag
	}
}

// WithArgumentFilters returns Option that sets filters.
// filters is a slice of functions that takes a string and returns a bool.
// If the string not matched all filters (means filter returns false), it will be added to the arguments.
func WithArgumentFilters(filters ...func(string) bool) Option {
	return func(p *parser) {
		if len(filters) == 0 {
			return
		}
		if p.filters == nil {
			p.filters = make([]func(string) bool, 0, len(filters))
		}
		p.filters = append(p.filters, filters...)
	}
}
