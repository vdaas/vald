// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/encoding/json"
	"sigs.k8s.io/yaml"
)

// Regression: the ValdBenchmarkJob CRD and chart schemas define
// object_config.filter_config as {host, port}, which the Go struct (a
// "target" list) used to ignore.
func TestFilterConfig_UnmarshalSchemaForm(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		yaml string
		want []*FilterTarget
	}{
		{
			name: "schema host/port form becomes one target",
			yaml: "object_config:\n  filter_config:\n    host: filter.vald.svc\n    port: 8081\n",
			want: []*FilterTarget{{Host: "filter.vald.svc", Port: 8081}},
		},
		{
			name: "target list form is still read",
			yaml: "object_config:\n  filter_config:\n    target:\n    - host: a\n      port: 1\n    - host: b\n      port: 2\n",
			want: []*FilterTarget{{Host: "a", Port: 1}, {Host: "b", Port: 2}},
		},
		{
			name: "empty filter_config has no target",
			yaml: "object_config:\n  filter_config: {}\n",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got BenchmarkJob
			if err := yaml.Unmarshal([]byte(tt.yaml), &got); err != nil {
				t.Fatalf("yaml.Unmarshal() error = %v", err)
			}
			if got.ObjectConfig == nil {
				t.Fatal("ObjectConfig is nil")
			}
			if !reflect.DeepEqual(got.ObjectConfig.FilterConfig.Targets, tt.want) {
				t.Errorf("Targets = %#v, want %#v", got.ObjectConfig.FilterConfig.Targets, tt.want)
			}
		})
	}
}

// A single target must be written in the schema's host/port form, otherwise
// the API server prunes it when the operator creates a ValdBenchmarkJob.
func TestFilterConfig_MarshalSchemaForm(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want string
		cfg  FilterConfig
	}{
		{
			name: "single target is written as host/port",
			cfg:  FilterConfig{Targets: []*FilterTarget{{Host: "filter.vald.svc", Port: 8081}}},
			want: `{"host":"filter.vald.svc","port":8081}`,
		},
		{
			name: "multiple targets keep the target list",
			cfg:  FilterConfig{Targets: []*FilterTarget{{Host: "a", Port: 1}, {Host: "b", Port: 2}}},
			want: `{"target":[{"host":"a","port":1},{"host":"b","port":2}]}`,
		},
		{
			name: "no target is an empty object",
			cfg:  FilterConfig{},
			want: `{}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			b, err := json.Marshal(ObjectConfig{FilterConfig: tt.cfg})
			if err != nil {
				t.Fatalf("json.Marshal() error = %v", err)
			}
			if want := `{"filter_config":` + tt.want + `}`; string(b) != want {
				t.Errorf("json.Marshal() = %s, want %s", b, want)
			}
			var back ObjectConfig
			if err := json.Unmarshal(b, &back); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if !reflect.DeepEqual(back.FilterConfig.Targets, tt.cfg.Targets) {
				t.Errorf("round trip Targets = %#v, want %#v", back.FilterConfig.Targets, tt.cfg.Targets)
			}
		})
	}
}
