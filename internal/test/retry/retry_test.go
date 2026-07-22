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

package retry

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/strings"
)

func TestDo_ModeOnce(t *testing.T) {
	var calls int
	wantErr := errors.New("attempt failure")
	if err := Do(t.Context(), Policy{}, func(context.Context) error {
		calls++
		return wantErr
	}); !errors.Is(err, wantErr) {
		t.Errorf("ModeOnce must return fn's error verbatim, got: %v", err)
	}
	if calls != 1 {
		t.Errorf("ModeOnce must run fn exactly once, ran %d times", calls)
	}
}

func TestDo_ModeSuccess(t *testing.T) {
	var calls, retried int
	if err := Do(t.Context(), Policy{
		Mode:    ModeSuccess,
		OnRetry: func(uint64, error) { retried++ },
	}, func(context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("not ready yet")
		}
		return nil
	}); err != nil {
		t.Errorf("ModeSuccess must return nil once an attempt succeeds, got: %v", err)
	}
	if calls != 3 {
		t.Errorf("ModeSuccess must stop at the first success, ran %d times", calls)
	}
	if retried != 2 {
		t.Errorf("OnRetry must observe each continued failure, observed %d", retried)
	}
}

func TestDo_ModeSuccess_BoundedWait(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()
	var calls int
	if err := Do(ctx, Policy{Mode: ModeSuccess, Interval: 10 * time.Millisecond}, func(context.Context) error {
		calls++
		return errors.New("never succeeds")
	}); err != nil {
		t.Errorf("ModeSuccess must end a bounded wait without error on context expiry, got: %v", err)
	}
	if calls < 1 {
		t.Error("ModeSuccess must attempt at least once before the window expires")
	}
}

func TestDo_ModeSuccess_AttemptDeadline(t *testing.T) {
	// An attempt surfacing the context expiry itself must also end the
	// bounded wait without error instead of retrying forever.
	if err := Do(t.Context(), Policy{Mode: ModeSuccess}, func(context.Context) error {
		return context.DeadlineExceeded
	}); err != nil {
		t.Errorf("ModeSuccess must treat an expiry-surfacing attempt as the end of the wait, got: %v", err)
	}
}

func TestDo_ModeCount(t *testing.T) {
	var calls int
	err := Do(t.Context(), Policy{Mode: ModeCount, Count: 4}, func(context.Context) error {
		calls++
		if calls == 2 {
			return errors.New("second attempt failure")
		}
		return nil
	})
	if calls != 4 {
		t.Errorf("ModeCount must run exactly Count attempts regardless of outcomes, ran %d/4", calls)
	}
	if err == nil || !strings.Contains(err.Error(), "attempt 2 failed") {
		t.Errorf("ModeCount must join failures tagged with their attempt number, got: %v", err)
	}
}

func TestDo_ModeCount_ContextExpiry(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	var calls int
	err := Do(ctx, Policy{Mode: ModeCount, Count: 5}, func(context.Context) error {
		calls++
		if calls == 2 {
			cancel()
		}
		return nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("ModeCount must report an expiry that prevented the promised count, got: %v", err)
	}
	if calls != 2 {
		t.Errorf("no attempt may start after cancellation, ran %d", calls)
	}
}

func TestDo_ModeTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 40*time.Millisecond)
	defer cancel()
	failOnce := true
	var calls int
	err := Do(ctx, Policy{Mode: ModeTimeout}, func(context.Context) error {
		calls++
		if failOnce {
			failOnce = false
			return errors.New("transient failure")
		}
		return nil
	})
	if calls < 2 {
		t.Errorf("ModeTimeout must keep repeating after successes until expiry, ran %d times", calls)
	}
	if err == nil || !strings.Contains(err.Error(), "transient failure") {
		t.Errorf("ModeTimeout must join non-expiry failures observed on the way, got: %v", err)
	}

	// A clean window returns nil: expiry is the normal exit.
	ctx2, cancel2 := context.WithTimeout(t.Context(), 30*time.Millisecond)
	defer cancel2()
	if err := Do(ctx2, Policy{Mode: ModeTimeout}, func(context.Context) error {
		return nil
	}); err != nil {
		t.Errorf("ModeTimeout with no failures must return nil on expiry, got: %v", err)
	}
}

func TestDo_IntervalIsContextAware(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	var calls int
	start := time.Now()
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()
	if err := Do(ctx, Policy{Mode: ModeSuccess, Interval: 10 * time.Second}, func(context.Context) error {
		calls++
		return errors.New("force the interval wait")
	}); err != nil {
		t.Errorf("cancellation during the interval must follow the mode's expiry rule, got: %v", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Errorf("interval wait must be cut short by cancellation, waited %s", elapsed)
	}
	if calls != 1 {
		t.Errorf("no second attempt may start after cancellation, ran %d", calls)
	}
}
