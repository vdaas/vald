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
package errors

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type fitosFakeFileInfo struct {
	name string
	size int64
	mode fs.FileMode
}

func (f fitosFakeFileInfo) Name() string       { return f.name }
func (f fitosFakeFileInfo) Size() int64        { return f.size }
func (f fitosFakeFileInfo) Mode() fs.FileMode  { return f.mode }
func (f fitosFakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f fitosFakeFileInfo) IsDir() bool        { return f.mode.IsDir() }
func (f fitosFakeFileInfo) Sys() any           { return nil }

func TestFitos(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	file := filepath.Join(dir, "f.txt")
	if err := os.WriteFile(file, []byte("abc"), 0o600); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(dir, "missing")
	const format = "{name: %s, size: %d, mode: %s, mode_int: %d, is_dir: %v}"

	tests := []struct {
		fi   os.FileInfo
		name string
		path string
		want string
	}{
		{
			// a given fi must be used as-is: the (nonexistent) path is never stat'ed,
			// otherwise the stat failure would produce the "unknown file info" string.
			name: "given fi is used without stat",
			path: missing,
			fi:   fitosFakeFileInfo{name: "fake", size: 7, mode: 0o640},
			want: fmt.Sprintf(format, "fake", 7, fs.FileMode(0o640).String(), fs.FileMode(0o640), false),
		},
		{
			name: "given dir fi",
			path: missing,
			fi:   fitosFakeFileInfo{name: "d", size: 0, mode: fs.ModeDir | 0o755},
			want: fmt.Sprintf(format, "d", 0, (fs.ModeDir | 0o755).String(), fs.ModeDir|0o755, true),
		},
		{
			name: "nil fi and existing path is stat'ed",
			path: file,
			want: fmt.Sprintf(format, "f.txt", 3, fs.FileMode(0o600).String(), fs.FileMode(0o600), false),
		},
		{
			name: "nil fi and missing path",
			path: missing,
			want: fmt.Sprintf("unknown file info: for %s\t%v", missing, os.FileInfo(nil)),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := fitos(tc.path, tc.fi); got != tc.want {
				t.Errorf("fitos() = %q, want %q", got, tc.want)
			}
		})
	}
}
