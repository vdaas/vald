//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package compress provides compress functions
package compress

import (
	"github.com/vdaas/vald/internal/compress/zstd"
)

// ZstdOption represents the functional option for zstdCompressor.
type ZstdOption func(c *zstdCompressor) error

var defaultZstdOpts = []ZstdOption{
	WithZstdGob(),
	WithZstdCompressionLevel(3),
}

// WithZstdGob represents the option to set the GobOption to initialize Gob.
func WithZstdGob(opts ...GobOption) ZstdOption {
	return func(c *zstdCompressor) error {
		gobc, err := NewGob(opts...)
		if err != nil {
			return err
		}
		c.gobc = gobc
		return nil
	}
}

// WithZstdCompressionLevel represents the option to set the compress level for zstd.
func WithZstdCompressionLevel(level int) ZstdOption {
	return func(c *zstdCompressor) error {
		c.eoptions = append(c.eoptions, zstd.WithEncoderLevel(level))
		return nil
	}
}
