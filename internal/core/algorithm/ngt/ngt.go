//
// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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
#include <stdlib.h>
*/
import "C"

import (
	"os"
	"reflect"
	"sync"
	"unsafe"

	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
)

type (

	// NGT is core interface.
	NGT interface {
		// Search returns search result as []SearchResult
		Search(vec []float32, size int, epsilon, radius float32) ([]SearchResult, error)

		// Linear Search returns linear search result as []SearchResult
		LinearSearch(vec []float32, size int) ([]SearchResult, error)

		// Insert returns NGT object id.
		// This only stores not indexing, you must call CreateIndex and SaveIndex.
		Insert(vec []float32) (uint, error)

		// InsertCommit returns NGT object id.
		// This stores and indexes at the same time.
		InsertCommit(vec []float32, poolSize uint32) (uint, error)

		// BulkInsert returns NGT object ids.
		// This only stores not indexing, you must call CreateIndex and SaveIndex.
		BulkInsert(vecs [][]float32) ([]uint, []error)

		// BulkInsertCommit returns NGT object ids.
		// This stores and indexes at the same time.
		BulkInsertCommit(vecs [][]float32, poolSize uint32) ([]uint, []error)

		// CreateAndSaveIndex call  CreateIndex and SaveIndex in a row.
		CreateAndSaveIndex(poolSize uint32) error

		// CreateIndex creates NGT index.
		CreateIndex(poolSize uint32) error

		// SaveIndex stores NGT index to storage.
		SaveIndex() error

		// SaveIndexWithPath stores NGT index to specified storage.
		SaveIndexWithPath(path string) error

		// Remove removes from NGT index.
		Remove(id uint) error

		// BulkRemove removes multiple NGT index
		BulkRemove(ids ...uint) error

		// GetVector returns vector stored in NGT index.
		GetVector(id uint) ([]float32, error)

		// Close Without save index.
		CloseWithoutSaveIndex()

		// Close NGT index.
		Close()
	}

	ngt struct {
		idxPath             string
		inMemory            bool
		bulkInsertChunkSize int
		dimension           C.int32_t
		objectType          objectType
		radius              float32
		epsilon             float32
		poolSize            uint32
		prop                C.NGTProperty
		epool               sync.Pool
		index               C.NGTIndex
		ospace              C.NGTObjectSpace
		mu                  *sync.RWMutex
	}
)

// ObjectType is alias of object type in NGT.
type objectType int

// DistanceType is alias of distance type in NGT.
type distanceType int

const (
	// -------------------------------------------------------------
	// Object Type Definition
	// -------------------------------------------------------------
	// ObjectNone is unknown object type.
	ObjectNone objectType = iota
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
	// Angle is angle distance.
	Angle
	// Hamming is hamming distance.
	Hamming
	// Cosine is cosine distance.
	Cosine
	// Poincare is poincare distance.
	Poincare
	// Lorentz is lorenz distance.
	Lorentz
	// Jaccard is jaccard distance.
	Jaccard
	// SparseJaccard is sparse jaccard distance.
	SparseJaccard
	// NormalizedL2 is l2 distance with normalization.
	NormalizedL2
	// NormalizedAngle is angle distance with normalization.
	NormalizedAngle
	// NormalizedCosine is cosine distance with normalization.
	NormalizedCosine

	// -------------------------------------------------------------.

	// -------------------------------------------------------------
	// ErrorCode is false
	// -------------------------------------------------------------.
	ErrorCode = C._Bool(false)
	// -------------------------------------------------------------.
)

