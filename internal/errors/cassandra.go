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

	// Cassandra
	ErrCassandraInvalidConsistencyType = func(consistency string) error {
		return Errorf("consistetncy type %q is not defined", consistency)
	}

	NewErrCassandraNotFoundIdentity = func() error {
		return &ErrCassandraNotFoundIdentity{
			err: New("error cassandra entry not found"),
		}
	}

	NewErrCassandraUnavailableIdentity = func() error {
		return &ErrCassandraUnavailableIdentity{
			err: New("error cassandra unavailable"),
		}
	}

	ErrCassandraUnavailable = func() error {
		return NewErrCassandraUnavailableIdentity()
	}

	ErrCassandraNotFound = func(keys ...string) error {
		switch {
		case len(keys) == 1:
			return Wrapf(NewErrCassandraNotFoundIdentity(), "error cassandra key '%s' not found", keys[0])
		case len(keys) > 1:
			return Wrapf(NewErrCassandraNotFoundIdentity(), "error cassandra keys '%s' not found", keys)
		default:
			return nil
		}
	}

	ErrCassandraGetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to fetch key (%s)", key)
	}

	ErrCassandraSetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to set key (%s)", key)
	}

	ErrCassandraDeleteOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to delete key (%s)", key)
	}
)

type ErrCassandraNotFoundIdentity struct {
	err error
}

func (e *ErrCassandraNotFoundIdentity) Error() string {
	return e.err.Error()
}

func IsErrCassandraNotFound(err error) bool {
	switch err.(type) {
	case *ErrCassandraNotFoundIdentity:
		return true
	default:
		return false
	}
}

type ErrCassandraUnavailableIdentity struct {
	err error
}

func (e *ErrCassandraUnavailableIdentity) Error() string {
	return e.err.Error()
}

func IsErrCassandraUnavailable(err error) bool {
	switch err.(type) {
	case *ErrCassandraUnavailableIdentity:
		return true
	default:
		return false
	}
}
