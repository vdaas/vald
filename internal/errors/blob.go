//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	// BlobStorage
	NewErrBlobNoSuchBucket = func(err error, name string) error {
		return &ErrBlobNoSuchBucket{
			err: Wrap(err, Errorf("bucket %s not found", name).Error()),
		}
	}

	NewErrBlobNoSuchKey = func(err error, key string) error {
		return &ErrBlobNoSuchKey{
			err: Wrap(err, Errorf("key %s not found", key).Error()),
		}
	}
)

type ErrBlobNoSuchBucket struct {
	err error
}

func (e *ErrBlobNoSuchBucket) Error() string {
	return e.err.Error()
}

func (e *ErrBlobNoSuchBucket) Unwrap() error {
	return e.err
}

func IsErrBlobNoSuchBucket(err error) bool {
	var target error = new(ErrBlobNoSuchBucket)
	return As(err, &target)
}

type ErrBlobNoSuchKey struct {
	err error
}

func (e *ErrBlobNoSuchKey) Error() string {
	return e.err.Error()
}

func (e *ErrBlobNoSuchKey) Unwrap() error {
	return e.err
}

func IsErrBlobNoSuchKey(err error) bool {
	var target error = new(ErrBlobNoSuchKey)
	return As(err, &target)
}
