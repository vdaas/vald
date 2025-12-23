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

package tikv

import "google.golang.org/protobuf/encoding/protojson"

// MarshalJSON implements json.Marshaler
func (msg *Error) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *Error) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *GetAllStoresRequest) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *GetAllStoresRequest) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *GetAllStoresResponse) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *GetAllStoresResponse) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *Region) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *Region) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *KeyRange) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *KeyRange) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *BatchScanRegionsRequest) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *BatchScanRegionsRequest) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{}.Unmarshal(b, msg)
}

// MarshalJSON implements json.Marshaler
func (msg *BatchScanRegionsResponse) MarshalJSON() ([]byte, error) {
	return protojson.MarshalOptions{}.Marshal(msg)
}

// UnmarshalJSON implements json.Unmarshaler
func (msg *BatchScanRegionsResponse) UnmarshalJSON(b []byte) error {
	return protojson.UnmarshalOptions{}.Unmarshal(b, msg)
}
