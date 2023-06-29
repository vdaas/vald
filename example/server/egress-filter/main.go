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
	"strings"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/apis/grpc/v1/filter/egress"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/net/grpc/codes"
	"github.com/vdaas/vald/internal/net/grpc/status"
	"google.golang.org/grpc"
)

var (
	egressServerPort uint
	dimension        uint
)

func init() {
	/**
	Port option specifies grpc server port of your egress filter. Default value is `8083`.
	Dimension option specifies dimension size of vectors. Default value is  `784`.
	**/
	flag.UintVar(&egressServerPort, "port", 8083, "server port")
	flag.UintVar(&dimension, "dimension", 784, "dimension size of vectors")
	flag.Parse()
}

func getSplitValue(str string, sep string, pos uint) (string, bool) {
	ss := strings.Split(str, sep)
	if len(ss) == int(pos+1) {
		return ss[pos], true
	}

	return "", false
}

type myEgressServer struct {
	egress.UnimplementedFilterServer
}

func (s *myEgressServer) FilterDistance(ctx context.Context, in *payload.Filter_DistanceRequest) (*payload.Filter_DistanceResponse, error) {
	// Write your own logic
	glg.Log("filtering vector %#v", in)
	qCategory, ok := getSplitValue(in.GetQuery().GetQuery(), "=", 1)
	if !ok {
		return &payload.Filter_DistanceResponse{
			Distance: in.GetDistance(),
		}, nil
	}

	filteredDis := []*payload.Object_Distance{}
	for _, d := range in.GetDistance() {
		iCategory, ok := getSplitValue(d.GetId(), "_", 1)
		if !ok {
			return &payload.Filter_DistanceResponse{
				Distance: in.GetDistance(),
			}, nil
		}

		glg.Infof("qCategory: %v, iCategory: %v", qCategory, iCategory)
		if qCategory == iCategory {
			filteredDis = append(filteredDis, d)
		}
	}

	if len(filteredDis) == 0 {
		return nil, status.Error(codes.NotFound, "FilterDistance results not found.")
	}

	return &payload.Filter_DistanceResponse{
		Distance: filteredDis,
	}, nil
}

func (s *myEgressServer) FilterVector(ctx context.Context, in *payload.Filter_VectorRequest) (*payload.Filter_VectorResponse, error) {
	// Write your own logic
	glg.Logf("filtering the vector %#v", in)
	return &payload.Filter_VectorResponse{
		Vector: in.GetVector(),
	}, nil
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", egressServerPort))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	egress.RegisterFilterServer(s, &myEgressServer{})

	go func() {
		glg.Infof("start gRPC server port: %v", egressServerPort)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	glg.Infof("stopping gRPC server...")
	s.GracefulStop()
}
