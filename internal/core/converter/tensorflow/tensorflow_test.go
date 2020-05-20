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

// Package tensorflow provides implementation of Go API for extract data to vector
package tensorflow

import (
	"reflect"
	"testing"

	tf "github.com/tensorflow/tensorflow/tensorflow/go"
	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	type want struct {
		want TF
		err  error
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, TF, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got TF, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return (t, nil): default options",
			args: args{
				opts: nil,
			},
			want: want{
				want: &tensorflow{
					graph:   nil,
					session: (&tf.SavedModel{}).Session,
				},
				err: nil,
			},
			checkFunc: defaultCheckFunc,
			beforeFunc: func(args args) {
				defaultOpts = []Option{}
				loadFunc = func(s string, ss []string, o *SessionOptions) (*tf.SavedModel, error) {
					return &tf.SavedModel{}, nil
				}
			},
		},
		{
			name: "return (t, nil): args options",
			args: args{
				opts: []Option{
					WithSessionTarget("test"),
					WithSessionConfig([]byte{}),
					WithNdim(1),
				},
			},
			want: want{
				want: &tensorflow{
					options: &tf.SessionOptions{
						Target: "test",
						Config: []byte{},
					},
					graph:   nil,
					session: (&tf.SavedModel{}).Session,
					ndim:    1,
				},
				err: nil,
			},
			checkFunc: defaultCheckFunc,
			beforeFunc: func(args args) {
				defaultOpts = []Option{}
				loadFunc = func(s string, ss []string, o *SessionOptions) (*tf.SavedModel, error) {
					return &tf.SavedModel{}, nil
				}
			},
		},
		{
			name: "return (nil, error)",
			args: args{
				nil,
			},
			want: want{
				want: nil,
				err:  errors.New("load error"),
			},
			checkFunc: defaultCheckFunc,
			beforeFunc: func(args args) {
				defaultOpts = []Option{}
				loadFunc = func(s string, ss []string, o *SessionOptions) (*tf.SavedModel, error) {
					return nil, errors.New("load error")
				}
			},
		},
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

			got, err := New(test.args.opts...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_tensorflow_Close(t *testing.T) {
	type fields struct {
		exportDir     string
		tags          []string
		feeds         []OutputSpec
		fetches       []OutputSpec
		operations    []*Operation
		options       *SessionOptions
		graph         *tf.Graph
		session       session
		ndim          uint8
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
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		return nil
	}
	tests := []test{
		{
			name: "return nil",
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					CloseFunc: func() error {
						return nil
					},
				},
				ndim: 0,
			},
			want: want{
				err: nil,
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "return error",
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					CloseFunc: func() error {
						return errors.New("fail")
					},
				},
				ndim: 0,
			},
			want: want{
				err: errors.New("fail"),
			},
			checkFunc: defaultCheckFunc,
		},
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
			t := &tensorflow{
				exportDir:     test.fields.exportDir,
				tags:          test.fields.tags,
				feeds:         test.fields.feeds,
				fetches:       test.fields.fetches,
				operations:    test.fields.operations,
				options:       test.fields.options,
				graph:         test.fields.graph,
				session:       test.fields.session,
				ndim:          test.fields.ndim,
			}

			err := t.Close()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_tensorflow_run(t *testing.T) {
	type args struct {
		inputs []string
	}
	type fields struct {
		exportDir     string
		tags          []string
		feeds         []OutputSpec
		fetches       []OutputSpec
		operations    []*Operation
		options       *SessionOptions
		graph         *tf.Graph
		session       session
		ndim          uint8
	}
	type want struct {
		want []*tf.Tensor
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []*tf.Tensor, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []*tf.Tensor, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return ([], nil): inputs == nil",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return []*tf.Tensor{}, nil
					},
				},
				ndim: 0,
			},
			want: want{
				want: []*tf.Tensor{},
				err:  nil,
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "return ([], nil): inputs == {\"test\"}",
			args: args{
				inputs: []string{
					"test",
				},
			},
			fields: fields{
				exportDir: "",
				tags:      nil,
				feeds: []OutputSpec{
					OutputSpec{
						operationName: "test",
						outputIndex:   0,
					},
				},
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         tf.NewGraph(),
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return []*tf.Tensor{}, nil
					},
				},
				ndim: 0,
			},
			want: want{
				want: []*tf.Tensor{},
				err:  nil,
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "length error",
			args: args{
				inputs: []string{
					"",
				},
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session:       nil,
				ndim:          0,
			},
			want: want{
				err: errors.ErrInputLength(1, 0),
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "session.Run() error",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         tf.NewGraph(),
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, errors.New("session.Run() error")
					},
				},
				ndim: 0,
			},
			want: want{
				err: errors.New("session.Run() error"),
			},
			checkFunc: defaultCheckFunc,
		},
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
			t := &tensorflow{
				exportDir:     test.fields.exportDir,
				tags:          test.fields.tags,
				feeds:         test.fields.feeds,
				fetches:       test.fields.fetches,
				operations:    test.fields.operations,
				options:       test.fields.options,
				graph:         test.fields.graph,
				session:       test.fields.session,
				ndim:          test.fields.ndim,
			}

			got, err := t.run(test.args.inputs...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_tensorflow_GetVector(t *testing.T) {
	type args struct {
		inputs []string
	}
	type fields struct {
		exportDir     string
		tags          []string
		feeds         []OutputSpec
		fetches       []OutputSpec
		operations    []*Operation
		options       *SessionOptions
		graph         *tf.Graph
		session       session
		ndim          uint8
	}
	type want struct {
		want []float64
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []float64, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got []float64, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return (vector, nil)",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor([]float64{1, 2, 3})
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor}, nil
					},
				},
				ndim: 0,
			},
			want: want{
				want: []float64{1, 2, 3},
				err:  nil,
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "run() error",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, errors.New("session.Run() error")
					},
				},
				ndim: 0,
			},
			want: want{
				want: nil,
				err:  errors.New("session.Run() error"),
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "nil tensor error: run() return nil",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, nil
					},
				},
				ndim: 0,
			},
			want: want{
				want: nil,
				err:  errors.ErrNilTensorTF([]*tf.Tensor{}),
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "nil tensor error: run() return [nil]",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return []*tf.Tensor{nil}, nil
					},
				},
				ndim: 0,
			},
			want: want{
				want: nil,
				err:  errors.ErrNilTensorTF([]*tf.Tensor{nil}),
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "failed to cast error: ndim == TwoDim",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor("test")
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor}, nil
					},
				},
				ndim: 2,
			},
			want: want{
				want: nil,
				err:  errors.ErrFailedToCastTF("test"),
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "failed to cast error: ndim == ThreeDim",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor("test")
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor}, nil
					},
				},
				ndim: 3,
			},
			want: want{
				want: nil,
				err:  errors.ErrFailedToCastTF("test"),
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "failed to cast error: ndim == default",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor("test")
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor}, nil
					},
				},
				ndim: 0,
			},
			want: want{
				want: nil,
				err:  errors.ErrFailedToCastTF("test"),
			},
			checkFunc: defaultCheckFunc,
		},
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
			t := &tensorflow{
				exportDir:     test.fields.exportDir,
				tags:          test.fields.tags,
				feeds:         test.fields.feeds,
				fetches:       test.fields.fetches,
				operations:    test.fields.operations,
				options:       test.fields.options,
				graph:         test.fields.graph,
				session:       test.fields.session,
				ndim:          test.fields.ndim,
			}

			got, err := t.GetVector(test.args.inputs...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_tensorflow_GetValue(t *testing.T) {
	type args struct {
		inputs []string
	}
	type fields struct {
		exportDir     string
		tags          []string
		feeds         []OutputSpec
		fetches       []OutputSpec
		operations    []*Operation
		options       *SessionOptions
		graph         *tf.Graph
		session       session
		ndim          uint8
	}
	type want struct {
		want interface{}
		err  error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, interface{}, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got interface{}, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "return (value, nil)",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor("test")
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor}, nil
					},
				},
				ndim: 0,
			},
			want: want{
				want: "test",
				err:  nil,
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "run() error",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, errors.New("session.Run() error")
					},
				},
				ndim: 0,
			},
			want: want{
				want: nil,
				err:  errors.New("session.Run() error"),
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "nil tensor error: run() return nil",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, nil
					},
				},
				ndim: 0,
			},
			want: want{
				want: nil,
				err:  errors.ErrNilTensorTF([]*tf.Tensor{}),
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "nil tensor error: run() return [nil]",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return []*tf.Tensor{nil}, nil
					},
				},
				ndim: 0,
			},
			want: want{
				want: nil,
				err:  errors.ErrNilTensorTF([]*tf.Tensor{nil}),
			},
			checkFunc: defaultCheckFunc,
		},
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
			t := &tensorflow{
				exportDir:     test.fields.exportDir,
				tags:          test.fields.tags,
				feeds:         test.fields.feeds,
				fetches:       test.fields.fetches,
				operations:    test.fields.operations,
				options:       test.fields.options,
				graph:         test.fields.graph,
				session:       test.fields.session,
				ndim:          test.fields.ndim,
			}

			got, err := t.GetValue(test.args.inputs...)
			if err := test.checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}

