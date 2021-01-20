// +build e2e

//
// Copyright (C) 2019-2021 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// package e2e provides e2e tests using ann-benchmarks datasets
package e2e

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/v1/agent/core"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"github.com/vdaas/vald/tests/e2e/portforward"

	"gonum.org/v1/hdf5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	host string
	port int
	ds   *dataset

	insertNum     int
	searchNum     int
	searchByIDNum int
	getObjectNum  int
	updateNum     int
	removeNum     int

	forwarder *portforward.Portforward
)

func init() {
	testing.Init()

	flag.StringVar(&host, "host", "localhost", "hostname")
	flag.IntVar(&port, "port", 8081, "gRPC port")

	flag.IntVar(&insertNum, "insert-num", 10000, "number of id-vector pairs used for insert")
	flag.IntVar(&searchNum, "search-num", 10000, "number of id-vector pairs used for search")

	datasetName := flag.String("dataset", "fashion-mnist-784-euclidean.hdf5", "dataset")

	pf := flag.Bool("portforward", false, "enable port forwarding")
	pfNamespace := flag.String("portforward-ns", "default", "namespace (only for port forward)")
	pfPodName := flag.String("portforward-pod-name", "vald-gateway-0", "pod name (only for port forward)")
	pfPodPort := flag.Int("portforward-pod-port", port, "pod gRPC port (only for port forward)")
	kubeConfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "kubeconfig path (only for port forward)")

	flag.Parse()

	var err error
	if *pf {
		forwarder, err = portforward.NewPortforward(*kubeConfig, *pfNamespace, *pfPodName, port, *pfPodPort)
		if err != nil {
			panic(err)
		}

		err = forwarder.Start()
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("loading dataset: %s", *datasetName)
	ds, err = hdf5ToDataset(*datasetName)
	if err != nil {
		panic(err)
	}
	fmt.Println("loading finished")
}

func teardown() {
	if forwarder != nil {
		forwarder.Close()
	}
}

func TestMain(m *testing.M) {
	ret := m.Run()
	teardown()
	os.Exit(ret)
}

type dataset struct {
	train     map[string][]float32
	test      map[string][]float32
	neighbors map[string][]string
}

func hdf5ToDataset(name string) (*dataset, error) {
	file, err := hdf5.OpenFile(name, hdf5.F_ACC_RDONLY)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	train, err := readDatasetF32(file, "train")
	if err != nil {
		return nil, err
	}

	test, err := readDatasetF32(file, "test")
	if err != nil {
		return nil, err
	}

	nbors, err := readDatasetI32(file, "neighbors")
	if err != nil {
		return nil, err
	}
	neighbors := make(map[string][]string, len(nbors))
	for k, vs := range nbors {
		vss := make([]string, len(vs))
		for i, v := range vs {
			vss[i] = strconv.Itoa(int(v))
		}
		neighbors[k] = vss
	}

	return &dataset{
		train:     train,
		test:      test,
		neighbors: neighbors,
	}, nil
}

func readDatasetF32(file *hdf5.File, name string) (map[string][]float32, error) {
	data, err := file.OpenDataset(name)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	dataspace := data.Space()
	defer dataspace.Close()

	dims, _, err := dataspace.SimpleExtentDims()
	if err != nil {
		return nil, err
	}
	height, width := int(dims[0]), int(dims[1])

	rawFloats := make([]float32, dataspace.SimpleExtentNPoints())
	if err := data.Read(&rawFloats); err != nil {
		return nil, err
	}

	vecs := make(map[string][]float32, height)
	for i := 0; i < height; i++ {
		vecs[strconv.Itoa(i)] = rawFloats[i*width : i*width+width]
	}

	return vecs, nil
}

func readDatasetI32(file *hdf5.File, name string) (map[string][]int32, error) {
	data, err := file.OpenDataset(name)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	dataspace := data.Space()
	defer dataspace.Close()

	dims, _, err := dataspace.SimpleExtentDims()
	if err != nil {
		return nil, err
	}
	height, width := int(dims[0]), int(dims[1])

	rawFloats := make([]int32, dataspace.SimpleExtentNPoints())
	if err := data.Read(&rawFloats); err != nil {
		return nil, err
	}

	vecs := make(map[string][]int32, height)
	for i := 0; i < height; i++ {
		vecs[strconv.Itoa(i)] = rawFloats[i*width : i*width+width]
	}

	return vecs, nil
}

