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
	"cmp"
	"context"
	"io/fs"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/kpango/fastime"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/internal/sync/errgroup"
)

const (
	sep = string(os.PathSeparator)
)

// rename is os.Rename; a variable so tests can inject rename failures into doMoveDir.
//
//nolint:gochecknoglobals // test seam for rename fault injection; only tests reassign it
var rename = os.Rename

// unixNanoNow is fastime.UnixNanoNow; a variable so tests can force doMoveDir's
// tmp name to collide with an existing path.
//
//nolint:gochecknoglobals // test seam for tmp name collisions; only tests reassign it
var unixNanoNow = fastime.UnixNanoNow

// Open opens the file with the given path, flag and permission.
// When flg contains O_CREATE or O_APPEND and the file does not exist,
// the parent folder and the file are created.
func Open(path string, flg int, perm fs.FileMode) (file *os.File, err error) {
	if path == "" {
		return nil, errors.ErrPathNotSpecified
	}

	defer func() {
		if err != nil && file != nil {
			err = errors.Join(file.Close(), err)
			file = nil
		}
	}()
	// The parent directory is created only when the caller asks for the file to be created,
	// so that read-only opens do not have side effects on the file system.
	if _, err = os.Stat(path); err != nil && flg&(os.O_CREATE|os.O_APPEND) != 0 {
		dir := filepath.Dir(path)
		fi, err := os.Stat(dir)
		if err != nil {
			err = MkdirAll(dir, perm)
			if err != nil {
				err = errors.ErrFailedToMkdir(err, dir, fi)
				log.Debug(err)
				return nil, err
			}
		}
		// O_APPEND has always created a missing file, so keep that behavior by adding O_CREATE
		// to the flags of the single OpenFile call below instead of pre-creating the file with
		// O_EXCL, which made the second open fail with EEXIST whenever the caller passed O_EXCL.
		flg |= os.O_CREATE
	}

	file, err = os.OpenFile(path, flg, perm)
	if err != nil {
		err = errors.ErrFailedToOpenFile(err, path, flg, perm)
		log.Warn(err)
		return nil, err
	}

	return file, nil
}

func MoveDir(ctx context.Context, src, dst string) (err error) {
	return doMoveDir(ctx, src, dst, true)
}

func doMoveDir(ctx context.Context, src, dst string, rollback bool) (err error) {
	if len(src) == 0 || len(dst) == 0 || src == dst {
		return nil
	}
	exits, fi, err := ExistsWithDetail(src)
	if !exits || fi == nil || !fi.IsDir() || err != nil {
		return errors.ErrDirectoryNotFound(err, src, fi)
	}

	err = rename(src, dst)
	if err != nil {
		log.Debug(errors.ErrFailedToRenameDir(err, src, dst, nil, nil))
		var (
			tmpPath string
			// keepTmp is set when the rollback failed: tmpPath then holds the only
			// copy of the previous dst and must not be removed.
			keepTmp bool
		)
		exits, fi, err := ExistsWithDetail(dst)
		if exits && fi.IsDir() && err == nil {
			// The nanos come from fastime, refreshed only every few ms, so a nested
			// rollback can compute the same name as the tmp dir it is restoring (its
			// src). Never pre-clean an existing name; pick one that does not exist.
			base := Join(filepath.Dir(dst), "tmp-"+strconv.FormatInt(unixNanoNow(), 10))
			tmpPath = base
			for i := 1; Exists(tmpPath); i++ {
				tmpPath = base + "-" + strconv.Itoa(i)
			}
			err = rename(dst, tmpPath)
			defer func() {
				if keepTmp {
					log.Warnf("keeping temporary path %s for manual recovery of %s", tmpPath, dst)
					return
				}
				if rerr := os.RemoveAll(tmpPath); rerr != nil {
					log.Warnf("failed to remove temporary path %s: %v", tmpPath, rerr)
				}
			}()
			if err != nil {
				log.Debugf("err: %v\t now trying to move file with I/O copy and Remove old index", errors.ErrFailedToRenameDir(err, dst, tmpPath, fi, nil))
				err := CopyDir(ctx, dst, tmpPath)
				if err != nil {
					err = errors.ErrFailedToCopyDir(err, dst, tmpPath, nil, nil)
					log.Warn(err)
					return err
				}
				err = os.RemoveAll(dst)
				if err != nil && Exists(dst) {
					err = errors.ErrFailedToRemoveDir(err, dst, nil)
					if rollback {
						rerr := doMoveDir(ctx, tmpPath, dst, false)
						keepTmp = rerr != nil
						err = errors.Join(rerr, errors.Wrapf(err, "trying to recover temporary file %s to rollback previous operation", tmpPath))
					}
					log.Warn(err)
					return err
				}
			}
			log.Debugf("directory %s successfully moved to tmp location %s", dst, tmpPath)
		}
		exits, fi, err = ExistsWithDetail(src)
		if exits && fi != nil && fi.IsDir() && err == nil {
			err = rename(src, dst)
			if err != nil {
				log.Debugf("err: %v\t now trying to move file with I/O copy and Remove old index", errors.ErrFailedToRenameDir(err, src, dst, fi, nil))
				// Copy into a sibling staging dir on dst's filesystem and rename it
				// into place, so dst never appears partially copied. The name is
				// derived from dst (callers moving to the same dst are serialized), so
				// the pre-clean below also removes a stage left behind by an attempt
				// killed mid-copy. A failed pre-clean aborts the move, since leftover
				// files would otherwise be merged into dst.
				stage := Join(filepath.Dir(dst), "tmp-copy-"+filepath.Base(dst))
				err = os.RemoveAll(stage)
				// Create the stage first so an empty src still moves (CopyDir creates
				// nothing for it and the rename below would fail).
				if err == nil {
					err = os.MkdirAll(stage, fi.Mode().Perm())
				}
				if err == nil {
					err = CopyDir(ctx, src, stage)
				}
				if err != nil {
					err = errors.ErrFailedToCopyDir(err, src, stage, fi, nil)
				} else if err = rename(stage, dst); err != nil {
					err = errors.ErrFailedToRenameDir(err, stage, dst, nil, nil)
				}
				if err != nil {
					if rerr := os.RemoveAll(stage); rerr != nil {
						log.Warnf("failed to remove staging path %s: %v", stage, rerr)
					}
					if rollback {
						rerr := doMoveDir(ctx, tmpPath, dst, false)
						keepTmp = rerr != nil
						err = errors.Join(rerr, errors.Wrapf(err, "trying to recover temporary file %s to rollback previous operation", tmpPath))
					}
					log.Warn(err)
					return err
				}
				err = os.RemoveAll(src)
				if err != nil && Exists(src) {
					err = errors.ErrFailedToRemoveDir(err, src, fi)
					log.Warn(err)
					return err
				}
			}
		} else {
			return errors.ErrDirectoryNotFound(err, src, fi)
		}
	}
	log.Infof("directory %s successfully moved to destination directory %s", src, dst)
	return nil
}

