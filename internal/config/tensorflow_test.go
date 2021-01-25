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

// Package config providers configuration type and load configuration logic
package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"go.uber.org/goleak"
)

func TestTensorflow_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		SessiontOption        *SessionOption
		ExportPath            string
		Tags                  []string
		Feeds                 []*OutputSpec
		FeedsMap              map[string]int
		Fetches               []*OutputSpec
		FetchesMap            map[string]int
		WarmupInputs          []string
		ResultNestedDimension uint8
	}
	type want struct {
		want *Tensorflow
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *Tensorflow) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *Tensorflow) error {
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
		       fields: fields {
		           SessiontOption: SessionOption{},
		           ExportPath: "",
		           Tags: nil,
		           Feeds: nil,
		           FeedsMap: nil,
		           Fetches: nil,
		           FetchesMap: nil,
		           WarmupInputs: nil,
		           ResultNestedDimension: 0,
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
		           SessiontOption: SessionOption{},
		           ExportPath: "",
		           Tags: nil,
		           Feeds: nil,
		           FeedsMap: nil,
		           Fetches: nil,
		           FetchesMap: nil,
		           WarmupInputs: nil,
		           ResultNestedDimension: 0,
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
			tf := &Tensorflow{
				SessiontOption:        test.fields.SessiontOption,
				ExportPath:            test.fields.ExportPath,
				Tags:                  test.fields.Tags,
				Feeds:                 test.fields.Feeds,
				FeedsMap:              test.fields.FeedsMap,
				Fetches:               test.fields.Fetches,
				FetchesMap:            test.fields.FetchesMap,
				WarmupInputs:          test.fields.WarmupInputs,
				ResultNestedDimension: test.fields.ResultNestedDimension,
			}

			got := tf.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestSessionOption_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		Target       string
		Base64Config string
		Config       []byte
	}
	type want struct {
		want *SessionOption
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *SessionOption) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *SessionOption) error {
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
		       fields: fields {
		           Target: "",
		           Base64Config: "",
		           Config: nil,
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
		           Target: "",
		           Base64Config: "",
		           Config: nil,
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
			s := &SessionOption{
				Target:       test.fields.Target,
				Base64Config: test.fields.Base64Config,
				Config:       test.fields.Config,
			}

			got := s.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestOutputSpec_Bind(t *testing.T) {
	t.Parallel()
	type fields struct {
		OperationName string
		OutputIndex   int
	}
	type want struct {
		want *OutputSpec
	}
	type test struct {
		name       string
		fields     fields
		want       want
		checkFunc  func(want, *OutputSpec) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got *OutputSpec) error {
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
		       fields: fields {
		           OperationName: "",
		           OutputIndex: 0,
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
		           OperationName: "",
		           OutputIndex: 0,
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
			o := &OutputSpec{
				OperationName: test.fields.OperationName,
				OutputIndex:   test.fields.OutputIndex,
			}

			got := o.Bind()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
