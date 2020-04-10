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

package kvs

import (
	"context"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want BidiMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bidi_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		b     *bidi
		args  args
		want  uint32
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.b.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("bidi.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("bidi.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_bidi_GetInverse(t *testing.T) {
	type args struct {
		val uint32
	}
	tests := []struct {
		name  string
		b     *bidi
		args  args
		want  string
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.b.GetInverse(tt.args.val)
			if got != tt.want {
				t.Errorf("bidi.GetInverse() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("bidi.GetInverse() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_bidi_Set(t *testing.T) {
	type args struct {
		key string
		val uint32
	}
	tests := []struct {
		name string
		b    *bidi
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Set(tt.args.key, tt.args.val)
		})
	}
}

func Test_bidi_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		b       *bidi
		args    args
		wantVal uint32
		wantOk  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotOk := tt.b.Delete(tt.args.key)
			if gotVal != tt.wantVal {
				t.Errorf("bidi.Delete() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotOk != tt.wantOk {
				t.Errorf("bidi.Delete() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_bidi_DeleteInverse(t *testing.T) {
	type args struct {
		val uint32
	}
	tests := []struct {
		name    string
		b       *bidi
		args    args
		wantKey string
		wantOk  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotOk := tt.b.DeleteInverse(tt.args.val)
			if gotKey != tt.wantKey {
				t.Errorf("bidi.DeleteInverse() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotOk != tt.wantOk {
				t.Errorf("bidi.DeleteInverse() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_bidi_Range(t *testing.T) {
	type args struct {
		ctx context.Context
		f   func(string, uint32) bool
	}
	tests := []struct {
		name string
		b    *bidi
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.b.Range(tt.args.ctx, tt.args.f)
		})
	}
}

func Test_bidi_Len(t *testing.T) {
	tests := []struct {
		name string
		b    *bidi
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Len(); got != tt.want {
				t.Errorf("bidi.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}
