//go:build e2e

//
// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
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

// package sidecar provides e2e tests for sidecar
package sidecar

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/tests/e2e/hdf5"
	"github.com/vdaas/vald/tests/e2e/kubernetes/client"
	"github.com/vdaas/vald/tests/e2e/kubernetes/portforward"
	"github.com/vdaas/vald/tests/e2e/operation"
)

var (
	host string
	port int
	ds   *hdf5.Dataset

	insertNum int
	searchNum int

	insertFrom int
	searchFrom int

	kubeClient client.Client
	namespace  string

	forwarder *portforward.Portforward
	pfEnabled bool
	pfPodPort int
)

func init() {
	testing.Init()

	flag.StringVar(&host, "host", "localhost", "hostname")
	flag.IntVar(&port, "port", 8081, "gRPC port")

	flag.IntVar(&insertNum, "insert-num", 10000, "number of id-vector pairs used for insert")
	flag.IntVar(&searchNum, "search-num", 10000, "number of id-vector pairs used for search")

	flag.IntVar(&insertFrom, "insert-from", 0, "first index of id-vector pairs used for insert")
	flag.IntVar(&searchFrom, "search-from", 0, "first index of id-vector pairs used for search")

	datasetName := flag.String("dataset", "fashion-mnist-784-euclidean.hdf5", "dataset")

	flag.BoolVar(&pfEnabled, "portforward", false, "enable port forwarding")
	pfPodName := flag.String("portforward-pod-name", "vald-gateway-0", "pod name (only for port forward)")
	flag.IntVar(&pfPodPort, "portforward-pod-port", port, "pod gRPC port (only for port forward)")

	kubeConfig := flag.String("kubeconfig", file.Join(os.Getenv("HOME"), ".kube", "config"), "kubeconfig path")
	flag.StringVar(&namespace, "namespace", "default", "namespace")

	flag.Parse()

	var err error
	if pfEnabled {
		kubeClient, err = client.New(*kubeConfig)
		if err != nil {
			panic(err)
		}

		forwarder = kubeClient.Portforward(namespace, *pfPodName, port, pfPodPort)

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
}

func teardown() {
	if forwarder != nil {
		forwarder.Close()
	}
}

func sleep(t *testing.T, dur time.Duration) {
	t.Logf("sleep for %s", dur)
	time.Sleep(dur)
	t.Log("sleep finished.")
}

func TestE2EForSidecar(t *testing.T) {
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

	err = op.CreateIndex(t, ctx)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	res, err := op.IndexInfo(t, ctx)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	if insertNum != int(res.GetStored()) {
		t.Errorf("Stored index count is invalid, expected: %d, stored: %d", insertNum, res.GetStored())
	}

	err = op.Search(t, ctx, operation.Dataset{
		Test:      ds.Test[searchFrom : searchFrom+searchNum],
		Neighbors: ds.Neighbors[searchFrom : searchFrom+searchNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	err = op.SaveIndex(t, ctx)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	pods, err := kubeClient.GetPods(ctx, namespace, "app=vald-agent-ngt")
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	if len(pods) == 0 {
		t.Fatalf("there's no Agent pods")
	}

	if forwarder != nil {
		forwarder.Close()
		forwarder = nil
	}

	podName := pods[0].Name

	t.Logf("pod: %s", podName)

	err = kubeClient.DeletePod(ctx, namespace, podName)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	t.Logf("pod %s deleted", podName)

	sleep(t, time.Minute)

	pods, err = kubeClient.GetPods(ctx, namespace, "app=vald-agent-ngt")
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	if len(pods) == 0 {
		t.Fatalf("there's no Agent pods")
	}

	podName = pods[0].Name

	t.Logf("pod: %s", podName)

	ok, err := kubeClient.WaitForPodReady(ctx, namespace, podName, 10*time.Minute)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}
	if !ok {
		t.Fatalf("pod didn't become ready")
	}

	if pfEnabled {
		pf := kubeClient.Portforward(namespace, podName, port, pfPodPort)
		err = pf.Start()
		if err != nil {
			t.Fatalf("an error occurred: %s", err)
		}
		defer pf.Close()
	}

	res, err = op.IndexInfo(t, ctx)
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	if insertNum != int(res.GetStored()) {
		t.Errorf("Stored index count is invalid, expected: %d, stored: %d", insertNum, res.GetStored())
	}

	err = op.Search(t, ctx, operation.Dataset{
		Test:      ds.Test[searchFrom : searchFrom+searchNum],
		Neighbors: ds.Neighbors[searchFrom : searchFrom+searchNum],
	})
	if err != nil {
		t.Fatalf("an error occurred: %s", err)
	}

	teardown()
}
