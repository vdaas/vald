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

// Package errors provides error types and function
package errors

var (
	// ErrInvalidCompressionLevel represents a function to generate an error of invalid compression level.
	ErrInvalidCompressionLevel = func(level int) error {
		return Errorf("invalid compression level: %d", level)
	}

	// ErrCompressorNameNotFound represents a function to generate an error of compressor not found.
	ErrCompressorNameNotFound = func(name string) error {
		return Errorf("compressor %s not found", name)
	}

	// ErrCompressedDataNotFound returns an error of compressed data is not found.
	ErrCompressedDataNotFound = New("compressed data not found")

	// ErrDecompressedDataNotFound returns an error of decompressed data is not found.
	ErrDecompressedDataNotFound = New("decompressed data not found")

	// ErrCompressFailed returns an error of compress failed.
	ErrCompressFailed = New("compress failed")

	// ErrDecompressFailed returns an error of decompressing failed.
	ErrDecompressFailed = New("decompress failed")

	// ErrCompressorRegistererIsNotRunning represents a function to generate an error of compressor registerers is not running.
	ErrCompressorRegistererIsNotRunning = func() error {
		return New("compressor registerers is not running")
	}

	// ErrCompressorRegistererChannelIsFull represents a function to generate an error that compressor registerer channel is full.
	ErrCompressorRegistererChannelIsFull = func() error {
		return New("compressor registerer channel is full")
	}
)
