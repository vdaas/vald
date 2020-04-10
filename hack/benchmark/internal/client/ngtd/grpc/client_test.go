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

// Package grpc provides grpc client functions
package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/client"
	proto "github.com/yahoojapan/ngtd/proto"
)

func TestNew(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		want    Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.ctx, tt.args.opts...)
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

func Test_ngtdClient_Exists(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectID
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		want    *client.ObjectID
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Exists(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.Exists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ngtdClient.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngtdClient_Search(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.SearchRequest
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		want    *client.SearchResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Search(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ngtdClient.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngtdClient_SearchByID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.SearchIDRequest
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		want    *client.SearchResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.SearchByID(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.SearchByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ngtdClient.SearchByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngtdClient_StreamSearch(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.SearchRequest
		f            func(*client.SearchResponse, error)
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamSearch(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.StreamSearch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_StreamSearchByID(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.SearchIDRequest
		f            func(*client.SearchResponse, error)
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamSearchByID(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.StreamSearchByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_Insert(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectVector
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Insert(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_StreamInsert(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectVector
		f            func(error)
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamInsert(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.StreamInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_MultiInsert(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectVectors
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.MultiInsert(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.MultiInsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_Update(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectVector
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Update(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_StreamUpdate(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectVector
		f            func(error)
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamUpdate(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.StreamUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_MultiUpdate(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectVectors
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.MultiUpdate(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.MultiUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_Remove(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectID
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Remove(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_StreamRemove(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectID
		f            func(error)
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamRemove(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.StreamRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_MultiRemove(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectIDs
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.MultiRemove(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.MultiRemove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_GetObject(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ObjectID
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		want    *client.ObjectVector
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetObject(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ngtdClient.GetObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngtdClient_StreamGetObject(t *testing.T) {
	type args struct {
		ctx          context.Context
		dataProvider func() *client.ObjectID
		f            func(*client.ObjectVector, error)
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.StreamGetObject(tt.args.ctx, tt.args.dataProvider, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.StreamGetObject() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_CreateIndex(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ControlCreateIndexRequest
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.CreateIndex(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_SaveIndex(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.SaveIndex(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.SaveIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_CreateAndSaveIndex(t *testing.T) {
	type args struct {
		ctx context.Context
		req *client.ControlCreateIndexRequest
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.CreateAndSaveIndex(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.CreateAndSaveIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngtdClient_IndexInfo(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		c       *ngtdClient
		args    args
		want    *client.InfoIndex
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.IndexInfo(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ngtdClient.IndexInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ngtdClient.IndexInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchRequestToNgtdSearchRequest(t *testing.T) {
	type args struct {
		in *client.SearchRequest
	}
	tests := []struct {
		name string
		args args
		want *proto.SearchRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchRequestToNgtdSearchRequest(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchRequestToNgtdSearchRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchIDRequestToNgtdSearchRequest(t *testing.T) {
	type args struct {
		in *client.SearchIDRequest
	}
	tests := []struct {
		name string
		args args
		want *proto.SearchRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchIDRequestToNgtdSearchRequest(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchIDRequestToNgtdSearchRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngtdSearchResponseToSearchResponse(t *testing.T) {
	type args struct {
		in *proto.SearchResponse
	}
	tests := []struct {
		name string
		args args
		want *client.SearchResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ngtdSearchResponseToSearchResponse(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ngtdSearchResponseToSearchResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngtdGetObjectResponseToObjectVector(t *testing.T) {
	type args struct {
		in *proto.GetObjectResponse
	}
	tests := []struct {
		name string
		args args
		want *client.ObjectVector
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ngtdGetObjectResponseToObjectVector(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ngtdGetObjectResponseToObjectVector() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_objectVectorToNGTDInsertRequest(t *testing.T) {
	type args struct {
		in *client.ObjectVector
	}
	tests := []struct {
		name string
		args args
		want *proto.InsertRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := objectVectorToNGTDInsertRequest(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objectVectorToNGTDInsertRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_objectIDToNGTDRemoveRequest(t *testing.T) {
	type args struct {
		in *client.ObjectID
	}
	tests := []struct {
		name string
		args args
		want *proto.RemoveRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := objectIDToNGTDRemoveRequest(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objectIDToNGTDRemoveRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_objectIDToNGTDGetObjectRequest(t *testing.T) {
	type args struct {
		in *client.ObjectID
	}
	tests := []struct {
		name string
		args args
		want *proto.GetObjectRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := objectIDToNGTDGetObjectRequest(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("objectIDToNGTDGetObjectRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_controlCreateIndexRequestToCreateIndexRequest(t *testing.T) {
	type args struct {
		in *client.ControlCreateIndexRequest
	}
	tests := []struct {
		name string
		args args
		want *proto.CreateIndexRequest
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := controlCreateIndexRequestToCreateIndexRequest(tt.args.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("controlCreateIndexRequestToCreateIndexRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSizeAndEpsilon(t *testing.T) {
	type args struct {
		cfg *client.SearchConfig
	}
	tests := []struct {
		name        string
		args        args
		wantSize    int32
		wantEpsilon float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSize, gotEpsilon := getSizeAndEpsilon(tt.args.cfg)
			if gotSize != tt.wantSize {
				t.Errorf("getSizeAndEpsilon() gotSize = %v, want %v", gotSize, tt.wantSize)
			}
			if gotEpsilon != tt.wantEpsilon {
				t.Errorf("getSizeAndEpsilon() gotEpsilon = %v, want %v", gotEpsilon, tt.wantEpsilon)
			}
		})
	}
}

func Test_tofloat64(t *testing.T) {
	type args struct {
		in []float32
	}
	tests := []struct {
		name    string
		args    args
		wantOut []float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := tofloat64(tt.args.in); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("tofloat64() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
