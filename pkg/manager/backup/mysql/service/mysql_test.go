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

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/pkg/manager/backup/mysql/model"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *config.MySQL
	}
	tests := []struct {
		name    string
		args    args
		wantMs  MySQL
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMs, err := New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMs, tt.wantMs) {
				t.Errorf("New() = %v, want %v", gotMs, tt.wantMs)
			}
		})
	}
}

func Test_client_Connect(t *testing.T) {
	type args struct {
		ctx context.Context
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
			if err := tt.c.Connect(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("client.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_Close(t *testing.T) {
	type args struct {
		ctx context.Context
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
			if err := tt.c.Close(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("client.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_GetMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    *model.MetaVector
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetMeta(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.GetMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.GetMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_GetIPs(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		c       *client
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetIPs(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.GetIPs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.GetIPs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_client_SetMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta *model.MetaVector
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
			if err := tt.c.SetMeta(tt.args.ctx, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("client.SetMeta() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_SetMetas(t *testing.T) {
	type args struct {
		ctx   context.Context
		metas []*model.MetaVector
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
			if err := tt.c.SetMetas(tt.args.ctx, tt.args.metas...); (err != nil) != tt.wantErr {
				t.Errorf("client.SetMetas() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_DeleteMeta(t *testing.T) {
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
			if err := tt.c.DeleteMeta(tt.args.ctx, tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("client.DeleteMeta() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_DeleteMetas(t *testing.T) {
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
			if err := tt.c.DeleteMetas(tt.args.ctx, tt.args.uuids...); (err != nil) != tt.wantErr {
				t.Errorf("client.DeleteMetas() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_client_SetIPs(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
		ips  []string
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
			if err := tt.c.SetIPs(tt.args.ctx, tt.args.uuid, tt.args.ips...); (err != nil) != tt.wantErr {
				t.Errorf("client.SetIPs() error = %v, wantErr %v", err, tt.wantErr)
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
			if err := tt.c.RemoveIPs(tt.args.ctx, tt.args.ips...); (err != nil) != tt.wantErr {
				t.Errorf("client.RemoveIPs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
