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

// Package service manages the main logic of server.
package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/pkg/agent/ngt/model"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg *config.NGT
	}
	tests := []struct {
		name    string
		args    args
		wantNn  NGT
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNn, err := New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNn, tt.wantNn) {
				t.Errorf("New() = %v, want %v", gotNn, tt.wantNn)
			}
		})
	}
}

func Test_ngt_Start(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		n    *ngt
		args args
		want <-chan error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Start(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ngt.Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngt_Search(t *testing.T) {
	type args struct {
		vec     []float32
		size    uint32
		epsilon float32
		radius  float32
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		want    []model.Distance
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.Search(tt.args.vec, tt.args.size, tt.args.epsilon, tt.args.radius)
			if (err != nil) != tt.wantErr {
				t.Errorf("ngt.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ngt.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngt_SearchByID(t *testing.T) {
	type args struct {
		uuid    string
		size    uint32
		epsilon float32
		radius  float32
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantDst []model.Distance
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDst, err := tt.n.SearchByID(tt.args.uuid, tt.args.size, tt.args.epsilon, tt.args.radius)
			if (err != nil) != tt.wantErr {
				t.Errorf("ngt.SearchByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDst, tt.wantDst) {
				t.Errorf("ngt.SearchByID() = %v, want %v", gotDst, tt.wantDst)
			}
		})
	}
}

func Test_ngt_Insert(t *testing.T) {
	type args struct {
		uuid string
		vec  []float32
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Insert(tt.args.uuid, tt.args.vec); (err != nil) != tt.wantErr {
				t.Errorf("ngt.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngt_insert(t *testing.T) {
	type args struct {
		uuid       string
		vec        []float32
		t          int64
		validation bool
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.insert(tt.args.uuid, tt.args.vec, tt.args.t, tt.args.validation); (err != nil) != tt.wantErr {
				t.Errorf("ngt.insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngt_InsertMultiple(t *testing.T) {
	type args struct {
		vecs map[string][]float32
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.InsertMultiple(tt.args.vecs); (err != nil) != tt.wantErr {
				t.Errorf("ngt.InsertMultiple() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngt_Update(t *testing.T) {
	type args struct {
		uuid string
		vec  []float32
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Update(tt.args.uuid, tt.args.vec); (err != nil) != tt.wantErr {
				t.Errorf("ngt.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngt_UpdateMultiple(t *testing.T) {
	type args struct {
		vecs map[string][]float32
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.UpdateMultiple(tt.args.vecs); (err != nil) != tt.wantErr {
				t.Errorf("ngt.UpdateMultiple() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngt_Delete(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.Delete(tt.args.uuid); (err != nil) != tt.wantErr {
				t.Errorf("ngt.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngt_delete(t *testing.T) {
	type args struct {
		uuid string
		t    int64
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.delete(tt.args.uuid, tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("ngt.delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngt_DeleteMultiple(t *testing.T) {
	type args struct {
		uuids []string
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.DeleteMultiple(tt.args.uuids...); (err != nil) != tt.wantErr {
				t.Errorf("ngt.DeleteMultiple() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngt_GetObject(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantVec []float32
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVec, err := tt.n.GetObject(tt.args.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ngt.GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotVec, tt.wantVec) {
				t.Errorf("ngt.GetObject() = %v, want %v", gotVec, tt.wantVec)
			}
		})
	}
}

func Test_ngt_CreateIndex(t *testing.T) {
	type args struct {
		poolSize uint32
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.CreateIndex(tt.args.poolSize); (err != nil) != tt.wantErr {
				t.Errorf("ngt.CreateIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngt_SaveIndex(t *testing.T) {
	tests := []struct {
		name    string
		n       *ngt
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.SaveIndex(); (err != nil) != tt.wantErr {
				t.Errorf("ngt.SaveIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngt_CreateAndSaveIndex(t *testing.T) {
	type args struct {
		poolSize uint32
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.n.CreateAndSaveIndex(tt.args.poolSize); (err != nil) != tt.wantErr {
				t.Errorf("ngt.CreateAndSaveIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ngt_Close(t *testing.T) {
	tests := []struct {
		name string
		n    *ngt
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.n.Close()
		})
	}
}

func Test_ngt_Exists(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name    string
		n       *ngt
		args    args
		wantOid uint32
		wantOk  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOid, gotOk := tt.n.Exists(tt.args.uuid)
			if gotOid != tt.wantOid {
				t.Errorf("ngt.Exists() gotOid = %v, want %v", gotOid, tt.wantOid)
			}
			if gotOk != tt.wantOk {
				t.Errorf("ngt.Exists() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_ngt_insertCache(t *testing.T) {
	type args struct {
		uuid string
	}
	tests := []struct {
		name  string
		n     *ngt
		args  args
		want  *vcache
		want1 bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.n.insertCache(tt.args.uuid)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ngt.insertCache() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ngt.insertCache() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_ngt_IsIndexing(t *testing.T) {
	tests := []struct {
		name string
		n    *ngt
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.IsIndexing(); got != tt.want {
				t.Errorf("ngt.IsIndexing() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngt_UUIDs(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		n         *ngt
		args      args
		wantUuids []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotUuids := tt.n.UUIDs(tt.args.ctx); !reflect.DeepEqual(gotUuids, tt.wantUuids) {
				t.Errorf("ngt.UUIDs() = %v, want %v", gotUuids, tt.wantUuids)
			}
		})
	}
}

func Test_ngt_UncommittedUUIDs(t *testing.T) {
	tests := []struct {
		name      string
		n         *ngt
		wantUuids []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotUuids := tt.n.UncommittedUUIDs(); !reflect.DeepEqual(gotUuids, tt.wantUuids) {
				t.Errorf("ngt.UncommittedUUIDs() = %v, want %v", gotUuids, tt.wantUuids)
			}
		})
	}
}

func Test_ngt_NumberOfCreateIndexExecution(t *testing.T) {
	tests := []struct {
		name string
		n    *ngt
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.NumberOfCreateIndexExecution(); got != tt.want {
				t.Errorf("ngt.NumberOfCreateIndexExecution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngt_Len(t *testing.T) {
	tests := []struct {
		name string
		n    *ngt
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Len(); got != tt.want {
				t.Errorf("ngt.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngt_InsertVCacheLen(t *testing.T) {
	tests := []struct {
		name string
		n    *ngt
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.InsertVCacheLen(); got != tt.want {
				t.Errorf("ngt.InsertVCacheLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ngt_DeleteVCacheLen(t *testing.T) {
	tests := []struct {
		name string
		n    *ngt
		want uint64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.DeleteVCacheLen(); got != tt.want {
				t.Errorf("ngt.DeleteVCacheLen() = %v, want %v", got, tt.want)
			}
		})
	}
}
