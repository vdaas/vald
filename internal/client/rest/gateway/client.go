package gateway

import (
	"context"
	"net/http"

	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/json"
)

type Client interface {
	client.Client
	client.MetaObjectReader
}

type gatewayClient struct {
	addr string
}

func New(opts ...Option) (Client, error) {
	c := new(gatewayClient)

	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}

	return c, nil
}

func (c *gatewayClient) Exists(ctx context.Context, req *client.ObjectID) (resp *client.ObjectID, err error) {
	resp = new(client.ObjectID)
	err = json.Request(ctx, http.MethodGet, c.addr+"/exists/"+req.GetId(), req, resp)
	return
}

func (c *gatewayClient) Search(ctx context.Context, req *client.SearchRequest) (resp *client.SearchResponse, err error) {
	resp = new(client.SearchResponse)
	err = json.Request(ctx, http.MethodPost, c.addr+"/search", req, resp)
	return
}

func (c *gatewayClient) SearchByID(ctx context.Context, req *client.SearchIDRequest) (resp *client.SearchResponse, err error) {
	resp = new(client.SearchResponse)
	err = json.Request(ctx, http.MethodPost, c.addr+"/search/id", req, resp)
	return
}

func (c *gatewayClient) StreamSearch(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) StreamSearchByID(ctx context.Context, dataProvider func() *client.SearchRequest, f func(*client.SearchResponse, error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) Insert(ctx context.Context, req *client.ObjectVector) error {
	return json.Request(ctx, http.MethodPost, c.addr+"/insert", req, nil)
}

func (c *gatewayClient) StreamInsert(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) MultiInsert(ctx context.Context, objectVectors *client.ObjectVectors) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) Update(ctx context.Context, req *client.ObjectVector) error {
	return json.Request(ctx, http.MethodPost, c.addr+"/update", req, nil)
}

func (c *gatewayClient) StreamUpdate(ctx context.Context, dataProvider func() *client.ObjectVector, f func(error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) MultiUpdate(ctx context.Context, objectVectors *client.ObjectVectors) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) Remove(ctx context.Context, objectID *client.ObjectID) error {
	return json.Request(ctx, http.MethodDelete, c.addr+"/remove/"+objectID.GetId(), nil, nil)
}

func (c *gatewayClient) StreamRemove(ctx context.Context, dataProvider func() *client.ObjectID, f func(error)) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) MultiRemove(ctx context.Context, objectIDs *client.ObjectIDs) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) GetObject(ctx context.Context, req *client.ObjectID) (resp *client.MetaObject, err error) {
	resp = new(client.MetaObject)
	err = json.Request(ctx, http.MethodGet, c.addr+"/object/"+req.GetId(), nil, nil)
	return
}

func (c *gatewayClient) StreamGetObject(ctx context.Context, dataProvider func() *client.ObjectID, f func(*client.MetaObject, error)) error {
	return errors.ErrUnsupportedClientMethod
}
