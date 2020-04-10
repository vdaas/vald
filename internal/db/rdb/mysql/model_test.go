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

package mysql

import (
	"reflect"
	"testing"
)

func Test_metaVector_GetUUID(t *testing.T) {
	tests := []struct {
		name string
		m    *metaVector
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetUUID(); got != tt.want {
				t.Errorf("metaVector.GetUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metaVector_GetVector(t *testing.T) {
	tests := []struct {
		name string
		m    *metaVector
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetVector(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("metaVector.GetVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metaVector_GetMeta(t *testing.T) {
	tests := []struct {
		name string
		m    *metaVector
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetMeta(); got != tt.want {
				t.Errorf("metaVector.GetMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_metaVector_GetIPs(t *testing.T) {
	tests := []struct {
		name string
		m    *metaVector
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.GetIPs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("metaVector.GetIPs() = %v, want %v", got, tt.want)
			}
		})
	}
}
