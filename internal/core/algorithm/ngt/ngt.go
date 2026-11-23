// Copyright (C) 2019-2026 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ngt

/*
#cgo LDFLAGS: -L/usr/local/lib -lngt
#cgo noescape vald_ngt_search_float
#cgo nocallback vald_ngt_search_float
#include <NGT/Capi.h>
#include <stdlib.h>

static inline bool vald_ngt_search_float(
	bool linear,
	NGTIndex index,
	float *vec,
	int32_t dim,
	size_t size,
	float epsilon,
	float radius,
	int edge_size,
	uint32_t *out_ids,
	float *out_dists,
	int out_cap,
	int *out_total,
	uint64_t *out_nobj,
	int *out_bad_idx,
	NGTError err
) {
	*out_total = 0;
	*out_nobj = 0;
	*out_bad_idx = -1;
	NGTObjectDistances results = ngt_create_empty_results(err);
	if (results == NULL) {
		return false;
	}
	bool ok = linear
		? ngt_linear_search_index_as_float(index, vec, dim, size, results, err)
		: ngt_search_index_as_float(index, vec, dim, size, epsilon, radius, results, edge_size, err);
	if (!ok) {
		ngt_destroy_results(results);
		return false;
	}
	int rsize = (int)ngt_get_result_size(results, err);
	*out_total = rsize;
	if (rsize <= 0) {
		*out_nobj = (uint64_t)ngt_get_number_of_objects(index, err);
		ngt_destroy_results(results);
		return true;
	}
	if (rsize > out_cap) {
		ngt_destroy_results(results);
		return true;
	}
	for (int i = 0; i < rsize; i++) {
		NGTObjectDistance d = ngt_get_result(results, (uint32_t)i, err);
		out_ids[i] = (uint32_t)d.id;
		out_dists[i] = d.distance;
		if (d.id == 0 && d.distance == 0 && *out_bad_idx < 0) {
			*out_bad_idx = i;
		}
	}
	ngt_destroy_results(results);
	return true;
}
*/
import "C"

import (
	"context"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/safety"
	"github.com/vdaas/vald/internal/sync"
	"github.com/vdaas/vald/internal/sync/atomic"
	"github.com/vdaas/vald/internal/sync/singleflight"
)

// errIndexClosed is returned instead of passing a NULL handle to an NGT C function
// that does not check for it (and would SIGSEGV) after Close or before the index is opened.
var errIndexClosed = errors.NewNGTError("ngt index is closed")

