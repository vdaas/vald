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

var (
	// ErrMetaDataAlreadyExists represents an error that vald metadata is already exists.
	ErrMetaDataAlreadyExists = func(meta string) error {
		return Errorf("vald metadata:\t%s\talready exists ", meta)
	}

	// ErrSameVectorAlreadyExists represents an error that vald already has same features vector data.
	ErrSameVectorAlreadyExists = func(meta string, n, o []float32) error {
		return Errorf("vald metadata:\t%s\talready exists reqested: %v, stored: %v", meta, n, o)
	}

	// ErrMetaDataCannotFetch represents an error that vald metadata cannot fetch.
	ErrMetaDataCannotFetch = Errorf("vald metadata cannot fetch")
)
