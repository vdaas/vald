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

// Package config providers configuration type and load configuration logic
package config

import "strings"

type compressAlgorithm uint8

const (
	GOB compressAlgorithm = iota
	GZIP
	LZ4
	ZSTD
	DDZSTD
)

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
	case DDZSTD:
		return "ddzstd"
	}
	return "unknown"
}

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
	case "ddzstd":
		return DDZSTD
	}
	return 0
}

type Compressor struct {
	CompressAlgorithm string `json:"compress_algorithm" yaml:"compress_algorithm"`
	CompressionLevel  int    `json:"compression_level" yaml:"compression_level"`
	ConcurrentLimit   int    `json:"concurrent_limit" yaml:"concurrent_limit"`
	Buffer            int    `json:"buffer" yaml:"buffer"`
}

func (c *Compressor) Bind() *Compressor {
	c.CompressAlgorithm = GetActualValue(c.CompressAlgorithm)
	return c
}
