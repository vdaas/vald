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

package json

import (
	json "github.com/goccy/go-json"
	"github.com/vdaas/vald/internal/io"
)

func Encode(w io.Writer, data interface{}) (err error) {
	return json.NewEncoder(w).Encode(data)
}

func Decode(r io.Reader, data interface{}) (err error) {
	return json.NewDecoder(r).Decode(data)
}

func Unmarshal(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}

func Marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func MarshalIndent(data interface{}, pref, ind string) ([]byte, error) {
	return json.MarshalIndent(data, pref, ind)
}
