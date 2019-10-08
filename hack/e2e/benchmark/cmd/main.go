//
// Copyright (C) 2019 kpango (Yusuke Kato)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package main

import (
	"context"
	"encoding/json"
	"os"
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
		log.Error(path, name)
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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "localhost:8082", grpc.WithInsecure())
	if err != nil {
		log.Error(err)
		return
	}
	client := agent.NewAgentClient(conn)

	train, err := load(datasetName, "train")
	if err != nil {
		log.Error(err)
		return
	}
	start := time.Now()
	log.Info("insert start")
	cstream, err := client.StreamInsert(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	for _, vector := range train {
		err = cstream.Send(&payload.Object_Vector{
			Id: &payload.Object_ID{
				Id: fuid.String(),
			},
			Vector: vector,
		})
		// _, err = client.Insert(ctx, &payload.Object_Vector{
		// 	Id: &payload.Object_ID{
		// 		Id: fuid.String(),
		// 	},
		// 	Vector: vector,
		// })
		if err != nil {
			log.Error(err)
		}
	}
	cstream.CloseSend()
	log.Info("insert finish", time.Now().Sub(start))

	start = time.Now()
	log.Info("indexing start")
	if _, err := client.CreateIndex(ctx, &payload.Controll_CreateIndexRequest{
		PoolSize: uint32(10000),
	}); err != nil {
		log.Error(err)
	}
	log.Info("indexing finish", time.Now().Sub(start))

	test, err := load(datasetName, "test")

	if err != nil {
		log.Error(err)
	}

	start = time.Now()
	log.Info("search start")
	all := make([]*payload.Object_Distance, 0, len(test)*10)
	for _, vector := range test {
		req := &payload.Search_Request{
			Vector: &payload.Object_Vector{
				Vector: vector,
			},
			Config: &payload.Search_Config{
				Num:     10,
				Radius:  -1.0,
				Epsilon: 0.01,
			},
		}
		res, err := client.Search(ctx, req)
		if err != nil {
			log.Error(err)
		}
		if res.GetResults() != nil {
			all = append(all, res.GetResults()...)
		}
	}
	b, _ := json.MarshalIndent(all, "", "\t")
	log.Info(string(b))
	log.Info("search finish", time.Now().Sub(start))
}
