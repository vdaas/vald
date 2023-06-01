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

// Package hdf5 is load hdf5 file
package hdf5

import (
	"os"
	"reflect"

	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/io"
	"github.com/vdaas/vald/internal/net/http/client"
	"gonum.org/v1/hdf5"
)

type Data interface {
	Download(url string) error
	Read() error
	GetName() DatasetName
	GetPath() string
	GetByGroupName(name string) [][]float32
	GetTrain() [][]float32
	GetTest() [][]float32
	GetNeighbors() [][]int
}

type DatasetName int

const (
	Original DatasetName = iota
	FashionMNIST784Euclidean
)

func (d DatasetName) String() string {
	switch d {
	case Original:
		return "original"
	case FashionMNIST784Euclidean:
		return "fashion-mnist-784-euc"
	default:
		return ""
	}
}

type DatasetUrl int

const (
	FashionMNIST784EuclideanUrl DatasetUrl = iota
)

func (d DatasetUrl) String() string {
	switch d {
	case FashionMNIST784EuclideanUrl:
		return "http://ann-benchmarks.com/fashion-mnist-784-euclidean.hdf5"
	default:
		return ""
	}
}

type hdf5Key int

const (
	Train hdf5Key = iota
	Test
	Neighors
)

func (h hdf5Key) String() string {
	switch h {
	case Train:
		return "train"
	case Test:
		return "test"
	case Neighors:
		return "neighbors"
	default:
		return ""
	}
}

type data struct {
	name      DatasetName
	path      string
	train     [][]float32
	test      [][]float32
	neighbors [][]int
}

func New(opts ...Option) (Data, error) {
	d := new(data)
	for _, opt := range append(defaultOptions, opts...) {
		if err := opt(d); err != nil {
			return nil, errors.ErrOptionFailed(err, reflect.ValueOf(opt))
		}
	}
	return d, nil
}

// Get downloads the hdf5 file.
// https://github.com/erikbern/ann-benchmarks/#data-sets
func (d *data) Download(url string) error {
	switch d.name {
	case Original:
		return downloadFile(url, d.path)
	case FashionMNIST784Euclidean:
		return downloadFile(FashionMNIST784EuclideanUrl.String(), d.path)
	default:
		return errors.NewErrInvalidOption("name", d.name)
	}
}

func (d *data) Read() error {
	f, err := hdf5.OpenFile(d.path, hdf5.F_ACC_RDONLY)
	if err != nil {
		return err
	}
	defer f.Close()

	// load training data
	train, err := ReadDatasetF32(f, Train)
	if err != nil {
		return err
	}
	d.train = train

	// load test data
	test, err := ReadDatasetF32(f, Test)
	if err != nil {
		return err
	}
	d.test = test

	// load neighbors
	neighbors32, err := ReadDatasetI32(f, Neighors)
	if err != nil {
		return err
	}
	neighbors := make([][]int, len(neighbors32))
	for i, ns := range neighbors32 {
		neighbors[i] = make([]int, len(ns))
		for j, n := range ns {
			neighbors[i][j] = int(n)
		}
	}
	d.neighbors = neighbors

	return nil
}

func (d *data) GetName() DatasetName {
	return d.name
}

func (d *data) GetPath() string {
	return d.path
}

// TODO: Apply generics
func (d *data) GetByGroupName(name string) [][]float32 {
	switch name {
	case "train":
		return d.GetTrain()
	case "test":
		return d.GetTest()
	case "neighbors":
		l := d.GetNeighbors()
		r := make([][]float32, 0)
		for x := range l {
			for y, z := range l[x] {
				r[x][y] = float32(z)
			}
		}
		return r
	default:
		return nil
	}
}

func (d *data) GetTrain() [][]float32 {
	return d.train
}

func (d *data) GetTest() [][]float32 {
	return d.test
}

func (d *data) GetNeighbors() [][]int {
	return d.neighbors
}

func downloadFile(url, path string) error {
	if len(path) == 0 {
		return errors.NewErrInvalidOption("no path is specified", path)
	}
	cli, err := client.New()
	if err != nil {
		return err
	}

	resp, err := cli.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.ErrInvalidStatusCode(resp.StatusCode)
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func ReadDatasetF32(file *hdf5.File, key hdf5Key) ([][]float32, error) {
	data, err := file.OpenDataset(key.String())
	if err != nil {
		return nil, err
	}
	defer data.Close()

	dataspace := data.Space()
	defer dataspace.Close()

	dims, _, err := dataspace.SimpleExtentDims()
	if err != nil {
		return nil, err
	}
	height, width := int(dims[0]), int(dims[1])

	rawFloats := make([]float32, dataspace.SimpleExtentNPoints())
	if err := data.Read(&rawFloats); err != nil {
		return nil, err
	}

	vecs := make([][]float32, height)
	for i := 0; i < height; i++ {
		vecs[i] = rawFloats[i*width : i*width+width]
	}

	return vecs, nil
}

func ReadDatasetI32(file *hdf5.File, key hdf5Key) ([][]int32, error) {
	data, err := file.OpenDataset(key.String())
	if err != nil {
		return nil, err
	}
	defer data.Close()

	dataspace := data.Space()
	defer dataspace.Close()

	dims, _, err := dataspace.SimpleExtentDims()
	if err != nil {
		return nil, err
	}
	height, width := int(dims[0]), int(dims[1])

	rawFloats := make([]int32, dataspace.SimpleExtentNPoints())
	if err := data.Read(&rawFloats); err != nil {
		return nil, err
	}

	vecs := make([][]int32, height)
	for i := 0; i < height; i++ {
		vecs[i] = rawFloats[i*width : i*width+width]
	}

	return vecs, nil
}
