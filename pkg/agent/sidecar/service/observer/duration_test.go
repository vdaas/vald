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

// time.NewTicker panics for a non-positive interval, and both values feed a ticker,
// so non-positive input must leave the previously configured (default) value untouched.
func TestObserverDurationOptions_NonPositive(t *testing.T) {
	t.Parallel()
	const initial = 7 * time.Second
	tests := []struct {
		name  string
		opt   func(string) Option
		field func(*observer) *time.Duration
		dur   string
		want  time.Duration
	}{
		{
			name:  "backup duration zero is ignored",
			opt:   WithBackupDuration,
			field: func(o *observer) *time.Duration { return &o.checkDuration },
			dur:   "0s",
			want:  initial,
		},
		{
			name:  "backup duration negative is ignored",
			opt:   WithBackupDuration,
			field: func(o *observer) *time.Duration { return &o.checkDuration },
			dur:   "-1s",
			want:  initial,
		},
		{
			name:  "backup duration positive is applied",
			opt:   WithBackupDuration,
			field: func(o *observer) *time.Duration { return &o.checkDuration },
			dur:   "2s",
			want:  2 * time.Second,
		},
		{
			name:  "post stop timeout zero is ignored",
			opt:   WithPostStopTimeout,
			field: func(o *observer) *time.Duration { return &o.postStopTimeout },
			dur:   "0s",
			want:  initial,
		},
		{
			name:  "post stop timeout negative is ignored",
			opt:   WithPostStopTimeout,
			field: func(o *observer) *time.Duration { return &o.postStopTimeout },
			dur:   "-1s",
			want:  initial,
		},
		{
			name:  "post stop timeout positive is applied",
			opt:   WithPostStopTimeout,
			field: func(o *observer) *time.Duration { return &o.postStopTimeout },
			dur:   "2s",
			want:  2 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			o := new(observer)
			*tt.field(o) = initial
			if err := tt.opt(tt.dur)(o); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := *tt.field(o); got != tt.want {
				t.Errorf("duration = %v, want %v", got, tt.want)
			}
		})
	}
}
