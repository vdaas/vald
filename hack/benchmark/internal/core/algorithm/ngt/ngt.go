//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package ngt provides ngt
package ngt

import (
	"os"

	c "github.com/vdaas/vald/hack/benchmark/internal/core/algorithm"
	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/file"
)

type ObjectType int

const (
	ObjectNone ObjectType = iota
	Uint8
	Float
)

type core struct {
	idxPath    string
	tmpdir     string
	objectType ObjectType
	dimension  int
	ngt.NGT
}

func New(opts ...Option) (c.Bit32, error) {
	c := new(core)
	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	tmpdir, err := file.MkdirTemp(c.idxPath)
	if err != nil {
		return nil, err
	}
	c.tmpdir = tmpdir

	typ := ngt.ObjectNone
	switch c.objectType {
	case Uint8:
		typ = ngt.Uint8
	case Float:
		typ = ngt.Float
	}

	n, err := ngt.New(
		ngt.WithIndexPath(tmpdir),
		ngt.WithDimension(c.dimension),
		ngt.WithObjectType(typ),
	)
	if err != nil {
		return nil, err
	}
	c.NGT = n

	return c, nil
}

func (c *core) Search(vec []float32, size int, epsilon, radius float32) (interface{}, error) {
	return c.NGT.Search(vec, size, epsilon, radius)
}

func (c *core) Close() {
	if len(c.tmpdir) != 0 {
		os.RemoveAll(c.tmpdir)
	}
	// c.NGT.Close()
}
