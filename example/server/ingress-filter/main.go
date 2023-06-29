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
	"net"
	"os"
	"os/signal"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/apis/grpc/v1/filter/ingress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/test/data/vector"
	"google.golang.org/grpc"
)

var (
	ingressServerPort uint
	dimension         uint
)

func init() {
	/**
	Port option specifies grpc server port of your ingress filter. Default value is `8082`.
	Dimension option specifies dimension size of vectors. Default value is  `784`.
	**/
	flag.UintVar(&ingressServerPort, "port", 8082, "server port")
	flag.UintVar(&dimension, "dimension", 784, "dimension size of vectors")
	flag.Parse()
}

type myIngressServer struct {
	ingress.UnimplementedFilterServer
}

func (s *myIngressServer) GenVector(ctx context.Context, in *payload.Object_Blob) (*payload.Object_Vector, error) {
	// Write your own logic
	glg.Logf("generating vector %#v", in)
	vec, err := vector.GenF32Vec(vector.Gaussian, 1, int(dimension))
	if err != nil {
		return nil, err
	}
	return &payload.Object_Vector{
		Id:     in.GetId(),
		Vector: vec[0],
	}, nil
}

func (s *myIngressServer) FilterVector(ctx context.Context, in *payload.Object_Vector) (*payload.Object_Vector, error) {
	// Write your own logic
	glg.Logf("filtering vector %#v", in)
	return in, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", ingressServerPort))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	ingress.RegisterFilterServer(s, &myIngressServer{})

	go func() {
		glg.Infof("start gRPC server port: %v", ingressServerPort)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	glg.Infof("stopping gRPC server...")
	s.GracefulStop()
}
