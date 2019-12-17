//
// Copyright (C) 2019 Vdaas.org Vald team ( kpango, kmrmt, rinx )
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
		WithConfigFilePathKey("f"),
		WithConfigFilePathDefault("/etc/server/config.yaml"),
		WithConfigFileDescription("config file path"),
		WithVersionKey("version"),
		WithVersionFlagDefault(false),
		WithVersionDescription("show server version"),
	}
)

func WithConfigFilePathKey(key string) Option {
	return func(p *parser) {
		p.filePath.key = key
	}
}

func WithConfigFilePathDefault(path string) Option {
	return func(p *parser) {
		p.filePath.defaultPath = path
	}
}

func WithConfigFileDescription(desc string) Option {
	return func(p *parser) {
		p.filePath.description = desc
	}
}

func WithVersionKey(key string) Option {
	return func(p *parser) {
		p.version.key = key
	}
}

func WithVersionFlagDefault(flag bool) Option {
	return func(p *parser) {
		p.version.defaultFlag = flag
	}
}

func WithVersionDescription(desc string) Option {
	return func(p *parser) {
		p.version.description = desc
	}
}
