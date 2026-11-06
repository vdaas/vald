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
	"strconv"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
)

// Fixture values shared by the CopyDir/MoveDir tests of this package: the dst
// directory name, the content of a file in the moved src (newContent) and in a
// pre-existing dst (oldContent), and a file kept at the top of a copied tree.
const (
	dstDirName = "dst"
	newContent = "new"
	oldContent = "old"
	topFile    = "top"
)

// DeleteDir used to remove every walked entry concurrently, so a parent's Remove
// raced its children and failed with ENOTEMPTY, leaving the tree behind.
func TestDeleteDir_nestedTree(t *testing.T) {
	t.Parallel()
	root := filepath.Join(t.TempDir(), "root")
	for _, d := range []string{"a/b/c", "a/d", "e"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{"a/b/c/f1", "a/d/f2", "e/f3", "f4"} {
		if err := os.WriteFile(filepath.Join(root, f), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := DeleteDir(t.Context(), root); err != nil {
		t.Fatalf("DeleteDir() error = %v", err)
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatalf("root still exists after DeleteDir: %v", err)
	}
}

// CopyDir used to flatten every file onto dst/<basename> (so src/a/x and src/b/x
// collided) and to create files and directories with the permission-less
// DirEntry.Type() mode.
func TestCopyDir_keepsLayoutAndPermissions(t *testing.T) {
	t.Parallel()
	src := filepath.Join(t.TempDir(), "src")
	dst := filepath.Join(src, "origin")
	files := map[string]string{"a/x": "ax", "b/x": "bx", topFile: topFile}
	for name, content := range files {
		p := filepath.Join(src, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o640); err != nil {
			t.Fatal(err)
		}
	}
	if err := CopyDir(t.Context(), src, dst); err != nil {
		t.Fatalf("CopyDir() error = %v", err)
	}
	for name, content := range files {
		p := filepath.Join(dst, name)
		got, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		if string(got) != content {
			t.Errorf("%s = %q, want %q", p, got, content)
		}
		fi, err := os.Stat(p)
		if err != nil {
			t.Fatal(err)
		}
		if perm := fi.Mode().Perm(); perm&0o600 != 0o600 {
			t.Errorf("%s perm = %v, want owner read/write", p, perm)
		}
	}
	fi, err := os.Stat(filepath.Join(dst, "a"))
	if err != nil {
		t.Fatal(err)
	}
	if perm := fi.Mode().Perm(); perm&0o700 != 0o700 {
		t.Errorf("dir perm = %v, want owner rwx", perm)
	}
}

// CopyDir skipped every path with dst as a plain string prefix, so with dst=src/origin
// the sibling src/origin2 was silently not copied.
func TestCopyDir_copiesSiblingsSharingDstPrefix(t *testing.T) {
	t.Parallel()
	src := filepath.Join(t.TempDir(), "src")
	dst := filepath.Join(src, "origin")
	files := map[string]string{"origin2/x": "o2x", "originfile": "of", topFile: topFile}
	for name, content := range files {
		p := filepath.Join(src, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// a pre-existing dst with content must not be copied into itself.
	if err := os.MkdirAll(filepath.Join(dst, "old"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dst, "old", "y"), []byte("y"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := CopyDir(t.Context(), src, dst); err != nil {
		t.Fatalf("CopyDir() error = %v", err)
	}
	for name, content := range files {
		p := filepath.Join(dst, name)
		got, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		if string(got) != content {
			t.Errorf("%s = %q, want %q", p, got, content)
		}
	}
	if _, err := os.Stat(filepath.Join(dst, "origin")); !os.IsNotExist(err) {
		t.Errorf("dst was copied into itself: stat error = %v", err)
	}
}

// CopyDir ignored the walk error reported for src itself, so a missing src
// returned nil with nothing copied.
func TestCopyDir_missingSrcReturnsError(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	src := filepath.Join(dir, "missing")
	dst := filepath.Join(dir, dstDirName)
	if err := CopyDir(t.Context(), src, dst); err == nil {
		t.Fatal("CopyDir() error = nil, want an error for a missing src")
	}
}

// CopyFileWithPerm rejected every empty regular file as not found, so CopyDir
// (and MoveDir's copy fallback) failed for any tree containing an empty file.
func TestCopyFile_emptyRegularFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	src := filepath.Join(dir, "src", "empty")
	if err := os.MkdirAll(filepath.Dir(src), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(src, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	dst := filepath.Join(dir, "copy")
	n, err := CopyFile(t.Context(), src, dst)
	if err != nil {
		t.Fatalf("CopyFile() error = %v, want nil", err)
	}
	if n != 0 {
		t.Errorf("CopyFile() n = %d, want 0", n)
	}
	if fi, err := os.Stat(dst); err != nil || !fi.Mode().IsRegular() || fi.Size() != 0 {
		t.Errorf("dst %s is not an empty regular file: fi=%v, err=%v", dst, fi, err)
	}

	dstDir := filepath.Join(dir, "dstdir")
	if err := CopyDir(t.Context(), filepath.Dir(src), dstDir); err != nil {
		t.Fatalf("CopyDir() error = %v, want nil", err)
	}
	if fi, err := os.Stat(filepath.Join(dstDir, "empty")); err != nil || fi.Size() != 0 {
		t.Errorf("empty file was not copied by CopyDir: fi=%v, err=%v", fi, err)
	}
}

// doMoveDir removed the moved-aside dst (tmpPath) in an unconditional defer, so
// when the rollback after a failed move also failed, the only copy of the
// previous dst was deleted. Not parallel: it swaps the package-level rename.
func TestMoveDir_failedRollbackKeepsPreviousDst(t *testing.T) {
	errInjected := errors.New("injected rename failure")
	tests := []struct {
		// fail reports whether rename(oldpath, newpath) must fail.
		fail     func(oldpath, newpath, dst string) bool
		name     string
		wantA    string // dst/a afterwards; "" means dst must not exist
		wantTmpA string // a in the kept tmp-<nanos> dir; "" means none is kept
		wantEnts int    // entries in dst's parent
	}{
		{
			name: "rollback succeeds: previous dst restored and tmp removed",
			fail: func(oldpath, newpath, dst string) bool {
				// only the forward moves into dst (from src or the staging dir) fail;
				// the rollback rename of tmp-<nanos> back to dst works.
				base := filepath.Base(oldpath)
				return newpath == dst && (!strings.HasPrefix(base, "tmp-") || base == "tmp-copy-dst")
			},
			wantA:    oldContent,
			wantEnts: 1,
		},
		{
			name: "rollback fails: previous dst is kept in tmp",
			fail: func(_, newpath, dst string) bool {
				return newpath == dst
			},
			wantTmpA: oldContent,
			wantEnts: 1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			src := filepath.Join(t.TempDir(), "src")
			parent := t.TempDir()
			dst := filepath.Join(parent, dstDirName)
			for p, content := range map[string]string{
				filepath.Join(src, "a"): newContent,
				filepath.Join(dst, "a"): oldContent,
			} {
				if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			orig := rename
			t.Cleanup(func() { rename = orig })
			rename = func(oldpath, newpath string) error {
				if tc.fail(oldpath, newpath, dst) {
					return errInjected
				}
				return orig(oldpath, newpath)
			}

			if err := MoveDir(t.Context(), src, dst); err == nil {
				t.Fatal("MoveDir() error = nil, want the injected failure")
			}
			ents := entryNames(t, parent)
			if len(ents) != tc.wantEnts {
				t.Fatalf("entries in dst parent = %v, want %d entries", ents, tc.wantEnts)
			}
			if tc.wantA != "" {
				got, err := os.ReadFile(filepath.Join(dst, "a"))
				if err != nil {
					t.Fatal(err)
				}
				if string(got) != tc.wantA {
					t.Errorf("dst/a = %q, want %q", got, tc.wantA)
				}
			} else if Exists(dst) {
				t.Errorf("dst %s exists, want it absent after the failed rollback", dst)
			}
			if tc.wantTmpA != "" {
				name := ents[0]
				if !strings.HasPrefix(name, "tmp-") || strings.HasPrefix(name, "tmp-copy-") {
					t.Fatalf("kept entry = %q, want the tmp-<nanos> dir holding the previous dst", name)
				}
				got, err := os.ReadFile(filepath.Join(parent, name, "a"))
				if err != nil {
					t.Fatalf("previous dst content lost: %v", err)
				}
				if string(got) != tc.wantTmpA {
					t.Errorf("%s/a = %q, want %q", name, got, tc.wantTmpA)
				}
			}
			if got, err := os.ReadFile(filepath.Join(src, "a")); err != nil || string(got) != newContent {
				t.Errorf("src/a = %q (err=%v), want src untouched after a failed move", got, err)
			}
		})
	}
}

// doMoveDir named its tmp dir "tmp-<fastime nanos>" and pre-cleaned that name.
// fastime refreshes only every few ms, so the nested rollback
// doMoveDir(tmpPath, dst) could compute its own tmpPath equal to its src and
// delete the only copy of the previous dst. Not parallel: it swaps package-level
// seams.
func TestDoMoveDir_tmpNameCollisionKeepsSrc(t *testing.T) {
	const nanos = 42
	parent := t.TempDir()
	src := filepath.Join(parent, "tmp-"+strconv.Itoa(nanos))
	dst := filepath.Join(parent, dstDirName)
	for p, content := range map[string]string{
		filepath.Join(src, "a"): oldContent,
		filepath.Join(dst, "a"): "busy",
	} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	origNow := unixNanoNow
	origRename := rename
	t.Cleanup(func() {
		unixNanoNow = origNow
		rename = origRename
	})
	unixNanoNow = func() int64 { return nanos }
	errInjected := errors.New("injected rename failure")
	rename = func(oldpath, newpath string) error {
		if oldpath == dst || newpath == dst {
			return errInjected
		}
		return origRename(oldpath, newpath)
	}

	if err := doMoveDir(t.Context(), src, dst, false); err == nil {
		t.Fatal("doMoveDir() error = nil, want the injected failure")
	}
	got, err := os.ReadFile(filepath.Join(src, "a"))
	if err != nil {
		t.Fatalf("src content lost: %v", err)
	}
	if string(got) != oldContent {
		t.Errorf("src/a = %q, want %q", got, oldContent)
	}
}

// writeFileWithContext opened every new file with fs.ModePerm instead of the
// requested perm, so WriteFile/OverWriteFile, and CopyDir through
// CopyFileWithPerm, created files as 0o777&^umask: restored and copied index
// files lost their mode and became world-readable and executable. A missing
// parent must still be created as a traversable directory.
func TestWriteFile_appliesPerm(t *testing.T) {
	t.Parallel()
	// Owner-only modes are not changed by any usual umask, so the result is exact.
	tests := []struct {
		write func(t *testing.T, target string, perm os.FileMode) error
		name  string
		perm  os.FileMode
	}{
		{
			name: "WriteFile",
			perm: 0o600,
			write: func(t *testing.T, target string, perm os.FileMode) error {
				t.Helper()
				_, err := WriteFile(t.Context(), target, strings.NewReader("data"), perm)
				return err
			},
		},
		{
			name: "WriteFile with a read-only mode still writes the new file",
			perm: 0o400,
			write: func(t *testing.T, target string, perm os.FileMode) error {
				t.Helper()
				_, err := WriteFile(t.Context(), target, strings.NewReader("data"), perm)
				return err
			},
		},
		{
			name: "OverWriteFile",
			perm: 0o600,
			write: func(t *testing.T, target string, perm os.FileMode) error {
				t.Helper()
				_, err := OverWriteFile(t.Context(), target, strings.NewReader("data"), perm)
				return err
			},
		},
		{
			name: "CopyDir keeps the source mode",
			perm: 0o600,
			write: func(t *testing.T, target string, perm os.FileMode) error {
				t.Helper()
				src := filepath.Join(t.TempDir(), "src")
				p := filepath.Join(src, filepath.Base(filepath.Dir(target)), filepath.Base(target))
				if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(p, []byte("data"), perm); err != nil {
					t.Fatal(err)
				}
				// os.WriteFile applies the umask; pin the exact source mode.
				if err := os.Chmod(p, perm); err != nil {
					t.Fatal(err)
				}
				return CopyDir(t.Context(), src, filepath.Dir(filepath.Dir(target)))
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(tt *testing.T) {
			tt.Parallel()
			dir := filepath.Join(tt.TempDir(), "parent")
			target := filepath.Join(dir, "f")
			if err := tc.write(tt, target, tc.perm); err != nil {
				tt.Fatalf("write error = %v", err)
			}
			fi, err := os.Stat(target)
			if err != nil {
				tt.Fatal(err)
			}
			if got := fi.Mode().Perm(); got != tc.perm {
				tt.Errorf("file perm = %v, want %v", got, tc.perm)
			}
			if fi.Size() != int64(len("data")) {
				tt.Errorf("file size = %d, want %d", fi.Size(), len("data"))
			}
			di, err := os.Stat(dir)
			if err != nil {
				tt.Fatal(err)
			}
			if perm := di.Mode().Perm(); !di.IsDir() || perm&0o700 != 0o700 {
				tt.Errorf("parent = %v, want a directory with owner rwx", di.Mode())
			}
		})
	}
}
