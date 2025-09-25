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

	// ErrCgroupModeDetectionFailed represents a function to generate an error that cgroup mode detection failed.
	ErrCgroupModeDetectionFailed = func() error {
		return New("unable to detect cgroups mode")
	}

	// ErrCgroupFirstSampleFailed represents a function to generate an error that the first cgroup sample failed.
	ErrCgroupFirstSampleFailed = func(err error) error {
		return Wrap(err, "failed to get first sample")
	}

	// ErrCgroupSecondSampleFailed represents a function to generate an error that the second cgroup sample failed.
	ErrCgroupSecondSampleFailed = func(err error) error {
		return Wrap(err, "failed to get second sample")
	}

	// ErrCgroupV2MemoryCurrentRead represents a function to generate an error that reading memory.current failed.
	ErrCgroupV2MemoryCurrentRead = func(err error) error {
		return Wrap(err, "v2 read memory.current")
	}

	// ErrCgroupV2MemoryCurrentParse represents a function to generate an error that parsing memory.current failed.
	ErrCgroupV2MemoryCurrentParse = func(err error) error {
		return Wrap(err, "v2 parse memory.current")
	}

	// ErrCgroupV2MemoryMaxRead represents a function to generate an error that reading memory.max failed.
	ErrCgroupV2MemoryMaxRead = func(err error) error {
		return Wrap(err, "v2 read memory.max")
	}

	// ErrCgroupV2MemoryMaxParse represents a function to generate an error that parsing memory.max failed.
	ErrCgroupV2MemoryMaxParse = func(err error) error {
		return Wrap(err, "v2 parse memory.max")
	}

	// ErrCgroupV2CPUStatRead represents a function to generate an error that reading cpu.stat failed.
	ErrCgroupV2CPUStatRead = func(err error) error {
		return Wrap(err, "v2 read cpu.stat")
	}

	// ErrCgroupV2CPUStatMissingUsage represents a function to generate an error that cpu.stat is missing usage_usec.
	ErrCgroupV2CPUStatMissingUsage = func() error {
		return New("v2 cpu.stat missing usage_usec")
	}

	// ErrCgroupV2CPUMaxRead represents a function to generate an error that reading cpu.max failed.
	ErrCgroupV2CPUMaxRead = func(err error) error {
		return Wrap(err, "v2 read cpu.max")
	}

	// ErrCgroupV2CPUMaxMalformed represents a function to generate an error that cpu.max is malformed.
	ErrCgroupV2CPUMaxMalformed = func(val string) error {
		return Errorf("v2 cpu.max malformed: %q", val)
	}

	// ErrCgroupV2CPUMaxParseQuota represents a function to generate an error that parsing cpu.max quota failed.
	ErrCgroupV2CPUMaxParseQuota = func(err error) error {
		return Wrap(err, "v2 cpu.max parse quota")
	}

	// ErrCgroupV2CPUMaxParsePeriod represents a function to generate an error that parsing cpu.max period failed.
	ErrCgroupV2CPUMaxParsePeriod = func(err error) error {
		return Wrap(err, "v2 cpu.max parse period")
	}

	// ErrCgroupV1MemoryUsageReadFailed represents a function to generate an error that reading memory usage failed.
	ErrCgroupV1MemoryUsageReadFailed = func(err error) error {
		return Wrap(err, "v1 memory usage read failed")
	}

	// ErrCgroupV1MemoryUsageParseFailed represents a function to generate an error that parsing memory usage failed.
	ErrCgroupV1MemoryUsageParseFailed = func(err error) error {
		return Wrap(err, "v1 memory usage parse failed")
	}

	// ErrCgroupV1CPUUsageReadFailed represents a function to generate an error that reading CPU usage failed.
	ErrCgroupV1CPUUsageReadFailed = func(err error) error {
		return Wrap(err, "v1 cpuacct.usage read failed")
	}
)
