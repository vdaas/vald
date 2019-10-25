package dataset

import (
	"testing"
)

type Dataset interface {
	Train() [][]float64
	Query() [][]float64
	IDs() []string
	Name() string
	Dimension() int
	DistanceType() string
	ObjectType() string
}

type dataset struct {
	train [][]float64
	query [][]float64
	ids []string
	name string
	dimension int
	distanceType string
	objectType string
}

const (
	datasetDir = "../../assets/dataset/"
)

var (
	Data = map[string]func(testing.TB) Dataset {
		"fashion-mnist": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, dim, err := LoadDataAndIDs(datasetDir + "fashion-mnist-784-euclidean.hdf5")
			if err != nil {
				tb.Error(err)
			}
			return &dataset{
				train: train,
				query: query,
				ids: ids,
				name: "fashion-mnist",
				dimension: dim,
				distanceType: "l2",
				objectType: "float",
			}
		},
		"mnist": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, dim, err := LoadDataAndIDs(datasetDir + "mnist-784-euclidean.hdf5")
			if err != nil {
				tb.Error(err)
			}
			return &dataset{
				train: train,
				query: query,
				ids: ids,
				name: "mnist",
				dimension: dim,
				distanceType: "l2",
				objectType: "float",
			}
		},
		"glove-25": func(tb testing.TB) Dataset {
			tb.Helper()
			ids, train, query, dim, err := LoadDataAndIDs(datasetDir + "glove-25-angular.hdf5")
			if err != nil {
				tb.Error(err)
			}
			return &dataset{
				train: train,
				query: query,
				ids: ids,
				name: "glove-25",
				dimension: dim,
				distanceType: "cosine",
				objectType: "float",
			}
		},
		"glove-50": func(tb testing.TB) Dataset {
			tb.Helper()
			return nil
		},
		"glove-100": func(tb testing.TB) Dataset {
			tb.Helper()
			return nil
		},
		"glove-200": func(tb testing.TB) Dataset {
			tb.Helper()
			return nil
		},
		"nytimes": func(tb testing.TB) Dataset {
			tb.Helper()
			return nil
		},
		"sift": func(tb testing.TB) Dataset {
			tb.Helper()
			return nil
		},
		"gist": func(tb testing.TB) Dataset {
			tb.Helper()
			return nil
		},
	}
)

func (d *dataset) Train() [][]float64 {
	return d.train
}

func (d *dataset) Query() [][]float64 {
	return d.query
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
