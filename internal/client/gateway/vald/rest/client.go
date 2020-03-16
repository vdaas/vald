//
// Copyright (C) 2019-2020 Vdaas.org Vald team ( kpango, rinx, kmrmt )
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

// Package rest provides vald REST client functions
package rest

import (
	"context"
	"net/http"

	"github.com/vdaas/vald/internal/client"
	"github.com/vdaas/vald/internal/errors"
	"github.com/vdaas/vald/internal/net/http/json"
)

// Client represents gateway client interface.
type Client interface {
	client.Client
	client.MetaObjectReader
	client.Upserter
}

type gatewayClient struct {
	addr string
}

// New returns Client implementation.
func New(opts ...Option) Client {
	c := new(gatewayClient)
	for _, opt := range append(defaultOptions, opts...) {
		opt(c)
	}
	return c
}

func (c *gatewayClient) Exists(
	ctx context.Context,
	req *client.ObjectID,
) (resp *client.ObjectID, err error) {
	resp = new(client.ObjectID)
	err = json.Request(ctx, http.MethodGet, c.addr+"/exists/"+req.GetId(), req, resp)
	return
}

func (c *gatewayClient) Search(
	ctx context.Context,
	req *client.SearchRequest,
) (resp *client.SearchResponse, err error) {
	resp = new(client.SearchResponse)
	err = json.Request(ctx, http.MethodPost, c.addr+"/search", req, resp)
	return
}

func (c *gatewayClient) SearchByID(
	ctx context.Context,
	req *client.SearchIDRequest,
) (resp *client.SearchResponse, err error) {
	resp = new(client.SearchResponse)
	err = json.Request(ctx, http.MethodPost, c.addr+"/search/id", req, resp)
	return
}

func (c *gatewayClient) StreamSearch(
	ctx context.Context,
	dataProvider func() *client.SearchRequest,
	f func(*client.SearchResponse, error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) StreamSearchByID(
	ctx context.Context,
	dataProvider func() *client.SearchIDRequest,
	f func(*client.SearchResponse, error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) Insert(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	return json.Request(ctx, http.MethodPost, c.addr+"/insert", req, nil)
}

func (c *gatewayClient) StreamInsert(
	ctx context.Context,
	dataProvider func() *client.ObjectVector,
	f func(error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) MultiInsert(
	ctx context.Context,
	req *client.ObjectVectors,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) Update(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	return json.Request(ctx, http.MethodPost, c.addr+"/update", req, nil)
}

func (c *gatewayClient) StreamUpdate(
	ctx context.Context,
	dataProvider func() *client.ObjectVector,
	f func(error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) MultiUpdate(
	ctx context.Context,
	req *client.ObjectVectors,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) Upsert(
	ctx context.Context,
	req *client.ObjectVector,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) MultiUpsert(
	context.Context,
	*client.ObjectVectors,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) StreamUpsert(
	context.Context,
	func() *client.ObjectVector,
	func(error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) Remove(
	ctx context.Context,
	req *client.ObjectID,
) error {
	return json.Request(ctx, http.MethodDelete, c.addr+"/remove/"+req.GetId(), nil, nil)
}

func (c *gatewayClient) StreamRemove(
	ctx context.Context,
	dataProvider func() *client.ObjectID,
	f func(error),
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) MultiRemove(
	ctx context.Context,
	req *client.ObjectIDs,
) error {
	return errors.ErrUnsupportedClientMethod
}

func (c *gatewayClient) GetObject(
	ctx context.Context,
	req *client.ObjectID,
) (resp *client.MetaObject, err error) {
	resp = new(client.MetaObject)
	err = json.Request(ctx, http.MethodGet, c.addr+"/object/"+req.GetId(), nil, nil)
	return
}

func (c *gatewayClient) StreamGetObject(
	ctx context.Context,
	dataProvider func() *client.ObjectID,
	f func(*client.MetaObject, error),
) error {
	return errors.ErrUnsupportedClientMethod
}
