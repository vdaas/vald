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

// Package rest provides agent ngt REST client functions
package rest

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/client"
)

func TestNew(t *testing.T) {
	type args struct {
		ctx  context.Context
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
			if got := New(tt.args.ctx, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_agentClient_Exists(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectID
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantRes *client.ObjectID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.c.Exists(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("agentClient.Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("agentClient.Exists() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_agentClient_Search(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.SearchRequest
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantRes *client.SearchResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.c.Search(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("agentClient.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("agentClient.Search() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_agentClient_SearchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.SearchIDRequest
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantRes *client.SearchResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.c.SearchByID(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("agentClient.SearchByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("agentClient.SearchByID() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_agentClient_StreamSearch(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.SearchRequest
		f            func(*client.SearchResponse, error)
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamSearch(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.StreamSearch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_StreamSearchByID(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.SearchIDRequest
		f            func(*client.SearchResponse, error)
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamSearchByID(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.StreamSearchByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_Insert(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectVector
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Insert(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_StreamInsert(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectVector
		f            func(error)
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamInsert(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.StreamInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_MultiInsert(t *testing.T) {
	type args struct {
		ctx           context.Context
		objectVectors *client.ObjectVectors
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.MultiInsert(tt.args.ctx, tt.args.objectVectors); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.MultiInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectVector
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Update(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_StreamUpdate(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectVector
		f            func(error)
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamUpdate(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.StreamUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_MultiUpdate(t *testing.T) {
	type args struct {
		ctx           context.Context
		objectVectors *client.ObjectVectors
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.MultiUpdate(tt.args.ctx, tt.args.objectVectors); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.MultiUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_Remove(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectID
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Remove(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_StreamRemove(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectID
		f            func(error)
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamRemove(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.StreamRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_MultiRemove(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectIDs
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.MultiRemove(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.MultiRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_GetObject(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectID
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantRes *client.ObjectVector
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.c.GetObject(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("agentClient.GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("agentClient.GetObject() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_agentClient_StreamGetObject(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectID
		f            func(*client.ObjectVector, error)
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamGetObject(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.StreamGetObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_CreateIndex(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ControlCreateIndexRequest
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.CreateIndex(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_SaveIndex(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SaveIndex(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.SaveIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_CreateAndSaveIndex(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ControlCreateIndexRequest
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.CreateAndSaveIndex(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("agentClient.CreateAndSaveIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_agentClient_IndexInfo(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *agentClient
		args    args
		wantRes *client.InfoIndex
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.c.IndexInfo(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("agentClient.IndexInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("agentClient.IndexInfo() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
