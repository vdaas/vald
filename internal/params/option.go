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

// Package params provides implementation of Go API for argument parser
package params

type Option func(*parser)

var (
	defaultOpts = []Option{
		WithConfigFilePathKeys("f", "file", "c", "config"),
		WithConfigFilePathDefault("/etc/server/config.yaml"),
		WithConfigFileDescription("config file path"),
		WithVersionKeys("v", "ver", "version"),
		WithVersionFlagDefault(false),
		WithVersionDescription("show server version"),
	}
)

func WithConfigFilePathKeys(keys ...string) Option {
	return func(p *parser) {
		if len(keys) != 0 {
			p.filePath.keys = append(p.filePath.keys, keys...)
		}
	}
}

func WithConfigFilePathDefault(path string) Option {
	return func(p *parser) {
		if path != "" {
			p.filePath.defaultPath = path
		}
	}
}

func WithConfigFileDescription(desc string) Option {
	return func(p *parser) {
		if desc != "" {
			p.filePath.description = desc
		}
	}
}

func WithVersionKeys(keys ...string) Option {
	return func(p *parser) {
		if len(keys) != 0 {
			p.version.keys = append(p.version.keys, keys...)
		}
	}
}

func WithVersionFlagDefault(flag bool) Option {
	return func(p *parser) {
		p.version.defaultFlag = flag
	}
}

func WithVersionDescription(desc string) Option {
	return func(p *parser) {
		if desc != "" {
			p.version.description = desc
		}
	}
}
