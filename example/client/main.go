package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/vdaas/vald/apis/grpc/vald"

	"gonum.org/v1/hdf5"
	"google.golang.org/grpc"
)

var (
	datasetDir     string
	grpcServerAddr string
)

func init() {
	flag.StringVar(&datasetDir, "path", "/tmp/", "set dataset path")
	flag.StringVar(&grpcServerAddr, "addr", ":8081", "set gRPC server address")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	conn, err := grpc.DialContext(context.Background(), grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	client := vald.NewValdClient(conn)
	fmt.Println(client)

	return nil
}

func load(path string) (train, test [][]float64, err error) {
	var f *hdf5.File
	f, err = hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, nil, err
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
		return nil, nil, err
	}

	test, err = readFn("train")
	if err != nil {
		return nil, nil, err
	}

	return
}
