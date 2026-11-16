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

package metadata

import (
	"io/fs"
	"path/filepath"

	"github.com/vdaas/vald/internal/encoding/json"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/os"
)

const (
	AgentMetadataFileName = "metadata.json"
)

type Metadata struct {
	NGT       *NGT   `json:"ngt,omitempty"   yaml:"ngt"`
	Faiss     *Faiss `json:"faiss,omitempty" yaml:"faiss"`
	IsInvalid bool   `json:"is_invalid"      yaml:"is_invalid"`
}

type NGT struct {
	IndexCount uint64 `json:"index_count" yaml:"index_count"`
}

type Faiss struct {
	IndexCount uint64 `json:"index_count" yaml:"index_count"`
}

func Load(path string) (meta *Metadata, err error) {
	var fi os.FileInfo
	exists, fi, err := file.ExistsWithDetail(path)
	switch {
	case !exists:
		return nil, errors.ErrMetadataFileNotFound
	case err != nil:
		return nil, err
	case fi == nil || fi.Size() == 0:
		return nil, errors.ErrMetadataFileEmpty
	}
	f, err := file.Open(path, os.O_RDONLY|os.O_SYNC, fs.ModePerm)
	if err != nil {
		return nil, err
	}
	defer func() {
		if f != nil {
			derr := f.Close()
			if derr != nil {
				err = errors.Join(err, derr)
			}
		}
	}()

	err = json.Decode(f, &meta)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	// A whitespace-only (io.EOF) or "null" file decodes to a nil *Metadata; report it as empty
	// so callers never dereference a nil result.
	if meta == nil {
		return nil, errors.ErrMetadataFileEmpty
	}
	return meta, nil
}

// Store writes meta to path atomically: it writes a temporary file in the same directory and
// renames it over path, so a crash never leaves an empty or partial metadata file and watchers
// see the new file (fsnotify Create) only once its content is complete.
// The temporary file name does not end with ".json", so it is never mistaken for the metadata file.
func Store(path string, meta *Metadata) (err error) {
	tmp := path + tmpFileSuffix
	var f *os.File
	f, err = file.Open(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fs.ModePerm)
	if err != nil {
		return err
	}
	defer func() {
		if f != nil {
			derr := f.Close()
			if derr != nil {
				err = errors.Join(err, derr)
			}
		}
		if err != nil {
			rerr := os.Remove(tmp)
			if rerr != nil && !errors.Is(rerr, fs.ErrNotExist) {
				err = errors.Join(err, rerr)
			}
		}
	}()

	err = json.Encode(f, &meta)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}
	err = f.Close()
	f = nil
	if err != nil {
		return err
	}

	err = os.Rename(tmp, path)
	if err != nil {
		return err
	}

	// Persist the rename itself. The metadata is already complete on disk at this point,
	// so a failure here is only logged.
	if serr := syncDir(filepath.Dir(path)); serr != nil {
		log.Warnf("failed to sync directory of metadata file %s: %v", path, serr)
	}
	return nil
}

const tmpFileSuffix = ".tmp"

func syncDir(dir string) (err error) {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := d.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}()
	return d.Sync()
}
