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

	// ErrRedisInvalidKVVKPrefix represents a function to generate an error that kv index and vk prefix are invalid.
	ErrRedisInvalidKVVKPrefix = func(kv, vk string) error {
		return Errorf("kv index and vk prefix must be defferent.\t(kv: %s,\tvk: %s)", kv, vk)
	}
	
	// NewErrRedisNotFoundIdentity represents a function to generate an ErrRedisNotFoundIdentity error.
	NewErrRedisNotFoundIdentity = func() error {
		return &ErrRedisNotFoundIdentity{
			err: New("error redis entry not found"),
		}
	}

	// ErrRedisNotFound represents a function to wrap redis key not found error and err.
	ErrRedisNotFound = func(key string) error {
		return Wrapf(NewErrRedisNotFoundIdentity(), "error redis key '%s' not found", key)
	}

	// ErrRedisInvalidOption generates a new error of redis invalid option.
	ErrRedisInvalidOption = New("error redis invalid option")

	// ErrRedisGetOperationFailed represents a function to wrap failed to fetch key error and err.
	ErrRedisGetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to fetch key (%s)", key)
	}

	// ErrRedisSetOperationFailed represents a function to wrap failed to set key error and err.
	ErrRedisSetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to set key (%s)", key)
	}

	// ErrRedisSetOperationFailed represents a function to wrap failed to delete key error and err.
	ErrRedisDeleteOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to delete key (%s)", key)
	}

	// ErrRedisSetOperationFailed represents a function to generate an error that invalid configuration version.
	ErrInvalidConfigVersion = func(cur, con string) error {
		return Errorf("invalid config version %s not satisfies version constraints %s", cur, con)
	}

	// ErrRedisAddrsNotFound generates a new error of address not found.
	ErrRedisAddrsNotFound = New("error redis addrs not found")

	// ErrRedisConnectionPingFailed generates a new error of redis connection ping failed.
	ErrRedisConnectionPingFailed = New("error redis connection ping failed")
)

type ErrRedisNotFoundIdentity struct {
	err error
}

// Error returns the string of ErrRedisNotFoundIdentity.error.
func (e *ErrRedisNotFoundIdentity) Error() string {
	if e.err != nil {
		return e.err.Error()
	}
	return ""
}

// Unwrap returns the error value of ErrRedisNotFoundIdentity.
func (e *ErrRedisNotFoundIdentity) Unwrap() error {
	return e.err
}

// IsErrRedisNotFound compares the input error and ErrRedisNotFoundIdentity.error and returns true or false that is the result of errors.As.
func IsErrRedisNotFound(err error) bool {
	target := new(ErrRedisNotFoundIdentity)
	return As(err, &target)
}
