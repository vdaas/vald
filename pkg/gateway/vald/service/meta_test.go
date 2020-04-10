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

// Package service provides meta service
package service

import (
	"context"
	"reflect"
	"testing"
)

func TestNewMeta(t *testing.T) {
	type args struct {
		opts []MetaOption
	}
	tests := []struct {
		name    string
		args    args
		wantMi  Meta
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMi, err := NewMeta(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotMi, tt.wantMi) {
				t.Errorf("NewMeta() = %v, want %v", gotMi, tt.wantMi)
			}
		})
	}
}

func Test_meta_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		want    <-chan error
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Start(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("meta.Start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("meta.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_meta_Exists(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta string
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Exists(tt.args.ctx, tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("meta.Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("meta.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_meta_GetMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		wantV   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV, err := tt.m.GetMeta(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("meta.GetMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotV != tt.wantV {
				t.Errorf("meta.GetMeta() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func Test_meta_GetMetas(t *testing.T) {
	type args struct {
		ctx   context.Context
		uuids []string
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetMetas(tt.args.ctx, tt.args.uuids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("meta.GetMetas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("meta.GetMetas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_meta_GetUUID(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta string
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		wantK   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotK, err := tt.m.GetUUID(tt.args.ctx, tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("meta.GetUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotK != tt.wantK {
				t.Errorf("meta.GetUUID() = %v, want %v", gotK, tt.wantK)
			}
		})
	}
}

func Test_meta_GetUUIDs(t *testing.T) {
	type args struct {
		ctx   context.Context
		metas []string
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.GetUUIDs(tt.args.ctx, tt.args.metas...)
			if (err != nil) != tt.wantErr {
				t.Errorf("meta.GetUUIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("meta.GetUUIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_meta_SetUUIDandMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
		meta string
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SetUUIDandMeta(tt.args.ctx, tt.args.uuid, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("meta.SetUUIDandMeta() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_meta_SetUUIDandMetas(t *testing.T) {
	type args struct {
		ctx context.Context
		kvs map[string]string
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.SetUUIDandMetas(tt.args.ctx, tt.args.kvs); (err != nil) != tt.wantErr {
				t.Errorf("meta.SetUUIDandMetas() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_meta_DeleteMeta(t *testing.T) {
	type args struct {
		ctx  context.Context
		uuid string
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		wantV   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV, err := tt.m.DeleteMeta(tt.args.ctx, tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("meta.DeleteMeta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotV != tt.wantV {
				t.Errorf("meta.DeleteMeta() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func Test_meta_DeleteMetas(t *testing.T) {
	type args struct {
		ctx   context.Context
		uuids []string
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.DeleteMetas(tt.args.ctx, tt.args.uuids...)
			if (err != nil) != tt.wantErr {
				t.Errorf("meta.DeleteMetas() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("meta.DeleteMetas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_meta_DeleteUUID(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta string
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.DeleteUUID(tt.args.ctx, tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("meta.DeleteUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("meta.DeleteUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_meta_DeleteUUIDs(t *testing.T) {
	type args struct {
		ctx   context.Context
		metas []string
	}
	tests := []struct {
		name    string
		m       *meta
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.DeleteUUIDs(tt.args.ctx, tt.args.metas...)
			if (err != nil) != tt.wantErr {
				t.Errorf("meta.DeleteUUIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("meta.DeleteUUIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}
