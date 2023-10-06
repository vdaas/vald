//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package errors provides error types and function
package errors

var (
	// ErrCassandraInvalidConsistencyType represents a function to generate an error of consistency type not defined.
	ErrCassandraInvalidConsistencyType = func(consistency string) error {
		return Errorf("consistetncy type %q is not defined", consistency)
	}

	// ErrCassandraNotFoundIdentity generates an error of cassandra entry not found.
	ErrCassandraNotFoundIdentity = &CassandraNotFoundIdentityError{
		err: New("cassandra entry not found"),
	}

	// ErrCassandraUnavailableIdentity generates an error of cassandra unavailable.
	ErrCassandraUnavailableIdentity = &CassandraUnavailableIdentityError{
		err: New("cassandra unavailable"),
	}

	// ErrCassandraUnavailable represents NewErrCassandraUnavailableIdentity.
	ErrCassandraUnavailable = ErrCassandraUnavailableIdentity

	// ErrCassandraNotFound represents a function to generate an error of cassandra keys not found.
	ErrCassandraNotFound = func(keys ...string) error {
		switch {
		case len(keys) == 1:
			return Wrapf(ErrCassandraNotFoundIdentity, "cassandra key '%s' not found", keys[0])
		case len(keys) > 1:
			return Wrapf(ErrCassandraNotFoundIdentity, "cassandra keys '%v' not found", keys)
		default:
			return nil
		}
	}

	// ErrCassandraGetOperationFailed represents a function to generate an error of fetch key failed.
	ErrCassandraGetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "error failed to fetch key (%s)", key)
	}

	// ErrCassandraSetOperationFailed represents a function to generate an error of set key failed.
	ErrCassandraSetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "error failed to set key (%s)", key)
	}

	// ErrCassandraDeleteOperationFailed represents a function to generate an error of delete key failed.
	ErrCassandraDeleteOperationFailed = func(key string, err error) error {
		return Wrapf(err, "error failed to delete key (%s)", key)
	}

	// ErrCassandraHostDownDetected represents a function to generate an error of cassandra host down detected.
	ErrCassandraHostDownDetected = func(err error, nodeInfo string) error {
		return Wrapf(err, "error cassandra host down detected\t%s", nodeInfo)
	}
	ErrCassandraFailedToCreateSession = func(err error, hosts []string, port int, cqlVersion string) error {
		return Wrapf(err, "error cassandra client failed to create session to hosts: %v\tport: %d\tcql_version: %s ", hosts, port, cqlVersion)
	}
)

// CassandraNotFoundIdentityError represents custom error for cassandra not found.
type CassandraNotFoundIdentityError struct {
	err error
}

// Error returns string of internal error.
func (e *CassandraNotFoundIdentityError) Error() string {
	return e.err.Error()
}

// Unwrap returns an internal error.
func (e *CassandraNotFoundIdentityError) Unwrap() error {
	return e.err
}

// IsCassandraNotFoundError reports whether any error in err's chain matches CassandraNotFoundError.
func IsCassandraNotFoundError(err error) bool {
	target := new(CassandraNotFoundIdentityError)
	return As(err, &target)
}

// CassandraUnavailableIdentityError represents custom error for cassandra unavailable.
type CassandraUnavailableIdentityError struct {
	err error
}

// Error returns string of internal error.
func (e *CassandraUnavailableIdentityError) Error() string {
	return e.err.Error()
}

// Unwrap returns internal error.
func (e *CassandraUnavailableIdentityError) Unwrap() error {
	return e.err
}

// IsCassandraUnavailableError reports whether any error in err's chain matches CassandraUnavailableIdentityError.
func IsCassandraUnavailableError(err error) bool {
	target := new(CassandraUnavailableIdentityError)
	return As(err, &target)
}
