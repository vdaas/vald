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

// Package errors provides error types and function
package errors

var (
	// ErrCreateIndexingIsInProgress represents an error that the indexing is in progress but search request received
	ErrCreateIndexingIsInProgress = New("create indexing is in progress")

	// ErrFlushingIsInProgress represents an error that the flushing is in progress but any request received
	ErrFlushingIsInProgress = New("flush is in progress")

	// ErrCreateProperty represents a function to generate an error that the property creation failed.
	ErrCreateProperty = func(err error) error {
		return Wrap(err, "failed to create property")
	}

	// ErrIndexFileNotFound represents an error that the index file is not found.
	ErrIndexFileNotFound = New("index file not found")

	// ErrIndicesAreTooFewComparedToMetadata represents an error that the index count is not enough to be compared by metadata.
	ErrIndicesAreTooFewComparedToMetadata = New("indices are too few compared to Metadata")

	// ErrIndexLoadTimeout represents an error that the index loading timeout.
	ErrIndexLoadTimeout = New("index load timeout")

	// ErrInvalidDimensionSize represents a function to generate an error that the dimension size is invalid.
	ErrInvalidDimensionSize = func(current, limit int) error {
		if limit == 0 {
			return Errorf("dimension size %d is invalid, the supporting dimension size must be bigger than 2", current)
		}
		return Errorf("dimension size %d is invalid, the supporting dimension size must be between 2 ~ %d", current, limit)
	}

	// ErrInvalidUUID represents a function to generate an error that the uuid is invalid.
	ErrInvalidUUID = func(uuid string) error {
		return Errorf("uuid \"%s\" is invalid", uuid)
	}

	// ErrDimensionLimitExceed represents a function to generate an error that the supported dimension limit exceeded.
	ErrDimensionLimitExceed = func(current, limit int) error {
		return Errorf("supported dimension limit exceed:\trequired = %d,\tlimit = %d", current, limit)
	}

	// ErrIncompatibleDimensionSize represents a function to generate an error that the incompatible dimension size detected.
	ErrIncompatibleDimensionSize = func(current, expected int) error {
		return Errorf("incompatible dimension size detected\trequested: %d,\tconfigured: %d", current, expected)
	}

	// ErrUnsupportedObjectType represents an error that the object type is unsupported.
	ErrUnsupportedObjectType = New("unsupported ObjectType")

	// ErrUnsupportedDistanceType represents an error that the distance type is unsupported.
	ErrUnsupportedDistanceType = New("unsupported DistanceType")

	// ErrFailedToSetDistanceType represents a function to generate an error that the set of distance type failed.
	ErrFailedToSetDistanceType = func(err error, distance string) error {
		return Wrap(err, "failed to set distance type "+distance)
	}

	// ErrFailedToSetObjectType represents a function to generate an error that the set of object type failed.
	ErrFailedToSetObjectType = func(err error, t string) error {
		return Wrap(err, "failed to set object type "+t)
	}

	// ErrFailedToSetDimension represents a function to generate an error that the set of dimension failed.
	ErrFailedToSetDimension = func(err error) error {
		return Wrap(err, "failed to set dimension")
	}

	// ErrFailedToSetCreationEdgeSize represents a function to generate an error that the set of creation edge size failed.
	ErrFailedToSetCreationEdgeSize = func(err error) error {
		return Wrap(err, "failed to set creation edge size")
	}

	// ErrFailedToSetSearchEdgeSize represents a function to generate an error that the set of search edge size failed.
	ErrFailedToSetSearchEdgeSize = func(err error) error {
		return Wrap(err, "failed to set search edge size")
	}

	// ErrUncommittedIndexExists represents a function to generate an error that the uncommitted indexes exist.
	ErrUncommittedIndexExists = func(num uint64) error {
		return Errorf("%d indexes are not committed", num)
	}

	// ErrUncommittedIndexNotFound represents an error that the uncommitted indexes are not found.
	ErrUncommittedIndexNotFound = New("uncommitted indexes are not found")

	// ErrCAPINotImplemented represents an error that the function is not implemented in C API.
	ErrCAPINotImplemented = New("not implemented in C API")

	// ErrUUIDAlreadyExists represents a function to generate an error that the uuid already exists.
	ErrUUIDAlreadyExists = func(uuid string) error {
		return Errorf("ngt uuid %s index already exists", uuid)
	}

	// ErrUUIDNotFound represents a function to generate an error that the uuid is not found.
	ErrUUIDNotFound = func(id uint32) error {
		if id == 0 {
			return New("ngt object uuid not found")
		}
		return Errorf("ngt object uuid %d's metadata not found", id)
	}

	// ErrObjectIDNotFound represents a function to generate an error that the object id is not found.
	ErrObjectIDNotFound = func(uuid string) error {
		return Errorf("ngt uuid %s's object id not found", uuid)
	}

	// ErrObjectNotFound represents a function to generate an error that the object is not found.
	ErrObjectNotFound = func(err error, uuid string) error {
		return Wrapf(err, "ngt uuid %s's object not found", uuid)
	}

	// ErrRemoveRequestedBeforeIndexing represents a function to generate an error that the object is not indexed so can not remove it.
	ErrRemoveRequestedBeforeIndexing = func(oid uint) error {
		return Errorf("object id %d is not indexed we cannot remove it", oid)
	}
)

type NGTError struct {
	Msg string
}

func NewNGTError(msg string) error {
	return NGTError{
		Msg: msg,
	}
}

func (n NGTError) Error() string {
	return n.Msg
}
