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

package faiss

import (
	"math"
	"slices"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test/data/vector"
)

// dim is the vector dimension of the faiss built by newTestFaiss.
const dim = 4

// newTestFaiss returns a faiss whose C index is nil (the state after Close),
// so that no call in these tests can reach the C library.
func newTestFaiss(t *testing.T) *faiss {
	t.Helper()
	f := &faiss{mu: new(sync.RWMutex)}
	if err := WithDimension(dim)(f); err != nil {
		t.Fatal(err)
	}
	return f
}

func TestFaiss_InvalidInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		call func(f *faiss) error
		name string
	}{
		{
			name: "Train returns error when nb is zero",
			call: func(f *faiss) error { return f.Train(0, nil) },
		},
		{
			name: "Train returns error when nb is negative",
			call: func(f *faiss) error { return f.Train(-1, make([]float32, dim)) },
		},
		{
			name: "Train returns error when xb is empty",
			call: func(f *faiss) error { return f.Train(1, []float32{}) },
		},
		{
			name: "Train returns error when xb is shorter than nb*dim",
			call: func(f *faiss) error { return f.Train(2, make([]float32, dim)) },
		},
		{
			name: "Add returns error when nb is zero",
			call: func(f *faiss) error { _, err := f.Add(0, nil, nil); return err },
		},
		{
			name: "Add returns error when nb is zero with non-empty xb",
			call: func(f *faiss) error { _, err := f.Add(0, make([]float32, dim), []int64{1}); return err },
		},
		{
			name: "Add returns error when xids is empty",
			call: func(f *faiss) error { _, err := f.Add(1, make([]float32, dim), nil); return err },
		},
		{
			name: "Add returns error when xb length mismatches dimension",
			call: func(f *faiss) error { _, err := f.Add(1, make([]float32, dim+1), []int64{1}); return err },
		},
		{
			name: "Search returns error when nq is zero",
			call: func(f *faiss) error { _, err := f.Search(1, 1, 0, nil); return err },
		},
		{
			name: "Search returns error when k is zero",
			call: func(f *faiss) error { _, err := f.Search(0, 1, 1, make([]float32, dim)); return err },
		},
		{
			name: "Search returns error when xq length mismatches dimension",
			call: func(f *faiss) error { _, err := f.Search(1, 1, 1, make([]float32, dim-1)); return err },
		},
		{
			name: "Remove returns error when size is zero",
			call: func(f *faiss) error { _, err := f.Remove(0, nil); return err },
		},
		{
			name: "Remove returns error when ids is empty",
			call: func(f *faiss) error { _, err := f.Remove(1, []int64{}); return err },
		},
		{
			name: "Remove returns error when ids is shorter than size",
			call: func(f *faiss) error { _, err := f.Remove(2, []int64{1}); return err },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			f := newTestFaiss(t)
			if err := test.call(f); err == nil {
				t.Error("expected an error, got nil")
			} else if errors.Is(err, errIndexClosed) {
				t.Errorf("validation must run before the closed check: %v", err)
			}
		})
	}
}

// TestFaiss_NaNOrInfVector verifies that vectors containing NaN or +-Inf are rejected
// with errors.ErrNaNOrInfVector before the closed-index check, i.e. without reaching C.
func TestFaiss_NaNOrInfVector(t *testing.T) {
	t.Parallel()

	bad := map[string]float32{
		"NaN":  float32(math.NaN()),
		"+Inf": float32(math.Inf(1)),
		"-Inf": float32(math.Inf(-1)),
	}
	calls := map[string]func(f *faiss, v []float32) error{
		"Train":  func(f *faiss, v []float32) error { return f.Train(1, v) },
		"Add":    func(f *faiss, v []float32) error { _, err := f.Add(1, v, []int64{1}); return err },
		"Search": func(f *faiss, v []float32) error { _, err := f.Search(1, 1, 1, v); return err },
	}

	for cname, call := range calls {
		for bname, b := range bad {
			t.Run(cname+"/"+bname, func(t *testing.T) {
				t.Parallel()
				vec := make([]float32, dim)
				vec[dim-1] = b
				if err := call(newTestFaiss(t), vec); !errors.Is(err, errors.ErrNaNOrInfVector) {
					t.Errorf("error = %v, want %v", err, errors.ErrNaNOrInfVector)
				}
			})
		}
	}
}

// TestFaiss_Closed verifies that valid requests against a closed index return an error
// instead of passing a NULL index pointer to the C library.
func TestFaiss_Closed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		call func(f *faiss) error
		name string
	}{
		{
			name: "SaveIndex",
			call: func(f *faiss) error { return f.SaveIndex() },
		},
		{
			name: "SaveIndexWithPath",
			call: func(f *faiss) error { return f.SaveIndexWithPath("dummy") },
		},
		{
			name: "Train",
			call: func(f *faiss) error { return f.Train(1, make([]float32, dim)) },
		},
		{
			name: "Add",
			call: func(f *faiss) error { _, err := f.Add(1, make([]float32, dim), []int64{1}); return err },
		},
		{
			name: "Search",
			call: func(f *faiss) error { _, err := f.Search(1, 1, 1, make([]float32, dim)); return err },
		},
		{
			name: "Remove",
			call: func(f *faiss) error { _, err := f.Remove(1, []int64{1}); return err },
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			f := newTestFaiss(t)
			if err := test.call(f); !errors.Is(err, errIndexClosed) {
				t.Errorf("expected errIndexClosed, got %v", err)
			}
		})
	}
}

