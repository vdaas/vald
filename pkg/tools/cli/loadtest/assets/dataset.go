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
	"os"
	"path/filepath"
	"strconv"

	"github.com/vdaas/vald/hack/benchmark/assets/x1b"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/strings"
)

var ErrOutOfBounds = x1b.ErrOutOfBounds

// Dataset is representation of train and test dataset.
type Dataset interface {
	Train(i int) (interface{}, error)
	TrainSize() int
	Query(i int) (interface{}, error)
	QuerySize() int
	Distance(i int) ([]float32, error)
	DistanceSize() int
	Neighbor(i int) ([]int, error)
	NeighborSize() int
	Name() string
	Dimension() int
	DistanceType() string
	ObjectType() string
}

type dataset struct {
	name         string
	dimension    int
	distanceType string
	objectType   string
}

// Name returns dataset name.
func (d *dataset) Name() string {
	return d.name
}

// Dimension returns vector dimension.
func (d *dataset) Dimension() int {
	return d.dimension
}

// DistanceType returns dataset distance type like l2, cosine, jaccard or etc.
func (d *dataset) DistanceType() string {
	return d.distanceType
}

// ObjectType returns dataset vector type like float or int.
func (d *dataset) ObjectType() string {
	return d.objectType
}

func findDir(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	root := func(cur string) string {
		for {
			if strings.HasSuffix(cur, "vald") {
				return cur
			} else {
				cur = filepath.Dir(cur)
			}
		}
	}(wd)
	return file.Join(root, path) + string(os.PathSeparator), nil
}

// Data loads specified dataset and returns it.
func Data(name string) func() (Dataset, error) {
	log.Debugf("start loading: %s", name)
	defer log.Debugf("finish loading: %s", name)
	if strings.HasPrefix(name, "identity-") {
		l := strings.Split(name, "-")
		i, _ := strconv.Atoi(l[1])
		return identity(i)
	}
	if strings.HasPrefix(name, "random-") {
		l := strings.Split(name, "-")
		d, _ := strconv.Atoi(l[1])
		s, _ := strconv.Atoi(l[2])
		return random(d, s)
	}
	if strings.HasPrefix(name, "gaussian-") {
		l := strings.Split(name, "-")
		d, _ := strconv.Atoi(l[1])
		s, _ := strconv.Atoi(l[2])
		m, _ := strconv.ParseFloat(l[3], 64)
		sd, _ := strconv.ParseFloat(l[4], 64)
		return gaussian(d, s, m, sd)
	}

	switch name {
	case "fashion-mnist":
		return loadSmallData("fashion-mnist-784-euclidean.hdf5", name, "l2", "float")
	case "mnist":
		return loadSmallData("mnist-784-euclidean.hdf5", name, "l2", "float")
	case "glove-25":
		return loadSmallData("glove-25-angular.hdf5", name, "cosine", "float")
	case "glove-50":
		return loadSmallData("glove-50-angular.hdf5", name, "cosine", "float")
	case "glove-100":
		return loadSmallData("glove-100-angular.hdf5", name, "cosine", "float")
	case "glove-200":
		return loadSmallData("glove-200-angular.hdf5", name, "cosine", "float")
	case "nytimes":
		return loadSmallData("nytimes-256-angular.hdf5", name, "cosine", "float")
	case "sift":
		return loadSmallData("sift-128-euclidean.hdf5", name, "l2", "float")
	case "gist":
		return loadSmallData("gist-960-euclidean.hdf5", name, "l2", "float")
	case "kosarak":
		return loadSmallData("kosarak-jaccard.hdf5", name, "jaccard", "float")
	case "sift1b":
		return loadLargeData("bigann_base.bvecs", "bigann_query.bvecs", "gnd/idx_1000M.ivecs", "gnd/dis_1000M.fvecs", name, "l2", "uint8")
	case "deep1b":
		return loadLargeData("deep1B_base.fvecs", "deep1B_query.fvecs", "deep1B_groundtruth.ivecs", "", name, "l2", "float")
	}
	return nil
}
