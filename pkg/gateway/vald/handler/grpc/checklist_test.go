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

package grpc

import "testing"

func Test_checkList_Exists(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    *checkList
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Exists(tt.args.key); got != tt.want {
				t.Errorf("checkList.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkList_Check(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    *checkList
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.Check(tt.args.key)
		})
	}
}

func Test_entryCheckList_tryStore(t *testing.T) {
	type args struct {
		i *struct{}
	}
	tests := []struct {
		name string
		e    *entryCheckList
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.tryStore(tt.args.i); got != tt.want {
				t.Errorf("entryCheckList.tryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entryCheckList_unexpungeLocked(t *testing.T) {
	tests := []struct {
		name            string
		e               *entryCheckList
		wantWasExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWasExpunged := tt.e.unexpungeLocked(); gotWasExpunged != tt.wantWasExpunged {
				t.Errorf("entryCheckList.unexpungeLocked() = %v, want %v", gotWasExpunged, tt.wantWasExpunged)
			}
		})
	}
}

func Test_checkList_missLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *checkList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.missLocked()
		})
	}
}

func Test_checkList_dirtyLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *checkList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.dirtyLocked()
		})
	}
}

func Test_entryCheckList_tryExpungeLocked(t *testing.T) {
	tests := []struct {
		name           string
		e              *entryCheckList
		wantIsExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsExpunged := tt.e.tryExpungeLocked(); gotIsExpunged != tt.wantIsExpunged {
				t.Errorf("entryCheckList.tryExpungeLocked() = %v, want %v", gotIsExpunged, tt.wantIsExpunged)
			}
		})
	}
}
