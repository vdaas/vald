// MIT License
//
// Copyright (c) 2019 kpango (Yusuke Kato)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
