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

package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/runner"
	"github.com/vdaas/vald/pkg/manager/index/config"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *config.Data
	}
	tests := []struct {
		name    string
		args    args
		wantR   runner.Runner
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, err := New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("New() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}

func Test_run_PreStart(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *run
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.PreStart(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("run.PreStart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_run_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *run
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
				t.Errorf("run.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("run.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_run_PreStop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *run
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.PreStop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("run.PreStop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_run_Stop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *run
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.Stop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("run.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_run_PostStop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		r       *run
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.r.PostStop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("run.PostStop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
