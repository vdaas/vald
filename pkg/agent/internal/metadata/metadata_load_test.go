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

package metadata

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
)

func TestLoad_nilDecode(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "returns ErrMetadataFileEmpty when the file contains only null",
			content: "null",
		},
		{
			name:    "returns ErrMetadataFileEmpty when the file contains only whitespace",
			content: " \n\t",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), AgentMetadataFileName)
			if err := os.WriteFile(path, []byte(tc.content), fs.ModePerm); err != nil {
				t.Fatalf("failed to write metadata: %v", err)
			}
			got, err := Load(path)
			if !errors.Is(err, errors.ErrMetadataFileEmpty) {
				t.Errorf("got error %v, want %v", err, errors.ErrMetadataFileEmpty)
			}
			if got != nil {
				t.Errorf("got metadata %#v, want nil", got)
			}
		})
	}
}
