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
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test"
)

func TestSlot_Record_And_Reset(t *testing.T) {
	type record struct {
		rr  *RequestResult
		win uint64
	}

	if err := test.Run(t.Context(), t, func(tt *testing.T, rs []record) (Slot, error) {
		// Using real histograms/exemplars to fully test integration
		h, err := NewHistogram(WithHistogramNumBuckets(10))
		if err != nil {
			return nil, err
		}
		e, _ := NewExemplar(WithExemplarCapacity(10))
		s := newSlot(1, h, h, e) // numCounters=1

		for _, rec := range rs {
			s.Record(rec.rr, rec.win)
		}
		return s, nil
	}, []test.Case[Slot, []record]{
		{
			Name: "record within same window",
			Args: []record{
				{
					rr:  &RequestResult{EndedAt: time.Now(), Latency: 100 * time.Millisecond},
					win: 0,
				},
				{
					rr:  &RequestResult{EndedAt: time.Now(), Latency: 200 * time.Millisecond},
					win: 0,
				},
			},
			CheckFunc: func(tt *testing.T, want test.Result[Slot], got test.Result[Slot]) error {
				if got.Err != nil {
					return got.Err
				}
				s := got.Val.(*slot)
				if s.Total.Load() != 2 {
					return errors.Errorf("expected total 2, got %d", s.Total.Load())
				}
				if s.WindowStart != 0 {
					return errors.Errorf("expected WindowStart 0, got %d", s.WindowStart)
				}
				return nil
			},
		},
		{
			Name: "record with window change (reset)",
			Args: []record{
				{
					rr:  &RequestResult{EndedAt: time.Now()},
					win: 0,
				},
				{
					rr:  &RequestResult{EndedAt: time.Now()},
					win: 1,
				},
			},
			CheckFunc: func(tt *testing.T, want test.Result[Slot], got test.Result[Slot]) error {
				if got.Err != nil {
					return got.Err
				}
				s := got.Val.(*slot)
				// Should have reset, so only the record for window 1 exists
				if s.Total.Load() != 1 {
					return errors.Errorf("expected total 1 (reset), got %d", s.Total.Load())
				}
				if s.WindowStart != 1 {
					return errors.Errorf("expected WindowStart 1, got %d", s.WindowStart)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestSlot_Merge(t *testing.T) {
	type args struct {
		s1Init func(*slot)
		s2Init func(*slot)
	}

	if err := test.Run(t.Context(), t, func(tt *testing.T, args args) (Slot, error) {
		s1 := newSlot(1, nil, nil, nil).(*slot)
		if args.s1Init != nil {
			args.s1Init(s1)
		}
		s2 := newSlot(1, nil, nil, nil).(*slot)
		if args.s2Init != nil {
			args.s2Init(s2)
		}
		if err := s1.Merge(s2); err != nil {
			return nil, err
		}
		return s1, nil
	}, []test.Case[Slot, args]{
		{
			Name: "merge slots",
			Args: args{
				s1Init: func(s *slot) {
					s.Total.Store(10)
					s.Errors.Store(1)
					s.updatedNS.Store(100)
				},
				s2Init: func(s *slot) {
					s.Total.Store(20)
					s.Errors.Store(2)
					s.updatedNS.Store(200)
				},
			},
			CheckFunc: func(tt *testing.T, want test.Result[Slot], got test.Result[Slot]) error {
				if got.Err != nil {
					return got.Err
				}
				s := got.Val.(*slot)
				if s.Total.Load() != 30 {
					return errors.Errorf("expected Total 30, got %d", s.Total.Load())
				}
				if s.Errors.Load() != 3 {
					return errors.Errorf("expected Errors 3, got %d", s.Errors.Load())
				}
				if s.updatedNS.Load() != 200 {
					return errors.Errorf("expected updatedNS 200, got %d", s.updatedNS.Load())
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestSlot_Concurrent_Record(t *testing.T) {
	type args struct {
		workers int
		loops   int
	}

	if err := test.Run(t.Context(), t, func(tt *testing.T, args args) (Slot, error) {
		s := newSlot(0, nil, nil, nil)
		var wg sync.WaitGroup
		start := time.Now()
		for i := 0; i < args.workers; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < args.loops; j++ {
					// Switch window index every 10 iterations to force rapid resets
					win := uint64((id*args.loops + j) / 10) //nolint:gosec // loop counters are positive
					s.Record(&RequestResult{
						EndedAt: start,
						Status:  codes.OK,
					}, win)
				}
			}(i)
		}
		wg.Wait()
		return s, nil
	}, []test.Case[Slot, args]{
		{
			Name: "concurrent stress",
			Args: args{
				workers: 10,
				loops:   100,
			},
			CheckFunc: func(tt *testing.T, want test.Result[Slot], got test.Result[Slot]) error {
				if got.Err != nil {
					return got.Err
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestSlot_Clone(t *testing.T) {
	if err := test.Run(t.Context(), t, func(tt *testing.T, total uint64) (Slot, error) {
		s1 := newSlot(1, nil, nil, nil).(*slot)
		s1.Total.Store(total)
		s2 := s1.Clone()
		// Modify s1
		s1.Total.Store(total + 1)
		return s2, nil
	}, []test.Case[Slot, uint64]{
		{
			Name: "clone independence",
			Args: 10,
			CheckFunc: func(tt *testing.T, want test.Result[Slot], got test.Result[Slot]) error {
				if got.Err != nil {
					return got.Err
				}
				s := got.Val.(*slot)
				if s.Total.Load() != 10 {
					return errors.Errorf("expected cloned Total 10, got %d", s.Total.Load())
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}
