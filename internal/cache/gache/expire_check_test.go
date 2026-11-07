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

package gache

import (
	"context"
	"testing"
	"time"

	gache "github.com/kpango/gache/v2"
)

// startExpiredRecorder records StartExpired calls without starting a goroutine.
type startExpiredRecorder struct {
	gache.Gache[any]
	durs []time.Duration
}

func (r *startExpiredRecorder) StartExpired(_ context.Context, dur time.Duration) gache.Gache[any] {
	r.durs = append(r.durs, dur)
	return r
}

// gache.StartExpired builds a time.NewTicker from the duration inside a goroutine,
// which panics (and kills the process) for a non-positive interval.
func TestCache_Start_NonPositiveExpireCheckDuration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		dur    time.Duration
		called bool
	}{
		{name: "unset (zero) does not start the expire daemon", dur: 0},
		{name: "negative does not start the expire daemon", dur: -time.Second},
		{name: "positive starts the expire daemon", dur: time.Second, called: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rec := new(startExpiredRecorder)
			c := &cache[any]{gache: rec, expireCheckDur: tt.dur}
			c.Start(t.Context())
			if got := len(rec.durs) != 0; got != tt.called {
				t.Fatalf("StartExpired called = %v, want %v", got, tt.called)
			}
			if tt.called && rec.durs[0] != tt.dur {
				t.Errorf("StartExpired dur = %v, want %v", rec.durs[0], tt.dur)
			}
		})
	}
}

func TestWithExpireCheckDuration_NonPositive(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		dur  time.Duration
		want time.Duration
	}{
		{name: "zero keeps previous value", dur: 0, want: 7 * time.Second},
		{name: "negative keeps previous value", dur: -time.Second, want: 7 * time.Second},
		{name: "positive is applied", dur: 2 * time.Second, want: 2 * time.Second},
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
