//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Dataset interface {
	Train() [][]float32
	TrainAsFloat64() [][]float64
	Query() [][]float32
	QueryAsFloat64() [][]float64
	Distances() [][]float32
	DistancesAsFloat64() [][]float64
	Neighbors() [][]int
	IDs() []string
	Name() string
	Dimension() int
	DistanceType() string
	ObjectType() string
}

type dataset struct {
	train              [][]float32
	trainAsFloat64     [][]float64
	trainOnce          sync.Once
	query              [][]float32
	queryAsFloat64     [][]float64
	queryOnce          sync.Once
	distances          [][]float32
	distancesAsFloat64 [][]float64
	distancesOnce      sync.Once
	neighbors          [][]int
	ids                []string
	name               string
	dimension          int
	distanceType       string
	objectType         string
}

var (
	data = map[string]func() (Dataset, error){
		"fashion-mnist": func() (Dataset, error) {
			dir, err := datasetDir()
			if err != nil {
				return nil, err
			}
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(dir + "fashion-mnist-784-euclidean.hdf5")
			if err != nil {
				return nil, err
			}
			return &dataset{
				train:        train,
				query:        query,
				distances:    distances,
				neighbors:    neighbors,
				ids:          ids,
				name:         "fashion-mnist",
				dimension:    dim,
				distanceType: "l2",
				objectType:   "float",
			}, nil
		},
		"mnist": func() (Dataset, error) {
			dir, err := datasetDir()
			if err != nil {
				return nil, err
			}
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(dir + "mnist-784-euclidean.hdf5")
			if err != nil {
				return nil, err
			}
			return &dataset{
				train:        train,
				query:        query,
				distances:    distances,
				neighbors:    neighbors,
				ids:          ids,
				name:         "mnist",
				dimension:    dim,
				distanceType: "l2",
				objectType:   "float",
			}, err
		},
		"glove-25": func() (Dataset, error) {
			dir, err := datasetDir()
			if err != nil {
				return nil, err
			}
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(dir + "glove-25-angular.hdf5")
			if err != nil {
				return nil, err
			}
			return &dataset{
				train:        train,
				query:        query,
				distances:    distances,
				neighbors:    neighbors,
				ids:          ids,
				name:         "glove-25",
				dimension:    dim,
				distanceType: "cosine",
				objectType:   "float",
			}, nil
		},
		"glove-50": func() (Dataset, error) {
			dir, err := datasetDir()
			if err != nil {
				return nil, err
			}
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(dir + "glove-50-angular.hdf5")
			if err != nil {
				return nil, err
			}
			return &dataset{
				train:        train,
				query:        query,
				distances:    distances,
				neighbors:    neighbors,
				ids:          ids,
				name:         "glove-50",
				dimension:    dim,
				distanceType: "cosine",
				objectType:   "float",
			}, nil
		},
		"glove-100": func() (Dataset, error) {
			dir, err := datasetDir()
			if err != nil {
				return nil, err
			}
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(dir + "glove-100-angular.hdf5")
			if err != nil {
				return nil, err
			}
			return &dataset{
				train:        train,
				query:        query,
				distances:    distances,
				neighbors:    neighbors,
				ids:          ids,
				name:         "glove-100",
				dimension:    dim,
				distanceType: "cosine",
				objectType:   "float",
			}, nil
		},
		"glove-200": func() (Dataset, error) {
			dir, err := datasetDir()
			if err != nil {
				return nil, err
			}
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(dir + "glove-200-angular.hdf5")
			if err != nil {
				return nil, err
			}
			return &dataset{
				train:        train,
				query:        query,
				distances:    distances,
				neighbors:    neighbors,
				ids:          ids,
				name:         "glove-200",
				dimension:    dim,
				distanceType: "cosine",
				objectType:   "float",
			}, nil
		},
		"nytimes": func() (Dataset, error) {
			dir, err := datasetDir()
			if err != nil {
				return nil, err
			}
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(dir + "nytimes-256-angular.hdf5")
			if err != nil {
				return nil, err
			}
			return &dataset{
				train:        train,
				query:        query,
				distances:    distances,
				neighbors:    neighbors,
				ids:          ids,
				name:         "nytimes",
				dimension:    dim,
				distanceType: "cosine",
				objectType:   "float",
			}, nil
		},
		"sift": func() (Dataset, error) {
			dir, err := datasetDir()
			if err != nil {
				return nil, err
			}
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(dir + "sift-128-euclidean.hdf5")
			if err != nil {
				return nil, err
			}
			return &dataset{
				train:        train,
				query:        query,
				distances:    distances,
				neighbors:    neighbors,
				ids:          ids,
				name:         "sift",
				dimension:    dim,
				distanceType: "l2",
				objectType:   "float",
			}, nil
		},
		"gist": func() (Dataset, error) {
			dir, err := datasetDir()
			if err != nil {
				return nil, err
			}
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(dir + "gist-960-euclidean.hdf5")
			if err != nil {
				return nil, err
			}
			return &dataset{
				train:        train,
				query:        query,
				distances:    distances,
				neighbors:    neighbors,
				ids:          ids,
				name:         "gist",
				dimension:    dim,
				distanceType: "l2",
				objectType:   "float",
			}, nil
		},
		"kosarak": func() (Dataset, error) {
			dir, err := datasetDir()
			if err != nil {
				return nil, err
			}
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(dir + "/kosarak-jaccard.hdf5")
			if err != nil {
				return nil, err
			}
			return &dataset{
				train:        train,
				query:        query,
				distances:    distances,
				neighbors:    neighbors,
				ids:          ids,
				name:         "kosarak",
				dimension:    dim,
				distanceType: "jaccard",
				objectType:   "float",
			}, nil
		},
	}
)

