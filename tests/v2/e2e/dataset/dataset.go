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

// Package dataset dispatches tests/v2/e2e/config.Dataset to the loader that
// matches its file format, converting the result to the common
// tests/v2/e2e/hdf5.Dataset shape (Train/Test/Neighbors) so that callers such
// as tests/v2/e2e/crud do not need to know which on-disk format a given
// dataset config uses.
//
// The x1b branch below eagerly decodes its files into [][]T slices via
// tests/v2/e2e/dataset/x1b, the same way hdf5.ToDataset eagerly reads an
// entire hdf5 file. That is appropriate for the small/medium ann-benchmarks
// fixtures this package targets, but not for the actual billion-scale
// SIFT1B/DEEP1B corpora hack/benchmark/internal/assets.loadLargeData streams
// via mmap; wiring a lazy, billion-scale-aware Dataset variant into e2e CI is
// out of scope here and left to whichever task migrates that lazy loader
// (see hack/benchmark/internal/assets/large_dataset.go).
package dataset

import (
	"path/filepath"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/internal/strings"
	"github.com/vdaas/vald/tests/v2/e2e/config"
	"github.com/vdaas/vald/tests/v2/e2e/dataset/x1b"
	"github.com/vdaas/vald/tests/v2/e2e/hdf5"
)

// ToDataset loads cfg into a *hdf5.Dataset, dispatching on the file extension
// of cfg.Name: ".hdf5"/".h5" is read via hdf5.ToDataset, while
// ".fvecs"/".bvecs" is read via tests/v2/e2e/dataset/x1b, additionally
// requiring cfg.Query (same encoding as cfg.Name) and cfg.Neighbors (an
// ivecs groundtruth file) to be set.
func ToDataset(cfg *config.Dataset) (*hdf5.Dataset, error) {
	if cfg == nil {
		return nil, errors.Wrap(errors.ErrInvalidConfig, "dataset: nil config")
	}
	switch ext := strings.ToLower(filepath.Ext(cfg.Name)); ext {
	case ".hdf5", ".h5":
		return hdf5.ToDataset(cfg.Name)
	case ".fvecs", ".bvecs":
		return x1bToDataset(cfg, ext)
	default:
		return nil, errors.Errorf("dataset: unsupported dataset file extension %q for %s", ext, cfg.Name)
	}
}

// x1bToDataset loads the multi-file fvecs/bvecs (x1b) format: cfg.Name is the
// train vectors file, cfg.Query the query vectors file (same ext as Name),
// and cfg.Neighbors an ivecs groundtruth file.
func x1bToDataset(cfg *config.Dataset, ext string) (*hdf5.Dataset, error) {
	if cfg.Query == "" || cfg.Neighbors == "" {
		return nil, errors.Errorf(
			"dataset: x1b format %q requires both query and neighbors files to be set, got query=%q neighbors=%q",
			ext, cfg.Query, cfg.Neighbors,
		)
	}

	train, err := loadX1BVectors(cfg.Name, ext)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load x1b train vectors from %s", cfg.Name)
	}
	query, err := loadX1BVectors(cfg.Query, ext)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load x1b query vectors from %s", cfg.Query)
	}
	neighbors, err := loadX1BNeighbors(cfg.Neighbors)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load x1b groundtruth neighbors from %s", cfg.Neighbors)
	}

	return &hdf5.Dataset{
		Train:     train,
		Test:      query,
		Neighbors: neighbors,
	}, nil
}

// loadX1BVectors opens name and parses it as fvecs or bvecs data (per ext),
// converting bvecs' uint8 elements to float32 so both formats share the
// [][]float32 shape hdf5.Dataset.Train/Test expect.
func loadX1BVectors(name, ext string) (vecs [][]float32, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	switch ext {
	case ".fvecs":
		return x1b.ParseFvecs(f)
	case ".bvecs":
		raw, perr := x1b.ParseBvecs(f)
		if perr != nil {
			return nil, perr
		}
		vecs = make([][]float32, len(raw))
		for i, v := range raw {
			fv := make([]float32, len(v))
			for j, e := range v {
				fv[j] = float32(e)
			}
			vecs[i] = fv
		}
		return vecs, nil
	default:
		return nil, errors.Errorf("dataset: unsupported x1b vector extension %q", ext)
	}
}

// loadX1BNeighbors opens name and parses it as ivecs groundtruth neighbor
// indices, converting int32 elements to int to match hdf5.Dataset.Neighbors.
func loadX1BNeighbors(name string) (neighbors [][]int, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	raw, err := x1b.ParseIvecs(f)
	if err != nil {
		return nil, err
	}
	neighbors = make([][]int, len(raw))
	for i, v := range raw {
		nv := make([]int, len(v))
		for j, e := range v {
			nv[j] = int(e)
		}
		neighbors[i] = nv
	}
	return neighbors, nil
}
