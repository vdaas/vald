//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package file provides file I/O functionality
package file

import (
	"os"
	"path/filepath"

	"github.com/vdaas/vald/internal/errors"
)

// Open opens the file with the given path, flag and permission.
// If the folder does not exists, create the folder.
// If the file does not exist, create the file.
func Open(path string, flg int, perm os.FileMode) (*os.File, error) {
	if path == "" {
		return nil, errors.ErrPathNotSpecified
	}

	var err error
	var file *os.File
	if _, err = os.Stat(path); err != nil {
		if _, err = os.Stat(filepath.Dir(path)); err != nil {
			err = os.MkdirAll(filepath.Dir(path), perm)
			if err != nil {
				return nil, err
			}
		}

		if flg&(os.O_CREATE|os.O_APPEND) > 0 {
			file, err = os.Create(path)
			if err != nil {
				return nil, err
			}
		}

		if file != nil {
			err = file.Close()
			if err != nil {
				return nil, err
			}
		}
	}

	file, err = os.OpenFile(path, flg, perm)
	if err != nil {
		return nil, err
	}

	return file, nil
}
