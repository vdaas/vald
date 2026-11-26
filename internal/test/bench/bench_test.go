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

package bench

import "testing"

func TestFull(t *testing.T) {
	oldFlag := *fullFlag
	*fullFlag = false
	t.Cleanup(func() {
		*fullFlag = oldFlag
	})

	tests := []struct {
		name    string
		envVal  string
		flagVal bool
		want    bool
	}{
		{name: "default empty", envVal: "", flagVal: false, want: false},
		{name: "env false", envVal: "false", flagVal: false, want: false},
		{name: "env zero", envVal: "0", flagVal: false, want: false},
		{name: "env invalid", envVal: "invalid", flagVal: false, want: false},
		{name: "env true", envVal: "true", flagVal: false, want: true},
		{name: "env one", envVal: "1", flagVal: false, want: true},
		{name: "env uppercase TRUE", envVal: "TRUE", flagVal: false, want: true},
		{name: "flag true overrides empty env", envVal: "", flagVal: true, want: true},
		{name: "flag true overrides false env", envVal: "false", flagVal: true, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(FullEnvVar, tt.envVal)
			*fullFlag = tt.flagVal
			if got := Full(); got != tt.want {
				t.Fatalf("Full() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipUnlessFull(t *testing.T) {
	oldFlag := *fullFlag
	*fullFlag = false
	t.Cleanup(func() {
		*fullFlag = oldFlag
	})

	t.Run("skips when fast mode", func(t *testing.T) {
		t.Setenv(FullEnvVar, "false")
		*fullFlag = false
		bodyRan := false
		testing.Benchmark(func(b *testing.B) {
			b.Helper()
			SkipUnlessFull(b)
			bodyRan = true
		})
		if bodyRan {
			t.Fatal("SkipUnlessFull did not skip benchmark in fast mode")
		}
	})

	t.Run("runs when full mode via env", func(t *testing.T) {
		t.Setenv(FullEnvVar, "true")
		*fullFlag = false
		bodyRan := false
		testing.Benchmark(func(b *testing.B) {
			b.Helper()
			SkipUnlessFull(b)
			bodyRan = true
		})
		if !bodyRan {
			t.Fatal("SkipUnlessFull skipped benchmark in full mode")
		}
	})
}
