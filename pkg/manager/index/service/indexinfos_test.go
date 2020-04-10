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
	"reflect"
	"testing"

	"github.com/vdaas/vald/apis/grpc/payload"
)

func Test_newEntryIndexInfos(t *testing.T) {
	type args struct {
		i *payload.Info_Index_Count
	}
	tests := []struct {
		name string
		args args
		want *entryIndexInfos
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEntryIndexInfos(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEntryIndexInfos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_indexInfos_Load(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		m         *indexInfos
		args      args
		wantValue *payload.Info_Index_Count
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.m.Load(tt.args.key)
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("indexInfos.Load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("indexInfos.Load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_entryIndexInfos_load(t *testing.T) {
	tests := []struct {
		name      string
		e         *entryIndexInfos
		wantValue *payload.Info_Index_Count
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.e.load()
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("entryIndexInfos.load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("entryIndexInfos.load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_indexInfos_Store(t *testing.T) {
	type args struct {
		key   string
		value *payload.Info_Index_Count
	}
	tests := []struct {
		name string
		m    *indexInfos
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Store(tt.args.key, tt.args.value)
		})
	}
}

func Test_entryIndexInfos_tryStore(t *testing.T) {
	type args struct {
		i **payload.Info_Index_Count
	}
	tests := []struct {
		name string
		e    *entryIndexInfos
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.tryStore(tt.args.i); got != tt.want {
				t.Errorf("entryIndexInfos.tryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entryIndexInfos_unexpungeLocked(t *testing.T) {
	tests := []struct {
		name            string
		e               *entryIndexInfos
		wantWasExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWasExpunged := tt.e.unexpungeLocked(); gotWasExpunged != tt.wantWasExpunged {
				t.Errorf("entryIndexInfos.unexpungeLocked() = %v, want %v", gotWasExpunged, tt.wantWasExpunged)
			}
		})
	}
}

func Test_entryIndexInfos_storeLocked(t *testing.T) {
	type args struct {
		i **payload.Info_Index_Count
	}
	tests := []struct {
		name string
		e    *entryIndexInfos
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.storeLocked(tt.args.i)
		})
	}
}

func Test_indexInfos_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    *indexInfos
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Delete(tt.args.key)
		})
	}
}

func Test_entryIndexInfos_delete(t *testing.T) {
	tests := []struct {
		name         string
		e            *entryIndexInfos
		wantHadValue bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHadValue := tt.e.delete(); gotHadValue != tt.wantHadValue {
				t.Errorf("entryIndexInfos.delete() = %v, want %v", gotHadValue, tt.wantHadValue)
			}
		})
	}
}

func Test_indexInfos_Range(t *testing.T) {
	type args struct {
		f func(key string, value *payload.Info_Index_Count) bool
	}
	tests := []struct {
		name string
		m    *indexInfos
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Range(tt.args.f)
		})
	}
}

func Test_indexInfos_missLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *indexInfos
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.missLocked()
		})
	}
}

func Test_indexInfos_dirtyLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *indexInfos
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.dirtyLocked()
		})
	}
}

func Test_entryIndexInfos_tryExpungeLocked(t *testing.T) {
	tests := []struct {
		name           string
		e              *entryIndexInfos
		wantIsExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsExpunged := tt.e.tryExpungeLocked(); gotIsExpunged != tt.wantIsExpunged {
				t.Errorf("entryIndexInfos.tryExpungeLocked() = %v, want %v", gotIsExpunged, tt.wantIsExpunged)
			}
		})
	}
}
