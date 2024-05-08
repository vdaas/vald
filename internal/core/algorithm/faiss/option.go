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

// Package faiss provides implementation of Go API for https://github.com/facebookresearch/faiss
package faiss

import "C"

import (
	"strings"

	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
)

// Option represents the functional option for faiss.
type Option func(*faiss) error

var defaultOptions = []Option{
	WithDimension(64),
	WithNlist(100),
	WithM(8),
	WithNbitsPerIdx(8),
	WithMetricType("l2"),
}

// WithDimension represents the option to set the dimension for faiss.
func WithDimension(dim int) Option {
	return func(f *faiss) error {
		if dim > algorithm.MaximumVectorDimensionSize || dim < algorithm.MinimumVectorDimensionSize {
			err := errors.ErrInvalidDimensionSize(dim, algorithm.MaximumVectorDimensionSize)
			return errors.NewErrCriticalOption("dimension", dim, err)
		}

		f.dimension = (C.int)(dim)
		return nil
	}
}

// WithNlist represents the option to set the nlist for faiss.
func WithNlist(nlist int) Option {
	return func(f *faiss) error {
		if nlist <= 0 {
			return errors.NewErrInvalidOption("nlist", nlist)
		}

		f.nlist = (C.int)(nlist)
		return nil
	}
}

// WithM represents the option to set the m for faiss.
func WithM(m int) Option {
	return func(f *faiss) error {
		if m <= 0 || int(f.dimension)%m != 0 {
			return errors.NewErrInvalidOption("m", m)
		}

		f.m = (C.int)(m)
		return nil
	}
}

// WithNbitsPerIdx represents the option to set the n bits per index for faiss.
func WithNbitsPerIdx(nbitsPerIdx int) Option {
	return func(f *faiss) error {
		if nbitsPerIdx <= 0 {
			return errors.NewErrInvalidOption("nbitsPerIdx", nbitsPerIdx)
		}

		f.nbitsPerIdx = (C.int)(nbitsPerIdx)
		return nil
	}
}

// WithMetricType represents the option to set the metric type for faiss.
func WithMetricType(metricType string) Option {
	return func(f *faiss) error {
		if len(metricType) == 0 {
			return errors.NewErrIgnoredOption("metricType")
		}

		switch strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToLower(metricType)) {
		case "innerproduct":
			f.metricType = InnerProduct
		case "l2":
			f.metricType = L2
		default:
			err := errors.ErrUnsupportedDistanceType
			return errors.NewErrCriticalOption("metricType", metricType, err)
		}

		return nil
	}
}

// WithIndexPath represents the option to set the index path for faiss.
func WithIndexPath(idxPath string) Option {
	return func(f *faiss) error {
		if len(idxPath) == 0 {
			return errors.NewErrIgnoredOption("indexPath")
		}

		f.idxPath = idxPath
		return nil
	}
}
