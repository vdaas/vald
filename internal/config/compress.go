//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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
	GOB compressAlgorithm = 1 + iota
	GZIP
	LZ4
	ZSTD
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
	}
	return 0
}

type CompressCore struct {
	// CompressorAlgorithm represents compression algorithm type
	CompressAlgorithm string `json:"compress_algorithm" yaml:"compress_algorithm"`

	// CompressionLevel represents compression level
	CompressionLevel int `json:"compression_level" yaml:"compression_level"`
}

func (c *CompressCore) Bind() *CompressCore {
	c.CompressAlgorithm = GetActualValue(c.CompressAlgorithm)

	return c
}

type Compressor struct {
	CompressCore `json:",inline" yaml:",inline"`

	// ConcurrentLimit represents limitation of compression worker concurrency
	ConcurrentLimit int `json:"concurrent_limit" yaml:"concurrent_limit"`

	// QueueCheckDuration represents duration of queue daemon block
	QueueCheckDuration string `json:"queue_check_duration" yaml:"queue_check_duration"`
}

func (c *Compressor) Bind() *Compressor {
	c.CompressCore = *c.CompressCore.Bind()

	c.QueueCheckDuration = GetActualValue(c.QueueCheckDuration)

	return c
}

type CompressorRegisterer struct {
	// ConcurrentLimit represents limitation of worker
	ConcurrentLimit int `json:"concurrent_limit" yaml:"concurrent_limit"`

	// QueueCheckDuration represents duration of queue daemon block
	QueueCheckDuration string `json:"queue_check_duration" yaml:"queue_check_duration"`

	// Compressor represents gRPC client config of compressor client (for forwarding use)
	Compressor *BackupManager `json:"compressor" yaml:"compressor"`
}

func (cr *CompressorRegisterer) Bind() *CompressorRegisterer {
	cr.QueueCheckDuration = GetActualValue(cr.QueueCheckDuration)

	if cr.Compressor != nil {
		cr.Compressor = cr.Compressor.Bind()
	} else {
		cr.Compressor = new(BackupManager)
	}

	return cr
}
