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
	"github.com/vdaas/vald/hack/benchmark/assets/x1b"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
)

const (
	largeDatasetPath = "hack/benchmark/assets/dataset/large"
)

type largeDataset struct {
	*dataset
	train       x1b.BillionScaleVectors
	query       x1b.BillionScaleVectors
	groundTruth [][]int
	distances   x1b.FloatVectors
}

func loadLargeData(trainFileName, queryFileName, groundTruthFileName, distanceFileName, name, distanceType, objectType string) func() (Dataset, error) {
	return func() (Dataset, error) {
		dir, err := findDir(largeDatasetPath)
		if err != nil {
			return nil, err
		}
		train, err := x1b.Open(file.Join(dir, trainFileName))
		if err != nil {
			return nil, err
		}
		query, err := x1b.Open(file.Join(dir, queryFileName))
		if err != nil {
			return nil, err
		}
		tdim := train.Dimension()
		qdim := query.Dimension()
		if tdim != qdim {
			return nil, errors.New("dimension must be same train and query.")
		}
		iv, err := x1b.NewInt32Vectors(file.Join(dir, groundTruthFileName))
		if err != nil {
			return nil, err
		}
		groundTruth := make([][]int, 0, iv.Size())
		for i := 0; ; i++ {
			gt32, err := iv.LoadInt32(i)
			if err == ErrOutOfBounds {
				break
			}
			gt := make([]int, 0, len(gt32))
			for _, v := range gt32 {
				gt = append(gt, int(v))
			}
			groundTruth = append(groundTruth, gt)
		}

		distances, err := x1b.NewFloatVectors(file.Join(dir, distanceFileName))
		if err != nil {
			return nil, err
		}
		return &largeDataset{
			dataset: &dataset{
				name:         name,
				dimension:    tdim,
				distanceType: distanceType,
				objectType:   objectType,
			},
			train:       train,
			query:       query,
			groundTruth: groundTruth,
			distances:   distances,
		}, nil
	}
}

func (d *largeDataset) Train(i int) (interface{}, error) {
	return d.train.Load(i)
}

func (d *largeDataset) TrainSize() int {
	return d.train.Size()
}

func (d *largeDataset) Query(i int) (interface{}, error) {
	return d.query.Load(i)
}

func (d *largeDataset) QuerySize() int {
	return d.query.Size()
}

func (d *largeDataset) Distance(i int) ([]float32, error) {
	return d.distances.LoadFloat32(i)
}

func (d *largeDataset) DistanceSize() int {
	return d.distances.Size()
}

func (d *largeDataset) Neighbor(i int) ([]int, error) {
	if i >= len(d.groundTruth) {
		return nil, ErrOutOfBounds
	}
	return d.groundTruth[i], nil
}

func (d *largeDataset) NeighborSize() int {
	return len(d.groundTruth)
}

func (d *largeDataset) Dimension() int {
	return d.dimension
}

func (d *largeDataset) DistanceType() string {
	return d.distanceType
}

func (d *largeDataset) ObjectType() string {
	return d.objectType
}

func (d *largeDataset) Name() string {
	return d.name
}
