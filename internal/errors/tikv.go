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
	ErrTiKVOptionFailed = func(err error) error {
		return Wrap(err, "TiKV option error")
	}

	ErrNewTiKVRawClientFailed = func(err error) error {
		return Wrap(err, "failed to create TiKV raw client")
	}

	ErrNewTiKVTxnClientFailed = func(err error) error {
		return Wrap(err, "failed to create TiKV txn client")
	}

	ErrTiKVBeginOperationFailed = func(err error) error {
		return Wrap(err, "failed to begin")
	}

	ErrTiKVSetOperationFailed = func(key, val []byte, err error) error {
		return Wrapf(err, "failed to set key (%s) - value (%s)", key, val)
	}

	ErrTiKVCommitOperationFailed = func(err error) error {
		return Wrap(err, "failed to commit")
	}

	ErrTiKVGetOperationFailed = func(key []byte, err error) error {
		return Wrapf(err, "failed to get key (%s)", key)
	}

	ErrTiKVDeleteOperationFailed = func(key []byte, err error) error {
		return Wrapf(err, "failed to delete key (%s)", key)
	}

	ErrTiKVRawClientCloseOperationFailed = func(err error) error {
		return Wrap(err, "failed to close TiKV raw client")
	}

	ErrTiKVTxnClientCloseOperationFailed = func(err error) error {
		return Wrap(err, "failed to close TiKV txn client")
	}
)
