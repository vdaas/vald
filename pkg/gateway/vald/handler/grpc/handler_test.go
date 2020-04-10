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

	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/net/grpc"
)

func TestNew(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want vald.ValdServer
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

func Test_server_Exists(t *testing.T) {
	type args struct {
		ctx  context.Context
		meta *payload.Object_ID
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		want    *payload.Object_ID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Exists(tt.args.ctx, tt.args.meta)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_Search(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_Request
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantRes *payload.Search_Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.Search(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.Search() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_SearchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *payload.Search_IDRequest
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantRes *payload.Search_Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.SearchByID(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.SearchByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.SearchByID() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_search(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *payload.Search_Config
		f   func(ctx context.Context, ac agent.AgentClient, copts ...grpc.CallOption) (*payload.Search_Response, error)
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantRes *payload.Search_Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes, err := tt.s.search(tt.args.ctx, tt.args.cfg, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.search() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_StreamSearch(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamSearchServer
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.StreamSearch(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("server.StreamSearch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_StreamSearchByID(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamSearchByIDServer
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.StreamSearchByID(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("server.StreamSearchByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_Insert(t *testing.T) {
	type args struct {
		ctx context.Context
		vec *payload.Object_Vector
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantCe  *payload.Empty
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCe, err := tt.s.Insert(tt.args.ctx, tt.args.vec)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCe, tt.wantCe) {
				t.Errorf("server.Insert() = %v, want %v", gotCe, tt.wantCe)
			}
		})
	}
}

func Test_server_StreamInsert(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamInsertServer
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.StreamInsert(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("server.StreamInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_MultiInsert(t *testing.T) {
	type args struct {
		ctx  context.Context
		vecs *payload.Object_Vectors
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
			gotRes, err := tt.s.MultiInsert(tt.args.ctx, tt.args.vecs)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.MultiInsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.MultiInsert() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		vec *payload.Object_Vector
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
			gotRes, err := tt.s.Update(tt.args.ctx, tt.args.vec)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.Update() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_StreamUpdate(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamUpdateServer
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.StreamUpdate(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("server.StreamUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_MultiUpdate(t *testing.T) {
	type args struct {
		ctx  context.Context
		vecs *payload.Object_Vectors
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
			gotRes, err := tt.s.MultiUpdate(tt.args.ctx, tt.args.vecs)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.MultiUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.MultiUpdate() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_Upsert(t *testing.T) {
	type args struct {
		ctx context.Context
		vec *payload.Object_Vector
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
			got, err := tt.s.Upsert(tt.args.ctx, tt.args.vec)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Upsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.Upsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_StreamUpsert(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamUpsertServer
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.StreamUpsert(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("server.StreamUpsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_MultiUpsert(t *testing.T) {
	type args struct {
		ctx  context.Context
		vecs *payload.Object_Vectors
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
			got, err := tt.s.MultiUpsert(tt.args.ctx, tt.args.vecs)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.MultiUpsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.MultiUpsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_Remove(t *testing.T) {
	type args struct {
		ctx context.Context
		id  *payload.Object_ID
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
			got, err := tt.s.Remove(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.Remove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("server.Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_server_StreamRemove(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamRemoveServer
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.StreamRemove(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("server.StreamRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_MultiRemove(t *testing.T) {
	type args struct {
		ctx context.Context
		ids *payload.Object_IDs
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
			gotRes, err := tt.s.MultiRemove(tt.args.ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.MultiRemove() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("server.MultiRemove() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_server_GetObject(t *testing.T) {
	type args struct {
		ctx context.Context
		id  *payload.Object_ID
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantVec *payload.Backup_MetaVector
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVec, err := tt.s.GetObject(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("server.GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVec, tt.wantVec) {
				t.Errorf("server.GetObject() = %v, want %v", gotVec, tt.wantVec)
			}
		})
	}
}

func Test_server_StreamGetObject(t *testing.T) {
	type args struct {
		stream vald.Vald_StreamGetObjectServer
	}
	tests := []struct {
		name    string
		s       *server
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.StreamGetObject(tt.args.stream); (err != nil) != tt.wantErr {
				t.Errorf("server.StreamGetObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
