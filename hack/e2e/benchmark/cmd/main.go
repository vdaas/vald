package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/kpango/fuid"
	"github.com/vdaas/vald/apis/grpc/agent"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/log"
	"gonum.org/v1/hdf5"
	"google.golang.org/grpc"
)

func load(path, name string) (vec [][]float64, err error) {
	f, err := hdf5.OpenFile(path, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = f.Close()
	}()
	dset, err := f.OpenDataset(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = dset.Close()
	}()
	space := dset.Space()
	defer func() {
		err = space.Close()
	}()
	dims, _, err := space.SimpleExtentDims()
	if err != nil {
		return nil, err
	}
	v := make([]float32, space.SimpleExtentNPoints())
	if err := dset.Read(&v); err != nil {
		return nil, err
	}

	row, col := int(dims[0]), int(dims[1])

	vec = make([][]float64, row)
	for i := 0; i < row; i++ {
		vec[i] = make([]float64, col)
		for j := 0; j < col; j++ {
			vec[i][j] = float64(v[i*col+j])
		}
	}
	return vec, nil
}

func main() {
	log.Init(log.DefaultGlg())

	datasetName := os.Args[1]

	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Hour)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "localhost:8082", grpc.WithInsecure())
	if err != nil {
		log.Error(err)
	}
	client := agent.NewAgentClient(conn)

	train, err := load(fmt.Sprintf("assets/%s.hdf5", datasetName), "train")
	if err != nil {
		log.Error(err)
	}
	log.Info("insert start")
	for _, vector := range train {
		if _, err := client.Insert(ctx, &payload.Object_Vector{
			Id: &payload.Object_ID{
				Id: fuid.String(),
			},
			Vector:vector,
		}); err != nil {
			log.Error(err)
		}
	}
	log.Info("insert finish")

	log.Info("indexing start")
	if _, err := client.CreateIndex(ctx, &payload.Controll_CreateIndexRequest{
		PoolSize: uint32(runtime.NumCPU()),
	}); err != nil {
		log.Error(err)
	}
	log.Info("indexing finish")

	test, err := load(fmt.Sprintf("assets/%s.hdf5", datasetName), "test")
	if err != nil {
		log.Error(err)
	}

	log.Info("search start")
	for _, vector := range test {
		req := &payload.Search_Request{
			Vector: &payload.Object_Vector{
				Vector: vector,
			},
			Config: &payload.Search_Config{
				Num: 10,
				Epsilon: 0.1,
			},
		}
		res, err := client.Search(ctx, req)
		if err != nil {
			log.Error(err)
		}
		log.Info(res.GetResults())
	}
	log.Info("search finish")
}
