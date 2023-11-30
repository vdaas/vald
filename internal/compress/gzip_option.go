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

// Package compress provides compress functions
package compress

import (
	"github.com/vdaas/vald/internal/compress/gzip"
	"github.com/vdaas/vald/internal/errors"
)

// GzipOption represents the functional option for gzipCompressor.
type GzipOption func(c *gzipCompressor) error

var defaultGzipOpts = []GzipOption{
	WithGzipGob(),
	WithGzipCompressionLevel(gzip.DefaultCompression),
}

// WithGzipGob represents the option to set the GobOption to initialize Gob.
func WithGzipGob(opts ...GobOption) GzipOption {
	return func(c *gzipCompressor) error {
		gobc, err := NewGob(opts...)
		if err != nil {
			return err
		}
		c.gobc = gobc
		return nil
	}
}

// WithGzipCompressionLevel represents the option to set the compress level for gzip.
func WithGzipCompressionLevel(level int) GzipOption {
	return func(c *gzipCompressor) error {
		if level < gzip.HuffmanOnly || level > gzip.BestCompression {
			return errors.ErrInvalidCompressionLevel(level)
		}
		c.compressionLevel = level
		return nil
	}
}
