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
package assets

import (
	"fmt"
	"math/rand"

	"github.com/vdaas/vald/internal/file"
)

const (
	smallDatasetPath = "hack/benchmark/assets/dataset"
)

type smallDataset struct {
	*dataset
	train     [][]float32
	query     [][]float32
	distances [][]float32
	neighbors [][]int
}

func loadSmallData(fileName, datasetName, distanceType, objectType string) func() (Dataset, error) {
	return func() (Dataset, error) {
		dir, err := findDir(smallDatasetPath)
		if err != nil {
			return nil, err
		}
		t, q, d, n, dim, err := Load(file.Join(dir, fileName))
		if err != nil {
			return nil, err
		}

		return &smallDataset{
			dataset: &dataset{
				name:         datasetName,
				dimension:    dim,
				distanceType: distanceType,
				objectType:   objectType,
			},
			train:     t,
			query:     q,
			distances: d,
			neighbors: n,
		}, nil
	}
}

func identity(dim int) func() (Dataset, error) {
	return func() (Dataset, error) {
		train := make([][]float32, dim)
		for i := range train {
			train[i] = make([]float32, dim)
			train[i][i] = 1
		}
		return &smallDataset{
			dataset: &dataset{
				name:         fmt.Sprintf("identity-%d", dim),
				dimension:    dim,
				distanceType: "l2",
				objectType:   "float",
			},
			train: train,
			query: train,
		}, nil
	}
}

func random(dim, size int) func() (Dataset, error) {
	return func() (Dataset, error) {
		train := make([][]float32, size)
		query := make([][]float32, size)
		for i := range train {
			train[i] = make([]float32, dim)
			query[i] = make([]float32, dim)
			for j := range train[i] {
				train[i][j] = rand.Float32()
				query[i][j] = rand.Float32()
			}
		}
		return &smallDataset{
			dataset: &dataset{
				name:         fmt.Sprintf("random-%d-%d", dim, size),
				dimension:    dim,
				distanceType: "l2",
				objectType:   "float",
			},
			train: train,
			query: query,
		}, nil
	}
}

func gaussian(dim, size int, mean, stdDev float64) func() (Dataset, error) {
	return func() (Dataset, error) {
		train := make([][]float32, size)
		query := make([][]float32, size)
		for i := range train {
			train[i] = make([]float32, dim)
			query[i] = make([]float32, dim)
			for j := range train[i] {
				train[i][j] = float32(rand.NormFloat64()*stdDev + mean)
				query[i][j] = float32(rand.NormFloat64()*stdDev + mean)
			}
		}
		return &smallDataset{
			dataset: &dataset{
				name:         fmt.Sprintf("gaussian-%d-%d-%f-%f", dim, size, mean, stdDev),
				dimension:    dim,
				distanceType: "l2",
				objectType:   "float",
			},
			train: train,
			query: query,
		}, nil
	}
}

// Train returns vectors for train.
func (s *smallDataset) Train(i int) (interface{}, error) {
	if i >= len(s.train) {
		return nil, ErrOutOfBounds
	}
	return s.train[i], nil
}

// TrainSize return size of vectors for train.
func (s *smallDataset) TrainSize() int {
	return len(s.train)
}

// Query returns vectors for test.
func (s *smallDataset) Query(i int) (interface{}, error) {
	if i >= len(s.query) {
		return nil, ErrOutOfBounds
	}
	return s.query[i], nil
}

// QuerySize return size of vectors for query.
func (s *smallDataset) QuerySize() int {
	return len(s.query)
}

// Distance returns distances between queries and answers.
func (s *smallDataset) Distance(i int) ([]float32, error) {
	if i >= len(s.distances) {
		return nil, ErrOutOfBounds
	}
	return s.distances[i], nil
}

// DistanceSize returns size of distances.
func (s *smallDataset) DistanceSize() int {
	return len(s.distances)
}

// Neighbors returns nearest vectors from queries.
func (s *smallDataset) Neighbor(i int) ([]int, error) {
	if i >= len(s.neighbors) {
		return nil, ErrOutOfBounds
	}
	return s.neighbors[i], nil
}

// NeighborSize returns size of neighbors.
func (s *smallDataset) NeighborSize() int {
	return len(s.neighbors)
}
