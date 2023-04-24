// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
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
	"flag"
	"fmt"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/apis/grpc/v1/filter/egress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net"
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

	addr := net.JoinHostPort(egressServerHost, uint16(egressServerPort))
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		glg.Error("Connection failed.")
		return
	}
	defer conn.Close()

	client = egress.NewFilterClient(conn)

	res, err := client.FilterDistance(context.Background(), &payload.Filter_DistanceRequest{})
	if err != nil {
		glg.Error(err)
	} else {
		glg.Info("FilterDistance Distance: ", res.GetDistance())
	}

	// TODO: Fix
	r, err := client.FilterVector(context.Background(), &payload.Filter_VectorRequest{
		Vector: []*payload.Object_Vector{
			{Id: "1", Vector: make([]float32, dimension)},
		},
		Query: &payload.Filter_Query{},
	})
	if err != nil {
		glg.Error(err)
	} else {
		glg.Info("FilterVector Vector: ", r.Vector)
	}
}
