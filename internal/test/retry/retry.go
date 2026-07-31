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

// Package retry implements the declarative repeat semantics scenario-driven
// tests use: run an attempt function under a policy whose exit condition is
// "first success", "fixed attempt count" or "until the context expires".
//
// It deliberately does NOT build on internal/backoff, whose contract is
// imperative transient-operation retrying: backoff treats reaching the
// deadline as a failure (ErrBackoffTimeout), stops on the first success
// unconditionally, and paces attempts with exponential intervals and
// jitter. The scenario semantics need the opposite on all three axes — a
// deadline can be the *successful* end of a bounded wait, ModeCount keeps
// going after successes, and pacing is a fixed configured interval.
package retry

import (
	"context"
	"time"

	"github.com/vdaas/vald/internal/errors"
)

// Mode selects the exit condition of Do's attempt loop.
type Mode uint8

const (
	// ModeOnce runs fn exactly once and returns its error verbatim; it is
	// the zero value so an empty Policy means "no repetition".
	ModeOnce Mode = iota
	// ModeSuccess retries fn until an attempt succeeds. Context expiry ends
	// the loop WITHOUT error: this mode doubles as a bounded wait, where
	// never succeeding within the window is an accepted outcome the caller
	// verifies by other means (e.g. max_vector_dim.yaml's ResourceExhausted
	// branch keeps returning NotFound and passes once the window elapses).
	ModeSuccess
	// ModeCount runs fn exactly Count times regardless of per-attempt
	// outcomes and returns every failure joined; context expiry before the
	// count is reached is reported as part of the joined error.
	ModeCount
	// ModeTimeout repeats fn until the context expires — expiry is the
	// normal exit, not an error — and returns the non-expiry failures
	// observed along the way joined together.
	ModeTimeout
)

// Policy declares how Do repeats an attempt function.
type Policy struct {
	// OnRetry, when non-nil, observes every failed attempt that the policy
	// continues past (attempt is 1-based). It exists so callers can log
	// per-attempt diagnostics without this package depending on a logger.
	OnRetry func(attempt uint64, err error)
	// Interval is the pause between consecutive attempts (never before the
	// first). The pause is context-aware: expiry during it is handled by
	// the Mode's expiry rule.
	Interval time.Duration
	// Count is the number of attempts for ModeCount; other modes ignore it.
	Count uint64
	// Mode selects the exit condition; the zero value is ModeOnce.
	Mode Mode
}

// Do runs fn under p against ctx. fn receives ctx so attempts observe the
// same cancellation the loop does. See the Mode constants for the exact
// exit and error semantics of each policy.
func Do(ctx context.Context, p Policy, fn func(context.Context) error) error {
	if p.Mode == ModeOnce {
		return fn(ctx)
	}
	var errs error
	for attempt := uint64(1); ; attempt++ {
		if p.Mode == ModeCount && attempt > p.Count {
			return errs
		}
		if pause(ctx, p.Interval, attempt) {
			return exitOnExpiry(p.Mode, errs, ctx.Err())
		}
		err := fn(ctx)
		switch {
		case err == nil:
			if p.Mode == ModeSuccess {
				return nil
			}
		case errors.IsAny(err, context.Canceled, context.DeadlineExceeded):
			return exitOnExpiry(p.Mode, errs, err)
		default:
			errs = p.observe(attempt, err, errs)
		}
	}
}

// pause sleeps Interval between consecutive attempts (never before the
// first) and reports whether ctx expired while waiting or before the next
// attempt could start.
func pause(ctx context.Context, interval time.Duration, attempt uint64) (expired bool) {
	if attempt > 1 && interval > 0 {
		t := time.NewTimer(interval)
		defer t.Stop()
		select {
		case <-ctx.Done():
			return true
		case <-t.C:
		}
	}
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

// observe reports a failed attempt the policy continues past to OnRetry
// and, for the failure-accumulating modes, joins it into errs tagged with
// its attempt number.
func (p Policy) observe(attempt uint64, err, errs error) error {
	if p.OnRetry != nil {
		p.OnRetry(attempt, err)
	}
	if p.Mode == ModeSuccess {
		return errs
	}
	return errors.Join(errs, errors.Wrapf(err, "attempt %d failed", attempt))
}

// exitOnExpiry applies each mode's context-expiry rule: ModeSuccess ends a
// bounded wait without error, ModeTimeout treats expiry as its normal exit
// (returning only the failures accumulated on the way), and ModeCount
// reports the expiry because it prevented the promised attempt count.
func exitOnExpiry(m Mode, errs, cause error) error {
	if m == ModeSuccess {
		return nil
	}
	if m == ModeTimeout {
		return errs
	}
	return errors.Join(errs, cause)
}
