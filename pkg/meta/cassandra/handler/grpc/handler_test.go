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

// Package grpc provides grpc server logic
package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/meta"
	"github.com/vdaas/vald/apis/grpc/payload"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want meta.MetaServer
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

func Test_server_GetMeta(t *testing.T) {
	type args struct {
		ctx context.Context
		key *payload.Meta_Key
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Meta_Val
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetMeta(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.GetMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.GetMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_GetMetas(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys *payload.Meta_Keys
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantMv  *payload.Meta_Vals
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMv, err := tt.s.GetMetas(tt.args.ctx, tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.GetMetas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMv, tt.wantMv) {
				t.Errorf("server.GetMetas() = %v, want %v", gotMv, tt.wantMv)
			}
		})
	}
}

func Test_server_GetMetaInverse(t *testing.T) {
	type args struct {
		ctx context.Context
		val *payload.Meta_Val
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Meta_Key
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetMetaInverse(tt.args.ctx, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.GetMetaInverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.GetMetaInverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_GetMetasInverse(t *testing.T) {
	type args struct {
		ctx  context.Context
		vals *payload.Meta_Vals
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantMk  *payload.Meta_Keys
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMk, err := tt.s.GetMetasInverse(tt.args.ctx, tt.args.vals)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.GetMetasInverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMk, tt.wantMk) {
				t.Errorf("server.GetMetasInverse() = %v, want %v", gotMk, tt.wantMk)
			}
		})
	}
}

func Test_server_SetMeta(t *testing.T) {
	type args struct {
		ctx context.Context
		kv  *payload.Meta_KeyVal
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.SetMeta(tt.args.ctx, tt.args.kv)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.SetMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.SetMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_SetMetas(t *testing.T) {
	type args struct {
		ctx context.Context
		kvs *payload.Meta_KeyVals
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.SetMetas(tt.args.ctx, tt.args.kvs)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.SetMetas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.SetMetas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_DeleteMeta(t *testing.T) {
	type args struct {
		ctx context.Context
		key *payload.Meta_Key
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Meta_Val
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.DeleteMeta(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.DeleteMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.DeleteMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_DeleteMetas(t *testing.T) {
	type args struct {
		ctx  context.Context
		keys *payload.Meta_Keys
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantMv  *payload.Meta_Vals
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMv, err := tt.s.DeleteMetas(tt.args.ctx, tt.args.keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.DeleteMetas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMv, tt.wantMv) {
				t.Errorf("server.DeleteMetas() = %v, want %v", gotMv, tt.wantMv)
			}
		})
	}
}

func Test_server_DeleteMetaInverse(t *testing.T) {
	type args struct {
		ctx context.Context
		val *payload.Meta_Val
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Meta_Key
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.DeleteMetaInverse(tt.args.ctx, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.DeleteMetaInverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.DeleteMetaInverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_DeleteMetasInverse(t *testing.T) {
	type args struct {
		ctx  context.Context
		vals *payload.Meta_Vals
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantMk  *payload.Meta_Keys
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMk, err := tt.s.DeleteMetasInverse(tt.args.ctx, tt.args.vals)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.DeleteMetasInverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMk, tt.wantMk) {
				t.Errorf("server.DeleteMetasInverse() = %v, want %v", gotMk, tt.wantMk)
			}
		})
	}
}
