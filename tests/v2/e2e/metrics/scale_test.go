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
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test"
)

func TestScale_Record_And_Reset(t *testing.T) {
	type args struct {
		name     string
		records  []*RequestResult
		width    uint64
		capacity uint64
		st       ScaleType
	}

	if err := test.Run(t.Context(), t, func(tt *testing.T, args args) (*ScaleSnapshot, error) {
		s, err := newScale(args.name, args.width, args.capacity, 0, args.st, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		for _, r := range args.records {
			s.Record(context.Background(), r)
		}
		return s.Snapshot(), nil
	}, []test.Case[*ScaleSnapshot, args]{
		{
			Name: "time scale basic logic",
			Args: args{
				name:     "time_test",
				width:    uint64(time.Second), // 1 second per slot
				capacity: 2,                   // 2 slots: 0 and 1
				st:       TimeScale,
				records: []*RequestResult{
					{EndedAt: time.Unix(0, 0)}, // Slot 0
					{EndedAt: time.Unix(1, 0)}, // Slot 1
					{EndedAt: time.Unix(2, 0)}, // Slot 0 (Reset)
				},
			},
			CheckFunc: func(tt *testing.T, want test.Result[*ScaleSnapshot], got test.Result[*ScaleSnapshot]) error {
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				// Slot 0 should have 1 (from time 2), Slot 1 should have 1 (from time 1)
				if snap.Slots[0].Total != 1 {
					return errors.Errorf("expected slot 0 total 1, got %d", snap.Slots[0].Total)
				}
				if snap.Slots[0].LastUpdated != time.Unix(2, 0).UnixNano() {
					return errors.Errorf("expected slot 0 update %d, got %d", time.Unix(2, 0).UnixNano(), snap.Slots[0].LastUpdated)
				}
				if snap.Slots[1].Total != 1 {
					return errors.Errorf("expected slot 1 total 1, got %d", snap.Slots[1].Total)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestScale_Concurrency(t *testing.T) {
	type args struct {
		workers  int
		loops    int
		capacity uint64
	}

	if err := test.Run(t.Context(), t, func(tt *testing.T, args args) (*ScaleSnapshot, error) {
		s, err := newScale("concurrent", uint64(time.Second), args.capacity, 0, TimeScale, nil, nil, nil)
		if err != nil {
			return nil, err
		}

		start := time.Now()
		var wg sync.WaitGroup

		for i := 0; i < args.workers; i++ {
			wg.Add(1)
			go func(offset int) {
				defer wg.Done()
				for j := 0; j < args.loops; j++ {
					// Use fake times that increment to force wrapping
					now := start.Add(time.Duration(offset+j) * time.Second)
					s.Record(context.Background(), &RequestResult{
						EndedAt: now,
					})
				}
			}(i)
		}
		wg.Wait()
		return s.Snapshot(), nil
	}, []test.Case[*ScaleSnapshot, args]{
		{
			Name: "concurrent write and wrap",
			Args: args{
				workers:  20,
				loops:    100,
				capacity: 5,
			},
			CheckFunc: func(tt *testing.T, want test.Result[*ScaleSnapshot], got test.Result[*ScaleSnapshot]) error {
				if got.Err != nil {
					return got.Err
				}
				if got.Val == nil {
					return errors.New("got nil snapshot")
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestScale_Merge(t *testing.T) {
	type args struct {
		s1Recs    []*RequestResult
		s2Recs    []*RequestResult
		width1    uint64
		capacity1 uint64
		width2    uint64
		capacity2 uint64
	}

	if err := test.Run(t.Context(), t, func(tt *testing.T, args args) (*ScaleSnapshot, error) {
		s1, err := newScale("test", args.width1, args.capacity1, 0, TimeScale, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		s2, err := newScale("test", args.width2, args.capacity2, 0, TimeScale, nil, nil, nil)
		if err != nil {
			return nil, err
		}

		for _, r := range args.s1Recs {
			s1.Record(context.Background(), r)
		}
		for _, r := range args.s2Recs {
			s2.Record(context.Background(), r)
		}

		if err := s1.Merge(s2); err != nil {
			return nil, err
		}
		return s1.Snapshot(), nil
	}, []test.Case[*ScaleSnapshot, args]{
		{
			Name: "merge compatible scales",
			Args: args{
				width1: uint64(time.Second), capacity1: 2,
				width2: uint64(time.Second), capacity2: 2,
				s1Recs: []*RequestResult{{EndedAt: time.Unix(0, 0)}},
				s2Recs: []*RequestResult{
					{EndedAt: time.Unix(0, 0)},
					{EndedAt: time.Unix(1, 0)},
				},
			},
			CheckFunc: func(tt *testing.T, want test.Result[*ScaleSnapshot], got test.Result[*ScaleSnapshot]) error {
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				// Slot 0: 1 (s1) + 1 (s2) = 2
				if snap.Slots[0].Total != 2 {
					return errors.Errorf("expected slot 0 total 2, got %d", snap.Slots[0].Total)
				}
				// Slot 1: 0 (s1) + 1 (s2) = 1
				if snap.Slots[1].Total != 1 {
					return errors.Errorf("expected slot 1 total 1, got %d", snap.Slots[1].Total)
				}
				return nil
			},
		},
		{
			Name: "merge incompatible scales",
			Args: args{
				width1: uint64(time.Second), capacity1: 2,
				width2: uint64(time.Second * 2), capacity2: 2, // different width
			},
			CheckFunc: func(tt *testing.T, want test.Result[*ScaleSnapshot], got test.Result[*ScaleSnapshot]) error {
				if got.Err == nil {
					return errors.New("expected error, got nil")
				}
				if got.Err.Error() != "incompatible scales" {
					return errors.Errorf("unexpected error message: %v", got.Err)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestScale_Clone(t *testing.T) {
	if err := test.Run(t.Context(), t, func(tt *testing.T, rr []*RequestResult) (*ScaleSnapshot, error) {
		s1, err := newScale("test", uint64(time.Second), 2, 0, TimeScale, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		for _, r := range rr {
			s1.Record(context.Background(), r)
		}
		s2 := s1.Clone()
		// Modify s1
		s1.Record(context.Background(), &RequestResult{EndedAt: time.Unix(0, 0)})
		return s2.Snapshot(), nil
	}, []test.Case[*ScaleSnapshot, []*RequestResult]{
		{
			Name: "clone independence",
			Args: []*RequestResult{{EndedAt: time.Unix(0, 0)}},
			CheckFunc: func(tt *testing.T, want test.Result[*ScaleSnapshot], got test.Result[*ScaleSnapshot]) error {
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				// s1 was 1, cloned. s1 became 2. s2 should remain 1.
				if snap.Slots[0].Total != 1 {
					return errors.Errorf("expected cloned slot 0 total 1, got %d", snap.Slots[0].Total)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestScale_RingBuffer_WrapAndReset(t *testing.T) {
	type args struct {
		records  []*RequestResult
		width    uint64
		capacity uint64
	}
	// This test verifies wrapping logic specifically
	if err := test.Run(t.Context(), t, func(tt *testing.T, args args) (*ScaleSnapshot, error) {
		s, err := newScale("test_wrap", args.width, args.capacity, 0, TimeScale, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		for _, r := range args.records {
			s.Record(context.Background(), r)
		}
		return s.Snapshot(), nil
	}, []test.Case[*ScaleSnapshot, args]{
		{
			Name: "wrap around and reset",
			Args: args{
				width:    uint64(time.Second),
				capacity: 3,
				records: []*RequestResult{
					{EndedAt: time.Unix(0, 0)}, // Slot 0
					{EndedAt: time.Unix(1, 0)}, // Slot 1
					{EndedAt: time.Unix(2, 0)}, // Slot 2
					{EndedAt: time.Unix(3, 0)}, // Slot 0 (Reset!)
				},
			},
			CheckFunc: func(tt *testing.T, want test.Result[*ScaleSnapshot], got test.Result[*ScaleSnapshot]) error {
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				// Slot 0 should have 1 (from sec 3), not 2
				if snap.Slots[0].Total != 1 {
					return errors.Errorf("Slot 0: expected 1, got %d", snap.Slots[0].Total)
				}
				if snap.Slots[0].LastUpdated != time.Unix(3, 0).UnixNano() {
					return errors.Errorf("Slot 0 time: expected %d, got %d", time.Unix(3, 0).UnixNano(), snap.Slots[0].LastUpdated)
				}
				// Slot 1 should have 1 (from sec 1)
				if snap.Slots[1].Total != 1 {
					return errors.Errorf("Slot 1: expected 1, got %d", snap.Slots[1].Total)
				}
				// Slot 2 should have 1 (from sec 2)
				if snap.Slots[2].Total != 1 {
					return errors.Errorf("Slot 2: expected 1, got %d", snap.Slots[2].Total)
				}
				return nil
			},
		},
		{
			Name: "wrap around skip slot",
			Args: args{
				width:    uint64(time.Second),
				capacity: 3,
				records: []*RequestResult{
					{EndedAt: time.Unix(0, 0)}, // Slot 0
					{EndedAt: time.Unix(1, 0)}, // Slot 1
					{EndedAt: time.Unix(2, 0)}, // Slot 2
					{EndedAt: time.Unix(3, 0)}, // Slot 0 (Reset)
					{EndedAt: time.Unix(5, 0)}, // Slot 2 (Reset) - skipping Slot 1
				},
			},
			CheckFunc: func(tt *testing.T, want test.Result[*ScaleSnapshot], got test.Result[*ScaleSnapshot]) error {
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				// Slot 2 should have 1 (from sec 5)
				if snap.Slots[2].Total != 1 {
					return errors.Errorf("Slot 2: expected 1, got %d", snap.Slots[2].Total)
				}
				if snap.Slots[2].LastUpdated != time.Unix(5, 0).UnixNano() {
					return errors.Errorf("Slot 2 time: expected %d, got %d", time.Unix(5, 0).UnixNano(), snap.Slots[2].LastUpdated)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func BenchmarkScale_Record(b *testing.B) {
	s, _ := newScale("bench", uint64(time.Second), 10, 0, TimeScale, nil, nil, nil)
	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			s.Record(ctx, &RequestResult{
				EndedAt: time.Unix(int64(i), 0),
			})
			i++
		}
	})
}