func CopyDir(ctx context.Context, src, dst string) (err error) {
	if len(src) == 0 || len(dst) == 0 || src == dst {
		return nil
	}
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}
	eg, _ := errgroup.New(ctx)
	err = filepath.WalkDir(src, func(childPath string, info fs.DirEntry, err error) error {
		// Check err first: WalkDir reports a failure to stat or read src itself on the
		// src path, and info is nil when the Lstat of the path failed.
		if err != nil {
			var fi fs.FileInfo
			if info != nil {
				var ierr error
				fi, ierr = info.Info()
				if ierr != nil {
					err = errors.Join(err, ierr)
				}
			}
			err = errors.ErrFailedToWalkDir(err, src, childPath, nil, fi)
			log.Warn(err)
			return err
		}
		if childPath == src {
			return nil
		}
		// Skip dst (when it is inside src) and its subtree, but not siblings that only
		// share its name as a prefix (src/origin2 when dst is src/origin).
		if childPath == dst {
			if info.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(childPath, dst+string(filepath.Separator)) {
			return nil
		}
		// Keep the path relative to src (Base would flatten src/a/x and src/b/x onto
		// dst/x) and the source permissions (DirEntry.Type carries no permission bits).
		rel, err := filepath.Rel(src, childPath)
		if err != nil {
			return err
		}
		fi, err := info.Info()
		if err != nil {
			return err
		}
		dstPath := Join(dst, rel)
		if info.IsDir() {
			err = MkdirAll(dstPath, fi.Mode().Perm())
			if err != nil {
				log.Warn(errors.ErrFailedToMkdir(err, dstPath, nil))
			}
			return nil
		}
		eg.Go(safety.RecoverFunc(func() (err error) {
			_, err = CopyFileWithPerm(ctx, childPath, dstPath, fi.Mode().Perm())
			return err
		}))
		return nil
	})
	if err != nil {
		return errors.Join(eg.Wait(), err)
	}
	return eg.Wait()
}

func CopyFile(ctx context.Context, src, dst string) (n int64, err error) {
	return CopyFileWithPerm(ctx, src, dst, fs.ModePerm)
}

