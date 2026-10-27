//go:build e2e

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
	"testing"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"sigs.k8s.io/yaml"
)

func TestAlgorithmOptions(t *testing.T) {
	tests := []struct {
		name      string
		options   AlgorithmOptions
		wantError bool
	}{
		{name: "empty"},
		{name: "document", options: AlgorithmOptions{{"@type": "type.googleapis.com/payload.v1.ZVec.DocumentOptions", "fields": map[string]any{"category": "news"}}}},
		{name: "search", options: AlgorithmOptions{{"@type": "type.googleapis.com/payload.v1.ZVec.SearchOptions", "pre_filter": "category = 'news'"}}},
		{name: "mixed unrelated options", options: AlgorithmOptions{{"@type": "type.googleapis.com/google.protobuf.StringValue", "value": "unrelated"}, {"@type": "type.googleapis.com/payload.v1.ZVec.SearchOptions", "pre_filter": "category = 'news'"}}},
		{name: "unknown type", options: AlgorithmOptions{{"@type": "type.googleapis.com/missing.Options"}}, wantError: true},
		{name: "missing type", options: AlgorithmOptions{{"pre_filter": "category = 'news'"}}, wantError: true},
		{name: "misspelled field", options: AlgorithmOptions{{"@type": "type.googleapis.com/payload.v1.ZVec.SearchOptions", "pre_filer": "category = 'news'"}}, wantError: true},
		{name: "wrong fields type", options: AlgorithmOptions{{"@type": "type.googleapis.com/payload.v1.ZVec.DocumentOptions", "fields": "news"}}, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.options.Any()
			if (err != nil) != tt.wantError {
				t.Fatalf("Any() error=%v, wantError=%t", err, tt.wantError)
			}
			if err == nil && len(got) != len(tt.options) {
				t.Errorf("got %d options, want %d", len(got), len(tt.options))
			}
		})
	}
}

func TestAlgorithmOptionsYAML(t *testing.T) {
	var query SearchQuery
	if err := yaml.Unmarshal([]byte(`k: 10
expected_ids: []
options:
  - "@type": type.googleapis.com/payload.v1.ZVec.SearchOptions
    pre_filter: "category = 'news'"
`), &query); err != nil {
		t.Fatal(err)
	}
	if _, err := query.Bind(); err != nil {
		t.Fatal(err)
	}
	if query.ExpectedIDs == nil || len(*query.ExpectedIDs) != 0 {
		t.Fatal("explicit empty expected_ids was lost")
	}
	options, err := query.Options.Any()
	if err != nil {
		t.Fatal(err)
	}
	value := new(payload.ZVec_SearchOptions)
	if err := options[0].UnmarshalTo(value); err != nil {
		t.Fatal(err)
	}
	if value.PreFilter != "category = 'news'" {
		t.Errorf("filter was lost: %q", value.PreFilter)
	}
}

func TestAlgorithmOptionsFixtureID(t *testing.T) {
	for _, tt := range []struct {
		name      string
		num       uint64
		id        string
		wantError bool
	}{
		{name: "single explicit ID", num: 1, id: "fixture"},
		{name: "normal generated IDs", num: 2},
		{name: "repeated explicit ID", num: 2, id: "fixture", wantError: true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			execution := &Execution{BaseConfig: &BaseConfig{Num: tt.num}, Type: OpInsert, Mode: OperationUnary, Modification: &ModificationConfig{ID: tt.id}}
			if _, err := execution.Bind(nil); (err != nil) != tt.wantError {
				t.Errorf("Bind() error=%v, wantError=%t", err, tt.wantError)
			}
		})
	}
}
