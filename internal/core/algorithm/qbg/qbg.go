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
#include <NGT/NGTQ/Capi.h>
#include <NGT/NGTQ/Quantizer.h>
#include <stdlib.h>
*/
import "C"

import (
	"sync"

	"github.com/vdaas/vald/internal/errors"
)

type (

	// QBG is core interface.
	QBG interface {
		// Search returns search result as []SearchResult
		Search(vec []float32, size int, epsilon, radius float32) ([]SearchResult, error)

		// Insert returns QBG object id.
		// This only stores not indexing, you must call CreateIndex and SaveIndex.
		Insert(vec []float32) (uint, error)

		// BulkInsert returns QBG object ids.
		// This only stores not indexing, you must call CreateIndex and SaveIndex.
		BulkInsert(vecs [][]float32) ([]uint, []error)

		// SaveIndex stores QBG index to storage.
		SaveIndex() error

		// GetVector returns vector stored in QBG index.
		GetVector(id uint) ([]float32, error)

		// Close QBG index.
		Close()
	}

	qbg struct {
		idxPath string
		params  *constructionParams
		epool   sync.Pool
		index   C.QBGIndex
		ospace  C.NGTObjectSpace
		mu      *sync.RWMutex
	}

	constructionParams struct {
		dimension          C.size_t
		extendedDimension  C.size_t
		numberOfSubVectors C.size_t
		numberOfBlobs      C.size_t
		dataType           dataType
		internalDataType   dataType
		distanceType       distanceType
	}

	buildParams struct {
		// hierarchical kmeans
		hierarchicalClusteringInitialiationMode int
		numberOfFirstClusters                   C.size_t
		numberOfFirstObjects                    C.size_t
		numberOfSecondClusters                  C.size_t
		numberOfSecondObjects                   C.size_t
		numberOfThirdClusters                   C.size_t
		// // optimization
		numberOfMatrices                         C.size_t
		numberOfObjects                          C.size_t
		numberOfSubvectors                       C.size_t
		optimizationClusteringInitializationMode int
		repositioning                            bool
		rotation                                 bool
		rotationIteration                        C.size_t
		subvectorIteration                       C.size_t
	}

	queryParams struct {
		query                    *C.float
		epsilon                  C.float
		blob_epsilon             C.float
		result_expansion         C.float
		number_of_results        C.size_t
		number_of_explored_blobs C.size_t
		number_of_edges          C.size_t
		radius                   C.float
	}

	qbgCore interface {
		// void qbg_initialize_construction_parameters(QBGConstructionParameters *parameters);
		InitializeConstructionParameters(params *C.QBGConstructionParameters)
		// void qbg_initialize_build_parameters(QBGBuildParameters *parameters);
		InitializeBuildParameters(params *C.QBGBuildParameters)
		// void qbg_initialize_query(QBGQuery *parameters);
		InitializeQuery(query *C.QBGQuery)
		// bool qbg_create(const char *indexPath, QBGConstructionParameters *parameters, QBGError error);
		Create(indexPath string, params *C.QBGConstructionParameters, err C.QBGError) bool
		// QBGIndex qbg_open_index(const char *index_path, QBGError error);
		Open(index_path string, err C.QBGError) C.QBGIndex
		// void qbg_close_index(QBGIndex index);
		Close(index C.QBGIndex)
		// bool qbg_save_index(QBGIndex index, QBGError error);
		Save(index C.QBGIndex, err C.QBGError) bool
		// ObjectID qbg_append_object(QBGIndex index, float *obj, uint32_t obj_dim, QBGError error);
		Append(index C.QBGIndex, vec []float32, obj_dim uint32, err C.QBGError) C.ObjectID
		// bool qbg_build_index(const char *index_path, QBGBuildParameters *parameters, QBGError error);
		Build(index_path string, params *C.QBGBuildParameters, err C.QBGError) bool
		// bool qbg_search_index(QBGIndex index, QBGQuery query, NGTObjectDistances results, QBGError error);
		Search(index C.QBGIndex, query C.QBGQuery, results C.NGTObjectDistances, err C.QBGError) bool
		// void qbg_destroy_results(QBGObjectDistances results);
		DestroyResults(results C.QBGObjectDistances)
	}
)

