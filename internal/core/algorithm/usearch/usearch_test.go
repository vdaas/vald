//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package usearch provides Go API implementation for USearch library. https://github.com/unum-cloud/usearch
package usearch

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

var (
	usearchComparator = []comparator.Option{
		comparator.AllowUnexported(usearch{}),
		comparator.RWMutexComparer,
		comparator.ErrorComparer,
		comparator.AtomicUint64Comparator,
	}

	searchResultComparator = []comparator.Option{
		comparator.CompareField("Distance", comparator.Comparer(func(s1, s2 float32) bool {
			if s1 == 0 { // if vec1 is same as vec2, the distance should be same
				return s2 == 0
			}
			// by setting non-zero value in test case, it will only check if both got/want is non-zero
			return s1 != 0 && s2 != 0
		})),
	}

	defaultAfterFunc = func(t *testing.T, u Usearch) error {
		t.Helper()

		if u == nil {
			return nil
		}

		u.Close()
		return nil
	}
)

func idxTempDir(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "index")
}

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	os.Exit(m.Run())
}

func Test_usearch_Search(t *testing.T) {
	type args struct {
		q []float32
		k int
	}
	type fields struct {
		idxPath          string
		quantizationType string
		metricType       string
		dimension        int
		connectivity     int
		expansionAdd     int
		expansionSearch  int
		multi            bool
	}
	type want struct {
		want []algorithm.SearchResult
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(t *testing.T, fields fields) (Usearch, error)
		want       want
		checkFunc  func(want, []algorithm.SearchResult, Usearch, error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, Usearch) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (Usearch, error) {
		t.Helper()

		return New(
			WithIndexPath(fields.idxPath),
			WithQuantizationType(fields.quantizationType),
			WithMetricType(fields.metricType),
			WithDimension(fields.dimension),
			WithConnectivity(fields.connectivity),
			WithExpansionAdd(fields.expansionAdd),
			WithExpansionSearch(fields.expansionSearch),
			WithMulti(fields.multi),
		)
	}
	defaultCheckFunc := func(w want, got []algorithm.SearchResult, n Usearch, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(got, w.want, searchResultComparator...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}

		return nil
	}
	insertCreateFunc := func(t *testing.T, fields fields, vecs [][]float32, poolSize uint32) (Usearch, error) { // create func with insert/index
		t.Helper()

		u, err := defaultCreateFunc(t, fields)
		if err != nil {
			return nil, err
		}

		err = u.Reserve(int(poolSize))
		if err != nil {
			return nil, err
		}

		for i, v := range vecs {
			if err := u.Add(uint64(i+1), v); err != nil {
				t.Error(err)
				return nil, err
			}
		}

		return u, nil
	}
	tests := []test{
		{
			name: "return vector id after the same vector inserted",
			args: args{
				q: []float32{0, 1, 2, 3, 4, 5, 6, 7, 8},
				k: 5,
			},
			fields: fields{
				idxPath:          idxTempDir(t),
				quantizationType: "F32",
				metricType:       "cosine",
				dimension:        9,
				connectivity:     0,
				expansionAdd:     0,
				expansionSearch:  0,
				multi:            false,
			},
			createFunc: func(t *testing.T, fields fields) (Usearch, error) {
				t.Helper()
				vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			want: want{
				want: []algorithm.SearchResult{
					{ID: uint32(1), Distance: 0},
				},
			},
		},
		{
			name: "resturn vector id after the nearby vector inserted",
			args: args{
				q: []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
				k: 5,
			},
			fields: fields{
				idxPath:          idxTempDir(t),
				quantizationType: "F32",
				metricType:       "cosine",
				dimension:        9,
				connectivity:     0,
				expansionAdd:     0,
				expansionSearch:  0,
				multi:            false,
			},
			createFunc: func(t *testing.T, fields fields) (Usearch, error) {
				t.Helper()
				iv := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{iv}, 1)
			},
			want: want{
				want: []algorithm.SearchResult{
					{ID: uint32(1), Distance: 1},
				},
			},
		},
		{
			name: "return limited result after insert 10 vectors with limited size 3",
			args: args{
				q: []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
				k: 3,
			},
			fields: fields{
				idxPath:          idxTempDir(t),
				quantizationType: "F32",
				metricType:       "cosine",
				dimension:        9,
				connectivity:     0,
				expansionAdd:     0,
				expansionSearch:  0,
				multi:            false,
			},
			createFunc: func(t *testing.T, fields fields) (Usearch, error) {
				t.Helper()
				ivs := [][]float32{ // insert 10 vec
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{2, 3, 4, 5, 6, 7, 8, 9, 10},
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{2, 3, 4, 5, 6, 7, 8, 9, 10},
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{2, 3, 4, 5, 6, 7, 8, 9, 10},
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{2, 3, 4, 5, 6, 7, 8, 9, 10},
					{2, 3, 4, 5, 6, 7, 8, 9, 10},
					{2, 3, 4, 5, 6, 7, 8, 9, math.MaxFloat32},
				}

				return insertCreateFunc(t, fields, ivs, 10)
			},
			want: want{
				want: []algorithm.SearchResult{
					{ID: uint32(10), Distance: 3},
					{ID: uint32(9), Distance: 3},
					{ID: uint32(8), Distance: 3},
				},
			},
		},
		{
			name: "return most accurate result after insert 10 vectors with limited size 5",
			args: args{
				q: []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
				k: 5,
			},
			fields: fields{
				idxPath:          idxTempDir(t),
				quantizationType: "F32",
				metricType:       "cosine",
				dimension:        9,
				connectivity:     0,
				expansionAdd:     0,
				expansionSearch:  0,
				multi:            false,
			},
			createFunc: func(t *testing.T, fields fields) (Usearch, error) {
				t.Helper()
				ivs := [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},    // vec id 1
					{2, 3, 4, 5, 6, 7, 8, 9, 10},   // vec id 2
					{0, 1, 2, 3, 4, 5, 6, 7, 8},    // vec id 3
					{2, 3, 4, 5, 6, 7, 8, 9, 10},   // vec id 4
					{0, 1, 2, 3, 4, 5, 6, 7, 8},    // vec id 5
					{2, 3, 4, 5, 6, 7, 8, 9, 10},   // vec id 6
					{2, 3, 4, 5, 6, 7, 8, 9, 9.04}, // vec id 7
					{2, 3, 4, 5, 6, 7, 8, 9, 9.03}, // vec id 8
					{1, 2, 3, 4, 5, 6, 7, 8, 9.01}, // vec id 9
					{1, 2, 3, 4, 5, 6, 7, 8, 9.02}, // vec id 10
				}

				return insertCreateFunc(t, fields, ivs, 10)
			},
			want: want{
				want: []algorithm.SearchResult{
					{ID: uint32(9), Distance: 2.384185791015625e-07},
					{ID: uint32(10), Distance: 5.364418029785156e-07},
					{ID: uint32(6), Distance: 3},
					{ID: uint32(4), Distance: 3},
					{ID: uint32(2), Distance: 3},
				},
			},
		},
		{
			name: "return nothing if the search dimension is less than the inserted vector",
			args: args{
				q: []float32{0, 1, 2, 3, 4, 5, 6, 7},
				k: 5,
			},
			fields: fields{
				idxPath:          idxTempDir(t),
				quantizationType: "F32",
				metricType:       "cosine",
				dimension:        9,
				connectivity:     0,
				expansionAdd:     0,
				expansionSearch:  0,
				multi:            false,
			},
			createFunc: func(t *testing.T, fields fields) (Usearch, error) {
				t.Helper()
				vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			want: want{
				err: errors.New("incompatible dimension size detected\trequested: 8,\tconfigured: 9"),
			},
		},
		{
			name: "return nothing if the search dimension is more than the inserted vector",
			args: args{
				q: []float32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				k: 5,
			},
			fields: fields{
				idxPath:          idxTempDir(t),
				quantizationType: "F32",
				metricType:       "cosine",
				dimension:        9,
				connectivity:     0,
				expansionAdd:     0,
				expansionSearch:  0,
				multi:            false,
			},
			createFunc: func(t *testing.T, fields fields) (Usearch, error) {
				t.Helper()
				vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			want: want{
				err: errors.New("incompatible dimension size detected\trequested: 10,\tconfigured: 9"),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()

			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc == nil {
				test.afterFunc = defaultAfterFunc
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			if test.createFunc == nil {
				test.createFunc = defaultCreateFunc
			}

			u, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}

			got, err := u.Search(test.args.q, test.args.k)
			if err := checkFunc(test.want, got, u, err); err != nil {
				tt.Errorf("error = %v", err)
			}

			if err := test.afterFunc(tt, u); err != nil {
				tt.Error(err)
			}
		})
	}
}