func Test_tensorflow_GetValues(t *testing.T) {
	type args struct {
		inputs []string
	}
	type fields struct {
		exportDir     string
		tags          []string
		feeds         []OutputSpec
		fetches       []OutputSpec
		operations    []*Operation
		options       *SessionOptions
		graph         *tf.Graph
		session       session
		ndim          uint8
	}
	type want struct {
		wantValues []interface{}
		err        error
	}
	type test struct {
		name       string
		args       args
		fields     fields
		want       want
		checkFunc  func(want, []interface{}, error) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotValues []interface{}, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got error = %v, want %v", err, w.err)
		}
		if !reflect.DeepEqual(gotValues, w.wantValues) {
			return errors.Errorf("got = %v, want %v", gotValues, w.wantValues)
		}
		return nil
	}
	tests := []test{
		{
			name: "return (values, nil)",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor("test")
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor, tensor}, nil
					},
				},
				ndim: 0,
			},
			want: want{
				wantValues: []interface{}{"test", "test"},
				err:        nil,
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "run() error",
			args: args{
				inputs: nil,
			},
			fields: fields{
				exportDir:     "",
				tags:          nil,
				feeds:         nil,
				fetches:       nil,
				operations:    nil,
				options:       nil,
				graph:         nil,
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, errors.New("session.Run() error")
					},
				},
				ndim: 0,
			},
			want: want{
				wantValues: nil,
				err:        errors.New("session.Run() error"),
			},
			checkFunc: defaultCheckFunc,
		},
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
			t := &tensorflow{
				exportDir:     test.fields.exportDir,
				tags:          test.fields.tags,
				feeds:         test.fields.feeds,
				fetches:       test.fields.fetches,
				operations:    test.fields.operations,
				options:       test.fields.options,
				graph:         test.fields.graph,
				session:       test.fields.session,
				ndim:          test.fields.ndim,
			}

			gotValues, err := t.GetValues(test.args.inputs...)
			if err := test.checkFunc(test.want, gotValues, err); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
