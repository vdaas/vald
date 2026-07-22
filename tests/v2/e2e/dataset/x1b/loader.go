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

package x1b

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/vdaas/vald/internal/errors"
)

const (
	// headerSize is the width, in bytes, of the little-endian int32
	// dimension header that precedes every vector record in
	// fvecs/bvecs/ivecs data.
	headerSize = 4
	// int32Size and float32Size are both 4 bytes; ivecs and fvecs payload
	// elements share this width, distinct from bvecs' 1-byte uint8 elements.
	int32Size   = 4
	float32Size = 4
	uint8Size   = 1
)

// ErrTruncatedData is returned when the input ends before a vector's
// declared dimension can be fully read, including the case where a
// trailing dimension header itself is shorter than 4 bytes.
var ErrTruncatedData = errors.New("x1b: truncated vector data")

// readDimension reads the next record's 4-byte little-endian int32 dimension
// header from r. ok is false with a nil error only when r is exhausted
// exactly on a record boundary, which is the sole condition under which
// parsing is allowed to stop cleanly; any other short read is a truncated
// header and is reported as ErrTruncatedData.
func readDimension(r io.Reader) (dim int32, ok bool, err error) {
	var header [headerSize]byte
	n, err := io.ReadFull(r, header[:])
	if err != nil {
		if errors.Is(err, io.EOF) && n == 0 {
			return 0, false, nil
		}
		return 0, false, ErrTruncatedData
	}
	// The on-disk header is itself a little-endian int32; reinterpreting its
	// bits via Uint32 and converting back to int32 recovers the exact same
	// value, so this is a bit-reinterpretation, not a lossy numeric narrowing.
	return int32(binary.LittleEndian.Uint32(header[:])), true, nil //nolint:gosec // bit-for-bit reinterpretation of an on-disk int32 header
}

// readPayload reads the dim*elementSize bytes of vector payload that follow a
// dimension header, reporting ErrTruncatedData if r ends early.
func readPayload(r io.Reader, dim int32, elementSize int) ([]byte, error) {
	buf := make([]byte, int(dim)*elementSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, ErrTruncatedData
	}
	return buf, nil
}

// ParseFvecs decodes fvecs-formatted data read from r: a sequence of records,
// each a 4-byte little-endian int32 dimension header followed by that many
// little-endian float32 elements, repeated until io.EOF.
func ParseFvecs(r io.Reader) ([][]float32, error) {
	vecs := make([][]float32, 0)
	for {
		dim, ok, err := readDimension(r)
		if err != nil {
			return nil, err
		}
		if !ok {
			return vecs, nil
		}
		buf, err := readPayload(r, dim, float32Size)
		if err != nil {
			return nil, err
		}
		v := make([]float32, dim)
		for i := range v {
			v[i] = math.Float32frombits(binary.LittleEndian.Uint32(buf[i*float32Size:]))
		}
		vecs = append(vecs, v)
	}
}

// ParseBvecs decodes bvecs-formatted data read from r: a sequence of records,
// each a 4-byte little-endian int32 dimension header followed by that many
// uint8 elements, repeated until io.EOF.
func ParseBvecs(r io.Reader) ([][]uint8, error) {
	vecs := make([][]uint8, 0)
	for {
		dim, ok, err := readDimension(r)
		if err != nil {
			return nil, err
		}
		if !ok {
			return vecs, nil
		}
		buf, err := readPayload(r, dim, uint8Size)
		if err != nil {
			return nil, err
		}
		v := make([]uint8, dim)
		copy(v, buf)
		vecs = append(vecs, v)
	}
}

// ParseIvecs decodes ivecs-formatted data read from r: a sequence of records,
// each a 4-byte little-endian int32 dimension header followed by that many
// little-endian int32 elements, repeated until io.EOF. ivecs files are
// typically used to store groundtruth neighbor indices.
func ParseIvecs(r io.Reader) ([][]int32, error) {
	vecs := make([][]int32, 0)
	for {
		dim, ok, err := readDimension(r)
		if err != nil {
			return nil, err
		}
		if !ok {
			return vecs, nil
		}
		buf, err := readPayload(r, dim, int32Size)
		if err != nil {
			return nil, err
		}
		v := make([]int32, dim)
		for i := range v {
			// Bit-for-bit reinterpretation of an on-disk int32 element, not a
			// lossy numeric narrowing; see readDimension's comment above.
			v[i] = int32(binary.LittleEndian.Uint32(buf[i*int32Size:])) //nolint:gosec // bit-for-bit reinterpretation of an on-disk int32 element
		}
		vecs = append(vecs, v)
	}
}