func getAgentClient(ctx context.Context) (core.AgentClient, error) {
	conn, err := grpc.DialContext(
		ctx,
		host+":"+strconv.Itoa(port),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                time.Second,
				Timeout:             5 * time.Second,
				PermitWithoutStream: true,
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return core.NewAgentClient(conn), nil
}

func getClient(ctx context.Context) (vald.Client, error) {
	conn, err := grpc.DialContext(
		ctx,
		host+":"+strconv.Itoa(port),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                time.Second,
				Timeout:             5 * time.Second,
				PermitWithoutStream: true,
			},
		),
	)
	if err != nil {
		return nil, err
	}

	return vald.NewValdClient(conn), nil
}

func TestE2EInsert(t *testing.T) {
	ctx := context.Background()

	client, err := getClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	sc, err := client.StreamInsert(ctx)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		count := 0
		for {
			_, err := sc.Recv()
			if err == io.EOF {
				t.Logf("%d items inserted.", count)
				return
			} else if err != nil {
				t.Fatal(err)
			}

			count++

			if count%1000 == 0 {
				t.Logf("inserted: %d", count)
			}
		}
	}()

	t.Log("insert start")
	for i := 0; i < len(ds.train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:     id,
				Vector: ds.train[id],
			},
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: false,
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		if (i+1)%1000 == 0 {
			t.Logf("sent: %d", i+1)
		}

		if i+1 >= insertNum {
			t.Logf("%d items sent.", i+1)
			break
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("insert finished.")
}

func TestE2ECreateIndex(t *testing.T) {
	ctx := context.Background()

	client, err := getAgentClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.CreateAndSaveIndex(ctx, &payload.Control_CreateIndexRequest{
		PoolSize: 10000,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestE2ESearch(t *testing.T) {
	ctx := context.Background()

	client, err := getClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	sc, err := client.StreamSearch(ctx)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		k := 0
		for {
			res, err := sc.Recv()
			if err == io.EOF {
				t.Logf("%d items searched.", k)
				return
			} else if err != nil {
				t.Fatal(err)
			}

			resp := res.GetResponse()
			if resp == nil {
				err := res.GetStatus()
				if err != nil {
					t.Errorf("an error returned: %s", err.GetMessage())
				}
			} else {
				topKIDs := make([]string, len(resp.GetResults()))
				for i, d := range resp.GetResults() {
					topKIDs[i] = d.Id
				}

				if len(topKIDs) == 0 {
					t.Errorf("empty result is returned for ID %d: %#v", k, topKIDs)
				}

				// TODO: validation
				// calculate recall?
				// t.Logf("result: %#v", topKIDs)
				// t.Logf("expected: %#v", ds.neighbors[strconv.Itoa(k)][:len(topKIDs)])
			}

			k++

			if k%1000 == 0 {
				t.Logf("searched: %d", k)
			}
		}
	}()

	t.Log("search start")
	for i := 0; i < len(ds.test); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Search_Request{
			Vector: ds.test[id],
			Config: &payload.Search_Config{
				Num:     100,
				Radius:  -1.0,
				Epsilon: 0.01,
				Timeout: 3000000000,
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		if (i+1)%1000 == 0 {
			t.Logf("sent: %d", i+1)
		}

		if i+1 >= searchNum {
			t.Logf("%d items sent.", i+1)
			break
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("search finished.")
}

func TestE2EIndexInfo(t *testing.T) {
	ctx := context.Background()

	client, err := getAgentClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.IndexInfo(ctx, &payload.Empty{})
	if err != nil {
		t.Fatal(err)
	}

	if res.GetStored() != uint32(insertNum) {
		t.Errorf("stored index number: %d, not equals to expected: %d", res.GetStored(), insertNum)
	}
}
