//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package ngt provides implementation of Go API for https://github.com/yahoojapan/NGT
package ngt

import (
	//"C"
	"math"
	"os"
	"reflect"
	"strings"
	"sync"

	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/test/comparator"
	"go.uber.org/goleak"
)

var (
	ngtComparator = []comparator.Option{
		comparator.AllowUnexported(ngt{}),
		// ignore C dependencies
		comparator.IgnoreFields(ngt{},
			"dimension", "prop", "ebuf", "index", "ospace"),
		comparator.RWMutexComparer,
	}
)

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want NGT
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, NGT, error) error
		beforeFunc func(args)
		afterFunc  func(t *testing.T, args args, w want, got NGT)
	}
	defaultCheckFunc := func(w want, got NGT, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		// comparator for idxPath
		comparators := append(ngtComparator, comparator.CompareField("idxPath", cmp.Comparer(func(s1, s2 string) bool {
			return strings.HasPrefix(s1, "/tmp/ngt-") || strings.HasPrefix(s1, "/tmp/ngt-")
		})))

		if diff := comparator.Diff(got, w.want, comparators...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}

		// check file is created in idxPath
		if got != nil {
			if ngt, ok := got.(*ngt); ok {
				if _, err := os.Stat(ngt.idxPath); err == os.ErrNotExist {
					return errors.Errorf("index file not exists, path: %s", ngt.idxPath)
				}
			}
		}

		return nil
	}
	defaultAfterFunc := func(t *testing.T, args args, w want, got NGT) {
		if got == nil {
			return
		}

		if ngt, ok := got.(*ngt); ok {
			if _, err := os.Stat(ngt.idxPath); err != os.ErrNotExist {
				_ = os.RemoveAll(ngt.idxPath)
			}
		}
	}
	tests := []test{
		{
			name: "return NGT when no option is set",
			args: args{
				opts: nil,
			},
			want: want{
				want: &ngt{
					// these options is in defaultOpts list, but these fields are ignored because of cgo dependencies
					// WithDimension(minimumDimensionSize),
					// WithCreationEdgeSize(10),
					// WithSearchEdgeSize(40),
					// WithDistanceType(L2),
					idxPath:             "/tmp/ngt-",
					radius:              DefaultRadius,
					epsilon:             DefaultEpsilon,
					poolSize:            DefaultPoolSize,
					bulkInsertChunkSize: 100,
					objectType:          Float,
					mu:                  &sync.RWMutex{},
				},
			},
		},
		{
			name: "return NGT when option is set",
			args: args{
				opts: []Option{
					WithObjectType(Uint8),
				},
			},
			want: want{
				want: &ngt{
					// these options is in defaultOpts list, but these fields are ignored because of cgo dependencies
					// WithDimension(minimumDimensionSize),
					// WithCreationEdgeSize(10),
					// WithSearchEdgeSize(40),
					// WithDistanceType(L2),
					idxPath:             "/tmp/ngt-",
					radius:              DefaultRadius,
					epsilon:             DefaultEpsilon,
					poolSize:            DefaultPoolSize,
					bulkInsertChunkSize: 100,
					objectType:          Uint8,
					mu:                  &sync.RWMutex{},
				},
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
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

			test.afterFunc(tt, test.args, test.want, got)
		})
	}
}

