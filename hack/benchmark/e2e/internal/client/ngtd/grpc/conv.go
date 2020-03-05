package grpc

import (
	"github.com/vdaas/vald/internal/client"
	proto "github.com/yahoojapan/ngtd/proto"
)

func SearchRequestToNGTDSearchRequest(in *client.SearchRequest) *proto.SearchRequest {
	size, epsilon := getSizeAndEpsilon(in.GetConfig())
	return &proto.SearchRequest{
		Vector:  tofloat64(in.GetVector()),
		Size_:   size,
		Epsilon: epsilon,
	}
}

func SearchIDRequestToNGTDSearchRequest(in *client.SearchIDRequest) *proto.SearchRequest {
	size, epsilon := getSizeAndEpsilon(in.GetConfig())
	return &proto.SearchRequest{
		Id:      []byte(in.GetId()),
		Size_:   size,
		Epsilon: epsilon,
	}
}

func NGTDSearchResponseToSearchResponse(in *proto.SearchResponse) *client.SearchResponse {
	if len(in.GetError()) != 0 {
		return nil
	}

	results := make([]*client.ObjectDistance, 0, len(in.GetResult()))
	for i, _ := range results {
		if len(in.Result[i].GetError()) == 0 {
			results = append(results, &client.ObjectDistance{
				Id:       string(in.Result[i].GetId()),
				Distance: in.Result[i].GetDistance(),
			})
		}
	}
	return &client.SearchResponse{
		Results: results,
	}
}

func NGTDGetObjectResponseToObjectVector(in *proto.GetObjectResponse) *client.ObjectVector {
	if len(in.GetError()) != 0 {
		return nil
	}

	return &client.ObjectVector{
		Id:     string(in.GetId()),
		Vector: in.GetVector(),
	}
}

func ObjectVectorToNGTDInsertRequest(in *client.ObjectVector) *proto.InsertRequest {
	return &proto.InsertRequest{
		Id:     []byte(in.GetId()),
		Vector: tofloat64(in.GetVector()),
	}
}

func ObjectIDToNGTDRemoveRequest(in *client.ObjectID) *proto.RemoveRequest {
	return &proto.RemoveRequest{
		Id: []byte(in.GetId()),
	}
}

func ObjectIDToNGTDGetObjectRequest(in *client.ObjectID) *proto.GetObjectRequest {
	return &proto.GetObjectRequest{
		Id: []byte(in.GetId()),
	}
}

func ControlCreateIndexRequestToCreateIndexRequest(in *client.ControlCreateIndexRequest) *proto.CreateIndexRequest {
	return &proto.CreateIndexRequest{
		PoolSize: in.GetPoolSize(),
	}
}

func getSizeAndEpsilon(cfg *client.SearchConfig) (size int32, epsilon float32) {
	if cfg != nil {
		size = int32(cfg.GetNum())
		epsilon = float32(cfg.GetEpsilon())
	}
	return
}

func tofloat64(in []float32) (out []float64) {
	out = make([]float64, len(in))
	for i := range in {
		out[i] = float64(in[i])
	}
	return
}
