package assets

import (
	"runtime"

	"github.com/vdaas/vald/internal/core/ngt"
	"gonum.org/v1/hdf5"
)

type Dataset struct {
	Path string
	f    *hdf5.File
}

func LoadDataset(path string) (*Dataset, error) {
	f, err := hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, err
	}
	return &Dataset{
		Path: path,
		f:    f,
	}, nil
}

func (d *Dataset) Close() error {
	return d.f.Close()
}

func (d *Dataset) load(name string) ([][]float64, error) {
	dset, err := d.f.OpenDataset(name)
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

func (d *Dataset) LoadTrain() ([][]float64, error) {
	return d.load("train")
}

func (d *Dataset) LoadTest() ([][]float64, error) {
	return d.load("test")
}

func CreateIndex(indexName string, loader func() ([][]float64, error), opts ...ngt.Option) error {
	vectors, err := loader()
	if err != nil {
		return err
	}
	opts = append(opts, ngt.WithIndexPath(indexName))
	opts = append(opts, ngt.WithDimension(len(vectors[0])))
	opts = append(opts, ngt.WithObjectType(ngt.Float))
	n, err := ngt.New(opts...)
	if err != nil {
		return err
	}
	defer n.Close()

	for _, v := range vectors {
		n.Insert(v)
	}
	return n.CreateAndSaveIndex(runtime.NumCPU())
}
