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

// Package service
package service

import (
	"github.com/vdaas/vald/internal/errgroup"
)

type CompressorOption func(c *compressor) error

var (
	defaultCompressorOpts = []CompressorOption{
		WithLimitation(0),
		WithBuffer(10),
		WithCompressAlgorithm("gob"),
		WithErrGroup(errgroup.Get()),
	}
)

func WithLimitation(limit int) CompressorOption {
	return func(c *compressor) error {
		c.limitation = limit
		return nil
	}
}

func WithBuffer(b int) CompressorOption {
	return func(c *compressor) error {
		c.buffer = b
		return nil
	}
}

func WithCompressAlgorithm(name string) CompressorOption {
	return func(c *compressor) error {
		c.algorithm = name
		return nil
	}
}

func WithCompressionLevel(level int) CompressorOption {
	return func(c *compressor) error {
		c.compressionLevel = level
		return nil
	}
}

func WithErrGroup(eg errgroup.Group) CompressorOption {
	return func(c *compressor) error {
		if eg != nil {
			c.eg = eg
		}
		return nil
	}
}
