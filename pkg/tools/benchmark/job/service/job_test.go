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
	"os"
	"os/signal"
	"slices"
	"syscall"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/client/v1/client/vald"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/grpc"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/errgroup"
	"github.com/vdaas/vald/internal/test/data/hdf5" //nolint:depguard // fake hdf5.Data for job tests
)

// testTimeout bounds every blocking wait in this file.
const testTimeout = 10 * time.Second

type fakeLimiter struct{}

func (fakeLimiter) Wait(ctx context.Context) error { return ctx.Err() }

type fakeHdf5 struct {
	hdf5.Data
	vecs [][]float32
}

func (f fakeHdf5) GetByGroupName(string) [][]float32 { return f.vecs }

type fakeClient struct {
	vald.Client
	startCh     <-chan error
	getReqs     []*payload.Object_VectorRequest
	searchCalls atomic.Int32
	linearCalls atomic.Int32
	mu          sync.Mutex
}

func (f *fakeClient) Start(context.Context) (<-chan error, error) { return f.startCh, nil }

func (f *fakeClient) Search(
	ctx context.Context, _ *payload.Search_Request, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	f.searchCalls.Add(1)
	return searchResponse("a"), nil
}

func (f *fakeClient) LinearSearch(
	ctx context.Context, _ *payload.Search_Request, _ ...grpc.CallOption,
) (*payload.Search_Response, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	f.linearCalls.Add(1)
	return searchResponse("a"), nil
}

func (f *fakeClient) GetObject(
	_ context.Context, in *payload.Object_VectorRequest, _ ...grpc.CallOption,
) (*payload.Object_Vector, error) {
	f.mu.Lock()
	f.getReqs = append(f.getReqs, in)
	f.mu.Unlock()
	return &payload.Object_Vector{}, nil
}

func newTestJob(cli *fakeClient, ds *config.BenchmarkDataset, vecs [][]float32) *job {
	eg, _ := errgroup.New(context.Background())
	return &job{
		eg:               eg,
		limiter:          fakeLimiter{},
		client:           cli,
		hdf5:             fakeHdf5{vecs: vecs},
		dataset:          ds,
		concurrencyLimit: 4,
	}
}

func Test_New_rangeCoversExactlyIndexes(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name    string
		start   int
		end     int
		indexes int
		wantEnd int
	}{
		{name: "range smaller than indexes is extended to exactly indexes entries", start: 1, end: 2, indexes: 5, wantEnd: 5},
		{name: "range starting at zero", start: 0, end: 0, indexes: 3, wantEnd: 2},
		{name: "range equal to indexes is untouched", start: 1, end: 5, indexes: 5, wantEnd: 5},
		{name: "range larger than indexes is untouched", start: 1, end: 9, indexes: 5, wantEnd: 9},
	} {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			ds := &config.BenchmarkDataset{
				Name:    "test-dataset",
				Indexes: tc.indexes,
				Range:   &config.BenchmarkDatasetRange{Start: tc.start, End: tc.end},
			}
			if _, err := New(WithJobType(EXISTS), WithDataset(ds)); err != nil {
				tt.Fatalf("New returned an error: %v", err)
			}
			if ds.Range.End != tc.wantEnd {
				tt.Errorf("Range.End = %d, want %d", ds.Range.End, tc.wantEnd)
			}
		})
	}
}

func Test_job_search_linearSearch(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name  string
		start int
		end   int
		vecs  int
	}{
		{name: "one-based range, one vector per iteration", start: 1, end: 3, vecs: 3},
		{name: "one-based range wrapping over fewer vectors", start: 1, end: 4, vecs: 2},
		{name: "range larger than indexes", start: 1, end: 5, vecs: 3},
		{name: "range not starting at one", start: 3, end: 6, vecs: 4},
	} {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			vecs := make([][]float32, tc.vecs)
			for i := range vecs {
				vecs[i] = []float32{float32(i)}
			}
			cli := new(fakeClient)
			j := newTestJob(cli, &config.BenchmarkDataset{
				Group:   "test",
				Indexes: 2,
				Range:   &config.BenchmarkDatasetRange{Start: tc.start, End: tc.end},
			}, vecs)
			j.searchConfig = &config.SearchConfig{Num: 1, EnableLinearSearch: true}

			if err := j.search(tt.Context(), make(chan error, 8)); err != nil {
				tt.Fatalf("search returned an error: %v", err)
			}
			want := tc.end - tc.start + 1
			if got := int(cli.searchCalls.Load()); got != want {
				tt.Errorf("Search calls = %d, want %d", got, want)
			}
			if got := int(cli.linearCalls.Load()); got != want {
				tt.Errorf("LinearSearch calls = %d, want %d", got, want)
			}
		})
	}
}

