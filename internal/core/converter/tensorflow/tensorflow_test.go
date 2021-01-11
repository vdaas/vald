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

// Package tensorflow provides implementation of Go API for extract data to vector
package tensorflow

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
	savedModel := &tf.SavedModel{}
	loadFunc := func(exportDir string, tags []string, options *SessionOptions) (*tf.SavedModel, error) {
		return savedModel, nil
	}
	defaultCheckFunc := func(w want, got TF, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		opts := []cmp.Option{
			cmp.AllowUnexported(tensorflow{}),
			cmp.AllowUnexported(OutputSpec{}),
			cmpopts.IgnoreFields(tensorflow{}, "loadFunc"),
			cmp.Comparer(func(want, got TF) bool {
				p1 := reflect.ValueOf(want).Elem().FieldByName("loadFunc").Pointer()
				p2 := reflect.ValueOf(got).Elem().FieldByName("loadFunc").Pointer()
				return p1 == p2
			}),
		}
		if diff := cmp.Diff(w.want, got, opts...); diff != "" {
			return errors.Errorf("err: %s", diff)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns (t, nil) when opts is nil",
			args: args{
				opts: []Option{
					withLoadFunc(loadFunc),
				},
			},
			want: want{
				want: &tensorflow{
					loadFunc: loadFunc,
					graph:    savedModel.Graph,
					session:  savedModel.Session,
				},
			},
			beforeFunc: func(args args) {
				defaultOptions = []Option{}
			},
		},
		{
			name: "returns (t, nil) when args is not nil",
			args: args{
				opts: []Option{
					withLoadFunc(loadFunc),
					WithFeed("test", 0),
					WithFetch("test", 0),
					WithSessionTarget("test"),
					WithSessionConfig([]byte{}),
					WithWarmupInputs(),
					WithNdim(1),
				},
			},
			want: want{
				want: &tensorflow{
					loadFunc: loadFunc,
					feeds: []OutputSpec{
						{
							operationName: "test",
							outputIndex:   0,
						},
					},
					fetches: []OutputSpec{
						{
							operationName: "test",
							outputIndex:   0,
						},
					},
					options: &tf.SessionOptions{
						Target: "test",
						Config: []byte{},
					},
					graph:        savedModel.Graph,
					session:      savedModel.Session,
					warmupInputs: nil,
					ndim:         1,
				},
			},
			beforeFunc: func(args args) {
				defaultOptions = []Option{}
			},
		},
		{
			name: "returns (nil, error) when loadFunc function returns error",
			args: args{
				opts: []Option{
					withLoadFunc(func(exportDir string, tags []string, options *SessionOptions) (*tf.SavedModel, error) {
						return nil, errors.New("load error")
					}),
				},
			},
			want: want{
				err: errors.New("load error"),
			},
			beforeFunc: func(args args) {
				defaultOptions = []Option{}
			},
		},
		{
			name: "returns (nil, error) when warmup error",
			args: args{
				opts: []Option{
					withLoadFunc(loadFunc),
					WithWarmupInputs("test"),
				},
			},
			want: want{
				err: errors.ErrInputLength(1, 0),
			},
			beforeFunc: func(args args) {
				defaultOptions = []Option{}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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

func Test_tensorflow_warmup(t *testing.T) {
	type fields struct {
		feeds        []OutputSpec
		graph        *tf.Graph
		session      session
		warmupInputs []string
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
		{
			name: "return nil when warmupInputs is nil",
			want: want{
				err: nil,
			},
		},
		{
			name: "return nil when warmupInputs is not nil",
			fields: fields{
				feeds: []OutputSpec{
					{
						operationName: "test",
						outputIndex:   0,
					},
				},
				graph: tf.NewGraph(),
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return []*tf.Tensor{}, nil
					},
				},
				warmupInputs: []string{
					"test",
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "return error",
			fields: fields{
				warmupInputs: []string{
					"test",
				},
			},
			want: want{
				err: errors.ErrInputLength(1, 0),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
				feeds:        test.fields.feeds,
				graph:        test.fields.graph,
				session:      test.fields.session,
				warmupInputs: test.fields.warmupInputs,
			}

			err := t.warmup()
			if err := test.checkFunc(test.want, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func Test_tensorflow_Close(t *testing.T) {
	type fields struct {
		session session
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
		{
			name: "return nil",
			fields: fields{
				session: &mockSession{
					CloseFunc: func() error {
						return nil
					},
				},
			},
		},
		{
			name: "return error",
			fields: fields{
				session: &mockSession{
					CloseFunc: func() error {
						return errors.New("fail")
					},
				},
			},
			want: want{
				err: errors.New("fail"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
				session: test.fields.session,
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
		feeds   []OutputSpec
		graph   *tf.Graph
		session session
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns ([], nil) when inputs is nil",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return []*tf.Tensor{}, nil
					},
				},
			},
			want: want{
				want: []*tf.Tensor{},
			},
		},
		{
			name: "returns ([], nil) when inputs is []string{`test`}",
			args: args{
				inputs: []string{
					"test",
				},
			},
			fields: fields{
				feeds: []OutputSpec{
					{
						operationName: "test",
						outputIndex:   0,
					},
				},
				graph: tf.NewGraph(),
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return []*tf.Tensor{}, nil
					},
				},
			},
			want: want{
				want: []*tf.Tensor{},
			},
		},
		{
			name: "returns (nil, error) when length of inputs and feeds field are different",
			args: args{
				inputs: []string{
					"",
				},
			},
			want: want{
				err: errors.ErrInputLength(1, 0),
			},
		},
		{
			name: "returns (nil, error) when Run function returns (nil, error)",
			fields: fields{
				graph: tf.NewGraph(),
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, errors.New("session.Run() error")
					},
				},
			},
			want: want{
				err: errors.New("session.Run() error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
				feeds:   test.fields.feeds,
				graph:   test.fields.graph,
				session: test.fields.session,
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
		session session
		ndim    uint8
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns (vector, nil) when run function returns (tensors, nil) and ndim is default",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor([]float64{
							1,
							2,
							3,
						})
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor}, nil
					},
				},
			},
			want: want{
				want: []float64{
					1,
					2,
					3,
				},
			},
		},
		{
			name: "returns (vector, nil) when run function returns (tensors, nil) and ndim is 2",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor([][]float64{
							{
								1,
								2,
								3,
							},
						})
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor}, nil
					},
				},
				ndim: 2,
			},
			want: want{
				want: []float64{
					1,
					2,
					3,
				},
			},
		},
		{
			name: "returns (vector, nil) when run function returns (tensors, nil) and ndim is 3",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor([][][]float64{
							{
								{
									1,
									2,
									3,
								},
							},
						})
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor}, nil
					},
				},
				ndim: 3,
			},
			want: want{
				want: []float64{
					1,
					2,
					3,
				},
			},
		},
		{
			name: "returns (nil, error) when run function returns (nil, error)",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, errors.New("session.Run() error")
					},
				},
			},
			want: want{
				err: errors.New("session.Run() error"),
			},
		},
		{
			name: "returns (nil, error) when tensors returned by the run funcion is nil",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, nil
					},
				},
			},
			want: want{
				err: errors.ErrNilTensorTF([]*tf.Tensor{}),
			},
		},
		{
			name: "returns (nil, error) when element of tensors returned by the run funcion is nil",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return []*tf.Tensor{nil}, nil
					},
				},
			},
			want: want{
				err: errors.ErrNilTensorTF([]*tf.Tensor{nil}),
			},
		},
		{
			name: "returns (nil, error) when ndim is `TwoDim` and returns error of `ErrFailedToCastTF`",
			fields: fields{
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
				err: errors.ErrFailedToCastTF("test"),
			},
		},
		{
			name: "returns (nil, error) when ndim is `ThreeDim` and returns error of `ErrFailedToCastTF`",
			fields: fields{
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
				err: errors.ErrFailedToCastTF("test"),
			},
		},
		{
			name: "returns (nil, error) when ndim is `default` and returns error of `ErrFailedToCastTF`",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor("test")
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor}, nil
					},
				},
			},
			want: want{
				err: errors.ErrFailedToCastTF("test"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
				session: test.fields.session,
				ndim:    test.fields.ndim,
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
		session session
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns (value, nil) when run function returns (tensors, nil)",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor("test")
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor}, nil
					},
				},
			},
			want: want{
				want: "test",
			},
		},
		{
			name: "returns (nil, error) when run function returns (nil, error)",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, errors.New("session.Run() error")
					},
				},
			},
			want: want{
				err: errors.New("session.Run() error"),
			},
		},
		{
			name: "returns (nil, error) when tensors returned by the run funcion is nil",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, nil
					},
				},
			},
			want: want{
				err: errors.ErrNilTensorTF([]*tf.Tensor{}),
			},
		},
		{
			name: "returns (nil, error) when element of tensors returned by the run funcion is nil",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return []*tf.Tensor{nil}, nil
					},
				},
			},
			want: want{
				err: errors.ErrNilTensorTF([]*tf.Tensor{nil}),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
				session: test.fields.session,
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
		session session
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
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotValues, w.wantValues) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotValues, w.wantValues)
		}
		return nil
	}
	tests := []test{
		{
			name: "return (values, nil) when run function returns (tensors, nil)",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						tensor, err := tf.NewTensor("test")
						if err != nil {
							return nil, errors.New("NewTensor error")
						}
						return []*tf.Tensor{tensor, tensor}, nil
					},
				},
			},
			want: want{
				wantValues: []interface{}{
					"test",
					"test",
				},
			},
		},
		{
			name: "returns (nil, error) when run function returns (nil, error)",
			fields: fields{
				session: &mockSession{
					RunFunc: func(feeds map[tf.Output]*tf.Tensor, fetches []tf.Output, operations []*tf.Operation) ([]*tf.Tensor, error) {
						return nil, errors.New("session.Run() error")
					},
				},
			},
			want: want{
				err: errors.New("session.Run() error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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
				session: test.fields.session,
			}

			gotValues, err := t.GetValues(test.args.inputs...)
			if err := test.checkFunc(test.want, gotValues, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
