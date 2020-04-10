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

// Package node provides kubernetes node information and preriodically update
package node

import (
	"reflect"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func TestWithControllerName(t *testing.T) {
	type args struct {
		name string
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
			if got := WithControllerName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithControllerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithManager(t *testing.T) {
	type args struct {
		mgr manager.Manager
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
			if got := WithManager(tt.args.mgr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithOnErrorFunc(t *testing.T) {
	type args struct {
		f func(err error)
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
			if got := WithOnErrorFunc(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOnErrorFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithOnReconcileFunc(t *testing.T) {
	type args struct {
		f func(nodes map[string]Node)
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
			if got := WithOnReconcileFunc(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithOnReconcileFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
