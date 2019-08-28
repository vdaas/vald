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
		key         string
		defaultPath string
		description string
	}
	version struct {
		key         string
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

func (p *parser) Parse() (*Data, error) {

	f := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ContinueOnError)

	d := new(Data)
	f.StringVar(&d.configFilePath,
		p.filePath.key,
		p.filePath.defaultPath,
		p.filePath.description,
	)

	f.BoolVar(&d.showVersion,
		p.version.key,
		p.version.defaultFlag,
		p.version.description,
	)

	err := f.Parse(os.Args[1:])
	if err != nil {
		return nil, errors.ErrArgumentParseFailed(err)
	}

	return d, nil
}

func (d *Data) ConfigFilePath() string {
	return d.configFilePath
}

func (d *Data) ShowVersion() bool {
	return d.showVersion
}
