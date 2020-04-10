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

// Package service manages the main logic of server.
package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantDsc Discoverer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDsc, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDsc, tt.wantDsc) {
				t.Errorf("New() = %v, want %v", gotDsc, tt.wantDsc)
			}
		})
	}
}

func Test_discoverer_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		d       *discoverer
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("discoverer.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("discoverer.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_discoverer_GetPods(t *testing.T) {
	type args struct {
		req *payload.Discoverer_Request
	}
	tests := []struct {
		name     string
		d        *discoverer
		args     args
		wantPods *payload.Info_Pods
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPods, err := tt.d.GetPods(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("discoverer.GetPods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPods, tt.wantPods) {
				t.Errorf("discoverer.GetPods() = %v, want %v", gotPods, tt.wantPods)
			}
		})
	}
}

func Test_discoverer_GetNodes(t *testing.T) {
	type args struct {
		req *payload.Discoverer_Request
	}
	tests := []struct {
		name      string
		d         *discoverer
		args      args
		wantNodes *payload.Info_Nodes
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNodes, err := tt.d.GetNodes(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("discoverer.GetNodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNodes, tt.wantNodes) {
				t.Errorf("discoverer.GetNodes() = %v, want %v", gotNodes, tt.wantNodes)
			}
		})
	}
}
