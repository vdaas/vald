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

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
)

// The usecase dereferences IndexCreation.Discoverer unconditionally, so a missing discoverer must be rejected here.
func TestNew_discoverer(t *testing.T) {
	t.Parallel()
	tests := []struct {
		err  error
		name string
		data string
	}{
		{
			name: "return ErrInvalidConfig when discoverer is missing",
			data: `{"server_config": {}, "creator": {}}`,
			err:  errors.ErrInvalidConfig,
		},
		{
			name: "return Data when discoverer is set",
			data: `{"server_config": {}, "creator": {"discoverer": {}}}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			path := file.Join(tt.TempDir(), "config.json")
			if err := os.WriteFile(path, []byte(tc.data), 0o600); err != nil {
				tt.Fatal(err)
			}
			cfg, err := New(path)
			if !errors.Is(err, tc.err) {
				tt.Fatalf("New() error = %v, want %v", err, tc.err)
			}
			if tc.err != nil {
				return
			}
			if cfg == nil || cfg.Creation == nil || cfg.Creation.Discoverer == nil || cfg.Creation.Discoverer.Client == nil {
				tt.Errorf("New() = %+v, want a bound discoverer client", cfg)
			}
		})
	}
}
