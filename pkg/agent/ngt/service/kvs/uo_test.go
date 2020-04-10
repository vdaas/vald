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

package kvs

import (
	"reflect"
	"testing"
)

func Test_newEntryUo(t *testing.T) {
	type args struct {
		i uint32
	}
	tests := []struct {
		name string
		args args
		want *entryUo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEntryUo(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEntryUo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_uo_Load(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		m         *uo
		args      args
		wantValue uint32
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.m.Load(tt.args.key)
			if gotValue != tt.wantValue {
				t.Errorf("uo.Load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("uo.Load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_entryUo_load(t *testing.T) {
	tests := []struct {
		name      string
		e         *entryUo
		wantValue uint32
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.e.load()
			if gotValue != tt.wantValue {
				t.Errorf("entryUo.load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("entryUo.load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_uo_Store(t *testing.T) {
	type args struct {
		key   string
		value uint32
	}
	tests := []struct {
		name string
		m    *uo
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

func Test_entryUo_tryStore(t *testing.T) {
	type args struct {
		i *uint32
	}
	tests := []struct {
		name string
		e    *entryUo
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.tryStore(tt.args.i); got != tt.want {
				t.Errorf("entryUo.tryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entryUo_unexpungeLocked(t *testing.T) {
	tests := []struct {
		name            string
		e               *entryUo
		wantWasExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWasExpunged := tt.e.unexpungeLocked(); gotWasExpunged != tt.wantWasExpunged {
				t.Errorf("entryUo.unexpungeLocked() = %v, want %v", gotWasExpunged, tt.wantWasExpunged)
			}
		})
	}
}

func Test_entryUo_storeLocked(t *testing.T) {
	type args struct {
		i *uint32
	}
	tests := []struct {
		name string
		e    *entryUo
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

func Test_uo_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    *uo
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

func Test_entryUo_delete(t *testing.T) {
	tests := []struct {
		name         string
		e            *entryUo
		wantHadValue bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHadValue := tt.e.delete(); gotHadValue != tt.wantHadValue {
				t.Errorf("entryUo.delete() = %v, want %v", gotHadValue, tt.wantHadValue)
			}
		})
	}
}

func Test_uo_Range(t *testing.T) {
	type args struct {
		f func(uuid string, oid uint32) bool
	}
	tests := []struct {
		name string
		m    *uo
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

func Test_uo_missLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *uo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.missLocked()
		})
	}
}

func Test_uo_dirtyLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *uo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.dirtyLocked()
		})
	}
}

func Test_entryUo_tryExpungeLocked(t *testing.T) {
	tests := []struct {
		name           string
		e              *entryUo
		wantIsExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsExpunged := tt.e.tryExpungeLocked(); gotIsExpunged != tt.wantIsExpunged {
				t.Errorf("entryUo.tryExpungeLocked() = %v, want %v", gotIsExpunged, tt.wantIsExpunged)
			}
		})
	}
}
