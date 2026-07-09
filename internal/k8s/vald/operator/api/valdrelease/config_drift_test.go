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

package valdrelease

import (
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/agent"
	"github.com/vdaas/vald/internal/k8s/vald/operator/api/valdrelease/defaults"
	"github.com/vdaas/vald/internal/strings"
)

// TestConfigTagAlignment guards the sparse overlay projections against drift
// from their internal/config sources of truth. The projection types cannot
// reference internal/config directly — the overlay relies on per-field
// omitempty to emit only operator-managed keys, while internal/config types
// are full bind targets without that contract — so this test asserts instead
// that every projected field keeps the same json key and underlying kind as
// its internal/config counterpart.
func TestConfigTagAlignment(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		api  any
		cfg  any
	}{
		{name: "Logging", api: defaults.Logging{}, cfg: config.Logging{}},
		{name: "Observability", api: defaults.Observability{}, cfg: config.Observability{}},
		{name: "Trace", api: defaults.Trace{}, cfg: config.Trace{}},
		{name: "OTLP", api: defaults.OTLP{}, cfg: config.OTLP{}},
		{name: "NGT", api: agent.NGT{}, cfg: config.NGT{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			cfgKinds := make(map[string]reflect.Kind)
			for f := range reflect.TypeOf(tc.cfg).Fields() {
				key, _, _ := strings.Cut(f.Tag.Get("json"), ",")
				if key == "" || key == "-" {
					continue
				}
				cfgKinds[key] = baseKind(f.Type)
			}
			for f := range reflect.TypeOf(tc.api).Fields() {
				key, _, _ := strings.Cut(f.Tag.Get("json"), ",")
				if key == "" || key == "-" {
					continue
				}
				ck, ok := cfgKinds[key]
				if !ok {
					t.Errorf("field %s: json key %q has no counterpart in %T", f.Name, key, tc.cfg)
					continue
				}
				if ak := baseKind(f.Type); ak != ck && ak != reflect.Struct && ck != reflect.Struct {
					t.Errorf("field %s (%q): kind %s diverges from internal/config kind %s", f.Name, key, ak, ck)
				}
			}
		})
	}
}

// baseKind dereferences pointers so *float32 and float32 compare equal.
func baseKind(t reflect.Type) reflect.Kind {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Kind()
}
