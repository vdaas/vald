package agent

import (
	"context"
	"net/http"

	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/config"
	igrpc "github.com/vdaas/vald/internal/net/grpc"
	ijson "github.com/vdaas/vald/internal/net/http/json"
)

type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type agentClient struct {
	addr              string
	cfg               *config.GRPCClient
	streamConcurrency int
	grpcClient        igrpc.Client
	client            *http.Client
}

func New(ctx context.Context, opts ...Option) Client {
	c := new(agentClient)

	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	c.grpcClient = igrpc.New(c.cfg.Opts()...)

	return c
}

func (c *agentClient) Exists(ctx context.Context, req *client.ObjectID) (res *client.ObjectID, err error) {
	res = new(client.ObjectID)
	err = ijson.Request(ctx, http.MethodGet, "", req, res)
	return
}

func (c *agentClient) Search(ctx context.Context, searchRequest *client.SearchRequest) (*client.SearchResponse, error) {
	return nil, nil
}

func (c *agentClient) SearchByID(ctx context.Context, req *client.SearchIDRequest) (res *client.SearchResponse, err error) {
	res = new(client.SearchResponse)
	err = ijson.Request(ctx, http.MethodPost, "", req, res)
	return
}

func (c *agentClient) StreamSearch(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	return nil
}

func (c *agentClient) StreamSearchByID(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	return nil
}

func (c *agentClient) Insert(ctx context.Context, objectVector *client.ObjectVector) error {
	return nil
}

func (c *agentClient) StreamInsert(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) error {
	return nil
}

func (c *agentClient) MultiInsert(ctx context.Context, objectVectors *client.ObjectVectors) error {
	return nil
}

func (c *agentClient) Update(ctx context.Context, objectVector *client.ObjectVector) error {
	return nil
}

func (c *agentClient) StreamUpdate(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) error {
	return nil
}

func (c *agentClient) MultiUpdate(ctx context.Context, objectVectors *client.ObjectVectors) error {
	return nil
}

func (c *agentClient) Remove(ctx context.Context, objectID *client.ObjectID) error {
	return nil
}

func (c *agentClient) StreamRemove(ctx context.Context, dataProvider func() *client.ObjectID, f func(error)) error {
	return nil
}

func (c *agentClient) MultiRemove(ctx context.Context, objectIDs *client.ObjectIDs) error {
	return nil
}

func (c *agentClient) GetObject(ctx context.Context, objectID *client.ObjectID) (*client.ObjectVector, error) {
	return nil, nil
}

func (c *agentClient) StreamGetObject(ctx context.Context, dataProvider func() *client.ObjectID, f func(*client.ObjectVector, error)) error {
	return nil
}

func (c *agentClient) CreateIndex(ctx context.Context, controlCreateIndexRequest *client.ControlCreateIndexRequest) error {
	return nil
}

func (c *agentClient) SaveIndex(ctx context.Context) error {
	return nil
}

func (c *agentClient) CreateAndSaveIndex(ctx context.Context, controlCreateIndexRequest *client.ControlCreateIndexRequest) error {
	return nil
}

func (c *agentClient) IndexInfo(ctx context.Context) (*client.InfoIndex, error) {
	return nil, nil
}