func CopyFileWithPerm(ctx context.Context, src, dst string, perm fs.FileMode) (n int64, err error) {
	if len(src) == 0 || len(dst) == 0 || src == dst {
		return 0, nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	defer func() {
		if err != nil {
			log.Warn(err)
		}
	}()

	exist, fi, err := ExistsWithDetail(src)
	switch {
	case !exist, fi == nil, fi.IsDir():
		return 0, errors.Join(err, errors.ErrFileNotFound(src))
	case err != nil && (!errors.Is(err, fs.ErrExist) || errors.Is(err, fs.ErrNotExist)):
		return 0, errors.Join(errors.ErrFileNotFound(src), err)
	case fi != nil && !fi.Mode().IsRegular():
		return 0, errors.ErrNonRegularFile(src, fi)
	case err != nil:
		log.Warn(err)
	}

	flg := os.O_RDONLY | os.O_SYNC
	sf, err := Open(src, flg, perm)
	if err != nil || sf == nil {
		err = errors.ErrFailedToCopyFile(errors.ErrFailedToOpenFile(err, src, flg, perm), src, dst, fi, nil)
		return 0, err
	}
	defer func() {
		if sf != nil {
			derr := sf.Close()
			if derr != nil {
				err = errors.Join(err, errors.ErrFailedToCloseFile(derr, src, fi))
			}
		}
	}()
	n, err = OverWriteFile(ctx, dst, sf, perm)
	if errors.IsNot(err, io.EOF) {
		err = errors.ErrFailedToCopyFile(err, src, dst, fi, nil)
		return 0, err
	}
	return n, nil
}

func WriteFile(
	ctx context.Context, target string, r io.Reader, perm fs.FileMode,
) (n int64, err error) {
	return writeFileWithContext(ctx, target, r, os.O_CREATE|os.O_WRONLY, perm)
}

func OverWriteFile(
	ctx context.Context, target string, r io.Reader, perm fs.FileMode,
) (n int64, err error) {
	return writeFileWithContext(ctx, target, r, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, perm)
}

func writeFileWithContext(
	ctx context.Context, target string, r io.Reader, flg int, perm fs.FileMode,
) (n int64, err error) {
	if len(target) == 0 || r == nil {
		return 0, nil
	}

	exist, fi, err := ExistsWithDetail(target)
	switch {
	case err == nil, exist, fi != nil && fi.Size() != 0, fi != nil && fi.IsDir():
		err = errors.ErrFileAlreadyExists(target)
	case err != nil && !errors.Is(err, fs.ErrNotExist), err != nil && errors.Is(err, fs.ErrExist):
		err = errors.Join(errors.ErrFileAlreadyExists(target), err)
	}

	// open flag is not O_TRUNC or O_APPEND this function returns AlreadyExists error
	if err != nil && flg&(os.O_TRUNC|os.O_APPEND) <= 0 &&
		errors.Is(err, errors.ErrFileAlreadyExists(target)) {
		return 0, err
	}

	// Create a missing parent with directory permissions here: Open would create it
	// with perm, and a file mode such as 0o640 makes a directory untraversable.
	if dir := filepath.Dir(target); !Exists(dir) {
		if err = MkdirAll(dir, fs.ModePerm); err != nil {
			err = errors.ErrFailedToOpenFile(errors.ErrFailedToMkdir(err, dir, nil), target, flg, perm)
			return 0, err
		}
	}

	f, err := Open(target, flg, perm)
	if err != nil || f == nil {
		err = errors.ErrFailedToOpenFile(err, target, flg, perm)
		return 0, err
	}
	defer func() {
		if f != nil {
			derr := f.Close()
			if derr != nil {
				err = errors.Join(err, errors.ErrFailedToCloseFile(derr, target, fi))
			}
		}
	}()
	w, err := io.NewWriterWithContext(ctx, f)
	if err != nil {
		return 0, err
	}
	cr, err := io.NewReaderWithContext(ctx, r)
	if err != nil {
		return 0, err
	}
	n, err = io.Copy(w, cr)
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, err
	}
	err = f.Sync()
	if err != nil {
		return 0, err
	}
	return n, nil
}

func ReadDir(path string) (dirs []fs.DirEntry, err error) {
	f, err := Open(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer func() {
		if f != nil {
			derr := f.Close()
			if derr != nil {
				err = errors.Join(err, errors.ErrFailedToCloseFile(derr, path, nil))
			}
		}
	}()

	dirs, err = f.ReadDir(-1)
	slices.SortFunc(dirs, func(a, b fs.DirEntry) int { return cmp.Compare(a.Name(), b.Name()) })
	return dirs, err
}

func ReadFile(path string) (n []byte, err error) {
	f, err := Open(path, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		return nil, err
	}
	defer func() {
		if f != nil {
			derr := f.Close()
			if derr != nil {
				err = errors.Join(err, errors.ErrFailedToCloseFile(derr, path, nil))
			}
		}
	}()
	return io.ReadAll(f)
}

func AbsolutePath(path string) string {
	if path == "" {
		return ""
	}
	if !filepath.IsAbs(path) {
		root, err := os.Getwd()
		if err == nil {
			path = joinFilePaths(root, path)
		}
		if !filepath.IsAbs(path) {
			absPath, err := filepath.Abs(path)
			if err == nil {
				path = absPath
			}
		}
	}
	return filepath.Clean(path)
}

// Exists returns file existence.
func Exists(path string) (e bool) {
	e, _, _ = ExistsWithDetail(path)
	return e
}

// ExistsWithDetail returns file existence with detailed information.
func ExistsWithDetail(path string) (e bool, fi fs.FileInfo, err error) {
	if path == "" {
		return false, nil, fs.ErrInvalid
	}
	fi, err = os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true, fi, nil
		}
		if os.IsNotExist(err) {
			return false, fi, err
		}
		return false, fi, err
	}
	return true, fi, nil
}

