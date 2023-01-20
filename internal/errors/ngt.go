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
