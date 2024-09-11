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

// Package usearch provides implementation of Go API for https://github.com/unum-cloud/usearch
package usearch

import (
	"strconv"
	"strings"

	"github.com/kpango/fastime"
	core "github.com/unum-cloud/usearch/golang"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
)

// Option represents the functional option for usearch.
type Option func(*usearch) error

var defaultOptions = []Option{
	WithIndexPath("/tmp/usearch-" + strconv.FormatInt(fastime.UnixNanoNow(), 10)),
	WithQuantizationType("F32"),
	WithMetricType("cosine"),
	WithDimension(64),
	WithConnectivity(0),
	WithExpansionAdd(0),
	WithExpansionSearch(0),
	WithMulti(false),
}

// WithIndexPath represents the option to set the index path for usearch.
func WithIndexPath(path string) Option {
	return func(u *usearch) error {
		if len(path) == 0 {
			return errors.NewErrIgnoredOption("indexPath")
		}
		u.idxPath = path
		return nil
	}
}

// WithQuantizationType represents the option to set the quantizationType for usearch.
func WithQuantizationType(quantizationType string) Option {
	return func(u *usearch) error {
		switch quantizationType{
		case "BF16":
			u.quantizationType = core.BF16
		case "F16":
			u.quantizationType = core.F16
		case "F32":
			u.quantizationType = core.F32
		case "F64":
			u.quantizationType = core.F64
		case "I8":
			u.quantizationType = core.I8
		case "B1":
			u.quantizationType = core.B1
		default:
			err := errors.NewUsearchError("unsupported QuantizationType")
			return errors.NewErrCriticalOption("QuantizationType", quantizationType, err)
		}
		return nil
	}
}

// WithMetricType represents the option to set the metricType for usearch.
func WithMetricType(metricType string) Option {
	return func(u *usearch) error {
		switch strings.NewReplacer("-", "", "_", "", " ", "").Replace(strings.ToLower(metricType)) {
		case "l2sq":
			u.metricType = core.L2sq
		case "ip":
			u.metricType = core.InnerProduct
		case "cosine":
			u.metricType = core.Cosine
		case "haversine":
			u.metricType = core.Haversine
		case "divergence":
			u.metricType = core.Divergence
		case "pearson":
			u.metricType = core.Pearson
		case "hamming":
			u.metricType = core.Hamming
		case "tanimoto":
			u.metricType = core.Tanimoto
		case "sorensen":
			u.metricType = core.Sorensen
		default:
			err := errors.NewUsearchError("unsupported MetricType")
			return errors.NewErrCriticalOption("MetricType", metricType, err)
		}
		return nil
	}
}

// WithDimension represents the option to set the dimension for usearch.
func WithDimension(dim int) Option {
	return func(u *usearch) error {
		if dim > algorithm.MaximumVectorDimensionSize || dim < algorithm.MinimumVectorDimensionSize {
			err := errors.ErrInvalidDimensionSize(dim, algorithm.MaximumVectorDimensionSize)
			return errors.NewErrCriticalOption("dimension", dim, err)
		}

		u.dimension = uint(dim)
		return nil
	}
}

// WithConnectivity represents the option to set the connectivity for usearch.
func WithConnectivity(connectivity int) Option {
	return func(u *usearch) error {
		if connectivity < 0 {
			return errors.NewErrInvalidOption("Connectivity", connectivity)
		}

		u.connectivity = uint(connectivity)
		return nil
	}
}

// WithExpansionAdd represents the option to set the expansion add for usearch.
func WithExpansionAdd(expansionAdd int) Option {
	return func(u *usearch) error {
		if expansionAdd < 0 {
			return errors.NewErrInvalidOption("Expansion Add", expansionAdd)
		}

		u.expansionAdd = uint(expansionAdd)
		return nil
	}
}

// WithExpansionSearch represents the option to set the expansion search for usearch.
func WithExpansionSearch(expansionSearch int) Option {
	return func(u *usearch) error {
		if expansionSearch < 0 {
			return errors.NewErrInvalidOption("Expansion Search", expansionSearch)
		}

		u.expansionSearch = uint(expansionSearch)
		return nil
	}
}

// WithMulti represents the option to set the multi for usearch.
func WithMulti(multi bool) Option {
	return func(u *usearch) error {
		u.multi = multi
		return nil
	}
}
