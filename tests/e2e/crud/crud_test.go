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

// package crud provides e2e tests using ann-benchmarks datasets
package crud

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vdaas/vald/tests/e2e/hdf5"
	"github.com/vdaas/vald/tests/e2e/kubernetes/client"
	"github.com/vdaas/vald/tests/e2e/kubernetes/portforward"
	"github.com/vdaas/vald/tests/e2e/operation"
)

var (
	host string
	port int
	ds   *hdf5.Dataset

	insertNum     int
	searchNum     int
	searchByIDNum int
	getObjectNum  int
	updateNum     int
	upsertNum     int
	removeNum     int

	insertFrom     int
	searchFrom     int
	searchByIDFrom int
	getObjectFrom  int
	updateFrom     int
	upsertFrom     int
	removeFrom     int

	waitAfterInsertDuration time.Duration

	kubeClient client.Client
	namespace  string

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
	flag.IntVar(&upsertNum, "upsert-num", 10000, "number of id-vector pairs used for upsert")
	flag.IntVar(&removeNum, "remove-num", 10000, "number of id-vector pairs used for remove")

	flag.IntVar(&insertFrom, "insert-from", 0, "first index of id-vector pairs used for insert")
	flag.IntVar(&searchFrom, "search-from", 0, "first index of id-vector pairs used for search")
	flag.IntVar(&searchByIDFrom, "search-by-id-from", 0, "first index of id-vector pairs used for search-by-id")
	flag.IntVar(&getObjectFrom, "get-object-from", 0, "first index of id-vector pairs used for get-object")
	flag.IntVar(&updateFrom, "update-from", 0, "first index of id-vector pairs used for update")
	flag.IntVar(&upsertFrom, "upsert-from", 0, "first index of id-vector pairs used for upsert")
	flag.IntVar(&removeFrom, "remove-from", 0, "first index of id-vector pairs used for remove")

	datasetName := flag.String("dataset", "fashion-mnist-784-euclidean.hdf5", "dataset")
	waitAfterInsert := flag.String("wait-after-insert", "3m", "wait duration after inserting vectors")

	pf := flag.Bool("portforward", false, "enable port forwarding")
	pfPodName := flag.String("portforward-pod-name", "vald-gateway-0", "pod name (only for port forward)")
	pfPodPort := flag.Int("portforward-pod-port", port, "pod gRPC port (only for port forward)")

	kubeConfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "kubeconfig path")
	flag.StringVar(&namespace, "namespace", "default", "namespace")

	flag.Parse()

	var err error
	if *pf {
		kubeClient, err = client.New(*kubeConfig)
		if err != nil {
			panic(err)
		}

		forwarder = kubeClient.Portforward(namespace, *pfPodName, port, *pfPodPort)

		err = forwarder.Start()
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("loading dataset: %s ", *datasetName)
	ds, err = hdf5.HDF5ToDataset(*datasetName)
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

func sleep(t *testing.T, dur time.Duration) {
	t.Logf("%v sleep for %s.", time.Now(), dur)
	time.Sleep(dur)
	t.Logf("%v sleep finished.", time.Now())
}

func TestE2EInsertOnly(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()

	op, err := operation.New(host, port)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.Insert(t, ctx, operation.Dataset{
		Train: ds.Train[insertFrom : insertFrom+insertNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}
}

func TestE2ESearchOnly(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()

	op, err := operation.New(host, port)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.Search(t, ctx, operation.Dataset{
		Test:      ds.Test[searchFrom : searchFrom+searchNum],
		Neighbors: ds.Neighbors[searchFrom : searchFrom+searchNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}
}

func TestE2EUpdateOnly(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()

	op, err := operation.New(host, port)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.Update(t, ctx, operation.Dataset{
		Train: ds.Train[updateFrom : updateFrom+updateNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}
}

func TestE2EUpsertOnly(t *testing.T) {
	ctx := context.Background()

	op, err := operation.New(host, port)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.Upsert(t, ctx, operation.Dataset{
		Train: ds.Train[upsertFrom : upsertFrom+upsertNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	teardown()
}

func TestE2ERemoveOnly(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()

	op, err := operation.New(host, port)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.Remove(t, ctx, operation.Dataset{
		Train: ds.Train[removeFrom : removeFrom+removeNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}
}

func TestE2EInsertAndSearch(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()

	op, err := operation.New(host, port)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.Insert(t, ctx, operation.Dataset{
		Train: ds.Train[insertFrom : insertFrom+insertNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	sleep(t, waitAfterInsertDuration)

	err = op.Search(t, ctx, operation.Dataset{
		Test:      ds.Test[searchFrom : searchFrom+searchNum],
		Neighbors: ds.Neighbors[searchFrom : searchFrom+searchNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}
}

func TestE2EStandardCRUD(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()

	op, err := operation.New(host, port)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.Insert(t, ctx, operation.Dataset{
		Train: ds.Train[insertFrom : insertFrom+insertNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	sleep(t, waitAfterInsertDuration)

	err = op.Search(t, ctx, operation.Dataset{
		Test:      ds.Test[searchFrom : searchFrom+searchNum],
		Neighbors: ds.Neighbors[searchFrom : searchFrom+searchNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.SearchByID(t, ctx, operation.Dataset{
		Train: ds.Train[searchByIDFrom : searchByIDFrom+searchByIDNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.Exists(t, ctx, "0")
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.GetObject(t, ctx, operation.Dataset{
		Train: ds.Train[getObjectFrom : getObjectFrom+getObjectNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.Update(t, ctx, operation.Dataset{
		Train: ds.Train[updateFrom : updateFrom+updateNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.Upsert(t, ctx, operation.Dataset{
		Train: ds.Train[upsertFrom : upsertFrom+upsertNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.Remove(t, ctx, operation.Dataset{
		Train: ds.Train[removeFrom : removeFrom+removeNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}
}
