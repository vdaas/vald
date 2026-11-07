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

package mysql

import (
	"testing"
	"time"
)

// time.NewTicker panics for a non-positive interval, so the option must never store one.
func TestWithInitialPingDuration_NonPositive(t *testing.T) {
	t.Parallel()
	// invalidDuration does not parse as a time.Duration.
	const invalidDuration = "invalid"
	tests := []struct {
		name string
		dur  string
		want time.Duration
	}{
		{name: "zero falls back to 50ms", dur: "0s", want: 50 * time.Millisecond},
		{name: "negative falls back to 50ms", dur: "-1s", want: 50 * time.Millisecond},
		{name: "invalid falls back to 50ms", dur: invalidDuration, want: 50 * time.Millisecond},
		{name: "positive is kept", dur: "10ms", want: 10 * time.Millisecond},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := new(mySQLClient)
			if err := WithInitialPingDuration(tt.dur)(m); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if m.initialPingDuration != tt.want {
				t.Errorf("initialPingDuration = %v, want %v", m.initialPingDuration, tt.want)
			}
		})
	}
}
