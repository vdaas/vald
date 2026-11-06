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

package errgroup

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

// fanOutProbe records how many calls to hold overlap. A call stays in flight
// until more than limit calls overlap (then every call is released at once, so
// an unbounded fan-out is seen without waiting) or a short grace period ends
// (so a bounded fan-out still makes progress).
type fanOutProbe struct {
	over   chan struct{}
	mu     sync.Mutex
	limit  int
	cur    int
	peak   int
	calls  int
	closed bool
}

func (p *fanOutProbe) hold() {
	p.mu.Lock()
	p.cur++
	p.calls++
	p.peak = max(p.peak, p.cur)
	if p.cur > p.limit && !p.closed {
		p.closed = true
		close(p.over)
	}
	p.mu.Unlock()
	timer := time.NewTimer(20 * time.Millisecond)
	select {
	case <-p.over:
	case <-timer.C:
	}
	timer.Stop()
	p.mu.Lock()
	p.cur--
	p.mu.Unlock()
}

func (p *fanOutProbe) stats() (peak, calls int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.peak, p.calls
}

// TestNewBoundedFanOut pins the group the Multi* handlers fan out on: it never
// runs more than limit items at once, limit <= 0 means unbounded (not a
// deadlock), and neither an already-cancelled request nor a failing item makes
// it silently skip the items that had to wait for a slot; the item error is
// still returned by Wait.
func TestNewBoundedFanOut(t *testing.T) {
	t.Parallel()
	const n = 6
	errItem := errors.New("item failed")
	tests := []struct {
		name      string
		limit     int
		cancelled bool
		fails     bool
	}{
		{name: "bounded", limit: 2},
		{name: "bounded with a cancelled request ctx", limit: 2, cancelled: true},
		{name: "bounded with a failing item", limit: 2, fails: true},
		{name: "zero limit is unbounded", limit: 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			if test.cancelled {
				cancel()
			}
			p := &fanOutProbe{over: make(chan struct{}), limit: 2}
			done := make(chan error, 1)
			go func() {
				eg := NewBoundedFanOut(ctx, test.limit)
				for i := range n {
					eg.Go(func() error {
						p.hold()
						if test.fails && i == 0 {
							return errItem
						}
						return nil
					})
				}
				done <- eg.Wait()
			}()
			select {
			case err := <-done:
				if test.fails {
					if !errors.Is(err, errItem) {
						t.Errorf("Wait() = %v, want %v", err, errItem)
					}
				} else if err != nil {
					t.Fatalf("Wait() = %v, want nil", err)
				}
			case <-time.After(10 * time.Second):
				t.Fatal("timed out waiting for the group")
			}
			peak, calls := p.stats()
			if calls != n {
				t.Errorf("ran %d items, want all %d", calls, n)
			}
			if test.limit > 0 && peak > test.limit {
				t.Errorf("peak concurrency = %d, want <= %d", peak, test.limit)
			}
		})
	}
}
