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

// This file encodes the TDD RED-phase test contract for the not-yet-implemented
// (*GlobalSnapshot).AchievedQPS() method and intentionally does not compile yet.
// It is excluded from the default build via the `ignore` tag below for the same
// reason documented in recall_test.go: this repository's Stop-hook / golangci-lint
// gate has no TDD RED-phase exception for typecheck failures. The next Maker
// should delete the `//go:build ignore` line (only) once AchievedQPS() is
// implemented on GlobalSnapshot.

package metrics

// RED PHASE (TDD Test Maker): this file locks the contract for a new method,
// (*GlobalSnapshot).AchievedQPS() float64, which does not exist yet.
//
// Design decision: AchievedQPS is computed as Total / LastUpdated.Sub(StartTime).Seconds().
// It must be defensive (nil receiver, zero Total, zero/degenerate time window) because
// GlobalSnapshot values routinely flow through JSON (de)serialization from files written
// by earlier/partial test runs (see the Pareto chart tooling in
// tests/v2/e2e/metrics/chart), where StartTime/LastUpdated may be zero-valued or, in
// pathological clock-skew cases, LastUpdated could even precede StartTime. In every one
// of those degenerate cases AchievedQPS must return 0, never NaN or +/-Inf, so that
// downstream chart rendering (log-scale Y axis) never receives a non-finite value.
import (
	"math"
	"testing"
	"time"
)

func TestGlobalSnapshot_AchievedQPS(t *testing.T) {
	t.Parallel()

	base := time.Date(2026, 7, 18, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		snap *GlobalSnapshot
		name string
		want float64
	}{
		{
			name: "nil receiver returns 0",
			snap: nil,
			want: 0,
		},
		{
			name: "zero Total returns 0 even with a valid time window",
			snap: &GlobalSnapshot{
				Total:       0,
				StartTime:   base,
				LastUpdated: base.Add(10 * time.Second),
			},
			want: 0,
		},
		{
			name: "zero StartTime (never started) returns 0",
			snap: &GlobalSnapshot{
				Total:       100,
				StartTime:   time.Time{},
				LastUpdated: base.Add(10 * time.Second),
			},
			want: 0,
		},
		{
			name: "zero LastUpdated (never updated) returns 0",
			snap: &GlobalSnapshot{
				Total:       100,
				StartTime:   base,
				LastUpdated: time.Time{},
			},
			want: 0,
		},
		{
			name: "LastUpdated == StartTime (zero duration) returns 0, not +Inf",
			snap: &GlobalSnapshot{
				Total:       100,
				StartTime:   base,
				LastUpdated: base,
			},
			want: 0,
		},
		{
			name: "LastUpdated before StartTime (negative duration / clock skew) returns 0",
			snap: &GlobalSnapshot{
				Total:       100,
				StartTime:   base,
				LastUpdated: base.Add(-5 * time.Second),
			},
			want: 0,
		},
		{
			name: "normal case: 1000 requests over 10s => 100 QPS",
			snap: &GlobalSnapshot{
				Total:       1000,
				StartTime:   base,
				LastUpdated: base.Add(10 * time.Second),
			},
			want: 100,
		},
		{
			name: "sub-second window: 1 request over 500ms => 2 QPS",
			snap: &GlobalSnapshot{
				Total:       1,
				StartTime:   base,
				LastUpdated: base.Add(500 * time.Millisecond),
			},
			want: 2,
		},
		{
			name: "fractional QPS: 3 requests over 2s => 1.5 QPS",
			snap: &GlobalSnapshot{
				Total:       3,
				StartTime:   base,
				LastUpdated: base.Add(2 * time.Second),
			},
			want: 1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.snap.AchievedQPS()
			if math.IsNaN(got) {
				t.Fatalf("AchievedQPS() = NaN, want a finite value (this would poison a log-scale chart axis)")
			}
			if math.IsInf(got, 0) {
				t.Fatalf("AchievedQPS() = %v (Inf), want a finite value", got)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("AchievedQPS() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("very small duration does not overflow to Inf/NaN", func(t *testing.T) {
		t.Parallel()
		snap := &GlobalSnapshot{
			Total:       1_000_000,
			StartTime:   base,
			LastUpdated: base.Add(1 * time.Nanosecond),
		}
		got := snap.AchievedQPS()
		if math.IsNaN(got) || math.IsInf(got, 0) {
			t.Fatalf("AchievedQPS() = %v, want a finite (if very large) value", got)
		}
		if got <= 0 {
			t.Errorf("AchievedQPS() = %v, want > 0 for a positive Total over a positive (if tiny) duration", got)
		}
	})

	t.Run("is a pure/idempotent computation over the same snapshot", func(t *testing.T) {
		t.Parallel()
		snap := &GlobalSnapshot{
			Total:       500,
			StartTime:   base,
			LastUpdated: base.Add(5 * time.Second),
		}
		first := snap.AchievedQPS()
		second := snap.AchievedQPS()
		if first != second {
			t.Errorf("AchievedQPS() is not idempotent: first=%v second=%v", first, second)
		}
		if snap.Total != 500 || !snap.StartTime.Equal(base) {
			t.Errorf("AchievedQPS() must not mutate the snapshot: %+v", snap)
		}
	})
}

// TestGlobalSnapshot_AchievedQPS_IntegrationWithCollector exercises AchievedQPS()
// against a snapshot produced by a real Collector (not a hand-built struct literal),
// to make sure StartTime/LastUpdated wiring from Record() feeds AchievedQPS()
// correctly end to end.
func TestGlobalSnapshot_AchievedQPS_IntegrationWithCollector(t *testing.T) {
	t.Parallel()

	c, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	start := time.Now()
	for i := 0; i < 5; i++ {
		c.Record(t.Context(), 0, &RequestResult{
			StartedAt: start.Add(time.Duration(i) * time.Second),
			EndedAt:   start.Add(time.Duration(i)*time.Second + time.Millisecond),
		})
	}

	snap := c.GlobalSnapshot()
	if snap.Total != 5 {
		t.Fatalf("precondition failed: Total = %d, want 5", snap.Total)
	}

	got := snap.AchievedQPS()
	if math.IsNaN(got) || math.IsInf(got, 0) {
		t.Fatalf("AchievedQPS() = %v, want finite", got)
	}
	// 5 requests spread across ~4 seconds (StartedAt at t+0..t+4) => ~1.25 QPS.
	// Use a loose tolerance since the exact window depends on EndedAt of the last
	// request rather than a fixed clock.
	if got <= 0 {
		t.Errorf("AchievedQPS() = %v, want > 0", got)
	}
}
