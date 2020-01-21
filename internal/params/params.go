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

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/vdaas/vald/internal/errors"
)

type Data struct {
	configFilePath string
	showVersion    bool
}

type parser struct {
	filePath struct {
		keys        []string
		defaultPath string
		description string
	}
	version struct {
		keys         []string
		defaultFlag bool
		description string
	}
}

func New(opts ...Option) *parser {
	p := new(parser)
	for _, opt := range append(defaultOpts, opts...) {
		opt(p)
	}
	return p
}

func (p *parser) Parse() (*Data, bool, error) {
	f := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)

	d := new(Data)
	for _, key := range p.filePath.keys{
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
		return d, true, nil
	}

	return d, false, nil
}

func (d *Data) ConfigFilePath() string {
	return d.configFilePath
}

func (d *Data) ShowVersion() bool {
	return d.showVersion
}
