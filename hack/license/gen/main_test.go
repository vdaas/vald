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

package main

import (
	"io"
	"path/filepath"
	"testing"
	"testing/iotest"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
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

func TestReadAndRewrite(t *testing.T) {
	t.Parallel()

	// longer than bufio.MaxScanTokenSize (64KiB), which truncated the file silently.
	longLine := "key: " + strings.Repeat("a", 3*64*1024)
	tests := []struct {
		name    string
		file    string
		content string
		perm    os.FileMode
	}{
		{
			name:    "keep a line longer than the default scanner token limit",
			file:    "long.yaml",
			content: "first: 1\n" + longLine + "\nlast: 2\n",
			perm:    0o644,
		},
		{
			name:    "keep the file mode of the rewritten file",
			file:    "script.sh",
			content: "#!/bin/sh\n\necho ok\n",
			perm:    0o755,
		},
		{
			// the file is only read; the atomic rename needs no write permission on it.
			name:    "rewrite a read-only file",
			file:    "readonly.yaml",
			content: "key: value\n",
			perm:    0o444,
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			dir := tt.TempDir()
			path := filepath.Join(dir, test.file)
			if err := os.WriteFile(path, []byte(test.content), test.perm); err != nil {
				tt.Fatal(err)
			}
			if err := os.Chmod(path, test.perm); err != nil {
				tt.Fatal(err)
			}

			if err := readAndRewrite(path); err != nil {
				tt.Fatal(err)
			}

			got, err := os.ReadFile(path)
			if err != nil {
				tt.Fatal(err)
			}
			if !strings.Contains(string(got), "Licensed under the Apache License") {
				tt.Errorf("license header is not applied: %q", got[:min(len(got), 200)])
			}
			for _, line := range strings.Split(test.content, "\n") {
				if line != "" && !strings.HasPrefix(line, "#!") && !strings.Contains(string(got), line+"\n") {
					tt.Errorf("line was lost by rewrite: %.40q", line)
				}
			}
			fi, err := os.Stat(path)
			if err != nil {
				tt.Fatal(err)
			}
			if fi.Mode().Perm() != test.perm {
				tt.Errorf("file mode expected %v, got %v", test.perm, fi.Mode().Perm())
			}
			ents, err := os.ReadDir(dir)
			if err != nil {
				tt.Fatal(err)
			}
			if len(ents) != 1 {
				tt.Errorf("temporary file is left behind: %v", ents)
			}
		})
	}
}

func TestReadAndRewriteKeepsOpenErrorCause(t *testing.T) {
	t.Parallel()
	if os.Geteuid() == 0 {
		t.Skip("root can open a file without read permission")
	}
	path := filepath.Join(t.TempDir(), "unreadable.yaml")
	if err := os.WriteFile(path, []byte("key: value\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o000); err != nil {
		t.Fatal(err)
	}
	err := readAndRewrite(path)
	if !errors.Is(err, os.ErrPermission) {
		t.Errorf("readAndRewrite() error = %v, want it to wrap %v", err, os.ErrPermission)
	}
}

func TestGenerateReturnsScanError(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "data.yaml")
	if err := os.WriteFile(path, []byte("a: 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	boom := errors.New("boom")
	r := io.MultiReader(strings.NewReader("a: 1\n"), iotest.ErrReader(boom))
	got, err := generate(r, path, fi, Data{Escape: sharpEscape, Maintainer: defaultMaintainer, Year: 2026})
	if !errors.Is(err, boom) {
		t.Errorf("expected %v, got %v", boom, err)
	}
	if got != "" {
		t.Errorf("expected no content on error, got %q", got)
	}
}
