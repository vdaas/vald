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

// Package rest provides vald REST client functions
package rest

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/client"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want Client
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

func Test_gatewayClient_Exists(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectID
	}
	tests := []struct {
		name     string
		c        *gatewayClient
		args     args
		wantResp *client.ObjectID
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.c.Exists(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("gatewayClient.Exists() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_gatewayClient_Search(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.SearchRequest
	}
	tests := []struct {
		name     string
		c        *gatewayClient
		args     args
		wantResp *client.SearchResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.c.Search(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("gatewayClient.Search() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_gatewayClient_SearchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.SearchIDRequest
	}
	tests := []struct {
		name     string
		c        *gatewayClient
		args     args
		wantResp *client.SearchResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.c.SearchByID(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.SearchByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("gatewayClient.SearchByID() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_gatewayClient_StreamSearch(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.SearchRequest
		f            func(*client.SearchResponse, error)
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamSearch(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.StreamSearch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_StreamSearchByID(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.SearchIDRequest
		f            func(*client.SearchResponse, error)
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamSearchByID(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.StreamSearchByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_Insert(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectVector
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Insert(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_StreamInsert(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectVector
		f            func(error)
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamInsert(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.StreamInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_MultiInsert(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectVectors
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.MultiInsert(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.MultiInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectVector
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Update(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_StreamUpdate(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectVector
		f            func(error)
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamUpdate(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.StreamUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_MultiUpdate(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectVectors
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.MultiUpdate(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.MultiUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_Upsert(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectVector
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Upsert(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.Upsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_MultiUpsert(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 *client.ObjectVectors
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.MultiUpsert(tt.args.in0, tt.args.in1); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.MultiUpsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_StreamUpsert(t *testing.T) {
	type args struct {
		in0 context.Context
		in1 func() *client.ObjectVector
		in2 func(error)
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamUpsert(tt.args.in0, tt.args.in1, tt.args.in2); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.StreamUpsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_Remove(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectID
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Remove(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_StreamRemove(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectID
		f            func(error)
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamRemove(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.StreamRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_MultiRemove(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectIDs
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.MultiRemove(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.MultiRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_gatewayClient_GetObject(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectID
	}
	tests := []struct {
		name     string
		c        *gatewayClient
		args     args
		wantResp *client.MetaObject
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResp, err := tt.c.GetObject(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("gatewayClient.GetObject() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}

func Test_gatewayClient_StreamGetObject(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectID
		f            func(*client.MetaObject, error)
	}
	tests := []struct {
		name    string
		c       *gatewayClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamGetObject(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("gatewayClient.StreamGetObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