func identity(dim int) func() (Dataset, error) {
	return func() (Dataset, error) {
		ids := CreateSequentialIDs(dim * 1000)
		train := make([][]float32, dim)
		for i := range train {
			train[i] = make([]float32, dim)
			train[i][i] = 1
		}
		return &dataset{
			train:        train,
			query:        train,
			ids:          ids,
			name:         fmt.Sprintf("identity-%d", dim),
			dimension:    dim,
			distanceType: "l2",
			objectType:   "float",
		}, nil
	}
}

func random(dim, size int) func() (Dataset, error) {
	return func() (Dataset, error) {
		ids := CreateRandomIDs(size)
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
		return &dataset{
			train: train,
			query: query,
			ids: ids,
			name: fmt.Sprintf("random-%d-%d", dim, size),
			dimension: dim,
			distanceType: "l2",
			objectType: "float",
		}, nil
	}
}

func datasetDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	root := func(cur string) string {
		for {
			parent := filepath.Dir(cur)
			if strings.HasSuffix(parent, "vald/hack") {
				return parent
			} else {
				cur = parent
			}
		}
	}(wd)
	return filepath.Join(root, "benchmark/assets/dataset") + "/", nil
}

func Data(name string) func() (Dataset, error) {
	if strings.HasPrefix(name, "identity-") {
		i, _ := strconv.Atoi(name[9:])
		return identity(i)
	}
	if strings.HasPrefix(name, "random-") {
		l := strings.Split(name[9:], "-")
		i, _ := strconv.Atoi(l[0])
		j, _ := strconv.Atoi(l[1])
		return random(i, j)
	}
	if d, ok := data[name]; ok {
		return d
	}
	return nil
}

func (d *dataset) Train() [][]float32 {
	return d.train
}

func (d *dataset) TrainAsFloat64() [][]float64 {
	d.trainOnce.Do(func() {
		d.trainAsFloat64 = float32To64(d.train)
	})
	return d.trainAsFloat64
}

func (d *dataset) Query() [][]float32 {
	return d.query
}

func (d *dataset) QueryAsFloat64() [][]float64 {
	d.queryOnce.Do(func() {
		d.queryAsFloat64 = float32To64(d.query)
	})
	return d.queryAsFloat64
}

func (d *dataset) Distances() [][]float32 {
	return d.distances
}

func (d *dataset) DistancesAsFloat64() [][]float64 {
	d.distancesOnce.Do(func() {
		d.distancesAsFloat64 = float32To64(d.distances)
	})
	return d.distancesAsFloat64
}

func (d *dataset) Neighbors() [][]int {
	return d.neighbors
}

func (d *dataset) IDs() []string {
	return d.ids
}

func (d *dataset) Name() string {
	return d.name
}

func (d *dataset) Dimension() int {
	return d.dimension
}

func (d *dataset) DistanceType() string {
	return d.distanceType
}

func (d *dataset) ObjectType() string {
	return d.objectType
}

func float32To64(x [][]float32) (y [][]float64) {
	y = make([][]float64, len(x))
	for i, z := range x {
		y[i] = make([]float64, len(z))
		for j, a := range z {
			y[i][j] = float64(a)
		}
	}
	return y
}
