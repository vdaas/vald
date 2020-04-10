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

func Test_newEntryOu(t *testing.T) {
	type args struct {
		i string
	}
	tests := []struct {
		name string
		args args
		want *entryOu
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEntryOu(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEntryOu() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ou_Load(t *testing.T) {
	type args struct {
		key uint32
	}
	tests := []struct {
		name      string
		m         *ou
		args      args
		wantValue string
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.m.Load(tt.args.key)
			if gotValue != tt.wantValue {
				t.Errorf("ou.Load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("ou.Load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_entryOu_load(t *testing.T) {
	tests := []struct {
		name      string
		e         *entryOu
		wantValue string
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.e.load()
			if gotValue != tt.wantValue {
				t.Errorf("entryOu.load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("entryOu.load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_ou_Store(t *testing.T) {
	type args struct {
		key   uint32
		value string
	}
	tests := []struct {
		name string
		m    *ou
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

func Test_entryOu_tryStore(t *testing.T) {
	type args struct {
		i *string
	}
	tests := []struct {
		name string
		e    *entryOu
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.tryStore(tt.args.i); got != tt.want {
				t.Errorf("entryOu.tryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entryOu_unexpungeLocked(t *testing.T) {
	tests := []struct {
		name            string
		e               *entryOu
		wantWasExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWasExpunged := tt.e.unexpungeLocked(); gotWasExpunged != tt.wantWasExpunged {
				t.Errorf("entryOu.unexpungeLocked() = %v, want %v", gotWasExpunged, tt.wantWasExpunged)
			}
		})
	}
}

func Test_entryOu_storeLocked(t *testing.T) {
	type args struct {
		i *string
	}
	tests := []struct {
		name string
		e    *entryOu
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

func Test_ou_Delete(t *testing.T) {
	type args struct {
		key uint32
	}
	tests := []struct {
		name string
		m    *ou
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

func Test_entryOu_delete(t *testing.T) {
	tests := []struct {
		name         string
		e            *entryOu
		wantHadValue bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHadValue := tt.e.delete(); gotHadValue != tt.wantHadValue {
				t.Errorf("entryOu.delete() = %v, want %v", gotHadValue, tt.wantHadValue)
			}
		})
	}
}

func Test_ou_missLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *ou
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.missLocked()
		})
	}
}

func Test_ou_dirtyLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *ou
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.dirtyLocked()
		})
	}
}

func Test_entryOu_tryExpungeLocked(t *testing.T) {
	tests := []struct {
		name           string
		e              *entryOu
		wantIsExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsExpunged := tt.e.tryExpungeLocked(); gotIsExpunged != tt.wantIsExpunged {
				t.Errorf("entryOu.tryExpungeLocked() = %v, want %v", gotIsExpunged, tt.wantIsExpunged)
			}
		})
	}
}
