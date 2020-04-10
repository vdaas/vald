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

package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
)

func TestNewBackup(t *testing.T) {
	type args struct {
		opts []BackupOption
	}
	tests := []struct {
		name    string
		args    args
		wantBu  Backup
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBu, err := NewBackup(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBackup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBu, tt.wantBu) {
				t.Errorf("NewBackup() = %v, want %v", gotBu, tt.wantBu)
			}
		})
	}
}

func Test_backup_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		b       *backup
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.b.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("backup.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("backup.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_backup_GetObject(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		b       *backup
		args    args
		wantVec *payload.Backup_Compressed_MetaVector
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVec, err := tt.b.GetObject(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("backup.GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVec, tt.wantVec) {
				t.Errorf("backup.GetObject() = %v, want %v", gotVec, tt.wantVec)
			}
		})
	}
}

func Test_backup_GetLocation(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name       string
		b          *backup
		args       args
		wantIpList []string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIpList, err := tt.b.GetLocation(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("backup.GetLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIpList, tt.wantIpList) {
				t.Errorf("backup.GetLocation() = %v, want %v", gotIpList, tt.wantIpList)
			}
		})
	}
}

func Test_backup_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		vec *payload.Backup_Compressed_MetaVector
	}
	tests := []struct {
		name    string
		b       *backup
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Register(tt.args.ctx, tt.args.vec); (err != nil) != tt.wantErr {
				t.Errorf("backup.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_backup_RegisterMultiple(t *testing.T) {
	type args struct {
		ctx  context.Context
		vecs *payload.Backup_Compressed_MetaVectors
	}
	tests := []struct {
		name    string
		b       *backup
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.RegisterMultiple(tt.args.ctx, tt.args.vecs); (err != nil) != tt.wantErr {
				t.Errorf("backup.RegisterMultiple() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_backup_Remove(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		b       *backup
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.Remove(tt.args.ctx, tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("backup.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_backup_RemoveMultiple(t *testing.T) {
	type args struct {
		ctx   context.Context
		uuids []string
	}
	tests := []struct {
		name    string
		b       *backup
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.RemoveMultiple(tt.args.ctx, tt.args.uuids...); (err != nil) != tt.wantErr {
				t.Errorf("backup.RemoveMultiple() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_backup_RegisterIPs(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
		ips  []string
	}
	tests := []struct {
		name    string
		b       *backup
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.RegisterIPs(tt.args.ctx, tt.args.uuid, tt.args.ips); (err != nil) != tt.wantErr {
				t.Errorf("backup.RegisterIPs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_backup_RemoveIPs(t *testing.T) {
	type args struct {
		ctx context.Context
		ips []string
	}
	tests := []struct {
		name    string
		b       *backup
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.b.RemoveIPs(tt.args.ctx, tt.args.ips); (err != nil) != tt.wantErr {
				t.Errorf("backup.RemoveIPs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
