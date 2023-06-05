//
// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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

// Package faiss provides implementation of Go API for https://github.com/facebookresearch/faiss
package faiss

/*
#cgo LDFLAGS: -lfaiss
#include <Capi.h>
*/
import "C"

import (
	"sync"
	"unsafe"

	"github.com/vdaas/vald/internal/errors"
)

type (
	// Faiss is core interface.
	Faiss interface {
		// SaveIndex stores faiss index to strage.
		SaveIndex() error

		// SaveIndexWithPath stores faiss index to specified storage.
		SaveIndexWithPath(idxPath string) error

		// Train trains faiss index.
		Train(nb int, xb []float32) error

		// Add returns faiss ntotal.
		Add(nb int, xb []float32, xids []int64) (int, error)

		// Search returns search result as []SearchResult.
		Search(k, nq int, xq []float32) ([]SearchResult, error)

		// Remove removes from faiss index.
		Remove(size int, ids []int64) (int, error)

		// Close faiss index.
		Close()
	}

	faiss struct {
		st          *C.FaissStruct
		dimension   C.int
		nlist       C.int
		m           C.int
		nbitsPerIdx C.int
		metricType  metricType
		idxPath     string
		mu          *sync.RWMutex
	}

	SearchResult struct {
		ID       uint32
		Distance float32
		Error    error
	}
)

// metricType is alias of metric type in Faiss.
type metricType int

const (
	// -------------------------------------------------------------
	// Metric Type Definition
	// (https://github.com/facebookresearch/faiss/wiki/MetricType-and-distances)
	// -------------------------------------------------------------
	// DistanceNone is unknown distance type.
	DistanceNone metricType = iota - 1
	// InnerProduct is inner product.
	InnerProduct
	// L2 is l2 norm.
	L2
	// -------------------------------------------------------------.

	// -------------------------------------------------------------
	// ErrorCode is false
	// -------------------------------------------------------------.
	ErrorCode = C._Bool(false)
	// -------------------------------------------------------------.
)

// New returns Faiss instance with recreating empty index file.
func New(opts ...Option) (Faiss, error) {
	return gen(false, opts...)
}

func Load(opts ...Option) (Faiss, error) {
	return gen(true, opts...)
}

func gen(isLoad bool, opts ...Option) (Faiss, error) {
	var (
		f   = new(faiss)
		err error
	)
	f.mu = new(sync.RWMutex)

	defer func() {
		if err != nil {
			f.Close()
		}
	}()

	for _, opt := range append(defaultOptions, opts...) {
		if err = opt(f); err != nil {
			return nil, errors.NewFaissError("faiss option error")
		}
	}

	if isLoad {
		path := C.CString(f.idxPath)
		defer C.free(unsafe.Pointer(path))
		f.st = C.faiss_read_index(path)
		if f.st == nil {
			return nil, errors.NewFaissError("faiss load index error")
		}
	} else {
		switch f.metricType {
		case InnerProduct:
			f.st = C.faiss_create_index(f.dimension, f.nlist, f.m, f.nbitsPerIdx, C.int(InnerProduct))
		case L2:
			f.st = C.faiss_create_index(f.dimension, f.nlist, f.m, f.nbitsPerIdx, C.int(L2))
		default:
			return nil, errors.NewFaissError("faiss create index error: no metric type")
		}
		if f.st == nil {
			return nil, errors.NewFaissError("faiss create index error: nil pointer")
		}
	}

	return f, nil
}

// SaveIndex stores faiss index to storage.
func (f *faiss) SaveIndex() error {
	path := C.CString(f.idxPath)
	defer C.free(unsafe.Pointer(path))

	f.mu.Lock()
	ret := C.faiss_write_index(f.st, path)
	f.mu.Unlock()
	if ret == ErrorCode {
		return errors.NewFaissError("failed to faiss_write_index")
	}

	return nil
}

// SaveIndexWithPath stores faiss index to specified storage.
func (f *faiss) SaveIndexWithPath(idxPath string) error {
	path := C.CString(idxPath)
	defer C.free(unsafe.Pointer(path))

	f.mu.Lock()
	ret := C.faiss_write_index(f.st, path)
	f.mu.Unlock()
	if ret == ErrorCode {
		return errors.NewFaissError("failed to faiss_write_index")
	}

	return nil
}

// Train trains faiss index.
func (f *faiss) Train(nb int, xb []float32) error {
	f.mu.Lock()
	ret := C.faiss_train(f.st, (C.int)(nb), (*C.float)(&xb[0]))
	f.mu.Unlock()
	if ret == ErrorCode {
		return errors.NewFaissError("failed to faiss_train")
	}

	return nil
}

// Add returns faiss ntotal.
func (f *faiss) Add(nb int, xb []float32, xids []int64) (int, error) {
	dim := int(f.dimension)
	if len(xb) != dim*nb || len(xb) != dim*len(xids) {
		return -1, errors.ErrIncompatibleDimensionSize(len(xb)/nb, dim)
	}

	f.mu.Lock()
	ntotal := int(C.faiss_add(f.st, (C.int)(nb), (*C.float)(&xb[0]), (*C.long)(&xids[0])))
	f.mu.Unlock()
	if ntotal < 0 {
		return ntotal, errors.NewFaissError("failed to faiss_add")
	}

	return ntotal, nil
}

// Search returns search result as []SearchResult.
func (f *faiss) Search(k, nq int, xq []float32) ([]SearchResult, error) {
	if len(xq) != nq*int(f.dimension) {
		return nil, errors.ErrIncompatibleDimensionSize(len(xq), int(f.dimension))
	}

	I := make([]int64, k*nq)
	D := make([]float32, k*nq)
	f.mu.RLock()
	ret := C.faiss_search(f.st, (C.int)(k), (C.int)(nq), (*C.float)(&xq[0]), (*C.long)(&I[0]), (*C.float)(&D[0]))
	f.mu.RUnlock()
	if ret == ErrorCode {
		return nil, errors.NewFaissError("failed to faiss_search")
	}

	if len(I) == 0 || len(D) == 0 {
		return nil, errors.ErrEmptySearchResult
	}

	result := make([]SearchResult, k)
	for i := range result {
		result[i] = SearchResult{uint32(I[i]), D[i], nil}
	}

	return result, nil
}

// Remove removes from faiss index.
func (f *faiss) Remove(size int, ids []int64) (int, error) {
	f.mu.Lock()
	ntotal := int(C.faiss_remove(f.st, (C.int)(size), (*C.long)(&ids[0])))
	f.mu.Unlock()
	if ntotal < 0 {
		return ntotal, errors.NewFaissError("failed to faiss_remove")
	}

	return ntotal, nil
}

// Close faiss index.
func (f *faiss) Close() {
	if f.st != nil {
		C.faiss_free(f.st)
		f.st = nil
	}
}
