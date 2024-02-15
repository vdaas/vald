//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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

// Package qbg provides implementation of Go API for https://github.com/yahoojapan/QBG
package qbg

/*
#cgo LDFLAGS: -lngt
#include <NGT/Capi.h>
*/
import "C"

import (
	"strconv"

	"github.com/kpango/fastime"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/strings"
)

// Option represents the functional option for QBG.
type Option func(*qbg) error

var (
	DefaultPoolSize = uint32(10000)
	DefaultRadius   = float32(-1.0)
	DefaultEpsilon  = float32(0.1)

	defaultOptions = []Option{
		WithIndexPath("/tmp/qbg-" + strconv.FormatInt(fastime.UnixNanoNow(), 10)),
		WithDimension(algorithm.MinimumVectorDimensionSize),
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

// WithInMemoryMode represents the option to set to start in memory mode or not for QBG.
func WithInMemoryMode(flg bool) Option {
	return func(n *qbg) error {
		n.inMemory = flg
		return nil
	}
}

// WithIndexPath represents the option to set the index path for QBG.
func WithIndexPath(path string) Option {
	return func(n *qbg) error {
		if len(path) == 0 {
			return errors.NewErrIgnoredOption("indexPath")
		}
		n.idxPath = path
		return nil
	}
}

// WithBulkInsertChunkSize represents the option to set the bulk insert chunk size for QBG.
func WithBulkInsertChunkSize(size int) Option {
	return func(n *qbg) error {
		if size < 0 {
			return errors.NewErrInvalidOption("BulkInsertChunkSize", size)
		}
		n.bulkInsertChunkSize = size
		return nil
	}
}

// WithDimension represents the option to set the dimension for QBG.
func WithDimension(size int) Option {
	return func(n *qbg) error {
		if size > algorithm.MaximumVectorDimensionSize || size < algorithm.MinimumVectorDimensionSize {
			err := errors.ErrInvalidDimensionSize(size, algorithm.MaximumVectorDimensionSize)
			return errors.NewErrCriticalOption("dimension", size, err)
		}

		ebuf := n.GetErrorBuffer()
		if C.qbg_set_property_dimension(n.prop, C.int32_t(size), ebuf) == ErrorCode {
			err := errors.ErrFailedToSetDimension(n.newGoError(ebuf))
			return errors.NewErrCriticalOption("dimension", size, err)
		}
		n.PutErrorBuffer(ebuf)

		n.dimension = C.int32_t(size)

		return nil
	}
}

// WithDistanceTypeByString represents the option to set the distance type for QBG.
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
	case "poincare", "poinc", "poi", "po", "pc":
		d = Poincare
	case "lorenz", "loren", "lor", "lo", "lz":
		d = Lorentz
	case "jaccard", "jac":
		d = Jaccard
	case "sparsejaccard", "sparsejac", "spjac", "sjc", "sj":
		d = SparseJaccard
	case "normalizedl2", "norml2", "nol2", "nl2":
		d = NormalizedL2
	case "normalizedangle", "normalizedang", "normang", "nang", "nangle":
		d = NormalizedAngle
	case "normalizedcosine", "normalizedcos", "normcos", "ncos", "ncosine":
		d = NormalizedCosine
	}
	return WithDistanceType(d)
}

