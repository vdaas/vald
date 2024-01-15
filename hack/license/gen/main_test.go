//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsSymlink(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	symlinkPath := filepath.Join(dir, "target")
	filePath := filepath.Join(dir, "file")

	_, err := os.Create(filePath)
	if err != nil {
		t.Error(err)
	}

	err = os.Symlink(filePath, symlinkPath)
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "return true when it is a symlink",
			path:     symlinkPath,
			expected: true,
		},
		{
			name:     "return false when it is a normal file",
			path:     filePath,
			expected: false,
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			isSymlink, err := isSymlink(test.path)
			if err != nil {
				tt.Error(err)
			}
			if isSymlink != test.expected {
				tt.Errorf("expected %v, got %v", test.expected, isSymlink)
			}
		})
	}
}
