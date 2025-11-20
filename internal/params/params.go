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
	"flag"
	"slices"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
)

type ErrorHandling = flag.ErrorHandling

const (
	ContinueOnError ErrorHandling = flag.ContinueOnError
	PanicOnError    ErrorHandling = flag.PanicOnError
	ExitOnError     ErrorHandling = flag.ExitOnError
)

// Data is an interface to get the configuration path and flag.
type Data interface {
	ConfigFilePath() string
	ShowVersion() bool
}

type data struct {
	configFilePath string
	showVersion    bool
}

// Parser is an interface to parse commnad-line argument.
type Parser interface {
	Parse() (Data, bool, error)
}

type parser struct {
	overrideDefault bool
	name            string
	filters         []func(string) bool
	f               *flag.FlagSet
	defaults        *flag.FlagSet
	filePath        struct {
		keys        []string
		defaultPath string
		description string
	}
	version struct {
		keys        []string
		defaultFlag bool
		description string
	}
	ErrorHandler ErrorHandling
}

// New returns parser object.
func New(opts ...Option) Parser {
	p := new(parser)
	for _, opt := range append(defaultOptions, opts...) {
		opt(p)
	}
	p.defaults = flag.CommandLine
	p.f = flag.NewFlagSet(p.name, p.ErrorHandler)
	if p.overrideDefault {
		p.Override()
	}
	return p
}

// Parse parses command-line argument and returns parsed data and whether there is a help option or not and error.
func (p *parser) Parse() (Data, bool, error) {
	if p == nil || p.f == nil {
		return nil, false, errors.ErrArgumentParserNotFound
	}
	d := new(data)
	for _, key := range p.filePath.keys {
		p.f.StringVar(&d.configFilePath,
			key,
			p.filePath.defaultPath,
			p.filePath.description,
		)
	}

	for _, key := range p.version.keys {
		p.f.BoolVar(&d.showVersion,
			key,
			p.version.defaultFlag,
			p.version.description,
		)
	}

	args := os.Args[1:]
	if p.filters != nil {
		args = slices.DeleteFunc(args, func(s string) bool {
			for _, filter := range p.filters {
				if filter != nil && filter(s) {
					return true
				}
			}
			return false
		})
	}

	err := p.f.Parse(args)
	if err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			return nil, false, errors.ErrArgumentParseFailed(err)
		}
		return nil, true, nil
	}

	d.configFilePath = file.AbsolutePath(d.configFilePath)

	if exist, _, err := file.ExistsWithDetail(d.configFilePath); !d.showVersion &&
		(!exist || d.configFilePath == "") {
		p.f.Usage()
		return nil, true, err
	}

	return d, false, nil
}

func (p *parser) Restore() {
	if p.defaults != nil {
		flag.CommandLine = p.defaults
	}
}

func (p *parser) Override() {
	if p.f != nil {
		flag.CommandLine = p.f
	}
}

// ConfigFilePath returns configFilePath.
func (d *data) ConfigFilePath() string {
	return d.configFilePath
}

// ShowVersion returns showVersion.
func (d *data) ShowVersion() bool {
	return d.showVersion
}
