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

package service

import (
	"context"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/data/hdf5" //nolint:depguard // fake hdf5.Data for job tests
)

const testDatasetName = "fashion-mnist"

// fakeLoadHdf5 records whether PreStart downloads the dataset and the key it
// asks to load.
type fakeLoadHdf5 struct {
	hdf5.Data
	readKey    *hdf5.Key
	downloaded *bool
}

func (f fakeLoadHdf5) Download(string) error {
	*f.downloaded = true
	return nil
}

func (fakeLoadHdf5) GetName() hdf5.DatasetName { return hdf5.Original }
func (f fakeLoadHdf5) Read(key hdf5.Key) error {
	*f.readKey = key
	return nil
}

func Test_job_PreStart_datasetGroup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		group    string
		wantKey  hdf5.Key
		wantRead bool
		wantErr  bool
	}{
		{name: "train group loads train", group: "train", wantKey: hdf5.Train, wantRead: true},
		{name: "test group loads test", group: "test", wantKey: hdf5.Test, wantRead: true},
		{name: "unknown group is rejected before loading", group: "unknown", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			const notRead = hdf5.Key(-1)
			key := notRead
			downloaded := false
			j := &job{
				jobType: SEARCH,
				hdf5:    fakeLoadHdf5{readKey: &key, downloaded: &downloaded},
				dataset: &config.BenchmarkDataset{Group: tt.group},
			}
			err := j.PreStart(t.Context())
			if tt.wantErr {
				var e *errors.ErrInvalidOption
				if !errors.As(err, &e) {
					t.Errorf("PreStart() error = %v, want ErrInvalidOption", err)
				}
			} else if err != nil {
				t.Errorf("PreStart() unexpected error: %v", err)
			}
			// An unknown group must be rejected before the (large) dataset
			// download, not only before Read.
			if downloaded != tt.wantRead {
				t.Errorf("Download called = %v, want %v", downloaded, tt.wantRead)
			}
			if !tt.wantRead {
				if key != notRead {
					t.Errorf("Read called with %v, want no call", key)
				}
				return
			}
			if key != tt.wantKey {
				t.Errorf("Read key = %v, want %v", key, tt.wantKey)
			}
		})
	}
}

func TestWithDataset_nilRange(t *testing.T) {
	t.Parallel()
	tests := []struct {
		ds      *config.BenchmarkDataset
		name    string
		wantErr bool
	}{
		{name: "nil range is rejected", ds: &config.BenchmarkDataset{Name: testDatasetName}, wantErr: true},
		{
			name: "non-nil range is accepted",
			ds:   &config.BenchmarkDataset{Name: testDatasetName, Range: &config.BenchmarkDatasetRange{Start: 1, End: 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			j := new(job)
			err := WithDataset(tt.ds)(j)
			if tt.wantErr {
				var e *errors.ErrInvalidOption
				if !errors.As(err, &e) {
					t.Errorf("error = %v, want ErrInvalidOption", err)
				}
				if j.dataset != nil {
					t.Errorf("dataset = %v, want unset", j.dataset)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestNew_missingJobConfigError(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		wantMsg string
		opts    []Option
	}{
		{
			name:    "remove without remove config names the remove config",
			opts:    []Option{WithJobType(REMOVE)},
			wantMsg: "name: remove config, val: <nil>",
		},
		{
			name:    "upsert without upsert config prints the upsert config",
			opts:    []Option{WithJobType(UPSERT), WithInsertConfig(&config.InsertConfig{Timestamp: "1"})},
			wantMsg: "name: upsert config, val: <nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ds := &config.BenchmarkDataset{Name: testDatasetName, Range: &config.BenchmarkDatasetRange{Start: 1, End: 1}}
			_, err := New(append(tt.opts, WithDataset(ds))...)
			if err == nil || !strings.Contains(err.Error(), tt.wantMsg) {
				t.Errorf("New() error = %v, want message containing %q", err, tt.wantMsg)
			}
		})
	}
}

// PreStart must not wait for unrelated tasks of the shared errgroup while
// polling the before job.
func Test_job_PreStart_beforeJobDoesNotWaitForSharedGroup(t *testing.T) {
	t.Parallel()
	eg, _ := errgroup.New(context.Background())
	release := make(chan struct{})
	eg.Go(func() error {
		<-release
		return nil
	})
	t.Cleanup(func() {
		close(release)
		if err := eg.Wait(); err != nil {
			t.Errorf("unrelated task error: %v", err)
		}
	})

	j := &job{
		jobType:       EXISTS,
		eg:            eg,
		beforeJobName: "before",
		beforeJobDur:  time.Hour,
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel() // the poll loop returns nil on ctx.Done before the first tick

	done := make(chan error, 1)
	go func() { done <- j.PreStart(ctx) }()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("PreStart() error = %v, want nil", err)
		}
	case <-time.After(testTimeout):
		t.Fatal("PreStart blocked on an unrelated errgroup task")
	}
}
