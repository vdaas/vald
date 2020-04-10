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
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantRp  Replicator
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRp, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRp, tt.wantRp) {
				t.Errorf("New() = %v, want %v", gotRp, tt.wantRp)
			}
		})
	}
}

func Test_replicator_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *replicator
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("replicator.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("replicator.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replicator_GetCurrentPodIPs(t *testing.T) {
	tests := []struct {
		name  string
		r     *replicator
		want  []string
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.r.GetCurrentPodIPs()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("replicator.GetCurrentPodIPs() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("replicator.GetCurrentPodIPs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_replicator_SendRecoveryRequest(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *replicator
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.SendRecoveryRequest(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("replicator.SendRecoveryRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
