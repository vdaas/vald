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

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/model"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Server
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

func Test_server_GetVector(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Backup_GetVector_Request
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantRes *payload.Backup_Compressed_MetaVector
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.GetVector(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.GetVector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.GetVector() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_Locations(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Backup_Locations_Request
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantRes *payload.Info_IPs
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.Locations(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Locations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.Locations() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_Register(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta *payload.Backup_Compressed_MetaVector
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantRes *payload.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.Register(tt.args.ctx, tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.Register() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_RegisterMulti(t *testing.T) {
	type args struct {
		ctx   context.Context
		metas *payload.Backup_Compressed_MetaVectors
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantRes *payload.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.RegisterMulti(tt.args.ctx, tt.args.metas)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.RegisterMulti() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.RegisterMulti() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_Remove(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Backup_Remove_Request
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantRes *payload.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.Remove(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.Remove() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_RemoveMulti(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Backup_Remove_RequestMulti
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantRes *payload.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.RemoveMulti(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.RemoveMulti() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.RemoveMulti() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_RegisterIPs(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Backup_IP_Register_Request
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantRes *payload.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.RegisterIPs(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.RegisterIPs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.RegisterIPs() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_RemoveIPs(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Backup_IP_Remove_Request
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantRes *payload.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.RemoveIPs(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.RemoveIPs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.RemoveIPs() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_toBackupMetaVector(t *testing.T) {
	type args struct {
		meta *model.MetaVector
	}
	tests := []struct {
		name    string
		args    args
		wantRes *payload.Backup_Compressed_MetaVector
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := toBackupMetaVector(tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("toBackupMetaVector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("toBackupMetaVector() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_toModelMetaVector(t *testing.T) {
	type args struct {
		obj *payload.Backup_Compressed_MetaVector
	}
	tests := []struct {
		name    string
		args    args
		wantRes *model.MetaVector
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := toModelMetaVector(tt.args.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("toModelMetaVector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("toModelMetaVector() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