type (

	//nolint:interfacebloat // core NGT wrapper interface
	// NGT is core interface.
	NGT interface {
		// Search returns search result as []algorithm.SearchResult
		Search(ctx context.Context, vec []float32, size int, epsilon, radius float32, edgeSize int32) ([]algorithm.SearchResult, error)

		// SearchValidated returns search result as []algorithm.SearchResult without re-checking vec for NaN/Inf.
		SearchValidated(ctx context.Context, vec []float32, size int, epsilon, radius float32, edgeSize int32) ([]algorithm.SearchResult, error)

		// Linear Search returns linear search result as []algorithm.SearchResult
		LinearSearch(ctx context.Context, vec []float32, size int) ([]algorithm.SearchResult, error)

		// LinearSearchValidated returns linear search result as []algorithm.SearchResult without re-checking vec for NaN/Inf.
		LinearSearchValidated(ctx context.Context, vec []float32, size int) ([]algorithm.SearchResult, error)

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

		// GetVectorInto writes the vector stored in NGT index into dst when cap(dst) is sufficient.
		GetVectorInto(id uint, dst []float32) ([]float32, error)

		GetGraphStatistics(ctx context.Context, m statisticsType) (stats *GraphStatistics, err error)

		// GetProperty returns NGT Index Property.
		GetProperty() (*Property, error)

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
		cnt                 atomic.Uint64
		prop                C.NGTProperty
		ces                 uint64               // NGT edge size for creation
		epool               sync.Pool[*ngtError] // NGT error buffer pool
		eps                 atomic.Int64         // NGT error buffer pool size (signed: concurrent borrowers may drive it below zero)
		epl                 uint64               // NGT error buffer pool size limit
		index               C.NGTIndex
		ospace              C.NGTObjectSpace
		group               singleflight.Group[*GraphStatistics]
		mu                  *sync.RWMutex
		smu                 *sync.Mutex
	}

	ngtError struct {
		err       C.NGTError
		destroyed atomic.Bool
	}

	// GraphStatistics represents the statistics of the NGT graph index.
	GraphStatistics struct {
		// Valid indicates whether the statistics are valid.
		Valid bool
		// MedianIndegree represents the median of indegrees.
		MedianIndegree int32
		// MedianOutdegree represents the median of outdegrees.
		MedianOutdegree int32
		// MaxNumberOfIndegree represents the maximum number of indegrees.
		MaxNumberOfIndegree uint64
		// MaxNumberOfOutdegree represents the maximum number of outdegrees.
		MaxNumberOfOutdegree uint64
		// MinNumberOfIndegree represents the minimum number of indegrees.
		MinNumberOfIndegree uint64
		// MinNumberOfOutdegree represents the minimum number of outdegrees.
		MinNumberOfOutdegree uint64
		// ModeIndegree represents the mode of indegrees.
		ModeIndegree uint64
		// ModeOutdegree represents the mode of outdegrees.
		ModeOutdegree uint64
		// NodesSkippedFor10Edges represents the number of nodes skipped for 10 edges.
		NodesSkippedFor10Edges uint64
		// NodesSkippedForIndegreeDistance represents the number of nodes skipped for indegree distance.
		NodesSkippedForIndegreeDistance uint64
		// NumberOfEdges represents the total number of edges.
		NumberOfEdges uint64
		// NumberOfIndexedObjects represents the number of indexed objects.
		NumberOfIndexedObjects uint64
		// NumberOfNodes represents the number of nodes.
		NumberOfNodes uint64
		// NumberOfNodesWithoutEdges represents the number of nodes without edges.
		NumberOfNodesWithoutEdges uint64
		// NumberOfNodesWithoutIndegree represents the number of nodes without indegree.
		NumberOfNodesWithoutIndegree uint64
		// NumberOfObjects represents the total number of objects.
		NumberOfObjects uint64
		// NumberOfRemovedObjects represents the number of removed objects.
		NumberOfRemovedObjects uint64
		// SizeOfObjectRepository represents the size of the object repository.
		SizeOfObjectRepository uint64
		// SizeOfRefinementObjectRepository represents the size of the refinement object repository.
		SizeOfRefinementObjectRepository uint64
		// VarianceOfIndegree represents the variance of indegrees.
		VarianceOfIndegree float64
		// VarianceOfOutdegree represents the variance of outdegrees.
		VarianceOfOutdegree float64
		// MeanEdgeLength represents the mean edge length.
		MeanEdgeLength float64
		// MeanEdgeLengthFor10Edges represents the mean edge length for 10 edges.
		MeanEdgeLengthFor10Edges float64
		// MeanIndegreeDistanceFor10Edges represents the mean indegree distance for 10 edges.
		MeanIndegreeDistanceFor10Edges float64
		// MeanNumberOfEdgesPerNode represents the mean number of edges per node.
		MeanNumberOfEdgesPerNode float64
		// C1Indegree represents the C1 indegree.
		C1Indegree float64
		// C5Indegree represents the C5 indegree.
		C5Indegree float64
		// C95Outdegree represents the C95 outdegree.
		C95Outdegree float64
		// C99Outdegree represents the C99 outdegree.
		C99Outdegree float64
		// IndegreeCount represents the indegree count distribution.
		IndegreeCount []int64
		// OutdegreeHistogram represents the outdegree histogram.
		OutdegreeHistogram []uint64
		// IndegreeHistogram represents the indegree histogram.
		IndegreeHistogram []uint64
	}

	// Property represents the NGT index property.
	Property struct {
		// Dimension represents the dimension of the vector.
		Dimension int32
		// ThreadPoolSize represents the size of the thread pool.
		ThreadPoolSize int32
		// ObjectType represents the type of the object.
		ObjectType objectType
		// DistanceType represents the type of the distance.
		DistanceType distanceType
		// IndexType represents the type of the index.
		IndexType indexType
		// DatabaseType represents the type of the database.
		DatabaseType databaseType
		// ObjectAlignment represents the object alignment.
		ObjectAlignment objectAlignment
		// PathAdjustmentInterval represents the path adjustment interval.
		PathAdjustmentInterval int32
		// GraphSharedMemorySize represents the size of the graph shared memory.
		GraphSharedMemorySize int32
		// TreeSharedMemorySize represents the size of the tree shared memory.
		TreeSharedMemorySize int32
		// ObjectSharedMemorySize represents the size of the object shared memory.
		ObjectSharedMemorySize int32
		// PrefetchOffset represents the prefetch offset.
		PrefetchOffset int32
		// PrefetchSize represents the prefetch size.
		PrefetchSize int32
		// AccuracyTable represents the accuracy table.
		AccuracyTable string
		// SearchType represents the search type.
		SearchType string
		// MaxMagnitude represents the maximum magnitude.
		MaxMagnitude float32
		// NOfNeighborsForInsertionOrder represents the number of neighbors for insertion order.
		NOfNeighborsForInsertionOrder int32
		// EpsilonForInsertionOrder represents the epsilon for insertion order.
		EpsilonForInsertionOrder float32
		// EpsilonForCreation represents the epsilon for creation.
		EpsilonForCreation float32
		// RefinementObjectType represents the refinement object type.
		RefinementObjectType objectType
		// TruncationThreshold represents the truncation threshold.
		TruncationThreshold int32
		// EdgeSizeForCreation represents the edge size for creation.
		EdgeSizeForCreation int32
		// EdgeSizeForSearch represents the edge size for search.
		EdgeSizeForSearch int32
		// EdgeSizeLimitForCreation represents the upper bound of the edge size for
		// creation. NGT >= 2.8.0 replaced edge_size_limit_for_creation with
		// min/max_edge_size_for_creation; this reports max_edge_size_for_creation
		// to keep the Info.Index.Property API unchanged.
		EdgeSizeLimitForCreation int32
		// InsertionRadiusCoefficient represents the insertion radius coefficient.
		InsertionRadiusCoefficient float64
		// SeedSize represents the seed size.
		SeedSize int32
		// SeedType represents the seed type.
		SeedType seedType
		// TruncationThreadPoolSize represents the truncation thread pool size.
		TruncationThreadPoolSize int32
		// BatchSizeForCreation represents the batch size for creation.
		BatchSizeForCreation int32
		// GraphType represents the graph type.
		GraphType graphType
		// DynamicEdgeSizeBase represents the dynamic edge size base.
		DynamicEdgeSizeBase int32
		// DynamicEdgeSizeRate represents the dynamic edge size rate.
		DynamicEdgeSizeRate int32
		// BuildTimeLimit represents the build time limit.
		BuildTimeLimit float32
		// OutgoingEdge represents the outgoing edge.
		OutgoingEdge int32
		// IncomingEdge represents the incoming edge.
		IncomingEdge int32
	}
)

func newNGTError() (n *ngtError) {
	n = &ngtError{
		err: C.ngt_create_error_object(),
	}
	n.destroyed.Store(false)
	runtime.SetFinalizer(n, func(ne *ngtError) {
		ne.close()
	})
	return n
}

