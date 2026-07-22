//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

package params

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
)

// Test_parser_Parse_FiltersDoNotMutateOSArgs is the regression test for the
// slices.Clone fix in parser.Parse: slices.DeleteFunc compacts in place and
// zeroes the removed tail, so filtering os.Args[1:] without cloning first
// destroyed the filtered arguments (e.g. go test's -test.* flags) inside the
// global os.Args for every later reader, such as testing's own flag parsing.
func Test_parser_Parse_FiltersDoNotMutateOSArgs(t *testing.T) {
	orig := os.Args
	t.Cleanup(func() { os.Args = orig })
	os.Args = []string{
		"parser_test",
		"-test.run=TestNothing",
		"-config", "/tmp/params-test-nonexistent.yaml",
		"-test.v=true",
	}
	want := make([]string, len(os.Args))
	copy(want, os.Args)

	// Parse fails on the nonexistent config path, which is irrelevant here:
	// the filter must have run (and must not have corrupted os.Args) either way.
	_, _, err := New(
		WithName("parser_test"),
		WithArgumentFilters(func(s string) bool {
			return strings.HasPrefix(s, "-test.")
		}),
	).Parse()
	t.Logf("Parse returned err: %v", err)

	if !reflect.DeepEqual(os.Args, want) {
		t.Errorf("os.Args was mutated by Parse with argument filters:\n got: %#v\nwant: %#v", os.Args, want)
	}
}
