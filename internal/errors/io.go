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
	// NewErrContextNotProvided represents a function to generate an error that the context is not provided.
	NewErrContextNotProvided = func() error {
		return New("context not provided")
	}

	// NewErrReaderNotProvided represents a function to generate an error that the io.Reader is not provided.
	NewErrReaderNotProvided = func() error {
		return New("io.Reader not provided")
	}

	// NewErrWriterNotProvided represents a function to generate an error that the io.Writer is not provided.
	NewErrWriterNotProvided = func() error {
		return New("io.Writer not provided")
	}
)
