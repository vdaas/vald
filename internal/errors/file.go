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

// Package errors provides error types and function
package errors

var (
	// ErrWatchDirNotFound represents an error that the watch directory is not found.
	ErrWatchDirNotFound = New("fs watcher watch dir not found")

	// ErrFileAlreadyExists represents a function to generate an error that the file already exists.
	ErrFileAlreadyExists = func(path string) error {
		return Errorf("file already exists: %s", path)
	}

	// ErrPathNotSpecified represents an error that the path is not specified.
	ErrPathNotSpecified = New("the path is not specified")
)
