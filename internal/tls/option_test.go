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

// Package tls provides implementation of Go API for tls certificate provider
package tls

import (
	"crypto/tls"
	"reflect"
	"testing"
)

func TestWithCert(t *testing.T) {
	type args struct {
		cert string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCert(tt.args.cert); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithCa(t *testing.T) {
	type args struct {
		ca string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCa(tt.args.ca); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCa() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithTLSConfig(t *testing.T) {
	type args struct {
		cfg *tls.Config
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithTLSConfig(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithTLSConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
