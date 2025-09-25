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
	// ErrCgroupModeDetectionFailed represents a function to generate an error that cgroup mode detection failed.
	ErrCgroupModeDetectionFailed = New("unable to detect cgroups mode")

	// ErrCgroupFirstSampleFailed represents a function to generate an error that the first cgroup sample failed.
	ErrCgroupFirstSampleFailed = func(err error) error {
		return Wrap(err, "failed to get first sample")
	}

	// ErrCgroupSecondSampleFailed represents a function to generate an error that the second cgroup sample failed.
	ErrCgroupSecondSampleFailed = func(err error) error {
		return Wrap(err, "failed to get second sample")
	}

	// ErrCgroupV2MemoryCurrentReadFailed represents a function to generate an error that reading memory.current failed.
	ErrCgroupV2MemoryCurrentReadFailed = func(err error) error {
		return Wrap(err, "v2 read memory.current")
	}

	// ErrCgroupV2MemoryCurrentParseFailed represents a function to generate an error that parsing memory.current failed.
	ErrCgroupV2MemoryCurrentParseFailed = func(err error) error {
		return Wrap(err, "v2 parse memory.current")
	}

	// ErrCgroupV2MemoryMaxReadFailed represents a function to generate an error that reading memory.max failed.
	ErrCgroupV2MemoryMaxReadFailed = func(err error) error {
		return Wrap(err, "v2 read memory.max")
	}

	// ErrCgroupV2MemoryMaxParseFailed represents a function to generate an error that parsing memory.max failed.
	ErrCgroupV2MemoryMaxParseFailed = func(err error) error {
		return Wrap(err, "v2 parse memory.max")
	}

	// ErrCgroupV2CPUStatReadFailed represents a function to generate an error that reading cpu.stat failed.
	ErrCgroupV2CPUStatReadFailed = func(err error) error {
		return Wrap(err, "v2 read cpu.stat")
	}

	// ErrCgroupV2CPUStatMissingUsage represents a function to generate an error that cpu.stat is missing usage_usec.
	ErrCgroupV2CPUStatMissingUsage = New("v2 cpu.stat missing usage_usec")

	// ErrCgroupV2CPUMaxReadFailed represents a function to generate an error that reading cpu.max failed.
	ErrCgroupV2CPUMaxReadFailed = func(err error) error {
		return Wrap(err, "v2 read cpu.max")
	}

	// ErrCgroupV2CPUMaxMalformed represents a function to generate an error that cpu.max is malformed.
	ErrCgroupV2CPUMaxMalformed = func(val string) error {
		return Errorf("v2 cpu.max malformed: %q", val)
	}

	// ErrCgroupV2CPUMaxParseQuotaFailed represents a function to generate an error that parsing cpu.max quota failed.
	ErrCgroupV2CPUMaxParseQuotaFailed = func(err error) error {
		return Wrap(err, "v2 cpu.max parse quota")
	}

	// ErrCgroupV2CPUMaxParsePeriodFailed represents a function to generate an error that parsing cpu.max period failed.
	ErrCgroupV2CPUMaxParsePeriodFailed = func(err error) error {
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
