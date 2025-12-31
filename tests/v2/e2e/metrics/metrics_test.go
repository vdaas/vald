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

package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test"
)

func TestNewCollector(t *testing.T) {
	if err := test.Run(t.Context(), t, func(t *testing.T, opts []Option) (Collector, error) {
		t.Helper()
		return NewCollector(opts...)
	}, []test.Case[Collector, []Option]{
		{
			Name: "initialize with default options",
			Args: nil,
			CheckFunc: func(t *testing.T, want test.Result[Collector], got test.Result[Collector]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				if got.Val == nil {
					return errors.New("got nil Collector")
				}
				return nil
			},
		},
		{
			Name: "initialize with time scale",
			Args: []Option{
				WithTimeScale("test_scale", time.Second, 10),
			},
			CheckFunc: func(t *testing.T, want test.Result[Collector], got test.Result[Collector]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				c := got.Val.(*collector)
				if len(c.scales) != 1 {
					return errors.Errorf("expected 1 scale, got %d", len(c.scales))
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestCollector_Record_And_Snapshot(t *testing.T) {
	type args struct {
		opts    []Option
		records []*RequestResult
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (*GlobalSnapshot, error) {
		t.Helper()
		c, err := NewCollector(args.opts...)
		if err != nil {
			return nil, err
		}
		for _, r := range args.records {
			c.Record(context.Background(), 0, r)
		}
		return c.GlobalSnapshot(), nil
	}, []test.Case[*GlobalSnapshot, args]{
		{
			Name: "record single success",
			Args: args{
				records: []*RequestResult{
					{
						Latency:   100 * time.Millisecond,
						QueueWait: 20 * time.Millisecond,
					},
				},
			},
			CheckFunc: func(t *testing.T, want test.Result[*GlobalSnapshot], got test.Result[*GlobalSnapshot]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				if snap.Total != 1 {
					return errors.Errorf("expected total 1, got %d", snap.Total)
				}
				if snap.Errors != 0 {
					return errors.Errorf("expected errors 0, got %d", snap.Errors)
				}
				if snap.Latencies.Total != 1 {
					return errors.Errorf("expected latencies total 1, got %d", snap.Latencies.Total)
				}
				return nil
			},
		},
		{
			Name: "record single error",
			Args: args{
				records: []*RequestResult{
					{
						Latency:   100 * time.Millisecond,
						QueueWait: 20 * time.Millisecond,
						Err:       errors.New("test error"),
					},
				},
			},
			CheckFunc: func(t *testing.T, want test.Result[*GlobalSnapshot], got test.Result[*GlobalSnapshot]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				if snap.Total != 1 {
					return errors.Errorf("expected total 1, got %d", snap.Total)
				}
				if snap.Errors != 1 {
					return errors.Errorf("expected errors 1, got %d", snap.Errors)
				}
				return nil
			},
		},
		{
			Name: "record multiple requests",
			Args: args{
				records: []*RequestResult{
					{Latency: 100 * time.Millisecond, QueueWait: 20 * time.Millisecond},
					{Latency: 200 * time.Millisecond, QueueWait: 30 * time.Millisecond, Err: errors.New("err")},
					{Latency: 300 * time.Millisecond, QueueWait: 40 * time.Millisecond},
				},
			},
			CheckFunc: func(t *testing.T, want test.Result[*GlobalSnapshot], got test.Result[*GlobalSnapshot]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				snap := got.Val
				if snap.Total != 3 {
					return errors.Errorf("expected total 3, got %d", snap.Total)
				}
				if snap.Errors != 1 {
					return errors.Errorf("expected errors 1, got %d", snap.Errors)
				}
				if snap.Latencies.Total != 3 {
					return errors.Errorf("expected latencies total 3, got %d", snap.Latencies.Total)
				}
				if snap.Latencies.Mean != float64(200*time.Millisecond) {
					return errors.Errorf("expected latency mean %v, got %v", float64(200*time.Millisecond), snap.Latencies.Mean)
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestCollector_Merge(t *testing.T) {
	type args struct {
		c1Opts    []Option
		c1Records []*RequestResult
		c2Opts    []Option
		c2Records []*RequestResult
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, args args) (Collector, error) {
		t.Helper()
		c1, err := NewCollector(args.c1Opts...)
		if err != nil {
			return nil, err
		}
		for _, r := range args.c1Records {
			c1.Record(context.Background(), 0, r)
		}

		c2, err := NewCollector(args.c2Opts...)
		if err != nil {
			return nil, err
		}
		for _, r := range args.c2Records {
			c2.Record(context.Background(), 0, r)
		}

		if err := c2.MergeInto(c1); err != nil {
			return nil, err
		}
		// c1 is modified
		h1, _ := c1.CounterHandle("c1")
		h1.Inc()
		return c1, nil
	}, []test.Case[Collector, args]{
		{
			Name: "merge two collectors",
			Args: args{
				c1Opts: []Option{WithCustomCounters("c1")},
				c1Records: []*RequestResult{
					{Latency: 100 * time.Millisecond},
				},
				c2Opts: []Option{WithCustomCounters("c1", "c2")},
				c2Records: []*RequestResult{
					{Latency: 200 * time.Millisecond, Err: errors.New("err")},
				},
			},
			CheckFunc: func(t *testing.T, want test.Result[Collector], got test.Result[Collector]) error {
				t.Helper()
				if got.Err != nil {
					return got.Err
				}
				c := got.Val
				snap := c.GlobalSnapshot()
				if snap.Total != 2 {
					return errors.Errorf("expected total 2, got %d", snap.Total)
				}
				if snap.Errors != 1 {
					return errors.Errorf("expected errors 1, got %d", snap.Errors)
				}
				h1, _ := c.CounterHandle("c1")
				if h1.value.Load() != 1 {
					return errors.Errorf("expected c1 counter to be 1, got %d", h1.value.Load())
				}
				h2, err := c.CounterHandle("c2")
				if err != nil {
					return errors.New("expected c2 counter to exist in c1 after merge")
				}
				if h2.value.Load() != 0 {
					return errors.Errorf("expected c2 counter to be 0, got %d", h2.value.Load())
				}
				return nil
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}
