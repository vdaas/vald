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

package observer

import (
	"testing"
	"time"
)

// time.NewTicker panics for a non-positive interval; timeout/5 is 0 for 1ns-4ns,
// and timeutil.Parse accepts such values.
func TestPostStopCheckInterval(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		timeout time.Duration
		want    time.Duration
	}{
		{name: "negative is clamped to 1ns", timeout: -time.Second, want: time.Nanosecond},
		{name: "zero is clamped to 1ns", timeout: 0, want: time.Nanosecond},
		{name: "1ns is clamped to 1ns", timeout: time.Nanosecond, want: time.Nanosecond},
		{name: "4ns is clamped to 1ns", timeout: 4 * time.Nanosecond, want: time.Nanosecond},
		{name: "5ns gives 1ns", timeout: 5 * time.Nanosecond, want: time.Nanosecond},
		{name: "10s gives a fifth", timeout: 10 * time.Second, want: 2 * time.Second},
		{name: "2m is capped at 2s", timeout: 2 * time.Minute, want: 2 * time.Second},
		{name: "1s gives a fifth", timeout: time.Second, want: 200 * time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := postStopCheckInterval(tt.timeout)
			if got != tt.want {
				t.Errorf("postStopCheckInterval(%v) = %v, want %v", tt.timeout, got, tt.want)
			}
			if got <= 0 {
				t.Errorf("interval must be positive, got %v", got)
			}
		})
	}
}