func (n *ngtError) close() {
	if n.destroyed.CompareAndSwap(false, true) {
		C.ngt_destroy_error_object(n.err)
	}
}

// ObjectType is alias of object type in NGT.
type objectType int

// DistanceType is alias of distance type in NGT.
type distanceType int

type statisticsType int

// IndexType is alias of index type in NGT.
type indexType int

// DatabaseType is alias of database type in NGT.
type databaseType int

// ObjectAlignment is alias of object alignment in NGT.
type objectAlignment int

// SeedType is alias of seed type in NGT.
type seedType int

// GraphType is alias of graph type in NGT.
type graphType int

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
)

const (
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
	// Lorentz is lorentz distance.
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
	// InnerProduct is inner product distance.
	InnerProduct
)

const (
	// -------------------------------------------------------------
	// IndexType Definition
	// -------------------------------------------------------------
	IndexTypeNone indexType = iota
	GraphAndTree
	Graph
)

const (
	// -------------------------------------------------------------
	// DatabaseType Definition
	// -------------------------------------------------------------
	DatabaseTypeNone databaseType = iota
	Memory
	MemoryMappedFile
)

const (
	// -------------------------------------------------------------
	// ObjectAlignment Definition
	// -------------------------------------------------------------
	ObjectAlignmentNone objectAlignment = iota
	ObjectAlignmentTrue
	ObjectAlignmentFalse
)

const (
	// -------------------------------------------------------------
	// SeedType Definition
	// -------------------------------------------------------------
	// SeedTypeNone is unknown seed type.
	SeedTypeNone seedType = iota
	RandomNodes
	FixedNodes
	FirstNode
	AllLeafNodes
)

const (
	// -------------------------------------------------------------
	// GraphType Definition
	// -------------------------------------------------------------
	GraphTypeNone graphType = iota
	ANNG
	KNNG
	BKNNG
	ONNG
	IANNG
	DNNG
	RANNG
	RIANNG
)

const (
	// -------------------------------------------------------------
	// ErrorCode is false
	// -------------------------------------------------------------.
	ErrorCode = C._Bool(false)
	// -------------------------------------------------------------.

	NormalStatistics statisticsType = iota - 1
	AdditionalStatistics
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
	case InnerProduct:
		return "InnerProduct"
	}
	return "Unknown"
}

func (i indexType) String() string {
	switch i {
	case GraphAndTree:
		return "GraphAndTree"
	case Graph:
		return "Graph"
	}
	return "Unknown"
}

func (d databaseType) String() string {
	switch d {
	case Memory:
		return "Memory"
	case MemoryMappedFile:
		return "MemoryMappedFile"
	}
	return "Unknown"
}

func (o objectAlignment) String() string {
	switch o {
	case ObjectAlignmentTrue:
		return "true"
	case ObjectAlignmentFalse:
		return "false"
	}
	return "Unknown"
}

func (s seedType) String() string {
	switch s {
	case RandomNodes:
		return "RandomNodes"
	case FixedNodes:
		return "FixedNodes"
	case FirstNode:
		return "FirstNode"
	case AllLeafNodes:
		return "AllLeafNodes"
	}
	return "Unknown"
}

