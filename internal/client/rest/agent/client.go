package agent

import (
	"context"
	"net/http"

	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/json"
)

type Client interface {
	client.Client
	client.ObjectReader
	client.Indexer
}

type agentClient struct {
	addr string
}

func New(ctx context.Context, opts ...Option) Client {
	c := new(agentClient)

	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	return c
}

func (c *agentClient) Exists(ctx context.Context, req *client.ObjectID) (res *client.ObjectID, err error) {
	res = new(client.ObjectID)
	err = json.Request(ctx, http.MethodGet, c.addr+"/exists/"+req.GetId(), req, res)
	return
}

func (c *agentClient) Search(ctx context.Context, req *client.SearchRequest) (res *client.SearchResponse, err error) {
	res = new(client.SearchResponse)
	err = json.Request(ctx, http.MethodPost, "/search", req, res)
	return
}

func (c *agentClient) SearchByID(ctx context.Context, req *client.SearchIDRequest) (res *client.SearchResponse, err error) {
	res = new(client.SearchResponse)
	err = json.Request(ctx, http.MethodPost, "/search/id", req, res)
	return
}

func (c *agentClient) StreamSearch(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *agentClient) StreamSearchByID(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *agentClient) Insert(ctx context.Context, req *client.ObjectVector) error {
	return json.Request(ctx, http.MethodPost, "/insert", req, nil)
}

func (c *agentClient) StreamInsert(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *agentClient) MultiInsert(ctx context.Context, objectVectors *client.ObjectVectors) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *agentClient) Update(ctx context.Context, req *client.ObjectVector) error {
	return json.Request(ctx, http.MethodPost, "/update", req, nil)
}

func (c *agentClient) StreamUpdate(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *agentClient) MultiUpdate(ctx context.Context, objectVectors *client.ObjectVectors) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *agentClient) Remove(ctx context.Context, req *client.ObjectID) error {
	return json.Request(ctx, http.MethodDelete, "/remove/"+req.GetId(), req, nil)
}

func (c *agentClient) StreamRemove(ctx context.Context, dataProvider func() *client.ObjectID, f func(error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *agentClient) MultiRemove(ctx context.Context, objectIDs *client.ObjectIDs) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *agentClient) GetObject(ctx context.Context, req *client.ObjectID) (res *client.ObjectVector, err error) {
	res = new(client.ObjectVector)
	err = json.Request(ctx, http.MethodGet, "/object/"+req.GetId(), req, res)
	return
}

func (c *agentClient) StreamGetObject(ctx context.Context, dataProvider func() *client.ObjectID, f func(*client.ObjectVector, error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *agentClient) CreateIndex(ctx context.Context, req *client.ControlCreateIndexRequest) error {
	return json.Request(ctx, http.MethodGet, "/index/create", req, nil)
}

func (c *agentClient) SaveIndex(ctx context.Context) error {
	return json.Request(ctx, http.MethodGet, "/index/save", nil, nil)
}

func (c *agentClient) CreateAndSaveIndex(ctx context.Context, controlCreateIndexRequest *client.ControlCreateIndexRequest) error {
	return json.Request(ctx, http.MethodGet, "/index/createandsave", nil, nil)
}

func (c *agentClient) IndexInfo(ctx context.Context) (res *client.InfoIndex, err error) {
	res = new(client.InfoIndex)
	err = json.Request(ctx, http.MethodGet, "/index/info", nil, res)
	return
}
