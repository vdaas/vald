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

// generate MultiInsert request with the same vector
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
