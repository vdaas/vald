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
package targets

import (
	"context"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
)

// messages flattens a joined error into its leaf messages in join order.
func messages(err error) []string {
	if err == nil {
		return nil
	}
	if u, ok := err.(interface{ Unwrap() []error }); ok {
		var ms []string
		for _, e := range u.Unwrap() {
			ms = append(ms, messages(e)...)
		}
		return ms
	}
	return []string{err.Error()}
}

func TestOutcomes(t *testing.T) {
	t.Parallel()

	type record struct {
		err    error
		target string
	}
	notFound := func(addr string) string { return errors.ErrGRPCClientConnNotFound(addr).Error() }
	tests := []struct {
		clientErr error
		name      string
		addrs     []string
		records   []record
		want      []string
	}{
		{
			name:    "all targets succeed",
			addrs:   []string{"a", "b"},
			records: []record{{target: "a"}, {target: "b"}},
		},
		{
			name:    "target errors are joined in target order, not record order",
			addrs:   []string{"a", "b", "c"},
			records: []record{{target: "c", err: errors.New("c failed")}, {target: "b"}, {target: "a", err: errors.New("a failed")}},
			want:    []string{"a failed", "c failed"},
		},
		{
			name:    "a successful retry clears the earlier failure",
			addrs:   []string{"a"},
			records: []record{{target: "a", err: errors.New("a failed")}, {target: "a"}},
		},
		{
			name:    "a failing retry keeps only the latest failure",
			addrs:   []string{"a"},
			records: []record{{target: "a", err: errors.New("a first")}, {target: "a", err: errors.New("a second")}},
			want:    []string{"a second"},
		},
		{
			name:    "skipped targets are reported in target order after target errors",
			addrs:   []string{"c", "a", "b"},
			records: []record{{target: "b", err: errors.New("b failed")}},
			want:    []string{"b failed", notFound("a"), notFound("c")},
		},
		{
			name:      "client error comes first and suppresses the skipped targets",
			addrs:     []string{"a", "b"},
			records:   []record{{target: "b", err: errors.New("b failed")}},
			clientErr: errors.New("client failed"),
			want:      []string{"client failed", "b failed"},
		},
		{
			name:  "no targets",
			addrs: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			o := New(tc.addrs)
			for _, r := range tc.records {
				o.Record(r.target, r.err)
			}
			err := o.Err(tc.clientErr)
			if got := messages(err); !slices.Equal(got, tc.want) {
				t.Errorf("Err() = %q, want %q", got, tc.want)
			}
			if len(tc.want) == 0 && err != nil {
				t.Errorf("Err() = %v, want nil", err)
			}
		})
	}
}

// TestOutcomesConcurrentRecord records from one goroutine per target, as the
// OrderedRangeConcurrent callbacks do; run under -race it pins the locking.
func TestOutcomesConcurrentRecord(t *testing.T) {
	t.Parallel()

	const n = 64
	addrs := make([]string, n)
	for i := range addrs {
		addrs[i] = "t" + strconv.Itoa(1000+i)
	}
	o := New(addrs)
	var wg sync.WaitGroup
	for i, addr := range addrs {
		wg.Go(func() {
			var err error
			if i%2 == 0 {
				err = errors.New(addr + " failed")
			}
			o.Record(addr, err)
		})
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("Record goroutines did not finish")
	}

	var want []string
	for i, addr := range addrs {
		if i%2 == 0 {
			want = append(want, addr+" failed")
		}
	}
	if got := messages(o.Err(nil)); !slices.Equal(got, want) {
		t.Errorf("Err() = %q, want %q", got, want)
	}
}
