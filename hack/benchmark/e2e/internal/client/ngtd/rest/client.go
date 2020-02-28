package rest

import (
	"context"

	"github.com/vdaas/vald/apis/grpc/payload"
	"github.com/vdaas/vald/internal/client"

	proto "github.com/yahoojapan/ngtd/proto"
)

type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type ngtdClient struct {
	addr string
}

func New(ctx context.Context, opts ...Option) (Client, error) {
	c := new(ngtdClient)

	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	return c, nil
}

func (c *ngtdClient) Exists(ctx context.Context, req *client.ObjectID) (*client.ObjectID, error) {
	// TODO: errors.NotSupportedClientMethod
	return nil, nil
}

func tofloat64(s []float32) []float64 {
	return nil
}

func toObjectDistances(in []*proto.ObjectDistance) (to []*payload.Object_Distance) {
	to = make([]*payload.Object_Distance, 0, len(in))

	for _, elm := range in {
		to = append(to, &payload.Object_Distance{
			Id:       string(elm.GetId()),
			Distance: elm.GetDistance(),
		})
	}
	return nil
}

func (c *ngtdClient) Search(ctx context.Context, req *client.SearchRequest) (*client.SearchResponse, error) {
	resp := new(proto.SearchResponse)

	/**
	err := json.Request(ctx, http.MethodPost, c.addr+"/search", &proto.SearchRequest {
		Vector:  tofloat64(req.GetVector()),
		Epsilon: req.GetConfig().GetEpsilon(),
		Size_:   int32(req.GetConfig().GetNum()),
	}, res)
	if err != nil {
		return nil, err
	}
	**/

	return &client.SearchResponse{
		Results: toObjectDistances(resp.GetResult()),
	}, nil
}

func (c *ngtdClient) SearchByID(ctx context.Context, req *client.SearchIDRequest) (*client.SearchResponse, error) {
	resp := new(proto.SearchResponse)

	/**
	err := json.Request(ctx, http.MethodPost, c.addr+"/searchbyid", &proto.SearchRequest {
		Id:      []byte(req.GetId()),
		Epsilon: req.GetConfig().GetEpsilon(),
		Size_:   int32(req.GetConfig().GetNum()),
	}, res)
	if err != nil {
		return nil, err
	}
	**/

	return &client.SearchResponse{
		Results: toObjectDistances(resp.GetResult()),
	}, nil
}

func (c *ngtdClient) StreamSearch(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	// TODO: errors.NotSupportedClientMethod
	return nil
}

func (c *ngtdClient) StreamSearchByID(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	// TODO: errors.NotSupportedClientMethod
	return nil
}

func (c *ngtdClient) Insert(ctx context.Context, req *client.ObjectVector) error {
	/**
	return json.Request(ctx, http.MethodPost, c.addr+"/insert", &proto.InsertRequest {
		Id:     []byte(req.GetId()),
		Vector: tofloat64(req.GetVector()),
	}, nil)
	**/
	return nil
}

func (c *ngtdClient) StreamInsert(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) {
	// TODO: errors.NotSupportedClientMethod
	return
}

func (c *ngtdClient) MultiInsert(ctx context.Context, req *client.ObjectVectors) error {
	// TODO: errors.NotSupportedClientMethod
	return nil
}

func (c *ngtdClient) Update(ctx context.Context, req *client.ObjectVector) error {
	// TODO: errors.NotSupportedClientMethod
	return nil
}

func (c *ngtdClient) StreamUpdate(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) {
	// TODO: errors.NotSupportedClientMethod
}

func (c *ngtdClient) MultiUpdate(ctx context.Context, req *client.ObjectVectors) error {
	// TODO: errors.NotSupportedClientMethod
	return nil
}

func (c *ngtdClient) Remove(ctx context.Context, req *client.ObjectID) error {
	/**
	return json.Request(ctx, http.MethodGet, c.addr+"/remove/"+req.GetId(), nil, nil)
	**/
	return nil
}

func (c *ngtdClient) StreamRemove(ctx context.Context, dataProvider func() *client.ObjectID, f func(error)) {
	// TODO: errors.NotSupportedClientMethod
	return
}

func (c *ngtdClient) MultiRemove(ctx context.Context, req *client.ObjectIDs) error {
	// TODO: errors.NotSupportedClientMethod
	return nil
}

func (c *ngtdClient) GetObject(ctx context.Context, req *client.ObjectID) (*client.ObjectVector, error) {
	// TODO: errors.NotSupportedClientMethod
	return nil, nil
}

func (c *ngtdClient) StreamGetObject(ctx context.Context, dataProvider func() *client.ObjectID, f func(*client.ObjectVector, error)) error {
	// TODO: errors.NotSupportedClientMethod
	return nil
}

func (c *ngtdClient) CreateIndex(ctx context.Context, req *client.ControlCreateIndexRequest) error {
	/**
	return json.Request(ctx, http.MethodGet, c.addr+"/index/create/"+req.GetPoolSize(), nil, nil)
	**/
	return nil
}

func (c *ngtdClient) SaveIndex(ctx context.Context) error {
	/**
	return json.Request(ctx, http.MethodGet, c.addr+"/index/save", nil, nil)
	**/
	return nil
}

func (c *ngtdClient) CreateAndSaveIndex(ctx context.Context, req *client.ControlCreateIndexRequest) error {
	// TODO: errors.NotSupportedClientMethod
	return nil
}

func (c *ngtdClient) IndexInfo(ctx context.Context) (*client.InfoIndex, error) {
	// TODO: errors.NotSupportedClientMethod
	return nil, nil
}
