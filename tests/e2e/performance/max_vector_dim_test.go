//go:build e2e
// +build e2e

// Copyright (C) 2019-2022 vdaas.org vald team <vald@vdaas.org>
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

// package performance provides e2e tests for check max bit size
package performance

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"

	"github.com/vdaas/vald/internal/core/algorithm/ngt"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"github.com/vdaas/vald/internal/test/data/vector"
	"github.com/vdaas/vald/tests/e2e/kubernetes/client"
	"github.com/vdaas/vald/tests/e2e/kubernetes/portforward"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	maxBit       = 32
	freeMemLimit = 500 // Limit of free memory remaining(MB)
)

var (
	host       string
	port       int
	namespace  string
	bit        int
	dim        int
	kubeClient client.Client
	forwarder  *portforward.Portforward
)

func init() {
	testing.Init()

	flag.StringVar(&host, "host", "localhost", "hostname")
	flag.IntVar(&port, "port", 8081, "gRPC port")
	flag.StringVar(&namespace, "namespace", "default", "namespace")
	flag.IntVar(&bit, "bit", 2, "bit")
	pf := flag.Bool("portforward", false, "enable port forwarding")
	pfPodName := flag.String("portforward-pod-name", "vald-lb-gateway", "pod name (only for port forward)")
	pfPodPort := flag.Int("portforward-pod-port", port, "pod gRPC port (only for port forward)")

	kubeConfig := flag.String("kubeconfig", file.Join(os.Getenv("HOME"), ".kube", "config"), "kubeconfig path")

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

	if bit < 2 || maxBit < bit {
		err := errors.New("Invalid argument: bit should be 2 ~ 32. set bit was " + strconv.Itoa(bit))
		panic(err)
	}
}

func teardown() {
	if forwarder != nil {
		forwarder.Close()
	}
}

func TestE2EInsertOnlyWithOneVector(t *testing.T) {
	t.Helper()
	t.Cleanup(teardown)
	dim := 1 << bit
	if bit == maxBit {
		dim--
	}
	if dim > ngt.VectorDimensionSizeLimit {
		t.Fatalf("Invalid argument: dimension should be equal or under than " + strconv.Itoa(ngt.VectorDimensionSizeLimit) + ". set dim was " + strconv.Itoa(dim))
	}
	ctx := context.Background()
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
		t.Fatalf("Failed to create grpc conn interface: %#v", err)
	}

	cli := vald.NewValdClient(conn)
	vec := vector.GaussianDistributedFloat32VectorGenerator(1, dim)[0]
	req := &payload.Insert_Request{
		Vector: &payload.Object_Vector{
			// Id should be named the unique name in the production environment.
			// In this case, it is the simple value because the main purpose of this test is checking the max vector dimension.
			Id:     "1",
			Vector: vec,
		},
		Config: &payload.Insert_Config{
			SkipStrictExistCheck: false,
		},
	}
	_, err = cli.Insert(ctx, req)
	if err != nil {
		st, _ := status.FromError(err)
		if st.Code() == codes.Code(code.Code_RESOURCE_EXHAUSTED) {
			// For checking code in the step of the github actions
			fmt.Println("Code: " + st.Code().String())
			// Output: Code: ResourceExhausted
			return
		}
		t.Fatal(err)
	}
	// For checking code in the step of the github actions
	fmt.Println("Code: OK")
	// Output: Code: OK
}
