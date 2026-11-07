// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compress

// LZ4Option represents the functional option for lz4Compressor.
type LZ4Option func(c *lz4Compressor) error

var defaultLZ4Opts = []LZ4Option{
	WithLZ4Gob(),
	WithLZ4CompressionLevel(0),
}

// WithLZ4Gob returns the option to set gobc for lz4Compressor.
func WithLZ4Gob(opts ...GobOption) LZ4Option {
	return func(c *lz4Compressor) error {
		gobc, err := NewGob(opts...)
		if err != nil {
			return err
		}
		c.gobc = gobc
		return nil
	}
}

// WithLZ4CompressionLevel returns the option to set compressionLevel for lz4Compressor.
// Every int is accepted: pierrec/lz4 v3 uses 0 for the fast compressor and a
// nonzero level as the HC search depth (<= 0 means unlimited).
func WithLZ4CompressionLevel(level int) LZ4Option {
	return func(c *lz4Compressor) error {
		c.compressionLevel = level
		return nil
	}
}
