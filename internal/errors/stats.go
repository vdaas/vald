//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package errors

var (
	// ErrCPULineNotFound represents a function to generate an error that the CPU line is not found in /proc/stat.
	ErrCPULineNotFound = func() error {
		return New("cpu line not found in /proc/stat")
	}

	// ErrInvalidCPULineFormat represents a function to generate an error that the CPU line format is invalid in /proc/stat.
	ErrInvalidCPULineFormat = func() error {
		return New("invalid cpu line format in /proc/stat")
	}

	// ErrCPUFieldParseFailed represents a function to generate an error that the CPU field parse failed.
	ErrCPUFieldParseFailed = func(fieldIndex int, err error) error {
		return Wrapf(err, "failed to parse CPU field %d", fieldIndex)
	}

	// ErrTotalMemoryNotFound represents a function to generate an error that the MemTotal is not found in /proc/meminfo.
	ErrTotalMemoryNotFound = func() error {
		return New("MemTotal not found in /proc/meminfo")
	}

	// ErrMemTotalParseFailed represents a function to generate an error that the MemTotal parse failed.
	ErrMemTotalParseFailed = func(err error) error {
		return Wrap(err, "failed to parse MemTotal")
	}

	// ErrProcessMemoryNotFound represents a function to generate an error that the VmRSS is not found in /proc/self/status.
	ErrProcessMemoryNotFound = func() error {
		return New("VmRSS not found in /proc/self/status")
	}

	// ErrVmRSSParseFailed represents a function to generate an error that the VmRSS parse failed.
	ErrVmRSSParseFailed = func(err error) error {
		return Wrap(err, "failed to parse VmRSS")
	}
)
