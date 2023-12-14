//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package params provides implementation of Go API for argument parser
package params

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
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
	filePath struct {
		keys        []string
		defaultPath string
		description string
	}
	version struct {
		keys        []string
		defaultFlag bool
		description string
	}
}

// New returns parser object.
func New(opts ...Option) Parser {
	p := new(parser)
	for _, opt := range append(defaultOptions, opts...) {
		opt(p)
	}
	return p
}

// Parse parses command-line argument and returns parsed data and whether there is a help option or not and error.
func (p *parser) Parse() (Data, bool, error) {
	f := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)

	d := new(data)
	for _, key := range p.filePath.keys {
		f.StringVar(&d.configFilePath,
			key,
			p.filePath.defaultPath,
			p.filePath.description,
		)
	}

	for _, key := range p.version.keys {
		f.BoolVar(&d.showVersion,
			key,
			p.version.defaultFlag,
			p.version.description,
		)
	}

	err := f.Parse(os.Args[1:])
	if err != nil {
		if err != flag.ErrHelp {
			return nil, false, errors.ErrArgumentParseFailed(err)
		}
		return nil, true, nil
	}

	if exist, _, err := file.ExistsWithDetail(d.configFilePath); !d.showVersion &&
		(!exist ||
			d.configFilePath == "") {
		f.Usage()
		return nil, true, err
	}

	return d, false, nil
}

// ConfigFilePath returns configFilePath.
func (d *data) ConfigFilePath() string {
	return d.configFilePath
}

// ShowVersion returns showVersion.
func (d *data) ShowVersion() bool {
	return d.showVersion
}
