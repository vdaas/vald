//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

package logger

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestType_String(t *testing.T) {
	type want struct {
		want string
	}
	type test struct {
		name       string
		m          Type
		want       want
		checkFunc  func(want, string) error
		beforeFunc func()
		afterFunc  func()
	}
	defaultCheckFunc := func(w want, got string) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns glg when m is GLG",
			m:    GLG,
			want: want{
				want: "glg",
			},
		},

		{
			name: "returns zap when m is ZAP",
			m:    ZAP,
			want: want{
				want: "zap",
			},
		},

		{
			name: "returns nop when m is NOP",
			m:    NOP,
			want: want{
				want: "nop",
			},
		},

		// {
		// 	name: "returns zerolog when m is ZEROLOG",
		// 	m:    ZEROLOG,
		// 	want: want{
		// 		want: "zerolog",
		// 	},
		// },
		//
		// {
		// 	name: "returns logrus when m is LOGRUS",
		// 	m:    LOGRUS,
		// 	want: want{
		// 		want: "logrus",
		// 	},
		// },
		//
		// {
		// 	name: "returns klog when m is KLOG",
		// 	m:    KLOG,
		// 	want: want{
		// 		want: "klog",
		// 	},
		// },

		{
			name: "returns unknown when m is unknown",
			m:    Type(100),
			want: want{
				want: "unknown",
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc()
			}
			if test.afterFunc != nil {
				defer test.afterFunc()
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			got := test.m.String()
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestAtot(t *testing.T) {
	type args struct {
		str string
	}
	type want struct {
		want Type
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, Type) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got Type) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns GLG when str is `glg`",
			args: args{
				str: "glg",
			},
			want: want{
				want: GLG,
			},
		},

		{
			name: "returns GLG when str is `GLg`",
			args: args{
				str: "GLg",
			},
			want: want{
				want: GLG,
			},
		},

		{
			name: "returns ZAP when str is `zap`",
			args: args{
				str: "zap",
			},
			want: want{
				want: ZAP,
			},
		},

		{
			name: "returns ZAP when str is `ZAp`",
			args: args{
				str: "ZAp",
			},
			want: want{
				want: ZAP,
			},
		},

		{
			name: "returns NOP when str is `nop`",
			args: args{
				str: "nop",
			},
			want: want{
				want: NOP,
			},
		},

		{
			name: "returns NOP when str is `NOp`",
			args: args{
				str: "NOp",
			},
			want: want{
				want: NOP,
			},
		},

		{
			name: "returns NOP when str is `EMpty`",
			args: args{
				str: "EMpty",
			},
			want: want{
				want: NOP,
			},
		},

		{
			name: "returns NOP when str is `discard`",
			args: args{
				str: "discard",
			},
			want: want{
				want: NOP,
			},
		},

		{
			name: "returns NOP when str is `DIscard`",
			args: args{
				str: "DIscard",
			},
			want: want{
				want: NOP,
			},
		},

		// {
		// 	name: "returns ZEROLOG when str is `zerolog`",
		// 	args: args{
		// 		str: "zerolog",
		// 	},
		// 	want: want{
		// 		want: ZEROLOG,
		// 	},
		// },
		//
		// {
		// 	name: "returns ZEROLOG when str is `ZEROLOg`",
		// 	args: args{
		// 		str: "ZEROLOg",
		// 	},
		// 	want: want{
		// 		want: ZEROLOG,
		// 	},
		// },
		//
		// {
		// 	name: "returns LOGRUS when str is `logrus`",
		// 	args: args{
		// 		str: "logrus",
		// 	},
		// 	want: want{
		// 		want: LOGRUS,
		// 	},
		// },
		//
		// {
		// 	name: "returns LOGRUS when str is `LOGRUs`",
		// 	args: args{
		// 		str: "LOGRUs",
		// 	},
		// 	want: want{
		// 		want: LOGRUS,
		// 	},
		// },
		//
		// {
		// 	name: "returns KLOG when str is `klog`",
		// 	args: args{
		// 		str: "klog",
		// 	},
		// 	want: want{
		// 		want: KLOG,
		// 	},
		// },
		//
		// {
		// 	name: "returns KLOG when str is `KLOg`",
		// 	args: args{
		// 		str: "KLog",
		// 	},
		// 	want: want{
		// 		want: KLOG,
		// 	},
		// },

		{
			name: "returns unknown when str is `Vald`",
			args: args{
				str: "Vald",
			},
			want: want{
				want: Unknown,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
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

			got := Atot(test.args.str)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
