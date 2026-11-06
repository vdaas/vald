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

// time.NewTicker panics for a non-positive interval, so the option must ignore it
// and leave the previously configured (default) value untouched without failing.
func TestWithExportIndexInfoDuration_NonPositive(t *testing.T) {
	t.Parallel()
	const initial = 7 * time.Second
	tests := []struct {
		name string
		dur  string
		want time.Duration
	}{
		{name: "zero is ignored", dur: "0s", want: initial},
		{name: "negative is ignored", dur: "-1s", want: initial},
		{name: "positive is applied", dur: "2s", want: 2 * time.Second},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			o := new(ngt)
			o.exportIndexInfoDuration = initial
			if err := WithExportIndexInfoDuration(tt.dur)(o); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if o.exportIndexInfoDuration != tt.want {
				t.Errorf("exportIndexInfoDuration = %v, want %v", o.exportIndexInfoDuration, tt.want)
			}
		})
	}
}
