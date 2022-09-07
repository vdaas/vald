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

package timeutil

import (
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/test/goleak"
)

var goleakIgnoreOptions = []goleak.Option{
	goleak.IgnoreTopFunction("github.com/kpango/fastime.(*fastime).StartTimerD.func1"),
}

func TestParse(t *testing.T) {
	type test struct {
		name    string
		t       string
		want    time.Duration
		wantErr bool
	}

	tests := []test{
		{
			name:    "returns time.Nanosecond and nil when t is 1ns",
			t:       "1ns",
			want:    time.Nanosecond,
			wantErr: false,
		},
		{
			name:    "returns time.Millisecond and nil when t is 1ms",
			t:       "1ms",
			want:    time.Millisecond,
			wantErr: false,
		},
		{
			name:    "returns time.Second and nil when t is 1s",
			t:       "1s",
			want:    time.Second,
			wantErr: false,
		},
		{
			name:    "returns tme.Minute and nil when t is 1m",
			t:       "1m",
			want:    time.Minute,
			wantErr: false,
		},
		{
			name:    "returns time.Hour and nil when t is 1h",
			t:       "1h",
			want:    time.Hour,
			wantErr: false,
		},
		{
			name:    "returns 0 and nil when t is empty",
			t:       "",
			want:    0,
			wantErr: false,
		},
		{
			name:    "returns 0 and incorrect string error when t is invalid",
			t:       "dummystring",
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t, goleakIgnoreOptions...)
			got, err := Parse(tt.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseWithDefault(t *testing.T) {
	type args struct {
		t string
		d time.Duration
	}
	type want struct {
		want time.Duration
	}
	type test struct {
		name       string
		args       args
		want       want
		checkFunc  func(want, time.Duration) error
		beforeFunc func(args)
		afterFunc  func(args)
	}
	defaultCheckFunc := func(w want, got time.Duration) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}
	tests := []test{
		{
			name: "returns parsed result when t is a valid string",
			args: args{
				t: "1s",
				d: time.Hour,
			},
			want: want{
				want: time.Second,
			},
		},
		{
			name: "returns default value when t is empty string",
			args: args{
				t: "",
				d: time.Hour,
			},
			want: want{
				want: time.Hour,
			},
		},
		{
			name: "returns default value when t is invalid string",
			args: args{
				t: "hoge",
				d: time.Hour,
			},
			want: want{
				want: time.Hour,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleakIgnoreOptions...)
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

			got := ParseWithDefault(test.args.t, test.args.d)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}
