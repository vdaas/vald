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

// Package worker provides worker processes
package worker

import (
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []WorkerOption
	}
	tests := []struct {
		name    string
		args    args
		want    Worker
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_worker_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		w       *worker
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.w.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("worker.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("worker.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_worker_startJobLoop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		w    *worker
		args args
		want <-chan error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.startJobLoop(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("worker.startJobLoop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_worker_Pause(t *testing.T) {
	tests := []struct {
		name string
		w    *worker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.Pause()
		})
	}
}

func Test_worker_Resume(t *testing.T) {
	tests := []struct {
		name string
		w    *worker
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.w.Resume()
		})
	}
}

func Test_worker_IsRunning(t *testing.T) {
	tests := []struct {
		name string
		w    *worker
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.IsRunning(); got != tt.want {
				t.Errorf("worker.IsRunning() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_worker_Name(t *testing.T) {
	tests := []struct {
		name string
		w    *worker
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.Name(); got != tt.want {
				t.Errorf("worker.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_worker_Len(t *testing.T) {
	tests := []struct {
		name string
		w    *worker
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.Len(); got != tt.want {
				t.Errorf("worker.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_worker_TotalRequested(t *testing.T) {
	tests := []struct {
		name string
		w    *worker
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.TotalRequested(); got != tt.want {
				t.Errorf("worker.TotalRequested() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_worker_TotalCompleted(t *testing.T) {
	tests := []struct {
		name string
		w    *worker
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.w.TotalCompleted(); got != tt.want {
				t.Errorf("worker.TotalCompleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_worker_Dispatch(t *testing.T) {
	type args struct {
		ctx context.Context
		f   JobFunc
	}
	tests := []struct {
		name    string
		w       *worker
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.w.Dispatch(tt.args.ctx, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("worker.Dispatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
