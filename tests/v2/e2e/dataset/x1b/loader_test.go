//
// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
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

// This is the RED phase of a TDD hand-off (see AGENTS.md / SWARM.md Fixer
// pattern): every case below is written to fail against the loader.go stubs
// (which unconditionally return errNotImplemented) and is expected to start
// passing once ParseFvecs/ParseBvecs/ParseIvecs gain a real implementation.
// No `//go:build ignore` gate is needed because loader.go already provides
// buildable stub signatures for every symbol this file references.
//
// Package placement/tag rationale (task step 4): this file intentionally has
// no `//go:build e2e` tag. The fvecs/bvecs/ivecs binary format is a pure,
// cluster-independent decoding problem -- exactly like the sibling
// tests/v2/e2e/metrics package, which is also untagged and tested via
// `go test ./...` without a cluster or the `-tags e2e` gonum/hdf5 cgo
// dependency pulled in by tests/v2/e2e/hdf5 and tests/v2/e2e/config. Keeping
// this parser untagged lets it run in ordinary CI unit-test jobs; only the
// eventual glue code that wires it into an e2e Dataset needs the e2e tag.
package x1b

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/vdaas/vald/internal/test"
)

// encodeFvecs builds a raw fvecs byte buffer from vectors: each vector is
// serialized as its own 4-byte little-endian int32 dimension header followed
// by len(vector) little-endian float32 elements, matching the layout
// hack/benchmark/assets/x1b.openFile mmaps directly (skipcq: GSC-G103 there).
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

// encodeBvecs builds a raw bvecs byte buffer from vectors: each vector is
// serialized as its own 4-byte little-endian int32 dimension header followed
// by len(vector) uint8 elements.
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

// encodeIvecs builds a raw ivecs byte buffer from vectors: each vector is
// serialized as its own 4-byte little-endian int32 dimension header followed
// by len(vector) little-endian int32 elements. ivecs is typically used for
// groundtruth neighbor indices (see hack/benchmark/internal/assets.loadLargeData).
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

func TestParseFvecs(t *testing.T) {
	type args struct {
		data []byte
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, a args) ([][]float32, error) {
		t.Helper()
		return ParseFvecs(bytes.NewReader(a.data))
	}, []test.Case[[][]float32, args]{
		{
			Name: "parses multiple fvecs vectors and preserves dimension and values",
			Args: args{
				data: encodeFvecs(t, [][]float32{
					{1.5, 2.5, 3.5, 4.5},
					{-1, 0, 1, 2},
					{0.125, 0.25, 0.5, 1},
				}),
			},
			Want: test.Result[[][]float32]{
				Val: [][]float32{
					{1.5, 2.5, 3.5, 4.5},
					{-1, 0, 1, 2},
					{0.125, 0.25, 0.5, 1},
				},
			},
		},
		{
			Name: "returns ErrTruncatedData when the declared dimension exceeds the remaining bytes",
			Args: args{
				// header declares dim=4 (16 bytes of float32 payload) but the
				// buffer is sliced down to only 8 bytes (2 float32s) of it.
				data: encodeFvecs(t, [][]float32{{1, 2, 3, 4}})[:4+8],
			},
			Want: test.Result[[][]float32]{
				Err: ErrTruncatedData,
			},
		},
		{
			Name: "returns ErrTruncatedData when a trailing partial dimension header remains",
			Args: args{
				// 2 complete vectors followed by 2 stray bytes: too short to
				// be a 4-byte dimension header and too short to be data.
				data: append(encodeFvecs(t, [][]float32{{1, 2}, {3, 4}}), 0x01, 0x02),
			},
			Want: test.Result[[][]float32]{
				Err: ErrTruncatedData,
			},
		},
		{
			Name: "returns an empty, non-nil result for empty input",
			Args: args{
				data: []byte{},
			},
			Want: test.Result[[][]float32]{
				Val: [][]float32{},
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestParseBvecs(t *testing.T) {
	type args struct {
		data []byte
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, a args) ([][]uint8, error) {
		t.Helper()
		return ParseBvecs(bytes.NewReader(a.data))
	}, []test.Case[[][]uint8, args]{
		{
			Name: "parses multiple bvecs vectors and preserves dimension and values",
			Args: args{
				data: encodeBvecs(t, [][]uint8{
					{1, 2, 3, 4},
					{255, 0, 128, 64},
				}),
			},
			Want: test.Result[[][]uint8]{
				Val: [][]uint8{
					{1, 2, 3, 4},
					{255, 0, 128, 64},
				},
			},
		},
		{
			Name: "returns ErrTruncatedData when the declared dimension exceeds the remaining bytes",
			Args: args{
				// header declares dim=4 (4 bytes of uint8 payload) but the
				// buffer is sliced down to only 2 bytes of it.
				data: encodeBvecs(t, [][]uint8{{1, 2, 3, 4}})[:4+2],
			},
			Want: test.Result[[][]uint8]{
				Err: ErrTruncatedData,
			},
		},
		{
			Name: "returns an empty, non-nil result for empty input",
			Args: args{
				data: []byte{},
			},
			Want: test.Result[[][]uint8]{
				Val: [][]uint8{},
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}

func TestParseIvecs(t *testing.T) {
	type args struct {
		data []byte
	}

	if err := test.Run(t.Context(), t, func(t *testing.T, a args) ([][]int32, error) {
		t.Helper()
		return ParseIvecs(bytes.NewReader(a.data))
	}, []test.Case[[][]int32, args]{
		{
			Name: "parses multiple ivecs groundtruth neighbor-index vectors",
			Args: args{
				data: encodeIvecs(t, [][]int32{
					{10, 42, 7},
					{0, 1, 2, 3},
				}),
			},
			Want: test.Result[[][]int32]{
				Val: [][]int32{
					{10, 42, 7},
					{0, 1, 2, 3},
				},
			},
		},
		{
			Name: "returns ErrTruncatedData when the declared dimension exceeds the remaining bytes",
			Args: args{
				// header declares dim=3 (12 bytes of int32 payload) but the
				// buffer is sliced down to only 4 bytes (1 int32) of it.
				data: encodeIvecs(t, [][]int32{{10, 42, 7}})[:4+4],
			},
			Want: test.Result[[][]int32]{
				Err: ErrTruncatedData,
			},
		},
		{
			Name: "returns an empty, non-nil result for empty input",
			Args: args{
				data: []byte{},
			},
			Want: test.Result[[][]int32]{
				Val: [][]int32{},
			},
		},
	}...); err != nil {
		t.Error(err)
	}
}
