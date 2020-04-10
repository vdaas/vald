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

// package config providers configuration type and load configuration logic
package config

import (
	"reflect"
	"testing"
)

func TestGlobalConfig_Bind(t *testing.T) {
	tests := []struct {
		name string
		c    *GlobalConfig
		want *GlobalConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Bind(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GlobalConfig.Bind() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGlobalConfig_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		c       *GlobalConfig
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("GlobalConfig.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRead(t *testing.T) {
	type args struct {
		path string
		cfg  interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Read(tt.args.path, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetActualValue(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := GetActualValue(tt.args.val); gotRes != tt.wantRes {
				t.Errorf("GetActualValue() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestGetActualValues(t *testing.T) {
	type args struct {
		vals []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetActualValues(tt.args.vals); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetActualValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkPrefixAndSuffix(t *testing.T) {
	type args struct {
		str  string
		pref string
		suf  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkPrefixAndSuffix(tt.args.str, tt.args.pref, tt.args.suf); got != tt.want {
				t.Errorf("checkPrefixAndSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToRawYaml(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToRawYaml(tt.args.data); got != tt.want {
				t.Errorf("ToRawYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}
