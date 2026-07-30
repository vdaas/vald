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

package dataset

import (
	"bytes"
	"encoding/binary"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/os"
	"github.com/vdaas/vald/tests/v2/e2e/config"
)

// encodeFvecs serializes vectors into the fvecs on-disk layout: each vector
// is a 4-byte little-endian int32 dimension header followed by that many
// little-endian float32 elements.
func encodeFvecs(tb testing.TB, vectors [][]float32) []byte {
	tb.Helper()
	buf := new(bytes.Buffer)
	for _, v := range vectors {
		if err := binary.Write(buf, binary.LittleEndian, int32(len(v))); err != nil { //nolint:gosec // test vectors are small, len(v) always fits in int32
			tb.Fatalf("failed to write fvecs dimension header: %v", err)
		}
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			tb.Fatalf("failed to write fvecs vector data: %v", err)
		}
	}
	return buf.Bytes()
}

// encodeBvecs serializes vectors into the bvecs on-disk layout: each vector
// is a 4-byte little-endian int32 dimension header followed by that many
// uint8 elements.
func encodeBvecs(tb testing.TB, vectors [][]uint8) []byte {
	tb.Helper()
	buf := new(bytes.Buffer)
	for _, v := range vectors {
		if err := binary.Write(buf, binary.LittleEndian, int32(len(v))); err != nil { //nolint:gosec // test vectors are small, len(v) always fits in int32
			tb.Fatalf("failed to write bvecs dimension header: %v", err)
		}
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			tb.Fatalf("failed to write bvecs vector data: %v", err)
		}
	}
	return buf.Bytes()
}

// encodeIvecs serializes vectors into the ivecs on-disk layout: each vector
// is a 4-byte little-endian int32 dimension header followed by that many
// little-endian int32 elements.
func encodeIvecs(tb testing.TB, vectors [][]int32) []byte {
	tb.Helper()
	buf := new(bytes.Buffer)
	for _, v := range vectors {
		if err := binary.Write(buf, binary.LittleEndian, int32(len(v))); err != nil { //nolint:gosec // test vectors are small, len(v) always fits in int32
			tb.Fatalf("failed to write ivecs dimension header: %v", err)
		}
		if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
			tb.Fatalf("failed to write ivecs vector data: %v", err)
		}
	}
	return buf.Bytes()
}

// writeFile writes data to dir/name and returns the full path.
func writeFile(tb testing.TB, dir, name string, data []byte) string {
	tb.Helper()
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		tb.Fatalf("failed to create %s: %v", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			tb.Fatalf("failed to close %s: %v", path, err)
		}
	}()
	if _, err := f.Write(data); err != nil {
		tb.Fatalf("failed to write %s: %v", path, err)
	}
	return path
}

func TestToDataset(t *testing.T) {
	t.Parallel()

	t.Run("returns an error for a nil config", func(t *testing.T) {
		t.Parallel()
		_, err := ToDataset(nil)
		if !errors.Is(err, errors.ErrInvalidConfig) {
			t.Errorf("got err = %v, want wrapping %v", err, errors.ErrInvalidConfig)
		}
	})

	t.Run("returns an error for an unsupported extension", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		path := writeFile(t, dir, "dataset.bin", []byte{0x00})
		_, err := ToDataset(&config.Dataset{Name: path})
		if err == nil {
			t.Error("expected an error for an unsupported extension, got nil")
		}
	})

	t.Run("loads fvecs train/query/neighbors into a hdf5.Dataset", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		train := [][]float32{{1, 2, 3}, {4, 5, 6}}
		query := [][]float32{{7, 8, 9}}
		neighbors := [][]int32{{1, 0}}

		trainPath := writeFile(t, dir, "train.fvecs", encodeFvecs(t, train))
		queryPath := writeFile(t, dir, "query.fvecs", encodeFvecs(t, query))
		neighborsPath := writeFile(t, dir, "neighbors.ivecs", encodeIvecs(t, neighbors))

		got, err := ToDataset(&config.Dataset{
			Name:      trainPath,
			Query:     queryPath,
			Neighbors: neighborsPath,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got.Train, train) {
			t.Errorf("Train: got %v, want %v", got.Train, train)
		}
		if !reflect.DeepEqual(got.Test, query) {
			t.Errorf("Test: got %v, want %v", got.Test, query)
		}
		want := [][]int{{1, 0}}
		if !reflect.DeepEqual(got.Neighbors, want) {
			t.Errorf("Neighbors: got %v, want %v", got.Neighbors, want)
		}
	})

	t.Run("loads bvecs train/query, converting uint8 elements to float32", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		train := [][]uint8{{1, 2, 3}, {255, 0, 128}}
		query := [][]uint8{{9, 8, 7}}
		neighbors := [][]int32{{0}}

		trainPath := writeFile(t, dir, "train.bvecs", encodeBvecs(t, train))
		queryPath := writeFile(t, dir, "query.bvecs", encodeBvecs(t, query))
		neighborsPath := writeFile(t, dir, "neighbors.ivecs", encodeIvecs(t, neighbors))

		got, err := ToDataset(&config.Dataset{
			Name:      trainPath,
			Query:     queryPath,
			Neighbors: neighborsPath,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		wantTrain := [][]float32{{1, 2, 3}, {255, 0, 128}}
		if !reflect.DeepEqual(got.Train, wantTrain) {
			t.Errorf("Train: got %v, want %v", got.Train, wantTrain)
		}
		wantQuery := [][]float32{{9, 8, 7}}
		if !reflect.DeepEqual(got.Test, wantQuery) {
			t.Errorf("Test: got %v, want %v", got.Test, wantQuery)
		}
	})

	t.Run("returns an error when query is missing for an x1b format", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		trainPath := writeFile(t, dir, "train.fvecs", encodeFvecs(t, [][]float32{{1, 2}}))
		neighborsPath := writeFile(t, dir, "neighbors.ivecs", encodeIvecs(t, [][]int32{{0}}))
		_, err := ToDataset(&config.Dataset{Name: trainPath, Neighbors: neighborsPath})
		if err == nil {
			t.Error("expected an error when Query is unset, got nil")
		}
	})

	t.Run("returns an error when neighbors is missing for an x1b format", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		trainPath := writeFile(t, dir, "train.fvecs", encodeFvecs(t, [][]float32{{1, 2}}))
		queryPath := writeFile(t, dir, "query.fvecs", encodeFvecs(t, [][]float32{{3, 4}}))
		_, err := ToDataset(&config.Dataset{Name: trainPath, Query: queryPath})
		if err == nil {
			t.Error("expected an error when Neighbors is unset, got nil")
		}
	})

	t.Run("returns an error when the train file does not exist", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		queryPath := writeFile(t, dir, "query.fvecs", encodeFvecs(t, [][]float32{{1, 2}}))
		neighborsPath := writeFile(t, dir, "neighbors.ivecs", encodeIvecs(t, [][]int32{{0}}))
		_, err := ToDataset(&config.Dataset{
			Name:      filepath.Join(dir, "does-not-exist.fvecs"),
			Query:     queryPath,
			Neighbors: neighborsPath,
		})
		if err == nil {
			t.Error("expected an error for a missing train file, got nil")
		}
	})
}
