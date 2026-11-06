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

package cache

import (
	"testing"
	"time"
)

// A zero/negative expire check duration must not overwrite the default (5m):
// it is later handed to a ticker, which panics for a non-positive interval.
func TestWithExpireCheckDuration_NonPositive(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		dur  string
		want time.Duration
	}{
		{name: "zero keeps previous value", dur: "0s", want: 7 * time.Second},
		{name: "negative keeps previous value", dur: "-1s", want: 7 * time.Second},
		{name: "positive is applied", dur: "2s", want: 2 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := new(cache[any])
			c.expireCheckDur = 7 * time.Second
			WithExpireCheckDuration[any](tt.dur)(c)
			if c.expireCheckDur != tt.want {
				t.Errorf("expireCheckDur = %v, want %v", c.expireCheckDur, tt.want)
			}
		})
	}
}
