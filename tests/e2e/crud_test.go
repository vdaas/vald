// +build e2e

//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/vdaas/vald/apis/grpc/gateway/vald"
	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/tests/e2e/portforward"

	"gonum.org/v1/hdf5"
	"google.golang.org/grpc"
)

var (
	host        string
	port        int
	datasetName string

	forwarder *portforward.Portforward
)

func init() {
	testing.Init()

	flag.StringVar(&host, "host", "localhost", "hostname")
	flag.IntVar(&port, "port", 8081, "gRPC port")
	flag.StringVar(&datasetName, "dataset", "fashion-mnist-784-euclidean.hdf5", "dataset")

	pf := flag.Bool("portforward", false, "enable port forwarding")
	pfNamespace := flag.String("portforward-ns", "default", "namespace (only for port forward)")
	pfPodName := flag.String("portforward-pod-name", "vald-gateway-0", "pod name (only for port forward)")
	pfPodPort := flag.Int("portforward-pod-port", port, "pod gRPC port (only for port forward)")
	kubeConfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "kubeconfig path (only for port forward)")

	flag.Parse()

	if *pf {
		var err error
		forwarder, err = portforward.NewPortforward(*kubeConfig, *pfNamespace, *pfPodName, port, *pfPodPort)
		if err != nil {
			panic(err)
		}

		err = forwarder.Start()
		if err != nil {
			panic(err)
		}
	}
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

func getClient(ctx context.Context) (vald.ValdClient, error) {
	conn, err := grpc.DialContext(
		ctx,
		host+":"+strconv.Itoa(port),
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	return vald.NewValdClient(conn), nil
}

func sleep(t *testing.T, dur time.Duration) {
	t.Logf("sleep for %s", dur)
	time.Sleep(dur)
	t.Log("sleep finished.")
}

func TestE2EInsert(t *testing.T) {
	t.Log("loading")
	ds, err := hdf5ToDataset(datasetName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("loading finished")

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

		for {
			_, err := sc.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				t.Fatal(err)
			}
		}
	}()

	t.Log("insert start")
	count := 0
	for k, v := range ds.train {
		err := sc.Send(&payload.Object_Vector{
			Id:     k,
			Vector: v,
		})
		if err != nil {
			t.Fatal(err)
		}

		count++

		if count%1000 == 0 {
			t.Logf("inserted: %d", count)
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("insert finished.")

	// wait for creating index.
	// TODO: too long?
	sleep(t, 3*time.Minute)
}

func TestE2ESearch(t *testing.T) {
	t.Log("loading")
	ds, err := hdf5ToDataset(datasetName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("loading finished")

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
				return
			} else if err != nil {
				t.Fatal(err)
			}

			topKIDs := make([]string, len(res.GetResults()))
			for i, d := range res.GetResults() {
				topKIDs[i] = d.Id
			}

			if len(topKIDs) == 0 {
				t.Errorf("empty result is returned for ID %d: %#v", k, topKIDs)
			}

			// TODO: validation
			// calculate recall?
			// t.Logf("result: %#v", topKIDs)
			// t.Logf("expected: %#v", ds.neighbors[strconv.Itoa(k)][:len(topKIDs)])

			k++
		}
	}()

	t.Log("search start")
	count := 0
	for _, v := range ds.test {
		err := sc.Send(&payload.Search_Request{
			Vector: v,
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

		count++

		if count%1000 == 0 {
			t.Logf("searched: %d", count)
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("search finished.")
}

func TestE2ESearchByID(t *testing.T) {
	t.Log("loading")
	ds, err := hdf5ToDataset(datasetName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("loading finished")

	ctx := context.Background()

	client, err := getClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	sc, err := client.StreamSearchByID(ctx)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				t.Fatal(err)
			}

			topKIDs := make([]string, len(res.GetResults()))
			for i, d := range res.GetResults() {
				topKIDs[i] = d.Id
			}

			if len(topKIDs) == 0 {
				t.Errorf("empty result is returned: %#v", topKIDs)
			}
		}
	}()

	t.Log("search-by-id start")
	count := 0
	for k, _ := range ds.test {
		err := sc.Send(&payload.Search_IDRequest{
			Id: k,
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

		count++

		if count%1000 == 0 {
			t.Logf("searched: %d", count)
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("search-by-id finished.")
}

func TestE2EGetObject(t *testing.T) {
	t.Log("loading")
	ds, err := hdf5ToDataset(datasetName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("loading finished")

	ctx := context.Background()

	client, err := getClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	sc, err := client.StreamGetObject(ctx)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			res, err := sc.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(res.GetVector(), ds.train[res.GetMeta()]) {
				t.Errorf(
					"result: %#v, expected: %#v",
					res.GetVector(),
					ds.train[res.GetMeta()],
				)
			}
		}
	}()

	t.Log("get object start")
	count := 0
	for k, _ := range ds.train {
		err := sc.Send(&payload.Object_ID{
			Id: k,
		})
		if err != nil {
			t.Fatal(err)
		}

		count++

		if count%1000 == 0 {
			t.Logf("get object: %d", count)
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("get object finished.")
}

func TestE2EUpdate(t *testing.T) {
	t.Log("loading")
	ds, err := hdf5ToDataset(datasetName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("loading finished")

	ctx := context.Background()

	client, err := getClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	sc, err := client.StreamUpdate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			_, err := sc.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				t.Fatal(err)
			}
		}
	}()

	t.Log("update start")
	count := 0
	for k, v := range ds.train {
		err := sc.Send(&payload.Object_Vector{
			Id:     k,
			Vector: append(v[1:], v[0]), // shift
		})
		if err != nil {
			t.Fatal(err)
		}

		count++

		if count%1000 == 0 {
			t.Logf("updated: %d", count)
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("update finished.")
}

func TestE2ERemove(t *testing.T) {
	t.Log("loading")
	ds, err := hdf5ToDataset(datasetName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("loading finished")

	ctx := context.Background()

	client, err := getClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	sc, err := client.StreamRemove(ctx)
	if err != nil {
		t.Fatal(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			_, err := sc.Recv()
			if err == io.EOF {
				return
			} else if err != nil {
				t.Fatal(err)
			}
		}
	}()

	t.Log("remove start")
	count := 0
	for k, _ := range ds.train {
		err := sc.Send(&payload.Object_ID{
			Id: k,
		})
		if err != nil {
			t.Fatal(err)
		}

		count++

		if count%1000 == 0 {
			t.Logf("removed: %d", count)
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("remove finished.")
}
