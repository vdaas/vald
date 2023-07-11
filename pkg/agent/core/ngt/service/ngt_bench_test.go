package service

import (
	"testing"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/internal/config"
	"gonum.org/v1/hdf5"
)

// func init() {
// 	/**
// 	Path option specifies hdf file by path. Default value is `fashion-mnist-784-euclidean.hdf5`.
// 	Addr option specifies grpc server address. Default value is `127.0.0.1:8080`.
// 	Wait option specifies indexing wait time (in seconds). Default value is  `60`.
// 	**/
// 	flag.StringVar(&datasetPath, "path", "fashion-mnist-784-euclidean.hdf5", "dataset path")
// 	flag.StringVar(&grpcServerAddr, "addr", "127.0.0.1:8081", "gRPC server address")
// 	flag.UintVar(&indexingWaitSeconds, "wait", 60, "indexing wait seconds")
// 	flag.Parse()
// }

// load function loads training and test vector from hdf file. The size of ids is same to the number of training data.
// Each id, which is an element of ids, will be set a random number.
func load(path string) (ids []string, train, test [][]float32, err error) {
	var f *hdf5.File
	f, err = hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	// readFn function reads vectors of the hierarchy with the given the name.
	readFn := func(name string) ([][]float32, error) {
		// Opens and returns a named Dataset.
		// The returned dataset must be closed by the user when it is no longer needed.
		d, err := f.OpenDataset(name)
		if err != nil {
			return nil, err
		}
		defer d.Close()

		// Space returns an identifier for a copy of the dataspace for a dataset.
		sp := d.Space()
		defer sp.Close()

		// SimpleExtentDims returns dataspace dimension size and maximum size.
		dims, _, _ := sp.SimpleExtentDims()
		row, dim := int(dims[0]), int(dims[1])

		// Gets the stored vector. All are represented as one-dimensional arrays.
		// The type of the slice depends on your dataset.
		// For fashion-mnist-784-euclidean.hdf5, the datatype is float32.
		vec := make([]float32, sp.SimpleExtentNPoints())
		if err := d.Read(&vec); err != nil {
			return nil, err
		}

		// Converts a one-dimensional array to a two-dimensional array.
		// Use the `dim` variable as a separator.
		vecs := make([][]float32, row)
		for i := 0; i < row; i++ {
			vecs[i] = make([]float32, dim)
			for j := 0; j < dim; j++ {
				vecs[i][j] = float32(vec[i*dim+j])
			}
		}

		return vecs, nil
	}

	// Gets vector of `train` hierarchy.
	train, err = readFn("train")
	if err != nil {
		return nil, nil, nil, err
	}

	// Gets vector of `test` hierarchy.
	test, err = readFn("test")
	if err != nil {
		return nil, nil, nil, err
	}

	// Generate as many random ids for training vectors.
	ids = make([]string, 0, len(train))
	for i := 0; i < len(train); i++ {
		ids = append(ids, fuid.String())
	}

	return
}



func BenchmarkGetObject(b *testing.B) {
	datasetPath := "fashion-mnist-784-euclidean.hdf5"
	insertCount := 40000
	defaultConfig := config.NGT{
		Dimension:           784,
		DistanceType:        "l2",
		ObjectType:          "float",
		BulkInsertChunkSize: 10,
		CreationEdgeSize:    20,
		SearchEdgeSize:      10,
		EnableProactiveGC:   false,
		EnableCopyOnWrite:   false,
		KVSDB: &config.KVSDB{
			Concurrency: 10,
		},
		BrokenIndexHistoryLimit: 1,
	}

	/**
	Gets training data, test data and ids based on the dataset path.
	the number of ids is equal to that of training dataset.
	**/
	ids, train, _, err := load(datasetPath)
	if err != nil {
		b.Error(err)
	}

	n, err := New(&defaultConfig)
	if err != nil {
		b.Error(err)
	}

	for i := range ids[:insertCount] {
		n.Insert(ids[i], train[i])
	}


	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		n.GetObject(ids[i%insertCount])
	}
	b.StopTimer()
}


func BenchmarkGetObjectWithUpdate(b *testing.B) {
	datasetPath := "fashion-mnist-784-euclidean.hdf5"
	insertCount := 40000
	updateCount := 20000
	defaultConfig := config.NGT{
		Dimension:           784,
		DistanceType:        "l2",
		ObjectType:          "float",
		BulkInsertChunkSize: 10,
		CreationEdgeSize:    20,
		SearchEdgeSize:      10,
		EnableProactiveGC:   false,
		EnableCopyOnWrite:   false,
		KVSDB: &config.KVSDB{
			Concurrency: 10,
		},
		BrokenIndexHistoryLimit: 1,
	}

	/**
	Gets training data, test data and ids based on the dataset path.
	the number of ids is equal to that of training dataset.
	**/
	ids, train, _, err := load(datasetPath)
	if err != nil {
		b.Error(err)
	}

	n, err := New(&defaultConfig)
	if err != nil {
		b.Error(err)
	}

	for i := range ids[:insertCount] {
		n.Insert(ids[i], train[i])
	}

	// update objects
	for i := range ids[:updateCount] {
		n.Update(ids[i], train[i + 1])
	}


	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		n.GetObject(ids[i%insertCount])
	}
	b.StopTimer()
}
