//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// Package config providers configuration type and load configuration logic
package config

import "github.com/vdaas/vald/internal/strings"

type compressAlgorithm uint8

const (
	// GOB represents gob algorithm.
	GOB compressAlgorithm = 1 + iota
	// GZIP represents gzip algorithm.
	GZIP
	// LZ4 represents lz4 algorithm.
	LZ4
	// ZSTD represents zstd algorithm.
	ZSTD
)

// String returns compress algorithm.
func (ca compressAlgorithm) String() string {
	switch ca {
	case GOB:
		return "gob"
	case GZIP:
		return "gzip"
	case LZ4:
		return "lz4"
	case ZSTD:
		return "zstd"
	}
	return "unknown"
}

// CompressAlgorithm returns compressAlgorithm converted from string.
func CompressAlgorithm(ca string) compressAlgorithm {
	switch strings.ToLower(ca) {
	case "gob":
		return GOB
	case "gzip":
		return GZIP
	case "lz4":
		return LZ4
	case "zstd":
		return ZSTD
	}
	return 0
}

// CompressCore represents CompressCore configuration.
type CompressCore struct {
	// CompressorAlgorithm represents compression algorithm type
	// compression algorithm. must be `gob`, `gzip`, `lz4` or `zstd`
	CompressAlgorithm string `json:"compress_algorithm" yaml:"compress_algorithm"`

	// CompressionLevel represents compression level
	// compression level. value range relies on which algorithm is used. `gob`: level will be ignored. `gzip`: -1 (default compression), 0 (no compression), or 1 (best speed) to 9 (best compression). `lz4`: >= 0, higher is better compression. `zstd`: 1 (fastest) to 22 (best), however implementation relies on klauspost/compress.
	CompressionLevel int `json:"compression_level" yaml:"compression_level"`
}

// Bind binds the actual data from the receiver field.
func (c *CompressCore) Bind() *CompressCore {
	c.CompressAlgorithm = GetActualValue(c.CompressAlgorithm)

	return c
}
