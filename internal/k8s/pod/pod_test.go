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

// Package pod provides kubernetes pod information and preriodically update
package pod

import (
	"context"
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want PodWatcher
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reconciler_Reconcile(t *testing.T) {
	type args struct {
		req reconcile.Request
	}
	tests := []struct {
		name    string
		r       *reconciler
		args    args
		wantRes reconcile.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.r.Reconcile(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("reconciler.Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("reconciler.Reconcile() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_reconciler_GetName(t *testing.T) {
	tests := []struct {
		name string
		r    *reconciler
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.GetName(); got != tt.want {
				t.Errorf("reconciler.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reconciler_NewReconciler(t *testing.T) {
	type args struct {
		ctx context.Context
		mgr manager.Manager
	}
	tests := []struct {
		name string
		r    *reconciler
		args args
		want reconcile.Reconciler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.NewReconciler(tt.args.ctx, tt.args.mgr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reconciler.NewReconciler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reconciler_For(t *testing.T) {
	tests := []struct {
		name string
		r    *reconciler
		want runtime.Object
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.For(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reconciler.For() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reconciler_Owns(t *testing.T) {
	tests := []struct {
		name string
		r    *reconciler
		want runtime.Object
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Owns(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reconciler.Owns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_reconciler_Watches(t *testing.T) {
	tests := []struct {
		name  string
		r     *reconciler
		want  *source.Kind
		want1 handler.EventHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.r.Watches()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("reconciler.Watches() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("reconciler.Watches() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
