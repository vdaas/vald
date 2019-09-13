//
// Copyright (C) 2019-2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//


// Package ngt provides implementation of Go API for https://github.com/yahoojapan/NGT
package ngt

/*
#cgo LDFLAGS: -lngt
#include <NGT/Capi.h>
*/
import "C"
import (
	"strings"
	"time"

	"github.com/vdaas/vald/internal/errors"
)

type Option func(*ngt) error

var (
	defaultOpts = []Option{
		WithIndexPath("/tmp/ngt-" + time.Now().Format(time.RFC3339)),
		WithDimension(0),
		WithCreationEdgeSize(10),
		WithSearchEdgeSize(40),
		WithObjectType(Float),
		WithDistanceType(L2),
		WithBulkInsertChunkSize(100),
	}
)

func WithIndexPath(path string) Option {
	return func(n *ngt) error {
		n.idxPath = path
		return nil
	}
}

func WithBulkInsertChunkSize(size int) Option {
	return func(n *ngt) error {
		n.bulkInsertChunkSize = size
		return nil
	}
}

func WithDimension(size int) Option {
	return func(n *ngt) error {
		if C.ngt_set_property_dimension(n.prop, C.int32_t(size), n.ebuf) == ErrorCode {
			return errors.ErrFailedToSetDimension(newGoError(n.ebuf))
		}
		n.dimension = C.int32_t(size)
		return nil
	}
}

func WithDistanceTypeByString(dt string) Option {
	var d distanceType
	switch strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToLower(dt)) {
	case "l1":
		d = L1
	case "l2":
		d = L2
	case "angle":
		d = Angle
	case "hamming":
		d = Hamming
	case "cosine":
		d = Cosine
	case "normalizedangle":
		d = NormalizedAngle
	case "normalizedcosine":
		d = NormalizedCosine
	}
	return WithDistanceType(d)
}

func WithDistanceType(t distanceType) Option {
	return func(n *ngt) error {
		switch t {
		case L1:
			if C.ngt_set_property_distance_type_l1(n.prop, n.ebuf) == ErrorCode {
				return errors.ErrFailedToSetDistanceType(newGoError(n.ebuf), "L1")
			}
		case L2:
			if C.ngt_set_property_distance_type_l2(n.prop, n.ebuf) == ErrorCode {
				return errors.ErrFailedToSetDistanceType(newGoError(n.ebuf), "L2")
			}
		case Angle:
			if C.ngt_set_property_distance_type_angle(n.prop, n.ebuf) == ErrorCode {
				return errors.ErrFailedToSetDistanceType(newGoError(n.ebuf), "Angle")
			}
		case Hamming:
			if C.ngt_set_property_distance_type_hamming(n.prop, n.ebuf) == ErrorCode {
				return errors.ErrFailedToSetDistanceType(newGoError(n.ebuf), "Hamming")
			}
		case Cosine:
			if C.ngt_set_property_distance_type_cosine(n.prop, n.ebuf) == ErrorCode {
				return errors.ErrFailedToSetDistanceType(newGoError(n.ebuf), "Cosine")
			}
		case NormalizedAngle:
			// TODO: not implemented in C API
			return errors.ErrFailedToSetDistanceType(newGoError(n.ebuf), "NormalizedAngle")
		case NormalizedCosine:
			// TODO: not implemented in C API
			return errors.ErrFailedToSetDistanceType(newGoError(n.ebuf), "NormalizedCosine")
		default:
			return errors.ErrUnsupportedDistanceType
		}
		return nil
	}
}

func WithObjectTypeByString(ot string) Option {
	var o objectType
	switch strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToLower(ot)) {
	case "uint8":
		o = Uint8
	case "float":
		o = Float
	}
	return WithObjectType(o)
}

func WithObjectType(t objectType) Option {
	return func(n *ngt) error {
		switch t {
		case Uint8:
			if C.ngt_set_property_object_type_integer(n.prop, n.ebuf) == ErrorCode {
				return errors.ErrFailedToSetObjectType(newGoError(n.ebuf), "Uint8")
			}
		case Float:
			if C.ngt_set_property_object_type_float(n.prop, n.ebuf) == ErrorCode {
				return errors.ErrFailedToSetObjectType(newGoError(n.ebuf), "Float")
			}
		default:
			return errors.ErrUnsupportedObjectType
		}
		n.objectType = t
		return nil
	}
}

func WithCreationEdgeSize(size int) Option {
	return func(n *ngt) error {
		if C.ngt_set_property_edge_size_for_creation(n.prop, C.int16_t(size), n.ebuf) == ErrorCode {
			return errors.ErrFailedToSetCreationEdgeSize(newGoError(n.ebuf))
		}
		return nil
	}
}

func WithSearchEdgeSize(size int) Option {
	return func(n *ngt) error {
		if C.ngt_set_property_edge_size_for_search(n.prop, C.int16_t(size), n.ebuf) == ErrorCode {
			return errors.ErrFailedToSetSearchEdgeSize(newGoError(n.ebuf))
		}
		return nil
	}
}
