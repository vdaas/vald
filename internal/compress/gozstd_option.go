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

// Package compress provides compress functions
package compress

import (
	"github.com/klauspost/compress/zstd"
	"github.com/vdaas/vald/internal/errors"
)

type GoZstdOption func(c *goZstdCompressor) error

var (
	defaultGoZstdOpts = []GoZstdOption{
		WithGoZstdGob(),
		WithGoZstdCompressionLevel(10),
	}
)

func WithGoZstdGob(opts ...GobOption) GoZstdOption {
	return func(c *goZstdCompressor) error {
		gobc, err := NewGob(opts...)
		if err != nil {
			return err
		}
		c.gobc = gobc
		return nil
	}
}

func WithGoZstdCompressionLevel(level int) GoZstdOption {
	return func(c *goZstdCompressor) error {
		if level < BestSpeed || level > BestCompression {
			return errors.ErrInvalidCompressionLevel(level)
		}
		c.eoptions = append(c.eoptions, zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(level)))
		return nil
	}
}
