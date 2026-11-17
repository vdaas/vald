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
	"io/fs"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"syscall"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/test/goleak"
)

func TestOpen(t *testing.T) {
	type args struct {
		path string
		flg  int
		perm os.FileMode
	}
	type want struct {
		want *os.File
		err  error
	}
	type test struct {
		want       want
		checkFunc  func(want, *os.File, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args, *os.File)
		name       string
		args       args
	}

	defaultCheckFunc := func(w want, got *os.File, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}

		if w.want == nil {
			if got != nil {
				return errors.New("got is not nil")
			}
		} else {
			if got, want := got.Name(), w.want.Name(); got != want {
				return errors.Errorf("got name = %s, want: %s")
			}
		}
		return nil
	}
	defaultAfterFunc := func(t *testing.T, args args, f *os.File) {
		t.Helper()

		if f != nil {
			if err := f.Close(); err != nil {
				t.Error(err)
			}
		}
	}
	tests := []test{
		{
			name: "returns *os.File when path is `test/data`",
			args: args{
				path: "test/data",
				flg:  os.O_CREATE,
				perm: fs.ModePerm,
			},
			checkFunc: func(_ want, got *os.File, gotErr error) error {
				file, err := os.OpenFile("test/data", os.O_CREATE, fs.ModePerm)
				if err != nil {
					return err
				}
				return defaultCheckFunc(want{
					want: file,
					err:  nil,
				}, got, gotErr)
			},
			afterFunc: func(t *testing.T, args args, f *os.File) {
				t.Helper()
				defaultAfterFunc(t, args, f)
				if err := os.RemoveAll("test"); err != nil {
					t.Fatal(err)
				}
			},
		},

		{
			name: "returns *os.File when path is `test/test/data`",
			args: args{
				path: "test/test/data",
				flg:  os.O_CREATE,
				perm: fs.ModePerm,
			},
			checkFunc: func(_ want, got *os.File, gotErr error) error {
				file, err := os.OpenFile("test/test/data", os.O_CREATE, fs.ModePerm)
				if err != nil {
					return err
				}
				return defaultCheckFunc(want{
					want: file,
					err:  nil,
				}, got, gotErr)
			},
			afterFunc: func(t *testing.T, args args, f *os.File) {
				t.Helper()
				defaultAfterFunc(t, args, f)
				if err := os.RemoveAll("test"); err != nil {
					t.Fatal(err)
				}
			},
		},

		{
			name: "returns *os.File when path is `file.go`",
			args: args{
				path: "file.go",
				flg:  os.O_RDONLY,
				perm: fs.ModePerm,
			},
			want: want{
				want: func() *os.File {
					f, _ := os.OpenFile("file.go", os.O_RDONLY, fs.ModePerm)
					return f
				}(),
				err: nil,
			},
		},

		{
			name: "returns (nil, error) when path is empty",
			args: args{
				flg:  os.O_CREATE,
				perm: fs.ModeDir,
			},
			want: want{
				want: nil,
				err:  errors.ErrPathNotSpecified,
			},
		},

		{
			name: "returns (nil, error) when file does not exists and flag is not CREATE or APPEND",
			args: args{
				path: "dummy",
				flg:  os.O_RDONLY,
				perm: fs.ModePerm,
			},
			want: want{
				want: nil,
				err: &os.PathError{
					Op:   "open",
					Path: "dummy",
					Err: func() error {
						_, err := syscall.Open("dummy", syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
						return err
					}(),
				},
			},
		},

		{
			name: "returns (nil, error) when the folder does not exists and flag is not CREATE or APPEND",
			args: args{
				path: "dummy/dummy",
				flg:  os.O_RDONLY,
				perm: fs.ModePerm,
			},
			want: want{
				want: nil,
				err: &os.PathError{
					Op:   "open",
					Path: "dummy/dummy",
					Err: func() error {
						_, err := syscall.Open("dummy/", syscall.O_RDONLY|syscall.O_CLOEXEC, 0)
						return err
					}(),
				},
			},
			afterFunc: func(t *testing.T, args args, f *os.File) {
				t.Helper()
				defaultAfterFunc(t, args, f)
				if err := os.RemoveAll("dummy"); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt)
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}
			if test.afterFunc == nil {
				test.afterFunc = defaultAfterFunc
			}

			got, err := Open(test.args.path, test.args.flg, test.args.perm)
			defer test.afterFunc(tt, test.args, got)

			if err := checkFunc(test.want, got, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestExists(t *testing.T) {
	type args struct {
		path string
	}
	type want struct {
		want bool
	}
	type test struct {
		checkFunc  func(want, bool) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
		want       want
	}
	defaultCheckFunc := func(w want, got bool) error {
		if !reflect.DeepEqual(got, w.want) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
		}
		return nil
	}

	const (
		baseDir      = "./utiltest"
		testDirPath  = baseDir + "/index"
		testFilePath = baseDir + "/ngt-meta.kvsdb"
		testSym      = "sym"

		testFailsDirPath  = baseDir + "/fails-index"
		testFailsFilePath = baseDir + "/ngt-meta-fails.kvsdb"
	)

	defaultAfterFunc := func(t *testing.T, args args) {
		t.Helper()
		if err := os.RemoveAll(baseDir); err != nil {
			t.Error(err)
		}
		if err := os.RemoveAll(args.path); err != nil {
			t.Error(err)
		}
	}

	tests := []test{
		{
			name: "return true when the directory exists",
			args: args{
				path: testDirPath,
			},
			beforeFunc: func(t *testing.T, args args) {
				t.Helper()
				if err := os.MkdirAll(args.path, fs.ModePerm); err != nil {
					t.Fatal(err)
				}
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when the directory exists and the type is symbolic link",
			args: args{
				path: testSym,
			},
			beforeFunc: func(t *testing.T, args args) {
				t.Helper()
				if err := os.MkdirAll(testDirPath, fs.ModePerm); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(testDirPath, testSym); err != nil {
					t.Error(err)
				}
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when the file exists",
			args: args{
				path: testFilePath,
			},
			beforeFunc: func(t *testing.T, args args) {
				t.Helper()
				if err := os.MkdirAll(baseDir, fs.ModePerm); err != nil {
					t.Fatal(err)
				}

				f, err := os.Create(args.path)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := f.Close(); err != nil {
						t.Error(err)
					}
				}()
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return true when the file exists and the type is symbolic link",
			args: args{
				path: testSym,
			},
			beforeFunc: func(t *testing.T, args args) {
				t.Helper()
				if err := os.MkdirAll(baseDir, fs.ModePerm); err != nil {
					t.Fatal(err)
				}

				f, err := os.Create(testFilePath)
				if err != nil {
					t.Fatal(err)
				}
				defer func() {
					if err := f.Close(); err != nil {
						t.Error(err)
					}
				}()

				if err := os.Symlink(testFilePath, testSym); err != nil {
					t.Error(err)
				}
			},
			want: want{
				want: true,
			},
		},
		{
			name: "return false when the directory does not exist",
			args: args{
				path: testFailsDirPath,
			},
			want: want{
				want: false,
			},
		},
		{
			name: "return false when the file exists",
			args: args{
				path: testFailsFilePath,
			},
			want: want{
				want: false,
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc == nil {
				test.afterFunc = defaultAfterFunc
			}
			defer test.afterFunc(tt, test.args)

			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			// TODO we have to check more patter about file or dir
			got := Exists(test.args.path)
			if err := checkFunc(test.want, got); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

func TestExistsWithDetail(t *testing.T) {
	type args struct {
		path string
	}
	type want struct {
		wantFi fs.FileInfo
		err    error
		wantE  bool
	}
	type test struct {
		want       want
		checkFunc  func(want, bool, fs.FileInfo, error) error
		beforeFunc func(*testing.T, args)
		afterFunc  func(*testing.T, args)
		name       string
		args       args
	}
	defaultCheckFunc := func(w want, gotE bool, gotFi fs.FileInfo, err error) error {
		if !errors.Is(err, w.err) {
			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
		}
		if !reflect.DeepEqual(gotE, w.wantE) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotE, w.wantE)
		}
		if !reflect.DeepEqual(gotFi, w.wantFi) {
			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotFi, w.wantFi)
		}
		return nil
	}
	tests := []test{
		{
			name: "return fs.ErrInvalid when path is empty",
			args: args{
				path: "",
			},
			want: want{
				wantE:  false,
				wantFi: nil,
				err:    fs.ErrInvalid,
			},
			checkFunc: defaultCheckFunc,
			beforeFunc: func(t *testing.T, args args) {
				t.Helper()
			},
			afterFunc: func(t *testing.T, args args) {
				t.Helper()
			},
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.name, func(tt *testing.T) {
			tt.Parallel()
			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
			if test.beforeFunc != nil {
				test.beforeFunc(tt, test.args)
			}
			if test.afterFunc != nil {
				defer test.afterFunc(tt, test.args)
			}
			checkFunc := test.checkFunc
			if test.checkFunc == nil {
				checkFunc = defaultCheckFunc
			}

			gotE, gotFi, err := ExistsWithDetail(test.args.path)
			if err := checkFunc(test.want, gotE, gotFi, err); err != nil {
				tt.Errorf("error = %v", err)
			}
		})
	}
}

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

// NOT IMPLEMENTED BELOW
//
// func TestMoveDir(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		src string
// 		dst string
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           src:"",
// 		           dst:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           src:"",
// 		           dst:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			err := MoveDir(test.args.ctx, test.args.src, test.args.dst)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_doMoveDir(t *testing.T) {
// 	type args struct {
// 		ctx      context.Context
// 		src      string
// 		dst      string
// 		rollback bool
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           src:"",
// 		           dst:"",
// 		           rollback:false,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           src:"",
// 		           dst:"",
// 		           rollback:false,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			err := doMoveDir(test.args.ctx, test.args.src, test.args.dst, test.args.rollback)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestCopyDir(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		src string
// 		dst string
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           src:"",
// 		           dst:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           src:"",
// 		           dst:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			err := CopyDir(test.args.ctx, test.args.src, test.args.dst)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestCopyFile(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 		src string
// 		dst string
// 	}
// 	type want struct {
// 		wantN int64
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, int64, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotN int64, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotN, w.wantN) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           src:"",
// 		           dst:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           src:"",
// 		           dst:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotN, err := CopyFile(test.args.ctx, test.args.src, test.args.dst)
// 			if err := checkFunc(test.want, gotN, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestCopyFileWithPerm(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		src  string
// 		dst  string
// 		perm fs.FileMode
// 	}
// 	type want struct {
// 		wantN int64
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, int64, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotN int64, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotN, w.wantN) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           src:"",
// 		           dst:"",
// 		           perm:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           src:"",
// 		           dst:"",
// 		           perm:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotN, err := CopyFileWithPerm(test.args.ctx, test.args.src, test.args.dst, test.args.perm)
// 			if err := checkFunc(test.want, gotN, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestWriteFile(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		target string
// 		r      io.Reader
// 		perm   fs.FileMode
// 	}
// 	type want struct {
// 		wantN int64
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, int64, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotN int64, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotN, w.wantN) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           target:"",
// 		           r:nil,
// 		           perm:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           target:"",
// 		           r:nil,
// 		           perm:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotN, err := WriteFile(test.args.ctx, test.args.target, test.args.r, test.args.perm)
// 			if err := checkFunc(test.want, gotN, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestOverWriteFile(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		target string
// 		r      io.Reader
// 		perm   fs.FileMode
// 	}
// 	type want struct {
// 		wantN int64
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, int64, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotN int64, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotN, w.wantN) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           target:"",
// 		           r:nil,
// 		           perm:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           target:"",
// 		           r:nil,
// 		           perm:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotN, err := OverWriteFile(test.args.ctx, test.args.target, test.args.r, test.args.perm)
// 			if err := checkFunc(test.want, gotN, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestAppendFile(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		target string
// 		r      io.Reader
// 		perm   fs.FileMode
// 	}
// 	type want struct {
// 		wantN int64
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, int64, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotN int64, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotN, w.wantN) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           target:"",
// 		           r:nil,
// 		           perm:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           target:"",
// 		           r:nil,
// 		           perm:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotN, err := AppendFile(test.args.ctx, test.args.target, test.args.r, test.args.perm)
// 			if err := checkFunc(test.want, gotN, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_writeFileWithContext(t *testing.T) {
// 	type args struct {
// 		ctx    context.Context
// 		target string
// 		r      io.Reader
// 		flg    int
// 		perm   fs.FileMode
// 	}
// 	type want struct {
// 		wantN int64
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, int64, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotN int64, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotN, w.wantN) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           target:"",
// 		           r:nil,
// 		           flg:0,
// 		           perm:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           target:"",
// 		           r:nil,
// 		           flg:0,
// 		           perm:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotN, err := writeFileWithContext(test.args.ctx, test.args.target, test.args.r, test.args.flg, test.args.perm)
// 			if err := checkFunc(test.want, gotN, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestReadDir(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	type want struct {
// 		wantDirs []fs.DirEntry
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []fs.DirEntry, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotDirs []fs.DirEntry, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotDirs, w.wantDirs) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotDirs, w.wantDirs)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           path:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           path:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotDirs, err := ReadDir(test.args.path)
// 			if err := checkFunc(test.want, gotDirs, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestReadFile(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	type want struct {
// 		wantN []byte
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []byte, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotN []byte, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotN, w.wantN) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotN, w.wantN)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           path:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           path:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotN, err := ReadFile(test.args.path)
// 			if err := checkFunc(test.want, gotN, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestAbsolutePath(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	type want struct {
// 		want string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got string) error {
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           path:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           path:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got := AbsolutePath(test.args.path)
// 			if err := checkFunc(test.want, got); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestMkdirAll(t *testing.T) {
// 	type args struct {
// 		path string
// 		perm fs.FileMode
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           path:"",
// 		           perm:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           path:"",
// 		           perm:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			err := MkdirAll(test.args.path, test.args.perm)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestMkdirTemp(t *testing.T) {
// 	type args struct {
// 		baseDir string
// 	}
// 	type want struct {
// 		wantPath string
// 		err      error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotPath string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotPath, w.wantPath) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPath, w.wantPath)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           baseDir:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           baseDir:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotPath, err := MkdirTemp(test.args.baseDir)
// 			if err := checkFunc(test.want, gotPath, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestCreateTemp(t *testing.T) {
// 	type args struct {
// 		baseDir string
// 	}
// 	type want struct {
// 		wantF *os.File
// 		err   error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, *os.File, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotF *os.File, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(gotF, w.wantF) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotF, w.wantF)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           baseDir:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           baseDir:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotF, err := CreateTemp(test.args.baseDir)
// 			if err := checkFunc(test.want, gotF, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestListInDir(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	type want struct {
// 		want []string
// 		err  error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, []string, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, got []string, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		if !reflect.DeepEqual(got, w.want) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", got, w.want)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           path:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           path:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			got, err := ListInDir(test.args.path)
// 			if err := checkFunc(test.want, got, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestDeleteDir(t *testing.T) {
// 	type args struct {
// 		ctx  context.Context
// 		path string
// 	}
// 	type want struct {
// 		err error
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, error) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, err error) error {
// 		if !errors.Is(err, w.err) {
// 			return errors.Errorf("got_error: \"%#v\",\n\t\t\t\twant: \"%#v\"", err, w.err)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           ctx:nil,
// 		           path:"",
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           ctx:nil,
// 		           path:"",
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			err := DeleteDir(test.args.ctx, test.args.path)
// 			if err := checkFunc(test.want, err); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func TestJoin(t *testing.T) {
// 	type args struct {
// 		paths []string
// 	}
// 	type want struct {
// 		wantPath string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotPath string) error {
// 		if !reflect.DeepEqual(gotPath, w.wantPath) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPath, w.wantPath)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           paths:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           paths:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotPath := Join(test.args.paths...)
// 			if err := checkFunc(test.want, gotPath); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }
//
// func Test_joinFilePaths(t *testing.T) {
// 	type args struct {
// 		paths []string
// 	}
// 	type want struct {
// 		wantPath string
// 	}
// 	type test struct {
// 		name       string
// 		args       args
// 		want       want
// 		checkFunc  func(want, string) error
// 		beforeFunc func(*testing.T, args)
// 		afterFunc  func(*testing.T, args)
// 	}
// 	defaultCheckFunc := func(w want, gotPath string) error {
// 		if !reflect.DeepEqual(gotPath, w.wantPath) {
// 			return errors.Errorf("got: \"%#v\",\n\t\t\t\twant: \"%#v\"", gotPath, w.wantPath)
// 		}
// 		return nil
// 	}
// 	tests := []test{
// 		// TODO test cases
// 		/*
// 		   {
// 		       name: "test_case_1",
// 		       args: args {
// 		           paths:nil,
// 		       },
// 		       want: want{},
// 		       checkFunc: defaultCheckFunc,
// 		       beforeFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		       afterFunc: func(t *testing.T, args args) {
// 		           t.Helper()
// 		       },
// 		   },
// 		*/
//
// 		// TODO test cases
// 		/*
// 		   func() test {
// 		       return test {
// 		           name: "test_case_2",
// 		           args: args {
// 		           paths:nil,
// 		           },
// 		           want: want{},
// 		           checkFunc: defaultCheckFunc,
// 		           beforeFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		           afterFunc: func(t *testing.T, args args) {
// 		               t.Helper()
// 		           },
// 		       }
// 		   }(),
// 		*/
// 	}
//
// 	for _, tc := range tests {
// 		test := tc
// 		t.Run(test.name, func(tt *testing.T) {
// 			tt.Parallel()
// 			defer goleak.VerifyNone(tt, goleak.IgnoreCurrent())
// 			if test.beforeFunc != nil {
// 				test.beforeFunc(tt, test.args)
// 			}
// 			if test.afterFunc != nil {
// 				defer test.afterFunc(tt, test.args)
// 			}
// 			checkFunc := test.checkFunc
// 			if test.checkFunc == nil {
// 				checkFunc = defaultCheckFunc
// 			}
//
// 			gotPath := joinFilePaths(test.args.paths...)
// 			if err := checkFunc(test.want, gotPath); err != nil {
// 				tt.Errorf("error = %v", err)
// 			}
// 		})
// 	}
// }

func TestOpen_flagHandling(t *testing.T) {
	t.Parallel()
	type test struct {
		wantErr error
		setup   func(t *testing.T, root string)
		check   func(t *testing.T, root string)
		name    string
		path    string
		flg     int
	}
	tests := []test{
		{
			name: "O_CREATE|O_EXCL on a new path in a missing directory succeeds",
			path: filepath.Join("a", "b", "new"),
			flg:  os.O_CREATE | os.O_EXCL | os.O_RDWR | os.O_TRUNC,
			check: func(t *testing.T, root string) {
				t.Helper()
				if !Exists(filepath.Join(root, "a", "b", "new")) {
					t.Error("the file must be created")
				}
			},
		},
		{
			name: "O_CREATE|O_EXCL on an existing path fails with ErrExist and keeps the content",
			path: "exist",
			flg:  os.O_CREATE | os.O_EXCL | os.O_RDWR,
			setup: func(t *testing.T, root string) {
				t.Helper()
				if err := os.WriteFile(filepath.Join(root, "exist"), []byte("keep"), fs.ModePerm); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: fs.ErrExist,
			check: func(t *testing.T, root string) {
				t.Helper()
				b, err := os.ReadFile(filepath.Join(root, "exist"))
				if err != nil {
					t.Fatal(err)
				}
				if string(b) != "keep" {
					t.Errorf("content = %q, want %q", b, "keep")
				}
			},
		},
		{
			name: "O_APPEND without O_CREATE still creates a missing file",
			path: filepath.Join("dir", "append"),
			flg:  os.O_APPEND | os.O_WRONLY,
			check: func(t *testing.T, root string) {
				t.Helper()
				if !Exists(filepath.Join(root, "dir", "append")) {
					t.Error("the file must be created")
				}
			},
		},
		{
			name:    "O_RDONLY on a missing path does not create the parent directory",
			path:    filepath.Join("missing", "file"),
			flg:     os.O_RDONLY,
			wantErr: fs.ErrNotExist,
			check: func(t *testing.T, root string) {
				t.Helper()
				if Exists(filepath.Join(root, "missing")) {
					t.Error("the parent directory must not be created by a read-only open")
				}
			},
		},
		{
			name:    "O_WRONLY without O_CREATE on a missing path does not create the parent directory",
			path:    filepath.Join("missing", "file"),
			flg:     os.O_WRONLY,
			wantErr: fs.ErrNotExist,
			check: func(t *testing.T, root string) {
				t.Helper()
				if Exists(filepath.Join(root, "missing")) {
					t.Error("the parent directory must not be created without O_CREATE or O_APPEND")
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := t.TempDir()
			if tc.setup != nil {
				tc.setup(t, root)
			}
			f, err := Open(filepath.Join(root, tc.path), tc.flg, fs.ModePerm)
			if f != nil {
				defer func() {
					if cerr := f.Close(); cerr != nil {
						t.Error(cerr)
					}
				}()
			}
			if tc.wantErr == nil {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			} else if !errors.Is(err, tc.wantErr) {
				t.Errorf("error = %v, want %v", err, tc.wantErr)
			}
			tc.check(t, root)
		})
	}
}

func TestCreateTemp_createsExactlyOneFile(t *testing.T) {
	t.Parallel()
	f, err := CreateTemp(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			t.Error(cerr)
		}
	}()
	entries, err := os.ReadDir(filepath.Dir(f.Name()))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("got %d files in the temporary directory, want 1", len(entries))
	}
}

func TestMkdirAll_doesNotDeleteExistingData(t *testing.T) {
	t.Parallel()
	type test struct {
		wantErr error
		setup   func(t *testing.T, path string)
		check   func(t *testing.T, path string)
		name    string
	}
	readKeep := func(t *testing.T, name string) {
		t.Helper()
		b, err := os.ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "keep" {
			t.Errorf("content = %q, want %q", b, "keep")
		}
	}
	tests := []test{
		{
			name: "existing regular file is not removed and an error is returned",
			setup: func(t *testing.T, path string) {
				t.Helper()
				if err := os.WriteFile(path, []byte("keep"), fs.ModePerm); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: fs.ErrExist,
			check: func(t *testing.T, path string) {
				t.Helper()
				readKeep(t, path)
			},
		},
		{
			name: "existing directory keeps its contents",
			setup: func(t *testing.T, path string) {
				t.Helper()
				if err := os.MkdirAll(path, fs.ModePerm); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(path, "child"), []byte("keep"), fs.ModePerm); err != nil {
					t.Fatal(err)
				}
			},
			check: func(t *testing.T, path string) {
				t.Helper()
				readKeep(t, filepath.Join(path, "child"))
			},
		},
		{
			name: "existing directory without permission keeps its contents",
			setup: func(t *testing.T, path string) {
				t.Helper()
				if os.Geteuid() == 0 {
					t.Skip("permission bits are not enforced for root")
				}
				if err := os.MkdirAll(path, fs.ModePerm); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(path, "child"), []byte("keep"), fs.ModePerm); err != nil {
					t.Fatal(err)
				}
				if err := os.Chmod(path, 0); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if err := os.Chmod(path, fs.ModePerm); err != nil {
						t.Error(err)
					}
				})
			},
			check: func(t *testing.T, path string) {
				t.Helper()
				if err := os.Chmod(path, fs.ModePerm); err != nil {
					t.Fatal(err)
				}
				readKeep(t, filepath.Join(path, "child"))
			},
		},
		{
			name: "missing nested directory is created",
			check: func(t *testing.T, path string) {
				t.Helper()
				fi, err := os.Stat(path)
				if err != nil {
					t.Fatal(err)
				}
				if !fi.IsDir() {
					t.Error("the path must be a directory")
				}
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path := filepath.Join(t.TempDir(), "x", "target")
			if tc.setup != nil {
				// setup creates the target itself, so make sure the parent exists
				if err := os.MkdirAll(filepath.Dir(path), fs.ModePerm); err != nil {
					t.Fatal(err)
				}
				tc.setup(t, path)
			}
			err := MkdirAll(path, fs.ModePerm)
			if tc.wantErr == nil {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			} else if !errors.Is(err, tc.wantErr) {
				t.Errorf("error = %v, want %v", err, tc.wantErr)
			}
			tc.check(t, path)
		})
	}
}
