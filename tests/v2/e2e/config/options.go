//go:build e2e

// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
)

// AlgorithmOptions contains protobuf JSON Any values, including their @type URLs.
// These are algorithm payloads, separate from the scenario's gRPC metadata.
type AlgorithmOptions []map[string]any

// Any validates and encodes algorithm-specific options without dropping unknown fields.
func (options AlgorithmOptions) Any() ([]*anypb.Any, error) {
	if len(options) == 0 {
		return nil, nil
	}
	out := make([]*anypb.Any, 0, len(options))
	for i, option := range options {
		data, err := json.Marshal(option)
		if err != nil {
			return nil, errors.Wrapf(err, "marshal algorithm option %d", i)
		}
		value := new(anypb.Any)
		if err := protojson.Unmarshal(data, value); err != nil {
			return nil, errors.Wrapf(err, "decode algorithm option %d", i)
		}
		out = append(out, value)
	}
	return out, nil
}
