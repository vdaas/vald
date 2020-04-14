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

func TestNewQueue(t *testing.T) {
	type args struct {
		opts []QueueOption
	}
	tests := []struct {
		name    string
		args    args
		want    Queue
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQueue(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewQueue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queue_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		q       *queue
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("queue.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queue.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queue_pause(t *testing.T) {
	tests := []struct {
		name string
		q    *queue
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.pause()
		})
	}
}

func Test_queue_isRunning(t *testing.T) {
	tests := []struct {
		name string
		q    *queue
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.isRunning(); got != tt.want {
				t.Errorf("queue.isRunning() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queue_Push(t *testing.T) {
	type args struct {
		ctx context.Context
		job JobFunc
	}
	tests := []struct {
		name    string
		q       *queue
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.q.Push(tt.args.ctx, tt.args.job); (err != nil) != tt.wantErr {
				t.Errorf("queue.Push() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_queue_Pop(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		q       *queue
		args    args
		want    JobFunc
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.Pop(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("queue.Pop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queue.Pop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queue_Len(t *testing.T) {
	tests := []struct {
		name string
		q    *queue
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Len(); got != tt.want {
				t.Errorf("queue.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
