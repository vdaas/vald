package assets

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vdaas/vald/internal/log"
)

type smallDataset struct {
	*dataset
	train     [][]float32
	query     [][]float32
	distances [][]float32
	neighbors [][]int
}

func loadSmallData(fileName, datasetName, distanceType, objectType string) func() (Dataset, error) {
	return func() (Dataset, error) {
		dir, err := findDir("hack/benchmark/assets/dataset")
		if err != nil {
			return nil, err
		}
		t, q, d, n, dim, err := Load(filepath.Join(dir, fileName))
		if err != nil {
			return nil, err
		}

		return &smallDataset{
			dataset: &dataset {
				name:         datasetName,
				dimension:    dim,
				distanceType: distanceType,
				objectType:   objectType,
			},
			train:     t,
			query:     q,
			distances: d,
			neighbors: n,
		}, nil
	}
}

func identity(dim int) func() (Dataset, error) {
	return func() (Dataset, error) {
		train := make([][]float32, dim)
		for i := range train {
			train[i] = make([]float32, dim)
			train[i][i] = 1
		}
		return &smallDataset{
			dataset: &dataset{
				name:         fmt.Sprintf("identity-%d", dim),
				dimension:    dim,
				distanceType: "l2",
				objectType:   "float",
			},
			train: train,
			query: train,
		}, nil
	}
}

func random(dim, size int) func() (Dataset, error) {
	return func() (Dataset, error) {
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
		return &smallDataset{
			dataset: &dataset{
				name:         fmt.Sprintf("random-%d-%d", dim, size),
				dimension:    dim,
				distanceType: "l2",
				objectType:   "float",
			},
			train: train,
			query: query,
		}, nil
	}
}

func gaussian(dim, size int, mean, stdDev float64) func() (Dataset, error) {
	return func() (Dataset, error) {
		train := make([][]float32, size)
		query := make([][]float32, size)
		for i := range train {
			train[i] = make([]float32, dim)
			query[i] = make([]float32, dim)
			for j := range train[i] {
				train[i][j] = float32(rand.NormFloat64()*stdDev + mean)
				query[i][j] = float32(rand.NormFloat64()*stdDev + mean)
			}
		}
		return &smallDataset{
			dataset: &dataset{
				name:         fmt.Sprintf("gaussian-%d-%d-%f-%f", dim, size, mean, stdDev),
				dimension:    dim,
				distanceType: "l2",
				objectType:   "float",
			},
			train: train,
			query: query,
		}, nil
	}
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

// Train returns vectors for train.
func (s *smallDataset) Train(i int) (interface{}, error) {
	if i >= len(s.train) {
		return nil, ErrOutOfBounds
	}
	return s.train[i], nil
}

// Query returns vectors for test.
func (s *smallDataset) Query(i int) (interface{}, error) {
	if i >= len(s.query) {
		return nil, ErrOutOfBounds
	}
	return s.query[i], nil
}

// Distance returns distances between queries and answers.
func (s *smallDataset) Distance(i int) ([]float32, error) {
	if i >= len(s.distances) {
		return nil, ErrOutOfBounds
	}
	return s.distances[i], nil
}

// Neighbors returns nearest vectors from queries.
func (s *smallDataset) Neighbor(i int) ([]int, error) {
	if i >= len(s.neighbors) {
		return nil, ErrOutOfBounds
	}
	return s.neighbors[i], nil
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

