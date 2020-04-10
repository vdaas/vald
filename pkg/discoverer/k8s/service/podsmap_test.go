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

	"github.com/vdaas/vald/internal/k8s/pod"
)

func Test_newEntryPodsMap(t *testing.T) {
	type args struct {
		i []pod.Pod
	}
	tests := []struct {
		name string
		args args
		want *entryPodsMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEntryPodsMap(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEntryPodsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_podsMap_Load(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		m         *podsMap
		args      args
		wantValue []pod.Pod
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.m.Load(tt.args.key)
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("podsMap.Load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("podsMap.Load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_entryPodsMap_load(t *testing.T) {
	tests := []struct {
		name      string
		e         *entryPodsMap
		wantValue []pod.Pod
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.e.load()
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("entryPodsMap.load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("entryPodsMap.load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_podsMap_Store(t *testing.T) {
	type args struct {
		key   string
		value []pod.Pod
	}
	tests := []struct {
		name string
		m    *podsMap
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

func Test_entryPodsMap_tryStore(t *testing.T) {
	type args struct {
		i *[]pod.Pod
	}
	tests := []struct {
		name string
		e    *entryPodsMap
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.tryStore(tt.args.i); got != tt.want {
				t.Errorf("entryPodsMap.tryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entryPodsMap_unexpungeLocked(t *testing.T) {
	tests := []struct {
		name            string
		e               *entryPodsMap
		wantWasExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWasExpunged := tt.e.unexpungeLocked(); gotWasExpunged != tt.wantWasExpunged {
				t.Errorf("entryPodsMap.unexpungeLocked() = %v, want %v", gotWasExpunged, tt.wantWasExpunged)
			}
		})
	}
}

func Test_entryPodsMap_storeLocked(t *testing.T) {
	type args struct {
		i *[]pod.Pod
	}
	tests := []struct {
		name string
		e    *entryPodsMap
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

func Test_podsMap_LoadOrStore(t *testing.T) {
	type args struct {
		key   string
		value []pod.Pod
	}
	tests := []struct {
		name       string
		m          *podsMap
		args       args
		wantActual []pod.Pod
		wantLoaded bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActual, gotLoaded := tt.m.LoadOrStore(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(gotActual, tt.wantActual) {
				t.Errorf("podsMap.LoadOrStore() gotActual = %v, want %v", gotActual, tt.wantActual)
			}
			if gotLoaded != tt.wantLoaded {
				t.Errorf("podsMap.LoadOrStore() gotLoaded = %v, want %v", gotLoaded, tt.wantLoaded)
			}
		})
	}
}

func Test_entryPodsMap_tryLoadOrStore(t *testing.T) {
	type args struct {
		i []pod.Pod
	}
	tests := []struct {
		name       string
		e          *entryPodsMap
		args       args
		wantActual []pod.Pod
		wantLoaded bool
		wantOk     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActual, gotLoaded, gotOk := tt.e.tryLoadOrStore(tt.args.i)
			if !reflect.DeepEqual(gotActual, tt.wantActual) {
				t.Errorf("entryPodsMap.tryLoadOrStore() gotActual = %v, want %v", gotActual, tt.wantActual)
			}
			if gotLoaded != tt.wantLoaded {
				t.Errorf("entryPodsMap.tryLoadOrStore() gotLoaded = %v, want %v", gotLoaded, tt.wantLoaded)
			}
			if gotOk != tt.wantOk {
				t.Errorf("entryPodsMap.tryLoadOrStore() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_podsMap_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    *podsMap
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

func Test_entryPodsMap_delete(t *testing.T) {
	tests := []struct {
		name         string
		e            *entryPodsMap
		wantHadValue bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHadValue := tt.e.delete(); gotHadValue != tt.wantHadValue {
				t.Errorf("entryPodsMap.delete() = %v, want %v", gotHadValue, tt.wantHadValue)
			}
		})
	}
}

func Test_podsMap_Range(t *testing.T) {
	type args struct {
		f func(key string, value []pod.Pod) bool
	}
	tests := []struct {
		name string
		m    *podsMap
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

func Test_podsMap_missLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *podsMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.missLocked()
		})
	}
}

func Test_podsMap_dirtyLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *podsMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.dirtyLocked()
		})
	}
}

func Test_entryPodsMap_tryExpungeLocked(t *testing.T) {
	tests := []struct {
		name           string
		e              *entryPodsMap
		wantIsExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsExpunged := tt.e.tryExpungeLocked(); gotIsExpunged != tt.wantIsExpunged {
				t.Errorf("entryPodsMap.tryExpungeLocked() = %v, want %v", gotIsExpunged, tt.wantIsExpunged)
			}
		})
	}
}
