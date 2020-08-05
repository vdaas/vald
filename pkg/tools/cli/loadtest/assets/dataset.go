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
	"os"
	"path/filepath"
	"strings"

	"github.com/vdaas/vald/hack/benchmark/assets/x1b"
)

var (
	ErrOutOfBounds = x1b.ErrOutOfBounds
)

// Dataset is representation of train and test dataset.
type Dataset interface {
	Train(i int) (interface{}, error)
	Query(i int) (interface{}, error)
	Distance(i int) ([]float32, error)
	Neighbor(i int) ([]int, error)
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
	return filepath.Join(root, path) + "/", nil
}