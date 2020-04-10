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

	mpod "github.com/vdaas/vald/internal/k8s/metrics/pod"
)

func Test_newEntryPodMetricsMap(t *testing.T) {
	type args struct {
		i mpod.Pod
	}
	tests := []struct {
		name string
		args args
		want *entryPodMetricsMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEntryPodMetricsMap(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEntryPodMetricsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_podMetricsMap_Load(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		m         *podMetricsMap
		args      args
		wantValue mpod.Pod
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.m.Load(tt.args.key)
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("podMetricsMap.Load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("podMetricsMap.Load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_entryPodMetricsMap_load(t *testing.T) {
	tests := []struct {
		name      string
		e         *entryPodMetricsMap
		wantValue mpod.Pod
		wantOk    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, gotOk := tt.e.load()
			if !reflect.DeepEqual(gotValue, tt.wantValue) {
				t.Errorf("entryPodMetricsMap.load() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("entryPodMetricsMap.load() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_podMetricsMap_Store(t *testing.T) {
	type args struct {
		key   string
		value mpod.Pod
	}
	tests := []struct {
		name string
		m    *podMetricsMap
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

func Test_entryPodMetricsMap_tryStore(t *testing.T) {
	type args struct {
		i *mpod.Pod
	}
	tests := []struct {
		name string
		e    *entryPodMetricsMap
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.tryStore(tt.args.i); got != tt.want {
				t.Errorf("entryPodMetricsMap.tryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_entryPodMetricsMap_unexpungeLocked(t *testing.T) {
	tests := []struct {
		name            string
		e               *entryPodMetricsMap
		wantWasExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWasExpunged := tt.e.unexpungeLocked(); gotWasExpunged != tt.wantWasExpunged {
				t.Errorf("entryPodMetricsMap.unexpungeLocked() = %v, want %v", gotWasExpunged, tt.wantWasExpunged)
			}
		})
	}
}

func Test_entryPodMetricsMap_storeLocked(t *testing.T) {
	type args struct {
		i *mpod.Pod
	}
	tests := []struct {
		name string
		e    *entryPodMetricsMap
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

func Test_podMetricsMap_LoadOrStore(t *testing.T) {
	type args struct {
		key   string
		value mpod.Pod
	}
	tests := []struct {
		name       string
		m          *podMetricsMap
		args       args
		wantActual mpod.Pod
		wantLoaded bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActual, gotLoaded := tt.m.LoadOrStore(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(gotActual, tt.wantActual) {
				t.Errorf("podMetricsMap.LoadOrStore() gotActual = %v, want %v", gotActual, tt.wantActual)
			}
			if gotLoaded != tt.wantLoaded {
				t.Errorf("podMetricsMap.LoadOrStore() gotLoaded = %v, want %v", gotLoaded, tt.wantLoaded)
			}
		})
	}
}

func Test_entryPodMetricsMap_tryLoadOrStore(t *testing.T) {
	type args struct {
		i mpod.Pod
	}
	tests := []struct {
		name       string
		e          *entryPodMetricsMap
		args       args
		wantActual mpod.Pod
		wantLoaded bool
		wantOk     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActual, gotLoaded, gotOk := tt.e.tryLoadOrStore(tt.args.i)
			if !reflect.DeepEqual(gotActual, tt.wantActual) {
				t.Errorf("entryPodMetricsMap.tryLoadOrStore() gotActual = %v, want %v", gotActual, tt.wantActual)
			}
			if gotLoaded != tt.wantLoaded {
				t.Errorf("entryPodMetricsMap.tryLoadOrStore() gotLoaded = %v, want %v", gotLoaded, tt.wantLoaded)
			}
			if gotOk != tt.wantOk {
				t.Errorf("entryPodMetricsMap.tryLoadOrStore() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_podMetricsMap_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		m    *podMetricsMap
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

func Test_entryPodMetricsMap_delete(t *testing.T) {
	tests := []struct {
		name         string
		e            *entryPodMetricsMap
		wantHadValue bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHadValue := tt.e.delete(); gotHadValue != tt.wantHadValue {
				t.Errorf("entryPodMetricsMap.delete() = %v, want %v", gotHadValue, tt.wantHadValue)
			}
		})
	}
}

func Test_podMetricsMap_Range(t *testing.T) {
	type args struct {
		f func(key string, value mpod.Pod) bool
	}
	tests := []struct {
		name string
		m    *podMetricsMap
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

func Test_podMetricsMap_missLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *podMetricsMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.missLocked()
		})
	}
}

func Test_podMetricsMap_dirtyLocked(t *testing.T) {
	tests := []struct {
		name string
		m    *podMetricsMap
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.dirtyLocked()
		})
	}
}

func Test_entryPodMetricsMap_tryExpungeLocked(t *testing.T) {
	tests := []struct {
		name           string
		e              *entryPodMetricsMap
		wantIsExpunged bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotIsExpunged := tt.e.tryExpungeLocked(); gotIsExpunged != tt.wantIsExpunged {
				t.Errorf("entryPodMetricsMap.tryExpungeLocked() = %v, want %v", gotIsExpunged, tt.wantIsExpunged)
			}
		})
	}
}
