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

// Package k8s provides kubernetes control functionality
package k8s

import (
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantCl  Controller
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCl, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCl, tt.wantCl) {
				t.Errorf("New() = %v, want %v", gotCl, tt.wantCl)
			}
		})
	}
}

func Test_controller_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *controller
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("controller.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("controller.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}
