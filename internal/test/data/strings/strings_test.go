// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"testing"

	"github.com/vdaas/vald/internal/test/goleak"
)

func TestRandom(t *testing.T) {
	type args struct {
		l int
	}
	type want struct {
		want string
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, string) bool
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
	}
	defaultCheckFunc := func(w want, got string) bool {
		return w.want == got
	}
	tests := func() []test {
		tests := make([]test, 1000)
		length := 10
		for idx := range tests {
			tests[idx] = test{
				name: fmt.Sprintf("random string #%05d", idx),
				args: args{
					l: length,
				},
				want: want{
					want: Random(length),
				},
			}
		}
		return tests
	}()
	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := Random(test.args.l)
			if checkFunc(test.want, got) {
				tt.Errorf("both data is equal even random string want: %s,\tgot: %s", test.want.want, got)
			}
		})
	}
}

// NOT IMPLEMENTED BELOW
