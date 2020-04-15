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

package errors

var (
	// internal compressor
	ErrInvalidCompressionLevel = func(level int) error {
		return Errorf("invalid compression level: %d", level)
	}

	// Compressor
	ErrCompressorNameNotFound = func(name string) error {
		return Errorf("compressor %s not found", name)
	}

	ErrCompressedDataNotFound = New("compressed data not found")

	ErrDecompressedDataNotFound = New("decompressed data not found")

	ErrCompressFailed = New("compress failed")

	ErrDecompressFailed = New("decompress failed")

	ErrCompressorRegistererIsNotRunning = func() error {
		return Errorf("compressor registerers is not running")
	}

	ErrCompressorRegistererChannelIsFull = func() error {
		return Errorf("compressor registerer channel is full")
	}
)