// DataType is alias of data type in QBG.
type dataType int

// DistanceType is alias of distance type in QBG.
type distanceType int

const (
	// -------------------------------------------------------------
	// Data Type Definition
	// -------------------------------------------------------------
	// DataTypeNone is unknown object type.
	DataTypeNone dataType = iota
	// Uint8 is 8bit unsigned integer.
	Uint8
	// Float is 32bit floating point number.
	Float
	// HalfFloat is 16bit floating point number.
	HalfFloat
	// -------------------------------------------------------------.

	// -------------------------------------------------------------
	// Distance Type Definition
	// -------------------------------------------------------------
	// DistanceNone is unknown distance type.
	DistanceNone distanceType = iota - 1
	// L1 is l1 norm.
	L1
	// L2 is l2 norm.
	L2
	// Hamming is hamming distance.
	Hamming
	// Angle is angle distance.
	Angle
	// Cosine is cosine distance.
	Cosine
	// NormalizedAngle is angle distance with normalization.
	NormalizedAngle
	// NormalizedCosine is cosine distance with normalization.
	NormalizedCosine
	// Jaccard is jaccard distance.
	Jaccard
	// SparseJaccard is sparse jaccard distance.
	SparseJaccard
	// NormalizedL2 is l2 distance with normalization.
	NormalizedL2

	// Poincare is poincare distance.
	Poincare distanceType = 100
	// Lorentz is lorenz distance.
	Lorentz distanceType = 101
	// -------------------------------------------------------------.

	// -------------------------------------------------------------
	// ErrorCode is false
	// -------------------------------------------------------------.
	ErrorCode = C._Bool(false)
	// -------------------------------------------------------------.
)

func (o dataType) String() string {
	switch o {
	case Uint8:
		return "Uint8"
	case HalfFloat:
		return "HalfFloat"
	case Float:
		return "Float"
	}
	return "Unknown"
}

func (d distanceType) String() string {
	switch d {
	case L1:
		return "L1"
	case L2:
		return "L2"
	case Angle:
		return "Angle"
	case Hamming:
		return "Hamming"
	case Cosine:
		return "Cosine"
	case Poincare:
		return "Poincare"
	case Lorentz:
		return "Lorentz"
	case Jaccard:
		return "Jaccard"
	case SparseJaccard:
		return "SparseJaccard"
	case NormalizedL2:
		return "NormalizedL2"
	case NormalizedAngle:
		return "NormalizedAngle"
	case NormalizedCosine:
		return "NormalizedCosine"
	}
	return "Unknown"
}

// New returns QBG instance with recreating empty index file.
func New(opts ...Option) (QBG, error) {
	return nil, nil
}

func (n *qbg) newGoError(ebuf C.QBGError) (err error) {
	msg := C.GoString(C.ngt_get_error_string(ebuf))
	if len(msg) == 0 {
		n.PutErrorBuffer(ebuf)
		return nil
	}
	n.PutErrorBuffer(C.ngt_create_error_object())
	C.qbg_destroy_error_object(ebuf)
	return errors.NewQBGError(msg)
}

// Close QBG index.
func (n *qbg) Close() {
	if n.index != nil {
		C.qbg_close_index(n.index)
		n.index = nil
		n.params = nil
		n.ospace = nil
	}
}

func (n *qbg) GetErrorBuffer() (ebuf C.QBGError) {
	var ok bool
	ebuf, ok = n.epool.Get().(C.QBGError)
	if !ok {
		ebuf = C.ngt_create_error_object()
	}
	return ebuf
}

func (n *qbg) PutErrorBuffer(ebuf C.QBGError) {
	n.epool.Put(ebuf)
}
