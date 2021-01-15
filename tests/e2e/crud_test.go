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
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

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

	insertFrom     int
	searchFrom     int
	searchByIDFrom int
	getObjectFrom  int
	updateFrom     int
	removeFrom     int

	waitAfterInsertDuration time.Duration

	forwarder *portforward.Portforward
)

func init() {
	testing.Init()

	flag.StringVar(&host, "host", "localhost", "hostname")
	flag.IntVar(&port, "port", 8081, "gRPC port")

	flag.IntVar(&insertNum, "insert-num", 10000, "number of id-vector pairs used for insert")
	flag.IntVar(&searchNum, "search-num", 10000, "number of id-vector pairs used for search")
	flag.IntVar(&searchByIDNum, "search-by-id-num", 100, "number of id-vector pairs used for search-by-id")
	flag.IntVar(&getObjectNum, "get-object-num", 100, "number of id-vector pairs used for get-object")
	flag.IntVar(&updateNum, "update-num", 10000, "number of id-vector pairs used for update")
	flag.IntVar(&removeNum, "remove-num", 10000, "number of id-vector pairs used for remove")

	flag.IntVar(&insertFrom, "insert-from", 0, "first index of id-vector pairs used for insert")
	flag.IntVar(&searchFrom, "search-from", 0, "first index of id-vector pairs used for search")
	flag.IntVar(&searchByIDFrom, "search-by-id-from", 0, "first index of id-vector pairs used for search-by-id")
	flag.IntVar(&getObjectFrom, "get-object-from", 0, "first index of id-vector pairs used for get-object")
	flag.IntVar(&updateFrom, "update-from", 0, "first index of id-vector pairs used for update")
	flag.IntVar(&removeFrom, "remove-from", 0, "first index of id-vector pairs used for remove")

	datasetName := flag.String("dataset", "fashion-mnist-784-euclidean.hdf5", "dataset")
	waitAfterInsert := flag.String("wait-after-insert", "3m", "wait duration after inserting vectors")

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

	waitAfterInsertDuration, err = time.ParseDuration(*waitAfterInsert)
	if err != nil {
		panic(err)
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

func sleep(t *testing.T, dur time.Duration) {
	t.Logf("sleep for %s", dur)
	time.Sleep(dur)
	t.Log("sleep finished.")
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
	for i := insertFrom; i < len(ds.train); i++ {
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
			t.Fatalf("send failed at %d: %s", i+1, err)
		}

		if (i+1)%1000 == 0 {
			t.Logf("sent: %d", i+1)
		}

		if i+1 >= insertFrom+insertNum {
			t.Logf("%d items sent.", i+1)
			break
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("insert finished.")

	sleep(t, waitAfterInsertDuration)
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
				err := res.GetError()
				if err != nil {
					t.Errorf("an error returned: %s", err)
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
	for i := searchFrom; i < len(ds.test); i++ {
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

		if i+1 >= searchFrom+searchNum {
			t.Logf("%d items sent.", i+1)
			break
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("search finished.")
}

func TestE2ESearchByID(t *testing.T) {
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

		count := 0
		for {
			res, err := sc.Recv()
			if err == io.EOF {
				t.Logf("%d items searched.", count)
				return
			} else if err != nil {
				t.Fatal(err)
			}

			resp := res.GetResponse()
			if resp == nil {
				err := res.GetError()
				if err != nil {
					t.Errorf("an error returned: %s", err)
				}
			} else {
				topKIDs := make([]string, len(resp.GetResults()))
				for i, d := range resp.GetResults() {
					topKIDs[i] = d.Id
				}

				if len(topKIDs) == 0 {
					t.Errorf("empty result is returned: %#v", topKIDs)
				}
			}

			count++

			if count%1000 == 0 {
				t.Logf("searched: %d", count)
			}
		}
	}()

	t.Log("search-by-id start")
	for i := searchByIDFrom; i < len(ds.train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Search_IDRequest{
			Id: id,
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

		if i+1 >= searchByIDFrom+searchByIDNum {
			t.Logf("%d items sent.", i+1)
			break
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("search-by-id finished.")
}

func TestE2EGetObject(t *testing.T) {
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

		count := 0
		for {
			res, err := sc.Recv()
			if err == io.EOF {
				t.Logf("%d items got.", count)
				return
			} else if err != nil {
				t.Fatal(err)
			}

			resp := res.GetVector()
			if resp == nil {
				err := res.GetError()
				if err != nil {
					t.Errorf("an error returned: %s", err)
				}
			} else {
				if !reflect.DeepEqual(res.GetVector(), ds.train[resp.GetId()]) {
					t.Errorf(
						"result: %#v, expected: %#v",
						res.GetVector(),
						ds.train[resp.GetId()],
					)
				}
			}

			count++

			if count%1000 == 0 {
				t.Logf("get object: %d", count)
			}
		}
	}()

	t.Log("get object start")
	for i := getObjectFrom; i < len(ds.train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Object_ID{
			Id: id,
		})
		if err != nil {
			t.Fatal(err)
		}

		if (i+1)%1000 == 0 {
			t.Logf("sent: %d", i+1)
		}

		if i+1 >= getObjectFrom+getObjectNum {
			t.Logf("%d items sent.", i+1)
			break
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("get object finished.")
}

func TestE2EUpdate(t *testing.T) {
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

		count := 0
		for {
			_, err := sc.Recv()
			if err == io.EOF {
				t.Logf("%d items updated.", count)
				return
			} else if err != nil {
				t.Fatal(err)
			}

			count++

			if count%1000 == 0 {
				t.Logf("updated: %d", count)
			}
		}
	}()

	t.Log("update start")
	for i := updateFrom; i < len(ds.train); i++ {
		id := strconv.Itoa(i)
		v := ds.train[id]
		err := sc.Send(&payload.Update_Request{
			Vector: &payload.Object_Vector{
				Id:     id,
				Vector: append(v[1:], v[0]), // shift
			},
			Config: &payload.Update_Config{
				SkipStrictExistCheck: false,
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		if (i+1)%1000 == 0 {
			t.Logf("sent: %d", i+1)
		}

		if i+1 >= updateFrom+updateNum {
			t.Logf("%d items sent.", i+1)
			break
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("update finished.")
}

func TestE2ERemove(t *testing.T) {
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

		count := 0
		for {
			_, err := sc.Recv()
			if err == io.EOF {
				t.Logf("%d items removed.", count)
				return
			} else if err != nil {
				t.Fatal(err)
			}

			count++

			if count%1000 == 0 {
				t.Logf("removed: %d", count)
			}
		}
	}()

	t.Log("remove start")
	for i := removeFrom; i < len(ds.train); i++ {
		id := strconv.Itoa(i)
		err := sc.Send(&payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: id,
			},
			Config: &payload.Remove_Config{
				SkipStrictExistCheck: false,
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		if (i+1)%1000 == 0 {
			t.Logf("sent: %d", i+1)
		}

		if i+1 >= removeFrom+removeNum {
			t.Logf("%d items sent.", i+1)
			break
		}
	}

	sc.CloseSend()

	wg.Wait()

	t.Log("remove finished.")
}