func Test_job_getObject_filterTargets(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		objectConfig *config.ObjectConfig
		name         string
		want         []*payload.Filter_Target
	}{
		{name: "no object config", objectConfig: nil, want: []*payload.Filter_Target{}},
		{
			name:         "empty targets",
			objectConfig: &config.ObjectConfig{},
			want:         []*payload.Filter_Target{},
		},
		{
			name: "targets are built once, nil entries skipped",
			objectConfig: &config.ObjectConfig{FilterConfig: config.FilterConfig{Targets: []*config.FilterTarget{
				{Host: "a", Port: 1},
				nil,
				{Host: "b", Port: 2},
			}}},
			want: []*payload.Filter_Target{{Host: "a", Port: 1}, {Host: "b", Port: 2}},
		},
	} {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			cli := new(fakeClient)
			j := newTestJob(cli, &config.BenchmarkDataset{
				Range: &config.BenchmarkDatasetRange{Start: 1, End: 3},
			}, nil)
			j.objectConfig = tc.objectConfig

			if err := j.getObject(tt.Context(), make(chan error, 8)); err != nil {
				tt.Fatalf("getObject returned an error: %v", err)
			}
			cli.mu.Lock()
			defer cli.mu.Unlock()
			if len(cli.getReqs) != 3 {
				tt.Fatalf("GetObject calls = %d, want 3", len(cli.getReqs))
			}
			for _, req := range cli.getReqs {
				got := req.GetFilters().GetTargets()
				if len(got) != len(tc.want) {
					tt.Fatalf("filter targets = %v, want %v", got, tc.want)
				}
				for i := range got {
					if got[i].GetHost() != tc.want[i].GetHost() || got[i].GetPort() != tc.want[i].GetPort() {
						tt.Errorf("filter target[%d] = %v, want %v", i, got[i], tc.want[i])
					}
				}
			}
		})
	}
}

func Test_addNoiseToVec(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		in   []float32
	}{
		{name: "nil", in: nil},
		{name: "empty", in: []float32{}},
		{name: "single element", in: []float32{1}},
		{name: "multiple elements", in: []float32{1, 2, 3, 4}},
	} {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			orig := slices.Clone(tc.in)
			got := addNoiseToVec(tc.in)
			if !slices.Equal(tc.in, orig) {
				tt.Errorf("input slice was mutated: got %v, want %v", tc.in, orig)
			}
			if len(got) != len(tc.in) {
				tt.Fatalf("len(result) = %d, want %d", len(got), len(tc.in))
			}
			if len(got) > 0 && &got[0] == &tc.in[0] {
				tt.Error("result aliases the input slice")
			}
		})
	}
}

func Test_job_Start(t *testing.T) {
	// The job goroutine sends SIGTERM to its own process when it finishes. The
	// signal is process-directed and delivered asynchronously, so every subtest
	// receives exactly one before returning and signal.Stop only runs after all
	// of them were consumed; otherwise a pending SIGTERM would kill the test
	// binary. internal/os has no signal aliases, hence stdlib os/signal.
	sigs := make(chan os.Signal, 16)
	signal.Notify(sigs, syscall.SIGTERM)
	defer signal.Stop(sigs)

	for _, tc := range []struct {
		monitor func() chan error
		jobFunc func(ctx context.Context, ech chan error) error
		name    string
		// wantNoValues requires that nothing is forwarded on ech.
		wantNoValues bool
	}{
		{
			// `case ech <- <-cech` blocked on the receive before select was entered,
			// so a silent monitor channel kept the goroutine (and ech) alive forever.
			name:    "silent monitor channel does not outlive ctx",
			monitor: func() chan error { return make(chan error) },
			jobFunc: func(ctx context.Context, _ chan error) error {
				<-ctx.Done()
				return nil
			},
			wantNoValues: true,
		},
		{
			// A closed monitor channel used to yield nil forever via the same expression.
			name: "closed monitor channel does not spam nil",
			monitor: func() chan error {
				ch := make(chan error)
				close(ch)
				return ch
			},
			jobFunc: func(ctx context.Context, _ chan error) error {
				<-ctx.Done()
				return nil
			},
			wantNoValues: true,
		},
		{
			// The job's deferred error send used to run against an ech the monitor
			// goroutine had already closed after ctx was cancelled (send on closed
			// channel). jobFunc waits (bounded) for ech to be closed: on the old code it
			// is closed early and the job then panics; on the new code ech stays open
			// until the job finishes, so the wait times out and the job returns its error.
			name: "job error after ctx cancel does not send on closed ech",
			monitor: func() chan error {
				// closed so the old monitor returned on ctx.Done and closed ech
				// (a silent channel kept it blocked in `<-cech`, never closing ech).
				ch := make(chan error)
				close(ch)
				return ch
			},
			jobFunc: func(ctx context.Context, ech chan error) error {
				<-ctx.Done()
				timer := time.NewTimer(200 * time.Millisecond)
				defer timer.Stop()
				for {
					select {
					case _, ok := <-ech:
						if !ok {
							return errors.New("job failed after cancel")
						}
					case <-timer.C:
						return errors.New("job failed after cancel")
					}
				}
			},
		},
	} {
		t.Run(tc.name, func(tt *testing.T) {
			cli := &fakeClient{startCh: tc.monitor()}
			j := newTestJob(cli, &config.BenchmarkDataset{}, nil)
			j.jobFunc = tc.jobFunc

			ctx, cancel := context.WithCancel(tt.Context())
			defer cancel()
			ech, err := j.Start(ctx)
			if err != nil {
				tt.Fatalf("Start returned an error: %v", err)
			}
			cancel()

			timeout := time.After(testTimeout)
			var got []error
			for open := true; open; {
				select {
				case v, ok := <-ech:
					if !ok {
						open = false
						continue
					}
					got = append(got, v)
				case <-timeout:
					tt.Fatal("ech was not closed after ctx was cancelled")
				}
			}
			if tc.wantNoValues && len(got) != 0 {
				tt.Errorf("unexpected values forwarded on ech: %v", got)
			}
			// The SIGTERM is sent before ech is closed; wait for its delivery.
			select {
			case <-sigs:
			case <-time.After(testTimeout):
				tt.Fatal("SIGTERM from the finished job was not received")
			}
		})
	}
}
