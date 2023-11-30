//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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

// Package errors provides error types and function
package errors

import (
	"fmt"
	"os"
)

var (
	// ErrWatchDirNotFound represents an error that the watch directory is not found.
	ErrWatchDirNotFound = New("fs watcher watch dir not found")

	// ErrFileAlreadyExists represents a function to generate an error that the file already exists.
	ErrFileAlreadyExists = func(path string) error {
		return Errorf("file already exists: %s", path)
	}

	// ErrFileNotFound represents a function to generate an error that the file not found.
	ErrFileNotFound = func(path string) error {
		return Errorf("file not found: %s", path)
	}

	// ErrPathNotSpecified represents an error that the path is not specified.
	ErrPathNotSpecified = New("the path is not specified")

	// ErrPathNotAllowed represents a function to generate an error indicates that the specified path is not allowed.
	ErrPathNotAllowed = func(path string) error {
		return Errorf("the specified file path is not allowed: %s", path)
	}

	ErrDirectoryNotFound = func(err error, dir string, fi os.FileInfo) error {
		return Wrapf(err, "directory not found %s\tdir info: %s", dir, fitos(dir, fi))
	}

	ErrFailedToGetAbsPath = func(err error, path string) error {
		return Wrapf(err, "failed to get Absolute path %s using filepath.Abs operation \tfile info: %s", path, fitos(path, nil))
	}

	ErrFailedToMkdir = func(err error, dir string, fi os.FileInfo) error {
		return Wrapf(err, "failed to make directory %s using os.MkdirAll operation\tdir info: %s", dir, fitos(dir, fi))
	}

	ErrFailedToMkTmpDir = func(err error, tmpDir string, fi os.FileInfo) error {
		return Wrapf(err, "failed to make temporary directory %s\tdir info: %s", tmpDir, fitos(tmpDir, fi))
	}

	ErrFailedToCreateFile = func(err error, path string, fi os.FileInfo) error {
		return Wrapf(err, "failed to create file %s using os.Create operation\tfile info: %s", path, fitos(path, fi))
	}

	ErrFailedToRemoveFile = func(err error, path string, fi os.FileInfo) error {
		return Wrapf(err, "failed to remove file %s using os.RemoveAll operation\tfile info: %s", path, fitos(path, fi))
	}

	ErrFailedToRemoveDir = func(err error, path string, fi os.FileInfo) error {
		return Wrapf(err, "failed to remove directory %s using os.RemoveAll operation\tfile info: %s", path, fitos(path, fi))
	}

	ErrFailedToOpenFile = func(err error, path string, flg int, perm os.FileMode) error {
		return Wrapf(err, "failed to open file %s using os.OpenFile operation with configuration, flg: %d, perm: %s\tfile info: %s", path, flg, perm.String(), fitos(path, nil))
	}

	ErrFailedToCloseFile = func(err error, path string, fi os.FileInfo) error {
		return Wrapf(err, "failed to close file %s\tfile info: %s", path, fitos(path, fi))
	}

	ErrFailedToRenameDir = func(err error, src, dst string, sfi, dfi os.FileInfo) error {
		return Wrapf(err, "failed to rename directory %s to %s using os.Rename operation,\tfile info: {source: %s, destination: %s}", src, dst, fitos(src, sfi), fitos(dst, dfi))
	}

	ErrFailedToCopyFile = func(err error, src, dst string, sfi, dfi os.FileInfo) error {
		return Wrapf(err, "failed to copy file %s to %s using os.Rename operation,\tfile info: {source: %s, destination: %s}", src, dst, fitos(src, sfi), fitos(dst, dfi))
	}

	ErrFailedToCopyDir = func(err error, src, dst string, sfi, dfi os.FileInfo) error {
		return Wrapf(err, "failed to copy directory %s to %s using os.Rename operation,\tfile info: {source: %s, destination: %s}", src, dst, fitos(src, sfi), fitos(dst, dfi))
	}

	ErrFailedToWalkDir = func(err error, root, dir string, rfi, cfi os.FileInfo) error {
		return Wrapf(err, "failed to walk directory %s in %s using filepath.WalkDir operation,\tfile info: {root: %s, current: %s}", dir, root, fitos(root, rfi), fitos(dir, cfi))
	}

	ErrNonRegularFile = func(path string, fi os.FileInfo) error {
		return Errorf("error file is not a regular file %s", fitos(path, fi))
	}
)

func fitos(path string, fi os.FileInfo) string {
	if fi == nil {
		var err error
		fi, err = os.Stat(path)
		if err != nil || fi == nil {
			return fmt.Sprintf("unknown file info: %v", fi)
		}
	}
	if fi != nil {
		return fmt.Sprintf("{name: %s, size: %d, mode: %s, mode_int: %d, is_dir: %v}",
			fi.Name(),
			fi.Size(),
			fi.Mode().String(),
			fi.Mode(),
			fi.IsDir(),
		)
	}
	return fmt.Sprintf("unknown file info: %v", fi)
}
