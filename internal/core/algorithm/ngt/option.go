//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
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

// Package ngt provides implementation of Go API for https://github.com/yahoojapan/NGT
package ngt

/*
#cgo LDFLAGS: -lngt
#include <NGT/Capi.h>
*/
import "C"

import (
	"strings"

	"github.com/kpango/fastime"
	"github.com/vdaas/vald/internal/errors"
)

// Option represents the functional option for NGT.
type Option func(*ngt) error

var (
	DefaultPoolSize = uint32(10000)
	DefaultRadius   = float32(-1.0)
	DefaultEpsilon  = float32(0.01)

	defaultOptions = []Option{
		WithIndexPath("/tmp/ngt-" + string(fastime.FormattedNow())),
		WithDimension(minimumDimensionSize),
		WithDefaultRadius(DefaultRadius),
		WithDefaultEpsilon(DefaultEpsilon),
		WithDefaultPoolSize(DefaultPoolSize),
		WithCreationEdgeSize(10),
		WithSearchEdgeSize(40),
		WithObjectType(Float),
		WithDistanceType(L2),
		WithBulkInsertChunkSize(100),
	}
)

// WithInMemoryMode represents the option to set to start in memory mode or not for NGT.
func WithInMemoryMode(flg bool) Option {
	return func(n *ngt) error {
		n.inMemory = flg
		return nil
	}
}

// WithIndexPath represents the option to set the index path for NGT.
func WithIndexPath(path string) Option {
	return func(n *ngt) error {
		if len(path) == 0 {
			return errors.NewErrUnsetOption("indexPath")
		}
		n.idxPath = path
		return nil
	}
}

// WithBulkInsertChunkSize represents the option to set the bulk insert chunk size for NGT.
func WithBulkInsertChunkSize(size int) Option {
	return func(n *ngt) error {
		if size < 0 {
			return errors.NewErrInvalidOption("BulkInsertChunkSize", size)
		}
		n.bulkInsertChunkSize = size
		return nil
	}
}

// WithDimension represents the option to set the dimension for NGT.
func WithDimension(size int) Option {
	return func(n *ngt) error {
		if size > ngtVectorDimensionSizeLimit || size < minimumDimensionSize {
			err := errors.ErrInvalidDimensionSize(size, ngtVectorDimensionSizeLimit)
			return errors.NewErrCriticalOption("dimension", size, err)
		}

		if C.ngt_set_property_dimension(n.prop, C.int32_t(size), n.ebuf) == ErrorCode {
			err := errors.ErrFailedToSetDimension(n.newGoError(n.ebuf))
			return errors.NewErrCriticalOption("dimension", size, err)
		}

		n.dimension = C.int32_t(size)

		return nil
	}
}

// WithDistanceTypeByString represents the option to set the distance type for NGT.
func WithDistanceTypeByString(dt string) Option {
	var d distanceType
	switch strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToLower(dt)) {
	case "l1":
		d = L1
	case "l2":
		d = L2
	case "angle", "ang":
		d = Angle
	case "hamming", "ham":
		d = Hamming
	case "cosine", "cos":
		d = Cosine
	case "normalizedangle", "normalizedang", "nang", "nangle":
		d = NormalizedAngle
	case "normalizedcosine", "normalizedcos", "ncos", "ncosine":
		d = NormalizedCosine
	case "jaccard", "jac":
		d = Jaccard
	}
	return WithDistanceType(d)
}