func (o objectType) String() string {
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

// New returns NGT instance with recreating empty index file.
func New(opts ...Option) (NGT, error) {
	return gen(false, opts...)
}

// Load returns NGT instance from existing index file.
func Load(opts ...Option) (NGT, error) {
	return gen(true, opts...)
}

func gen(isLoad bool, opts ...Option) (NGT, error) {
	var (
		n   = new(ngt)
		err error
	)
	n.mu = new(sync.RWMutex)

	defer func() {
		if err != nil {
			n.Close()
		}
	}()

	err = n.setup()
	if err != nil {
		log.Warnf("failed to setup ngt core index\terr: %v", err)
		return nil, err
	}

	defer C.ngt_destroy_property(n.prop)

	err = n.loadOptions(opts...)
	if err != nil {
		log.Warnf("failed to load ngt core options\terr: %v", err)
		return nil, err
	}

	if isLoad {
		err = n.open()
		if err != nil {
			log.Warnf("failed to load ngt core index\terr: %v", err)
			return nil, err
		}
	} else {
		err = n.create()
		if err != nil {
			log.Warnf("failed to create new ngt core index\terr: %v", err)
			return nil, err
		}
	}

	err = n.loadObjectSpace()
	if err != nil {
		log.Warnf("failed to load ngt object space\terr: %v", err)
		return nil, err
	}

	return n, nil
}

func (n *ngt) setup() error {
	n.epool = sync.Pool{
		New: func() interface{} {
			return C.ngt_create_error_object()
		},
	}

	for i := 0; i < 20; i++ {
		n.PutErrorBuffer(C.ngt_create_error_object())
	}

	ebuf := n.GetErrorBuffer()
	n.prop = C.ngt_create_property(ebuf)
	if n.prop == nil {
		return errors.ErrCreateProperty(n.newGoError(ebuf))
	}
	n.PutErrorBuffer(ebuf)
	return nil
}

func (n *ngt) loadOptions(opts ...Option) (err error) {
	for _, opt := range append(defaultOptions, opts...) {
		err = opt(n)
		if err != nil {
			werr := errors.ErrOptionFailed(err, reflect.ValueOf(opt))
			e := new(errors.ErrCriticalOption)
			if errors.As(err, &e) {
				log.Error(werr)
				return werr
			}
			ue := new(errors.ErrIgnoredOption)
			if errors.As(err, &ue) {
				log.Debug(werr)
			} else {
				log.Warn(werr)
			}
		}
	}
	return nil
}

func (n *ngt) create() (err error) {
	files, err := file.ListInDir(n.idxPath)
	if err == nil && len(files) != 0 {
		log.Warnf("index path exists, will remove the directories: %v", files)
		for _, f := range files {
			err = os.RemoveAll(f)
			if err != nil {
				return err
			}
		}
	} else if err != nil {
		log.Debug(err)
	}
	path := C.CString(n.idxPath)
	defer C.free(unsafe.Pointer(path))

	ebuf := n.GetErrorBuffer()
	if !n.inMemory {
		n.index = C.ngt_create_graph_and_tree(path, n.prop, ebuf)
		if n.index == nil {
			return n.newGoError(ebuf)
		}
		if C.ngt_save_index(n.index, path, ebuf) == ErrorCode {
			return n.newGoError(ebuf)
		}
	} else {
		n.index = C.ngt_create_graph_and_tree_in_memory(n.prop, ebuf)
		if n.index == nil {
			return n.newGoError(ebuf)
		}
	}
	n.PutErrorBuffer(ebuf)

	return nil
}

func (n *ngt) open() error {
	if !file.Exists(n.idxPath) {
		return errors.ErrIndexFileNotFound
	}

	path := C.CString(n.idxPath)
	defer C.free(unsafe.Pointer(path))

	ebuf := n.GetErrorBuffer()
	n.index = C.ngt_open_index(path, ebuf)
	if n.index == nil {
		return n.newGoError(ebuf)
	}

	if C.ngt_get_property(n.index, n.prop, ebuf) == ErrorCode {
		return n.newGoError(ebuf)
	}

	n.dimension = C.ngt_get_property_dimension(n.prop, ebuf)
	if int(n.dimension) == -1 {
		return n.newGoError(ebuf)
	}
	n.PutErrorBuffer(ebuf)
	return nil
}

func (n *ngt) loadObjectSpace() error {
	ebuf := n.GetErrorBuffer()
	n.ospace = C.ngt_get_object_space(n.index, ebuf)
	if n.ospace == nil {
		return n.newGoError(ebuf)
	}
	n.PutErrorBuffer(ebuf)
	return nil
}

// Search returns search result as []SearchResult.
func (n *ngt) Search(vec []float32, size int, epsilon, radius float32) (result []SearchResult, err error) {
	if len(vec) != int(n.dimension) {
		return nil, errors.ErrIncompatibleDimensionSize(len(vec), int(n.dimension))
	}

	ebuf := n.GetErrorBuffer()
	results := C.ngt_create_empty_results(ebuf)
	defer C.ngt_destroy_results(results)
	if results == nil {
		return nil, n.newGoError(ebuf)
	}

	if epsilon == 0 {
		epsilon = n.epsilon
	}

	if radius == 0 {
		radius = n.radius
	}

	n.mu.RLock()
	ret := C.ngt_search_index_as_float(
		n.index,
		(*C.float)(&vec[0]),
		n.dimension,
		*(*C.size_t)(unsafe.Pointer(&size)),
		*(*C.float)(unsafe.Pointer(&epsilon)),
		*(*C.float)(unsafe.Pointer(&radius)),
		results,
		ebuf)

	if ret == ErrorCode {
		ne := ebuf
		n.mu.RUnlock()
		return nil, n.newGoError(ne)
	}
	n.mu.RUnlock()

	rsize := int(C.ngt_get_result_size(results, ebuf))
	if rsize <= 0 {
		err = n.newGoError(ebuf)
		if err == nil {
			err = errors.ErrEmptySearchResult
		}
		return nil, err
	}
	result = make([]SearchResult, rsize)

	for i := range result {
		d := C.ngt_get_result(results, C.uint32_t(i), ebuf)
		if d.id == 0 && d.distance == 0 {
			result[i] = SearchResult{0, 0, n.newGoError(ebuf)}
			ebuf = n.GetErrorBuffer()
		} else {
			result[i] = SearchResult{uint32(d.id), float32(d.distance), nil}
		}
	}
	n.PutErrorBuffer(ebuf)

	return result, nil
}

// Linear Search returns linear search result as []SearchResult.
func (n *ngt) LinearSearch(vec []float32, size int) (result []SearchResult, err error) {
	if len(vec) != int(n.dimension) {
		return nil, errors.ErrIncompatibleDimensionSize(len(vec), int(n.dimension))
	}

	ebuf := n.GetErrorBuffer()
	results := C.ngt_create_empty_results(ebuf)
	// defer C.free(unsafe.Pointer(results))
	defer C.ngt_destroy_results(results)
	if results == nil {
		return nil, n.newGoError(ebuf)
	}

	n.mu.RLock()
	ret := C.ngt_linear_search_index_as_float(
		n.index,
		(*C.float)(&vec[0]),
		n.dimension,
		// C.size_t(size),
		*(*C.size_t)(unsafe.Pointer(&size)),
		results,
		ebuf)

	if ret == ErrorCode {
		ne := ebuf
		n.mu.RUnlock()
		return nil, n.newGoError(ne)
	}
	n.mu.RUnlock()

	rsize := int(C.ngt_get_result_size(results, ebuf))
	if rsize <= 0 {
		err = n.newGoError(ebuf)
		if err == nil {
			err = errors.ErrEmptySearchResult
		}
		return nil, err
	}
	result = make([]SearchResult, rsize)

	for i := range result {
		d := C.ngt_get_result(results, C.uint32_t(i), ebuf)
		if d.id == 0 && d.distance == 0 {
			result[i] = SearchResult{0, 0, n.newGoError(ebuf)}
			ebuf = n.GetErrorBuffer()
		} else {
			result[i] = SearchResult{uint32(d.id), float32(d.distance), nil}
		}
	}
	n.PutErrorBuffer(ebuf)

	return result, nil
}

// Insert returns NGT object id.
// This only stores not indexing, you must call CreateIndex and SaveIndex.
func (n *ngt) Insert(vec []float32) (uint, error) {
	dim := int(n.dimension)
	if len(vec) != dim {
		return 0, errors.ErrIncompatibleDimensionSize(len(vec), dim)
	}

	ebuf := n.GetErrorBuffer()
	n.mu.Lock()
	id := C.ngt_insert_index_as_float(n.index, (*C.float)(&vec[0]), C.uint32_t(n.dimension), ebuf)
	n.mu.Unlock()
	if id == 0 {
		return 0, n.newGoError(ebuf)
	}
	n.PutErrorBuffer(ebuf)

	return uint(id), nil
}

// InsertCommit returns NGT object id.
// This stores and indexes at the same time.
func (n *ngt) InsertCommit(vec []float32, poolSize uint32) (uint, error) {
	id, err := n.Insert(vec)
	if err != nil {
		return id, err
	}

	err = n.CreateIndex(poolSize)
	if err != nil {
		return id, err
	}

	err = n.SaveIndex()
	if err != nil {
		return id, err
	}

	return id, nil
}

// BulkInsert returns NGT object ids.
// This only stores not indexing, you must call CreateIndex and SaveIndex.
func (n *ngt) BulkInsert(vecs [][]float32) ([]uint, []error) {
	ids := make([]uint, 0, len(vecs))
	errs := make([]error, 0, len(vecs))

	log.Infof("started to bulk insert %d of vectors", len(vecs))
	for i, vec := range vecs {
		id, err := n.Insert(vec)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "bulkinsert error detected index number: %d,\tid: %d", i, id))
		} else {
			ids = append(ids, id)
		}
	}

	return ids, errs
}

