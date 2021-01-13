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
	// NewErrBlobNoSuchBucket represents a function to create no such bucket error.
	NewErrBlobNoSuchBucket = func(err error, name string) error {
		return &ErrBlobNoSuchBucket{
			err: Wrap(err, Errorf("bucket %s not found", name).Error()),
		}
	}

	// NewErrBlobNoSuchKey represents a function to create no such key error.
	NewErrBlobNoSuchKey = func(err error, key string) error {
		return &ErrBlobNoSuchKey{
			err: Wrap(err, Errorf("key %s not found", key).Error()),
		}
	}
)

// ErrBlobNoSuchBucket represents no such bucket error of S3.
type ErrBlobNoSuchBucket struct {
	err error
}

// Error returns the string representation of the internal error.
func (e *ErrBlobNoSuchBucket) Error() string {
	return e.err.Error()
}

// Unwrap unwraps and returns the internal error.
func (e *ErrBlobNoSuchBucket) Unwrap() error {
	return e.err
}

// IsErrBlobNoSuchBucket returns if the error is ErrBlobNoSuchBucket.
func IsErrBlobNoSuchBucket(err error) bool {
	target := new(ErrBlobNoSuchBucket)
	return As(err, &target)
}

// ErrBlobNoSuchKey represents no such key error of S3.
type ErrBlobNoSuchKey struct {
	err error
}

// Error returns the string representation of the internal error.
func (e *ErrBlobNoSuchKey) Error() string {
	return e.err.Error()
}

// Unwrap unwraps and returns the internal error.
func (e *ErrBlobNoSuchKey) Unwrap() error {
	return e.err
}

// IsErrBlobNoSuchKey returns if the error is ErrBlobNoSuchKey.
func IsErrBlobNoSuchKey(err error) bool {
	target := new(ErrBlobNoSuchKey)
	return As(err, &target)
}
