package service

import (
	"context"
	"os"
	"runtime/trace"
	"testing"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/internal/config"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/rand"
	"gonum.org/v1/hdf5"
)

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

	return ids, train, test, err
}

func BenchmarkGetObjectWithUpdate(b *testing.B) {
	log.Init(log.WithLoggerType(logger.NOP.String()))

	datasetPath := "fashion-mnist-784-euclidean.hdf5"
	insertCount := 10000
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
	n.CreateIndex(context.Background(), uint32(insertCount))

	f, err := os.Create("trace.out")
	if err != nil {
		b.Error(err)
	}
	defer f.Close()

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()

	trace.Start(f)
	defer trace.Stop()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			go func() {
				id := rand.LimitedUint32(uint64(insertCount))
				id2 := rand.LimitedUint32(uint64(insertCount))
				n.Update(ids[id], train[id2])
			}()
			_, _ = n.GetObject(ids[rand.LimitedUint32(uint64(insertCount))])
		}
	})

	b.StopTimer()
}

func BenchmarkSeachWithCreateIndex(b *testing.B) {
	// stdlog.Println("starting benchmark")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Init(log.WithLoggerType(logger.NOP.String()))

	datasetPath := "fashion-mnist-784-euclidean.hdf5"
	insertCount := 10000
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
	maxindex := len(ids)

	n, err := New(&defaultConfig)
	if err != nil {
		b.Error(err)
	}

	f, err := os.Create("trace.out")
	if err != nil {
		b.Error(err)
	}
	defer f.Close()

	b.StopTimer()
	b.ResetTimer()
	b.ReportAllocs()
	b.StartTimer()

	trace.Start(f)
	defer trace.Stop()

	b.RunParallel(func(pb *testing.PB) {
		// Keep creating and deleting index
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					// stdlog.Println("done")
					return
				default:
					var err2 error
					idxs := make([]uint32, 0, insertCount)

					for i := 0; i < insertCount; i++ {
						idx := rand.LimitedUint32(uint64(maxindex))
						idxs = append(idxs, idx)

						err2 = n.Insert(ids[idx], train[idx])
						if err2 != nil {
							if errors.Is(err2, errors.ErrUUIDAlreadyExists(ids[idx])) {
								continue
							}
							// stdlog.Println(err2, ": insert error")
						}
					}

					err2 = n.CreateIndex(ctx, uint32(insertCount))
					if err2 != nil {
						// stdlog.Println(err2, "create index error")
					}

					for _, idx := range idxs {
						err2 = n.Delete(ids[idx])
						if err2 != nil {
							if errors.Is(err2, errors.ErrObjectIDNotFound(ids[idx])) {
								continue
							}
							// stdlog.Println(err2, "delete error")
						}
					}
				}
			}
		}(ctx)

		// Bench serch
		for pb.Next() {
			_, err := n.Search(
				ctx,
				train[12345], // use this random vector as query
				uint32(defaultConfig.SearchEdgeSize),
				defaultConfig.DefaultEpsilon,
				defaultConfig.DefaultRadius,
			)
			if err != nil {
				if errors.Is(err, errors.ErrCreateIndexingIsInProgress) {
					// Check if it returns is_indexing error
					// stdlog.Println(err)
				}
			}
		}
	})

	b.StopTimer()
}
