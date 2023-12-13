// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"math"
	"time"

	"github.com/kpango/fuid"
	"github.com/kpango/glg"
	agent "github.com/vdaas/vald-client-go/v1/agent/core"
	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"
	"gonum.org/v1/hdf5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	insertCount = 400
	removeCount = 200
	testCount   = 20
)

var (
	datasetPath         string
	grpcServerAddr      string
	indexingWaitSeconds uint
)

func init() {
	/**
	Path option specifies hdf file by path. Default value is `fashion-mnist-784-euclidean.hdf5`.
	Addr option specifies grpc server address. Default value is `127.0.0.1:8080`.
	Wait option specifies indexing wait time (in seconds). Default value is  `60`.
	**/
	flag.StringVar(&datasetPath, "path", "fashion-mnist-784-euclidean.hdf5", "dataset path")
	flag.StringVar(&grpcServerAddr, "addr", "127.0.0.1:8081", "gRPC server address")
	flag.UintVar(&indexingWaitSeconds, "wait", 60, "indexing wait seconds")
	flag.Parse()
}

func main() {
	/**
	Gets training data, test data and ids based on the dataset path.
	the number of ids is equal to that of training dataset.
	**/
	ids, train, test, err := load(datasetPath)
	if err != nil {
		glg.Fatal(err)
	}

	ctx := context.Background()

	// Create a Vald Agent client for connecting to the Vald cluster.
	conn, err := grpc.DialContext(ctx, grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glg.Fatal(err)
	}

	// Creates Vald Agent client for gRPC.
	client := vald.NewValdClient(conn)

	glg.Infof("Start Inserting %d training vector to Vald Agent", insertCount)
	// Insert 400 example vectors into Vald cluster
	for i := range ids[:insertCount] {
		// Calls `Insert` function of Vald Agent client.
		// Sends set of vector and id to server via gRPC.
		_, err := client.Insert(ctx, &payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:     ids[i],
				Vector: train[i],
			},
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: true,
			},
		})
		if err != nil {
			glg.Fatal(err)
		}
		if i%10 == 0 {
			glg.Infof("Inserted: %d", i+10)
		}
	}
	glg.Infof("Finish Inserting %d training vector to Vald Agent", insertCount)
	/**
	Optional: Run Indexing instead of Auto Indexing
	If you run client.CreateAndSaveIndex, it costs less time for search
	When default_pool_size is not set, the below codes are required.
	**/
	glg.Info("Start Indexing dataset.")
	_, err = agent.NewAgentClient(conn).CreateAndSaveIndex(ctx, &payload.Control_CreateIndexRequest{
		PoolSize: uint32(insertCount),
	})
	if err != nil {
		glg.Fatal(err)
	}
	glg.Info("Finish Indexing dataset. \n\n")

	// Vald Agent starts indexing automatically after insert. It needs to wait until the indexing is completed before a search action is performed.
	wt := time.Duration(indexingWaitSeconds) * time.Second
	glg.Infof("Wait %s for indexing to finish", wt)
	time.Sleep(wt)

	glg.Infof("Start getting object")
	for i := range ids[:insertCount] {
		// Call `GetObject` function of Vald client.
		// Sends id to server via gRPC.
		res, err := client.GetObject(ctx, &payload.Object_VectorRequest{
			Id: &payload.Object_ID{
				Id: ids[i],
			},
		})
		if err != nil {
			glg.Fatal(err)
		}
		glg.Infof("ID: %s, Vector: %v", res.GetId(), res.GetVector())

		// calc Euclidean distance of r and t
		r := res.GetVector()
		t := train[i]
		var sum float64
		for i := range r {
			sum += math.Pow(float64(t[i]-r[i]), 2)
		}
		glg.Infof("Euclidean distance of r and t: %v", sum)
	}
	glg.Info("Finish getting object")

	/**
	Gets approximate vectors, which is based on the value of `SearchConfig`, from the indexed tree based on the training data.
	In this example, Vald Agent gets 10 approximate vectors each search vector.
	**/
	glg.Infof("Start searching %d times", testCount)
	for i, vec := range test[:testCount] {
		// Send searching vector and configuration object to the Vald Agent server via gRPC.
		res, err := client.Search(ctx, &payload.Search_Request{
			Vector: vec,
			// Conditions for hitting the search.
			Config: &payload.Search_Config{
				Num:     10,        // the number of search results
				Radius:  -1,        // Radius is used to determine the space of search candidate radius for neighborhood vectors. -1 means infinite circle.
				Epsilon: 0.1,       // Epsilon is used to determines how much to expand from search candidate radius.
				Timeout: 100000000, // Timeout is used for search time deadline. The unit is nano-seconds.
			},
		})
		if err != nil {
			glg.Fatal(err)
		}

		b, _ := json.MarshalIndent(res.GetResults(), "", " ")
		glg.Infof("%d - Results : %s\n\n", i+1, string(b))
		time.Sleep(1 * time.Second)
	}

	glg.Info("Start removing vector")
	// Remove indexed 200 vectors from vald cluster.
	for i := range ids[:removeCount] {
		// Call `Remove` function of Vald client.
		// Sends id to server via gRPC.
		_, err := client.Remove(ctx, &payload.Remove_Request{
			// Conditions for removing the vector.
			Id: &payload.Object_ID{
				Id: ids[i],
			},
		})
		if err != nil {
			glg.Fatal(err)
		}
		if i%10 == 0 {
			glg.Infof("Removed: %d", i+10)
		}
	}
	glg.Info("Finish removing vector")
	glg.Info("Start removing indexed vector from backup")
	/**
	Run Indexing instead of Auto Indexing.
	Before calling the SaveIndex (or CreateAndSaveIndex) API, the vectors you inserted before still exist in the NGT graph index even the Remove API is called due to the design of NGT.
	So at this moment, the neighbor vectors will be returned from the SearchByID API.
	To remove the vectors from the NGT graph completely, the SaveIndex API will be used here instead of waiting auto CreateIndex phase.
	**/
	_, err = agent.NewAgentClient(conn).SaveIndex(ctx, &payload.Empty{})
	if err != nil {
		glg.Fatal(err)
	}
	glg.Info("Finish removing indexed vector from backup")
	glg.Info("Start flushing vector")
	_, err = client.Flush(ctx, &payload.Flush_Request{})
	if err != nil {
		glg.Fatal(err)
	}
	glg.Info("Finish flushing vector")
}

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
