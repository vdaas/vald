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
package location

import (
	"reflect"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/errors"
)

func TestGMT(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns GMT location",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GMT()
			if got == nil {
				t.Error("got is nil")
			} else if got, want := got.String(), locationGMT; got != want {
				t.Errorf("String() not equals. want: %v, but got: %v", want, got)
			}
		})
	}
}

func TestUTC(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns UTC location",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UTC()
			if !reflect.DeepEqual(got, time.UTC) {
				t.Errorf("not equals. want: %v, but got: %v", time.UTC, got)
			}
		})
	}
}

func TestJST(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "returns JST location",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JST()
			if got == nil {
				t.Error("got is nil")
			} else if got, want := got.String(), locationJST; got != want {
				t.Errorf("String() not equals. want: %v, but got: %v", want, got)
			}
		})
	}
}

func Test_location(t *testing.T) {
	type args struct {
		zone   string
		offset int
	}

	type test struct {
		name      string
		args      args
		checkFunc func(got *time.Location) error
	}

	tests := []test{
		{
			name: "returns UTC location when zone is UTC",
			args: args{
				zone:   locationUTC,
				offset: 0,
			},
			checkFunc: func(got *time.Location) error {
				if !reflect.DeepEqual(got, time.UTC) {
					return errors.Errorf("not equals. want: %v, but got: %v", time.UTC, got)
				}
				return nil
			},
		},

		{
			name: "returns invalid location when zone is invalid",
			args: args{
				zone:   "invalid",
				offset: 0,
			},
			checkFunc: func(got *time.Location) error {
				if got == nil {
					return errors.New("got is nil")
				} else if got, want := got.String(), "invalid"; got != want {
					return errors.Errorf("String() not equals. want: %v, but got: %v", want, got)
				}
				return nil
			},
		},

		{
			name: "returns UTC location when zone is empty",
			args: args{
				zone:   "",
				offset: 0,
			},
			checkFunc: func(got *time.Location) error {
				if !reflect.DeepEqual(got, time.UTC) {
					return errors.Errorf("not equals. want: %v, but got: %v", time.UTC, got)
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := location(tt.args.zone, tt.args.offset)
			if err := tt.checkFunc(got); err != nil {
				t.Error(err)
			}
		})
	}
}
