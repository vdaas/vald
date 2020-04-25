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
package logger

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
)

func TestString(t *testing.T) {
	type test struct {
		name    string
		logType Type
		want    string
	}

	tests := []test{
		{
			name:    "returns glg",
			logType: GLG,
			want:    "glg",
		},

		{
			name:    "returns zap",
			logType: ZAP,
			want:    "zap",
		},

		{
			name:    "returns zerolog",
			logType: ZEROLOG,
			want:    "zerolog",
		},

		{
			name:    "returns logrus",
			logType: LOGRUS,
			want:    "logrus",
		},

		{
			name:    "returns klog",
			logType: KLOG,
			want:    "klog",
		},

		{
			name:    "returns unknown",
			logType: Type(100),
			want:    "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.logType.String()
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}

func TestAtot(t *testing.T) {
	type test struct {
		name string
		str  string
		want Type
	}

	tests := []test{
		{
			name: "returns GLG when str is glg",
			str:  "glg",
			want: GLG,
		},

		{
			name: "returns GLG when str is GLg",
			str:  "GLg",
			want: GLG,
		},

		{
			name: "returns ZAP when str is zap",
			str:  "zap",
			want: ZAP,
		},

		{
			name: "returns ZAP when str is ZAp",
			str:  "ZAp",
			want: ZAP,
		},

		{
			name: "returns ZEROLOG when str is zerolog",
			str:  "zerolog",
			want: ZEROLOG,
		},

		{
			name: "returns ZEROLOG when str is ZEROLOg",
			str:  "ZEROLOg",
			want: ZEROLOG,
		},

		{
			name: "returns LOGRUS when str is logrus",
			str:  "logrus",
			want: LOGRUS,
		},

		{
			name: "returns LOGRUS when str is LOGRUs",
			str:  "LOGRUs",
			want: LOGRUS,
		},

		{
			name: "returns KLOG when str is klog",
			str:  "klog",
			want: KLOG,
		},

		{
			name: "returns KLOG when str is KLOg",
			str:  "KLog",
			want: KLOG,
		},

		{
			name: "returns unknown when str is Vald",
			str:  "Vald",
			want: Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atot(tt.str)
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}

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
			return errors.Errorf("got = %v, want %v", got, w.want)
		}
		return nil
	}
	tests := []test{
		// TODO test cases
		/*
		   {
		       name: "test_case_1",
		       want: want{},
		       checkFunc: defaultCheckFunc,
		   },
		*/

		// TODO test cases
		/*
		   func() test {
		       return test {
		           name: "test_case_2",
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

			got := test.m.String()
			if err := test.checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}

		})
	}
}
