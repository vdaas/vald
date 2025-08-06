//
// Copyright (C) 2019-2025 vdaas.org vald team <vald@vdaas.org>
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

package errors

var (
	// NewBlobNoSuchBucketError represents a function to create no such bucket error.
	NewBlobNoSuchBucketError = func(err error, name string) error {
		return &BlobNoSuchBucketError{
			err: Wrap(err, Errorf("bucket %s not found", name).Error()),
		}
	}

	// NewBlobNoSuchKeyError represents a function to create no such key error.
	NewBlobNoSuchKeyError = func(err error, key string) error {
		return &BlobNoSuchKeyError{
			err: Wrap(err, Errorf("key %s not found", key).Error()),
		}
	}

	// NewBlobInvalidChunkRangeError represents a function to create invalid chunk range error.
	NewBlobInvalidChunkRangeError = func(err error, rng string) error {
		return &BlobInvalidChunkRangeError{
			err: Wrap(err, Errorf("chunk range %s is invalid", rng).Error()),
		}
	}
)

// BlobNoSuchBucketError represents no such bucket error of S3.
type BlobNoSuchBucketError struct {
	err error
}

// Error returns the string representation of the internal error.
func (e *BlobNoSuchBucketError) Error() string {
	return e.err.Error()
}

// Unwrap unwraps and returns the internal error.
func (e *BlobNoSuchBucketError) Unwrap() error {
	return e.err
}

// IsErrBlobNoSuchBucket returns true if the error is ErrBlobNoSuchBucket.
func IsErrBlobNoSuchBucket(err error) bool {
	var target *BlobNoSuchBucketError
	return As(err, &target)
}

// BlobNoSuchKeyError represents no such key error of S3.
type BlobNoSuchKeyError struct {
	err error
}

// Error returns the string representation of the internal error.
func (e *BlobNoSuchKeyError) Error() string {
	return e.err.Error()
}

// Unwrap unwraps and returns the internal error.
func (e *BlobNoSuchKeyError) Unwrap() error {
	return e.err
}

// IsErrBlobNoSuchKey returns true if the error is ErrBlobNoSuchKey.
func IsErrBlobNoSuchKey(err error) bool {
	var target *BlobNoSuchKeyError
	return As(err, &target)
}

// BlobInvalidChunkRangeError represents no invalid chunk range error of S3.
type BlobInvalidChunkRangeError struct {
	err error
}

// Error returns the string representation of the internal error.
func (e *BlobInvalidChunkRangeError) Error() string {
	return e.err.Error()
}

// Unwrap unwraps and returns the internal error.
func (e *BlobInvalidChunkRangeError) Unwrap() error {
	return e.err
}
