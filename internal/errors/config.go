//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

import "reflect"

var (
	ErrInvalidConfig = New("component config is invalid")

	ErrUnsupportedConfigFileType = func(ext string) error {
		return Errorf("unsupported file type: %s", ext)
	}

	ErrNotMatchFieldType = func(path string, dType, sType reflect.Type) error {
		return Errorf("types do not match at %s: %v vs %v", path, dType, sType)
	}

	ErrNotMatchFieldNum = func(path string, dNum, sNum int) error {
		return Errorf("number of fields do not match at %s, dst: %d, src: %d", path, dNum, sNum)
	}

	ErrNotMatchArrayLength = func(path string, dLen, sLen int) error {
		return Errorf("array length do not match at %s, dst: %d, src: %d", path, dLen, sLen)
	}

	ErrDeepMergeKind = func(kind string, nf string, err error) error {
		return Errorf("error in %s at %s: %w", kind, nf, err)
	}
)
