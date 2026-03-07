//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

// package jsonpath provides utilities for working with JSONPath in Go.
package jsonpath

import (
	"reflect"
	"testing"
)

func TestJSONPathEval(t *testing.T) {
	jsonStr := []byte(`{
		"test": [1, 2, 3],
		"obj": {"key": "value"},
		"counts": {
			"A": { "stored": 10 },
			"B": { "stored": 5, "uncommitted": 5 }
		},
		"array": [
			{"A": "value1"},
			{"B": "value2"}
		],
		"info": {
			"stored": 100
		}
	}`)

	tests := []struct {
		expected any
		name     string
		path     string
		json     []byte
		wantErr  bool
	}{
		{
			json:     []byte("{}"),
			name:     "empty length",
			path:     "$.length()",
			expected: 0,
			wantErr:  false,
		},
		{
			json:     []byte("{}"),
			name:     "root",
			path:     "$",
			expected: map[string]any{},
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "access map length",
			path:     "$.length()",
			expected: 5,
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "access array",
			path:     ".test",
			expected: []any{float64(1), float64(2), float64(3)},
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "access array length",
			path:     ".test.length()",
			expected: 3,
			wantErr:  false,
		},
		{
			json:    jsonStr,
			name:    "missing key",
			path:    ".missing",
			wantErr: true,
		},
		{
			json:     jsonStr,
			name:     "length on non-array",
			path:     "$.obj.length()",
			expected: 1,
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "nested key access",
			path:     ".obj.key",
			expected: "value",
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "* for map",
			path:     ".obj.*",
			expected: []any{"value"},
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "*.key",
			path:     ".counts.*.stored",
			expected: []any{10.0, 5.0},
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "*.*",
			path:     ".counts.*.*",
			expected: []any{10.0, 5.0, 5.0},
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "*.* for array",
			path:     ".array.*.*",
			expected: []any{"value1", "value2"},
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "array index",
			path:     ".test.0",
			expected: 1.0,
			wantErr:  false,
		},
		{
			json:     jsonStr,
			name:     "sum",
			path:     ".info.sum()",
			expected: 100.0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONPathEval(tt.json, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error = %v, wantErr = %v", err, tt.wantErr)
				if err != nil {
					t.Errorf("error details: %v", err)
				}
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got = %#v, expected = %#v", got, tt.expected)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
//
// func Test_flatten(t *testing.T) {
// 	type args struct {
// 		input []any
// 	}
// 	type want struct {
// 		want []any
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []any) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []any) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           input:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           input:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := flatten(test.args.input)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_recEval(t *testing.T) {
// 	type args struct {
// 		data  any
// 		parts []string
// 	}
// 	type want struct {
// 		want any
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, any, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got any, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           data:nil,
// 		           parts:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           data:nil,
// 		           parts:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := recEval(test.args.data, test.args.parts)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
