// Copyright (C) 2019-2024 vdaas.org vald team <vald@vdaas.org>
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
	// NOTE:
	// The correct approach is to use "github.com/vdaas/vald-client-go/v1/payload" and "github.com/vdaas/vald-client-go/v1/vald" in the "example/client".
	// However, the "vald-client-go" module is not available in the filter client example
	// because the changes to the filter query have not been released. (current version is v1.7.12)
	// Therefore, the root module is used until it is released.
	// The import path and go.mod will be changed after release.
	"context"
	"flag"
	"net"
	"strconv"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/apis/grpc/v1/filter/egress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client           egress.FilterClient
	egressServerHost string
	egressServerPort uint
	dimension        uint
)

func init() {
	/**
	  Ingresshost option specifies grpc server host of your egress filter. Default value is `127.0.0.1`.
	  Ingressport option specifies grpc server port of your egress filter. Default value is `8083`.
	  Dimension option specifies dimension size of vectors. Default value is  `784`.
	  **/
	flag.StringVar(&egressServerHost, "host", "127.0.0.1", "ingress server host")
	flag.UintVar(&egressServerPort, "port", 8083, "ingress server port")
	flag.UintVar(&dimension, "dimension", 784, "dimension size of vectors")
	flag.Parse()
}

func main() {
	glg.Println("start gRPC Client.")

	addr := net.JoinHostPort(egressServerHost, strconv.Itoa(int(egressServerPort)))
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		glg.Error("Connection failed.")
		return
	}
	defer conn.Close()

	client = egress.NewFilterClient(conn)

	fdr := &payload.Filter_DistanceRequest{
		Distance: []*payload.Object_Distance{
			{
				Id:       "1_fashion",
				Distance: 0.01,
			},
			{
				Id:       "2_food",
				Distance: 0.02,
			},
			{
				Id:       "3_fashion",
				Distance: 0.03,
			},
			{
				Id:       "4_pet",
				Distance: 0.04,
			},
		},
		Query: &payload.Filter_Query{
			Query: "category=fashion",
		},
	}
	res, err := client.FilterDistance(context.Background(), fdr)
	if err != nil {
		glg.Error(err)
		return
	}
	glg.Info("FilterDistance Distance: ", res.GetDistance())

	r, err := client.FilterVector(context.Background(), &payload.Filter_VectorRequest{
		Vector: &payload.Object_Vector{
			Id: "1", Vector: make([]float32, dimension),
		},
		Query: &payload.Filter_Query{},
	})
	if err != nil {
		glg.Error(err)
		return
	}
	glg.Info("FilterVector Vector: ", r.GetVector())
}
