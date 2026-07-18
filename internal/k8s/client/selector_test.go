//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

package client

import (
	"testing"

	"k8s.io/apimachinery/pkg/selection"
)

func TestNewLabelSelector(t *testing.T) {
	t.Parallel()

	const labelKey = "app"

	//nolint:govet // test table struct, field order kept readable rather than memory-optimal
	tests := []struct {
		vals    []string
		name    string
		key     string
		want    string
		op      selection.Operator
		wantErr bool
	}{
		{
			name: "equals",
			key:  labelKey,
			op:   selection.Equals,
			vals: []string{"vald"},
			want: "app=vald",
		},
		{
			name: "exists",
			key:  labelKey,
			op:   selection.Exists,
			vals: []string{},
			want: labelKey,
		},
		{
			name:    "invalid operator returns error",
			key:     labelKey,
			op:      selection.Operator("bogus"),
			vals:    []string{"vald"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			sel, err := NewLabelSelector(tc.key, tc.op, tc.vals)
			if tc.wantErr {
				if err == nil {
					t.Fatal("NewLabelSelector() error = nil, want error")
				}
				return
			}
			if err != nil {
				t.Fatalf("NewLabelSelector() error = %v, want nil", err)
			}
			if sel.String() != tc.want {
				t.Errorf("NewLabelSelector() = %q, want %q", sel.String(), tc.want)
			}
		})
	}
}
