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
)

func Test_identity(t *testing.T) {
	type args struct {
		dim int
	}
	tests := []struct {
		name string
		args args
		want func(tb testing.TB) Dataset
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := identity(tt.args.dim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("identity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_datasetDir(t *testing.T) {
	type args struct {
		tb testing.TB
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := datasetDir(tt.args.tb); got != tt.want {
				t.Errorf("datasetDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want func(testing.TB) Dataset
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Data(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Data() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_Train(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want [][]float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Train(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataset.Train() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_TrainAsFloat64(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want [][]float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.TrainAsFloat64(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataset.TrainAsFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_Query(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want [][]float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Query(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataset.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_QueryAsFloat64(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want [][]float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.QueryAsFloat64(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataset.QueryAsFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_Distances(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want [][]float32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Distances(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataset.Distances() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_DistancesAsFloat64(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want [][]float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.DistancesAsFloat64(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataset.DistancesAsFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_Neighbors(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want [][]int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Neighbors(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataset.Neighbors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_IDs(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.IDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataset.IDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_Name(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Name(); got != tt.want {
				t.Errorf("dataset.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_Dimension(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Dimension(); got != tt.want {
				t.Errorf("dataset.Dimension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_DistanceType(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.DistanceType(); got != tt.want {
				t.Errorf("dataset.DistanceType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dataset_ObjectType(t *testing.T) {
	tests := []struct {
		name string
		d    *dataset
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.ObjectType(); got != tt.want {
				t.Errorf("dataset.ObjectType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_float32To64(t *testing.T) {
	type args struct {
		x [][]float32
	}
	tests := []struct {
		name  string
		args  args
		wantY [][]float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotY := float32To64(tt.args.x); !reflect.DeepEqual(gotY, tt.wantY) {
				t.Errorf("float32To64() = %v, want %v", gotY, tt.wantY)
			}
		})
	}
}
