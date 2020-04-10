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

package mysql

import (
	"context"
	"reflect"
	"testing"

	dbr "github.com/gocraft/dbr/v2"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    MySQL
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

func Test_mySQLClient_Open(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		m       *mySQLClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Open(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("mySQLClient.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mySQLClient_Ping(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		m       *mySQLClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Ping(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("mySQLClient.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mySQLClient_Close(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		m       *mySQLClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Close(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("mySQLClient.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mySQLClient_GetMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		m       *mySQLClient
		args    args
		want    MetaVector
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetMeta(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("mySQLClient.GetMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLClient.GetMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mySQLClient_GetIPs(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		m       *mySQLClient
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetIPs(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("mySQLClient.GetIPs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mySQLClient.GetIPs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateMeta(t *testing.T) {
	type args struct {
		meta MetaVector
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateMeta(tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("validateMeta() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mySQLClient_SetMeta(t *testing.T) {
	type args struct {
		ctx context.Context
		mv  MetaVector
	}
	tests := []struct {
		name    string
		m       *mySQLClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SetMeta(tt.args.ctx, tt.args.mv); (err != nil) != tt.wantErr {
				t.Errorf("mySQLClient.SetMeta() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mySQLClient_SetMetas(t *testing.T) {
	type args struct {
		ctx   context.Context
		metas []MetaVector
	}
	tests := []struct {
		name    string
		m       *mySQLClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SetMetas(tt.args.ctx, tt.args.metas...); (err != nil) != tt.wantErr {
				t.Errorf("mySQLClient.SetMetas() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_deleteMetaWithTx(t *testing.T) {
	type args struct {
		ctx  context.Context
		tx   *dbr.Tx
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deleteMetaWithTx(tt.args.ctx, tt.args.tx, tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("deleteMetaWithTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mySQLClient_DeleteMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		m       *mySQLClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.DeleteMeta(tt.args.ctx, tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("mySQLClient.DeleteMeta() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mySQLClient_DeleteMetas(t *testing.T) {
	type args struct {
		ctx   context.Context
		uuids []string
	}
	tests := []struct {
		name    string
		m       *mySQLClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.DeleteMetas(tt.args.ctx, tt.args.uuids...); (err != nil) != tt.wantErr {
				t.Errorf("mySQLClient.DeleteMetas() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mySQLClient_SetIPs(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
		ips  []string
	}
	tests := []struct {
		name    string
		m       *mySQLClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SetIPs(tt.args.ctx, tt.args.uuid, tt.args.ips...); (err != nil) != tt.wantErr {
				t.Errorf("mySQLClient.SetIPs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_mySQLClient_RemoveIPs(t *testing.T) {
	type args struct {
		ctx context.Context
		ips []string
	}
	tests := []struct {
		name    string
		m       *mySQLClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.RemoveIPs(tt.args.ctx, tt.args.ips...); (err != nil) != tt.wantErr {
				t.Errorf("mySQLClient.RemoveIPs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
