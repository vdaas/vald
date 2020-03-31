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

type Compressor struct {
	// CompressorAlgorithm represents compression algorithm type
	CompressAlgorithm string `json:"compress_algorithm" yaml:"compress_algorithm"`

	// CompressionLevel represents compression level
	CompressionLevel int `json:"compression_level" yaml:"compression_level"`

	// ConcurrentLimit represents limitation of compression worker concurrency
	ConcurrentLimit int `json:"concurrent_limit" yaml:"concurrent_limit"`

	// Buffer represents capacity of buffer for compression
	Buffer int `json:"buffer" yaml:"buffer"`

	// PodIP represents pod ip of compressor instance. it is recommended to use status.podIP field of k8s pod
	PodIP string `json:"pod_ip" yaml:"pod_ip"`

	// Registerer represents registerer options
	Registerer *CompressorRegisterer `json:"registerer" yaml:"registerer"`
}

type CompressorRegisterer struct {
	// Buffer represents capacity of buffer for registerer
	Buffer int `json:"buffer" yaml:"buffer"`

	// Backoff represents backoff configuration of registerer
	Backoff *Backoff `json:"backoff" yaml:"backoff"`

	// Worker represents worker options
	Worker *Worker `json:"worker" yaml:"worker"`

	// Compressor represents gRPC client config of compressor client (for forwarding use)
	Compressor *BackupManager `json:"compressor" yaml:"compressor"`
}

type Worker struct {
	// ConcurrentLimit represents limitation of worker
	ConcurrentLimit int `json:"concurrent_limit" yaml:"concurrent_limit"`

	// Buffer represents capacity of buffer for worker
	Buffer int `json:"buffer" yaml:"buffer"`
}

func (c *Compressor) Bind() *Compressor {
	c.CompressAlgorithm = GetActualValue(c.CompressAlgorithm)

	c.PodIP = GetActualValue(c.PodIP)

	if c.Registerer == nil {
		c.Registerer = new(CompressorRegisterer)
	}

	if c.Registerer.Backoff != nil {
		c.Registerer.Backoff = c.Registerer.Backoff.Bind()
	} else {
		c.Registerer.Backoff = new(Backoff)
	}

	if c.Registerer.Worker == nil {
		c.Registerer.Worker = new(Worker)
	}

	if c.Registerer.Compressor != nil {
		c.Registerer.Compressor = c.Registerer.Compressor.Bind()
	} else {
		c.Registerer.Compressor = new(BackupManager)
	}

	return c
}
