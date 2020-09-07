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
			err: New("cassandra entry not found"),
		}
	}

	NewErrCassandraUnavailableIdentity = func() error {
		return &ErrCassandraUnavailableIdentity{
			err: New("cassandra unavailable"),
		}
	}

	ErrCassandraUnavailable = NewErrCassandraUnavailableIdentity

	ErrCassandraNotFound = func(keys ...string) error {
		switch {
		case len(keys) == 1:
			return Wrapf(NewErrCassandraNotFoundIdentity(), "cassandra key '%s' not found", keys[0])
		case len(keys) > 1:
			return Wrapf(NewErrCassandraNotFoundIdentity(), "cassandra keys '%v' not found", keys)
		default:
			return nil
		}
	}

	ErrCassandraGetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "error failed to fetch key (%s)", key)
	}

	ErrCassandraSetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "error failed to set key (%s)", key)
	}

	ErrCassandraDeleteOperationFailed = func(key string, err error) error {
		return Wrapf(err, "error failed to delete key (%s)", key)
	}

	ErrCassandraHostDownDetected = func(err error, nodeInfo string) error {
		return Wrapf(err, "error cassandra host down detected\t%s", nodeInfo)
	}
)

type ErrCassandraNotFoundIdentity struct {
	err error
}

func (e *ErrCassandraNotFoundIdentity) Error() string {
	return e.err.Error()
}

func (e *ErrCassandraNotFoundIdentity) Unwrap() error {
	return e.err
}

func IsErrCassandraNotFound(err error) bool {
	return As(err, &ErrCassandraNotFoundIdentity{})
}

type ErrCassandraUnavailableIdentity struct {
	err error
}

func (e *ErrCassandraUnavailableIdentity) Error() string {
	return e.err.Error()
}

func (e *ErrCassandraUnavailableIdentity) Unwrap() error {
	return e.err
}

func IsErrCassandraUnavailable(err error) bool {
	return As(err, &ErrCassandraUnavailableIdentity{})
}
