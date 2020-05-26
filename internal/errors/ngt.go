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
	//NGT

	ErrCreateProperty = func(err error) error {
		return Wrap(err, "failed to create property")
	}

	ErrIndexNotFound = New("index file not found")

	ErrInvalidDimensionSize = func(current, limit int) error {
		if limit == 0 {
			return Errorf("dimension size %d is invalid, the supporting dimension size must be bigger than 2", current)
		}
		return Errorf("dimension size %d is invalid, the supporting dimension size must be between 2 ~ %d", current, limit)
	}

	ErrIncompatibleDimensionSize = func(req, dim int) error {
		return Errorf("incompatible dimension size detected\trequested: %d,\tconfigured: %d", req, dim)
	}

	ErrUnsupportedObjectType = New("unsupported ObjectType")

	ErrUnsupportedDistanceType = New("unsupported DistanceType")

	ErrFailedToSetDistanceType = func(err error, distance string) error {
		return Wrap(err, "failed to set distance type "+distance)
	}

	ErrFailedToSetObjectType = func(err error, t string) error {
		return Wrap(err, "failed to set object type "+t)
	}

	ErrFailedToSetDimension = func(err error) error {
		return Wrap(err, "failed to set dimension")
	}

	ErrFailedToSetCreationEdgeSize = func(err error) error {
		return Wrap(err, "failed to set creation edge size")
	}

	ErrFailedToSetSearchEdgeSize = func(err error) error {
		return Wrap(err, "failed to set search edge size")
	}

	ErrUncommittedIndexExists = func(num uint64) error {
		return Errorf("%d indexes are not committed", num)
	}

	ErrUncommittedIndexNotFound = New("uncommitted indexes are not found")

	// ErrCAPINotImplemented raises using not implemented function in C API
	ErrCAPINotImplemented = New("not implemented in C API")

	ErrUUIDAlreadyExists = func(uuid string, oid uint) error {
		return Errorf("ngt uuid %s object id %d already exists ", uuid, oid)
	}

	ErrUUIDNotFound = func(id uint32) error {
		if id == 0 {
			return Errorf("ngt object uuid not found", id)
		}
		return Errorf("ngt object uuid %d's metadata not found", id)
	}

	ErrObjectIDNotFound = func(uuid string) error {
		return Errorf("ngt uuid %s's object id not found", uuid)
	}

	ErrObjectNotFound = func(err error, uuid string) error {
		return Wrapf(err, "ngt uuid %s's object not found", uuid)
	}

	ErrRemoveRequestedBeforeIndexing = func(oid uint) error {
		return Errorf("object id %d is not indexed we cannot remove it", oid)
	}
)
