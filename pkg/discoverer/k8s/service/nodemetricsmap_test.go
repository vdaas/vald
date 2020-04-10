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

	mnode "github.com/vdaas/vald/internal/k8s/metrics/node"
)

func Test_newEntryNodeMetricsMap(t *testing.T) {
	type args struct {
		i mnode.Node
	}
	tests := []struct {
		name string
		args args
		want *entryNodeMetricsMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEntryNodeMetricsMap(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEntryNodeMetricsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nodeMetricsMap_Load(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		m         *nodeMetricsMap
		args      args
		wantValue mnode.Node
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.m.Load(tt.args.key)
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("nodeMetricsMap.Load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("nodeMetricsMap.Load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_entryNodeMetricsMap_load(t *testing.T) {
	tests := []struct {
		name      string
		e         *entryNodeMetricsMap
		wantValue mnode.Node
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.e.load()
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("entryNodeMetricsMap.load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("entryNodeMetricsMap.load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_nodeMetricsMap_Store(t *testing.T) {
	type args struct {
		key   string
		value mnode.Node
	}
	tests := []struct {
		name string
		m    *nodeMetricsMap
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

func Test_entryNodeMetricsMap_tryStore(t *testing.T) {
	type args struct {
		i *mnode.Node
	}
	tests := []struct {
		name string
		e    *entryNodeMetricsMap
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.tryStore(tt.args.i); got != tt.want {
				t.Errorf("entryNodeMetricsMap.tryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entryNodeMetricsMap_unexpungeLocked(t *testing.T) {
	tests := []struct {
		name            string
		e               *entryNodeMetricsMap
		wantWasExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWasExpunged := tt.e.unexpungeLocked(); gotWasExpunged != tt.wantWasExpunged {
				t.Errorf("entryNodeMetricsMap.unexpungeLocked() = %v, want %v", gotWasExpunged, tt.wantWasExpunged)
			}
		})
	}
}

func Test_entryNodeMetricsMap_storeLocked(t *testing.T) {
	type args struct {
		i *mnode.Node
	}
	tests := []struct {
		name string
		e    *entryNodeMetricsMap
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

func Test_nodeMetricsMap_LoadOrStore(t *testing.T) {
	type args struct {
		key   string
		value mnode.Node
	}
	tests := []struct {
		name       string
		m          *nodeMetricsMap
		args       args
		wantActual mnode.Node
		wantLoaded bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActual, gotLoaded := tt.m.LoadOrStore(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(gotActual, tt.wantActual) {
				t.Errorf("nodeMetricsMap.LoadOrStore() gotActual = %v, want %v", gotActual, tt.wantActual)
			}
			if gotLoaded != tt.wantLoaded {
				t.Errorf("nodeMetricsMap.LoadOrStore() gotLoaded = %v, want %v", gotLoaded, tt.wantLoaded)
			}
		})
	}
}

func Test_entryNodeMetricsMap_tryLoadOrStore(t *testing.T) {
	type args struct {
		i mnode.Node
	}
	tests := []struct {
		name       string
		e          *entryNodeMetricsMap
		args       args
		wantActual mnode.Node
		wantLoaded bool
		wantOk     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActual, gotLoaded, gotOk := tt.e.tryLoadOrStore(tt.args.i)
			if !reflect.DeepEqual(gotActual, tt.wantActual) {
				t.Errorf("entryNodeMetricsMap.tryLoadOrStore() gotActual = %v, want %v", gotActual, tt.wantActual)
			}
			if gotLoaded != tt.wantLoaded {
				t.Errorf("entryNodeMetricsMap.tryLoadOrStore() gotLoaded = %v, want %v", gotLoaded, tt.wantLoaded)
			}
			if gotOk != tt.wantOk {
				t.Errorf("entryNodeMetricsMap.tryLoadOrStore() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_nodeMetricsMap_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    *nodeMetricsMap
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

func Test_entryNodeMetricsMap_delete(t *testing.T) {
	tests := []struct {
		name         string
		e            *entryNodeMetricsMap
		wantHadValue bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHadValue := tt.e.delete(); gotHadValue != tt.wantHadValue {
				t.Errorf("entryNodeMetricsMap.delete() = %v, want %v", gotHadValue, tt.wantHadValue)
			}
		})
	}
}

func Test_nodeMetricsMap_Range(t *testing.T) {
	type args struct {
		f func(key string, value mnode.Node) bool
	}
	tests := []struct {
		name string
		m    *nodeMetricsMap
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

func Test_nodeMetricsMap_missLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *nodeMetricsMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.missLocked()
		})
	}
}

func Test_nodeMetricsMap_dirtyLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *nodeMetricsMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.dirtyLocked()
		})
	}
}

func Test_entryNodeMetricsMap_tryExpungeLocked(t *testing.T) {
	tests := []struct {
		name           string
		e              *entryNodeMetricsMap
		wantIsExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsExpunged := tt.e.tryExpungeLocked(); gotIsExpunged != tt.wantIsExpunged {
				t.Errorf("entryNodeMetricsMap.tryExpungeLocked() = %v, want %v", gotIsExpunged, tt.wantIsExpunged)
			}
		})
	}
}
