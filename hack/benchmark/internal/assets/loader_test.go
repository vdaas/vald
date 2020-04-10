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
	"reflect"
	"testing"

	"gonum.org/v1/hdf5"
)

func Test_loadFloat32(t *testing.T) {
	type args struct {
		dset    *hdf5.Dataset
		npoints int
		row     int
		dim     int
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadFloat32(tt.args.dset, tt.args.npoints, tt.args.row, tt.args.dim)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFloat32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadInt(t *testing.T) {
	type args struct {
		dset    *hdf5.Dataset
		npoints int
		row     int
		dim     int
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadInt(tt.args.dset, tt.args.npoints, tt.args.row, tt.args.dim)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadDataset(t *testing.T) {
	type args struct {
		file *hdf5.File
		name string
		f    loaderFunc
	}
	tests := []struct {
		name    string
		args    args
		wantDim int
		wantVec interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDim, gotVec, err := loadDataset(tt.args.file, tt.args.name, tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadDataset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDim != tt.wantDim {
				t.Errorf("loadDataset() gotDim = %v, want %v", gotDim, tt.wantDim)
			}
			if !reflect.DeepEqual(gotVec, tt.wantVec) {
				t.Errorf("loadDataset() gotVec = %v, want %v", gotVec, tt.wantVec)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name          string
		args          args
		wantTrain     [][]float32
		wantTest      [][]float32
		wantDistances [][]float32
		wantNeighbors [][]int
		wantDim       int
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTrain, gotTest, gotDistances, gotNeighbors, gotDim, err := Load(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTrain, tt.wantTrain) {
				t.Errorf("Load() gotTrain = %v, want %v", gotTrain, tt.wantTrain)
			}
			if !reflect.DeepEqual(gotTest, tt.wantTest) {
				t.Errorf("Load() gotTest = %v, want %v", gotTest, tt.wantTest)
			}
			if !reflect.DeepEqual(gotDistances, tt.wantDistances) {
				t.Errorf("Load() gotDistances = %v, want %v", gotDistances, tt.wantDistances)
			}
			if !reflect.DeepEqual(gotNeighbors, tt.wantNeighbors) {
				t.Errorf("Load() gotNeighbors = %v, want %v", gotNeighbors, tt.wantNeighbors)
			}
			if gotDim != tt.wantDim {
				t.Errorf("Load() gotDim = %v, want %v", gotDim, tt.wantDim)
			}
		})
	}
}

func TestCreateRandomIDs(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateRandomIDs(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRandomIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateSequentialIDs(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateSequentialIDs(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSequentialIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadDataWithRandomIDs(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name          string
		args          args
		wantIds       []string
		wantTrain     [][]float32
		wantTest      [][]float32
		wantDistances [][]float32
		wantNeighbors [][]int
		wantDim       int
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIds, gotTrain, gotTest, gotDistances, gotNeighbors, gotDim, err := LoadDataWithRandomIDs(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadDataWithRandomIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIds, tt.wantIds) {
				t.Errorf("LoadDataWithRandomIDs() gotIds = %v, want %v", gotIds, tt.wantIds)
			}
			if !reflect.DeepEqual(gotTrain, tt.wantTrain) {
				t.Errorf("LoadDataWithRandomIDs() gotTrain = %v, want %v", gotTrain, tt.wantTrain)
			}
			if !reflect.DeepEqual(gotTest, tt.wantTest) {
				t.Errorf("LoadDataWithRandomIDs() gotTest = %v, want %v", gotTest, tt.wantTest)
			}
			if !reflect.DeepEqual(gotDistances, tt.wantDistances) {
				t.Errorf("LoadDataWithRandomIDs() gotDistances = %v, want %v", gotDistances, tt.wantDistances)
			}
			if !reflect.DeepEqual(gotNeighbors, tt.wantNeighbors) {
				t.Errorf("LoadDataWithRandomIDs() gotNeighbors = %v, want %v", gotNeighbors, tt.wantNeighbors)
			}
			if gotDim != tt.wantDim {
				t.Errorf("LoadDataWithRandomIDs() gotDim = %v, want %v", gotDim, tt.wantDim)
			}
		})
	}
}

func TestLoadDataWithSequentialIDs(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name          string
		args          args
		wantIds       []string
		wantTrain     [][]float32
		wantTest      [][]float32
		wantDistances [][]float32
		wantNeighbors [][]int
		wantDim       int
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIds, gotTrain, gotTest, gotDistances, gotNeighbors, gotDim, err := LoadDataWithSequentialIDs(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadDataWithSequentialIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIds, tt.wantIds) {
				t.Errorf("LoadDataWithSequentialIDs() gotIds = %v, want %v", gotIds, tt.wantIds)
			}
			if !reflect.DeepEqual(gotTrain, tt.wantTrain) {
				t.Errorf("LoadDataWithSequentialIDs() gotTrain = %v, want %v", gotTrain, tt.wantTrain)
			}
			if !reflect.DeepEqual(gotTest, tt.wantTest) {
				t.Errorf("LoadDataWithSequentialIDs() gotTest = %v, want %v", gotTest, tt.wantTest)
			}
			if !reflect.DeepEqual(gotDistances, tt.wantDistances) {
				t.Errorf("LoadDataWithSequentialIDs() gotDistances = %v, want %v", gotDistances, tt.wantDistances)
			}
			if !reflect.DeepEqual(gotNeighbors, tt.wantNeighbors) {
				t.Errorf("LoadDataWithSequentialIDs() gotNeighbors = %v, want %v", gotNeighbors, tt.wantNeighbors)
			}
			if gotDim != tt.wantDim {
				t.Errorf("LoadDataWithSequentialIDs() gotDim = %v, want %v", gotDim, tt.wantDim)
			}
		})
	}
}