// MkdirAll creates directory like mkdir -p.
func MkdirAll(path string, perm fs.FileMode) (err error) {
	var (
		exist bool
		fi    fs.FileInfo
		merr  error
	)
	exist, fi, err = ExistsWithDetail(path)
	if exist {
		if err == nil && fi != nil && fi.IsDir() {
			return nil
		}
		// never delete whatever already exists at the path: it may be user data (e.g. an index file).
		return errors.ErrFailedToMkdir(fs.ErrExist, path, fi)
	}
	if err != nil {
		err = errors.ErrDirectoryNotFound(err, path, fi)
	}
	merr = os.MkdirAll(path, perm)
	if merr == nil {
		return nil
	}
	err = errors.Join(merr, err)
	if err != nil {
		if errors.Is(err, fs.ErrPermission) {
			merr = os.MkdirAll(path, fs.ModePerm)
			if merr != nil {
				err = errors.Join(err, errors.ErrFailedToMkdir(merr, path, fi))
			}
		}
		log.Warn(err)
		return err
	}
	return nil
}

// MkdirTemp create temporary directory from given base path
// if base path is nil temporary directory will create from Go's standard library.
func MkdirTemp(baseDir string) (path string, err error) {
	if len(baseDir) == 0 {
		baseDir = os.TempDir()
	}
	path = Join(baseDir, strconv.FormatInt(time.Now().UnixNano(), 10))
	err = MkdirAll(path, fs.ModePerm)
	if err != nil {
		err = errors.ErrFailedToMkTmpDir(err, path, nil)
		log.Debug(err)
		return "", err
	}
	return path, nil
}

// CreateTemp create temporary file from given base path
// if base path is nil temporary directory will create from Go's standard library.
func CreateTemp(baseDir string) (f *os.File, err error) {
	dir, err := MkdirTemp(baseDir)
	if err != nil {
		return nil, err
	}
	var path string
	for range 10000 {
		path = Join(dir, strconv.FormatInt(time.Now().UnixNano(), 10))
		f, err = Open(path, os.O_CREATE|os.O_EXCL|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		if err == nil && f != nil {
			return f, nil
		}
		if f != nil {
			f.Close()
		}
	}
	return nil, errors.ErrFailedToCreateFile(err, path, nil)
}

// ListInDir returns file list in directory.
func ListInDir(path string) ([]string, error) {
	exists, fi, err := ExistsWithDetail(path)
	if !exists {
		return nil, err
	}
	if fi.Mode().IsDir() && !strings.HasSuffix(path, sep) {
		path += sep
	}
	path = filepath.Dir(path)
	files, err := filepath.Glob(Join(path, "*"))
	if err != nil {
		return nil, err
	}
	return files, nil
}

// DeleteDir recursively deletes each file or directory based on the path specified as an argument.
func DeleteDir(_ context.Context, path string) (err error) {
	exists, _, err := ExistsWithDetail(path)
	if err != nil {
		return err
	}
	if !exists {
		log.Debugf("path %s does not exist, so deleting the directory is not executed.", path)
		return nil
	}
	// RemoveAll deletes children before their parents; removing every entry
	// concurrently made a parent's Remove race its children and fail with ENOTEMPTY.
	if err = os.RemoveAll(path); err != nil {
		return errors.ErrFailedToRemoveDir(err, path, nil)
	}
	return nil
}

func Join(paths ...string) (path string) {
	if paths == nil {
		return ""
	}
	if len(paths) > 1 {
		path = joinFilePaths(paths...)
	} else {
		path = replacer.Replace(paths[0])
	}
	return AbsolutePath(path)
}

var replacer = strings.NewReplacer(
	sep+sep+sep,
	sep,
	sep+sep,
	sep,
)

func joinFilePaths(paths ...string) (path string) {
	for i, path := range paths {
		if path != "" {
			return replacer.Replace(strings.Join(paths[i:], sep))
		}
	}
	return ""
}
