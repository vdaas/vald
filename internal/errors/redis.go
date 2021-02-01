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

	// Redis.
	ErrRedisInvalidKVVKPrefix = func(kv, vk string) error {
		return Errorf("kv index and vk prefix must be defferent.\t(kv: %s,\tvk: %s)", kv, vk)
	}

	NewErrRedisNotFoundIdentity = func() error {
		return &ErrRedisNotFoundIdentity{
			err: New("error redis entry not found"),
		}
	}

	ErrRedisNotFound = func(key string) error {
		return Wrapf(NewErrRedisNotFoundIdentity(), "error redis key '%s' not found", key)
	}

	ErrRedisInvalidOption = New("error redis invalid option")

	ErrRedisGetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to fetch key (%s)", key)
	}

	ErrRedisSetOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to set key (%s)", key)
	}

	ErrRedisDeleteOperationFailed = func(key string, err error) error {
		return Wrapf(err, "Failed to delete key (%s)", key)
	}

	ErrInvalidConfigVersion = func(cur, con string) error {
		return Errorf("invalid config version %s not satisfies version constraints %s", cur, con)
	}

	ErrRedisAddrsNotFound = New("addrs not found")

	ErrRedisConnectionPingFailed = New("error Redis connection ping failed")
)

type ErrRedisNotFoundIdentity struct {
	err error
}

func (e *ErrRedisNotFoundIdentity) Error() string {
	return e.err.Error()
}

func (e *ErrRedisNotFoundIdentity) Unwrap() error {
	return e.err
}

func IsErrRedisNotFound(err error) bool {
	target := new(ErrRedisNotFoundIdentity)
	return As(err, &target)
}
