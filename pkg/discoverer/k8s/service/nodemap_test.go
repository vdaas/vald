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

	"github.com/vdaas/vald/internal/k8s/node"
)

func Test_newEntryNodeMap(t *testing.T) {
	type args struct {
		i node.Node
	}
	tests := []struct {
		name string
		args args
		want *entryNodeMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEntryNodeMap(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEntryNodeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeMap_Load(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		m         *nodeMap
		args      args
		wantValue node.Node
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.m.Load(tt.args.key)
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("nodeMap.Load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("nodeMap.Load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_entryNodeMap_load(t *testing.T) {
	tests := []struct {
		name      string
		e         *entryNodeMap
		wantValue node.Node
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.e.load()
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("entryNodeMap.load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("entryNodeMap.load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_nodeMap_Store(t *testing.T) {
	type args struct {
		key   string
		value node.Node
	}
	tests := []struct {
		name string
		m    *nodeMap
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

func Test_entryNodeMap_tryStore(t *testing.T) {
	type args struct {
		i *node.Node
	}
	tests := []struct {
		name string
		e    *entryNodeMap
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.tryStore(tt.args.i); got != tt.want {
				t.Errorf("entryNodeMap.tryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entryNodeMap_unexpungeLocked(t *testing.T) {
	tests := []struct {
		name            string
		e               *entryNodeMap
		wantWasExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWasExpunged := tt.e.unexpungeLocked(); gotWasExpunged != tt.wantWasExpunged {
				t.Errorf("entryNodeMap.unexpungeLocked() = %v, want %v", gotWasExpunged, tt.wantWasExpunged)
			}
		})
	}
}

func Test_entryNodeMap_storeLocked(t *testing.T) {
	type args struct {
		i *node.Node
	}
	tests := []struct {
		name string
		e    *entryNodeMap
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

func Test_nodeMap_LoadOrStore(t *testing.T) {
	type args struct {
		key   string
		value node.Node
	}
	tests := []struct {
		name       string
		m          *nodeMap
		args       args
		wantActual node.Node
		wantLoaded bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActual, gotLoaded := tt.m.LoadOrStore(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(gotActual, tt.wantActual) {
				t.Errorf("nodeMap.LoadOrStore() gotActual = %v, want %v", gotActual, tt.wantActual)
			}
			if gotLoaded != tt.wantLoaded {
				t.Errorf("nodeMap.LoadOrStore() gotLoaded = %v, want %v", gotLoaded, tt.wantLoaded)
			}
		})
	}
}

func Test_entryNodeMap_tryLoadOrStore(t *testing.T) {
	type args struct {
		i node.Node
	}
	tests := []struct {
		name       string
		e          *entryNodeMap
		args       args
		wantActual node.Node
		wantLoaded bool
		wantOk     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActual, gotLoaded, gotOk := tt.e.tryLoadOrStore(tt.args.i)
			if !reflect.DeepEqual(gotActual, tt.wantActual) {
				t.Errorf("entryNodeMap.tryLoadOrStore() gotActual = %v, want %v", gotActual, tt.wantActual)
			}
			if gotLoaded != tt.wantLoaded {
				t.Errorf("entryNodeMap.tryLoadOrStore() gotLoaded = %v, want %v", gotLoaded, tt.wantLoaded)
			}
			if gotOk != tt.wantOk {
				t.Errorf("entryNodeMap.tryLoadOrStore() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_nodeMap_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    *nodeMap
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

func Test_entryNodeMap_delete(t *testing.T) {
	tests := []struct {
		name         string
		e            *entryNodeMap
		wantHadValue bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHadValue := tt.e.delete(); gotHadValue != tt.wantHadValue {
				t.Errorf("entryNodeMap.delete() = %v, want %v", gotHadValue, tt.wantHadValue)
			}
		})
	}
}

func Test_nodeMap_Range(t *testing.T) {
	type args struct {
		f func(key string, value node.Node) bool
	}
	tests := []struct {
		name string
		m    *nodeMap
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

func Test_nodeMap_missLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *nodeMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.missLocked()
		})
	}
}

func Test_nodeMap_dirtyLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *nodeMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.dirtyLocked()
		})
	}
}

func Test_entryNodeMap_tryExpungeLocked(t *testing.T) {
	tests := []struct {
		name           string
		e              *entryNodeMap
		wantIsExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsExpunged := tt.e.tryExpungeLocked(); gotIsExpunged != tt.wantIsExpunged {
				t.Errorf("entryNodeMap.tryExpungeLocked() = %v, want %v", gotIsExpunged, tt.wantIsExpunged)
			}
		})
	}
}
