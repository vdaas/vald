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

// Package config providers configuration type and load configuration logic
package config

import "github.com/vdaas/vald/internal/strings"

// CompressAlgorithm is an enum for compress algorithm.
type CompressAlgorithm uint8

const (
	// GOB represents gob algorithm.
	GOB CompressAlgorithm = 1 + iota
	// GZIP represents gzip algorithm.
	GZIP
	// LZ4 represents lz4 algorithm.
	LZ4
	// ZSTD represents zstd algorithm.
	ZSTD
)

// String returns compress algorithm.
func (ca CompressAlgorithm) String() string {
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

// AToCompressAlgorithm returns CompressAlgorithm converted from string.
func AToCompressAlgorithm(ca string) CompressAlgorithm {
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
	CompressAlgorithm string `json:"compress_algorithm" yaml:"compress_algorithm"`

	// CompressionLevel represents compression level
	CompressionLevel int `json:"compression_level" yaml:"compression_level"`
}

// Bind binds the actual data from the receiver field.
func (c *CompressCore) Bind() *CompressCore {
	c.CompressAlgorithm = GetActualValue(c.CompressAlgorithm)

	return c
}

// Compressor represents Compressor configuration.
type Compressor struct {
	CompressCore `json:",inline" yaml:",inline"`

	// ConcurrentLimit represents limitation of compression worker concurrency
	ConcurrentLimit int `json:"concurrent_limit" yaml:"concurrent_limit"`

	// QueueCheckDuration represents duration of queue daemon block
	QueueCheckDuration string `json:"queue_check_duration" yaml:"queue_check_duration"`
}

// Bind binds the actual data from the Compressor receiver field.
func (c *Compressor) Bind() *Compressor {
	c.CompressCore = *c.CompressCore.Bind()

	c.QueueCheckDuration = GetActualValue(c.QueueCheckDuration)

	return c
}

// CompressorRegisterer represents CompressorRegisterer configuration.
type CompressorRegisterer struct {
	// ConcurrentLimit represents limitation of worker
	ConcurrentLimit int `json:"concurrent_limit" yaml:"concurrent_limit"`

	// QueueCheckDuration represents duration of queue daemon block
	QueueCheckDuration string `json:"queue_check_duration" yaml:"queue_check_duration"`

	// Compressor represents gRPC client config of compressor client (for forwarding use)
	Compressor *BackupManager `json:"compressor" yaml:"compressor"`
}

// Bind binds the actual data from the CompressorRegisterer receiver field.
func (cr *CompressorRegisterer) Bind() *CompressorRegisterer {
	cr.QueueCheckDuration = GetActualValue(cr.QueueCheckDuration)

	if cr.Compressor != nil {
		cr.Compressor = cr.Compressor.Bind()
	} else {
		cr.Compressor = new(BackupManager)
	}

	return cr
}
