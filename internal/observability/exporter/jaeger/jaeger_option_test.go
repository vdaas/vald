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

// Package jaeger provides a jaeger exporter.
package jaeger

import (
	"reflect"
	"testing"
)

func TestWithCollectorEndpoint(t *testing.T) {
	type args struct {
		cep string
	}
	tests := []struct {
		name string
		args args
		want JaegerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithCollectorEndpoint(tt.args.cep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithCollectorEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithAgentEndpoint(t *testing.T) {
	type args struct {
		aep string
	}
	tests := []struct {
		name string
		args args
		want JaegerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithAgentEndpoint(tt.args.aep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithAgentEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithUsername(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name string
		args args
		want JaegerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithUsername(tt.args.username); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithUsername() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want JaegerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithPassword(tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithServiceName(t *testing.T) {
	type args struct {
		serviceName string
	}
	tests := []struct {
		name string
		args args
		want JaegerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithServiceName(tt.args.serviceName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithServiceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithBufferMaxCount(t *testing.T) {
	type args struct {
		cnt int
	}
	tests := []struct {
		name string
		args args
		want JaegerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithBufferMaxCount(tt.args.cnt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithBufferMaxCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithOnErrorFunc(t *testing.T) {
	type args struct {
		f func(error)
	}
	tests := []struct {
		name string
		args args
		want JaegerOption
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithOnErrorFunc(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOnErrorFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
