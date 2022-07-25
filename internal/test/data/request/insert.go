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
package request

import (
	"strconv"

	"github.com/vdaas/vald/apis/grpc/v1/payload"
	"github.com/vdaas/vald/internal/test/data/vector"
)

type ObjectType int

const (
	Uint8 ObjectType = iota
	Float
)

func GenMultiInsertReq(t ObjectType, dist vector.Distribution, num int, dim int, cfg *payload.Insert_Config) (*payload.Insert_MultiRequest, error) {
	var vecs [][]float32
	var err error
	switch t {
	case Float:
		vecs, err = vector.GenF32Vec(dist, num, dim)
	case Uint8:
		vecs, err = vector.GenUint8Vec(dist, num, dim)
	}
	if err != nil {
		return nil, err
	}

	req := &payload.Insert_MultiRequest{
		Requests: make([]*payload.Insert_Request, num),
	}
	for i, vec := range vecs {
		req.Requests[i] = &payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:     "uuid-" + strconv.Itoa(i+1),
				Vector: vec,
			},
			Config: cfg,
		}
	}

	return req, nil
}

// GenSameVecMultiInsertReq generates Insert_MultiRequest with the same vector.
func GenSameVecMultiInsertReq(num int, vec []float32, cfg *payload.Insert_Config) *payload.Insert_MultiRequest {
	req := &payload.Insert_MultiRequest{
		Requests: make([]*payload.Insert_Request, num),
	}
	for i := 0; i < num; i++ {
		req.Requests[i] = &payload.Insert_Request{
			Vector: &payload.Object_Vector{
				Id:     "uuid-" + strconv.Itoa(i+1),
				Vector: vec,
			},
			Config: cfg,
		}
	}

	return req
}
