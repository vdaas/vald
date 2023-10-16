// Copyright (C) 2019-2023 vdaas.org vald team <vald@vdaas.org>
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
	"context"
	"encoding/json"
	"flag"
	"time"

	"github.com/kpango/glg"
	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/apis/grpc/v1/vald"
	"google.golang.org/grpc"
)

type dataset struct {
	id     string
	vector []float32
}

var (
	grpcServerAddr      string
	ingressServerHost   string
	ingressServerPort   uint
	egressServerHost    string
	egressServerPort    uint
	indexingWaitSeconds uint
	dimension           uint
)

func init() {
	/**
	Addr option specifies grpc server address of filter gateway. Default value is `127.0.0.1:8081`.
	Ingresshost option specifies grpc server host of your ingress filter. Default value is `127.0.0.1`.
	Ingressport option specifies grpc server port of your ingress filter. Default value is `8082`.
	Egresshost option specifies grpc server host of your egress filter. Default value is `127.0.0.1`.
	Egressport option specifies grpc server port of your egress filter. Default value is `8083`.
	Wait option specifies indexing wait time (in seconds). Default value is  `240`.
	Dimension option specifies dimension size of vectors. Default value is  `784`.
	**/
	flag.StringVar(&grpcServerAddr, "addr", "127.0.0.1:8081", "gRPC server address of filter gateway")
	flag.StringVar(&ingressServerHost, "ingresshost", "127.0.0.1", "ingress server host")
	flag.UintVar(&ingressServerPort, "ingressport", 8082, "ingress server port")
	flag.StringVar(&egressServerHost, "egresshost", "127.0.0.1", "egress server host")
	flag.UintVar(&egressServerPort, "egressport", 8083, "egress server port")
	flag.UintVar(&indexingWaitSeconds, "wait", 240, "indexing wait seconds")
	flag.UintVar(&dimension, "dimension", 784, "dimension size of vectors")
	flag.Parse()
}

// Please execute after setting up the server of vald cluster and ingress/egress filter
func main() {
	dataset := genDataset()
	query := "category=fashion"

	// connect to the Vald cluster
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		glg.Error(err)
		return
	}

	// create a filter client
	glg.Info("Start inserting object via vald filter client")
	var object []byte
	fclient := vald.NewFilterClient(conn)

	for _, ds := range dataset {
		icfg := &payload.Insert_ObjectRequest{
			// object data to pass to GenVector function of your ingress filter
			Object: &payload.Object_Blob{
				Id:     ds.id,
				Object: object,
			},
			// insert config
			Config: &payload.Insert_Config{
				SkipStrictExistCheck: false,
				// config to call FilterVector function of your ingress filter
				Filters: []*payload.Filter_Config{
					{
						Target: &payload.Filter_Target{
							Host: ingressServerHost,
							Port: uint32(ingressServerPort),
						},
						Query: &payload.Filter_Query{},
					},
				},
			},
			// specify vectorizer component location
			Vectorizer: &payload.Filter_Target{
				Host: ingressServerHost,
				Port: uint32(ingressServerPort),
			},
		}

		// send InsertObject request
		res, err := fclient.InsertObject(ctx, icfg)
		if err != nil {
			glg.Error(err)
			return
		}

		glg.Infof("location: %#v", res.Ips)
	}

	// Vald Agent starts indexing automatically after insert. It needs to wait until the indexing is completed before a search action is performed.
	wt := time.Duration(indexingWaitSeconds) * time.Second
	glg.Infof("Wait %s for indexing to finish", wt)
	time.Sleep(wt)

	// create a search client
	glg.Log("Start searching dataset")
	sclient := vald.NewSearchClient(conn)

	for _, ds := range dataset {
		scfg := &payload.Search_Config{
			Num:     10,
			Epsilon: 0.1,
			Radius:  -1,
			// config to call DistanceVector function of your egress filter
			EgressFilters: []*payload.Filter_Config{
				{
					Target: &payload.Filter_Target{
						Host: egressServerHost,
						Port: uint32(egressServerPort),
					},
					Query: &payload.Filter_Query{
						Query: query,
					},
				},
			},
		}

		// send Search request
		res, err := sclient.Search(ctx, &payload.Search_Request{
			Vector: ds.vector,
			Config: scfg,
		})
		if err != nil {
			glg.Error(err)
			return
		}
		b, _ := json.MarshalIndent(res.GetResults(), "", " ")
		glg.Infof("Results : %s\n\n", string(b))
	}

	// create an object client
	glg.Info("Start GetObject")
	oclient := vald.NewObjectClient(conn)

	for _, ds := range dataset {
		vreq := &payload.Object_VectorRequest{
			Id: &payload.Object_ID{Id: ds.id},
			// config to call FilterVector function of your egress filter
			Filters: []*payload.Filter_Config{
				{
					Target: &payload.Filter_Target{
						Host: egressServerHost,
						Port: uint32(egressServerPort),
					},
					Query: &payload.Filter_Query{},
				},
			},
		}

		// send GetObject request
		r, err := oclient.GetObject(ctx, vreq)
		if err != nil {
			glg.Error(err)
			return
		}
		b, _ := json.Marshal(r.GetVector())
		glg.Infof("Get Object result: %s\n", string(b))
	}

	// send remove request
	glg.Info("Start removing data")
	rclient := vald.NewRemoveClient(conn)

	for _, ds := range dataset {
		rreq := &payload.Remove_Request{
			Id: &payload.Object_ID{
				Id: ds.id,
			},
		}
		if _, err := rclient.Remove(ctx, rreq); err != nil {
			glg.Errorf("Failed to remove, ID: %v", ds.id)
		} else {
			glg.Infof("Remove ID %v successed", ds.id)
		}
	}
}

func genDataset() []dataset {
	// create a data set for operation confirmation
	makeVecFn := func(dim int, value float32) []float32 {
		vec := make([]float32, dim)
		for i := 0; i < dim; i++ {
			vec[i] = value
		}
		return vec
	}
	return []dataset{
		{
			id:     "1_fashion",
			vector: makeVecFn(int(dimension), 0.1),
		},
		{
			id:     "2_food",
			vector: makeVecFn(int(dimension), 0.2),
		},
		{
			id:     "3_fashion",
			vector: makeVecFn(int(dimension), 0.3),
		},
		{
			id:     "4_pet",
			vector: makeVecFn(int(dimension), 0.4),
		},
	}
}
