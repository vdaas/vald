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

// Package tensorflow provides implementation of Go API for extract data to vector
package tensorflow

import (
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
		want    TF
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

func Test_tensorflow_Close(t *testing.T) {
	tests := []struct {
		name    string
		t       *tensorflow
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.t.Close(); (err != nil) != tt.wantErr {
				t.Errorf("tensorflow.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tensorflow_GetVector(t *testing.T) {
	type args struct {
		feeds   []Feed
		fetches []Fetch
		targets []*Operation
	}
	tests := []struct {
		name       string
		t          *tensorflow
		args       args
		wantValues [][][]float64
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValues, err := tt.t.GetVector(tt.args.feeds, tt.args.fetches, tt.args.targets...)
			if (err != nil) != tt.wantErr {
				t.Errorf("tensorflow.GetVector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotValues, tt.wantValues) {
				t.Errorf("tensorflow.GetVector() = %v, want %v", gotValues, tt.wantValues)
			}
		})
	}
}
