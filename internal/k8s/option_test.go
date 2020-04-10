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

// Package k8s provides kubernetes control functionality
package k8s

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errgroup"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func TestWithErrGroup(t *testing.T) {
	type args struct {
		eg errgroup.Group
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithErrGroup(tt.args.eg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithErrGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithControllerName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithControllerName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithControllerName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithResourceController(t *testing.T) {
	type args struct {
		rc ResourceController
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithResourceController(tt.args.rc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithResourceController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithManager(t *testing.T) {
	type args struct {
		mgr manager.Manager
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithManager(tt.args.mgr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithMetricsAddress(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithMetricsAddress(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithMetricsAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithEnableLeaderElection(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithEnableLeaderElection(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithEnableLeaderElection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDisableLeaderElection(t *testing.T) {
	tests := []struct {
		name string
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDisableLeaderElection(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithDisableLeaderElection() = %v, want %v", got, tt.want)
			}
		})
	}
}
