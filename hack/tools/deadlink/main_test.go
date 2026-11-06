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

package main

import (
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func newStubClient(status int, calls *atomic.Int32) *http.Client {
	return &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			calls.Add(1)
			return &http.Response{
				StatusCode: status,
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    r,
			}, nil
		}),
	}
}

func TestGetFiles(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		setup   func(t *testing.T, dir string) string
		want    []string
		wantErr bool
	}{
		{
			name: "collect html files",
			setup: func(t *testing.T, dir string) string {
				t.Helper()
				for _, f := range []string{"a.html", "b.txt"} {
					if err := os.WriteFile(filepath.Join(dir, f), nil, 0o600); err != nil {
						t.Fatal(err)
					}
				}
				return dir
			},
			want: []string{"a.html"},
		},
		{
			// a walk that stops early must not be reported as a complete file list.
			name: "walk failure is returned",
			setup: func(t *testing.T, dir string) string {
				t.Helper()
				return filepath.Join(dir, "missing")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			root := tt.setup(t, dir)
			got, err := getFiles(root)
			if (err != nil) != tt.wantErr {
				t.Fatalf("getFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			want := make([]string, 0, len(tt.want))
			for _, f := range tt.want {
				want = append(want, filepath.Join(dir, f))
			}
			if !slices.Equal(got, want) {
				t.Errorf("getFiles() = %v, want %v", got, want)
			}
		})
	}
}

func TestCheckLink(t *testing.T) {
	t.Parallel()

	const n = 8
	tests := []struct {
		name        string
		urlPrefix   string
		status      int
		wantOK      bool
		wantErrCode int
	}{
		{
			name:      "alive links to BASE_URL do not deadlock",
			urlPrefix: BASE_URL + "/",
			status:    http.StatusOK,
			wantOK:    true,
		},
		{
			name:        "dead links to BASE_URL are recorded in errLinks",
			urlPrefix:   BASE_URL + "/",
			status:      http.StatusNotFound,
			wantErrCode: http.StatusNotFound,
		},
		{
			name:      "alive external links do not deadlock",
			urlPrefix: "https://example.com/",
			status:    http.StatusOK,
			wantOK:    true,
		},
		{
			name:        "dead external links are recorded in errLinks",
			urlPrefix:   "https://example.com/",
			status:      http.StatusNotFound,
			wantErrCode: http.StatusNotFound,
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			var (
				calls    atomic.Int32
				mu       sync.Mutex
				errLinks = map[string]int{}
				exLinks  = map[string]int{}
				cli      = newStubClient(test.status, &calls)
				results  = make(chan bool, n)
			)
			keys := make([]string, n)
			for i := range keys {
				keys[i] = string(rune('a' + i))
			}
			for _, k := range keys {
				go func() {
					results <- checkLink(&mu, errLinks, exLinks, k, test.urlPrefix+k, cli)
				}()
			}
			timeout := time.After(10 * time.Second)
			for range n {
				select {
				case got := <-results:
					if got != test.wantOK {
						tt.Errorf("checkLink expected %v, got %v", test.wantOK, got)
					}
				case <-timeout:
					tt.Fatal("checkLink did not return in time: the mutex is likely never released")
				}
			}

			mu.Lock()
			defer mu.Unlock()
			if test.wantOK {
				if len(errLinks) != 0 {
					tt.Errorf("unexpected errLinks: %v", errLinks)
				}
				return
			}
			for _, k := range keys {
				if errLinks[k] != test.wantErrCode {
					tt.Errorf("errLinks[%q] expected %d, got %d", k, test.wantErrCode, errLinks[k])
				}
			}
		})
	}
}

func TestCheckLinkSkipsRecordedErrLink(t *testing.T) {
	t.Parallel()

	var (
		calls    atomic.Int32
		mu       sync.Mutex
		errLinks = map[string]int{"k": http.StatusNotFound}
	)
	done := make(chan bool, 1)
	go func() {
		done <- checkLink(&mu, errLinks, map[string]int{}, "k", BASE_URL+"/k", newStubClient(http.StatusOK, &calls))
	}()
	select {
	case got := <-done:
		if got {
			t.Error("checkLink expected false for a link recorded in errLinks")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("checkLink did not return in time")
	}
	if c := calls.Load(); c != 0 {
		t.Errorf("expected no request, got %d", c)
	}
}
