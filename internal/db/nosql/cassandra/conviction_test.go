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

package cassandra

import (
	"reflect"
	"testing"

	"github.com/gocql/gocql"
)

func TestNewConvictionPolicy(t *testing.T) {
	tests := []struct {
		name string
		want gocql.ConvictionPolicy
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConvictionPolicy(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConvictionPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convictionPolicy_AddFailure(t *testing.T) {
	type args struct {
		err  error
		host *gocql.HostInfo
	}
	tests := []struct {
		name string
		c    *convictionPolicy
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.AddFailure(tt.args.err, tt.args.host); got != tt.want {
				t.Errorf("convictionPolicy.AddFailure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convictionPolicy_Reset(t *testing.T) {
	type args struct {
		host *gocql.HostInfo
	}
	tests := []struct {
		name string
		c    *convictionPolicy
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.Reset(tt.args.host)
		})
	}
}
