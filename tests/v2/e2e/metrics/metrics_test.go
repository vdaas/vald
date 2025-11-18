//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
)

func TestCollector(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		collector func() (Collector, error)
		records   []*RequestResult
		merge     func() (Collector, error)
		check     func(t *testing.T, c Collector)
	}

	tests := []testCase{
		{
			name: "record a single successful request",
			collector: func() (Collector, error) {
				return NewCollector()
			},
			records: []*RequestResult{
				{
					Latency:   100 * time.Millisecond,
					QueueWait: 20 * time.Millisecond,
				},
			},
			check: func(t *testing.T, c Collector) {
				snap := c.GlobalSnapshot()
				if snap.Total != 1 {
					t.Errorf("expected total 1, got %d", snap.Total)
				}
				if snap.Errors != 0 {
					t.Errorf("expected errors 0, got %d", snap.Errors)
				}
				if snap.Latencies.Total != 1 {
					t.Errorf("expected latencies total 1, got %d", snap.Latencies.Total)
				}
			},
		},
		{
			name: "record a single errored request",
			collector: func() (Collector, error) {
				return NewCollector()
			},
			records: []*RequestResult{
				{
					Latency:   100 * time.Millisecond,
					QueueWait: 20 * time.Millisecond,
					Err:       errors.New("test error"),
				},
			},
			check: func(t *testing.T, c Collector) {
				snap := c.GlobalSnapshot()
				if snap.Total != 1 {
					t.Errorf("expected total 1, got %d", snap.Total)
				}
				if snap.Errors != 1 {
					t.Errorf("expected errors 1, got %d", snap.Errors)
				}
			},
		},
		{
			name: "merge two collectors",
			collector: func() (Collector, error) {
				return NewCollector(WithCustomCounters("c1"))
			},
			records: []*RequestResult{
				{
					Latency: 100 * time.Millisecond,
				},
			},
			merge: func() (Collector, error) {
				c, err := NewCollector(WithCustomCounters("c1", "c2"))
				if err != nil {
					return nil, err
				}
				c.Record(context.Background(), &RequestResult{Latency: 200, Err: errors.New("err")})
				h1, _ := c.CounterHandle("c1")
				h1.Inc()
				h2, _ := c.CounterHandle("c2")
				h2.Inc()
				return c, nil
			},
			check: func(t *testing.T, c Collector) {
				snap := c.GlobalSnapshot()
				if snap.Total != 2 {
					t.Errorf("expected total 2, got %d", snap.Total)
				}
				if snap.Errors != 1 {
					t.Errorf("expected errors 1, got %d", snap.Errors)
				}
				h1, _ := c.CounterHandle("c1")
				if h1.value.Load() != 1 {
					t.Errorf("expected c1 counter to be 1, got %d", h1.value.Load())
				}
				h2, _ := c.CounterHandle("c2")
				if h2.value.Load() != 1 {
					t.Errorf("expected c2 counter to be 1, got %d", h2.value.Load())
				}
			},
		},
		{
			name: "record multiple successful and errored requests",
			collector: func() (Collector, error) {
				return NewCollector()
			},
			records: []*RequestResult{
				{
					Latency:   100 * time.Millisecond,
					QueueWait: 20 * time.Millisecond,
				},
				{
					Latency:   200 * time.Millisecond,
					QueueWait: 30 * time.Millisecond,
					Err:       errors.New("test error"),
				},
				{
					Latency:   300 * time.Millisecond,
					QueueWait: 40 * time.Millisecond,
				},
			},
			check: func(t *testing.T, c Collector) {
				snap := c.GlobalSnapshot()
				if snap.Total != 3 {
					t.Errorf("expected total 3, got %d", snap.Total)
				}
				if snap.Errors != 1 {
					t.Errorf("expected errors 1, got %d", snap.Errors)
				}
				if snap.Latencies.Total != 3 {
					t.Errorf("expected latencies total 3, got %d", snap.Latencies.Total)
				}
				if snap.QueueWaits.Total != 3 {
					t.Errorf("expected queue waits total 3, got %d", snap.QueueWaits.Total)
				}
			},
		},
		{
			name: "merge with an empty collector",
			collector: func() (Collector, error) {
				return NewCollector()
			},
			records: []*RequestResult{
				{
					Latency: 100 * time.Millisecond,
				},
			},
			merge: func() (Collector, error) {
				return NewCollector()
			},
			check: func(t *testing.T, c Collector) {
				snap := c.GlobalSnapshot()
				if snap.Total != 1 {
					t.Errorf("expected total 1, got %d", snap.Total)
				}
				if snap.Errors != 0 {
					t.Errorf("expected errors 0, got %d", snap.Errors)
				}
			},
		},
		{
			name: "global snapshot aggregates data correctly",
			collector: func() (Collector, error) {
				return NewCollector()
			},
			records: []*RequestResult{
				{
					Latency:   100 * time.Millisecond,
					QueueWait: 10 * time.Millisecond,
				},
				{
					Latency:   200 * time.Millisecond,
					QueueWait: 20 * time.Millisecond,
					Err:       errors.New("error"),
				},
			},
			check: func(t *testing.T, c Collector) {
				snap := c.GlobalSnapshot()
				if snap.Total != 2 {
					t.Errorf("expected total 2, got %d", snap.Total)
				}
				if snap.Errors != 1 {
					t.Errorf("expected errors 1, got %d", snap.Errors)
				}
				if snap.Latencies.Mean != float64(150*time.Millisecond) {
					t.Errorf("expected latency mean %v, got %v", 150*time.Millisecond, time.Duration(snap.Latencies.Mean))
				}
				if snap.QueueWaits.Mean != float64(15*time.Millisecond) {
					t.Errorf("expected queue wait mean %v, got %v", 15*time.Millisecond, time.Duration(snap.QueueWaits.Mean))
				}
			},
		},
		{
			name: "NewCollector with WithTimeScale does not panic",
			collector: func() (Collector, error) {
				return NewCollector(WithTimeScale("test", 1, 10))
			},
			records: []*RequestResult{},
			check: func(t *testing.T, c Collector) {
				if c == nil {
					t.Fatal("collector should not be nil")
				}
				coll := c.(*collector)
				if len(coll.scales) != 1 {
					t.Errorf("expected 1 scale, got %d", len(coll.scales))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := tt.collector()
			if err != nil {
				t.Fatalf("failed to create collector: %v", err)
			}

			for _, r := range tt.records {
				c.Record(context.Background(), r)
			}

			if tt.merge != nil {
				mc, err := tt.merge()
				if err != nil {
					t.Fatalf("failed to create collector for merge: %v", err)
				}
				if err := c.Merge(mc); err != nil {
					t.Fatalf("failed to merge collectors: %v", err)
				}
			}

			tt.check(t, c)
		})
	}
}