func TestFaiss_CloseIdempotent(t *testing.T) {
	t.Parallel()
	f := newTestFaiss(t)
	f.Close()
	f.Close()
	if f.st != nil {
		t.Error("st must be nil after Close")
	}
}

// TestNew_optionErrorKeepsCause verifies that an option error is returned with its
// cause instead of an opaque "faiss option error". It fails before any C call.
func TestNew_optionErrorKeepsCause(t *testing.T) {
	t.Parallel()
	// m must divide dimension.
	_, err := New(WithDimension(4), WithM(3))
	var target *errors.ErrInvalidOption
	if !errors.As(err, &target) {
		t.Errorf("error = %v, want it to wrap *errors.ErrInvalidOption", err)
	}
}

// TestToSearchResult verifies that the -1 labels FAISS pads missing neighbours
// with are dropped instead of being reported as ID uint32(-1).
func TestToSearchResult(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		I    []int64
		D    []float32
		want []algorithm.SearchResult
		k    int
	}{
		{
			name: "all k neighbours found",
			k:    2,
			I:    []int64{3, 7},
			D:    []float32{0.1, 0.2},
			want: []algorithm.SearchResult{{ID: 3, Distance: 0.1}, {ID: 7, Distance: 0.2}},
		},
		{
			name: "padded tail is dropped",
			k:    3,
			I:    []int64{5, -1, -1},
			D:    []float32{0.5, 3.4e38, 3.4e38},
			want: []algorithm.SearchResult{{ID: 5, Distance: 0.5}},
		},
		{
			name: "no neighbour found",
			k:    2,
			I:    []int64{-1, -1},
			D:    []float32{3.4e38, 3.4e38},
			want: []algorithm.SearchResult{},
		},
		{
			name: "only the first query's k results are used",
			k:    1,
			I:    []int64{4, 9},
			D:    []float32{0.4, 0.9},
			want: []algorithm.SearchResult{{ID: 4, Distance: 0.4}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			if got := toSearchResult(test.k, test.I, test.D); !slices.Equal(got, test.want) {
				t.Errorf("got %+v, want %+v", got, test.want)
			}
		})
	}
}

// TestFaiss_SearchConcurrentNprobe runs searches with different nprobe values at
// the same time on one real index. faiss_search used to write nprobe into the
// shared index before searching, so a search could run with the nprobe of another
// concurrent call. A full-probe search must always return its sequential result.
// The final Close frees the real index and its quantizer through faiss_free.
func TestFaiss_SearchConcurrentNprobe(t *testing.T) {
	t.Parallel()
	const (
		dim     = 8
		nlist   = 4
		nb      = 1000
		k       = 10
		queries = 8
		rounds  = 50
	)
	tests := []struct {
		name       string
		methodType string
	}{
		{name: "ivfpq", methodType: "ivfpq"},
		{name: "binaryindex", methodType: "binaryindex"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			fi, err := New(
				WithDimension(dim),
				WithNlist(nlist),
				WithM(2),
				WithNbitsPerIdx(8),
				WithMethodType(test.methodType),
				WithMetricType("l2"),
			)
			if err != nil {
				t.Fatal(err)
			}
			defer fi.Close()

			vecs := vector.GaussianDistributedFloat32VectorGenerator(nb, dim)
			xb := make([]float32, 0, nb*dim)
			ids := make([]int64, nb)
			for i, v := range vecs {
				xb = append(xb, v...)
				ids[i] = int64(i)
			}
			if err = fi.Train(nb, xb); err != nil {
				t.Fatal(err)
			}
			if _, err = fi.Add(nb, xb, ids); err != nil {
				t.Fatal(err)
			}

			want := make([][]algorithm.SearchResult, queries)
			for i := range queries {
				if want[i], err = fi.Search(k, nlist, 1, vecs[i]); err != nil {
					t.Fatal(err)
				}
			}

			errs := make(chan error, 2*queries)
			var wg sync.WaitGroup
			for i := range queries {
				q := vecs[i]
				wg.Go(func() {
					for range rounds {
						got, err := fi.Search(k, nlist, 1, q)
						if err != nil {
							errs <- err
							return
						}
						if !slices.Equal(got, want[i]) {
							errs <- errors.Errorf("query %d: full-probe search got %+v, want %+v", i, got, want[i])
							return
						}
					}
				})
				wg.Go(func() {
					for range rounds {
						if _, err := fi.Search(k, 1, 1, q); err != nil {
							errs <- err
							return
						}
					}
				})
			}
			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(time.Minute):
				t.Fatal("concurrent searches did not finish within 1m")
			}
			close(errs)
			for err := range errs {
				t.Error(err)
			}

			fi.Close()
		})
	}
}
