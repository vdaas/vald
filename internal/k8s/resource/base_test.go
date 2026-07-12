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

package resource

import (
	"bytes"
	"encoding/json"
	"maps"
	"slices"
	"testing"
)

// nested owns a reference field so that the synthetic type exercises the
// DeepCopyIntoer-based helpers.
type nested struct {
	Names []string `json:"names,omitempty"`
}

func (in *nested) DeepCopyInto(out *nested) {
	*out = *in
	out.Names = slices.Clone(in.Names)
}

// synthetic exercises every helper: slice, map, ptr, nested references.
type synthetic struct {
	Base[synthetic, *synthetic] `json:"-"`

	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
	Count  *int              `json:"count,omitempty"`
	Nested *nested           `json:"nested,omitempty"`
	Items  []nested          `json:"items,omitempty"`
}

func (in *synthetic) DeepCopyInto(out *synthetic) {
	*out = *in
	out.Labels = maps.Clone(in.Labels)
	out.Count = CopyPtr(in.Count)
	out.Nested = CopyPtrInto(in.Nested)
	out.Items = CopySliceInto(in.Items)
}

// syntheticNoBase mirrors synthetic without the Base embed for the JSON
// invariance check.
type syntheticNoBase struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
	Count  *int              `json:"count,omitempty"`
	Nested *nested           `json:"nested,omitempty"`
	Items  []nested          `json:"items,omitempty"`
}

func newSynthetic() *synthetic {
	count := 7
	return &synthetic{
		Name:   "origin",
		Labels: map[string]string{"app": "vald"},
		Count:  &count,
		Nested: &nested{Names: []string{"a", "b"}},
		Items:  []nested{{Names: []string{"c"}}, {Names: []string{"d"}}},
	}
}

func TestBase_DeepCopy(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		mutate func(cp *synthetic)
		check  func(t *testing.T, orig *synthetic)
	}

	tests := []test{
		{
			name:   "mutating the copied map does not affect the original",
			mutate: func(cp *synthetic) { cp.Labels["app"] = "mutated" },
			check: func(t *testing.T, orig *synthetic) {
				t.Helper()
				if orig.Labels["app"] != "vald" {
					t.Errorf("original map mutated: %v", orig.Labels)
				}
			},
		},
		{
			name:   "mutating the copied pointer does not affect the original",
			mutate: func(cp *synthetic) { *cp.Count = 99 },
			check: func(t *testing.T, orig *synthetic) {
				t.Helper()
				if *orig.Count != 7 {
					t.Errorf("original pointer mutated: %d", *orig.Count)
				}
			},
		},
		{
			name:   "mutating the copied nested reference does not affect the original",
			mutate: func(cp *synthetic) { cp.Nested.Names[0] = "mutated" },
			check: func(t *testing.T, orig *synthetic) {
				t.Helper()
				if orig.Nested.Names[0] != "a" {
					t.Errorf("original nested slice mutated: %v", orig.Nested.Names)
				}
			},
		},
		{
			name:   "mutating the copied slice elements does not affect the original",
			mutate: func(cp *synthetic) { cp.Items[1].Names[0] = "mutated" },
			check: func(t *testing.T, orig *synthetic) {
				t.Helper()
				if orig.Items[1].Names[0] != "d" {
					t.Errorf("original items mutated: %v", orig.Items)
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			orig := newSynthetic()
			cp := orig.DeepCopy()
			if cp == nil {
				t.Fatal("DeepCopy() = nil, want copy")
			}
			tc.mutate(cp)
			tc.check(t, orig)
		})
	}
}

func TestBase_NilReceiver(t *testing.T) {
	t.Parallel()

	var b *Base[synthetic, *synthetic]
	if got := b.DeepCopy(); got != nil {
		t.Errorf("nil receiver DeepCopy() = %v, want nil", got)
	}
	if got := b.DeepCopyObject(); got != nil {
		t.Errorf("nil receiver DeepCopyObject() = %v, want nil", got)
	}
}

func TestBase_DeepCopyObject(t *testing.T) {
	t.Parallel()

	// synthetic does not implement runtime.Object (no ObjectKind), so
	// DeepCopyObject must return nil. The positive case is covered by the
	// scheme-registered API types (e.g. ValdOperatorRelease) in their own package.
	if got := newSynthetic().DeepCopyObject(); got != nil {
		t.Errorf("DeepCopyObject() = %v, want nil for non API type", got)
	}
}

func TestBase_JSONMarshalInvariance(t *testing.T) {
	t.Parallel()

	count := 7
	with := newSynthetic()
	without := &syntheticNoBase{
		Name:   "origin",
		Labels: map[string]string{"app": "vald"},
		Count:  &count,
		Nested: &nested{Names: []string{"a", "b"}},
		Items:  []nested{{Names: []string{"c"}}, {Names: []string{"d"}}},
	}

	got, err := json.Marshal(with)
	if err != nil {
		t.Fatalf("marshal with Base: %v", err)
	}
	want, err := json.Marshal(without)
	if err != nil {
		t.Fatalf("marshal without Base: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("json output changed by Base embed:\nwith:    %s\nwithout: %s", got, want)
	}
}

func TestCopyHelpers_Nil(t *testing.T) {
	t.Parallel()

	if CopyPtr[int](nil) != nil {
		t.Error("CopyPtr(nil) != nil")
	}
	if CopyPtrInto[nested](nil) != nil {
		t.Error("CopyPtrInto(nil) != nil")
	}
	if CopySliceInto[nested](nil) != nil {
		t.Error("CopySliceInto(nil) != nil")
	}
	if CopyPtrSliceInto[nested](nil) != nil {
		t.Error("CopyPtrSliceInto(nil) != nil")
	}
}

// misplacedBase violates the Base contract: the embed is not the first field.
type misplacedBase struct {
	Name string

	Base[misplacedBase, *misplacedBase] `json:"-"`
}

func (in *misplacedBase) DeepCopyInto(out *misplacedBase) { *out = *in }

// missingBase violates the Base contract: no Base embed at all.
type missingBase struct {
	Name string
}

func (in *missingBase) DeepCopyInto(out *missingBase) { *out = *in }

func TestVerifyBase(t *testing.T) {
	t.Parallel()

	type test struct {
		name    string
		verify  func() error
		wantErr bool
	}

	tests := []test{
		{
			name:   "valid first-field embed passes",
			verify: VerifyBase[synthetic, *synthetic],
		},
		{
			name:    "embed not at offset 0 fails",
			verify:  VerifyBase[misplacedBase, *misplacedBase],
			wantErr: true,
		},
		{
			name:    "missing embed fails",
			verify:  VerifyBase[missingBase, *missingBase],
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if err := tc.verify(); tc.wantErr != (err != nil) {
				t.Errorf("VerifyBase() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