func Test_ngt_Search(t *testing.T) {
	type args struct {
		vec     []float32
		size    int
		epsilon float32
		radius  float32
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		//	dimension           C.int32_t
		dimension  int
		objectType objectType
		radius     float32
		epsilon    float32
		poolSize   uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
		want []SearchResult
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(t *testing.T, fields fields) (NGT, error)
		want       want
		checkFunc  func(want, []SearchResult, NGT, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(WithInMemoryMode(true),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(int(fields.dimension)),
		)
	}
	defaultCheckFunc := func(w want, got []SearchResult, n NGT, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		comparators := []comparator.Option{
			comparator.CompareField("Distance", cmp.Comparer(func(s1, s2 float32) bool {
				if s1 == 0 { // if vec1 is same as vec2, the distance should be same
					return s2 == 0
				}
				// by setting non-zero value in test case, it will only check if both got/want is non-zero
				return s1 != 0 && s2 != 0
			}))}

		if diff := comparator.Diff(got, w.want, comparators...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}

		return nil
	}
	insertCreateFunc := func(t *testing.T, fields fields, vecs [][]float32, poolSize uint32) (NGT, error) { // create func with insert/index
		t.Helper()

		ngt, err := defaultCreateFunc(t, fields)
		if err != nil {
			return nil, err
		}

		for _, vec := range vecs {
			if _, err := ngt.Insert(vec); err != nil {
				return nil, err
			}
		}
		if err := ngt.CreateIndex(poolSize); err != nil {
			return nil, err
		}

		return ngt, nil
	}
	tests := []test{
		func() test {
			vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

			return test{
				name: "return result after insert with same vec",
				args: args{
					vec:     vec,
					size:    5,
					epsilon: 0,
					radius:  0,
				},
				fields: fields{
					inMemory:            false,
					bulkInsertChunkSize: 100,
					dimension:           9,
					objectType:          Uint8,
					radius:              float32(-1.0),
					epsilon:             float32(0.01),
				},
				createFunc: func(t *testing.T, fields fields) (NGT, error) {
					t.Helper()

					return insertCreateFunc(t, fields, [][]float32{vec}, 1)
				},
				want: want{
					want: []SearchResult{
						{ID: uint32(1), Distance: 0},
					},
				},
			}
		}(),
		func() test {
			iv := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}  // insert vec
			vec := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9} // search vec

			return test{
				name: "resturn result after insert with nearby vec",
				args: args{
					vec:  vec,
					size: 5,
				},
				fields: fields{
					inMemory:            false,
					bulkInsertChunkSize: 100,
					dimension:           9,
					objectType:          Uint8,
					radius:              float32(-1.0),
					epsilon:             float32(0.01),
				},
				createFunc: func(t *testing.T, fields fields) (NGT, error) {
					t.Helper()

					return insertCreateFunc(t, fields, [][]float32{iv}, 1)
				},
				want: want{
					want: []SearchResult{
						{ID: uint32(1), Distance: 1},
					},
				},
			}
		}(),
		func() test {
			ivs := [][]float32{ // insert vec
				{0, 1, 2, 3, 4, 5, 6, 7, 8},
				{2, 3, 4, 5, 6, 7, 8, 9, 10},
				{2, 3, 4, 5, math.MaxFloat32 / 2, 7, 8, 9, math.MaxFloat32},
			}
			vec := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9} // search vec

			return test{
				name: "return result after insert with multiple vecs",
				args: args{
					vec:  vec,
					size: 5,
				},
				fields: fields{
					inMemory:            false,
					bulkInsertChunkSize: 100,
					dimension:           9,
					objectType:          Uint8,
					radius:              float32(-1.0),
					epsilon:             float32(0.01),
				},
				createFunc: func(t *testing.T, fields fields) (NGT, error) {
					t.Helper()

					return insertCreateFunc(t, fields, ivs, 1)
				},
				want: want{
					want: []SearchResult{
						{ID: uint32(1), Distance: 3},
						{ID: uint32(2), Distance: 3},
						{ID: uint32(3), Distance: 3},
					},
				},
			}
		}(),
		func() test {
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
			vec := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9} // search vec

			return test{
				name: "return limited result after insert 10 vecs with limited size 5",
				args: args{
					vec:  vec,
					size: 5,
				},
				fields: fields{
					inMemory:            false,
					bulkInsertChunkSize: 100,
					dimension:           9,
					objectType:          Uint8,
					radius:              float32(-1.0),
					epsilon:             float32(0.01),
				},
				createFunc: func(t *testing.T, fields fields) (NGT, error) {
					t.Helper()

					return insertCreateFunc(t, fields, ivs, 1)
				},
				want: want{
					want: []SearchResult{
						{ID: uint32(1), Distance: 3},
						{ID: uint32(2), Distance: 3},
						{ID: uint32(3), Distance: 3},
						{ID: uint32(4), Distance: 3},
						{ID: uint32(5), Distance: 3},
					},
				},
			}
		}(),
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			if test.createFunc == nil {
				test.createFunc = defaultCreateFunc
			}

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}

			got, err := n.Search(test.args.vec, test.args.size, test.args.epsilon, test.args.radius)
			if err := test.checkFunc(test.want, got, n, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_Insert(t *testing.T) {
	type args struct {
		vec []float32
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		// dimension           C.int32_t
		objectType objectType
		radius     float32
		epsilon    float32
		poolSize   uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
		want uint
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got uint, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           vec: nil,
		       },
		       fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           vec: nil,
		           },
		           fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				// dimension:           test.fields.dimension,
				objectType: test.fields.objectType,
				radius:     test.fields.radius,
				epsilon:    test.fields.epsilon,
				poolSize:   test.fields.poolSize,
				// prop:                test.fields.prop,
				// ebuf:                test.fields.ebuf,
				// index:               test.fields.index,
				// ospace:              test.fields.ospace,
				mu: test.fields.mu,
			}

			got, err := n.Insert(test.args.vec)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_InsertCommit(t *testing.T) {
	type args struct {
		vec      []float32
		poolSize uint32
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		//	dimension           C.int32_t
		objectType objectType
		radius     float32
		epsilon    float32
		poolSize   uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
		want uint
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, uint, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got uint, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           vec: nil,
		           poolSize: 0,
		       },
		       fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           vec: nil,
		           poolSize: 0,
		           },
		           fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				// dimension:           test.fields.dimension,
				objectType: test.fields.objectType,
				radius:     test.fields.radius,
				epsilon:    test.fields.epsilon,
				poolSize:   test.fields.poolSize,
				// prop:                test.fields.prop,
				// ebuf:                test.fields.ebuf,
				// index:               test.fields.index,
				// ospace:              test.fields.ospace,
				mu: test.fields.mu,
			}

			got, err := n.InsertCommit(test.args.vec, test.args.poolSize)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_BulkInsert(t *testing.T) {
	type args struct {
		vecs [][]float32
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		// dimension           C.int32_t
		objectType objectType
		radius     float32
		epsilon    float32
		// poolSize            uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
		want  []uint
		want1 []error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []uint, []error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []uint, got1 []error) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           vecs: nil,
		       },
		       fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           vecs: nil,
		           },
		           fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				// dimension:           test.fields.dimension,
				objectType: test.fields.objectType,
				radius:     test.fields.radius,
				epsilon:    test.fields.epsilon,
				// poolSize:            test.fields.poolSize,
				// prop:                test.fields.prop,
				// ebuf:                test.fields.ebuf,
				// index:               test.fields.index,
				// ospace:              test.fields.ospace,
				mu: test.fields.mu,
			}

			got, got1 := n.BulkInsert(test.args.vecs)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_BulkInsertCommit(t *testing.T) {
	type args struct {
		vecs     [][]float32
		poolSize uint32
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		//	dimension           C.int32_t
		objectType objectType
		radius     float32
		epsilon    float32
		poolSize   uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
		want  []uint
		want1 []error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []uint, []error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []uint, got1 []error) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		if !reflect.DeepEqual(got1, w.want1) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got1, w.want1)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           vecs: nil,
		           poolSize: 0,
		       },
		       fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           vecs: nil,
		           poolSize: 0,
		           },
		           fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				// dimension:           test.fields.dimension,
				objectType: test.fields.objectType,
				radius:     test.fields.radius,
				epsilon:    test.fields.epsilon,
				poolSize:   test.fields.poolSize,
				// prop:                test.fields.prop,
				// ebuf:                test.fields.ebuf,
				// index:               test.fields.index,
				// ospace:              test.fields.ospace,
				mu: test.fields.mu,
			}

			got, got1 := n.BulkInsertCommit(test.args.vecs, test.args.poolSize)
			if err := test.checkFunc(test.want, got, got1); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_CreateAndSaveIndex(t *testing.T) {
	type args struct {
		poolSize uint32
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		// dimension           C.int32_t
		objectType objectType
		radius     float32
		epsilon    float32
		poolSize   uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           poolSize: 0,
		       },
		       fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           poolSize: 0,
		           },
		           fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				// dimension:           test.fields.dimension,
				objectType: test.fields.objectType,
				radius:     test.fields.radius,
				epsilon:    test.fields.epsilon,
				poolSize:   test.fields.poolSize,
				// prop:                test.fields.prop,
				// ebuf:                test.fields.ebuf,
				// index:               test.fields.index,
				// ospace:              test.fields.ospace,
				mu: test.fields.mu,
			}

			err := n.CreateAndSaveIndex(test.args.poolSize)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_CreateIndex(t *testing.T) {
	type args struct {
		poolSize uint32
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		// dimension           C.int32_t
		objectType objectType
		radius     float32
		epsilon    float32
		poolSize   uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           poolSize: 0,
		       },
		       fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           poolSize: 0,
		           },
		           fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				// dimension:           test.fields.dimension,
				objectType: test.fields.objectType,
				radius:     test.fields.radius,
				epsilon:    test.fields.epsilon,
				poolSize:   test.fields.poolSize,
				// prop:                test.fields.prop,
				// ebuf:                test.fields.ebuf,
				// index:               test.fields.index,
				// ospace:              test.fields.ospace,
				mu: test.fields.mu,
			}

			err := n.CreateIndex(test.args.poolSize)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_SaveIndex(t *testing.T) {
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		// dimension           C.int32_t
		objectType objectType
		radius     float32
		epsilon    float32
		poolSize   uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				// dimension:           test.fields.dimension,
				objectType: test.fields.objectType,
				radius:     test.fields.radius,
				epsilon:    test.fields.epsilon,
				poolSize:   test.fields.poolSize,
				// prop:                test.fields.prop,
				// ebuf:                test.fields.ebuf,
				// index:               test.fields.index,
				// ospace:              test.fields.ospace,
				mu: test.fields.mu,
			}

			err := n.SaveIndex()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_Remove(t *testing.T) {
	type args struct {
		id uint
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		// dimension           C.int32_t
		objectType objectType
		radius     float32
		epsilon    float32
		poolSize   uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           id: 0,
		       },
		       fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           id: 0,
		           },
		           fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				// dimension:           test.fields.dimension,
				objectType: test.fields.objectType,
				radius:     test.fields.radius,
				epsilon:    test.fields.epsilon,
				poolSize:   test.fields.poolSize,
				// prop:                test.fields.prop,
				// ebuf:                test.fields.ebuf,
				// index:               test.fields.index,
				// ospace:              test.fields.ospace,
				mu: test.fields.mu,
			}

			err := n.Remove(test.args.id)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_BulkRemove(t *testing.T) {
	type args struct {
		ids []uint
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		// dimension           C.int32_t
		objectType objectType
		radius     float32
		epsilon    float32
		poolSize   uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ids: nil,
		       },
		       fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           ids: nil,
		           },
		           fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				// dimension:           test.fields.dimension,
				objectType: test.fields.objectType,
				radius:     test.fields.radius,
				epsilon:    test.fields.epsilon,
				poolSize:   test.fields.poolSize,
				// prop:                test.fields.prop,
				// ebuf:                test.fields.ebuf,
				// index:               test.fields.index,
				// ospace:              test.fields.ospace,
				mu: test.fields.mu,
			}

			err := n.BulkRemove(test.args.ids...)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_GetVector(t *testing.T) {
	type args struct {
		id uint
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		// dimension           C.int32_t
		objectType objectType
		radius     float32
		epsilon    float32
		poolSize   uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
		want []float32
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []float32, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []float32, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           id: 0,
		       },
		       fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           args: args {
		           id: 0,
		           },
		           fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				// dimension:           test.fields.dimension,
				objectType: test.fields.objectType,
				radius:     test.fields.radius,
				epsilon:    test.fields.epsilon,
				poolSize:   test.fields.poolSize,
				// prop:                test.fields.prop,
				// ebuf:                test.fields.ebuf,
				// index:               test.fields.index,
				// ospace:              test.fields.ospace,
				mu: test.fields.mu,
			}

			got, err := n.GetVector(test.args.id)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_ngt_Close(t *testing.T) {
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		// dimension           C.int32_t
		objectType objectType
		radius     float32
		epsilon    float32
		poolSize   uint32
		// prop                C.NGTProperty
		// ebuf                C.NGTError
		// index               C.NGTIndex
		// ospace              C.NGTObjectSpace
		mu *sync.RWMutex
	}
	type want struct {
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		       },
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
		           fields: fields {
		           idxPath: "",
		           inMemory: false,
		           bulkInsertChunkSize: 0,
		           dimension: nil,
		           objectType: nil,
		           radius: 0,
		           epsilon: 0,
		           poolSize: 0,
		           prop: nil,
		           ebuf: nil,
		           index: nil,
		           ospace: nil,
		           mu: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				// dimension:           test.fields.dimension,
				objectType: test.fields.objectType,
				radius:     test.fields.radius,
				epsilon:    test.fields.epsilon,
				poolSize:   test.fields.poolSize,
				// prop:                test.fields.prop,
				// ebuf:                test.fields.ebuf,
				// index:               test.fields.index,
				// ospace:              test.fields.ospace,
				mu: test.fields.mu,
			}

			n.Close()
			if err := test.checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