// BulkInsertCommit returns NGT object ids.
// This stores and indexes at the same time.
func (n *ngt) BulkInsertCommit(vecs [][]float32, poolSize uint32) ([]uint, []error) {
	ids := make([]uint, 0, len(vecs))
	errs := make([]error, 0, len(vecs))

	idx := 0
	var id uint
	var err error

	for i, vec := range vecs {
		if id, err = n.Insert(vec); err == nil {
			ids = append(ids, id)
			idx++
			if idx >= n.bulkInsertChunkSize {
				err = n.CreateAndSaveIndex(poolSize)
				if err != nil {
					errs = append(errs, err)
				}
				idx = 0
			}
		} else {
			errs = append(errs, errors.Wrapf(err, "bulkinsert error detected index number: %d,\tid: %d", i, id))
		}
	}

	if idx > 0 {
		err = n.CreateAndSaveIndex(poolSize)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return ids, errs
}

// CreateAndSaveIndex call  CreateIndex and SaveIndex in a row.
func (n *ngt) CreateAndSaveIndex(poolSize uint32) error {
	err := n.CreateIndex(poolSize)
	if err != nil {
		return err
	}
	return n.SaveIndex()
}

// CreateIndex creates NGT index.
func (n *ngt) CreateIndex(poolSize uint32) error {
	if poolSize == 0 {
		poolSize = n.poolSize
	}
	ebuf := n.GetErrorBuffer()
	n.mu.Lock()
	ret := C.ngt_create_index(n.index, C.uint32_t(poolSize), ebuf)
	n.mu.Unlock()
	if ret == ErrorCode {
		return n.newGoError(ebuf)
	}
	n.PutErrorBuffer(ebuf)

	return nil
}

// SaveIndex stores NGT index to storage.
func (n *ngt) SaveIndex() error {
	if !n.inMemory {
		path := C.CString(n.idxPath)
		defer C.free(unsafe.Pointer(path))
		ebuf := n.GetErrorBuffer()
		n.mu.Lock()
		ret := C.ngt_save_index(n.index, path, ebuf)
		n.mu.Unlock()
		if ret == ErrorCode {
			return n.newGoError(ebuf)
		}
		n.PutErrorBuffer(ebuf)
	}

	return nil
}

// SaveIndexWithPath stores NGT index to specified storage.
func (n *ngt) SaveIndexWithPath(idxPath string) error {
	if !n.inMemory && len(idxPath) != 0 {
		path := C.CString(idxPath)
		defer C.free(unsafe.Pointer(path))
		ebuf := n.GetErrorBuffer()
		n.mu.Lock()
		ret := C.ngt_save_index(n.index, path, ebuf)
		n.mu.Unlock()
		if ret == ErrorCode {
			return n.newGoError(ebuf)
		}
		n.PutErrorBuffer(ebuf)
	}

	return nil
}

// Remove removes from NGT index.
func (n *ngt) Remove(id uint) error {
	ebuf := n.GetErrorBuffer()
	n.mu.Lock()
	ret := C.ngt_remove_index(n.index, C.ObjectID(id), ebuf)
	n.mu.Unlock()
	if ret == ErrorCode {
		return n.newGoError(ebuf)
	}
	n.PutErrorBuffer(ebuf)

	return nil
}

// BulkRemove removes multiple index from NGT index.
func (n *ngt) BulkRemove(ids ...uint) (errs error) {
	for i, id := range ids {
		err := n.Remove(id)
		if err != nil {
			errs = errors.Wrapf(errs, "bulkremove error detected index number: %d,\tid: %d\terr: %v", i, id, err)
		}
	}
	return errs
}

// GetVector returns vector stored in NGT index.
func (n *ngt) GetVector(id uint) ([]float32, error) {
	dimension := int(n.dimension)
	var ret []float32
	ebuf := n.GetErrorBuffer()
	switch n.objectType {
	case Float, HalfFloat:
		n.mu.RLock()
		results := C.ngt_get_object_as_float(n.ospace, C.ObjectID(id), ebuf)
		n.mu.RUnlock()
		if results == nil {
			return nil, n.newGoError(ebuf)
		}
		ret = (*[algorithm.MaximumVectorDimensionSize]float32)(unsafe.Pointer(results))[:dimension:dimension]
		// for _, elem := range (*[algorithm.MaximumVectorDimensionSize]C.float)(unsafe.Pointer(results))[:dimension:dimension]{
		// 	ret = append(ret, float32(elem))
		// }
	case Uint8:
		n.mu.RLock()
		results := C.ngt_get_object_as_integer(n.ospace, C.ObjectID(id), ebuf)
		n.mu.RUnlock()
		if results == nil {
			return nil, n.newGoError(ebuf)
		}
		ret = make([]float32, 0, dimension)
		for _, elem := range (*[algorithm.MaximumVectorDimensionSize]C.uint8_t)(unsafe.Pointer(results))[:dimension:dimension] {
			ret = append(ret, float32(elem))
		}
	default:
		n.PutErrorBuffer(ebuf)
		return nil, errors.ErrUnsupportedObjectType
	}
	n.PutErrorBuffer(ebuf)
	return ret, nil
}

func (n *ngt) newGoError(ebuf C.NGTError) (err error) {
	msg := C.GoString(C.ngt_get_error_string(ebuf))
	if len(msg) == 0 {
		n.PutErrorBuffer(ebuf)
		return nil
	}
	n.PutErrorBuffer(C.ngt_create_error_object())
	C.ngt_destroy_error_object(ebuf)
	return errors.NewNGTError(msg)
}

// Close NGT without save index.
func (n *ngt) CloseWithoutSaveIndex() {
	n.index = nil
	n.prop = nil
	n.ospace = nil
}

// Close NGT index.
func (n *ngt) Close() {
	if n.index != nil {
		C.ngt_close_index(n.index)
		n.index = nil
		n.prop = nil
		n.ospace = nil
	}
}

func (n *ngt) GetErrorBuffer() (ebuf C.NGTError) {
	var ok bool
	ebuf, ok = n.epool.Get().(C.NGTError)
	if !ok {
		ebuf = C.ngt_create_error_object()
	}
	return ebuf
}

func (n *ngt) PutErrorBuffer(ebuf C.NGTError) {
	n.epool.Put(ebuf)
}
