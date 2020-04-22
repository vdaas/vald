package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/apis/grpc/vald"

	"gonum.org/v1/hdf5"
	"google.golang.org/grpc"
)

var (
	datasetPath    string
	grpcServerAddr string
)

var searchConfig = payload.Search_Config{
	Num:     10,
	Radius:  -1,
	Epsilon: 0.01,
}

func init() {
	flag.StringVar(&datasetPath, "path", "fashion-mnist-784-euclidean.hdf5", "set dataset path")
	flag.StringVar(&grpcServerAddr, "addr", "TODO", "set gRPC server address")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ids, train, _, err := load(datasetPath)
	if err != nil {
		return err
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := vald.NewValdClient(conn)

	errCh := make(chan error, len(ids)*10)
	wg := new(sync.WaitGroup)

	limCh := make(chan struct{}, 10000)

	for i := 0; i < 100; i++ {
		fmt.Printf("%d - \n", i)
		for j := range ids {
			wg.Add(1)
			limCh <- struct{}{}
			go func(j int, id string, vec []float64) {
				defer wg.Done()
				defer func() {
					<-limCh
				}()

				if j%100 == 0 {
					fmt.Printf("%d end \n", j)
				}

				_, err := client.Insert(ctx, &payload.Object_Vector{
					Id:     id,
					Vector: vec,
				})
				if err != nil {
					fmt.Println(err)
					errCh <- err
				}
			}(j, ids[j], train[j%len(train)])
		}
	}

	fmt.Println("WAIT..")
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func load(path string) (ids []string, train, test [][]float64, err error) {
	var f *hdf5.File
	f, err = hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	readFn := func(name string) ([][]float64, error) {
		d, err := f.OpenDataset(name)
		if err != nil {
			return nil, err
		}
		defer d.Close()

		sp := d.Space()
		defer sp.Close()

		dims, _, _ := sp.SimpleExtentDims()
		row, dim := int(dims[0]), int(dims[1])

		vec := make([]float64, sp.SimpleExtentNPoints())
		if err := d.Read(&vec); err != nil {
			return nil, err
		}

		vecs := make([][]float64, row)
		for i := 0; i < row; i++ {
			vecs[i] = make([]float64, dim)
			for j := 0; j < dim; j++ {
				vecs[i][j] = vec[i*dim+j]
			}
		}

		return vecs, nil
	}

	train, err = readFn("train")
	if err != nil {
		return nil, nil, nil, err
	}

	test, err = readFn("train")
	if err != nil {
		return nil, nil, nil, err
	}

	ids = make([]string, 0, len(train)*150)
	for i := 0; i < len(train)*150; i++ {
		ids = append(ids, fuid.String())
	}

	return
}
