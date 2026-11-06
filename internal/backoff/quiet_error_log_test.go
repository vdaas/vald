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

package backoff

import (
	"context"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/log/nop"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

// recordingLogger discards everything except Errord/Debugd messages, which it
// records so a test can count them by a unique backoff name.
type recordingLogger struct {
	logger.Interface
	errors []string
	debugs []string
	mu     sync.Mutex
}

//nolint:gochecknoglobals // TestMain installs this as the process-wide log.Init logger, so tests must share the one instance to read what Do logged
var recLog = &recordingLogger{Interface: nop.New()}

func (r *recordingLogger) Errord(msg string, _ ...any) {
	r.mu.Lock()
	r.errors = append(r.errors, msg)
	r.mu.Unlock()
}

func (r *recordingLogger) Debugd(msg string, _ ...any) {
	r.mu.Lock()
	r.debugs = append(r.debugs, msg)
	r.mu.Unlock()
}

// count returns how many Errord and Debugd messages mention name.
func (r *recordingLogger) count(name string) (errs, debugs int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range r.errors {
		if strings.Contains(m, "for: "+name+",") {
			errs++
		}
	}
	for _, m := range r.debugs {
		if strings.Contains(m, "for: "+name+",") {
			debugs++
		}
	}
	return errs, debugs
}

func TestWithQuietErrorLog(t *testing.T) {
	t.Parallel()
	tests := []struct {
		ctx  func() context.Context
		name string
		want bool
	}{
		{
			name: "unmarked context is not quiet",
			ctx:  context.Background,
			want: false,
		},
		{
			name: "marked context is quiet",
			ctx:  func() context.Context { return WithQuietErrorLog(context.Background()) },
			want: true,
		},
		{
			name: "derived contexts keep the flag",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(WithQuietErrorLog(context.Background()))
				t.Cleanup(cancel)
				return WithBackoffName(ctx, "derived")
			},
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := isQuietErrorLog(tc.ctx()); got != tc.want {
				t.Errorf("isQuietErrorLog() = %v, want %v", got, tc.want)
			}
		})
	}
}

// A quiet context demotes Do's per-attempt retryable error log from Error to
// Debug; an unmarked context keeps logging each attempt at Error.
func Test_backoff_Do_quietErrorLog(t *testing.T) {
	t.Parallel()
	const attempts = 3
	tests := []struct {
		name       string
		wantErrs   int
		wantDebugs int
		quiet      bool
	}{
		{
			name:     "unmarked context logs each retryable error at Error",
			quiet:    false,
			wantErrs: attempts,
		},
		{
			name:       "quiet context logs each retryable error at Debug",
			quiet:      true,
			wantDebugs: attempts,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			name := "Test_backoff_Do_quietErrorLog/" + tc.name
			ctx := WithBackoffName(t.Context(), name)
			if tc.quiet {
				ctx = WithQuietErrorLog(ctx)
			}
			b := New(
				WithInitialDuration("1ms"),
				WithMaximumDuration("2ms"),
				WithJitterLimit("1ms"),
				WithBackOffTimeLimit("1m"),
				WithRetryCount(attempts),
				WithEnableErrorLog(),
			)
			t.Cleanup(b.Close)
			retryErr := errors.New("transient failure")
			_, err := b.Do(ctx, func(context.Context) (any, bool, error) {
				return nil, true, retryErr
			})
			if !errors.Is(err, retryErr) {
				t.Fatalf("Do() error = %v, want %v", err, retryErr)
			}
			gotErrs, gotDebugs := recLog.count(name)
			if gotErrs != tc.wantErrs || gotDebugs != tc.wantDebugs {
				t.Errorf("Errord/Debugd count = %d/%d, want %d/%d", gotErrs, gotDebugs, tc.wantErrs, tc.wantDebugs)
			}
		})
	}
}
