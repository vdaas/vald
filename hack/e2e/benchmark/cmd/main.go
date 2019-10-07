package main

import (
	"gonum.org/v1/hdf5"
)

func load(path, name string) ([][]float64, error) {
	f, err := hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dset, err := f.OpenDataset(name)
	if err != nil {
		return nil, err
	}
	defer dset.Close()
	space := dset.Space()
	defer space.Close()
	dims, _, err := space.SimpleExtentDims()
	if err != nil {
		return nil, err
	}
	v := make([]float32, space.SimpleExtentNPoints())
	if err := dset.Read(&v); err != nil {
		return nil, err
	}

	row := int(dims[0])
	col := int(dims[1])

	vec := make([][]float64, row)
	for i := 0; i < row; i++ {
		vec[i] = make([]float64, col)
		for j := 0; j < col; j++ {
			vec[i][j] = float64(v[i*col+j])
		}
	}
	return vec, nil
}

func main() {

}
