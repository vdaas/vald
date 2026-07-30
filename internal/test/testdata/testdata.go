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

// Package testdata resolves fixture files under internal/test/data. It is
// deliberately separate from the framework layers (tb, table): path
// resolution needs the filesystem wrappers (internal/file, internal/os),
// and keeping those dependencies out of the framework keeps the latter a
// leaf package.
package testdata

import (
	"path/filepath"

	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
)

// ValidIndex is the fixture path of a valid NGT agent index.
const ValidIndex = "agent/ngt/validIndex"

// GetTestdataPath returns the test data file path under `internal/test/data`.
func GetTestdataPath(filename string) string {
	return file.Join(baseDir(), "/internal/test/data/", filename)
}

// baseDir walks up from the working directory to the repository root,
// identified as the directory whose go.mod declares the vald module path.
// A directory-name suffix match would also accept forks or unrelated
// directories named *vald, and a bare go.mod-presence check would stop at
// nested modules such as example/client.
func baseDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	for cur := wd; cur != string(os.PathSeparator); cur = filepath.Dir(cur) {
		b, err := os.ReadFile(file.Join(cur, "go.mod"))
		if err == nil && strings.Contains(string(b), "module github.com/vdaas/vald\n") {
			return cur
		}
	}
	return ""
}