// WithDistanceType represents the option to set the distance type for NGT.
func WithDistanceType(t distanceType) Option {
	return func(n *ngt) error {
		switch t {
		case L1:
			if C.ngt_set_property_distance_type_l1(n.prop, n.ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(n.ebuf), "L1")
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case L2:
			if C.ngt_set_property_distance_type_l2(n.prop, n.ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(n.ebuf), "L2")
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case Angle:
			if C.ngt_set_property_distance_type_angle(n.prop, n.ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(n.ebuf), "Angle")
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case Hamming:
			if C.ngt_set_property_distance_type_hamming(n.prop, n.ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(n.ebuf), "Hamming")
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case Cosine:
			if C.ngt_set_property_distance_type_cosine(n.prop, n.ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(n.ebuf), "Cosine")
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case NormalizedAngle:
			if C.ngt_set_property_distance_type_normalized_angle(n.prop, n.ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(n.ebuf), "NormalizedAngle")
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case NormalizedCosine:
			if C.ngt_set_property_distance_type_normalized_cosine(n.prop, n.ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(n.ebuf), "NormalizedCosine")
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case Jaccard:
			if C.ngt_set_property_distance_type_jaccard(n.prop, n.ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(n.ebuf), "Jaccard")
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		default:
			err := errors.ErrUnsupportedDistanceType
			return errors.NewErrCriticalOption("distanceType", t, err)
		}
		return nil
	}
}

// WithObjectTypeByString represents the option to set the object type by string for NGT.
func WithObjectTypeByString(ot string) Option {
	var o objectType
	switch strings.NewReplacer("-", "", "_", "", " ", "", "double", "float").Replace(strings.ToLower(ot)) {
	case "uint8":
		o = Uint8
	case "float":
		o = Float
	}
	return WithObjectType(o)
}

// WithObjectType represents the option to set the object type for NGT.
func WithObjectType(t objectType) Option {
	return func(n *ngt) error {
		switch t {
		case Uint8:
			if C.ngt_set_property_object_type_integer(n.prop, n.ebuf) == ErrorCode {
				err := errors.ErrFailedToSetObjectType(n.newGoError(n.ebuf), "Uint8")
				return errors.NewErrCriticalOption("objectType", t, err)
			}
		case Float:
			if C.ngt_set_property_object_type_float(n.prop, n.ebuf) == ErrorCode {
				err := errors.ErrFailedToSetObjectType(n.newGoError(n.ebuf), "Float")
				return errors.NewErrCriticalOption("objectType", t, err)
			}
		default:
			err := errors.ErrUnsupportedObjectType
			return errors.NewErrCriticalOption("objectType", t, err)
		}
		n.objectType = t
		return nil
	}
}

// WithCreationEdgeSize represents the option to set the creation edge size for NGT.
func WithCreationEdgeSize(size int) Option {
	return func(n *ngt) error {
		if C.ngt_set_property_edge_size_for_creation(n.prop, C.int16_t(size), n.ebuf) == ErrorCode {
			err := errors.ErrFailedToSetCreationEdgeSize(n.newGoError(n.ebuf))
			return errors.NewErrCriticalOption("creationEdgeSize", size, err)
		}
		return nil
	}
}

// WithSearchEdgeSize represents the option to set the search edge size for NGT.
func WithSearchEdgeSize(size int) Option {
	return func(n *ngt) error {
		if C.ngt_set_property_edge_size_for_search(n.prop, C.int16_t(size), n.ebuf) == ErrorCode {
			err := errors.ErrFailedToSetSearchEdgeSize(n.newGoError(n.ebuf))
			return errors.NewErrCriticalOption("searchEdgeSize", size, err)
		}
		return nil
	}
}

// WithDefaultPoolSize represents the option to set the default pool size for NGT.
func WithDefaultPoolSize(poolSize uint32) Option {
	return func(n *ngt) error {
		if poolSize == 0 {
			return errors.NewErrInvalidOption("defaultPoolSize", poolSize)
		}
		n.poolSize = poolSize
		return nil
	}
}

// WithDefaultRadius represents the option to set the default radius for NGT.
func WithDefaultRadius(radius float32) Option {
	return func(n *ngt) error {
		if radius == 0 {
			return errors.NewErrInvalidOption("defaultRadius", radius)
		}
		n.radius = radius
		return nil
	}
}

// WithDefaultEpsilon represents the option to set the default epsilon for NGT.
func WithDefaultEpsilon(epsilon float32) Option {
	return func(n *ngt) error {
		if epsilon == 0 {
			return errors.NewErrInvalidOption("defaultEpsilon", epsilon)
		}
		n.epsilon = epsilon
		return nil
	}
}
