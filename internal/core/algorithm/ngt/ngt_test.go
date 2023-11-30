//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"context"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/test/comparator"
	"github.com/vdaas/vald/internal/test/goleak"
)

var (
	ngtComparator = []comparator.Option{
		comparator.AllowUnexported(ngt{}),
		// !!! These fields will not be verified in the entire test
		// Do not validate C dependencies
		comparator.IgnoreFields(ngt{},
			"dimension", "prop", "epool", "index", "ospace"),
		comparator.RWMutexComparer,
		comparator.ErrorComparer,
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

	defaultAfterFunc = func(t *testing.T, n NGT) error {
		t.Helper()

		if n == nil {
			return nil
		}

		n.Close()
		if ngt, ok := n.(*ngt); ok {
			if !ngt.inMemory {
				return os.RemoveAll(ngt.idxPath)
			}
		}
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

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want NGT
		err  error
	}
	type test struct {
		name        string
		args        args
		want        want
		comparators []comparator.Option
		checkFunc   func(want, NGT, error, ...comparator.Option) error
		beforeFunc  func(args)
		afterFunc   func(*testing.T, NGT) error
	}
	defaultComprators := append(ngtComparator, comparator.CompareField("idxPath", comparator.Comparer(func(s1, s2 string) bool {
		return s1 == s2
	})))
	defaultCheckFunc := func(w want, got NGT, err error, comparators ...comparator.Option) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		if diff := comparator.Diff(got, w.want, comparators...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}

		// check file is created in idxPath
		if ngt, ok := got.(*ngt); ok {
			if _, err := os.Stat(ngt.idxPath); errors.Is(err, fs.ErrNotExist) {
				return errors.Errorf("index file not exists, path: %s", ngt.idxPath)
			}
		}

		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "return NGT when no option is set",
				args: args{
					opts: nil,
				},
				want: want{
					want: &ngt{
						idxPath:             "/tmp/ngt-",
						radius:              DefaultRadius,
						epsilon:             DefaultEpsilon,
						poolSize:            DefaultPoolSize,
						bulkInsertChunkSize: 100,
						objectType:          Float,
						mu:                  &sync.RWMutex{},
					},
				},
				comparators: append(ngtComparator, comparator.CompareField("idxPath", comparator.Comparer(func(s1, s2 string) bool {
					return strings.HasPrefix(s1, "/tmp/ngt-") || strings.HasPrefix(s2, "/tmp/ngt-")
				}))),
			}
		}(),
		func() test {
			idxPath := idxTempDir(t)
			return test{
				name: "return NGT when 1 option is set",
				args: args{
					opts: []Option{
						WithIndexPath(idxPath),
					},
				},
				want: want{
					want: &ngt{
						idxPath:             idxPath,
						radius:              DefaultRadius,
						epsilon:             DefaultEpsilon,
						poolSize:            DefaultPoolSize,
						bulkInsertChunkSize: 100,
						objectType:          Float,
						mu:                  &sync.RWMutex{},
					},
				},
			}
		}(),
		func() test {
			idxPath := idxTempDir(t)
			return test{
				name: "return NGT when multiple options are set",
				args: args{
					opts: []Option{
						WithObjectType(Uint8),
						WithDefaultPoolSize(100),
						WithIndexPath(idxPath),
					},
				},
				want: want{
					want: &ngt{
						idxPath:             idxPath,
						radius:              DefaultRadius,
						epsilon:             DefaultEpsilon,
						poolSize:            100,
						bulkInsertChunkSize: 100,
						objectType:          Uint8,
						mu:                  &sync.RWMutex{},
					},
				},
			}
		}(),
		func() test {
			return test{
				name: "return error when option return error",
				args: args{
					opts: []Option{
						WithDimension(1),
					},
				},
				want: want{
					err: errors.NewErrCriticalOption("dimension", 1, errors.ErrInvalidDimensionSize(1, algorithm.MaximumVectorDimensionSize)),
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
			if test.afterFunc == nil {
				test.afterFunc = defaultAfterFunc
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			comparators := test.comparators
			if test.comparators == nil || len(test.comparators) == 0 {
				comparators = defaultComprators
			}

			got, err := New(test.args.opts...)
			defer func() {
				if err := test.afterFunc(tt, got); err != nil {
					tt.Error(err)
				}
			}()
			if err := checkFunc(test.want, got, err, comparators...); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestLoad(t *testing.T) {
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
		checkFunc  func(context.Context, want, NGT, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCheckFunc := func(_ context.Context, w want, got NGT, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		// comparator for idxPath
		comparators := append(ngtComparator, comparator.CompareField("idxPath", comparator.Comparer(func(s1, s2 string) bool {
			return s1 == s2
		})))

		if diff := comparator.Diff(got, w.want, comparators...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}

		// check file is created in idxPath
		if ngt, ok := got.(*ngt); ok {
			if _, err := os.Stat(ngt.idxPath); errors.Is(err, fs.ErrNotExist) {
				return errors.Errorf("index file not exists, path: %s", ngt.idxPath)
			}
		}

		return nil
	}
	tests := []test{
		// uint
		func() test {
			idxPath := idxTempDir(t)
			opts := []Option{
				WithDimension(9),
				WithIndexPath(idxPath),
				WithObjectType(Uint8),
			}

			return test{
				name: "Load success to restore backup file with no vector inserted (uint)",
				args: args{
					opts: opts,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()

					n, err := New(opts...)
					if err != nil {
						t.Error(err)
					}

					if err = n.CreateAndSaveIndex(1); err != nil {
						t.Error(err)
					}
					n.Close()
				},
				want: want{
					want: &ngt{
						idxPath:             idxPath,
						radius:              DefaultRadius,
						epsilon:             DefaultEpsilon,
						poolSize:            DefaultPoolSize,
						bulkInsertChunkSize: 100,
						objectType:          Uint8,
						mu:                  &sync.RWMutex{},
					},
				},
				checkFunc: func(ctx context.Context, w want, n NGT, e error) error {
					if err := defaultCheckFunc(ctx, w, n, e); err != nil {
						return err
					}

					// check no vector exists
					v, err := n.GetVector(1)
					if err == nil || len(v) > 0 {
						return errors.Errorf("vector exists but not inserted, vec: %s", v)
					}

					// check no vector can be searched
					vs, err := n.Search(ctx, []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}, 10, 0, 0)
					if err != nil && !errors.Is(err, errors.ErrEmptySearchResult) {
						return err
					}
					if len(vs) != 0 {
						t.Errorf("got vec is not the same as inserted vec, got: %v", vs)
					}

					return nil
				},
			}
		}(),
		func() test {
			idxPath := idxTempDir(t)
			vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}
			opts := []Option{
				WithDimension(9),
				WithIndexPath(idxPath),
				WithObjectType(Uint8),
			}

			return test{
				name: "Load success to restore backup file with 1 vector inserted (uint)",
				args: args{
					opts: opts,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()

					n, err := New(opts...)
					if err != nil {
						t.Error(err)
					}

					if _, err = n.Insert(vec); err != nil {
						t.Error(err)
					}
					if err = n.CreateAndSaveIndex(1); err != nil {
						t.Error(err)
					}
					n.Close()
				},
				want: want{
					want: &ngt{
						idxPath:             idxPath,
						radius:              DefaultRadius,
						epsilon:             DefaultEpsilon,
						poolSize:            DefaultPoolSize,
						bulkInsertChunkSize: 100,
						objectType:          Uint8,
						mu:                  &sync.RWMutex{},
					},
				},
				checkFunc: func(ctx context.Context, w want, n NGT, e error) error {
					if err := defaultCheckFunc(ctx, w, n, e); err != nil {
						return err
					}

					// check inserted vector exists
					v, err := n.GetVector(1)
					if err != nil {
						return err
					}
					if !reflect.DeepEqual(v, vec) {
						t.Errorf("got vec is not the same as inserted vec, got: %v, want: %v", v, vec)
					}

					// check inserted vector can be searched
					vs, err := n.Search(ctx, vec, 10, 0, 0)
					if err != nil && !errors.Is(err, errors.ErrEmptySearchResult) {
						return err
					}
					if len(vs) != 1 || vs[0].ID != 1 || vs[0].Distance != 0 {
						t.Errorf("got vec is not the same as inserted vec, got: %v, want: %v", vs, vec)
					}

					return nil
				},
			}
		}(),
		// float
		func() test {
			idxPath := idxTempDir(t)
			opts := []Option{
				WithDimension(9),
				WithIndexPath(idxPath),
				WithObjectType(Float),
			}

			return test{
				name: "Load success to restore backup file with no vector inserted (float)",
				args: args{
					opts: opts,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()

					n, err := New(opts...)
					if err != nil {
						t.Error(err)
					}

					if err = n.CreateAndSaveIndex(1); err != nil {
						t.Error(err)
					}
					n.Close()
				},
				want: want{
					want: &ngt{
						idxPath:             idxPath,
						radius:              DefaultRadius,
						epsilon:             DefaultEpsilon,
						poolSize:            DefaultPoolSize,
						bulkInsertChunkSize: 100,
						objectType:          Float,
						mu:                  &sync.RWMutex{},
					},
				},
				checkFunc: func(ctx context.Context, w want, n NGT, e error) error {
					if err := defaultCheckFunc(ctx, w, n, e); err != nil {
						return err
					}

					// check no vector exists
					v, err := n.GetVector(1)
					if err == nil || len(v) > 0 {
						return errors.Errorf("vector exists but not inserted, vec: %s", v)
					}

					// check no vector can be searched
					vs, err := n.Search(ctx, []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}, 10, 0, 0)
					if err != nil && !errors.Is(err, errors.ErrEmptySearchResult) {
						return err
					}
					if len(vs) != 0 {
						t.Errorf("got vec is not the same as inserted vec, got: %v", vs)
					}

					return nil
				},
			}
		}(),
		func() test {
			idxPath := idxTempDir(t)
			vec := []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}
			opts := []Option{
				WithDimension(9),
				WithIndexPath(idxPath),
				WithObjectType(Float),
			}

			return test{
				name: "Load success to restore backup file with 1 vector inserted (float)",
				args: args{
					opts: opts,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()

					n, err := New(opts...)
					if err != nil {
						t.Error(err)
					}

					if _, err = n.Insert(vec); err != nil {
						t.Error(err)
					}
					if err = n.CreateAndSaveIndex(1); err != nil {
						t.Error(err)
					}
					n.Close()
				},
				want: want{
					want: &ngt{
						idxPath:             idxPath,
						radius:              DefaultRadius,
						epsilon:             DefaultEpsilon,
						poolSize:            DefaultPoolSize,
						bulkInsertChunkSize: 100,
						objectType:          Float,
						mu:                  &sync.RWMutex{},
					},
				},
				checkFunc: func(ctx context.Context, w want, n NGT, e error) error {
					if err := defaultCheckFunc(ctx, w, n, e); err != nil {
						return err
					}

					// check inserted vector exists
					v, err := n.GetVector(1)
					if err != nil {
						return err
					}
					if !reflect.DeepEqual(v, vec) {
						t.Errorf("got vec is not the same as inserted vec, got: %v, want: %v", v, vec)
					}

					// check inserted vector can be searched
					vs, err := n.Search(ctx, vec, 10, 0, 0)
					if err != nil && !errors.Is(err, errors.ErrEmptySearchResult) {
						return err
					}
					if len(vs) != 1 || vs[0].ID != 1 || vs[0].Distance != 0 {
						t.Errorf("got vec is not the same as inserted vec, got: %v, want: %v", vs, vec)
					}

					return nil
				},
			}
		}(),
		// other
		func() test {
			idxPath := idxTempDir(t)
			vec := []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}
			opts := []Option{
				WithDimension(9),
				WithIndexPath(idxPath),
				WithObjectType(Float),
			}

			return test{
				name: "Load failed with ErrIndexFileNotFound if the index path is not exists",
				args: args{
					opts: opts,
				},
				want: want{
					want: nil,
					err:  errors.ErrIndexFileNotFound,
				},
				checkFunc: func(ctx context.Context, w want, n NGT, e error) error {
					if err := defaultCheckFunc(ctx, w, n, e); err != nil {
						return err
					}

					// check no object returned
					if n != nil {
						// check no vector exists
						v, err := n.GetVector(1)
						if err == nil || len(v) > 0 {
							return errors.Errorf("vector exists but not inserted, vec: %s", v)
						}

						// check no vector can be searched
						vs, err := n.Search(ctx, vec, 10, 0, 0)
						if err != nil && !errors.Is(err, errors.ErrEmptySearchResult) {
							return err
						}
						if len(vs) != 0 {
							t.Errorf("got vec is not the same as inserted vec, got: %v", vs)
						}
					}

					return nil
				},
			}
		}(),
		func() test {
			idxPath := idxTempDir(t)
			opts := []Option{
				WithDimension(9),
				WithIndexPath(idxPath),
				WithObjectType(Float),
			}

			return test{
				name: "Load failed if the index path dir exists but no index file found which returns NGTError",
				args: args{
					opts: opts,
				},
				want: want{
					want: nil,
					err:  new(errors.NGTError),
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()
					if err := file.MkdirAll(idxPath, fs.ModePerm); err != nil {
						t.Error(err)
					}
				},
				checkFunc: func(_ context.Context, w want, n NGT, e error) error {
					if e != nil && !errors.As(e, w.err) {
						t.Error(e)
						return e
					}
					return nil
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
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc == nil {
				test.afterFunc = defaultAfterFunc
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got, err := Load(test.args.opts...)
			defer func() {
				if err := test.afterFunc(tt, got); err != nil {
					tt.Error(err)
				}
			}()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := checkFunc(ctx, test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_gen(t *testing.T) {
	type args struct {
		isLoad bool
		opts   []Option
	}
	type want struct {
		want NGT
		err  error
	}
	type test struct {
		name        string
		args        args
		want        want
		comparators []comparator.Option
		checkFunc   func(context.Context, want, NGT, error, ...comparator.Option) error
		beforeFunc  func(*testing.T, args)
		afterFunc   func(*testing.T, NGT) error
	}
	defaultComprators := append(ngtComparator, comparator.CompareField("idxPath", comparator.Comparer(func(s1, s2 string) bool {
		return s1 == s2
	})))
	defaultCheckFunc := func(_ context.Context, w want, got NGT, err error, comparators ...comparator.Option) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		if diff := comparator.Diff(got, w.want, comparators...); diff != "" {
			return errors.Errorf("diff: %s", diff)
		}

		// check file is created in idxPath
		if ngt, ok := got.(*ngt); ok {
			if _, err := os.Stat(ngt.idxPath); errors.Is(err, fs.ErrNotExist) {
				return errors.Errorf("index file not exists, path: %s", ngt.idxPath)
			}
		}

		return nil
	}
	tests := []test{
		{
			name: "gen success and do not load backup data",
			args: args{
				opts: nil,
			},
			want: want{
				want: &ngt{
					idxPath:             "/tmp/ngt-",
					radius:              DefaultRadius,
					epsilon:             DefaultEpsilon,
					poolSize:            DefaultPoolSize,
					bulkInsertChunkSize: 100,
					objectType:          Float,
					mu:                  &sync.RWMutex{},
				},
			},
			comparators: append(ngtComparator, comparator.CompareField("idxPath", comparator.Comparer(func(s1, s2 string) bool {
				return strings.HasPrefix(s1, "/tmp/ngt-") || strings.HasPrefix(s2, "/tmp/ngt-")
			}))),
		},
		func() test {
			idxPath := idxTempDir(t)
			vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}
			opts := []Option{
				WithDimension(9),
				WithIndexPath(idxPath),
				WithObjectType(Uint8),
			}

			return test{
				name: "gen success and load backup data success",
				args: args{
					isLoad: true,
					opts:   opts,
				},
				beforeFunc: func(t *testing.T, a args) {
					t.Helper()

					n, err := New(opts...)
					if err != nil {
						t.Error(err)
					}

					if _, err = n.Insert(vec); err != nil {
						t.Error(err)
					}
					if err = n.CreateAndSaveIndex(1); err != nil {
						t.Error(err)
					}
					n.Close()
				},
				want: want{
					want: &ngt{
						idxPath:             idxPath,
						radius:              DefaultRadius,
						epsilon:             DefaultEpsilon,
						poolSize:            DefaultPoolSize,
						bulkInsertChunkSize: 100,
						objectType:          Uint8,
						mu:                  &sync.RWMutex{},
					},
				},
				checkFunc: func(ctx context.Context, w want, n NGT, e error, comparators ...comparator.Option) error {
					if err := defaultCheckFunc(ctx, w, n, e, comparators...); err != nil {
						return err
					}

					// check inserted vector exists
					v, err := n.GetVector(1)
					if err != nil {
						return err
					}
					if !reflect.DeepEqual(v, vec) {
						t.Errorf("got vec is not the same as inserted vec, got: %v, want: %v", v, vec)
					}

					// check inserted vector can be searched
					vs, err := n.Search(ctx, vec, 10, 0, 0)
					if err != nil && !errors.Is(err, errors.ErrEmptySearchResult) {
						return err
					}
					if len(vs) != 1 || vs[0].ID != 1 || vs[0].Distance != 0 {
						t.Errorf("got vec is not the same as inserted vec, got: %v, want: %v", vs, vec)
					}

					return nil
				},
			}
		}(),
		{
			name: "return error when option return error",
			args: args{
				opts: []Option{
					WithDimension(1),
				},
			},
			want: want{
				err: errors.NewErrCriticalOption("dimension", 1, errors.ErrInvalidDimensionSize(1, algorithm.MaximumVectorDimensionSize)),
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc == nil {
				test.afterFunc = defaultAfterFunc
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			comparators := test.comparators
			if test.comparators == nil || len(test.comparators) == 0 {
				comparators = defaultComprators
			}

			got, err := gen(test.args.isLoad, test.args.opts...)
			defer func() {
				if err := test.afterFunc(tt, got); err != nil {
					tt.Error(err)
				}
			}()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if err := checkFunc(ctx, test.want, got, err, comparators...); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_setup(t *testing.T) {
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
		mu                  *sync.RWMutex
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
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "return nil when object type is uint8",
			fields: fields{
				objectType: Uint8,
			},
			want: want{},
		},
		{
			name: "return nil when object type is float",
			fields: fields{
				objectType: Float,
			},
			want: want{},
		},
		{
			name: "return nil when object type is invalid value",
			fields: fields{
				objectType: 999,
			},
			want: want{},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc == nil {
				test.afterFunc = defaultAfterFunc
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			n := &ngt{
				idxPath:             test.fields.idxPath,
				inMemory:            test.fields.inMemory,
				bulkInsertChunkSize: test.fields.bulkInsertChunkSize,
				objectType:          test.fields.objectType,
				radius:              test.fields.radius,
				epsilon:             test.fields.epsilon,
				poolSize:            test.fields.poolSize,
				mu:                  test.fields.mu,
			}
			defer func() {
				if err := test.afterFunc(tt, n); err != nil {
					tt.Error(err)
				}
			}()

			err := n.setup()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_loadOptions(t *testing.T) {
	type args struct {
		opts []Option
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
		mu                  *sync.RWMutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(t *testing.T, fields fields) (NGT, error)
		want       want
		checkFunc  func(want, *ngt, error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		n := &ngt{
			idxPath:             fields.idxPath,
			inMemory:            fields.inMemory,
			bulkInsertChunkSize: fields.bulkInsertChunkSize,
			objectType:          fields.objectType,
			radius:              fields.radius,
			epsilon:             fields.epsilon,
			poolSize:            fields.poolSize,
			mu:                  fields.mu,
		}
		if err := n.setup(); err != nil {
			return nil, err
		}
		return n, nil
	}
	defaultCheckFunc := func(w want, n *ngt, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		func() test {
			idxPath := "tmp/ngt-41"

			return test{
				name: "load option success",
				args: args{
					opts: []Option{
						WithIndexPath(idxPath),
					},
				},
				fields: fields{
					idxPath:    idxPath,
					objectType: Uint8,
					dimension:  9,
				},
				checkFunc: func(w want, n *ngt, e error) error {
					if err := defaultCheckFunc(w, n, e); err != nil {
						return err
					}

					if n.idxPath != idxPath {
						return errors.New("index path does not set")
					}

					return nil
				},
			}
		}(),
		{
			name: "load option failed with critical error",
			args: args{
				opts: []Option{
					func(n *ngt) error {
						return errors.NewErrCriticalOption("objectType", 1)
					},
				},
			},
			fields: fields{
				idxPath:    "tmp/ngt-42",
				objectType: Uint8,
				dimension:  9,
			},
			want: want{
				err: errors.NewErrCriticalOption("objectType", 1),
			},
		},
		{
			name: "load option failed with Ignoreable error",
			args: args{
				opts: []Option{
					func(n *ngt) error {
						return errors.NewErrIgnoredOption("object")
					},
				},
			},
			fields: fields{
				idxPath:    "tmp/ngt-43",
				objectType: Uint8,
				dimension:  9,
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

			obj, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			defer func() {
				if err := test.afterFunc(tt, obj); err != nil {
					tt.Error(err)
				}
			}()

			n := obj.(*ngt)

			err = n.loadOptions(test.args.opts...)
			if err := checkFunc(test.want, n, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_create(t *testing.T) {
	// This test is skipped because it requires ngt.prop to be set probably.
	// We cannot initialize ngt.prop since it is C dependencies.
	// This function is called by New(), and the ngt.prop is destoried in New(), so we cannot test this function individually.
	t.SkipNow()
}

func Test_ngt_open(t *testing.T) {
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
		mu                  *sync.RWMutex
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		createFunc func(*testing.T, fields) (NGT, error)
		want       want
		checkFunc  func(want, error) error
		beforeFunc func(*testing.T, fields)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		n := &ngt{
			idxPath:             fields.idxPath,
			inMemory:            fields.inMemory,
			bulkInsertChunkSize: fields.bulkInsertChunkSize,
			objectType:          fields.objectType,
			radius:              fields.radius,
			epsilon:             fields.epsilon,
			poolSize:            fields.poolSize,
			mu:                  fields.mu,
		}
		if err := n.setup(); err != nil {
			t.Error(err)
		}

		return n, nil
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "return nil when index exists",
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Float,
				mu:         &sync.RWMutex{},
			},
			beforeFunc: func(t *testing.T, fields fields) {
				t.Helper()

				n, err := New(
					WithIndexPath(fields.idxPath),
					WithDimension(9),
					WithObjectType(Float),
				)
				if err != nil {
					t.Error(err)
				}

				if _, err = n.Insert([]float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}); err != nil {
					t.Error(err)
				}

				if err = n.CreateAndSaveIndex(1); err != nil {
					t.Error(err)
				}
				n.Close()
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "return error when index path is not exists",
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Float,
				mu:         &sync.RWMutex{},
			},
			want: want{
				err: errors.ErrIndexFileNotFound,
			},
		},
		{
			name: "return error when index path contains no file",
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Float,
				mu:         &sync.RWMutex{},
			},
			beforeFunc: func(t *testing.T, fields fields) {
				t.Helper()
				_ = file.MkdirAll(fields.idxPath, fs.ModePerm)
			},
			checkFunc: func(w want, e error) error {
				if e == nil {
					return errors.New("error should be returned")
				}
				return nil
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.fields)
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

			obj, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			defer func() {
				if err := test.afterFunc(tt, obj); err != nil {
					tt.Error(err)
				}
			}()

			n, ok := obj.(*ngt)
			if !ok {
				tt.Fatal("cannot cast ngt")
			}

			err = n.open()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_loadObjectSpace(t *testing.T) {
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		fields     fields
		createFunc func(t *testing.T, fields fields) (NGT, error)
		want       want
		checkFunc  func(want, error) error
		beforeFunc func()
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(w want, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "return nil when load object space success",
			fields: fields{
				idxPath:    idxTempDir(t),
				dimension:  9,
				objectType: Float,
			},
			want: want{},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
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

			obj, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			defer func() {
				if err := test.afterFunc(tt, obj); err != nil {
					tt.Error(err)
				}
			}()
			n, ok := obj.(*ngt)
			if !ok {
				tt.Fatal("cannot cast ngt")
			}

			err = n.loadObjectSpace()
			if err := checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_Search(t *testing.T) {
	type args struct {
		ctx     context.Context
		vec     []float32
		size    int
		epsilon float32
		radius  float32
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
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
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(w want, got []SearchResult, n NGT, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if diff := comparator.Diff(got, w.want, searchResultComparator...); diff != "" {
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

		if _, err := ngt.BulkInsertCommit(vecs, poolSize); len(err) != 0 {
			t.Error(err)
			return nil, err[0]
		}

		return ngt, nil
	}
	tests := []test{
		// object type uint8
		{
			name: "return vector id after the same vector inserted (uint8)",
			args: args{
				ctx:     context.Background(),
				vec:     []float32{0, 1, 2, 3, 4, 5, 6, 7, 8},
				size:    5,
				epsilon: 0,
				radius:  0,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			want: want{
				want: []SearchResult{
					{ID: uint32(1), Distance: 0},
				},
			},
		},
		{
			name: "resturn vector id after the nearby vector inserted (uint8)",
			args: args{
				ctx:  context.Background(),
				vec:  []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
				size: 5,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				iv := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{iv}, 1)
			},
			want: want{
				want: []SearchResult{
					{ID: uint32(1), Distance: 1},
				},
			},
		},
		{
			name: "return vector ids after insert with multiple vectors (uint8)",
			args: args{
				ctx:  context.Background(),
				vec:  []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
				size: 5,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				ivs := [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{2, 3, 4, 5, 6, 7, 8, 9, 10},
					{2, 3, 4, 5, math.MaxFloat32 / 2, 7, 8, 9, math.MaxFloat32},
				}

				return insertCreateFunc(t, fields, ivs, 1)
			},
			want: want{
				want: []SearchResult{
					{ID: uint32(1), Distance: 3},
					{ID: uint32(2), Distance: 3},
					{ID: uint32(3), Distance: 3},
				},
			},
		},
		{
			name: "return limited result after insert 10 vectors with limited size 3 (uint8)",
			args: args{
				ctx:  context.Background(),
				vec:  []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
				size: 3,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
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

				return insertCreateFunc(t, fields, ivs, 1)
			},
			want: want{
				want: []SearchResult{
					{ID: uint32(1), Distance: 3},
					{ID: uint32(2), Distance: 3},
					{ID: uint32(3), Distance: 3},
				},
			},
		},
		{
			name: "return most accurate result after insert 10 vectors with limited size 5 (uint8)",
			args: args{
				ctx:  context.Background(),
				vec:  []float32{1, 2, 3, 4, 5, 6, 7, 8, 9},
				size: 5,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
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

				return insertCreateFunc(t, fields, ivs, 1)
			},
			want: want{
				want: []SearchResult{
					{ID: uint32(9), Distance: 0},
					{ID: uint32(10), Distance: 0},
					{ID: uint32(7), Distance: 3},
					{ID: uint32(8), Distance: 3},
					{ID: uint32(1), Distance: 3},
				},
			},
		},
		// object type float
		{
			name: "return vector id after the same vector inserted (float)",
			args: args{
				ctx:     context.Background(),
				vec:     []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
				size:    5,
				epsilon: 0,
				radius:  0,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			want: want{
				want: []SearchResult{
					{ID: uint32(1), Distance: 0},
				},
			},
		},
		{
			name: "resturn vector id after the nearby vector inserted (float)",
			args: args{
				ctx:  context.Background(),
				vec:  []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.91},
				size: 5,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				iv := []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}

				return insertCreateFunc(t, fields, [][]float32{iv}, 1)
			},
			want: want{
				want: []SearchResult{
					{ID: uint32(1), Distance: 1},
				},
			},
		},
		{
			name: "return vector ids after insert with multiple vectors (float)",
			args: args{
				ctx:  context.Background(),
				vec:  []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
				size: 5,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				ivs := [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
					{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10},
					{0.2, 0.3, 0.4, 0.5, math.MaxFloat32, 0.7, 0.8, 0.9, math.MaxFloat32},
				}

				return insertCreateFunc(t, fields, ivs, 1)
			},
			want: want{
				want: []SearchResult{
					{ID: uint32(1), Distance: 3},
					{ID: uint32(2), Distance: 3},
				},
			},
		},
		{
			name: "return limited result after insert 10 vectors with limited size 3 (float)",
			args: args{
				ctx:  context.Background(),
				vec:  []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
				size: 3,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				ivs := [][]float32{ // insert 10 vec
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},           // vec id 1
					{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},         // vec id 2
					{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10},        // vec id 3
					{0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11},       // vec id 4
					{0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12},      // vec id 5
					{0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13},     // vec id 6
					{0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14},    // vec id 7
					{0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15},   // vec id 8
					{0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16},  // vec id 9
					{0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17}, // vec id 10
				}

				return insertCreateFunc(t, fields, ivs, 1)
			},
			want: want{
				want: []SearchResult{
					{ID: uint32(2), Distance: 0},
					{ID: uint32(1), Distance: 3},
					{ID: uint32(3), Distance: 3},
				},
			},
		},
		{
			name: "return most accurate result after insert 10 vectors with limited size 5 (float)",
			args: args{
				ctx:  context.Background(),
				vec:  []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
				size: 5,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				ivs := [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},           // vec id 1
					{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},         // vec id 2
					{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10},        // vec id 3
					{0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11},       // vec id 4
					{0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12},      // vec id 5
					{0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13},     // vec id 6
					{0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14},    // vec id 7
					{0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15},   // vec id 8
					{0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16},  // vec id 9
					{0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17}, // vec id 10
				}

				return insertCreateFunc(t, fields, ivs, 1)
			},
			want: want{
				want: []SearchResult{
					{ID: uint32(2), Distance: 0},
					{ID: uint32(1), Distance: 3},
					{ID: uint32(3), Distance: 3},
					{ID: uint32(4), Distance: 3},
					{ID: uint32(5), Distance: 3},
				},
			},
		},
		// other cases
		{
			name: "return nothing if the search dimension is less than the inserted vector",
			args: args{
				ctx:     context.Background(),
				vec:     []float32{0, 1, 2, 3, 4, 5, 6, 7},
				size:    5,
				epsilon: 0,
				radius:  0,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
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
				ctx:     context.Background(),
				vec:     []float32{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
				size:    5,
				epsilon: 0,
				radius:  0,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			want: want{
				err: errors.New("incompatible dimension size detected\trequested: 10,\tconfigured: 9"),
			},
		},
		{
			name: "return ErrEmptySearchResult error if there is no inserted vector",
			args: args{
				ctx:  context.Background(),
				vec:  []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
				size: 3,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: defaultCreateFunc,
			want: want{
				err: errors.ErrEmptySearchResult,
			},
		},
		{
			name: "return ErrEmptySearchResult error if the context is canceled",
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					return ctx
				}(),
				vec:  []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
				size: 3,
			},
			fields: fields{
				inMemory:            false,
				idxPath:             idxTempDir(t),
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: defaultCreateFunc,
			want: want{
				err: errors.ErrEmptySearchResult,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			ctx, cancel := context.WithCancel(test.args.ctx)
			defer cancel()

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

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}

			got, err := n.Search(ctx, test.args.vec, test.args.size, test.args.epsilon, test.args.radius)
			if err := checkFunc(test.want, got, n, err); err != nil {
				tt.Errorf("error = %v", err)
			}

			if err := test.afterFunc(tt, n); err != nil {
				tt.Error(err)
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
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct {
		want uint
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(*testing.T, fields) (NGT, error)
		want       want
		checkFunc  func(context.Context, want, uint, NGT, args, error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(ctx context.Context, w want, got uint, n NGT, args args, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if (w.want == 0 && got != 0) || (w.want != 0 && got == 0) {
			return errors.Errorf("want: %d, got: %d", w.want, got)
		}

		if got == 0 || err != nil {
			return nil
		}

		// search before indexing, it should return nothing
		r, err := n.Search(ctx, args.vec, 5, 0, 0)
		if err != nil && !errors.Is(err, errors.ErrEmptySearchResult) {
			return err
		}
		if len(r) != 0 {
			return errors.Errorf("search return before index, result: %#v", r)
		}

		// search after indexing, it should return result
		if err := n.CreateIndex(1); err != nil {
			return err
		}
		r, err = n.Search(ctx, args.vec, 5, 0, 0)
		if err != nil {
			return err
		}
		if len(r) != 1 {
			return errors.Errorf("search return invalid result after index, result: %#v", r)
		}
		if r[0].Distance != 0 {
			return errors.Errorf("distance is not 0")
		}

		return nil
	}
	tests := []test{
		{
			name: "return object id when object type is uint8 and the vector is valid",
			args: args{
				vec: []float32{0, 1, 2, 3, 4, 5, 6, 7, 8},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Uint8,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is uint8 and all vector elem are 0",
			args: args{
				vec: []float32{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Uint8,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is uint8 and all vector elem are min value",
			args: args{
				vec: []float32{
					math.MinInt8, math.MinInt8, math.MinInt8, math.MinInt8,
					math.MinInt8, math.MinInt8, math.MinInt8, math.MinInt8, math.MinInt8,
				},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Uint8,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is uint8 and all vector elem are max value",
			args: args{
				vec: []float32{
					math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
					math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
				},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Uint8,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is float",
			args: args{
				vec: []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Float,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is float and all vector elem are 0",
			args: args{
				vec: []float32{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Float,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is float and all vector elem are min value",
			args: args{
				vec: []float32{
					math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32,
					math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32,
				},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Float,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is float and all vector elem are max value",
			args: args{
				vec: []float32{
					math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32,
					math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32,
				},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Float,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return error if dimension is not the same as insert vector",
			args: args{
				vec: []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  5,
				objectType: Float,
			},
			want: want{
				err: errors.New("incompatible dimension size detected\trequested: 9,\tconfigured: 5"),
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

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			got, err := n.Insert(test.args.vec)
			if err := checkFunc(ctx, test.want, got, n, test.args, err); err != nil {
				tt.Errorf("error = %v", err)
			}

			if err := test.afterFunc(tt, n); err != nil {
				tt.Error(err)
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
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct {
		want uint
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(*testing.T, fields) (NGT, error)
		want       want
		checkFunc  func(context.Context, want, uint, NGT, args, error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(ctx context.Context, w want, got uint, n NGT, args args, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if (w.want == 0 && got != 0) || (w.want != 0 && got == 0) {
			return errors.Errorf("want: %d, got: %d", w.want, got)
		}

		if got == 0 {
			return nil
		}

		r, err := n.Search(ctx, args.vec, 5, 0, 0)
		if err != nil {
			return err
		}
		if len(r) != 1 {
			return errors.Errorf("search return invalid result, result: %#v", r)
		}
		if r[0].Distance != 0 {
			return errors.Errorf("distance is not 0")
		}

		return nil
	}
	tests := []test{ // copied from insert tests
		{
			name: "return object id when object type is uint8",
			args: args{
				vec: []float32{0, 1, 2, 3, 4, 5, 6, 7, 8},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Uint8,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is uint8 and all vector elem are 0",
			args: args{
				vec: []float32{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Uint8,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is uint8 and all vector elem are min value",
			args: args{
				vec: []float32{
					math.MinInt8, math.MinInt8, math.MinInt8, math.MinInt8,
					math.MinInt8, math.MinInt8, math.MinInt8, math.MinInt8, math.MinInt8,
				},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Uint8,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is uint8 and all vector elem are max value",
			args: args{
				vec: []float32{
					math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
					math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8, math.MaxUint8,
				},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Uint8,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is float",
			args: args{
				vec: []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Float,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is float and all vector elem are 0",
			args: args{
				vec: []float32{0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Float,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is float and all vector elem are min value",
			args: args{
				vec: []float32{
					math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32,
					math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32, math.SmallestNonzeroFloat32,
				},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Float,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return object id when object type is float and all vector elem are max value",
			args: args{
				vec: []float32{
					math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32,
					math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32, math.MaxFloat32,
				},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  9,
				objectType: Float,
			},
			want: want{
				want: 1,
			},
		},
		{
			name: "return error if dimension is not the same as insert vector",
			args: args{
				vec: []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
			},
			fields: fields{
				idxPath:    idxTempDir(t),
				inMemory:   false,
				dimension:  5,
				objectType: Float,
			},
			want: want{
				err: errors.New("incompatible dimension size detected\trequested: 9,\tconfigured: 5"),
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

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			got, err := n.InsertCommit(test.args.vec, test.args.poolSize)
			if err := checkFunc(ctx, test.want, got, n, test.args, err); err != nil {
				tt.Errorf("error = %v", err)
			}

			if err := test.afterFunc(tt, n); err != nil {
				tt.Error(err)
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
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct {
		want  []uint
		want1 []error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(*testing.T, fields) (NGT, error)
		want       want
		checkFunc  func(context.Context, want, []uint, NGT, fields, args, []error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(ctx context.Context, w want, got []uint, n NGT, fields fields, args args, got1 []error) error {
		if diff := comparator.Diff(w.want1, got1, comparator.ErrorComparer); diff != "" {
			return errors.New(diff)
		}
		if len(w.want) != len(got) {
			return errors.Errorf("got length not match with want length")
		}

		// check all the vectors can not be get even before indexing
		for _, vid := range got {
			r, err := n.GetVector(vid)
			if err != nil {
				return err
			}
			if len(r) == 0 {
				return errors.Errorf("get object cannot return the result, vec id: %d", r)
			}
		}

		// check all the vectors cannot not be searched before indexing
		for _, vec := range args.vecs {
			if len(vec) != fields.dimension {
				continue
			}
			r, err := n.Search(ctx, vec, 1, 0, 0)
			if err != nil && !errors.Is(err, errors.ErrEmptySearchResult) {
				return err
			}
			if len(r) != 0 {
				return errors.Errorf("search return before create index, result: %#v", r)
			}
		}

		// check all vectors can be searched after indexing
		if err := n.CreateIndex(uint32(len(args.vecs))); err != nil {
			return err
		}
		for _, vec := range args.vecs {
			if len(vec) != fields.dimension {
				continue
			}
			r, err := n.Search(ctx, vec, 1, 0, 0)
			if err != nil {
				return err
			}
			if len(r) != 1 {
				return errors.Errorf("search return invalid result, result: %#v", r)
			}
			if r[0].Distance != 0 {
				return errors.Errorf("vector distance is invalid, got: %d, want: %d", r[0].Distance, 0)
			}
		}
		return nil
	}
	tests := []test{
		// int
		{
			name: "return 1 object id when insert 1 vector (uint8)",
			args: args{
				vecs: [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Uint8,
			},
			want: want{
				want:  []uint{1},
				want1: []error{},
			},
		},
		{
			name: "return 5 object id when insert 5 vectors (uint8)",
			args: args{
				vecs: [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{1, 2, 3, 4, 5, 6, 7, 8, 9},
					{2, 3, 4, 5, 6, 7, 8, 9, 10},
					{3, 4, 5, 6, 7, 8, 9, 10, 11},
					{4, 5, 6, 7, 8, 9, 10, 11, 12},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Uint8,
			},
			want: want{
				want:  []uint{1, 2, 3, 4, 5},
				want1: []error{},
			},
		},
		{
			name: "return 2 object id when insert 2 same vectors (uint8)",
			args: args{
				vecs: [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Uint8,
			},
			want: want{
				want:  []uint{1, 2},
				want1: []error{},
			},
		},
		{
			name: "return 2 object id and 2 errors when insert 2 vectors with same dimension and 2 with invalid dimension (uint8)",
			args: args{
				vecs: [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{0, 1, 2, 3, 4, 5, 6, 7, 9},
					{0, 1, 2, 3, 4, 5, 6, 7},
					{0, 1, 2, 3, 4, 5, 6, 7, 8, 10},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Uint8,
			},
			want: want{
				want: []uint{1, 2},
				want1: []error{
					errors.New("bulkinsert error detected index number: 2,	id: 0: incompatible dimension size detected	requested: 8,	configured: 9"),
					errors.New("bulkinsert error detected index number: 3,	id: 0: incompatible dimension size detected	requested: 10,	configured: 9"),
				},
			},
		},
		// float
		{
			name: "return 1 object id when insert 1 vector (float)",
			args: args{
				vecs: [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
			want: want{
				want:  []uint{1},
				want1: []error{},
			},
		},
		{
			name: "return 5 object id when insert 5 vectors (float)",
			args: args{
				vecs: [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
					{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
					{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10},
					{0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11},
					{0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
			want: want{
				want:  []uint{1, 2, 3, 4, 5},
				want1: []error{},
			},
		},
		{
			name: "return 2 object id when insert 2 same vectors (float)",
			args: args{
				vecs: [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
			want: want{
				want:  []uint{1, 2},
				want1: []error{},
			},
		},
		{
			name: "return 2 object id and 2 errors when insert 2 vectors with same dimension and 2 with invalid dimension (float)",
			args: args{
				vecs: [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.9},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.1},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
			want: want{
				want: []uint{1, 2},
				want1: []error{
					errors.New("bulkinsert error detected index number: 2,	id: 0: incompatible dimension size detected	requested: 8,	configured: 9"),
					errors.New("bulkinsert error detected index number: 3,	id: 0: incompatible dimension size detected	requested: 10,	configured: 9"),
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			if test.createFunc == nil {
				test.createFunc = defaultCreateFunc
			}

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			defer func() {
				if err := test.afterFunc(tt, n); err != nil {
					tt.Error(err)
				}
			}()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			got, got1 := n.BulkInsert(test.args.vecs)
			if err := checkFunc(ctx, test.want, got, n, test.fields, test.args, got1); err != nil {
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
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct {
		want  []uint
		want1 []error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(*testing.T, fields) (NGT, error)
		want       want
		checkFunc  func(context.Context, want, []uint, NGT, fields, args, []error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(ctx context.Context, w want, got []uint, n NGT, fields fields, args args, got1 []error) error {
		if diff := comparator.Diff(w.want1, got1, comparator.ErrorComparer); diff != "" {
			return errors.New(diff)
		}
		if len(w.want) != len(got) || len(w.want) != len(args.vecs)-len(got1) {
			return errors.Errorf("got length not match with want length")
		}

		for _, vid := range got {
			r, err := n.GetVector(vid)
			if err != nil {
				return err
			}
			if len(r) == 0 {
				return errors.Errorf("get object cannot return the result, vec id: %d", r)
			}
		}

		for _, vec := range args.vecs {
			if len(vec) != fields.dimension {
				continue
			}
			r, err := n.Search(ctx, vec, 1, 0, 0)
			if err != nil {
				return err
			}
			if len(r) != 1 {
				return errors.Errorf("search return invalid result, result: %#v", r)
			}
			if r[0].Distance != 0 {
				return errors.Errorf("vector distance is invalid, got: %d, want: %d", r[0].Distance, 0)
			}
		}
		return nil
	}
	tests := []test{ // copy from bulk insert test
		// int
		{
			name: "return 1 object id when insert 1 vector (uint8)",
			args: args{
				vecs: [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Uint8,
			},
			want: want{
				want:  []uint{1},
				want1: []error{},
			},
		},
		{
			name: "return 5 object id when insert 5 vectors (uint8)",
			args: args{
				vecs: [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{1, 2, 3, 4, 5, 6, 7, 8, 9},
					{2, 3, 4, 5, 6, 7, 8, 9, 10},
					{3, 4, 5, 6, 7, 8, 9, 10, 11},
					{4, 5, 6, 7, 8, 9, 10, 11, 12},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Uint8,
			},
			want: want{
				want:  []uint{1, 2, 3, 4, 5},
				want1: []error{},
			},
		},
		{
			name: "return 2 object id when insert 2 same vectors (uint8)",
			args: args{
				vecs: [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Uint8,
			},
			want: want{
				want:  []uint{1, 2},
				want1: []error{},
			},
		},
		{
			name: "return 2 object id and 2 errors when insert 2 vectors with same dimension and 2 with invalid dimension (uint8)",
			args: args{
				vecs: [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{0, 1, 2, 3, 4, 5, 6, 7, 9},
					{0, 1, 2, 3, 4, 5, 6, 7},
					{0, 1, 2, 3, 4, 5, 6, 7, 8, 10},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Uint8,
			},
			want: want{
				want: []uint{1, 2},
				want1: []error{
					errors.New("bulkinsert error detected index number: 2,	id: 0: incompatible dimension size detected	requested: 8,	configured: 9"),
					errors.New("bulkinsert error detected index number: 3,	id: 0: incompatible dimension size detected	requested: 10,	configured: 9"),
				},
			},
		},
		// float
		{
			name: "return 1 object id when insert 1 vector (float)",
			args: args{
				vecs: [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
			want: want{
				want:  []uint{1},
				want1: []error{},
			},
		},
		{
			name: "return 5 object id when insert 5 vectors (float)",
			args: args{
				vecs: [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
					{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},
					{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10},
					{0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11},
					{0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
			want: want{
				want:  []uint{1, 2, 3, 4, 5},
				want1: []error{},
			},
		},
		{
			name: "return 2 object id when insert 2 same vectors (float)",
			args: args{
				vecs: [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
			want: want{
				want:  []uint{1, 2},
				want1: []error{},
			},
		},
		{
			name: "return 2 object id and 2 errors when insert 2 vectors with same dimension and 2 with invalid dimension (float)",
			args: args{
				vecs: [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.9},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.10},
				},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
			want: want{
				want: []uint{1, 2},
				want1: []error{
					errors.New("bulkinsert error detected index number: 2,	id: 0: incompatible dimension size detected	requested: 8,	configured: 9"),
					errors.New("bulkinsert error detected index number: 3,	id: 0: incompatible dimension size detected	requested: 10,	configured: 9"),
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			if test.createFunc == nil {
				test.createFunc = defaultCreateFunc
			}

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			defer func() {
				if err := test.afterFunc(tt, n); err != nil {
					tt.Error(err)
				}
			}()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			got, got1 := n.BulkInsertCommit(test.args.vecs, test.args.poolSize)
			if err := checkFunc(ctx, test.want, got, n, test.fields, test.args, got1); err != nil {
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
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(*testing.T, fields) (NGT, error)
		want       want
		checkFunc  func(context.Context, want, NGT, args, error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(_ context.Context, w want, n NGT, args args, got error) error {
		if diff := comparator.Diff(w.err, got); diff != "" {
			return errors.New(diff)
		}
		if ngt, ok := n.(*ngt); ok {
			_, err := os.Stat(ngt.idxPath)
			// if ngt is in-memory mode, the index file should not be created
			if ngt.inMemory {
				if !errors.Is(err, fs.ErrNotExist) {
					return errors.Errorf("NGT index file created, err: %s", err)
				}
			} else { // if ngt is not in-memory mode, the file should be created
				if err != nil {
					return errors.Errorf("NGT index file error, err: %s", err)
				}
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "return nil when no index is required and pool size is 0",
			args: args{
				poolSize: 0,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
		},
		{
			name: "return nil when no index is required and pool size > 0",
			args: args{
				poolSize: 100,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
		},
		func() test {
			ivs := [][]float32{
				{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},           // vec id 1
				{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},         // vec id 2
				{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10},        // vec id 3
				{0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11},       // vec id 4
				{0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12},      // vec id 5
				{0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13},     // vec id 6
				{0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14},    // vec id 7
				{0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15},   // vec id 8
				{0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16},  // vec id 9
				{0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17}, // vec id 10
			}

			return test{
				name: "return nil when index is required and pool size = 0",
				args: args{
					poolSize: 0,
				},
				fields: fields{
					idxPath:             idxTempDir(t),
					inMemory:            false,
					bulkInsertChunkSize: 100,
					dimension:           9,
					objectType:          Float,
				},
				createFunc: func(t *testing.T, f fields) (NGT, error) {
					t.Helper()

					ngt, err := defaultCreateFunc(t, f)
					if err != nil {
						return nil, err
					}

					if _, err := ngt.BulkInsert(ivs); len(err) != 0 {
						t.Error(err)
						return nil, err[0]
					}

					return ngt, err
				},
				checkFunc: func(ctx context.Context, w want, n NGT, a args, e error) error {
					if err := defaultCheckFunc(ctx, w, n, a, e); err != nil {
						return err
					}

					// search the inserted vector exists after create index
					for _, v := range ivs {
						if rs, err := n.Search(ctx, v, 1, 0, 0); err != nil {
							if rs[0].Distance != 0 {
								return errors.Errorf("vector distance is invalid, got: %d, want: %d", rs[0].Distance, 0)
							}
						}
					}

					return nil
				},
			}
		}(),
		func() test {
			ivs := [][]float32{
				{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},           // vec id 1
				{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},         // vec id 2
				{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10},        // vec id 3
				{0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11},       // vec id 4
				{0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12},      // vec id 5
				{0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13},     // vec id 6
				{0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14},    // vec id 7
				{0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15},   // vec id 8
				{0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16},  // vec id 9
				{0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17}, // vec id 10
			}

			return test{
				name: "return nil when index is required and pool size > 0",
				args: args{
					poolSize: 100,
				},
				fields: fields{
					idxPath:             idxTempDir(t),
					inMemory:            false,
					bulkInsertChunkSize: 5,
					dimension:           9,
					objectType:          Float,
				},
				createFunc: func(t *testing.T, f fields) (NGT, error) {
					t.Helper()

					ngt, err := defaultCreateFunc(t, f)
					if err != nil {
						return nil, err
					}

					if _, err := ngt.BulkInsert(ivs); len(err) != 0 {
						t.Error(err)
						return nil, err[0]
					}

					return ngt, err
				},
				checkFunc: func(ctx context.Context, w want, n NGT, a args, e error) error {
					if err := defaultCheckFunc(ctx, w, n, a, e); err != nil {
						return err
					}

					// search the inserted vector exists after create index
					for _, v := range ivs {
						if rs, err := n.Search(ctx, v, 1, 0, 0); err != nil {
							if rs[0].Distance != 0 {
								return errors.Errorf("vector distance is invalid, got: %d, want: %d", rs[0].Distance, 0)
							}
						}
					}

					return nil
				},
			}
		}(),
		{
			name: "return nil when no index is required and in memory mode",
			args: args{
				poolSize: 100,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            true,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
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

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			defer func() {
				if err := test.afterFunc(tt, n); err != nil {
					tt.Error(err)
				}
			}()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err = n.CreateAndSaveIndex(test.args.poolSize)
			if err := checkFunc(ctx, test.want, n, test.args, err); err != nil {
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
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(*testing.T, fields) (NGT, error)
		want       want
		checkFunc  func(context.Context, want, NGT, args, error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(_ context.Context, w want, n NGT, args args, got error) error {
		if diff := comparator.Diff(w.err, got); diff != "" {
			return errors.New(diff)
		}
		if ngt, ok := n.(*ngt); ok {
			_, err := os.Stat(ngt.idxPath)
			// if ngt is in-memory mode, the index file should not be created
			if ngt.inMemory {
				if !errors.Is(err, fs.ErrNotExist) {
					return errors.Errorf("NGT index file created, err: %s", err)
				}
			} else { // if ngt is not in-memory mode, the file should be created
				if err != nil {
					return errors.Errorf("NGT index file error, err: %s", err)
				}
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "return nil when no index is required and pool size is 0",
			args: args{
				poolSize: 0,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
		},
		{
			name: "return nil when no index is required and pool size > 0",
			args: args{
				poolSize: 100,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
		},
		func() test {
			ivs := [][]float32{
				{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},           // vec id 1
				{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},         // vec id 2
				{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10},        // vec id 3
				{0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11},       // vec id 4
				{0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12},      // vec id 5
				{0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13},     // vec id 6
				{0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14},    // vec id 7
				{0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15},   // vec id 8
				{0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16},  // vec id 9
				{0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17}, // vec id 10
			}

			return test{
				name: "return nil when index is required and pool size = 0",
				args: args{
					poolSize: 0,
				},
				fields: fields{
					idxPath:             idxTempDir(t),
					inMemory:            false,
					bulkInsertChunkSize: 100,
					dimension:           9,
					objectType:          Float,
				},
				createFunc: func(t *testing.T, f fields) (NGT, error) {
					t.Helper()

					ngt, err := defaultCreateFunc(t, f)
					if err != nil {
						return nil, err
					}

					if _, err := ngt.BulkInsert(ivs); len(err) != 0 {
						t.Error(err)
						return nil, err[0]
					}

					return ngt, err
				},
				checkFunc: func(ctx context.Context, w want, n NGT, a args, e error) error {
					if err := defaultCheckFunc(ctx, w, n, a, e); err != nil {
						return err
					}

					// search the inserted vector exists after create index
					for _, v := range ivs {
						if rs, err := n.Search(ctx, v, 1, 0, 0); err != nil {
							if rs[0].Distance != 0 {
								return errors.Errorf("vector distance is invalid, got: %d, want: %d", rs[0].Distance, 0)
							}
						}
					}

					return nil
				},
			}
		}(),
		func() test {
			ivs := [][]float32{
				{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},           // vec id 1
				{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},         // vec id 2
				{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10},        // vec id 3
				{0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11},       // vec id 4
				{0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12},      // vec id 5
				{0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13},     // vec id 6
				{0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14},    // vec id 7
				{0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15},   // vec id 8
				{0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16},  // vec id 9
				{0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17}, // vec id 10
			}

			return test{
				name: "return nil when index is required and pool size > 0",
				args: args{
					poolSize: 100,
				},
				fields: fields{
					idxPath:             idxTempDir(t),
					inMemory:            false,
					bulkInsertChunkSize: 5,
					dimension:           9,
					objectType:          Float,
				},
				createFunc: func(t *testing.T, f fields) (NGT, error) {
					t.Helper()

					ngt, err := defaultCreateFunc(t, f)
					if err != nil {
						return nil, err
					}

					if _, err := ngt.BulkInsert(ivs); len(err) != 0 {
						t.Error(err)
						return nil, err[0]
					}

					return ngt, err
				},
				checkFunc: func(ctx context.Context, w want, n NGT, a args, e error) error {
					if err := defaultCheckFunc(ctx, w, n, a, e); err != nil {
						return err
					}

					// search the inserted vector exists after create index
					for _, v := range ivs {
						if rs, err := n.Search(ctx, v, 1, 0, 0); err != nil {
							if rs[0].Distance != 0 {
								return errors.Errorf("vector distance is invalid, got: %d, want: %d", rs[0].Distance, 0)
							}
						}
					}

					return nil
				},
			}
		}(),
		{
			name: "return nil when no index is required and in memory mode",
			args: args{
				poolSize: 100,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            true,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
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

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			defer func() {
				if err := test.afterFunc(tt, n); err != nil {
					tt.Error(err)
				}
			}()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err = n.CreateIndex(test.args.poolSize)
			if err := checkFunc(ctx, test.want, n, test.args, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_ngt_SaveIndex(t *testing.T) {
	type args struct {
		poolSize uint32
	}
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(*testing.T, fields) (NGT, error)
		want       want
		checkFunc  func(context.Context, want, NGT, args, error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(_ context.Context, w want, n NGT, args args, e error) error {
		if ngt, ok := n.(*ngt); ok {
			_, err := os.Stat(ngt.idxPath)
			// if ngt is in-memory mode, the index file should not be created
			if ngt.inMemory {
				if !errors.Is(err, fs.ErrNotExist) {
					return errors.Errorf("NGT index file created, err: %s", err)
				}
			} else { // if ngt is not in-memory mode, the file should be created
				if err != nil {
					return errors.Errorf("NGT index file error, err: %s", err)
				}
			}
		}
		return nil
	}
	tests := []test{
		{
			name: "return nil when no index is required and pool size is 0",
			args: args{
				poolSize: 0,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
		},
		{
			name: "return nil when no index is required and pool size > 0",
			args: args{
				poolSize: 100,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
			},
		},
		func() test {
			ivs := [][]float32{
				{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},           // vec id 1
				{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},         // vec id 2
				{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10},        // vec id 3
				{0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11},       // vec id 4
				{0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12},      // vec id 5
				{0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13},     // vec id 6
				{0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14},    // vec id 7
				{0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15},   // vec id 8
				{0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16},  // vec id 9
				{0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17}, // vec id 10
			}

			return test{
				name: "return nil when index is required and pool size = 0",
				args: args{
					poolSize: 0,
				},
				fields: fields{
					idxPath:             idxTempDir(t),
					inMemory:            false,
					bulkInsertChunkSize: 100,
					dimension:           9,
					objectType:          Float,
				},
				createFunc: func(t *testing.T, f fields) (NGT, error) {
					t.Helper()

					ngt, err := defaultCreateFunc(t, f)
					if err != nil {
						return nil, err
					}

					if _, err := ngt.BulkInsert(ivs); len(err) != 0 {
						t.Error(err)
						return nil, err[0]
					}

					return ngt, err
				},
				checkFunc: func(ctx context.Context, w want, n NGT, a args, e error) error {
					if err := defaultCheckFunc(ctx, w, n, a, e); err != nil {
						return err
					}

					// search the inserted vector exists after create index
					for _, v := range ivs {
						if rs, err := n.Search(ctx, v, 1, 0, 0); err != nil && !errors.Is(err, errors.ErrEmptySearchResult) {
							if rs[0].Distance != 0 {
								return errors.Errorf("vector distance is invalid, got: %d, want: %d", rs[0].Distance, 0)
							}
						}
					}

					return nil
				},
			}
		}(),
		func() test {
			ivs := [][]float32{
				{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},           // vec id 1
				{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9},         // vec id 2
				{0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10},        // vec id 3
				{0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11},       // vec id 4
				{0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12},      // vec id 5
				{0.5, 0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13},     // vec id 6
				{0.6, 0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14},    // vec id 7
				{0.7, 0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15},   // vec id 8
				{0.8, 0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16},  // vec id 9
				{0.9, 0.10, 0.11, 0.12, 0.13, 0.14, 0.15, 0.16, 0.17}, // vec id 10
			}

			return test{
				name: "return nil when index is required and pool size > 0",
				args: args{
					poolSize: 100,
				},
				fields: fields{
					idxPath:             idxTempDir(t),
					inMemory:            false,
					bulkInsertChunkSize: 5,
					dimension:           9,
					objectType:          Float,
				},
				createFunc: func(t *testing.T, f fields) (NGT, error) {
					t.Helper()

					ngt, err := defaultCreateFunc(t, f)
					if err != nil {
						return nil, err
					}

					if _, err := ngt.BulkInsert(ivs); len(err) != 0 {
						t.Error(err)
						return nil, err[0]
					}

					return ngt, err
				},
				checkFunc: func(ctx context.Context, w want, n NGT, a args, e error) error {
					if err := defaultCheckFunc(ctx, w, n, a, e); err != nil {
						return err
					}

					// search the inserted vector exists after create index
					for _, v := range ivs {
						if rs, err := n.Search(ctx, v, 1, 0, 0); err != nil && !errors.Is(err, errors.ErrEmptySearchResult) {
							if rs[0].Distance != 0 {
								return errors.Errorf("vector distance is invalid, got: %d, want: %d", rs[0].Distance, 0)
							}
						}
					}

					return nil
				},
			}
		}(),
		{
			name: "return nil when no index is required and in memory mode",
			args: args{
				poolSize: 100,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            true,
				bulkInsertChunkSize: 5,
				dimension:           9,
				objectType:          Float,
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

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			defer func() {
				if err := test.afterFunc(tt, n); err != nil {
					tt.Error(err)
				}
			}()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err = n.SaveIndex()
			if err := checkFunc(ctx, test.want, n, test.args, err); err != nil {
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
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(t *testing.T, fields fields) (NGT, error)
		want       want
		checkFunc  func(want, NGT, error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(w want, n NGT, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		return nil
	}
	insertCreateFunc := func(t *testing.T, fields fields, vecs [][]float32, poolSize uint32) (NGT, error) { // create func with insert/index
		t.Helper()

		ngt, err := defaultCreateFunc(t, fields)
		if err != nil {
			return nil, err
		}

		if _, err := ngt.BulkInsertCommit(vecs, poolSize); len(err) != 0 {
			t.Error(err)
			return nil, err[0]
		}

		return ngt, nil
	}
	tests := []test{
		// int
		{
			name: "remove success when id exists (int)",
			args: args{
				id: 1,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			checkFunc: func(w want, n NGT, e error) error {
				if err := defaultCheckFunc(w, n, e); err != nil {
					return err
				}

				if v, err := n.GetVector(1); err == nil || len(v) > 0 {
					return errors.Errorf("vector removed but returned, vec: %s, err: %s", v, err)
				}

				return nil
			},
		},
		{
			name: "return error when id do not exists (int)",
			args: args{
				id: 999,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			checkFunc: func(w want, n NGT, e error) error {
				if e == nil {
					return errors.New("no error returned")
				}

				// ensure the inserted vector exists
				if v, err := n.GetVector(1); err != nil || len(v) == 0 {
					return errors.Errorf("vector do not return, vec: %s, err: %s", v, err)
				}

				return nil
			},
		},
		// float
		{
			name: "remove success when id exists (float)",
			args: args{
				id: 1,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			checkFunc: func(w want, n NGT, e error) error {
				if err := defaultCheckFunc(w, n, e); err != nil {
					return err
				}

				if v, err := n.GetVector(1); err == nil || len(v) > 0 {
					return errors.Errorf("vector removed but returned, vec: %s, err: %s", v, err)
				}

				return nil
			},
		},
		{
			name: "return error when id do not exists (float)",
			args: args{
				id: 999,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			checkFunc: func(w want, n NGT, e error) error {
				if e == nil {
					return errors.New("no error returned")
				}

				// ensure the inserted vector exists
				if v, err := n.GetVector(1); err != nil || len(v) == 0 {
					return errors.Errorf("vector do not return, vec: %s, err: %s", v, err)
				}

				return nil
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

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}

			err = n.Remove(test.args.id)
			if err := checkFunc(test.want, n, err); err != nil {
				tt.Errorf("error = %v", err)
			}

			if err := test.afterFunc(tt, n); err != nil {
				tt.Error(err)
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
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct {
		err error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(t *testing.T, fields fields) (NGT, error)
		want       want
		checkFunc  func(want, NGT, error) error
		beforeFunc func(args)
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(w want, n NGT, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		return nil
	}
	insertCreateFunc := func(t *testing.T, fields fields, vecs [][]float32, poolSize uint32) (NGT, error) { // create func with insert/index
		t.Helper()

		ngt, err := defaultCreateFunc(t, fields)
		if err != nil {
			return nil, err
		}

		if _, err := ngt.BulkInsertCommit(vecs, poolSize); len(err) != 0 {
			t.Error(err)
			return nil, err[0]
		}

		return ngt, nil
	}
	tests := []test{
		// int
		{
			name: "remove success when id exists (int)",
			args: args{
				ids: []uint{1},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			checkFunc: func(w want, n NGT, e error) error {
				if err := defaultCheckFunc(w, n, e); err != nil {
					return err
				}

				if v, err := n.GetVector(1); err == nil || len(v) > 0 {
					return errors.Errorf("vector removed but returned, vec: %s, err: %s", v, err)
				}

				return nil
			},
		},
		{
			name: "return error when id do not exists (int)",
			args: args{
				ids: []uint{999},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			checkFunc: func(w want, n NGT, e error) error {
				if e == nil { // the error contains random character, so check with nil
					return errors.New("error should be returned, but nil returned")
				}

				// ensure the inserted vector exists
				if v, err := n.GetVector(1); err != nil || len(v) == 0 {
					return errors.Errorf("vector do not return, vec: %s, err: %s", v, err)
				}

				return nil
			},
		},
		{
			name: "return error when some id do not exists (int)",
			args: args{
				ids: []uint{1, 999},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vecs := [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},  // vec id 1
					{0, 1, 2, 3, 4, 5, 6, 7, 9},  // vec id 2
					{0, 1, 2, 3, 4, 5, 6, 7, 10}, // vec id 3
				}

				return insertCreateFunc(t, fields, vecs, 1)
			},
			checkFunc: func(w want, n NGT, e error) error {
				if e == nil { // the error contains random character, so check with nil
					return errors.New("error should be returned, but nil returned")
				}

				// ensure the vector is deleted
				if v, err := n.GetVector(1); err == nil || len(v) > 0 {
					return errors.Errorf("vector removed but returned, vec: %s, err: %s", v, err)
				}

				// ensure the inserted vector exists
				if v, err := n.GetVector(2); err != nil || len(v) == 0 {
					return errors.Errorf("vector do not return, vec: %s, err: %s", v, err)
				}
				if v, err := n.GetVector(3); err != nil || len(v) == 0 {
					return errors.Errorf("vector do not return, vec: %s, err: %s", v, err)
				}

				return nil
			},
		},
		// float
		{
			name: "remove success when id exists (float)",
			args: args{
				ids: []uint{1},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			checkFunc: func(w want, n NGT, e error) error {
				if err := defaultCheckFunc(w, n, e); err != nil {
					return err
				}

				if v, err := n.GetVector(1); err == nil || len(v) > 0 {
					return errors.Errorf("vector removed but returned, vec: %s, err: %s", v, err)
				}

				return nil
			},
		},
		{
			name: "return error when id do not exists (float)",
			args: args{
				ids: []uint{999},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			checkFunc: func(w want, n NGT, e error) error {
				if e == nil { // the error contains random character, so check with nil
					return errors.New("no error returned")
				}

				// ensure the inserted vector exists
				if v, err := n.GetVector(1); err != nil || len(v) == 0 {
					return errors.Errorf("vector do not return, vec: %s, err: %s", v, err)
				}

				return nil
			},
		},
		{
			name: "return error when some id do not exists (float)",
			args: args{
				ids: []uint{1, 999},
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vecs := [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},  // vec id 1
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.9},  // vec id 2
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.10}, // vec id 3
				}

				return insertCreateFunc(t, fields, vecs, 1)
			},
			checkFunc: func(w want, n NGT, e error) error {
				if e == nil { // the error contains random character, so check with nil
					return errors.New("error should be returned, but nil returned")
				}

				// ensure the vector is deleted
				if v, err := n.GetVector(1); err == nil || len(v) > 0 {
					return errors.Errorf("vector removed but returned, vec: %s, err: %s", v, err)
				}

				// ensure the inserted vector exists
				if v, err := n.GetVector(2); err != nil || len(v) == 0 {
					return errors.Errorf("vector do not return, vec: %s, err: %s", v, err)
				}
				if v, err := n.GetVector(3); err != nil || len(v) == 0 {
					return errors.Errorf("vector do not return, vec: %s, err: %s", v, err)
				}

				return nil
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

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}

			err = n.BulkRemove(test.args.ids...)
			if err := checkFunc(test.want, n, err); err != nil {
				tt.Errorf("error = %v", err)
			}

			if err := test.afterFunc(tt, n); err != nil {
				tt.Error(err)
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
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct {
		want []float32
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		createFunc func(t *testing.T, fields fields) (NGT, error)
		want       want
		checkFunc  func(w want, got []float32, n NGT, err error) error
		beforeFunc func(args)
		afterFunc  func(t *testing.T, n NGT) error
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	defaultCheckFunc := func(w want, got []float32, n NGT, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	insertCreateFunc := func(t *testing.T, fields fields, vecs [][]float32, poolSize uint32) (NGT, error) { // create func with insert/index
		t.Helper()

		ngt, err := defaultCreateFunc(t, fields)
		if err != nil {
			return nil, err
		}

		if _, err := ngt.BulkInsertCommit(vecs, poolSize); len(err) != 0 {
			t.Error(err)
			return nil, err[0]
		}

		return ngt, nil
	}
	tests := []test{
		// int
		{
			name: "return vector when id exists (int)",
			args: args{
				id: 1,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vecs := [][]float32{
					{0, 1, 2, 3, 4, 5, 6, 7, 8},
					{0, 1, 2, 3, 4, 5, 6, 7, 9},
					{0, 1, 2, 3, 4, 5, 6, 7, 10},
					{0, 1, 2, 3, 4, 5, 6, 7, 11},
					{0, 1, 2, 3, 4, 5, 6, 7, 12},
				}

				return insertCreateFunc(t, fields, vecs, 1)
			},
			want: want{
				want: []float32{0, 1, 2, 3, 4, 5, 6, 7, 8},
			},
		},
		{
			name: "return error when id do not exists (int)",
			args: args{
				id: 10,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Uint8,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 1, 2, 3, 4, 5, 6, 7, 8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			checkFunc: func(w want, got []float32, n NGT, err error) error {
				if err == nil {
					return errors.New("no error return when vector not exists")
				}
				if !reflect.DeepEqual(got, w.want) {
					return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
				}
				return nil
			},
		},
		// float
		{
			name: "return float vector when id exists",
			args: args{
				id: 1,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vecs := [][]float32{
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.9},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.10},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.11},
					{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.12},
				}

				return insertCreateFunc(t, fields, vecs, 1)
			},
			want: want{
				want: []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
			},
		},
		{
			name: "return error when id do not exists (float)",
			args: args{
				id: 10,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()
				vec := []float32{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8}

				return insertCreateFunc(t, fields, [][]float32{vec}, 1)
			},
			checkFunc: func(w want, got []float32, n NGT, err error) error {
				if err == nil {
					return errors.New("no error return when vector not exists")
				}
				if !reflect.DeepEqual(got, w.want) {
					return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
				}
				return nil
			},
		},
		// other
		{
			name: "return error when object type is invalid",
			args: args{
				id: 10,
			},
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          999,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			createFunc: func(t *testing.T, fields fields) (NGT, error) {
				t.Helper()

				n := &ngt{
					idxPath:             fields.idxPath,
					inMemory:            fields.inMemory,
					bulkInsertChunkSize: fields.bulkInsertChunkSize,
					objectType:          fields.objectType,
					radius:              fields.radius,
					epsilon:             fields.epsilon,
					poolSize:            fields.poolSize,
				}
				return n, nil
			},
			want: want{
				err: errors.ErrUnsupportedObjectType,
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

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			defer func() {
				if err := test.afterFunc(tt, n); err != nil {
					tt.Error(err)
				}
			}()

			got, err := n.GetVector(test.args.id)
			if err := checkFunc(test.want, got, n, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

// Skip to test this function because of it contains C dependencies in the argument list,
// and we cannot test it without importing C dependencies, but gotest does not support it.
// Keep this test function to avoid generating from gotests command
func Test_ngt_newGoError(t *testing.T) {
	t.SkipNow()
}

func Test_ngt_Close(t *testing.T) {
	type fields struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		dimension           int
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
	}
	type want struct{}
	type test struct {
		name       string
		fields     fields
		want       want
		createFunc func(t *testing.T, fields fields) (NGT, error)
		checkFunc  func(want) error
		beforeFunc func()
		afterFunc  func(*testing.T, NGT) error
	}
	defaultCheckFunc := func(w want) error {
		return nil
	}
	defaultCreateFunc := func(t *testing.T, fields fields) (NGT, error) {
		t.Helper()

		return New(
			WithInMemoryMode(fields.inMemory),
			WithIndexPath(fields.idxPath),
			WithBulkInsertChunkSize(fields.bulkInsertChunkSize),
			WithObjectType(fields.objectType),
			WithDefaultRadius(fields.radius),
			WithDefaultEpsilon(fields.epsilon),
			WithDefaultPoolSize(fields.poolSize),
			WithDimension(fields.dimension),
		)
	}
	tests := []test{
		{
			name: "close success",
			fields: fields{
				idxPath:             idxTempDir(t),
				inMemory:            false,
				bulkInsertChunkSize: 100,
				dimension:           9,
				objectType:          Float,
				radius:              float32(-1.0),
				epsilon:             float32(0.1),
			},
			want: want{},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc()
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

			n, err := test.createFunc(tt, test.fields)
			if err != nil {
				tt.Fatal(err)
			}
			defer func() {
				if err := test.afterFunc(tt, n); err != nil {
					tt.Error(err)
				}
			}()

			n.Close()
			if err := checkFunc(test.want); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
