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
		sessionTarget string
		sessionConfig []byte
		options       *SessionOptions
		graph         *tf.Graph
		session       *tf.Session
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       fields: fields {
		           exportDir: "",
		           tags: nil,
		           feeds: nil,
		           fetches: nil,
		           operations: nil,
		           sessionTarget: "",
		           sessionConfig: nil,
		           options: nil,
		           graph: nil,
		           session: nil,
		           ndim: 0,
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
		           exportDir: "",
		           tags: nil,
		           feeds: nil,
		           fetches: nil,
		           operations: nil,
		           sessionTarget: "",
		           sessionConfig: nil,
		           options: nil,
		           graph: nil,
		           session: nil,
		           ndim: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				sessionTarget: test.fields.sessionTarget,
				sessionConfig: test.fields.sessionConfig,
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
		sessionTarget string
		sessionConfig []byte
		options       *SessionOptions
		graph         *tf.Graph
		session       *tf.Session
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           inputs: nil,
		       },
		       fields: fields {
		           exportDir: "",
		           tags: nil,
		           feeds: nil,
		           fetches: nil,
		           operations: nil,
		           sessionTarget: "",
		           sessionConfig: nil,
		           options: nil,
		           graph: nil,
		           session: nil,
		           ndim: 0,
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
		           inputs: nil,
		           },
		           fields: fields {
		           exportDir: "",
		           tags: nil,
		           feeds: nil,
		           fetches: nil,
		           operations: nil,
		           sessionTarget: "",
		           sessionConfig: nil,
		           options: nil,
		           graph: nil,
		           session: nil,
		           ndim: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				sessionTarget: test.fields.sessionTarget,
				sessionConfig: test.fields.sessionConfig,
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
		sessionTarget string
		sessionConfig []byte
		options       *SessionOptions
		graph         *tf.Graph
		session       *tf.Session
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           inputs: nil,
		       },
		       fields: fields {
		           exportDir: "",
		           tags: nil,
		           feeds: nil,
		           fetches: nil,
		           operations: nil,
		           sessionTarget: "",
		           sessionConfig: nil,
		           options: nil,
		           graph: nil,
		           session: nil,
		           ndim: 0,
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
		           inputs: nil,
		           },
		           fields: fields {
		           exportDir: "",
		           tags: nil,
		           feeds: nil,
		           fetches: nil,
		           operations: nil,
		           sessionTarget: "",
		           sessionConfig: nil,
		           options: nil,
		           graph: nil,
		           session: nil,
		           ndim: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				sessionTarget: test.fields.sessionTarget,
				sessionConfig: test.fields.sessionConfig,
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
		sessionTarget string
		sessionConfig []byte
		options       *SessionOptions
		graph         *tf.Graph
		session       *tf.Session
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           inputs: nil,
		       },
		       fields: fields {
		           exportDir: "",
		           tags: nil,
		           feeds: nil,
		           fetches: nil,
		           operations: nil,
		           sessionTarget: "",
		           sessionConfig: nil,
		           options: nil,
		           graph: nil,
		           session: nil,
		           ndim: 0,
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
		           inputs: nil,
		           },
		           fields: fields {
		           exportDir: "",
		           tags: nil,
		           feeds: nil,
		           fetches: nil,
		           operations: nil,
		           sessionTarget: "",
		           sessionConfig: nil,
		           options: nil,
		           graph: nil,
		           session: nil,
		           ndim: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				sessionTarget: test.fields.sessionTarget,
				sessionConfig: test.fields.sessionConfig,
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
		sessionTarget string
		sessionConfig []byte
		options       *SessionOptions
		graph         *tf.Graph
		session       *tf.Session
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
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       args: args {
		           inputs: nil,
		       },
		       fields: fields {
		           exportDir: "",
		           tags: nil,
		           feeds: nil,
		           fetches: nil,
		           operations: nil,
		           sessionTarget: "",
		           sessionConfig: nil,
		           options: nil,
		           graph: nil,
		           session: nil,
		           ndim: 0,
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
		           inputs: nil,
		           },
		           fields: fields {
		           exportDir: "",
		           tags: nil,
		           feeds: nil,
		           fetches: nil,
		           operations: nil,
		           sessionTarget: "",
		           sessionConfig: nil,
		           options: nil,
		           graph: nil,
		           session: nil,
		           ndim: 0,
		           },
		           want: want{},
		           checkFunc: defaultCheckFunc,
		       }
		   }(),
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
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
				sessionTarget: test.fields.sessionTarget,
				sessionConfig: test.fields.sessionConfig,
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
