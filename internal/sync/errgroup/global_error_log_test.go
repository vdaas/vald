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
	"fmt"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/log/nop"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
)

// recordingLogger discards everything except Errorf messages, which it records
// so a test can count them by a unique error message.
type recordingLogger struct {
	logger.Interface
	errors []string
	mu     sync.Mutex
}

//nolint:gochecknoglobals // TestMain installs this as the process-wide log.Init logger, so tests must share the one instance to read what the group logged
var recLog = &recordingLogger{Interface: nop.New()}

func (r *recordingLogger) Errorf(format string, vals ...any) {
	msg := fmt.Sprintf(format, vals...)
	r.mu.Lock()
	r.errors = append(r.errors, msg)
	r.mu.Unlock()
}

// count returns how many recorded Errorf messages contain s.
func (r *recordingLogger) count(s string) (n int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, m := range r.errors {
		if strings.Contains(m, s) {
			n++
		}
	}
	return n
}

// A task error on the process-wide group is logged when it occurs (its
// cancellation is watched by nobody and Wait only runs at shutdown); a local
// group keeps leaving reporting to its caller, and cancellation is never logged.
// Not parallel: it resets the package-level instance.
func TestGlobalGroupLogsTaskErrors(t *testing.T) {
	tests := []struct {
		group   func(t *testing.T) Group
		taskErr error
		name    string
		want    int
	}{
		{
			name: "global group logs a failed task",
			group: func(t *testing.T) Group {
				t.Helper()
				instance, once = nil, sync.Once{}
				t.Cleanup(func() { instance, once = nil, sync.Once{} })
				Init(t.Context())
				return Get()
			},
			taskErr: errors.New("TestGlobalGroupLogsTaskErrors/global failure"),
			want:    1,
		},
		{
			name: "local group does not log a failed task",
			group: func(t *testing.T) Group {
				t.Helper()
				g, _ := New(t.Context())
				return g
			},
			taskErr: errors.New("TestGlobalGroupLogsTaskErrors/local failure"),
			want:    0,
		},
		{
			name: "global group does not log a canceled task",
			group: func(t *testing.T) Group {
				t.Helper()
				instance, once = nil, sync.Once{}
				t.Cleanup(func() { instance, once = nil, sync.Once{} })
				Init(t.Context())
				return Get()
			},
			taskErr: errors.Errorf("TestGlobalGroupLogsTaskErrors/canceled: %w", context.Canceled),
			want:    0,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := tc.group(t)
			g.Go(func() error { return tc.taskErr })
			// Wait returns once the task has run; its result is not under test.
			if err := g.Wait(); err != nil {
				t.Logf("Wait() error = %v", err)
			}
			if got := recLog.count(tc.taskErr.Error()); got != tc.want {
				t.Errorf("logged %d times, want %d", got, tc.want)
			}
		})
	}
}
