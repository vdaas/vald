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

package service

import (
	"testing"
	"time"
)

// time.NewTicker panics for a non-positive interval, so the option must never store one.
func TestWithDiscoverDuration_NonPositive(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		dur  string
		want time.Duration
	}{
		{name: "zero falls back to 1s", dur: "0s", want: time.Second},
		{name: "negative falls back to 1s", dur: "-1s", want: time.Second},
		{name: "invalid falls back to 1s", dur: "invalid", want: time.Second},
		{name: "positive is kept", dur: "3s", want: 3 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			d := new(discoverer)
			if err := WithDiscoverDuration(tt.dur)(d); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if d.csd != tt.want {
				t.Errorf("csd = %v, want %v", d.csd, tt.want)
			}
		})
	}
}
