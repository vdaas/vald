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

// Package metadata provides agent metadata structs and info.
package metadata

import (
	"io"
	"os"

	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/file"
)

const (
	AgentMetadataFileName = "metadata.json"
)

type Metadata struct {
	IsInvalid bool `json:"is_invalid" yaml:"is_invalid"`
	NGT       *NGT `json:"ngt,omitempty" yaml:"ngt"`
}

type NGT struct {
	IndexCount uint64 `json:"index_count" yaml:"index_count"`
}

func Load(path string) (*Metadata, error) {
	f, err := file.Open(path, os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var meta Metadata
	err = json.Decode(f, &meta)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return &meta, nil
}

func Store(path string, meta *Metadata) error {
	f, err := file.Open(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.Encode(f, &meta)
}
