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

// Package compressor represents compressor client
package compressor

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
		wantC   Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := New(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("New() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func Test_client_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_GetVector(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantVec *payload.Backup_MetaVector
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVec, err := tt.c.GetVector(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.GetVector() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVec, tt.wantVec) {
				t.Errorf("client.GetVector() = %v, want %v", gotVec, tt.wantVec)
			}
		})
	}
}

func Test_client_GetLocation(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name       string
		c          *client
		args       args
		wantIpList []string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIpList, err := tt.c.GetLocation(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.GetLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIpList, tt.wantIpList) {
				t.Errorf("client.GetLocation() = %v, want %v", gotIpList, tt.wantIpList)
			}
		})
	}
}

func Test_client_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		vec *payload.Backup_MetaVector
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Register(tt.args.ctx, tt.args.vec); (err != nil) != tt.wantErr {
				t.Errorf("client.Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_RegisterMultiple(t *testing.T) {
	type args struct {
		ctx  context.Context
		vecs *payload.Backup_MetaVectors
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.RegisterMultiple(tt.args.ctx, tt.args.vecs); (err != nil) != tt.wantErr {
				t.Errorf("client.RegisterMultiple() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Remove(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Remove(tt.args.ctx, tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("client.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_RemoveMultiple(t *testing.T) {
	type args struct {
		ctx   context.Context
		uuids []string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.RemoveMultiple(tt.args.ctx, tt.args.uuids...); (err != nil) != tt.wantErr {
				t.Errorf("client.RemoveMultiple() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_RegisterIPs(t *testing.T) {
	type args struct {
		ctx context.Context
		ips []string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.RegisterIPs(tt.args.ctx, tt.args.ips); (err != nil) != tt.wantErr {
				t.Errorf("client.RegisterIPs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_RemoveIPs(t *testing.T) {
	type args struct {
		ctx context.Context
		ips []string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.RemoveIPs(tt.args.ctx, tt.args.ips); (err != nil) != tt.wantErr {
				t.Errorf("client.RemoveIPs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
