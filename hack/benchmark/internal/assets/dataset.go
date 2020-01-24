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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
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
	data = map[string]func(testing.TB) Dataset{
		"fashion-mnist": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(datasetDir(tb) + "fashion-mnist-784-euclidean.hdf5")
			if err != nil {
				tb.Error(err)
				return nil
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
			}
		},
		"mnist": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(datasetDir(tb) + "mnist-784-euclidean.hdf5")
			if err != nil {
				tb.Error(err)
				return nil
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
			}
		},
		"glove-25": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(datasetDir(tb) + "glove-25-angular.hdf5")
			if err != nil {
				tb.Error(err)
				return nil
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
			}
		},
		"glove-50": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(datasetDir(tb) + "glove-50-angular.hdf5")
			if err != nil {
				tb.Error(err)
				return nil
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
			}
		},
		"glove-100": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(datasetDir(tb) + "glove-100-angular.hdf5")
			if err != nil {
				tb.Error(err)
				return nil
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
			}
		},
		"glove-200": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(datasetDir(tb) + "glove-200-angular.hdf5")
			if err != nil {
				tb.Error(err)
				return nil
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
			}
		},
		"nytimes": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(datasetDir(tb) + "nytimes-256-angular.hdf5")
			if err != nil {
				tb.Error(err)
				return nil
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
			}
		},
		"sift": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(datasetDir(tb) + "sift-128-euclidean.hdf5")
			if err != nil {
				tb.Error(err)
				return nil
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
			}
		},
		"gist": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(datasetDir(tb) + "gist-960-euclidean.hdf5")
			if err != nil {
				tb.Error(err)
				return nil
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
			}
		},
		"kosarak": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, distances, neighbors, dim, err := LoadDataWithSequentialIDs(datasetDir(tb) + "/kosarak-jaccard.hdf5")
			if err != nil {
				tb.Error(err)
				return nil
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
			}
		},
	}
)

func identity(dim int) func(tb testing.TB) Dataset {
	return func(tb testing.TB) Dataset {
		tb.Helper()
		ids := CreateSequentialIDs(dim)
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
		}
	}
}

func datasetDir(tb testing.TB) string {
	tb.Helper()
	wd, err := os.Getwd()
	if err != nil {
		tb.Error(err)
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
	return filepath.Join(root, "benchmark/assets/dataset") + "/"
}

func Data(name string) func(testing.TB) Dataset {
	if strings.HasPrefix(name, "identity-") {
		i, _ := strconv.Atoi(name[9:])
		return identity(i)
	}
	return data[name]
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
