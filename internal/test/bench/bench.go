// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package bench provides shared controls for selecting fast or full benchmark suites.
package bench

import (
	"flag"
	"os"
	"strconv"
	"testing"
)

// FullEnvVar is the environment variable name that enables full benchmark suites.
const FullEnvVar = "VALD_BENCH_FULL"

//nolint:gochecknoglobals // shared benchmark CLI flag
var fullFlag = flag.Bool("bench-full", false, "run full benchmark matrices and heavy datasets")

// Full reports whether heavy / full-matrix benchmarks are enabled via -bench-full or VALD_BENCH_FULL.
func Full() bool {
	if fullFlag != nil && *fullFlag {
		return true
	}
	v, err := strconv.ParseBool(os.Getenv(FullEnvVar))
	return err == nil && v
}

// SkipUnlessFull skips b when running in fast (default) benchmark mode.
func SkipUnlessFull(b *testing.B) {
	b.Helper()
	if !Full() {
		b.Skipf("skipping heavy benchmark in fast mode (set %s=true or BENCH_FULL=true to run)", FullEnvVar)
	}
}
