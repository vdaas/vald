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

package file

import (
	"path/filepath"
	"slices"
	"syscall"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
)

// crossFSDir returns a directory on a filesystem other than t.TempDir()'s, or
// skips the test when none is available (rename between them must be EXDEV).
func crossFSDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("/dev/shm", "vald-movedir-")
	if err != nil {
		t.Skipf("/dev/shm unavailable: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			t.Logf("cleanup %s: %v", dir, err)
		}
	})
	probe := filepath.Join(dir, "probe")
	if err := os.Mkdir(probe, 0o755); err != nil {
		t.Skipf("cannot create probe dir: %v", err)
	}
	target := filepath.Join(t.TempDir(), "probe")
	if err := os.Rename(probe, target); !errors.Is(err, syscall.EXDEV) {
		t.Skipf("rename between /dev/shm and %s is not EXDEV (err=%v)", target, err)
	}
	return dir
}

func entryNames(t *testing.T, dir string) []string {
	t.Helper()
	des, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	names := make([]string, 0, len(des))
	for _, de := range des {
		names = append(names, de.Name())
	}
	slices.Sort(names)
	return names
}

// MoveDir's cross-filesystem fallback used to CopyDir straight into dst, so a
// failed (or interrupted) copy left a partially populated dst. It now copies
// into a staging sibling and renames it into place.
func TestMoveDir_crossFilesystemCopyFallback(t *testing.T) {
	t.Parallel()
	tests := []struct {
		files      map[string]string // src content; an empty value creates a FIFO, which CopyFileWithPerm rejects
		name       string
		oldDst     string   // existing dst/a content; "" means dst does not exist
		wantA      string   // dst/a content afterwards; "" means dst must not exist
		wantEntry  []string // entries left in dst's parent
		wantErr    bool
		staleStage bool // leave a staging dir from an attempt killed mid-copy
	}{
		{
			name:      "copy succeeds over an existing dst",
			files:     map[string]string{"a": newContent},
			oldDst:    oldContent,
			wantEntry: []string{dstDirName},
			wantA:     newContent,
		},
		{
			name:      "copy succeeds without a dst",
			files:     map[string]string{"a": newContent},
			wantEntry: []string{dstDirName},
			wantA:     newContent,
		},
		{
			name:      "copy fails without a dst: dst must not appear",
			files:     map[string]string{"a": newContent, "empty": ""},
			wantErr:   true,
			wantEntry: []string{},
		},
		{
			name:      "copy fails over an existing dst: old dst is rolled back",
			files:     map[string]string{"a": newContent, "empty": ""},
			oldDst:    oldContent,
			wantErr:   true,
			wantEntry: []string{dstDirName},
			wantA:     oldContent,
		},
		{
			name:       "stale stage of a killed attempt is removed and not merged into dst",
			files:      map[string]string{"a": newContent},
			oldDst:     oldContent,
			staleStage: true,
			wantEntry:  []string{dstDirName},
			wantA:      newContent,
		},
		{
			name:       "stale stage is removed even when the copy fails",
			files:      map[string]string{"a": newContent, "blank": ""},
			staleStage: true,
			wantErr:    true,
			wantEntry:  []string{},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			src := filepath.Join(crossFSDir(t), "src")
			if err := os.MkdirAll(src, 0o755); err != nil {
				t.Fatal(err)
			}
			for name, content := range tc.files {
				p := filepath.Join(src, name)
				if content == "" {
					if err := syscall.Mkfifo(p, 0o600); err != nil {
						t.Fatal(err)
					}
					continue
				}
				if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			parent := t.TempDir()
			dst := filepath.Join(parent, dstDirName)
			if tc.oldDst != "" {
				if err := os.MkdirAll(dst, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(dst, "a"), []byte(tc.oldDst), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			if tc.staleStage {
				stage := filepath.Join(parent, "tmp-copy-dst")
				if err := os.MkdirAll(stage, 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(stage, "partial"), []byte("stale"), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			err := MoveDir(t.Context(), src, dst)
			if (err != nil) != tc.wantErr {
				t.Fatalf("MoveDir() error = %v, wantErr %v", err, tc.wantErr)
			}
			if got := entryNames(t, parent); !slices.Equal(got, tc.wantEntry) {
				t.Errorf("entries in dst parent = %v, want %v", got, tc.wantEntry)
			}
			if tc.wantA == "" {
				if Exists(dst) {
					t.Errorf("dst %s exists after a failed copy", dst)
				}
				return
			}
			got, err := os.ReadFile(filepath.Join(dst, "a"))
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tc.wantA {
				t.Errorf("dst/a = %q, want %q", got, tc.wantA)
			}
			if Exists(filepath.Join(dst, "partial")) {
				t.Errorf("stale staging content was merged into dst %s", dst)
			}
			if !tc.wantErr && Exists(src) {
				t.Errorf("src %s still exists after a successful move", src)
			}
		})
	}
}

// TestMoveDir_crossFilesystemEmptySrc pins that an empty src is still moved
// across filesystems: CopyDir creates nothing for it, so the staging dir must
// exist before it is renamed into place.
func TestMoveDir_crossFilesystemEmptySrc(t *testing.T) {
	t.Parallel()
	src := filepath.Join(crossFSDir(t), "src")
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}
	parent := t.TempDir()
	dst := filepath.Join(parent, dstDirName)

	if err := MoveDir(t.Context(), src, dst); err != nil {
		t.Fatalf("MoveDir() error = %v, want nil", err)
	}
	if fi, err := os.Stat(dst); err != nil || !fi.IsDir() {
		t.Fatalf("dst %s is not a directory after the move: %v", dst, err)
	}
	if got := entryNames(t, parent); !slices.Equal(got, []string{dstDirName}) {
		t.Errorf("entries in dst parent = %v, want [dst]", got)
	}
	if Exists(src) {
		t.Errorf("src %s still exists after a successful move", src)
	}
}
