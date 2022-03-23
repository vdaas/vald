//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
package test

import (
	"os"
	"path/filepath"

	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/strings"
)

// GetTestdataPath returns the test data file path under `internal/test/data`.
func GetTestdataPath(filename string) string {
	return file.Join(baseDir(), "/internal/test/data/", filename)
}

func baseDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}

	for cur := filepath.Dir(wd); cur != string(os.PathSeparator); cur = filepath.Dir(cur) {
		if strings.HasSuffix(cur, "vald") {
			return cur
		}
	}
	return ""
}
