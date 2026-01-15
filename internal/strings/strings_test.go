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
package strings

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/test/goleak"
)

func TestJoin(t *testing.T) {
	type args struct {
		elems []string
		sep   string
	}
	type want struct {
		wantStr string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotStr string) error {
		if !reflect.DeepEqual(gotStr, w.wantStr) {
			return fmt.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotStr, w.wantStr)
		}
		return nil
	}
	tests := []test{
		func() test {
			return test{
				name: "test_case_2",
				args: args{
					elems: []string{"a", "b", "c"},
					sep:   "/",
				},
				want: want{
					wantStr: "a/b/c",
				},
				checkFunc: defaultCheckFunc,
			}
		}(),
		// This test case consumes a significant memory (maybe 300GB~) due to its characteristic of generating long strings and concatenating multiple strings, and takes a very long time to execute.
		// func() test {
		// 	l := 2 << 24
		// 	c := 1024
		// 	sep := "/"
		// 	rs := randStr(l)
		// 	ss := make([]string, 0, c)
		// 	for i := 0; i < c; i++ {
		// 		ss = append(ss, rs)
		// 	}
		// 	return test{
		// 		name: "test_long_string",
		// 		args: args{
		// 			elems: ss,
		// 			sep:   sep,
		// 		},
		// 		want: want{
		// 			wantStr: strings.Join(ss, sep),
		// 		},
		// 		checkFunc: defaultCheckFunc,
		// 	}
		// }(),
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
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotStr := Join(test.args.elems, test.args.sep)
			if err := checkFunc(test.want, gotStr); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestTrimForCompare(t *testing.T) {
	type args struct {
		origin string
	}
	type want struct {
		wantStr string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, gotStr string) error {
		if !reflect.DeepEqual(gotStr, w.wantStr) {
			return fmt.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotStr, w.wantStr)
		}
		return nil
	}
	tests := []test{
		{
			name: "convert chars",
			args: args{
				origin: "test string_is-kind-of|like/this\r\n\t=?",
			},
			want: want{
				wantStr: "teststringiskindoflikethis?",
			},
			checkFunc: defaultCheckFunc,
		},
		{
			name: "convert upper cases",
			args: args{
				origin: "test string_IS-kind-OF|like/THIS\r\n\t=?",
			},
			want: want{
				wantStr: "teststringiskindoflikethis?",
			},
			checkFunc: defaultCheckFunc,
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
			if test.afterFunc != nil {
				defer test.afterFunc(test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotStr := TrimForCompare(test.args.origin)
			if err := checkFunc(test.want, gotStr); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestReplacerForTrimForCompare(t *testing.T) {
	algo := checkReplacerAlgorithm(t, oldnew...)
	if algo != "ByteStringReplacer" {
		t.Errorf("unexpected algorithm (ByteStringReplacer) for replacer detected %s", algo)
	}
	t.Log(algo)
}

func checkReplacerAlgorithm(t *testing.T, on ...string) string {
	t.Helper()
	if len(on) == 2 && len(on[0]) > 1 {
		return "SingleStringReplacer"
	}
	allNewBytes := true
	for i := 0; i < len(on); i += 2 {
		if len(on[i]) != 1 {
			return "GenericReplacer"
		}
		if len(on[i+1]) != 1 {
			allNewBytes = false
		}
	}
	if allNewBytes {
		return "ByteReplacer"
	}
	return "ByteStringReplacer"
}

// NOT IMPLEMENTED BELOW
