//go:build e2e

//
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
	"encoding/json"
	"flag"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/vdaas/vald-client-go/v1/payload"
	"github.com/vdaas/vald-client-go/v1/vald"
	"github.com/vdaas/vald/internal/core/algorithm"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/file"
	"github.com/vdaas/vald/internal/log"
	"github.com/vdaas/vald/internal/log/logger"
	"github.com/vdaas/vald/internal/net"
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
	maxBit = 32
	id     = "1"
)

var (
	host                string
	port                int
	namespace           string
	bit                 int
	dim                 int
	kubeClient          client.Client
	forwarder           *portforward.Portforward
	indexingWaitSeconds uint
	kubeConfig          string
	pf                  bool
	pfPodName           string
	pfPodPort           int
	fileName            string
)

func init() {
	testing.Init()

	flag.StringVar(&host, "host", "localhost", "hostname")
	flag.IntVar(&port, "port", 8081, "gRPC port")
	flag.StringVar(&namespace, "namespace", "default", "namespace")
	flag.IntVar(&bit, "bit", 1, "bit")
	flag.UintVar(&indexingWaitSeconds, "wait", 30, "indexing wait seconds")
	flag.StringVar(&fileName, "file", "tmp.log", "output file name")

	flag.BoolVar(&pf, "portforward", false, "enable port forwarding")
	flag.StringVar(&pfPodName, "portforward-pod-name", "vald-lb-gateway", "pod name (only for port forward)")
	flag.IntVar(&pfPodPort, "portforward-pod-port", port, "pod gRPC port (only for port forward)")

	flag.StringVar(&kubeConfig, "kubeconfig", file.Join(os.Getenv("HOME"), ".kube", "config"), "kubeconfig path")

	flag.Parse()
}

func TestMain(m *testing.M) {
	log.Init(log.WithLoggerType(logger.NOP.String()))
	var err error
	if pf {
		kubeClient, err = client.New(kubeConfig)
		if err != nil {
			os.WriteFile(fileName, []byte(err.Error()), os.ModePerm)
			os.Exit(1)
		}
		forwarder = kubeClient.Portforward(namespace, pfPodName, port, pfPodPort)
		err = forwarder.Start()
		if err != nil {
			os.WriteFile(fileName, []byte(err.Error()), os.ModePerm)
			os.Exit(1)
		}
	}

	if bit < 1 || maxBit < bit {
		err = errors.New("Invalid argument: bit should be 0 ~ 32. set bit was " + strconv.Itoa(bit))
		os.WriteFile(fileName, []byte(err.Error()), os.ModePerm)
		os.Exit(1)
	}
	_ = m.Run()
	if forwarder != nil {
		forwarder.Close()
	}
	os.Exit(0)
}

func TestE2EInsertOnlyWithOneVectorAndSearch(t *testing.T) {
	t.Helper()
	dim := 1 << bit
	if bit == maxBit {
		dim--
	}
	if dim > algorithm.MaximumVectorDimensionSize {
		t.Fatalf("Invalid argument: dimension should be equal or under than " + strconv.Itoa(algorithm.MaximumVectorDimensionSize) + ". set dim was " + strconv.Itoa(dim))
	}
	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		net.JoinHostPort(host, uint16(port)),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(
			keepalive.ClientParameters{
				Time:                time.Second,
				Timeout:             5 * time.Minute,
				PermitWithoutStream: true,
			},
		),
	)
	if err != nil {
		t.Fatalf("Failed to create grpc conn interface: %v", err)
	}

	cli := vald.NewValdClient(conn)
	vec := vector.GaussianDistributedFloat32VectorGenerator(1, dim)[0]
	req := &payload.Insert_Request{
		Vector: &payload.Object_Vector{
			// Id should be named the unique name in the production environment.
			// In this case, it is the simple value because the main purpose of this test is checking the max vector dimension.
			Id:     id,
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
			os.WriteFile(fileName, []byte(st.Code().String()), os.ModePerm)
			return
		}
		t.Fatalf("TestE2EInsertOnlyWithOneVectorAndSearch\t Insert Error: %v", err)
	}
	t.Logf("[Pass] Insert process (Bit = %d)", bit)
	wt := time.Duration(indexingWaitSeconds) * time.Second
	for cnt := 0; cnt < 10; cnt++ {
		t.Logf("[Pause] Wait %#vs for Finish createIndex process (Bit = %d, cnt = %d)", wt.Seconds(), bit, cnt+1)
		time.Sleep(wt)
		res, err := cli.SearchByID(
			ctx,
			&payload.Search_IDRequest{
				Id: id,
				Config: &payload.Search_Config{
					Num: 1,
				},
			},
		)
		if errors.Is(nil, err) {
			b, err := json.MarshalIndent(res.GetResults(), "", " ")
			if err != nil {
				t.Fatalf("TestE2EInsertOnlyWithOneVectorAndSearch\tMarshalIndent Error: %v", err)
			}
			t.Logf("[Pass] SearchByID process (Bit = %d)", bit)
			if string(b) != "" {
				os.WriteFile(fileName, []byte("OK"), os.ModePerm)
				return
			}
		}
		st, _ := status.FromError(err)
		if st.Code() != codes.Code(code.Code_NOT_FOUND) {
			t.Fatalf("TestE2EInsertOnlyWithOneVectorAndSearch\t SearchById Error: %v", err)
			break
		}
	}
	t.Fatal("TestE2EInsertOnlyWithOneVectorAndSearch\tError: No Result")
}