func (g graphType) String() string {
	switch g {
	case ANNG:
		return "ANNG"
	case KNNG:
		return "KNNG"
	case BKNNG:
		return "BKNNG"
	case ONNG:
		return "ONNG"
	case IANNG:
		return "IANNG"
	case DNNG:
		return "DNNG"
	case RANNG:
		return "RANNG"
	case RIANNG:
		return "RIANNG"
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
	n.smu = new(sync.Mutex)
	n.group = singleflight.New[*GraphStatistics]()
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
	n.epool = sync.Pool[*ngtError]{New: newNGTError}

	ne := n.GetErrorBuffer()
	n.prop = C.ngt_create_property(ne.err)
	if n.prop == nil {
		return errors.ErrCreateProperty(n.newGoError(ne))
	}
	n.PutErrorBuffer(ne)
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
	path := C.CString(n.idxPath)
	defer C.free(unsafe.Pointer(path))

	n.smu.Lock()
	defer n.smu.Unlock()
	ne := n.GetErrorBuffer()
	if !n.inMemory {
		n.index = C.ngt_create_graph_and_tree(path, n.prop, ne.err)
		if n.index == nil {
			return n.newGoError(ne)
		}
		if C.ngt_save_index(n.index, path, ne.err) == ErrorCode {
			return n.newGoError(ne)
		}
	} else {
		n.index = C.ngt_create_graph_and_tree_in_memory(n.prop, ne.err)
		if n.index == nil {
			return n.newGoError(ne)
		}
	}
	n.PutErrorBuffer(ne)

	return nil
}

func (n *ngt) open() error {
	if !file.Exists(n.idxPath) {
		return errors.ErrIndexFileNotFound
	}

	path := C.CString(n.idxPath)
	defer C.free(unsafe.Pointer(path))

	n.smu.Lock()
	defer n.smu.Unlock()
	ne := n.GetErrorBuffer()
	n.index = C.ngt_open_index(path, ne.err)
	if n.index == nil {
		return n.newGoError(ne)
	}

	if C.ngt_get_property(n.index, n.prop, ne.err) == ErrorCode {
		return n.newGoError(ne)
	}

	n.dimension = C.ngt_get_property_dimension(n.prop, ne.err)
	if int(n.dimension) == -1 {
		return n.newGoError(ne)
	}
	n.PutErrorBuffer(ne)
	return nil
}

func (n *ngt) loadObjectSpace() error {
	ne := n.GetErrorBuffer()
	n.ospace = C.ngt_get_object_space(n.index, ne.err)
	if n.ospace == nil {
		return n.newGoError(ne)
	}
	n.PutErrorBuffer(ne)
	return nil
}

// Search returns search result as []algorithm.SearchResult.
// It searches for the k-nearest neighbors of the given vector.
//
// Rationale:
// - Uses a read lock to allow concurrent searches while blocking writes (like Insert/Remove).
// - Calls CGO (C.ngt_search_index_as_float) which requires careful memory management for the results object.
func (n *ngt) Search(
	ctx context.Context, vec []float32, size int, epsilon, radius float32, edgeSize int32,
) (result []algorithm.SearchResult, err error) {
	if len(vec) != int(n.dimension) {
		return nil, errors.ErrIncompatibleDimensionSize(len(vec), int(n.dimension))
	}
	if algorithm.HasNaNOrInf(vec) {
		return nil, errors.ErrNaNOrInfVector
	}
	return n.search(ctx, vec, size, epsilon, radius, edgeSize)
}

// SearchValidated returns search result as []algorithm.SearchResult without re-checking vec for NaN/Inf.
func (n *ngt) SearchValidated(
	ctx context.Context, vec []float32, size int, epsilon, radius float32, edgeSize int32,
) (result []algorithm.SearchResult, err error) {
	return n.SearchValidatedInto(ctx, vec, size, epsilon, radius, edgeSize, nil)
}

// SearchValidatedInto performs approximate search without re-checking vec for NaN/Inf,
// writing results into dst when cap(dst) is sufficient.
func (n *ngt) SearchValidatedInto(
	ctx context.Context,
	vec []float32,
	size int,
	epsilon, radius float32,
	edgeSize int32,
	dst []algorithm.SearchResult,
) (result []algorithm.SearchResult, err error) {
	if len(vec) != int(n.dimension) {
		return nil, errors.ErrIncompatibleDimensionSize(len(vec), int(n.dimension))
	}
	return n.searchInto(ctx, false, vec, size, epsilon, radius, edgeSize, dst)
}

const stackSearchResultCap = 64

func (n *ngt) search(
	ctx context.Context, vec []float32, size int, epsilon, radius float32, edgeSize int32,
) ([]algorithm.SearchResult, error) {
	return n.searchInto(ctx, false, vec, size, epsilon, radius, edgeSize, nil)
}

func (n *ngt) searchInto(
	ctx context.Context,
	linear bool,
	vec []float32,
	size int,
	epsilon, radius float32,
	edgeSize int32,
	dst []algorithm.SearchResult,
) (result []algorithm.SearchResult, err error) {
	if size < 0 {
		size = 0
	}
	if !linear {
		if epsilon == 0 {
			epsilon = n.epsilon
		}
		if radius == 0 {
			radius = n.radius
		}
	}

	outCap := max(size, stackSearchResultCap)
	var (
		stackIDs   [stackSearchResultCap]C.uint32_t
		stackDists [stackSearchResultCap]C.float
		outIDs     []C.uint32_t
		outDists   []C.float
	)
	if outCap == stackSearchResultCap {
		outIDs = stackIDs[:]
		outDists = stackDists[:]
	} else {
		outIDs = make([]C.uint32_t, outCap)
		outDists = make([]C.float, outCap)
	}

	var cvec *C.float
	if len(vec) > 0 {
		cvec = (*C.float)(&vec[0])
	}

	ne := n.GetErrorBuffer()
	var (
		rsize  C.int
		nobj   C.uint64_t
		badIdx C.int
	)
	n.mu.RLock()
	ok := C.vald_ngt_search_float(
		C.bool(linear),
		n.index,
		cvec,
		n.dimension,
		C.size_t(size),
		C.float(epsilon),
		C.float(radius),
		C.int(edgeSize),
		&outIDs[0],
		&outDists[0],
		C.int(outCap),
		&rsize,
		&nobj,
		&badIdx,
		ne.err,
	)
	n.mu.RUnlock()

	if !ok {
		return nil, n.newGoError(ne)
	}
	if int(rsize) <= 0 {
		if err = n.newGoError(ne); err != nil {
			return nil, err
		}
		if nobj == 0 {
			return nil, errors.ErrSearchResultEmptyButNoDataStored
		}
		return nil, errors.ErrEmptySearchResult
	}
	if cerr := ctx.Err(); cerr != nil {
		n.PutErrorBuffer(ne)
		return nil, cerr
	}
	if int(rsize) > outCap {
		n.PutErrorBuffer(ne)
		return n.searchSlow(ctx, linear, vec, size, epsilon, radius, edgeSize, dst)
	}

	total := int(rsize)
	if cap(dst) >= total {
		result = dst[:total]
	} else {
		result = make([]algorithm.SearchResult, total)
	}
	for i := range result {
		result[i] = algorithm.SearchResult{
			ID:       uint32(outIDs[i]),
			Distance: float32(outDists[i]),
			Error:    nil,
		}
	}
	if badIdx >= 0 {
		result[int(badIdx)].Error = n.newGoError(ne)
		return result, nil
	}
	n.PutErrorBuffer(ne)
	return result, nil
}

func (n *ngt) searchSlow(
	ctx context.Context,
	linear bool,
	vec []float32,
	size int,
	epsilon, radius float32,
	edgeSize int32,
	dst []algorithm.SearchResult,
) (result []algorithm.SearchResult, err error) {
	ne := n.GetErrorBuffer()
	results := C.ngt_create_empty_results(ne.err)
	if results == nil {
		return nil, n.newGoError(ne)
	}

	var (
		p    runtime.Pinner
		cvec *C.float
	)
	if len(vec) > 0 {
		p.Pin(&vec[0])
		cvec = (*C.float)(&vec[0])
	}

	n.mu.RLock()
	var ret C.bool
	if linear {
		ret = C.ngt_linear_search_index_as_float(
			n.index,
			cvec,
			n.dimension,
			C.size_t(size),
			results,
			ne.err)
	} else {
		ret = C.ngt_search_index_as_float(
			n.index,
			cvec,
			n.dimension,
			C.size_t(size),
			C.float(epsilon),
			C.float(radius),
			results,
			C.int(edgeSize),
			ne.err)
	}
	n.mu.RUnlock()
	if len(vec) > 0 {
		p.Unpin()
	}

	if ret == ErrorCode {
		C.ngt_destroy_results(results)
		return nil, n.newGoError(ne)
	}

	result, err = n.toSearchResultsInto(ctx, results, ne, dst)
	C.ngt_destroy_results(results)
	return result, err
}

// Linear Search returns linear search result as []algorithm.SearchResult.
func (n *ngt) LinearSearch(
	ctx context.Context, vec []float32, size int,
) (result []algorithm.SearchResult, err error) {
	if len(vec) != int(n.dimension) {
		return nil, errors.ErrIncompatibleDimensionSize(len(vec), int(n.dimension))
	}
	if algorithm.HasNaNOrInf(vec) {
		return nil, errors.ErrNaNOrInfVector
	}
	return n.searchInto(ctx, true, vec, size, 0, 0, 0, nil)
}

// LinearSearchValidated returns linear search result as []algorithm.SearchResult without re-checking vec for NaN/Inf.
func (n *ngt) LinearSearchValidated(
	ctx context.Context, vec []float32, size int,
) (result []algorithm.SearchResult, err error) {
	return n.LinearSearchValidatedInto(ctx, vec, size, nil)
}

// LinearSearchValidatedInto performs linear search without re-checking vec for NaN/Inf,
// writing results into dst when cap(dst) is sufficient.
func (n *ngt) LinearSearchValidatedInto(
	ctx context.Context, vec []float32, size int, dst []algorithm.SearchResult,
) (result []algorithm.SearchResult, err error) {
	if len(vec) != int(n.dimension) {
		return nil, errors.ErrIncompatibleDimensionSize(len(vec), int(n.dimension))
	}
	return n.searchInto(ctx, true, vec, size, 0, 0, 0, dst)
}

// toSearchResults converts the C results of Search and LinearSearch into
// []algorithm.SearchResult. It takes over ne: every return path releases it
// via PutErrorBuffer or newGoError. The caller still owns and destroys results.
func (n *ngt) toSearchResults(
	ctx context.Context, results C.NGTObjectDistances, ne *ngtError,
) (result []algorithm.SearchResult, err error) {
	return n.toSearchResultsInto(ctx, results, ne, nil)
}

func (n *ngt) toSearchResultsInto(
	ctx context.Context, results C.NGTObjectDistances, ne *ngtError, dst []algorithm.SearchResult,
) (result []algorithm.SearchResult, err error) {
	rsize := int(C.ngt_get_result_size(results, ne.err))
	if rsize <= 0 {
		// n.cnt only counts Insert/Remove in this process and stays 0 for an index
		// restored by Load, so ask NGT how many objects are stored.
		n.mu.RLock()
		nobj := C.ngt_get_number_of_objects(n.index, ne.err)
		n.mu.RUnlock()
		if err = n.newGoError(ne); err != nil {
			return nil, err
		}
		if nobj == 0 {
			return nil, errors.ErrSearchResultEmptyButNoDataStored
		}
		return nil, errors.ErrEmptySearchResult
	}
	if err := ctx.Err(); err != nil {
		n.PutErrorBuffer(ne)
		return nil, err
	}
	if cap(dst) >= rsize {
		result = dst[:rsize]
	} else {
		result = make([]algorithm.SearchResult, rsize)
	}

	for i := range result {
		d := C.ngt_get_result(results, C.uint32_t(i), ne.err)
		if d.id == 0 && d.distance == 0 {
			result[i] = algorithm.SearchResult{
				ID:       0,
				Distance: 0,
				Error:    n.newGoError(ne),
			}
			ne = n.GetErrorBuffer()
		} else {
			result[i] = algorithm.SearchResult{
				ID:       uint32(d.id),
				Distance: float32(d.distance),
				Error:    nil,
			}
		}
	}
	n.PutErrorBuffer(ne)

	return result, nil
}

// HasVector reports whether id is stored in the NGT object space without copying its vector.
func (n *ngt) HasVector(id uint) (bool, error) {
	if id == 0 {
		return false, nil
	}
	ne := n.GetErrorBuffer()
	var exists bool
	switch n.objectType {
	case Float:
		n.mu.RLock()
		if n.ospace == nil {
			n.mu.RUnlock()
			n.PutErrorBuffer(ne)
			return false, errIndexClosed
		}
		exists = C.ngt_get_object_as_float(n.ospace, C.ObjectID(id), ne.err) != nil
		n.mu.RUnlock()
	case HalfFloat:
		n.mu.RLock()
		if n.ospace == nil {
			n.mu.RUnlock()
			n.PutErrorBuffer(ne)
			return false, errIndexClosed
		}
		p := C.ngt_get_allocated_object_as_float(n.ospace, C.ObjectID(id), ne.err)
		n.mu.RUnlock()
		if p != nil {
			exists = true
			C.free(unsafe.Pointer(p))
		}
	case Uint8:
		n.mu.RLock()
		if n.ospace == nil {
			n.mu.RUnlock()
			n.PutErrorBuffer(ne)
			return false, errIndexClosed
		}
		exists = C.ngt_get_object_as_integer(n.ospace, C.ObjectID(id), ne.err) != nil
		n.mu.RUnlock()
	case ObjectNone:
		n.PutErrorBuffer(ne)
		return false, errors.ErrUnsupportedObjectType
	default:
		n.PutErrorBuffer(ne)
		return false, errors.ErrUnsupportedObjectType
	}
	if !exists {
		return false, n.newGoError(ne)
	}
	n.PutErrorBuffer(ne)
	return true, nil
}

// Insert returns NGT object id.
// It inserts a vector into the NGT object repository but does not build the index immediately.
//
// Rationale:
// - Acquire a write lock to ensure thread safety during CGO call.
// - You must call CreateIndex and SaveIndex to make the vector searchable.
func (n *ngt) Insert(vec []float32) (id uint, err error) {
	if len(vec) != int(n.dimension) {
		return 0, errors.ErrIncompatibleDimensionSize(len(vec), int(n.dimension))
	}
	if algorithm.HasNaNOrInf(vec) {
		return 0, errors.ErrNaNOrInfVector
	}

	var (
		p    runtime.Pinner
		cvec *C.float
	)
	if len(vec) > 0 {
		p.Pin(&vec[0])
		// cvec must be obtained after p.Pin to ensure the backing array is pinned before passing to CGo.
		cvec = (*C.float)(&vec[0])
	}
	dim := C.uint32_t(n.dimension)
	ne := n.GetErrorBuffer()
	n.mu.Lock()
	oid := C.ngt_insert_index_as_float(n.index, cvec, dim, ne.err)
	n.mu.Unlock()
	if len(vec) > 0 {
		p.Unpin()
	}
	id = uint(oid)
	if id == 0 {
		return 0, n.newGoError(ne)
	}
	n.PutErrorBuffer(ne)
	n.cnt.Add(1)

	return id, nil
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
	n.smu.Lock()
	defer n.smu.Unlock()
	ne := n.GetErrorBuffer()
	n.mu.Lock()
	ret := C.ngt_create_index(n.index, C.uint32_t(poolSize), ne.err)
	n.mu.Unlock()
	if ret == ErrorCode {
		return n.newGoError(ne)
	}
	n.PutErrorBuffer(ne)

	return nil
}

// SaveIndex stores NGT index to storage.
func (n *ngt) SaveIndex() error {
	return n.saveIndex(n.idxPath)
}

// SaveIndexWithPath stores NGT index to specified storage.
func (n *ngt) SaveIndexWithPath(idxPath string) error {
	if len(idxPath) == 0 {
		return nil
	}
	return n.saveIndex(idxPath)
}

func (n *ngt) saveIndex(idxPath string) error {
	if !n.inMemory {
		path := C.CString(idxPath)
		defer C.free(unsafe.Pointer(path))
		n.smu.Lock()
		defer n.smu.Unlock()
		n.mu.RLock()
		if n.index == nil {
			n.mu.RUnlock()
			return errIndexClosed
		}
		ne := n.GetErrorBuffer()
		ret := C.ngt_save_index(n.index, path, ne.err)
		n.mu.RUnlock()
		if ret == ErrorCode {
			return n.newGoError(ne)
		}
		n.PutErrorBuffer(ne)
	}

	return nil
}

// Remove removes from NGT index.
func (n *ngt) Remove(id uint) error {
	ne := n.GetErrorBuffer()
	n.mu.Lock()
	ret := C.ngt_remove_index(n.index, C.ObjectID(id), ne.err)
	n.mu.Unlock()
	if ret == ErrorCode {
		return n.newGoError(ne)
	}
	n.PutErrorBuffer(ne)
	n.cnt.Add(^uint64(0))
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
	return n.GetVectorInto(id, nil)
}

// GetVectorInto writes the vector stored in NGT index into dst when cap(dst) is sufficient.
func (n *ngt) GetVectorInto(id uint, dst []float32) (ret []float32, err error) {
	dimension := int(n.dimension)
	allocVec := func() []float32 {
		if dimension > 0 && cap(dst) >= dimension {
			return dst[:dimension]
		}
		return make([]float32, dimension)
	}
	ne := n.GetErrorBuffer()
	switch n.objectType {
	case Float:
		// results points into NGT's object repository, which Remove frees under the
		// write lock, so copy it out before releasing the read lock.
		n.mu.RLock()
		if n.ospace == nil {
			n.mu.RUnlock()
			n.PutErrorBuffer(ne)
			return nil, errIndexClosed
		}
		results := C.ngt_get_object_as_float(n.ospace, C.ObjectID(id), ne.err)
		if results != nil {
			ret = allocVec()
			copy(ret, unsafe.Slice((*float32)(unsafe.Pointer(results)), dimension))
		}
		n.mu.RUnlock()
		if results == nil {
			return nil, n.newGoError(ne)
		}
	case HalfFloat:
		n.mu.RLock()
		if n.ospace == nil {
			n.mu.RUnlock()
			n.PutErrorBuffer(ne)
			return nil, errIndexClosed
		}
		results := C.ngt_get_allocated_object_as_float(n.ospace, C.ObjectID(id), ne.err)
		n.mu.RUnlock()
		if results == nil {
			return nil, n.newGoError(ne)
		}
		ret = allocVec()
		copy(ret, unsafe.Slice((*float32)(unsafe.Pointer(results)), dimension))
		C.free(unsafe.Pointer(results))
	case Uint8:
		n.mu.RLock()
		if n.ospace == nil {
			n.mu.RUnlock()
			n.PutErrorBuffer(ne)
			return nil, errIndexClosed
		}
		results := C.ngt_get_object_as_integer(n.ospace, C.ObjectID(id), ne.err)
		if results != nil {
			ret = allocVec()
			for i, elem := range unsafe.Slice((*C.uint8_t)(unsafe.Pointer(results)), dimension) {
				ret[i] = float32(elem)
			}
		}
		n.mu.RUnlock()
		if results == nil {
			return nil, n.newGoError(ne)
		}
	default:
		n.PutErrorBuffer(ne)
		return nil, errors.ErrUnsupportedObjectType
	}
	n.PutErrorBuffer(ne)
	return ret, nil
}

func (n *ngt) newGoError(ne *ngtError) (err error) {
	msg := C.GoString(C.ngt_get_error_string(ne.err))
	C.ngt_clear_error_string(ne.err)
	n.PutErrorBuffer(ne)
	if len(msg) == 0 {
		return nil
	}
	return errors.NewNGTError(msg)
}

// Close NGT without save index.
// ngt_close_index only frees the index (it does not save), so this is the same as Close.
func (n *ngt) CloseWithoutSaveIndex() {
	n.Close()
}

func (n *ngt) GetErrorBuffer() (ne *ngtError) {
	ne = n.epool.Get()
	if ne == nil {
		ne = newNGTError()
	}
	return ne
}

func (n *ngt) PutErrorBuffer(ne *ngtError) {
	n.epool.Put(ne)
}

// Close NGT index.
func (n *ngt) Close() {
	// gen sets the mutexes before it opens the index, so an ngt without them has
	// no index to close.
	if n.smu == nil || n.mu == nil {
		return
	}
	// Same lock order as CreateIndex/SaveIndex (smu -> mu), so in-flight
	// Search/Insert/Remove/GetVector/statistics finish before the index is freed.
	n.smu.Lock()
	defer n.smu.Unlock()
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.index != nil {
		C.ngt_close_index(n.index)
		n.index = nil
		n.prop = nil
		n.ospace = nil
	}
}

func fromCGraphStatistics(cstats *C.NGTGraphStatistics) *GraphStatistics {
	if cstats == nil {
		return nil
	}
	goStats := &GraphStatistics{
		C1Indegree:                       float64(cstats.c1Indegree),
		C5Indegree:                       float64(cstats.c5Indegree),
		C95Outdegree:                     float64(cstats.c95Outdegree),
		C99Outdegree:                     float64(cstats.c99Outdegree),
		MaxNumberOfIndegree:              uint64(cstats.maxNumberOfIndegree),
		MaxNumberOfOutdegree:             uint64(cstats.maxNumberOfOutdegree),
		MeanEdgeLength:                   float64(cstats.meanEdgeLength),
		MeanEdgeLengthFor10Edges:         float64(cstats.meanEdgeLengthFor10Edges),
		MeanIndegreeDistanceFor10Edges:   float64(cstats.meanIndegreeDistanceFor10Edges),
		MeanNumberOfEdgesPerNode:         float64(cstats.meanNumberOfEdgesPerNode),
		MedianIndegree:                   int32(cstats.medianIndegree),
		MedianOutdegree:                  int32(cstats.medianOutdegree),
		MinNumberOfIndegree:              uint64(cstats.minNumberOfIndegree),
		MinNumberOfOutdegree:             uint64(cstats.minNumberOfOutdegree),
		ModeIndegree:                     uint64(cstats.modeIndegree),
		ModeOutdegree:                    uint64(cstats.modeOutdegree),
		NodesSkippedFor10Edges:           uint64(cstats.nodesSkippedFor10Edges),
		NodesSkippedForIndegreeDistance:  uint64(cstats.nodesSkippedForIndegreeDistance),
		NumberOfEdges:                    uint64(cstats.numberOfEdges),
		NumberOfIndexedObjects:           uint64(cstats.numberOfIndexedObjects),
		NumberOfNodes:                    uint64(cstats.numberOfNodes),
		NumberOfNodesWithoutEdges:        uint64(cstats.numberOfNodesWithoutEdges),
		NumberOfNodesWithoutIndegree:     uint64(cstats.numberOfNodesWithoutIndegree),
		NumberOfObjects:                  uint64(cstats.numberOfObjects),
		NumberOfRemovedObjects:           uint64(cstats.numberOfRemovedObjects),
		SizeOfObjectRepository:           uint64(cstats.sizeOfObjectRepository),
		SizeOfRefinementObjectRepository: uint64(cstats.sizeOfRefinementObjectRepository),
		VarianceOfIndegree:               float64(cstats.varianceOfIndegree),
		VarianceOfOutdegree:              float64(cstats.varianceOfOutdegree),
		Valid:                            bool(cstats.valid),
	}

	// Convert indegreeCount
	indegreeCountSize := int(cstats.indegreeCountSize)
	goStats.IndegreeCount = make([]int64, indegreeCountSize)
	cIndegreeCount := (*[1 << 30]C.size_t)(unsafe.Pointer(cstats.indegreeCount))[:indegreeCountSize:indegreeCountSize]
	for i := range indegreeCountSize {
		goStats.IndegreeCount[i] = int64(cIndegreeCount[i])
	}

	// Convert outdegreeHistogram
	outdegreeHistogramSize := int(cstats.outdegreeHistogramSize)
	goStats.OutdegreeHistogram = make([]uint64, outdegreeHistogramSize)
	cOutdegreeHistogram := (*[1 << 30]C.size_t)(unsafe.Pointer(cstats.outdegreeHistogram))[:outdegreeHistogramSize:outdegreeHistogramSize]
	for i := range outdegreeHistogramSize {
		goStats.OutdegreeHistogram[i] = uint64(cOutdegreeHistogram[i])
	}

	// Convert indegreeHistogram
	indegreeHistogramSize := int(cstats.indegreeHistogramSize)
	goStats.IndegreeHistogram = make([]uint64, indegreeHistogramSize)
	cIndegreeHistogram := (*[1 << 30]C.size_t)(unsafe.Pointer(cstats.indegreeHistogram))[:indegreeHistogramSize:indegreeHistogramSize]
	for i := range indegreeHistogramSize {
		goStats.IndegreeHistogram[i] = uint64(cIndegreeHistogram[i])
	}

	return goStats
}

func (n *ngt) GetGraphStatistics(
	ctx context.Context, m statisticsType,
) (stats *GraphStatistics, err error) {
	return n.getGraphStatistics(ctx, m, 10)
}

func (n *ngt) getGraphStatistics(
	ctx context.Context, m statisticsType, cnt int,
) (stats *GraphStatistics, err error) {
	var shared bool
	stats, shared, err = n.group.Do(ctx, "GetGraphStatistics", func(context.Context) (stats *GraphStatistics, err error) {
		n.smu.Lock()
		defer n.smu.Unlock()
		var mode rune
		switch m {
		case NormalStatistics:
			mode = '-'
		case AdditionalStatistics:
			mode = 'a'
		}
		err = safety.RecoverFunc(func() (err error) {
			ne := n.GetErrorBuffer()
			// The walk reads the whole graph, so exclude Insert/Remove like Search does.
			n.mu.RLock()
			cstats := C.ngt_get_graph_statistics(n.index, C.char(mode), C.size_t(n.ces), ne.err)
			n.mu.RUnlock()
			// On the NULL-index and exception paths NGT sets the error string and
			// returns cstats uninitialized, so it must be neither read nor freed.
			if err := n.newGoError(ne); err != nil {
				return err
			}
			// Otherwise the arrays are allocated even when the graph is not valid.
			defer C.ngt_free_graph_statistics(&cstats)
			if !cstats.valid {
				return nil
			}
			stats = fromCGraphStatistics(&cstats)
			if stats == nil {
				return errors.ErrNGTIndexStatisticsNotReady
			}
			return nil
		})()
		if err != nil {
			return nil, err
		}
		return stats, nil
	})
	if err != nil {
		if shared && cnt > 0 && !errors.Is(err, errors.ErrNGTIndexStatisticsNotReady) {
			cnt--
			return n.getGraphStatistics(ctx, m, cnt)
		}
		return nil, err
	}
	return stats, nil
}

func (n *ngt) GetProperty() (prop *Property, err error) {
	n.mu.RLock()
	if n.index == nil {
		n.mu.RUnlock()
		return nil, errIndexClosed
	}
	ne := n.GetErrorBuffer()
	cprop := C.ngt_get_property_info(n.index, ne.err)
	n.mu.RUnlock()
	// newGoError returns ne to the pool only when it holds no error message.
	if err := n.newGoError(ne); err != nil {
		return nil, err
	}
	// AccuracyTable and SearchType are left empty on purpose: ngt_get_property_info
	// (NGT Capi.cpp) fills cprop.accuracy_table / cprop.search_type with c_str() of
	// std::string members of a local NGT::Property that is destroyed before it returns,
	// so both pointers dangle and must not be read. Vald never sets either value
	// (NGT defaults both to ""), and the C API offers no other way to read them.
	return &Property{
		Dimension:                     int32(cprop.dimension),
		ThreadPoolSize:                int32(cprop.thread_pool_size),
		ObjectType:                    objectType(cprop.object_type),
		DistanceType:                  distanceType(cprop.distance_type),
		IndexType:                     indexType(cprop.index_type),
		DatabaseType:                  databaseType(cprop.database_type),
		ObjectAlignment:               objectAlignment(cprop.object_alignment),
		PathAdjustmentInterval:        int32(cprop.path_adjustment_interval),
		GraphSharedMemorySize:         int32(cprop.graph_shared_memory_size),
		TreeSharedMemorySize:          int32(cprop.tree_shared_memory_size),
		ObjectSharedMemorySize:        int32(cprop.object_shared_memory_size),
		PrefetchOffset:                int32(cprop.prefetch_offset),
		PrefetchSize:                  int32(cprop.prefetch_size),
		MaxMagnitude:                  float32(cprop.max_magnitude),
		NOfNeighborsForInsertionOrder: int32(cprop.n_of_neighbors_for_insertion_order),
		EpsilonForInsertionOrder:      float32(cprop.epsilon_for_insertion_order),
		EpsilonForCreation:            float32(cprop.epsilon_for_creation),
		RefinementObjectType:          objectType(cprop.refinement_object_type),
		TruncationThreshold:           int32(cprop.truncation_threshold),
		EdgeSizeForCreation:           int32(cprop.edge_size_for_creation),
		EdgeSizeForSearch:             int32(cprop.edge_size_for_search),
		EdgeSizeLimitForCreation:      int32(cprop.max_edge_size_for_creation),
		InsertionRadiusCoefficient:    float64(cprop.insertion_radius_coefficient),
		SeedSize:                      int32(cprop.seed_size),
		SeedType:                      seedType(cprop.seed_type),
		TruncationThreadPoolSize:      int32(cprop.truncation_thread_pool_size),
		BatchSizeForCreation:          int32(cprop.batch_size_for_creation),
		GraphType:                     graphType(cprop.graph_type),
		DynamicEdgeSizeBase:           int32(cprop.dynamic_edge_size_base),
		DynamicEdgeSizeRate:           int32(cprop.dynamic_edge_size_rate),
		BuildTimeLimit:                float32(cprop.build_time_limit),
		OutgoingEdge:                  int32(cprop.outgoing_edge),
		IncomingEdge:                  int32(cprop.incoming_edge),
	}, nil
}
