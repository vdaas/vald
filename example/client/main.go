package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald-client-go/gateway/vald"
	"github.com/vdaas/vald-client-go/payload"

	"github.com/cheggaaa/pb/v3"
	"gonum.org/v1/hdf5"
	"google.golang.org/grpc"
)

const (
	insertCont = 400
	testCount  = 20
)

var (
	datasetPath         string
	grpcServerAddr      string
	indexingWaitSeconds uint
)

var searchConfig = payload.Search_Config{
	Num:     10,
	Radius:  -1,
	Epsilon: 0.01,
}

func init() {
	flag.StringVar(&datasetPath, "path", "fashion-mnist-784-euclidean.hdf5", "set dataset path")
	flag.StringVar(&grpcServerAddr, "addr", "127.0.0.1:8081", "set gRPC server address")
	flag.UintVar(&indexingWaitSeconds, "wait", 60, "set indexing wait seconds")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	ids, train, test, err := load(datasetPath)
	if err != nil {
		return err
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := vald.NewValdClient(conn)

	bar := pb.StartNew(insertCont)
	fmt.Println("Start Inserting: ")

	for i := range ids[:insertCont] {
		bar.Increment()

		_, err := client.Insert(ctx, &payload.Object_Vector{
			Id:     ids[i],
			Vector: train[i],
		})
		if err != nil {
			return err
		}
	}

	bar.Finish()
	fmt.Printf("Finish Inserting. \n\n")
	fmt.Println("Wait for indexing to finish")
	time.Sleep(time.Duration(indexingWaitSeconds) * time.Second)

	for _, vec := range test[:testCount] {
		res, err := client.Search(ctx, &payload.Search_Request{
			Vector: vec,
			Config: &searchConfig,
		})
		if err != nil {
			return err
		}

		b, _ := json.MarshalIndent(res.GetResults(), "", " ")
		fmt.Printf("results : %v\n\n", string(b))
	}

	return nil
}

func load(path string) (ids []string, train, test [][]float32, err error) {
	var f *hdf5.File
	f, err = hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, nil, nil, err
	}
	defer f.Close()

	readFn := func(name string) ([][]float32, error) {
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

		vecs := make([][]float32, row)
		for i := 0; i < row; i++ {
			vecs[i] = make([]float32, dim)
			for j := 0; j < dim; j++ {
				vecs[i][j] = float32(vec[i*dim+j])
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

	ids = make([]string, 0, len(train))
	for i := 0; i < len(train); i++ {
		ids = append(ids, fuid.String())
	}

	return
}