// WithDistanceType represents the option to set the distance type for QBG.
func WithDistanceType(t distanceType) Option {
	return func(n *qbg) error {
		ebuf := n.GetErrorBuffer()
		switch t {
		case L1:
			if C.qbg_set_property_distance_type_l1(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case L2:
			if C.qbg_set_property_distance_type_l2(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case Angle:
			if C.qbg_set_property_distance_type_angle(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case Hamming:
			if C.qbg_set_property_distance_type_hamming(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case Cosine:
			if C.qbg_set_property_distance_type_cosine(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case Poincare:
			if C.qbg_set_property_distance_type_poincare(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case Lorentz:
			if C.qbg_set_property_distance_type_lorentz(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case Jaccard:
			if C.qbg_set_property_distance_type_jaccard(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case SparseJaccard:
			if C.qbg_set_property_distance_type_sparse_jaccard(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case NormalizedL2:
			if C.qbg_set_property_distance_type_normalized_l2(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case NormalizedAngle:
			if C.qbg_set_property_distance_type_normalized_angle(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		case NormalizedCosine:
			if C.qbg_set_property_distance_type_normalized_cosine(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetDistanceType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("distanceType", t, err)
			}
		default:
			err := errors.ErrUnsupportedDistanceType
			n.PutErrorBuffer(ebuf)
			return errors.NewErrCriticalOption("distanceType", t, err)
		}
		n.PutErrorBuffer(ebuf)
		return nil
	}
}

// WithObjectTypeByString represents the option to set the object type by string for QBG.
func WithObjectTypeByString(ot string) Option {
	var o objectType
	switch strings.NewReplacer("-", "", "_", "", " ", "", "double", "float").Replace(strings.ToLower(ot)) {
	case "uint8", "ui8", "u8":
		o = Uint8
	case "float", "float32", "f", "f32", "fp32":
		o = Float
	case "float16", "halfFloat", "half_float", "f16", "fp16":
		o = HalfFloat
	}
	return WithObjectType(o)
}

// WithObjectType represents the option to set the object type for QBG.
func WithObjectType(t objectType) Option {
	return func(n *qbg) error {
		ebuf := n.GetErrorBuffer()
		switch t {
		case Uint8:
			if C.qbg_set_property_object_type_integer(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetObjectType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("objectType", t, err)
			}
		case HalfFloat:
			if C.qbg_set_property_object_type_float16(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetObjectType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("objectType", t, err)
			}
		case Float:
			if C.qbg_set_property_object_type_float(n.prop, ebuf) == ErrorCode {
				err := errors.ErrFailedToSetObjectType(n.newGoError(ebuf), t.String())
				return errors.NewErrCriticalOption("objectType", t, err)
			}
		default:
			n.PutErrorBuffer(ebuf)
			err := errors.ErrUnsupportedObjectType
			return errors.NewErrCriticalOption("objectType", t, err)
		}
		n.PutErrorBuffer(ebuf)
		n.objectType = t
		return nil
	}
}

// WithCreationEdgeSize represents the option to set the creation edge size for QBG.
func WithCreationEdgeSize(size int) Option {
	return func(n *qbg) error {
		ebuf := n.GetErrorBuffer()
		if C.qbg_set_property_edge_size_for_creation(n.prop, C.int16_t(size), ebuf) == ErrorCode {
			err := errors.ErrFailedToSetCreationEdgeSize(n.newGoError(ebuf))
			return errors.NewErrCriticalOption("creationEdgeSize", size, err)
		}
		n.PutErrorBuffer(ebuf)
		return nil
	}
}

// WithSearchEdgeSize represents the option to set the search edge size for QBG.
func WithSearchEdgeSize(size int) Option {
	return func(n *qbg) error {
		ebuf := n.GetErrorBuffer()
		if C.qbg_set_property_edge_size_for_search(n.prop, C.int16_t(size), ebuf) == ErrorCode {
			err := errors.ErrFailedToSetSearchEdgeSize(n.newGoError(ebuf))
			return errors.NewErrCriticalOption("searchEdgeSize", size, err)
		}
		n.PutErrorBuffer(ebuf)
		return nil
	}
}

// WithDefaultPoolSize represents the option to set the default pool size for QBG.
func WithDefaultPoolSize(poolSize uint32) Option {
	return func(n *qbg) error {
		if poolSize == 0 {
			return errors.NewErrInvalidOption("defaultPoolSize", poolSize)
		}
		n.poolSize = poolSize
		return nil
	}
}

// WithDefaultRadius represents the option to set the default radius for QBG.
func WithDefaultRadius(radius float32) Option {
	return func(n *qbg) error {
		if radius == 0 {
			return errors.NewErrInvalidOption("defaultRadius", radius)
		}
		n.radius = radius
		return nil
	}
}

// WithDefaultEpsilon represents the option to set the default epsilon for QBG.
func WithDefaultEpsilon(epsilon float32) Option {
	return func(n *qbg) error {
		if epsilon == 0 {
			return errors.NewErrInvalidOption("defaultEpsilon", epsilon)
		}
		n.epsilon = epsilon
		return nil
	}
}
