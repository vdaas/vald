//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/compress"
	"github.com/vdaas/vald/internal/errgroup"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/worker"
	"go.uber.org/goleak"
)

func TestNewCompressor(t *testing.T) {
	type args struct {
		opts []CompressorOption
	}
	type want struct {
		want Compressor
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Compressor, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Compressor, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           opts: nil,
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
		           opts: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}

			got, err := NewCompressor(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_compressor_PreStart(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		algorithm        string
		compressionLevel int
		compressor       compress.Compressor
		worker           worker.Worker
		workerOpts       []worker.WorkerOption
		eg               errgroup.Group
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		       },
		       fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
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
		           ctx: nil,
		           },
		           fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &compressor{
				algorithm:        test.fields.algorithm,
				compressionLevel: test.fields.compressionLevel,
				compressor:       test.fields.compressor,
				worker:           test.fields.worker,
				workerOpts:       test.fields.workerOpts,
				eg:               test.fields.eg,
			}

			err := c.PreStart(test.args.ctx)
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_compressor_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type fields struct {
		algorithm        string
		compressionLevel int
		compressor       compress.Compressor
		worker           worker.Worker
		workerOpts       []worker.WorkerOption
		eg               errgroup.Group
	}
	type want struct {
		want <-chan error
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, <-chan error, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got <-chan error, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		       },
		       fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
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
		           ctx: nil,
		           },
		           fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &compressor{
				algorithm:        test.fields.algorithm,
				compressionLevel: test.fields.compressionLevel,
				compressor:       test.fields.compressor,
				worker:           test.fields.worker,
				workerOpts:       test.fields.workerOpts,
				eg:               test.fields.eg,
			}

			got, err := c.Start(test.args.ctx)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_compressor_dispatchCompress(t *testing.T) {
	type args struct {
		ctx     context.Context
		vectors [][]float32
	}
	type fields struct {
		algorithm        string
		compressionLevel int
		compressor       compress.Compressor
		worker           worker.Worker
		workerOpts       []worker.WorkerOption
		eg               errgroup.Group
	}
	type want struct {
		wantResults [][]byte
		err         error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, [][]byte, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotResults [][]byte, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotResults, w.wantResults) {
			return errors.Errorf("got = %v, want %v", gotResults, w.wantResults)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           vectors: nil,
		       },
		       fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
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
		           ctx: nil,
		           vectors: nil,
		           },
		           fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &compressor{
				algorithm:        test.fields.algorithm,
				compressionLevel: test.fields.compressionLevel,
				compressor:       test.fields.compressor,
				worker:           test.fields.worker,
				workerOpts:       test.fields.workerOpts,
				eg:               test.fields.eg,
			}

			gotResults, err := c.dispatchCompress(test.args.ctx, test.args.vectors...)
			if err := test.checkFunc(test.want, gotResults, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_compressor_dispatchDecompress(t *testing.T) {
	type args struct {
		ctx    context.Context
		bytess [][]byte
	}
	type fields struct {
		algorithm        string
		compressionLevel int
		compressor       compress.Compressor
		worker           worker.Worker
		workerOpts       []worker.WorkerOption
		eg               errgroup.Group
	}
	type want struct {
		wantResults [][]float32
		err         error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, [][]float32, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotResults [][]float32, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotResults, w.wantResults) {
			return errors.Errorf("got = %v, want %v", gotResults, w.wantResults)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           bytess: nil,
		       },
		       fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
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
		           ctx: nil,
		           bytess: nil,
		           },
		           fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &compressor{
				algorithm:        test.fields.algorithm,
				compressionLevel: test.fields.compressionLevel,
				compressor:       test.fields.compressor,
				worker:           test.fields.worker,
				workerOpts:       test.fields.workerOpts,
				eg:               test.fields.eg,
			}

			gotResults, err := c.dispatchDecompress(test.args.ctx, test.args.bytess...)
			if err := test.checkFunc(test.want, gotResults, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_compressor_Compress(t *testing.T) {
	type args struct {
		ctx    context.Context
		vector []float32
	}
	type fields struct {
		algorithm        string
		compressionLevel int
		compressor       compress.Compressor
		worker           worker.Worker
		workerOpts       []worker.WorkerOption
		eg               errgroup.Group
	}
	type want struct {
		want []byte
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []byte, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []byte, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           vector: nil,
		       },
		       fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
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
		           ctx: nil,
		           vector: nil,
		           },
		           fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &compressor{
				algorithm:        test.fields.algorithm,
				compressionLevel: test.fields.compressionLevel,
				compressor:       test.fields.compressor,
				worker:           test.fields.worker,
				workerOpts:       test.fields.workerOpts,
				eg:               test.fields.eg,
			}

			got, err := c.Compress(test.args.ctx, test.args.vector)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_compressor_Decompress(t *testing.T) {
	type args struct {
		ctx   context.Context
		bytes []byte
	}
	type fields struct {
		algorithm        string
		compressionLevel int
		compressor       compress.Compressor
		worker           worker.Worker
		workerOpts       []worker.WorkerOption
		eg               errgroup.Group
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           bytes: nil,
		       },
		       fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
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
		           ctx: nil,
		           bytes: nil,
		           },
		           fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &compressor{
				algorithm:        test.fields.algorithm,
				compressionLevel: test.fields.compressionLevel,
				compressor:       test.fields.compressor,
				worker:           test.fields.worker,
				workerOpts:       test.fields.workerOpts,
				eg:               test.fields.eg,
			}

			got, err := c.Decompress(test.args.ctx, test.args.bytes)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_compressor_MultiCompress(t *testing.T) {
	type args struct {
		ctx     context.Context
		vectors [][]float32
	}
	type fields struct {
		algorithm        string
		compressionLevel int
		compressor       compress.Compressor
		worker           worker.Worker
		workerOpts       []worker.WorkerOption
		eg               errgroup.Group
	}
	type want struct {
		want [][]byte
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, [][]byte, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got [][]byte, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           vectors: nil,
		       },
		       fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
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
		           ctx: nil,
		           vectors: nil,
		           },
		           fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &compressor{
				algorithm:        test.fields.algorithm,
				compressionLevel: test.fields.compressionLevel,
				compressor:       test.fields.compressor,
				worker:           test.fields.worker,
				workerOpts:       test.fields.workerOpts,
				eg:               test.fields.eg,
			}

			got, err := c.MultiCompress(test.args.ctx, test.args.vectors)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_compressor_MultiDecompress(t *testing.T) {
	type args struct {
		ctx    context.Context
		bytess [][]byte
	}
	type fields struct {
		algorithm        string
		compressionLevel int
		compressor       compress.Compressor
		worker           worker.Worker
		workerOpts       []worker.WorkerOption
		eg               errgroup.Group
	}
	type want struct {
		want [][]float32
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, [][]float32, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got [][]float32, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           ctx: nil,
		           bytess: nil,
		       },
		       fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
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
		           ctx: nil,
		           bytess: nil,
		           },
		           fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc(test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &compressor{
				algorithm:        test.fields.algorithm,
				compressionLevel: test.fields.compressionLevel,
				compressor:       test.fields.compressor,
				worker:           test.fields.worker,
				workerOpts:       test.fields.workerOpts,
				eg:               test.fields.eg,
			}

			got, err := c.MultiDecompress(test.args.ctx, test.args.bytess)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_compressor_Len(t *testing.T) {
	type fields struct {
		algorithm        string
		compressionLevel int
		compressor       compress.Compressor
		worker           worker.Worker
		workerOpts       []worker.WorkerOption
		eg               errgroup.Group
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint64) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
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
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &compressor{
				algorithm:        test.fields.algorithm,
				compressionLevel: test.fields.compressionLevel,
				compressor:       test.fields.compressor,
				worker:           test.fields.worker,
				workerOpts:       test.fields.workerOpts,
				eg:               test.fields.eg,
			}

			got := c.Len()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_compressor_TotalRequested(t *testing.T) {
	type fields struct {
		algorithm        string
		compressionLevel int
		compressor       compress.Compressor
		worker           worker.Worker
		workerOpts       []worker.WorkerOption
		eg               errgroup.Group
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint64) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
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
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &compressor{
				algorithm:        test.fields.algorithm,
				compressionLevel: test.fields.compressionLevel,
				compressor:       test.fields.compressor,
				worker:           test.fields.worker,
				workerOpts:       test.fields.workerOpts,
				eg:               test.fields.eg,
			}

			got := c.TotalRequested()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_compressor_TotalCompleted(t *testing.T) {
	type fields struct {
		algorithm        string
		compressionLevel int
		compressor       compress.Compressor
		worker           worker.Worker
		workerOpts       []worker.WorkerOption
		eg               errgroup.Group
	}
	type want struct {
		want uint64
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, uint64) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got uint64) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
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
		           algorithm: "",
		           compressionLevel: 0,
		           compressor: nil,
		           worker: nil,
		           workerOpts: nil,
		           eg: nil,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(t)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			if test.checkFunc == nil {
				test.checkFunc = defaultCheckFunc
			}
			c := &compressor{
				algorithm:        test.fields.algorithm,
				compressionLevel: test.fields.compressionLevel,
				compressor:       test.fields.compressor,
				worker:           test.fields.worker,
				workerOpts:       test.fields.workerOpts,
				eg:               test.fields.eg,
			}

			got := c.TotalCompleted()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
